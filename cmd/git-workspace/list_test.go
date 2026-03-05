package main

import (
	"testing"
)

func TestFilterFlagsOnListCmd(t *testing.T) {
	// Verify all filter flags are registered on listCmd with correct shorthands
	flags := []struct {
		name      string
		shorthand string
	}{
		{"type", "y"},
		{"tag", "t"},
		{"name", "n"},
		{"path", "p"},
		{"output", "o"},
		{"status", "s"},
		{"show-user", "u"},
		{"remote", "r"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Filter flag --%s not found on listCmd", f.name)
				return
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Filter flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
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

func TestListCmdFlagStacking(t *testing.T) {
	// Verify POSIX flag stacking: -su should set both showStatus and showUser
	origStatus := showStatus
	origUser := showUser
	defer func() {
		showStatus = origStatus
		showUser = origUser
	}()

	// Reset
	showStatus = false
	showUser = false

	// Parse stacked flags
	err := listCmd.Flags().Parse([]string{"-su"})
	if err != nil {
		t.Fatalf("Failed to parse -su: %v", err)
	}

	if !showStatus {
		t.Error("Expected showStatus to be true after -su")
	}
	if !showUser {
		t.Error("Expected showUser to be true after -su")
	}
}

func TestListCmdFlagStackingWithRemote(t *testing.T) {
	// Verify POSIX flag stacking: -rsu should set showRemote, showStatus, and showUser
	origStatus := showStatus
	origUser := showUser
	origRemote := showRemote
	defer func() {
		showStatus = origStatus
		showUser = origUser
		showRemote = origRemote
	}()

	// Reset
	showStatus = false
	showUser = false
	showRemote = false

	// Parse stacked flags
	err := listCmd.Flags().Parse([]string{"-rsu"})
	if err != nil {
		t.Fatalf("Failed to parse -rsu: %v", err)
	}

	if !showRemote {
		t.Error("Expected showRemote to be true after -rsu")
	}
	if !showStatus {
		t.Error("Expected showStatus to be true after -rsu")
	}
	if !showUser {
		t.Error("Expected showUser to be true after -rsu")
	}
}
