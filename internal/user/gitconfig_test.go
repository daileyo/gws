package user

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseGitconfig_BasicKeyValuePairs(t *testing.T) {
	// Create a temporary gitconfig file
	tmpDir := t.TempDir()
	gitconfigPath := filepath.Join(tmpDir, ".gitconfig")

	content := `[user]
	name = John Doe
	email = john@example.com
[core]
	editor = vim
`
	if err := os.WriteFile(gitconfigPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test gitconfig: %v", err)
	}

	cfg, err := ParseGitconfig(gitconfigPath)
	if err != nil {
		t.Fatalf("ParseGitconfig failed: %v", err)
	}

	// Check user section
	userSection, ok := cfg.Sections["user"]
	if !ok {
		t.Fatal("Expected user section not found")
	}

	if userSection.Values["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", userSection.Values["name"])
	}

	if userSection.Values["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", userSection.Values["email"])
	}

	// Check core section
	coreSection, ok := cfg.Sections["core"]
	if !ok {
		t.Fatal("Expected core section not found")
	}

	if coreSection.Values["editor"] != "vim" {
		t.Errorf("Expected editor 'vim', got '%s'", coreSection.Values["editor"])
	}
}

func TestParseGitconfig_QuotedValues(t *testing.T) {
	tmpDir := t.TempDir()
	gitconfigPath := filepath.Join(tmpDir, ".gitconfig")

	content := `[user]
	name = "John Doe"
	email = 'john@example.com'
`
	if err := os.WriteFile(gitconfigPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test gitconfig: %v", err)
	}

	cfg, err := ParseGitconfig(gitconfigPath)
	if err != nil {
		t.Fatalf("ParseGitconfig failed: %v", err)
	}

	userSection := cfg.Sections["user"]
	if userSection.Values["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", userSection.Values["name"])
	}

	if userSection.Values["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", userSection.Values["email"])
	}
}

func TestParseGitconfig_Comments(t *testing.T) {
	tmpDir := t.TempDir()
	gitconfigPath := filepath.Join(tmpDir, ".gitconfig")

	content := `# This is a comment
[user]
	; Another comment
	name = John Doe
	# Comment in section
	email = john@example.com
`
	if err := os.WriteFile(gitconfigPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test gitconfig: %v", err)
	}

	cfg, err := ParseGitconfig(gitconfigPath)
	if err != nil {
		t.Fatalf("ParseGitconfig failed: %v", err)
	}

	userSection := cfg.Sections["user"]
	if userSection.Values["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", userSection.Values["name"])
	}

	if userSection.Values["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", userSection.Values["email"])
	}
}

func TestExtractIncludeIfs_GitdirPatterns(t *testing.T) {
	tmpDir := t.TempDir()
	gitconfigPath := filepath.Join(tmpDir, ".gitconfig")

	// Create included config files
	workConfigPath := filepath.Join(tmpDir, ".gitconfig-work")
	personalConfigPath := filepath.Join(tmpDir, ".gitconfig-personal")

	if err := os.WriteFile(workConfigPath, []byte("[user]\n\tname = Work User\n"), 0644); err != nil {
		t.Fatalf("Failed to write work config: %v", err)
	}
	if err := os.WriteFile(personalConfigPath, []byte("[user]\n\tname = Personal User\n"), 0644); err != nil {
		t.Fatalf("Failed to write personal config: %v", err)
	}

	content := `[user]
	name = Default User
	email = default@example.com

[includeIf "gitdir:~/work/"]
	path = ` + workConfigPath + `

[includeIf "gitdir:~/personal/"]
	path = ` + personalConfigPath + `
`
	if err := os.WriteFile(gitconfigPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test gitconfig: %v", err)
	}

	cfg, err := ParseGitconfig(gitconfigPath)
	if err != nil {
		t.Fatalf("ParseGitconfig failed: %v", err)
	}

	includeIfs := ExtractIncludeIfs(cfg)
	if len(includeIfs) != 2 {
		t.Fatalf("Expected 2 includeIf directives, got %d", len(includeIfs))
	}

	if includeIfs[0].Condition != "gitdir:~/work/" {
		t.Errorf("Expected condition 'gitdir:~/work/', got '%s'", includeIfs[0].Condition)
	}

	if includeIfs[0].Path != workConfigPath {
		t.Errorf("Expected path '%s', got '%s'", workConfigPath, includeIfs[0].Path)
	}

	if includeIfs[1].Condition != "gitdir:~/personal/" {
		t.Errorf("Expected condition 'gitdir:~/personal/', got '%s'", includeIfs[1].Condition)
	}
}

