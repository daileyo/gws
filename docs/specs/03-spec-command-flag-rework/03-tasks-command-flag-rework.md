# 03-tasks-command-flag-rework

## Relevant Files

- `cmd/gws/main.go` - Root command definition, version command, flag registration, entry point. Will be heavily restructured to host all command flags and dispatch logic.
- `cmd/gws/main_test.go` - Tests for root command and version. Must be updated to test new flag-based interface instead of subcommands.
- `cmd/gws/list.go` - List command, filter flags, table/JSON display. Logic moves into root command dispatch; display functions remain as helpers.
- `cmd/gws/init.go` - Init command and `pluralize` helper. Logic moves into root command dispatch.
- `cmd/gws/tag.go` - Tag command and `findRepositories` helper. Logic moves into root command dispatch; helper function remains.
- `cmd/gws/tag_test.go` - Tests for `findRepositories` and tag management. Helper tests remain; command invocation tests updated.
- `cmd/gws/untag.go` - Untag command. Logic moves into root command dispatch.
- `cmd/gws/refresh.go` - Refresh command. Logic moves into root command dispatch.
- `internal/filter/filter.go` - Filter criteria and matching logic. Updated to support wildcard pattern matching with `*` and `?`.
- `internal/filter/filter_test.go` - Filter tests. Extended with wildcard pattern test cases.

### Notes

- Unit tests are placed alongside the code files they test (e.g., `main.go` and `main_test.go` in the same directory).
- Use `make test` to run all tests and `make ci` for full CI checks (vet, lint, test-race).
- Follow existing Go patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests.
- Follow conventional commits for version history.

## Tasks

### [x] 1.0 Convert subcommands to root-level flags with dispatch logic

Remove all Cobra subcommands (`list`, `init`, `tag`, `untag`, `refresh`, `version`) and replace them with flags on the root command. Implement a dispatch handler in `rootCmd.RunE` that checks which command flag is set, enforces mutual exclusivity (only one command flag active at a time), and routes to the appropriate logic. Replace the `version` subcommand with Cobra's built-in `--version`/`-v`. Remove the tag filter shorthand (`gws personal`). Add `-w` shorthand to `--print-workspace`. Update `rootCmd.Use` and help text to reflect the new flag-based interface.

#### 1.0 Proof Artifact(s)

- Test: `main_test.go` updated tests pass — demonstrates root command accepts new flags and old subcommands are removed
- Test: `tag_test.go` updated tests pass — demonstrates `findRepositories` and tag logic still works with new dispatch
- CLI: `gws -h` output shows flag-based interface with all command flags and shorthands

#### 1.0 Tasks

- [x] 1.1 Register command flags on the root command in `main.go`: `--list`/`-l` (bool), `--init`/`-i` (string), `--add-tag`/`-a` (bool), `--remove-tag`/`-u` (bool), `--refresh`/`-r` (bool), `--print-workspace`/`-w` (bool). Set Cobra's built-in version via `rootCmd.Version` using the `version` variable.
- [x] 1.2 Register all list filter flags on the root command in `main.go`: `--type` (string), `--tag` (string slice), `--name` (string), `--path` (string), `--output`/`-o` (string, default "table"), `--status`/`-s` (bool). Remove the filter flag registrations from `list.go`'s `init()` function.
- [x] 1.3 Implement a mutual exclusivity check at the top of `rootCmd.RunE` that counts how many command flags (`--list`, `--init`, `--add-tag`, `--remove-tag`, `--refresh`, `--print-workspace`) are set and returns an error if more than one is active.
- [x] 1.4 Implement dispatch logic in `rootCmd.RunE`: check each command flag and call the corresponding handler function. Extract the existing `RunE` logic from each subcommand into standalone functions (e.g., `runList`, `runInit`, `runTag`, `runUntag`, `runRefresh`) that accept `cmd *cobra.Command` and `args []string`.
- [x] 1.5 Remove all `rootCmd.AddCommand()` calls from the `init()` functions in `list.go`, `init.go`, `tag.go`, `untag.go`, `refresh.go`, and `main.go` (for versionCmd). Remove `versionCmd`, `listCmd`, `initCmd`, `tagCmd`, `untagCmd`, and `refreshCmd` variable declarations.
- [x] 1.6 Remove the tag filter shorthand from `rootCmd.RunE` (the `if len(args) > 0` block that treats bare args as tag filters). Update `rootCmd.Use` from `"gws [tag]"` to `"gws"` and update the `Long` description to reflect the new flag-based interface.
- [x] 1.7 Update `main_test.go`: replace `TestVersionCommandExists` (which uses `rootCmd.Find`) with a test verifying `rootCmd.Version` is set. Update `TestRootCommand` to check the new `Use` string. Remove or update any tests that reference subcommands.
- [x] 1.8 Verify `tag_test.go` still passes — the `findRepositories` helper and tag management tests should be unaffected since they test business logic, not CLI structure.

