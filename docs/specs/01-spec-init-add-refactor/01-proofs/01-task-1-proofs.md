# Task 1.0 Proof Artifacts — Refactor `--init` Command and Reassign Flag Short Forms

## CLI: `gws --help` — Flag Signature Changes

`--init` is now a boolean flag (no path argument). `--add-tag` uses `-d` as its short form.

```
$ gws --help

Commands:
  -d, --add-tag           Add a tag to repositories (args: <repo> <tag>)
  -i, --init              Initialize a gws workspace in the current directory
  -l, --list              List all tracked repositories
  -w, --print-workspace   Print workspace path (for shell integration)
  -r, --refresh           Refresh repository metadata and git status cache
  -u, --remove-tag        Remove a tag from repositories (args: <repo> <tag>)
```

**Demonstrates:** `--init` has no path argument; `--add-tag` short form is now `-d`.

---

## CLI: `gws --init` — Successful Initialization

Run in a directory containing a git repository.

```
$ gws --init
Initialized workspace at: /tmp/proof-workspace
Found 1 repository.
```

**Demonstrates:** Boolean-flag init discovers repos in current directory and prints concise confirmation.

---

## CLI: `gws --init` — Already-Initialized Guard

Run a second time in the same directory.

```
$ gws --init
Workspace already initialized at: /tmp/proof-workspace

To add more repositories:  gws --add
To re-scan the workspace:  gws --refresh
```

**Demonstrates:** Already-initialized guard exits 0 with workspace path and suggestions for `--add` and `--refresh`.

---

## Test Results: `go test -race -v ./cmd/git-workspace/... -run "TestRunInit|TestCommandFlagsRegistered|TestMutualExclusivity"`

```
=== RUN   TestRunInit_HappyPath
--- PASS: TestRunInit_HappyPath (0.01s)
=== RUN   TestRunInit_AlreadyInitialized
--- PASS: TestRunInit_AlreadyInitialized (0.00s)
=== RUN   TestRunInit_EmptyDirectory
--- PASS: TestRunInit_EmptyDirectory (0.00s)
=== RUN   TestCommandFlagsRegistered
=== RUN   TestCommandFlagsRegistered/list
=== RUN   TestCommandFlagsRegistered/init
=== RUN   TestCommandFlagsRegistered/add-tag
=== RUN   TestCommandFlagsRegistered/remove-tag
=== RUN   TestCommandFlagsRegistered/refresh
=== RUN   TestCommandFlagsRegistered/print-workspace
=== RUN   TestCommandFlagsRegistered/go
=== RUN   TestCommandFlagsRegistered/quiet
--- PASS: TestCommandFlagsRegistered (0.00s)
    --- PASS: TestCommandFlagsRegistered/list (0.00s)
    --- PASS: TestCommandFlagsRegistered/init (0.00s)
    --- PASS: TestCommandFlagsRegistered/add-tag (0.00s)
    --- PASS: TestCommandFlagsRegistered/remove-tag (0.00s)
    --- PASS: TestCommandFlagsRegistered/refresh (0.00s)
    --- PASS: TestCommandFlagsRegistered/print-workspace (0.00s)
    --- PASS: TestCommandFlagsRegistered/go (0.00s)
    --- PASS: TestCommandFlagsRegistered/quiet (0.00s)
=== RUN   TestMutualExclusivity
--- PASS: TestMutualExclusivity (0.00s)
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	1.429s
```

**Demonstrates:** All init behavioral tests pass; flag registration tests confirm `-d` for `--add-tag`; mutual exclusivity test updated and passing.

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
