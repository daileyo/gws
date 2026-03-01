package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
)

// runList handles the --list flag logic
func runList(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Repositories) == 0 {
		if outputFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("No repositories found. Run 'gws --init <directory>' to discover repositories.")
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
}

// displayTable shows repositories in a formatted table
func displayTable(repos []config.Repository, statusCache *git.Cache) {
	fmt.Printf("Found %d %s:\n\n", len(repos), pluralize(len(repos), "repository", "repositories"))

	// Calculate column widths
	maxNameLen := 4   // "NAME"
	maxTypeLen := 4   // "TYPE"
	maxVisLen := 10   // "VISIBILITY"
	maxTagsLen := 4   // "TAGS"
	maxStatusLen := 6 // "STATUS"
	maxUserLen := 4   // "USER"
	maxEmailLen := 5  // "EMAIL"

	// Pre-compute user info and drift status for width calculation and display
	type repoUserInfo struct {
		userDisplay string
		email       string
		signStr     string
		hasDrift    bool
	}
	userInfoMap := make(map[string]repoUserInfo)

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

		// Calculate user column widths and pre-compute display values
		if showUser {
			info := repoUserInfo{}

			// Start with stored config values
			displayUser := repo.User
			displayEmail := repo.Email
			displaySign := repo.SigningEnabled
			displaySource := repo.UserSource

			// Read effective config at display time to detect local overrides
			currentCfg, err := git.GetUserConfig(repo.Path)
			if err == nil && currentCfg != nil {
				if currentCfg.Source == config.UserSourceLocal {
					// Local config detected - show effective local values
					displayUser = currentCfg.Name
					displayEmail = currentCfg.Email
					displaySign = currentCfg.SignCommits
					displaySource = config.UserSourceLocal
				} else if (repo.User != "" && repo.User != currentCfg.Name) ||
					(repo.Email != "" && repo.Email != currentCfg.Email) {
					info.hasDrift = true
				}
			}

			// Format user display with source marker
			info.userDisplay = displayUser
			if info.userDisplay == "" {
				info.userDisplay = "-"
			} else if displaySource == config.UserSourceLocal {
				info.userDisplay += " (local)"
			}

			info.email = displayEmail
			if info.email == "" {
				info.email = "-"
			}

			if displaySign {
				info.signStr = "✓"
			} else {
				info.signStr = ""
			}

			userInfoMap[repo.Path] = info

			if len(info.userDisplay) > maxUserLen {
				maxUserLen = len(info.userDisplay)
			}
			if len(info.email) > maxEmailLen {
				maxEmailLen = len(info.email)
			}
		}
	}

	// Build and print header
	var headerParts []string
	var separatorParts []string

	headerParts = append(headerParts, fmt.Sprintf("%-*s", maxNameLen, "NAME"))
	separatorParts = append(separatorParts, strings.Repeat("-", maxNameLen))

	if statusCache != nil {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxStatusLen, "STATUS"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxStatusLen))
	}

	if showUser {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxUserLen, "USER"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxUserLen))
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxEmailLen, "EMAIL"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxEmailLen))
		headerParts = append(headerParts, fmt.Sprintf("%-4s", "SIGN"))
		separatorParts = append(separatorParts, strings.Repeat("-", 4))
	}

	headerParts = append(headerParts, fmt.Sprintf("%-*s", maxTypeLen, "TYPE"))
	separatorParts = append(separatorParts, strings.Repeat("-", maxTypeLen))
	headerParts = append(headerParts, fmt.Sprintf("%-*s", maxVisLen, "VISIBILITY"))
	separatorParts = append(separatorParts, strings.Repeat("-", maxVisLen))
	headerParts = append(headerParts, fmt.Sprintf("%-*s", maxTagsLen, "TAGS"))
	separatorParts = append(separatorParts, strings.Repeat("-", maxTagsLen))
	headerParts = append(headerParts, "PATH")
	separatorParts = append(separatorParts, strings.Repeat("-", 4))

	fmt.Println(strings.Join(headerParts, "  "))
	fmt.Println(strings.Join(separatorParts, "  "))

	// Print repositories
	for _, repo := range repos {
		var rowParts []string

		// Name column - add drift indicator if applicable
		nameDisplay := repo.Name
		if showUser {
			if info, ok := userInfoMap[repo.Path]; ok && info.hasDrift {
				nameDisplay += " ⚠"
			}
		}
		rowParts = append(rowParts, fmt.Sprintf("%-*s", maxNameLen, nameDisplay))

		// Status column (optional)
		if statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			statusStr := "-"
			if status != nil {
				statusStr = formatStatusShort(status)
			}
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxStatusLen, statusStr))
		}

		// User columns (optional)
		if showUser {
			info := userInfoMap[repo.Path]
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxUserLen, info.userDisplay))
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxEmailLen, info.email))
			rowParts = append(rowParts, fmt.Sprintf("%-4s", info.signStr))
		}

		// Type column
		typeStr := string(repo.Type)
		if typeStr == "" {
			typeStr = "unknown"
		}
		rowParts = append(rowParts, fmt.Sprintf("%-*s", maxTypeLen, typeStr))

		// Visibility column
		visStr := string(repo.Visibility)
		if visStr == "" {
			visStr = "unknown"
		}
		rowParts = append(rowParts, fmt.Sprintf("%-*s", maxVisLen, visStr))

		// Tags column
		tags := strings.Join(repo.Tags, ", ")
		if tags == "" {
			tags = "-"
		}
		rowParts = append(rowParts, fmt.Sprintf("%-*s", maxTagsLen, tags))

		// Path column
		rowParts = append(rowParts, repo.Path)

		fmt.Println(strings.Join(rowParts, "  "))
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
