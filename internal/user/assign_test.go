package user

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestAssignLocal(t *testing.T) {
	// Create a temporary git repository
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Initialize git repo
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	profile := config.Profile{
		Name:    "test",
		GitName: "Test User",
		Email:   "test@example.com",
	}

	// Assign profile
	if err := AssignLocal(repoPath, profile); err != nil {
		t.Fatalf("AssignLocal failed: %v", err)
	}

	// Verify git config was set
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Test User") {
		t.Error("Expected user.name 'Test User' in git config")
	}
	if !strings.Contains(content, "test@example.com") {
		t.Error("Expected user.email 'test@example.com' in git config")
	}
}

func TestAssignLocalWithSigning(t *testing.T) {
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

	profile := config.Profile{
		Name:        "test",
		GitName:     "Test User",
		Email:       "test@example.com",
		SigningKey:  "ABCD1234",
		SignCommits: true,
	}

	if err := AssignLocal(repoPath, profile); err != nil {
		t.Fatalf("AssignLocal failed: %v", err)
	}

	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		t.Fatalf("Failed to read git config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "signingkey") {
		t.Error("Expected signingkey in git config")
	}
	if !strings.Contains(content, "gpgsign") {
		t.Error("Expected gpgsign in git config")
	}
}

func TestAssignLocalNonRepo(t *testing.T) {
	tmpDir := t.TempDir()

	profile := config.Profile{
		Name:    "test",
		GitName: "Test User",
		Email:   "test@example.com",
	}

	err := AssignLocal(tmpDir, profile)
	if err == nil {
		t.Error("Expected error for non-git directory")
	}
}

func TestCreateProfileSubdir(t *testing.T) {
	tmpDir := t.TempDir()

	subdirPath, err := CreateProfileSubdir(tmpDir, "work")
	if err != nil {
		t.Fatalf("CreateProfileSubdir failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "work")
	if subdirPath != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, subdirPath)
	}

	// Verify directory exists
	if _, err := os.Stat(subdirPath); os.IsNotExist(err) {
		t.Error("Subdirectory was not created")
	}
}

func TestMoveRepository(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source directory with a file
	srcPath := filepath.Join(tmpDir, "source-repo")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	testFile := filepath.Join(srcPath, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Move to destination
	dstPath := filepath.Join(tmpDir, "dest", "repo")
	if err := MoveRepository(srcPath, dstPath); err != nil {
		t.Fatalf("MoveRepository failed: %v", err)
	}

	// Verify source doesn't exist
	if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
		t.Error("Source directory still exists after move")
	}

	// Verify destination exists with content
	movedFile := filepath.Join(dstPath, "test.txt")
	data, err := os.ReadFile(movedFile)
	if err != nil {
		t.Fatalf("Failed to read moved file: %v", err)
	}
	if string(data) != "test content" {
		t.Error("Moved file content doesn't match")
	}
}

func TestMoveRepositoryDestinationExists(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source")
	dstPath := filepath.Join(tmpDir, "dest")

	if err := os.MkdirAll(srcPath, 0755); err != nil {
		t.Fatalf("Failed to create source: %v", err)
	}
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		t.Fatalf("Failed to create dest: %v", err)
	}

	err := MoveRepository(srcPath, dstPath)
	if err == nil {
		t.Error("Expected error when destination exists")
	}
}

func TestCreateProfileGitconfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gitconfig-test")

	profile := config.Profile{
		Name:        "test",
		GitName:     "Test User",
		Email:       "test@example.com",
		SigningKey:  "KEY123",
		SignCommits: true,
	}

	if err := CreateProfileGitconfig(configPath, profile); err != nil {
		t.Fatalf("CreateProfileGitconfig failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	content := string(data)

	checks := []string{
		"[user]",
		"name = Test User",
		"email = test@example.com",
		"signingkey = KEY123",
		"[commit]",
		"gpgsign = true",
	}

	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("Expected '%s' in gitconfig", check)
		}
	}
}

