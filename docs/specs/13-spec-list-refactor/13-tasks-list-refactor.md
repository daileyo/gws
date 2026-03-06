# 13 Tasks - List Refactor

## Relevant Files

- `cmd/git-workspace/list.go` - Main list command: flag registration, `runList()`, `displayTable()`, `displayJSON()`. Primary file for this refactor.
- `cmd/git-workspace/list_test.go` - Tests for list command flags and behavior.
- `cmd/git-workspace/main.go` - Root command: default no-args behavior, navigation, usage template.
- `cmd/git-workspace/main_test.go` - Tests for root command behavior.
- `cmd/git-workspace/deprecated.go` - Deprecated flag registration and dispatch. Must remain backward-compatible.
- `cmd/git-workspace/deprecated_test.go` - Tests for deprecated flag behavior.
- `internal/filter/filter.go` - Filter logic (`Criteria`, `Apply`, `MatchesPattern`). May need `Visibility` field added to `Criteria`.
- `internal/filter/filter_test.go` - Filter unit tests.
- `go.mod` / `go.sum` - May need `golang.org/x/term` dependency for terminal width detection.

### Notes

- Unit tests are co-located with source files (e.g., `list.go` and `list_test.go`).
- Use `go test ./...` to run all tests; `make ci` for full CI checks (vet + lint + race).
- Follow existing table-driven test patterns with `t.Run()` subtests.
- Follow `gofmt` formatting and Conventional Commits for any commits.

## Tasks

### [x] 1.0 Implement Dual-Purpose Flag Infrastructure

Rework the `gws list` flag system from the current mix of `Bool`/`String`/`StringSlice` flags to dual-purpose flags that support "no value = show column" and "with value = filter + show column". This is the foundational change that all other tasks depend on.

#### 1.0 Proof Artifact(s)

- CLI: `gws list -y` shows NAME and TYPE columns (display-only, no filter)
- CLI: `gws list -y github` filters by type=github AND shows NAME and TYPE columns
- CLI: `gws list -yVtp` shows NAME, TYPE, VISIBILITY, TAGS, PATH columns (POSIX stacking)
- CLI: `gws list -y github -V private -tp` filters type+visibility, displays all four columns plus name
- CLI: `gws list -n "api-*"` filters by name pattern (filter-only flag, no column toggle)
- Test: All existing `list_test.go` tests pass; new tests cover dual-purpose flag parsing

#### 1.0 Tasks

- [x] 1.1 Design the dual-purpose flag mechanism using Cobra's `NoOptDefVal` feature. Each flag (type, visibility, tag, path, status, show-user, remote) should be a `String` flag with `NoOptDefVal` set to a sentinel value (e.g., `""` or a constant like `showColumnOnly`). When the flag is present with the sentinel, it means "show column"; when present with a real value, it means "filter + show column"; when absent, neither show nor filter. Prototype this approach with `-y`/`--type` first.
- [x] 1.2 Add the new `-V`/`--visibility` flag to `listCmd`. Currently there is no visibility flag — it's always shown. Register it as a `String` flag with shorthand `V` and `NoOptDefVal` for display-only behavior.
- [x] 1.3 Convert existing `Bool` flags (`-s`/`--status`, `-u`/`--show-user`, `-r`/`--remote`) to `String` flags with `NoOptDefVal`. Update the variables from `bool` to `string` and adjust all references in `runList()` and `displayTable()`.
- [x] 1.4 Convert `-t`/`--tag` from `StringSlice` to a `String` flag with `NoOptDefVal`. For multiple tag filters, support repeating the flag (`-t personal -t work`). Handle the `filterTags` shared variable carefully to maintain backward compatibility with `deprecated.go`.
- [x] 1.5 Convert `-p`/`--path` from a pure filter `String` flag to a dual-purpose flag with `NoOptDefVal` (display-only shows PATH column, with value filters + shows).
- [x] 1.6 Update the `ListOptions` struct to track which columns to display separately from filter values. Add fields like `ShowType`, `ShowVisibility`, `ShowTags`, `ShowPath`, `ShowStatus`, `ShowUser`, `ShowRemote` (booleans) alongside the existing filter fields.
- [x] 1.7 Update `runList()` to parse the new dual-purpose flags: for each flag, check if it was `Changed` on the command, then check its value to determine show-only vs filter+show. Populate `ListOptions` accordingly.
- [x] 1.8 Update `displayTable()` to conditionally show columns based on the `Show*` booleans in `ListOptions` instead of always showing TYPE, VISIBILITY, TAGS, PATH. Maintain the fixed column order: NAME, STATUS, USER, EMAIL, SIGN, TYPE, VISIBILITY, TAGS, PATH, REMOTE.
- [x] 1.9 Add `Visibility` field to `filter.Criteria` in `internal/filter/filter.go` and implement `ByVisibility()` filter function with case-insensitive matching (same pattern as `ByType`).
- [x] 1.10 Update existing tests in `list_test.go` to reflect the new flag types (String with NoOptDefVal instead of Bool). Add new tests for: dual-purpose flag parsing (display-only vs filter+display), POSIX flag stacking with new flags (`-yVtp`), visibility filtering, and combined filter+display scenarios.
- [x] 1.11 Add unit tests for visibility filtering in `internal/filter/filter_test.go`.
- [x] 1.12 Run `go test ./...` and `make ci` to verify all tests pass and no lint/vet issues.

