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
