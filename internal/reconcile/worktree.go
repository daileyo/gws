package reconcile

import (
	"os"
	"path/filepath"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
	"github.com/daileyo/gws/internal/git"
)

// worktreeCounts tallies the worktree phase.
type worktreeCounts struct {
	reposWithWorktrees int
	total              int
	aligned            int
	unaligned          int
}

// discoverWorktrees populates the Worktrees field on each repository by asking
// git, which is the authority on what worktrees exist. Repositories are
// repaired and pruned first so entries with fixable links survive and truly
// dead ones are dropped.
func discoverWorktrees(repos []config.Repository, progress ProgressReporter) worktreeCounts {
	var counts worktreeCounts

	for i := range repos {
		if progress != nil {
			progress.Increment()
		}

		// Repair broken links first, then prune truly dead entries.
		_ = git.RepairWorktrees(repos[i].Path)
		_ = git.PruneWorktrees(repos[i].Path)

		entries, err := git.ListWorktrees(repos[i].Path)
		if err != nil {
			continue
		}

		var wts []config.Worktree
		for _, e := range entries {
			// Skip worktrees whose path no longer exists (prunable).
			if _, statErr := os.Stat(e.Path); statErr != nil {
				continue
			}
			aligned := git.IsAligned(e.Path, repos[i].Path)
			if aligned {
				counts.aligned++
			} else {
				counts.unaligned++
			}
			counts.total++
			wts = append(wts, config.Worktree{
				Path:    e.Path,
				Branch:  e.Branch,
				Aligned: aligned,
			})
		}

		if len(wts) == 0 {
			repos[i].Worktrees = nil
			continue
		}
		repos[i].Worktrees = wts
		counts.reposWithWorktrees++
	}

	return counts
}

// registerOrphanMainRepos adds the main repository of any linked worktree that
// the scan found but whose repository is not tracked. Without this, worktrees
// discovered under the workspace root would belong to nothing and never appear
// in "gws worktree list".
//
// A main repository may live outside the workspace root; it is tracked by its
// real path and receives no workspace symlink.
func registerOrphanMainRepos(repos []config.Repository, mainRepoPaths map[string]bool) ([]config.Repository, int) {
	if len(mainRepoPaths) == 0 {
		return repos, 0
	}

	tracked := make(map[string]bool, len(repos))
	for _, repo := range repos {
		tracked[repo.Path] = true
	}

	added := 0
	for mainPath := range mainRepoPaths {
		resolved := resolveRealPath(mainPath)
		if tracked[resolved] || tracked[mainPath] {
			continue
		}
		if !repositoryExists(resolved) {
			continue
		}

		repo, err := discovery.BuildRepository(resolved)
		if err != nil {
			continue
		}

		tracked[resolved] = true
		repos = append(repos, *repo)
		added++
	}

	return repos, added
}

// resolveRealPath resolves symlinks, falling back to a cleaned path when the
// target cannot be resolved.
func resolveRealPath(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return resolved
}
