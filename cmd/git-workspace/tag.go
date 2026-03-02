package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

// Tag subcommand flags for short-flag invocation (tag -a / tag -d).
var (
	tagFlagAdd    bool
	tagFlagDelete bool
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage repository tags",
	Long: `Manage tags on tracked repositories.

Examples:
  gws tag add my-repo work          # Add the "work" tag
  gws tag -a my-repo work           # Same, using short flag
  gws tag remove my-repo work       # Remove the "work" tag
  gws tag -d my-repo work           # Same, using short flag
  gws tag                           # Show this help`,
	Args:              cobra.ArbitraryArgs,
	ValidArgsFunction: completeRepoThenNone,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Mutual exclusivity: -a and -d cannot both be set
		if tagFlagAdd && tagFlagDelete {
			return fmt.Errorf("-a (add) and -d (delete) are mutually exclusive")
		}

		if tagFlagAdd {
			if len(args) != 2 {
				return fmt.Errorf("tag -a requires exactly 2 arguments: <repository> <tag>")
			}
			return runAddTag(args[0], args[1])
		}

		if tagFlagDelete {
			if len(args) != 2 {
				return fmt.Errorf("tag -d requires exactly 2 arguments: <repository> <tag>")
			}
			return runRemoveTag(args[0], args[1])
		}

		// No sub-operation specified — show help
		return cmd.Help()
	},
}

var tagAddCmd = &cobra.Command{
	Use:   "add <repo> <tag>",
	Short: "Add a tag to repositories",
	Long: `Add a tag to all repositories matching the identifier.

The repository can be specified by name (partial match, case-insensitive) or exact path.

Examples:
  gws tag add my-repo work
  gws tag add api backend`,
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeRepoThenNone,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAddTag(args[0], args[1])
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove <repo> <tag>",
	Short: "Remove a tag from repositories",
	Long: `Remove a tag from all repositories matching the identifier.

The repository can be specified by name (partial match, case-insensitive) or exact path.

Examples:
  gws tag remove my-repo work
  gws tag remove api backend`,
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeRepoThenTags,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoveTag(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)

	// Short flags on the tag command for quick invocation
	tagCmd.Flags().BoolVarP(&tagFlagAdd, "add", "a", false, "Add a tag (equivalent to 'tag add')")
	tagCmd.Flags().BoolVarP(&tagFlagDelete, "delete", "d", false, "Remove a tag (equivalent to 'tag remove')")

	// Reset to Cobra's default template so tag help doesn't inherit root's custom template
	tagCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

// completeRepoThenNone completes repo names for the first arg, nothing for the second.
func completeRepoThenNone(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeRepoNames(toComplete)
}

// completeRepoThenTags completes repo names for the first arg, existing tags for the second.
func completeRepoThenTags(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeRepoNames(toComplete)
	}
	if len(args) == 1 {
		return completeRepoTags(args[0], toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// completeRepoNames returns repo names matching the prefix.
func completeRepoNames(toComplete string) ([]string, cobra.ShellCompDirective) {
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

// completeRepoTags returns existing tags for repos matching the identifier.
func completeRepoTags(repoIdentifier, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, err := config.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	repos := findRepositories(cfg, repoIdentifier)
	seen := make(map[string]bool)
	var tags []string
	for _, repo := range repos {
		for _, tag := range repo.Tags {
			if strings.HasPrefix(strings.ToLower(tag), strings.ToLower(toComplete)) && !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	return tags, cobra.ShellCompDirectiveNoFileComp
}

// runAddTag adds a tag to all repositories matching the identifier.
func runAddTag(repoIdentifier, tag string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Find all matching repositories
	repos := findRepositories(cfg, repoIdentifier)
	if len(repos) == 0 {
		return fmt.Errorf("no repositories found matching: %s", repoIdentifier)
	}

	// Apply tag to all matching repositories
	taggedCount := 0
	skippedCount := 0
	for _, repo := range repos {
		// Check if tag already exists
		hasTag := false
		for _, existingTag := range repo.Tags {
			if existingTag == tag {
				hasTag = true
				break
			}
		}

		if hasTag {
			skippedCount++
			continue
		}

		// Add tag
		repo.Tags = append(repo.Tags, tag)
		taggedCount++
	}

	// Save configuration
	if taggedCount > 0 {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	// Report results
	if taggedCount > 0 {
		fmt.Printf("Added tag '%s' to %d %s\n", tag, taggedCount, pluralize(taggedCount, "repository", "repositories"))
		if taggedCount <= 5 {
			for _, repo := range repos {
				// Check if this repo was tagged (doesn't already have the tag)
				hasTag := false
				for _, existingTag := range repo.Tags {
					if existingTag == tag {
						hasTag = true
						break
					}
				}
				if hasTag && skippedCount < len(repos) {
					// This repo was just tagged
					fmt.Printf("  - %s\n", repo.Name)
				}
			}
		}
	}

	if skippedCount > 0 {
		fmt.Printf("%d %s already had tag '%s'\n", skippedCount, pluralize(skippedCount, "repository", "repositories"), tag)
	}

	return nil
}

// findRepositories locates all repositories matching the identifier by name or path
func findRepositories(cfg *config.Config, identifier string) []*config.Repository {
	var matched []*config.Repository

	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]
		// Match by exact path or partial name (case-insensitive)
		if repo.Path == identifier || strings.Contains(strings.ToLower(repo.Name), strings.ToLower(identifier)) {
			matched = append(matched, repo)
		}
	}

	return matched
}

// findRepositoriesWithFilters locates repositories matching the given filters.
// repoFilter performs a partial, case-insensitive name match.
// pathFilter performs a case-sensitive prefix match; falls back to substring if no prefix matches.
// When both are non-empty, a repo must satisfy both conditions (AND logic).
func findRepositoriesWithFilters(cfg *config.Config, repoFilter, pathFilter string) []*config.Repository {
	candidates := make([]*config.Repository, 0, len(cfg.Repositories))
	for i := range cfg.Repositories {
		candidates = append(candidates, &cfg.Repositories[i])
	}

	// Apply name filter (partial, case-insensitive)
	if repoFilter != "" {
		var filtered []*config.Repository
		for _, repo := range candidates {
			if strings.Contains(strings.ToLower(repo.Name), strings.ToLower(repoFilter)) {
				filtered = append(filtered, repo)
			}
		}
		candidates = filtered
	}

	// Apply path filter (case-sensitive prefix, with substring fallback)
	if pathFilter != "" {
		var prefixMatches []*config.Repository
		for _, repo := range candidates {
			if strings.HasPrefix(repo.Path, pathFilter) {
				prefixMatches = append(prefixMatches, repo)
			}
		}
		if len(prefixMatches) > 0 {
			candidates = prefixMatches
		} else {
			var substringMatches []*config.Repository
			for _, repo := range candidates {
				if strings.Contains(repo.Path, pathFilter) {
					substringMatches = append(substringMatches, repo)
				}
			}
			candidates = substringMatches
		}
	}

	return candidates
}
