# 09-spec-command-structure-refactor

## Introduction/Overview

The gws CLI currently uses `--` prefixed flags on the root command for all operations (`gws --list`, `gws --init`, `gws --add`). This pattern, introduced in v2.0 to help tab completion distinguish commands from repo names, has made the interface less intuitive and harder to extend. This spec refactors the CLI back to a subcommand-based structure (`gws list`, `gws init`, `gws add`) while adding single-letter root aliases (`gws -l`, `gws -i`, `gws -a`) for power users. Filter flags move from the root command to their respective subcommands, enabling reusable short flags across commands without conflicts. A deprecation layer preserves the old `--` flag forms with warnings for one release cycle.

This is spec 1 of 3. It covers the core command structure, `list` subcommand with filters, and the deprecation layer. Tag and User command refactors are handled in separate follow-up specs (10 and 11).

## Goals

- Convert root-level `--` flag commands to Cobra subcommands with single-letter root flag aliases
- Scope filter flags to the `list` subcommand so short flags (`-t`, `-s`, etc.) can be reused by other subcommands without conflicts
- Enable POSIX-style flag stacking within subcommands (e.g., `gws list -tsv`)
- Maintain `gws <repo>` positional navigation with tab completion that shows both subcommands and matching repo names
- Provide a deprecation layer that keeps old `--flag` forms working with warnings for one release cycle
- Update shell integration templates to route subcommands correctly

## User Stories

- **As a CLI user**, I want to use intuitive subcommands (`gws list`, `gws init`) so that the interface feels natural and discoverable.
- **As a power user**, I want single-letter shortcuts (`gws -l`, `gws -i`) so that I can work quickly without typing full command names.
- **As a user filtering repositories**, I want to combine short flags within the list command (`gws list -tsv`) so that I can compose filters efficiently.
- **As a user navigating to repos**, I want `gws <repo>` to continue working and tab completion to suggest both commands and matching repo names.
- **As an existing user**, I want the old `--list`, `--init` style to keep working (with deprecation warnings) so that my muscle memory and scripts don't break immediately.

## Demoable Units of Work

### Unit 1: Core Subcommand Structure

**Purpose:** Replace root-level `--` flag commands with Cobra subcommands and single-letter root aliases. This establishes the new CLI skeleton that all subsequent work (tag refactor, user refactor) builds on.

**Functional Requirements:**
- The system shall register the following as Cobra subcommands on the root command: `list`, `init`, `add`, `refresh`, `print-workspace`
- Each subcommand shall have a single-letter alias registered as a root-level boolean flag that triggers the subcommand's logic (e.g., `-l` triggers `list`, `-i` triggers `init`, `-a` triggers `add`, `-r` triggers `refresh`, `-w` triggers `print-workspace`)
- The system shall enforce mutual exclusivity: only one base command (subcommand or root alias flag) may be active at a time
- The `init` subcommand shall accept an optional positional argument for the directory path (defaults to `.`)
- The `add` subcommand shall accept an optional positional argument for the repo path (defaults to `.`) and a `--recursive`/`-v` flag
- The `refresh` subcommand shall take no arguments and refresh all repo metadata
- The `print-workspace` subcommand shall take no arguments and print the workspace root path
- The root command with no arguments shall display workspace info (current default behavior)
- The `shell-init` and `completion` subcommands shall remain unchanged
- The `--version`/`-V` flag shall remain unchanged

**Proof Artifacts:**
- CLI: `gws list` produces the same output as current `gws --list` demonstrates subcommand works
- CLI: `gws -l` produces the same output as `gws list` demonstrates root alias works
- CLI: `gws init` in a new directory creates workspace demonstrates init subcommand works
- Test: All existing command tests updated and passing demonstrates no regressions

### Unit 2: List Subcommand with Scoped Filters

**Purpose:** Move all filter flags from the root command to the `list` subcommand, enabling flag reuse across commands and POSIX stacking within `list`.

**Functional Requirements:**
- The `list` subcommand shall register the following flags scoped to itself:
  - `--tag`/`-t` (string slice): filter by tag(s), AND logic, supports wildcards
  - `--type`/`-y` (string): filter by repo type
  - `--name`/`-n` (string): filter by repo name, supports wildcards
  - `--path`/`-p` (string): filter by repo path, supports wildcards
  - `--output`/`-o` (string): output format (`table` or `json`)
  - `--status`/`-s` (bool): show git status columns
  - `--show-user`/`-u` (bool): show user/email/sign columns
- The short flags `-t`, `-s`, `-u`, `-n`, `-p`, `-y`, `-o` shall be stackable in POSIX style when boolean (e.g., `gws list -su` shows both status and user columns)
- All existing filter behavior (wildcard patterns, AND logic, partial matching) shall remain unchanged
- The `list` subcommand with no flags shall list all repositories in table format (current default)

**Proof Artifacts:**
- CLI: `gws list -t work -y github -su` produces filtered output with status and user columns demonstrates scoped filters and stacking work
- CLI: `gws list --tag "wo*" --status` produces wildcard-filtered output demonstrates long-form filters work
- Test: Filter flags on `list` subcommand resolve correctly demonstrates flag registration

### Unit 3: Navigation and Tab Completion

**Purpose:** Maintain the `gws <repo>` navigation experience and update tab completion to suggest both subcommands and repo names.

