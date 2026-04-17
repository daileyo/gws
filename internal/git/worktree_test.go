package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// initBareGitRepo creates a git repo with an initial commit so worktrees can be added.
func initBareGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v failed: %s\n%s", args, err, out)
		}
	}
}

func TestListWorktrees(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create two worktrees
	wt1 := filepath.Join(t.TempDir(), "wt-feature")
	wt2 := filepath.Join(t.TempDir(), "wt-bugfix")

	cmd := exec.Command("git", "worktree", "add", "-b", "feature", wt1)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add feature failed: %s\n%s", err, out)
	}

	cmd = exec.Command("git", "worktree", "add", "-b", "bugfix", wt2)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add bugfix failed: %s\n%s", err, out)
	}

	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 worktrees, got %d", len(entries))
	}

	// Collect branches
	branches := map[string]string{}
	for _, e := range entries {
		branches[e.Branch] = e.Path
	}

	if _, ok := branches["feature"]; !ok {
		t.Error("Expected to find worktree with branch 'feature'")
	}
	if _, ok := branches["bugfix"]; !ok {
		t.Error("Expected to find worktree with branch 'bugfix'")
	}
}

func TestListWorktrees_NoWorktrees(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 worktrees, got %d", len(entries))
	}
}

func TestListWorktrees_MainWorktreeExcluded(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Add one worktree
	wtPath := filepath.Join(t.TempDir(), "wt-test")
	cmd := exec.Command("git", "worktree", "add", "-b", "test-branch", wtPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add failed: %s\n%s", err, out)
	}

	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	// Main worktree should be excluded, only the added one should appear
	for _, e := range entries {
		if filepath.Clean(e.Path) == filepath.Clean(repoDir) {
			t.Error("Main worktree should be excluded from results")
		}
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 worktree entry, got %d", len(entries))
	}
}

func TestParseWorktreeListPorcelain(t *testing.T) {
	repoPath := "/workspace/my-repo"

	porcelain := `worktree /workspace/my-repo
HEAD abc123
branch refs/heads/main

worktree /workspace/my-repo.wt/feature-x
HEAD def456
branch refs/heads/feature-x

worktree /tmp/hotfix
HEAD 789abc
branch refs/heads/hotfix/urgent
`

	entries := parseWorktreeListPorcelain(porcelain, repoPath)

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	if entries[0].Path != "/workspace/my-repo.wt/feature-x" {
		t.Errorf("Expected path '/workspace/my-repo.wt/feature-x', got '%s'", entries[0].Path)
	}
	if entries[0].Branch != "feature-x" {
		t.Errorf("Expected branch 'feature-x', got '%s'", entries[0].Branch)
	}

	if entries[1].Path != "/tmp/hotfix" {
		t.Errorf("Expected path '/tmp/hotfix', got '%s'", entries[1].Path)
	}
	if entries[1].Branch != "hotfix/urgent" {
		t.Errorf("Expected branch 'hotfix/urgent', got '%s'", entries[1].Branch)
	}
}

func TestParseWorktreeListPorcelain_Empty(t *testing.T) {
	entries := parseWorktreeListPorcelain("", "/some/repo")
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty output, got %d", len(entries))
	}
}

func TestIsAligned(t *testing.T) {
	tests := []struct {
		name         string
		worktreePath string
		repoPath     string
		expected     bool
	}{
		{
			name:         "aligned - directly inside .wt dir",
			worktreePath: "/workspace/my-repo.wt/feature-x",
			repoPath:     "/workspace/my-repo",
			expected:     true,
		},
		{
			name:         "aligned - nested inside .wt dir",
			worktreePath: "/workspace/my-repo.wt/feature/auth-flow",
			repoPath:     "/workspace/my-repo",
			expected:     true,
		},
		{
			name:         "unaligned - different location",
			worktreePath: "/tmp/some-worktree",
			repoPath:     "/workspace/my-repo",
			expected:     false,
		},
		{
			name:         "unaligned - similar name but not .wt suffix",
			worktreePath: "/workspace/my-repo-wt/feature",
			repoPath:     "/workspace/my-repo",
			expected:     false,
		},
		{
			name:         "unaligned - sibling directory with similar prefix",
			worktreePath: "/workspace/my-repo.wtx/feature",
			repoPath:     "/workspace/my-repo",
			expected:     false,
		},
		{
			name:         "aligned - trailing slash on repo path",
			worktreePath: "/workspace/my-repo.wt/bugfix",
			repoPath:     "/workspace/my-repo/",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAligned(tt.worktreePath, tt.repoPath)
			if got != tt.expected {
				t.Errorf("IsAligned(%q, %q) = %v, want %v", tt.worktreePath, tt.repoPath, got, tt.expected)
			}
		})
	}
}

