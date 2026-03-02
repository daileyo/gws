package main

import (
	"fmt"

	"github.com/daileyo/gws/internal/config"
)

// runRemoveTagWithFilters removes a tag from all repositories matching the given filters.
func runRemoveTagWithFilters(repoFilter, pathFilter, tag string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	repos := findRepositoriesWithFilters(cfg, repoFilter, pathFilter)
	if len(repos) == 0 {
		return noMatchError(repoFilter, pathFilter)
	}

	untaggedCount := 0
	skippedCount := 0
	for _, repo := range repos {
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
		repo.Tags = newTags
		untaggedCount++
	}

	if untaggedCount > 0 {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	if untaggedCount > 0 {
		fmt.Printf("Removed tag '%s' from %d %s\n", tag, untaggedCount, pluralize(untaggedCount, "repository", "repositories"))
		if untaggedCount <= 5 {
			for _, repo := range repos {
				hadTag := false
				for _, existingTag := range repo.Tags {
					if existingTag == tag {
						hadTag = true
						break
					}
				}
				if !hadTag && skippedCount < len(repos) {
					if pathFilter != "" {
						fmt.Printf("  - %s  %s\n", repo.Name, repo.Path)
					} else {
						fmt.Printf("  - %s\n", repo.Name)
					}
				}
			}
		}
	}

	if skippedCount > 0 {
		fmt.Printf("%d %s did not have tag '%s'\n", skippedCount, pluralize(skippedCount, "repository", "repositories"), tag)
	}

	return nil
}

// runRemoveTag removes a tag from all repositories matching the identifier.
func runRemoveTag(repoIdentifier, tag string) error {
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
}
