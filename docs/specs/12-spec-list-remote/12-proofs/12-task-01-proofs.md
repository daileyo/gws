# Task 1.0 Proof Artifacts - Add `--remote` flag and table output

## CLI Output

### `gws list --remote` (no workspace initialized)

```
$ gws list --remote
No repositories found. Run 'gws init' to discover repositories.
```

### `gws list -r` short flag

```
$ gws list -r
No repositories found. Run 'gws init' to discover repositories.
```

### `gws list -rsu` flag stacking

```
$ gws list -rsu
No repositories found. Run 'gws init' to discover repositories.
```

All three flag forms parse without error, confirming flag registration and stacking work correctly.

### `gws list --help` shows new example

```
Examples:
  gws list                          # List all repositories
  gws list --type github            # List only GitHub repositories
  gws list -t personal -s           # List repos tagged "personal" with git status
  gws list -n "api-*" -o json       # Filter by name with wildcard, output as JSON
  gws list -su                      # List with status and user columns
  gws list -r                       # List with remote URL column
```

## Test Results

```
$ go test ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	0.483s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Verification

- [x] `--remote`/`-r` flag registered on listCmd
- [x] Flag stacking `-rsu` works without error
- [x] `ShowRemote` field added to `ListOptions`
- [x] `displayTable` accepts `showRemote` parameter
- [x] REMOTE column rendered as last column (after PATH) when flag is set
- [x] Deprecated dispatch updated to pass `ShowRemote: false`
- [x] Help text includes `-r` example
- [x] All existing tests pass
