package user

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daileyo/gws/internal/config"
)

// GitConfig represents a parsed gitconfig file
type GitConfig struct {
	Path       string            // Path to the gitconfig file
	Sections   map[string]Section
	IncludeIfs []IncludeIf
}

// Section represents a gitconfig section with key-value pairs
type Section struct {
	Name   string
	Values map[string]string
}

// IncludeIf represents an includeIf directive in gitconfig
type IncludeIf struct {
	Condition string // The condition (e.g., "gitdir:~/work/")
	Path      string // The path to the included config file
}

// ParseGitconfig parses a gitconfig file and returns its structure
func ParseGitconfig(path string) (*GitConfig, error) {
	// Expand ~ to home directory
	expandedPath := expandPath(path)

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gitconfig: %w", err)
	}

	return parseGitconfigContent(string(data), expandedPath)
}

// parseGitconfigContent parses gitconfig content and extracts sections and includeIfs
func parseGitconfigContent(content, path string) (*GitConfig, error) {
	cfg := &GitConfig{
		Path:       path,
		Sections:   make(map[string]Section),
		IncludeIfs: []IncludeIf{},
	}

	lines := strings.Split(content, "\n")
	var currentSection string
	var currentIncludeIf *IncludeIf

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionLine := strings.Trim(line, "[]")

			// Check if this is an includeIf section
			if strings.HasPrefix(strings.ToLower(sectionLine), "includeif ") {
				condition := extractQuotedValue(sectionLine[10:]) // Skip "includeIf "
				currentIncludeIf = &IncludeIf{
					Condition: condition,
				}
				currentSection = ""
			} else {
				// Regular section
				currentSection = strings.ToLower(sectionLine)
				if _, exists := cfg.Sections[currentSection]; !exists {
					cfg.Sections[currentSection] = Section{
						Name:   currentSection,
						Values: make(map[string]string),
					}
				}
				currentIncludeIf = nil
			}
			continue
		}

		// Parse key-value pairs
		if currentIncludeIf != nil {
			// We're in an includeIf section
			key, value := parseKeyValue(line)
			if strings.ToLower(key) == "path" {
				currentIncludeIf.Path = expandPath(value)
				cfg.IncludeIfs = append(cfg.IncludeIfs, *currentIncludeIf)
			}
		} else if currentSection != "" {
			// Regular section
			key, value := parseKeyValue(line)
			if key != "" {
				section := cfg.Sections[currentSection]
				section.Values[strings.ToLower(key)] = value
				cfg.Sections[currentSection] = section
			}
		}
	}

	return cfg, nil
}

// ExtractIncludeIfs returns all includeIf directives from a GitConfig
func ExtractIncludeIfs(cfg *GitConfig) []IncludeIf {
	if cfg == nil {
		return nil
	}
	return cfg.IncludeIfs
}

// ParseIncludedConfig reads a gitconfig file referenced by an includeIf and extracts a Profile
func ParseIncludedConfig(includeIf IncludeIf) (*config.Profile, error) {
	expandedPath := expandPath(includeIf.Path)

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read included config: %w", err)
	}

	gitCfg, err := parseGitconfigContent(string(data), expandedPath)
	if err != nil {
		return nil, err
	}

	// Extract user section values
	userSection, hasUser := gitCfg.Sections["user"]
	if !hasUser {
		return nil, fmt.Errorf("no [user] section found in %s", expandedPath)
	}

	// Derive profile name from the includeIf condition or file path
	profileName := deriveProfileName(includeIf)

	profile := &config.Profile{
		Name:    profileName,
		GitName: userSection.Values["name"],
		Email:   userSection.Values["email"],
	}

	// Check for signing configuration
	if signingKey, ok := userSection.Values["signingkey"]; ok {
		profile.SigningKey = signingKey
	}

	// Check commit section for gpgsign
	if commitSection, hasCommit := gitCfg.Sections["commit"]; hasCommit {
		if gpgsign, ok := commitSection.Values["gpgsign"]; ok {
			profile.SignCommits = gpgsign == "true" || gpgsign == "1" || gpgsign == "yes"
		}
	}

	return profile, nil
}

