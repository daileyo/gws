# 13-spec-list-refactor

## Introduction/Overview

The `gws list` command currently always displays a full table with NAME, TYPE, VISIBILITY, TAGS, and PATH columns. This refactor transforms the list experience: `gws` with no arguments defaults to `gws list`, the default view shows only repo names in a multi-column layout (like `ls`), and individual columns can be toggled on/off via short flags that also double as filters when given a value.

## Goals

- Make `gws` (no arguments) default to `gws list` output while preserving positional-arg navigation
- Change the default `gws list` view to a compact multi-column display of repo names only
- Introduce dual-purpose flags that toggle column display (no value) or filter+display (with value)
- Add `-v`/`--verbose` for stored-data columns and `-vv` for all columns including live-fetched data
- Ensure `gws list -h` shows all available flags (fix any currently missing)
- Support JSON output that respects display flag selection

## User Stories

- **As a developer**, I want `gws` with no arguments to show my repos so that I get useful output immediately without remembering subcommands.
- **As a developer**, I want a compact default view showing just repo names so that I can quickly scan my workspace without visual clutter.
- **As a developer**, I want to selectively show columns (e.g., `-tp` for tags and path) so that I see only the information I need.
- **As a developer**, I want to filter and display in one flag (e.g., `-y github`) so that I can combine filtering and column selection concisely.
- **As a developer**, I want `-v` for a quick overview and `-vv` for full details so that I have tiered verbosity levels.

## Demoable Units of Work

### Unit 1: Default No-Args Behavior and Multi-Column View

**Purpose:** Make `gws` default to listing repos and display names in a compact multi-column layout.

**Functional Requirements:**
- The system shall run `gws list` when `gws` is invoked with no arguments and no flags (navigation with positional args continues to work unchanged)
- The system shall display repo names in a multi-column layout that auto-detects terminal width, filling left-to-right then top-to-bottom (like `ls`)
- The system shall sort repo names alphabetically in the multi-column view
- The system shall fall back to a single-column layout when output is not a TTY (piped or redirected)
- The system shall show the "Found N repositories" header before the multi-column output

**Proof Artifacts:**
- CLI: `gws` with no arguments displays multi-column repo names
- CLI: `gws my-repo` still navigates to the repo (unchanged behavior)
- CLI: `gws list | cat` outputs one repo name per line (non-TTY fallback)

### Unit 2: Dual-Purpose Column/Filter Flags

**Purpose:** Allow each data column to be toggled on via a flag, and filtered when a value is provided.

**Functional Requirements:**
- The system shall support these dual-purpose flags on `gws list`:
  - `-y`/`--type` — show TYPE column (no value) or filter by type (with value)
  - `-V`/`--visibility` — show VISIBILITY column (no value) or filter by visibility (with value)
  - `-t`/`--tag` — show TAGS column (no value) or filter by tag (with value, repeatable, AND logic)
  - `-p`/`--path` — show PATH column (no value) or filter by path pattern (with value)
  - `-s`/`--status` — show STATUS column (no value) or filter by status (with value)
  - `-u`/`--show-user` — show USER/EMAIL/SIGN columns (no value) or filter by user (with value)
  - `-r`/`--remote` — show REMOTE column (no value) or filter by remote URL (with value)
- The system shall always display the NAME column regardless of flags
- The system shall support POSIX flag stacking for display-only flags (e.g., `-yVtp` shows type, visibility, tags, path columns)
- The `-n`/`--name` flag shall remain filter-only (name column is always shown)
- When a flag provides a value, it shall both filter results AND display that column
- When `-v`/`--verbose` is combined with other flags/filters, all columns shall be displayed (verbose overrides selective display)

**Proof Artifacts:**
- CLI: `gws list -yVtp` shows NAME, TYPE, VISIBILITY, TAGS, PATH columns
- CLI: `gws list -y github -tp` filters on type=github, displays NAME, TYPE, TAGS, PATH
- CLI: `gws list -y github -V private -tp` filters on type=github AND visibility=private, displays NAME, TYPE, VISIBILITY, TAGS, PATH
- CLI: `gws list -v -y github` filters on type=github, displays all stored-data columns

### Unit 3: Verbose Levels (`-v` and `-vv`)

**Purpose:** Provide tiered verbosity for quick overview vs. full detail.