// TestListWorktrees_BranchWithSlashes verifies branches like feature/auth are handled.
func TestListWorktrees_BranchWithSlashes(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	wtPath := filepath.Join(t.TempDir(), "wt-feat")
	cmd := exec.Command("git", "worktree", "add", "-b", "feature/auth-flow", wtPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add failed: %s\n%s", err, out)
	}

	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 worktree, got %d", len(entries))
	}
	if entries[0].Branch != "feature/auth-flow" {
		t.Errorf("Expected branch 'feature/auth-flow', got '%s'", entries[0].Branch)
	}
}

func TestAddWorktree(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	destPath := filepath.Join(t.TempDir(), "new-wt")
	resolved, _ := filepath.EvalSymlinks(destPath)
	if resolved != "" {
		destPath = resolved
	}

	err := AddWorktree(repoDir, "new-feature", destPath)
	if err != nil {
		t.Fatalf("AddWorktree failed: %v", err)
	}

	// Verify the worktree directory exists
	if _, err := os.Stat(destPath); err != nil {
		t.Errorf("worktree directory should exist at %s: %v", destPath, err)
	}

	// Verify it appears in ListWorktrees
	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	found := false
	for _, e := range entries {
		if e.Branch == "new-feature" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 'new-feature' worktree after AddWorktree")
	}
}

func TestAddWorktree_ExistingBranch(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create a branch first
	cmd := exec.Command("git", "branch", "existing-branch")
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git branch failed: %s\n%s", err, out)
	}

	destPath := filepath.Join(t.TempDir(), "existing-wt")

	err := AddWorktree(repoDir, "existing-branch", destPath)
	if err != nil {
		t.Fatalf("AddWorktree for existing branch failed: %v", err)
	}

	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	found := false
	for _, e := range entries {
		if e.Branch == "existing-branch" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find 'existing-branch' worktree")
	}
}

