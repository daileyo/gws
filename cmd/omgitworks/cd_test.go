package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/omgitworks/internal/config"
)

// withCdWorkspace initializes an isolated workspace config pointing at dir.
func withCdWorkspace(t *testing.T, dir string) {
	t.Helper()
	withTempHome(t)
	cfg := config.New(dir)
	if err := config.Save(cfg); err != nil {
		t.Fatalf("withCdWorkspace: failed to save config: %v", err)
	}
}

// withStdoutTerminal forces the stdout TTY check for the duration of the test.
func withStdoutTerminal(t *testing.T, isTTY bool) {
	t.Helper()
	orig := stdoutIsTerminalFunc
	stdoutIsTerminalFunc = func() bool { return isTTY }
	t.Cleanup(func() { stdoutIsTerminalFunc = orig })
}

func TestRunCd_PrintsPathToStdout(t *testing.T) {
	workspace := t.TempDir()
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, false)

	var stdout, stderr bytes.Buffer
	if err := runCd(false, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := strings.TrimSpace(stdout.String()); got != workspace {
		t.Errorf("stdout = %q, want %q", got, workspace)
	}
	if !strings.Contains(stderr.String(), "workspace → "+workspace) {
		t.Errorf("stderr missing informational line, got %q", stderr.String())
	}
}

func TestRunCd_StdoutHoldsOnlyThePath(t *testing.T) {
	workspace := t.TempDir()
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, true)

	var stdout, stderr bytes.Buffer
	if err := runCd(false, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Even with the hint printed, stdout must stay capturable by $(...).
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("stdout should hold exactly one line, got %d: %q", len(lines), stdout.String())
	}
	if lines[0] != workspace {
		t.Errorf("stdout = %q, want %q", lines[0], workspace)
	}
}

func TestRunCd_QuietSuppressesStderr(t *testing.T) {
	workspace := t.TempDir()
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, true)

	var stdout, stderr bytes.Buffer
	if err := runCd(true, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stderr.Len() != 0 {
		t.Errorf("quiet mode should emit no stderr, got %q", stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != workspace {
		t.Errorf("stdout = %q, want %q", got, workspace)
	}
}

func TestRunCd_HintShownWhenStdoutIsTerminal(t *testing.T) {
	workspace := t.TempDir()
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, true)

	var stdout, stderr bytes.Buffer
	if err := runCd(false, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(stderr.String(), "shell integration") {
		t.Errorf("expected shell-integration hint on stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "shell-init") {
		t.Errorf("hint should name shell-init, got %q", stderr.String())
	}
}

func TestRunCd_NoHintWhenStdoutIsPiped(t *testing.T) {
	workspace := t.TempDir()
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, false)

	var stdout, stderr bytes.Buffer
	if err := runCd(false, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(stderr.String(), "shell integration") {
		t.Errorf("hint should be suppressed when stdout is piped, got %q", stderr.String())
	}
}

func TestRunCd_UninitializedWorkspace(t *testing.T) {
	withTempHome(t) // no config saved
	withStdoutTerminal(t, false)

	var stdout, stderr bytes.Buffer
	err := runCd(false, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when the workspace is not initialized")
	}

	out := stderr.String()
	if !strings.Contains(out, "workspace not initialized") {
		t.Errorf("stderr should report the workspace is not initialized, got %q", out)
	}
	if !strings.Contains(out, "gws init") {
		t.Errorf("stderr should suggest 'gws init', got %q", out)
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should stay empty on error, got %q", stdout.String())
	}
}

func TestRunCd_WorkspaceDirectoryMissing(t *testing.T) {
	workspace := filepath.Join(t.TempDir(), "gone")
	withCdWorkspace(t, workspace)
	withStdoutTerminal(t, false)

	var stdout, stderr bytes.Buffer
	err := runCd(false, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when the workspace directory does not exist")
	}
	if !strings.Contains(err.Error(), workspace) {
		t.Errorf("error should name the missing path, got %q", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should stay empty rather than emit a dead path, got %q", stdout.String())
	}
}

func TestRunCd_WorkspaceIsNotADirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "notadir")
	if err := os.WriteFile(file, []byte("x"), 0600); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	withCdWorkspace(t, file)
	withStdoutTerminal(t, false)

	var stdout, stderr bytes.Buffer
	err := runCd(false, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when the workspace path is not a directory")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("error should say the path is not a directory, got %q", err.Error())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should stay empty, got %q", stdout.String())
	}
}

func TestCdCmd_RejectsArguments(t *testing.T) {
	if err := cdCmd.Args(cdCmd, []string{"somewhere"}); err == nil {
		t.Error("cd should reject positional arguments so 'gws cd foo' is not silently treated as the root")
	}
	if err := cdCmd.Args(cdCmd, []string{}); err != nil {
		t.Errorf("cd should accept zero arguments, got %v", err)
	}
}

func TestIsCharDevice_NonTerminalFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "regular")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer f.Close() //nolint:errcheck // test cleanup

	if isCharDevice(f) {
		t.Error("a regular file should not be reported as a character device")
	}
	if isCharDevice(nil) {
		t.Error("a nil file should not be reported as a character device")
	}
}
