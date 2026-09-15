package main

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/daileyo/omgitworks/internal/config"
	"github.com/daileyo/omgitworks/internal/xdg"
)

// isolateGitHooks stops a developer's global git hooks from running against the
// throwaway repositories these tests create.
//
// Fixtures commit with scaffolding messages like "init". A global
// core.hooksPath — a conventional-commit validator, say — applies to every repo
// on the machine, including fixtures in a temp dir, and rejects those messages.
// Fixture setup then fails before a single assertion runs.
//
// GIT_CONFIG_* is used rather than a repo-local setting because it propagates
// to every git subprocess the tests spawn, wherever that repo lives.
func isolateGitHooks() {
	os.Setenv("GIT_CONFIG_COUNT", "1")
	os.Setenv("GIT_CONFIG_KEY_0", "core.hooksPath")
	os.Setenv("GIT_CONFIG_VALUE_0", "")
}

// isolateXDGDirs clears the XDG variables for the whole package.
//
// Tests isolate themselves by pointing HOME at a temp directory, but XDG_CONFIG_HOME
// and XDG_DATA_HOME take precedence over HOME when resolving paths. On a machine
// where either is set — CI runners set XDG_CONFIG_HOME — that isolation silently
// stops working and tests write to the developer's real config.
//
// Clearing them here makes HOME authoritative, so every existing helper that
// overrides HOME isolates fully again.
func isolateXDGDirs() {
	os.Unsetenv(xdg.EnvConfigHome)
	os.Unsetenv(xdg.EnvDataHome)
}

// TestMain guards against tests corrupting the real user config.
// It snapshots the config file before all tests and fails if
// the content has changed after the suite finishes.
func TestMain(m *testing.M) {
	isolateGitHooks()
	isolateXDGDirs()

	configPath, err := config.GetConfigPath()
	if err != nil {
		// If we can't resolve the path, just run tests without the guard.
		os.Exit(m.Run())
	}

	var before []byte
	existed := false
	if data, err := os.ReadFile(configPath); err == nil {
		before = data
		existed = true
	}

	code := m.Run()

	if existed {
		after, err := os.ReadFile(configPath)
		if err != nil {
			os.Stderr.WriteString("FAIL: real config was deleted by tests: " + configPath + "\n")
			os.Exit(1)
		}
		if !bytes.Equal(before, after) {
			os.Stderr.WriteString("FAIL: real config was modified by tests: " + configPath + "\n")
			os.Exit(1)
		}
	} else {
		if _, err := os.Stat(configPath); err == nil {
			os.Stderr.WriteString("FAIL: real config was created by tests: " + configPath + "\n")
			// Clean up the file that shouldn't exist.
			os.Remove(configPath)
			os.Exit(1)
		}
	}

	os.Exit(code)
}

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
	if rootCmd.Use != "omgitworks" {
		t.Errorf("Expected root command Use to be 'omgitworks', got '%s'", rootCmd.Use)
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

func TestRootVersionFlagOutput(t *testing.T) {
	origVersion := rootCmd.Version
	var out bytes.Buffer
	rootCmd.Version = "v1.2.3\n  commit: abc1234\n  built:  2026-01-01T00:00:00Z"
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"--version"})
	defer func() {
		rootCmd.Version = origVersion
		rootCmd.SetOut(nil)
		rootCmd.SetArgs(nil)
		if f := rootCmd.Flags().Lookup("version"); f != nil {
			_ = f.Value.Set("false")
			f.Changed = false
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("omgitworks --version: %v", err)
	}
	want := "omgitworks version v1.2.3\n  commit: abc1234\n  built:  2026-01-01T00:00:00Z\n"
	if out.String() != want {
		t.Errorf("omgitworks --version output:\ngot:  %q\nwant: %q", out.String(), want)
	}
}

// helpFlagLine matches a flag entry in rendered help and captures its long name.
var helpFlagLine = regexp.MustCompile(`^\s+(?:-\S, )?--([\w-]+)`)

