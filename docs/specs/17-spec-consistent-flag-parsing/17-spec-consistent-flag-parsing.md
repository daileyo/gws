# 17-spec-consistent-flag-parsing

## Introduction/Overview

The `list` command's flag behavior is inconsistent: filter-only flags like `-n` accept values with a space (`-n dks`), while dual-purpose flags like `-t` require an equals sign (`-t=ai`). This inconsistency confuses users and violates the principle of least surprise. This spec standardizes all flags to accept values with a space (no `=` required), while also introducing a new lowercase/uppercase convention that separates filtering from column display control.

## Goals

- Eliminate the `=` requirement for all flag values — space-separated values work everywhere
- Introduce a consistent lowercase/uppercase flag convention: lowercase filters only, uppercase filters and shows the column
- Maintain backward compatibility — `=` syntax continues to work
- Remap `-V` (visibility) to `-i`/`-I` and `-R` (remote-raw) to `-b`/`-B` to support the new pattern
- Keep `-v` verbose count flag unchanged
- Apply consistent behavior across all commands with filter flags

## User Stories

- **As a CLI user**, I want to type `-t ai` instead of `-t=ai` so that filtering feels natural and consistent with other flags like `-n`.
- **As a power user**, I want to filter without forcing a column to display (e.g., `-t ai` filters by tag without showing the tags column) so that I can keep my output focused.
- **As a power user**, I want to show a column without filtering by using the uppercase flag bare (e.g., `-T`) so that I have explicit control over column visibility.
- **As an existing user**, I want my current `=` syntax to continue working so that my scripts and muscle memory are not broken.

## Demoable Units of Work

### Unit 1: Remove `=` Requirement from Dual-Purpose Flags

**Purpose:** Make all flags accept space-separated values consistently, eliminating the most visible inconsistency.

**Functional Requirements:**
- The system shall accept space-separated values for all dual-purpose flags (e.g., `-t ai`, `-y github`, `-s clean`)
- The system shall continue to accept `=` syntax for all flags (e.g., `-t=ai`, `-y=github`)
- The system shall treat `-t ai` and `-t=ai` as equivalent operations
- The system shall correctly distinguish a bare flag at end-of-args (e.g., `gws list -t`) from a flag with a value (e.g., `gws list -t ai`)
- The system shall support POSIX flag stacking where only the last flag in a stack takes a value (e.g., `-yt ai` means show type column + filter tags by "ai")

**Proof Artifacts:**
- CLI: `gws list -t ai` returns filtered results (same as current `gws list -t=ai`) demonstrates space-separated values work
- CLI: `gws list -t=ai` continues to return filtered results demonstrates backward compatibility
- CLI: `gws list -yt ai` filters by tag and shows type column demonstrates stacking with values
- Test: Updated list command tests pass demonstrates all parsing paths are covered

### Unit 2: Implement Lowercase/Uppercase Flag Convention

**Purpose:** Give users explicit control over whether a filter flag also displays its associated column, following a predictable case-based pattern.

**Functional Requirements:**
- The system shall use lowercase short flags for filter-only behavior (no column displayed):
  - `-t value` filters by tag without showing the tags column
  - `-y value` filters by type without showing the type column
  - `-s value` filters by status without showing the status column
  - `-p value` filters by path without showing the path column
  - `-u value` filters by user without showing the user column
  - `-r value` filters by remote without showing the remote column
  - `-i value` filters by visibility without showing the visibility column
  - `-b value` filters by remote-raw without showing the remote-raw column
  - `-n value` filters by name (name column always displays — no uppercase variant needed)
- The system shall use uppercase short flags for show-column behavior (with optional filter):
  - `-T` (bare) shows the tags column; `-T value` shows the column and filters
  - `-Y` (bare) shows the type column; `-Y value` shows the column and filters
  - `-S` (bare) shows the status column; `-S value` shows the column and filters
  - `-P` (bare) shows the path column; `-P value` shows the column and filters
  - `-U` (bare) shows the user column; `-U value` shows the column and filters
  - `-R` (bare) shows the remote column; `-R value` shows the column and filters
  - `-I` (bare) shows the visibility column; `-I value` shows the column and filters
  - `-B` (bare) shows the remote-raw column; `-B value` shows the column and filters
- The system shall keep `-v` as the verbose count flag (unchanged behavior)
- The system shall remove the `NoOptDefVal` sentinel pattern from all flags

**Proof Artifacts:**
- CLI: `gws list -t ai` filters by tag without showing tag column demonstrates lowercase=filter-only
- CLI: `gws list -T ai` filters by tag and shows tag column demonstrates uppercase=show+filter
- CLI: `gws list -T` shows tag column without filtering demonstrates uppercase bare=show column
- CLI: `gws list -i private` filters by visibility without showing column demonstrates new `-i` flag
- CLI: `gws list -b github` filters by remote-raw without showing column demonstrates new `-b` flag
- Test: All list command tests updated and passing demonstrates comprehensive coverage

### Unit 3: Remap Existing Short Flags and Update Help Text

**Purpose:** Remap `-V` to `-i`/`-I` and `-R` to `-b`/`-B`, update deprecated flag layer, and ensure help text reflects the new convention.

