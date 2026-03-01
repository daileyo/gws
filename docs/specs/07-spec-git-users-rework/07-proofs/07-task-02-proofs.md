# 07 Task 2.0 Proof Artifacts - Replace Heuristic IncludeIf with Proper Parsing

## Test Results

```
=== RUN   TestMatchesGitdirCondition
=== RUN   TestMatchesGitdirCondition/exact_directory_match_with_trailing_slash
=== RUN   TestMatchesGitdirCondition/nested_repo_under_matching_directory
=== RUN   TestMatchesGitdirCondition/trailing_**_glob
=== RUN   TestMatchesGitdirCondition/non-matching_path
=== RUN   TestMatchesGitdirCondition/tilde_expansion
=== RUN   TestMatchesGitdirCondition/tilde_expansion_non-match
=== RUN   TestMatchesGitdirCondition/case-insensitive_gitdir/i_match
=== RUN   TestMatchesGitdirCondition/case-sensitive_gitdir_does_not_match_different_case
=== RUN   TestMatchesGitdirCondition/not_a_gitdir_condition
=== RUN   TestMatchesGitdirCondition/directory_without_trailing_slash
--- PASS: TestMatchesGitdirCondition (0.00s)

=== RUN   TestParseIncludeIfs
--- PASS: TestParseIncludeIfs (0.00s)

=== RUN   TestParseIncludeIfs_Empty
--- PASS: TestParseIncludeIfs_Empty (0.00s)

All tests pass: make test exits cleanly.
```

## Changes Summary

| File | Change |
|------|--------|
| `internal/git/user.go` | Added `MatchesGitdirCondition()`, `checkIncludeIfMatch()`, `parseIncludeIfs()`, `parseIncludeIfKeyValue()`. Removed `isLikelyIncludeIfPath()`. Updated `GetUserConfig()` to use proper includeIf matching |
| `internal/git/user_test.go` | Replaced `TestIsLikelyIncludeIfPath` with `TestMatchesGitdirCondition` (10 cases) and `TestParseIncludeIfs` tests |

## Implementation Notes

- Avoided import cycle (`git` → `user` → `git`) by implementing gitconfig includeIf parsing directly in the `git` package via `parseIncludeIfs()` rather than importing `internal/user`
- `MatchesGitdirCondition()` supports: `gitdir:`, `gitdir/i:` (case-insensitive), `~/` expansion, trailing `**` globs, with/without trailing slash
- `checkIncludeIfMatch()` reads `~/.gitconfig`, parses includeIf directives, matches against repo path, and reads the included config for user info
