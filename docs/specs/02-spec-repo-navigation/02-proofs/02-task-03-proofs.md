# 02-task-03 Proof Artifacts: No Match with Suggestions

## Test Results

All suggestion tests pass:

```
--- PASS: TestRunNavigate_NoMatch (0.00s)
--- PASS: TestRunNavigate_NoMatch_WithSuggestions (0.00s)
--- PASS: TestRunNavigate_NoMatch_SuggestsPartialOverlap (0.00s)
--- PASS: TestRunNavigate_NoMatch_NoSuggestions (0.00s)
--- PASS: TestFindSuggestions_MaxLimit (0.00s)
--- PASS: TestFindSuggestions_SubstringMatching (0.00s)
--- PASS: TestFindSuggestions_ShortQuery (0.00s)
```

All packages pass:

```
ok  	github.com/daileyo/gws/cmd/gws         0.346s
ok  	github.com/daileyo/gws/internal/classifier  (cached)
ok  	github.com/daileyo/gws/internal/config      (cached)
ok  	github.com/daileyo/gws/internal/discovery    (cached)
ok  	github.com/daileyo/gws/internal/filter       (cached)
ok  	github.com/daileyo/gws/internal/git          (cached)
```

## Implementation Summary

### `findSuggestions()`
- Checks if any 2+ character substring of the query appears in each repo name (case-insensitive)
- Returns up to `max` (default 5) suggestions
- Single-character queries return no suggestions (no meaningful overlap)

### `handleNoMatch()` Updated
- Prints error message: `No repositories found matching '<query>'`
- Calls `findSuggestions()` and if non-empty, prints `Did you mean:` with indented suggestions
- Returns error for non-zero exit status

### `isSimilar()` Helper
- Extracts all substrings of length >= 2 from the query
- Checks if any appear in the candidate name
- Case-insensitive via pre-lowered strings

## Expected Output Examples

With suggestions:
```
No repositories found matching 'frx'

Did you mean:
  frontend
```

Without suggestions:
```
No repositories found matching 'zq'
```

## Quality Gates

```
$ go vet ./...
(no output — clean)
```

## Verification

- [x] Error message displayed for no match
- [x] Suggestions shown when similar names exist (substring overlap)
- [x] No suggestions shown when query has no overlap
- [x] Max 5 suggestions enforced
- [x] Single-character queries produce no suggestions
- [x] Non-zero exit status on no match
