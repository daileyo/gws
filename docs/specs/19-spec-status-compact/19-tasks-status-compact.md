# 19 Tasks - Status Compact Display

## Relevant Files

- `cmd/git-workspace/list.go` - Main list command: flag definitions, `parseDualPurposeFlags()`, `runList()`, `displayTable()`, `formatStatusIcons()`, `displayWidth()`, column width calculations, and row rendering. Most changes happen here.
- `cmd/git-workspace/list_test.go` - Tests for list command including flag parsing, display formatting, and column rendering.
- `internal/git/status.go` - `Status` struct and `formatStatusShort()`. Reference only — no changes needed.
- `docs/site/commands-core.md` - CLI documentation for the `list` command. Updated `-s` flag description.

### Notes

- Unit tests use table-driven patterns with `TestFunctionName_Scenario` naming.
- Run `go vet ./... && go test -race ./...` to validate all changes.
- Follow Conventional Commits for commit messages.
- Status filtering happens post-fetch in `runList()` using `filter.MatchesPattern()` directly (no changes to `filter.Criteria` needed since status data isn't available at initial filter time).

## Tasks

### [x] 1.0 Enable `-s` to trigger status fetching and add compact display mode

Change `-s` from a filter-only flag to a flag that also triggers status data fetching and displays compact status icons (ahead/behind/clean/dirty) right-justified in the NAME column. This is the core feature: bare `gws list -s` shows repo names with status icons inline, no separate STATUS column.

#### 1.0 Proof Artifact(s)

- CLI: `gws list -s` output shows repo names with right-justified status icons in the NAME column, no STATUS column header present
- CLI: `gws list -s` with color enabled shows correctly colored icons (magenta ↓, cyan ↑, green ✓, red ✗)
- CLI: `gws list -s` shows `✓` for clean repos with no ahead/behind
- CLI: `gws list -S` output is unchanged from current behavior (regression check)

#### 1.0 Tasks

- [x] 1.1 Add `NoOptDefVal` to the `-s`/`--status` flag so bare `-s` is distinguished from `-s <value>`. Added `'s': &flagStatus` to `shortToShowFlag` for trailing arg reassignment.
- [x] 1.2 Add a `CompactStatus` bool field to the `ListOptions` struct.
- [x] 1.3 Update `parseDualPurposeFlags()` to set `opts.CompactStatus = true` when `--status` flag is changed and `ShowStatus` is false. Clean sentinel from `FilterStatus`.
- [x] 1.4 Update `AnyColumnSelected()` to include `CompactStatus`.
- [x] 1.5 Update status fetching in `runList()` to trigger on `opts.ShowStatus || opts.CompactStatus`.
- [x] 1.6 Add compact icon width calculation in `displayTable()` and expand `maxNameLen`.
- [x] 1.7 Add compact icon rendering in row loop with right-justified icons in NAME field.
- [x] 1.8 Changed STATUS column header/row conditions from `statusMap != nil` to `opts.ShowStatus && statusMap != nil`.

### [x] 2.0 Implement status filtering with compact display

Wire up the `-s <pattern>` filter so that `gws list -s dirty` both filters repos by status pattern AND shows compact icons in the NAME column.

#### 2.0 Proof Artifact(s)

- CLI: `gws list -s dirty` shows only dirty repos with compact status icons
- CLI: `gws list -s clean` shows only clean repos with `✓` in the NAME column
- CLI: `gws list -s ahead` shows only repos with unpushed commits

#### 2.0 Tasks

- [x] 2.1 Skipped adding `Status` field to `filter.Criteria` — not needed since status filtering is post-fetch.
- [x] 2.2 Added `statusFilterText()` helper that generates searchable text ("clean", "dirty", "ahead", "behind") from a `git.Status`.
- [x] 2.3 Added post-fetch status filter in `runList()` using `filter.MatchesPattern()` — runs after status data is fetched, before display.
- [x] 2.4 Verified via unit tests (`TestStatusFilterText`).

### [x] 3.0 Handle `-s`/`-S` combination and edge cases

When both `-s` and `-S` are specified, `-S` wins (full STATUS column shown, no compact icons in NAME). JSON output ignores compact status. `-vv` continues to show full status column.

#### 3.0 Proof Artifact(s)

- CLI: `gws list -s -S` shows full STATUS column without compact icons in NAME column
- CLI: `gws list -s dirty -S` shows full STATUS column filtered to dirty repos only
- CLI: `gws list -s -o json` does NOT include compact status fields in JSON output
- CLI: `gws list -vv` continues to show full STATUS column (no compact icons)

#### 3.0 Tasks

- [x] 3.1 `CompactStatus` is only set when `-s` is changed AND `ShowStatus` is false — `-S` wins by design.
- [x] 3.2 `-vv` sets `ShowStatus=true` but never sets `CompactStatus` — verified by code structure.
- [x] 3.3 `displayJSON` checks `opts.ShowStatus` not `CompactStatus` — JSON unaffected. Verified via `TestDisplayJSON_CompactStatus_NoStatusField`.
- [x] 3.4 All combinations verified via unit tests.

### [x] 4.0 Update help text, documentation, and tests

Updated `-s` flag description, command examples in help text, and added tests.

#### 4.0 Proof Artifact(s)

- CLI: `gws list --help` shows updated description for `-s` flag reflecting compact display behavior
- Test: All new and existing tests pass via `go vet ./... && go test -race ./...`

#### 4.0 Tasks

- [x] 4.1 Updated `-s` flag description to "Show compact status in name column, or show and filter by status pattern".
- [x] 4.2 Updated `Long` description with `-s` and `-s dirty` examples.
- [x] 4.3 Updated `docs/site/commands-core.md` with compact status section and filter examples.
- [x] 4.4 Added tests: `TestDisplayTable_CompactStatus`, `TestDisplayTable_CompactStatus_NoStatusColumn`, `TestDisplayTable_CompactStatusAlignment`.
- [x] 4.5 Added tests: `TestStatusFilterText` covers clean/dirty/ahead/behind filtering.
- [x] 4.6 Added tests: `TestDisplayTable_ShowStatusWithCompactStatus_ShowStatusWins`, `TestDisplayJSON_CompactStatus_NoStatusField`.
- [x] 4.7 All tests pass: `go vet ./...` clean, `go test -race ./...` all pass.
