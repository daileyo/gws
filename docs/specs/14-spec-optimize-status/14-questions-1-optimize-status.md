# 14 Questions Round 1 - Optimize Status

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Scale & Baseline

How many repositories are you typically managing, and roughly how long does a `gws` or `gws list` call take currently?

- [ ] (A) 50-100 repos, takes a few seconds (annoying but tolerable)
- [ ] (B) 100-300 repos, takes 10-30 seconds (painful)
- [ ] (C) 300-500 repos, takes 30-60 seconds (unusable for interactive use)
- [ ] (D) 500+ repos, takes over a minute
- [X] (E) Other (describe) A is typical; but C is not uncommon for me. I would like this to account for A by default... but perhaps be a configurable (if that is even needed) to account for C.

## 2. Concurrency Strategy

The biggest optimization opportunity is fetching status for multiple repos in parallel using goroutines. What's your preference?

- [ ] (A) Worker pool with configurable concurrency (e.g., `--workers N` flag, sensible default like 8-16)
- [ ] (B) Worker pool with a fixed sensible default (no flag, just make it fast)
- [ ] (C) Fully concurrent (one goroutine per repo, no limit) — fastest but may hit OS/git limits
- [X] (D) Other (describe) I would like B, but with a configurable concurrency (e.g., `--workers N` flag that can be used to make it faster if needed. bonus if it is something that can be saved to the meta data as a prefernece.

## 3. Cache Strategy Improvements

The current cache has a 5-minute TTL and persists to disk. Would you like to adjust caching behavior?

- [X] (A) Keep current TTL approach but add background prefetch (warm cache ahead of use)
- [ ] (B) Make TTL configurable (e.g., `--cache-ttl` flag or config setting)
- [ ] (C) Add a `--no-cache` / `--fresh` flag to force re-fetch when needed
- [ ] (D) Current cache strategy is fine, focus only on parallel fetching
- [ ] (E) Other (describe)

## 4. Git Backend

The current implementation uses `go-git` (pure Go). Shelling out to `git` CLI can be significantly faster for status operations. Should we consider this?

- [X] (A) Switch to `git` CLI for status operations (faster, requires git installed — which it always is for gws users)
- [ ] (B) Keep `go-git` but optimize the operations (fewer allocations, smarter queries)
- [ ] (C) Hybrid — use `git` CLI for status, keep `go-git` for other things
- [ ] (D) Open to whichever is measurably faster
- [ ] (E) Other (describe)

## 5. User Feedback During Long Operations

When status fetching takes time, should the user see progress?

- [ ] (A) Yes, show a progress bar or counter (e.g., "Checking repos: 45/200")
- [X] (B) Yes, show a simple spinner with count
- [ ] (C) No, just make it fast enough that progress isn't needed
- [ ] (D) Other (describe)

## 6. Scope Boundaries

Which of these should be IN scope for this optimization? (Select all that apply)

- [X] (A) Parallel status fetching (goroutine worker pool)
- [X] (B) Faster git operations (CLI vs go-git)
- [X] (C) Cache improvements (TTL, prefetch, flags)
- [X] (D) Progress feedback during status checks
- [ ] (E) Benchmarking/profiling infrastructure to measure improvements
- [ ] (F) Other (describe)

## 7. Proof Artifacts

What would convince you the optimization is successful?

- [X] (A) Before/after timing comparison (e.g., `time gws list` output)
- [ ] (B) Benchmark tests in the codebase (Go benchmarks)
- [X] (C) A specific target time (e.g., "200 repos in under 2 seconds")
- [ ] (D) Just noticeably faster in daily use
- [ ] (E) Other (describe)
