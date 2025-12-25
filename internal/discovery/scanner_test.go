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

	result, err := Scan(tmpDir)
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

	result, err := Scan(tmpDir)
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

	result, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should find both parent and child repos
	if result.Count != 2 {
		t.Errorf("Expected 2 repositories, got %d", result.Count)
	}
}

func TestScan_NonExistentDirectory(t *testing.T) {
	result, err := Scan("/nonexistent/directory/path")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for nonexistent directory")
	}
}

func TestScan_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := Scan(tmpDir)
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

	result, err := Scan(tmpDir)
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

	result, err := Scan(tmpDir)
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
		{".git", false}, // .git is handled separately
	}

	for _, tt := range tests {
		result := shouldSkipDir(tt.dirName)
		if result != tt.skip {
			t.Errorf("shouldSkipDir(%s) = %v, want %v", tt.dirName, result, tt.skip)
		}
	}
}
