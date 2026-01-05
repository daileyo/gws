package main

import (
	"fmt"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
	"github.com/daileyo/gws/internal/git"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh repository metadata and git status cache",
	Long: `Refresh the gws workspace by:
- Re-scanning for new repositories in the workspace
- Updating classification for existing repositories
- Clearing and rebuilding the git status cache

This is useful when:
- You've added new repositories to your workspace
- You want to update cached git status information
- Remote URLs have changed

Examples:
  gws refresh                 # Refresh everything
  gws refresh && gws list -s  # Refresh and show updated status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load existing configuration
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Printf("Refreshing workspace at: %s\n", cfg.Workspace)
		fmt.Println("Re-scanning for repositories...")

		// Re-scan the workspace
		result, err := discovery.Scan(cfg.Workspace)
		if err != nil {
			return fmt.Errorf("failed to scan workspace: %w", err)
		}

		// Display scan errors if any
		if len(result.Errors) > 0 {
			fmt.Printf("\nWarning: %d errors occurred during scanning:\n", len(result.Errors))
			for i, err := range result.Errors {
				if i < 3 { // Limit to first 3 errors
					fmt.Printf("  - %v\n", err)
				}
			}
			if len(result.Errors) > 3 {
				fmt.Printf("  ... and %d more errors\n", len(result.Errors)-3)
			}
		}

		// Preserve existing tags by merging old and new
		oldRepoMap := make(map[string]config.Repository)
		for _, repo := range cfg.Repositories {
			oldRepoMap[repo.Path] = repo
		}

		// Merge tags from old repos to new repos
		for i := range result.Repositories {
			newRepo := &result.Repositories[i]
			if oldRepo, exists := oldRepoMap[newRepo.Path]; exists {
				// Preserve tags from old repo
				newRepo.Tags = oldRepo.Tags
			}
		}

		// Update configuration
		cfg.Repositories = result.Repositories

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// Clear git status cache
		cachePath, err := git.GetCachePath()
		if err == nil {
			statusCache := git.NewCache(git.DefaultTTL)
			statusCache.Clear()
			_ = statusCache.Save(cachePath) // Clear the cache file
			fmt.Println("Cleared git status cache")
		}

		// Report results
		fmt.Printf("\nRefresh complete!\n")
		fmt.Printf("Total repositories: %d\n", result.Count)

		// Show new vs existing
		newCount := 0
		for _, newRepo := range result.Repositories {
			if _, exists := oldRepoMap[newRepo.Path]; !exists {
				newCount++
			}
		}

		if newCount > 0 {
			fmt.Printf("New repositories found: %d\n", newCount)
		}

		removedCount := len(cfg.Repositories) - result.Count + newCount
		if removedCount > 0 {
			fmt.Printf("Repositories no longer found: %d\n", removedCount)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
