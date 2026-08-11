package discovery

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// gitWorktreeAdd creates a linked worktree for repoDir on a new branch.
func gitWorktreeAdd(t *testing.T, repoDir, branch, destPath string) {
	t.Helper()
	cmd := exec.Command("git", "worktree", "add", "-b", branch, destPath)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add %s failed: %s\n%s", destPath, err, out)
	}
}

// createTestRepo creates a test git repository in the given directory
func createTestRepo(t *testing.T, path string, remoteURL string) {
	t.Helper()

	// Initialize git repository
	repo, err := git.PlainInit(path, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Add remote if specified
	if remoteURL != "" {
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{remoteURL},
		})
		if err != nil {
			t.Fatalf("Failed to create remote: %v", err)
		}
	}

	// Create a dummy commit to make it a valid repo
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(path, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	if _, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	}); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
}

func TestScan_SingleRepository(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	createTestRepo(t, repoDir, "https://github.com/test/repo.git")

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Expected 1 repository, got %d", result.Count)
	}

	if len(result.Repositories) != 1 {
		t.Fatalf("Expected 1 repository in results, got %d", len(result.Repositories))
	}

	repo := result.Repositories[0]
	if repo.Name != "test-repo" {
		t.Errorf("Expected repository name 'test-repo', got '%s'", repo.Name)
	}

	if repo.RemoteURL != "https://github.com/test/repo.git" {
		t.Errorf("Expected remote URL 'https://github.com/test/repo.git', got '%s'", repo.RemoteURL)
	}
}

func TestScan_MultipleRepositories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test repositories
	repos := []struct {
		name      string
		remoteURL string
	}{
		{"repo1", "https://github.com/test/repo1.git"},
		{"repo2", "https://github.com/test/repo2.git"},
		{"repo3", "https://gitlab.com/test/repo3.git"},
	}

	for _, r := range repos {
		repoDir := filepath.Join(tmpDir, r.name)
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			t.Fatalf("Failed to create repo dir: %v", err)
		}
		createTestRepo(t, repoDir, r.remoteURL)
	}

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Count != 3 {
		t.Errorf("Expected 3 repositories, got %d", result.Count)
	}

	if len(result.Repositories) != 3 {
		t.Fatalf("Expected 3 repositories in results, got %d", len(result.Repositories))
	}

	// Verify repository names
	repoNames := make(map[string]bool)
	for _, repo := range result.Repositories {
		repoNames[repo.Name] = true
	}

	for _, r := range repos {
		if !repoNames[r.name] {
			t.Errorf("Repository '%s' not found in scan results", r.name)
		}
	}
}

func TestScan_NestedRepositories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested structure: tmpDir/parent/child
	parentDir := filepath.Join(tmpDir, "parent")
	childDir := filepath.Join(parentDir, "child")

	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	createTestRepo(t, parentDir, "https://github.com/test/parent.git")
	createTestRepo(t, childDir, "https://github.com/test/child.git")

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Strict boundary pruning: traversal stops at the parent repository root,
	// so the nested clone is never registered.
	if result.Count != 1 {
		t.Errorf("Expected 1 repository (nested clone pruned), got %d", result.Count)
	}

	if len(result.Repositories) > 0 && result.Repositories[0].Name != "parent" {
		t.Errorf("Expected to find 'parent', got '%s'", result.Repositories[0].Name)
	}
}

func TestScan_NonExistentDirectory(t *testing.T) {
	result, err := Scan("/nonexistent/directory/path", Options{})
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for nonexistent directory")
	}
}

func TestScan_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Count != 0 {
		t.Errorf("Expected 0 repositories in empty directory, got %d", result.Count)
	}

	if len(result.Repositories) != 0 {
		t.Errorf("Expected empty repositories slice, got %d items", len(result.Repositories))
	}
}

