package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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
