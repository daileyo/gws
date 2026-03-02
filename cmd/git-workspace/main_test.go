package main

import (
	"strings"
	"testing"
)

func TestVersionVariablesAreDefined(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

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
	origVersion := version
	defer func() { version = origVersion }()

	version = "v1.0.0-test"
	rootCmd.Version = version

	if rootCmd.Version != "v1.0.0-test" {
		t.Errorf("Expected rootCmd.Version to be 'v1.0.0-test', got '%s'", rootCmd.Version)
	}
}

func TestCommandFlagsRegistered(t *testing.T) {
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

func TestTagFlagOnRoot(t *testing.T) {
	// --tag (deprecated filter, no shorthand — shorthand moved to tag-cmd alias)
	flag := rootCmd.Flags().Lookup("tag")
	if flag == nil {
		t.Fatal("--tag flag not found on root command")
	}
	if flag.Shorthand != "" {
		t.Errorf("--tag shorthand: expected '' (removed), got '%s'", flag.Shorthand)
	}
}

func TestTagCmdAliasOnRoot(t *testing.T) {
	// -t on root is the alias for the tag subcommand
	flag := rootCmd.Flags().Lookup("tag-cmd")
	if flag == nil {
		t.Fatal("--tag-cmd flag not found on root command")
	}
	if flag.Shorthand != "t" {
		t.Errorf("--tag-cmd shorthand: expected 't', got '%s'", flag.Shorthand)
	}
	if !flag.Hidden {
		t.Error("--tag-cmd should be hidden")
	}
}

func TestSubcommandsRegistered(t *testing.T) {
	expectedSubcommands := []string{"list", "init", "add", "refresh", "print-workspace", "tag", "user"}

	for _, name := range expectedSubcommands {
		t.Run(name, func(t *testing.T) {
			found := false
			for _, cmd := range rootCmd.Commands() {
				if cmd.Name() == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Subcommand '%s' should be registered on rootCmd", name)
			}
		})
	}
}

func TestAliasFlagsAreHidden(t *testing.T) {
	hiddenFlags := []string{"list", "init", "add", "recursive", "refresh", "print-workspace", "go",
		"type", "tag", "name", "path", "output", "status", "show-user", "tag-cmd"}

	for _, name := range hiddenFlags {
		t.Run(name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Deprecated flag --%s not found on root command", name)
				return
			}
			if !flag.Hidden {
				t.Errorf("Deprecated flag --%s should be hidden", name)
			}
		})
	}
}

func TestMutualExclusivity(t *testing.T) {
	origList := depList
	origInit := depInit
	defer func() {
		depList = origList
		depInit = origInit
	}()

	depList = true
	depInit = true

	err := rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("Expected error when multiple command flags are set")
		return
	}

	expected := "only one command can be used at a time"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing '%s', got: %s", expected, err.Error())
	}
}

func TestTagRequiresListOrUser(t *testing.T) {
	origList := depList
	origUser := flagUser
	defer func() {
		depList = origList
		flagUser = origUser
		rootCmd.Flags().Lookup("tag").Changed = false
	}()

	depList = false
	flagUser = false
	rootCmd.Flags().Lookup("tag").Changed = true

	err := rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("Expected error when --tag used without --list or --user")
		return
	}
	if !strings.Contains(err.Error(), "--tag requires") {
		t.Errorf("Expected error about --tag, got: %s", err.Error())
	}
}

func TestUserFlagValidation(t *testing.T) {
	resetUserFlags := func() {
		flagUser = false
		flagUpdate = false
		flagDelete = false
		flagAll = false
		flagVerbose = false
		flagInlineName = ""
		flagInlineEmail = ""
		depList = false
		depInit = false
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

	t.Run("list-users is mutually exclusive with list", func(t *testing.T) {
		defer resetUserFlags()
		flagListUsers = true
		depList = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if !strings.Contains(err.Error(), "only one command can be used at a time") {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		flagListUsers = false
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
		depList = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		expected := "only one command can be used at a time"
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
