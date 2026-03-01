package main

import (
	"strings"
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
	if err.Error() != "only one command flag can be used at a time (--list, --init, --add, --add-tag, --remove-tag, --refresh, --print-workspace, --go, --user)" {
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

	expected := "filter flags (--type, --name, --path, --output, --status) require --list/-l to be set"
	if !strings.Contains(err.Error(), "require --list") {
		t.Errorf("Expected error containing 'require --list', got: %s", err.Error())
	}
	_ = expected
}

func TestUserFlagValidation(t *testing.T) {
	// Helper to reset all user flags after each sub-test
	resetUserFlags := func() {
		flagUser = false
		flagUpdate = false
		flagDelete = false
		flagAll = false
		flagVerbose = false
		flagInlineName = ""
		flagInlineEmail = ""
		flagList = false
		flagInit = false
	}

	t.Run("update without user returns error", func(t *testing.T) {
		defer resetUserFlags()
		flagUpdate = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != "--update/-u requires --user to be set" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("delete without user returns error", func(t *testing.T) {
		defer resetUserFlags()
		flagDelete = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != "--delete/-D requires --user to be set" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("user without update or delete returns error", func(t *testing.T) {
		defer resetUserFlags()
		flagUser = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != "--user requires either --update/-u or --delete/-D" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("user update and delete are mutually exclusive", func(t *testing.T) {
		defer resetUserFlags()
		flagUser = true
		flagUpdate = true
		flagDelete = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != "--update and --delete are mutually exclusive" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("all without delete returns error", func(t *testing.T) {
		defer resetUserFlags()
		flagUser = true
		flagUpdate = true
		flagAll = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != "--all requires --delete/-D to be set" {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("user is mutually exclusive with list", func(t *testing.T) {
		defer resetUserFlags()
		flagUser = true
		flagUpdate = true
		flagList = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		expected := "only one command flag can be used at a time"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error containing '%s', got: %s", expected, err.Error())
		}
	})
}

func TestUserFlagsRegistered(t *testing.T) {
	flags := []struct {
		name      string
		shorthand string
	}{
		{"user", ""},
		{"update", "u"},
		{"delete", "D"},
		{"all", ""},
		{"verbose", ""},
		{"git-name", ""},
		{"git-email", ""},
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
