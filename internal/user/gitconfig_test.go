package user

import (
	"os"
	"path/filepath"
	"testing"
)

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

// isolateGitConfig points git at an empty home directory so a test sees only
// the config it writes, never the developer's own. Returns the home directory.
func isolateGitConfig(t *testing.T) string {
	t.Helper()

	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to resolve temp home: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	// An empty GIT_CONFIG_GLOBAL disables global config entirely, so unset it.
	t.Setenv("GIT_CONFIG_GLOBAL", "")
	os.Unsetenv("GIT_CONFIG_GLOBAL")
	return home
}

func TestDetectProfiles_Integration(t *testing.T) {
	home := isolateGitConfig(t)

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
	if err := os.WriteFile(filepath.Join(home, ".gitconfig-work"), []byte(workConfig), 0644); err != nil {
		t.Fatalf("Failed to write work config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".gitconfig-personal"), []byte(personalConfig), 0644); err != nil {
		t.Fatalf("Failed to write personal config: %v", err)
	}

	// A relative path resolves against the including file, as git does
	mainConfig := `[user]
	name = Default User
	email = default@example.com

[includeIf "gitdir:~/work/"]
	path = .gitconfig-work

[includeIf "gitdir:~/personal/"]
	path = ~/.gitconfig-personal
`
	if err := os.WriteFile(filepath.Join(home, ".gitconfig"), []byte(mainConfig), 0644); err != nil {
		t.Fatalf("Failed to write main config: %v", err)
	}

	profiles, err := DetectProfiles()
	if err != nil {
		t.Fatalf("DetectProfiles failed: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("Expected 2 profiles, got %d: %+v", len(profiles), profiles)
	}

	if profiles[0].Name != "work" || profiles[0].Email != "work@company.com" || profiles[0].SignCommits {
		t.Errorf("Unexpected work profile: %+v", profiles[0])
	}
	if profiles[1].Name != "personal" || profiles[1].Email != "personal@gmail.com" ||
		profiles[1].SigningKey != "PERSONAL123" || !profiles[1].SignCommits {
		t.Errorf("Unexpected personal profile: %+v", profiles[1])
	}
}

// Regression: includeIf directives in git's XDG config were never detected.
func TestDetectProfiles_XDGConfig(t *testing.T) {
	home := isolateGitConfig(t)
	gitDir := filepath.Join(home, ".config", "git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create git config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "gitconfig-work"), []byte("[user]\n\tname = Work User\n\temail = work@company.com\n"), 0644); err != nil {
		t.Fatalf("Failed to write work config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("[includeIf \"gitdir:~/work/\"]\n\tpath = gitconfig-work\n"), 0644); err != nil {
		t.Fatalf("Failed to write XDG config: %v", err)
	}

	profiles, err := DetectProfiles()
	if err != nil {
		t.Fatalf("DetectProfiles failed: %v", err)
	}
	if len(profiles) != 1 || profiles[0].Email != "work@company.com" || profiles[0].Name != "work" {
		t.Errorf("Expected work profile from XDG config, got %+v", profiles)
	}
}

func TestDetectProfiles_NoConfig(t *testing.T) {
	isolateGitConfig(t)

	profiles, err := DetectProfiles()
	if err != nil {
		t.Fatalf("DetectProfiles failed: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("Expected no profiles, got %+v", profiles)
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