**Functional Requirements:**
- The system shall map `--visibility` to `-i` (filter) and `-I` (show+filter), replacing `-V`
- The system shall map `--remote-raw` to `-b` (filter) and `-B` (show+filter), replacing `-R`
- The system shall map `--remote` to `-r` (filter) and `-R` (show+filter), freeing `-R` from remote-raw
- The system shall update help text for all flags to clearly describe the lowercase/uppercase convention
- The system shall update the deprecated flag layer to map old flag usage to new behavior where applicable
- The system shall display a deprecation warning when `-V` or old `-R` (for remote-raw) is used, guiding users to the new flags

**Proof Artifacts:**
- CLI: `gws list --help` shows updated flag descriptions with lowercase/uppercase convention documented demonstrates help text clarity
- CLI: `gws list -V private` emits deprecation warning and still works demonstrates backward compat with migration path
- Test: Deprecated flag mapping tests pass demonstrates backward compatibility layer works

### Unit 4: Apply Consistent Flag Behavior to Other Commands

**Purpose:** Extend the consistent flag pattern to `tag` and any other commands with filter flags, if their scope is reasonable.

**Functional Requirements:**
- The system shall ensure `tag --repo` / `-r` and `tag --path` / `-p` accept space-separated values consistently (verify current behavior)
- The system shall apply the same space-separated value behavior to any other command flags that currently require `=`
- The system shall not introduce the uppercase/lowercase column convention to commands that have no column display concept (e.g., `tag`)

**Proof Artifacts:**
- CLI: `gws tag add -r myrepo -p /path mytag` works with space-separated values demonstrates cross-command consistency
- Test: Tag command tests pass with space-separated flag values demonstrates coverage

## Non-Goals (Out of Scope)

1. **Adding new filter capabilities**: This spec only changes how existing filters are invoked, not what they can filter on
2. **Changing the `-v` verbose flag**: The verbose count flag stays as-is; this spec does not alter its behavior or column mappings
3. **Adding a `-N` uppercase flag for name**: Name is always displayed as the minimum column, so uppercase show+filter is unnecessary
4. **Changing long flag names**: Only short flag letters are being remapped; `--visibility`, `--remote-raw`, etc. keep their long names
5. **Modifying output format or column rendering**: How columns are displayed is unchanged; only which columns appear is affected by the new convention

## Design Considerations

No specific design requirements identified. This is a CLI flag parsing change with no visual design component.

## Repository Standards

- Go project using Cobra v1.10.2 with spf13/pflag v1.0.10
- Commit convention: `feat:` for minor version bump, `fix:` for patch, `docs:` for documentation
- Existing test patterns in `cmd/git-workspace/list_test.go` should be extended
- Deprecated flags are managed in `cmd/git-workspace/deprecated.go` with `MarkHidden()` and stderr warnings
- Flag completion functions are registered via `RegisterFlagCompletionFunc`

## Technical Considerations

- **Removing `NoOptDefVal`**: The current dual-purpose pattern relies on Cobra's `NoOptDefVal` to distinguish bare flags from flags with values. Removing it means Cobra will always expect a value for string flags. The new pattern uses separate lowercase/uppercase flag pairs instead of a single dual-purpose flag, so `NoOptDefVal` is no longer needed on lowercase flags. Uppercase flags will still need `NoOptDefVal` to support bare usage (show column without filter value).
- **Flag pair registration**: Each filterable column will require two flag registrations — one lowercase (standard `StringVarP`, no `NoOptDefVal`) and one uppercase (with `NoOptDefVal` for bare show-column support). This doubles the flag count but makes behavior explicit and predictable.
- **`parseDualPurposeFlags` refactor**: The existing function must be rewritten to handle the new two-flag-per-column pattern. Each column's visibility and filter value come from checking both the lowercase and uppercase flag.
- **Tag flag is repeatable**: The `-t`/`-T` flags support multiple values (AND logic). Both lowercase and uppercase variants must support this.
- **Stacking behavior**: POSIX stacking (`-YT`) works naturally with Cobra for bare `NoOptDefVal` flags. Stacking with a trailing value (`-YT ai`) will apply the value to the last flag in the stack (standard Cobra behavior).
- **Deprecated flag layer**: `deprecated.go` currently maps old root-level flags to the list subcommand. The old `-V` and `-R` (remote-raw) short flags need deprecation warnings pointing to `-i`/`-I` and `-b`/`-B` respectively.

## Security Considerations

No specific security considerations identified.

## Success Metrics

1. **Consistency**: All filter flags accept space-separated values without `=` — zero flags require `=`
2. **Backward compatibility**: All existing `=` syntax continues to work with no behavior change
3. **Convention clarity**: `--help` output clearly documents the lowercase/uppercase pattern
4. **Test coverage**: All new flag combinations have test coverage, including stacking and edge cases
5. **Zero regressions**: All existing tests continue to pass

## Open Questions

No open questions at this time.

## Resolved Questions

1. **`-v` verbose levels**: The `-v` verbose count flag will continue to auto-show columns at each level, unchanged. This provides a convenient shorthand alongside the new explicit uppercase flags.
2. **Deprecated `-V` mapping**: The deprecated `-V` flag will map to `-I` (show+filter) to retain current behavior, since `-V` bare currently shows the visibility column. `-V value` will map to `-I value` (show column + filter). A deprecation warning will guide users to the new flag.
