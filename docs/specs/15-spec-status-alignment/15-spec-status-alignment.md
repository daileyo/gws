# 15-spec-status-alignment

## Introduction/Overview

The `gws list` STATUS column currently displays branch name and status icons (clean/dirty, ahead/behind) as a single left-justified string (e.g., `main ✓`). This makes it hard to visually scan status across many repos because the icons shift depending on branch name length. This feature restructures the STATUS column into two aligned sub-columns (branch left-justified, icons right-justified) and adds color to the status icons for faster visual parsing.

## Goals

- Right-justify status icons (`✓`/`✗`, `↑N`/`↓N`) within the STATUS column so they align vertically across rows
- Left-justify branch names within the STATUS column
- Add ANSI color to status icons: green for clean, red for dirty, cyan for ahead, magenta for behind
- Auto-detect TTY and disable colors when output is piped, with a `--color` flag to override
- Truncate long branch names with an ellipsis to keep the column manageable
- Fix the existing column width bug where `len()` counts bytes instead of visual character width for Unicode characters

## User Stories

- **As a developer managing multiple repos**, I want the status icons to be vertically aligned so that I can quickly scan which repos need attention without my eyes jumping around.
- **As a developer**, I want color-coded status icons so that I can instantly distinguish clean repos from dirty ones and see ahead/behind state at a glance.
- **As a developer piping `gws list` output**, I want colors to be automatically disabled so that ANSI codes don't pollute my downstream tools, but I want a `--color` flag to force colors on or off when needed.
- **As a developer with long branch names**, I want them truncated with an ellipsis so the STATUS column stays readable.

## Demoable Units of Work

### Unit 1: Status Column Split and Alignment

**Purpose:** Restructure the STATUS column so branch names are left-justified and status icons are right-justified, making the output scannable across many repos.

**Functional Requirements:**
- The system shall split the STATUS column into two visual sub-columns: branch name (left-justified) and status icons (right-justified)
- The system shall calculate the branch sub-column width based on the longest branch name across all displayed repos
- The system shall calculate the icons sub-column width based on the widest icon string (e.g., `✗ ↑12↓34`) across all displayed repos
- The system shall use Unicode-aware width calculation (rune count or display width) instead of `len()` for column width computation to fix the existing alignment bug
- The system shall pad between the branch name and icons with spaces so icons right-align within the STATUS column
- The separator line under the STATUS header shall match the total visual width of the combined sub-columns

**Proof Artifacts:**
- CLI: `gws list -v` output demonstrates branch names left-justified and icons right-justified with vertical alignment across repos
- CLI: Status column separator line (`---`) matches the visual width of the column content

### Unit 2: Branch Name Truncation

**Purpose:** Prevent long branch names from making the STATUS column excessively wide.

**Functional Requirements:**
- The system shall truncate branch names that exceed a configurable maximum length (default: 30 characters)
- The system shall append an ellipsis (`...`) to truncated branch names, so the displayed length is max_length characters total (including the ellipsis)
- The system shall use the truncated length for column width calculation
- The system shall not truncate branch names shorter than the maximum length

**Proof Artifacts:**
- CLI: `gws list -v` with a repo on a long branch name (e.g., `feature/very-long-branch-name-that-exceeds-limit`) demonstrates truncation with ellipsis
- CLI: Column alignment remains correct when some branches are truncated and others are not

### Unit 3: Colored Status Icons

**Purpose:** Add color to status icons so developers can visually distinguish repo states at a glance.

**Functional Requirements:**
- The system shall color the `✓` icon green (ANSI code 32)
- The system shall color the `✗` icon red (ANSI code 31)
- The system shall color ahead indicators (`↑N`) cyan (ANSI code 36)
- The system shall color behind indicators (`↓N`) magenta (ANSI code 35)
- The system shall use raw ANSI escape codes without adding a color library dependency
- The system shall not include ANSI escape codes in column width calculations (only visible characters count toward width)

**Proof Artifacts:**
- CLI: `gws list -v` output in a terminal demonstrates colored status icons with correct alignment
- CLI: Colors appear correctly for all combinations: clean, dirty, ahead, behind, ahead+behind

