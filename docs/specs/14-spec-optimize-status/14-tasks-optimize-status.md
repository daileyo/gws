# 14-tasks-optimize-status

## Relevant Files

- `internal/git/status.go` - Core status fetching logic; rewrite `GetStatus()`, `getAheadBehind()`, `countCommitsAhead()` to use git CLI
- `internal/git/status_test.go` - Unit tests for status; add tests for the new git CLI-based implementation
- `internal/git/cache.go` - Cache with TTL; add `FetchAll()` for parallel fetching, `GetStale()` and background prefetch support
- `internal/git/cache_test.go` - Tests for new cache methods (parallel fetch, stale retrieval, background prefetch)
- `internal/git/progress.go` - New file for spinner/progress feedback during status fetching
- `internal/git/progress_test.go` - Tests for progress feedback (TTY detection, threshold behavior)
- `internal/config/config.go` - Add `Preferences` struct with `StatusWorkers` field to `Config`
- `cmd/git-workspace/list.go` - Add `--workers` flag, integrate parallel fetch, spinner, and background prefetch
- `go.mod` - No changes needed (go-git stays; used by scanner, remote, user, add)

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `status.go` and `status_test.go` in the same directory).
- Use `go test -race ./...` to run tests with the race detector enabled.
- Follow Conventional Commits format for all commits (e.g., `perf:`, `feat:`, `fix:`).
- Follow the existing `internal/` package structure and naming conventions.

## Tasks

### [x] 1.0 Switch Status Backend from go-git to git CLI

Replace the `go-git` library calls in `internal/git/status.go` with native `git` CLI subprocess calls. This provides an immediate per-repo speed improvement by using optimized C-based git instead of the pure-Go implementation. The `Status` struct, `IsStale()`, and `String()` methods remain unchanged. Only the `GetStatus()`, `getAheadBehind()`, and `countCommitsAhead()` functions are rewritten to use `exec.Command`. Edge cases (no commits, detached HEAD, no remote tracking branch) must behave identically.

#### 1.0 Proof Artifact(s)

- CLI: `gws list --status` output is identical before and after the change, demonstrating backward compatibility
- CLI: `time gws list --status` with cold cache shows measurable improvement, demonstrating faster git operations
- Test: `go test -race ./internal/git/...` passes, demonstrating correctness and thread safety

#### 1.0 Tasks

- [x] 1.1 Rewrite `GetStatus()` in `internal/git/status.go` to use `exec.Command` instead of `go-git`. Use `git rev-parse --abbrev-ref HEAD` for branch name, `git status --porcelain` for clean/dirty state, and `git rev-list --count` for ahead/behind. Keep the `Status` struct, `IsStale()`, and `String()` methods unchanged.
- [x] 1.2 Handle edge cases in the new CLI-based `GetStatus()`: no commits yet (empty repo), detached HEAD, no remote tracking branch (no upstream set). Each case must return the same `Status` values as the current `go-git` implementation.
- [x] 1.3 Remove the `go-git` imports from `status.go` (the `go-git` dependency stays in `go.mod` since other files use it). Replace with `os/exec`, `strings`, and `strconv` imports.
- [x] 1.4 Update `internal/git/status_test.go` with tests for the new CLI-based `GetStatus()`. Create temporary git repos in tests (using `git init`, `git commit`, etc. via `exec.Command`) to verify branch detection, clean/dirty state, and ahead/behind counts. Existing `IsStale` and `String` tests remain unchanged.
- [x] 1.5 Run `go test -race ./internal/git/...` and `go test -race ./...` to verify all tests pass with no regressions.

### [x] 2.0 Parallel Status Fetching with Worker Pool

Add a goroutine-based worker pool to fetch repository statuses concurrently. Introduce a `--workers N` flag on the `list` command (default: 8). Add a `preferences` section to `~/.gws/config.json` to persist the worker count so users can set-and-forget their preferred concurrency. All statuses are collected into a map before rendering begins, ensuring output is never interleaved. Individual repo errors are handled gracefully (skipped, not fatal).

#### 2.0 Proof Artifact(s)

- CLI: `time gws list --status` with 50+ repos shows significant speedup vs v2.13.0 baseline, demonstrating parallel fetching works
- CLI: `gws list --workers 16` completes faster than `gws list --workers 2`, demonstrating configurable concurrency
- CLI: 200 repos with status in under 2 seconds, demonstrating target performance at scale
- Test: `go test -race ./...` passes, demonstrating no race conditions in concurrent code

#### 2.0 Tasks

