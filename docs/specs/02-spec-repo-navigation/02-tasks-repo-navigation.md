# 02-tasks-repo-navigation

Task list for implementing repository navigation as defined in [02-spec-repo-navigation.md](./02-spec-repo-navigation.md).

## Relevant Files

- `cmd/gws/main.go` - Root command, flag registration, mutual exclusivity check, dispatch logic. Navigation flag (`--go`/`-g`, `--quiet`/`-q`) registration and positional arg routing added here.
- `cmd/gws/main_test.go` - Tests for flag registration, mutual exclusivity, and navigation dispatch logic.
- `cmd/gws/navigate.go` - New file. Contains `runNavigate()` handler with matching, output formatting, interactive selection, and suggestion logic.
- `cmd/gws/navigate_test.go` - New file. Tests for navigation handler: single match, multiple match, no match, suggestions, TTY behavior, wildcard support.
- `internal/filter/filter.go` - Export `matchesPattern()` to `MatchesPattern()` so navigation can reuse wildcard/partial matching.
- `internal/filter/filter_test.go` - Update references if `matchesPattern` is renamed/exported.
- `README.md` - Add repository navigation section with shell wrapper documentation.

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `navigate.go` and `navigate_test.go` in `cmd/gws/`).
- Use `go test ./...` or `make test` to run all tests.
- Follow existing Go conventions: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests.
- Follow conventional commits for all changes.
- Run `make lint` before committing.

## Tasks

### [x] 1.0 Basic Navigation with Single Match

Implement core navigation: `--go`/`-g` flag and positional argument support, case-insensitive partial matching via the existing filter package, verbose/quiet output modes, and mutual exclusivity with other command flags. When exactly one repository matches, print the path to stdout. Verbose output (repo name and type) goes to stderr so shell command substitution always works.

#### 1.0 Proof Artifact(s)

- CLI: `gws my-repo` outputs verbose info to stderr and path to stdout demonstrates positional navigation
- CLI: `gws --go my-repo` outputs the path demonstrates flag-based navigation
- CLI: `gws -g my-repo -q` outputs only the path demonstrates quiet mode
- CLI: `gws --go my-repo --list` returns mutual exclusivity error demonstrates flag validation
- Test: Navigation unit tests in `cmd/gws/navigate_test.go` pass demonstrates core matching and output logic

#### 1.0 Tasks

- [x] 1.1 Export `matchesPattern()` in `internal/filter/filter.go` — rename to `MatchesPattern()` (exported). Update all internal call sites (`matchesCriteria`) to use the new name. Run `make test` to verify no regressions.
- [x] 1.2 Register `--go`/`-g` as a string flag and `--quiet`/`-q` as a boolean flag on the root command in `cmd/gws/main.go`. Add `flagGo` and `flagQuiet` to the command flag variables.
- [x] 1.3 Update the mutual exclusivity check in `rootCmd.RunE` to include `flagGo` in the active count. Update the error message to list `--go`.
- [x] 1.4 Add positional argument routing: when no command flag is set and `len(args) > 0`, treat `args[0]` as the navigation query. If `--go` is set and positional args are also provided, return an error. Add dispatch to `runNavigate()`.
- [x] 1.5 Create `cmd/gws/navigate.go` with `runNavigate(query string, quiet bool, repos []config.Repository, stderr io.Writer, stdout io.Writer) error`. Implement single-match logic: use `filter.MatchesPattern()` to find matching repos by name. When exactly one match, print path to stdout. In verbose mode, print `name (type) → path` to stderr. In quiet mode, print only the path to stdout.
- [x] 1.6 Write tests in `cmd/gws/navigate_test.go`: single exact match (verbose), single exact match (quiet), single partial match, flag registration for `--go` and `--quiet`, mutual exclusivity with `--go`, positional arg navigation, error when both `--go` and positional arg provided.
- [x] 1.7 Update `TestCommandFlagsRegistered` in `cmd/gws/main_test.go` to include `{"go", "g"}` and `{"quiet", "q"}`. Update `TestMutualExclusivity` error message assertion to include `--go`.
- [x] 1.8 Run `make test` and `make lint` to verify all tests pass and no lint errors.

### [x] 2.0 Multiple Match Interactive Selection

When a navigation query matches multiple repositories, display a numbered list and prompt the user to select one. Includes TTY detection — when stdin is not a TTY, print all matching paths and exit with non-zero status. Add wildcard pattern support (`*`, `?`) to navigation queries by reusing the exported `MatchesPattern()` from the filter package.

#### 2.0 Proof Artifact(s)

- CLI: `gws "api"` with multiple matches displays numbered list and accepts selection demonstrates interactive selection
- CLI: `echo "" | gws "api"` (non-TTY) prints all matches and exits non-zero demonstrates non-TTY behavior
- CLI: `gws "my-*"` with wildcard matches correct repos demonstrates wildcard support
- Test: Multiple match, selection, TTY detection, and wildcard tests pass demonstrates selection and pattern logic