func TestScan_SkipsCommonDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create repos in directories that should be skipped
	skipDirs := []string{"node_modules", "vendor", ".venv", "target"}
	for _, dir := range skipDirs {
		repoDir := filepath.Join(tmpDir, dir, "test-repo")
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		createTestRepo(t, repoDir, "")
	}

	// Create a repo that should be found
	validRepoDir := filepath.Join(tmpDir, "valid-repo")
	if err := os.MkdirAll(validRepoDir, 0755); err != nil {
		t.Fatalf("Failed to create valid repo dir: %v", err)
	}
	createTestRepo(t, validRepoDir, "")

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should only find the valid repo, not the ones in skipped directories
	if result.Count != 1 {
		t.Errorf("Expected 1 repository (skipped dirs should be ignored), got %d", result.Count)
	}

	if len(result.Repositories) > 0 && result.Repositories[0].Name != "valid-repo" {
		t.Errorf("Expected to find 'valid-repo', got '%s'", result.Repositories[0].Name)
	}
}

func TestScan_RepositoryWithoutRemote(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "no-remote-repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	createTestRepo(t, repoDir, "") // No remote URL

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Expected 1 repository, got %d", result.Count)
	}

	repo := result.Repositories[0]
	if repo.Name != "no-remote-repo" {
		t.Errorf("Expected repository name 'no-remote-repo', got '%s'", repo.Name)
	}

	if repo.RemoteURL != "" {
		t.Errorf("Expected empty remote URL, got '%s'", repo.RemoteURL)
	}
}

// repoNameSet returns the set of repository names in a scan result.
func repoNameSet(result *ScanResult) map[string]bool {
	names := make(map[string]bool, len(result.Repositories))
	for _, repo := range result.Repositories {
		names[repo.Name] = true
	}
	return names
}

// mkRepo creates a git repository at the given path, creating parents as needed.
func mkRepo(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create dir %s: %v", path, err)
	}
	createTestRepo(t, path, "")
}

func TestScan_BoundaryPruningThroughContainers(t *testing.T) {
	tmpDir := t.TempDir()

	// Container directories at varying depths, each holding a repository.
	// The repository under "org-a/team-one" also contains a vendored clone
	// that must never be registered.
	mkRepo(t, filepath.Join(tmpDir, "top-level-repo"))
	mkRepo(t, filepath.Join(tmpDir, "org-a", "team-one", "service-api"))
	mkRepo(t, filepath.Join(tmpDir, "org-a", "team-two", "service-web"))
	mkRepo(t, filepath.Join(tmpDir, "org-b", "tooling"))
	mkRepo(t, filepath.Join(tmpDir, "org-a", "team-one", "service-api", "third_party", "embedded"))

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	want := []string{"top-level-repo", "service-api", "service-web", "tooling"}
	if result.Count != len(want) {
		t.Errorf("Expected %d repositories, got %d: %v", len(want), result.Count, repoNameSet(result))
	}

	names := repoNameSet(result)
	for _, name := range want {
		if !names[name] {
			t.Errorf("Expected repository %q in results, got %v", name, names)
		}
	}
	if names["embedded"] {
		t.Error("Nested clone 'embedded' should have been pruned, but was registered")
	}
}

func TestScan_SiblingsScannedAfterPruning(t *testing.T) {
	tmpDir := t.TempDir()

	// A repository and a sibling container directory at the same level.
	// Pruning below the repository must not stop the sibling being scanned.
	mkRepo(t, filepath.Join(tmpDir, "group", "first-repo"))
	mkRepo(t, filepath.Join(tmpDir, "group", "first-repo", "nested"))
	mkRepo(t, filepath.Join(tmpDir, "group", "second-repo"))
	mkRepo(t, filepath.Join(tmpDir, "other-group", "third-repo"))

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	names := repoNameSet(result)
	for _, want := range []string{"first-repo", "second-repo", "third-repo"} {
		if !names[want] {
			t.Errorf("Expected repository %q, got %v", want, names)
		}
	}
	if names["nested"] {
		t.Error("Nested repository should have been pruned")
	}
	if result.Count != 3 {
		t.Errorf("Expected 3 repositories, got %d: %v", result.Count, names)
	}
}