func TestMoveWorktree_StaleDestinationRetry(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create a worktree
	origPath := filepath.Join(t.TempDir(), "retry-wt")
	resolvedOrig, _ := filepath.EvalSymlinks(filepath.Dir(origPath))
	origPath = filepath.Join(resolvedOrig, "retry-wt")

	cmd := exec.Command("git", "worktree", "add", "-b", "retry-branch", origPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	newDir := t.TempDir()
	resolvedNew, _ := filepath.EvalSymlinks(newDir)
	newPath := filepath.Join(resolvedNew, "retry-wt")

	// Simulate a stale destination from a prior failed attempt: create a file
	// at the exact destination path so git worktree move fails with "exists".
	if err := os.WriteFile(newPath, []byte("stale"), 0644); err != nil {
		t.Fatalf("failed to create stale destination: %v", err)
	}

	// MoveWorktree should clean up the stale destination and retry successfully
	err := MoveWorktree(repoDir, origPath, newPath)
	if err != nil {
		t.Fatalf("MoveWorktree should recover from stale destination, got: %v", err)
	}

	// Source should be gone
	if _, err := os.Stat(origPath); err == nil {
		t.Error("original path should not exist after successful retry")
	}

	// Destination should exist with worktree contents
	if _, err := os.Stat(newPath); err != nil {
		t.Errorf("destination should exist after successful retry: %v", err)
	}

	// Git should know the worktree at the new location
	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	found := false
	for _, e := range entries {
		if e.Branch == "retry-branch" {
			found = true
			if filepath.Clean(e.Path) != filepath.Clean(newPath) {
				t.Errorf("expected worktree at %s, got %s", newPath, e.Path)
			}
		}
	}
	if !found {
		t.Error("expected to find 'retry-branch' worktree after retry")
	}
}

func TestMoveWorktree_PartialMoveRollback(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create a worktree
	origPath := filepath.Join(t.TempDir(), "partial-wt")
	resolvedOrig, _ := filepath.EvalSymlinks(filepath.Dir(origPath))
	origPath = filepath.Join(resolvedOrig, "partial-wt")

	cmd := exec.Command("git", "worktree", "add", "-b", "partial-branch", origPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	// Simulate a partial move: physically move the directory to the destination
	// but DON'T update git internals. This means git worktree move will fail
	// because the source is gone, and MoveWorktree should roll back.
	newDir := t.TempDir()
	resolvedNew, _ := filepath.EvalSymlinks(newDir)
	newPath := filepath.Join(resolvedNew, "partial-dest")

	// Physically move the directory (simulating a partial git worktree move)
	if err := os.Rename(origPath, newPath); err != nil {
		t.Fatalf("failed to simulate partial move: %v", err)
	}

	// Now call MoveWorktree — git worktree move will fail because source is gone.
	// The recovery code should detect source gone + dest exists and roll back.
	err := MoveWorktree(repoDir, origPath, newPath)
	if err == nil {
		t.Fatal("MoveWorktree should return an error for partial move scenario")
	}

	// The error message should mention rollback
	if !strings.Contains(err.Error(), "rolled back") {
		t.Errorf("error should mention rollback, got: %v", err)
	}

	// The worktree should be back at the original location
	if _, err := os.Stat(origPath); err != nil {
		t.Errorf("worktree should be rolled back to original path: %v", err)
	}

	// Destination should be gone after rollback
	if _, err := os.Stat(newPath); err == nil {
		t.Error("destination should not exist after rollback")
	}
}

func TestMoveWorktree_LockedWorktreeReturnsError(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create a worktree
	origPath := filepath.Join(t.TempDir(), "locked-wt")
	resolvedOrig, _ := filepath.EvalSymlinks(filepath.Dir(origPath))
	origPath = filepath.Join(resolvedOrig, "locked-wt")

	cmd := exec.Command("git", "worktree", "add", "-b", "locked-branch", origPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	// Lock the worktree
	cmd = exec.Command("git", "worktree", "lock", origPath, "--reason", "test lock")
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree lock failed: %s\n%s", err, out)
	}

	newDir := t.TempDir()
	resolvedNew, _ := filepath.EvalSymlinks(newDir)
	newPath := filepath.Join(resolvedNew, "locked-dest")

	// MoveWorktree should fail with a descriptive error
	err := MoveWorktree(repoDir, origPath, newPath)
	if err == nil {
		t.Fatal("MoveWorktree should fail for locked worktree")
	}

	// Error should contain git's reason
	if !strings.Contains(err.Error(), "lock") {
		t.Errorf("error should mention lock, got: %v", err)
	}

	// Original should still be intact
	if _, err := os.Stat(origPath); err != nil {
		t.Errorf("locked worktree should remain at original path: %v", err)
	}
}

func TestMoveWorktree(t *testing.T) {
	repoDir := t.TempDir()
	initBareGitRepo(t, repoDir)

	// Create a worktree
	origPath := filepath.Join(t.TempDir(), "orig-wt")
	cmd := exec.Command("git", "worktree", "add", "-b", "move-me", origPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	// Move it
	newDir := t.TempDir()
	resolvedNew, _ := filepath.EvalSymlinks(newDir)
	newPath := filepath.Join(resolvedNew, "moved-wt")

	err := MoveWorktree(repoDir, origPath, newPath)
	if err != nil {
		t.Fatalf("MoveWorktree failed: %v", err)
	}

	// Old path should not exist
	if _, err := os.Stat(origPath); err == nil {
		t.Error("old worktree path should not exist after move")
	}

	// New path should exist
	if _, err := os.Stat(newPath); err != nil {
		t.Errorf("new worktree path should exist after move: %v", err)
	}

	// Verify git still knows about it
	entries, err := ListWorktrees(repoDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	found := false
	for _, e := range entries {
		if e.Branch == "move-me" {
			found = true
			if filepath.Clean(e.Path) != filepath.Clean(newPath) {
				t.Errorf("expected worktree at %s, got %s", newPath, e.Path)
			}
		}
	}
	if !found {
		t.Error("expected to find 'move-me' worktree after move")
	}
}