### Unit 4: TTY Detection and --color Flag

**Purpose:** Ensure colored output works correctly in both interactive and scripted contexts.

**Functional Requirements:**
- The system shall auto-detect whether stdout is a TTY and disable ANSI colors when output is piped or redirected
- The system shall provide a `--color` flag on the `list` command with values: `auto` (default), `always`, `never`
- When `--color=always`, the system shall emit ANSI colors regardless of TTY detection
- When `--color=never`, the system shall suppress all ANSI colors regardless of TTY detection
- When `--color=auto` (default), the system shall emit colors only when stdout is a TTY
- The existing `golang.org/x/term` dependency shall be used for TTY detection (already in go.mod)

**Proof Artifacts:**
- CLI: `gws list -v | cat` demonstrates colors are stripped when piped
- CLI: `gws list -v --color=always | cat` demonstrates colors are preserved when forced
- CLI: `gws list -v --color=never` demonstrates no color output in a terminal
- CLI: `gws list -v --help` demonstrates the `--color` flag in help text

## Non-Goals (Out of Scope)

1. **Configurable color themes**: Users cannot customize which colors map to which states. The colors are hardcoded.
2. **Color support for other columns**: Only the STATUS column icons receive color treatment. NAME, TYPE, TAGS, etc. remain uncolored.
3. **`NO_COLOR` environment variable support**: While a good practice, support for the `NO_COLOR` convention is deferred to a future enhancement.
4. **Branch name max-length configuration flag**: The truncation max length is a code constant, not a user-facing CLI flag.

## Design Considerations

The STATUS column visual layout changes from:

```
STATUS
------
main ✓
feature ✗ ↑3
develop ✓ ↓2
bugfix/long-name ✗ ↑1↓5
```

To (with right-justified icons, colors not shown):

```
STATUS
---------------------
main                ✓
feature          ✗ ↑3
develop          ✓ ↓2
bugfix/long-name ✗ ↑1↓5
```

The icons sub-column aligns to the right edge, and branch names align to the left edge. The combined width is `max_branch_len + gap + max_icons_len`.

## Repository Standards

- Go project using `cobra` for CLI commands
- Tests in `*_test.go` files alongside source
- Commit messages follow conventional commits: `feat:`, `fix:`, `docs:`, etc.
- Existing TTY detection uses `golang.org/x/term` (`term.IsTerminal()`)
- Table output uses `fmt.Sprintf("%-*s", width, value)` pattern with two-space column separators
- Unicode characters already used in output (`✓`, `✗`, `↑`, `↓`, `⚠`, `*`)

## Technical Considerations

- Column width calculation currently uses `len()` (byte count) which is incorrect for multi-byte Unicode. Must switch to a rune-aware or display-width-aware calculation. Since all current Unicode characters are single display-width, `utf8.RuneCountInString()` is sufficient (no full-width CJK characters in output).
- ANSI escape codes are zero-width in terminals but have non-zero byte and rune length. Width calculation must strip or exclude ANSI codes.
- The `formatStatusShort` function will need to be refactored into separate branch and icon components so they can be independently measured and padded.
- Raw ANSI codes: `\033[32m` (green), `\033[31m` (red), `\033[36m` (cyan), `\033[35m` (magenta), `\033[0m` (reset).
- The color decision (on/off) must be resolved once at command execution time and passed through the rendering pipeline, not checked per-cell.

## Security Considerations

No specific security considerations identified.

## Success Metrics

1. **Visual alignment**: Status icons align vertically across all rows when branch names vary in length
2. **Color correctness**: All four icon types display in their specified colors in a terminal
3. **Pipe safety**: `gws list -v | cat` produces clean output with no ANSI escape codes (by default)
4. **No regressions**: Existing tests pass; column alignment is correct for all column combinations
5. **No new dependencies**: Implementation uses raw ANSI codes with no new Go module dependencies

## Open Questions

1. What should the default branch name truncation length be? (Spec proposes 30 characters — is this appropriate for your typical branch naming conventions?)
