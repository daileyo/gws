# Task 3.0 Proof Artifacts - Cache Improvements and Background Prefetch

## CLI Output: Cache Hit Performance

```
# Cold cache (first run)
/tmp/gws-new list -s > /dev/null 2>&1  0.19s user 0.51s system 44% cpu 1.595 total

# Warm cache (second run, within TTL)
/tmp/gws-new list -s > /dev/null 2>&1  0.00s user 0.00s system 120% cpu 0.003 total
```

Warm cache returns in **3ms** — near-instant as specified.

## Background Prefetch Behavior

After TTL expires, `gws list --status` still returns quickly using stale cached data,
then refreshes entries in the background. The `PrefetchInBackground()` method returns
a channel that the caller waits on (`defer <-done`) to ensure the background goroutine
completes before process exit, so the updated cache is saved to disk.

## Test Results

```
$ go test -race ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	(cached)
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	2.234s
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Verification

- GetStale() returns cached entries regardless of TTL (nil only for missing)
- FetchAllStale() only fetches repos with stale cache entries
- PrefetchInBackground() spawns goroutine, refreshes stale entries, saves to disk
- runList() two-phase approach: display from stale cache, then prefetch in background
- Background goroutine waited on before process exit via defer <-done
- Race detector passes with no issues