// DetectProfiles reads ~/.gitconfig and extracts profiles from includeIf directives
func DetectProfiles() ([]config.Profile, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	gitconfigPath := filepath.Join(home, ".gitconfig")

	// Parse the main gitconfig
	gitCfg, err := ParseGitconfig(gitconfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No global config, no profiles to detect
		}
		return nil, fmt.Errorf("failed to parse gitconfig: %w", err)
	}

	// Extract includeIf directives
	includeIfs := ExtractIncludeIfs(gitCfg)
	if len(includeIfs) == 0 {
		return nil, nil // No includeIf directives
	}

	// Parse each included config to extract profiles
	var profiles []config.Profile
	seen := make(map[string]bool)

	for _, includeIf := range includeIfs {
		profile, err := ParseIncludedConfig(includeIf)
		if err != nil {
			// Skip files that can't be parsed or don't have user info
			continue
		}

		// Avoid duplicate profiles (by email since that's more unique)
		if seen[profile.Email] {
			continue
		}
		seen[profile.Email] = true

		profiles = append(profiles, *profile)
	}

	return profiles, nil
}

// deriveProfileName extracts a meaningful profile name from the includeIf directive
func deriveProfileName(includeIf IncludeIf) string {
	// First, check if the config file path has a clear profile name pattern
	// (e.g., .gitconfig-work, gitconfig-ado)
	configName := extractNameFromConfigPath(includeIf.Path)
	if configName != "" && configName != "default" && !isGenericDirName(configName) {
		return configName
	}

	// Fall back to extracting from the gitdir condition
	condition := includeIf.Condition

	// Handle gitdir: and gitdir/i: conditions
	if strings.HasPrefix(strings.ToLower(condition), "gitdir:") {
		path := strings.TrimPrefix(condition, "gitdir:")
		path = strings.TrimPrefix(path, "gitdir/i:")
		return extractNameFromPath(path)
	}

	if strings.HasPrefix(strings.ToLower(condition), "gitdir/i:") {
		path := strings.TrimPrefix(condition, "gitdir/i:")
		return extractNameFromPath(path)
	}

	// Final fallback
	if configName != "" {
		return configName
	}
	return "default"
}

// isGenericDirName checks if a name is too generic to be a profile name
func isGenericDirName(name string) bool {
	generic := []string{"home", "users", "user", "repos", "repositories", "git", "src", "code"}
	nameLower := strings.ToLower(name)
	for _, g := range generic {
		if nameLower == g {
			return true
		}
	}
	return false
}

// extractNameFromPath extracts a profile name from a directory path
func extractNameFromPath(path string) string {
	// Expand and clean the path
	path = expandPath(path)
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, "/**")

	// Get the last meaningful directory name
	base := filepath.Base(path)
	if base == "." || base == "/" || base == "" {
		// Try parent directory
		parent := filepath.Dir(path)
		base = filepath.Base(parent)
	}

	// Clean up and return
	base = strings.ToLower(base)
	base = strings.TrimPrefix(base, ".")

	if base == "" || base == "home" || base == "users" {
		return "default"
	}

	return base
}

// extractNameFromConfigPath extracts a profile name from a gitconfig file path
func extractNameFromConfigPath(path string) string {
	base := filepath.Base(path)

	// Handle .gitconfig-<name> pattern
	if strings.HasPrefix(base, ".gitconfig-") {
		return strings.TrimPrefix(base, ".gitconfig-")
	}

	// Handle gitconfig-<name> pattern
	if strings.HasPrefix(base, "gitconfig-") {
		return strings.TrimPrefix(base, "gitconfig-")
	}

	// Handle .gitconfig.<name> pattern
	if strings.HasPrefix(base, ".gitconfig.") {
		return strings.TrimPrefix(base, ".gitconfig.")
	}

	// Handle <name>.gitconfig pattern
	if strings.HasSuffix(base, ".gitconfig") {
		return strings.TrimSuffix(base, ".gitconfig")
	}

	// Fall back to directory name
	return extractNameFromPath(filepath.Dir(path))
}

// expandPath expands ~ to home directory and cleans the path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	return filepath.Clean(path)
}

// parseKeyValue parses a "key = value" line
func parseKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Remove quotes if present
	value = strings.Trim(value, "\"'")

	return key, value
}

// extractQuotedValue extracts a value that may be quoted
func extractQuotedValue(s string) string {
	s = strings.TrimSpace(s)
	// Remove surrounding quotes
	s = strings.Trim(s, "\"'")
	return s
}
