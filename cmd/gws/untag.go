package main

import (
	"fmt"

	"github.com/daileyo/gws/internal/config"
	"github.com/spf13/cobra"
)

var untagCmd = &cobra.Command{
	Use:   "untag <repository> <tag>",
	Short: "Remove a custom tag from a repository",
	Long: `Remove a custom tag from a repository.

The repository can be identified by name or path.

Examples:
  gws untag my-project personal
  gws untag /path/to/repo work
  gws untag my-app production`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoIdentifier := args[0]
		tag := args[1]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Find repository
		repo, err := findRepository(cfg, repoIdentifier)
		if err != nil {
			return err
		}

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
			return fmt.Errorf("repository '%s' does not have tag '%s'", repo.Name, tag)
		}

		// Update tags
		repo.Tags = newTags

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Removed tag '%s' from repository '%s'\n", tag, repo.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(untagCmd)
}
