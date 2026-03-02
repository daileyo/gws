# 08-spec-user-update-delete.md

## Introduction/Overview

This feature exposes git user configuration management through the gws flag-based CLI interface. Users will be able to set (update) and remove (delete) local git user configuration for repositories via `--user -u` (update) and `--user -d` (delete) root flags, with support for batch operations via `--tag` filtering. The existing `gws user assign` subcommand logic is already implemented but hidden from the CLI — this spec surfaces that functionality through the established flag-based interface pattern and adds delete capability.

## Goals

- Expose user assignment functionality as `--user -u` (update) flag on the root command, consistent with the existing flag-based CLI pattern
- Add user config deletion via `--user -d` (delete) flag, ensuring repos fall back to the default (global/includeIf) user after deletion
- Support batch update and delete operations via `--tag` filtering
- Support both named profiles and inline values (`--name`, `--email`) for updates, with inline values taking precedence
- Update `config.json` after all operations to keep stored state consistent

## User Stories

- **As a developer**, I want to set the local git user config for a repository via `gws --user -u my-repo work` so that I can quickly configure the correct identity for commits without manually editing `.git/config`.
- **As a developer**, I want to delete the local git user config from a repository via `gws --user -d my-repo` so that the repository falls back to my default (global or includeIf) user identity.
- **As a developer**, I want to batch update user config for all repos with a specific tag via `gws --user -u --tag work-projects work` so that I can configure multiple repositories at once.
- **As a developer**, I want to batch delete local user config from tagged repos via `gws --user -d --tag old-project` so that those repos revert to using my default identity.
- **As a developer**, I want to optionally provide inline `--name` and `--email` values instead of a profile name so that I can set user config without first creating a stored profile.

## Demoable Units of Work

### Unit 1: Single Repository User Update (`--user -u`)

**Purpose:** Allow users to set local git user config on a single repository using the flag-based CLI interface, leveraging the existing `user assign` logic.

**Functional Requirements:**
- The system shall accept `gws --user -u <repo-identifier> <profile-name>` to set local user config from a stored or auto-detected profile
- The system shall accept `gws --user -u <repo-identifier> --name "Name" --email "email@example.com"` to set local user config from inline values
- The system shall accept both a profile name AND inline values, with inline `--name`/`--email` overriding the profile's values
- The system shall match repositories by name (partial, case-insensitive) or path, using the existing `findRepositories` logic
- The system shall update all matching repositories when multiple repos match the identifier (interactive picker is deferred to a future spec)
- The system shall write `user.name` and `user.email` to the repository's local `.git/config` `[user]` section
- The system shall update `config.json` with the new User, Email, SigningEnabled, and UserSource (set to "local") after a successful update
- The system shall display moderate output by default: each repo name and what changed (e.g., `my-project: user.name "Old" → "New"`)
- The system shall support `--verbose` to show full before/after details for each repo
- The system shall support `--quiet` to suppress all output

**Proof Artifacts:**
- CLI: `gws --user -u <repo> <profile>` output demonstrates single repo update with moderate summary
- CLI: `git config --local --list` in the updated repo demonstrates actual `.git/config` changes
- CLI: `gws -l --show-user` demonstrates the repo reflects the new user info

### Unit 2: Single Repository User Delete (`--user -d`)

**Purpose:** Allow users to remove local git user config from a repository, causing it to fall back to the default (global or includeIf) user identity.

**Functional Requirements:**
- The system shall accept `gws --user -d <repo-identifier>` to remove local user config from matching repositories
- The system shall remove the `[user]` section (`name` and `email` keys) from the repository's local `.git/config` by default
- The system shall NOT remove signing config (`[user] signingkey`, `[commit] gpgsign`) by default
- The system shall accept an `--all` flag (`gws --user -d <repo> --all`) to also remove signing config (`signingkey` and `gpgsign`)
- The system shall match repositories using the same `findRepositories` logic as update (name/path, update all on multiple match)
- The system shall update `config.json` after deletion: re-detect the effective user config (which will now be global or includeIf) and store the updated values
- The system shall display moderate output by default, `--verbose` for detail, `--quiet` for silent operation
- After deletion, `gws -l --show-user` shall show the repository using its default (global/includeIf) user, not the previously-set local user

**Proof Artifacts:**
- CLI: `gws --user -d <repo>` output demonstrates local config removal
- CLI: `git config --local --list` confirms `[user]` section removed from `.git/config`
- CLI: `gws -l --show-user` demonstrates the repo now shows the default user (no "(local)" indicator)

### Unit 3: Batch Update and Delete via Tags (`--tag`)

**Purpose:** Allow users to update or delete local git user config for multiple repositories at once, filtered by tag.

