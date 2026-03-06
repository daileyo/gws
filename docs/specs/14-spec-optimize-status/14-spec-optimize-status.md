# 14-spec-optimize-status

## Introduction/Overview

The `gws list` command fetches git status for each repository sequentially using the `go-git` pure-Go library, resulting in slow performance that scales linearly with repository count. For workspaces with 50-500 repositories, this makes the command painfully slow or unusable for interactive use. This spec defines optimizations to make status fetching fast through parallel execution, switching to the native `git` CLI for status operations, improving caching, and providing user feedback during longer operations.

## Goals

- Reduce `gws list` execution time by 5-10x for typical workspaces (50-100 repos)
- Handle large workspaces (300-500 repos) in a reasonable time with configurable concurrency
- Maintain backward-compatible output and cache format
- Provide visual feedback (spinner with count) during status fetching
- Allow users to tune concurrency via flag and persisted preference

## User Stories

- **As a developer with many repositories**, I want `gws list` to return quickly so that I can use it as part of my regular workflow without waiting.
- **As a developer with a large workspace**, I want to configure the number of concurrent workers so that I can balance speed against system resources.
- **As a developer**, I want to see progress feedback when status checks take more than a moment so that I know the tool is working and roughly how long I'll wait.
- **As a developer**, I want the cache to warm itself in the background so that subsequent runs are even faster.

## Demoable Units of Work

### Unit 1: Switch Status Backend from go-git to git CLI

**Purpose:** Replace the slow `go-git` library calls with native `git` CLI commands for status operations, providing an immediate per-repo speed improvement.

**Functional Requirements:**
- The system shall use `git status --porcelain` (or equivalent) to determine clean/dirty state instead of `go-git`'s `worktree.Status()`
- The system shall use `git rev-list --count` to calculate ahead/behind counts instead of iterating full commit histories via `go-git`
- The system shall use `git rev-parse --abbrev-ref HEAD` (or equivalent) to determine the current branch
- The system shall produce identical `Status` struct output as the current `go-git` implementation
- The system shall handle edge cases (no commits, detached HEAD, no remote tracking branch) identically to the current implementation
- The system shall remove `go-git` as a dependency if it is no longer used elsewhere in the codebase

**Proof Artifacts:**
- CLI: `gws list --status` output is identical before and after the change, demonstrating backward compatibility
- CLI: `time gws list` shows measurable per-repo improvement with cold cache, demonstrating faster git operations

### Unit 2: Parallel Status Fetching with Worker Pool

**Purpose:** Fetch status for multiple repositories concurrently using a goroutine worker pool, providing the largest overall performance improvement.

**Functional Requirements:**
- The system shall fetch repository statuses concurrently using a configurable worker pool (default: 8 workers)
- The system shall accept a `--workers N` flag on the `list` command to override the default concurrency
- The system shall support persisting the worker count preference in `~/.gws/config.json` (e.g., a `preferences` section with `"status_workers": N`)
- The system shall use the persisted preference as default when `--workers` is not explicitly provided
- The system shall collect all statuses before rendering output (results must not interleave or corrupt display)
- The system shall handle errors from individual repo status checks gracefully without stopping the entire batch
- The system shall use a bounded worker pool (semaphore or buffered channel pattern) to avoid exhausting OS resources

**Proof Artifacts:**
- CLI: `time gws list` with 50+ repos shows significant speedup vs sequential, demonstrating parallel fetching works
- CLI: `gws list --workers 16` runs faster than `--workers 2`, demonstrating configurable concurrency
- CLI: Target of 200 repos with status in under 2 seconds (warm or cold cache), demonstrating acceptable performance at scale

### Unit 3: Cache Improvements and Background Prefetch

**Purpose:** Improve cache behavior so that repeated `gws list` calls return near-instantly, and stale entries are refreshed proactively.

