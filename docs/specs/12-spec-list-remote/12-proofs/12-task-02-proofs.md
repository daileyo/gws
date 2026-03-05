# Task 2.0 Proof Artifacts - Asterisk indicator via live remote inspection

## Implementation

### `internal/git/remote.go` - GetRemoteInfo

Created `GetRemoteInfo(repoPath string) (*RemoteInfo, error)` that:
- Opens repo with `git.PlainOpen`
- Iterates all remotes to find origin URL and count non-origin remotes
- Returns `RemoteInfo{OriginURL, HasMultiple}`

### `displayTable` - Asterisk logic

Pre-computes `remoteDisplayMap` for all repos before rendering:
- **Origin only** → `https://github.com/user/repo.git`
- **Origin + others** → `* https://github.com/user/repo.git`
- **No origin + others** → `*`
- **Error/inaccessible** → falls back to stored `repo.RemoteURL` (no asterisk)
- **No remotes at all** → falls back to stored `repo.RemoteURL`

Width calculation uses the pre-computed display strings (accounting for `* ` prefix).

## Test Results

```
$ go test ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	0.463s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	0.046s
ok  	github.com/daileyo/gws/internal/user	0.158s
```

## Verification

- [x] `GetRemoteInfo` created in `internal/git/remote.go`
- [x] Origin-only repos show URL without asterisk
- [x] Origin + additional remotes prefix with `* `
- [x] No-origin + other remotes show standalone `*`
- [x] Inaccessible repo paths fall back to stored URL without asterisk
- [x] Width calculation accounts for `* ` prefix
- [x] All existing tests pass
