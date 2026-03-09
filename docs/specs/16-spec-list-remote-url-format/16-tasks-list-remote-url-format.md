# 16 Tasks - List Remote URL Format

## Relevant Files

- `internal/git/remote.go` - Contains `RemoteInfo` struct and `GetRemoteInfo()`. Add `FormatRemoteURL()` here.
- `internal/git/remote_test.go` - Existing remote tests. Add `TestFormatRemoteURL` table-driven tests here.
- `cmd/git-workspace/list.go` - List command: flag registration (`init()`), `ListOptions` struct, `parseDualPurposeFlags()`, table rendering (`remoteDisplayMap` block ~L450-471), JSON rendering (~L932-945), help text (~L99-123).
- `cmd/git-workspace/list_test.go` - List command tests: flag registration tests, flag stacking tests. Add `--remote-raw`/`-R` tests here.

### Notes

- Unit tests are placed alongside the code files they test (e.g., `remote.go` and `remote_test.go` in the same package).
- Run tests with `go test ./...` or target specific packages with `go test ./internal/git/` or `go test ./cmd/git-workspace/`.
- Follow the existing dual-purpose flag pattern using `NoOptDefVal = showColumnSentinel`.
- Follow conventional commits: `feat:` for new functionality, `test:` for test-only changes.

## Tasks

### [x] 1.0 URL Formatting Utility

Add a `FormatRemoteURL` function to `internal/git/remote.go` that converts raw git remote URLs into clean HTTPS URLs. This is the pure-logic foundation — no CLI or display changes yet.

Conversions:
- SSH (`git@host:path`) → `https://host/path`
- HTTPS with user info (`https://user@host/path`) → `https://host/path`
- HTTPS with user:pass (`https://user:pass@host/path`) → `https://host/path`
- Azure DevOps SSH (`git@ssh.dev.azure.com:v3/org/project/repo`) → `https://dev.azure.com/org/project/_git/repo`
- Unrecognized formats → return original URL unchanged

#### 1.0 Proof Artifact(s)

- Test: `go test ./internal/git/ -run TestFormatRemoteURL` passes demonstrating correct URL conversion for SSH, HTTPS, HTTPS-with-user-info, Azure DevOps, and unformattable URLs

#### 1.0 Tasks

- [x] 1.1 Add `FormatRemoteURL(rawURL string) string` function to `internal/git/remote.go`. Use `net/url.Parse` for HTTPS URLs (strip `Userinfo`, return clean URL). For SSH format (`git@host:path`), extract host and path, construct `https://host/path`. For Azure DevOps SSH (`git@ssh.dev.azure.com:v3/org/project/repo`), convert to `https://dev.azure.com/org/project/_git/repo`. Return the original URL unchanged for `file://`, empty strings, or any URL that fails parsing.
- [x] 1.2 Add `TestFormatRemoteURL` table-driven tests to `internal/git/remote_test.go` covering: (a) standard SSH `git@github.com:owner/repo.git`, (b) HTTPS without user info (passthrough), (c) HTTPS with `user@`, (d) HTTPS with `user:pass@`, (e) Azure DevOps SSH, (f) Azure DevOps HTTPS with user info, (g) `file:///local/path` (unchanged), (h) empty string (unchanged), (i) custom SSH alias like `myhost:repo.git` (unchanged or best-effort), (j) GitLab SSH `git@gitlab.com:group/subgroup/repo.git`.
- [x] 1.3 Run `go test ./internal/git/ -run TestFormatRemoteURL` and verify all tests pass.

### [x] 2.0 `--remote-raw`/`-R` Flag and Display Integration

Register the `--remote-raw`/`-R` flag on the `list` command following the existing dual-purpose pattern. Update `ListOptions`, `parseDualPurposeFlags`, and the table/JSON rendering to:
- Default `-r` to formatted URLs (via `FormatRemoteURL`)
- `-R` shows raw remote URLs
- `-R` without `-r` implicitly enables remote display
- `-rR` together shows raw (override)
- Support flag stacking (`-Rsu`)

#### 2.0 Proof Artifact(s)

- CLI: `gws list -r` output shows formatted HTTPS URLs demonstrates default formatting
- CLI: `gws list -R` output shows raw remote URLs demonstrates raw flag
- CLI: `gws list -rR` output shows raw URLs demonstrates override behavior
- CLI: `gws list -r -o json` output shows formatted `remote_url` demonstrates JSON formatting
- CLI: `gws list -R -o json` output shows raw `remote_url` demonstrates JSON raw mode
- Test: Flag registration and stacking tests pass for `--remote-raw`/`-R`

