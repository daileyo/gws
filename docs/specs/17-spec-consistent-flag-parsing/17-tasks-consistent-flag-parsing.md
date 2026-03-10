# 17-tasks-consistent-flag-parsing

## Relevant Files

- `cmd/git-workspace/list.go` - Flag registration, `parseDualPurposeFlags`, `ListOptions` struct, `showColumnSentinel` constant, custom help function. Primary file for this spec.
- `cmd/git-workspace/list_test.go` - Existing tests for list command flag parsing, stacking, and output. Must be updated extensively.
- `cmd/git-workspace/deprecated.go` - Deprecated root-level flags and deprecation warning system. Needs updates for `-V` and old `-R` deprecation.
- `cmd/git-workspace/tag.go` - Tag command flag registration (`-r`, `-p`). Verify space-separated values work (should already work since no `NoOptDefVal`).
- `cmd/git-workspace/tag_test.go` - Tag command tests. Verify or add space-separated value tests.
- `cmd/git-workspace/main.go` - Root command definition. May need minor updates if root-level flag references change.

### Notes

- Unit tests are in the same directory as the code files (`cmd/git-workspace/`).
- Run tests with `go test ./cmd/git-workspace/` or target specific tests with `-run TestName`.
- Follow existing table-driven test patterns with `t.Run()` subtests.
- Preserve global flag state in tests using defer closures (save/restore pattern already used).
- The `showColumnSentinel` constant (`"\x00show"`) will still be used by uppercase flags' `NoOptDefVal`.
- Cobra flag completion functions (`RegisterFlagCompletionFunc`) may need updates for new flag names.

## Tasks

### [x] 1.0 Implement Lowercase/Uppercase Flag Pair Infrastructure

Replace the current single dual-purpose flag pattern (using `NoOptDefVal` sentinel) with explicit lowercase/uppercase flag pairs. Lowercase flags are standard `StringVarP` (filter only, no `NoOptDefVal`), uppercase flags use `NoOptDefVal` (show column bare, or show+filter with value). Rewrite `parseDualPurposeFlags` to merge both flag sources into `ListOptions`.

#### 1.0 Proof Artifact(s)

- CLI: `gws list -t ai` filters by tag without showing tag column demonstrates lowercase=filter-only with space-separated value
- CLI: `gws list -T` shows tag column without filtering demonstrates uppercase bare=show column
- CLI: `gws list -T ai` shows tag column and filters demonstrates uppercase=show+filter with space-separated value
- CLI: `gws list -t=ai` continues to work demonstrates backward compatibility with `=` syntax
- CLI: `gws list -YT` shows both type and tag columns demonstrates POSIX stacking with bare uppercase flags
- CLI: `gws list -YT ai` shows type column + filters/shows tags demonstrates stacking with trailing value

#### 1.0 Tasks

- [x] 1.1 Add new uppercase flag variables alongside existing lowercase ones in `list.go`.
- [x] 1.2 Refactor flag registration in `init()` to register two flags per column.
- [x] 1.3 Rewrite `parseDualPurposeFlags` to merge lowercase and uppercase flag state.
- [x] 1.4 Handle the tag flag's repeatable nature.
- [x] 1.5 Build and smoke test.

### [x] 2.0 Remap Visibility and Remote-Raw Short Flags

Remap `--visibility` from `-V` to `-i` (filter) / `-I` (show+filter). Remap `--remote-raw` from `-R` to `-b` (filter) / `-B` (show+filter). Free `-R` to become the uppercase show+filter variant of `--remote` (currently `-r`). Add deprecation warnings for old `-V` usage. Update the deprecated flag layer in `deprecated.go`.

#### 2.0 Proof Artifact(s)

- CLI: `gws list -i private` filters by visibility without showing column demonstrates new `-i` flag
- CLI: `gws list -I` shows visibility column demonstrates new `-I` flag
- CLI: `gws list -b github` filters by remote-raw without showing column demonstrates new `-b` flag
- CLI: `gws list -B` shows remote-raw column demonstrates new `-B` flag
- CLI: `gws list -R` shows remote column (not remote-raw) demonstrates `-R` remapped to remote show+filter
- CLI: Old `-V` usage emits deprecation warning demonstrates migration path

