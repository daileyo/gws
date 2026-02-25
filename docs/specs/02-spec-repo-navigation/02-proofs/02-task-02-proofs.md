# 02-task-02 Proof Artifacts: Multiple Match Interactive Selection

## Test Results

All multiple match and wildcard tests pass:

```
--- PASS: TestRunNavigate_MultipleMatches_Selection (0.00s)
--- PASS: TestRunNavigate_MultipleMatches_SelectionQuiet (0.00s)
--- PASS: TestRunNavigate_MultipleMatches_InvalidThenValid (0.00s)
--- PASS: TestRunNavigate_MultipleMatches_OutOfRange (0.00s)
--- PASS: TestRunNavigate_MultipleMatches_NonTTY (0.00s)
--- PASS: TestRunNavigate_WildcardSingleMatch (0.00s)
--- PASS: TestRunNavigate_WildcardMultipleMatches (0.00s)
--- PASS: TestRunNavigate_WildcardQuestionMark (0.00s)
--- PASS: TestRunNavigate_WildcardInteractiveSelection (0.00s)
--- PASS: TestIsTerminal_NonFile (0.00s)
--- PASS: TestIsTerminal_BytesBuffer (0.00s)
```

All packages pass:

```
ok  	github.com/daileyo/gws/cmd/gws         (cached)
ok  	github.com/daileyo/gws/internal/classifier  (cached)
ok  	github.com/daileyo/gws/internal/config      (cached)
ok  	github.com/daileyo/gws/internal/discovery    (cached)
ok  	github.com/daileyo/gws/internal/filter       (cached)
ok  	github.com/daileyo/gws/internal/git          (cached)
```

## Implementation Summary

### TTY Detection
- `isTerminal()` checks if reader is `*os.File` with `ModeCharDevice` flag
- `isTerminalFunc` package variable allows test overrides without adding dependencies
- `withTTY(t)` and `withNonTTY(t)` test helpers use `t.Cleanup()` for safe restoration

### Interactive Selection (TTY)
- Displays numbered list to stderr: `  1) name (type) path`
- Prompts `Select repository [1-N]:` on stderr
- Validates input: re-prompts on invalid (non-numeric or out-of-range)
- Max 3 attempts before returning error
- Selection respects verbose/quiet mode for final output

### Non-TTY Behavior
- Prints all matching paths (one per line) to stdout
- Returns error with match count for non-zero exit status
- No prompting or interactive elements

### Wildcard Support
- Reuses `filter.MatchesPattern()` for `*` and `?` glob patterns
- Works for single match (direct navigation) and multiple match (selection)
- Case-insensitive matching consistent with filter behavior

## Quality Gates

```
$ go vet ./...
(no output — clean)
```

## Verification

- [x] TTY detection via `os.File.Stat()` with `ModeCharDevice`
- [x] No new external dependencies added
- [x] Numbered list displayed to stderr
- [x] User selection validated with re-prompt on invalid input
- [x] Max 3 attempts enforced
- [x] Non-TTY prints all paths to stdout with non-zero exit
- [x] Wildcard `*` and `?` patterns work in navigation queries
- [x] Quiet mode respected after selection
