package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

// Deprecated flag variables — these back the hidden root flags that provide
// backward compatibility with the old `gws --list`, `gws --init`, etc. forms.
var (
	depList           bool
	depInit           bool
	depAdd            string
	depRecursive      bool
	depRefresh        bool
	depPrintWorkspace bool
	depGo             string

	// Deprecated tag flags (migrated to tag subcommand in spec 10).
	depAddTag    bool
	depRemoveTag bool

	// Deprecated filter flags (for compound usage like `gws --list --type github`).
	// filterTags is shared with user operations (which still use --tag on root).
	depFilterType   string
	depFilterName   string
	depFilterPath   string
	depOutputFormat string
	depShowStatus   bool
	depShowUser     bool
)

// registerDeprecatedFlags registers hidden flags on rootCmd that provide
// backward compatibility with the old flag-based interface.
func registerDeprecatedFlags(root *cobra.Command) {
	root.Flags().BoolVarP(&depList, "list", "l", false, "Shorthand for 'list' subcommand")
	root.Flags().BoolVarP(&depInit, "init", "i", false, "Shorthand for 'init' subcommand")
	root.Flags().StringVarP(&depAdd, "add", "a", "", "Shorthand for 'add' subcommand")
	root.Flags().Lookup("add").NoOptDefVal = "."
	root.Flags().BoolVarP(&depRecursive, "recursive", "v", false, "Recursively add all git repositories (use with --add/-a)")
	root.Flags().BoolVarP(&depRefresh, "refresh", "r", false, "Shorthand for 'refresh' subcommand")
	root.Flags().BoolVarP(&depPrintWorkspace, "print-workspace", "w", false, "Shorthand for 'print-workspace' subcommand")
	root.Flags().StringVarP(&depGo, "go", "g", "", "Navigate to a repository by name (prints path)")

	// Deprecated tag flags (migrated to tag subcommand in spec 10)
	root.Flags().BoolVarP(&depAddTag, "add-tag", "d", false, "Add a tag to repositories")
	root.Flags().BoolVarP(&depRemoveTag, "remove-tag", "x", false, "Remove a tag from repositories")

	// Deprecated filter flags (compound usage: gws --list --type github)
	root.Flags().StringVarP(&depFilterType, "type", "y", "", "Filter by repository type")
	root.Flags().StringSliceVar(&filterTags, "tag", []string{}, "Filter by custom tag(s)")
	root.Flags().StringVarP(&depFilterName, "name", "n", "", "Filter by repository name")
	root.Flags().StringVarP(&depFilterPath, "path", "p", "", "Filter by repository path")
	root.Flags().StringVarP(&depOutputFormat, "output", "o", "table", "Output format: table, json")
	root.Flags().BoolVarP(&depShowStatus, "status", "s", false, "Show git status")
	root.Flags().BoolVar(&depShowUser, "show-user", false, "Show git user info")

	// Hide all deprecated flags from help output
	hiddenFlags := []string{
		"list", "init", "add", "recursive", "refresh", "print-workspace", "go",
		"add-tag", "remove-tag",
		"type", "tag", "name", "path", "output", "status", "show-user",
	}
	for _, name := range hiddenFlags {
		_ = root.Flags().MarkHidden(name)
	}
}

// depWarnings maps deprecated flag names to their replacement forms.
var depWarnings = map[string]string{
	"list":            "gws list",
	"init":            "gws init",
	"add":             "gws add [path]",
	"recursive":       "gws add --recursive",
	"refresh":         "gws refresh",
	"print-workspace": "gws print-workspace",
	"go":              "gws <repo-name>",
	"add-tag":         "gws tag add <repo> <tag>",
	"remove-tag":      "gws tag remove <repo> <tag>",
	"type":            "gws list --type",
	"name":            "gws list --name",
	"path":            "gws list --path",
	"output":          "gws list --output",
	"status":          "gws list --status",
	"show-user":       "gws list --show-user",
}

// emitDeprecationWarnings prints warnings for all deprecated flags that were set.
func emitDeprecationWarnings(cmd *cobra.Command) {
	for flagName, newForm := range depWarnings {
		if cmd.Flags().Changed(flagName) {
			fmt.Fprintf(os.Stderr, "Warning: --%s is deprecated, use '%s' instead\n", flagName, newForm)
		}
	}
}

