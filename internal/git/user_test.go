package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/daileyo/omgitworks/internal/config"
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

// writeConfigFile writes a config file, creating parent directories.
func writeConfigFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
}

// initBareTestRepo creates a repository with no commits and no user config.
func initBareTestRepo(t *testing.T, path string) {
	t.Helper()
	if _, err := git.PlainInit(path, false); err != nil {
		t.Fatalf("Failed to init test repo: %v", err)
	}
}

func TestParseConfigEntries(t *testing.T) {
	out := []byte("global\x00file:/h/.gitconfig\x00user.name\nJane Doe\x00" +
		"global\x00file:/h/.gitconfig\x00commit.gpgsign\x00" +
		"local\x00file:.git/config\x00user.email\njane@example.com\x00")

	entries := parseConfigEntries(out)
	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d: %+v", len(entries), entries)
	}

	want := []configEntry{
		{scope: "global", origin: "file:/h/.gitconfig", key: "user.name", value: "Jane Doe"},
		{scope: "global", origin: "file:/h/.gitconfig", key: "commit.gpgsign", bare: true},
		{scope: "local", origin: "file:.git/config", key: "user.email", value: "jane@example.com"},
	}
	for i, w := range want {
		if entries[i] != w {
			t.Errorf("Entry %d: expected %+v, got %+v", i, w, entries[i])
		}
	}

	if got := parseConfigEntries(nil); len(got) != 0 {
		t.Errorf("Expected no entries for empty output, got %+v", got)
	}
}

func TestEntryBool(t *testing.T) {
	tests := []struct {
		entry configEntry
		want  bool
	}{
		{configEntry{bare: true}, true},
		{configEntry{value: "true"}, true},
		{configEntry{value: "Yes"}, true},
		{configEntry{value: "on"}, true},
		{configEntry{value: "1"}, true},
		{configEntry{value: "false"}, false},
		{configEntry{value: "0"}, false},
		{configEntry{}, false},
	}

	for _, tt := range tests {
		if got := entryBool(tt.entry); got != tt.want {
			t.Errorf("entryBool(%+v) = %v, want %v", tt.entry, got, tt.want)
		}
	}
}

func TestResolveIncludePath(t *testing.T) {
	home := isolateGitConfig(t)

	tests := []struct {
		name   string
		path   string
		origin string
		want   string
	}{
		{"tilde", "~/.gitconfig-work", "file:/etc/gitconfig", filepath.Join(home, ".gitconfig-work")},
		{"absolute", "/opt/git/work.inc", "file:/etc/gitconfig", "/opt/git/work.inc"},
		{"relative to including file", "ids/work.inc", "file:/home/me/.config/git/config", "/home/me/.config/git/ids/work.inc"},
		{"relative without file origin", "ids/work.inc", "command line:", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveIncludePath(tt.path, tt.origin); got != tt.want {
				t.Errorf("resolveIncludePath(%q, %q) = %q, want %q", tt.path, tt.origin, got, tt.want)
			}
		})
	}
}

func TestGetUserConfig_HomeGitconfig(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), "[user]\n\tname = \"Home User\"\n\temail = home@example.com\n")

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.Name != "Home User" || userCfg.Email != "home@example.com" {
		t.Errorf("Expected Home User <home@example.com>, got %s <%s>", userCfg.Name, userCfg.Email)
	}
	if userCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected source 'global', got '%s'", userCfg.Source)
	}
}

// Regression: identity kept in git's XDG config was invisible to omgw.
func TestGetUserConfig_XDGGlobalConfig(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".config", "git", "config"), "[user]\n\tname = XDG User\n\temail = xdg@example.com\n")

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.Name != "XDG User" || userCfg.Email != "xdg@example.com" {
		t.Errorf("Expected XDG User <xdg@example.com>, got %s <%s>", userCfg.Name, userCfg.Email)
	}
	if userCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected source 'global', got '%s'", userCfg.Source)
	}
}

// Regression: a relative [include] path resolved against the working directory
// instead of the including file.
func TestGetUserConfig_RelativeInclude(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig-id"), "[user]\n\tname = Included User\n\temail = included@example.com\n")
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), "[include]\n\tpath = .gitconfig-id\n")

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.Email != "included@example.com" {
		t.Errorf("Expected email 'included@example.com', got '%s'", userCfg.Email)
	}
	if userCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected unconditional include to be 'global', got '%s'", userCfg.Source)
	}
}

