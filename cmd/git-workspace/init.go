package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

// runInit handles the --init flag logic
func runInit(_ *cobra.Command, _ []string) error {
	workspacePath := flagInit

	// Convert to absolute path
	absPath, err := filepath.Abs(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to resolve workspace path: %w", err)
	}

	fmt.Printf("Initializing gws workspace at: %s\n", absPath)
	fmt.Println("Scanning for git repositories...")

	// Scan for repositories
	result, err := discovery.Scan(absPath)
	if err != nil {
		return fmt.Errorf("failed to scan workspace: %w", err)
	}

	// Display scan errors if any
	if len(result.Errors) > 0 {
		fmt.Printf("\nWarning: %d errors occurred during scanning:\n", len(result.Errors))
		for i, err := range result.Errors {
			if i < 5 { // Limit to first 5 errors
				fmt.Printf("  - %v\n", err)
			}
		}
		if len(result.Errors) > 5 {
			fmt.Printf("  ... and %d more errors\n", len(result.Errors)-5)
		}
	}

	// Create and save configuration
	cfg := config.New(absPath)
	cfg.Repositories = result.Repositories

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Display results
	fmt.Printf("\nWorkspace initialized successfully!\n")
	fmt.Printf("Found %d git %s:\n\n", result.Count, pluralize(result.Count, "repository", "repositories"))

	for i, repo := range result.Repositories {
		if i < 10 { // Limit to first 10 repositories
			fmt.Printf("  - %s\n", repo.Name)
			if repo.RemoteURL != "" {
				fmt.Printf("    %s\n", repo.RemoteURL)
			}
		}
	}

	if result.Count > 10 {
		fmt.Printf("  ... and %d more %s\n", result.Count-10, pluralize(result.Count-10, "repository", "repositories"))
	}

	configPath, _ := config.GetConfigPath()
	fmt.Printf("\nConfiguration saved to: %s\n", configPath)

	return nil
}

func pluralize(count int, singular, plural string) string { //nolint:unparam // singular is always "repository" but kept for readability
	if count == 1 {
		return singular
	}
	return plural
}
