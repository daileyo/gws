package main

import (
	"fmt"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag <repository> <tag>",
	Short: "Add a custom tag to repositories",
	Long: `Add a custom tag to all repositories matching the given name or path.

The repository identifier can match by:
- Name (case-insensitive, partial match)
- Path (exact match)

If multiple repositories match, the tag will be applied to all of them.

Examples:
  gws tag my-project personal        # Tag all repos with "my-project" in name
  gws tag api work                   # Tag all repos with "api" in name
  gws tag /path/to/repo production   # Tag specific repo by exact path`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoIdentifier := args[0]
		tag := args[1]

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
			fmt.Printf("Added tag '%s' to %d %s\n", tag, taggedCount, pluralize("repository", "repositories", taggedCount))
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
			fmt.Printf("%d %s already had tag '%s'\n", skippedCount, pluralize("repository", "repositories", skippedCount), tag)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
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
