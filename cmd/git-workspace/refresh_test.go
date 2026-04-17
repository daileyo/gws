package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// saveConfig writes a config with the given repos so runRefresh() picks it up.
func saveConfigWithRepos(t *testing.T, workspace string, repos []config.Repository) {
	t.Helper()
	cfg := config.New(workspace)
	cfg.Repositories = repos
	if err := config.Save(cfg); err != nil {
		t.Fatalf("saveConfigWithRepos: failed to save: %v", err)
	}
}

// TestRunRefresh_ValidPathRetained verifies that a repo at a still-valid path is
// kept in config and its metadata is refreshed.
func TestRunRefresh_ValidPathRetained(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := filepath.Join(workspaceDir, "valid-repo")
	createInitTestRepo(t, repoDir)

	// Pre-populate config with the repo including a tag that must be preserved.
	saveConfigWithRepos(t, workspaceDir, []config.Repository{
		{Name: "valid-repo", Path: repoDir, Tags: []string{"keep-me"}},
	})

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Fatalf("expected 1 repository, got %d", len(cfg.Repositories))
	}
	repo := cfg.Repositories[0]
	if repo.Path != repoDir {
		t.Errorf("path: expected %s, got %s", repoDir, repo.Path)
	}
	if len(repo.Tags) != 1 || repo.Tags[0] != "keep-me" {
		t.Errorf("tags not preserved: got %v", repo.Tags)
	}
}

// TestRunRefresh_MissingPathRemoved verifies that a repo whose path no longer
// contains a .git directory is removed from config, and its workspace symlink
// is also removed.
func TestRunRefresh_MissingPathRemoved(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := setupExternalRepo(t, "gone-repo")

	// Create the workspace symlink for this external repo.
	symlinkPath := filepath.Join(workspaceDir, "gone-repo")
	if err := os.Symlink(repoDir, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// Track the repo in config.
	saveConfigWithRepos(t, workspaceDir, []config.Repository{
		{Name: "gone-repo", Path: repoDir},
	})

	// Remove the repo directory to simulate a deleted repo.
	if err := os.RemoveAll(repoDir); err != nil {
		t.Fatalf("failed to remove repo dir: %v", err)
	}

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 0 {
		t.Errorf("expected 0 repositories after removal, got %d", len(cfg.Repositories))
	}

	// Symlink must be removed.
	if _, err := os.Lstat(symlinkPath); err == nil {
		t.Errorf("expected workspace symlink to be removed, but it still exists: %s", symlinkPath)
	}
}

// TestRunRefresh_NewSymlinkDiscovered verifies that a symlink placed in the
// workspace that points to an external repo is discovered and added with the
// real path (not the symlink path).
func TestRunRefresh_NewSymlinkDiscovered(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := setupExternalRepo(t, "sym-repo")

	// Manually place a symlink in the workspace — simulates user doing this
	// outside of gws add.
	symlinkPath := filepath.Join(workspaceDir, "sym-repo")
	if err := os.Symlink(repoDir, symlinkPath); err != nil {
		t.Fatalf("failed to create workspace symlink: %v", err)
	}

	// Config starts empty.
	saveConfigWithRepos(t, workspaceDir, nil)

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Fatalf("expected 1 repository, got %d", len(cfg.Repositories))
	}
	// Real path must be stored, not the symlink path.
	if cfg.Repositories[0].Path != repoDir {
		t.Errorf("path: expected real path %s, got %s", repoDir, cfg.Repositories[0].Path)
	}
}

// TestRunRefresh_RealDirDiscovered verifies that a real git repository placed
// directly inside the workspace directory is discovered and added.
func TestRunRefresh_RealDirDiscovered(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := filepath.Join(workspaceDir, "new-repo")
	createInitTestRepo(t, repoDir)

	// Config starts empty.
	saveConfigWithRepos(t, workspaceDir, nil)

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Fatalf("expected 1 repository, got %d", len(cfg.Repositories))
	}
	if cfg.Repositories[0].Name != "new-repo" {
		t.Errorf("name: expected 'new-repo', got %s", cfg.Repositories[0].Name)
	}
}

