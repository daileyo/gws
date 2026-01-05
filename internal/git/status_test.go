package git

import (
	"testing"
	"time"
)

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
