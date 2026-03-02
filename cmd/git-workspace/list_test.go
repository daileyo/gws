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

func TestFilterFlagsNotOnRoot(t *testing.T) {
	// Filter flags (except --tag) should NOT be on rootCmd
	rootOnlyAbsent := []string{"type", "name", "path", "output", "status", "show-user"}

	for _, name := range rootOnlyAbsent {
		t.Run(name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(name)
			if flag != nil {
				t.Errorf("Filter flag --%s should NOT be on rootCmd (scoped to listCmd)", name)
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
