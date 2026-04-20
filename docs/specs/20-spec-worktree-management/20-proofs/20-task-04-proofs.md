# Task 4.0 Proof Artifacts — Worktree Subcommands (`list`, `add`, `align`)

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestRunWorktreeList|TestRunWorktreeAdd|TestRunWorktreeAlign" -v
=== RUN   TestRunWorktreeList_AllWorktrees
--- PASS: TestRunWorktreeList_AllWorktrees (0.00s)
=== RUN   TestRunWorktreeList_FilterByRepo
--- PASS: TestRunWorktreeList_FilterByRepo (0.00s)
=== RUN   TestRunWorktreeList_NoWorktrees
--- PASS: TestRunWorktreeList_NoWorktrees (0.00s)
=== RUN   TestRunWorktreeList_NoWorktreesForRepo
--- PASS: TestRunWorktreeList_NoWorktreesForRepo (0.00s)
=== RUN   TestRunWorktreeList_UnalignedIndicator
--- PASS: TestRunWorktreeList_UnalignedIndicator (0.00s)
=== RUN   TestRunWorktreeAdd_Success
--- PASS: TestRunWorktreeAdd_Success (0.09s)
=== RUN   TestRunWorktreeAdd_CreatesWtDir
--- PASS: TestRunWorktreeAdd_CreatesWtDir (0.08s)
=== RUN   TestRunWorktreeAdd_UnknownRepo
--- PASS: TestRunWorktreeAdd_UnknownRepo (0.04s)
=== RUN   TestRunWorktreeAdd_DuplicateBranch
--- PASS: TestRunWorktreeAdd_DuplicateBranch (0.09s)
=== RUN   TestRunWorktreeAdd_BranchWithSlash
--- PASS: TestRunWorktreeAdd_BranchWithSlash (0.09s)
=== RUN   TestRunWorktreeAlign_MovesUnaligned
--- PASS: TestRunWorktreeAlign_MovesUnaligned (0.10s)
=== RUN   TestRunWorktreeAlign_SkipsAligned
--- PASS: TestRunWorktreeAlign_SkipsAligned (0.00s)
=== RUN   TestRunWorktreeAlign_DryRun
--- PASS: TestRunWorktreeAlign_DryRun (0.07s)
=== RUN   TestRunWorktreeAlign_DuplicateNames
--- PASS: TestRunWorktreeAlign_DuplicateNames (0.11s)
=== RUN   TestRunWorktreeAlign_FilterByRepo
--- PASS: TestRunWorktreeAlign_FilterByRepo (0.09s)

$ go test ./internal/git/ -run "TestAddWorktree|TestMoveWorktree" -v
=== RUN   TestAddWorktree
--- PASS: TestAddWorktree (0.10s)
=== RUN   TestAddWorktree_ExistingBranch
--- PASS: TestAddWorktree_ExistingBranch (0.09s)
=== RUN   TestMoveWorktree
--- PASS: TestMoveWorktree (0.09s)
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

### `gws worktree list`
- Lists all worktrees across repos in a table with REPO, BRANCH, PATH, STATUS columns
- Filters by repo name when argument provided
- Shows `(unaligned)` for worktrees outside `.wt/`
- Shows "No worktrees found" when none exist

### `gws worktree add`
- Creates worktree at `<repo>.wt/<branch>` using `git worktree add`
- Creates `.wt/` directory automatically if needed
- Handles both new branches (`-b`) and existing branches
- Handles branch names with slashes (e.g., `feature/auth-flow`)
- Updates config with new worktree entry
- Errors on unknown repo or duplicate branch

### `gws worktree align`
- Moves unaligned worktrees into `<repo>.wt/` using `git worktree move`
- Skips already-aligned worktrees
- `--dry-run` previews moves without executing
- Handles naming conflicts with `-dup-NN` suffix
- Filters by repo name when argument provided
- Updates config after moves
- Creates `.wt/` directory and parent directories for branches with slashes

### `internal/git` operations
- `AddWorktree` creates worktrees for new and existing branches
- `MoveWorktree` moves worktrees and updates git internal pointers
