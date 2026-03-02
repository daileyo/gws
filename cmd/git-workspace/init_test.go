package main

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	gitcfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/daileyo/gws/internal/config"
)

// createInitTestRepo creates a minimal valid git repository for init tests.
func createInitTestRepo(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("createInitTestRepo: failed to create dir %s: %v", path, err)
	}
	repo, err := gogit.PlainInit(path, false)
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to init git repo: %v", err)
	}
	_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to create remote: %v", err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to get worktree: %v", err)
	}
	f := filepath.Join(path, "README.md")
	if err := os.WriteFile(f, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("createInitTestRepo: failed to write file: %v", err)
	}
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("createInitTestRepo: failed to stage file: %v", err)
	}
	if _, err := worktree.Commit("Initial commit", &gogit.CommitOptions{
		Author: &object.Signature{Name: "Test User", Email: "test@example.com"},
	}); err != nil {
		t.Fatalf("createInitTestRepo: failed to commit: %v", err)
	}
}

// withTempHome redirects HOME to a temp directory so config.Save/Load/Exists
// operate on an isolated file instead of the real ~/.gws/config.json.
func withTempHome(t *testing.T) {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
}

// withTempWorkdir changes the process working directory to dir for the duration
// of the test, restoring it in a defer.
func withTempWorkdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("withTempWorkdir: failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("withTempWorkdir: failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestRunInit_HappyPath(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir()
	repoDir := filepath.Join(workspaceDir, "my-repo")
	createInitTestRepo(t, repoDir)

	withTempWorkdir(t, workspaceDir)

	// Resolve symlinks so the comparison works on macOS (/var → /private/var)
	resolvedWorkspace, err := filepath.EvalSymlinks(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for workspace dir: %v", err)
	}

	if err := runInit(""); err != nil {
		t.Fatalf("runInit returned unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after init: %v", err)
	}
	if cfg.Workspace != resolvedWorkspace {
		t.Errorf("Workspace: expected %s, got %s", resolvedWorkspace, cfg.Workspace)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Repository count: expected 1, got %d", len(cfg.Repositories))
	}
	if len(cfg.Repositories) > 0 && cfg.Repositories[0].Name != "my-repo" {
		t.Errorf("Repository name: expected 'my-repo', got %s", cfg.Repositories[0].Name)
	}
}

func TestRunInit_AlreadyInitialized(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir()
	withTempWorkdir(t, workspaceDir)

	// First init creates the workspace
	if err := runInit(""); err != nil {
		t.Fatalf("First runInit failed: %v", err)
	}

	originalCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after first init: %v", err)
	}

	// Second init should return nil (exit 0) without overwriting the config
	if err := runInit(""); err != nil {
		t.Errorf("Second runInit should return nil (already-initialized guard), got: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after second init: %v", err)
	}
	if cfg.Workspace != originalCfg.Workspace {
		t.Errorf("Workspace changed by second init: %s -> %s", originalCfg.Workspace, cfg.Workspace)
	}
}

func TestRunInit_EmptyDirectory(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir() // no git repos inside
	withTempWorkdir(t, workspaceDir)

	if err := runInit(""); err != nil {
		t.Fatalf("runInit on empty directory returned unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 0 {
		t.Errorf("Expected 0 repositories for empty directory, got %d", len(cfg.Repositories))
	}
}
