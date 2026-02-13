package filter

import (
	"testing"

	"github.com/daileyo/gws/internal/config"
)

// Test data
var testRepos = []config.Repository{
	{
		Name:       "my-project",
		Path:       "/home/user/projects/my-project",
		RemoteURL:  "https://github.com/user/my-project.git",
		Type:       config.TypeGitHub,
		Visibility: config.VisibilityUnknown,
		Tags:       []string{"personal", "web", "go"},
	},
	{
		Name:       "work-api",
		Path:       "/home/user/work/work-api",
		RemoteURL:  "git@gitlab.com:company/work-api.git",
		Type:       config.TypeGitLab,
		Visibility: config.VisibilityPrivate,
		Tags:       []string{"work", "backend"},
	},
	{
		Name:       "client-site",
		Path:       "/home/user/clients/client-site",
		RemoteURL:  "https://bitbucket.org/user/client-site.git",
		Type:       config.TypeBitbucket,
		Visibility: config.VisibilityUnknown,
		Tags:       []string{"client", "archived", "web"},
	},
	{
		Name:       "internal-tool",
		Path:       "/home/user/work/internal-tool",
		RemoteURL:  "https://dev.azure.com/org/project/_git/internal-tool",
		Type:       config.TypeADO,
		Visibility: config.VisibilityUnknown,
		Tags:       []string{"work", "tool"},
	},
	{
		Name:       "local-repo",
		Path:       "/home/user/projects/local-repo",
		RemoteURL:  "",
		Type:       config.TypeUnknown,
		Visibility: config.VisibilityUnknown,
		Tags:       []string{"personal"},
	},
}

func TestByType(t *testing.T) {
	tests := []struct {
		name     string
		repoType string
		expected int
	}{
		{
			name:     "Filter by GitHub",
			repoType: "github",
			expected: 1,
		},
		{
			name:     "Filter by GitLab",
			repoType: "gitlab",
			expected: 1,
		},
		{
			name:     "Filter by Bitbucket",
			repoType: "bitbucket",
			expected: 1,
		},
		{
			name:     "Filter by ADO",
			repoType: "ado",
			expected: 1,
		},
		{
			name:     "Filter by Unknown",
			repoType: "unknown",
			expected: 1,
		},
		{
			name:     "Filter by non-existent type",
			repoType: "notexist",
			expected: 0,
		},
		{
			name:     "Case insensitive type filter",
			repoType: "GITHUB",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByType(testRepos, tt.repoType)
			if len(result) != tt.expected {
				t.Errorf("ByType(%q) returned %d repos, expected %d", tt.repoType, len(result), tt.expected)
			}
		})
	}
}

func TestByTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected int
	}{
		{
			name:     "Filter by single tag 'work'",
			tags:     []string{"work"},
			expected: 2, // work-api and internal-tool
		},
		{
			name:     "Filter by single tag 'personal'",
			tags:     []string{"personal"},
			expected: 2, // my-project and local-repo
		},
		{
			name:     "Filter by single tag 'web'",
			tags:     []string{"web"},
			expected: 2, // my-project and client-site
		},
		{
			name:     "Filter by multiple tags (AND) - 'work' AND 'backend'",
			tags:     []string{"work", "backend"},
			expected: 1, // only work-api has both
		},
		{
			name:     "Filter by multiple tags (AND) - 'personal' AND 'web'",
			tags:     []string{"personal", "web"},
			expected: 1, // only my-project has both
		},
		{
			name:     "Filter by non-existent tag",
			tags:     []string{"nonexistent"},
			expected: 0,
		},
		{
			name:     "Filter by multiple tags where none match all",
			tags:     []string{"work", "web"},
			expected: 0, // no repo has both work and web
		},
		{
			name:     "Case insensitive tag filter",
			tags:     []string{"WORK"},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByTags(testRepos, tt.tags)
			if len(result) != tt.expected {
				t.Errorf("ByTags(%v) returned %d repos, expected %d", tt.tags, len(result), tt.expected)
			}
		})
	}
}

