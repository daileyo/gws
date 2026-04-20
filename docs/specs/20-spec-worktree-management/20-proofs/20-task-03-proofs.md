# Task 3.0 Proof Artifacts — Worktree Navigation (`-wt` flag)

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestRunWorktreeNavigate|TestWorktreeFlag" -v
=== RUN   TestRunWorktreeNavigate_SingleMatch
--- PASS: TestRunWorktreeNavigate_SingleMatch (0.00s)
=== RUN   TestRunWorktreeNavigate_SingleMatch_Quiet
--- PASS: TestRunWorktreeNavigate_SingleMatch_Quiet (0.00s)
=== RUN   TestRunWorktreeNavigate_PartialMatch
--- PASS: TestRunWorktreeNavigate_PartialMatch (0.00s)
=== RUN   TestRunWorktreeNavigate_MultipleMatches_Selection
--- PASS: TestRunWorktreeNavigate_MultipleMatches_Selection (0.00s)
=== RUN   TestRunWorktreeNavigate_MultipleMatches_NonTTY
--- PASS: TestRunWorktreeNavigate_MultipleMatches_NonTTY (0.00s)
=== RUN   TestRunWorktreeNavigate_NoMatch
--- PASS: TestRunWorktreeNavigate_NoMatch (0.00s)
=== RUN   TestRunWorktreeNavigate_BareFlag_ListAll
--- PASS: TestRunWorktreeNavigate_BareFlag_ListAll (0.00s)
=== RUN   TestRunWorktreeNavigate_UnalignedIndicator
--- PASS: TestRunWorktreeNavigate_UnalignedIndicator (0.00s)
=== RUN   TestRunWorktreeNavigate_RepoNoWorktrees
--- PASS: TestRunWorktreeNavigate_RepoNoWorktrees (0.00s)
=== RUN   TestRunWorktreeNavigate_RepoNotFound
--- PASS: TestRunWorktreeNavigate_RepoNotFound (0.00s)
=== RUN   TestRunWorktreeNavigate_BranchWithSlashes
--- PASS: TestRunWorktreeNavigate_BranchWithSlashes (0.00s)
=== RUN   TestWorktreeFlagRegistered
--- PASS: TestWorktreeFlagRegistered (0.00s)
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

- `--worktree` flag registered on root command with NoOptDefVal sentinel
- No shorthand used (`-w` is taken by deprecated `--print-workspace`)
- Shell function will provide `-wt` shorthand (Task 5.0)
- Single worktree match prints path to stdout
- Multiple matches display numbered selection list with branch name, path, and `(unaligned)` indicator
- No match displays error with list of available worktrees
- Bare `--worktree` (no value) lists all worktrees for interactive selection
- Partial/substring matching works for branch names (e.g., "auth" matches "feature-auth")
- Branch names with slashes work (e.g., "urgent" matches "hotfix/urgent")
- Repo with no worktrees returns clear error
- Non-existent repo returns no-match error with suggestions
- Non-TTY mode prints all matching paths to stdout
