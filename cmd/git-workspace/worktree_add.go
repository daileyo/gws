package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
	"github.com/daileyo/gws/internal/xdg"
)

var worktreeAddCmd = &cobra.Command{
	Use:   "add <repo> <branch>",
	Short: "Create a new worktree in the projects root",
	Long: `Create a new git worktree for the given branch.

Worktrees live under the XDG projects root, one directory per repository:

  $XDG_DATA_HOME/gws/projects/<repo>/<branch>
  (default: ~/.local/share/gws/projects/<repo>/<branch>)

The directory is created automatically if it does not exist.

Examples:
  gws worktree add my-repo feature-auth
  gws worktree add my-repo feature/new-api`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorktreeAdd(args[0], args[1])
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeAddCmd)
}

func runWorktreeAdd(repoName, branch string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Find matching repo
	var repoIdx = -1
	for i, repo := range cfg.Repositories {
		if filter.MatchesPattern(repo.Name, repoName) {
			if repoIdx >= 0 {
				return fmt.Errorf("multiple repositories match '%s', narrow your query", repoName)
			}
			repoIdx = i
		}
	}
	if repoIdx < 0 {
		return fmt.Errorf("no repository found matching '%s'", repoName)
	}

	repo := &cfg.Repositories[repoIdx]

	// Check for duplicate branch
	for _, wt := range repo.Worktrees {
		if wt.Branch == branch {
			return fmt.Errorf("worktree for branch '%s' already exists at %s", branch, wt.Path)
		}
	}

	// Compute destination path under the XDG projects root
	wtDir, err := xdg.RepoProjectsDir(repo.Name)
	if err != nil {
		return err
	}
	destPath := filepath.Join(wtDir, branch)

	// Create the repo's projects directory if needed. A branch name containing
	// slashes nests, so create the destination's parent rather than wtDir.
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create projects directory: %w", err)
	}

	// Create the worktree
	if err := git.AddWorktree(repo.Path, branch, destPath); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// Re-discover worktrees for this repo and save config
	entries, err := git.ListWorktrees(repo.Path)
	if err == nil {
		wts := make([]config.Worktree, len(entries))
		for j, e := range entries {
			wts[j] = config.Worktree{
				Path:    e.Path,
				Branch:  e.Branch,
				Aligned: git.IsAligned(e.Path, repo.Name),
			}
		}
		repo.Worktrees = wts
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Created worktree for branch '%s' at %s\n", branch, destPath)
	return nil
}