#### 2.0 Tasks

- [x] 2.1 Add TTY detection helper in `cmd/gws/navigate.go`: `func isTerminal(r io.Reader) bool` using `os.File` type assertion and `r.(*os.File).Stat()` to check `ModeCharDevice`, avoiding a new dependency on `golang.org/x/term`.
- [x] 2.2 Extend `runNavigate()` to accept an `stdin io.Reader` parameter. When multiple repos match and stdin is a TTY, display the numbered list to stderr (format: `  1) name (type) path`) and prompt `Select repository [1-N]:` on stderr. Read the user's selection from stdin, validate it, and print the selected repo's path (respecting verbose/quiet mode).
- [x] 2.3 Handle invalid selection input: if the user enters a non-numeric value or out-of-range number, print an error to stderr and re-prompt (up to 3 attempts, then exit with error).
- [x] 2.4 Handle non-TTY stdin: when multiple repos match and stdin is not a TTY, print all matching paths (one per line) to stdout and return an error (non-zero exit) to indicate ambiguity.
- [x] 2.5 Verify wildcard support works end-to-end: `MatchesPattern()` already supports `*` and `?`, so navigation queries like `"my-*"` should match correctly. Add test cases for wildcard queries with single and multiple results.
- [x] 2.6 Write tests in `cmd/gws/navigate_test.go`: multiple matches with simulated stdin selection, invalid selection re-prompt, non-TTY prints all paths, wildcard single match, wildcard multiple matches. Use `strings.NewReader` to simulate stdin input in tests.
- [x] 2.7 Run `make test` and `make lint`.

### [x] 3.0 No Match with Suggestions

When no repositories match the navigation query, display an error with up to 5 substring-based suggestions of similar repository names. Exit with non-zero status.

#### 3.0 Proof Artifact(s)

- CLI: `gws nonexistent` shows error with suggestions (when similar names exist) demonstrates suggestion behavior
- CLI: `gws xyzzy123` shows error without suggestions (when nothing is similar) demonstrates graceful no-suggestion case
- Test: No-match and suggestion tests pass demonstrates suggestion logic

#### 3.0 Tasks

- [x] 3.1 Add `findSuggestions(query string, repos []config.Repository, max int) []string` in `cmd/gws/navigate.go`. For each repository, check if any substring of the query (minimum 2 characters) appears in the repo name, or if any substring of the repo name appears in the query, using case-insensitive comparison. Return up to `max` matching names.
- [x] 3.2 Extend the no-match path in `runNavigate()`: after printing the error message `No repositories found matching '<query>'` to stderr, call `findSuggestions()` and if results are non-empty, print `Did you mean:` followed by each suggestion indented with `  ` to stderr. Return an error for non-zero exit.
- [x] 3.3 Write tests in `cmd/gws/navigate_test.go`: no match with suggestions found (verify up to 5 shown), no match with zero suggestions (verify only error message), suggestion relevance (verify substring matching logic).
- [x] 3.4 Run `make test` and `make lint`.

### [ ] 4.0 Shell Integration and Documentation

Update README with repository navigation documentation including shell wrapper functions (`cdg`), eval-based alternative, and usage examples. Update the root command `--help` long description to include navigation examples. Ensure existing `cdgws`/`gcd` documentation remains unchanged.

#### 4.0 Proof Artifact(s)

- CLI: `cd "$(gws -g my-repo -q)"` in a shell successfully changes directory demonstrates end-to-end shell integration
- CLI: `gws --help` output includes navigation examples demonstrates help text
- README: Updated navigation section with `cdg` wrapper function demonstrates documentation completeness

#### 4.0 Tasks

- [ ] 4.1 Update the `Long` description in `rootCmd` (`cmd/gws/main.go`) to include navigation examples: `gws my-repo` (positional), `gws --go my-repo` (flag), `gws -g my-repo -q` (quiet). Update the shell integration section to show both `cdgws` (workspace) and `cdg` (repo navigation).
- [ ] 4.2 Add a "Repository Navigation" section to `README.md` after the existing "Shell Navigation" section. Include: basic usage (`gws <name>`), flag usage (`gws --go <name>`), quiet mode (`gws -g <name> -q`), wildcard examples, multiple match behavior, and shell wrapper setup.
- [ ] 4.3 Document the `cdg` shell wrapper function in README: `function cdg() { cd "$(gws -g "$1" -q)"; }` with alias `alias cdr=cdg`. Document the eval alternative: `eval "$(gws -g my-repo -q | xargs -I{} echo cd {})"`.
- [ ] 4.4 Verify existing `cdgws`/`gcd` documentation in README is unchanged and still references `--print-workspace`.
- [ ] 4.5 Run `make test` and `make lint` to verify no regressions from help text changes.