// handleDeprecatedFlags checks if any deprecated flag is set and dispatches
// to the corresponding logic. Returns (handled bool, err error).
// If handled is true, the caller should return err and skip further dispatch.
func handleDeprecatedFlags(cmd *cobra.Command, args []string) (bool, error) {
	// Count active deprecated command flags for mutual exclusivity
	activeCount := 0
	if depList {
		activeCount++
	}
	if depInit {
		activeCount++
	}
	if depAdd != "" {
		activeCount++
	}
	if depRefresh {
		activeCount++
	}
	if depPrintWorkspace {
		activeCount++
	}
	if depGo != "" {
		activeCount++
	}

	// Also count tag alias, deprecated tag flags, and user flags on root
	if flagTagAlias {
		activeCount++
	}
	if depAddTag {
		activeCount++
	}
	if depRemoveTag {
		activeCount++
	}
	if flagUser || flagListUsers {
		activeCount++
	}

	if activeCount > 1 {
		return true, fmt.Errorf("only one command can be used at a time")
	}

	// No deprecated command flag set — check if filter-only flags were used without --list
	if !depList && !flagUser && hasDeprecatedFilterFlags(cmd) {
		return true, fmt.Errorf("filter flags (--type, --name, --path, --output, --status, --show-user) require --list/-l or 'gws list'")
	}

	// --tag on root requires --list or --user
	if !depList && !flagUser && cmd.Flags().Changed("tag") {
		return true, fmt.Errorf("--tag requires --list/-l or --user to be set")
	}

	// Validate --recursive requires --add
	if depRecursive && depAdd == "" {
		return true, fmt.Errorf("--recursive/-v requires --add/-a to be set")
	}

	// Validate --go and positional args are not both provided
	if depGo != "" && len(args) > 0 {
		return true, fmt.Errorf("cannot use both --go flag and positional argument for navigation")
	}

	// Dispatch deprecated init
	if depInit {
		emitDeprecationWarnings(cmd)
		return true, runInit("")
	}

	// Dispatch deprecated add
	if depAdd != "" {
		emitDeprecationWarnings(cmd)
		return true, runAdd(depAdd, depRecursive)
	}

	// Validate arg counts before workspace check (fail fast on bad input)
	if depAddTag && len(args) != 2 {
		return true, fmt.Errorf("--add-tag requires exactly 2 arguments: <repository> <tag>")
	}
	if depRemoveTag && len(args) != 2 {
		return true, fmt.Errorf("--remove-tag requires exactly 2 arguments: <repository> <tag>")
	}

	// These require workspace to exist
	if depList || depRefresh || depPrintWorkspace || depGo != "" || depAddTag || depRemoveTag {
		exists, err := config.Exists()
		if err != nil {
			return true, fmt.Errorf("failed to check workspace status: %w", err)
		}
		if !exists {
			fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To get started, navigate to your projects directory and run:")
			fmt.Fprintln(os.Stderr, "  gws init")
			return true, fmt.Errorf("workspace not initialized")
		}
	}

	// Dispatch deprecated --add-tag (arg count already validated above)
	if depAddTag {
		emitDeprecationWarnings(cmd)
		return true, runAddTag(args[0], args[1])
	}

	// Dispatch deprecated --remove-tag (arg count already validated above)
	if depRemoveTag {
		emitDeprecationWarnings(cmd)
		return true, runRemoveTag(args[0], args[1])
	}

	// Dispatch deprecated print-workspace
	if depPrintWorkspace {
		emitDeprecationWarnings(cmd)
		return true, runPrintWorkspace()
	}

	// Dispatch deprecated list (with optional compound filter flags)
	if depList {
		emitDeprecationWarnings(cmd)
		return true, runList(ListOptions{
			FilterType:   depFilterType,
			FilterTags:   filterTags,
			FilterName:   depFilterName,
			FilterPath:   depFilterPath,
			OutputFormat: depOutputFormat,
			ShowStatus:   depShowStatus,
			ShowUser:     depShowUser,
		})
	}

	// Dispatch deprecated refresh
	if depRefresh {
		emitDeprecationWarnings(cmd)
		return true, runRefresh()
	}

	// Dispatch deprecated --go
	if depGo != "" {
		emitDeprecationWarnings(cmd)
		cfg, err := config.Load()
		if err != nil {
			return true, fmt.Errorf("failed to load workspace configuration: %w", err)
		}
		return true, runNavigate(depGo, flagQuiet, cfg.Repositories, os.Stderr, os.Stdout, os.Stdin)
	}

	return false, nil
}

// hasDeprecatedFilterFlags checks if any deprecated filter flag (other than --tag) was set.
func hasDeprecatedFilterFlags(cmd *cobra.Command) bool {
	filterFlags := []string{"type", "name", "path", "status", "show-user"}
	for _, name := range filterFlags {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return cmd.Flags().Changed("output") && depOutputFormat != "table"
}
