# Task 4.0 Proof Artifacts - Progress Feedback (Spinner with Count)

## CLI Output: Spinner with Cold Cache

When running `gws list -s` with cold cache on 10 repos, the spinner appears on stderr
with an incrementing count (e.g., `⠋ Checking repos... 3/10`). The spinner uses
Unicode braille characters for smooth animation, renders to stderr, and clears
before the final table output.

The spinner only appears after a 200ms delay threshold, avoiding flicker for
fast operations.

## CLI Output: No Spinner with Warm Cache

```
# Warm cache - completes in 3ms, no spinner shown
$ gws list -s
Found 10 repositories:

NAME                  STATUS
...
```

No spinner is displayed because all results come from cache.

## CLI Output: Piped Output (TTY Detection)

```
$ gws list -s 2>/dev/null | head -5
Found 10 repositories:

NAME                  STATUS
--------------------  -----------------------
my-side-project  master ✓
```

Clean output with no spinner artifacts when piped or redirected.

## Test Results

```
$ go test -race ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	1.442s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	2.064s
ok  	github.com/daileyo/gws/internal/user	1.174s
```

## Verification

- Progress struct uses atomic.Int32 for thread-safe increment
- Spinner renders to stderr with Unicode braille frames (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- 200ms delay threshold before spinner appears
- Stop() clears spinner line with \r\033[K (ANSI erase line)
- TTY detection via term.IsTerminal(os.Stderr.Fd())
- Spinner only created when uncached repos exist
- Double-stop is safe (select on closed channel)
- Race detector passes with no issues
