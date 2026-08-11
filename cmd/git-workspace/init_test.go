package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	gitcfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/daileyo/gws/internal/config"
)

// createInitTestRepo creates a minimal valid git repository for init tests.
func createInitTestRepo(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("createInitTestRepo: failed to create dir %s: %v", path, err)
	}
	repo, err := gogit.PlainInit(path, false)
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to init git repo: %v", err)
	}
	_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to create remote: %v", err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("createInitTestRepo: failed to get worktree: %v", err)
	}
	f := filepath.Join(path, "README.md")
	if err := os.WriteFile(f, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("createInitTestRepo: failed to write file: %v", err)
	}
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("createInitTestRepo: failed to stage file: %v", err)
	}
	if _, err := worktree.Commit("Initial commit", &gogit.CommitOptions{
		Author: &object.Signature{Name: "Test User", Email: "test@example.com"},
	}); err != nil {
		t.Fatalf("createInitTestRepo: failed to commit: %v", err)
	}
}

// withTempHome redirects HOME to a temp directory so config.Save/Load/Exists
// operate on an isolated file instead of the real ~/.gws/config.json.
func withTempHome(t *testing.T) {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
}

// withTempWorkdir changes the process working directory to dir for the duration
// of the test, restoring it in a defer.
func withTempWorkdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("withTempWorkdir: failed to get current dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("withTempWorkdir: failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestRunInit_HappyPath(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir()
	repoDir := filepath.Join(workspaceDir, "my-repo")
	createInitTestRepo(t, repoDir)

	withTempWorkdir(t, workspaceDir)

	// Resolve symlinks so the comparison works on macOS (/var → /private/var)
	resolvedWorkspace, err := filepath.EvalSymlinks(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for workspace dir: %v", err)
	}

	if err := runInit("", io.Discard); err != nil {
		t.Fatalf("runInit returned unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after init: %v", err)
	}
	if cfg.Workspace != resolvedWorkspace {
		t.Errorf("Workspace: expected %s, got %s", resolvedWorkspace, cfg.Workspace)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Repository count: expected 1, got %d", len(cfg.Repositories))
	}
	if len(cfg.Repositories) > 0 && cfg.Repositories[0].Name != "my-repo" {
		t.Errorf("Repository name: expected 'my-repo', got %s", cfg.Repositories[0].Name)
	}
}

func TestRunInit_AlreadyInitialized(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir()
	withTempWorkdir(t, workspaceDir)

	// First init creates the workspace
	if err := runInit("", io.Discard); err != nil {
		t.Fatalf("First runInit failed: %v", err)
	}

	originalCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after first init: %v", err)
	}

	// Second init should return nil (exit 0) without overwriting the config
	if err := runInit("", io.Discard); err != nil {
		t.Errorf("Second runInit should return nil (already-initialized guard), got: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config after second init: %v", err)
	}
	if cfg.Workspace != originalCfg.Workspace {
		t.Errorf("Workspace changed by second init: %s -> %s", originalCfg.Workspace, cfg.Workspace)
	}
}

