package main

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
)

func TestFindRepositories(t *testing.T) {
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
				Name:      "my-project-api",
				Path:      "/home/user/projects/my-project-api",
				RemoteURL: "git@github.com:user/my-project-api.git",
				Type:      config.TypeGitHub,
				Tags:      []string{"personal", "backend"},
			},
			{
				Name:      "work-repo",
				Path:      "/home/user/projects/work-repo",
				RemoteURL: "git@gitlab.com:company/work-repo.git",
				Type:      config.TypeGitLab,
				Tags:      []string{"work"},
			},
			{
				Name:      "Another-API",
				Path:      "/home/user/projects/Another-API",
				RemoteURL: "https://bitbucket.org/user/another.git",
				Type:      config.TypeBitbucket,
				Tags:      []string{},
			},
		},
	}

	tests := []struct {
		name       string
		identifier string
		expected   int // number of matching repositories
		names      []string
	}{
		{
			name:       "Find by partial name - multiple matches",
			identifier: "my-project",
			expected:   2,
			names:      []string{"my-project", "my-project-api"},
		},
		{
			name:       "Find by partial name - single match",
			identifier: "work",
			expected:   1,
			names:      []string{"work-repo"},
		},
		{
			name:       "Find by exact path",
			identifier: "/home/user/projects/work-repo",
			expected:   1,
			names:      []string{"work-repo"},
		},
		{
			name:       "Find by partial name - case insensitive",
			identifier: "API",
			expected:   2,
			names:      []string{"my-project-api", "Another-API"},
		},
		{
			name:       "Find by exact name match",
			identifier: "work-repo",
			expected:   1,
			names:      []string{"work-repo"},
		},
		{
			name:       "No matches",
			identifier: "nonexistent",
			expected:   0,
			names:      []string{},
		},
		{
			name:       "Case insensitive partial match",
			identifier: "PROJECT",
			expected:   2,
			names:      []string{"my-project", "my-project-api"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := findRepositories(cfg, tt.identifier)

			if len(repos) != tt.expected {
				t.Errorf("findRepositories(%q) returned %d repos, expected %d", tt.identifier, len(repos), tt.expected)
			}

			// Verify the names match
			if len(repos) > 0 {
				foundNames := make(map[string]bool)
				for _, repo := range repos {
					foundNames[repo.Name] = true
				}

				for _, expectedName := range tt.names {
					if !foundNames[expectedName] {
						t.Errorf("Expected to find repository %q but didn't", expectedName)
					}
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

func TestTagMultipleRepositories(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "api-service", Path: "/path/api-service", Tags: []string{}},
			{Name: "api-gateway", Path: "/path/api-gateway", Tags: []string{"production"}},
			{Name: "web-app", Path: "/path/web-app", Tags: []string{}},
		},
	}

	// Find all repos with "api" in name
	repos := findRepositories(cfg, "api")
	if len(repos) != 2 {
		t.Errorf("Expected 2 repos with 'api', got %d", len(repos))
	}

	// Simulate tagging all matching repos
	newTag := "backend"
	for _, repo := range repos {
		// Check if tag already exists
		hasTag := false
		for _, tag := range repo.Tags {
			if tag == newTag {
				hasTag = true
				break
			}
		}

		if !hasTag {
			repo.Tags = append(repo.Tags, newTag)
		}
	}

	// Verify both repos now have the backend tag
	apiServiceHasTag := false
	apiGatewayHasTag := false

	for _, tag := range cfg.Repositories[0].Tags {
		if tag == "backend" {
			apiServiceHasTag = true
		}
	}

	for _, tag := range cfg.Repositories[1].Tags {
		if tag == "backend" {
			apiGatewayHasTag = true
		}
	}

	if !apiServiceHasTag {
		t.Error("api-service should have backend tag")
	}

	if !apiGatewayHasTag {
		t.Error("api-gateway should have backend tag")
	}

	// Verify web-app doesn't have the tag
	webAppHasTag := false
	for _, tag := range cfg.Repositories[2].Tags {
		if tag == "backend" {
			webAppHasTag = true
		}
	}

	if webAppHasTag {
		t.Error("web-app should not have backend tag")
	}
}

func TestTagSubcommandsRegistered(t *testing.T) {
	expected := []string{"add", "remove"}
	for _, name := range expected {
		t.Run(name, func(t *testing.T) {
			found := false
			for _, cmd := range tagCmd.Commands() {
				if cmd.Name() == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Sub-subcommand '%s' should be registered on tagCmd", name)
			}
		})
	}
}

func TestTagFlagsMutualExclusivity(t *testing.T) {
	origAdd := tagFlagAdd
	origDel := tagFlagDelete
	defer func() {
		tagFlagAdd = origAdd
		tagFlagDelete = origDel
	}()

	tagFlagAdd = true
	tagFlagDelete = true

	err := tagCmd.RunE(tagCmd, []string{"repo", "tag"})
	if err == nil {
		t.Fatal("Expected error when both -a and -d are set")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("Expected mutual exclusivity error, got: %s", err.Error())
	}
}

func TestTagFlagAdd(t *testing.T) {
	origAdd := tagFlagAdd
	defer func() { tagFlagAdd = origAdd }()

	tagFlagAdd = true

	// With wrong number of args should error
	err := tagCmd.RunE(tagCmd, []string{"repo-only"})
	if err == nil {
		t.Fatal("Expected error with wrong number of args")
	}
	if !strings.Contains(err.Error(), "exactly 2 arguments") {
		t.Errorf("Expected arg count error, got: %s", err.Error())
	}
}

func TestTagFlagDelete(t *testing.T) {
	origDel := tagFlagDelete
	defer func() { tagFlagDelete = origDel }()

	tagFlagDelete = true

	// With wrong number of args should error
	err := tagCmd.RunE(tagCmd, []string{})
	if err == nil {
		t.Fatal("Expected error with wrong number of args")
	}
	if !strings.Contains(err.Error(), "exactly 2 arguments") {
		t.Errorf("Expected arg count error, got: %s", err.Error())
	}
}

func TestTagCompletionFunctionsRegistered(t *testing.T) {
	// Verify ValidArgsFunction is set on all tag commands
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"tagCmd", tagCmd},
		{"tagAddCmd", tagAddCmd},
		{"tagRemoveCmd", tagRemoveCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.ValidArgsFunction == nil {
				t.Errorf("%s should have ValidArgsFunction set", tt.name)
			}
		})
	}
}

