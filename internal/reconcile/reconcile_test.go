package reconcile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// TestReconcileWorkspace_InitPath covers the "gws init" case: no existing
// metadata, so everything discovered is new.
func TestReconcileWorkspace_InitPath(t *testing.T) {
	workspace := t.TempDir()

	repoA := mkRepo(t, filepath.Join(workspace, "alpha"))
	repoB := mkRepo(t, filepath.Join(workspace, "group", "beta"))
	repoC := mkRepo(t, filepath.Join(workspace, "group", "sub", "gamma"))

	result, err := ReconcileWorkspace(workspace, nil, Options{})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	if result.TotalRepositories != 3 {
		t.Errorf("TotalRepositories = %d, want 3 (%v)", result.TotalRepositories, repoPathSet(result))
	}
	if result.Added != 3 {
		t.Errorf("Added = %d, want 3", result.Added)
	}
	if result.Removed != 0 {
		t.Errorf("Removed = %d, want 0", result.Removed)
	}
	if result.Updated != 0 {
		t.Errorf("Updated = %d, want 0", result.Updated)
	}

	paths := repoPathSet(result)
	for _, want := range []string{repoA, repoB, repoC} {
		if !paths[want] {
			t.Errorf("Expected repository %s in result, got %v", want, paths)
		}
	}
}

func TestReconcileWorkspace_RequiresWorkspaceRoot(t *testing.T) {
	if _, err := ReconcileWorkspace("", nil, Options{}); err == nil {
		t.Error("Expected an error for an empty workspace root")
	}
}

func TestReconcileWorkspace_ResultsAreSortedByPath(t *testing.T) {
	workspace := t.TempDir()

	mkRepo(t, filepath.Join(workspace, "zulu"))
	mkRepo(t, filepath.Join(workspace, "alpha"))
	mkRepo(t, filepath.Join(workspace, "mike"))

	result, err := ReconcileWorkspace(workspace, nil, Options{})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	paths := make([]string, len(result.Repositories))
	for i, repo := range result.Repositories {
		paths[i] = repo.Path
	}
	if !sort.StringsAreSorted(paths) {
		t.Errorf("Expected repositories sorted by path, got %v", paths)
	}
}

// TestReconcileWorkspace_DeterministicAcrossRuns guards the idempotence
// property that config.json depends on.
func TestReconcileWorkspace_DeterministicAcrossRuns(t *testing.T) {
	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "one"))
	mkRepo(t, filepath.Join(workspace, "group", "two"))
	mkRepo(t, filepath.Join(workspace, "three"))

	first, err := ReconcileWorkspace(workspace, nil, Options{})
	if err != nil {
		t.Fatalf("first reconcile failed: %v", err)
	}

	cfg := &config.Config{Workspace: workspace, Repositories: first.Repositories}
	second, err := ReconcileWorkspace(workspace, cfg, Options{})
	if err != nil {
		t.Fatalf("second reconcile failed: %v", err)
	}

	if len(first.Repositories) != len(second.Repositories) {
		t.Fatalf("Repository count changed between runs: %d then %d",
			len(first.Repositories), len(second.Repositories))
	}
	for i := range first.Repositories {
		if first.Repositories[i].Path != second.Repositories[i].Path {
			t.Errorf("Order changed at index %d: %s then %s",
				i, first.Repositories[i].Path, second.Repositories[i].Path)
		}
	}
	if second.Added != 0 || second.Removed != 0 {
		t.Errorf("Second run should be a no-op, got added=%d removed=%d",
			second.Added, second.Removed)
	}
}

// TestReconcileWorkspace_NoSideEffects asserts the engine stays pure: it must
// not touch the status cache or write anything into the workspace.
func TestReconcileWorkspace_NoSideEffects(t *testing.T) {
	workspace := t.TempDir()
	external := t.TempDir()

	mkRepo(t, filepath.Join(workspace, "inside"))
	extRepo := mkRepo(t, filepath.Join(external, "outside"))
	symlink(t, extRepo, filepath.Join(workspace, "outside"))

	before := listDirRecursive(t, workspace)

	// A sentinel standing in for the status cache file.
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "status-cache.json")
	if err := os.WriteFile(cachePath, []byte(`{"sentinel":true}`), 0600); err != nil {
		t.Fatalf("Failed to write cache sentinel: %v", err)
	}
	cacheBefore, err := os.ReadFile(cachePath) // #nosec G304 -- test-controlled path
	if err != nil {
		t.Fatalf("Failed to read cache sentinel: %v", err)
	}

	if _, err := ReconcileWorkspace(workspace, nil, Options{}); err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	after := listDirRecursive(t, workspace)
	if len(before) != len(after) {
		t.Errorf("Workspace directory changed.\nbefore: %v\nafter:  %v", before, after)
	} else {
		for i := range before {
			if before[i] != after[i] {
				t.Errorf("Workspace entry %d changed: %q -> %q", i, before[i], after[i])
			}
		}
	}

	cacheAfter, err := os.ReadFile(cachePath) // #nosec G304 -- test-controlled path
	if err != nil {
		t.Fatalf("Failed to re-read cache sentinel: %v", err)
	}
	if string(cacheBefore) != string(cacheAfter) {
		t.Error("Reconciliation modified the status cache sentinel")
	}
}

