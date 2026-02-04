package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/daileyo/gws/internal/config"
)

// UserConfig represents the git user configuration for a repository
type UserConfig struct {
	Name        string            // user.name value
	Email       string            // user.email value
	SigningKey  string            // user.signingkey value
	SignCommits bool              // commit.gpgsign setting
	Source      config.UserSource // Where the config comes from
}

// GetUserConfig reads the effective git user configuration for a repository.
// It checks the local .git/config first, then falls back to global config.
func GetUserConfig(repoPath string) (*UserConfig, error) {
	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	userConfig := &UserConfig{
		Source: config.UserSourceUnknown,
	}

	// Try to get local repository config first
	localCfg, err := repo.Config()
	if err == nil && localCfg != nil {
		// Check if user is configured locally
		if localCfg.User.Name != "" || localCfg.User.Email != "" {
			userConfig.Name = localCfg.User.Name
			userConfig.Email = localCfg.User.Email
			userConfig.Source = config.UserSourceLocal
		}
	}

	// If no local user config, try global config
	if userConfig.Source == config.UserSourceUnknown {
		globalCfg, err := loadGlobalConfig()
		if err == nil && globalCfg != nil {
			userConfig.Name = globalCfg.Name
			userConfig.Email = globalCfg.Email
			userConfig.SigningKey = globalCfg.SigningKey
			userConfig.SignCommits = globalCfg.SignCommits
			userConfig.Source = config.UserSourceGlobal
		}
	}

	// Check for signing configuration in local config
	// (may override global signing settings)
	if localCfg != nil {
		signingKey, signCommits := getSigningFromRawConfig(repoPath)
		if signingKey != "" {
			userConfig.SigningKey = signingKey
		}
		if signCommits {
			userConfig.SignCommits = signCommits
		}
	}

	// If still no user found, check if this might be an includeIf scenario
	// by checking if the path matches common patterns
	if userConfig.Source == config.UserSourceGlobal {
		if isLikelyIncludeIfPath(repoPath) {
			userConfig.Source = config.UserSourceIncludeIf
		}
	}

	return userConfig, nil
}

// GlobalUserConfig represents user config from global gitconfig
type GlobalUserConfig struct {
	Name        string
	Email       string
	SigningKey  string
	SignCommits bool
}

// loadGlobalConfig reads the global git configuration
func loadGlobalConfig() (*GlobalUserConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	gitconfigPath := filepath.Join(home, ".gitconfig")

	// Read and parse the gitconfig file, following includes
	return parseGitConfigWithIncludes(gitconfigPath, home)
}

// parseGitConfigWithIncludes parses a gitconfig file and follows [include] directives
func parseGitConfigWithIncludes(configPath, home string) (*GlobalUserConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No global config
		}
		return nil, fmt.Errorf("failed to read gitconfig: %w", err)
	}

	cfg, includes := parseGitConfigAndIncludes(string(data))

	// Follow include directives to find user info
	for _, includePath := range includes {
		// Expand ~ to home directory
		if strings.HasPrefix(includePath, "~/") {
			includePath = filepath.Join(home, includePath[2:])
		}

		includeData, err := os.ReadFile(includePath)
		if err != nil {
			continue // Skip includes that can't be read
		}

		includeCfg, _ := parseGitConfigAndIncludes(string(includeData))

		// Merge: included config values take precedence if not already set
		if cfg.Name == "" && includeCfg.Name != "" {
			cfg.Name = includeCfg.Name
		}
		if cfg.Email == "" && includeCfg.Email != "" {
			cfg.Email = includeCfg.Email
		}
		if cfg.SigningKey == "" && includeCfg.SigningKey != "" {
			cfg.SigningKey = includeCfg.SigningKey
		}
		if !cfg.SignCommits && includeCfg.SignCommits {
			cfg.SignCommits = includeCfg.SignCommits
		}
	}

	return cfg, nil
}

// parseGitConfigAndIncludes parses a gitconfig file content and extracts user info and include paths
func parseGitConfigAndIncludes(content string) (*GlobalUserConfig, []string) {
	cfg := &GlobalUserConfig{}
	var includes []string

	lines := strings.Split(content, "\n")
	inUserSection := false
	inCommitSection := false
	inIncludeSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.ToLower(strings.Trim(line, "[]"))
			inUserSection = section == "user"
			inCommitSection = section == "commit"
			inIncludeSection = section == "include"
			continue
		}

		// Parse key-value pairs
		if inUserSection {
			if strings.HasPrefix(strings.ToLower(line), "name") {
				cfg.Name = extractValue(line)
			} else if strings.HasPrefix(strings.ToLower(line), "email") {
				cfg.Email = extractValue(line)
			} else if strings.HasPrefix(strings.ToLower(line), "signingkey") {
				cfg.SigningKey = extractValue(line)
			}
		}

		if inCommitSection {
			if strings.HasPrefix(strings.ToLower(line), "gpgsign") {
				value := strings.ToLower(extractValue(line))
				cfg.SignCommits = value == "true" || value == "1" || value == "yes"
			}
		}

		if inIncludeSection {
			if strings.HasPrefix(strings.ToLower(line), "path") {
				includes = append(includes, extractValue(line))
			}
		}
	}

	return cfg, includes
}

// parseGitConfig parses a gitconfig file content and extracts user info (for backward compatibility)
func parseGitConfig(content string) *GlobalUserConfig {
	cfg, _ := parseGitConfigAndIncludes(content)
	return cfg
}

// extractValue extracts the value from a "key = value" line
func extractValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return ""
	}
	value := strings.TrimSpace(parts[1])
	// Remove quotes if present
	value = strings.Trim(value, "\"'")
	return value
}

// getSigningFromRawConfig reads signing configuration from repo's .git/config
func getSigningFromRawConfig(repoPath string) (signingKey string, signCommits bool) {
	configPath := filepath.Join(repoPath, ".git", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", false
	}

	cfg := parseGitConfig(string(data))
	return cfg.SigningKey, cfg.SignCommits
}

// isLikelyIncludeIfPath checks if the repo path matches common includeIf patterns
// This is a heuristic - we look for paths that contain directory patterns
// commonly used with includeIf (e.g., ~/work/, ~/personal/, etc.)
func isLikelyIncludeIfPath(repoPath string) bool {
	// Common directory patterns that suggest includeIf usage
	patterns := []string{
		"/work/",
		"/personal/",
		"/ado/",
		"/gh/",
		"/github/",
		"/gitlab/",
		"/bitbucket/",
		"/horizon/",
		"/gws/",
	}

	lowerPath := strings.ToLower(repoPath)
	for _, pattern := range patterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}

	return false
}

// GetEffectiveUserConfig reads the effective git user for a repo by running git config
// This is a more accurate method that respects all includeIf directives
func GetEffectiveUserConfig(repoPath string) (*UserConfig, error) {
	// First try the standard method
	cfg, err := GetUserConfig(repoPath)
	if err != nil {
		return nil, err
	}

	// If we detected a user, return it
	if cfg.Name != "" || cfg.Email != "" {
		return cfg, nil
	}

	// Fallback: return empty config with unknown source
	return &UserConfig{
		Source: config.UserSourceUnknown,
	}, nil
}
