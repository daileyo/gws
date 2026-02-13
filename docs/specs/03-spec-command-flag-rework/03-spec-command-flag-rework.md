# 03-spec-command-flag-rework

## Introduction/Overview

This feature reworks the gws CLI from a subcommand-based pattern (`gws list`, `gws init`) to a flag-based pattern (`gws --list`, `gws --init`). All top-level commands become flags on the root command with `--long-name` and `-shorthand` forms, while positional arguments and sub-flags remain unchanged. Additionally, all filter values gain wildcard pattern support (e.g., `wi*` matches anything starting with `wi`).

## Goals

- Convert all existing subcommands to root-level flags with consistent `--flag`/`-f` naming
- Add shorthand flags for all commands and filter options
- Add wildcard pattern matching support to all filter values
- Remove the tag filter shorthand (`gws personal`) in favor of explicit flag usage
- Replace the `version` subcommand with Cobra's built-in `--version`/`-v` flag
- Maintain all existing business logic and internal packages unchanged

## User Stories

- **As a CLI user**, I want to use flag-style options (`gws -l`, `gws --list`) so that the CLI follows a consistent, predictable pattern for all operations.
- **As a power user**, I want shorthand flags for all commands and filters (`gws -l -y github -s`) so that I can work more efficiently.
- **As a user filtering repositories**, I want to use wildcard patterns in filter values (e.g., `gws -l --name "my*"`) so that I can match repositories flexibly without typing full names.
- **As a user managing tags**, I want clear, non-conflicting flag names (`--add-tag`, `--remove-tag`) so that tagging operations are unambiguous alongside tag-based filtering.

## Demoable Units of Work

### Unit 1: Root Command Flag Conversion

**Purpose:** Convert all subcommands to root-level flags and remove the old subcommand pattern. This establishes the new CLI interface.

**Functional Requirements:**
- The system shall accept `--list` / `-l` as a boolean flag to trigger repository listing
- The system shall accept `--init` / `-i` as a string flag taking a directory path to initialize a workspace
- The system shall accept `--add-tag` / `-t` as a flag to add a tag to repositories (takes repo and tag as positional args)
- The system shall accept `--remove-tag` / `-u` as a flag to remove a tag from repositories (takes repo and tag as positional args)
- The system shall accept `--refresh` / `-r` as a boolean flag to refresh repository metadata
- The system shall use Cobra's built-in `--version` / `-v` flag showing version string
- The system shall accept `--print-workspace` / `-w` (adding shorthand) to print workspace path
- The system shall remove all existing subcommands (`list`, `init`, `tag`, `untag`, `refresh`, `version`)
- The system shall remove the tag filter shorthand (bare `gws personal` no longer works)
- The system shall display workspace info when invoked with no flags (current default behavior)
- The root command help (`gws -h`) shall clearly document all available flags and their shorthands

**Proof Artifacts:**
- Test: All command tests updated and passing demonstrates flag-based invocation works
- Test: Old subcommand patterns no longer registered demonstrates clean break

### Unit 2: List Filter Shorthands

**Purpose:** Add shorthand flags to all list filter options for faster CLI usage.

**Functional Requirements:**
- The `--type` filter shall have shorthand `-y`
- The `--name` filter shall have shorthand `-n`
- The `--path` filter shall have shorthand `-p`
- The `--tag` filter shall have shorthand `-t` (scoped to list context)
- The `--output` / `-o` shorthand shall remain unchanged
- The `--status` / `-s` shorthand shall remain unchanged
- All filter flags shall only apply when `--list` / `-l` is present
- All filter flags shall remain combinable (AND logic)

**Proof Artifacts:**
- Test: Filter shorthand flags resolve correctly demonstrates shorthand registration
- Test: Combined filters with shorthands produce correct results demonstrates composability

### Unit 3: Wildcard Pattern Matching for Filters

