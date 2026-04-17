package main

import "github.com/spf13/cobra"

// worktreeCmd is the parent command for worktree management subcommands.
var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees",
	Long: `Manage git worktrees across your workspace.

Subcommands:
  gws worktree list [repo]          # List all worktrees (optionally filtered by repo)
  gws worktree add <repo> <branch>  # Create a new worktree in <repo>.wt/<branch>
  gws worktree align [repo]         # Move unaligned worktrees into <repo>.wt/`,
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
}
