package main

import (
	"fmt"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/spf13/cobra"
)

var (
	filterType string
	filterTag  string
	filterName string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked repositories with their metadata",
	Long: `List all tracked repositories with their classification metadata including
type (GitHub, GitLab, ADO, Bitbucket), visibility (public/private), and custom tags.

Supports filtering by type, tag, or name.

Examples:
  gws list                          # List all repositories
  gws list --type github            # List only GitHub repositories
  gws list --tag personal           # List repositories tagged as "personal"
  gws list --name myproject         # List repositories matching "myproject"
  gws list --type gitlab --tag work # Combine multiple filters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Repositories) == 0 {
			fmt.Println("No repositories found. Run 'gws init <directory>' to discover repositories.")
			return nil
		}

		// Filter repositories
		filtered := filterRepositories(cfg.Repositories)

		if len(filtered) == 0 {
			fmt.Println("No repositories match the specified filters.")
			return nil
		}

		// Display results
		fmt.Printf("Found %d %s:\n\n", len(filtered), pluralize("repository", "repositories", len(filtered)))
		displayRepositories(filtered)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&filterType, "type", "", "Filter by repository type (github, gitlab, ado, bitbucket)")
	listCmd.Flags().StringVar(&filterTag, "tag", "", "Filter by custom tag")
	listCmd.Flags().StringVar(&filterName, "name", "", "Filter by repository name (partial match)")
}

// filterRepositories applies the specified filters to the repository list
func filterRepositories(repos []config.Repository) []config.Repository {
	filtered := []config.Repository{}

	for _, repo := range repos {
		// Apply type filter
		if filterType != "" && !strings.EqualFold(string(repo.Type), filterType) {
			continue
		}

		// Apply tag filter
		if filterTag != "" {
			hasTag := false
			for _, tag := range repo.Tags {
				if strings.EqualFold(tag, filterTag) {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		// Apply name filter (case-insensitive partial match)
		if filterName != "" && !strings.Contains(strings.ToLower(repo.Name), strings.ToLower(filterName)) {
			continue
		}

		filtered = append(filtered, repo)
	}

	return filtered
}

// displayRepositories shows repositories in a formatted table
func displayRepositories(repos []config.Repository) {
	// Calculate column widths
	maxNameLen := 10
	maxTypeLen := 10
	maxVisLen := 10
	maxTagsLen := 10

	for _, repo := range repos {
		if len(repo.Name) > maxNameLen {
			maxNameLen = len(repo.Name)
		}
		if len(repo.Type) > maxTypeLen {
			maxTypeLen = len(repo.Type)
		}
		if len(repo.Visibility) > maxVisLen {
			maxVisLen = len(repo.Visibility)
		}
		tags := strings.Join(repo.Tags, ", ")
		if len(tags) > maxTagsLen {
			maxTagsLen = len(tags)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
		maxNameLen, "NAME",
		maxTypeLen, "TYPE",
		maxVisLen, "VISIBILITY",
		maxTagsLen, "TAGS",
		"PATH")
	fmt.Printf("%s  %s  %s  %s  %s\n",
		strings.Repeat("-", maxNameLen),
		strings.Repeat("-", maxTypeLen),
		strings.Repeat("-", maxVisLen),
		strings.Repeat("-", maxTagsLen),
		strings.Repeat("-", 4))

	// Print repositories
	for _, repo := range repos {
		tags := strings.Join(repo.Tags, ", ")
		if tags == "" {
			tags = "-"
		}

		typeStr := string(repo.Type)
		if typeStr == "" {
			typeStr = "unknown"
		}

		visStr := string(repo.Visibility)
		if visStr == "" {
			visStr = "unknown"
		}

		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
			maxNameLen, repo.Name,
			maxTypeLen, typeStr,
			maxVisLen, visStr,
			maxTagsLen, tags,
			repo.Path)
	}
}