**Purpose:** Enable wildcard pattern matching in filter values so users can flexibly match repositories.

**Functional Requirements:**
- The `--name` filter shall support wildcard patterns (e.g., `my*` matches `myproject`, `myapp`)
- The `--type` filter shall support wildcard patterns (e.g., `git*` matches `github`, `gitlab`)
- The `--path` filter shall support wildcard patterns (e.g., `/home/*/projects` matches paths)
- The `--tag` filter shall support wildcard patterns (e.g., `wo*` matches `work`, `workflows`)
- Wildcard matching shall use `*` (zero or more characters) and `?` (single character) as glob characters
- Filter values without wildcards shall continue to use current matching behavior (partial, case-insensitive)
- Wildcard matching shall be case-insensitive

**Proof Artifacts:**
- Test: Wildcard patterns match expected repositories demonstrates glob matching works
- Test: Non-wildcard values preserve existing behavior demonstrates backward compatibility
- Test: Wildcard patterns work across all filter types demonstrates consistent implementation

## Non-Goals (Out of Scope)

1. **Changes to internal packages**: No modifications to `internal/config`, `internal/discovery`, `internal/classifier`, or `internal/git` packages beyond what's needed for wildcard support in `internal/filter`
2. **New commands or features**: No new functionality beyond the flag pattern rework and wildcard support
3. **Shell completion**: Shell completion scripts are not part of this rework
4. **Interactive mode**: No interactive selection or prompting
5. **Backward compatibility layer**: No deprecation warnings or aliases for old subcommand patterns

## Design Considerations

No specific design requirements identified. The CLI output format (tables, JSON) remains unchanged.

## Repository Standards

- Follow existing Go code patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests
- Use Cobra/pflag conventions for flag registration
- Maintain the existing `cmd/gws/` and `internal/` package separation
- Follow conventional commits for version history
- Use `make test` and `make lint` for validation

## Technical Considerations

- **Cobra flag scoping**: Since all commands become root-level flags, filter flags (`--type`, `--tag`, etc.) must also be registered on the root command. Custom validation logic will be needed to ensure filter flags are only meaningful when `--list`/`-l` is active.
- **Mutual exclusivity**: The root command must enforce that only one command flag is active at a time (e.g., `--list` and `--init` cannot both be set). Cobra does not natively support mutual exclusion of flags, so this requires custom validation in the `RunE` handler.
- **Tag flag conflict resolution**: `--tag`/`-t` is used as a list filter flag. The tag command uses `--add-tag`/`-t`... wait, `-t` can't be shared. `--add-tag` needs a different shorthand. Use `-a` for `--add-tag` instead. Updated mapping:
  - `--add-tag` / `-a`: Add tag to repository
  - `--remove-tag` / `-u`: Remove tag from repository
  - `--tag` / `-t`: Filter by tag (list context)
- **Wildcard implementation**: Use Go's `filepath.Match` or `path.Match` for glob pattern matching, falling back to current partial match behavior when no wildcard characters are present.
- **Positional args routing**: When `--add-tag` or `--remove-tag` is set, positional args are `<repo> <tag>`. When `--init` is set, positional arg is `<directory>`. When `--list` is set with `--tag`, tag values are provided via the flag. Custom arg validation per active command flag is required.

## Security Considerations

No specific security considerations identified. No changes to credential handling, file permissions, or external service communication.

## Success Metrics

1. **All existing tests pass** with updated assertions reflecting the new flag pattern
2. **No regressions** in repository listing, filtering, tagging, or initialization functionality
3. **Wildcard patterns work** across all filter types with case-insensitive matching
4. **`make ci` passes** (vet, lint, test-race)

## Open Questions

No open questions at this time. All questions resolved:
- `-t` assigned to `--tag` (list filter), `-a` assigned to `--add-tag` (tag command) — confirmed acceptable.
- Both `*` (zero or more characters) and `?` (single character) wildcards will be supported.
