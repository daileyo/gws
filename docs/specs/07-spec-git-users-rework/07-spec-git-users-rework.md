# 07-spec-git-users-rework

## Introduction/Overview

The git user functionality in `gws` needs to be reworked to improve user visibility, smarter configuration detection, and CLI consistency. This spec covers marking the default (global gitconfig) user in list output, evaluating includeIf and local git configurations during init/add/refresh operations, updating CLI flag shorthands and names, and adding unit test coverage for all changes.

## Goals

- Mark repos using the global gitconfig user with a visual indicator (`*`) in list output so users can quickly identify which repos use the default vs. overridden identity
- Evaluate `includeIf` gitdir configurations during init, add, and refresh to detect effective user identities and auto-link repos to matching stored profiles
- Detect local `.git/config` author info during init, add, and refresh for display in `--show-user` output
- Update `--remove-tag` shorthand from `-u` to a more intuitive shorthand (keeping `-d` with `--add-tag`)
- Rename `--user` flag to `--show-user` for clarity
- Add unit tests covering all new and modified functionality

## User Stories

**As a developer managing multiple git identities**, I want to see which repos are using my default (global) git user so that I can quickly identify repos that may need a specific identity assigned.

**As a developer with includeIf gitdir configurations**, I want `gws` to automatically detect and display the effective user for each repo based on my gitconfig setup so that I don't have to manually assign profiles to repos already configured via includeIf.

**As a developer**, I want repos with local `.git/config` author overrides to be reflected in the user display so that I can see the actual effective identity for each repo.

**As a CLI user**, I want consistent, intuitive flag names and shorthands so that the tool is easy to use and remember.

## Demoable Units of Work

### Unit 1: Default User Indicator in List Output

**Purpose:** Allow users to visually identify repos using the global/default gitconfig user when viewing repository lists.

**Functional Requirements:**
- The system shall read the global `~/.gitconfig` user.name and user.email (excluding includeIf and local overrides) to determine the "default" user identity
- The system shall append an asterisk marker (`*`) after the user name in `--show-user` list output for repos whose effective user matches the global default (e.g., `John Doe *`)
- The system shall not display the marker for repos using a local or includeIf-sourced user identity
- The system shall handle the case where no global user is configured (no marker shown for any repo)

**Proof Artifacts:**
- CLI: `gws -l --show-user` output showing `*` marker on repos using global default user
- Test: Unit tests verifying default user detection and marker logic pass

### Unit 2: IncludeIf Evaluation During Init/Add/Refresh

**Purpose:** Automatically detect effective user identities from gitconfig includeIf directives and link repos to matching stored profiles during workspace operations.

**Functional Requirements:**
- The system shall parse `~/.gitconfig` for `[includeIf "gitdir:..."]` directives during init, add, and refresh operations
- The system shall evaluate each repo's path against includeIf gitdir patterns to determine if an alternate user config applies
- The system shall store the effective user/email with `UserSource: includeif` when a repo matches an includeIf directive
- The system shall attempt to match the detected includeIf user to a stored profile by comparing email and git name, and auto-link the repo to that profile if a match is found
- The system shall handle nested includes and multiple includeIf directives

**Proof Artifacts:**
- CLI: `gws -l --show-user` output showing includeIf-resolved users for repos under gitdir paths
- CLI: `gws --refresh` output showing repos updated with includeIf-detected user info
- Test: Unit tests for includeIf path matching, user resolution, and profile auto-linking pass

### Unit 3: Local Config Detection for Display

**Purpose:** Detect and display local `.git/config` author information during init, add, and refresh so users see the actual effective identity per repo.

**Functional Requirements:**
- The system shall read `user.name` and `user.email` from each repo's `.git/config` during init, add, and refresh operations
- The system shall display local config values in `--show-user` output when present, showing them as the effective user for that repo
- The system shall not persist local config values into the gws config file (detect for display only)
- The system shall prioritize local config over includeIf over global config when determining the displayed user (matching git's own precedence)

**Proof Artifacts:**
- CLI: `gws -l --show-user` output showing locally configured users with `(local)` source marker
- Test: Unit tests for local config detection and display priority pass

### Unit 4: CLI Flag Updates

**Purpose:** Improve CLI consistency by updating flag shorthands and names.

**Functional Requirements:**
- The system shall change the `--remove-tag` shorthand from `-u` to `-x` (mnemonic: "x out" / delete)
- The system shall rename the `--user` filter flag to `--show-user` with no shorthand
- The system shall update all help text, error messages, and documentation referencing the old flag names/shorthands
- The system shall update mutual exclusivity error messages to reflect the new flag names

**Proof Artifacts:**
- CLI: `gws --help` output showing updated flag names and shorthands
- CLI: `gws -x repo-name tag-name` successfully removes a tag
- Test: Unit tests verifying flag registration and shorthand mappings pass

## Non-Goals (Out of Scope)

1. **Adding a `gws user default` command** - No new subcommand to set a default profile; the default is simply the global gitconfig user
2. **Persisting local config into gws config** - Local `.git/config` values are detected for display only, not stored
3. **Changing `--add-tag` shorthand** - The `-d` shorthand stays with `--add-tag`
4. **Adding new user subcommands** - No new `user` subcommands beyond the existing set
5. **Integration or end-to-end tests** - Only unit tests are in scope

## Design Considerations

No specific design requirements identified. The `*` marker for default user follows the existing pattern of inline indicators (e.g., `(local)` for user source, `✓` for signing).

## Repository Standards

- **Language/Build:** Go 1.23.0 with `make build`, `make test`, `make lint`
- **Testing:** Table-driven tests with `t.Run()` subtests, `t.TempDir()` for filesystem tests
- **Package structure:** CLI commands in `cmd/git-workspace/`, business logic in `internal/`
- **Error handling:** Explicit error returns with `fmt.Errorf` wrapping
- **Naming:** Test files use `*_test.go` suffix, test functions use `TestFunctionName` pattern

## Technical Considerations

- The `git.GetUserConfig()` function already handles reading effective config; it may need enhancement to distinguish global vs. includeIf sources more precisely
- The `internal/user/gitconfig.go` already has `ParseGitconfig()` and `ExtractIncludeIfs()` which can be leveraged for includeIf evaluation
- The list display logic in `cmd/git-workspace/list.go` already computes `repoUserInfo` with source markers; the default indicator extends this pattern
- Flag changes in `cmd/git-workspace/main.go` require updating the `init()` function, help template functions, and error message strings
- Local config detection should use git's own precedence: local > includeIf > global

## Security Considerations

No specific security considerations identified. All operations read existing git configurations and do not handle credentials or sensitive data.

## Success Metrics

1. **Default user visibility**: Repos using global gitconfig identity are clearly marked with `*` in list output
2. **IncludeIf accuracy**: Repos under includeIf gitdir paths show the correct effective user after init/add/refresh
3. **Local config accuracy**: Repos with local `.git/config` user overrides display the correct effective identity
4. **CLI consistency**: All flag names, shorthands, help text, and error messages reflect the updated conventions
5. **Test coverage**: All new and modified functions have unit tests with table-driven patterns

## Open Questions

1. Should the `--remove-tag` shorthand be `-x` (suggested above) or would you prefer `-X`, `-D`, or `-R`?
-x works for now.
