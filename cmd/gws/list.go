package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/spf13/cobra"
)

var (
	filterType   string
	filterTags   []string
	filterName   string
	filterPath   string
	outputFormat string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked repositories with their metadata",
	Long: `List all tracked repositories with their classification metadata including
type (GitHub, GitLab, ADO, Bitbucket), visibility (public/private), and custom tags.

Supports filtering by type, tags, name, and path. Multiple filters can be combined.

Examples:
  gws list                                    # List all repositories
  gws list --type github                      # List only GitHub repositories
  gws list --tag personal                     # List repositories tagged as "personal"
  gws list --tag work --tag backend           # List repositories with both "work" AND "backend" tags
  gws list --name myproject                   # List repositories matching "myproject"
  gws list --path /home/user/projects         # List repositories in specific path
  gws list --type gitlab --tag work           # Combine multiple filters
  gws list --output json                      # Output in JSON format
  gws list -o json                            # Short form for JSON output`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Repositories) == 0 {
			if outputFormat == "json" {
				fmt.Println("[]")
			} else {
				fmt.Println("No repositories found. Run 'gws init <directory>' to discover repositories.")
			}
			return nil
		}

		// Build filter criteria
		criteria := filter.Criteria{
			Type: filterType,
			Tags: filterTags,
			Name: filterName,
			Path: filterPath,
		}

		// Apply filters
		filtered := filter.Apply(cfg.Repositories, criteria)

		if len(filtered) == 0 {
			if outputFormat == "json" {
				fmt.Println("[]")
			} else {
				fmt.Println("No repositories match the specified filters.")
			}
			return nil
		}

		// Display results based on format
		switch outputFormat {
		case "json":
			return displayJSON(filtered)
		default:
			displayTable(filtered)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&filterType, "type", "", "Filter by repository type (github, gitlab, ado, bitbucket)")
	listCmd.Flags().StringSliceVar(&filterTags, "tag", []string{}, "Filter by custom tag(s) - can be specified multiple times for AND logic")
	listCmd.Flags().StringVar(&filterName, "name", "", "Filter by repository name (partial match)")
	listCmd.Flags().StringVar(&filterPath, "path", "", "Filter by repository path (partial match)")
	listCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")
}

// displayTable shows repositories in a formatted table
func displayTable(repos []config.Repository) {
	fmt.Printf("Found %d %s:\n\n", len(repos), pluralize("repository", "repositories", len(repos)))

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

// displayJSON outputs repositories in JSON format
func displayJSON(repos []config.Repository) error {
	data, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
