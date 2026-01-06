package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

var untagCmd = &cobra.Command{
	Use:   "untag <repository> <tag>",
	Short: "Remove a custom tag from repositories",
	Long: `Remove a custom tag from all repositories matching the given name or path.

The repository identifier can match by:
- Name (case-insensitive, partial match)
- Path (exact match)

If multiple repositories match, the tag will be removed from all of them.

Examples:
  gws untag my-project personal        # Remove "personal" tag from all repos with "my-project" in name
  gws untag api work                   # Remove "work" tag from all repos with "api" in name
  gws untag /path/to/repo production   # Remove tag from specific repo by exact path`,
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

		// Remove tag from all matching repositories
		untaggedCount := 0
		skippedCount := 0
		for _, repo := range repos {
			// Find and remove tag
			found := false
			newTags := []string{}
			for _, existingTag := range repo.Tags {
				if existingTag == tag {
					found = true
				} else {
					newTags = append(newTags, existingTag)
				}
			}

			if !found {
				skippedCount++
				continue
			}

			// Update tags
			repo.Tags = newTags
			untaggedCount++
		}

		// Save configuration
		if untaggedCount > 0 {
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}
		}

		// Report results
		if untaggedCount > 0 {
			fmt.Printf("Removed tag '%s' from %d %s\n", tag, untaggedCount, pluralize(untaggedCount, "repository", "repositories"))
			if untaggedCount <= 5 {
				for _, repo := range repos {
					// Check if tag was removed from this repo
					hadTag := false
					for _, existingTag := range repo.Tags {
						if existingTag == tag {
							hadTag = true
							break
						}
					}
					if !hadTag && skippedCount < len(repos) {
						// This repo had the tag removed
						fmt.Printf("  - %s\n", repo.Name)
					}
				}
			}
		}

		if skippedCount > 0 {
			fmt.Printf("%d %s did not have tag '%s'\n", skippedCount, pluralize(skippedCount, "repository", "repositories"), tag)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(untagCmd)
}
