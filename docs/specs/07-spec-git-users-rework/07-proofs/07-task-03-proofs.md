# 07 Task 3.0 Proof Artifacts - Wire User Detection into Init and Add

## Changes Summary

| File | Change |
|------|--------|
| `cmd/git-workspace/userdetect.go` | New file: shared `detectUserForRepos()` helper function |
| `cmd/git-workspace/refresh.go` | Refactored to use shared helper instead of inline user detection loop |
| `cmd/git-workspace/init.go` | Added user detection after repo discovery, reports count |
| `cmd/git-workspace/add.go` | Added user detection for single add and recursive add |

## Test Results

All existing tests continue to pass after refactoring:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.098s
ok  	github.com/daileyo/gws/internal/git	0.024s
ok  	github.com/daileyo/gws/internal/user	0.105s
```

## Implementation Details

- Created `cmd/git-workspace/userdetect.go` with `detectUserForRepos(repos []config.Repository) int`
- The helper iterates over repos by index (pointer semantics), calls `git.GetUserConfig()`, and populates User/Email/SigningEnabled/UserSource
- `refresh.go` now separates tag preservation from user detection, calling the shared helper
- `init.go` calls the helper after `cfg.Repositories = result.Repositories` and prints count if > 0
- `add.go` single add: creates a temporary slice, calls helper, copies back to pointer
- `add.go` recursive add: calls helper on `newRepos` slice before appending to config
