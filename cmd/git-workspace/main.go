package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/daileyo/gws/internal/config"
)

var (
	// Version information (will be set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Command flags that remain on root (tag, user — migrated in specs 10/11).
var (
	flagAddTag      bool
	flagRemoveTag   bool
	flagQuiet       bool
	flagUser        bool
	flagUpdate      bool
	flagDelete      bool
	flagAll         bool
	flagVerbose     bool
	flagInlineName  string
	flagInlineEmail string
	flagListUsers   bool
)

// filterTags is shared between the deprecated --tag flag on root (registered in deprecated.go)
// and user operations (userupdate.go, userdelete.go). It will move to the user subcommand in spec 11.
var filterTags []string

var rootCmd = &cobra.Command{
	Use:              "git-workspace",
	Short:            "Git Workspace - Discover, organize, and navigate git repositories",
	Args:             cobra.ArbitraryArgs,
	TraverseChildren: true,
	Long: `git-workspace is a lightweight, cross-platform CLI tool for discovering, organizing,
and navigating git repositories on your local system. It provides an intelligent
repository index and navigation layer with powerful search and filtering capabilities.

Commands:
  gws list                              # List all repositories
  gws list --type github                # List only GitHub repositories
  gws list -t personal -s               # List repos tagged "personal" with git status
  gws init                              # Initialize workspace in current directory
  gws add                               # Add current directory to workspace
  gws add ~/elsewhere/my-repo           # Add a specific repo to workspace
  gws refresh                           # Refresh repository metadata
  gws print-workspace                   # Print workspace root path

Navigation:
  gws my-repo                           # Navigate to repository by name
  gws "api-*"                           # Wildcard match with interactive selection

Shell integration (add to ~/.bashrc or ~/.zshrc):

  # Recommended — sets up 'gws' and tab completion, always up to date:
  export PATH="$HOME/.local/bin:$PATH"
  eval "$(git-workspace shell-init zsh)"   # or: shell-init bash

  # Navigate to workspace root
  function cdgws() { cd "$(gws print-workspace)"; }
  alias gcd=cdgws`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle deprecated flags first (--list, --init, --add, --refresh, etc.)
		handled, err := handleDeprecatedFlags(cmd, args)
		if handled {
			return err
		}

		// Validate --user flag dependencies
		if flagUpdate && !flagUser {
			return fmt.Errorf("--update/-u requires --user to be set")
		}
		if flagDelete && !flagUser {
			return fmt.Errorf("--delete/-D requires --user to be set")
		}
		if flagUpdate && flagDelete {
			return fmt.Errorf("--update and --delete are mutually exclusive")
		}
		if flagAll && !flagDelete {
			return fmt.Errorf("--all requires --delete/-D to be set")
		}
		if (flagInlineName != "" || flagInlineEmail != "") && !flagUpdate {
			return fmt.Errorf("--git-name and --git-email require --user --update to be set")
		}
		if flagVerbose && !flagUser {
			return fmt.Errorf("--verbose requires --user to be set")
		}

		// Validate --quiet applies to navigation and user operations
		if flagQuiet && len(args) == 0 && !flagUser {
			return fmt.Errorf("--quiet/-q can only be used with navigation or --user operations")
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

		// Tag operations (remain on root until spec 10)
		if flagAddTag {
			if len(args) != 2 {
				return fmt.Errorf("--add-tag requires exactly 2 arguments: <repository> <tag>")
			}
			return runAddTag(args[0], args[1])
		}
		if flagRemoveTag {
			if len(args) != 2 {
				return fmt.Errorf("--remove-tag requires exactly 2 arguments: <repository> <tag>")
			}
			return runRemoveTag(args[0], args[1])
		}

		// User operations (remain on root until spec 11)
		if flagListUsers {
			return runListUsers(cmd, args)
		}
		if flagUser {
			if flagUpdate {
				return runUserUpdate(cmd, args)
			}
			if flagDelete {
				return runUserDelete(cmd, args)
			}
			// --user alone: list profiles
			return runListUsers(cmd, args)
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
	// Register print-workspace subcommand
	rootCmd.AddCommand(printWorkspaceCmd)

	// Register deprecated flags (hidden, emit warnings when used)
	registerDeprecatedFlags(rootCmd)

	// Tag command flags (remain on root — migrated to tag subcommand in spec 10)
	rootCmd.Flags().BoolVarP(&flagAddTag, "add-tag", "d", false, "Add a tag to repositories (args: <repo> <tag>)")
	rootCmd.Flags().BoolVarP(&flagRemoveTag, "remove-tag", "x", false, "Remove a tag from repositories (args: <repo> <tag>)")

	// Navigation flags
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress verbose output, print only the path (navigation only)")

	// User operation flags (remain on root — migrated to user subcommand in spec 11)
	rootCmd.Flags().BoolVar(&flagUser, "user", false, "List profiles, or use with --update/-u or --delete/-D")
	rootCmd.Flags().BoolVarP(&flagUpdate, "update", "u", false, "Update local git user config for repositories (requires --user)")
	rootCmd.Flags().BoolVarP(&flagDelete, "delete", "D", false, "Delete local git user config from repositories (requires --user)")
	rootCmd.Flags().BoolVar(&flagAll, "all", false, "Also remove signing config when deleting (requires --delete)")
	rootCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show detailed output for user operations")
	rootCmd.Flags().StringVar(&flagInlineName, "git-name", "", "Inline git user.name for --user --update")
	rootCmd.Flags().StringVar(&flagInlineEmail, "git-email", "", "Inline git user.email for --user --update")
	rootCmd.Flags().BoolVar(&flagListUsers, "list-users", false, "List all available user profiles")

	// Register template helpers for grouped flag display
	flagGroup := func(names []string) func(*cobra.Command) string {
		return func(cmd *cobra.Command) string {
			fs := pflag.NewFlagSet("", pflag.ContinueOnError)
			for _, name := range names {
				if f := cmd.Flags().Lookup(name); f != nil {
					fs.AddFlag(f)
				}
			}
			return fs.FlagUsages()
		}
	}
	cobra.AddTemplateFunc("userFlagUsages", flagGroup([]string{"user", "update", "delete", "all", "verbose", "git-name", "git-email", "list-users"}))
	cobra.AddTemplateFunc("navigationFlagUsages", flagGroup([]string{"quiet"}))

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

Tag Operations (migrating to subcommand in future release):
  gws --add-tag <repo> <tag>                        # Add tag to matching repos
  gws --remove-tag <repo> <tag>                     # Remove tag from matching repos

User Operations (migrating to subcommand in future release):
  gws --user                                        # List available user profiles
  gws --list-users                                  # Same as above
  gws --user -u <repo> <profile>                    # Set local user from a stored profile
  gws --user -D <repo>                              # Remove local user config

  Flags:
{{userFlagUsages . | trimRightSpace}}

Navigation:
{{navigationFlagUsages . | trimRightSpace}}

Other:
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
