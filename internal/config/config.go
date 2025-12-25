package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Version of the config format for future migrations
const ConfigVersion = "1.0.0"

// Config represents the gws workspace configuration
type Config struct {
	Version      string       `json:"version"`
	Workspace    string       `json:"workspace"`
	Repositories []Repository `json:"repositories"`
}

// Repository represents a discovered git repository
type Repository struct {
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	RemoteURL string   `json:"remote_url,omitempty"`
	Tags      []string `json:"tags,omitempty"`
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

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
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
