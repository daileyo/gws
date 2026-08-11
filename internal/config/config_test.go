package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	workspacePath := "/test/workspace"
	cfg := New(workspacePath)

	if cfg.Version != ConfigVersion {
		t.Errorf("Expected version %s, got %s", ConfigVersion, cfg.Version)
	}

	if cfg.Workspace != workspacePath {
		t.Errorf("Expected workspace %s, got %s", workspacePath, cfg.Workspace)
	}

	if cfg.Repositories == nil {
		t.Error("Expected repositories slice to be initialized")
	}

	if len(cfg.Repositories) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(cfg.Repositories))
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory for test config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create test configuration
	cfg := &Config{
		Version:   ConfigVersion,
		Workspace: "/test/workspace",
		Repositories: []Repository{
			{
				Name:      "test-repo",
				Path:      "/test/workspace/test-repo",
				RemoteURL: "https://github.com/test/repo.git",
				Tags:      []string{"personal"},
			},
		},
	}

	// Save configuration
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load and verify
	loadedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var loadedCfg Config
	if err := json.Unmarshal(loadedData, &loadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loadedCfg.Version != cfg.Version {
		t.Errorf("Expected version %s, got %s", cfg.Version, loadedCfg.Version)
	}

	if loadedCfg.Workspace != cfg.Workspace {
		t.Errorf("Expected workspace %s, got %s", cfg.Workspace, loadedCfg.Workspace)
	}

	if len(loadedCfg.Repositories) != len(cfg.Repositories) {
		t.Errorf("Expected %d repositories, got %d", len(cfg.Repositories), len(loadedCfg.Repositories))
	}

	if len(loadedCfg.Repositories) > 0 {
		repo := loadedCfg.Repositories[0]
		if repo.Name != "test-repo" {
			t.Errorf("Expected repository name 'test-repo', got '%s'", repo.Name)
		}
		if repo.Path != "/test/workspace/test-repo" {
			t.Errorf("Expected repository path '/test/workspace/test-repo', got '%s'", repo.Path)
		}
		if repo.RemoteURL != "https://github.com/test/repo.git" {
			t.Errorf("Expected remote URL 'https://github.com/test/repo.git', got '%s'", repo.RemoteURL)
		}
		if len(repo.Tags) != 1 || repo.Tags[0] != "personal" {
			t.Errorf("Expected tags ['personal'], got %v", repo.Tags)
		}
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	if path == "" {
		t.Error("Config path should not be empty")
	}

	if !filepath.IsAbs(path) {
		t.Error("Config path should be absolute")
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected config filename 'config.json', got '%s'", filepath.Base(path))
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Failed to get config dir: %v", err)
	}

	if dir == "" {
		t.Error("Config dir should not be empty")
	}

	if !filepath.IsAbs(dir) {
		t.Error("Config dir should be absolute")
	}

	if filepath.Base(dir) != ".gws" {
		t.Errorf("Expected config dirname '.gws', got '%s'", filepath.Base(dir))
	}
}

func TestRepositoryStruct(t *testing.T) {
	repo := Repository{
		Name:      "test-repo",
		Path:      "/path/to/repo",
		RemoteURL: "https://github.com/test/repo.git",
		Tags:      []string{"work", "client"},
	}

	if repo.Name != "test-repo" {
		t.Errorf("Expected name 'test-repo', got '%s'", repo.Name)
	}

	if repo.Path != "/path/to/repo" {
		t.Errorf("Expected path '/path/to/repo', got '%s'", repo.Path)
	}

	if repo.RemoteURL != "https://github.com/test/repo.git" {
		t.Errorf("Expected remote URL 'https://github.com/test/repo.git', got '%s'", repo.RemoteURL)
	}

	if len(repo.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(repo.Tags))
	}
}

