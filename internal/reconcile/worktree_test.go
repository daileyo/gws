package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestDiscoverWorktrees_AlignmentClassification(t *testing.T) {
	workspace := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(workspace, "svc"))

	// Aligned: inside <repo>.wt/. Unaligned: anywhere else.
	addWorktree(t, repoPath, "feature", filepath.Join(workspace, "svc.wt", "feature"))
	addWorktree(t, repoPath, "hotfix", filepath.Join(workspace, "elsewhere", "hotfix"))

	repos := []config.Repository{{Name: "svc", Path: repoPath}}
	counts := discoverWorktrees(repos, nil)

	if counts.reposWithWorktrees != 1 {
		t.Errorf("reposWithWorktrees = %d, want 1", counts.reposWithWorktrees)
	}
	if counts.total != 2 {
		t.Fatalf("total = %d, want 2", counts.total)
	}
	if counts.aligned != 1 {
		t.Errorf("aligned = %d, want 1", counts.aligned)
	}
	if counts.unaligned != 1 {
		t.Errorf("unaligned = %d, want 1", counts.unaligned)
	}

	if len(repos[0].Worktrees) != 2 {
		t.Fatalf("Expected 2 worktrees stored, got %d", len(repos[0].Worktrees))
	}

	byBranch := map[string]config.Worktree{}
	for _, wt := range repos[0].Worktrees {
		byBranch[wt.Branch] = wt
	}
	if !byBranch["feature"].Aligned {
		t.Error("Expected the worktree under svc.wt/ to be aligned")
	}
	if byBranch["hotfix"].Aligned {
		t.Error("Expected the worktree outside svc.wt/ to be unaligned")
	}
}

func TestDiscoverWorktrees_NoWorktrees(t *testing.T) {
	workspace := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(workspace, "plain"))

	repos := []config.Repository{{Name: "plain", Path: repoPath}}
	counts := discoverWorktrees(repos, nil)

	if counts.total != 0 || counts.reposWithWorktrees != 0 {
		t.Errorf("Expected no worktrees, got total=%d repos=%d", counts.total, counts.reposWithWorktrees)
	}
	if repos[0].Worktrees != nil {
		t.Errorf("Expected Worktrees to be nil, got %v", repos[0].Worktrees)
	}
}

// TestDiscoverWorktrees_ClearsStaleEntries verifies that a worktree removed
// from disk since the last pass is dropped from metadata.
func TestDiscoverWorktrees_ClearsStaleEntries(t *testing.T) {
	workspace := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(workspace, "svc"))
	wtPath := filepath.Join(workspace, "svc.wt", "feature")
	addWorktree(t, repoPath, "feature", wtPath)

	repos := []config.Repository{{Name: "svc", Path: repoPath}}
	if counts := discoverWorktrees(repos, nil); counts.total != 1 {
		t.Fatalf("Setup failed: expected 1 worktree, got %d", counts.total)
	}

	if err := os.RemoveAll(wtPath); err != nil {
		t.Fatalf("Failed to remove worktree: %v", err)
	}

	counts := discoverWorktrees(repos, nil)
	if counts.total != 0 {
		t.Errorf("total = %d, want 0 after the worktree directory was deleted", counts.total)
	}
	if repos[0].Worktrees != nil {
		t.Errorf("Expected Worktrees cleared, got %v", repos[0].Worktrees)
	}
}

func TestDiscoverWorktrees_IncrementsProgressPerRepository(t *testing.T) {
	workspace := t.TempDir()
	repos := []config.Repository{
		{Name: "one", Path: mkRepo(t, filepath.Join(workspace, "one"))},
		{Name: "two", Path: mkRepo(t, filepath.Join(workspace, "two"))},
		{Name: "three", Path: mkRepo(t, filepath.Join(workspace, "three"))},
	}

	spy := &progressSpy{}
	discoverWorktrees(repos, spy)

	if spy.increments != 3 {
		t.Errorf("Increment called %d times, want 3", spy.increments)
	}
}

