# 07 Task 5.0 Proof Artifacts - Local Config Detection for Display

## Test Results

```
=== RUN   TestGetNonLocalUserConfig_SkipsLocalConfig
--- PASS: TestGetNonLocalUserConfig_SkipsLocalConfig (0.01s)

All tests pass: make test exits cleanly.
```

## Changes Summary

| File | Change |
|------|--------|
| `internal/git/user.go` | Added `GetNonLocalUserConfig()` - skips local `.git/config` and returns global/includeIf config |
| `internal/git/user_test.go` | Added `TestGetNonLocalUserConfig_SkipsLocalConfig` |
| `cmd/git-workspace/userdetect.go` | Updated `detectUserForRepos()` to call `GetNonLocalUserConfig()` when local source detected, stores underlying config |
| `cmd/git-workspace/list.go` | Updated `displayTable()` to read `GetUserConfig()` at display time; if local, shows effective values with `(local)` marker |

## Display Priority Logic

1. At display time, `GetUserConfig()` is called for each repo
2. If source is `local` → display local values with `(local)` suffix
3. If source is not local → use stored config values (global or includeIf)
4. Drift detection still works for non-local repos

## Persistence Logic

1. `detectUserForRepos()` calls `GetUserConfig()`
2. If result is `UserSourceLocal`, calls `GetNonLocalUserConfig()` instead
3. Stores the underlying global/includeIf values in config
4. Local config is never persisted, only read at display time