**Functional Requirements:**
- The system shall continue to support `gws <repo-name>` for positional navigation (unchanged behavior)
- The `--go`/`-g` flag shall continue to work as an explicit navigation flag
- Tab completion shall suggest both subcommand names and matching repo names when the user presses tab after `gws `
- Cobra's subcommand matching shall take priority: if the input exactly matches a subcommand, it routes to that subcommand
- Partial matches shall include both subcommand and repo name suggestions (e.g., typing `gws us<tab>` suggests both the `user` command and any repos starting with "us")
- The shell integration function (`gws` shell function in zsh/bash) shall be updated to route subcommand words to the binary directly
- The shell function shall detect subcommands by name (not just the `-*` prefix pattern) for proper routing

**Proof Artifacts:**
- CLI: `gws my-repo` navigates to the repo (via shell function `cd`) demonstrates positional navigation unchanged
- CLI: Tab completion after `gws ` shows both commands and repo names demonstrates dual completion
- Test: Shell integration templates updated and generate correct function demonstrates shell routing

### Unit 4: Deprecation Layer

**Purpose:** Keep old `--` flag forms working with deprecation warnings so existing users and scripts aren't broken immediately.

**Functional Requirements:**
- The system shall register hidden root-level flags for all old forms: `--list`, `--init`, `--add`, `--refresh`, `--print-workspace`, `--go`
- When an old `--` flag is used, the system shall print a deprecation warning to stderr: `Warning: --<flag> is deprecated, use '<subcommand>' instead`
- The old flags shall invoke the same logic as the new subcommands (no behavior difference)
- The old filter flags (`--type`, `--tag`, `--name`, `--path`, `--output`, `--status`, `--show-user`) shall continue to work when combined with `--list`, with deprecation warnings
- Hidden flags shall not appear in `--help` output
- The deprecation layer shall be removable in a future release (clean separation in code)
- The system shall NOT register deprecated forms for `--add-tag`, `--remove-tag`, `--user`, `--update`, `--delete` — those will be handled in specs 10 and 11

**Proof Artifacts:**
- CLI: `gws --list` produces correct output plus a deprecation warning to stderr demonstrates backward compatibility
- CLI: `gws --list --type github` produces filtered output plus deprecation warnings demonstrates filter compat
- CLI: `gws --help` does NOT show deprecated flags demonstrates clean help output
- Test: Deprecated flags route to correct subcommand logic demonstrates deprecation layer works

## Non-Goals (Out of Scope)

1. **Tag command refactor**: The `tag`/`--add-tag`/`--remove-tag` restructuring is deferred to spec 10
2. **User command refactor**: The `user` subcommand and `--user`/`--update`/`--delete` restructuring is deferred to spec 11
3. **New features**: No new functionality beyond the structural refactor
4. **Internal package changes**: No modifications to `internal/config`, `internal/filter`, `internal/git`, `internal/discovery`, `internal/classifier`, or `internal/user`
5. **Removing the deprecation layer**: This spec adds it; removal is a future task after the deprecation period

## Design Considerations

No specific design requirements identified. The CLI output format (tables, JSON) remains unchanged. Help text for each subcommand should clearly document available flags and usage examples.

## Repository Standards

- Follow existing Go code patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests
- Use Cobra/pflag conventions for subcommand and flag registration
- Maintain the existing `cmd/git-workspace/` and `internal/` package separation
- Follow conventional commits for version history (`feat!:` for breaking change = major version bump)
- Use `make test` and `make lint` for validation
- Test files co-located with source, using `t.TempDir()` and `t.Setenv()` for isolation
- TTY simulation via injectable `isTerminalFunc` for interactive behavior tests

## Technical Considerations

- **Cobra subcommand + root flag aliases**: Each base command is a Cobra subcommand (`rootCmd.AddCommand(listCmd)`). The single-letter alias (e.g., `-l`) is a hidden boolean flag on root that, when set, delegates to the subcommand's `RunE`. This allows both `gws list` and `gws -l` to work.
- **Flag scoping**: Since filter flags are registered on the `list` subcommand (not root), the same short flag letters (e.g., `-t`, `-s`) can be reused by other subcommands (tag, user) in later specs without conflicts.
- **Mutual exclusivity**: The root `RunE` must check that at most one alias flag is set and that no subcommand was also invoked. Cobra's `TraverseChildren` may need adjustment.
- **Shell function update**: The current shell function routes based on `case "$1" in -*|__*|completion|shell-init|help)`. This must be updated to also match subcommand names (`list|init|add|refresh|print-workspace|tag|user|completion|shell-init|help`), or changed to a simpler pattern (route everything to the binary, only `cd` for unrecognized args).
- **Tab completion**: Cobra's `ValidArgsFunction` on the root command currently returns repo names. With subcommands registered, Cobra will automatically include subcommand names in completions. The `ValidArgsFunction` should continue to return repo names so both appear. Testing is needed to confirm Cobra merges these correctly.
- **Deprecation code isolation**: Deprecated flag registration and warning logic should live in a single file (e.g., `deprecated.go`) for easy removal later.

## Security Considerations

No specific security considerations identified. No changes to credential handling, file permissions, or external service communication.

## Success Metrics

1. **All existing tests pass** with updated assertions reflecting the new subcommand pattern
2. **`make ci` passes** (vet, lint, test-race) with no regressions
3. **Tab completion** shows both subcommands and repo names in zsh and bash
4. **Deprecation warnings** appear on stderr when old `--flag` forms are used
5. **Old `--flag` forms** produce identical output to new subcommands (minus the warning)

## Open Questions

1. Should the major version bump to v3.0 happen with this spec (core structure), or wait until all 3 specs are complete (after tag and user refactors)?
2. Should the `--quiet`/`-q` flag (currently used with navigation) remain on root or move to specific subcommands?
