# 19-spec-status-compact

## Introduction/Overview

The `gws list` command currently uses `-s` as a filter-only flag for status and `-S` to display a full STATUS column (branch + icons). Users want a lightweight way to see repo health at a glance without the overhead of a separate column. This feature changes `-s` to show compact status icons (ahead/behind/clean/dirty) right-justified within the existing NAME column, while `-S` retains its current full-column behavior.

## Goals

- Provide a quick, at-a-glance view of repository status without adding a new column
- Keep the existing `-S` (full status column) behavior unchanged
- Reuse existing status icon format and colors for consistency
- Maintain filtering capability when `-s` is given a value (e.g., `-s dirty`)
- Minimize visual clutter — icons only, no branch name in compact mode

## User Stories

- **As a developer managing many repos**, I want to quickly see which repos are dirty or have unpushed commits so that I can prioritize my work without scrolling through a wide table.
- **As a developer who prefers minimal output**, I want status indicators embedded in the name column so that I don't need a separate column just to see clean/dirty state.
- **As a developer who needs full details**, I want `-S` to continue working exactly as it does today so that I can still see branch names and full status when needed.

## Demoable Units of Work

### Unit 1: Compact Status Display in NAME Column

**Purpose:** When `-s` is used without a value, display compact status icons (ahead/behind/clean/dirty) right-justified within the NAME column instead of showing a separate STATUS column.

**Functional Requirements:**
- The system shall display status icons right-justified in the NAME column when `-s` is specified
- The system shall use the same icon format as the current full status: `↓N` (behind), `↑N` (ahead), `✓` (clean), `✗` (dirty), with spaces between icons
- The system shall apply the same ANSI colors: magenta for behind, cyan for ahead, green for clean, red for dirty
- The system shall show `✓` for repos that are clean with no ahead/behind
- The system shall expand the NAME column width to accommodate the longest repo name + padding + longest icon string across all displayed repos
- The system shall NOT display the branch name in compact mode — only status icons
- The system shall NOT add compact status data to JSON output; JSON status data requires `-S`
- The system shall NOT display a separate STATUS column header when only `-s` is used

**Proof Artifacts:**
- CLI: `gws list -s` output shows repo names with status icons right-justified in the NAME column, no STATUS column present
- CLI: `gws list -s` with color enabled shows correctly colored icons
- CLI: `gws list -s` shows `✓` for clean repos with no ahead/behind

### Unit 2: Filtering with Compact Display

**Purpose:** Preserve and activate the existing `-s` filter behavior so that `-s dirty` both shows compact status AND filters results.

**Functional Requirements:**
- The system shall filter repos by status pattern when `-s` is given a value (e.g., `-s dirty`, `-s ahead`)
- The system shall display compact status icons in the NAME column when filtering with `-s <pattern>`
- The system shall support the same partial-match filtering semantics as other lowercase filter flags

**Proof Artifacts:**
- CLI: `gws list -s dirty` shows only dirty repos with compact status icons in the NAME column
- CLI: `gws list -s clean` shows only clean repos with `✓` in the NAME column

### Unit 3: `-S` Unchanged and `-s`/`-S` Combination

**Purpose:** Ensure `-S` continues to work exactly as it does today, and handle the case where both flags are specified.

**Functional Requirements:**
- The system shall display the full STATUS column (branch + icons) when `-S` is specified, identical to current behavior
- The system shall treat `-S` as the winner when both `-s` and `-S` are specified — showing the full STATUS column and suppressing compact icons in the NAME column
- The system shall still apply any filter value from `-s` when both flags are present (e.g., `-s dirty -S` shows full status column filtered to dirty repos)

**Proof Artifacts:**
- CLI: `gws list -S` output is identical to current behavior (full STATUS column with branch + icons)
- CLI: `gws list -s -S` shows full STATUS column without compact icons in NAME column
- CLI: `gws list -s dirty -S` shows full STATUS column filtered to dirty repos only

## Non-Goals (Out of Scope)

1. **Changing the `-S` display format**: The full status column behavior remains exactly as-is
2. **Adding branch name to compact display**: Compact mode shows only status icons, not branch names
3. **JSON output for compact status**: No new JSON fields; status in JSON requires `-S`
4. **New icons or color schemes**: Reuse existing icon and color conventions

## Design Considerations

Example output for `gws list -s`:

```
NAME
------------------------------
my-repo                ↓2 ↑1 ✗
another-repo                  ✓
some-project              ↑3 ✗
clean-repo                    ✓
```

- Icons are right-justified within the NAME column
- The NAME column expands to fit the widest combination of name + icons
- Minimum padding of at least 2 spaces between the longest repo name and its icons

## Repository Standards

- Follow Conventional Commits for all commit messages
- Use table-driven tests with the `TestFunctionName_Scenario` naming pattern
- Run `make ci` (vet, lint, test with race detector) to validate changes
- Follow existing patterns in `list.go` for column formatting and width calculations

## Technical Considerations

- The `-s` flag currently uses `NoOptDefVal` similar to `-S` to distinguish "flag present without value" from "flag not present." This mechanism should be extended or adapted so that bare `-s` triggers compact display while `-s dirty` triggers compact display + filtering.
- Status fetching (cache lookup, `FetchAll`) must be triggered when `-s` is used, same as when `-S` is used — the status data source is the same.
- The `formatStatusIcons` function can be reused directly for generating the compact icon string.
- The NAME column width calculation in the display loop must account for icon display widths (including ANSI escape codes) when `-s` is active.
- The `displayWidth` helper (which strips ANSI codes) should be used for correct padding calculations with colored output.

## Security Considerations

No specific security considerations identified.

## Success Metrics

1. **Functional correctness**: `gws list -s` shows compact status icons right-justified in NAME column for all repos
2. **No regression**: `gws list -S` output is byte-identical to current behavior
3. **Performance**: No additional git operations — reuses existing status cache and fetch mechanism
4. **All existing tests pass**: No regressions in `make ci`

## Open Questions

No open questions at this time.
