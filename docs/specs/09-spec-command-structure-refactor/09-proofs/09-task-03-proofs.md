# Task 3.0 Proof Artifacts - Scope Filter Flags to List Subcommand

## Filter Flags Scoped to listCmd

Filter flags are now registered only on `listCmd`, not on `rootCmd`:

| Flag | Short | Registered On |
|------|-------|---------------|
| `--type` | `-y` | `listCmd` only |
| `--tag` | `-t` | `listCmd` AND `rootCmd` (root keeps it for user operations) |
| `--name` | `-n` | `listCmd` only |
| `--path` | `-p` | `listCmd` only |
| `--output` | `-o` | `listCmd` only |
| `--status` | `-s` | `listCmd` only |
| `--show-user` | `-u` | `listCmd` only |

## POSIX Flag Stacking

```
$ go test -run TestListCmdFlagStacking -v
=== RUN   TestListCmdFlagStacking
--- PASS: TestListCmdFlagStacking (0.00s)
```

`gws list -su` correctly sets both `showStatus` and `showUser` to true.

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestFilterFlags|TestTagFlag|TestTagRequires|TestListCmd" -v

=== RUN   TestFilterFlagsOnListCmd
--- PASS: TestFilterFlagsOnListCmd (0.00s)  (7 sub-tests)
=== RUN   TestFilterFlagsNotOnRoot
--- PASS: TestFilterFlagsNotOnRoot (0.00s)  (6 sub-tests)
=== RUN   TestListCmdFlagStacking
--- PASS: TestListCmdFlagStacking (0.00s)
=== RUN   TestTagFlagOnRoot
--- PASS: TestTagFlagOnRoot (0.00s)
=== RUN   TestTagRequiresListOrUser
--- PASS: TestTagRequiresListOrUser (0.00s)
PASS
```

## Changes Made

1. **Moved filter variables** (`filterType`, `filterName`, `filterPath`, `outputFormat`, `showStatus`, `showUser`) from `main.go` to `list.go`
2. **Kept `filterTags`** on root (`main.go`) — still needed by user operations (specs 10/11)
3. **Removed `hasNonTagFilterFlags`** helper function from `main.go`
4. **Removed filter validation** from root `RunE` (only `--tag` validation remains)
5. **Removed `filterFlagUsages`** template function and "List Filters" section from usage template
6. **Simplified root `--list` alias** dispatch to only pass `filterTags` and default `outputFormat`
7. **Created `list_test.go`** with `TestFilterFlagsOnListCmd`, `TestFilterFlagsNotOnRoot`, `TestListCmdFlagStacking`
8. **Updated `main_test.go`**: replaced `TestFilterFlagsRegistered` with `TestTagFlagOnRoot`, replaced `TestFilterFlagsRequireList` with `TestTagRequiresListOrUser`

## Full Suite

```
$ make vet && go test ./...
go vet ./...  ✓
ok  github.com/daileyo/gws/cmd/git-workspace   0.407s
ok  github.com/daileyo/gws/internal/classifier  (cached)
ok  github.com/daileyo/gws/internal/config      (cached)
ok  github.com/daileyo/gws/internal/discovery    (cached)
ok  github.com/daileyo/gws/internal/filter       (cached)
ok  github.com/daileyo/gws/internal/git          (cached)
ok  github.com/daileyo/gws/internal/user         (cached)
```