// TestRunRefresh_BrokenSymlinkSkipped verifies that a broken symlink in the
// workspace does not cause a crash and is silently skipped.
func TestRunRefresh_BrokenSymlinkSkipped(t *testing.T) {
	workspaceDir := setupWorkspace(t)

	// Create a symlink pointing to a path that does not exist.
	brokenPath := filepath.Join(workspaceDir, "broken-link")
	if err := os.Symlink("/nonexistent/path/to/nowhere", brokenPath); err != nil {
		t.Fatalf("failed to create broken symlink: %v", err)
	}

	saveConfigWithRepos(t, workspaceDir, nil)

	// Must not error.
	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error for broken symlink: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 0 {
		t.Errorf("expected 0 repositories (broken symlink should be skipped), got %d", len(cfg.Repositories))
	}
}

// TestRunRefresh_DiscoverWorktrees verifies that refresh populates the Worktrees
// field for repos that have worktrees, and leaves it empty for those that don't.
func TestRunRefresh_DiscoverWorktrees(t *testing.T) {
	workspaceDir := setupWorkspace(t)

	// Create two repos — one with worktrees, one without.
	repoWithWT := filepath.Join(workspaceDir, "repo-with-wt")
	createInitTestRepo(t, repoWithWT)

	repoWithoutWT := filepath.Join(workspaceDir, "repo-without-wt")
	createInitTestRepo(t, repoWithoutWT)

	// Add a worktree to the first repo.
	wtDir := repoWithWT + ".wt"
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatalf("failed to create .wt dir: %v", err)
	}
	wtPath := filepath.Join(wtDir, "feature-x")
	cmd := exec.Command("git", "worktree", "add", "-b", "feature-x", wtPath)
	cmd.Dir = repoWithWT
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	// Resolve paths for comparison (macOS symlink resolution).
	resolvedWithWT, _ := filepath.EvalSymlinks(repoWithWT)

	saveConfigWithRepos(t, workspaceDir, []config.Repository{
		{Name: "repo-with-wt", Path: resolvedWithWT},
		{Name: "repo-without-wt", Path: repoWithoutWT},
	})

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	var withWT, withoutWT *config.Repository
	for i := range cfg.Repositories {
		switch cfg.Repositories[i].Name {
		case "repo-with-wt":
			withWT = &cfg.Repositories[i]
		case "repo-without-wt":
			withoutWT = &cfg.Repositories[i]
		}
	}

	if withWT == nil {
		t.Fatal("repo-with-wt not found in config")
	}
	if len(withWT.Worktrees) != 1 {
		t.Fatalf("expected 1 worktree for repo-with-wt, got %d", len(withWT.Worktrees))
	}
	if withWT.Worktrees[0].Branch != "feature-x" {
		t.Errorf("expected worktree branch 'feature-x', got '%s'", withWT.Worktrees[0].Branch)
	}
	if !withWT.Worktrees[0].Aligned {
		t.Error("expected worktree to be aligned (inside .wt/ dir)")
	}

	if withoutWT == nil {
		t.Fatal("repo-without-wt not found in config")
	}
	if len(withoutWT.Worktrees) != 0 {
		t.Errorf("expected 0 worktrees for repo-without-wt, got %d", len(withoutWT.Worktrees))
	}
}

// TestRunRefresh_CorrectSymlinkNoDuplicate verifies that when a tracked external
// repo already has a correct workspace symlink, refresh does not create a
// duplicate entry in config.
func TestRunRefresh_CorrectSymlinkNoDuplicate(t *testing.T) {
	workspaceDir := setupWorkspace(t)
	repoDir := setupExternalRepo(t, "ext-repo")

	// Create the correct workspace symlink.
	symlinkPath := filepath.Join(workspaceDir, "ext-repo")
	if err := os.Symlink(repoDir, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// Track the repo in config.
	saveConfigWithRepos(t, workspaceDir, []config.Repository{
		{Name: "ext-repo", Path: repoDir},
	})

	if err := runRefresh(); err != nil {
		t.Fatalf("runRefresh returned error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("expected exactly 1 repository (no duplicate), got %d", len(cfg.Repositories))
	}

	// Symlink must still be correct.
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if target != repoDir {
		t.Errorf("symlink target: expected %s, got %s", repoDir, target)
	}
}
