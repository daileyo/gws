package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

func TestMergeRepositories(t *testing.T) {
	tests := []struct {
		name        string
		existing    []config.Repository
		discovered  []config.Repository
		wantMerged  []string
		wantAdded   int
		wantRemoved int
	}{
		{
			name:       "empty existing adds everything discovered",
			existing:   nil,
			discovered: []config.Repository{{Name: "a", Path: "PATH_A"}, {Name: "b", Path: "PATH_B"}},
			wantMerged: []string{"PATH_A", "PATH_B"},
			wantAdded:  2,
		},
		{
			name:       "already tracked repository is retained, not re-added",
			existing:   []config.Repository{{Name: "a", Path: "PATH_A"}},
			discovered: []config.Repository{{Name: "a", Path: "PATH_A"}},
			wantMerged: []string{"PATH_A"},
			wantAdded:  0,
		},
		{
			name:       "newly discovered repository is added alongside tracked ones",
			existing:   []config.Repository{{Name: "a", Path: "PATH_A"}},
			discovered: []config.Repository{{Name: "a", Path: "PATH_A"}, {Name: "b", Path: "PATH_B"}},
			wantMerged: []string{"PATH_A", "PATH_B"},
			wantAdded:  1,
		},
		{
			name:       "duplicate entries in stored config collapse to one",
			existing:   []config.Repository{{Name: "a", Path: "PATH_A"}, {Name: "a", Path: "PATH_A"}},
			discovered: []config.Repository{{Name: "a", Path: "PATH_A"}},
			wantMerged: []string{"PATH_A"},
			wantAdded:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			pathA := mkRepo(t, filepath.Join(base, "a"))
			pathB := mkRepo(t, filepath.Join(base, "b"))

			subst := func(repos []config.Repository) []config.Repository {
				out := make([]config.Repository, len(repos))
				for i, r := range repos {
					switch r.Path {
					case "PATH_A":
						r.Path = pathA
					case "PATH_B":
						r.Path = pathB
					}
					out[i] = r
				}
				return out
			}
			resolve := func(p string) string {
				if p == "PATH_A" {
					return pathA
				}
				return pathB
			}

			merged, _, counts := mergeRepositories(subst(tt.existing), subst(tt.discovered))

			if len(merged) != len(tt.wantMerged) {
				t.Fatalf("merged %d repositories, want %d", len(merged), len(tt.wantMerged))
			}
			gotPaths := make(map[string]bool, len(merged))
			for _, repo := range merged {
				gotPaths[repo.Path] = true
			}
			for _, want := range tt.wantMerged {
				if !gotPaths[resolve(want)] {
					t.Errorf("Expected %s in merged set, got %v", resolve(want), gotPaths)
				}
			}
			if counts.added != tt.wantAdded {
				t.Errorf("added = %d, want %d", counts.added, tt.wantAdded)
			}
			if counts.removed != tt.wantRemoved {
				t.Errorf("removed = %d, want %d", counts.removed, tt.wantRemoved)
			}
		})
	}
}

func TestMergeRepositories_RemovesVanishedPaths(t *testing.T) {
	base := t.TempDir()
	alive := mkRepo(t, filepath.Join(base, "alive"))

	dead := filepath.Join(base, "dead")
	mkRepo(t, dead)
	deadReal := realPath(t, dead)
	if err := os.RemoveAll(dead); err != nil {
		t.Fatalf("Failed to remove repository: %v", err)
	}

	existing := []config.Repository{
		{Name: "alive", Path: alive},
		{Name: "dead", Path: deadReal, Tags: []string{"stale"}},
	}

	merged, removedRepos, counts := mergeRepositories(existing, []config.Repository{{Name: "alive", Path: alive}})

	if counts.removed != 1 {
		t.Errorf("removed = %d, want 1", counts.removed)
	}
	if len(merged) != 1 || merged[0].Path != alive {
		t.Errorf("Expected only the surviving repository, got %v", merged)
	}
	if len(removedRepos) != 1 || removedRepos[0].Path != deadReal {
		t.Errorf("Expected the removed repository to be reported, got %v", removedRepos)
	}
}

