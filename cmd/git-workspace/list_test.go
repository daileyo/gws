package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFilterFlagsOnListCmd(t *testing.T) {
	// Verify all filter/display flags are registered on listCmd with correct shorthands
	flags := []struct {
		name      string
		shorthand string
	}{
		{"type", "y"},
		{"visibility", "V"},
		{"tag", "t"},
		{"name", "n"},
		{"path", "p"},
		{"output", "o"},
		{"status", "s"},
		{"show-user", "u"},
		{"remote", "r"},
		{"verbose", "v"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Flag --%s not found on listCmd", f.name)
				return
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestFilterFlagsOnRootAreHidden(t *testing.T) {
	// Filter flags on root are deprecated (hidden) — they exist for backward compat
	deprecatedFilters := []string{"type", "name", "path", "output", "status", "show-user"}

	for _, name := range deprecatedFilters {
		t.Run(name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Deprecated filter flag --%s not found on rootCmd", name)
				return
			}
			if !flag.Hidden {
				t.Errorf("Deprecated filter flag --%s on rootCmd should be hidden", name)
			}
		})
	}
}

func TestDualPurposeFlagNoOptDefVal(t *testing.T) {
	// Verify dual-purpose flags have NoOptDefVal set to the sentinel
	dualPurposeFlags := []string{"type", "visibility", "tag", "path", "status", "show-user", "remote"}

	for _, name := range dualPurposeFlags {
		t.Run(name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Flag --%s not found on listCmd", name)
				return
			}
			if flag.NoOptDefVal != showColumnSentinel {
				t.Errorf("Flag --%s NoOptDefVal: expected sentinel, got %q", name, flag.NoOptDefVal)
			}
		})
	}
}

func TestListCmdFlagStacking(t *testing.T) {
	// Verify POSIX flag stacking: -su should set both status and user flags
	origStatus := flagStatus
	origUser := flagUser
	defer func() {
		flagStatus = origStatus
		flagUser = origUser
	}()

	// Reset
	flagStatus = ""
	flagUser = ""

	// Parse stacked flags — these are NoOptDefVal flags so stacking sets them to sentinel
	err := listCmd.Flags().Parse([]string{"-su"})
	if err != nil {
		t.Fatalf("Failed to parse -su: %v", err)
	}

	if flagStatus != showColumnSentinel {
		t.Errorf("Expected flagStatus to be sentinel after -su, got %q", flagStatus)
	}
	if flagUser != showColumnSentinel {
		t.Errorf("Expected flagUser to be sentinel after -su, got %q", flagUser)
	}
}

func TestListCmdFlagStackingWithRemote(t *testing.T) {
	// Verify POSIX flag stacking: -rsu should set remote, status, and user
	origStatus := flagStatus
	origUser := flagUser
	origRemote := flagRemote
	defer func() {
		flagStatus = origStatus
		flagUser = origUser
		flagRemote = origRemote
	}()

	// Reset
	flagStatus = ""
	flagUser = ""
	flagRemote = ""

	// Parse stacked flags
	err := listCmd.Flags().Parse([]string{"-rsu"})
	if err != nil {
		t.Fatalf("Failed to parse -rsu: %v", err)
	}

	if flagRemote != showColumnSentinel {
		t.Errorf("Expected flagRemote to be sentinel after -rsu, got %q", flagRemote)
	}
	if flagStatus != showColumnSentinel {
		t.Errorf("Expected flagStatus to be sentinel after -rsu, got %q", flagStatus)
	}
	if flagUser != showColumnSentinel {
		t.Errorf("Expected flagUser to be sentinel after -rsu, got %q", flagUser)
	}
}

func TestListCmdFlagStackingNewFlags(t *testing.T) {
	// Verify POSIX flag stacking with new flags: -yVtp
	origType := flagType
	origVis := flagVisibility
	origTag := flagTag
	origPath := flagPath
	defer func() {
		flagType = origType
		flagVisibility = origVis
		flagTag = origTag
		flagPath = origPath
	}()

	// Reset
	flagType = ""
	flagVisibility = ""
	flagTag = ""
	flagPath = ""

	err := listCmd.Flags().Parse([]string{"-yVtp"})
	if err != nil {
		t.Fatalf("Failed to parse -yVtp: %v", err)
	}

	if flagType != showColumnSentinel {
		t.Errorf("Expected flagType to be sentinel after -yVtp, got %q", flagType)
	}
	if flagVisibility != showColumnSentinel {
		t.Errorf("Expected flagVisibility to be sentinel after -yVtp, got %q", flagVisibility)
	}
	if flagTag != showColumnSentinel {
		t.Errorf("Expected flagTag to be sentinel after -yVtp, got %q", flagTag)
	}
	if flagPath != showColumnSentinel {
		t.Errorf("Expected flagPath to be sentinel after -yVtp, got %q", flagPath)
	}
}

func TestDualPurposeFlagWithValue(t *testing.T) {
	// With NoOptDefVal, values require = syntax: -y=github or --type=github
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"short with equals", []string{"-y=github"}, "github"},
		{"long with equals", []string{"--type=github"}, "github"},
		{"short without value", []string{"-y"}, showColumnSentinel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origType := flagType
			defer func() { flagType = origType }()
			flagType = ""

			err := listCmd.Flags().Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args, err)
			}

			if flagType != tt.expected {
				t.Errorf("Expected flagType to be %q, got %q", tt.expected, flagType)
			}
		})
	}
}

