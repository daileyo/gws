# 12-spec-list-remote

## Introduction/Overview

Add a `--remote`/`-r` flag to `gws list` that displays the origin remote URL for each repository. When a repository has additional remotes beyond origin, an asterisk prefix indicates this. When a repository has no origin but has other remotes, only an asterisk is displayed. This gives users quick visibility into where their repositories are hosted without leaving the list view.

## Goals

- Display the origin remote URL in list output when `--remote`/`-r` is specified
- Indicate the presence of non-origin remotes with an asterisk prefix
- Handle the no-origin case gracefully with a standalone asterisk
- Support both table and JSON output formats
- Follow existing flag patterns (`--status`, `--show-user`) for consistency

## User Stories

- **As a developer managing multiple repositories**, I want to see which remote URLs my repos point to so that I can quickly identify hosting providers and verify correct remote configurations.
- **As a developer with complex remote setups**, I want to know at a glance which repos have remotes beyond origin so that I can investigate further when needed.

## Demoable Units of Work

### Unit 1: Add `--remote` Flag and Table Output

**Purpose:** Enable users to see origin remote URLs in the standard table output with the new flag.

**Functional Requirements:**
- The system shall add a `--remote`/`-r` boolean flag to the `list` command
- The system shall display a REMOTE column as the last column in table output when `--remote` is specified
- The system shall display the stored `remote_url` from config for each repository
- The system shall display an empty value when a repository has no `remote_url` in config and no remotes detected
- The flag shall be stackable with existing flags (`--status`, `--show-user`, e.g., `gws list -rsu`)
- The REMOTE column width shall be dynamically calculated like other columns

**Proof Artifacts:**
- CLI: `gws list --remote` output shows REMOTE column with URLs as the last column
- CLI: `gws list -r` short flag works identically
- CLI: `gws list -rsu` demonstrates flag stacking with status and user columns

### Unit 2: Asterisk Indicator for Multiple Remotes

**Purpose:** Provide visual indication when repositories have remotes beyond origin, helping users identify complex remote configurations.

**Functional Requirements:**
- The system shall perform live git inspection of each repository's remotes when `--remote` is specified
- When a repository has origin plus one or more additional remotes, the system shall prefix the origin URL with `* ` (asterisk space)
- When a repository has no origin remote but has one or more other remotes, the system shall display just `*`
- When a repository has only origin (no additional remotes), the system shall display the origin URL without any prefix
- When a repository path is inaccessible or not a valid git repo, the system shall display the stored `remote_url` from config without an asterisk (graceful fallback)

**Proof Artifacts:**
- CLI: Repository with only origin shows URL without asterisk
- CLI: Repository with origin + upstream shows `* https://github.com/user/repo.git`
- CLI: Repository with no origin but has other remotes shows `*`

### Unit 3: JSON Output Support

**Purpose:** Ensure the remote information is available in JSON output for scripting and automation.

**Functional Requirements:**
- When `--remote` and `--output json` are both specified, the system shall include `remote_url` (string) in the JSON output for each repository
- When `--remote` and `--output json` are both specified, the system shall include `has_multiple_remotes` (boolean) in the JSON output for each repository
- The `remote_url` field shall contain the origin URL (or empty string if no origin)
- The `has_multiple_remotes` field shall be `true` when non-origin remotes exist, `false` otherwise

**Proof Artifacts:**
- CLI: `gws list -r -o json` output includes `remote_url` and `has_multiple_remotes` fields
- CLI: JSON output for a repo with multiple remotes shows `"has_multiple_remotes": true`

### Unit 4: Tests

**Purpose:** Ensure the feature is covered by unit tests following existing test patterns.

**Functional Requirements:**
- Tests shall verify the `--remote`/`-r` flag is registered on the list command
- Tests shall verify flag stacking works with `-r` combined with `-s` and `-u`
- Tests shall verify remote display logic: origin-only, origin+others, no-origin+others, no remotes

**Proof Artifacts:**
- Test: `go test ./cmd/git-workspace/ -run TestRemote` passes
- Test: All existing tests continue to pass (`go test ./...`)

## Non-Goals (Out of Scope)

1. **Displaying non-origin remote URLs**: Only the origin URL is shown; users must investigate other remotes manually
2. **Remote URL editing or management**: This is read-only display; no commands to change remotes
3. **Caching remote information**: The multi-remote detection is done live each time; no caching layer is added
4. **Filtering by remote URL**: No `--remote-url` filter flag is being added

## Design Considerations

No specific design requirements identified. The REMOTE column follows the existing dynamic-width table pattern used by other columns (NAME, TYPE, TAGS, PATH, etc.).

**Table output example:**
```
NAME          TYPE      TAGS       PATH                              REMOTE
----          ----      ----       ----                              ------
my-project    github    personal   /home/user/projects/my-project    https://github.com/user/my-project.git
work-api      gitlab    work       /home/user/work/work-api          * https://gitlab.com/team/work-api.git
local-only    unknown              /home/user/local-only             *
bare-repo     unknown              /home/user/bare-repo
```

## Repository Standards

- Follow Cobra flag registration pattern: `listCmd.Flags().BoolVarP(&showRemote, "remote", "r", false, "...")`
- Add `ShowRemote` field to `ListOptions` struct
- Follow existing column rendering pattern in `displayTable`
- Follow existing test patterns in `list_test.go`
- Use `go-git/go-git/v5` for live remote inspection (already a dependency)

## Technical Considerations

- **Remote URL source**: Use `RemoteURL` from config (`config.Repository.RemoteURL` field) — this is populated during `gws add` / discovery
- **Multi-remote detection**: Requires live git inspection via `go-git` to enumerate remotes. Open each repo with `git.PlainOpen()` and call `repo.Remotes()` to check for non-origin remotes
- **Performance**: Live git inspection adds I/O overhead similar to `--status`. Consider that users with many repos may notice a delay
- **Graceful fallback**: If a repo path doesn't exist or can't be opened, fall back to displaying the stored `remote_url` without asterisk indicators
- **Short flag `-r`**: Available on the `list` command (currently used by `--recursive` on `add`, but not on `list`)

## Security Considerations

No specific security considerations identified. Remote URLs are already stored in config and displayed elsewhere in the tool.

## Success Metrics

1. **Feature completeness**: `gws list -r` displays remote URLs with correct asterisk indicators for all repositories
2. **Backward compatibility**: No existing flags or output formats are broken
3. **Test coverage**: New flag and display logic covered by unit tests

## Open Questions

No open questions at this time.
