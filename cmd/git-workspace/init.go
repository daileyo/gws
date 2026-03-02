package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

// initCmd is the Cobra subcommand for initializing a workspace.
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a gws workspace",
	Long: `Initialize a gws workspace in the specified directory (defaults to current directory).
Scans for existing git repositories and creates the workspace configuration.

Examples:
  gws init                # Initialize in current directory
  gws init ~/projects     # Initialize in a specific directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := ""
		if len(args) > 0 {
			dir = args[0]
		}
		return runInit(dir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// runInit handles the init logic.
// dir is the directory to initialize (empty string means current directory).
func runInit(dir string) error {
	// Guard: if a workspace is already initialized, notify and exit cleanly
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check workspace status: %w", err)
	}
	if exists {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load workspace configuration: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Workspace already initialized at: %s\n", cfg.Workspace)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "To add more repositories:  gws add")
		fmt.Fprintln(os.Stderr, "To re-scan the workspace:  gws refresh")
		return nil
	}

	// Resolve workspace path
	var absPath string
	if dir == "" || dir == "." {
		absPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to resolve current directory: %w", err)
		}
	} else {
		absPath, err = filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve directory path: %w", err)
		}
	}

	// Scan for repositories
	result, err := discovery.Scan(absPath)
	if err != nil {
		return fmt.Errorf("failed to scan workspace: %w", err)
	}

	// Display scan errors if any
	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: %d errors occurred during scanning:\n", len(result.Errors))
		for i, scanErr := range result.Errors {
			if i < 5 {
				fmt.Fprintf(os.Stderr, "  - %v\n", scanErr)
			}
		}
		if len(result.Errors) > 5 {
			fmt.Fprintf(os.Stderr, "  ... and %d more errors\n", len(result.Errors)-5)
		}
	}

	// Create and save configuration
	cfg := config.New(absPath)
	cfg.Repositories = result.Repositories

	// Detect git user configuration for discovered repos
	userDetectedCount := detectUserForRepos(cfg.Repositories)

	// Sync profiles from detected repo users
	syncProfilesFromRepos(cfg)

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Initialized workspace at: %s\n", absPath)
	fmt.Printf("Found %d %s.\n", result.Count, pluralize(result.Count, "repository", "repositories"))
	if userDetectedCount > 0 {
		fmt.Printf("Repositories with user configuration: %d\n", userDetectedCount)
	}

	return nil
}

func pluralize(count int, singular, plural string) string { //nolint:unparam // singular is always "repository" but kept for readability
	if count == 1 {
		return singular
	}
	return plural
}
