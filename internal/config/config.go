package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Version of the config format for future migrations
const ConfigVersion = "1.1.0"

// DefaultScanMaxDepth is the number of directory levels below the workspace
// root that repository discovery traverses when no preference is set.
// A repository found at exactly this depth is registered; nothing below it
// is traversed.
const DefaultScanMaxDepth = 6

// Config represents the gws workspace configuration
type Config struct {
	Version      string       `json:"version"`
	Workspace    string       `json:"workspace"`
	Profiles     []Profile    `json:"profiles,omitempty"`
	Repositories []Repository `json:"repositories"`
	Preferences  *Preferences `json:"preferences,omitempty"`
}

// Preferences holds user-configurable settings
type Preferences struct {
	StatusWorkers int `json:"status_workers,omitempty"`
	ScanMaxDepth  int `json:"scan_max_depth,omitempty"`
}

// EffectiveScanMaxDepth returns the configured scan depth, falling back to
// DefaultScanMaxDepth when unset or non-positive.
func (c *Config) EffectiveScanMaxDepth() int {
	if c != nil && c.Preferences != nil && c.Preferences.ScanMaxDepth > 0 {
		return c.Preferences.ScanMaxDepth
	}
	return DefaultScanMaxDepth
}

// RepositoryType represents the type of git hosting provider
type RepositoryType string

const (
	TypeGitHub    RepositoryType = "github"
	TypeGitLab    RepositoryType = "gitlab"
	TypeADO       RepositoryType = "ado"
	TypeBitbucket RepositoryType = "bitbucket"
	TypeUnknown   RepositoryType = "unknown"
)

// RepositoryVisibility represents whether a repository is public or private
type RepositoryVisibility string

const (
	VisibilityPublic  RepositoryVisibility = "public"
	VisibilityPrivate RepositoryVisibility = "private"
	VisibilityUnknown RepositoryVisibility = "unknown"
)

// UserSource represents where the git user configuration comes from
type UserSource string

const (
	UserSourceGlobal    UserSource = "global"
	UserSourceLocal     UserSource = "local"
	UserSourceIncludeIf UserSource = "includeif"
	UserSourceUnknown   UserSource = "unknown"
)

// Profile represents a git user profile with identity and signing configuration
type Profile struct {
	Name        string `json:"name"`                   // Profile identifier (e.g., "work", "personal")
	GitName     string `json:"git_name"`               // user.name value
	Email       string `json:"email"`                  // user.email value
	SigningKey  string `json:"signing_key,omitempty"`  // user.signingkey value (optional)
	SignCommits bool   `json:"sign_commits,omitempty"` // commit.gpgsign setting
}

// Worktree represents a git worktree associated with a repository
type Worktree struct {
	Path    string `json:"path"`              // Absolute filesystem path to the worktree
	Branch  string `json:"branch"`            // Branch checked out in the worktree
	Aligned bool   `json:"aligned,omitempty"` // Whether the worktree is inside <repo>.wt/
}

// Repository represents a discovered git repository
type Repository struct {
	Name           string               `json:"name"`
	Path           string               `json:"path"`
	RemoteURL      string               `json:"remote_url,omitempty"`
	Type           RepositoryType       `json:"type,omitempty"`
	Visibility     RepositoryVisibility `json:"visibility,omitempty"`
	Tags           []string             `json:"tags,omitempty"`
	User           string               `json:"user,omitempty"`            // Git user.name for this repo
	Email          string               `json:"email,omitempty"`           // Git user.email for this repo
	SigningEnabled bool                 `json:"signing_enabled,omitempty"` // Whether commit signing is configured
	UserSource     UserSource           `json:"user_source,omitempty"`     // Where the user config comes from
	Worktrees      []Worktree           `json:"worktrees,omitempty"`       // Git worktrees for this repo
}

// GetConfigPath returns the path to the gws config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".gws", "config.json"), nil
}

// GetConfigDir returns the directory containing the config file
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".gws"), nil
}

// Load reads the configuration from ~/.gws/config.json
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workspace not initialized: run 'gws init <directory>' first")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to ~/.gws/config.json
func Save(cfg *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with secure permissions (user read/write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// New creates a new configuration with the specified workspace path
func New(workspacePath string) *Config {
	return &Config{
		Version:      ConfigVersion,
		Workspace:    workspacePath,
		Repositories: []Repository{},
	}
}

// Exists checks if a configuration file exists
func Exists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
