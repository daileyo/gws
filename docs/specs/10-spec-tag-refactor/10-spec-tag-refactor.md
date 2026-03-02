# 10-spec-tag-refactor

## Introduction/Overview

The gws CLI currently handles tagging via root-level flags `--add-tag`/`-d` and `--remove-tag`/`-x`, which are non-intuitive and inconsistent with the new subcommand structure established in spec 09. This spec refactors tagging into a proper `tag` subcommand with reusable sub-flags (`-a` for add, `-d` for delete) and word-form aliases (`tag add`, `tag remove`). The `-t` short flag on root serves as the `tag` command alias, while `-t` on the `list` subcommand remains the tag filter — no conflict since they're scoped to different Cobra commands.

This is spec 2 of 3 in the command structure refactor. It depends on spec 09 (core structure) being complete.

## Goals

- Consolidate tag operations under a single `tag` subcommand with `add`/`remove` sub-operations
- Provide both word-form (`tag add`) and reusable short-flag (`tag -a`) interfaces for sub-operations
- Register `-t` as the root-level alias for the `tag` command (no conflict with `list -t` filter)
- Deprecate old `--add-tag`/`--remove-tag` root flags with warnings for one release cycle
- Show command help when `gws tag` is invoked with no sub-operation

## User Stories

- **As a CLI user**, I want to manage tags with `gws tag add <repo> <tag>` so that the command reads naturally and is discoverable.
- **As a power user**, I want to use `gws tag -a <repo> <tag>` or `gws -t -a <repo> <tag>` so that I can work quickly with short flags.
- **As a user**, I want `gws tag` to show me available sub-operations so that I can learn the command without reading docs.
- **As an existing user**, I want `gws --add-tag` to keep working (with a warning) so that my workflows don't break immediately.

## Demoable Units of Work

### Unit 1: Tag Subcommand with Add/Remove

**Purpose:** Create the `tag` subcommand with `add` and `remove` as both word-form sub-subcommands and reusable short flags.

**Functional Requirements:**
- The system shall register `tag` as a Cobra subcommand on the root command
- The system shall register `-t` as a hidden boolean flag on the root command that delegates to the `tag` subcommand
- The `tag` subcommand shall register `add` as a sub-subcommand accepting `<repo> <tag>` positional arguments
- The `tag` subcommand shall register `remove` as a sub-subcommand accepting `<repo> <tag>` positional arguments
- The `tag` subcommand shall register `-a` as a boolean flag that triggers the `add` operation (equivalent to `tag add`)
- The `tag` subcommand shall register `-d` as a boolean flag that triggers the `remove` operation (equivalent to `tag remove`)
- When `-a` or `-d` is used, the `<repo>` and `<tag>` shall be provided as positional arguments: `gws tag -a <repo> <tag>`
- The system shall enforce mutual exclusivity: `-a` and `-d` cannot both be set, and neither can be combined with word-form sub-subcommands
- When `gws tag` is invoked with no sub-operation and no flag, the system shall display the tag command's help text
- Tab completion for the `<repo>` argument shall suggest repo names from the workspace config
- Tab completion for the `<tag>` argument (in `remove`) shall suggest existing tags for the specified repo
- All existing tag business logic (adding to config, removing from config, validation) shall remain unchanged

**Proof Artifacts:**
- CLI: `gws tag add my-repo work` adds the "work" tag demonstrates word-form add works
- CLI: `gws tag -a my-repo work` adds the "work" tag demonstrates short-flag add works
- CLI: `gws tag remove my-repo work` removes the "work" tag demonstrates word-form remove works
- CLI: `gws tag -d my-repo work` removes the "work" tag demonstrates short-flag remove works
- CLI: `gws tag` displays help text demonstrates default behavior
- Test: Tag add/remove tests pass with both invocation styles demonstrates dual interface

### Unit 2: Tag Deprecation Layer

**Purpose:** Keep old `--add-tag` and `--remove-tag` root flags working with deprecation warnings.

**Functional Requirements:**
- The system shall register `--add-tag` as a hidden boolean flag on root that delegates to `tag add` logic
- The system shall register `--remove-tag` as a hidden boolean flag on root that delegates to `tag remove` logic
- When `--add-tag` is used, the system shall print to stderr: `Warning: --add-tag is deprecated, use 'tag add' or 'tag -a' instead`
- When `--remove-tag` is used, the system shall print to stderr: `Warning: --remove-tag is deprecated, use 'tag remove' or 'tag -d' instead`
- Deprecated flags shall not appear in `--help` output
- The deprecated flags shall produce identical results to the new subcommand forms
- The deprecation code shall be co-located with the spec 09 deprecation layer (e.g., in `deprecated.go`)

**Proof Artifacts:**
- CLI: `gws --add-tag my-repo work` adds the tag plus prints deprecation warning demonstrates backward compat
- CLI: `gws --remove-tag my-repo work` removes the tag plus prints deprecation warning demonstrates backward compat
- CLI: `gws --help` does NOT show `--add-tag` or `--remove-tag` demonstrates clean help output
- Test: Deprecated tag flags route correctly demonstrates deprecation layer

## Non-Goals (Out of Scope)

1. **Tag listing/querying**: This spec does not add a `tag list` or `tag -l` operation — tags are viewed via `gws list --tag` or in the list output
2. **Bulk tag operations**: No support for tagging multiple repos at once
3. **Tag renaming**: No rename operation
4. **Changes to the list `--tag`/`-t` filter**: The filter behavior on the `list` subcommand is unchanged (established in spec 09)

## Design Considerations

No specific design requirements identified. Tag operations produce the same output format as today.

## Repository Standards

- Follow existing Go code patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests
- Use Cobra/pflag conventions for subcommand and flag registration
- Maintain the existing `cmd/git-workspace/` and `internal/` package separation
- Follow conventional commits for version history
- Use `make test` and `make lint` for validation
- Test files co-located with source, using `t.TempDir()` and `t.Setenv()` for isolation

## Technical Considerations

- **Dual interface pattern**: The `tag` subcommand needs to handle both `tag add <repo> <tag>` (sub-subcommand) and `tag -a <repo> <tag>` (flag + positional args). The `RunE` on the `tag` command should check if `-a` or `-d` flags are set and route accordingly, while `add`/`remove` sub-subcommands have their own `RunE`.
- **Root alias `-t`**: Registered as a hidden boolean flag on root. When set, root's `RunE` delegates to the `tag` command's execution with remaining args passed through.
- **No `-t` conflict**: Since `-t` on root and `-t` on `list` are on different Cobra commands, pflag treats them as independent flag sets. No collision.
- **Existing tag logic**: The `runAddTag()` and `runRemoveTag()` functions in the current codebase contain the business logic. These should be extracted/reused by the new subcommand handlers.

## Security Considerations

No specific security considerations identified.

## Success Metrics

1. **All tag tests pass** with both word-form and short-flag invocations
2. **`make ci` passes** with no regressions
3. **Deprecated `--add-tag`/`--remove-tag`** produce correct results with warnings
4. **`gws tag`** displays useful help text

## Open Questions

No open questions at this time. All decisions resolved in spec 09 Q&A (tag sub-flag pattern: Q3 = both word and flag forms; default behavior: Q5 = show help).
