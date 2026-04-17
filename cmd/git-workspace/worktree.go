package main

import (
	"os"

	"github.com/spf13/cobra"
)

// worktreeCmd is the parent command for worktree management subcommands.
// When called with a repo name argument (e.g., "gws worktree my-repo"),
// it delegates to "worktree list <repo>".
var worktreeCmd = &cobra.Command{
	Use:   "worktree [repo]",
	Short: "Manage git worktrees",
	Long: `Manage git worktrees across your workspace.

When called with a repo name, lists worktrees for that repo:
  gws worktree my-repo              # Same as: gws worktree list my-repo

Subcommands:
  gws worktree list [repo]          # List all worktrees (optionally filtered by repo)
  gws worktree add <repo> <branch>  # Create a new worktree in <repo>.wt/<branch>
  gws worktree align [repo]         # Move unaligned worktrees into <repo>.wt/`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return runWorktreeList(args[0], os.Stdout)
		}
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
}
