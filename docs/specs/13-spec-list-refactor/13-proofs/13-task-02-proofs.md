# Task 2.0 Proof Artifacts - Multi-Column Default View

## Test Results

All tests pass:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.417s
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Multi-Column Tests

```
--- PASS: TestDisplayMultiColumn_NonTTY/empty_list
--- PASS: TestDisplayMultiColumn_NonTTY/single_repo
--- PASS: TestDisplayMultiColumn_NonTTY/multiple_repos_sorted (alphabetical: charlie,alpha,bravo → alpha,bravo,charlie)
--- PASS: TestDisplayMultiColumn_NonTTY/already_sorted
--- PASS: TestGetTerminalWidth (returns 80 default in non-TTY)
```

## Key Implementation Details

- `displayMultiColumn()` sorts names alphabetically, fills left-to-right with 2-space gap
- `term.IsTerminal()` detects TTY; non-TTY outputs one name per line
- `getTerminalWidth()` returns terminal width or 80 as default
- Column count calculated as `termWidth / (maxNameLen + 2)`
- `displayTable()` header "Found N repositories" shown before multi-column output
- `runList()` routes to `displayMultiColumn()` when no `Show*` flags and `VerboseLevel == 0`
- Added `golang.org/x/term` dependency
- Fixed pre-existing `go vet` issue in `userupdate.go` (Go 1.24 stricter checks)