// TestRegisterOrphanMainRepos_InsideWorkspace covers a worktree found under the
// workspace root whose main repository is not yet tracked.
func TestRegisterOrphanMainRepos_InsideWorkspace(t *testing.T) {
	workspace := t.TempDir()
	mainRepo := mkRepo(t, filepath.Join(workspace, "orphan-main"))

	repos, added := registerOrphanMainRepos(nil, map[string]bool{mainRepo: true})

	if added != 1 {
		t.Fatalf("added = %d, want 1", added)
	}
	if len(repos) != 1 || repos[0].Path != mainRepo {
		t.Errorf("Expected the main repository registered, got %v", repos)
	}
	if repos[0].Name != "orphan-main" {
		t.Errorf("Expected name 'orphan-main', got %q", repos[0].Name)
	}
}

// TestRegisterOrphanMainRepos_OutsideWorkspace covers the case the spec calls
// out explicitly: the main repository may live outside the workspace root and
// is still tracked, by its real path.
func TestRegisterOrphanMainRepos_OutsideWorkspace(t *testing.T) {
	external := t.TempDir()
	mainRepo := mkRepo(t, filepath.Join(external, "external-main"))

	repos, added := registerOrphanMainRepos(nil, map[string]bool{mainRepo: true})

	if added != 1 {
		t.Fatalf("added = %d, want 1", added)
	}
	if repos[0].Path != mainRepo {
		t.Errorf("Expected path %s, got %s", mainRepo, repos[0].Path)
	}
}

func TestRegisterOrphanMainRepos_SkipsAlreadyTracked(t *testing.T) {
	workspace := t.TempDir()
	mainRepo := mkRepo(t, filepath.Join(workspace, "tracked"))

	existing := []config.Repository{{Name: "tracked", Path: mainRepo, Tags: []string{"mine"}}}
	repos, added := registerOrphanMainRepos(existing, map[string]bool{mainRepo: true})

	if added != 0 {
		t.Errorf("added = %d, want 0 — the repository is already tracked", added)
	}
	if len(repos) != 1 {
		t.Fatalf("Expected 1 repository, got %d", len(repos))
	}
	if len(repos[0].Tags) != 1 || repos[0].Tags[0] != "mine" {
		t.Errorf("Expected the tracked entry untouched, got %v", repos[0])
	}
}

func TestRegisterOrphanMainRepos_SkipsMissingPaths(t *testing.T) {
	base := t.TempDir()
	missing := filepath.Join(base, "not-there")

	repos, added := registerOrphanMainRepos(nil, map[string]bool{missing: true})

	if added != 0 {
		t.Errorf("added = %d, want 0", added)
	}
	if len(repos) != 0 {
		t.Errorf("Expected no repositories, got %v", repos)
	}
}

func TestRegisterOrphanMainRepos_EmptyInput(t *testing.T) {
	existing := []config.Repository{{Name: "a", Path: "/somewhere"}}
	repos, added := registerOrphanMainRepos(existing, nil)

	if added != 0 {
		t.Errorf("added = %d, want 0", added)
	}
	if len(repos) != 1 {
		t.Errorf("Expected the input returned unchanged, got %v", repos)
	}
}

// TestReconcileWorkspace_RegistersOrphanMainRepoEndToEnd exercises orphan
// registration through the full engine rather than the helper alone.
func TestReconcileWorkspace_RegistersOrphanMainRepoEndToEnd(t *testing.T) {
	workspace := t.TempDir()
	external := t.TempDir()

	// Main repository outside the workspace; only its worktree is inside.
	mainRepo := mkRepo(t, filepath.Join(external, "proj"))
	addWorktree(t, mainRepo, "feature", filepath.Join(workspace, "proj.wt", "feature"))

	result, err := ReconcileWorkspace(workspace, nil, Options{})
	if err != nil {
		t.Fatalf("ReconcileWorkspace failed: %v", err)
	}

	repo := repoByPath(result, mainRepo)
	if repo == nil {
		t.Fatalf("Expected the external main repository %s to be registered, got %v",
			mainRepo, repoPathSet(result))
	}
	if len(repo.Worktrees) != 1 {
		t.Fatalf("Expected 1 worktree on the registered repository, got %d", len(repo.Worktrees))
	}
	if repo.Worktrees[0].Branch != "feature" {
		t.Errorf("Expected branch 'feature', got %q", repo.Worktrees[0].Branch)
	}
	if result.TotalWorktrees != 1 {
		t.Errorf("TotalWorktrees = %d, want 1", result.TotalWorktrees)
	}
}