func TestParseIncludedConfig_UserInfo(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gitconfig-work")

	content := `[user]
	name = Work User
	email = work@company.com
	signingkey = ABCD1234

[commit]
	gpgsign = true
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	includeIf := IncludeIf{
		Condition: "gitdir:~/work/",
		Path:      configPath,
	}

	profile, err := ParseIncludedConfig(includeIf)
	if err != nil {
		t.Fatalf("ParseIncludedConfig failed: %v", err)
	}

	if profile.Name != "work" {
		t.Errorf("Expected profile name 'work', got '%s'", profile.Name)
	}

	if profile.GitName != "Work User" {
		t.Errorf("Expected git name 'Work User', got '%s'", profile.GitName)
	}

	if profile.Email != "work@company.com" {
		t.Errorf("Expected email 'work@company.com', got '%s'", profile.Email)
	}

	if profile.SigningKey != "ABCD1234" {
		t.Errorf("Expected signing key 'ABCD1234', got '%s'", profile.SigningKey)
	}

	if !profile.SignCommits {
		t.Error("Expected SignCommits to be true")
	}
}

func TestParseIncludedConfig_NoUserSection(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gitconfig-empty")

	content := `[core]
	editor = vim
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	includeIf := IncludeIf{
		Condition: "gitdir:~/empty/",
		Path:      configPath,
	}

	_, err := ParseIncludedConfig(includeIf)
	if err == nil {
		t.Error("Expected error for config without user section")
	}
}

func TestDeriveProfileName_FromGitdirCondition(t *testing.T) {
	tests := []struct {
		name      string
		includeIf IncludeIf
		expected  string
	}{
		{
			name: "gitdir with trailing slash",
			includeIf: IncludeIf{
				Condition: "gitdir:~/work/",
				Path:      "/home/user/.gitconfig-work",
			},
			expected: "work",
		},
		{
			name: "gitdir with wildcard",
			includeIf: IncludeIf{
				Condition: "gitdir:~/personal/**",
				Path:      "/home/user/.gitconfig-personal",
			},
			expected: "personal",
		},
		{
			name: "gitdir case insensitive",
			includeIf: IncludeIf{
				Condition: "gitdir/i:~/GitHub/",
				Path:      "/home/user/.gitconfig-github",
			},
			expected: "github",
		},
		{
			name: "derive from config filename",
			includeIf: IncludeIf{
				Condition: "gitdir:~/repos/",
				Path:      "/home/user/.gitconfig-ado",
			},
			expected: "ado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveProfileName(tt.includeIf)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestExtractNameFromConfigPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{".gitconfig-work", "work"},
		{".gitconfig-personal", "personal"},
		{"gitconfig-ado", "ado"},
		{".gitconfig.github", "github"},
		{"work.gitconfig", "work"},
		{"/home/user/.gitconfig-horizon", "horizon"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractNameFromConfigPath(tt.path)
			if result != tt.expected {
				t.Errorf("extractNameFromConfigPath(%s): expected '%s', got '%s'", tt.path, tt.expected, result)
			}
		})
	}
}