### [~] 2.0 Implement Multi-Column Default View

Change the default `gws list` output from the full table (NAME, TYPE, VISIBILITY, TAGS, PATH) to a compact multi-column layout showing only repo names, sorted alphabetically. Detect terminal width for column sizing, with single-column fallback for non-TTY.

#### 2.0 Proof Artifact(s)

- CLI: `gws list` displays repo names in multi-column layout (ls-style, left-to-right)
- CLI: `gws list | cat` outputs one repo name per line (non-TTY fallback)
- CLI: `gws list -n "api-*"` filters names and displays in multi-column layout
- Test: Unit tests cover column calculation, alphabetical sorting, and TTY detection

#### 2.0 Tasks

- [x] 2.1 Add `golang.org/x/term` dependency (`go get golang.org/x/term`) for terminal width detection. This provides `term.GetSize()` to query terminal dimensions.
- [x] 2.2 Implement a `displayMultiColumn()` function in `list.go` that takes a list of repo names, sorts them alphabetically, detects terminal width, calculates column count and widths, and prints names left-to-right then top-to-bottom (like `ls`). Use 2-space padding between columns.
- [x] 2.3 Implement TTY detection: check if stdout is a terminal using `term.IsTerminal()`. If not a TTY (piped/redirected), output one name per line. Default terminal width to 80 if detection fails.
- [x] 2.4 Update `runList()` to call `displayMultiColumn()` when no column display flags are set (no `-y`, `-V`, `-t`, `-p`, `-s`, `-u`, `-r`, `-v`). The table format with headers is only used when at least one column flag or verbose flag is present.
- [x] 2.5 Keep the "Found N repositories" header line before the multi-column output.
- [x] 2.6 Add unit tests for `displayMultiColumn()`: test column count calculation for various terminal widths, alphabetical sorting, single-column fallback, and edge cases (0 repos, 1 repo, more repos than fit in one row).
- [x] 2.7 Run `go test ./...` and `make ci` to verify.

### [ ] 3.0 Implement Verbose Levels (`-v` and `-vv`)

Add `-v`/`--verbose` (count flag) to `gws list`. Single `-v` shows all stored-data columns (NAME, TYPE, VISIBILITY, TAGS, PATH). Double `-vv` shows everything including live-fetched columns (STATUS, USER, EMAIL, SIGN, REMOTE). When combined with other flags/filters, verbose overrides selective column display.

#### 3.0 Proof Artifact(s)

- CLI: `gws list -v` shows NAME, TYPE, VISIBILITY, TAGS, PATH table (equivalent to old default)
- CLI: `gws list -vv` shows all columns including STATUS, USER, EMAIL, SIGN, REMOTE
- CLI: `gws list -v -y github` filters by github, displays all stored-data columns
- CLI: `gws list -vv -y github` filters by github, displays all columns
- Test: Tests cover `-v` vs `-vv` column selection and verbose + filter interaction

#### 3.0 Tasks

- [ ] 3.1 Register a count flag `-v`/`--verbose` on `listCmd` using Cobra's `CountVarP()`. This tracks how many times `-v` appears: 0 = not set, 1 = `-v`, 2+ = `-vv`.
- [ ] 3.2 Add a `VerboseLevel` field (int) to `ListOptions`.
- [ ] 3.3 Update the column display logic in `runList()`: if `VerboseLevel >= 1`, set all stored-data `Show*` flags to true (ShowType, ShowVisibility, ShowTags, ShowPath). If `VerboseLevel >= 2`, also set ShowStatus, ShowUser, ShowRemote to true. Verbose overrides any selective column flags but does NOT override filter flags.
- [ ] 3.4 When `-v` or `-vv` is combined with filter flags (e.g., `-v -y github`), the filter applies but all columns for that verbosity level are displayed.
- [ ] 3.5 Add unit tests: `-v` sets verbose level 1 and shows stored-data columns; `-vv` sets verbose level 2 and shows all columns; `-v -y github` filters by type but shows all stored-data columns.
- [ ] 3.6 Run `go test ./...` and `make ci` to verify.