func TestByName(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected int
	}{
		{
			name:     "Filter by full name",
			filter:   "my-project",
			expected: 1,
		},
		{
			name:     "Filter by partial name",
			filter:   "project",
			expected: 1,
		},
		{
			name:     "Filter by partial name matching one",
			filter:   "work",
			expected: 1, // work-api (only this has 'work' in name)
		},
		{
			name:     "Filter by name substring",
			filter:   "api",
			expected: 1,
		},
		{
			name:     "Filter by non-existent name",
			filter:   "nonexistent",
			expected: 0,
		},
		{
			name:     "Case insensitive name filter",
			filter:   "MY-PROJECT",
			expected: 1,
		},
		{
			name:     "Filter by single character",
			filter:   "a",
			expected: 3, // work-api, client-site, local-repo
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByName(testRepos, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("ByName(%q) returned %d repos, expected %d", tt.filter, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s", r.Name)
				}
			}
		})
	}
}

func TestByPath(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected int
	}{
		{
			name:     "Filter by full path",
			filter:   "/home/user/projects/my-project",
			expected: 1,
		},
		{
			name:     "Filter by partial path - 'projects'",
			filter:   "projects",
			expected: 2, // my-project and local-repo
		},
		{
			name:     "Filter by partial path - 'work'",
			filter:   "work",
			expected: 2, // work-api and internal-tool
		},
		{
			name:     "Filter by partial path - 'clients'",
			filter:   "clients",
			expected: 1,
		},
		{
			name:     "Filter by non-existent path",
			filter:   "/nonexistent",
			expected: 0,
		},
		{
			name:     "Case insensitive path filter",
			filter:   "PROJECTS",
			expected: 2,
		},
		{
			name:     "Filter by user directory",
			filter:   "/home/user",
			expected: 5, // all repos
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByPath(testRepos, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("ByPath(%q) returned %d repos, expected %d", tt.filter, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s", r.Path)
				}
			}
		})
	}
}

func TestApply_MultipleCriteria(t *testing.T) {
	tests := []struct {
		name     string
		criteria Criteria
		expected int
	}{
		{
			name: "Type and tag filter",
			criteria: Criteria{
				Type: "github",
				Tags: []string{"personal"},
			},
			expected: 1, // my-project
		},
		{
			name: "Type and name filter",
			criteria: Criteria{
				Type: "gitlab",
				Name: "api",
			},
			expected: 1, // work-api
		},
		{
			name: "Tag and path filter",
			criteria: Criteria{
				Tags: []string{"work"},
				Path: "work",
			},
			expected: 2, // work-api and internal-tool
		},
		{
			name: "All criteria combined",
			criteria: Criteria{
				Type: "github",
				Tags: []string{"personal", "web"},
				Name: "project",
				Path: "projects",
			},
			expected: 1, // my-project
		},
		{
			name:     "No criteria (return all)",
			criteria: Criteria{},
			expected: 5,
		},
		{
			name: "Criteria with no matches",
			criteria: Criteria{
				Type: "github",
				Tags: []string{"work"}, // GitHub repo doesn't have work tag
			},
			expected: 0,
		},
		{
			name: "Multiple tags where only one repo matches all",
			criteria: Criteria{
				Tags: []string{"personal", "web", "go"},
			},
			expected: 1, // only my-project has all three
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Apply(testRepos, tt.criteria)
			if len(result) != tt.expected {
				t.Errorf("Apply() with criteria %+v returned %d repos, expected %d", tt.criteria, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s (type: %s, tags: %v)", r.Name, r.Type, r.Tags)
				}
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  string
		expected bool
	}{
		// Non-wildcard (partial, case-insensitive)
		{"partial match", "my-project", "project", true},
		{"case insensitive partial", "My-Project", "PROJECT", true},
		{"no partial match", "my-project", "other", false},

		// Star wildcard
		{"star prefix", "my-project", "my*", true},
		{"star suffix", "my-project", "*project", true},
		{"star both", "my-project", "*-*", true},
		{"star middle", "my-project", "my*ject", true},
		{"star no match", "my-project", "foo*", false},
		{"star match all", "anything", "*", true},
		{"star case insensitive", "GitHub", "git*", true},

		// Question mark wildcard
		{"question mark single char", "web", "w?b", true},
		{"question mark no match", "web", "w??b", false},
		{"question mark in name", "my-project", "m?-project", true},
		{"question mark case insensitive", "GitHub", "?ithub", true},

		// Combined wildcards
		{"star and question", "github", "g?t*", true},
		{"star and question match gitlab", "gitlab", "g?t*b", true},
		{"star and question no match", "github", "g?t*z", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPattern(tt.value, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchesPattern(%q, %q) = %v, expected %v", tt.value, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestByType_Wildcard(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected int
	}{
		{"Wildcard git* matches github and gitlab", "git*", 2},
		{"Wildcard ?ithub matches github", "?ithub", 1},
		{"Wildcard *bucket matches bitbucket", "*bucket", 1},
		{"Wildcard * matches all", "*", 5},
		{"Wildcard no match", "x*", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByType(testRepos, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("ByType(%q) returned %d repos, expected %d", tt.filter, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s (type: %s)", r.Name, r.Type)
				}
			}
		})
	}
}

func TestByName_Wildcard(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected int
	}{
		{"Wildcard my* matches my-project", "my*", 1},
		{"Wildcard *api matches work-api", "*api", 1},
		{"Wildcard *-* matches all with hyphen", "*-*", 5},
		{"Wildcard m?-project matches my-project", "m?-project", 1},
		{"Wildcard *tool matches internal-tool", "*tool", 1},
		{"Wildcard case insensitive", "MY*", 1},
		{"Wildcard no match", "z*", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByName(testRepos, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("ByName(%q) returned %d repos, expected %d", tt.filter, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s", r.Name)
				}
			}
		})
	}
}

func TestByTags_Wildcard(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected int
	}{
		{"Wildcard wo* matches work", []string{"wo*"}, 2},
		{"Wildcard per* matches personal", []string{"per*"}, 2},
		{"Wildcard w?b matches web", []string{"w?b"}, 2},
		{"Wildcard * matches all tagged repos", []string{"*"}, 5},
		{"Wildcard with AND logic", []string{"wo*", "back*"}, 1},
		{"Wildcard no match", []string{"z*"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByTags(testRepos, tt.tags)
			if len(result) != tt.expected {
				t.Errorf("ByTags(%v) returned %d repos, expected %d", tt.tags, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s (tags: %v)", r.Name, r.Tags)
				}
			}
		})
	}
}

