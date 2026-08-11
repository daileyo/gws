package reconcile

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// mkRepo creates a git repository with one commit at path, creating parents.
func mkRepo(t *testing.T, path string) string {
	t.Helper()

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create dir %s: %v", path, err)
	}

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...) // #nosec G204 -- fixed test commands
		cmd.Dir = path
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v in %s failed: %s\n%s", args, path, err, out)
		}
	}

	return realPath(t, path)
}

// mkRepoWithRemote creates a repository with an origin remote.
func mkRepoWithRemote(t *testing.T, path, remoteURL string) string {
	t.Helper()
	resolved := mkRepo(t, path)
	gitIn(t, path, "remote", "add", "origin", remoteURL)
	return resolved
}

// gitIn runs a git command in dir, failing the test on error.
func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...) // #nosec G204 -- test-controlled arguments
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s failed: %s\n%s", args, dir, err, out)
	}
}

// addWorktree creates a linked worktree on a new branch.
func addWorktree(t *testing.T, repoDir, branch, destPath string) {
	t.Helper()
	gitIn(t, repoDir, "worktree", "add", "-b", branch, destPath)
}

// realPath resolves symlinks for comparison against reconcile output.
func realPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return resolved
}

// symlink creates a symlink at linkPath pointing to target.
func symlink(t *testing.T, target, linkPath string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		t.Fatalf("Failed to create parent for symlink: %v", err)
	}
	if err := os.Symlink(target, linkPath); err != nil {
		t.Fatalf("Failed to create symlink %s -> %s: %v", linkPath, target, err)
	}
}

// repoByPath finds a repository in a result by path.
func repoByPath(result *Result, path string) *config.Repository {
	for i := range result.Repositories {
		if result.Repositories[i].Path == path {
			return &result.Repositories[i]
		}
	}
	return nil
}

// repoPathSet returns the set of repository paths in a result.
func repoPathSet(result *Result) map[string]bool {
	paths := make(map[string]bool, len(result.Repositories))
	for _, repo := range result.Repositories {
		paths[repo.Path] = true
	}
	return paths
}