func TestCompleteRepoThenNone(t *testing.T) {
	// When 1+ args already provided, should return no completions
	names, directive := completeRepoThenNone(tagAddCmd, []string{"my-repo"}, "")
	if names != nil {
		t.Errorf("Expected nil completions after first arg, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("Expected NoFileComp directive, got %d", directive)
	}
}

func TestCompleteRepoThenTags_SecondArg(t *testing.T) {
	// When 2+ args already provided, should return no completions
	names, directive := completeRepoThenTags(tagRemoveCmd, []string{"repo", "tag"}, "")
	if names != nil {
		t.Errorf("Expected nil completions after second arg, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("Expected NoFileComp directive, got %d", directive)
	}
}

func TestFindRepositoriesByPath(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "api-service", Path: "/home/user/work/api-service"},
			{Name: "api-gateway", Path: "/home/user/work/api-gateway"},
			{Name: "personal-app", Path: "/home/user/personal/personal-app"},
			{Name: "Other-Repo", Path: "/Home/User/Other-Repo"},
		},
	}

	tests := []struct {
		name       string
		pathFilter string
		wantNames  []string
	}{
		{
			name:       "Prefix match returns two repos",
			pathFilter: "/home/user/work",
			wantNames:  []string{"api-service", "api-gateway"},
		},
		{
			name:       "Prefix match single repo",
			pathFilter: "/home/user/work/api-service",
			wantNames:  []string{"api-service"},
		},
		{
			name:       "Substring fallback when no prefix match",
			pathFilter: "user/personal",
			wantNames:  []string{"personal-app"},
		},
		{
			name:       "Case-sensitive no match on wrong case",
			pathFilter: "/HOME/USER/WORK",
			wantNames:  []string{},
		},
		{
			name:       "Exact case match for mixed-case path",
			pathFilter: "/Home/User/Other-Repo",
			wantNames:  []string{"Other-Repo"},
		},
		{
			name:       "No match returns empty",
			pathFilter: "/nonexistent",
			wantNames:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := findRepositoriesWithFilters(cfg, "", tt.pathFilter)

			if len(repos) != len(tt.wantNames) {
				t.Errorf("got %d repos, want %d", len(repos), len(tt.wantNames))
				return
			}

			found := make(map[string]bool)
			for _, r := range repos {
				found[r.Name] = true
			}
			for _, name := range tt.wantNames {
				if !found[name] {
					t.Errorf("expected repo %q in results but it was missing", name)
				}
			}
		})
	}
}

func TestFindRepositoriesWithFilters(t *testing.T) {
	cfg := &config.Config{
		Repositories: []config.Repository{
			{Name: "api-service", Path: "/home/user/work/api-service"},
			{Name: "api-gateway", Path: "/home/user/work/api-gateway"},
			{Name: "personal-app", Path: "/home/user/personal/personal-app"},
			{Name: "Other-Repo", Path: "/Home/User/Other-Repo"},
		},
	}

	tests := []struct {
		name       string
		repoFilter string
		pathFilter string
		wantNames  []string
	}{
		{
			name:       "Repo filter only",
			repoFilter: "api",
			pathFilter: "",
			wantNames:  []string{"api-service", "api-gateway"},
		},
		{
			name:       "Path filter only",
			repoFilter: "",
			pathFilter: "/home/user/work",
			wantNames:  []string{"api-service", "api-gateway"},
		},
		{
			name:       "AND logic: intersection of repo and path",
			repoFilter: "api",
			pathFilter: "/home/user/work/api-service",
			wantNames:  []string{"api-service"},
		},
		{
			name:       "AND logic: no intersection returns empty",
			repoFilter: "personal",
			pathFilter: "/home/user/work",
			wantNames:  []string{},
		},
		{
			name:       "Both empty returns all repos",
			repoFilter: "",
			pathFilter: "",
			wantNames:  []string{"api-service", "api-gateway", "personal-app", "Other-Repo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := findRepositoriesWithFilters(cfg, tt.repoFilter, tt.pathFilter)

			if len(repos) != len(tt.wantNames) {
				t.Errorf("got %d repos, want %d", len(repos), len(tt.wantNames))
				return
			}

			found := make(map[string]bool)
			for _, r := range repos {
				found[r.Name] = true
			}
			for _, name := range tt.wantNames {
				if !found[name] {
					t.Errorf("expected repo %q in results but it was missing", name)
				}
			}
		})
	}
}

func TestTagNoOperationShowsHelp(t *testing.T) {
	origAdd := tagFlagAdd
	origDel := tagFlagDelete
	defer func() {
		tagFlagAdd = origAdd
		tagFlagDelete = origDel
	}()

	tagFlagAdd = false
	tagFlagDelete = false

	// When neither flag nor subcommand is specified, RunE should show help (returns nil)
	err := tagCmd.RunE(tagCmd, []string{})
	if err != nil {
		t.Errorf("Expected nil error (help display), got: %s", err.Error())
	}
}
