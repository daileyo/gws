# Task 3.0 Proof Artifacts — `--add --recursive` / `-v` (Batch Repository Add)

## CLI: `gws --add --recursive` — Batch Add from Current Directory

```
$ cd /tmp/scan-dir          # directory containing repo-a, repo-b, repo-c
$ gws --add --recursive
Added 3 repositories: repo-a, repo-b, repo-c.
```

**Demonstrates:** Recursive scan of current directory adds all discovered repos with count + names.

---

## CLI: `gws -a -v` — Short-Form Flag Equivalence

```
$ gws -a -v
Added 3 repositories: repo-a, repo-b, repo-c.
```

**Demonstrates:** `-a -v` short forms produce identical results to `--add --recursive`.

---

## CLI: `gws --add --recursive` — Idempotent (All Already Tracked)

```
$ gws --add --recursive
repo-a is already tracked, skipping.
repo-b is already tracked, skipping.
repo-c is already tracked, skipping.
No new repositories found.
```

**Demonstrates:** Second run warns for each tracked repo and exits 0 without adding duplicates.

---

## CLI: `gws --recursive` Without `--add` — Validation Error

```
$ gws --recursive
Error: --recursive/-v requires --add/-a to be set
exit: 1
```

**Demonstrates:** Using `--recursive` without `--add` returns a clear error and exits non-zero.

---

## Test Results: `go test -race -v ./cmd/git-workspace/... -run TestRunAdd`

```
=== RUN   TestRunAdd_CurrentDirectory
Added my-repo to workspace.
--- PASS: TestRunAdd_CurrentDirectory (0.01s)
=== RUN   TestRunAdd_ExplicitExternalPath
Added external-repo to workspace.
Created symlink: [workspace]/external-repo → [external-dir]/external-repo
--- PASS: TestRunAdd_ExplicitExternalPath (0.01s)
=== RUN   TestRunAdd_NonGitDirectory
Error: [path] is not a git repository
--- PASS: TestRunAdd_NonGitDirectory (0.00s)
=== RUN   TestRunAdd_AlreadyTracked
Added my-repo to workspace.
my-repo is already tracked, skipping.
--- PASS: TestRunAdd_AlreadyTracked (0.01s)
=== RUN   TestRunAdd_SymlinkPathConflict
Warning: symlink path [workspace]/conflict-repo already exists, skipping symlink creation
Added conflict-repo to workspace.
--- PASS: TestRunAdd_SymlinkPathConflict (0.01s)
=== RUN   TestRunAdd_NoWorkspace
Error: workspace not initialized
Run gws --init first to create a workspace.
--- PASS: TestRunAdd_NoWorkspace (0.00s)
=== RUN   TestRunAdd_ConfigMetadata
Added typed-repo to workspace.
--- PASS: TestRunAdd_ConfigMetadata (0.01s)
=== RUN   TestRunAddRecursive_MultipleRepos
Added 3 repositories: repo-a, repo-b, repo-c.
--- PASS: TestRunAddRecursive_MultipleRepos (0.03s)
=== RUN   TestRunAddRecursive_AllAlreadyTracked
Added repo-a to workspace.
Added repo-b to workspace.
repo-a is already tracked, skipping.
repo-b is already tracked, skipping.
No new repositories found.
--- PASS: TestRunAddRecursive_AllAlreadyTracked (0.02s)
=== RUN   TestRunAddRecursive_NoRepos
No new repositories found.
--- PASS: TestRunAddRecursive_NoRepos (0.00s)
=== RUN   TestRunAddRecursive_RequiresAdd
--- PASS: TestRunAddRecursive_RequiresAdd (0.00s)
PASS
ok  github.com/daileyo/gws/cmd/git-workspace  1.432s
```

**Demonstrates:** All 11 add tests pass including all 4 recursive-specific cases.

---

## Full Test Suite: `go vet ./... && go test -race ./...`

```
ok  github.com/daileyo/gws/cmd/git-workspace
ok  github.com/daileyo/gws/internal/classifier
ok  github.com/daileyo/gws/internal/config
ok  github.com/daileyo/gws/internal/discovery
ok  github.com/daileyo/gws/internal/filter
ok  github.com/daileyo/gws/internal/git
```

**Demonstrates:** No regressions across the full test suite.
