package main

import (
	"os"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestUserShortFlagsRegistered(t *testing.T) {
	flags := []struct {
		name      string
		shorthand string
	}{
		{"add", "a"},
		{"delete", "d"},
		{"list", "l"},
		{"show", "s"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := userCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Fatalf("Flag --%s not found on userCmd", f.name)
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestUserShortFlagMutualExclusivity(t *testing.T) {
	tests := []struct {
		name string
		setA bool
		setD bool
		setL bool
		setS bool
	}{
		{"add and delete", true, true, false, false},
		{"add and list", true, false, true, false},
		{"add and show", true, false, false, true},
		{"delete and list", false, true, true, false},
		{"delete and show", false, true, false, true},
		{"list and show", false, false, true, true},
		{"all four", true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origA := userFlagAdd
			origD := userFlagDelete
			origL := userFlagList
			origS := userFlagShow
			defer func() {
				userFlagAdd = origA
				userFlagDelete = origD
				userFlagList = origL
				userFlagShow = origS
			}()

			userFlagAdd = tt.setA
			userFlagDelete = tt.setD
			userFlagList = tt.setL
			userFlagShow = tt.setS

			err := userCmd.RunE(userCmd, []string{"some-arg"})
			if err == nil {
				t.Fatal("Expected mutual exclusivity error")
			}
			expected := "only one of -a, -d, -l, -s can be used at a time"
			if err.Error() != expected {
				t.Errorf("Expected error %q, got %q", expected, err.Error())
			}
		})
	}
}

func TestUserSubcommandsRegistered(t *testing.T) {
	expected := []string{"list", "add", "show", "remove", "assign", "sync"}

	for _, name := range expected {
		t.Run(name, func(t *testing.T) {
			found := false
			for _, cmd := range userCmd.Commands() {
				if cmd.Name() == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Subcommand '%s' should be registered on userCmd", name)
			}
		})
	}
}

func TestUserNoOperationShowsHelp(t *testing.T) {
	origA := userFlagAdd
	origD := userFlagDelete
	origL := userFlagList
	origS := userFlagShow
	defer func() {
		userFlagAdd = origA
		userFlagDelete = origD
		userFlagList = origL
		userFlagShow = origS
	}()

	userFlagAdd = false
	userFlagDelete = false
	userFlagList = false
	userFlagShow = false

	// Redirect stdout to avoid printing help in test output
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	err := userCmd.RunE(userCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	// Help returns nil error
	if err != nil {
		t.Errorf("Expected no error when showing help, got: %v", err)
	}
}

func TestUserShortFlagAddRequiresArgs(t *testing.T) {
	origA := userFlagAdd
	defer func() { userFlagAdd = origA }()

	userFlagAdd = true

	err := userCmd.RunE(userCmd, []string{})
	if err == nil {
		t.Fatal("Expected error when -a used without args")
	}
	if err.Error() != "user -a requires a profile name: gws user -a <name> --email <email>" {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestUserShortFlagShowRequiresOneArg(t *testing.T) {
	origS := userFlagShow
	defer func() { userFlagShow = origS }()

	userFlagShow = true

	// No args
	err := userCmd.RunE(userCmd, []string{})
	if err == nil {
		t.Fatal("Expected error when -s used without args")
	}

	// Too many args
	err = userCmd.RunE(userCmd, []string{"a", "b"})
	if err == nil {
		t.Fatal("Expected error when -s used with too many args")
	}
}

func TestUserShortFlagDeleteRequiresOneArg(t *testing.T) {
	origD := userFlagDelete
	defer func() { userFlagDelete = origD }()

	userFlagDelete = true

	// No args
	err := userCmd.RunE(userCmd, []string{})
	if err == nil {
		t.Fatal("Expected error when -d used without args")
	}

	// Too many args
	err = userCmd.RunE(userCmd, []string{"a", "b"})
	if err == nil {
		t.Fatal("Expected error when -d used with too many args")
	}
}

func TestUserAddFlagsOnUserCmd(t *testing.T) {
	// Verify that add-related flags are registered on userCmd
	// so they work with gws user -a <name> --email <email>
	flags := []string{"email", "name", "signing-key", "sign-commits"}
	for _, name := range flags {
		t.Run(name, func(t *testing.T) {
			flag := userCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Flag --%s not found on userCmd", name)
			}
		})
	}
}

func TestUserProfileCompletion(t *testing.T) {
	// Set up a temporary workspace with profiles
	tmpDir := t.TempDir()
	t.Setenv("GWS_CONFIG_DIR", tmpDir)

	cfg := &config.Config{
		Workspace: tmpDir,
		Profiles: []config.Profile{
			{Name: "work", GitName: "Worker", Email: "work@example.com"},
			{Name: "personal", GitName: "Personal", Email: "personal@example.com"},
			{Name: "oss", GitName: "OSS", Email: "oss@example.com"},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test completion with "w" prefix
	names, _ := completeProfileNames("w")
	if len(names) != 1 || names[0] != "work" {
		t.Errorf("Expected [work], got %v", names)
	}

	// Test completion with empty prefix (returns all)
	names, _ = completeProfileNames("")
	if len(names) != 3 {
		t.Errorf("Expected 3 profiles, got %d: %v", len(names), names)
	}

	// Test completion with no match
	names, _ = completeProfileNames("xyz")
	if len(names) != 0 {
		t.Errorf("Expected no matches, got %v", names)
	}
}