**Functional Requirements:**
- The system shall support `-v`/`--verbose` to show all stored-data columns: NAME, TYPE, VISIBILITY, TAGS, PATH (equivalent to current default table view)
- The system shall support `-vv` (double verbose) to show all columns including live-fetched data: NAME, STATUS, USER, EMAIL, SIGN, TYPE, VISIBILITY, TAGS, PATH, REMOTE
- The `-v` and `-vv` output shall use the existing table format (headers, separator lines, padded columns)
- When `-v` or `-vv` is combined with filter flags, the filter shall apply but all columns for that verbosity level shall display

**Proof Artifacts:**
- CLI: `gws list -v` shows NAME, TYPE, VISIBILITY, TAGS, PATH table
- CLI: `gws list -vv` shows all columns including STATUS, USER, EMAIL, SIGN, REMOTE
- CLI: `gws list -vv -y github` filters by github, shows all columns

### Unit 4: Help Text and JSON Output Updates

**Purpose:** Ensure help text is complete and JSON output respects column selection.

**Functional Requirements:**
- The system shall display all `gws list` flags when running `gws list -h` (fix any currently hidden flags)
- The system shall include the new `-V`/`--visibility` and `-v`/`--verbose` flags in help output
- JSON output (`-o json`) shall include only fields corresponding to displayed columns (respecting `-y`, `-V`, `-t`, `-p`, etc.)
- JSON output shall always include the `name` field
- Filter flags shall apply to JSON output (filtering which repos appear)
- The `-o`/`--output` flag shall continue to support `table` and `json` formats
- Deprecated root-level flags (`gws --list`, `gws --status`, etc.) shall continue working with deprecation warnings, mapping to new behavior

**Proof Artifacts:**
- CLI: `gws list -h` shows all flags including new ones
- CLI: `gws list -yp -o json` outputs JSON with only name, type, and path fields
- CLI: `gws list -y github -o json` filters by github and includes only name and type fields in JSON

## Non-Goals (Out of Scope)

1. **Changing navigation behavior**: `gws <repo-name>` positional-arg navigation remains unchanged
2. **Adding new data columns**: No new data fields beyond what exists (name, type, visibility, tags, path, status, user, email, sign, remote)
3. **Removing deprecated flags**: Old `gws --list` etc. continue with deprecation warnings; removal is a separate effort
4. **Color or styling**: No color/theme changes in this refactor
5. **Interactive/TUI features**: No interactive selection or scrolling

## Design Considerations

The multi-column default view should visually resemble `ls` output — clean, compact, and familiar to CLI users. When any column flag is passed, the output switches to the existing table format with headers and separator lines.

## Repository Standards

- Follow Conventional Commits format: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `style:`, `chore:`
- Table-driven tests in `_test.go` files co-located with source
- Use `go test ./...` to validate; `make ci` for full CI checks (vet + lint + race)
- Follow existing patterns in `list.go` for flag registration and `runList()` structure
- Keep `gofmt` formatting compliance

## Technical Considerations

- **Optional-value flags in Cobra**: Cobra does not natively support flags that accept an optional value (present = display, present with value = filter). Implementation will likely require custom flag types or `NoOptDefVal` to distinguish "flag present without value" from "flag not present". This is the primary technical challenge.
- **Terminal width detection**: Use `golang.org/x/term` or similar to detect terminal width for multi-column layout. Must handle non-TTY gracefully (single column fallback).
- **Flag counting for `-vv`**: Cobra's `CountP` or similar mechanism can track how many times `-v` is passed to distinguish `-v` from `-vv`.
- **Column ordering**: When selective columns are shown, maintain a consistent column order: NAME, STATUS, USER, EMAIL, SIGN, TYPE, VISIBILITY, TAGS, PATH, REMOTE (same as current `-rsu` full view).
- **Backward compatibility**: The `filterTags` variable is shared between list and deprecated flag paths — ensure changes don't break deprecated flag behavior.

## Security Considerations

No specific security considerations identified. No new credentials, tokens, or sensitive data handling is introduced.

## Success Metrics

1. **`gws` with no args** produces multi-column repo name output (replaces current summary message)
2. **All flags visible** in `gws list -h` output
3. **Flag stacking works**: `-yVtp` shows four columns; `-y github -V private` filters correctly
4. **Backward compatibility**: All deprecated flags continue working with deprecation warnings
5. **Tests pass**: All existing tests continue to pass; new tests cover dual-purpose flag behavior

## Open Questions

No open questions at this time. Resolved:

1. **Tag flag parsing**: `-t` without a value shows the TAGS column (no filtering). `-t personal` filters by tag AND shows the column. Implementation will use `NoOptDefVal` or custom flag type.
2. **Multi-column sort order**: Alphabetical by default. Other orderings are future work.
3. **Status cache**: Use cached data with existing 5-min TTL behavior (no change).
