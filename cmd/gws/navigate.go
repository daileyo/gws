package main

import (
	"fmt"
	"io"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
)

// runNavigate handles repository navigation by matching a query against repository names.
// It prints the matching repository path to stdout. In verbose mode (quiet=false),
// it also prints repo details to stderr so shell command substitution works correctly.
func runNavigate(query string, quiet bool, repos []config.Repository, stderr io.Writer, stdout io.Writer, stdin io.Reader) error {
	// Find matching repositories by name
	var matches []config.Repository
	for _, repo := range repos {
		if filter.MatchesPattern(repo.Name, query) {
			matches = append(matches, repo)
		}
	}

	switch len(matches) {
	case 0:
		return handleNoMatch(query, repos, stderr)
	case 1:
		return printMatch(matches[0], quiet, stderr, stdout)
	default:
		return handleMultipleMatches(matches, query, quiet, stderr, stdout, stdin)
	}
}

// printMatch outputs a single matched repository
func printMatch(repo config.Repository, quiet bool, stderr io.Writer, stdout io.Writer) error {
	if !quiet {
		typeStr := string(repo.Type)
		if typeStr == "" {
			typeStr = "unknown"
		}
		fmt.Fprintf(stderr, "%s (%s) → %s\n", repo.Name, typeStr, repo.Path)
	}
	fmt.Fprintln(stdout, repo.Path)
	return nil
}

// handleNoMatch displays an error when no repositories match (placeholder for task 3.0)
func handleNoMatch(query string, repos []config.Repository, stderr io.Writer) error {
	fmt.Fprintf(stderr, "No repositories found matching '%s'\n", query)
	return fmt.Errorf("no repositories found matching '%s'", query)
}

// handleMultipleMatches handles the case when multiple repos match (placeholder for task 2.0)
func handleMultipleMatches(matches []config.Repository, query string, quiet bool, stderr io.Writer, stdout io.Writer, stdin io.Reader) error {
	// Placeholder: will be fully implemented in task 2.0
	// For now, return the first match
	return printMatch(matches[0], quiet, stderr, stdout)
}