#### 2.0 Tasks

- [x] 2.1 Change `--visibility` flag registration to `-i` / `-I`.
- [x] 2.2 Change `--remote-raw` flag registration to `-b` / `-B`.
- [x] 2.3 Change `--remote` uppercase flag to use `-R`.
- [x] 2.4 Update `parseDualPurposeFlags` for the remapped flags.
- [x] 2.5 Add deprecation handling for old `-V` flag.
- [x] 2.6 Build and smoke test remapped flags.

### [x] 3.0 Update Help Text and Command Examples

Update the `list` command's `Long` description, flag descriptions, and custom help function to reflect the new lowercase/uppercase convention. Remove references to `=` requirement. Add clear documentation of the pattern in the help output.

#### 3.0 Proof Artifact(s)

- CLI: `gws list --help` shows updated examples with space-separated values and lowercase/uppercase convention demonstrates help text clarity
- CLI: Flag descriptions clearly indicate "filter only" vs "show column / filter+show" demonstrates convention is documented

#### 3.0 Tasks

- [x] 3.1 Update the `Long` description on `listCmd`.
- [x] 3.2 Update flag description strings.
- [x] 3.3 Update the custom help function.
- [x] 3.4 No flag completion functions needed updating (none registered on list command).

### [x] 4.0 Comprehensive Test Coverage

Update existing list command tests and add new tests covering: space-separated values for all flags, lowercase filter-only behavior, uppercase show+filter behavior, bare uppercase show-column behavior, flag stacking, `=` backward compatibility, deprecated flag warnings, and edge cases.

#### 4.0 Proof Artifact(s)

- Test: `go test ./cmd/git-workspace/` all pass demonstrates comprehensive flag parsing coverage
- Test: Coverage includes space-separated values, bare uppercase, stacking, and deprecated flag warnings demonstrates all new behaviors are tested

#### 4.0 Tasks

- [x] 4.1 Update existing flag registration tests with new shorthands.
- [x] 4.2 Add tests for lowercase filter-only behavior (TestLowercaseFilterSpaceSeparated).
- [x] 4.3 Add tests for uppercase show+filter behavior (TestUppercaseBareShowColumn, TestUppercaseWithValueEqualsOnly, TestUppercaseSpaceSeparatedViaReassignment).
- [x] 4.4 Add tests for `=` backward compatibility (TestLowercaseFilterEqualsCompatibility, TestUppercaseEqualsCompatibility).
- [x] 4.5 Update flag stacking tests (TestUppercaseFlagStacking, TestUppercaseFlagStackingMultiple, TestUppercaseFlagStackingWithTrailingValue, TestUppercaseStackingWithRemote).
- [x] 4.6 Add tests for remapped flags (TestDeprecatedVisibilityFlag).
- [x] 4.7 Add tests for NoOptDefVal presence/absence (TestLowercaseFlagsNoOptDefVal, TestUppercaseFlagsHaveNoOptDefVal).
- [x] 4.8 Full test suite passes (`go test ./...`).

### [x] 5.0 Verify Cross-Command Flag Consistency

Verify that `tag --repo` / `-r` and `tag --path` / `-p` already accept space-separated values (they should, since they're standard `StringVarP` without `NoOptDefVal`). Fix any inconsistencies found. Confirm no other commands have `=`-required flags.

#### 5.0 Proof Artifact(s)

- CLI: Tag command flags already use standard `StringVarP` - space-separated values work
- Test: `go test ./...` passes demonstrates tag command flag behavior is correct

#### 5.0 Tasks

- [x] 5.1 Verified tag command flags accept space-separated values (standard StringVarP, no NoOptDefVal).
- [x] 5.2 Audited all commands for NoOptDefVal - only list uppercase flags and deprecated --add use it.
- [x] 5.3 Tag command tests already exist and pass.
- [x] 5.4 Full test suite passes with zero regressions.
