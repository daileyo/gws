package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DefaultTTL is the default cache time-to-live (5 minutes)
const DefaultTTL = 5 * time.Minute

// Cache manages cached git status information
type Cache struct {
	data map[string]*Status // repo path -> status
	mu   sync.RWMutex
	ttl  time.Duration
}

// cacheData is used for JSON serialization
type cacheData struct {
	Statuses map[string]*Status `json:"statuses"`
}

// NewCache creates a new status cache with the given TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data: make(map[string]*Status),
		ttl:  ttl,
	}
}

// Get retrieves cached status for a repository
// Returns nil if not cached or if the cache is stale
func (c *Cache) Get(repoPath string) *Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status, ok := c.data[repoPath]
	if !ok {
		return nil
	}

	// Check if cache is stale
	if status.IsStale(c.ttl) {
		return nil
	}

	return status
}

// GetStale retrieves cached status for a repository even if it's past TTL.
// Returns nil only if no entry exists at all. Used to show stale data immediately
// while background refresh happens.
func (c *Cache) GetStale(repoPath string) *Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data[repoPath]
}

// Set stores status for a repository
func (c *Cache) Set(repoPath string, status *Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[repoPath] = status
}

// GetOrFetch retrieves cached status or fetches fresh status if not cached/stale
func (c *Cache) GetOrFetch(repoPath string) (*Status, error) {
	// Try to get from cache
	if status := c.Get(repoPath); status != nil {
		return status, nil
	}

	// Fetch fresh status
	status, err := GetStatus(repoPath)
	if err != nil {
		return nil, err
	}

	// Cache it
	c.Set(repoPath, status)

	return status, nil
}

// Clear removes all cached entries
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*Status)
}

// Save writes the cache to disk
func (c *Cache) Save(cachePath string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Ensure cache directory exists
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data := &cacheData{
		Statuses: c.data,
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Write to file with secure permissions (user read/write only)
	if err := os.WriteFile(cachePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Load reads the cache from disk
func (c *Cache) Load(cachePath string) error {
	// Read file
	jsonData, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No cache file exists yet, that's okay
			return nil
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	// Unmarshal JSON
	var data cacheData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = data.Statuses
	if c.data == nil {
		c.data = make(map[string]*Status)
	}

	return nil
}

// FetchAll fetches status for all given repo paths concurrently using a worker pool.
// It returns cached (fresh) entries without re-fetching. Errors for individual repos
// are silently skipped. The results map contains all successfully resolved statuses.
func (c *Cache) FetchAll(repoPaths []string, workers int) map[string]*Status {
	if workers <= 0 {
		workers = 8
	}

	results := make(map[string]*Status, len(repoPaths))
	var mu sync.Mutex

	// Separate cached from needs-fetch
	var toFetch []string
	for _, p := range repoPaths {
		if s := c.Get(p); s != nil {
			mu.Lock()
			results[p] = s
			mu.Unlock()
		} else {
			toFetch = append(toFetch, p)
		}
	}

	if len(toFetch) == 0 {
		return results
	}

	// Worker pool using buffered channel as semaphore
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for _, repoPath := range toFetch {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			status, err := GetStatus(p)
			if err != nil {
				return // skip this repo
			}

			c.Set(p, status)

			mu.Lock()
			results[p] = status
			mu.Unlock()
		}(repoPath)
	}

	wg.Wait()
	return results
}

// FetchAllStale fetches status only for repos whose cache entries are stale (past TTL).
// Fresh entries and entries with no prior cache are skipped. Used for background prefetch.
func (c *Cache) FetchAllStale(repoPaths []string, workers int) {
	if workers <= 0 {
		workers = 8
	}

	// Find only stale entries that need refreshing
	var toFetch []string
	for _, p := range repoPaths {
		if s := c.Get(p); s == nil {
			// Not fresh — check if there's a stale entry to refresh
			if c.GetStale(p) != nil {
				toFetch = append(toFetch, p)
			}
		}
	}

	if len(toFetch) == 0 {
		return
	}

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for _, repoPath := range toFetch {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			status, err := GetStatus(p)
			if err != nil {
				return
			}
			c.Set(p, status)
		}(repoPath)
	}

	wg.Wait()
}

// PrefetchInBackground spawns a goroutine that refreshes stale cache entries
// and saves the updated cache to disk. Returns a channel that is closed when
// the prefetch completes (callers can optionally wait on it).
func (c *Cache) PrefetchInBackground(repoPaths []string, workers int, cachePath string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.FetchAllStale(repoPaths, workers)
		if cachePath != "" {
			_ = c.Save(cachePath)
		}
	}()
	return done
}

// GetCachePath returns the path to the cache file
func GetCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".gws", "status-cache.json"), nil
}
