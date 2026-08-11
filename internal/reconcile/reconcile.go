package reconcile

import (
	"fmt"
	"sort"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

// ProgressReporter receives progress during the per-repository phases.
// *git.Progress satisfies it.
type ProgressReporter interface {
	Start()
	Increment()
	Stop()
}

// Options controls a reconciliation pass.
type Options struct {
	// MaxDepth is the traversal depth passed to discovery. Zero means the
	// configured default.
	MaxDepth int
	// Progress, when non-nil, is started before the per-repository phases and
	// stopped before ReconcileWorkspace returns.
	Progress ProgressReporter
}

// ReconcileWorkspace brings workspace metadata into agreement with the
// filesystem and reports what changed.
//
// existing may be nil, which is the "gws init" case: there is no prior
// metadata, so every discovered repository is new. When existing is supplied —
// the "gws refresh" case — its repositories are merged with what the scan
// found, preserving tags and detected git identity.
//
// The pass runs three phases: discover repositories, merge them against
// existing metadata, then discover each repository's worktrees.
func ReconcileWorkspace(workspaceRoot string, existing *config.Config, opts Options) (*Result, error) {
	if workspaceRoot == "" {
		return nil, fmt.Errorf("workspace root is required")
	}

	// Phase 1: discover what is on disk.
	scan, err := discovery.Scan(workspaceRoot, discovery.Options{MaxDepth: opts.MaxDepth})
	if err != nil {
		return nil, fmt.Errorf("failed to scan workspace: %w", err)
	}

	// Phase 2: merge against existing metadata.
	var existingRepos []config.Repository
	if existing != nil {
		existingRepos = existing.Repositories
	}

	merged, removedRepos, counts := mergeRepositories(existingRepos, scan.Repositories)

	// Worktrees discovered under the workspace root may belong to a repository
	// that is not tracked yet. Register those main repositories so their
	// worktrees have an owner.
	merged, orphansAdded := registerOrphanMainRepos(merged, scan.WorktreeMainRepos)
	counts.added += orphansAdded

	// Deterministic ordering keeps config.json from churning between runs.
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Path < merged[j].Path
	})

	// Phase 3: discover worktrees for every tracked repository.
	if opts.Progress != nil {
		opts.Progress.Start()
		defer opts.Progress.Stop()
	}

	wtCounts := discoverWorktrees(merged, opts.Progress)

	return &Result{
		Repositories:        merged,
		TotalRepositories:   len(merged),
		Added:               counts.added,
		Removed:             counts.removed,
		Updated:             counts.updated,
		ReposWithWorktrees:  wtCounts.reposWithWorktrees,
		TotalWorktrees:      wtCounts.total,
		AlignedWorktrees:    wtCounts.aligned,
		UnalignedWorktrees:  wtCounts.unaligned,
		RemovedRepositories: removedRepos,
		Errors:              scan.Errors,
		ReachablePaths:      scan.ReachablePaths,
	}, nil
}
