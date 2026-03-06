package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// createTestRepo creates a temporary git repo and returns its path.
// The caller is responsible for cleanup via t.TempDir().
func createTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	run("init", "-b", "main")
	return dir
}

func TestGetStatus_EmptyRepo(t *testing.T) {
	dir := createTestRepo(t)

	status, err := GetStatus(dir)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.Branch != "" {
		t.Errorf("expected empty branch, got %q", status.Branch)
	}
	if !status.IsClean {
		t.Error("expected empty repo to be clean")
	}
}

func TestGetStatus_CleanBranch(t *testing.T) {
	dir := createTestRepo(t)

	// Create a commit
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v\n%s", err, out)
	}

	status, err := GetStatus(dir)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.Branch != "main" {
		t.Errorf("expected branch 'main', got %q", status.Branch)
	}
	if !status.IsClean {
		t.Error("expected clean status")
	}
	if status.HasChanges {
		t.Error("expected no changes")
	}
}

func TestGetStatus_DirtyRepo(t *testing.T) {
	dir := createTestRepo(t)

	// Create initial commit
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v\n%s", err, out)
	}

	// Make a dirty change
	if err := os.WriteFile(f, []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	status, err := GetStatus(dir)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.IsClean {
		t.Error("expected dirty status")
	}
	if !status.HasChanges {
		t.Error("expected changes")
	}
}

func TestGetStatus_DetachedHead(t *testing.T) {
	dir := createTestRepo(t)

	// Create a commit
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	run("add", ".")
	run("commit", "-m", "initial")

	// Detach HEAD
	run("checkout", "--detach", "HEAD")

	status, err := GetStatus(dir)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.Branch != "HEAD detached" {
		t.Errorf("expected 'HEAD detached', got %q", status.Branch)
	}
}

