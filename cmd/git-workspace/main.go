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

// Command flags
var (
	flagList           bool
	flagInit           bool
	flagAdd            string
	flagRecursive      bool
	flagAddTag         bool
	flagRemoveTag      bool
	flagRefresh        bool
	flagPrintWorkspace bool
	flagGo             string
	flagQuiet          bool
	flagUser           bool
	flagUpdate         bool
	flagDelete         bool
	flagAll            bool
	flagVerbose        bool
	flagInlineName     string
	flagInlineEmail    string
)

// Filter flags (for --list)
var (
	filterType   string
	filterTags   []string
	filterName   string
	filterPath   string
	outputFormat string
	showStatus   bool
	showUser     bool
)

var rootCmd = &cobra.Command{
	Use:   "git-workspace",
	Short: "Git Workspace - Discover, organize, and navigate git repositories",
	Long: `git-workspace is a lightweight, cross-platform CLI tool for discovering, organizing,
and navigating git repositories on your local system. It provides an intelligent
repository index and navigation layer with powerful search and filtering capabilities.

Commands (flags):
  git-workspace --list                         # List all repositories
  git-workspace -l --type github               # List only GitHub repositories
  git-workspace -l --tag personal --status     # List repos tagged "personal" with git status
  git-workspace --init                         # Initialize workspace in current directory
  git-workspace --add                          # Add current directory to workspace
  git-workspace --add ~/elsewhere/my-repo     # Add a specific repo to workspace
  git-workspace --add-tag my-project personal  # Add tag to matching repos
  git-workspace --remove-tag api work          # Remove tag from matching repos
  git-workspace --refresh                      # Refresh repository metadata

Navigation:
  git-workspace my-repo                        # Navigate to repository by name
  git-workspace --go my-repo                   # Navigate using flag (same as above)
  git-workspace -g my-repo -q                  # Quiet mode: print only the path
  git-workspace "api-*"                        # Wildcard match with interactive selection

Shell integration (add to ~/.bashrc or ~/.zshrc):

  # Recommended — sets up 'gws' and tab completion, always up to date:
  export PATH="$HOME/.local/bin:$PATH"
  eval "$(git-workspace shell-init zsh)"   # or: shell-init bash

  # Navigate to workspace root
  function cdgws() { cd "$(git-workspace --print-workspace)"; }
  alias gcd=cdgws`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Count active command flags for mutual exclusivity
		activeCount := 0
		if flagList {
			activeCount++
		}
		if flagInit {
			activeCount++
		}
		if flagAdd != "" {
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
		if flagUser {
			activeCount++
		}

		if activeCount > 1 {
			return fmt.Errorf("only one command flag can be used at a time (--list, --init, --add, --add-tag, --remove-tag, --refresh, --print-workspace, --go, --user)")
		}

		// Validate --user flag dependencies
		if flagUpdate && !flagUser {
			return fmt.Errorf("--update/-u requires --user to be set")
		}
		if flagDelete && !flagUser {
			return fmt.Errorf("--delete/-D requires --user to be set")
		}
		if flagUser && !flagUpdate && !flagDelete {
			return fmt.Errorf("--user requires either --update/-u or --delete/-D")
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

		// Validate filter flags: --tag is allowed with --list or --user; other filters require --list
		if !flagList && hasNonTagFilterFlags(cmd) {
			return fmt.Errorf("filter flags (--type, --name, --path, --output, --status) require --list/-l to be set")
		}
		if !flagList && !flagUser && cmd.Flags().Changed("tag") {
			return fmt.Errorf("--tag requires --list/-l or --user to be set")
		}

		// Validate --recursive requires --add
		if flagRecursive && flagAdd == "" {
			return fmt.Errorf("--recursive/-v requires --add/-a to be set")
		}

		// Validate --quiet applies to navigation and user operations
		if flagQuiet && flagGo == "" && len(args) == 0 && !flagUser {
			return fmt.Errorf("--quiet/-q can only be used with navigation or --user operations")
		}

		// Validate --go and positional args are not both provided
		if flagGo != "" && len(args) > 0 {
			return fmt.Errorf("cannot use both --go flag and positional argument for navigation")
		}

		// Dispatch to command handlers
		if flagInit {
			return runInit(cmd, args)
		}
		if flagAdd != "" {
			return runAdd(cmd, args)
		}

		// All other commands require workspace to be initialized
		exists, err := config.Exists()
		if err != nil {
			return fmt.Errorf("failed to check workspace status: %w", err)
		}

		if !exists {
			fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To get started, navigate to your projects directory and run:")
			fmt.Fprintln(os.Stderr, "  gws --init")
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

		if flagUser {
			if flagUpdate {
				return runUserUpdate(cmd, args)
			}
			return runUserDelete(cmd, args)
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
	rootCmd.Flags().BoolVarP(&flagInit, "init", "i", false, "Initialize a gws workspace in the current directory")
	rootCmd.Flags().StringVarP(&flagAdd, "add", "a", "", "Add a git repository to the workspace (defaults to current directory)")
	rootCmd.Flags().Lookup("add").NoOptDefVal = "."
	rootCmd.Flags().BoolVarP(&flagRecursive, "recursive", "v", false, "Recursively add all git repositories found in the current directory (use with --add)")
	rootCmd.Flags().BoolVarP(&flagAddTag, "add-tag", "d", false, "Add a tag to repositories (args: <repo> <tag>)")
	rootCmd.Flags().BoolVarP(&flagRemoveTag, "remove-tag", "x", false, "Remove a tag from repositories (args: <repo> <tag>)")
	rootCmd.Flags().BoolVarP(&flagRefresh, "refresh", "r", false, "Refresh repository metadata and git status cache")
	rootCmd.Flags().BoolVarP(&flagPrintWorkspace, "print-workspace", "w", false, "Print workspace path (for shell integration)")
	rootCmd.Flags().StringVarP(&flagGo, "go", "g", "", "Navigate to a repository by name (prints path)")
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress verbose output, print only the path (navigation only)")

	// User operation flags
	rootCmd.Flags().BoolVar(&flagUser, "user", false, "Enable user operations (requires --update or --delete)")
	rootCmd.Flags().BoolVarP(&flagUpdate, "update", "u", false, "Update local git user config for repositories (requires --user)")
	rootCmd.Flags().BoolVarP(&flagDelete, "delete", "D", false, "Delete local git user config from repositories (requires --user)")
	rootCmd.Flags().BoolVar(&flagAll, "all", false, "Also remove signing config when deleting (requires --delete)")
	rootCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show detailed output for user operations")
	rootCmd.Flags().StringVar(&flagInlineName, "git-name", "", "Inline git user.name for --user --update")
	rootCmd.Flags().StringVar(&flagInlineEmail, "git-email", "", "Inline git user.email for --user --update")

	// Filter flags (apply when --list is active)
	rootCmd.Flags().StringVarP(&filterType, "type", "y", "", "Filter by repository type (github, gitlab, ado, bitbucket)")
	rootCmd.Flags().StringSliceVarP(&filterTags, "tag", "t", []string{}, "Filter by custom tag(s) - can be specified multiple times for AND logic")
	rootCmd.Flags().StringVarP(&filterName, "name", "n", "", "Filter by repository name (partial match)")
	rootCmd.Flags().StringVarP(&filterPath, "path", "p", "", "Filter by repository path (partial match)")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")
	rootCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show git status (branch, clean/dirty, ahead/behind)")
	rootCmd.Flags().BoolVar(&showUser, "show-user", false, "Show git user info (USER, EMAIL, SIGN columns)")

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
	cobra.AddTemplateFunc("commandFlagUsages", flagGroup([]string{"list", "init", "add", "add-tag", "remove-tag", "refresh", "print-workspace"}))
	cobra.AddTemplateFunc("userFlagUsages", flagGroup([]string{"user", "update", "delete", "all", "verbose", "git-name", "git-email"}))
	cobra.AddTemplateFunc("filterFlagUsages", flagGroup([]string{"type", "tag", "name", "path", "output", "status", "show-user"}))
	cobra.AddTemplateFunc("navigationFlagUsages", flagGroup([]string{"go", "quiet"}))

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

Commands:
{{commandFlagUsages . | trimRightSpace}}
  gws completion [bash|zsh|fish|powershell]    Generate shell completion script
  gws shell-init [zsh|bash]                    Output shell integration code to eval

User Operations (require --user):
{{userFlagUsages . | trimRightSpace}}

List Filters (require --list / -l):
{{filterFlagUsages . | trimRightSpace}}

Navigation:
{{navigationFlagUsages . | trimRightSpace}}

Other:
  -h, --help      help for {{.Name}}
      --version   version for {{.Name}}
`)
}

// hasFilterFlags checks if any filter flag has been explicitly set by the user
func hasFilterFlags(cmd *cobra.Command) bool {
	filterFlagNames := []string{"type", "tag", "name", "path", "status", "show-user"}
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

// hasNonTagFilterFlags checks if any filter flag other than --tag has been set
func hasNonTagFilterFlags(cmd *cobra.Command) bool {
	nonTagFilterNames := []string{"type", "name", "path", "status", "show-user"}
	for _, name := range nonTagFilterNames {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
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
