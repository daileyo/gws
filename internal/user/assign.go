package user

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
)

// AssignLocal sets user.name and user.email in a repository's local .git/config
func AssignLocal(repoPath string, profile config.Profile) error {
	gitConfigPath := filepath.Join(repoPath, ".git", "config")

	// Verify the repo exists
	if _, err := os.Stat(gitConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Read existing config
	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read git config: %w", err)
	}

	content := string(data)

	// Update or add [user] section
	content = setGitConfigValue(content, "user", "name", profile.GitName)
	content = setGitConfigValue(content, "user", "email", profile.Email)

	// Set signing configuration if provided
	if profile.SigningKey != "" {
		content = setGitConfigValue(content, "user", "signingkey", profile.SigningKey)
	}
	if profile.SignCommits {
		content = setGitConfigValue(content, "commit", "gpgsign", "true")
	}

	// Write updated config
	if err := os.WriteFile(gitConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write git config: %w", err)
	}

	return nil
}

// setGitConfigValue sets a value in a gitconfig content string
func setGitConfigValue(content, section, key, value string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inSection := false
	sectionFound := false
	keySet := false
	sectionHeader := fmt.Sprintf("[%s]", section)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for section header
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			// If we were in the target section and didn't set the key, add it now
			if inSection && !keySet {
				result = append(result, fmt.Sprintf("\t%s = %s", key, value))
				keySet = true
			}

			sectionName := strings.ToLower(strings.Trim(trimmed, "[]"))
			// Handle sections with subsections like [remote "origin"]
			sectionName = strings.Split(sectionName, " ")[0]
			inSection = sectionName == section
			if inSection {
				sectionFound = true
			}
		}

		// Check for key in target section
		if inSection && !keySet {
			keyLower := strings.ToLower(key)
			lineLower := strings.ToLower(trimmed)
			if strings.HasPrefix(lineLower, keyLower+" ") || strings.HasPrefix(lineLower, keyLower+"=") {
				// Replace the key
				result = append(result, fmt.Sprintf("\t%s = %s", key, value))
				keySet = true
				continue
			}
		}

		result = append(result, line)

		// If this is the last line and we're still in the section, add the key
		if i == len(lines)-1 && inSection && !keySet {
			result = append(result, fmt.Sprintf("\t%s = %s", key, value))
			keySet = true
		}
	}

	// If section wasn't found, add it at the end
	if !sectionFound {
		result = append(result, sectionHeader)
		result = append(result, fmt.Sprintf("\t%s = %s", key, value))
	}

	return strings.Join(result, "\n")
}

// AssignWithSubdirs moves a repository to a profile subdirectory and sets up includeIf
func AssignWithSubdirs(cfg *config.Config, repoPath string, profile config.Profile, workspace string) (string, error) {
	// Create profile subdirectory
	subdirPath, err := CreateProfileSubdir(workspace, profile.Name)
	if err != nil {
		return "", err
	}

	// Calculate new path
	repoName := filepath.Base(repoPath)
	newPath := filepath.Join(subdirPath, repoName)

	// Move repository
	if err := MoveRepository(repoPath, newPath); err != nil {
		return "", err
	}

	// Create profile gitconfig
	profileGitconfigPath := filepath.Join(subdirPath, fmt.Sprintf(".gitconfig-%s", profile.Name))
	if err := CreateProfileGitconfig(profileGitconfigPath, profile); err != nil {
		// Try to move repo back on failure
		_ = MoveRepository(newPath, repoPath)
		return "", err
	}

	// Add includeIf to global gitconfig
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	globalGitconfig := filepath.Join(home, ".gitconfig")

	if err := AddIncludeIf(globalGitconfig, subdirPath, profileGitconfigPath); err != nil {
		// Try to move repo back on failure
		_ = MoveRepository(newPath, repoPath)
		return "", err
	}

	// Update repository path in config
	for i := range cfg.Repositories {
		if cfg.Repositories[i].Path == repoPath {
			cfg.Repositories[i].Path = newPath
			cfg.Repositories[i].User = profile.GitName
			cfg.Repositories[i].Email = profile.Email
			cfg.Repositories[i].SigningEnabled = profile.SignCommits
			cfg.Repositories[i].UserSource = config.UserSourceIncludeIf
			break
		}
	}

	return newPath, nil
}

// CreateProfileSubdir creates a subdirectory for a profile under the workspace
func CreateProfileSubdir(workspace string, profileName string) (string, error) {
	subdirPath := filepath.Join(workspace, profileName)

	if err := os.MkdirAll(subdirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create profile subdirectory: %w", err)
	}

	return subdirPath, nil
}

