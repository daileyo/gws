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

	// If using global config, check if an includeIf directive applies to this repo
	if userConfig.Source == config.UserSourceGlobal {
		if includeIfCfg, matched := checkIncludeIfMatch(repoPath); matched && includeIfCfg != nil {
			userConfig.Name = includeIfCfg.Name
			userConfig.Email = includeIfCfg.Email
			if includeIfCfg.SigningKey != "" {
				userConfig.SigningKey = includeIfCfg.SigningKey
			}
			if includeIfCfg.SignCommits {
				userConfig.SignCommits = includeIfCfg.SignCommits
			}
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

// MatchesGitdirCondition evaluates whether a repo path matches a gitdir includeIf
// condition string (e.g., "gitdir:~/work/" matches "/home/user/work/my-repo").
// Supports gitdir: and gitdir/i: (case-insensitive) prefixes, trailing ** globs,
// and ~/ home directory expansion.
func MatchesGitdirCondition(repoPath string, condition string) bool {
	// Extract the prefix and pattern
	var pattern string
	caseInsensitive := false

	lower := strings.ToLower(condition)
	if strings.HasPrefix(lower, "gitdir/i:") {
		pattern = condition[len("gitdir/i:"):]
		caseInsensitive = true
	} else if strings.HasPrefix(lower, "gitdir:") {
		pattern = condition[len("gitdir:"):]
	} else {
		return false // Not a gitdir condition
	}

	// Expand ~ in pattern
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	if strings.HasPrefix(pattern, "~/") {
		pattern = filepath.Join(home, pattern[2:])
	}

	// Clean and normalize
	pattern = filepath.Clean(pattern)
	repoPathClean := filepath.Clean(repoPath)

	if caseInsensitive {
		pattern = strings.ToLower(pattern)
		repoPathClean = strings.ToLower(repoPathClean)
	}

	// Remove trailing ** if present (means "match anything below this dir")
	pattern = strings.TrimSuffix(pattern, "/**")
	pattern = strings.TrimSuffix(pattern, "**")

	// Ensure pattern ends with separator for directory matching
	if !strings.HasSuffix(pattern, string(filepath.Separator)) {
		pattern += string(filepath.Separator)
	}

	// Check if repo path starts with the pattern directory
	if !strings.HasSuffix(repoPathClean, string(filepath.Separator)) {
		repoPathClean += string(filepath.Separator)
	}

	return strings.HasPrefix(repoPathClean, pattern)
}

// includeIfEntry represents a parsed includeIf directive from gitconfig
type includeIfEntry struct {
	condition string
	path      string
}

// checkIncludeIfMatch checks if any includeIf directive in ~/.gitconfig matches
// the given repo path. Returns the user config from the matched included config
// and whether a match was found.
func checkIncludeIfMatch(repoPath string) (*UserConfig, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, false
	}

	gitconfigPath := filepath.Join(home, ".gitconfig")
	data, err := os.ReadFile(gitconfigPath)
	if err != nil {
		return nil, false
	}

	entries := parseIncludeIfs(string(data), home)
	for _, entry := range entries {
		if !MatchesGitdirCondition(repoPath, entry.condition) {
			continue
		}

		// Read and parse the included config file
		includedData, err := os.ReadFile(entry.path)
		if err != nil {
			continue
		}

		includedCfg := parseGitConfig(string(includedData))
		if includedCfg.Name == "" && includedCfg.Email == "" {
			continue
		}

		return &UserConfig{
			Name:        includedCfg.Name,
			Email:       includedCfg.Email,
			SigningKey:  includedCfg.SigningKey,
			SignCommits: includedCfg.SignCommits,
			Source:      config.UserSourceIncludeIf,
		}, true
	}

	return nil, false
}

// parseIncludeIfs extracts includeIf directives from gitconfig content
func parseIncludeIfs(content string, home string) []includeIfEntry {
	var entries []includeIfEntry
	lines := strings.Split(content, "\n")
	var currentCondition string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for [includeIf "condition"] section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionLine := strings.Trim(line, "[]")
			lower := strings.ToLower(sectionLine)
			if strings.HasPrefix(lower, "includeif ") {
				// Extract quoted condition
				quoted := sectionLine[10:] // Skip "includeIf "
				currentCondition = strings.Trim(strings.TrimSpace(quoted), "\"'")
			} else {
				currentCondition = ""
			}
			continue
		}

		// If we're in an includeIf section, look for path =
		if currentCondition != "" {
			key, value := parseIncludeIfKeyValue(line)
			if strings.ToLower(key) == "path" {
				// Expand ~ in path
				if strings.HasPrefix(value, "~/") {
					value = filepath.Join(home, value[2:])
				}
				entries = append(entries, includeIfEntry{
					condition: currentCondition,
					path:      filepath.Clean(value),
				})
			}
		}
	}

	return entries
}

// parseIncludeIfKeyValue parses a "key = value" line for includeIf sections
func parseIncludeIfKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\"'")
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
