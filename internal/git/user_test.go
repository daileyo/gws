package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/daileyo/gws/internal/config"
)

// createTestRepoWithUser creates a test git repository with specific user config
func createTestRepoWithUser(t *testing.T, path string, userName, userEmail string, local bool) {
	t.Helper()

	// Initialize the repository
	repo, err := git.PlainInit(path, false)
	if err != nil {
		t.Fatalf("Failed to init test repo: %v", err)
	}

	if local && (userName != "" || userEmail != "") {
		// Set local user config
		cfg, err := repo.Config()
		if err != nil {
			t.Fatalf("Failed to get repo config: %v", err)
		}

		cfg.User.Name = userName
		cfg.User.Email = userEmail

		if err := repo.SetConfig(cfg); err != nil {
			t.Fatalf("Failed to set repo config: %v", err)
		}
	}

	// Create an initial commit so the repo is valid
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(path, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add the file
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	// Commit with the specified user (or default)
	commitUser := userName
	commitEmail := userEmail
	if commitUser == "" {
		commitUser = "Test User"
	}
	if commitEmail == "" {
		commitEmail = "test@example.com"
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitUser,
			Email: commitEmail,
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
}

// createTestRepoWithSigningConfig creates a test repo with signing configuration
func createTestRepoWithSigningConfig(t *testing.T, path string, signingKey string, signCommits bool) {
	t.Helper()

	// First create a basic repo
	createTestRepoWithUser(t, path, "Test User", "test@example.com", true)

	// Now add signing config to .git/config
	configPath := filepath.Join(path, ".git", "config")

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Append signing configuration
	appendConfig := "\n[user]\n"
	if signingKey != "" {
		appendConfig += "\tsigningkey = " + signingKey + "\n"
	}
	if signCommits {
		appendConfig += "[commit]\n\tgpgsign = true\n"
	}

	newData := append(data, []byte(appendConfig)...)
	if err := os.WriteFile(configPath, newData, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
}

func TestGetUserConfig_LocalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create repo with local user config
	createTestRepoWithUser(t, repoPath, "Local User", "local@example.com", true)

	// Get user config
	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}

	// Verify local config is detected
	if userCfg.Name != "Local User" {
		t.Errorf("Expected name 'Local User', got '%s'", userCfg.Name)
	}
	if userCfg.Email != "local@example.com" {
		t.Errorf("Expected email 'local@example.com', got '%s'", userCfg.Email)
	}
	if userCfg.Source != config.UserSourceLocal {
		t.Errorf("Expected source 'local', got '%s'", userCfg.Source)
	}
}

func TestGetUserConfig_GlobalFallback(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create repo without local user config
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("Failed to init test repo: %v", err)
	}

	// Create a commit (required for valid repo)
	worktree, _ := repo.Worktree()
	testFile := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
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

	// Get user config
	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}

	// Source should be global or unknown (depending on system gitconfig)
	// We can't guarantee a global config exists, so just check it doesn't error
	if userCfg == nil {
		t.Error("Expected non-nil user config")
	}
}

func TestGetUserConfig_WithSigningConfig(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create repo with signing config
	createTestRepoWithSigningConfig(t, repoPath, "ABCD1234", true)

	// Get user config
	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}

	// Verify signing config is detected
	if userCfg.SigningKey != "ABCD1234" {
		t.Errorf("Expected signing key 'ABCD1234', got '%s'", userCfg.SigningKey)
	}
	if !userCfg.SignCommits {
		t.Error("Expected SignCommits to be true")
	}
}

func TestGetUserConfig_EmptyRepo(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create an empty repo (no commits)
	_, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("Failed to init test repo: %v", err)
	}

	// Get user config - should not error even with empty repo
	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}

	if userCfg == nil {
		t.Error("Expected non-nil user config for empty repo")
	}
}

func TestGetUserConfig_InvalidPath(t *testing.T) {
	// Try to get user config from non-existent path
	_, err := GetUserConfig("/non/existent/path")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestParseGitConfig(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantName    string
		wantEmail   string
		wantSignKey string
		wantSign    bool
	}{
		{
			name: "basic user config",
			content: `[user]
	name = John Doe
	email = john@example.com`,
			wantName:  "John Doe",
			wantEmail: "john@example.com",
		},
		{
			name: "user config with signing",
			content: `[user]
	name = Jane Doe
	email = jane@example.com
	signingkey = ABC123
[commit]
	gpgsign = true`,
			wantName:    "Jane Doe",
			wantEmail:   "jane@example.com",
			wantSignKey: "ABC123",
			wantSign:    true,
		},
		{
			name: "config with equals in value",
			content: `[user]
	name = John = Doe
	email = john@example.com`,
			wantName:  "John = Doe",
			wantEmail: "john@example.com",
		},
		{
			name: "config with quotes",
			content: `[user]
	name = "Quoted Name"
	email = 'quoted@example.com'`,
			wantName:  "Quoted Name",
			wantEmail: "quoted@example.com",
		},
		{
			name: "config with comments",
			content: `# This is a comment
[user]
	name = Test User
	; Another comment
	email = test@example.com`,
			wantName:  "Test User",
			wantEmail: "test@example.com",
		},
		{
			name:      "empty config",
			content:   "",
			wantName:  "",
			wantEmail: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := parseGitConfig(tt.content)

			if cfg.Name != tt.wantName {
				t.Errorf("Name: expected '%s', got '%s'", tt.wantName, cfg.Name)
			}
			if cfg.Email != tt.wantEmail {
				t.Errorf("Email: expected '%s', got '%s'", tt.wantEmail, cfg.Email)
			}
			if cfg.SigningKey != tt.wantSignKey {
				t.Errorf("SigningKey: expected '%s', got '%s'", tt.wantSignKey, cfg.SigningKey)
			}
			if cfg.SignCommits != tt.wantSign {
				t.Errorf("SignCommits: expected %v, got %v", tt.wantSign, cfg.SignCommits)
			}
		})
	}
}

