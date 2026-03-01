package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/user"
)

func createTestGitRepoWithUser(t *testing.T, name, email string) string {
	t.Helper()
	repoPath := createTestGitRepo(t)

	// Set local user config
	profile := config.Profile{
		GitName: name,
		Email:   email,
	}
	if err := user.AssignLocal(repoPath, profile); err != nil {
		t.Fatalf("Failed to set user config: %v", err)
	}
	return repoPath
}

func createTestGitRepoWithUserAndSigning(t *testing.T, name, email, signingKey string) string {
	t.Helper()
	repoPath := createTestGitRepo(t)

	profile := config.Profile{
		GitName:     name,
		Email:       email,
		SigningKey:  signingKey,
		SignCommits: true,
	}
	if err := user.AssignLocal(repoPath, profile); err != nil {
		t.Fatalf("Failed to set user config: %v", err)
	}
	return repoPath
}

func TestRunUserDelete_SingleRepo(t *testing.T) {
	repoPath := createTestGitRepoWithUser(t, "Local User", "local@test.com")

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, User: "Local User", Email: "local@test.com", UserSource: config.UserSourceLocal},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	origAll := flagAll
	defer func() {
		flagQuiet = origQuiet
		flagAll = origAll
		filterTags = []string{}
	}()

	flagQuiet = true

	err := runUserDelete(nil, []string{"test-repo"})
	if err != nil {
		t.Fatalf("runUserDelete failed: %v", err)
	}

	// Verify .git/config no longer has user.name and user.email
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}
	content := string(data)
	if strings.Contains(content, "name = Local User") {
		t.Error("user.name should have been removed")
	}
	if strings.Contains(content, "email = local@test.com") {
		t.Error("user.email should have been removed")
	}
}

func TestRunUserDelete_WithAll(t *testing.T) {
	repoPath := createTestGitRepoWithUserAndSigning(t, "Test User", "test@test.com", "KEY123")

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, User: "Test User", Email: "test@test.com", SigningEnabled: true, UserSource: config.UserSourceLocal},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	origAll := flagAll
	defer func() {
		flagQuiet = origQuiet
		flagAll = origAll
		filterTags = []string{}
	}()

	flagQuiet = true
	flagAll = true

	err := runUserDelete(nil, []string{"test-repo"})
	if err != nil {
		t.Fatalf("runUserDelete failed: %v", err)
	}

	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, _ := os.ReadFile(gitConfigPath)
	content := string(data)

	if strings.Contains(content, "signingkey") {
		t.Error("signingkey should have been removed with --all")
	}
	if strings.Contains(content, "gpgsign") {
		t.Error("gpgsign should have been removed with --all")
	}
}

func TestRunUserDelete_NoLocalConfig(t *testing.T) {
	// Create a repo with NO local user config
	repoPath := createTestGitRepo(t)

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, UserSource: config.UserSourceGlobal},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		filterTags = []string{}
	}()

	flagQuiet = true

	// Should not error — just skip repos without local config
	err := runUserDelete(nil, []string{"test-repo"})
	if err != nil {
		t.Fatalf("runUserDelete should not fail for repos without local config: %v", err)
	}
}

func TestRunUserDelete_NoMatch(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "my-project", Path: "/nonexistent"},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	defer func() {
		filterTags = []string{}
	}()

	err := runUserDelete(nil, []string{"nonexistent-repo"})
	if err == nil {
		t.Fatal("expected error for no matching repos")
	}
	if !strings.Contains(err.Error(), "no repositories found") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRunUserDelete_NoArgs(t *testing.T) {
	cfg := &config.Config{}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	defer func() {
		filterTags = []string{}
	}()

	err := runUserDelete(nil, []string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
	if !strings.Contains(err.Error(), "provide a repository name or path") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRunUserDelete_ConfigJsonUpdated(t *testing.T) {
	repoPath := createTestGitRepoWithUser(t, "Local User", "local@test.com")

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, User: "Local User", Email: "local@test.com", UserSource: config.UserSourceLocal},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		flagAll = false
		filterTags = []string{}
	}()

	flagQuiet = true

	err := runUserDelete(nil, []string{"test-repo"})
	if err != nil {
		t.Fatalf("runUserDelete failed: %v", err)
	}

	// Load config and verify it was updated
	updatedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	repo := updatedCfg.Repositories[0]
	// After deletion, the user source should no longer be "local"
	if repo.UserSource == config.UserSourceLocal {
		t.Error("UserSource should no longer be 'local' after deletion")
	}
}

