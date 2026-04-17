package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// captureStdoutStr captures stdout output from a function and returns it as a string.
func captureStdoutStr(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// setupAlignTestRepo creates a workspace with a git repo that has an unaligned worktree.
func setupAlignTestRepo(t *testing.T) (repoDir string) {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	workspaceDir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(workspaceDir)
	if err != nil {
		t.Fatalf("EvalSymlinks failed: %v", err)
	}
	workspaceDir = resolved

	repoDir = filepath.Join(workspaceDir, "my-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s\n%s", args, err, out)
		}
	}

	// Create an unaligned worktree outside the .wt dir
	unalignedDir := filepath.Join(t.TempDir(), "unaligned-wt")
	resolvedUnaligned, _ := filepath.EvalSymlinks(filepath.Dir(unalignedDir))
	unalignedDir = filepath.Join(resolvedUnaligned, "unaligned-wt")

	cmd := exec.Command("git", "worktree", "add", "-b", "feature-unaligned", unalignedDir)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add failed: %s\n%s", err, out)
	}

	cfg := config.New(workspaceDir)
	cfg.Repositories = []config.Repository{
		{
			Name: "my-repo",
			Path: repoDir,
			Worktrees: []config.Worktree{
				{Path: unalignedDir, Branch: "feature-unaligned", Aligned: false},
			},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	return repoDir
}

func TestRunWorktreeAlign_MovesUnaligned(t *testing.T) {
	repoDir := setupAlignTestRepo(t)

	output := captureStdoutStr(func() {
		if err := runWorktreeAlign("", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// Verify the worktree was moved to .wt/
	expectedPath := filepath.Join(repoDir+".wt", "feature-unaligned")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("worktree should exist at %s after align: %v", expectedPath, err)
	}

	// Verify output mentions the move
	if !strings.Contains(output, "Moving") {
		t.Errorf("output should mention 'Moving', got:\n%s", output)
	}
	if !strings.Contains(output, "Aligned 1 worktree") {
		t.Errorf("output should show alignment count, got:\n%s", output)
	}

	// Verify config was updated
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	repo := cfg.Repositories[0]
	if len(repo.Worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(repo.Worktrees))
	}
	if !repo.Worktrees[0].Aligned {
		t.Error("worktree should be aligned after move")
	}
}

func TestRunWorktreeAlign_SkipsAligned(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	workspaceDir := t.TempDir()
	resolved, _ := filepath.EvalSymlinks(workspaceDir)

	cfg := config.New(resolved)
	cfg.Repositories = []config.Repository{
		{
			Name: "my-repo",
			Path: filepath.Join(resolved, "my-repo"),
			Worktrees: []config.Worktree{
				{Path: filepath.Join(resolved, "my-repo.wt", "feat"), Branch: "feat", Aligned: true},
			},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	output := captureStdoutStr(func() {
		if err := runWorktreeAlign("", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "All worktrees are already aligned") {
		t.Errorf("should report all aligned, got:\n%s", output)
	}
}

func TestRunWorktreeAlign_DryRun(t *testing.T) {
	repoDir := setupAlignTestRepo(t)

	// Get the unaligned worktree path before dry-run
	cfg, _ := config.Load()
	unalignedPath := cfg.Repositories[0].Worktrees[0].Path

	output := captureStdoutStr(func() {
		if err := runWorktreeAlign("", true); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// Should say "Would move" not "Moving"
	if !strings.Contains(output, "Would move") {
		t.Errorf("dry-run should say 'Would move', got:\n%s", output)
	}
	if !strings.Contains(output, "Dry run") {
		t.Errorf("should mention dry run, got:\n%s", output)
	}

	// Worktree should NOT have been moved
	if _, err := os.Stat(unalignedPath); err != nil {
		t.Errorf("worktree should still exist at original path after dry-run: %v", err)
	}

	expectedDest := filepath.Join(repoDir+".wt", "feature-unaligned")
	if _, err := os.Stat(expectedDest); err == nil {
		t.Error("worktree should NOT exist at .wt/ path after dry-run")
	}
}

func TestRunWorktreeAlign_DuplicateNames(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	workspaceDir := t.TempDir()
	resolved, _ := filepath.EvalSymlinks(workspaceDir)

	repoDir := filepath.Join(resolved, "my-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s\n%s", args, err, out)
		}
	}

	// Create .wt dir with an existing aligned worktree named "feat"
	wtDir := repoDir + ".wt"
	alignedPath := filepath.Join(wtDir, "feat")
	cmd := exec.Command("git", "worktree", "add", "-b", "feat", alignedPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add aligned failed: %s\n%s", err, out)
	}

	// Create an unaligned worktree also named "feat" (different branch but same dir name)
	unalignedDir := filepath.Join(t.TempDir(), "feat")
	resolvedUnaligned, _ := filepath.EvalSymlinks(filepath.Dir(unalignedDir))
	unalignedDir = filepath.Join(resolvedUnaligned, "feat")
	cmd = exec.Command("git", "worktree", "add", "-b", "feat-v2", unalignedDir)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add unaligned failed: %s\n%s", err, out)
	}

	cfg := config.New(resolved)
	cfg.Repositories = []config.Repository{
		{
			Name: "my-repo",
			Path: repoDir,
			Worktrees: []config.Worktree{
				{Path: alignedPath, Branch: "feat", Aligned: true},
				{Path: unalignedDir, Branch: "feat-v2", Aligned: false},
			},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	output := captureStdoutStr(func() {
		if err := runWorktreeAlign("", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// The unaligned worktree "feat-v2" should be moved to .wt/feat-v2 (no conflict since branch name is different)
	if !strings.Contains(output, "Aligned 1 worktree") {
		t.Errorf("should align 1 worktree, got:\n%s", output)
	}
}

func TestRunWorktreeAlign_FilterByRepo(t *testing.T) {
	repoDir := setupAlignTestRepo(t)

	// Try aligning a different repo — should find nothing to align
	output := captureStdoutStr(func() {
		if err := runWorktreeAlign("other-repo", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "All worktrees are already aligned") {
		t.Errorf("filtering to non-matching repo should show nothing to align, got:\n%s", output)
	}

	// Now align the actual repo
	output = captureStdoutStr(func() {
		if err := runWorktreeAlign("my-repo", false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	expectedPath := filepath.Join(repoDir+".wt", "feature-unaligned")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("worktree should be moved when filtered to correct repo: %v", err)
	}
}
