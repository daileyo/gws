package main

import (
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
