# Task 2.0 Proof Artifacts — `--add` Command (Single Repository with Symlink Support)

## CLI: `gws --help` — `--add` Flag Registered

```
Commands:
  -a, --add string[="."]   Add a git repository to the workspace (defaults to current directory)
  -d, --add-tag            Add a tag to repositories (args: <repo> <tag>)
  -i, --init               Initialize a gws workspace in the current directory
  -l, --list               List all tracked repositories
```

**Demonstrates:** `--add` / `-a` registered with optional path value (default `.`); `--add-tag` correctly carries `-d`.

---

## CLI: `gws --add` — Add from Current Directory (Internal Repo)

```
$ cd /workspace/internal-repo
$ gws --add
Added internal-repo to workspace.
```

**Demonstrates:** Adding a repo from the current directory (inside workspace — no symlink created).

---

## CLI: `gws --add /path/to/external/repo` — External Repo with Symlink

```
$ gws --add /tmp/external/external-repo
Added external-repo to workspace.
Created symlink: /workspace/external-repo → /tmp/external/external-repo
```

**Demonstrates:** Adding an external repo creates a symlink in the workspace directory.

---

## Filesystem: Symlink Visible in Workspace

```
$ ls -la /workspace
lrwxr-xr-x  external-repo -> /tmp/external/external-repo
drwxr-xr-x  internal-repo
```

**Demonstrates:** Symlink exists on disk pointing to external repo; internal repo is a real directory.

---

## CLI: `gws --add` — Duplicate Detection

```
$ gws --add /tmp/external/external-repo
external-repo is already tracked, skipping.
```

**Demonstrates:** Second add of same repo warns and exits 0 without duplicating the config entry.

---

## CLI: `gws --add` — Non-Git Directory Error

```
$ gws --add /tmp/not-a-git-dir
Error: /tmp/not-a-git-dir is not a git repository
exit: 1
```

**Demonstrates:** Adding a non-git directory prints error and exits non-zero.

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
PASS
ok  github.com/daileyo/gws/cmd/git-workspace  1.581s
```

**Demonstrates:** All 7 add tests pass including: current dir add, external path with symlink, non-git error, duplicate detection, symlink conflict, no-workspace guard, and metadata (type/remote URL) extraction.

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
