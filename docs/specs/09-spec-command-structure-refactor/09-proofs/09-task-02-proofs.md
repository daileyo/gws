# Task 2.0 Proof Artifacts - Create Core Subcommands

## Subcommands Registered

The following Cobra subcommands are now registered on `rootCmd`:

- `list` — `cmd/git-workspace/list.go`
- `init` — `cmd/git-workspace/init.go`
- `add` — `cmd/git-workspace/add.go`
- `refresh` — `cmd/git-workspace/refresh.go`
- `print-workspace` — `cmd/git-workspace/main.go`

## Root Alias Flags (Hidden)

Hidden alias flags on root delegate to subcommand logic:

| Flag | Short | Hidden | Delegates To |
|------|-------|--------|-------------|
| `--list` | `-l` | Yes | `listCmd` |
| `--init` | `-i` | Yes | `initCmd` |
| `--add` | `-a` | Yes | `addCmd` |
| `--recursive` | `-v` | Yes | `addCmd` (with `--add`) |
| `--refresh` | `-r` | Yes | `refreshCmd` |
| `--print-workspace` | `-w` | Yes | `printWorkspaceCmd` |

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestSubcommandsRegistered|TestAliasFlagsAreHidden|TestMutualExclusivity|TestUserFlagValidation|TestCommandFlagsRegistered" -v

=== RUN   TestCommandFlagsRegistered
--- PASS: TestCommandFlagsRegistered (0.00s)
=== RUN   TestSubcommandsRegistered
=== RUN   TestSubcommandsRegistered/list
=== RUN   TestSubcommandsRegistered/init
=== RUN   TestSubcommandsRegistered/add
=== RUN   TestSubcommandsRegistered/refresh
=== RUN   TestSubcommandsRegistered/print-workspace
--- PASS: TestSubcommandsRegistered (0.00s)
=== RUN   TestAliasFlagsAreHidden
=== RUN   TestAliasFlagsAreHidden/list
=== RUN   TestAliasFlagsAreHidden/init
=== RUN   TestAliasFlagsAreHidden/add
=== RUN   TestAliasFlagsAreHidden/recursive
=== RUN   TestAliasFlagsAreHidden/refresh
=== RUN   TestAliasFlagsAreHidden/print-workspace
--- PASS: TestAliasFlagsAreHidden (0.00s)
=== RUN   TestMutualExclusivity
--- PASS: TestMutualExclusivity (0.00s)
=== RUN   TestUserFlagValidation
--- PASS: TestUserFlagValidation (0.00s)
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	0.009s
```

## Full Test Suite

```
$ make vet && make test

Running go vet...
go vet ./...
Running tests...
go test -v ./...
ok  github.com/daileyo/gws/cmd/git-workspace   0.358s (67 tests, all pass)
ok  github.com/daileyo/gws/internal/classifier  (cached)
ok  github.com/daileyo/gws/internal/config      (cached)
ok  github.com/daileyo/gws/internal/discovery    (cached)
ok  github.com/daileyo/gws/internal/filter       (cached)
ok  github.com/daileyo/gws/internal/git          (cached)
ok  github.com/daileyo/gws/internal/user         (cached)
```

## Key Changes

1. **New test: `TestSubcommandsRegistered`** — Verifies `list`, `init`, `add`, `refresh`, `print-workspace` are registered as Cobra subcommands
2. **New test: `TestAliasFlagsAreHidden`** — Verifies alias flags `--list`, `--init`, `--add`, `--recursive`, `--refresh`, `--print-workspace` are hidden
3. **Updated `TestMutualExclusivity`** — Error message updated to "only one command can be used at a time"
4. **Updated `TestUserFlagValidation`** — Mutual exclusivity sub-tests updated to match new error message
5. **Each subcommand** has its own `init()` function registering it with `rootCmd.AddCommand()`
6. **Filter flags** registered on both `rootCmd` and `listCmd` (sharing same global variables) for backward compatibility
7. **`addCmd`** has its own `addRecursive` variable and `--recursive`/`-v` flag (separate from root's `flagRecursive`)
8. **Usage template** updated with `Available Commands` section showing subcommands
