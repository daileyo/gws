package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
	"github.com/daileyo/gws/internal/reconcile"
)

// initCmd is the Cobra subcommand for initializing a workspace.
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a gws workspace",
	Long: `Initialize a gws workspace in the specified directory (defaults to current directory).
Recursively scans for git repositories, discovers their worktrees, and creates
the workspace configuration.

Examples:
  gws init                # Initialize in current directory
  gws init ~/projects     # Initialize in a specific directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := ""
		if len(args) > 0 {
			dir = args[0]
		}
		return runInit(dir, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// runInit creates the workspace metadata library and performs the first full
// reconciliation against it.
//
// dir is the directory to initialize (empty string means current directory).
// init and refresh share one reconciliation engine, so a freshly initialized
// workspace has the same repository and worktree model a refreshed one does.
func runInit(dir string, stdout io.Writer) error {
	// Guard: if a workspace is already initialized, notify and exit cleanly
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check workspace status: %w", err)
	}
	if exists {
		cfg, loadErr := config.Load()
		if loadErr != nil {
			return fmt.Errorf("failed to load workspace configuration: %w", loadErr)
		}
		fmt.Fprintf(os.Stderr, "Workspace already initialized at: %s\n", cfg.Workspace)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "To re-scan the workspace and pick up changes:  gws refresh")
		fmt.Fprintln(os.Stderr, "To add a single repository:                    gws add")
		return nil
	}

	absPath, err := resolveWorkspacePath(dir)
	if err != nil {
		return err
	}

	// Create the metadata library, then reconcile it. Passing a nil existing
	// configuration is what marks this as the init case.
	cfg := config.New(absPath)

	progress := git.NewProgressWithLabel(0, "Scanning workspace...")
	result, err := reconcile.ReconcileWorkspace(absPath, nil, reconcile.Options{
		MaxDepth: cfg.EffectiveScanMaxDepth(),
		Progress: progress,
	})
	if err != nil {
		return err
	}

	reconcile.ReportScanErrors(os.Stderr, result.Errors)

	cfg.Repositories = result.Repositories

	// Detect git user configuration for discovered repos
	userDetectedCount := detectUserForRepos(cfg.Repositories)

	// Sync profiles from detected repo users
	syncProfilesFromRepos(cfg)

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Fprintf(stdout, "Initialized workspace at: %s\n", absPath)
	fmt.Fprintf(stdout, "Found %d %s.\n",
		result.TotalRepositories,
		pluralize(result.TotalRepositories, "repository", "repositories"))
	if userDetectedCount > 0 {
		fmt.Fprintf(stdout, "Repositories with user configuration: %d\n", userDetectedCount)
	}
	writeWorktreeSummary(stdout, result)

	return nil
}

// resolveWorkspacePath turns a user-supplied directory into an absolute path,
// defaulting to the current directory.
func resolveWorkspacePath(dir string) (string, error) {
	if dir == "" || dir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to resolve current directory: %w", err)
		}
		return cwd, nil
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve directory path: %w", err)
	}
	return absPath, nil
}

// writeWorktreeSummary prints the shared worktree lines used by both init and
// refresh, so the two commands never describe the same workspace differently.
// Nothing is printed when no worktrees were found.
func writeWorktreeSummary(w io.Writer, result *reconcile.Result) {
	if result.ReposWithWorktrees > 0 {
		fmt.Fprintf(w, "Repositories with worktrees: %d\n", result.ReposWithWorktrees)
	}
	if result.TotalWorktrees > 0 {
		fmt.Fprintf(w, "Worktrees: %d (%d aligned, %d unaligned)\n",
			result.TotalWorktrees, result.AlignedWorktrees, result.UnalignedWorktrees)
	}
}

func pluralize(count int, singular, plural string) string { //nolint:unparam // singular is always "repository" but kept for readability
	if count == 1 {
		return singular
	}
	return plural
}
