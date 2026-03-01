package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func createTestGitRepo(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	return repoPath
}

func TestResolveProfile(t *testing.T) {
	resetInlineFlags := func() {
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}

	t.Run("named profile found in stored profiles", func(t *testing.T) {
		defer resetInlineFlags()
		cfg := &config.Config{
			Profiles: []config.Profile{
				{Name: "work", GitName: "Work User", Email: "work@company.com"},
			},
		}

		profile, err := resolveProfile(cfg, []string{"my-repo", "work"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.GitName != "Work User" {
			t.Errorf("expected GitName 'Work User', got %q", profile.GitName)
		}
		if profile.Email != "work@company.com" {
			t.Errorf("expected Email 'work@company.com', got %q", profile.Email)
		}
	})

	t.Run("named profile not found", func(t *testing.T) {
		defer resetInlineFlags()
		cfg := &config.Config{
			Profiles: []config.Profile{},
		}

		_, err := resolveProfile(cfg, []string{"my-repo", "nonexistent"})
		if err == nil {
			t.Fatal("expected error for nonexistent profile")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected 'not found' error, got: %s", err.Error())
		}
	})

	t.Run("inline values only", func(t *testing.T) {
		defer resetInlineFlags()
		flagInlineName = "Inline User"
		flagInlineEmail = "inline@example.com"
		cfg := &config.Config{}

		profile, err := resolveProfile(cfg, []string{"my-repo"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.GitName != "Inline User" {
			t.Errorf("expected GitName 'Inline User', got %q", profile.GitName)
		}
		if profile.Email != "inline@example.com" {
			t.Errorf("expected Email 'inline@example.com', got %q", profile.Email)
		}
	})

	t.Run("inline email only uses email prefix as name", func(t *testing.T) {
		defer resetInlineFlags()
		flagInlineEmail = "john@example.com"
		cfg := &config.Config{}

		profile, err := resolveProfile(cfg, []string{"my-repo"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.GitName != "john" {
			t.Errorf("expected GitName 'john', got %q", profile.GitName)
		}
	})

	t.Run("profile name with inline overrides", func(t *testing.T) {
		defer resetInlineFlags()
		flagInlineName = "Override Name"
		cfg := &config.Config{
			Profiles: []config.Profile{
				{Name: "work", GitName: "Work User", Email: "work@company.com"},
			},
		}

		profile, err := resolveProfile(cfg, []string{"my-repo", "work"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.GitName != "Override Name" {
			t.Errorf("expected GitName 'Override Name', got %q", profile.GitName)
		}
		if profile.Email != "work@company.com" {
			t.Errorf("expected Email 'work@company.com' (from profile), got %q", profile.Email)
		}
	})

	t.Run("no profile name and no inline values returns error", func(t *testing.T) {
		defer resetInlineFlags()
		cfg := &config.Config{}

		_, err := resolveProfile(cfg, []string{"my-repo"})
		if err == nil {
			t.Fatal("expected error when no profile or inline values provided")
		}
		if !strings.Contains(err.Error(), "provide a profile name or --git-email") {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("tag mode uses first arg as profile name", func(t *testing.T) {
		defer resetInlineFlags()
		filterTags = []string{"work"}
		cfg := &config.Config{
			Profiles: []config.Profile{
				{Name: "work", GitName: "Work User", Email: "work@company.com"},
			},
		}

		profile, err := resolveProfile(cfg, []string{"work"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.GitName != "Work User" {
			t.Errorf("expected GitName 'Work User', got %q", profile.GitName)
		}
	})
}

func TestRunUserUpdate_SingleRepo(t *testing.T) {
	repoPath := createTestGitRepo(t)

	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "Work User", Email: "work@company.com"},
		},
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, Tags: []string{}},
		},
	}

	// Save config so runUserUpdate can load it
	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	// Reset flags
	origQuiet := flagQuiet
	origVerbose := flagVerbose
	defer func() {
		flagQuiet = origQuiet
		flagVerbose = origVerbose
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagQuiet = true // suppress output for test

	err := runUserUpdate(nil, []string{"test-repo", "work"})
	if err != nil {
		t.Fatalf("runUserUpdate failed: %v", err)
	}

	// Verify .git/config was updated
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Work User") {
		t.Error("Expected user.name 'Work User' in git config")
	}
	if !strings.Contains(content, "work@company.com") {
		t.Error("Expected user.email 'work@company.com' in git config")
	}
}

func TestRunUserUpdate_InlineValues(t *testing.T) {
	repoPath := createTestGitRepo(t)

	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repoPath, Tags: []string{}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagQuiet = true
	flagInlineName = "Inline User"
	flagInlineEmail = "inline@test.com"

	err := runUserUpdate(nil, []string{"test-repo"})
	if err != nil {
		t.Fatalf("runUserUpdate failed: %v", err)
	}

	// Verify .git/config
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Inline User") {
		t.Error("Expected user.name 'Inline User' in git config")
	}
	if !strings.Contains(content, "inline@test.com") {
		t.Error("Expected user.email 'inline@test.com' in git config")
	}
}

func TestRunUserUpdate_MultipleRepoMatch(t *testing.T) {
	repo1Path := createTestGitRepo(t)
	repo2Path := createTestGitRepo(t)

	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "Work User", Email: "work@company.com"},
		},
		Repositories: []config.Repository{
			{Name: "test-repo", Path: repo1Path, Tags: []string{}},
			{Name: "test-repo-api", Path: repo2Path, Tags: []string{}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagQuiet = true

	err := runUserUpdate(nil, []string{"test-repo", "work"})
	if err != nil {
		t.Fatalf("runUserUpdate failed: %v", err)
	}

	// Verify both repos were updated
	for _, repoPath := range []string{repo1Path, repo2Path} {
		gitConfigPath := filepath.Join(repoPath, ".git", "config")
		data, err := os.ReadFile(gitConfigPath)
		if err != nil {
			t.Fatalf("Failed to read git config for %s: %v", repoPath, err)
		}
		content := string(data)
		if !strings.Contains(content, "Work User") {
			t.Errorf("Expected user.name 'Work User' in git config for %s", repoPath)
		}
	}
}

func TestRunUserUpdate_NoMatch(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "my-project", Path: "/nonexistent"},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	defer func() {
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagInlineEmail = "test@test.com"

	err := runUserUpdate(nil, []string{"nonexistent-repo"})
	if err == nil {
		t.Fatal("expected error for no matching repos")
	}
	if !strings.Contains(err.Error(), "no repositories found") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRunUserUpdate_NoArgs(t *testing.T) {
	cfg := &config.Config{}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	defer func() {
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagInlineEmail = "test@test.com"

	err := runUserUpdate(nil, []string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
	if !strings.Contains(err.Error(), "provide a repository name or path") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRunUserUpdate_BatchByTag(t *testing.T) {
	repo1Path := createTestGitRepo(t)
	repo2Path := createTestGitRepo(t)
	repo3Path := createTestGitRepo(t)

	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "Work User", Email: "work@company.com"},
		},
		Repositories: []config.Repository{
			{Name: "repo-a", Path: repo1Path, Tags: []string{"work"}},
			{Name: "repo-b", Path: repo2Path, Tags: []string{"work", "backend"}},
			{Name: "repo-c", Path: repo3Path, Tags: []string{"personal"}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagQuiet = true
	filterTags = []string{"work"}

	err := runUserUpdate(nil, []string{"work"})
	if err != nil {
		t.Fatalf("runUserUpdate batch failed: %v", err)
	}

	// Verify tagged repos were updated
	for _, repoPath := range []string{repo1Path, repo2Path} {
		gitConfigPath := filepath.Join(repoPath, ".git", "config")
		data, err := os.ReadFile(gitConfigPath)
		if err != nil {
			t.Fatalf("Failed to read git config for %s: %v", repoPath, err)
		}
		content := string(data)
		if !strings.Contains(content, "Work User") {
			t.Errorf("Expected user.name 'Work User' in git config for %s", repoPath)
		}
	}

	// Verify non-tagged repo was NOT updated
	gitConfigPath := filepath.Join(repo3Path, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}
	if strings.Contains(string(data), "Work User") {
		t.Error("repo-c should not have been updated (no 'work' tag)")
	}
}

func TestRunUserUpdate_BatchByMultipleTags(t *testing.T) {
	repo1Path := createTestGitRepo(t)
	repo2Path := createTestGitRepo(t)

	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "Work User", Email: "work@company.com"},
		},
		Repositories: []config.Repository{
			{Name: "repo-a", Path: repo1Path, Tags: []string{"work", "backend"}},
			{Name: "repo-b", Path: repo2Path, Tags: []string{"work"}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	origQuiet := flagQuiet
	defer func() {
		flagQuiet = origQuiet
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	flagQuiet = true
	filterTags = []string{"work", "backend"} // AND logic

	err := runUserUpdate(nil, []string{"work"})
	if err != nil {
		t.Fatalf("runUserUpdate batch failed: %v", err)
	}

	// Only repo-a has both tags
	data1, _ := os.ReadFile(filepath.Join(repo1Path, ".git", "config"))
	if !strings.Contains(string(data1), "Work User") {
		t.Error("repo-a should have been updated (has both tags)")
	}

	data2, _ := os.ReadFile(filepath.Join(repo2Path, ".git", "config"))
	if strings.Contains(string(data2), "Work User") {
		t.Error("repo-b should not have been updated (missing 'backend' tag)")
	}
}

func TestRunUserUpdate_BatchNoTagMatch(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "Work User", Email: "work@company.com"},
		},
		Repositories: []config.Repository{
			{Name: "repo-a", Path: "/nonexistent", Tags: []string{"personal"}},
		},
	}

	origLoad := setupTestConfig(t, cfg)
	defer origLoad()

	defer func() {
		flagInlineName = ""
		flagInlineEmail = ""
		filterTags = []string{}
	}()

	filterTags = []string{"work"}

	err := runUserUpdate(nil, []string{"work"})
	if err == nil {
		t.Fatal("expected error for no matching tags")
	}
	if !strings.Contains(err.Error(), "no repositories found with tag(s)") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

// setupTestConfig saves a config to a temp location and patches config.Load
// to return it. Returns a cleanup function to restore the original.
func setupTestConfig(t *testing.T, cfg *config.Config) func() {
	t.Helper()
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".gws")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	// Save the config
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	return func() {
		os.Setenv("HOME", origHome)
	}
}
