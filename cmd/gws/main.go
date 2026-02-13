package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

var (
	// Version information (will be set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Command flags
var (
	flagList           bool
	flagInit           string
	flagAddTag         bool
	flagRemoveTag      bool
	flagRefresh        bool
	flagPrintWorkspace bool
)

// Filter flags (for --list)
var (
	filterType   string
	filterTags   []string
	filterName   string
	filterPath   string
	outputFormat string
	showStatus   bool
)

var rootCmd = &cobra.Command{
	Use:   "gws",
	Short: "Git Workspace - Discover, organize, and navigate git repositories",
	Long: `gws is a lightweight, cross-platform CLI tool for discovering, organizing,
and navigating git repositories on your local system. It provides an intelligent
repository index and navigation layer with powerful search and filtering capabilities.

Commands (flags):
  gws --list                         # List all repositories
  gws -l --type github               # List only GitHub repositories
  gws -l --tag personal --status     # List repos tagged "personal" with git status
  gws --init ~/projects              # Initialize workspace
  gws --add-tag my-project personal  # Add tag to matching repos
  gws --remove-tag api work          # Remove tag from matching repos
  gws --refresh                      # Refresh repository metadata

For navigation support, add this to your shell config:

  # Bash/Zsh
  function cdgws() { cd "$(gws --print-workspace)"; }
  alias gcd=cdgws`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Count active command flags for mutual exclusivity
		activeCount := 0
		if flagList {
			activeCount++
		}
		if flagInit != "" {
			activeCount++
		}
		if flagAddTag {
			activeCount++
		}
		if flagRemoveTag {
			activeCount++
		}
		if flagRefresh {
			activeCount++
		}
		if flagPrintWorkspace {
			activeCount++
		}

		if activeCount > 1 {
			return fmt.Errorf("only one command flag can be used at a time (--list, --init, --add-tag, --remove-tag, --refresh, --print-workspace)")
		}

		// Dispatch to command handlers
		if flagInit != "" {
			return runInit(cmd, args)
		}

		// All other commands require workspace to be initialized
		if activeCount > 0 || len(args) == 0 {
			exists, err := config.Exists()
			if err != nil {
				return fmt.Errorf("failed to check workspace status: %w", err)
			}

			if !exists {
				fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "To get started, initialize a workspace:")
				fmt.Fprintln(os.Stderr, "  gws --init <directory>")
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "Example:")
				fmt.Fprintln(os.Stderr, "  gws --init ~/projects")
				return fmt.Errorf("workspace not initialized")
			}
		}

		if flagPrintWorkspace {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load workspace configuration: %w", err)
			}
			fmt.Println(cfg.Workspace)
			return nil
		}

		if flagList {
			return runList(cmd, args)
		}

		if flagAddTag {
			return runTag(cmd, args)
		}

		if flagRemoveTag {
			return runUntag(cmd, args)
		}

		if flagRefresh {
			return runRefresh(cmd, args)
		}

		// Default behavior: display workspace information
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load workspace configuration: %w", err)
		}

		fmt.Printf("Workspace: %s\n", cfg.Workspace)
		fmt.Printf("Repositories: %d\n", len(cfg.Repositories))
		fmt.Println("")
		fmt.Println("Use 'gws --help' to see available commands")

		return nil
	},
}

func init() {
	// Command flags
	rootCmd.Flags().BoolVarP(&flagList, "list", "l", false, "List all tracked repositories")
	rootCmd.Flags().StringVarP(&flagInit, "init", "i", "", "Initialize workspace by scanning directory")
	rootCmd.Flags().BoolVarP(&flagAddTag, "add-tag", "a", false, "Add a tag to repositories (args: <repo> <tag>)")
	rootCmd.Flags().BoolVarP(&flagRemoveTag, "remove-tag", "u", false, "Remove a tag from repositories (args: <repo> <tag>)")
	rootCmd.Flags().BoolVarP(&flagRefresh, "refresh", "r", false, "Refresh repository metadata and git status cache")
	rootCmd.Flags().BoolVarP(&flagPrintWorkspace, "print-workspace", "w", false, "Print workspace path (for shell integration)")

	// Filter flags (apply when --list is active)
	rootCmd.Flags().StringVar(&filterType, "type", "", "Filter by repository type (github, gitlab, ado, bitbucket)")
	rootCmd.Flags().StringSliceVar(&filterTags, "tag", []string{}, "Filter by custom tag(s) - can be specified multiple times for AND logic")
	rootCmd.Flags().StringVar(&filterName, "name", "", "Filter by repository name (partial match)")
	rootCmd.Flags().StringVar(&filterPath, "path", "", "Filter by repository path (partial match)")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")
	rootCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show git status (branch, clean/dirty, ahead/behind)")
}

func main() {
	rootCmd.Version = fmt.Sprintf("%s\n  commit: %s\n  built:  %s", version, commit, date)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
