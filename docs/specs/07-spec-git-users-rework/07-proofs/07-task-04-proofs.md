# 07 Task 4.0 Proof Artifacts - IncludeIf Profile Auto-Linking

## Test Results

```
=== RUN   TestMatchProfileByUser
=== RUN   TestMatchProfileByUser/exact_email_match
=== RUN   TestMatchProfileByUser/email_match_with_different_name
=== RUN   TestMatchProfileByUser/case-insensitive_email_match
=== RUN   TestMatchProfileByUser/name_match_when_no_email
=== RUN   TestMatchProfileByUser/no_match
=== RUN   TestMatchProfileByUser/empty_profiles_list
=== RUN   TestMatchProfileByUser/nil_profiles
--- PASS: TestMatchProfileByUser (0.00s)

All tests pass: make test exits cleanly.
```

## Changes Summary

| File | Change |
|------|--------|
| `internal/user/profile.go` | Added `MatchProfileByUser()` - matches profiles by email (primary) or git name (secondary), case-insensitive |
| `internal/user/profile_test.go` | Added `TestMatchProfileByUser` with 7 table-driven cases |
| `cmd/git-workspace/userdetect.go` | Updated `detectUserForRepos()` to accept variadic profiles param, calls `MatchProfileByUser()` for includeIf users |
| `cmd/git-workspace/refresh.go` | Passes `cfg.Profiles` to `detectUserForRepos()` |
| `cmd/git-workspace/add.go` | Passes `cfg.Profiles` to `detectUserForRepos()` for single and recursive add |
