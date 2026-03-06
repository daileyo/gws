# 15 Tasks - Status Alignment

## Relevant Files

- `cmd/git-workspace/list.go` - Contains `formatStatusShort`, column width calculation, table rendering, and `listCmd` flag registration. Primary file for all changes.
- `cmd/git-workspace/list_test.go` - Unit tests covering `formatStatusShort` refactor, color utilities, truncation, and alignment logic.
- `internal/git/status.go` - Status data model (`Status` struct, `String()` method). Read-only reference; no changes expected.
- `internal/git/status_test.go` - Existing status tests. Must continue to pass (regression check).

### Notes

- Unit tests should be placed in `cmd/git-workspace/list_test.go` (same package as the functions under test).
- Run tests with `go test ./cmd/git-workspace/... -run <TestName>`.
- Follow existing Go conventions: `snake_case` test names with `Test` prefix, table-driven tests.
- Commit messages follow conventional commits: `feat:` for new functionality.

## Tasks

### [x] 1.0 ANSI Color Utilities and TTY Detection

Build the foundational color infrastructure: raw ANSI color helper functions, a color-enabled flag resolved from TTY detection, and the `--color` flag on the `list` command. This is the prerequisite for all colored output.

#### 1.0 Proof Artifact(s)

- CLI: `gws list -v --help` demonstrates the `--color` flag with `auto`, `always`, `never` values in help text
- Test: Color utility unit tests pass demonstrating correct ANSI wrapping and strip functions
- CLI: `gws list -v --color=never` produces output with no ANSI escape codes in a terminal

#### 1.0 Tasks

- [x] 1.1 Create ANSI color constants and a `colorize(text, code string) string` helper function in `list.go` that wraps text with `\033[<code>m...\033[0m`. Define constants for green (32), red (31), cyan (36), magenta (35), and reset (0).
- [x] 1.2 Create a `stripANSI(s string) string` function that removes all ANSI escape sequences from a string, for use in width calculations. Use a simple regex or manual scan for `\033[...m` patterns.
- [x] 1.3 Create a `displayWidth(s string) int` function that returns the visible character count of a string by stripping ANSI codes and then using `utf8.RuneCountInString()`.
- [x] 1.4 Add a `--color` flag to `listCmd` in `init()` with a package-level `flagColor` variable. Accept values `auto` (default), `always`, `never`. Register it as a regular string flag (not dual-purpose).
- [x] 1.5 Add a `resolveColorEnabled(cmd *cobra.Command) bool` function that reads the `--color` flag value and returns true/false. For `auto`, use `term.IsTerminal(int(os.Stdout.Fd()))`. For `always`, return true. For `never`, return false.
- [x] 1.6 Call `resolveColorEnabled` at the start of `runList` and store the result in a variable that will be passed to rendering functions (e.g., add a `ColorEnabled bool` field to `ListOptions` or use a package-level variable).
- [x] 1.7 Write unit tests in `list_test.go` for `colorize`, `stripANSI`, and `displayWidth` covering: plain text, text with ANSI codes, multi-byte Unicode characters (`✓`, `✗`, `↑`, `↓`), and empty strings.

### [x] 2.0 Refactor formatStatusShort into Split Components

Refactor `formatStatusShort` to return separate branch name and icons strings instead of a single combined string. Switch column width calculation from `len()` to `utf8.RuneCountInString()` for Unicode-correct measurement. This is the structural change that enables independent alignment of each sub-column.

#### 2.0 Proof Artifact(s)

- Test: Unit tests for the refactored format functions pass demonstrating correct branch/icon separation
- CLI: `gws list -v` output still displays correctly (no visual regression) after refactor
- Test: Existing `status_test.go` tests continue to pass

#### 2.0 Tasks

- [x] 2.1 Create a `formatStatusBranch(status *git.Status) string` function that returns only the branch name (or `"no commits"` when Branch is empty).
- [x] 2.2 Create a `formatStatusIcons(status *git.Status) string` function that returns only the icons string: clean/dirty indicator (`✓`/`✗`) followed by ahead/behind arrows (`↑N`, `↓N`). Return empty string when Branch is empty (the `"no commits"` case has no icons).
- [x] 2.3 Update `formatStatusShort` to call `formatStatusBranch` and `formatStatusIcons` internally, joining them with a space, so existing callers (including `displayJSON`) are unaffected.
- [x] 2.4 In the column width pre-computation loop (around line 444-452), replace the single `maxStatusLen` calculation with two separate calculations: `maxBranchLen` (using `displayWidth` on branch string) and `maxIconsLen` (using `displayWidth` on icons string). Keep `maxStatusLen` as `maxBranchLen + 1 + maxIconsLen` (1 for the space separator) for header/separator compatibility.
- [x] 2.5 Write unit tests for `formatStatusBranch` and `formatStatusIcons` covering: clean repo, dirty repo, ahead only, behind only, ahead+behind, and no-commits cases.

