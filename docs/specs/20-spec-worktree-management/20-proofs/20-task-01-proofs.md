# Task 1.0 Proof Artifacts — Worktree Data Model and Discovery

## Test Results

```
$ go test ./internal/git/ -run "TestListWorktrees|TestParseWorktree|TestIsAligned" -v
=== RUN   TestListWorktrees
--- PASS: TestListWorktrees (0.12s)
=== RUN   TestListWorktrees_NoWorktrees
--- PASS: TestListWorktrees_NoWorktrees (0.05s)
=== RUN   TestListWorktrees_MainWorktreeExcluded
--- PASS: TestListWorktrees_MainWorktreeExcluded (0.08s)
=== RUN   TestParseWorktreeListPorcelain
--- PASS: TestParseWorktreeListPorcelain (0.00s)
=== RUN   TestParseWorktreeListPorcelain_Empty
--- PASS: TestParseWorktreeListPorcelain_Empty (0.00s)
=== RUN   TestIsAligned
=== RUN   TestIsAligned/aligned_-_directly_inside_.wt_dir
=== RUN   TestIsAligned/aligned_-_nested_inside_.wt_dir
=== RUN   TestIsAligned/unaligned_-_different_location
=== RUN   TestIsAligned/unaligned_-_similar_name_but_not_.wt_suffix
=== RUN   TestIsAligned/unaligned_-_sibling_directory_with_similar_prefix
=== RUN   TestIsAligned/aligned_-_trailing_slash_on_repo_path
--- PASS: TestIsAligned (0.00s)
=== RUN   TestListWorktrees_BranchWithSlashes
--- PASS: TestListWorktrees_BranchWithSlashes (0.08s)
PASS

$ go test ./internal/config/ -run TestRepositoryWithWorktrees -v
=== RUN   TestRepositoryWithWorktreesMarshaling
=== RUN   TestRepositoryWithWorktreesMarshaling/worktrees_round-trip_through_JSON
=== RUN   TestRepositoryWithWorktreesMarshaling/worktrees_omitted_from_JSON_when_empty
=== RUN   TestRepositoryWithWorktreesMarshaling/backward_compat_-_old_config_without_worktrees_loads_cleanly
--- PASS: TestRepositoryWithWorktreesMarshaling (0.00s)

$ go test ./cmd/git-workspace/ -run TestRunRefresh_DiscoverWorktrees -v
=== RUN   TestRunRefresh_DiscoverWorktrees
Refreshing workspace at: ...
Detecting git user configuration...
Cleared git status cache
Refresh complete!
Total repositories: 2
Updated 2 repositories
Repositories with worktrees: 1
--- PASS: TestRunRefresh_DiscoverWorktrees (0.06s)
```

## CLI Output

```
$ gws refresh
...
Repositories with worktrees: 1
```

The summary line "Repositories with worktrees: N" is printed when N > 0.

## Full Test Suite

```
$ go test ./...
ok  github.com/daileyo/gws/cmd/git-workspace
ok  github.com/daileyo/gws/internal/classifier
ok  github.com/daileyo/gws/internal/config
ok  github.com/daileyo/gws/internal/discovery
ok  github.com/daileyo/gws/internal/filter
ok  github.com/daileyo/gws/internal/git
ok  github.com/daileyo/gws/internal/user
```

## Verification

- Worktree struct added to config with proper JSON tags and omitempty
- `ListWorktrees` correctly parses `git worktree list --porcelain` output
- `IsAligned` correctly identifies worktrees inside `<repo>.wt/`
- Refresh discovers worktrees and persists them in config
- Backward compatibility preserved (old configs without worktrees load cleanly)
- macOS symlink resolution handled (`/var` → `/private/var`)
