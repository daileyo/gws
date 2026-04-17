package git

import (
	"os/exec"
	"path/filepath"
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