// TestReconcileWorkspace_ResultContract asserts every count in Result is
// populated coherently for a workspace exercising all of them.
func TestReconcileWorkspace_ResultContract(t *testing.T) {
	workspace := t.TempDir()

	kept := mkRepo(t, filepath.Join(workspace, "kept"))
	fresh := mkRepo(t, filepath.Join(workspace, "fresh"))
	withWT := mkRepo(t, filepath.Join(workspace, "with-worktrees"))

	addWorktree(t, withWT, "feature", filepath.Join(workspace, "with-worktrees.wt", "feature"))
	addWorktree(t, withWT, "hotfix", filepath.Join(workspace, "loose", "hotfix"))

	gone := filepath.Join(workspace, "gone")
	mkRepo(t, gone)
	goneReal := realPath(t, gone)
	if err := os.RemoveAll(gone); err != nil {
		t.Fatalf("Failed to remove repository: %v", err)
	}

	existing := &config.Config{
		Workspace: workspace,
		Repositories: []config.Repository{
			{Name: "kept", Path: kept, Tags: []string{"keep-me"}},
			{Name: "gone", Path: goneReal, Tags: []string{"lost"}},
			{Name: "with-worktrees", Path: withWT},
		},
	}

	result, err := ReconcileWorkspace(workspace, existing, Options{})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	if result.TotalRepositories != 3 {
		t.Errorf("TotalRepositories = %d, want 3 (%v)", result.TotalRepositories, repoPathSet(result))
	}
	if result.Added != 1 {
		t.Errorf("Added = %d, want 1 (the 'fresh' repository)", result.Added)
	}
	if result.Removed != 1 {
		t.Errorf("Removed = %d, want 1 (the deleted repository)", result.Removed)
	}
	if len(result.RemovedRepositories) != 1 || result.RemovedRepositories[0].Path != goneReal {
		t.Errorf("RemovedRepositories = %v, want one entry for %s", result.RemovedRepositories, goneReal)
	}
	if result.ReposWithWorktrees != 1 {
		t.Errorf("ReposWithWorktrees = %d, want 1", result.ReposWithWorktrees)
	}
	if result.TotalWorktrees != 2 {
		t.Errorf("TotalWorktrees = %d, want 2", result.TotalWorktrees)
	}
	if result.AlignedWorktrees != 1 {
		t.Errorf("AlignedWorktrees = %d, want 1", result.AlignedWorktrees)
	}
	if result.UnalignedWorktrees != 1 {
		t.Errorf("UnalignedWorktrees = %d, want 1", result.UnalignedWorktrees)
	}
	if result.AlignedWorktrees+result.UnalignedWorktrees != result.TotalWorktrees {
		t.Error("aligned + unaligned must equal TotalWorktrees")
	}
	if !result.ReachablePaths[fresh] {
		t.Errorf("Expected %s in ReachablePaths, got %v", fresh, result.ReachablePaths)
	}
	if result.ReachablePaths[goneReal] {
		t.Error("A deleted repository must not appear in ReachablePaths")
	}
}

// TestReconcileWorkspace_ProgressLifecycle asserts the reporter is started,
// incremented once per repository, and stopped.
func TestReconcileWorkspace_ProgressLifecycle(t *testing.T) {
	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "one"))
	mkRepo(t, filepath.Join(workspace, "two"))

	spy := &progressSpy{}
	result, err := ReconcileWorkspace(workspace, nil, Options{Progress: spy})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	if spy.starts != 1 {
		t.Errorf("Start called %d times, want 1", spy.starts)
	}
	if spy.stops != 1 {
		t.Errorf("Stop called %d times, want 1", spy.stops)
	}
	if spy.increments != result.TotalRepositories {
		t.Errorf("Increment called %d times, want %d", spy.increments, result.TotalRepositories)
	}
}

func TestReconcileWorkspace_ScanErrorsAreReturnedNotPrinted(t *testing.T) {
	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "good"))
	symlink(t, filepath.Join(workspace, "nowhere"), filepath.Join(workspace, "dangling"))

	result, err := ReconcileWorkspace(workspace, nil, Options{})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected the dangling symlink to be reported in Result.Errors")
	}
	if result.TotalRepositories != 1 {
		t.Errorf("Expected the scan to continue, got %d repositories", result.TotalRepositories)
	}
}

// progressSpy records ProgressReporter calls.
type progressSpy struct {
	starts     int
	increments int
	stops      int
}

func (p *progressSpy) Start()     { p.starts++ }
func (p *progressSpy) Increment() { p.increments++ }
func (p *progressSpy) Stop()      { p.stops++ }

// listDirRecursive returns a sorted listing of every entry under root,
// recording symlinks as links rather than following them.
func listDirRecursive(t *testing.T, root string) []string {
	t.Helper()

	var entries []string
	var walk func(dir, prefix string)
	walk = func(dir, prefix string) {
		items, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, item := range items {
			name := filepath.Join(prefix, item.Name())
			if item.Type()&os.ModeSymlink != 0 {
				target, _ := os.Readlink(filepath.Join(dir, item.Name()))
				entries = append(entries, "link:"+name+" -> "+target)
				continue
			}
			entries = append(entries, name)
			if item.IsDir() {
				walk(filepath.Join(dir, item.Name()), name)
			}
		}
	}
	walk(root, "")

	sort.Strings(entries)
	return entries
}
