package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

var parentQuiet bool

var parentCmd = &cobra.Command{
	Use:   "parent <repo>",
	Short: "Print the parent directory of a repository",
	Long: `Print the parent directory path of a repository.

Useful for navigating to the directory that contains a repository.

Examples:
  gws parent my-repo             # Print parent directory path
  cd "$(gws parent my-repo)"     # Navigate to parent directory`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRepoThenNone,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load workspace configuration: %w", err)
		}
		return runNavigate(args[0], parentQuiet, cfg.Repositories, os.Stderr, &dirnameWriter{os.Stdout}, os.Stdin)
	},
}

func init() {
	rootCmd.AddCommand(parentCmd)
	parentCmd.Flags().BoolVarP(&parentQuiet, "quiet", "q", false, "Suppress verbose output, print only the path")
}

// dirnameWriter wraps an io.Writer and applies filepath.Dir to each path written.
// Each Write call is expected to contain a single path followed by a newline,
// as produced by runNavigate writing to stdout.
type dirnameWriter struct {
	w io.Writer
}

func (d *dirnameWriter) Write(b []byte) (int, error) {
	path := strings.TrimRight(string(b), "\n")
	if path != "" {
		fmt.Fprintln(d.w, filepath.Dir(path))
	}
	return len(b), nil
}
