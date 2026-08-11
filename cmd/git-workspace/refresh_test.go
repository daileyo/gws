package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
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

	if err := runRefresh(io.Discard); err != nil {
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

	if err := runRefresh(io.Discard); err != nil {
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

	if err := runRefresh(io.Discard); err != nil {
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

	if err := runRefresh(io.Discard); err != nil {
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
	if err := runRefresh(io.Discard); err != nil {
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

	if err := runRefresh(io.Discard); err != nil {
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

	if err := runRefresh(io.Discard); err != nil {
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

// --- Spec 21: unified reconciliation ---

// TestRunRefresh_DiscoversDeeplyNestedRepository proves the one-level-deep
// scan is gone: a repository added three container directories deep after
// initialization is found by a plain refresh.
func TestRunRefresh_DiscoversDeeplyNestedRepository(t *testing.T) {
	f := newWorkspaceFixture(t)

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	deep := fixtureRepo(t, filepath.Join(f.Root, "org-b", "team-x", "squad-y", "late-arrival"))

	var out bytes.Buffer
	if err := runRefresh(&out); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	found := false
	for _, repo := range cfg.Repositories {
		if repo.Path == deep {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected the deeply nested repository %s to be discovered", deep)
	}
	if !strings.Contains(out.String(), "Found 1 new repository") {
		t.Errorf("Expected the summary to report 1 new repository, got:\n%s", out.String())
	}
}

func TestRunRefresh_RemovesAndCountsDeletedRepository(t *testing.T) {
	f := newWorkspaceFixture(t)

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	doomed := f.TopLevel[0]
	if err := os.RemoveAll(doomed); err != nil {
		t.Fatalf("Failed to remove repository: %v", err)
	}

	var out bytes.Buffer
	if err := runRefresh(&out); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	if !strings.Contains(out.String(), "Removed 1 repository (path no longer valid)") {
		t.Errorf("Expected the removal to be reported, got:\n%s", out.String())
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	for _, repo := range cfg.Repositories {
		if repo.Path == doomed {
			t.Errorf("Expected %s to be removed from config", doomed)
		}
	}
}

// TestRunRefresh_NewWorktreeVisibleImmediately covers a worktree created after
// the previous refresh appearing after exactly one refresh.
func TestRunRefresh_NewWorktreeVisibleImmediately(t *testing.T) {
	f := newWorkspaceFixture(t)

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	fixtureWorktree(t, f.TopLevel[0], "brand-new", filepath.Join(f.Root, "alpha.wt", "brand-new"))

	if err := runRefresh(io.Discard); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	var listOut bytes.Buffer
	if err := runWorktreeList("", &listOut); err != nil {
		t.Fatalf("runWorktreeList failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "brand-new") {
		t.Errorf("Expected the new worktree in the list after one refresh, got:\n%s", listOut.String())
	}
}

// TestRunRefresh_IsIdempotent is the regression test for symlink accumulation:
// repeated refreshes of an unchanged curated workspace must change nothing.
func TestRunRefresh_IsIdempotent(t *testing.T) {
	root, _ := newSymlinkOnlyWorkspace(t, "apps", "kubernetes", "tooling")

	if err := runInit(root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := runRefresh(io.Discard); err != nil {
		t.Fatalf("First runRefresh failed: %v", err)
	}

	listingBefore := workspaceListing(t, root)
	configBefore, err := os.ReadFile(configPathForTest(t))
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	if err := runRefresh(io.Discard); err != nil {
		t.Fatalf("Second runRefresh failed: %v", err)
	}

	listingAfter := workspaceListing(t, root)
	configAfter, err := os.ReadFile(configPathForTest(t))
	if err != nil {
		t.Fatalf("Failed to re-read config: %v", err)
	}

	if len(listingBefore) != len(listingAfter) {
		t.Fatalf("Workspace changed between refreshes.\nbefore: %v\nafter:  %v", listingBefore, listingAfter)
	}
	for i := range listingBefore {
		if listingBefore[i] != listingAfter[i] {
			t.Errorf("Workspace entry %d changed: %q -> %q", i, listingBefore[i], listingAfter[i])
		}
	}
	if string(configBefore) != string(configAfter) {
		t.Error("config.json changed between two refreshes of an unchanged workspace")
	}
}

// TestRunRefresh_NoDuplicateLinkForNestedSymlink is the accumulation bug in its
// original form: a repository reached through a symlink below the workspace
// root has an external real path, and must not gain a second link at the root.
func TestRunRefresh_NoDuplicateLinkForNestedSymlink(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	external := resolvedTempDir(t)

	extRepo := fixtureRepo(t, filepath.Join(external, "nested-target"))
	if err := os.MkdirAll(filepath.Join(root, "grouping"), 0755); err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	if err := os.Symlink(extRepo, filepath.Join(root, "grouping", "nested-target")); err != nil {
		t.Fatalf("Failed to create nested symlink: %v", err)
	}

	if err := runInit(root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := runRefresh(io.Discard); err != nil {
			t.Fatalf("runRefresh %d failed: %v", i+1, err)
		}
	}

	if _, err := os.Lstat(filepath.Join(root, "nested-target")); err == nil {
		t.Error("refresh created a duplicate symlink at the workspace root for a repository already reachable through a nested symlink")
	}

	links := symlinkEntries(workspaceListing(t, root))
	if len(links) != 1 {
		t.Errorf("Expected exactly the original symlink, got %v", links)
	}
}

// TestRunRefresh_RepairsMissingLinkForTrackedExternalRepo covers the other half
// of the repair-only rule: a tracked external repository that became
// unreachable gets its link restored.
func TestRunRefresh_RepairsMissingLinkForTrackedExternalRepo(t *testing.T) {
	root, repoPaths := newSymlinkOnlyWorkspace(t, "linked-repo")

	if err := runInit(root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	linkPath := filepath.Join(root, "linked-repo")
	if err := os.Remove(linkPath); err != nil {
		t.Fatalf("Failed to remove symlink: %v", err)
	}

	if err := runRefresh(io.Discard); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("Expected the workspace symlink to be repaired, got: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("Expected a symlink at the repaired path")
	}
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("Failed to read repaired symlink: %v", err)
	}
	if target != repoPaths[0] {
		t.Errorf("Repaired symlink points at %s, want %s", target, repoPaths[0])
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Expected the repository to remain tracked, got %d", len(cfg.Repositories))
	}
}

func TestRunRefresh_ClearsStatusCache(t *testing.T) {
	f := newWorkspaceFixture(t)

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Seed the cache so clearing it is observable.
	cachePath, err := git.GetCachePath()
	if err != nil {
		t.Fatalf("Failed to resolve cache path: %v", err)
	}
	seeded := git.NewCache(git.DefaultTTL)
	seeded.Set(f.TopLevel[0], &git.Status{Branch: "main"})
	if err := seeded.Save(cachePath); err != nil {
		t.Fatalf("Failed to save seeded cache: %v", err)
	}

	var out bytes.Buffer
	if err := runRefresh(&out); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	if !strings.Contains(out.String(), "Cleared git status cache") {
		t.Errorf("Expected the cache-clear message, got:\n%s", out.String())
	}

	reloaded := git.NewCache(git.DefaultTTL)
	if err := reloaded.Load(cachePath); err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}
	if got := reloaded.GetStale(f.TopLevel[0]); got != nil {
		t.Errorf("Expected the cache to be cleared, still found %+v", got)
	}
}

// TestRunRefresh_ParentRescanSafetyNet covers a repository replaced in place:
// the old checkout is deleted and a different one appears at the same path.
func TestRunRefresh_ParentRescanSafetyNet(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	container := filepath.Join(root, "projects")
	original := fixtureRepo(t, filepath.Join(container, "original"))

	if err := runInit(root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	if err := os.RemoveAll(original); err != nil {
		t.Fatalf("Failed to remove repository: %v", err)
	}
	replacement := fixtureRepo(t, filepath.Join(container, "replacement"))

	if err := runRefresh(io.Discard); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	paths := map[string]bool{}
	for _, repo := range cfg.Repositories {
		paths[repo.Path] = true
	}
	if !paths[replacement] {
		t.Errorf("Expected the replacement repository %s to be tracked, got %v", replacement, paths)
	}
	if paths[original] {
		t.Errorf("Expected the removed repository %s to be gone, got %v", original, paths)
	}
}

// TestInitAndRefresh_ReportIdenticalTotals is the central parity claim of the
// spec: the same workspace yields the same model from either command.
func TestInitAndRefresh_ReportIdenticalTotals(t *testing.T) {
	f := newWorkspaceFixture(t)

	var initOut bytes.Buffer
	if err := runInit(f.Root, &initOut); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	initCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after init: %v", err)
	}
	initRepos := len(initCfg.Repositories)
	initWorktrees := countWorktrees(initCfg)

	var refreshOut bytes.Buffer
	if err := runRefresh(&refreshOut); err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	refreshCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after refresh: %v", err)
	}
	refreshRepos := len(refreshCfg.Repositories)
	refreshWorktrees := countWorktrees(refreshCfg)

	if initRepos != refreshRepos {
		t.Errorf("Repository totals disagree: init=%d refresh=%d", initRepos, refreshRepos)
	}
	if initWorktrees != refreshWorktrees {
		t.Errorf("Worktree totals disagree: init=%d refresh=%d", initWorktrees, refreshWorktrees)
	}

	// The shared worktree summary line must be byte-identical in both reports.
	wantLine := "Worktrees: 2 (1 aligned, 1 unaligned)"
	if !strings.Contains(initOut.String(), wantLine) {
		t.Errorf("init output missing %q:\n%s", wantLine, initOut.String())
	}
	if !strings.Contains(refreshOut.String(), wantLine) {
		t.Errorf("refresh output missing %q:\n%s", wantLine, refreshOut.String())
	}

	initTotal := "Found 5 repositories."
	refreshTotal := "Total repositories: 5"
	if !strings.Contains(initOut.String(), initTotal) {
		t.Errorf("init output missing %q:\n%s", initTotal, initOut.String())
	}
	if !strings.Contains(refreshOut.String(), refreshTotal) {
		t.Errorf("refresh output missing %q:\n%s", refreshTotal, refreshOut.String())
	}
}

// TestInitAndRefresh_IdenticalScanErrorOutput pins the unified reporter: the
// same scan problems must read identically from either command, on stderr.
func TestInitAndRefresh_IdenticalScanErrorOutput(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	fixtureRepo(t, filepath.Join(root, "good-repo"))

	// Several dangling symlinks so the truncation branch is exercised too.
	for i := 0; i < 7; i++ {
		name := fmt.Sprintf("dangling-%d", i)
		if err := os.Symlink(filepath.Join(root, "nowhere", name), filepath.Join(root, name)); err != nil {
			t.Fatalf("Failed to create dangling symlink: %v", err)
		}
	}

	var initOut bytes.Buffer
	initStderr, err := captureStderr(t, func() error {
		return runInit(root, &initOut)
	})
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	var refreshOut bytes.Buffer
	refreshStderr, err := captureStderr(t, func() error {
		return runRefresh(&refreshOut)
	})
	if err != nil {
		t.Fatalf("runRefresh failed: %v", err)
	}

	if initStderr != refreshStderr {
		t.Errorf("Scan-error output differs between commands.\ninit:\n%s\nrefresh:\n%s", initStderr, refreshStderr)
	}
	if !strings.Contains(initStderr, "Warning: 7 errors occurred during scanning:") {
		t.Errorf("Expected the shared header, got:\n%s", initStderr)
	}
	if !strings.Contains(initStderr, "... and 2 more errors") {
		t.Errorf("Expected truncation after 5 errors, got:\n%s", initStderr)
	}

	// Warnings belong on stderr, never stdout.
	for name, out := range map[string]string{"init": initOut.String(), "refresh": refreshOut.String()} {
		if strings.Contains(out, "errors occurred during scanning:") {
			t.Errorf("%s wrote scan warnings to stdout:\n%s", name, out)
		}
	}
}

// countWorktrees totals worktrees across a configuration.
func countWorktrees(cfg *config.Config) int {
	total := 0
	for _, repo := range cfg.Repositories {
		total += len(repo.Worktrees)
	}
	return total
}
