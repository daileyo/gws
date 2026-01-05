package main

import (
	"fmt"
	"os"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
	"github.com/spf13/cobra"
)

var (
	// Version information (will be set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gws [tag]",
	Short: "Git Workspace - Discover, organize, and navigate git repositories",
	Long: `gws is a lightweight, cross-platform CLI tool for discovering, organizing,
and navigating git repositories on your local system. It provides an intelligent
repository index and navigation layer with powerful search and filtering capabilities.

Usage:
  gws              # Show workspace information
  gws personal     # List repositories tagged as "personal" (shorthand for: gws list --tag personal)
  gws work         # List repositories tagged as "work"

For navigation support, add this to your shell config:

  # Bash/Zsh
  function cdgws() { cd "$(gws --print-workspace)"; }
  alias gcd=cdgws`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if workspace is initialized
		exists, err := config.Exists()
		if err != nil {
			return fmt.Errorf("failed to check workspace status: %w", err)
		}

		if !exists {
			fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To get started, initialize a workspace:")
			fmt.Fprintln(os.Stderr, "  gws init <directory>")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Example:")
			fmt.Fprintln(os.Stderr, "  gws init ~/projects")
			return fmt.Errorf("workspace not initialized")
		}

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load workspace configuration: %w", err)
		}

		// Check for --print-workspace flag
		printWorkspace, _ := cmd.Flags().GetBool("print-workspace")
		if printWorkspace {
			fmt.Println(cfg.Workspace)
			return nil
		}

		// If args provided, treat as tag filter shorthand
		if len(args) > 0 {
			tagFilter := args[0]

			// Filter repositories by tag
			criteria := filter.Criteria{
				Tags: []string{tagFilter},
			}
			filtered := filter.Apply(cfg.Repositories, criteria)

			if len(filtered) == 0 {
				fmt.Printf("No repositories found with tag '%s'\n", tagFilter)
				return nil
			}

			// Display results
			fmt.Printf("Repositories tagged '%s':\n\n", tagFilter)

			// Load git status cache
			statusCache := git.NewCache(git.DefaultTTL)
			cachePath, err := git.GetCachePath()
			if err == nil {
				_ = statusCache.Load(cachePath)
			}

			displayTable(filtered, statusCache)
			return nil
		}

		// Default behavior: display workspace information
		fmt.Printf("Workspace: %s\n", cfg.Workspace)
		fmt.Printf("Repositories: %d\n", len(cfg.Repositories))
		fmt.Println("")
		fmt.Println("Use 'gws --help' to see available commands")
		fmt.Println("Use 'gws <tag>' to filter repositories by tag (e.g., 'gws personal')")

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, commit hash, and build date of the gws CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gws version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().Bool("print-workspace", false, "Print workspace path (for shell integration)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
