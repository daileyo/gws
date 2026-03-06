# Task 5.0 Proof Artifacts - JSON Output Column Selection and Help Text Updates

## Test Results

All tests pass (including new JSON column selection and deprecated compat tests):

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.407s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)
```

Race detector clean:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	1.472s
```

## CLI: `gws list -h` shows all flags

```
List all repositories tracked by gws with optional filtering.

By default, only repository names are shown in a compact multi-column layout.
Use flags to add columns or filter results. A flag without a value shows the
column; a flag with a value filters by that value AND shows the column.

Examples:
  gws list                          # Multi-column repo names
  gws list -v                       # Table with type, visibility, tags, path
  gws list -vv                      # Table with all columns (status, user, remote, etc.)
  gws list -yVtp                    # Show type, visibility, tags, path columns
  gws list -y=github                # Filter by GitHub, show type column
  gws list -y=github -V=private     # Filter by GitHub + private
  gws list -y=github -tp            # Filter by GitHub, show type, tags, path
  gws list -n "api-*"               # Filter by name pattern
  gws list -su                      # Show status and user columns
  gws list -r                       # Show remote URL column
  gws list -o json                  # Output as JSON

Note: dual-purpose flags use = for values (e.g., -y=github, --type=github).
Without =, the flag shows the column without filtering.

Usage:
  git-workspace list [flags]

Flags:
  -h, --help                          help for list
  -n, --name string                   Filter by repository name (partial match, supports wildcards)
  -o, --output string                 Output format: table, json (default "table")
  -p, --path  Show path column, or filter by path pattern
  -r, --remote  Show remote URL column, or filter by remote URL pattern
  -u, --show-user  Show user/email/sign columns, or filter by user name
  -s, --status  Show git status column, or filter by status pattern
  -t, --tag  Show tags column, or filter by tag value (repeatable for AND logic)
  -y, --type  Show type column, or filter by type value
  -v, --verbose count                 Verbose output (-v stored data, -vv all columns)
  -V, --visibility  Show visibility column, or filter by visibility value
```

All flags visible: `-y`, `-V`, `-t`, `-p`, `-s`, `-u`, `-r`, `-n`, `-o`, `-v`. NoOptDefVal artifacts cleaned.

## Changes Made

### `list.go` - displayJSON() rewrite

Replaced struct marshaling with dynamic `map[string]interface{}` that respects `Show*` flags:
- Always includes `name`
- `ShowType` -> `type`, `ShowVisibility` -> `visibility`, `ShowTags` -> `tags`, `ShowPath` -> `path`
- `ShowStatus` -> `status` (formatted as string via `formatStatusShort`)
- `ShowUser` -> `user`, `email`, `signing_enabled`
- `ShowRemote` -> `remote_url`, `has_multiple_remotes`
- Default (no column flags) outputs JSON with only `name`

### `list.go` - Help text fix

Replaced `SetUsageFunc` with `SetHelpFunc` to avoid duplicate Long description bug. The custom help function:
- Strips `string[="..."]` NoOptDefVal artifacts from flag usage strings
- Renders Long description + Usage + Flags in a clean format
- No longer causes Cobra to double-print the description

### `list_test.go` - New tests

- `TestDisplayJSON_ColumnSelection`: 4 test cases verifying JSON field selection:
  - Default: only `name` in output
  - `-yp`: `name`, `type`, `path` only
  - `-v`: all stored-data fields
  - `-u`: `name`, `user`, `email`, `signing_enabled`
- `TestDeprecatedListDispatchPopulatesListOptions`: Verifies deprecated dispatch sets stored-data columns on, live columns off

### Deprecated backward compatibility

Verified `deprecated.go` dispatch at line 248 correctly populates:
- `ShowType/ShowVisibility/ShowTags/ShowPath = true` (old default)
- `ShowStatus/ShowUser/ShowRemote = false`
- `depWarnings` map unchanged (no flag names changed)
