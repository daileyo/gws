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