func TestByPath_Wildcard(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected int
	}{
		{"Wildcard full path pattern", "/home/user/projects/*", 2},
		{"Wildcard work directory", "/home/user/work/*", 2},
		{"Wildcard clients directory", "/home/user/clients/*", 1},
		{"Wildcard all paths", "/home/user/*", 0}, // filepath.Match doesn't match /
		{"Wildcard with question mark", "/home/user/projects/my-projec?", 1},
		{"Wildcard no match", "/nonexistent/*", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByPath(testRepos, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("ByPath(%q) returned %d repos, expected %d", tt.filter, len(result), tt.expected)
				for _, r := range result {
					t.Logf("  - %s", r.Path)
				}
			}
		})
	}
}

func TestMatchesCriteria(t *testing.T) {
	repo := config.Repository{
		Name:       "test-repo",
		Path:       "/home/user/projects/test-repo",
		Type:       config.TypeGitHub,
		Visibility: config.VisibilityPrivate,
		Tags:       []string{"personal", "web"},
	}

	tests := []struct {
		name     string
		criteria Criteria
		expected bool
	}{
		{
			name:     "Matches type",
			criteria: Criteria{Type: "github"},
			expected: true,
		},
		{
			name:     "Does not match type",
			criteria: Criteria{Type: "gitlab"},
			expected: false,
		},
		{
			name:     "Matches single tag",
			criteria: Criteria{Tags: []string{"personal"}},
			expected: true,
		},
		{
			name:     "Matches all tags",
			criteria: Criteria{Tags: []string{"personal", "web"}},
			expected: true,
		},
		{
			name:     "Does not match all tags",
			criteria: Criteria{Tags: []string{"personal", "backend"}},
			expected: false,
		},
		{
			name:     "Matches name",
			criteria: Criteria{Name: "test"},
			expected: true,
		},
		{
			name:     "Does not match name",
			criteria: Criteria{Name: "other"},
			expected: false,
		},
		{
			name:     "Matches path",
			criteria: Criteria{Path: "projects"},
			expected: true,
		},
		{
			name:     "Does not match path",
			criteria: Criteria{Path: "work"},
			expected: false,
		},
		{
			name:     "Matches all criteria",
			criteria: Criteria{Type: "github", Tags: []string{"personal"}, Name: "test", Path: "projects"},
			expected: true,
		},
		{
			name:     "Fails one criterion",
			criteria: Criteria{Type: "github", Tags: []string{"work"}}, // has type but not tag
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesCriteria(repo, tt.criteria)
			if result != tt.expected {
				t.Errorf("matchesCriteria() = %v, expected %v for criteria %+v", result, tt.expected, tt.criteria)
			}
		})
	}
}