func TestDetectProfiles_Integration(t *testing.T) {
	// Create a temporary directory structure that mimics a real setup
	tmpDir := t.TempDir()

	// Create included config files
	workDir := filepath.Join(tmpDir, "work")
	personalDir := filepath.Join(tmpDir, "personal")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		t.Fatalf("Failed to create work dir: %v", err)
	}
	if err := os.MkdirAll(personalDir, 0755); err != nil {
		t.Fatalf("Failed to create personal dir: %v", err)
	}

	workConfigPath := filepath.Join(tmpDir, ".gitconfig-work")
	personalConfigPath := filepath.Join(tmpDir, ".gitconfig-personal")

	workConfig := `[user]
	name = Work User
	email = work@company.com
`
	personalConfig := `[user]
	name = Personal User
	email = personal@gmail.com
	signingkey = PERSONAL123

[commit]
	gpgsign = true
`
	if err := os.WriteFile(workConfigPath, []byte(workConfig), 0644); err != nil {
		t.Fatalf("Failed to write work config: %v", err)
	}
	if err := os.WriteFile(personalConfigPath, []byte(personalConfig), 0644); err != nil {
		t.Fatalf("Failed to write personal config: %v", err)
	}

	// Create main gitconfig with absolute paths
	mainConfig := `[user]
	name = Default User
	email = default@example.com

[includeIf "gitdir:` + workDir + `/"]
	path = ` + workConfigPath + `

[includeIf "gitdir:` + personalDir + `/"]
	path = ` + personalConfigPath + `
`
	mainConfigPath := filepath.Join(tmpDir, ".gitconfig")
	if err := os.WriteFile(mainConfigPath, []byte(mainConfig), 0644); err != nil {
		t.Fatalf("Failed to write main config: %v", err)
	}

	// Test ParseGitconfig directly with our test file
	cfg, err := ParseGitconfig(mainConfigPath)
	if err != nil {
		t.Fatalf("ParseGitconfig failed: %v", err)
	}

	// Check that we have 2 includeIf directives
	includeIfs := ExtractIncludeIfs(cfg)
	if len(includeIfs) != 2 {
		t.Fatalf("Expected 2 includeIf directives, got %d", len(includeIfs))
	}

	// Parse each included config
	var profiles []ProfileInfo
	for _, includeIf := range includeIfs {
		profile, err := ParseIncludedConfig(includeIf)
		if err != nil {
			t.Errorf("Failed to parse included config: %v", err)
			continue
		}
		profiles = append(profiles, ProfileInfo{
			Name:        profile.Name,
			Email:       profile.Email,
			SignCommits: profile.SignCommits,
		})
	}

	if len(profiles) != 2 {
		t.Fatalf("Expected 2 profiles, got %d", len(profiles))
	}

	// Verify profiles
	workFound := false
	personalFound := false
	for _, p := range profiles {
		if p.Email == "work@company.com" {
			workFound = true
			if p.SignCommits {
				t.Error("Work profile should not have signing enabled")
			}
		}
		if p.Email == "personal@gmail.com" {
			personalFound = true
			if !p.SignCommits {
				t.Error("Personal profile should have signing enabled")
			}
		}
	}

	if !workFound {
		t.Error("Work profile not found")
	}
	if !personalFound {
		t.Error("Personal profile not found")
	}
}

// ProfileInfo is a helper struct for tests
type ProfileInfo struct {
	Name        string
	Email       string
	SignCommits bool
}

func TestParseGitconfig_FileNotFound(t *testing.T) {
	_, err := ParseGitconfig("/nonexistent/path/.gitconfig")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestExtractIncludeIfs_NilConfig(t *testing.T) {
	result := ExtractIncludeIfs(nil)
	if result != nil {
		t.Error("Expected nil for nil config")
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input    string
		expected string
	}{
		{"~/test", filepath.Join(home, "test")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%s): expected '%s', got '%s'", tt.input, tt.expected, result)
			}
		})
	}
}
