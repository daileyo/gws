package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDirnameWriter_SinglePath(t *testing.T) {
	var buf bytes.Buffer
	w := &dirnameWriter{&buf}

	_, err := w.Write([]byte("/home/user/projects/my-repo\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if got != "/home/user/projects" {
		t.Errorf("expected '/home/user/projects', got '%s'", got)
	}
}

func TestDirnameWriter_Empty(t *testing.T) {
	var buf bytes.Buffer
	w := &dirnameWriter{&buf}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty path, got '%s'", buf.String())
	}
}

func TestDirnameWriter_ReturnsInputLen(t *testing.T) {
	var buf bytes.Buffer
	w := &dirnameWriter{&buf}

	input := []byte("/some/path\n")
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write should return input length %d, got %d", len(input), n)
	}
}

func TestRunParent_SingleMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("frontend", true, testRepos, &stderr, &dirnameWriter{&stdout}, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(stdout.String())
	if got != "/home/user/projects" {
		t.Errorf("expected '/home/user/projects', got '%s'", got)
	}
}

func TestRunParent_VerboseOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("frontend", false, testRepos, &stderr, &dirnameWriter{&stdout}, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// stdout should have the parent path
	got := strings.TrimSpace(stdout.String())
	if got != "/home/user/projects" {
		t.Errorf("stdout: expected '/home/user/projects', got '%s'", got)
	}

	// stderr should still have the repo info line (unaffected by dirnameWriter)
	if !strings.Contains(stderr.String(), "frontend") {
		t.Errorf("stderr should contain repo info, got '%s'", stderr.String())
	}
}

func TestRunParent_NoMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	err := runNavigate("nonexistent", true, testRepos, &stderr, &dirnameWriter{&stdout}, stdin)
	if err == nil {
		t.Fatal("expected error for no match")
	}
}

func TestRunParent_MultipleMatches_Selection(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("1\n")

	repos := testRepos
	err := runNavigate("my-api", true, repos, &stderr, &dirnameWriter{&stdout}, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(stdout.String())
	if got != "/home/user/projects" {
		t.Errorf("expected '/home/user/projects', got '%s'", got)
	}
}

func TestParentCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "parent" {
			found = true
			break
		}
	}
	if !found {
		t.Error("parent command not registered on root")
	}
}

func TestParentCmdHasQuietFlag(t *testing.T) {
	flag := parentCmd.Flags().Lookup("quiet")
	if flag == nil {
		t.Error("parent command missing --quiet flag")
		return
	}
	if flag.Shorthand != "q" {
		t.Errorf("quiet flag shorthand: expected 'q', got '%s'", flag.Shorthand)
	}
}