func TestGetStatus_InvalidPath(t *testing.T) {
	_, err := GetStatus("/nonexistent/path")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestStatus_IsStale(t *testing.T) {
	tests := []struct {
		name        string
		lastChecked time.Time
		ttl         time.Duration
		expected    bool
	}{
		{
			name:        "Fresh status",
			lastChecked: time.Now(),
			ttl:         5 * time.Minute,
			expected:    false,
		},
		{
			name:        "Stale status",
			lastChecked: time.Now().Add(-10 * time.Minute),
			ttl:         5 * time.Minute,
			expected:    true,
		},
		{
			name:        "Just expired",
			lastChecked: time.Now().Add(-5*time.Minute - time.Second),
			ttl:         5 * time.Minute,
			expected:    true,
		},
		{
			name:        "Just fresh",
			lastChecked: time.Now().Add(-4*time.Minute - 59*time.Second),
			ttl:         5 * time.Minute,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &Status{
				LastChecked: tt.lastChecked,
			}

			result := status.IsStale(tt.ttl)
			if result != tt.expected {
				t.Errorf("IsStale() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{
			name: "No commits",
			status: Status{
				Branch: "",
			},
			expected: "no commits",
		},
		{
			name: "Clean branch",
			status: Status{
				Branch:     "main",
				IsClean:    true,
				HasChanges: false,
			},
			expected: "main (clean)",
		},
		{
			name: "Dirty branch",
			status: Status{
				Branch:     "main",
				IsClean:    false,
				HasChanges: true,
			},
			expected: "main (dirty)",
		},
		{
			name: "Ahead of remote",
			status: Status{
				Branch:     "feature",
				IsClean:    true,
				HasChanges: false,
				Ahead:      3,
			},
			expected: "feature (clean) ↑3",
		},
		{
			name: "Behind remote",
			status: Status{
				Branch:     "main",
				IsClean:    true,
				HasChanges: false,
				Behind:     2,
			},
			expected: "main (clean) ↓2",
		},
		{
			name: "Ahead and behind",
			status: Status{
				Branch:     "develop",
				IsClean:    false,
				HasChanges: true,
				Ahead:      5,
				Behind:     3,
			},
			expected: "develop (dirty) ↑5 ↓3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestCache(t *testing.T) {
	t.Run("Get and Set", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		// Get non-existent entry
		result := cache.Get("/path/to/repo")
		if result != nil {
			t.Error("Expected nil for non-existent entry")
		}

		// Set an entry
		status := &Status{
			Branch:      "main",
			IsClean:     true,
			LastChecked: time.Now(),
		}
		cache.Set("/path/to/repo", status)

		// Get the entry
		result = cache.Get("/path/to/repo")
		if result == nil {
			t.Fatal("Expected to get status from cache")
		}

		if result.Branch != "main" {
			t.Errorf("Expected branch 'main', got %q", result.Branch)
		}
	})

	t.Run("Stale cache returns nil", func(t *testing.T) {
		cache := NewCache(1 * time.Millisecond)

		// Set an entry
		status := &Status{
			Branch:      "main",
			LastChecked: time.Now().Add(-2 * time.Millisecond),
		}
		cache.Set("/path/to/repo", status)

		// Get should return nil because it's stale
		result := cache.Get("/path/to/repo")
		if result != nil {
			t.Error("Expected nil for stale cache entry")
		}
	})

	t.Run("Clear cache", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		// Set entries
		cache.Set("/repo1", &Status{Branch: "main", LastChecked: time.Now()})
		cache.Set("/repo2", &Status{Branch: "develop", LastChecked: time.Now()})

		// Clear cache
		cache.Clear()

		// Verify both entries are gone
		if cache.Get("/repo1") != nil {
			t.Error("Expected cache to be cleared")
		}
		if cache.Get("/repo2") != nil {
			t.Error("Expected cache to be cleared")
		}
	})
}

func TestFetchAll(t *testing.T) {
	// Helper to create a test repo with a commit
	makeRepo := func(t *testing.T) string {
		t.Helper()
		dir := createTestRepo(t)
		f := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(f, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git add: %v\n%s", err, out)
		}
		cmd = exec.Command("git", "commit", "-m", "init")
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=t@t.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git commit: %v\n%s", err, out)
		}
		return dir
	}

	t.Run("fetches multiple repos concurrently", func(t *testing.T) {
		repo1 := makeRepo(t)
		repo2 := makeRepo(t)
		repo3 := makeRepo(t)

		cache := NewCache(5 * time.Minute)
		results := cache.FetchAll([]string{repo1, repo2, repo3}, 2, nil)

		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
		for _, path := range []string{repo1, repo2, repo3} {
			if results[path] == nil {
				t.Errorf("missing result for %s", path)
			}
			if results[path].Branch != "main" {
				t.Errorf("expected branch 'main', got %q for %s", results[path].Branch, path)
			}
		}
	})

	t.Run("uses cached entries without re-fetching", func(t *testing.T) {
		repo1 := makeRepo(t)

		cache := NewCache(5 * time.Minute)
		// Pre-populate cache
		cache.Set(repo1, &Status{Branch: "cached-branch", IsClean: true, LastChecked: time.Now()})

		results := cache.FetchAll([]string{repo1}, 2, nil)
		if results[repo1].Branch != "cached-branch" {
			t.Errorf("expected cached branch, got %q", results[repo1].Branch)
		}
	})

	t.Run("handles bad repo gracefully", func(t *testing.T) {
		repo1 := makeRepo(t)
		badPath := "/nonexistent/repo/path"

		cache := NewCache(5 * time.Minute)
		results := cache.FetchAll([]string{repo1, badPath}, 2, nil)

		if results[repo1] == nil {
			t.Error("expected result for valid repo")
		}
		if results[badPath] != nil {
			t.Error("expected nil for bad repo")
		}
	})
}

func TestGetStale(t *testing.T) {
	t.Run("returns stale entry", func(t *testing.T) {
		cache := NewCache(1 * time.Millisecond)
		cache.Set("/repo", &Status{Branch: "main", LastChecked: time.Now().Add(-1 * time.Second)})

		// Get returns nil (stale)
		if cache.Get("/repo") != nil {
			t.Error("expected Get to return nil for stale entry")
		}

		// GetStale returns the entry
		s := cache.GetStale("/repo")
		if s == nil {
			t.Fatal("expected GetStale to return stale entry")
		}
		if s.Branch != "main" {
			t.Errorf("expected branch 'main', got %q", s.Branch)
		}
	})

	t.Run("returns nil for missing entry", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)
		if cache.GetStale("/nonexistent") != nil {
			t.Error("expected nil for missing entry")
		}
	})
}

func TestPrefetchInBackground(t *testing.T) {
	makeRepo := func(t *testing.T) string {
		t.Helper()
		dir := createTestRepo(t)
		f := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(f, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git add: %v\n%s", err, out)
		}
		cmd = exec.Command("git", "commit", "-m", "init")
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=t@t.com",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=t@t.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git commit: %v\n%s", err, out)
		}
		return dir
	}

	t.Run("refreshes stale entries in background", func(t *testing.T) {
		repo := makeRepo(t)
		cache := NewCache(5 * time.Minute)

		// Pre-populate with stale entry (timestamp far in the past)
		cache.Set(repo, &Status{Branch: "stale-branch", LastChecked: time.Now().Add(-10 * time.Minute)})

		// Verify it's stale via Get (returns nil) but exists via GetStale
		if cache.Get(repo) != nil {
			t.Error("expected entry to be stale")
		}
		if cache.GetStale(repo) == nil {
			t.Error("expected stale entry to exist")
		}

		// Run prefetch and wait for completion
		done := cache.PrefetchInBackground([]string{repo}, 2, "")
		<-done

		// Entry should now be fresh
		s := cache.Get(repo)
		if s == nil {
			t.Fatal("expected fresh entry after prefetch")
		}
		if s.Branch != "main" {
			t.Errorf("expected refreshed branch 'main', got %q", s.Branch)
		}
	})
}
