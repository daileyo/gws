package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
)

var flagDryRun bool

var worktreeAlignCmd = &cobra.Command{
	Use:   "align [repo]",
	Short: "Move unaligned worktrees into the <repo>.wt/ convention",
	Long: `Move all unaligned worktrees into the standard <repo>.wt/ directory structure
using git worktree move (requires Git 2.17+).

When a repo name argument is provided, only that repo's worktrees are aligned.
Without an argument, all repos are processed.

If two worktrees would produce the same directory name, a -dup-NN suffix is
appended (where NN is 00-99).

Examples:
  gws worktree align                # Align all repos
  gws worktree align my-repo        # Align only my-repo
  gws worktree align --dry-run      # Preview moves without executing`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoFilter := ""
		if len(args) == 1 {
			repoFilter = args[0]
		}
		return runWorktreeAlign(repoFilter, flagDryRun)
	},
}

func init() {
	worktreeAlignCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Preview moves without executing them")
	worktreeCmd.AddCommand(worktreeAlignCmd)
}

// alignPlan describes a single worktree move operation.
type alignPlan struct {
	RepoName string
	RepoPath string
	Branch   string
	From     string
	To       string
	Renamed  bool // true if a -dup-NN suffix was applied
}

func runWorktreeAlign(repoFilter string, dryRun bool) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var plans []alignPlan

	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]
		if repoFilter != "" && !filter.MatchesPattern(repo.Name, repoFilter) {
			continue
		}

		// Prune stale worktree entries before planning moves
		_ = git.PruneWorktrees(repo.Path)

		wtDir := repo.Path + ".wt"
		// Track names used in this repo's .wt dir to detect conflicts
		usedNames := make(map[string]bool)

		// Pre-populate with existing aligned worktrees
		for _, wt := range repo.Worktrees {
			if wt.Aligned {
				relPath, err := filepath.Rel(wtDir, wt.Path)
				if err == nil {
					usedNames[relPath] = true
				}
			}
		}

		for _, wt := range repo.Worktrees {
			if wt.Aligned {
				continue
			}

			targetName := wt.Branch
			targetPath := filepath.Join(wtDir, targetName)
			renamed := false

			// Handle naming conflicts
			if usedNames[targetName] {
				renamed = true
				for dup := 0; dup < 100; dup++ {
					candidate := fmt.Sprintf("%s-dup-%02d", targetName, dup)
					if !usedNames[candidate] {
						targetName = candidate
						targetPath = filepath.Join(wtDir, targetName)
						break
					}
				}
			}
			usedNames[targetName] = true

			plans = append(plans, alignPlan{
				RepoName: repo.Name,
				RepoPath: repo.Path,
				Branch:   wt.Branch,
				From:     wt.Path,
				To:       targetPath,
				Renamed:  renamed,
			})
		}
	}

	if len(plans) == 0 {
		fmt.Println("All worktrees are already aligned")
		return nil
	}

	// Display plan
	verb := "Moving"
	if dryRun {
		verb = "Would move"
		fmt.Println("Dry run — no changes will be made:")
		fmt.Println()
	}

	for _, p := range plans {
		suffix := ""
		if p.Renamed {
			suffix = "  (renamed to avoid conflict)"
		}
		fmt.Printf("%s [%s] %s\n  from: %s\n  to:   %s%s\n\n", verb, p.RepoName, p.Branch, p.From, p.To, suffix)
	}

	if dryRun {
		fmt.Printf("Total: %d %s to align\n", len(plans), pluralize(len(plans), "worktree", "worktrees"))
		return nil
	}

	// Execute moves
	var errors []string
	moved := 0
	for _, p := range plans {
		// Create .wt directory if needed
		wtDir := p.RepoPath + ".wt"
		if err := os.MkdirAll(wtDir, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("  %s/%s: failed to create .wt dir: %v", p.RepoName, p.Branch, err))
			continue
		}

		// Ensure parent directories exist for branches with slashes
		parentDir := filepath.Dir(p.To)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("  %s/%s: failed to create parent dir: %v", p.RepoName, p.Branch, err))
			continue
		}

		if err := git.MoveWorktree(p.RepoPath, p.From, p.To); err != nil {
			errors = append(errors, fmt.Sprintf("  %s/%s: %v", p.RepoName, p.Branch, err))
			continue
		}
		moved++
	}

	// Re-discover worktrees for affected repos and save
	affectedRepos := make(map[string]bool)
	for _, p := range plans {
		affectedRepos[p.RepoPath] = true
	}
	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]
		if !affectedRepos[repo.Path] {
			continue
		}
		// Prune any stale entries left by failed moves
		_ = git.PruneWorktrees(repo.Path)
		entries, err := git.ListWorktrees(repo.Path)
		if err != nil {
			continue
		}
		wts := make([]config.Worktree, len(entries))
		for j, e := range entries {
			wts[j] = config.Worktree{
				Path:    e.Path,
				Branch:  e.Branch,
				Aligned: git.IsAligned(e.Path, repo.Path),
			}
		}
		repo.Worktrees = wts
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Aligned %d %s\n", moved, pluralize(moved, "worktree", "worktrees"))
	if len(errors) > 0 {
		fmt.Printf("\n%d %s:\n%s\n", len(errors), pluralize(len(errors), "error", "errors"), strings.Join(errors, "\n"))
	}

	return nil
}