func TestSyncUserInfo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test git repo
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Set up repo with local user config
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	data, _ := os.ReadFile(gitConfigPath)
	content := string(data) + "\n[user]\n\tname = Local User\n\temail = local@test.com\n"
	if err := os.WriteFile(gitConfigPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write git config: %v", err)
	}

	// Create config with outdated info
	cfg := &config.Config{
		Repositories: []config.Repository{
			{
				Name:  "test-repo",
				Path:  repoPath,
				User:  "Old User",
				Email: "old@test.com",
			},
		},
	}

	// Sync
	updated, err := SyncUserInfo(cfg)
	if err != nil {
		t.Fatalf("SyncUserInfo failed: %v", err)
	}

	if updated != 1 {
		t.Errorf("Expected 1 updated repo, got %d", updated)
	}

	// Verify update
	if cfg.Repositories[0].User != "Local User" {
		t.Errorf("Expected user 'Local User', got '%s'", cfg.Repositories[0].User)
	}
	if cfg.Repositories[0].Email != "local@test.com" {
		t.Errorf("Expected email 'local@test.com', got '%s'", cfg.Repositories[0].Email)
	}
}

func TestSetGitConfigValue(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		section     string
		key         string
		value       string
		mustContain []string
	}{
		{
			name:        "add to empty config",
			content:     "",
			section:     "user",
			key:         "name",
			value:       "Test User",
			mustContain: []string{"[user]", "name = Test User"},
		},
		{
			name:        "add to existing section",
			content:     "[user]\n\temail = test@test.com\n",
			section:     "user",
			key:         "name",
			value:       "Test User",
			mustContain: []string{"[user]", "email = test@test.com", "name = Test User"},
		},
		{
			name:        "update existing key",
			content:     "[user]\n\tname = Old Name\n",
			section:     "user",
			key:         "name",
			value:       "New Name",
			mustContain: []string{"[user]", "name = New Name"},
		},
		{
			name:        "add new section",
			content:     "[core]\n\teditor = vim\n",
			section:     "user",
			key:         "name",
			value:       "Test User",
			mustContain: []string{"[core]", "editor = vim", "[user]", "name = Test User"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setGitConfigValue(tt.content, tt.section, tt.key, tt.value)
			for _, expected := range tt.mustContain {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestPreviewAssignLocal(t *testing.T) {
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

	profile := config.Profile{
		Name:    "test",
		GitName: "New User",
		Email:   "new@example.com",
	}

	changes, err := PreviewAssignLocal(repoPath, profile)
	if err != nil {
		t.Fatalf("PreviewAssignLocal failed: %v", err)
	}

	// Should have changes since repo has no user configured
	if len(changes) < 2 {
		t.Errorf("Expected at least 2 changes, got %d", len(changes))
	}
}

func TestRemoveGitConfigKey(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		section        string
		key            string
		mustContain    []string
		mustNotContain []string
	}{
		{
			name:           "remove key from section",
			content:        "[user]\n\tname = Test User\n\temail = test@test.com\n",
			section:        "user",
			key:            "name",
			mustContain:    []string{"[user]", "email = test@test.com"},
			mustNotContain: []string{"name = Test User"},
		},
		{
			name:           "remove last key removes section header",
			content:        "[core]\n\teditor = vim\n[user]\n\tname = Test User\n",
			section:        "user",
			key:            "name",
			mustContain:    []string{"[core]", "editor = vim"},
			mustNotContain: []string{"[user]", "name = Test User"},
		},
		{
			name:           "remove key preserves other sections",
			content:        "[core]\n\teditor = vim\n[user]\n\tname = Test\n\temail = t@t.com\n[commit]\n\tgpgsign = true\n",
			section:        "user",
			key:            "name",
			mustContain:    []string{"[core]", "editor = vim", "[user]", "email = t@t.com", "[commit]", "gpgsign = true"},
			mustNotContain: []string{"name = Test"},
		},
		{
			name:           "remove nonexistent key is no-op",
			content:        "[user]\n\tname = Test\n",
			section:        "user",
			key:            "email",
			mustContain:    []string{"[user]", "name = Test"},
			mustNotContain: []string{},
		},
		{
			name:           "remove from nonexistent section is no-op",
			content:        "[core]\n\teditor = vim\n",
			section:        "user",
			key:            "name",
			mustContain:    []string{"[core]", "editor = vim"},
			mustNotContain: []string{"[user]"},
		},
		{
			name:           "case insensitive key match",
			content:        "[user]\n\tName = Test User\n\temail = test@test.com\n",
			section:        "user",
			key:            "name",
			mustContain:    []string{"[user]", "email = test@test.com"},
			mustNotContain: []string{"Name = Test User"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeGitConfigKey(tt.content, tt.section, tt.key)
			for _, expected := range tt.mustContain {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', got:\n%s", expected, result)
				}
			}
			for _, notExpected := range tt.mustNotContain {
				if strings.Contains(result, notExpected) {
					t.Errorf("Expected result NOT to contain '%s', got:\n%s", notExpected, result)
				}
			}
		})
	}
}

func TestDeleteLocal(t *testing.T) {
	t.Run("remove name and email only", func(t *testing.T) {
		tmpDir := t.TempDir()
		repoPath := filepath.Join(tmpDir, "test-repo")
		os.MkdirAll(repoPath, 0755)

		cmd := exec.Command("git", "init")
		cmd.Dir = repoPath
		cmd.Run()

		// Set user config with signing
		profile := config.Profile{
			GitName:     "Test User",
			Email:       "test@example.com",
			SigningKey:  "KEY123",
			SignCommits: true,
		}
		AssignLocal(repoPath, profile)

		// Delete without --all
		if err := DeleteLocal(repoPath, false); err != nil {
			t.Fatalf("DeleteLocal failed: %v", err)
		}

		data, _ := os.ReadFile(filepath.Join(repoPath, ".git", "config"))
		content := string(data)

		if strings.Contains(content, "name = Test User") {
			t.Error("user.name should have been removed")
		}
		if strings.Contains(content, "email = test@example.com") {
			t.Error("user.email should have been removed")
		}
		// Signing should still be present
		if !strings.Contains(content, "signingkey = KEY123") {
			t.Error("signingkey should still be present")
		}
		if !strings.Contains(content, "gpgsign = true") {
			t.Error("gpgsign should still be present")
		}
	})

	t.Run("remove all including signing", func(t *testing.T) {
		tmpDir := t.TempDir()
		repoPath := filepath.Join(tmpDir, "test-repo")
		os.MkdirAll(repoPath, 0755)

		cmd := exec.Command("git", "init")
		cmd.Dir = repoPath
		cmd.Run()

		profile := config.Profile{
			GitName:     "Test User",
			Email:       "test@example.com",
			SigningKey:  "KEY123",
			SignCommits: true,
		}
		AssignLocal(repoPath, profile)

		// Delete with --all
		if err := DeleteLocal(repoPath, true); err != nil {
			t.Fatalf("DeleteLocal failed: %v", err)
		}

		data, _ := os.ReadFile(filepath.Join(repoPath, ".git", "config"))
		content := string(data)

		if strings.Contains(content, "name = Test User") {
			t.Error("user.name should have been removed")
		}
		if strings.Contains(content, "email = test@example.com") {
			t.Error("user.email should have been removed")
		}
		if strings.Contains(content, "signingkey") {
			t.Error("signingkey should have been removed")
		}
		if strings.Contains(content, "gpgsign") {
			t.Error("gpgsign should have been removed")
		}
	})

	t.Run("preserves non-user sections", func(t *testing.T) {
		tmpDir := t.TempDir()
		repoPath := filepath.Join(tmpDir, "test-repo")
		os.MkdirAll(repoPath, 0755)

		cmd := exec.Command("git", "init")
		cmd.Dir = repoPath
		cmd.Run()

		AssignLocal(repoPath, config.Profile{GitName: "Test", Email: "t@t.com"})

		if err := DeleteLocal(repoPath, false); err != nil {
			t.Fatalf("DeleteLocal failed: %v", err)
		}

		data, _ := os.ReadFile(filepath.Join(repoPath, ".git", "config"))
		content := string(data)

		// Core section should still exist
		if !strings.Contains(content, "[core]") {
			t.Error("[core] section should be preserved")
		}
	})

	t.Run("error for non-git repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := DeleteLocal(tmpDir, false)
		if err == nil {
			t.Error("Expected error for non-git directory")
		}
	})
}