func TestRunInit_EmptyDirectory(t *testing.T) {
	withTempHome(t)

	workspaceDir := t.TempDir() // no git repos inside
	withTempWorkdir(t, workspaceDir)

	if err := runInit("", io.Discard); err != nil {
		t.Fatalf("runInit on empty directory returned unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != 0 {
		t.Errorf("Expected 0 repositories for empty directory, got %d", len(cfg.Repositories))
	}
}

// --- Spec 21: unified reconciliation ---

// TestRunInit_FullReconciliationSummary asserts the exact summary lines for a
// workspace exercising every discovery path.
func TestRunInit_FullReconciliationSummary(t *testing.T) {
	f := newWorkspaceFixture(t)

	var out bytes.Buffer
	if err := runInit(f.Root, &out); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	got := out.String()

	wantLines := []string{
		"Initialized workspace at: " + f.Root,
		"Found 5 repositories.",
		"Repositories with worktrees: 1",
		"Worktrees: 2 (1 aligned, 1 unaligned)",
	}
	for _, want := range wantLines {
		if !strings.Contains(got, want) {
			t.Errorf("Expected output to contain %q.\nFull output:\n%s", want, got)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != f.ExpectedRepoCount() {
		t.Errorf("Config holds %d repositories, want %d", len(cfg.Repositories), f.ExpectedRepoCount())
	}

	// Every discovery path must be represented in the saved config.
	paths := map[string]bool{}
	for _, repo := range cfg.Repositories {
		paths[repo.Path] = true
	}
	for name, want := range map[string]string{
		"nested repository":   f.Nested,
		"symlinked external":  f.Symlinked,
		"repo with worktrees": f.WithWorktrees,
	} {
		if !paths[want] {
			t.Errorf("Expected %s (%s) in config, got %v", name, want, paths)
		}
	}
}

// TestRunInit_CuratedSymlinkWorkspace covers the headline behavior change: a
// workspace of nothing but symlinks previously yielded zero repositories.
func TestRunInit_CuratedSymlinkWorkspace(t *testing.T) {
	root, repoPaths := newSymlinkOnlyWorkspace(t, "dratlab-apps", "dratlab-kubernetes", "dratlab-tooling")

	var out bytes.Buffer
	if err := runInit(root, &out); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if len(cfg.Repositories) != len(repoPaths) {
		t.Fatalf("Expected %d repositories from a symlink-only workspace, got %d",
			len(repoPaths), len(cfg.Repositories))
	}

	stored := map[string]bool{}
	for _, repo := range cfg.Repositories {
		stored[repo.Path] = true
	}
	for _, want := range repoPaths {
		if !stored[want] {
			t.Errorf("Expected repository %s to be tracked by its real path, got %v", want, stored)
		}
	}
	if !strings.Contains(out.String(), "Found 3 repositories.") {
		t.Errorf("Expected 'Found 3 repositories.' in output:\n%s", out.String())
	}
}

// TestRunInit_WorktreesVisibleImmediately proves init alone produces a complete
// model: no refresh needed before worktree list works.
func TestRunInit_WorktreesVisibleImmediately(t *testing.T) {
	f := newWorkspaceFixture(t)

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	var listOut bytes.Buffer
	if err := runWorktreeList("", &listOut); err != nil {
		t.Fatalf("runWorktreeList failed: %v", err)
	}

	got := listOut.String()
	for _, want := range []string{"feature", "hotfix", "aligned", "(unaligned)"} {
		if !strings.Contains(got, want) {
			t.Errorf("Expected worktree list to contain %q.\nFull output:\n%s", want, got)
		}
	}
}

func TestRunInit_OmitsZeroCountLines(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	fixtureRepo(t, filepath.Join(root, "plain-repo"))

	var out bytes.Buffer
	if err := runInit(root, &out); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Found 1 repository.") {
		t.Errorf("Expected singular 'Found 1 repository.', got:\n%s", got)
	}
	for _, unwanted := range []string{"Repositories with worktrees:", "Worktrees:"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("Expected %q to be omitted when the count is zero, got:\n%s", unwanted, got)
		}
	}
}

// TestRunInit_GuardMessageRecommendsRefresh covers the reworded protective
// guard: stderr, exit zero, nothing modified.
func TestRunInit_GuardMessageRecommendsRefresh(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	fixtureRepo(t, filepath.Join(root, "repo-one"))

	if err := runInit(root, io.Discard); err != nil {
		t.Fatalf("First runInit failed: %v", err)
	}
	before, err := os.ReadFile(configPathForTest(t))
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	stderr, err := captureStderr(t, func() error {
		return runInit(root, io.Discard)
	})
	if err != nil {
		t.Errorf("Second runInit should return nil, got: %v", err)
	}

	if !strings.Contains(stderr, "Workspace already initialized at:") {
		t.Errorf("Expected the guard message on stderr, got:\n%s", stderr)
	}
	if !strings.Contains(stderr, "gws refresh") {
		t.Errorf("Expected the guard message to recommend 'gws refresh', got:\n%s", stderr)
	}

	// refresh must be presented ahead of add as the recommended next step.
	refreshIdx := strings.Index(stderr, "gws refresh")
	addIdx := strings.Index(stderr, "gws add")
	if refreshIdx == -1 || addIdx == -1 || refreshIdx > addIdx {
		t.Errorf("Expected 'gws refresh' to appear before 'gws add', got:\n%s", stderr)
	}

	after, err := os.ReadFile(configPathForTest(t))
	if err != nil {
		t.Fatalf("Failed to re-read config: %v", err)
	}
	if string(before) != string(after) {
		t.Error("The already-initialized guard must not modify the configuration")
	}
}

// TestRunInit_CreatesNoSymlinks pins round-1 Q8-B: symlink maintenance belongs
// to refresh, never to init.
func TestRunInit_CreatesNoSymlinks(t *testing.T) {
	f := newWorkspaceFixture(t)

	before := symlinkEntries(workspaceListing(t, f.Root))

	if err := runInit(f.Root, io.Discard); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	after := symlinkEntries(workspaceListing(t, f.Root))

	if len(before) != len(after) {
		t.Errorf("init changed the workspace symlinks.\nbefore: %v\nafter:  %v", before, after)
	}
	for i := range before {
		if before[i] != after[i] {
			t.Errorf("Symlink %d changed: %q -> %q", i, before[i], after[i])
		}
	}
}

func TestRunInit_ReportsScanErrorsOnStderr(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := resolvedTempDir(t)
	fixtureRepo(t, filepath.Join(root, "good-repo"))
	if err := os.Symlink(filepath.Join(root, "nowhere"), filepath.Join(root, "dangling")); err != nil {
		t.Fatalf("Failed to create dangling symlink: %v", err)
	}

	var out bytes.Buffer
	stderr, err := captureStderr(t, func() error {
		return runInit(root, &out)
	})
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	if !strings.Contains(stderr, "errors occurred during scanning:") {
		t.Errorf("Expected scan warnings on stderr, got:\n%s", stderr)
	}
	if strings.Contains(out.String(), "errors occurred during scanning:") {
		t.Errorf("Scan warnings must not appear on stdout, got:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Found 1 repository.") {
		t.Errorf("Expected the scan to continue, got:\n%s", out.String())
	}
}

// captureStderr redirects os.Stderr for the duration of fn.
func captureStderr(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("Failed to create pipe: %v", pipeErr)
	}
	orig := os.Stderr
	os.Stderr = w

	fnErr := fn()

	os.Stderr = orig
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("Failed to close pipe: %v", closeErr)
	}

	var buf bytes.Buffer
	if _, copyErr := io.Copy(&buf, r); copyErr != nil {
		t.Fatalf("Failed to read captured stderr: %v", copyErr)
	}
	return buf.String(), fnErr
}

// configPathForTest returns the config path under the test HOME.
func configPathForTest(t *testing.T) string {
	t.Helper()
	path, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("Failed to resolve config path: %v", err)
	}
	return path
}
