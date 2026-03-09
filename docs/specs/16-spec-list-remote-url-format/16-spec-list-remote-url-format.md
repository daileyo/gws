# 16-spec-list-remote-url-format

## Introduction/Overview

The `gws list -r` command currently displays raw git remote URLs (e.g., `git@github.com:daileyo/gws.git` or `https://user@github.com/daileyo/gws.git`). These raw URLs contain protocol-specific syntax and sometimes user credentials, making them unsuitable for sharing with team members. This feature changes the default remote display to a clean, shareable HTTPS URL and adds a `--remote-raw`/`-R` flag to access the original raw remote when needed.

## Goals

- Display formatted HTTPS URLs by default when using `-r|--remote`, stripping user info and normalizing SSH to HTTPS
- Provide a `--remote-raw`/`-R` flag to show the original raw remote URL
- Maintain backward compatibility for scripts by applying the same behavior to JSON output
- Handle non-standard remote URLs gracefully by falling back to raw display

## User Stories

- **As a developer**, I want `gws list -r` to show clean HTTPS URLs so that I can quickly copy and share repo URLs with team members without exposing credentials or confusing them with SSH syntax.
- **As a developer**, I want a `--remote-raw` flag so that I can still see the original git remote URL when I need it for debugging or configuration purposes.

## Demoable Units of Work

### Unit 1: URL Formatting Utility

**Purpose:** Create a reusable function that converts raw git remote URLs into clean HTTPS URLs, serving as the foundation for display changes.

**Functional Requirements:**
- The system shall convert SSH remote URLs (e.g., `git@github.com:owner/repo.git`) to HTTPS format (`https://github.com/owner/repo.git`)
- The system shall strip user info from HTTPS remote URLs (e.g., `https://user@github.com/owner/repo.git` becomes `https://github.com/owner/repo.git`)
- The system shall preserve the `.git` suffix in formatted URLs
- The system shall return the original URL unchanged if it cannot be parsed or uses an unsupported protocol (e.g., `file://`, custom SSH aliases)
- The system shall handle Azure DevOps URLs (e.g., `git@ssh.dev.azure.com:v3/org/project/repo` and `https://user@dev.azure.com/org/project/_git/repo`)

**Proof Artifacts:**
- Test: `remote_test.go` unit tests pass demonstrating correct URL conversion for SSH, HTTPS, HTTPS-with-user-info, Azure DevOps, and unformattable URLs

### Unit 2: `--remote-raw`/`-R` Flag Registration

**Purpose:** Add the new flag to the list command so users can toggle between formatted and raw remote display.

**Functional Requirements:**
- The system shall register a `--remote-raw` flag with short form `-R` on the `list` command
- The flag shall follow the existing dual-purpose pattern: `-R` (no value) shows the raw remote column, `-R=pattern` filters by raw remote pattern and shows the column
- When `-R` is used without `-r`, the system shall implicitly enable remote display (showing raw URLs)
- When both `-r` and `-R` are used together, the system shall show raw URLs (i.e., `-R` overrides the formatted default)
- The system shall support flag stacking (e.g., `-Rsu` for remote-raw + status + user)

**Proof Artifacts:**
- CLI: `gws list -r` output shows formatted HTTPS URLs demonstrates default formatting works
- CLI: `gws list -R` output shows raw remote URLs demonstrates raw flag works
- CLI: `gws list -rR` output shows raw remote URLs demonstrates override behavior
- Test: Flag registration tests pass demonstrating `--remote-raw`/`-R` is properly registered

### Unit 3: Table and JSON Output Integration

**Purpose:** Update both table and JSON output rendering to use the formatted URL by default and respect the `--remote-raw` flag.

**Functional Requirements:**
- The system shall display formatted HTTPS URLs in the REMOTE table column by default when `-r` is used
- The system shall display raw remote URLs in the REMOTE table column when `-R` is used
- The system shall output formatted URLs in the `remote_url` JSON field by default when `-r` is used
- The system shall output raw remote URLs in the `remote_url` JSON field when `-R` is used
- The system shall preserve the existing asterisk (`*`) prefix behavior for repositories with multiple remotes

**Proof Artifacts:**
- CLI: `gws list -r -o json` output shows formatted `remote_url` values demonstrates JSON formatting
- CLI: `gws list -R -o json` output shows raw `remote_url` values demonstrates JSON raw mode
- CLI: Table output with multi-remote repo shows `* https://...` format demonstrates asterisk preservation

### Unit 4: Help Text and Documentation

**Purpose:** Update command help text so users can discover and understand the new behavior.

**Functional Requirements:**
- The `list` command help text shall describe `-r` as showing formatted remote URLs
- The `list` command help text shall describe `-R`/`--remote-raw` as showing raw remote URLs
- The examples section shall include usage of both `-r` and `-R`

**Proof Artifacts:**
- CLI: `gws list --help` output shows updated descriptions for `-r` and `-R` demonstrates discoverability

## Non-Goals (Out of Scope)

1. **Stripping `.git` suffix**: The formatted URL preserves `.git` — no "browser-friendly" URL normalization
2. **Multiple remote URL display**: No changes to how multiple remotes are indicated (asterisk behavior unchanged)
3. **Remote filtering implementation**: The existing unimplemented remote filtering is out of scope for this spec
4. **New columns**: No separate "formatted" vs "raw" columns — one REMOTE column with flag-driven behavior

## Design Considerations

No specific design requirements identified. The REMOTE column layout and alignment behavior remain unchanged; only the content (formatted vs. raw URL) changes.

## Repository Standards

- **Flag registration**: Follow existing dual-purpose flag pattern with `NoOptDefVal` sentinel
- **Testing**: Table-driven tests in `_test.go` files using Go's `testing` package
- **Code location**: URL formatting utility in `internal/git/remote.go`; flag/display logic in `cmd/git-workspace/list.go`
- **Commit messages**: Conventional Commits (`feat:`, `fix:`, `test:`)
- **Error handling**: Graceful fallback, no panics

## Technical Considerations

- URL parsing should handle both `git@host:path` SSH format and standard `scheme://[user@]host/path` format
- Azure DevOps SSH URLs use a different path structure (`v3/org/project/repo`) that must be handled specially
- The formatting function should be a pure function with no side effects for easy testing
- Consider using Go's `net/url` package for HTTPS URL parsing, with custom logic for SSH format

## Security Considerations

- The primary motivation for this feature is to strip user credentials from displayed URLs, which is a security improvement
- The formatting function must reliably remove `user@` and `user:password@` from HTTPS URLs
- Raw mode (`-R`) intentionally shows the original URL including any embedded credentials — this is expected behavior

## Success Metrics

1. **Default output is shareable**: `gws list -r` produces URLs that can be pasted into a browser or shared without modification
2. **No information loss**: `gws list -R` preserves the original raw remote for debugging
3. **All existing tests pass**: No regression in current remote display functionality

## Open Questions

No open questions at this time.
