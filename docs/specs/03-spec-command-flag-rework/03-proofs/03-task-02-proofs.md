# 03 Task 2.0 Proof Artifacts - Add Shorthands to List Filter Flags

## CLI Output - Filter Flag Shorthands

```
$ gws -h (flags section)
  -a, --add-tag           Add a tag to repositories (args: <repo> <tag>)
  -h, --help              help for gws
  -i, --init string       Initialize workspace by scanning directory
  -l, --list              List all tracked repositories
  -n, --name string       Filter by repository name (partial match)
  -o, --output string     Output format: table, json (default "table")
  -p, --path string       Filter by repository path (partial match)
  -w, --print-workspace   Print workspace path (for shell integration)
  -r, --refresh           Refresh repository metadata and git status cache
  -u, --remove-tag        Remove a tag from repositories (args: <repo> <tag>)
  -s, --status            Show git status (branch, clean/dirty, ahead/behind)
  -t, --tag strings       Filter by custom tag(s) - can be specified multiple times for AND logic
  -y, --type string       Filter by repository type (github, gitlab, ado, bitbucket)
  -v, --version           version for gws
```

All filter flags now have shorthands: `-y` (type), `-t` (tag), `-n` (name), `-p` (path), `-o` (output), `-s` (status).

## Test Results

```
$ go test -v ./cmd/gws/ -run TestFilterFlags
=== RUN   TestFilterFlagsRegistered
=== RUN   TestFilterFlagsRegistered/type
=== RUN   TestFilterFlagsRegistered/tag
=== RUN   TestFilterFlagsRegistered/name
=== RUN   TestFilterFlagsRegistered/path
=== RUN   TestFilterFlagsRegistered/output
=== RUN   TestFilterFlagsRegistered/status
--- PASS: TestFilterFlagsRegistered (0.00s)
=== RUN   TestFilterFlagsRequireList
--- PASS: TestFilterFlagsRequireList (0.00s)
PASS
```

## Verification

- All 6 filter flags have correct shorthands registered (-y, -t, -n, -p, -o, -s)
- Filter flags without `--list` produce error: "filter flags (--type, --tag, --name, --path, --output, --status) require --list/-l to be set"
- `hasFilterFlags()` helper uses `cmd.Flags().Changed()` to detect explicitly set flags (ignores defaults)
- All existing tests continue to pass
