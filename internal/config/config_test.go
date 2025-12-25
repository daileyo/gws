package config

import (
	"encoding/json"
	"os"
	"path/filepath"
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