// TestMergeRepositories_PreservesTags is the guarantee users notice most: a
// routine refresh must never lose organizational work.
func TestMergeRepositories_PreservesTags(t *testing.T) {
	base := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(base, "tagged"))

	existing := []config.Repository{{
		Name:           "tagged",
		Path:           repoPath,
		Tags:           []string{"work", "priority", "backend"},
		User:           "Test User",
		Email:          "test@example.com",
		SigningEnabled: true,
		UserSource:     config.UserSourceLocal,
	}}
	discovered := []config.Repository{{Name: "tagged", Path: repoPath, Tags: []string{}}}

	merged, _, _ := mergeRepositories(existing, discovered)

	if len(merged) != 1 {
		t.Fatalf("Expected 1 repository, got %d", len(merged))
	}
	got := merged[0]

	if len(got.Tags) != 3 {
		t.Fatalf("Expected 3 tags preserved, got %v", got.Tags)
	}
	for i, want := range []string{"work", "priority", "backend"} {
		if got.Tags[i] != want {
			t.Errorf("Tag %d = %q, want %q", i, got.Tags[i], want)
		}
	}
	if got.User != "Test User" || got.Email != "test@example.com" {
		t.Errorf("Expected git identity preserved, got user=%q email=%q", got.User, got.Email)
	}
	if !got.SigningEnabled {
		t.Error("Expected SigningEnabled to be preserved")
	}
	if got.UserSource != config.UserSourceLocal {
		t.Errorf("Expected UserSource preserved, got %q", got.UserSource)
	}
}

func TestMergeRepositories_CountsMetadataUpdates(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(repo *config.Repository)
		wantUpdated int
	}{
		{
			name:        "unchanged metadata is not counted",
			mutate:      func(*config.Repository) {},
			wantUpdated: 0,
		},
		{
			name:        "changed remote URL is counted",
			mutate:      func(r *config.Repository) { r.RemoteURL = "https://github.com/old/url.git" },
			wantUpdated: 1,
		},
		{
			name:        "changed type is counted",
			mutate:      func(r *config.Repository) { r.Type = config.TypeGitLab },
			wantUpdated: 1,
		},
		{
			// An https remote classifies as VisibilityUnknown, so the stored
			// copy must differ from that to represent a real change.
			name:        "changed visibility is counted",
			mutate:      func(r *config.Repository) { r.Visibility = config.VisibilityPrivate },
			wantUpdated: 1,
		},
		{
			name:        "changed name is counted",
			mutate:      func(r *config.Repository) { r.Name = "renamed" },
			wantUpdated: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			repoPath := mkRepoWithRemote(t, filepath.Join(base, "svc"), "https://github.com/test/svc.git")

			// Baseline must be the metadata the repository actually has on
			// disk, otherwise every field would read as changed.
			built, err := discovery.BuildRepository(repoPath)
			if err != nil {
				t.Fatalf("BuildRepository failed: %v", err)
			}
			stored := *built
			tt.mutate(&stored)

			_, _, counts := mergeRepositories([]config.Repository{stored}, nil)

			if counts.updated != tt.wantUpdated {
				t.Errorf("updated = %d, want %d", counts.updated, tt.wantUpdated)
			}
		})
	}
}

// TestMergeRepositories_RetainsUnreachableButExistingRepository covers a
// repository that still exists on disk but is no longer reachable from the
// workspace root, for example after its curated symlink was deleted. Dropping
// it would silently discard the user's tags.
func TestMergeRepositories_RetainsUnreachableButExistingRepository(t *testing.T) {
	external := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(external, "detached"))

	existing := []config.Repository{{Name: "detached", Path: repoPath, Tags: []string{"keep"}}}

	// The scan found nothing: the repository is unreachable from the workspace.
	merged, removedRepos, counts := mergeRepositories(existing, nil)

	if counts.removed != 0 {
		t.Errorf("removed = %d, want 0 — the repository still exists on disk", counts.removed)
	}
	if len(removedRepos) != 0 {
		t.Errorf("Expected no removed repositories, got %v", removedRepos)
	}
	if len(merged) != 1 || merged[0].Path != repoPath {
		t.Fatalf("Expected the repository to be retained, got %v", merged)
	}
	if len(merged[0].Tags) != 1 || merged[0].Tags[0] != "keep" {
		t.Errorf("Expected tags preserved, got %v", merged[0].Tags)
	}
}

func TestRepositoryExists(t *testing.T) {
	base := t.TempDir()
	repoPath := mkRepo(t, filepath.Join(base, "present"))

	if !repositoryExists(repoPath) {
		t.Error("Expected an existing repository to report true")
	}
	if repositoryExists(filepath.Join(base, "absent")) {
		t.Error("Expected a missing path to report false")
	}

	plain := filepath.Join(base, "plain-dir")
	if err := os.MkdirAll(plain, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if repositoryExists(plain) {
		t.Error("Expected a directory with no .git to report false")
	}
}
