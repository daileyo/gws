# Task 3.0 Proof Artifacts - JSON output support for remote info

## Implementation

### `repoJSONWithRemote` struct

Embeds `config.Repository` and adds `HasMultipleRemotes bool` with JSON tag `has_multiple_remotes`.

### `displayJSON` updated

- Accepts `showRemote bool` parameter
- When false: marshals repos directly (existing behavior)
- When true: calls `GetRemoteInfo` for each repo, populates extended struct with `remote_url` and `has_multiple_remotes`

### JSON output format (verified via embedding test)

```json
{
  "name": "test",
  "remote_url": "https://github.com/x/y.git",
  "has_multiple_remotes": true
}
```

## Test Results

```
$ go test ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	0.413s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Verification

- [x] `repoJSONWithRemote` struct created with embedded `config.Repository`
- [x] `displayJSON` accepts `showRemote` parameter
- [x] `GetRemoteInfo` called per repo when `showRemote` is true
- [x] JSON output includes `remote_url` and `has_multiple_remotes` fields
- [x] No duplicate JSON keys from embedding
- [x] `runList` passes `ShowRemote` through to JSON path
- [x] All existing tests pass