func TestScan_SkipsFilteredDirectoriesAtEveryLevel(t *testing.T) {
	tmpDir := t.TempDir()

	// Skip-list and hidden directories at depth 1 and deeper.
	mkRepo(t, filepath.Join(tmpDir, "node_modules", "shallow-skipped"))
	mkRepo(t, filepath.Join(tmpDir, "group", "node_modules", "deep-skipped"))
	mkRepo(t, filepath.Join(tmpDir, ".hidden", "shallow-hidden"))
	mkRepo(t, filepath.Join(tmpDir, "group", ".hidden", "deep-hidden"))
	mkRepo(t, filepath.Join(tmpDir, ".dotfiles"))
	mkRepo(t, filepath.Join(tmpDir, "group", "visible-repo"))

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	names := repoNameSet(result)
	if !names["visible-repo"] {
		t.Errorf("Expected 'visible-repo' to be found, got %v", names)
	}
	for _, unwanted := range []string{"shallow-skipped", "deep-skipped", "shallow-hidden", "deep-hidden", ".dotfiles"} {
		if names[unwanted] {
			t.Errorf("Expected %q to be filtered out, got %v", unwanted, names)
		}
	}
	if result.Count != 1 {
		t.Errorf("Expected 1 repository, got %d: %v", result.Count, names)
	}
}

func TestScan_MaxDepth(t *testing.T) {
	tmpDir := t.TempDir()

	// Nest repositories at depths 1 through 8:
	//   tmpDir/d1                      (depth 1)
	//   tmpDir/c1/d2                   (depth 2)
	//   tmpDir/c1/c2/d3                (depth 3) ... and so on
	for depth := 1; depth <= 8; depth++ {
		parts := []string{tmpDir}
		for c := 1; c < depth; c++ {
			parts = append(parts, "c"+string(rune('0'+c)))
		}
		parts = append(parts, "d"+string(rune('0'+depth)))
		mkRepo(t, filepath.Join(parts...))
	}

	tests := []struct {
		name      string
		opts      Options
		wantDepth int
	}{
		{"default depth registers through depth 6", Options{}, 6},
		{"explicit depth 3 registers through depth 3", Options{MaxDepth: 3}, 3},
		{"explicit depth 1 registers only the top level", Options{MaxDepth: 1}, 1},
		{"zero falls back to the default", Options{MaxDepth: 0}, 6},
		{"negative falls back to the default", Options{MaxDepth: -5}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Scan(tmpDir, tt.opts)
			if err != nil {
				t.Fatalf("Scan failed: %v", err)
			}

			names := repoNameSet(result)
			for depth := 1; depth <= 8; depth++ {
				repoName := "d" + string(rune('0'+depth))
				want := depth <= tt.wantDepth
				if names[repoName] != want {
					t.Errorf("repository %q at depth %d: found=%v, want=%v (results: %v)",
						repoName, depth, names[repoName], want, names)
				}
			}
			if result.Count != tt.wantDepth {
				t.Errorf("Expected %d repositories, got %d: %v", tt.wantDepth, result.Count, names)
			}
		})
	}
}

func TestScan_RepositoryAtScanRoot(t *testing.T) {
	tmpDir := t.TempDir()
	mkRepo(t, tmpDir)
	mkRepo(t, filepath.Join(tmpDir, "inner"))

	result, err := Scan(tmpDir, Options{})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// The scan root is itself a repository, so it is registered and nothing
	// below it is traversed.
	if result.Count != 1 {
		t.Errorf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
	}
	if len(result.Repositories) > 0 && result.Repositories[0].Path != tmpDir {
		t.Errorf("Expected the scan root itself, got %s", result.Repositories[0].Path)
	}
}