### [x] 3.0 Status Column Split Alignment and Branch Truncation

Implement the two-sub-column layout: branch name left-justified, icons right-justified within the STATUS column. Add branch name truncation with ellipsis for names exceeding 30 characters. Fix separator line width to match visual column width.

#### 3.0 Proof Artifact(s)

- CLI: `gws list -v` output demonstrates branch names left-justified and status icons right-justified with vertical alignment across repos with varying branch name lengths
- CLI: A repo with a branch name longer than 30 characters displays truncated with `...` suffix
- CLI: Status column separator line (`---`) matches the visual width of the column content

#### 3.0 Tasks

- [x] 3.1 Add a `const maxBranchDisplayLen = 30` constant in `list.go`.
- [x] 3.2 Create a `truncateBranch(name string, maxLen int) string` function that truncates branch names exceeding `maxLen` by trimming to `maxLen - 3` characters and appending `...`. Return the name unchanged if it fits within `maxLen`.
- [x] 3.3 Apply `truncateBranch` in `formatStatusBranch` before returning the branch name, so all downstream width calculations and rendering use the truncated value.
- [x] 3.4 Update the STATUS column row rendering (around line 578-584) to format the status as two sub-columns: `fmt.Sprintf("%-*s %*s", maxBranchLen, branch, maxIconsLen, icons)` — branch left-justified, icons right-justified. Use `displayWidth` for padding calculations to account for future ANSI codes.
- [x] 3.5 Update the STATUS header and separator to use the new combined width (`maxBranchLen + 1 + maxIconsLen`), ensuring the separator dashes match the visual column width.
- [x] 3.6 Handle the edge case where status is `"-"` (no status data) or `"no commits"` — these should span the full STATUS column width left-justified without icon sub-column formatting.
- [x] 3.7 Write unit tests for `truncateBranch` covering: short name (no truncation), exactly at limit, one over limit, very long name, and name shorter than 4 characters (edge case where `maxLen - 3` could be problematic).

### [x] 4.0 Apply Colors to Status Icons

Wire the color utilities into the status icon rendering: green `✓`, red `✗`, cyan `↑N`, magenta `↓N`. Ensure ANSI codes are excluded from width calculations so alignment remains correct. Verify TTY auto-detection and `--color` flag work end-to-end.

#### 4.0 Proof Artifact(s)

- CLI: `gws list -v` in a terminal demonstrates colored status icons with correct alignment
- CLI: `gws list -v | cat` demonstrates colors are stripped when piped (auto mode)
- CLI: `gws list -v --color=always | cat` demonstrates colors are preserved when forced
- CLI: Colors appear correctly for all combinations: clean, dirty, ahead, behind, ahead+behind

#### 4.0 Tasks

- [x] 4.1 Update `formatStatusIcons` to accept a `colorEnabled bool` parameter. When true, wrap `✓` with green, `✗` with red, `↑N` with cyan, and `↓N` with magenta using the `colorize` helper.
- [x] 4.2 Thread the `colorEnabled` value from `runList` through to all call sites of `formatStatusIcons` — both in the width pre-computation loop and in the row rendering loop.
- [x] 4.3 Ensure the width pre-computation loop uses `displayWidth` (which strips ANSI) so that colored icon strings do not inflate column widths.
- [x] 4.4 Verify that `formatStatusShort` (used by `displayJSON`) does NOT pass `colorEnabled=true`, so JSON output never contains ANSI codes.
- [x] 4.5 Write unit tests verifying: (a) `formatStatusIcons` with `colorEnabled=true` produces strings containing ANSI codes, (b) `displayWidth` of colored output equals `displayWidth` of uncolored output, (c) `formatStatusIcons` with `colorEnabled=false` produces no ANSI codes.
- [x] 4.6 Manual end-to-end verification: run `gws list -v` in terminal (colors visible), `gws list -v | cat` (no colors), `gws list -v --color=always | cat` (colors forced), `gws list -v --color=never` (no colors in terminal).
