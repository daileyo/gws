package user

import (
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestAddProfile(t *testing.T) {
	tests := []struct {
		name        string
		profile     config.Profile
		existing    []config.Profile
		expectError bool
		errorMsg    string
	}{
		{
			name: "add valid profile",
			profile: config.Profile{
				Name:    "work",
				GitName: "John Doe",
				Email:   "john@company.com",
			},
			existing:    nil,
			expectError: false,
		},
		{
			name: "add profile with signing",
			profile: config.Profile{
				Name:        "personal",
				GitName:     "John Doe",
				Email:       "john@personal.com",
				SigningKey:  "ABCD1234",
				SignCommits: true,
			},
			existing:    nil,
			expectError: false,
		},
		{
			name: "duplicate name fails",
			profile: config.Profile{
				Name:    "work",
				GitName: "Jane Doe",
				Email:   "jane@company.com",
			},
			existing: []config.Profile{
				{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			},
			expectError: true,
			errorMsg:    "already exists",
		},
		{
			name: "duplicate name case insensitive",
			profile: config.Profile{
				Name:    "WORK",
				GitName: "Jane Doe",
				Email:   "jane@company.com",
			},
			existing: []config.Profile{
				{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			},
			expectError: true,
			errorMsg:    "already exists",
		},
		{
			name: "empty name fails",
			profile: config.Profile{
				Name:    "",
				GitName: "John Doe",
				Email:   "john@company.com",
			},
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name: "invalid email fails",
			profile: config.Profile{
				Name:    "work",
				GitName: "John Doe",
				Email:   "invalid-email",
			},
			expectError: true,
			errorMsg:    "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Profiles: tt.existing,
			}

			err := AddProfile(cfg, tt.profile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Verify profile was added
				if len(cfg.Profiles) != len(tt.existing)+1 {
					t.Errorf("Profile was not added to config")
				}
			}
		})
	}
}

func TestRemoveProfile(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		existing    []config.Profile
		repos       []config.Repository
		expectError bool
		expectRepos int // number of repos using the profile
	}{
		{
			name:        "remove existing profile",
			profileName: "work",
			existing: []config.Profile{
				{Name: "work", GitName: "John Doe", Email: "john@company.com"},
				{Name: "personal", GitName: "John Doe", Email: "john@personal.com"},
			},
			expectError: false,
			expectRepos: 0,
		},
		{
			name:        "remove non-existent profile fails",
			profileName: "nonexistent",
			existing: []config.Profile{
				{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			},
			expectError: true,
		},
		{
			name:        "remove profile with repos using it",
			profileName: "work",
			existing: []config.Profile{
				{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			},
			repos: []config.Repository{
				{Name: "repo1", User: "John Doe", Email: "john@company.com"},
				{Name: "repo2", User: "John Doe", Email: "john@company.com"},
				{Name: "repo3", User: "Jane Doe", Email: "jane@other.com"},
			},
			expectError: false,
			expectRepos: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Profiles:     tt.existing,
				Repositories: tt.repos,
			}

			reposUsing, err := RemoveProfile(cfg, tt.profileName)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(reposUsing) != tt.expectRepos {
					t.Errorf("Expected %d repos using profile, got %d", tt.expectRepos, len(reposUsing))
				}
				// Verify profile was removed
				for _, p := range cfg.Profiles {
					if p.Name == tt.profileName {
						t.Error("Profile was not removed from config")
					}
				}
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			{Name: "personal", GitName: "John Doe", Email: "john@personal.com"},
		},
	}

	tests := []struct {
		name        string
		profileName string
		expectError bool
	}{
		{"find existing profile", "work", false},
		{"find case insensitive", "WORK", false},
		{"find mixed case", "Personal", false},
		{"not found", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := GetProfile(cfg, tt.profileName)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if profile == nil {
					t.Error("Expected profile but got nil")
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email       string
		expectError bool
	}{
		{"john@example.com", false},
		{"john.doe@example.com", false},
		{"john+tag@example.com", false},
		{"john@sub.example.com", false},
		{"john_doe@example.co.uk", false},
		{"", true},
		{"invalid", true},
		{"@example.com", true},
		{"john@", true},
		{"john@.com", true},
		{"john@example", true},
		{"john example@test.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := ValidateEmail(tt.email)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for email '%s' but got none", tt.email)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for email '%s': %v", tt.email, err)
			}
		})
	}
}

func TestGetAllProfiles(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "John Doe", Email: "john@company.com"},
		},
	}

	profiles := GetAllProfiles(cfg)

	// Should have at least the stored profile
	if len(profiles) < 1 {
		t.Error("Expected at least 1 profile")
	}

	// First profile should be the stored one
	if profiles[0].Name != "work" {
		t.Errorf("Expected first profile to be 'work', got '%s'", profiles[0].Name)
	}
}

func TestListProfiles(t *testing.T) {
	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "work", GitName: "John Doe", Email: "john@company.com"},
			{Name: "personal", GitName: "John Doe", Email: "john@personal.com"},
		},
	}

	stored, detected := ListProfiles(cfg)

	if len(stored) != 2 {
		t.Errorf("Expected 2 stored profiles, got %d", len(stored))
	}

	// Detected might be empty or have some profiles depending on system gitconfig
	// Just verify it doesn't error
	_ = detected
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsString(s, substr))
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
