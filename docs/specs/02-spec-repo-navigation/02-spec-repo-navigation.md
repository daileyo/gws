# 02-spec-repo-navigation.md

## Introduction/Overview

This feature adds repository navigation to gws, allowing users to quickly look up a repository's path by name. Since a CLI tool cannot change the parent shell's working directory, gws prints the matching repository path to stdout, which users can wrap in a shell function (e.g., `cd "$(gws my-repo)"`) for seamless navigation. Navigation supports both positional arguments and a `--go`/`-g` flag, wildcard pattern matching, interactive selection when multiple repos match, and fuzzy suggestions when no match is found.

## Goals

- Enable users to quickly find a repository's path by name using `gws <name>` or `gws --go <name>`
- Support wildcard patterns (`*`, `?`) consistent with Spec 03's filter behavior
- Provide interactive selection when multiple repositories match a query
- Offer fuzzy/partial-match suggestions when no repositories match
- Support both verbose (default) and quiet output modes for human and scripting use cases

## User Stories

- **As a developer with many repositories**, I want to type `gws my-repo` and get the path to that repository so that I can quickly navigate to it using a shell wrapper.

- **As a power user**, I want to use wildcards like `gws "api-*"` to see all matching repositories and select one so that I don't need to remember exact names.

- **As a developer scripting with gws**, I want a `--quiet` flag that outputs only the path so that I can use gws in shell scripts and pipelines without parsing extra output.

- **As a developer who can't remember exact names**, I want gws to suggest similar repository names when my query doesn't match so that I can correct my input quickly.

## Demoable Units of Work

### Unit 1: Basic Navigation with Single Match

**Purpose:** Establish the core navigation capability where a user can look up a repository path by name using either positional arguments or the `--go` flag.

**Functional Requirements:**
- The system shall accept a positional argument as a repository name to navigate to (e.g., `gws my-repo`)
- The system shall accept `--go` / `-g` as a string flag taking a repository name (e.g., `gws --go my-repo` or `gws -g my-repo`)
- The system shall be mutually exclusive with other command flags (`--list`, `--init`, `--add-tag`, `--remove-tag`, `--refresh`, `--print-workspace`)
- The system shall match repository names using case-insensitive partial matching (consistent with existing filter behavior)
- When exactly one repository matches, the system shall print the repository path to stdout
- In verbose mode (default), the system shall display the repo name and type followed by the path (e.g., `my-repo (github) → /home/user/projects/my-repo`)
- The system shall accept `--quiet` / `-q` as a boolean flag that suppresses all output except the path
- In quiet mode, the system shall print only the absolute path to stdout with no additional text
- The `--quiet` flag shall only apply when navigation is active (positional arg or `--go` flag)
- The system shall display an error when no workspace is initialized

**Proof Artifacts:**
- CLI: `gws my-repo` outputs the path demonstrates positional navigation works
- CLI: `gws --go my-repo` outputs the path demonstrates flag-based navigation works
- CLI: `gws -g my-repo -q` outputs only the path demonstrates quiet mode
- Test: Navigation unit tests pass demonstrates core matching logic

### Unit 2: Multiple Match Selection

**Purpose:** When a navigation query matches multiple repositories, provide an interactive selection mechanism so the user can choose the intended repository.

**Functional Requirements:**
- When multiple repositories match the query, the system shall display a numbered list of all matches
- Each entry in the list shall show the index number, repository name, type, and path
- The system shall prompt the user to enter a number to select a repository
- After selection, the system shall print the selected repository's path (respecting verbose/quiet mode)
- If the user enters an invalid selection, the system shall display an error and re-prompt
- If stdin is not a TTY (piped input), the system shall print all matching paths (one per line) and exit with a non-zero status to indicate ambiguity
- The system shall support wildcard patterns (`*`, `?`) in the repository name query, using the same `matchesPattern()` function from the filter package

**Proof Artifacts:**
- CLI: `gws "my-*"` with multiple matches displays a numbered list and accepts selection demonstrates interactive selection
- CLI: Piping `gws "my-*"` shows all matches without prompting demonstrates non-TTY behavior
- Test: Multiple match and selection tests pass demonstrates selection logic

### Unit 3: No Match with Suggestions

**Purpose:** When no repositories match, provide helpful suggestions so the user can correct their query.

**Functional Requirements:**
- When no repositories match the query, the system shall display an error message: `No repositories found matching '<query>'`
- The system shall search for similar repository names using partial/substring matching against all tracked repositories
- If similar names are found, the system shall display them as suggestions (e.g., `Did you mean: my-api, my-app?`)
- Suggestions shall be limited to a maximum of 5 entries to avoid overwhelming output
- The system shall exit with a non-zero status code when no match is found