// symlink creates a symlink at linkPath pointing to target.
func symlink(t *testing.T, target, linkPath string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		t.Fatalf("Failed to create parent for symlink: %v", err)
	}
	if err := os.Symlink(target, linkPath); err != nil {
		t.Fatalf("Failed to create symlink %s -> %s: %v", linkPath, target, err)
	}
}

// realPath resolves a path for comparison against scanner output.
func realPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return resolved
}

// scanWithDeadline runs a scan and fails the test if it does not finish in
// time, so a traversal-loop regression fails rather than hanging CI.
func scanWithDeadline(t *testing.T, root string, opts Options) *ScanResult {
	t.Helper()

	type outcome struct {
		result *ScanResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := Scan(root, opts)
		done <- outcome{result, err}
	}()

	select {
	case got := <-done:
		if got.err != nil {
			t.Fatalf("Scan failed: %v", got.err)
		}
		return got.result
	case <-time.After(30 * time.Second):
		t.Fatal("Scan did not complete within 30s — likely a traversal loop")
		return nil
	}
}

// TestScan_CuratedSymlinkWorkspace covers the workspace model that the previous
// filepath.Walk implementation could not see at all: a workspace directory
// containing nothing but symlinks to repositories living elsewhere.
func TestScan_CuratedSymlinkWorkspace(t *testing.T) {
	external := t.TempDir()
	workspace := t.TempDir()

	names := []string{"dratlab-apps", "dratlab-kubernetes", "dratlab-tooling"}
	for _, name := range names {
		repoDir := filepath.Join(external, name)
		mkRepo(t, repoDir)
		symlink(t, repoDir, filepath.Join(workspace, name))
	}

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != len(names) {
		t.Fatalf("Expected %d repositories, got %d: %v", len(names), result.Count, repoNameSet(result))
	}

	found := repoNameSet(result)
	for _, name := range names {
		if !found[name] {
			t.Errorf("Expected symlinked repository %q to be discovered, got %v", name, found)
		}
	}

	// Stored paths must be the real external locations, not the symlink paths.
	for _, repo := range result.Repositories {
		wantPrefix := realPath(t, external)
		if !strings.HasPrefix(repo.Path, wantPrefix) {
			t.Errorf("Expected repository path under %s, got %s", wantPrefix, repo.Path)
		}
		if !result.ReachablePaths[repo.Path] {
			t.Errorf("Expected %s to be recorded in ReachablePaths", repo.Path)
		}
	}
}

// TestScan_SymlinkedRepoInsideContainer covers a symlink nested below the
// workspace root rather than a direct child.
func TestScan_SymlinkedRepoInsideContainer(t *testing.T) {
	external := t.TempDir()
	workspace := t.TempDir()

	repoDir := filepath.Join(external, "deep-repo")
	mkRepo(t, repoDir)
	symlink(t, repoDir, filepath.Join(workspace, "group", "sub", "deep-repo"))

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != 1 {
		t.Fatalf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
	}
	if got, want := result.Repositories[0].Path, realPath(t, repoDir); got != want {
		t.Errorf("Expected stored path %s, got %s", want, got)
	}
}

// TestScan_DeduplicatesByRealPath asserts that a repository reachable both
// physically and through a symlink registers exactly once, in either traversal
// order. Entry names are chosen so that ReadDir's lexical ordering visits the
// symlink first in one case and the physical directory first in the other.
func TestScan_DeduplicatesByRealPath(t *testing.T) {
	tests := []struct {
		name         string
		physicalName string
		linkName     string
	}{
		{"physical visited first", "aaa-repo", "zzz-link"},
		{"symlink visited first", "zzz-repo", "aaa-link"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()

			repoDir := filepath.Join(workspace, tt.physicalName)
			mkRepo(t, repoDir)
			symlink(t, repoDir, filepath.Join(workspace, tt.linkName))

			result := scanWithDeadline(t, workspace, Options{})

			if result.Count != 1 {
				t.Fatalf("Expected 1 repository after deduplication, got %d: %v",
					result.Count, repoNameSet(result))
			}
			if got, want := result.Repositories[0].Path, realPath(t, repoDir); got != want {
				t.Errorf("Expected the real path %s, got %s", want, got)
			}
			if got, want := result.Repositories[0].Name, tt.physicalName; got != want {
				t.Errorf("Expected name from the real path %q, got %q", want, got)
			}
		})
	}
}

