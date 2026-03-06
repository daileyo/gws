# Task 1.0 Proof Artifacts - Dual-Purpose Flag Infrastructure

## Test Results

All tests pass across all packages:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.400s
ok  	github.com/daileyo/gws/internal/filter	0.008s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Dual-Purpose Flag Tests

```
--- PASS: TestFilterFlagsOnListCmd (flags: type/y, visibility/V, tag/t, name/n, path/p, output/o, status/s, show-user/u, remote/r, verbose/v)
--- PASS: TestDualPurposeFlagNoOptDefVal (all 7 dual-purpose flags have NoOptDefVal sentinel)
--- PASS: TestListCmdFlagStacking (-su sets both status and user to sentinel)
--- PASS: TestListCmdFlagStackingWithRemote (-rsu sets remote, status, user)
--- PASS: TestListCmdFlagStackingNewFlags (-yVtp sets type, visibility, tag, path)
--- PASS: TestDualPurposeFlagWithValue (-y=github → "github", --type=github → "github", -y → sentinel)
--- PASS: TestAnyColumnSelected (correctly detects when any Show* flag is true)
```

## Visibility Filter Tests

```
--- PASS: TestByVisibility/Filter_by_private (1 match)
--- PASS: TestByVisibility/Filter_by_unknown (4 matches)
--- PASS: TestByVisibility/Filter_by_public (0 matches)
--- PASS: TestByVisibility/Case_insensitive (PRIVATE → 1 match)
--- PASS: TestByVisibility/Wildcard_priv* (1 match)
--- PASS: TestByVisibility/Wildcard_*known (4 matches)
```

## Deprecated Flag Backward Compatibility

```
--- PASS: TestFilterFlagsOnRootAreHidden (type, name, path, output, status, show-user all hidden)
--- PASS: TestDeprecatedFlagsAreHidden (all 19 deprecated flags hidden)
--- PASS: TestDeprecatedListEmitsWarning
--- PASS: TestDeprecatedMutualExclusivity
```

## Key Implementation Details

- Dual-purpose flags use `NoOptDefVal` with sentinel `\x00show`
- Value syntax requires `=`: `-y=github`, `--type=github`
- POSIX flag stacking works: `-yVtp`, `-rsu`
- New `-V`/`--visibility` flag added
- `filter.Criteria` has new `Visibility` field with `ByVisibility()` function
- `ListOptions` separates `Show*` booleans from `Filter*` values
- `displayTable()` conditionally renders columns based on `Show*` flags
- Deprecated flag dispatch updated to set `ShowType/ShowVisibility/ShowTags/ShowPath=true`