**Functional Requirements:**
- The system shall accept `gws --user -u --tag <tag> <profile-name>` to update all repos matching the tag filter
- The system shall accept `gws --user -u --tag <tag> --name "Name" --email "email@example.com"` for inline batch update
- The system shall accept `gws --user -d --tag <tag>` to delete local user config from all repos matching the tag filter
- The system shall accept `gws --user -d --tag <tag> --all` to also remove signing config from all matching repos
- The system shall support multiple `--tag` values with AND logic (consistent with existing `--list --tag` behavior)
- The system shall overwrite existing local config on matching repos without confirmation (batch operations are intentional)
- The system shall display moderate output by default: list each affected repo and changes made
- The system shall support `--verbose` and `--quiet` flags for batch operations
- The system shall update `config.json` for all affected repositories after the batch operation completes

**Proof Artifacts:**
- CLI: `gws --user -u --tag <tag> <profile>` output demonstrates batch update with summary of all changed repos
- CLI: `gws --user -d --tag <tag>` output demonstrates batch delete with summary
- CLI: `gws -l --show-user --tag <tag>` demonstrates all tagged repos reflect the changes

## Non-Goals (Out of Scope)

1. **Interactive multi-select picker**: When multiple repos match an identifier, all are updated/deleted. An interactive spacebar-selectable list is desired but deferred to a future spec (pending library decision)
2. **`--dry-run` support**: The existing `user assign` subcommand has `--dry-run`; this flag-based interface does not include it in this spec
3. **`--use-subdirs` support**: The existing `user assign` subcommand supports moving repos to profile subdirectories; this is not included in the flag-based interface
4. **Removing or deprecating the `gws user` subcommand tree**: The existing subcommands (`user list`, `user add`, `user show`, `user remove`, `user assign`, `user sync`) remain as-is
5. **Signing-only operations**: This spec does not support setting signing config independently of user config

## Design Considerations

No specific design requirements identified. The CLI follows the existing flag-based interface pattern with grouped flag sections in `--help` output.

## Repository Standards

- **CLI Pattern**: Flag-based interface on root command (consistent with `--add-tag`, `--remove-tag`, `--refresh`, etc.)
- **Flag Grouping**: New flags should be added to the appropriate section in the custom usage template (Commands section)
- **Go Conventions**: Follow existing patterns — Cobra for CLI, `internal/` packages for business logic, table-driven tests with `t.TempDir()`
- **Error Handling**: Wrapped errors with context (`fmt.Errorf("context: %w", err)`), early returns for validation
- **Commit Convention**: Conventional Commits format (`feat:`, `fix:`, etc.)
- **Case-Insensitive Matching**: Use `strings.EqualFold()` for profile and repo name comparisons

## Technical Considerations

- **Reuse `user.AssignLocal()`**: The existing function in `internal/user/assign.go` handles writing `user.name` and `user.email` to `.git/config`. The update operation should reuse this.
- **New `DeleteLocal()` function needed**: A new function in `internal/user/assign.go` to remove `[user]` section keys (and optionally signing config) from `.git/config`.
- **Config.json re-detection after delete**: After removing local config, call `git.GetNonLocalUserConfig()` (or `git.GetUserConfig()`) to detect the effective fallback user and update the stored Repository fields.
- **Flag mutual exclusivity**: `--user -u` and `--user -d` must be mutually exclusive. The `--user` flag may need to be a boolean that enables the user-operation mode, with `-u` and `-d` as the operation selectors.
- **`--tag` reuse**: The existing `--tag` flag (string slice) is currently restricted to `--list`. This spec extends its use to `--user` operations, which requires updating the validation logic in `main.go`.
- **Inline value flags**: `--name` and `--email` flags need to be added to the root command (scoped to `--user -u` operations).

## Security Considerations

No specific security considerations identified. Operations modify local `.git/config` files and the gws `config.json` (both user-owned, local files). No credentials or sensitive data are handled.

## Success Metrics

1. **Single update works**: `gws --user -u <repo> <profile>` sets local git config and updates config.json
2. **Single delete works**: `gws --user -d <repo>` removes local user config and repo falls back to default
3. **Batch update works**: `gws --user -u --tag <tag> <profile>` updates all matching repos
4. **Batch delete works**: `gws --user -d --tag <tag>` deletes local config from all matching repos
5. **`--all` flag works**: `gws --user -d <repo> --all` also removes signing config
6. **Output modes work**: Default (moderate), `--verbose`, and `--quiet` produce expected output levels
7. **Config.json stays consistent**: After all operations, `gws -l --show-user` reflects the correct state

## Open Questions

1. Interactive multi-select picker library choice — deferred, will be addressed in a future spec
2. Should `--user -u` also support `--signing-key` and `--sign-commits` inline flags (like `user add` does), or only `--name` and `--email`?