// TestScan_DeduplicatesExternalRepoReachedTwice covers the same repository
// reached through two different symlinks.
func TestScan_DeduplicatesExternalRepoReachedTwice(t *testing.T) {
	external := t.TempDir()
	workspace := t.TempDir()

	repoDir := filepath.Join(external, "shared-repo")
	mkRepo(t, repoDir)
	symlink(t, repoDir, filepath.Join(workspace, "first-alias"))
	symlink(t, repoDir, filepath.Join(workspace, "second-alias"))

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != 1 {
		t.Fatalf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
	}
}

func TestScan_SymlinkLoopSafety(t *testing.T) {
	t.Run("symlink pointing at the workspace root", func(t *testing.T) {
		workspace := t.TempDir()
		mkRepo(t, filepath.Join(workspace, "real-repo"))
		symlink(t, workspace, filepath.Join(workspace, "loop-to-root"))

		result := scanWithDeadline(t, workspace, Options{})

		if result.Count != 1 {
			t.Errorf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
		}
	})

	t.Run("symlink pointing at its own parent", func(t *testing.T) {
		workspace := t.TempDir()
		group := filepath.Join(workspace, "group")
		mkRepo(t, filepath.Join(group, "real-repo"))
		symlink(t, group, filepath.Join(group, "loop-to-parent"))

		result := scanWithDeadline(t, workspace, Options{})

		if result.Count != 1 {
			t.Errorf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
		}
	})

	t.Run("two directories symlinked to each other", func(t *testing.T) {
		workspace := t.TempDir()
		dirA := filepath.Join(workspace, "a")
		dirB := filepath.Join(workspace, "b")
		mkRepo(t, filepath.Join(dirA, "repo-a"))
		mkRepo(t, filepath.Join(dirB, "repo-b"))
		symlink(t, dirB, filepath.Join(dirA, "to-b"))
		symlink(t, dirA, filepath.Join(dirB, "to-a"))

		result := scanWithDeadline(t, workspace, Options{})

		names := repoNameSet(result)
		if !names["repo-a"] || !names["repo-b"] {
			t.Errorf("Expected both repo-a and repo-b, got %v", names)
		}
		if result.Count != 2 {
			t.Errorf("Expected 2 repositories, got %d: %v", result.Count, names)
		}
	})

	t.Run("symlink to an ancestor above the workspace root", func(t *testing.T) {
		base := t.TempDir()
		workspace := filepath.Join(base, "workspace")
		mkRepo(t, filepath.Join(workspace, "real-repo"))
		symlink(t, base, filepath.Join(workspace, "escape-hatch"))

		result := scanWithDeadline(t, workspace, Options{})

		if result.Count != 1 {
			t.Errorf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
		}
	})
}

func TestScan_WorktreesAreNotRegisteredAsRepositories(t *testing.T) {
	workspace := t.TempDir()
	repoDir := filepath.Join(workspace, "main-repo")
	mkRepo(t, repoDir)

	// An aligned worktree under <repo>.wt/ and an unaligned one elsewhere.
	alignedWT := filepath.Join(workspace, "main-repo.wt", "feature")
	unalignedWT := filepath.Join(workspace, "loose", "hotfix")
	gitWorktreeAdd(t, repoDir, "feature", alignedWT)
	gitWorktreeAdd(t, repoDir, "hotfix", unalignedWT)

	result := scanWithDeadline(t, workspace, Options{})

	names := repoNameSet(result)
	if result.Count != 1 {
		t.Fatalf("Expected only the main repository, got %d: %v", result.Count, names)
	}
	if !names["main-repo"] {
		t.Errorf("Expected 'main-repo', got %v", names)
	}
	for _, unwanted := range []string{"feature", "hotfix"} {
		if names[unwanted] {
			t.Errorf("Worktree %q must not be registered as a repository", unwanted)
		}
	}

	// The scan records the main repository inferred from each worktree so the
	// reconciliation layer can register untracked main repositories.
	wantMain := realPath(t, repoDir)
	if !result.WorktreeMainRepos[wantMain] {
		t.Errorf("Expected WorktreeMainRepos to contain %s, got %v", wantMain, result.WorktreeMainRepos)
	}
}

