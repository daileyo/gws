# 12 Tasks - List Remote

## Relevant Files

- `cmd/git-workspace/list.go` - Main list command: flag registration, `ListOptions`, `runList`, `displayTable`, `displayJSON`
- `cmd/git-workspace/list_test.go` - Existing list command tests (flag registration, stacking)
- `cmd/git-workspace/deprecated.go` - Deprecated root flags that dispatch to `runList` — must pass `ShowRemote` through
- `internal/git/remote.go` - New file: helper function to inspect remotes via go-git
- `internal/config/config.go` - `Repository` struct (has `RemoteURL` field, no changes needed)

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `list_test.go` in `cmd/git-workspace/`)
- Use `make test` or `go test ./...` to run the full test suite
- Follow existing Cobra flag registration patterns (`BoolVarP`, flag stacking)
- Follow existing column rendering patterns in `displayTable` (dynamic width, header/separator/row)
- Use `go-git/go-git/v5` for live remote inspection (already a project dependency)

## Tasks

### [x] 1.0 Add `--remote` flag and display origin URL in table output

#### 1.0 Proof Artifact(s)

- CLI: `gws list --remote` output shows REMOTE column with origin URLs as the last column
- CLI: `gws list -r` short flag works identically
- CLI: `gws list -rsu` demonstrates flag stacking with status, user, and remote columns

#### 1.0 Tasks

- [x] 1.1 Add `showRemote` package-level variable and register `--remote`/`-r` boolean flag on `listCmd` in `init()` (follow the `showStatus`/`showUser` pattern)
- [x] 1.2 Add `ShowRemote bool` field to the `ListOptions` struct and pass `showRemote` into it from the `listCmd.RunE` closure
- [x] 1.3 Update `displayTable` signature to accept `showRemote bool` parameter; update all call sites (`runList` and deprecated dispatch in `deprecated.go`)
- [x] 1.4 Add REMOTE column rendering in `displayTable`: calculate `maxRemoteLen` from `repo.RemoteURL` values, add REMOTE header/separator after PATH, and append remote URL to each row (when `showRemote` is true)
- [x] 1.5 Add an example line to the `listCmd.Long` help text showing `--remote` usage (e.g., `gws list -r`)

### [x] 2.0 Add asterisk indicator via live remote inspection

#### 2.0 Proof Artifact(s)

- CLI: Repository with only origin shows URL without asterisk
- CLI: Repository with origin + additional remote shows `* https://...`
- CLI: Repository with no origin but other remotes shows `*`
- CLI: Inaccessible repo path falls back to stored URL without asterisk

#### 2.0 Tasks

- [x] 2.1 Create `internal/git/remote.go` with a `GetRemoteInfo(repoPath string) (originURL string, hasMultiple bool, err error)` function that opens the repo with `git.PlainOpen`, iterates remotes, finds origin URL, and determines if non-origin remotes exist
- [x] 2.2 In `displayTable`, when `showRemote` is true, call `GetRemoteInfo` for each repo to build a `remoteDisplay` string: origin-only → URL, origin+others → `* URL`, no-origin+others → `*`, error/no-remotes → stored `repo.RemoteURL`
- [x] 2.3 Update `maxRemoteLen` calculation to account for the `* ` prefix (2 extra chars) when applicable

### [x] 3.0 JSON output support for remote info

#### 3.0 Proof Artifact(s)

- CLI: `gws list -r -o json` output includes `remote_url` and `has_multiple_remotes` fields per repo
- CLI: JSON for a repo with multiple remotes shows `"has_multiple_remotes": true`
- CLI: JSON for a repo with only origin shows `"has_multiple_remotes": false`

#### 3.0 Tasks

- [x] 3.1 Create a `listJSONEntry` struct (or similar) that embeds/mirrors `config.Repository` fields and adds `RemoteURL string` and `HasMultipleRemotes bool` JSON fields for remote-aware output
- [x] 3.2 Update `displayJSON` (or create a new `displayJSONWithRemote`) to accept a `showRemote bool` parameter; when true, call `GetRemoteInfo` for each repo and populate the extended struct before marshalling
- [x] 3.3 Update the `runList` JSON dispatch to pass `showRemote` through to the JSON display function

### [x] 4.0 Unit tests

#### 4.0 Proof Artifact(s)

- Test: `go test ./cmd/git-workspace/ -run TestRemote` passes
- Test: `go test ./...` passes (no regressions)

#### 4.0 Tasks

- [x] 4.1 Add `{"remote", "r"}` to the flag table in `TestFilterFlagsOnListCmd` to verify the flag is registered with correct shorthand
- [x] 4.2 Add a `TestListCmdFlagStackingWithRemote` test that parses `-rsu` and verifies all three booleans (`showRemote`, `showStatus`, `showUser`) are true
- [x] 4.3 Add `TestGetRemoteInfo` in `internal/git/remote_test.go` covering: origin-only repo, origin+upstream repo, no-origin repo with other remotes, no remotes, and invalid path (use `git.PlainInit` to create temp repos in tests)