**Proof Artifacts:**
- CLI: `gws nonexistent` shows error with suggestions demonstrates suggestion behavior
- CLI: `gws xyz` with no similar names shows error without suggestions demonstrates graceful no-suggestion case
- Test: No-match and suggestion tests pass demonstrates suggestion logic

### Unit 4: Shell Integration and Documentation

**Purpose:** Provide shell wrapper functions and documentation so users can use gws navigation for actual directory changes.

**Functional Requirements:**
- The system shall print only the path to stdout (verbose details to stderr or suppressed with `--quiet`) so that shell command substitution works correctly
- The README shall document a shell wrapper function for bash/zsh: `function cdg() { cd "$(gws -g "$1" -q)"; }`
- The README shall document an eval-based alternative: `eval "$(gws -g my-repo -q | xargs -I{} echo cd {})"`
- The `--help` output for the root command shall include a navigation example in the long description
- The existing `cdgws`/`gcd` shell function for workspace navigation shall remain unchanged

**Proof Artifacts:**
- CLI: `cd "$(gws -g my-repo -q)"` in a shell successfully changes directory demonstrates end-to-end shell integration
- README: Updated navigation section demonstrates documentation completeness

## Non-Goals (Out of Scope)

1. **Actual directory changing**: gws will NOT change the parent shell's directory directly; it prints paths for shell wrappers to consume
2. **Fuzzy matching library**: No external fuzzy matching library (e.g., `fzf`); suggestions use simple substring matching from existing code
3. **TUI framework**: No external TUI library (e.g., `bubbletea`, `promptui`); interactive selection uses standard I/O with numbered prompts
4. **Shell completion scripts**: Generating bash/zsh/fish completion scripts is a separate concern
5. **Navigation history**: Tracking or recalling previously navigated repositories

## Design Considerations

**Verbose Output Format:**
```
my-repo (github) → /home/user/projects/my-repo
```

**Quiet Output Format:**
```
/home/user/projects/my-repo
```

**Multiple Match Display:**
```
Multiple repositories match 'api':

  1) my-api (github) /home/user/projects/my-api
  2) my-api-v2 (github) /home/user/projects/my-api-v2
  3) work-api (gitlab) /home/user/projects/work-api

Select repository [1-3]:
```

**No Match with Suggestions:**
```
No repositories found matching 'aip'

Did you mean:
  my-api
  work-api
```

## Repository Standards

- Follow existing Go code patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests
- Use Cobra/pflag conventions for flag registration on the root command
- Reuse existing `matchesPattern()` from `internal/filter` for wildcard support
- Maintain mutual exclusivity validation in the root command's `RunE` handler
- Follow conventional commits for version history
- Use `make test` and `make lint` for validation

## Technical Considerations

- **Stdout vs Stderr**: Verbose output (repo name, type) must go to stderr so that stdout contains only the path for shell command substitution. This ensures `cd "$(gws my-repo)"` works correctly even in verbose mode. The `--quiet` flag suppresses stderr output entirely.
- **TTY Detection**: Use `os.Stdin` with `isatty` check (via `golang.org/x/term` or `os.Stdin.Stat()`) to determine if interactive prompting is appropriate. When stdin is not a TTY, skip prompting and output all matches.
- **Wildcard Reuse**: Leverage the existing `matchesPattern()` function from `internal/filter/filter.go` for wildcard matching. This may require exporting the function or extracting it to a shared utility.
- **Flag Registration**: `--go`/`-g` is a string flag on the root command. `--quiet`/`-q` is a boolean flag. Both are registered alongside existing command flags. The mutual exclusivity check must include `--go`.
- **Positional Arg vs Flag**: When a positional argument is provided and no command flag is set, treat it as navigation. When `--go` is explicitly set, use its value. If both are provided, return an error.

## Security Considerations

No specific security considerations identified. Navigation only reads repository metadata from the existing config file and prints paths. No new file access, network calls, or credential handling is introduced.

## Success Metrics

1. **Single match navigation**: `gws <exact-name>` returns the correct path in under 100ms
2. **Multiple match selection**: Users can select from matches with a single keystroke (number + Enter)
3. **Suggestion quality**: Substring-based suggestions surface the intended repo in the top 5 results
4. **Shell integration**: `cd "$(gws -g <name> -q)"` works correctly in both bash and zsh
5. **No regressions**: All existing tests pass, `make ci` succeeds

## Open Questions

No open questions at this time.

---

**Next Steps**: Once this spec is approved, run `/SDD-2-generate-task-list-from-spec` to create the implementation task list.
