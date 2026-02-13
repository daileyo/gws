package main

import (
	"bytes"
	"fmt"
	"io"
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

// withTTY overrides isTerminalFunc to return true for the duration of a test
func withTTY(t *testing.T) {
	t.Helper()
	orig := isTerminalFunc
	isTerminalFunc = func(_ io.Reader) bool { return true }
	t.Cleanup(func() { isTerminalFunc = orig })
}

// withNonTTY overrides isTerminalFunc to return false for the duration of a test
func withNonTTY(t *testing.T) {
	t.Helper()
	orig := isTerminalFunc
	isTerminalFunc = func(_ io.Reader) bool { return false }
	t.Cleanup(func() { isTerminalFunc = orig })
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

func TestRunNavigate_MultipleMatches_Selection(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	// Simulate user typing "2\n" to select second match
	stdin := strings.NewReader("2\n")

	// "my-api" partial matches "my-api", "my-api-v2", and "work-api"
	// but let's use a query that matches exactly 2 for clarity
	repos := []config.Repository{
		{Name: "my-api", Path: "/home/user/projects/my-api", Type: config.TypeGitHub},
		{Name: "my-api-v2", Path: "/home/user/projects/my-api-v2", Type: config.TypeGitHub},
	}

	err := runNavigate("my-api", false, repos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/my-api-v2" {
		t.Errorf("stdout: expected '/home/user/projects/my-api-v2', got '%s'", gotStdout)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "Multiple repositories match") {
		t.Errorf("stderr should contain multiple match header, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "1) my-api (github)") {
		t.Errorf("stderr should list first option, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "2) my-api-v2 (github)") {
		t.Errorf("stderr should list second option, got '%s'", gotStderr)
	}
}

func TestRunNavigate_MultipleMatches_SelectionQuiet(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("1\n")

	repos := []config.Repository{
		{Name: "my-api", Path: "/home/user/projects/my-api", Type: config.TypeGitHub},
		{Name: "my-api-v2", Path: "/home/user/projects/my-api-v2", Type: config.TypeGitHub},
	}

	err := runNavigate("my-api", true, repos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/my-api" {
		t.Errorf("stdout: expected '/home/user/projects/my-api', got '%s'", gotStdout)
	}

	// In quiet mode, stderr should still show the list (needed for selection)
	// but NOT the verbose printMatch output
	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "Multiple repositories match") {
		t.Errorf("stderr should contain multiple match header even in quiet mode, got '%s'", gotStderr)
	}
}

func TestRunNavigate_MultipleMatches_InvalidThenValid(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	// First input is invalid, second is valid
	stdin := strings.NewReader("abc\n1\n")

	repos := []config.Repository{
		{Name: "my-api", Path: "/home/user/projects/my-api", Type: config.TypeGitHub},
		{Name: "my-api-v2", Path: "/home/user/projects/my-api-v2", Type: config.TypeGitHub},
	}

	err := runNavigate("my-api", false, repos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/my-api" {
		t.Errorf("stdout: expected '/home/user/projects/my-api', got '%s'", gotStdout)
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "Invalid selection") {
		t.Errorf("stderr should contain invalid selection message, got '%s'", gotStderr)
	}
}

func TestRunNavigate_MultipleMatches_OutOfRange(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	// All 3 attempts are out of range
	stdin := strings.NewReader("99\n0\n-1\n")

	repos := []config.Repository{
		{Name: "my-api", Path: "/home/user/projects/my-api", Type: config.TypeGitHub},
		{Name: "my-api-v2", Path: "/home/user/projects/my-api-v2", Type: config.TypeGitHub},
	}

	err := runNavigate("my-api", false, repos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error after too many invalid attempts")
	}

	if !strings.Contains(err.Error(), "too many invalid selection attempts") {
		t.Errorf("error should mention too many attempts, got: %v", err)
	}
}

func TestRunNavigate_MultipleMatches_NonTTY(t *testing.T) {
	withNonTTY(t)
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "api" matches my-api, my-api-v2, work-api
	err := runNavigate("api", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for non-TTY multiple matches")
	}

	if !strings.Contains(err.Error(), "multiple repositories match") {
		t.Errorf("error should mention multiple matches, got: %v", err)
	}

	// Should print all matching paths to stdout
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 paths in stdout, got %d: %v", len(lines), lines)
	}
}

func TestRunNavigate_WildcardSingleMatch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "front*" should match only "frontend"
	err := runNavigate("front*", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}
}

