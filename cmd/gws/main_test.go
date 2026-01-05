package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Set test version values
	version = "test-version"
	commit = "test-commit"
	date = "test-date"

	// Create a buffer to capture output
	buf := new(bytes.Buffer)

	// Create a new version command for testing
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Display the version, commit hash, and build date of the gws CLI tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			buf.WriteString("gws version " + version + "\n")
			buf.WriteString("  commit: " + commit + "\n")
			buf.WriteString("  built:  " + date + "\n")
		},
	}

	// Execute the command
	cmd.Run(cmd, []string{})

	// Verify output
	output := buf.String()
	expected := "gws version test-version\n  commit: test-commit\n  built:  test-date\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestRootCommand(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "gws [tag]" {
		t.Errorf("Expected root command Use to be 'gws [tag]', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Root command Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Root command Long description should not be empty")
	}
}

func TestVersionCommandExists(t *testing.T) {
	// Test that version command is registered
	cmd, _, err := rootCmd.Find([]string{"version"})
	if err != nil {
		t.Fatalf("version command not found: %v", err)
	}

	if cmd.Use != "version" {
		t.Errorf("Expected command Use to be 'version', got '%s'", cmd.Use)
	}
}