func TestExtractValue(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"name = John Doe", "John Doe"},
		{"email=test@example.com", "test@example.com"},
		{"  name  =  Spaced Value  ", "Spaced Value"},
		{`name = "Quoted"`, "Quoted"},
		{"name = 'Single Quoted'", "Single Quoted"},
		{"key = value = with = equals", "value = with = equals"},
		{"noequals", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := extractValue(tt.line)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestMatchesGitdirCondition(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name      string
		repoPath  string
		condition string
		expected  bool
	}{
		{
			name:      "exact directory match with trailing slash",
			repoPath:  "/home/user/work/my-project",
			condition: "gitdir:/home/user/work/",
			expected:  true,
		},
		{
			name:      "nested repo under matching directory",
			repoPath:  "/home/user/work/team/my-project",
			condition: "gitdir:/home/user/work/",
			expected:  true,
		},
		{
			name:      "trailing ** glob",
			repoPath:  "/home/user/work/deep/nested/repo",
			condition: "gitdir:/home/user/work/**",
			expected:  true,
		},
		{
			name:      "non-matching path",
			repoPath:  "/home/user/personal/my-project",
			condition: "gitdir:/home/user/work/",
			expected:  false,
		},
		{
			name:      "tilde expansion",
			repoPath:  filepath.Join(home, "work", "my-project"),
			condition: "gitdir:~/work/",
			expected:  true,
		},
		{
			name:      "tilde expansion non-match",
			repoPath:  filepath.Join(home, "personal", "my-project"),
			condition: "gitdir:~/work/",
			expected:  false,
		},
		{
			name:      "case-insensitive gitdir/i match",
			repoPath:  "/home/user/Work/my-project",
			condition: "gitdir/i:/home/user/work/",
			expected:  true,
		},
		{
			name:      "case-sensitive gitdir does not match different case",
			repoPath:  "/home/user/Work/my-project",
			condition: "gitdir:/home/user/work/",
			expected:  false,
		},
		{
			name:      "not a gitdir condition",
			repoPath:  "/home/user/work/repo",
			condition: "onbranch:main",
			expected:  false,
		},
		{
			name:      "directory without trailing slash",
			repoPath:  "/home/user/work/my-project",
			condition: "gitdir:/home/user/work",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchesGitdirCondition(tt.repoPath, tt.condition)
			if result != tt.expected {
				t.Errorf("MatchesGitdirCondition(%q, %q) = %v, want %v",
					tt.repoPath, tt.condition, result, tt.expected)
			}
		})
	}
}

func TestParseIncludeIfs(t *testing.T) {
	home, _ := os.UserHomeDir()

	content := `[user]
	name = Default User
	email = default@example.com
[includeIf "gitdir:~/work/"]
	path = ~/.gitconfig-work
[includeIf "gitdir:~/personal/"]
	path = ~/.gitconfig-personal
[core]
	editor = vim
`
	entries := parseIncludeIfs(content, home)

	if len(entries) != 2 {
		t.Fatalf("Expected 2 includeIf entries, got %d", len(entries))
	}

	if entries[0].condition != "gitdir:~/work/" {
		t.Errorf("Expected condition 'gitdir:~/work/', got '%s'", entries[0].condition)
	}
	expectedPath := filepath.Clean(filepath.Join(home, ".gitconfig-work"))
	if entries[0].path != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, entries[0].path)
	}

	if entries[1].condition != "gitdir:~/personal/" {
		t.Errorf("Expected condition 'gitdir:~/personal/', got '%s'", entries[1].condition)
	}
}

func TestParseIncludeIfs_Empty(t *testing.T) {
	content := `[user]
	name = Default User
	email = default@example.com
`
	entries := parseIncludeIfs(content, "/home/test")

	if len(entries) != 0 {
		t.Errorf("Expected 0 includeIf entries, got %d", len(entries))
	}
}

func TestUserConfigStruct(t *testing.T) {
	cfg := &UserConfig{
		Name:        "Test User",
		Email:       "test@example.com",
		SigningKey:  "KEY123",
		SignCommits: true,
		Source:      config.UserSourceLocal,
	}

	if cfg.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", cfg.Name)
	}
	if cfg.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", cfg.Email)
	}
	if cfg.SigningKey != "KEY123" {
		t.Errorf("Expected signing key 'KEY123', got '%s'", cfg.SigningKey)
	}
	if !cfg.SignCommits {
		t.Error("Expected SignCommits to be true")
	}
	if cfg.Source != config.UserSourceLocal {
		t.Errorf("Expected source 'local', got '%s'", cfg.Source)
	}
}

// Ensure gitconfig import is used (to avoid unused import error)
var _ = gitconfig.Config{}
