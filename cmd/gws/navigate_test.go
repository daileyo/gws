package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

var testRepos = []config.Repository{
	{Name: "my-api", Path: "/home/user/projects/my-api", Type: config.TypeGitHub},
	{Name: "my-api-v2", Path: "/home/user/projects/my-api-v2", Type: config.TypeGitHub},
	{Name: "work-api", Path: "/home/user/projects/work-api", Type: config.TypeGitLab},
	{Name: "frontend", Path: "/home/user/projects/frontend", Type: config.TypeGitHub},
	{Name: "infra-tools", Path: "/home/user/projects/infra-tools", Type: config.TypeBitbucket},
}

func TestRunNavigate_SingleExactMatch_Verbose(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("frontend", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "frontend (github)") {
		t.Errorf("stderr should contain verbose info, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "→") {
		t.Errorf("stderr should contain arrow separator, got '%s'", gotStderr)
	}
}

func TestRunNavigate_SingleExactMatch_Quiet(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("frontend", true, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}

	if stderr.Len() != 0 {
		t.Errorf("stderr should be empty in quiet mode, got '%s'", stderr.String())
	}
}

func TestRunNavigate_SinglePartialMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("front", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}
}

func TestRunNavigate_CaseInsensitiveMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("FRONTEND", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}
}

func TestRunNavigate_NoMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("nonexistent", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for no match")
	}

	if !strings.Contains(err.Error(), "no repositories found matching") {
		t.Errorf("error should mention no match, got: %v", err)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "No repositories found matching 'nonexistent'") {
		t.Errorf("stderr should contain no-match message, got '%s'", gotStderr)
	}
}

func TestRunNavigate_UnknownType(t *testing.T) {
	repos := []config.Repository{
		{Name: "local-repo", Path: "/home/user/local-repo", Type: ""},
	}
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("local-repo", false, repos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "unknown") {
		t.Errorf("stderr should show 'unknown' for empty type, got '%s'", gotStderr)
	}
}

func TestNavigateFlagsRegistered(t *testing.T) {
	flags := []struct {
		name      string
		shorthand string
	}{
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
