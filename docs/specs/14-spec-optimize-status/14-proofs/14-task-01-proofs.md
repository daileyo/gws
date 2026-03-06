# Task 1.0 Proof Artifacts - Switch Status Backend from go-git to git CLI

## CLI Output: Backward Compatibility

Baseline (go-git) and new (git CLI) produce identical output:

```
Found 10 repositories:

NAME                  STATUS
--------------------  -----------------------
my-side-project  master ✗
claude-statusline     main ✗
dratdns               master ✗
dratreprox            main ✗ ↑1
gws                   fix-optimize-status ✗
hl-talos              no commits
nvim                  refactor ✗
reverse-proxy         master ✗
nvim-config           main ✗
pvt-dotfiles          minidrat ✗ ↑2
```

## CLI Output: Timing Comparison (Cold Cache, 10 repos)

```
# Baseline (go-git)
/tmp/gws-baseline list -s > /dev/null 2>&1  0.37s user 0.68s system 12% cpu 8.289 total

# New (git CLI)
/tmp/gws-new list -s > /dev/null 2>&1  0.54s user 0.38s system 12% cpu 7.653 total
```

Per-repo improvement is modest since most time is spent in sequential I/O. The major speedup will come from parallelization (Task 2.0).

## Test Results

```
$ go test -race ./internal/git/...
ok  	github.com/daileyo/gws/internal/git	1.457s

$ go test -race ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	1.432s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	1.138s
```

## Verification

- Output is byte-identical between baseline and new implementation
- All edge cases tested: empty repo, clean branch, dirty repo, detached HEAD, invalid path
- All existing tests (IsStale, String, Cache) continue to pass
- go-git imports removed from status.go; dependency stays in go.mod for other packages
- Race detector passes with no issues