func TestRunUserDelete_BatchByTag(t *testing.T) {
	repo1Path := createTestGitRepoWithUser(t, "Work User", "work@company.com")
	repo2Path := createTestGitRepoWithUser(t, "Work User", "work@company.com")
	repo3Path := createTestGitRepoWithUser(t, "Personal User", "personal@test.com")

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "repo-a", Path: repo1Path, User: "Work User", Email: "work@company.com", UserSource: config.UserSourceLocal, Tags: []string{"work"}},
			{Name: "repo-b", Path: repo2Path, User: "Work User", Email: "work@company.com", UserSource: config.UserSourceLocal, Tags: []string{"work", "backend"}},
			{Name: "repo-c", Path: repo3Path, User: "Personal User", Email: "personal@test.com", UserSource: config.UserSourceLocal, Tags: []string{"personal"}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	origAll := flagAll
	defer func() {
		flagQuiet = origQuiet
		flagAll = origAll
		filterTags = []string{}
	}()

	flagQuiet = true
	filterTags = []string{"work"}

	err := runUserDelete(nil, []string{})
	if err != nil {
		t.Fatalf("runUserDelete batch failed: %v", err)
	}

	// Verify tagged repos had local config removed
	for _, repoPath := range []string{repo1Path, repo2Path} {
		gitConfigPath := filepath.Join(repoPath, ".git", "config")
		data, _ := os.ReadFile(gitConfigPath)
		content := string(data)
		if strings.Contains(content, "name = Work User") {
			t.Errorf("Expected user.name removed from %s", repoPath)
		}
	}

	// Verify non-tagged repo was NOT touched
	gitConfigPath := filepath.Join(repo3Path, ".git", "config")
	data, _ := os.ReadFile(gitConfigPath)
	if !strings.Contains(string(data), "name = Personal User") {
		t.Error("repo-c should still have local user config (no 'work' tag)")
	}
}

func TestRunUserDelete_BatchNoTagMatch(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "repo-a", Path: "/nonexistent", Tags: []string{"personal"}, UserSource: config.UserSourceLocal},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		filterTags = []string{}
	}()

	filterTags = []string{"work"}

	err := runUserDelete(nil, []string{})
	if err == nil {
		t.Fatal("expected error for no matching tags")
	}
	if !strings.Contains(err.Error(), "no repositories found with tag(s)") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRunUserDelete_BatchQuietSuppressesOutput(t *testing.T) {
	repoPath := createTestGitRepoWithUser(t, "Work User", "work@company.com")

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "repo-a", Path: repoPath, User: "Work User", Email: "work@company.com", UserSource: config.UserSourceLocal, Tags: []string{"work"}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	origAll := flagAll
	defer func() {
		flagQuiet = origQuiet
		flagAll = origAll
		filterTags = []string{}
	}()

	flagQuiet = true
	filterTags = []string{"work"}

	// Should succeed with no output (quiet mode)
	err := runUserDelete(nil, []string{})
	if err != nil {
		t.Fatalf("runUserDelete batch quiet failed: %v", err)
	}

	// Verify config was still updated
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, _ := os.ReadFile(gitConfigPath)
	if strings.Contains(string(data), "name = Work User") {
		t.Error("Expected user.name removed")
	}
}

// Ensure createTestGitRepo helper is available (defined in userupdate_test.go)
// If running tests in isolation, this needs to be in the same package.
var _ = exec.Command // ensure exec import is used
