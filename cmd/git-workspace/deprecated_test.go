package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDeprecatedFlagsAreHidden(t *testing.T) {
	// All deprecated flags should be hidden from help
	deprecatedFlags := []string{
		"list", "init", "add", "recursive", "refresh", "print-workspace", "go",
		"type", "tag", "name", "path", "output", "status", "show-user",
	}

	for _, name := range deprecatedFlags {
		t.Run(name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(name)
			if flag == nil {
				t.Fatalf("Flag --%s not found", name)
			}
			if !flag.Hidden {
				t.Errorf("Flag --%s should be hidden", name)
			}
		})
	}
}

func TestDeprecatedListEmitsWarning(t *testing.T) {
	origList := depList
	defer func() { depList = origList }()

	depList = true
	rootCmd.Flags().Lookup("list").Changed = true
	defer func() { rootCmd.Flags().Lookup("list").Changed = false }()

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	emitDeprecationWarnings(rootCmd)

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Warning: --list is deprecated") {
		t.Errorf("Expected deprecation warning for --list, got: %s", output)
	}
	if !strings.Contains(output, "gws list") {
		t.Errorf("Expected suggestion to use 'gws list', got: %s", output)
	}
}

func TestDeprecatedMutualExclusivity(t *testing.T) {
	origList := depList
	origInit := depInit
	defer func() {
		depList = origList
		depInit = origInit
	}()

	depList = true
	depInit = true

	handled, err := handleDeprecatedFlags(rootCmd, []string{})
	if !handled {
		t.Fatal("Expected handleDeprecatedFlags to handle the error")
	}
	if err == nil {
		t.Fatal("Expected error for mutual exclusivity")
	}
	if !strings.Contains(err.Error(), "only one command can be used at a time") {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDeprecatedRecursiveRequiresAdd(t *testing.T) {
	origRecursive := depRecursive
	origAdd := depAdd
	defer func() {
		depRecursive = origRecursive
		depAdd = origAdd
	}()

	depRecursive = true
	depAdd = ""

	handled, err := handleDeprecatedFlags(rootCmd, []string{})
	if !handled {
		t.Fatal("Expected handleDeprecatedFlags to handle the error")
	}
	if err == nil || !strings.Contains(err.Error(), "--recursive/-v requires --add/-a") {
		t.Errorf("Expected recursive requires add error, got: %v", err)
	}
}

func TestDeprecatedGoConflictsWithArgs(t *testing.T) {
	origGo := depGo
	defer func() { depGo = origGo }()

	depGo = "my-repo"

	handled, err := handleDeprecatedFlags(rootCmd, []string{"other-repo"})
	if !handled {
		t.Fatal("Expected handleDeprecatedFlags to handle the error")
	}
	if err == nil || !strings.Contains(err.Error(), "cannot use both --go flag and positional argument") {
		t.Errorf("Expected go/args conflict error, got: %v", err)
	}
}

func TestDeprecatedNoFlagsSet(t *testing.T) {
	// When no deprecated flags are set, handleDeprecatedFlags should not handle
	origList := depList
	origInit := depInit
	origAdd := depAdd
	origRefresh := depRefresh
	origPW := depPrintWorkspace
	origGo := depGo
	defer func() {
		depList = origList
		depInit = origInit
		depAdd = origAdd
		depRefresh = origRefresh
		depPrintWorkspace = origPW
		depGo = origGo
	}()

	depList = false
	depInit = false
	depAdd = ""
	depRefresh = false
	depPrintWorkspace = false
	depGo = ""

	handled, _ := handleDeprecatedFlags(rootCmd, []string{})
	if handled {
		t.Error("Expected handleDeprecatedFlags to NOT handle when no deprecated flags are set")
	}
}

func TestDepWarningsMap(t *testing.T) {
	// Verify all expected deprecation mappings exist
	expected := map[string]string{
		"list":            "gws list",
		"init":            "gws init",
		"add":             "gws add [path]",
		"recursive":       "gws add --recursive",
		"refresh":         "gws refresh",
		"print-workspace": "gws print-workspace",
		"go":              "gws <repo-name>",
	}

	for flag, newForm := range expected {
		t.Run(flag, func(t *testing.T) {
			if actual, ok := depWarnings[flag]; !ok {
				t.Errorf("Missing deprecation mapping for --%s", flag)
			} else if actual != newForm {
				t.Errorf("Deprecation for --%s: expected '%s', got '%s'", flag, newForm, actual)
			}
		})
	}
}