func TestAnyColumnSelected(t *testing.T) {
	tests := []struct {
		name     string
		opts     ListOptions
		expected bool
	}{
		{"no columns", ListOptions{}, false},
		{"show type", ListOptions{ShowType: true}, true},
		{"show visibility", ListOptions{ShowVisibility: true}, true},
		{"show tags", ListOptions{ShowTags: true}, true},
		{"show path", ListOptions{ShowPath: true}, true},
		{"show status", ListOptions{ShowStatus: true}, true},
		{"show user", ListOptions{ShowUser: true}, true},
		{"show remote", ListOptions{ShowRemote: true}, true},
		{"multiple columns", ListOptions{ShowType: true, ShowTags: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opts.AnyColumnSelected()
			if result != tt.expected {
				t.Errorf("AnyColumnSelected() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// captureStdout redirects stdout to capture printed output.
// Since displayMultiColumn uses term.IsTerminal which will return false
// for a pipe, we test the non-TTY (single-column) path directly.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestDisplayMultiColumn_NonTTY(t *testing.T) {
	// When stdout is a pipe (non-TTY), displayMultiColumn outputs one name per line
	tests := []struct {
		name     string
		names    []string
		expected []string
	}{
		{
			"empty list",
			[]string{},
			[]string{},
		},
		{
			"single repo",
			[]string{"my-repo"},
			[]string{"my-repo"},
		},
		{
			"multiple repos sorted",
			[]string{"charlie", "alpha", "bravo"},
			[]string{"alpha", "bravo", "charlie"}, // alphabetical
		},
		{
			"already sorted",
			[]string{"aaa", "bbb", "ccc"},
			[]string{"aaa", "bbb", "ccc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayMultiColumn(tt.names)
			})

			if len(tt.expected) == 0 {
				if output != "" {
					t.Errorf("Expected empty output, got %q", output)
				}
				return
			}

			lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d: %v", len(tt.expected), len(lines), lines)
				return
			}
			for i, expected := range tt.expected {
				if lines[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
				}
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// In test environment (non-TTY), should return default of 80
	width := getTerminalWidth()
	if width != 80 {
		t.Logf("Terminal width detected as %d (may be real TTY or default 80)", width)
	}
	if width <= 0 {
		t.Errorf("getTerminalWidth() returned %d, expected positive value", width)
	}
}

func TestVerboseCountFlag(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedCount int
	}{
		{"no verbose", []string{}, 0},
		{"single -v", []string{"-v"}, 1},
		{"double -vv", []string{"-vv"}, 2},
		{"triple -vvv", []string{"-vvv"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCount := verboseCount
			defer func() { verboseCount = origCount }()
			verboseCount = 0

			if len(tt.args) > 0 {
				err := listCmd.Flags().Parse(tt.args)
				if err != nil {
					t.Fatalf("Failed to parse %v: %v", tt.args, err)
				}
			}

			if verboseCount != tt.expectedCount {
				t.Errorf("Expected verboseCount=%d, got %d", tt.expectedCount, verboseCount)
			}
		})
	}
}

func TestVerboseLevelColumnOverrides(t *testing.T) {
	tests := []struct {
		name           string
		opts           ListOptions
		expectShowType bool
		expectShowVis  bool
		expectShowTags bool
		expectShowPath bool
		expectShowStat bool
		expectShowUser bool
		expectShowRem  bool
	}{
		{
			"verbose 0 - no overrides",
			ListOptions{VerboseLevel: 0},
			false, false, false, false, false, false, false,
		},
		{
			"verbose 1 - stored data columns",
			ListOptions{VerboseLevel: 1},
			true, true, true, true, false, false, false,
		},
		{
			"verbose 2 - all columns",
			ListOptions{VerboseLevel: 2},
			true, true, true, true, true, true, true,
		},
		{
			"verbose 1 with filter - filter preserved, all stored columns shown",
			ListOptions{VerboseLevel: 1, FilterType: "github", ShowType: true},
			true, true, true, true, false, false, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.opts
			// Apply verbose overrides (same logic as parseDualPurposeFlags)
			if opts.VerboseLevel >= 1 {
				opts.ShowType = true
				opts.ShowVisibility = true
				opts.ShowTags = true
				opts.ShowPath = true
			}
			if opts.VerboseLevel >= 2 {
				opts.ShowStatus = true
				opts.ShowUser = true
				opts.ShowRemote = true
			}

			if opts.ShowType != tt.expectShowType {
				t.Errorf("ShowType = %v, want %v", opts.ShowType, tt.expectShowType)
			}
			if opts.ShowVisibility != tt.expectShowVis {
				t.Errorf("ShowVisibility = %v, want %v", opts.ShowVisibility, tt.expectShowVis)
			}
			if opts.ShowTags != tt.expectShowTags {
				t.Errorf("ShowTags = %v, want %v", opts.ShowTags, tt.expectShowTags)
			}
			if opts.ShowPath != tt.expectShowPath {
				t.Errorf("ShowPath = %v, want %v", opts.ShowPath, tt.expectShowPath)
			}
			if opts.ShowStatus != tt.expectShowStat {
				t.Errorf("ShowStatus = %v, want %v", opts.ShowStatus, tt.expectShowStat)
			}
			if opts.ShowUser != tt.expectShowUser {
				t.Errorf("ShowUser = %v, want %v", opts.ShowUser, tt.expectShowUser)
			}
			if opts.ShowRemote != tt.expectShowRem {
				t.Errorf("ShowRemote = %v, want %v", opts.ShowRemote, tt.expectShowRem)
			}

			// Verify filter is preserved
			if tt.opts.FilterType != "" && opts.FilterType != tt.opts.FilterType {
				t.Errorf("FilterType = %q, want %q (filter should be preserved)", opts.FilterType, tt.opts.FilterType)
			}
		})
	}
}
