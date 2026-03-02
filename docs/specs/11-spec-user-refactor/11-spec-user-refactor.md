# 11-spec-user-refactor

## Introduction/Overview

The gws CLI currently has two overlapping interfaces for user management: a proper subcommand tree (`gws user add`, `gws user list`, etc.) and root-level flags (`--user`, `--update`, `--delete`). This duplication is confusing and inconsistent with the new subcommand structure from spec 09. This spec refactors the `user` command to use reusable short flags (`-a` add, `-d` remove, `-l` list, `-s` show) alongside the existing word-form sub-subcommands, removes the root-level `--user`/`--update`/`--delete` flags (with deprecation warnings), and keeps `assign` and `sync` as word-only sub-subcommands since they don't fit the CRUD pattern.

This is spec 3 of 3 in the command structure refactor. It depends on spec 09 (core structure) being complete.

## Goals

- Add reusable short flags (`-a`, `-d`, `-l`, `-s`) to the `user` subcommand as aliases for `add`, `remove`, `list`, `show`
- Remove the root-level `--user`, `--update`, `--delete`, `--all`, `--list-users` flags (with deprecation layer)
- Keep `assign` and `sync` as word-form sub-subcommands only (no short-flag equivalents)
- Show command help when `gws user` is invoked with no sub-operation
- Maintain all existing user management business logic unchanged

## User Stories

- **As a CLI user**, I want to manage users with both `gws user add` and `gws user -a` so that I can choose the style that fits my workflow.
- **As a power user**, I want to quickly list profiles with `gws user -l` instead of typing `gws user list`.
- **As an existing user**, I want `gws --user` and `gws --user --update` to keep working (with warnings) so that my workflows don't break immediately.
- **As a user**, I want `gws user` to show me all available operations so that the interface is discoverable.

## Demoable Units of Work

### Unit 1: User Short-Flag Aliases

**Purpose:** Add reusable CRUD short flags to the `user` subcommand as aliases for existing word-form sub-subcommands.

**Functional Requirements:**
- The `user` subcommand shall register `-a` as a boolean flag that triggers the `add` operation (equivalent to `user add`)
- The `user` subcommand shall register `-d` as a boolean flag that triggers the `remove` operation (equivalent to `user remove`)
- The `user` subcommand shall register `-l` as a boolean flag that triggers the `list` operation (equivalent to `user list`)
- The `user` subcommand shall register `-s` as a boolean flag that triggers the `show` operation (equivalent to `user show`)
- When `-a` is used, the profile name shall be a positional argument: `gws user -a <name> --email <email> [--name <name>] [--signing-key <key>] [--sign-commits]`
- When `-d` is used, the profile name shall be a positional argument: `gws user -d <name>`
- When `-s` is used, the profile name shall be a positional argument: `gws user -s <name>`
- When `-l` is used, no positional arguments are required: `gws user -l`
- The system shall enforce mutual exclusivity: only one of `-a`, `-d`, `-l`, `-s`, or a word-form sub-subcommand may be active at a time
- Short flags shall be stackable where it makes sense (though CRUD flags are mutually exclusive, future modifier flags like `-v` verbose could stack: `gws user -lv`)
- The `assign` sub-subcommand shall remain word-form only: `gws user assign <repo> <profile> [--use-subdirs] [--dry-run]`
- The `sync` sub-subcommand shall remain word-form only: `gws user sync`
- When `gws user` is invoked with no sub-operation and no flag, the system shall display the user command's help text
- All existing user management business logic (profile CRUD, assignment, sync) shall remain unchanged
- The `--verbose` flag shall remain available on the `user` subcommand for detailed output

**Proof Artifacts:**
- CLI: `gws user -a dev-profile --email dev@example.com` adds a profile demonstrates short-flag add
- CLI: `gws user -l` lists all profiles demonstrates short-flag list
- CLI: `gws user -s dev-profile` shows profile details demonstrates short-flag show
- CLI: `gws user -d dev-profile` removes a profile demonstrates short-flag remove
- CLI: `gws user add dev-profile --email dev@example.com` still works demonstrates word-form preserved
- CLI: `gws user assign my-repo dev-profile` still works demonstrates assign unchanged
- CLI: `gws user` displays help text demonstrates default behavior
- Test: User CRUD tests pass with both invocation styles demonstrates dual interface