// helpListedFlags renders cmd's help as `<cmd> --help` would and returns the
// long names of the flags listed under its Flags and Global Flags sections.
// It also fails the test if a sentinel NoOptDefVal leaks into the output.
func helpListedFlags(t *testing.T, cmd *cobra.Command) map[string]bool {
	t.Helper()
	// Cobra adds these flags at execution time, just before rendering help.
	cmd.InitDefaultHelpFlag()
	cmd.InitDefaultVersionFlag()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	defer cmd.SetOut(nil)
	cmd.HelpFunc()(cmd, nil)

	if strings.Contains(buf.String(), "\x00") {
		t.Errorf("help for %q leaks a sentinel NoOptDefVal:\n%s", cmd.CommandPath(), buf.String())
	}

	listed := map[string]bool{}
	inFlags := false
	for _, line := range strings.Split(buf.String(), "\n") {
		switch {
		case line == "Flags:" || line == "Global Flags:":
			inFlags = true
		case strings.TrimSpace(line) == "" || !strings.HasPrefix(line, " "):
			inFlags = false
		case inFlags:
			if m := helpFlagLine.FindStringSubmatch(line); m != nil {
				listed[m[1]] = true
			}
		}
	}
	return listed
}

func TestHelpListsExactlyTheFlagsEachCommandAccepts(t *testing.T) {
	origVersion := rootCmd.Version
	rootCmd.Version = "v1.2.3"
	defer func() { rootCmd.Version = origVersion }()

	var walk func(cmd *cobra.Command)
	walk = func(cmd *cobra.Command) {
		t.Run(cmd.CommandPath(), func(t *testing.T) {
			listed := helpListedFlags(t, cmd)

			for name := range listed {
				if cmd.LocalFlags().Lookup(name) == nil && cmd.InheritedFlags().Lookup(name) == nil {
					t.Errorf("help lists --%s, which %q does not accept", name, cmd.CommandPath())
				}
			}

			for _, fs := range []*pflag.FlagSet{cmd.LocalFlags(), cmd.InheritedFlags()} {
				fs.VisitAll(func(f *pflag.Flag) {
					if !f.Hidden && !listed[f.Name] {
						t.Errorf("help omits --%s, which %q accepts", f.Name, cmd.CommandPath())
					}
				})
			}
		})
		for _, sub := range cmd.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)
}

func TestNavigationHelpOnlyOnRoot(t *testing.T) {
	usage := func(cmd *cobra.Command) string {
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		defer cmd.SetOut(nil)
		_ = cmd.Usage()
		return buf.String()
	}

	if !strings.Contains(usage(rootCmd), "Navigation:") {
		t.Error("root help should include the Navigation section")
	}
	if strings.Contains(usage(worktreeAlignCmd), "Navigation:") {
		t.Error("subcommand help should not include the root Navigation section")
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
		"add-tag", "remove-tag",
		"type", "tag", "name", "path", "output", "status", "show-user", "tag-cmd",
		"user", "update", "delete", "all", "verbose", "git-name", "git-email", "list-users"}

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
	origUser := depUser
	defer func() {
		depList = origList
		depUser = origUser
		rootCmd.Flags().Lookup("tag").Changed = false
	}()

	depList = false
	depUser = false
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
		depUser = false
		depUpdate = false
		depDelete = false
		depAll = false
		depVerbose = false
		depInlineName = ""
		depInlineEmail = ""
		depList = false
		depInit = false
	}

	t.Run("update without user returns error", func(t *testing.T) {
		defer resetUserFlags()
		depUpdate = true

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
		depDelete = true

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
		depListUsers = true
		depList = true

		err := rootCmd.RunE(rootCmd, []string{})
		if err == nil {
			t.Fatal("Expected error")
		}
		if !strings.Contains(err.Error(), "only one command can be used at a time") {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		depListUsers = false
	})

	t.Run("user update and delete are mutually exclusive", func(t *testing.T) {
		defer resetUserFlags()
		depUser = true
		depUpdate = true
		depDelete = true

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
		depUser = true
		depUpdate = true
		depAll = true

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
		depUser = true
		depUpdate = true
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
