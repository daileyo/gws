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
	flagGo             string
	flagQuiet          bool
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

Navigation:
  gws my-repo                        # Navigate to repository by name
  gws --go my-repo                   # Navigate using flag (same as above)
  gws -g my-repo -q                  # Quiet mode: print only the path
  gws "api-*"                        # Wildcard match with interactive selection

Shell integration (add to ~/.bashrc or ~/.zshrc):

  # Navigate to workspace root
  function cdgws() { cd "$(gws --print-workspace)"; }
  alias gcd=cdgws

  # Navigate to a repository by name
  function cdg() { cd "$(gws -g "$1" -q)"; }`,
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
		if flagGo != "" {
			activeCount++
		}

		if activeCount > 1 {
			return fmt.Errorf("only one command flag can be used at a time (--list, --init, --add-tag, --remove-tag, --refresh, --print-workspace, --go)")
		}

		// Validate filter flags require --list
		if !flagList && hasFilterFlags(cmd) {
			return fmt.Errorf("filter flags (--type, --tag, --name, --path, --output, --status) require --list/-l to be set")
		}

		// Validate --quiet only applies to navigation
		if flagQuiet && flagGo == "" && len(args) == 0 {
			return fmt.Errorf("--quiet/-q can only be used with navigation (--go or positional argument)")
		}

		// Validate --go and positional args are not both provided
		if flagGo != "" && len(args) > 0 {
			return fmt.Errorf("cannot use both --go flag and positional argument for navigation")
		}

		// Dispatch to command handlers
		if flagInit != "" {
			return runInit(cmd, args)
		}

		// All other commands require workspace to be initialized
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

		// Navigation via --go flag
		if flagGo != "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load workspace configuration: %w", err)
			}
			return runNavigate(flagGo, flagQuiet, cfg.Repositories, os.Stderr, os.Stdout, os.Stdin)
		}

		// Navigation via positional argument
		if len(args) > 0 {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load workspace configuration: %w", err)
			}
			return runNavigate(args[0], flagQuiet, cfg.Repositories, os.Stderr, os.Stdout, os.Stdin)
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
	rootCmd.Flags().StringVarP(&flagGo, "go", "g", "", "Navigate to a repository by name (prints path)")
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress verbose output, print only the path (navigation only)")

	// Filter flags (apply when --list is active)
	rootCmd.Flags().StringVarP(&filterType, "type", "y", "", "Filter by repository type (github, gitlab, ado, bitbucket)")
	rootCmd.Flags().StringSliceVarP(&filterTags, "tag", "t", []string{}, "Filter by custom tag(s) - can be specified multiple times for AND logic")
	rootCmd.Flags().StringVarP(&filterName, "name", "n", "", "Filter by repository name (partial match)")
	rootCmd.Flags().StringVarP(&filterPath, "path", "p", "", "Filter by repository path (partial match)")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")
	rootCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show git status (branch, clean/dirty, ahead/behind)")
}

// hasFilterFlags checks if any filter flag has been explicitly set by the user
func hasFilterFlags(cmd *cobra.Command) bool {
	filterFlagNames := []string{"type", "tag", "name", "path", "status"}
	for _, name := range filterFlagNames {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	// Check output only if it was explicitly changed from default
	if cmd.Flags().Changed("output") {
		return true
	}
	return false
}

func main() {
	rootCmd.Version = fmt.Sprintf("%s\n  commit: %s\n  built:  %s", version, commit, date)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