func TestConfigJSONMarshaling(t *testing.T) {
	cfg := &Config{
		Version:   "1.0.0",
		Workspace: "/test/workspace",
		Repositories: []Repository{
			{
				Name:      "repo1",
				Path:      "/test/workspace/repo1",
				RemoteURL: "https://github.com/test/repo1.git",
				Tags:      []string{"personal"},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal back
	var loadedCfg Config
	if err := json.Unmarshal(data, &loadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify fields
	if loadedCfg.Version != cfg.Version {
		t.Errorf("Version mismatch after marshal/unmarshal")
	}

	if loadedCfg.Workspace != cfg.Workspace {
		t.Errorf("Workspace mismatch after marshal/unmarshal")
	}

	if len(loadedCfg.Repositories) != len(cfg.Repositories) {
		t.Errorf("Repositories count mismatch after marshal/unmarshal")
	}
}

func TestProfileJSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
	}{
		{
			name: "full profile with signing",
			profile: Profile{
				Name:        "work",
				GitName:     "John Doe",
				Email:       "john@work.com",
				SigningKey:  "ABC123",
				SignCommits: true,
			},
		},
		{
			name: "profile without signing",
			profile: Profile{
				Name:        "personal",
				GitName:     "John Doe",
				Email:       "john@personal.com",
				SigningKey:  "",
				SignCommits: false,
			},
		},
		{
			name: "minimal profile",
			profile: Profile{
				Name:    "minimal",
				GitName: "Test User",
				Email:   "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.profile)
			if err != nil {
				t.Fatalf("Failed to marshal profile: %v", err)
			}

			// Unmarshal back
			var loaded Profile
			if err := json.Unmarshal(data, &loaded); err != nil {
				t.Fatalf("Failed to unmarshal profile: %v", err)
			}

			// Verify fields
			if loaded.Name != tt.profile.Name {
				t.Errorf("Name mismatch: expected %s, got %s", tt.profile.Name, loaded.Name)
			}
			if loaded.GitName != tt.profile.GitName {
				t.Errorf("GitName mismatch: expected %s, got %s", tt.profile.GitName, loaded.GitName)
			}
			if loaded.Email != tt.profile.Email {
				t.Errorf("Email mismatch: expected %s, got %s", tt.profile.Email, loaded.Email)
			}
			if loaded.SigningKey != tt.profile.SigningKey {
				t.Errorf("SigningKey mismatch: expected %s, got %s", tt.profile.SigningKey, loaded.SigningKey)
			}
			if loaded.SignCommits != tt.profile.SignCommits {
				t.Errorf("SignCommits mismatch: expected %v, got %v", tt.profile.SignCommits, loaded.SignCommits)
			}
		})
	}
}

func TestRepositoryWithUserFieldsMarshaling(t *testing.T) {
	tests := []struct {
		name string
		repo Repository
	}{
		{
			name: "repository with all user fields",
			repo: Repository{
				Name:           "test-repo",
				Path:           "/path/to/repo",
				RemoteURL:      "https://github.com/test/repo.git",
				Type:           TypeGitHub,
				Visibility:     VisibilityPrivate,
				Tags:           []string{"work"},
				User:           "John Doe",
				Email:          "john@work.com",
				SigningEnabled: true,
				UserSource:     UserSourceLocal,
			},
		},
		{
			name: "repository with global user source",
			repo: Repository{
				Name:           "another-repo",
				Path:           "/path/to/another",
				RemoteURL:      "https://gitlab.com/test/another.git",
				Type:           TypeGitLab,
				User:           "Jane Doe",
				Email:          "jane@example.com",
				SigningEnabled: false,
				UserSource:     UserSourceGlobal,
			},
		},
		{
			name: "repository with includeif user source",
			repo: Repository{
				Name:       "work-repo",
				Path:       "/path/to/work",
				User:       "Work User",
				Email:      "work@company.com",
				UserSource: UserSourceIncludeIf,
			},
		},
		{
			name: "repository without user fields (legacy)",
			repo: Repository{
				Name:      "legacy-repo",
				Path:      "/path/to/legacy",
				RemoteURL: "https://github.com/test/legacy.git",
				Tags:      []string{"old"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.repo)
			if err != nil {
				t.Fatalf("Failed to marshal repository: %v", err)
			}

			// Unmarshal back
			var loaded Repository
			if err := json.Unmarshal(data, &loaded); err != nil {
				t.Fatalf("Failed to unmarshal repository: %v", err)
			}

			// Verify user fields
			if loaded.User != tt.repo.User {
				t.Errorf("User mismatch: expected %s, got %s", tt.repo.User, loaded.User)
			}
			if loaded.Email != tt.repo.Email {
				t.Errorf("Email mismatch: expected %s, got %s", tt.repo.Email, loaded.Email)
			}
			if loaded.SigningEnabled != tt.repo.SigningEnabled {
				t.Errorf("SigningEnabled mismatch: expected %v, got %v", tt.repo.SigningEnabled, loaded.SigningEnabled)
			}
			if loaded.UserSource != tt.repo.UserSource {
				t.Errorf("UserSource mismatch: expected %s, got %s", tt.repo.UserSource, loaded.UserSource)
			}
		})
	}
}

func TestConfigWithProfilesMarshaling(t *testing.T) {
	cfg := &Config{
		Version:   ConfigVersion,
		Workspace: "/test/workspace",
		Profiles: []Profile{
			{
				Name:        "work",
				GitName:     "John Doe",
				Email:       "john@work.com",
				SigningKey:  "WORK123",
				SignCommits: true,
			},
			{
				Name:    "personal",
				GitName: "John Doe",
				Email:   "john@personal.com",
			},
		},
		Repositories: []Repository{
			{
				Name:           "work-repo",
				Path:           "/test/workspace/work-repo",
				User:           "John Doe",
				Email:          "john@work.com",
				SigningEnabled: true,
				UserSource:     UserSourceIncludeIf,
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal back
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify profiles
	if len(loaded.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(loaded.Profiles))
	}

	if len(loaded.Profiles) > 0 {
		if loaded.Profiles[0].Name != "work" {
			t.Errorf("Expected first profile name 'work', got '%s'", loaded.Profiles[0].Name)
		}
		if loaded.Profiles[0].Email != "john@work.com" {
			t.Errorf("Expected first profile email 'john@work.com', got '%s'", loaded.Profiles[0].Email)
		}
		if !loaded.Profiles[0].SignCommits {
			t.Error("Expected first profile to have SignCommits enabled")
		}
	}

	// Verify repository user fields
	if len(loaded.Repositories) > 0 {
		repo := loaded.Repositories[0]
		if repo.User != "John Doe" {
			t.Errorf("Expected repo user 'John Doe', got '%s'", repo.User)
		}
		if repo.UserSource != UserSourceIncludeIf {
			t.Errorf("Expected repo user source 'includeif', got '%s'", repo.UserSource)
		}
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Simulate loading a config file from older version (without profiles or user fields)
	oldConfigJSON := `{
		"version": "1.0.0",
		"workspace": "/test/workspace",
		"repositories": [
			{
				"name": "old-repo",
				"path": "/test/workspace/old-repo",
				"remote_url": "https://github.com/test/old.git",
				"type": "github",
				"visibility": "private",
				"tags": ["legacy"]
			}
		]
	}`

	var cfg Config
	if err := json.Unmarshal([]byte(oldConfigJSON), &cfg); err != nil {
		t.Fatalf("Failed to unmarshal old config format: %v", err)
	}

	// Verify old fields are loaded correctly
	if cfg.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", cfg.Version)
	}
	if cfg.Workspace != "/test/workspace" {
		t.Errorf("Expected workspace '/test/workspace', got '%s'", cfg.Workspace)
	}
	if len(cfg.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(cfg.Repositories))
	}

	// Verify profiles defaults to nil/empty
	if len(cfg.Profiles) != 0 {
		t.Errorf("Expected profiles to be nil/empty for old config, got %v", cfg.Profiles)
	}

	// Verify user fields default to zero values
	if len(cfg.Repositories) > 0 {
		repo := cfg.Repositories[0]
		if repo.User != "" {
			t.Errorf("Expected User to be empty for old config, got '%s'", repo.User)
		}
		if repo.Email != "" {
			t.Errorf("Expected Email to be empty for old config, got '%s'", repo.Email)
		}
		if repo.SigningEnabled {
			t.Error("Expected SigningEnabled to be false for old config")
		}
		if repo.UserSource != "" {
			t.Errorf("Expected UserSource to be empty for old config, got '%s'", repo.UserSource)
		}
	}
}

func TestEffectiveScanMaxDepth(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want int
	}{
		{
			name: "nil config falls back to default",
			cfg:  nil,
			want: DefaultScanMaxDepth,
		},
		{
			name: "nil preferences falls back to default",
			cfg:  &Config{Workspace: "/test"},
			want: DefaultScanMaxDepth,
		},
		{
			name: "preferences present but field unset falls back to default",
			cfg:  &Config{Preferences: &Preferences{StatusWorkers: 4}},
			want: DefaultScanMaxDepth,
		},
		{
			name: "explicit value is used",
			cfg:  &Config{Preferences: &Preferences{ScanMaxDepth: 3}},
			want: 3,
		},
		{
			name: "zero falls back to default",
			cfg:  &Config{Preferences: &Preferences{ScanMaxDepth: 0}},
			want: DefaultScanMaxDepth,
		},
		{
			name: "negative falls back to default",
			cfg:  &Config{Preferences: &Preferences{ScanMaxDepth: -1}},
			want: DefaultScanMaxDepth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.EffectiveScanMaxDepth(); got != tt.want {
				t.Errorf("EffectiveScanMaxDepth() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestScanMaxDepthBackwardCompatibility(t *testing.T) {
	t.Run("config without preferences key defaults to 6", func(t *testing.T) {
		oldConfigJSON := `{
			"version": "1.1.0",
			"workspace": "/test/workspace",
			"repositories": []
		}`

		var cfg Config
		if err := json.Unmarshal([]byte(oldConfigJSON), &cfg); err != nil {
			t.Fatalf("Failed to unmarshal config without preferences: %v", err)
		}

		if cfg.Preferences != nil {
			t.Errorf("Expected Preferences to be nil, got %+v", cfg.Preferences)
		}
		if got := cfg.EffectiveScanMaxDepth(); got != DefaultScanMaxDepth {
			t.Errorf("Expected default depth %d, got %d", DefaultScanMaxDepth, got)
		}
	})

	t.Run("config with existing preferences but no scan_max_depth defaults to 6", func(t *testing.T) {
		configJSON := `{
			"version": "1.1.0",
			"workspace": "/test/workspace",
			"preferences": {"status_workers": 8},
			"repositories": []
		}`

		var cfg Config
		if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
			t.Fatalf("Failed to unmarshal config: %v", err)
		}

		if cfg.Preferences == nil || cfg.Preferences.StatusWorkers != 8 {
			t.Fatalf("Expected StatusWorkers 8 to survive, got %+v", cfg.Preferences)
		}
		if got := cfg.EffectiveScanMaxDepth(); got != DefaultScanMaxDepth {
			t.Errorf("Expected default depth %d, got %d", DefaultScanMaxDepth, got)
		}
	})

	t.Run("scan_max_depth round-trips and is omitted when zero", func(t *testing.T) {
		cfg := &Config{
			Version:      ConfigVersion,
			Workspace:    "/test/workspace",
			Preferences:  &Preferences{ScanMaxDepth: 4},
			Repositories: []Repository{},
		}

		data, err := json.Marshal(cfg)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}
		if !strings.Contains(string(data), `"scan_max_depth":4`) {
			t.Errorf("Expected scan_max_depth in JSON, got %s", string(data))
		}

		var decoded Config
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal config: %v", err)
		}
		if got := decoded.EffectiveScanMaxDepth(); got != 4 {
			t.Errorf("Expected depth 4 after round-trip, got %d", got)
		}

		zeroCfg := &Config{Preferences: &Preferences{StatusWorkers: 2}}
		zeroData, err := json.Marshal(zeroCfg)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}
		if strings.Contains(string(zeroData), "scan_max_depth") {
			t.Errorf("Expected scan_max_depth to be omitted when zero, got %s", string(zeroData))
		}
	})
}

func TestRepositoryWithWorktreesMarshaling(t *testing.T) {
	t.Run("worktrees round-trip through JSON", func(t *testing.T) {
		repo := Repository{
			Name: "my-repo",
			Path: "/workspace/my-repo",
			Worktrees: []Worktree{
				{Path: "/workspace/my-repo.wt/feature-x", Branch: "feature-x", Aligned: true},
				{Path: "/tmp/other-worktree", Branch: "hotfix", Aligned: false},
			},
		}

		data, err := json.Marshal(repo)
		if err != nil {
			t.Fatalf("Failed to marshal repository: %v", err)
		}

		var loaded Repository
		if err := json.Unmarshal(data, &loaded); err != nil {
			t.Fatalf("Failed to unmarshal repository: %v", err)
		}

		if len(loaded.Worktrees) != 2 {
			t.Fatalf("Expected 2 worktrees, got %d", len(loaded.Worktrees))
		}
		if loaded.Worktrees[0].Path != "/workspace/my-repo.wt/feature-x" {
			t.Errorf("Expected worktree path '/workspace/my-repo.wt/feature-x', got '%s'", loaded.Worktrees[0].Path)
		}
		if loaded.Worktrees[0].Branch != "feature-x" {
			t.Errorf("Expected worktree branch 'feature-x', got '%s'", loaded.Worktrees[0].Branch)
		}
		if !loaded.Worktrees[0].Aligned {
			t.Error("Expected first worktree to be aligned")
		}
		if loaded.Worktrees[1].Aligned {
			t.Error("Expected second worktree to not be aligned")
		}
	})

	t.Run("worktrees omitted from JSON when empty", func(t *testing.T) {
		repo := Repository{
			Name: "no-wt-repo",
			Path: "/workspace/no-wt-repo",
		}

		data, err := json.Marshal(repo)
		if err != nil {
			t.Fatalf("Failed to marshal repository: %v", err)
		}

		jsonStr := string(data)
		if contains := "worktrees"; len(jsonStr) > 0 {
			for i := 0; i <= len(jsonStr)-len(contains); i++ {
				if jsonStr[i:i+len(contains)] == contains {
					t.Error("Expected 'worktrees' key to be omitted from JSON when empty")
					break
				}
			}
		}
	})

	t.Run("backward compat - old config without worktrees loads cleanly", func(t *testing.T) {
		oldJSON := `{"name":"old-repo","path":"/workspace/old-repo"}`

		var loaded Repository
		if err := json.Unmarshal([]byte(oldJSON), &loaded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if loaded.Worktrees != nil {
			t.Errorf("Expected nil worktrees for old config, got %v", loaded.Worktrees)
		}
	})
}

func TestUserSourceConstants(t *testing.T) {
	tests := []struct {
		source   UserSource
		expected string
	}{
		{UserSourceGlobal, "global"},
		{UserSourceLocal, "local"},
		{UserSourceIncludeIf, "includeif"},
		{UserSourceUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.source), func(t *testing.T) {
			if string(tt.source) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.source))
			}
		})
	}
}
