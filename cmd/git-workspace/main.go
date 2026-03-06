package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

var (
	// Version information (will be set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Root-level flags that are not deprecated.
var (
	flagTagAlias bool
	flagQuiet    bool
)

// filterTags is used by the deprecated root --tag flag (deprecated.go)
// for backward compat with --list and deprecated user operations (--user -u --tag).
var filterTags []string

var rootCmd = &cobra.Command{
	Use:              "git-workspace",
	Short:            "Git Workspace - Discover, organize, and navigate git repositories",
	Args:             cobra.ArbitraryArgs,
	TraverseChildren: true,
	Long: `git-workspace is a lightweight, cross-platform CLI tool for discovering, organizing,
and navigating git repositories on your local system. It provides an intelligent
repository index and navigation layer with powerful search and filtering capabilities.

Running gws with no arguments lists all repositories in a compact multi-column view.

Commands:
  gws                                    # List repos (default, same as gws list)
  gws list                               # List all repositories
  gws list -v                            # Show type, visibility, tags, path
  gws list -vv                           # Show all columns including status, user, remote
  gws list -y=github                     # Filter by GitHub, show type column
  gws init                               # Initialize workspace in current directory
  gws add                                # Add current directory to workspace
  gws refresh                            # Refresh repository metadata

Navigation:
  gws my-repo                            # Navigate to repository by name
  gws "api-*"                            # Wildcard match with interactive selection

Shell integration (add to ~/.bashrc or ~/.zshrc):

  # Recommended — sets up 'gws' and tab completion, always up to date:
  export PATH="$HOME/.local/bin:$PATH"
  eval "$(git-workspace shell-init zsh)"   # or: shell-init bash

  # Navigate to workspace root
  function cdgws() { cd "$(gws print-workspace)"; }
  alias gcd=cdgws`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle deprecated flags first (--list, --init, --add, --refresh, --user, etc.)
		handled, err := handleDeprecatedFlags(cmd, args)
		if handled {
			return err
		}

		// Validate --quiet applies to navigation only
		if flagQuiet && len(args) == 0 {
			return fmt.Errorf("--quiet/-q can only be used with navigation")
		}

		// All remaining commands require workspace to be initialized
		exists, err := config.Exists()
		if err != nil {
			return fmt.Errorf("failed to check workspace status: %w", err)
		}

		if !exists {
			fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To get started, navigate to your projects directory and run:")
			fmt.Fprintln(os.Stderr, "  gws init")
			return fmt.Errorf("workspace not initialized")
		}

		// Tag alias: -t delegates to the tag subcommand with remaining args
		if flagTagAlias {
			tagCmd.SetArgs(args)
			return tagCmd.Execute()
		}

		// Navigation via positional argument
		if len(args) > 0 {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load workspace configuration: %w", err)
			}
			return runNavigate(args[0], flagQuiet, cfg.Repositories, os.Stderr, os.Stdout, os.Stdin)
		}

		// Default behavior: run list with default options (multi-column names)
		return runList(ListOptions{OutputFormat: "table"})
	},
}

func init() {
	// Register print-workspace subcommand
	rootCmd.AddCommand(printWorkspaceCmd)

	// Register deprecated flags (hidden, emit warnings when used)
	registerDeprecatedFlags(rootCmd)

	// Tag alias: -t on root delegates to the tag subcommand
	rootCmd.Flags().BoolVarP(&flagTagAlias, "tag-cmd", "t", false, "Alias for 'tag' subcommand")
	_ = rootCmd.Flags().MarkHidden("tag-cmd")

	// Navigation flags
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress verbose output, print only the path (navigation only)")

	// Register Cobra's built-in completion subcommand (bash, zsh, fish, powershell)
	rootCmd.InitDefaultCompletionCmd()

	// Repo name completion for positional args (navigation)
	rootCmd.ValidArgsFunction = func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		cfg, err := config.Load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		seen := make(map[string]bool)
		var names []string
		for _, repo := range cfg.Repositories {
			if strings.HasPrefix(strings.ToLower(repo.Name), strings.ToLower(toComplete)) && !seen[repo.Name] {
				seen[repo.Name] = true
				names = append(names, repo.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}

	// Directory completion for --add flag value
	_ = rootCmd.RegisterFlagCompletionFunc("add", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	})

	rootCmd.SetUsageTemplate(`Usage:
  {{.UseLine}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Navigation:
  gws <repo-name>                                   # Navigate to repository by name
  gws "api-*"                                       # Wildcard match with interactive selection

Flags:
  -q, --quiet     Suppress verbose output, print only the path (navigation only)
  -h, --help      help for {{.Name}}
      --version   version for {{.Name}}
`)
}

// printWorkspaceCmd is the Cobra subcommand for printing the workspace path.
var printWorkspaceCmd = &cobra.Command{
	Use:   "print-workspace",
	Short: "Print workspace path (for shell integration)",
	Long: `Print the workspace root path. Useful for shell integration and scripting.

Examples:
  gws print-workspace
  cd "$(gws print-workspace)"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrintWorkspace()
	},
}

// runPrintWorkspace prints the workspace root path.
func runPrintWorkspace() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load workspace configuration: %w", err)
	}
	fmt.Println(cfg.Workspace)
	return nil
}

func main() {
	rootCmd.Version = fmt.Sprintf("%s\n  commit: %s\n  built:  %s", version, commit, date)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