### [x] 2.0 Add shorthands to list filter flags

Add shorthand flags to all list filter options: `--type`/`-y`, `--name`/`-n`, `--path`/`-p`, `--tag`/`-t`. Keep existing `--output`/`-o` and `--status`/`-s` shorthands. Since all flags are now on the root command, add validation that filter flags are only meaningful when `--list`/`-l` is active (warn or error if filter flags used without `--list`).

#### 2.0 Proof Artifact(s)

- Test: Filter flags with shorthands resolve to correct values demonstrates registration is correct
- Test: Filter flags without `--list` produce an error demonstrates validation works

#### 2.0 Tasks

- [x] 2.1 Update the filter flag registrations in `main.go` to use `StringVarP`/`BoolVarP` variants with shorthands: `--type`/`-y`, `--name`/`-n`, `--path`/`-p`. Update `--tag` from `StringSliceVar` to `StringSliceVarP` with shorthand `-t`.
- [x] 2.2 Add validation in `rootCmd.RunE` that checks if any filter flag (`--type`, `--tag`, `--name`, `--path`, `--output`, `--status`) is set when `--list` is not active. If so, return an error message like `"filter flags (--type, --tag, etc.) require --list/-l to be set"`.
- [x] 2.3 Add tests in `main_test.go` verifying: (a) filter flags with shorthands resolve correctly (e.g., `-y github` sets `filterType` to `"github"`), and (b) setting a filter flag without `--list` produces the expected error.

### [x] 3.0 Add wildcard pattern matching to filters

Implement wildcard pattern matching (`*` for zero or more characters, `?` for single character) in `internal/filter/filter.go`. When a filter value contains `*` or `?`, use glob-style matching (case-insensitive). When no wildcard characters are present, preserve existing behavior (partial, case-insensitive match). Apply wildcard support to all four filter types: `--type`, `--name`, `--path`, and `--tag`.

#### 3.0 Proof Artifact(s)

- Test: `filter_test.go` wildcard tests pass — demonstrates `*` and `?` patterns match correctly across all filter types
- Test: `filter_test.go` existing non-wildcard tests still pass — demonstrates backward compatibility

#### 3.0 Tasks

- [x] 3.1 Create a helper function `matchesPattern(value, pattern string) bool` in `internal/filter/filter.go` that: (a) checks if the pattern contains `*` or `?`; (b) if yes, lowercases both value and pattern and uses `filepath.Match` for glob matching; (c) if no, falls back to the existing partial case-insensitive `strings.Contains` match.
- [x] 3.2 Update `matchesCriteria` in `internal/filter/filter.go` to use `matchesPattern` for the Type, Name, and Path filters instead of the current inline matching logic. For the Type filter, replace `strings.EqualFold` with `matchesPattern` (note: non-wildcard behavior changes from exact-case-insensitive to partial-case-insensitive — this is acceptable per spec).
- [x] 3.3 Update the Tags filter loop in `matchesCriteria` to use `matchesPattern` when comparing each `filterTag` against each `repoTag`, so wildcard patterns like `wo*` match tags like `work` and `workflows`.
- [x] 3.4 Add wildcard-specific test cases to `filter_test.go`: (a) `TestByName` with patterns like `my*`, `*api`, `*project*`, `m?-project`; (b) `TestByType` with patterns like `git*`, `?ithub`; (c) `TestByPath` with patterns like `*/work/*`, `/home/*/projects/*`; (d) `TestByTags` with patterns like `wo*`, `per?onal`.
- [x] 3.5 Verify all existing `filter_test.go` tests still pass unchanged to confirm backward compatibility.

### [ ] 4.0 End-to-end validation and CI pass

Run `make ci` (vet, lint, test-race) and verify all tests pass with no regressions. Fix any linting issues or test failures introduced by the rework. Verify the complete flag mapping matches the spec.

#### 4.0 Proof Artifact(s)

- CLI: `make ci` passes with zero errors demonstrates all quality gates pass
- CLI: `make test` output shows all tests passing demonstrates no regressions

#### 4.0 Tasks

- [ ] 4.1 Run `make fmt` to ensure all modified files are properly formatted.
- [ ] 4.2 Run `make vet` and fix any issues reported by `go vet`.
- [ ] 4.3 Run `make lint` and fix any `golangci-lint` issues (unused variables, missing error checks, etc.).
- [ ] 4.4 Run `make test-race` and fix any race conditions or test failures.
- [ ] 4.5 Run `make ci` to confirm all quality gates pass in sequence.
- [ ] 4.6 Manually verify the complete flag mapping by running `gws -h` and confirming it shows: `--list`/`-l`, `--init`/`-i`, `--add-tag`/`-a`, `--remove-tag`/`-u`, `--refresh`/`-r`, `--print-workspace`/`-w`, `--version`/`-v`, plus filter flags with shorthands.