func TestRunNavigate_WildcardMultipleMatches(t *testing.T) {
	withNonTTY(t)
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "my-api*" should match "my-api" and "my-api-v2"
	err := runNavigate("my-api*", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for non-TTY multiple wildcard matches")
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 paths in stdout, got %d: %v", len(lines), lines)
	}
}

func TestRunNavigate_WildcardQuestionMark(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "?rontend" should match "frontend"
	err := runNavigate("?rontend", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/frontend" {
		t.Errorf("stdout: expected '/home/user/projects/frontend', got '%s'", gotStdout)
	}
}

func TestIsTerminal_NonFile(t *testing.T) {
	// strings.NewReader is not an *os.File
	r := strings.NewReader("test")
	if isTerminal(r) {
		t.Error("strings.NewReader should not be detected as terminal")
	}
}

func TestIsTerminal_BytesBuffer(t *testing.T) {
	var buf bytes.Buffer
	if isTerminal(&buf) {
		t.Error("bytes.Buffer should not be detected as terminal")
	}
}

func TestRunNavigate_WildcardInteractiveSelection(t *testing.T) {
	withTTY(t)
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("1\n")

	// "*-api" should match "my-api", "my-api-v2", "work-api"
	err := runNavigate("*-api", false, testRepos, &stderr, &stdout, stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotStdout := strings.TrimSpace(stdout.String())
	if gotStdout != "/home/user/projects/my-api" {
		t.Errorf("stdout: expected '/home/user/projects/my-api', got '%s'", gotStdout)
	}
}

func TestRunNavigate_NoMatch_WithSuggestions(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "aip" does not partial-match any repo name, but shares "ai" with nothing
	// and shares no substring. Use "frx" — "fr" appears in "frontend"
	err := runNavigate("frx", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for no match")
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "No repositories found matching 'frx'") {
		t.Errorf("stderr should contain no-match message, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "Did you mean:") {
		t.Errorf("stderr should contain suggestions, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "frontend") {
		t.Errorf("stderr should suggest 'frontend', got '%s'", gotStderr)
	}
}

func TestRunNavigate_NoMatch_SuggestsPartialOverlap(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "fron" shares substring with "frontend" but won't match exactly via wildcard
	// Actually "fron" will partial-match "frontend" via MatchesPattern. Need truly no match.
	// Use "frx" — shares "fr" with "frontend"
	err := runNavigate("frx", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for no match")
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "Did you mean:") {
		t.Errorf("stderr should contain suggestions, got '%s'", gotStderr)
	}
	if !strings.Contains(gotStderr, "frontend") {
		t.Errorf("stderr should suggest 'frontend', got '%s'", gotStderr)
	}
}

func TestRunNavigate_NoMatch_NoSuggestions(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	// "zq" has no 2-char substring overlap with any repo name
	err := runNavigate("zq", false, testRepos, &stderr, &stdout, stdin)
	if err == nil {
		t.Fatal("expected error for no match")
	}

	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "No repositories found matching 'zq'") {
		t.Errorf("stderr should contain no-match message, got '%s'", gotStderr)
	}
	if strings.Contains(gotStderr, "Did you mean:") {
		t.Errorf("stderr should NOT contain suggestions for completely unrelated query, got '%s'", gotStderr)
	}
}

func TestFindSuggestions_MaxLimit(t *testing.T) {
	// Create many repos that would all match
	repos := make([]config.Repository, 10)
	for i := range repos {
		repos[i] = config.Repository{Name: fmt.Sprintf("test-repo-%d", i), Path: fmt.Sprintf("/path/%d", i)}
	}

	suggestions := findSuggestions("test", repos, 5)
	if len(suggestions) > 5 {
		t.Errorf("expected at most 5 suggestions, got %d", len(suggestions))
	}
	if len(suggestions) != 5 {
		t.Errorf("expected exactly 5 suggestions (max), got %d", len(suggestions))
	}
}

func TestFindSuggestions_SubstringMatching(t *testing.T) {
	repos := []config.Repository{
		{Name: "my-api"},
		{Name: "frontend"},
		{Name: "backend"},
		{Name: "zzz-unrelated"},
	}

	// "pi" is a substring of "my-api"
	suggestions := findSuggestions("pi", repos, 5)
	found := false
	for _, s := range suggestions {
		if s == "my-api" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'my-api' in suggestions for 'pi', got %v", suggestions)
	}
}

func TestFindSuggestions_ShortQuery(t *testing.T) {
	repos := []config.Repository{
		{Name: "a-repo"},
	}

	// Single char query — no 2-char substrings possible
	suggestions := findSuggestions("x", repos, 5)
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions for single-char query, got %v", suggestions)
	}
}
