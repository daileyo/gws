package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// setupWorkspace creates a temporary home + an initialized workspace config,
// returning the workspace directory path.
func setupWorkspace(t *testing.T) string {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	workspaceDir := t.TempDir()
	resolvedWorkspace, err := filepath.EvalSymlinks(workspaceDir)
	if err != nil {
		t.Fatalf("setupWorkspace: EvalSymlinks failed: %v", err)
	}

	cfg := config.New(resolvedWorkspace)
	if err := config.Save(cfg); err != nil {
		t.Fatalf("setupWorkspace: failed to save config: %v", err)
	}
	return resolvedWorkspace
}

// setupExternalRepo creates a git repo in a directory that is NOT inside the workspace.
func setupExternalRepo(t *testing.T, name string) string {
	t.Helper()
	externalDir := t.TempDir()
	repoDir := filepath.Join(externalDir, name)
	createInitTestRepo(t, repoDir)
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("setupExternalRepo: EvalSymlinks failed: %v", err)
	}
	return resolved
}

func TestRunAdd_CurrentDirectory(t *testing.T) {
	workspaceDir := setupWorkspace(t)

	// Create a git repo inside the workspace (not external — no symlink expected)
	repoDir := filepath.Join(workspaceDir, "my-repo")
	createInitTestRepo(t, repoDir)

	// Set flagAdd to "." and CWD to the repo
	origFlag := flagAdd
	flagAdd = "."
	defer func() { flagAdd = origFlag }()
	withTempWorkdir(t, repoDir)

	if err := runAdd(nil, nil); err != nil {
		t.Fatalf("runAdd returned unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(cfg.Repositories))
	}
	if len(cfg.Repositories) > 0 && cfg.Repositories[0].Name != "my-repo" {
		t.Errorf("Expected repo name 'my-repo', got %s", cfg.Repositories[0].Name)
	}
}

func TestRunAdd_ExplicitExternalPath(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := setupExternalRepo(t, "external-repo")

	origFlag := flagAdd
	flagAdd = repoDir
	defer func() { flagAdd = origFlag }()

	if err := runAdd(nil, nil); err != nil {
		t.Fatalf("runAdd returned unexpected error: %v", err)
	}

	// Verify repo is tracked
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Fatalf("Expected 1 repository, got %d", len(cfg.Repositories))
	}
	if cfg.Repositories[0].Name != "external-repo" {
		t.Errorf("Expected repo name 'external-repo', got %s", cfg.Repositories[0].Name)
	}

	// Verify symlink was created in workspace
	symlinkPath := filepath.Join(workspaceDir, "external-repo")
	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		t.Fatalf("Expected symlink at %s, got error: %v", symlinkPath, err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Errorf("Expected a symlink at %s, got mode %v", symlinkPath, fi.Mode())
	}

	// Verify symlink points to the correct target
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if target != repoDir {
		t.Errorf("Symlink target: expected %s, got %s", repoDir, target)
	}
}

func TestRunAdd_NonGitDirectory(t *testing.T) {
	setupWorkspace(t)

	nonGitDir := t.TempDir()

	origFlag := flagAdd
	flagAdd = nonGitDir
	defer func() { flagAdd = origFlag }()

	err := runAdd(nil, nil)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
}

func TestRunAdd_AlreadyTracked(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := filepath.Join(workspaceDir, "my-repo")
	createInitTestRepo(t, repoDir)

	origFlag := flagAdd
	flagAdd = repoDir
	defer func() { flagAdd = origFlag }()

	// First add
	if err := runAdd(nil, nil); err != nil {
		t.Fatalf("First runAdd failed: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	countAfterFirst := len(cfg.Repositories)

	// Second add — should warn and skip, not error
	if err := runAdd(nil, nil); err != nil {
		t.Errorf("Second runAdd should return nil for already-tracked repo, got: %v", err)
	}

	cfg, err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after second add: %v", err)
	}
	if len(cfg.Repositories) != countAfterFirst {
		t.Errorf("Repository count changed on duplicate add: %d → %d", countAfterFirst, len(cfg.Repositories))
	}
}

func TestRunAdd_SymlinkPathConflict(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := setupExternalRepo(t, "conflict-repo")

	// Pre-create a file at the would-be symlink path
	conflictPath := filepath.Join(workspaceDir, "conflict-repo")
	if err := os.WriteFile(conflictPath, []byte("blocker"), 0644); err != nil {
		t.Fatalf("Failed to create conflict file: %v", err)
	}

	origFlag := flagAdd
	flagAdd = repoDir
	defer func() { flagAdd = origFlag }()

	// Should succeed (repo added to config) but skip symlink creation
	if err := runAdd(nil, nil); err != nil {
		t.Errorf("runAdd should not error on symlink conflict, got: %v", err)
	}

	// Verify repo IS tracked in config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Expected 1 repository in config, got %d", len(cfg.Repositories))
	}

	// Verify the conflict file is still a plain file (not replaced by symlink)
	fi, err := os.Lstat(conflictPath)
	if err != nil {
		t.Fatalf("Conflict path missing: %v", err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Error("Conflict path was replaced by a symlink; expected it to remain unchanged")
	}
}

func TestRunAdd_NoWorkspace(t *testing.T) {
	// Redirect HOME to an empty temp dir — no config file exists
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	origFlag := flagAdd
	flagAdd = "."
	defer func() { flagAdd = origFlag }()

	err := runAdd(nil, nil)
	if err == nil {
		t.Error("Expected error when no workspace initialized, got nil")
	}
}

func TestCreateSymlinkIfExternal_InsideWorkspace(t *testing.T) {
	workspaceDir := t.TempDir()
	resolved, _ := filepath.EvalSymlinks(workspaceDir)
	cfg := &config.Config{Workspace: resolved}

	// Repo inside workspace — no symlink should be created
	repoPath := filepath.Join(resolved, "my-repo")
	created, err := createSymlinkIfExternal(cfg, repoPath, "my-repo")
	if err != nil {
		t.Errorf("Unexpected error for internal repo: %v", err)
	}
	if created {
		t.Error("Expected no symlink for repo inside workspace")
	}
}

func TestRunAdd_ConfigMetadata(t *testing.T) {
	// Verify that repo type and remote URL are populated correctly
	workspaceDir := setupWorkspace(t)
	repoDir := filepath.Join(workspaceDir, "typed-repo")
	createInitTestRepo(t, repoDir) // createInitTestRepo sets origin to https://github.com/test/repo.git

	origFlag := flagAdd
	flagAdd = repoDir
	defer func() { flagAdd = origFlag }()

	if err := runAdd(nil, nil); err != nil {
		t.Fatalf("runAdd failed: %v", err)
	}

	configPath, _ := config.GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if len(cfg.Repositories) == 0 {
		t.Fatal("No repositories in config")
	}
	repo := cfg.Repositories[0]
	if repo.RemoteURL != "https://github.com/test/repo.git" {
		t.Errorf("RemoteURL: expected 'https://github.com/test/repo.git', got %q", repo.RemoteURL)
	}
	if repo.Type != config.TypeGitHub {
		t.Errorf("Type: expected 'github', got %q", repo.Type)
	}
}
