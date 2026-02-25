package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

// runTag handles the --add-tag flag logic
func runTag(_ *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("--add-tag requires exactly 2 arguments: <repository> <tag>")
	}

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