**Functional Requirements:**
- The system shall prefetch stale cache entries in background goroutines after returning cached results for display
- The system shall return cached (potentially slightly stale) results immediately when available, then update cache in background
- The system shall save the updated cache to disk after background prefetch completes
- The system shall use the existing TTL mechanism to determine which entries need background refresh
- The system shall ensure background prefetch uses the same worker pool concurrency limits as foreground fetching

**Proof Artifacts:**
- CLI: Second run of `gws list` within TTL window returns near-instantly, demonstrating cache hit path
- CLI: After TTL expires, `gws list` still returns quickly (using stale cache) and refreshes in background, demonstrating prefetch behavior

### Unit 4: Progress Feedback (Spinner with Count)

**Purpose:** Provide visual feedback during status fetching so users know the tool is working and can estimate remaining time.

**Functional Requirements:**
- The system shall display a spinner with progress count (e.g., "Checking repos... 45/200") when fetching status for more than a brief threshold (e.g., 200ms)
- The system shall update the spinner count as each repo status completes
- The system shall clear the spinner line before rendering the final output
- The system shall suppress the spinner when output is not a TTY (piped or redirected)
- The system shall not display the spinner when all results come from cache (no fetching needed)

**Proof Artifacts:**
- CLI: Running `gws list` with cold cache on 50+ repos shows spinner with incrementing count, demonstrating progress feedback
- CLI: Running `gws list` with warm cache shows no spinner, demonstrating it only appears when fetching

## Non-Goals (Out of Scope)

1. **Benchmarking infrastructure**: Formal Go benchmark tests or profiling tooling are not part of this spec
2. **Changing the cache storage format**: The existing JSON-based cache file format will be preserved
3. **Parallelizing non-status operations**: Discovery, filtering, and rendering remain sequential
4. **Configurable cache TTL**: The 5-minute TTL remains fixed (background prefetch makes this less impactful)

## Design Considerations

No specific design requirements identified. The spinner and progress count should follow the existing CLI output style and be unobtrusive.

## Repository Standards

- **Conventional Commits**: All commits must follow the format (e.g., `feat:`, `fix:`, `perf:`)
- **Go 1.24**: Target the current Go version in `go.mod`
- **Testing**: Unit tests with race detector (`go test -race`)
- **Linting**: Code must pass `gofmt` and existing lint checks
- **Architecture**: Follow the existing `internal/` package structure — status and cache logic in `internal/git/`, command integration in `cmd/git-workspace/`

## Technical Considerations

- **git CLI dependency**: The `git` binary must be available on PATH. This is a safe assumption since `gws` manages git repositories. Use `exec.Command` for subprocess calls.
- **Worker pool pattern**: Use a buffered channel or `sync.WaitGroup` + semaphore pattern. Results should be collected into a `map[string]*Status` before rendering.
- **Race conditions**: The cache's existing `sync.RWMutex` protects concurrent writes. Ensure the parallel worker pool writes to cache safely.
- **Config schema evolution**: Adding a `preferences` field to `Config` struct requires backward-compatible JSON unmarshaling (use `omitempty`).
- **go-git removal**: If `go-git` is no longer used after switching to git CLI, remove it from `go.mod` to reduce binary size and dependencies.

## Security Considerations

No specific security considerations identified. The git CLI commands operate on local repositories with the user's existing git credentials and permissions.

## Success Metrics

1. **Performance**: 200 repos status-checked in under 2 seconds (cold cache)
2. **Cache hit performance**: Warm cache `gws list` completes in under 500ms regardless of repo count
3. **No regressions**: All existing tests pass, output format unchanged
4. **Resource safety**: No file descriptor exhaustion or excessive memory usage with 500+ repos

## Open Questions

1. Should the `--workers` flag also apply to the `refresh` command's status operations, or only `list`?
2. What should the spinner look like — a simple character spinner (|/-\) or a Unicode spinner? Should it match any existing spinner patterns in the codebase?