### [ ] 4.0 Default No-Args Behavior (`gws` → `gws list`)

Change the root command so `gws` with no arguments runs `gws list` (multi-column repo names) instead of the current workspace summary. Navigation via positional args (`gws my-repo`) continues unchanged.

#### 4.0 Proof Artifact(s)

- CLI: `gws` with no arguments displays multi-column repo names (same as `gws list`)
- CLI: `gws my-repo` still navigates to the repo (unchanged)
- CLI: `gws -q my-repo` still works with quiet flag
- Test: Tests verify no-args defaults to list, positional-arg navigation unchanged

#### 4.0 Tasks

- [ ] 4.1 In `main.go`, update the `rootCmd.RunE` function: when `len(args) == 0` and no deprecated flags are active, call `runList()` with default options (no filters, no column flags, verbose level 0) instead of printing the workspace summary. This will trigger the multi-column default view from Task 2.0.
- [ ] 4.2 Remove the `--quiet` validation that errors when `len(args) == 0`. With the new default list behavior, `--quiet` with no args should either be silently ignored or still produce an error — match the spec's intent (quiet is navigation-only).
- [ ] 4.3 Verify the `--quiet` flag still works correctly with positional args: `gws -q my-repo`.
- [ ] 4.4 Update the root command's `Long` description and usage template to reflect that `gws` now defaults to listing repositories.
- [ ] 4.5 Add/update tests in `main_test.go`: `gws` with no args runs list (returns repo names), positional arg still navigates, `--quiet` with positional arg still works.
- [ ] 4.6 Run `go test ./...` and `make ci` to verify.

### [ ] 5.0 JSON Output Column Selection and Help Text Updates

Update JSON output (`-o json`) to respect display flags — only include fields for selected columns (always include `name`). Fix `gws list -h` to show all flags including new `-V`/`--visibility` and `-v`/`--verbose`. Ensure deprecated root-level flags continue working with deprecation warnings.

#### 5.0 Proof Artifact(s)

- CLI: `gws list -h` shows all flags including `-V`, `-v`, and all dual-purpose flags
- CLI: `gws list -yp -o json` outputs JSON with only name, type, and path fields
- CLI: `gws list -y github -o json` filters by github, includes only name and type in JSON
- CLI: `gws list -v -o json` outputs JSON with all stored-data fields
- CLI: `gws --list --type github` still works with deprecation warning
- Test: Tests cover JSON field selection and deprecated flag backward compatibility

#### 5.0 Tasks

- [ ] 5.1 Update `displayJSON()` to accept the `Show*` flags from `ListOptions`. Instead of marshaling the full `config.Repository` struct, build a dynamic `map[string]interface{}` for each repo that only includes fields for displayed columns. Always include `name`.
- [ ] 5.2 Map column flags to JSON fields: `ShowType` → `type`, `ShowVisibility` → `visibility`, `ShowTags` → `tags`, `ShowPath` → `path`, `ShowStatus` → `status` (format as string), `ShowUser` → `user`/`email`/`signing_enabled`, `ShowRemote` → `remote_url`/`has_multiple_remotes`. When no column flags are set (default multi-column view), JSON outputs only `name`.
- [ ] 5.3 When `-v` is used with `-o json`, include all stored-data fields. When `-vv` is used, include all fields including live-fetched ones.
- [ ] 5.4 Verify all flags appear in `gws list -h` output. Ensure none are marked hidden on `listCmd`. Update help text and examples in `listCmd.Long` to reflect the new flag behavior (dual-purpose, verbose levels).
- [ ] 5.5 Verify deprecated flags in `deprecated.go` still work: `gws --list --type github` should emit deprecation warnings and call `runList()` with appropriate options. Update the deprecated dispatch to populate the new `ListOptions` fields correctly (set `ShowType`, `ShowVisibility`, etc. to match old default behavior — show all stored-data columns).
- [ ] 5.6 Update `depWarnings` map if any deprecated flag names changed.
- [ ] 5.7 Add tests for JSON column selection: `-yp -o json` includes only name/type/path; `-v -o json` includes all stored fields; default `-o json` includes only name.
- [ ] 5.8 Add tests verifying deprecated flag backward compatibility with new `ListOptions` structure.
- [ ] 5.9 Run `go test ./...` and `make ci` to verify all tests pass.
