package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/daileyo/omgitworks/internal/config"
)

var cdQuiet bool

// stdoutIsTerminalFunc reports whether stdout is attached to a terminal.
// It is a variable so tests can force both branches of the shell-integration hint.
var stdoutIsTerminalFunc = func() bool { return isCharDevice(os.Stdout) }

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Navigate to the workspace root",
	Long: `Navigate to the workspace root directory.

Requires shell integration — the binary prints the path and the 'gws' shell
function performs the directory change. Without it, this command only prints
the path. See 'gws shell-init --help' to set that up.

For scripting, use 'gws print-workspace' instead.

Examples:
  gws cd                         # Navigate to the workspace root
  gws cd -q                      # Print only the path`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCd(cdQuiet, os.Stdout, os.Stderr)
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
	cdCmd.Flags().BoolVarP(&cdQuiet, "quiet", "q", false, "Suppress verbose output, print only the path")
}

// runCd resolves the workspace root and prints it to stdout so the shell
// function can cd to it. Human-facing output goes to stderr so that command
// substitution captures only the path.
func runCd(quiet bool, stdout, stderr io.Writer) error {
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check workspace status: %w", err)
	}
	if !exists {
		fmt.Fprintln(stderr, "Error: workspace not initialized")
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "To get started, navigate to your projects directory and run:")
		fmt.Fprintln(stderr, "  gws init")
		return fmt.Errorf("workspace not initialized")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load workspace configuration: %w", err)
	}

	// Don't hand the shell a dead path — report it here instead.
	info, err := os.Stat(cfg.Workspace)
	if err != nil {
		return fmt.Errorf("workspace directory not found: %s\n\nRe-run 'gws init <directory>' to point gws at the right place", cfg.Workspace)
	}
	if !info.IsDir() {
		return fmt.Errorf("workspace path is not a directory: %s", cfg.Workspace)
	}

	if !quiet {
		fmt.Fprintf(stderr, "workspace → %s\n", cfg.Workspace)
	}

	fmt.Fprintln(stdout, cfg.Workspace)

	// A piped stdout means the shell function captured us. A terminal means the
	// user ran the binary directly, where the cd silently cannot happen.
	if !quiet && stdoutIsTerminalFunc() {
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "Note: changing directory requires shell integration; this printed the path only.")
		fmt.Fprintln(stderr, "Add to your shell config:  eval \"$(git-workspace shell-init zsh)\"   # or: bash")
		fmt.Fprintln(stderr, "PowerShell ($PROFILE):     Invoke-Expression (& git-workspace shell-init powershell | Out-String)")
	}

	return nil
}