func TestGetUserConfig_IncludeIf(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".config", "git", "work.inc"), "[user]\n\temail = work@company.com\n\tsigningkey = WORK123\n")
	writeConfigFile(t, filepath.Join(home, ".config", "git", "config"), `[user]
	name = Default User
	email = default@example.com
[includeIf "gitdir:~/work/"]
	path = work.inc
`)

	workRepo := filepath.Join(home, "work", "repo")
	otherRepo := filepath.Join(home, "other", "repo")
	initBareTestRepo(t, workRepo)
	initBareTestRepo(t, otherRepo)

	workCfg, err := GetUserConfig(workRepo)
	if err != nil {
		t.Fatalf("Failed to get work user config: %v", err)
	}
	if workCfg.Name != "Default User" || workCfg.Email != "work@company.com" || workCfg.SigningKey != "WORK123" {
		t.Errorf("Unexpected work identity: %+v", workCfg)
	}
	if workCfg.Source != config.UserSourceIncludeIf {
		t.Errorf("Expected source 'includeif', got '%s'", workCfg.Source)
	}

	otherCfg, err := GetUserConfig(otherRepo)
	if err != nil {
		t.Fatalf("Failed to get other user config: %v", err)
	}
	if otherCfg.Email != "default@example.com" {
		t.Errorf("Expected email 'default@example.com', got '%s'", otherCfg.Email)
	}
	if otherCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected source 'global', got '%s'", otherCfg.Source)
	}
}

func TestGetUserConfig_LocalSigningOverridesGlobal(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), "[user]\n\tname = Global User\n\temail = global@example.com\n[commit]\n\tgpgsign\n")

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if !userCfg.SignCommits {
		t.Error("Expected bare commit.gpgsign to enable signing")
	}

	writeConfigFile(t, filepath.Join(repoPath, ".git", "config"), "[commit]\n\tgpgsign = false\n")

	userCfg, err = GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.SignCommits {
		t.Error("Expected local commit.gpgsign=false to disable signing")
	}
	if userCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected identity source 'global', got '%s'", userCfg.Source)
	}
}

func TestGetUserConfig_NameWithoutEmail(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), "[user]\n\tname = Name Only\n")

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.Name != "Name Only" || userCfg.Email != "" {
		t.Errorf("Expected name-only identity, got %s <%s>", userCfg.Name, userCfg.Email)
	}
}

func TestGetUserConfig_NoIdentity(t *testing.T) {
	home := isolateGitConfig(t)

	repoPath := filepath.Join(home, "repo")
	initBareTestRepo(t, repoPath)

	userCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if userCfg.Source != config.UserSourceUnknown {
		t.Errorf("Expected source 'unknown', got '%s'", userCfg.Source)
	}
}

func TestGetNonLocalUserConfig_ReturnsUnderlyingGlobal(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".config", "git", "config"), "[user]\n\tname = Global User\n\temail = global@example.com\n")

	repoPath := filepath.Join(home, "repo")
	createTestRepoWithUser(t, repoPath, "Local User", "local@example.com", true)

	userCfg, err := GetNonLocalUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get non-local user config: %v", err)
	}
	if userCfg.Name != "Global User" || userCfg.Email != "global@example.com" {
		t.Errorf("Expected Global User <global@example.com>, got %s <%s>", userCfg.Name, userCfg.Email)
	}
	if userCfg.Source != config.UserSourceGlobal {
		t.Errorf("Expected source 'global', got '%s'", userCfg.Source)
	}
}

func TestGetGlobalDefaultUser_IgnoresIncludeIf(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig-work"), "[user]\n\temail = work@company.com\n")
	writeConfigFile(t, filepath.Join(home, ".gitconfig-default"), "[user]\n\tname = Default User\n\temail = default@example.com\n")
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), `[include]
	path = ~/.gitconfig-default
[includeIf "gitdir:/"]
	path = ~/.gitconfig-work
`)

	cfg, err := GetGlobalDefaultUser()
	if err != nil {
		t.Fatalf("GetGlobalDefaultUser failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("Expected a global default user")
	}
	if cfg.Name != "Default User" || cfg.Email != "default@example.com" {
		t.Errorf("Expected Default User <default@example.com>, got %s <%s>", cfg.Name, cfg.Email)
	}
}

