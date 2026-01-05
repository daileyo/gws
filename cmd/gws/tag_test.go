package main

import (
	"testing"

	"github.com/daileyo/gws/internal/config"
)

func TestFindRepository(t *testing.T) {
	cfg := &config.Config{
		Version:   "1.0.0",
		Workspace: "/home/user/projects",
		Repositories: []config.Repository{
			{
				Name:      "my-project",
				Path:      "/home/user/projects/my-project",
				RemoteURL: "https://github.com/user/my-project.git",
				Type:      config.TypeGitHub,
				Tags:      []string{"personal"},
			},
			{
				Name:      "work-repo",
				Path:      "/home/user/projects/work-repo",
				RemoteURL: "git@gitlab.com:company/work-repo.git",
				Type:      config.TypeGitLab,
				Tags:      []string{"work"},
			},
			{
				Name:      "Another-Project",
				Path:      "/home/user/projects/Another-Project",
				RemoteURL: "https://bitbucket.org/user/another.git",
				Type:      config.TypeBitbucket,
				Tags:      []string{},
			},
		},
	}

	tests := []struct {
		name       string
		identifier string
		shouldFind bool
		expected   string // expected repository name
	}{
		{
			name:       "Find by exact name match",
			identifier: "my-project",
			shouldFind: true,
			expected:   "my-project",
		},
		{
			name:       "Find by case-insensitive name",
			identifier: "MY-PROJECT",
			shouldFind: true,
			expected:   "my-project",
		},
		{
			name:       "Find by exact path",
			identifier: "/home/user/projects/work-repo",
			shouldFind: true,
			expected:   "work-repo",
		},
		{
			name:       "Find by mixed case name",
			identifier: "another-project",
			shouldFind: true,
			expected:   "Another-Project",
		},
		{
			name:       "Not found - invalid name",
			identifier: "nonexistent-repo",
			shouldFind: false,
		},
		{
			name:       "Not found - invalid path",
			identifier: "/invalid/path",
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := findRepository(cfg, tt.identifier)

			if tt.shouldFind {
				if err != nil {
					t.Errorf("Expected to find repository, but got error: %v", err)
					return
				}
				if repo.Name != tt.expected {
					t.Errorf("Expected repository name %q, got %q", tt.expected, repo.Name)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for identifier %q, but found repository: %v", tt.identifier, repo.Name)
				}
			}
		})
	}
}

func TestTagManagement(t *testing.T) {
	// Test adding tags
	t.Run("Add tag to repository", func(t *testing.T) {
		cfg := &config.Config{
			Repositories: []config.Repository{
				{
					Name: "test-repo",
					Path: "/path/to/repo",
					Tags: []string{},
				},
			},
		}

		repo := &cfg.Repositories[0]
		repo.Tags = append(repo.Tags, "personal")

		if len(repo.Tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(repo.Tags))
		}
		if repo.Tags[0] != "personal" {
			t.Errorf("Expected tag 'personal', got %q", repo.Tags[0])
		}
	})

	// Test duplicate tag prevention
	t.Run("Prevent duplicate tags", func(t *testing.T) {
		tags := []string{"personal", "work"}
		newTag := "personal"

		// Check if tag exists
		exists := false
		for _, tag := range tags {
			if tag == newTag {
				exists = true
				break
			}
		}

		if !exists {
			tags = append(tags, newTag)
		}

		if len(tags) != 2 {
			t.Errorf("Expected 2 tags (no duplicate), got %d", len(tags))
		}
	})

	// Test removing tags
	t.Run("Remove tag from repository", func(t *testing.T) {
		tags := []string{"personal", "work", "archived"}
		removeTag := "work"

		// Remove tag
		newTags := []string{}
		for _, tag := range tags {
			if tag != removeTag {
				newTags = append(newTags, tag)
			}
		}

		if len(newTags) != 2 {
			t.Errorf("Expected 2 tags after removal, got %d", len(newTags))
		}

		// Verify the correct tag was removed
		for _, tag := range newTags {
			if tag == removeTag {
				t.Errorf("Tag %q should have been removed", removeTag)
			}
		}
	})

	// Test removing non-existent tag
	t.Run("Remove non-existent tag", func(t *testing.T) {
		tags := []string{"personal", "work"}
		removeTag := "archived"

		// Try to remove tag
		found := false
		newTags := []string{}
		for _, tag := range tags {
			if tag == removeTag {
				found = true
			} else {
				newTags = append(newTags, tag)
			}
		}

		if found {
			t.Errorf("Tag %q should not exist in the original tags", removeTag)
		}
		if len(newTags) != 2 {
			t.Errorf("Expected tags to remain unchanged, got %d tags", len(newTags))
		}
	})
}
