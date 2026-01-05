package classifier

import (
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestDetectType(t *testing.T) {
	tests := []struct {
		name      string
		remoteURL string
		expected  config.RepositoryType
	}{
		// GitHub URLs
		{
			name:      "GitHub HTTPS",
			remoteURL: "https://github.com/user/repo.git",
			expected:  config.TypeGitHub,
		},
		{
			name:      "GitHub SSH",
			remoteURL: "git@github.com:user/repo.git",
			expected:  config.TypeGitHub,
		},
		{
			name:      "GitHub HTTPS without .git",
			remoteURL: "https://github.com/user/repo",
			expected:  config.TypeGitHub,
		},

		// GitLab URLs
		{
			name:      "GitLab.com HTTPS",
			remoteURL: "https://gitlab.com/user/repo.git",
			expected:  config.TypeGitLab,
		},
		{
			name:      "GitLab.com SSH",
			remoteURL: "git@gitlab.com:user/repo.git",
			expected:  config.TypeGitLab,
		},
		{
			name:      "Self-hosted GitLab",
			remoteURL: "https://gitlab.mycompany.com/user/repo.git",
			expected:  config.TypeGitLab,
		},

		// Azure DevOps URLs
		{
			name:      "Azure DevOps HTTPS (dev.azure.com)",
			remoteURL: "https://dev.azure.com/organization/project/_git/repo",
			expected:  config.TypeADO,
		},
		{
			name:      "Azure DevOps HTTPS (visualstudio.com)",
			remoteURL: "https://organization.visualstudio.com/project/_git/repo",
			expected:  config.TypeADO,
		},
		{
			name:      "Azure DevOps SSH",
			remoteURL: "git@ssh.dev.azure.com:v3/organization/project/repo",
			expected:  config.TypeADO,
		},

		// Bitbucket URLs
		{
			name:      "Bitbucket HTTPS",
			remoteURL: "https://bitbucket.org/user/repo.git",
			expected:  config.TypeBitbucket,
		},
		{
			name:      "Bitbucket SSH",
			remoteURL: "git@bitbucket.org:user/repo.git",
			expected:  config.TypeBitbucket,
		},

		// Unknown/Edge cases
		{
			name:      "Empty URL",
			remoteURL: "",
			expected:  config.TypeUnknown,
		},
		{
			name:      "Unknown provider",
			remoteURL: "https://git.example.com/user/repo.git",
			expected:  config.TypeUnknown,
		},
		{
			name:      "Local path",
			remoteURL: "/path/to/repo",
			expected:  config.TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectType(tt.remoteURL)
			if result != tt.expected {
				t.Errorf("DetectType(%q) = %v, expected %v", tt.remoteURL, result, tt.expected)
			}
		})
	}
}

func TestDetectVisibility(t *testing.T) {
	tests := []struct {
		name      string
		remoteURL string
		expected  config.RepositoryVisibility
	}{
		// SSH URLs - typically private
		{
			name:      "SSH GitHub",
			remoteURL: "git@github.com:user/repo.git",
			expected:  config.VisibilityPrivate,
		},
		{
			name:      "SSH GitLab",
			remoteURL: "git@gitlab.com:user/repo.git",
			expected:  config.VisibilityPrivate,
		},
		{
			name:      "SSH with ssh:// prefix",
			remoteURL: "ssh://git@github.com/user/repo.git",
			expected:  config.VisibilityPrivate,
		},

		// HTTPS URLs - unknown (could be public or private)
		{
			name:      "HTTPS GitHub",
			remoteURL: "https://github.com/user/repo.git",
			expected:  config.VisibilityUnknown,
		},
		{
			name:      "HTTP GitLab",
			remoteURL: "http://gitlab.com/user/repo.git",
			expected:  config.VisibilityUnknown,
		},

		// Edge cases
		{
			name:      "Empty URL",
			remoteURL: "",
			expected:  config.VisibilityUnknown,
		},
		{
			name:      "Unknown protocol",
			remoteURL: "file:///path/to/repo",
			expected:  config.VisibilityUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectVisibility(tt.remoteURL)
			if result != tt.expected {
				t.Errorf("DetectVisibility(%q) = %v, expected %v", tt.remoteURL, result, tt.expected)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name             string
		remoteURL        string
		expectedType     config.RepositoryType
		expectedVisility config.RepositoryVisibility
	}{
		{
			name:             "GitHub SSH private repo",
			remoteURL:        "git@github.com:user/private-repo.git",
			expectedType:     config.TypeGitHub,
			expectedVisility: config.VisibilityPrivate,
		},
		{
			name:             "GitLab HTTPS repo",
			remoteURL:        "https://gitlab.com/user/public-repo.git",
			expectedType:     config.TypeGitLab,
			expectedVisility: config.VisibilityUnknown,
		},
		{
			name:             "Azure DevOps SSH repo",
			remoteURL:        "git@ssh.dev.azure.com:v3/org/project/repo",
			expectedType:     config.TypeADO,
			expectedVisility: config.VisibilityPrivate,
		},
		{
			name:             "Unknown repo type with HTTPS",
			remoteURL:        "https://git.custom.com/repo.git",
			expectedType:     config.TypeUnknown,
			expectedVisility: config.VisibilityUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &config.Repository{
				Name:      "test-repo",
				Path:      "/path/to/repo",
				RemoteURL: tt.remoteURL,
			}

			Classify(repo)

			if repo.Type != tt.expectedType {
				t.Errorf("Classify(%q) set Type = %v, expected %v", tt.remoteURL, repo.Type, tt.expectedType)
			}
			if repo.Visibility != tt.expectedVisility {
				t.Errorf("Classify(%q) set Visibility = %v, expected %v", tt.remoteURL, repo.Visibility, tt.expectedVisility)
			}
		})
	}
}