// MoveRepository moves a repository from oldPath to newPath
func MoveRepository(oldPath, newPath string) error {
	// Check if destination already exists
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("destination already exists: %s", newPath)
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(newPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Move the repository
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move repository: %w", err)
	}

	return nil
}

// CreateProfileGitconfig creates a .gitconfig-<profile> file for a profile
func CreateProfileGitconfig(path string, profile config.Profile) error {
	var content strings.Builder

	content.WriteString("[user]\n")
	content.WriteString(fmt.Sprintf("\tname = %s\n", profile.GitName))
	content.WriteString(fmt.Sprintf("\temail = %s\n", profile.Email))

	if profile.SigningKey != "" {
		content.WriteString(fmt.Sprintf("\tsigningkey = %s\n", profile.SigningKey))
	}

	if profile.SignCommits {
		content.WriteString("\n[commit]\n")
		content.WriteString("\tgpgsign = true\n")
	}

	if err := os.WriteFile(path, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write profile gitconfig: %w", err)
	}

	return nil
}

// AddIncludeIf adds an includeIf directive to the global gitconfig
func AddIncludeIf(gitconfigPath, subdirPath, profileGitconfigPath string) error {
	// Create backup first
	if err := backupFile(gitconfigPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Read existing config
	data, err := os.ReadFile(gitconfigPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read gitconfig: %w", err)
	}

	content := string(data)

	// Check if this includeIf already exists
	includeIfSection := fmt.Sprintf("[includeIf \"gitdir:%s/\"]", subdirPath)
	if strings.Contains(content, includeIfSection) {
		// Already exists, skip
		return nil
	}

	// Make path relative to home if possible
	home, _ := os.UserHomeDir()
	displayPath := profileGitconfigPath
	if strings.HasPrefix(profileGitconfigPath, home) {
		displayPath = "~" + profileGitconfigPath[len(home):]
	}

	displaySubdir := subdirPath
	if strings.HasPrefix(subdirPath, home) {
		displaySubdir = "~" + subdirPath[len(home):]
	}

	// Add the includeIf section
	var newContent strings.Builder
	newContent.WriteString(content)
	if !strings.HasSuffix(content, "\n") && content != "" {
		newContent.WriteString("\n")
	}
	newContent.WriteString(fmt.Sprintf("\n[includeIf \"gitdir:%s/\"]\n", displaySubdir))
	newContent.WriteString(fmt.Sprintf("\tpath = %s\n", displayPath))

	if err := os.WriteFile(gitconfigPath, []byte(newContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write gitconfig: %w", err)
	}

	return nil
}

// backupFile creates a backup of a file with timestamp
func backupFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	// Create backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", path, timestamp)

	// Copy file
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// SyncUserInfo updates stored user info for all repositories to match effective config
func SyncUserInfo(cfg *config.Config) (int, error) {
	updated := 0

	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]

		// Get current effective config
		userCfg, err := git.GetUserConfig(repo.Path)
		if err != nil {
			continue // Skip repos that can't be read
		}

		// Check if update needed
		if repo.User != userCfg.Name ||
			repo.Email != userCfg.Email ||
			repo.SigningEnabled != userCfg.SignCommits ||
			repo.UserSource != userCfg.Source {

			repo.User = userCfg.Name
			repo.Email = userCfg.Email
			repo.SigningEnabled = userCfg.SignCommits
			repo.UserSource = userCfg.Source
			updated++
		}
	}

	return updated, nil
}

// PreviewAssignLocal returns what changes would be made without applying them
func PreviewAssignLocal(repoPath string, profile config.Profile) ([]string, error) {
	var changes []string

	// Get current config
	currentCfg, err := git.GetUserConfig(repoPath)
	if err != nil {
		return nil, err
	}

	// Compare and note changes
	if currentCfg.Name != profile.GitName {
		changes = append(changes, fmt.Sprintf("user.name: %q -> %q", currentCfg.Name, profile.GitName))
	}
	if currentCfg.Email != profile.Email {
		changes = append(changes, fmt.Sprintf("user.email: %q -> %q", currentCfg.Email, profile.Email))
	}
	if profile.SigningKey != "" && currentCfg.SigningKey != profile.SigningKey {
		changes = append(changes, fmt.Sprintf("user.signingkey: %q -> %q", currentCfg.SigningKey, profile.SigningKey))
	}
	if profile.SignCommits && !currentCfg.SignCommits {
		changes = append(changes, "commit.gpgsign: false -> true")
	}

	return changes, nil
}
