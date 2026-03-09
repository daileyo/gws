package git

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
)

func TestGetRemoteInfo_OriginOnly(t *testing.T) {
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/user/repo.git"},
	})
	if err != nil {
		t.Fatalf("Failed to create remote: %v", err)
	}

	info, err := GetRemoteInfo(dir)
	if err != nil {
		t.Fatalf("GetRemoteInfo failed: %v", err)
	}

	if info.OriginURL != "https://github.com/user/repo.git" {
		t.Errorf("Expected origin URL 'https://github.com/user/repo.git', got '%s'", info.OriginURL)
	}
	if info.HasMultiple {
		t.Error("Expected HasMultiple to be false for origin-only repo")
	}
}

func TestGetRemoteInfo_OriginPlusUpstream(t *testing.T) {
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/user/fork.git"},
	})
	if err != nil {
		t.Fatalf("Failed to create origin: %v", err)
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "upstream",
		URLs: []string{"https://github.com/org/repo.git"},
	})
	if err != nil {
		t.Fatalf("Failed to create upstream: %v", err)
	}

	info, err := GetRemoteInfo(dir)
	if err != nil {
		t.Fatalf("GetRemoteInfo failed: %v", err)
	}

	if info.OriginURL != "https://github.com/user/fork.git" {
		t.Errorf("Expected origin URL 'https://github.com/user/fork.git', got '%s'", info.OriginURL)
	}
	if !info.HasMultiple {
		t.Error("Expected HasMultiple to be true for origin+upstream repo")
	}
}

func TestGetRemoteInfo_NoOriginWithOtherRemotes(t *testing.T) {
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "upstream",
		URLs: []string{"https://github.com/org/repo.git"},
	})
	if err != nil {
		t.Fatalf("Failed to create upstream: %v", err)
	}

	info, err := GetRemoteInfo(dir)
	if err != nil {
		t.Fatalf("GetRemoteInfo failed: %v", err)
	}

	if info.OriginURL != "" {
		t.Errorf("Expected empty origin URL, got '%s'", info.OriginURL)
	}
	if !info.HasMultiple {
		t.Error("Expected HasMultiple to be true for no-origin repo with other remotes")
	}
}

func TestGetRemoteInfo_NoRemotes(t *testing.T) {
	dir := t.TempDir()
	_, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	info, err := GetRemoteInfo(dir)
	if err != nil {
		t.Fatalf("GetRemoteInfo failed: %v", err)
	}

	if info.OriginURL != "" {
		t.Errorf("Expected empty origin URL, got '%s'", info.OriginURL)
	}
	if info.HasMultiple {
		t.Error("Expected HasMultiple to be false for repo with no remotes")
	}
}

func TestGetRemoteInfo_InvalidPath(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "nonexistent-repo-path-12345")

	_, err := GetRemoteInfo(dir)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestFormatRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SSH GitHub",
			input:    "git@github.com:owner/repo.git",
			expected: "https://github.com/owner/repo.git",
		},
		{
			name:     "SSH GitLab with subgroup",
			input:    "git@gitlab.com:group/subgroup/repo.git",
			expected: "https://gitlab.com/group/subgroup/repo.git",
		},
		{
			name:     "HTTPS without user info",
			input:    "https://github.com/owner/repo.git",
			expected: "https://github.com/owner/repo.git",
		},
		{
			name:     "HTTPS with user info",
			input:    "https://user@github.com/owner/repo.git",
			expected: "https://github.com/owner/repo.git",
		},
		{
			name:     "HTTPS with user and password",
			input:    "https://user:pass@github.com/owner/repo.git",
			expected: "https://github.com/owner/repo.git",
		},
		{
			name:     "Azure DevOps SSH",
			input:    "git@ssh.dev.azure.com:v3/org/project/repo",
			expected: "https://dev.azure.com/org/project/_git/repo",
		},
		{
			name:     "Azure DevOps HTTPS with user info",
			input:    "https://user@dev.azure.com/org/project/_git/repo",
			expected: "https://dev.azure.com/org/project/_git/repo",
		},
		{
			name:     "file protocol unchanged",
			input:    "file:///local/path/repo",
			expected: "file:///local/path/repo",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "SSH protocol URL unchanged",
			input:    "ssh://git@github.com/owner/repo.git",
			expected: "ssh://git@github.com/owner/repo.git",
		},
		{
			name:     "HTTP with user info",
			input:    "http://user@example.com/repo.git",
			expected: "http://example.com/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRemoteURL(tt.input)
			if got != tt.expected {
				t.Errorf("FormatRemoteURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
