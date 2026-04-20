# Task 2.0 Proof Artifacts — Worktree Indicator in List Output

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestRepoDisplayName|TestDisplayMultiColumn_Worktree|TestDisplayTable_Worktree|TestDisplayJSON_Worktree" -v
=== RUN   TestRepoDisplayName
=== RUN   TestRepoDisplayName/repo_with_worktrees_shows_(wt)
=== RUN   TestRepoDisplayName/repo_without_worktrees_has_plain_name
--- PASS: TestRepoDisplayName (0.00s)
=== RUN   TestDisplayMultiColumn_WorktreeIndicator
--- PASS: TestDisplayMultiColumn_WorktreeIndicator (0.00s)
=== RUN   TestDisplayTable_WorktreeIndicator
--- PASS: TestDisplayTable_WorktreeIndicator (0.00s)
=== RUN   TestDisplayJSON_WorktreeIndicator
=== RUN   TestDisplayJSON_WorktreeIndicator/repo_with_worktrees_has_has_worktrees_field
=== RUN   TestDisplayJSON_WorktreeIndicator/repo_without_worktrees_omits_has_worktrees_field
--- PASS: TestDisplayJSON_WorktreeIndicator (0.00s)
PASS
```

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

- `repoDisplayName()` returns `"name (wt)"` for repos with worktrees, plain `"name"` for others
- Compact multi-column mode: `(wt)` indicator appears inline next to repo name
- Table mode: `(wt)` indicator appears in NAME column with width accounted for
- JSON mode: `has_worktrees: true` field present for repos with worktrees, omitted otherwise
- Repos without worktrees never show the indicator in any mode
