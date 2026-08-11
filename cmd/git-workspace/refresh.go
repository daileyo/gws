package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
	"github.com/daileyo/gws/internal/git"
	"github.com/daileyo/gws/internal/reconcile"
)

// refreshCmd is the Cobra subcommand for refreshing repository metadata.
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh repository metadata and git status cache",
	Long: `Re-scan the workspace for repositories, update metadata, discover worktrees,
and clear the git status cache.

Refresh applies the same discovery rules as 'gws init': it recurses through
container directories, stops at repository boundaries, follows workspace
symlinks, and asks git for each repository's worktrees. New repositories are
added, repositories whose paths are gone are removed, and existing tags are
preserved.

Examples:
  gws refresh`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRefresh(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

// runRefresh reconciles the existing metadata library against the filesystem.
//
// It shares one reconciliation engine with init and adds the responsibilities
// that belong to refresh alone: clearing the git status cache, rescanning the
// parents of removed repositories, and repairing workspace symlinks.
func runRefresh(stdout io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Refreshing workspace at: %s\n", cfg.Workspace)

	// Capture which repositories gws is responsible for linking before
	// reconciliation changes the tracked set.
	linkOwned := workspaceLinkOwnedRepos(cfg)

	progress := git.NewProgressWithLabel(0, "Scanning workspace...")
	result, err := reconcile.ReconcileWorkspace(cfg.Workspace, cfg, reconcile.Options{
		MaxDepth: cfg.EffectiveScanMaxDepth(),
		Progress: progress,
	})
	if err != nil {
		return err
	}

	reconcile.ReportScanErrors(os.Stderr, result.Errors)

	repos := result.Repositories

	// Safety net: a repository replaced in place at the same location is
	// picked up by rescanning the parent of anything that disappeared.
	recovered := rescanRemovedParents(cfg, result)
	if len(recovered) > 0 {
		repos = append(repos, recovered...)
		sort.Slice(repos, func(i, j int) bool { return repos[i].Path < repos[j].Path })
	}

	// Symlink maintenance belongs to refresh, never to init.
	for _, removed := range result.RemovedRepositories {
		removeWorkspaceSymlink(cfg.Workspace, removed.Name, removed.Path)
	}
	repairWorkspaceSymlinks(cfg, repos, result, linkOwned)

	fmt.Fprintln(stdout, "Detecting git user configuration...")
	userDetectedCount := detectUserForRepos(repos, cfg.Profiles)

	cfg.Repositories = repos
	syncProfilesFromRepos(cfg)

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	cachePath, cacheErr := git.GetCachePath()
	if cacheErr == nil {
		statusCache := git.NewCache(git.DefaultTTL)
		statusCache.Clear()
		_ = statusCache.Save(cachePath)
		fmt.Fprintln(stdout, "Cleared git status cache")
	}

	writeRefreshSummary(stdout, result, len(repos), len(recovered), userDetectedCount)

	return nil
}

// writeRefreshSummary prints the refresh report. Counts come from the same
// reconciliation result init reports from, so the two commands cannot disagree
// about the same workspace.
func writeRefreshSummary(w io.Writer, result *reconcile.Result, totalRepos, recovered, userDetectedCount int) {
	added := result.Added + recovered

	fmt.Fprintf(w, "\nRefresh complete!\n")
	fmt.Fprintf(w, "Total repositories: %d\n", totalRepos)
	if result.Removed > 0 {
		fmt.Fprintf(w, "Removed %d %s (path no longer valid)\n",
			result.Removed, pluralize(result.Removed, "repository", "repositories"))
	}
	if added > 0 {
		fmt.Fprintf(w, "Found %d new %s\n", added, pluralize(added, "repository", "repositories"))
	}
	if result.Updated > 0 {
		fmt.Fprintf(w, "Updated %d %s\n", result.Updated, pluralize(result.Updated, "repository", "repositories"))
	}
	if userDetectedCount > 0 {
		fmt.Fprintf(w, "Repositories with user configuration: %d\n", userDetectedCount)
	}
	writeWorktreeSummary(w, result)
}

// rescanRemovedParents scans the parent directory of every repository that
// disappeared, so a checkout replaced in place is picked up in the same run.
// Returns repositories that are not already tracked.
func rescanRemovedParents(cfg *config.Config, result *reconcile.Result) []config.Repository {
	if len(result.RemovedRepositories) == 0 {
		return nil
	}

	tracked := make(map[string]bool, len(result.Repositories))
	for _, repo := range result.Repositories {
		tracked[repo.Path] = true
	}

	var recovered []config.Repository
	scannedParents := make(map[string]bool)

	for _, removed := range result.RemovedRepositories {
		parentDir := filepath.Dir(removed.Path)
		if scannedParents[parentDir] {
			continue
		}
		scannedParents[parentDir] = true

		scanResult, scanErr := discovery.Scan(parentDir, discovery.Options{MaxDepth: cfg.EffectiveScanMaxDepth()})
		if scanErr != nil {
			continue
		}
		for _, found := range scanResult.Repositories {
			if tracked[found.Path] {
				continue
			}
			tracked[found.Path] = true
			recovered = append(recovered, found)
		}
	}

	return recovered
}

// workspaceLinkOwnedRepos returns the set of tracked repository paths that gws
// may be responsible for linking into the workspace: repositories already in
// the metadata library that live outside the workspace root.
//
// Being tracked and external is a necessary condition, not a sufficient one.
// It holds for repositories the user added with 'gws add', which creates the
// workspace symlink — but it also holds for main repositories auto-registered
// from a worktree, which were never linked. repairWorkspaceSymlinks filters
// those out.
//
// Checking the filesystem for an existing link instead would be
// self-defeating: the case worth repairing is precisely the one where the link
// is missing.
func workspaceLinkOwnedRepos(cfg *config.Config) map[string]bool {
	owned := make(map[string]bool)
	for _, repo := range cfg.Repositories {
		if !isInsideWorkspace(cfg.Workspace, repo.Path) {
			owned[repo.Path] = true
		}
	}
	return owned
}

// repairWorkspaceSymlinks recreates missing workspace symlinks, and only those.
//
// Three conditions must all hold before a link is created:
//
//  1. The scan could not reach the repository. Anything reachable is already
//     findable from the workspace root, whether directly or through a symlink
//     someone else created. Linking it anyway is what used to accumulate
//     duplicate links on every refresh, because a repository discovered
//     through a nested symlink has an external real path and looked
//     link-worthy.
//
//  2. The repository was already tracked and lives outside the workspace, so
//     this is a repair rather than a discovery-driven creation.
//
//  3. The repository is not one gws registered on its own behalf because a
//     worktree of it turned up in the workspace. Those are tracked by real
//     path and deliberately get no symlink; creating one would put an entry
//     in the user's workspace that they never asked for.
//
// Condition 3 means a user-added external repository that also has a worktree
// inside the workspace will not have a deleted link auto-repaired. That is the
// conservative side of the trade: gws declines to act rather than risk
// creating a link it does not own.
func repairWorkspaceSymlinks(cfg *config.Config, repos []config.Repository, result *reconcile.Result, linkOwned map[string]bool) {
	for _, repo := range repos {
		if result.ReachablePaths[repo.Path] {
			continue
		}
		if !linkOwned[repo.Path] {
			continue
		}
		if result.WorktreeMainRepos[repo.Path] {
			continue
		}
		ensureWorkspaceSymlink(cfg, repo.Path, repo.Name)
	}
}

// isInsideWorkspace reports whether repoPath is the workspace root or below it.
func isInsideWorkspace(workspace, repoPath string) bool {
	workspacePath := filepath.Clean(workspace)
	repoClean := filepath.Clean(repoPath)
	return repoClean == workspacePath ||
		strings.HasPrefix(repoClean, workspacePath+string(filepath.Separator))
}

// removeWorkspaceSymlink removes workspace/name only when it is a symlink
// whose target matches repoPath. Leaves any non-symlink or differently-targeted
// entry untouched.
func removeWorkspaceSymlink(workspace, name, repoPath string) {
	symlinkPath := filepath.Join(workspace, name)
	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		return
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return
	}
	target, err := os.Readlink(symlinkPath)
	if err != nil || target != repoPath {
		return
	}
	_ = os.Remove(symlinkPath)
}

// ensureWorkspaceSymlink creates or corrects workspace/name → repoPath for
// repos that live outside the workspace directory. No-op for repos inside the
// workspace.
func ensureWorkspaceSymlink(cfg *config.Config, repoPath, repoName string) {
	if isInsideWorkspace(cfg.Workspace, repoPath) {
		return
	}

	symlinkPath := filepath.Join(cfg.Workspace, repoName)
	fi, err := os.Lstat(symlinkPath)
	if err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			target, readErr := os.Readlink(symlinkPath)
			if readErr == nil && target == repoPath {
				return // already correct
			}
			// Remove incorrect symlink so we can recreate it.
			_ = os.Remove(symlinkPath)
		} else {
			return // non-symlink; leave it alone
		}
	}

	_ = os.Symlink(repoPath, symlinkPath)
}