// TestScan_CloneInDotWtDirectoryIsRegistered confirms discovery never filters
// on directory naming: a real clone in a directory ending in .wt is a
// repository like any other.
func TestScan_CloneInDotWtDirectoryIsRegistered(t *testing.T) {
	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "important-repo.wt"))

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != 1 {
		t.Fatalf("Expected 1 repository, got %d: %v", result.Count, repoNameSet(result))
	}
	if got := result.Repositories[0].Name; got != "important-repo.wt" {
		t.Errorf("Expected 'important-repo.wt', got %q", got)
	}
}

func TestScan_BrokenSymlinkRecordedAsError(t *testing.T) {
	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "good-repo"))
	symlink(t, filepath.Join(workspace, "does-not-exist"), filepath.Join(workspace, "dangling"))

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != 1 {
		t.Errorf("Expected the scan to continue and find 1 repository, got %d", result.Count)
	}
	if len(result.Errors) == 0 {
		t.Error("Expected a scan error for the dangling symlink, got none")
	}
}

func TestScan_UnreadableDirectoryRecordedAsError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; permission bits are not enforced")
	}

	workspace := t.TempDir()
	mkRepo(t, filepath.Join(workspace, "good-repo"))

	locked := filepath.Join(workspace, "locked")
	if err := os.MkdirAll(locked, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.Chmod(locked, 0000); err != nil {
		t.Fatalf("Failed to chmod directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0755) })

	result := scanWithDeadline(t, workspace, Options{})

	if result.Count != 1 {
		t.Errorf("Expected the scan to continue and find 1 repository, got %d", result.Count)
	}
	if len(result.Errors) == 0 {
		t.Error("Expected a scan error for the unreadable directory, got none")
	}
}

func TestIsAncestorOrEqual(t *testing.T) {
	tests := []struct {
		candidate string
		path      string
		want      bool
	}{
		{"/a/b", "/a/b", true},
		{"/a", "/a/b", true},
		{"/a", "/a/b/c", true},
		{"/a/b", "/a", false},
		{"/a/b", "/a/bc", false},
		{"/a/b/", "/a/b/c", true},
		{"/x", "/a/b", false},
	}

	for _, tt := range tests {
		got := isAncestorOrEqual(tt.candidate, tt.path)
		if got != tt.want {
			t.Errorf("isAncestorOrEqual(%q, %q) = %v, want %v", tt.candidate, tt.path, got, tt.want)
		}
	}
}

func TestShouldSkipDir(t *testing.T) {
	tests := []struct {
		dirName string
		skip    bool
	}{
		{"node_modules", true},
		{"vendor", true},
		{".venv", true},
		{"venv", true},
		{"__pycache__", true},
		{"target", true},
		{"build", true},
		{"dist", true},
		{".cache", true},
		{".terraform", true},
		{"valid-dir", false},
		{"my-project", false},
		{".git", true},        // hidden directories are skipped at every level
		{".dotfiles", true},   // hidden, even though not in the skip list
		{".", false},          // current-directory marker is not "hidden"
		{"not.hidden", false}, // a dot elsewhere in the name is fine
	}

	for _, tt := range tests {
		result := shouldSkipDir(tt.dirName)
		if result != tt.skip {
			t.Errorf("shouldSkipDir(%s) = %v, want %v", tt.dirName, result, tt.skip)
		}
	}
}
