// Package reconcile brings gws workspace metadata back into agreement with
// what is actually on disk. It is the single engine behind both "gws init"
// and "gws refresh": the two commands differ in lifecycle, not in discovery
// or merge semantics.
//
// Reconciliation is deliberately free of side effects beyond reading the
// filesystem. It does not clear the git status cache and does not create,
// repair, or remove workspace symlinks; those responsibilities belong to the
// calling command.
package reconcile

import "github.com/daileyo/gws/internal/config"

// Result is the outcome of a reconciliation pass.
type Result struct {
	// Repositories is the reconciled repository set, sorted by path.
	Repositories []config.Repository

	// TotalRepositories is len(Repositories), carried explicitly so callers
	// can report it without recomputing.
	TotalRepositories int
	// Added counts repositories present after reconciliation that were not in
	// the existing metadata.
	Added int
	// Removed counts entries dropped because their stored path no longer exists.
	Removed int
	// Updated counts retained repositories whose name, remote URL, type, or
	// visibility changed since the last pass.
	Updated int

	// ReposWithWorktrees counts repositories that have at least one worktree.
	ReposWithWorktrees int
	// TotalWorktrees counts worktrees across every repository.
	TotalWorktrees int
	// AlignedWorktrees counts worktrees inside <repo>.wt/.
	AlignedWorktrees int
	// UnalignedWorktrees counts worktrees anywhere else.
	UnalignedWorktrees int

	// RemovedRepositories holds the entries counted by Removed, so callers can
	// act on them — for example by rescanning their parent directories.
	RemovedRepositories []config.Repository

	// Errors holds non-fatal problems encountered while scanning. Callers
	// report them with ReportScanErrors.
	Errors []error

	// ReachablePaths is the set of resolved repository paths reachable from
	// the workspace root during this pass. A tracked repository absent from
	// this set still exists on disk but can no longer be reached from within
	// the workspace.
	ReachablePaths map[string]bool

	// WorktreeMainRepos is the set of main repository paths inferred from
	// linked worktrees found during this pass. A repository listed here is
	// tracked because one of its worktrees lives in the workspace, not
	// because the user added it, so callers must not treat it as one they
	// are responsible for linking.
	WorktreeMainRepos map[string]bool
}
