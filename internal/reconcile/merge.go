package reconcile

import (
	"os"
	"path/filepath"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

// mergeCounts tallies what a merge did.
type mergeCounts struct {
	added   int
	removed int
	updated int
}

// mergeRepositories reconciles the discovered repository set against existing
// metadata.
//
// Retention is decided by whether a tracked path still exists on disk, not by
// whether the scan reached it. A repository can remain valid while becoming
// unreachable from the workspace root — for example when a curated symlink is
// deleted — and dropping it in that case would silently lose the user's tags.
func mergeRepositories(existing, discovered []config.Repository) ([]config.Repository, []config.Repository, mergeCounts) {
	var counts mergeCounts
	var merged []config.Repository
	var removedRepos []config.Repository

	discoveredByPath := make(map[string]config.Repository, len(discovered))
	for _, repo := range discovered {
		discoveredByPath[repo.Path] = repo
	}

	seen := make(map[string]bool, len(existing)+len(discovered))

	// Pass 1: decide the fate of every already-tracked repository.
	for _, repo := range existing {
		if seen[repo.Path] {
			// Duplicate entry in the stored config; keep only the first.
			continue
		}
		seen[repo.Path] = true

		if !repositoryExists(repo.Path) {
			counts.removed++
			removedRepos = append(removedRepos, repo)
			continue
		}

		retained, changed := refreshMetadata(repo)
		if changed {
			counts.updated++
		}
		merged = append(merged, retained)
	}

	// Pass 2: add repositories the scan found that were not already tracked.
	for _, repo := range discovered {
		if seen[repo.Path] {
			continue
		}
		seen[repo.Path] = true
		counts.added++
		merged = append(merged, repo)
	}

	return merged, removedRepos, counts
}

// repositoryExists reports whether a tracked path still holds a repository.
func repositoryExists(repoPath string) bool {
	_, err := os.Lstat(filepath.Join(repoPath, ".git"))
	return err == nil
}

// refreshMetadata re-reads the mutable fields of a tracked repository while
// preserving everything the user owns: tags and the resolved git identity.
// It reports whether any mutable field changed.
func refreshMetadata(repo config.Repository) (config.Repository, bool) {
	updated, err := discovery.BuildRepository(repo.Path)
	if err != nil {
		// Cannot re-read it — keep the existing entry rather than losing it.
		return repo, false
	}

	changed := updated.Name != repo.Name ||
		updated.RemoteURL != repo.RemoteURL ||
		updated.Type != repo.Type ||
		updated.Visibility != repo.Visibility

	// User-owned and separately-detected fields survive reconciliation.
	updated.Tags = repo.Tags
	updated.User = repo.User
	updated.Email = repo.Email
	updated.SigningEnabled = repo.SigningEnabled
	updated.UserSource = repo.UserSource
	updated.Worktrees = repo.Worktrees

	return *updated, changed
}