func TestGetGlobalDefaultUser_NoIdentity(t *testing.T) {
	home := isolateGitConfig(t)
	writeConfigFile(t, filepath.Join(home, ".gitconfig"), "[core]\n\teditor = vim\n")

	cfg, err := GetGlobalDefaultUser()
	if err != nil {
		t.Fatalf("GetGlobalDefaultUser failed: %v", err)
	}
	if cfg != nil {
		t.Errorf("Expected nil config without an identity, got %+v", cfg)
	}
}

func TestListIncludeIfDirectives(t *testing.T) {
	home := isolateGitConfig(t)
	xdgConfig := filepath.Join(home, ".config", "git", "config")
	writeConfigFile(t, xdgConfig, `[includeIf "gitdir:~/Work/"]
	path = ids/work.inc
[includeIf "gitdir/i:~/personal/"]
	path = ~/.gitconfig-personal
`)

	directives, err := ListIncludeIfDirectives()
	if err != nil {
		t.Fatalf("ListIncludeIfDirectives failed: %v", err)
	}

	want := []IncludeIfDirective{
		{Condition: "gitdir:~/Work/", Path: filepath.Join(home, ".config", "git", "ids", "work.inc")},
		{Condition: "gitdir/i:~/personal/", Path: filepath.Join(home, ".gitconfig-personal")},
	}
	if len(directives) != len(want) {
		t.Fatalf("Expected %d directives, got %d: %+v", len(want), len(directives), directives)
	}
	for i, w := range want {
		if directives[i] != w {
			t.Errorf("Directive %d: expected %+v, got %+v", i, w, directives[i])
		}
	}
}

func TestReadUserConfigFile(t *testing.T) {
	home := isolateGitConfig(t)
	path := filepath.Join(home, "ids", "work.inc")
	writeConfigFile(t, filepath.Join(home, "ids", "signing.inc"), "[user]\n\tsigningkey = KEY123\n[commit]\n\tgpgsign\n")
	writeConfigFile(t, path, "[user]\n\tname = \"Work User\"\n\temail = work@company.com\n[include]\n\tpath = signing.inc\n")

	cfg, err := ReadUserConfigFile(path)
	if err != nil {
		t.Fatalf("ReadUserConfigFile failed: %v", err)
	}
	if cfg.Name != "Work User" || cfg.Email != "work@company.com" {
		t.Errorf("Expected Work User <work@company.com>, got %s <%s>", cfg.Name, cfg.Email)
	}
	if cfg.SigningKey != "KEY123" || !cfg.SignCommits {
		t.Errorf("Expected signing from nested include, got key=%q sign=%v", cfg.SigningKey, cfg.SignCommits)
	}

	if _, err := ReadUserConfigFile(filepath.Join(home, "missing.inc")); err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestGetNonLocalUserConfig_SkipsLocalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create repo with local user config
	createTestRepoWithUser(t, repoPath, "Local User", "local@example.com", true)

	// GetUserConfig should return local config
	localCfg, err := GetUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get user config: %v", err)
	}
	if localCfg.Source != config.UserSourceLocal {
		t.Fatalf("Expected source 'local', got '%s'", localCfg.Source)
	}

	// GetNonLocalUserConfig should skip local and return global/includeIf/unknown
	nonLocalCfg, err := GetNonLocalUserConfig(repoPath)
	if err != nil {
		t.Fatalf("Failed to get non-local user config: %v", err)
	}

	// Should NOT return local source
	if nonLocalCfg.Source == config.UserSourceLocal {
		t.Error("GetNonLocalUserConfig should not return local source")
	}

	// The non-local config should differ from local (unless global happens to match)
	// At minimum, it should not have the local name if global differs
	if nonLocalCfg.Source == config.UserSourceGlobal || nonLocalCfg.Source == config.UserSourceIncludeIf || nonLocalCfg.Source == config.UserSourceUnknown {
		// Valid non-local source
	} else {
		t.Errorf("Expected non-local source, got '%s'", nonLocalCfg.Source)
	}
}

func TestGetGlobalDefaultUser(t *testing.T) {
	// GetGlobalDefaultUser should return the global user from ~/.gitconfig
	// We can't control the system gitconfig in tests, so just verify it doesn't error
	cfg, err := GetGlobalDefaultUser()
	if err != nil {
		// Not an error if no global config exists
		t.Logf("GetGlobalDefaultUser returned error (may be expected): %v", err)
	}
	if cfg != nil {
		t.Logf("Global default user: name=%q email=%q", cfg.Name, cfg.Email)
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
