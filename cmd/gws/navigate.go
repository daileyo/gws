package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
)

const maxSelectionAttempts = 3

// isTerminalFunc is the function used to check if stdin is a TTY.
// It can be overridden in tests to simulate TTY/non-TTY behavior.
var isTerminalFunc = isTerminal

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

// isTerminal checks if the given reader is connected to a terminal (TTY)
func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// handleMultipleMatches handles the case when multiple repos match the query.
// In TTY mode, it displays a numbered list and prompts for selection.
// In non-TTY mode, it prints all matching paths and returns an error.
func handleMultipleMatches(matches []config.Repository, query string, quiet bool, stderr io.Writer, stdout io.Writer, stdin io.Reader) error {
	if !isTerminalFunc(stdin) {
		// Non-TTY: print all matching paths and exit with error
		for _, repo := range matches {
			fmt.Fprintln(stdout, repo.Path)
		}
		return fmt.Errorf("multiple repositories match '%s' (%d matches)", query, len(matches))
	}

	// TTY: display numbered list and prompt for selection
	fmt.Fprintf(stderr, "Multiple repositories match '%s':\n\n", query)
	for i, repo := range matches {
		typeStr := string(repo.Type)
		if typeStr == "" {
			typeStr = "unknown"
		}
		fmt.Fprintf(stderr, "  %d) %s (%s) %s\n", i+1, repo.Name, typeStr, repo.Path)
	}
	fmt.Fprintln(stderr)

	scanner := bufio.NewScanner(stdin)
	for attempt := 0; attempt < maxSelectionAttempts; attempt++ {
		fmt.Fprintf(stderr, "Select repository [1-%d]: ", len(matches))

		if !scanner.Scan() {
			return fmt.Errorf("failed to read selection")
		}

		input := strings.TrimSpace(scanner.Text())
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(matches) {
			fmt.Fprintf(stderr, "Invalid selection '%s'. Please enter a number between 1 and %d.\n", input, len(matches))
			continue
		}

		return printMatch(matches[num-1], quiet, stderr, stdout)
	}

	return fmt.Errorf("too many invalid selection attempts")
}
