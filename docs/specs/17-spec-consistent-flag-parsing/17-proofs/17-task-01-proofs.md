# Task 1.0-5.0 Proof Artifacts: Consistent Flag Parsing

All tasks implemented together in batch mode.

## CLI Output

### Lowercase filter with space-separated value (-t ai)

```
$ gws list -t ai
Found 37 repositories:

action-tag-release
cct-ww-internal
danswer
...
```

Filter applied, no tag column displayed. Demonstrates lowercase=filter-only.

### Uppercase bare flag (-T)

```
$ gws list -T
Found 218 repositories:

NAME                                                     TAGS
-------------------------------------------------------  ---------------------------
bgs-feedback-service                                     client, aa
```

Tag column shown, no filter. Demonstrates uppercase bare=show column.

### Uppercase with value (-T ai)

```
$ gws list -T ai
Found 37 repositories:

NAME                                     TAGS
---------------------------------------  --------
productcore-shell                        grainger
```

Tag column shown AND filtered by "ai". Demonstrates uppercase=show+filter with space-separated value.

### Equals backward compatibility (-t=ai)

```
$ gws list -t=ai
Found 37 repositories:

action-tag-release
cct-ww-internal
...
```

Same filter result as `-t ai`. Demonstrates = syntax still works.

### Uppercase stacking (-YTSP)

```
$ gws list -YTSP
Found 218 repositories:

NAME                                      STATUS                                TYPE      TAGS                  PATH
```

All four columns displayed. Demonstrates POSIX stacking with bare uppercase flags.

### New visibility flag (-i private)

```
$ gws list -i private
Found 139 repositories:

CopilotChat.nvim
Hangar
...
```

Filtered by visibility without showing column. Demonstrates new -i flag.

### New remote-raw flag (-b github)

```
$ gws list -b github
Found 218 repositories:
```

Filtered by remote-raw pattern. Demonstrates new -b flag.

### Show remote column (-R)

```
$ gws list -R
Found 218 repositories:

NAME                                                     REMOTE
-------------------------------------------------------  --------
bgs-feedback-service                                     https://github-aa/...
```

-R now shows remote column (not remote-raw). Demonstrates remapping.

### Deprecated -V warning

```
$ gws list -V
Warning: -V is deprecated, use '-i' (filter) or '-I' (show column) instead
Found 218 repositories:

NAME                                                     VISIBILITY
```

Deprecation warning emitted, behavior preserved. Demonstrates migration path.

### Help text

```
$ gws list --help
Flag convention:
  LOWERCASE (-t, -y, -s, etc.) = filter only (no column displayed)
  UPPERCASE (-T, -Y, -S, etc.) = show column (bare) or show + filter (with value)

Examples:
  gws list -t ai                    # Filter by tag "ai" (no tag column shown)
  gws list -T ai                    # Filter by tag "ai" AND show tag column
  gws list -T                       # Show tag column (no filter)
```

Convention documented clearly in help output.

## Test Results

```
$ go test ./...
ok      github.com/daileyo/gws/cmd/git-workspace       0.868s
ok      github.com/daileyo/gws/internal/classifier      (cached)
ok      github.com/daileyo/gws/internal/config           (cached)
ok      github.com/daileyo/gws/internal/discovery        (cached)
ok      github.com/daileyo/gws/internal/filter           (cached)
ok      github.com/daileyo/gws/internal/git              (cached)
ok      github.com/daileyo/gws/internal/user             (cached)
```

All tests pass with zero regressions.

## Cross-Command Consistency

Tag command flags (`-r`, `-p`) are standard `StringVarP` without `NoOptDefVal` - they already accept space-separated values. No changes needed.

The only remaining `NoOptDefVal` usage outside list's uppercase flags is the deprecated `--add` flag on root (`NoOptDefVal = "."`), which is correct and unchanged.

## Flag Mapping Summary

| Long Flag | Filter (lowercase) | Show+Filter (uppercase) |
|-----------|-------------------|------------------------|
| `--type` | `-y` | `-Y` |
| `--visibility` | `-i` (was `-V`) | `-I` |
| `--tag` | `-t` | `-T` |
| `--path` | `-p` | `-P` |
| `--status` | `-s` | `-S` |
| `--show-user` | `-u` | `-U` |
| `--remote` | `-r` | `-R` (was remote-raw) |
| `--remote-raw` | `-b` (was `-R`) | `-B` |
| `--name` | `-n` | N/A (always shown) |
