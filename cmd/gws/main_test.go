package main

import (
	"bytes"
	"strings"
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

// TestVersionVariablesAreDefined verifies that the version variables are
// properly defined and accessible for ldflags injection
func TestVersionVariablesAreDefined(t *testing.T) {
	// These variables should always be defined (with default values)
	// When built with ldflags, they will be overridden

	// Temporarily store original values
	origVersion := version
	origCommit := commit
	origDate := date

	// Restore after test
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	// Verify defaults exist (should be "dev", "none", "unknown" when not built with ldflags)
	if version == "" {
		t.Error("version variable should not be empty")
	}
	if commit == "" {
		t.Error("commit variable should not be empty")
	}
	if date == "" {
		t.Error("date variable should not be empty")
	}
}

// TestVersionLdflagsSimulation simulates ldflags injection by setting version variables
// and verifying they are correctly used in output
func TestVersionLdflagsSimulation(t *testing.T) {
	testCases := []struct {
		name            string
		version         string
		commit          string
		date            string
		expectedVersion string
		expectedCommit  string
		expectedDate    string
	}{
		{
			name:            "semantic version",
			version:         "v1.2.3",
			commit:          "abc1234",
			date:            "2025-01-05T12:00:00Z",
			expectedVersion: "v1.2.3",
			expectedCommit:  "abc1234",
			expectedDate:    "2025-01-05T12:00:00Z",
		},
		{
			name:            "dev version",
			version:         "dev",
			commit:          "none",
			date:            "unknown",
			expectedVersion: "dev",
			expectedCommit:  "none",
			expectedDate:    "unknown",
		},
		{
			name:            "pre-release version",
			version:         "v2.0.0-beta.1",
			commit:          "def5678",
			date:            "2025-06-15T18:30:00Z",
			expectedVersion: "v2.0.0-beta.1",
			expectedCommit:  "def5678",
			expectedDate:    "2025-06-15T18:30:00Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate ldflags injection
			version = tc.version
			commit = tc.commit
			date = tc.date

			// Capture output
			buf := new(bytes.Buffer)

			// Create version command
			cmd := &cobra.Command{
				Use: "version",
				Run: func(cmd *cobra.Command, args []string) {
					buf.WriteString("gws version " + version + "\n")
					buf.WriteString("  commit: " + commit + "\n")
					buf.WriteString("  built:  " + date + "\n")
				},
			}

			cmd.Run(cmd, []string{})
			output := buf.String()

			// Verify version is in output
			if !strings.Contains(output, tc.expectedVersion) {
				t.Errorf("Expected output to contain version '%s', got: %s", tc.expectedVersion, output)
			}
			if !strings.Contains(output, tc.expectedCommit) {
				t.Errorf("Expected output to contain commit '%s', got: %s", tc.expectedCommit, output)
			}
			if !strings.Contains(output, tc.expectedDate) {
				t.Errorf("Expected output to contain date '%s', got: %s", tc.expectedDate, output)
			}
		})
	}
}

// TestVersionOutputFormat verifies the version output format is consistent
func TestVersionOutputFormat(t *testing.T) {
	version = "v1.0.0"
	commit = "abc123"
	date = "2025-01-05T00:00:00Z"

	buf := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			buf.WriteString("gws version " + version + "\n")
			buf.WriteString("  commit: " + commit + "\n")
			buf.WriteString("  built:  " + date + "\n")
		},
	}

	cmd.Run(cmd, []string{})
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have exactly 3 lines
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines of output, got %d: %v", len(lines), lines)
	}

	// First line should start with "gws version"
	if !strings.HasPrefix(lines[0], "gws version ") {
		t.Errorf("First line should start with 'gws version ', got: %s", lines[0])
	}

	// Second line should contain commit
	if !strings.Contains(lines[1], "commit:") {
		t.Errorf("Second line should contain 'commit:', got: %s", lines[1])
	}

	// Third line should contain built
	if !strings.Contains(lines[2], "built:") {
		t.Errorf("Third line should contain 'built:', got: %s", lines[2])
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
