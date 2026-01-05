package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
	"github.com/spf13/cobra"
)

var (
	filterType   string
	filterTags   []string
	filterName   string
	filterPath   string
	outputFormat string
	showStatus   bool
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

		// Load git status cache if status display is enabled
		var statusCache *git.Cache
		if showStatus {
			statusCache = git.NewCache(git.DefaultTTL)
			cachePath, err := git.GetCachePath()
			if err == nil {
				_ = statusCache.Load(cachePath) // Ignore errors, cache might not exist yet
			}
		}

		// Display results based on format
		switch outputFormat {
		case "json":
			return displayJSON(filtered)
		default:
			displayTable(filtered, statusCache)
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
	listCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show git status (branch, clean/dirty, ahead/behind)")
}

// displayTable shows repositories in a formatted table
func displayTable(repos []config.Repository, statusCache *git.Cache) {
	fmt.Printf("Found %d %s:\n\n", len(repos), pluralize("repository", "repositories", len(repos)))

	// Calculate column widths
	maxNameLen := 10
	maxTypeLen := 10
	maxVisLen := 10
	maxTagsLen := 10
	maxStatusLen := 6 // "STATUS"

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

		// Calculate status column width
		if statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			if status != nil {
				statusStr := formatStatusShort(status)
				if len(statusStr) > maxStatusLen {
					maxStatusLen = len(statusStr)
				}
			}
		}
	}

	// Print header
	if statusCache != nil {
		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
			maxNameLen, "NAME",
			maxStatusLen, "STATUS",
			maxTypeLen, "TYPE",
			maxVisLen, "VISIBILITY",
			maxTagsLen, "TAGS",
			"PATH")
		fmt.Printf("%s  %s  %s  %s  %s  %s\n",
			strings.Repeat("-", maxNameLen),
			strings.Repeat("-", maxStatusLen),
			strings.Repeat("-", maxTypeLen),
			strings.Repeat("-", maxVisLen),
			strings.Repeat("-", maxTagsLen),
			strings.Repeat("-", 4))
	} else {
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
	}

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

		if statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			statusStr := "-"
			if status != nil {
				statusStr = formatStatusShort(status)
			}

			fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
				maxNameLen, repo.Name,
				maxStatusLen, statusStr,
				maxTypeLen, typeStr,
				maxVisLen, visStr,
				maxTagsLen, tags,
				repo.Path)
		} else {
			fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
				maxNameLen, repo.Name,
				maxTypeLen, typeStr,
				maxVisLen, visStr,
				maxTagsLen, tags,
				repo.Path)
		}
	}

	// Save cache if we fetched any statuses
	if statusCache != nil {
		cachePath, err := git.GetCachePath()
		if err == nil {
			_ = statusCache.Save(cachePath) // Ignore save errors
		}
	}
}

// formatStatusShort returns a short status string with visual indicators
func formatStatusShort(status *git.Status) string {
	if status.Branch == "" {
		return "no commits"
	}

	result := status.Branch

	// Add clean/dirty indicator
	if status.HasChanges {
		result += " ✗"
	} else {
		result += " ✓"
	}

	// Add ahead/behind indicators
	if status.Ahead > 0 && status.Behind > 0 {
		result += fmt.Sprintf(" ↑%d↓%d", status.Ahead, status.Behind)
	} else if status.Ahead > 0 {
		result += fmt.Sprintf(" ↑%d", status.Ahead)
	} else if status.Behind > 0 {
		result += fmt.Sprintf(" ↓%d", status.Behind)
	}

	return result
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