#### 2.0 Tasks

- [x] 2.1 Add `flagRemoteRaw string` variable to the flag variables block in `list.go` (~L27-40).
- [x] 2.2 Register the `--remote-raw` flag in the `init()` function following the existing pattern: `listCmd.Flags().StringVarP(&flagRemoteRaw, "remote-raw", "R", "", "Show raw remote URL column, or filter by raw remote URL pattern")` with `listCmd.Flags().Lookup("remote-raw").NoOptDefVal = showColumnSentinel`.
- [x] 2.3 Add `ShowRemoteRaw bool` and `FilterRemoteRaw string` fields to the `ListOptions` struct (~L196-225).
- [x] 2.4 Update `parseDualPurposeFlags()` to parse the new flag: `opts.FilterRemoteRaw, opts.ShowRemoteRaw = parseDual("remote-raw", flagRemoteRaw)`. Add logic: if `ShowRemoteRaw` is true, implicitly set `ShowRemote = true` (so `-R` alone enables the REMOTE column).
- [x] 2.5 Update `AnyColumnSelected()` to include `o.ShowRemoteRaw` in the return expression.
- [x] 2.6 Update the `remoteDisplayMap` pre-computation block (~L450-471) to apply `git.FormatRemoteURL()` to the URL when `ShowRemoteRaw` is false. When `ShowRemoteRaw` is true, use the raw URL (current behavior). The asterisk prefix logic remains unchanged.
- [x] 2.7 Update the JSON rendering block (~L932-945) to apply `git.FormatRemoteURL()` to `remote_url` values when `ShowRemoteRaw` is false, and use raw URLs when `ShowRemoteRaw` is true.
- [x] 2.8 Add flag registration test for `--remote-raw`/`-R` to `TestFilterFlagsOnListCmd` in `list_test.go`.
- [x] 2.9 Add a flag stacking test `TestListCmdFlagStackingWithRemoteRaw` in `list_test.go` to verify `-Rsu` sets `flagRemoteRaw`, `flagStatus`, and `flagUser` to sentinel values.
- [x] 2.10 Run `go test ./cmd/git-workspace/` and verify all tests pass.

### [x] 3.0 Help Text Updates

Update the `list` command's long description, flag descriptions, and examples to document both `-r` (formatted) and `-R`/`--remote-raw` (raw) behavior.

#### 3.0 Proof Artifact(s)

- CLI: `gws list --help` output shows updated descriptions for `-r` and `-R` demonstrates discoverability

#### 3.0 Tasks

- [x] 3.1 Update the `-r`/`--remote` flag description from `"Show remote URL column, or filter by remote URL pattern"` to `"Show formatted remote URL column, or filter by remote URL pattern"` (~L155).
- [x] 3.2 Update the `listCmd.Long` examples section (~L109-121) to include examples for `-R` usage, e.g., `gws list -R                       # Show raw remote URL column` and `gws list -rR                      # Show raw remote URL (override)`.
- [x] 3.3 Run `go build ./cmd/git-workspace && ./build/git-workspace list --help` and verify the updated descriptions and examples appear correctly.

### [x] 4.0 End-to-End Verification and Cleanup

Run full test suite, verify no regressions, and confirm all proof artifacts from Units 1-3 are satisfied.

#### 4.0 Proof Artifact(s)

- Test: `go test ./...` passes with no failures demonstrates no regressions
- CLI: `gws list -r` on real workspace shows formatted URLs demonstrates end-to-end behavior

#### 4.0 Tasks

- [x] 4.1 Run `go vet ./...` and verify no issues.
- [x] 4.2 Run `go test ./...` and verify all tests pass with no failures.
- [x] 4.3 Build and run `gws list -r` on the real workspace, verify formatted HTTPS URLs appear in the REMOTE column.
- [x] 4.4 Run `gws list -R` on the real workspace, verify raw remote URLs appear.
- [x] 4.5 Run `gws list -rR` and verify raw URLs are shown (override behavior).
- [x] 4.6 Run `gws list -r -o json` and `gws list -R -o json`, verify JSON output uses formatted and raw URLs respectively.