### Unit 2: User Deprecation Layer

**Purpose:** Deprecate the root-level `--user`, `--update`, `--delete`, `--all`, `--list-users`, `--git-name`, `--git-email` flags with warnings.

**Functional Requirements:**
- The system shall register the following as hidden flags on root with deprecation warnings:
  - `--user` (bool): when used alone, delegates to `user list`; warns to use `user list` or `user -l`
  - `--user` + `--update` (bool): delegates to `user assign`; warns to use `user assign`
  - `--user` + `--delete` (bool): delegates to `user remove` logic (removing local git config); warns to use appropriate new form
  - `--list-users` (bool): delegates to `user list`; warns to use `user list` or `user -l`
  - `--all` (bool): modifier for `--delete`, warns about new form
  - `--git-name` (string): inline value for update, warns about new form
  - `--git-email` (string): inline value for update, warns about new form
- Each deprecated flag shall print a warning to stderr when used: `Warning: --<flag> is deprecated, use '<new form>' instead`
- Deprecated flags shall produce identical behavior to the new subcommand forms
- Deprecated flags shall not appear in `--help` output
- The deprecation code shall be co-located with the existing deprecation layer from specs 09 and 10 (in `deprecated.go`)

**Proof Artifacts:**
- CLI: `gws --user` lists profiles plus prints deprecation warning demonstrates backward compat
- CLI: `gws --list-users` lists profiles plus prints deprecation warning demonstrates backward compat
- CLI: `gws --user --update my-repo dev-profile` assigns profile plus prints deprecation warning demonstrates backward compat
- CLI: `gws --help` does NOT show deprecated user flags demonstrates clean help output
- Test: Deprecated user flags route correctly demonstrates deprecation layer

## Non-Goals (Out of Scope)

1. **New user operations**: No new functionality (e.g., user rename, user export)
2. **Changes to user assign/sync logic**: Business logic remains identical
3. **Profile auto-detection changes**: The gitconfig includeIf detection is unchanged
4. **Changes to the `--show-user`/`-u` filter on `list`**: That flag is scoped to the `list` subcommand (spec 09) and is unrelated to user management

## Design Considerations

No specific design requirements identified. User command output format remains unchanged.

## Repository Standards

- Follow existing Go code patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests
- Use Cobra/pflag conventions for subcommand and flag registration
- Maintain the existing `cmd/git-workspace/` and `internal/` package separation
- Follow conventional commits for version history
- Use `make test` and `make lint` for validation
- Test files co-located with source, using `t.TempDir()` and `t.Setenv()` for isolation

## Technical Considerations

- **Dual interface on existing subcommand**: Unlike `tag` (spec 10) which is being created fresh, the `user` subcommand already exists with word-form sub-subcommands. The short flags (`-a`, `-d`, `-l`, `-s`) are added to the existing `user` command and its `RunE` routes based on which flag is set. This must coexist with Cobra's sub-subcommand dispatch.
- **`-u` conflict avoidance**: The `list` subcommand uses `-u` for `--show-user`. The `user` subcommand does NOT need `-u` (update is handled via `assign` word-form). This avoids a cross-command flag letter collision in future refactors.
- **Deprecation complexity**: The root-level `--user` flag acts as a mode selector (combined with `--update` or `--delete`). The deprecation layer needs to inspect multiple flag combinations to determine which new command to suggest in the warning message.
- **Confirmation prompts**: `user remove` has an interactive confirmation if repos are using the profile. This behavior must be preserved regardless of whether invoked via `user remove <name>` or `user -d <name>`.

## Security Considerations

No specific security considerations identified. User profiles contain git configuration (name, email, signing key) but no secrets. Signing keys are references (key IDs), not private key material.

## Success Metrics

1. **All user tests pass** with both word-form and short-flag invocations
2. **`make ci` passes** with no regressions
3. **Deprecated root-level user flags** produce correct results with warnings
4. **`gws user`** displays useful help text showing all available operations
5. **`assign` and `sync`** continue to work unchanged

## Open Questions

No open questions at this time. All decisions resolved in spec 09 Q&A (user sub-flag mapping: Q4 = direct CRUD mapping with assign/sync as word-only; default behavior: Q5 = show help).
