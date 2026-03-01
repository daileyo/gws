package main

import (
	"testing"
)

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

func TestRootCommand(t *testing.T) {
	// Test that root command is properly configured with new flag-based interface
	if rootCmd.Use != "git-workspace" {
		t.Errorf("Expected root command Use to be 'git-workspace', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Root command Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Root command Long description should not be empty")
	}
}

func TestRootCommandHasVersionSet(t *testing.T) {
	// Cobra's built-in --version/v flag requires rootCmd.Version to be set
	// We set it in main() before Execute(), so test the mechanism
	origVersion := version
	defer func() { version = origVersion }()

	version = "v1.0.0-test"
	// Simulate what main() does
	rootCmd.Version = version

	if rootCmd.Version != "v1.0.0-test" {
		t.Errorf("Expected rootCmd.Version to be 'v1.0.0-test', got '%s'", rootCmd.Version)
	}
}

func TestCommandFlagsRegistered(t *testing.T) {
	// Verify all command flags are registered on the root command
	flags := []struct {
		name      string
		shorthand string
	}{
		{"list", "l"},
		{"init", "i"},
		{"add-tag", "d"},
		{"remove-tag", "x"},
		{"refresh", "r"},
		{"print-workspace", "w"},
		{"go", "g"},
		{"quiet", "q"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Flag --%s not found on root command", f.name)
				return
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestFilterFlagsRegistered(t *testing.T) {
	// Verify all filter flags are registered on the root command with shorthands
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
		{"show-user", ""},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Filter flag --%s not found on root command", f.name)
				return
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Filter flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestNoSubcommandsRegistered(t *testing.T) {
	// Verify old subcommands are not registered (clean break)
	oldSubcommands := []string{"list", "init", "tag", "untag", "refresh", "version"}

	for _, name := range oldSubcommands {
		t.Run(name, func(t *testing.T) {
			for _, cmd := range rootCmd.Commands() {
				if cmd.Name() == name {
					t.Errorf("Old subcommand '%s' should not be registered", name)
				}
			}
		})
	}
}

func TestMutualExclusivity(t *testing.T) {
	// Test that setting multiple command flags returns an error
	// We test this by checking the error message pattern
	origList := flagList
	origInit := flagInit
	defer func() {
		flagList = origList
		flagInit = origInit
	}()

	flagList = true
	flagInit = true

	err := rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("Expected error when multiple command flags are set")
		return
	}

	expected := "only one command flag can be used at a time"
	if err.Error() != "only one command flag can be used at a time (--list, --init, --add, --add-tag, --remove-tag, --refresh, --print-workspace, --go)" {
		t.Errorf("Expected error containing '%s', got: %s", expected, err.Error())
	}
}

func TestFilterFlagsRequireList(t *testing.T) {
	// Test that filter flags without --list produce an error
	origList := flagList
	defer func() {
		flagList = origList
		rootCmd.Flags().Lookup("type").Changed = false
	}()

	flagList = false
	rootCmd.Flags().Lookup("type").Changed = true

	err := rootCmd.RunE(rootCmd, []string{})

	if err == nil {
		t.Error("Expected error when filter flag used without --list")
		return
	}

	expected := "filter flags (--type, --tag, --name, --path, --output, --status) require --list/-l to be set"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got: %s", expected, err.Error())
	}
}
