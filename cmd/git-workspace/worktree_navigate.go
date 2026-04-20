package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
)

var worktreeNavigateCmd = &cobra.Command{
	Use:   "navigate <branch-pattern>",
	Short: "Navigate to a worktree by branch name",
	Long: `Navigate to a worktree matching the given branch name or pattern across all
tracked repositories. If multiple worktrees match, an interactive selection
prompt is displayed.

This is the canonical form of the shorthand: gws worktree <branch-pattern>

Examples:
  gws worktree navigate feat-auth       # Navigate to worktree with branch feat-auth
  gws worktree navigate "feat-*"        # Wildcard match across all repos
  gws worktree feat-auth                # Shorthand (same behavior)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return runWorktreeNavigateGlobal(args[0], flagWorktreeQuiet, cfg.Repositories, os.Stderr, os.Stdout, os.Stdin)
	},
}

func init() {
	worktreeNavigateCmd.ValidArgsFunction = completeWorktreeBranches
	worktreeCmd.AddCommand(worktreeNavigateCmd)
}

// globalWorktreeMatch pairs a worktree with its owning repository for display.
type globalWorktreeMatch struct {
	Repo     config.Repository
	Worktree config.Worktree
}

// runWorktreeNavigateGlobal searches all repos' worktrees for branches matching
// the query. Single match navigates directly; multiple matches prompt for selection.
func runWorktreeNavigateGlobal(query string, quiet bool, repos []config.Repository, stderr io.Writer, stdout io.Writer, stdin io.Reader) error {
	var matches []globalWorktreeMatch

	for _, repo := range repos {
		for _, wt := range repo.Worktrees {
			if filter.MatchesPattern(wt.Branch, query) {
				matches = append(matches, globalWorktreeMatch{Repo: repo, Worktree: wt})
			}
		}
	}

	switch len(matches) {
	case 0:
		return handleNoGlobalWorktreeMatch(query, repos, stderr)
	case 1:
		return printGlobalWorktreeMatch(matches[0], quiet, stderr, stdout)
	default:
		return handleGlobalWorktreeSelection(matches, query, quiet, stderr, stdout, stdin)
	}
}

func printGlobalWorktreeMatch(m globalWorktreeMatch, quiet bool, stderr io.Writer, stdout io.Writer) error {
	if !quiet {
		branch := m.Worktree.Branch
		if branch == "" {
			branch = "(detached)"
		}
		fmt.Fprintf(stderr, "%s / %s → %s\n", m.Repo.Name, branch, m.Worktree.Path)
	}
	fmt.Fprintln(stdout, m.Worktree.Path)
	return nil
}

func handleNoGlobalWorktreeMatch(query string, repos []config.Repository, stderr io.Writer) error {
	fmt.Fprintf(stderr, "No worktrees found matching '%s'\n", query)

	// Collect all worktree branches for suggestions
	var allBranches []string
	queryLower := strings.ToLower(query)
	for _, repo := range repos {
		for _, wt := range repo.Worktrees {
			if wt.Branch != "" && isSimilar(queryLower, strings.ToLower(wt.Branch)) {
				allBranches = append(allBranches, fmt.Sprintf("%s (%s)", wt.Branch, repo.Name))
				if len(allBranches) >= maxSuggestions {
					break
				}
			}
		}
		if len(allBranches) >= maxSuggestions {
			break
		}
	}

	if len(allBranches) > 0 {
		fmt.Fprintln(stderr, "\nDid you mean:")
		for _, s := range allBranches {
			fmt.Fprintf(stderr, "  %s\n", s)
		}
	}

	return fmt.Errorf("no worktrees found matching '%s'", query)
}

func handleGlobalWorktreeSelection(matches []globalWorktreeMatch, query string, quiet bool, stderr io.Writer, stdout io.Writer, stdin io.Reader) error {
	if !isTerminalFunc(stdin) {
		for _, m := range matches {
			fmt.Fprintln(stdout, m.Worktree.Path)
		}
		return fmt.Errorf("multiple worktrees match '%s' (%d matches)", query, len(matches))
	}

	fmt.Fprintf(stderr, "Multiple worktrees match '%s':\n\n", query)
	for i, m := range matches {
		branch := m.Worktree.Branch
		if branch == "" {
			branch = "(detached)"
		}
		fmt.Fprintf(stderr, "  %d) %s / %s  %s\n", i+1, m.Repo.Name, branch, m.Worktree.Path)
	}
	fmt.Fprintln(stderr)

	scanner := bufio.NewScanner(stdin)
	for attempt := 0; attempt < maxSelectionAttempts; attempt++ {
		fmt.Fprintf(stderr, "Select worktree [1-%d]: ", len(matches))

		if !scanner.Scan() {
			return fmt.Errorf("failed to read selection")
		}

		input := strings.TrimSpace(scanner.Text())
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(matches) {
			fmt.Fprintf(stderr, "Invalid selection '%s'. Please enter a number between 1 and %d.\n", input, len(matches))
			continue
		}

		return printGlobalWorktreeMatch(matches[num-1], quiet, stderr, stdout)
	}

	return fmt.Errorf("too many invalid selection attempts")
}