- [x] 2.1 Add a `FetchAll(repoPaths []string, workers int) map[string]*Status` method to `internal/git/cache.go`. This method uses a buffered channel (size = `workers`) as a semaphore, spawns goroutines for each repo that needs fetching (cache miss or stale), collects results into a `map[string]*Status` protected by the existing `sync.RWMutex`, and returns the complete map. Errors for individual repos are logged/skipped, not fatal.
- [x] 2.2 Add a `Preferences` struct to `internal/config/config.go` with a `StatusWorkers int` field (`json:"status_workers,omitempty"`). Add a `Preferences *Preferences` field to the `Config` struct (`json:"preferences,omitempty"`). This is backward-compatible — existing configs without the field will unmarshal with `nil` preferences.
- [x] 2.3 Add a `--workers` flag (type `int`, default `0` meaning "use config or default") to the `list` command in `cmd/git-workspace/list.go`. In `runList()`, resolve the effective worker count: use `--workers` value if provided > 0, else use `config.Preferences.StatusWorkers` if set > 0, else default to `8`. Add `Workers int` to `ListOptions`.
- [x] 2.4 Refactor `displayTable()` and `displayJSON()` in `list.go` to pre-fetch all statuses using `cache.FetchAll()` before the rendering loop. Replace the per-repo `statusCache.GetOrFetch()` calls in the column-width calculation loop (line ~369) and row rendering loop (line ~503) with lookups from the pre-fetched map. The JSON path (line ~719) should also use the pre-fetched map.
- [x] 2.5 Add tests for `FetchAll()` in `internal/git/cache_test.go`. Test with multiple repos, verify concurrency (results arrive for all repos), verify error handling (one bad repo doesn't stop others), and run with `-race` flag.
- [x] 2.6 Run `go test -race ./...` to verify no race conditions in the concurrent code.

### [x] 3.0 Cache Improvements and Background Prefetch

Enhance the cache to return stale entries immediately for display while refreshing them in background goroutines. After the user sees results, the cache is updated in the background using the same worker pool concurrency limits. The updated cache is saved to disk once background refresh completes. This makes repeated `gws list` calls near-instant even after TTL expiry.

#### 3.0 Proof Artifact(s)

- CLI: Second `gws list --status` run within TTL window returns near-instantly, demonstrating cache hit path
- CLI: After TTL expires, `gws list --status` still returns quickly (using stale cache), demonstrating background prefetch behavior
- Test: `go test -race ./internal/git/...` passes, demonstrating cache concurrency safety

#### 3.0 Tasks

- [x] 3.1 Add a `GetStale(repoPath string) *Status` method to `internal/git/cache.go` that returns a cached entry even if it's past TTL (returns `nil` only if no entry exists at all). This is used to show stale data immediately while background refresh happens.
- [x] 3.2 Add `FetchAllStale()` method that only fetches repos whose cache entries are stale (past TTL). This is used for background prefetch.
- [x] 3.3 Add a `PrefetchInBackground(repoPaths []string, workers int, cachePath string)` method to `internal/git/cache.go`. This spawns a goroutine that calls `FetchAllStale()`, then saves the cache to disk. Returns a channel for optional waiting.
- [x] 3.4 Update `runList()` in `list.go` to use a two-phase approach: (1) build the display map from fresh cache + stale fallback, fetch only uncached repos; (2) after rendering, `PrefetchInBackground()` refreshes stale entries with `defer <-done` to wait before exit.
- [x] 3.5 Add tests for `GetStale()` and `PrefetchInBackground()` in `status_test.go`. Verify that `GetStale()` returns expired entries, and that `PrefetchInBackground()` updates the cache asynchronously.
- [x] 3.6 Run `go test -race ./...` to verify cache concurrency safety.

### [ ] 4.0 Progress Feedback (Spinner with Count)

Display a spinner with progress count (e.g., "Checking repos... 45/200") when status fetching takes longer than 200ms. The spinner updates as each repo completes, then clears before final output renders. The spinner is suppressed when output is not a TTY (piped/redirected) and when all results come from cache.

#### 4.0 Proof Artifact(s)

- CLI: `gws list --status` with cold cache on 50+ repos shows spinner with incrementing count, demonstrating progress feedback
- CLI: `gws list --status` with warm cache shows no spinner, demonstrating it only appears when fetching
- CLI: `gws list --status | head` produces clean output without spinner artifacts, demonstrating TTY detection

#### 4.0 Tasks

- [ ] 4.1 Create `internal/git/progress.go` with a `Progress` struct that manages spinner display. It should accept a `total int` count, have an `Increment()` method (thread-safe via `atomic.Int32`), and a `Start()` method that begins a background goroutine rendering the spinner to stderr. The spinner only appears after a 200ms delay threshold.
- [ ] 4.2 Add a `Stop()` method to `Progress` that clears the spinner line (write `\r` + spaces + `\r` to stderr) and stops the background goroutine. This must be called before rendering the final table output.
- [ ] 4.3 Add TTY detection to `Progress.Start()`: check `term.IsTerminal(int(os.Stderr.Fd()))` and skip all spinner output if not a TTY. Use the existing `golang.org/x/term` dependency already in the project.
- [ ] 4.4 Integrate the `Progress` spinner into `FetchAll()` in `cache.go`. Accept an optional `*Progress` parameter (or nil to disable). Call `progress.Increment()` after each repo status completes. The caller in `list.go` creates the `Progress` instance and passes it to `FetchAll()`.
- [ ] 4.5 In `list.go`, create the `Progress` spinner only when there are repos to fetch (not all cached) and `ShowStatus` is true. Pass it to `FetchAll()`, then call `Stop()` before rendering. Skip creating the spinner entirely when all statuses come from cache.
- [ ] 4.6 Add tests for `Progress` in `internal/git/progress_test.go`. Test that `Increment()` is thread-safe (concurrent calls), that `Stop()` clears output, and that spinner is suppressed for non-TTY. Run with `-race` flag.
- [ ] 4.7 Run `go test -race ./...` to verify all tests pass with no regressions.
