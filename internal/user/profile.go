package user

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/daileyo/gws/internal/config"
)

// AddProfile adds a new profile to the configuration with duplicate name validation
func AddProfile(cfg *config.Config, profile config.Profile) error {
	// Validate profile name
	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	// Validate email
	if err := ValidateEmail(profile.Email); err != nil {
		return err
	}

	// Check for duplicate name
	for _, existing := range cfg.Profiles {
		if strings.EqualFold(existing.Name, profile.Name) {
			return fmt.Errorf("profile with name '%s' already exists", profile.Name)
		}
	}

	// Add profile
	cfg.Profiles = append(cfg.Profiles, profile)
	return nil
}

// RemoveProfile removes a profile from the configuration
// Returns the list of repositories using this profile (for confirmation)
func RemoveProfile(cfg *config.Config, name string) ([]string, error) {
	// Find profile index
	profileIdx := -1
	for i, p := range cfg.Profiles {
		if strings.EqualFold(p.Name, name) {
			profileIdx = i
			break
		}
	}

	if profileIdx == -1 {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}

	// Check if any repositories reference this profile
	var reposUsingProfile []string
	profile := cfg.Profiles[profileIdx]
	for _, repo := range cfg.Repositories {
		// Check if repo's user/email matches this profile
		if repo.User == profile.GitName && repo.Email == profile.Email {
			reposUsingProfile = append(reposUsingProfile, repo.Name)
		}
	}

	// Remove profile from slice
	cfg.Profiles = append(cfg.Profiles[:profileIdx], cfg.Profiles[profileIdx+1:]...)

	return reposUsingProfile, nil
}

// GetProfile finds a profile by name (case-insensitive)
func GetProfile(cfg *config.Config, name string) (*config.Profile, error) {
	for i := range cfg.Profiles {
		if strings.EqualFold(cfg.Profiles[i].Name, name) {
			return &cfg.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile '%s' not found", name)
}

// GetAllProfiles returns all stored profiles merged with auto-detected profiles
// Auto-detected profiles that match stored profiles (by email) are skipped
func GetAllProfiles(cfg *config.Config) []config.Profile {
	// Start with stored profiles
	profiles := make([]config.Profile, len(cfg.Profiles))
	copy(profiles, cfg.Profiles)

	// Create a set of existing emails for deduplication
	emailSet := make(map[string]bool)
	for _, p := range profiles {
		emailSet[strings.ToLower(p.Email)] = true
	}

	// Try to auto-detect profiles from gitconfig
	detected, err := DetectProfiles()
	if err == nil {
		for _, dp := range detected {
			// Only add if email is not already present
			if !emailSet[strings.ToLower(dp.Email)] {
				// Mark as auto-detected by prefixing name if needed
				profiles = append(profiles, dp)
				emailSet[strings.ToLower(dp.Email)] = true
			}
		}
	}

	return profiles
}

// ValidateEmail performs basic email format validation
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Basic email regex pattern
	// Allows: local@domain.tld format with reasonable characters
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// UpdateProfile updates an existing profile with new values
func UpdateProfile(cfg *config.Config, name string, updates config.Profile) error {
	profile, err := GetProfile(cfg, name)
	if err != nil {
		return err
	}

	// Update fields if provided
	if updates.GitName != "" {
		profile.GitName = updates.GitName
	}
	if updates.Email != "" {
		if err := ValidateEmail(updates.Email); err != nil {
			return err
		}
		profile.Email = updates.Email
	}
	if updates.SigningKey != "" {
		profile.SigningKey = updates.SigningKey
	}
	// SignCommits can be explicitly set
	profile.SignCommits = updates.SignCommits

	return nil
}

// MatchProfileByUser finds a stored profile matching the given email (primary)
// or git name (secondary). Returns nil if no match is found.
func MatchProfileByUser(profiles []config.Profile, userName string, email string) *config.Profile {
	// Primary: match by email
	if email != "" {
		for i := range profiles {
			if strings.EqualFold(profiles[i].Email, email) {
				return &profiles[i]
			}
		}
	}

	// Secondary: match by git name
	if userName != "" {
		for i := range profiles {
			if strings.EqualFold(profiles[i].GitName, userName) {
				return &profiles[i]
			}
		}
	}

	return nil
}

// ListProfiles returns a formatted list of all profiles (stored and detected)
func ListProfiles(cfg *config.Config) (stored []config.Profile, detected []config.Profile) {
	stored = cfg.Profiles

	// Get auto-detected profiles
	detectedList, err := DetectProfiles()
	if err != nil {
		return stored, nil
	}

	// Filter out detected profiles that match stored ones
	emailSet := make(map[string]bool)
	for _, p := range stored {
		emailSet[strings.ToLower(p.Email)] = true
	}

	for _, dp := range detectedList {
		if !emailSet[strings.ToLower(dp.Email)] {
			detected = append(detected, dp)
		}
	}

	return stored, detected
}
