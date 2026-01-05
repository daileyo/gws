package main

import (
	"fmt"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag <repository> <tag>",
	Short: "Add a custom tag to a repository",
	Long: `Add a custom tag to a repository for organization and filtering.

The repository can be identified by name or path.

Examples:
  gws tag my-project personal
  gws tag /path/to/repo work
  gws tag my-app production`,
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

		// Check if tag already exists
		for _, existingTag := range repo.Tags {
			if existingTag == tag {
				return fmt.Errorf("repository '%s' already has tag '%s'", repo.Name, tag)
			}
		}

		// Add tag
		repo.Tags = append(repo.Tags, tag)

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Added tag '%s' to repository '%s'\n", tag, repo.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

// findRepository locates a repository by name or path
func findRepository(cfg *config.Config, identifier string) (*config.Repository, error) {
	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]
		// Match by name (case-insensitive) or exact path
		if strings.EqualFold(repo.Name, identifier) || repo.Path == identifier {
			return repo, nil
		}
	}
	return nil, fmt.Errorf("repository not found: %s", identifier)
}
