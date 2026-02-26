# 06-spec-commit-formatting.md

## Introduction/Overview

Commit messages in this project use conventional commit types (`feat:`, `fix:`, `refactor:`, etc.) that vary in length, causing description text to appear ragged and misaligned in log views. This feature adds a `commit-msg` git hook that automatically pads the type token with spaces so that all commit descriptions begin at the same column. The hook is stored in the repository and installed locally via a Makefile target.

## Goals

- All conventional commit descriptions align at a consistent column when viewed in `git log` or `gws`-style log output.
- The alignment column is derived from the project's defined type set (including `!` breaking-change variants), so no manual padding is required from the author.
- The hook is transparent: authors write commits normally; reformatting happens silently.
- The hook and its installation are reproducible and self-documented via a Makefile target.
- Non-conforming messages (merge commits, WIPs, etc.) are warned about but never blocked.

## User Stories

- **As a developer**, I want my conventional commit messages automatically padded so that I don't have to manually count spaces to align descriptions.
- **As a developer**, I want to install the commit formatting hook in one command so that I don't have to manually copy files into `.git/hooks/`.
- **As a developer**, I want non-conventional commit messages (WIPs, merge commits) to pass through unmodified so that the hook never blocks legitimate work.

## Demoable Units of Work

### Unit 1: commit-msg Hook Script

**Purpose:** Implements the formatting logic as a shell script stored in the repository and installed into `.git/hooks/`.

**Functional Requirements:**
- The script shall be stored at `scripts/hooks/commit-msg` and be executable.
- The script shall accept the commit message file path as `$1` (standard git `commit-msg` hook contract).
- The script shall parse the first line of the commit message against the pattern `^(type)(!?): (.+)` where `type` is any project-defined conventional commit type.
- The script shall derive the alignment column width from the longest possible token across all project-defined types including their `!` (breaking change) variants. For the current type set (`ci`, `fix`, `feat`, `docs`, `perf`, `test`, `chore`, `style`, `refactor`), the longest token is `refactor!:` (10 characters).
- The script shall pad the type token with spaces so that all descriptions begin at column `(longest_token_width + 1)`.
- For a matching message, the script shall rewrite the commit message file in place with the padded format. The body and footer (lines after the first blank line) shall be preserved unchanged.
- For a non-matching message, the script shall print a warning to stderr and exit 0 (allowing the commit to proceed unchanged).
- The script shall be POSIX-compatible shell (`#!/bin/sh`) to avoid dependency on bash or any language runtime.

**Proof Artifacts:**
- `CLI`: Running `git commit -m "feat: add new command"` produces a commit whose first line reads `feat:      add new command` (padded to `refactor!:` width) — demonstrates auto-formatting works transparently.
- `CLI`: Running `git commit -m "WIP: not formatted"` proceeds without error and prints a warning to stderr — demonstrates non-conforming messages are warned but not blocked.
- `CLI`: Running `git commit -m "refactor!: breaking change"` produces `refactor!: breaking change` (no extra padding needed, already at max width) — demonstrates breaking-change syntax is handled correctly.
- `CLI`: Running `git log --oneline` shows all type tokens left-aligned with descriptions starting at the same column — demonstrates visual alignment in log output.

### Unit 2: Makefile Install Target and Documentation

**Purpose:** Makes the hook installable in one command and documents the setup step for contributors.

**Functional Requirements:**
- The Makefile shall include an `install-hooks` target that symlinks `scripts/hooks/commit-msg` into `.git/hooks/commit-msg`.
- The `install-hooks` target shall be idempotent: running it multiple times shall not produce an error.
- The `install-hooks` target shall print a confirmation message indicating the hook was installed (or was already in place).
- The `README.md` Development section shall document `make install-hooks` as a recommended one-time setup step after cloning.
- The `make help` output (or equivalent) shall include a description of the `install-hooks` target.

**Proof Artifacts:**
- `CLI`: Running `make install-hooks` on a fresh clone creates `.git/hooks/commit-msg` as a symlink to `scripts/hooks/commit-msg` and prints a success message — demonstrates the install target works.
- `CLI`: Running `make install-hooks` a second time completes without error — demonstrates idempotency.
- `Screenshot`: README Development section showing `make install-hooks` documented as a setup step — demonstrates contributor discoverability.

## Non-Goals (Out of Scope)

1. **Scoped commits (Pattern 2)**: Commits with an optional scope (`feat(scope): description`) are explicitly out of scope for this spec. The hook shall pass these through unchanged (treating them as non-conforming for now).
2. **Commit message linting**: This feature formats messages; it does not validate or enforce correctness. A wrong type or missing description is warned about but not blocked.
3. **Reformatting historical commits**: Only new commits are formatted. No tooling to rewrite existing history is included.
4. **Multi-line subject line handling**: Only the first line of the commit message is formatted. Multi-paragraph bodies are preserved as-is.
5. **CI enforcement**: No CI check that validates padding on existing commits in the repo.

## Design Considerations

No specific UI/UX requirements. The formatting output format is:

```
<type><padding>: <description>
```

Where `<type><padding>` is left-padded with spaces so that `:` lands at column `(longest_token_width)` and the description begins at column `(longest_token_width + 2)`.

Example with current type set (longest token `refactor!:` = 10 chars, description starts at column 12):

```
ci:         add workflow for docs deployment
fix:        correct requirements.txt path
feat:       add shell-init command
chore:      update dependencies
style:      fix linter warnings
refactor:   simplify discovery logic
refactor!:  breaking API change
```

## Repository Standards

- Shell scripts in this project use `#!/bin/sh` for POSIX compatibility.
- Makefile targets follow the existing pattern: short lowercase names, printed feedback on execution.
- Commit messages follow Conventional Commits as documented in `README.md`.
- Hook script lives under `scripts/` to keep tooling separate from source code.

## Technical Considerations

- The project-defined type list shall be maintained as a hardcoded list within the hook script itself, mirroring the types documented in `README.md`. If a new type is added to the project, both `README.md` and the hook script must be updated.
- The alignment width is computed at hook execution time from the hardcoded type list (including `!` variants), so adding a longer type automatically widens the column without requiring a magic number change.
- Git passes the commit message file path as `$1` to the `commit-msg` hook. The script must read and rewrite this file, not stdout.
- The symlink approach for `make install-hooks` is preferred over copying the file, so that updates to `scripts/hooks/commit-msg` (via `git pull`) are automatically reflected without re-running the install.
- The hook must handle the case where `.git/hooks/commit-msg` already exists (either as an old copy or a different hook). The install target shall warn the user if an existing non-symlink file is present and skip overwriting it.

## Security Considerations

No specific security considerations identified. The hook script operates only on local commit message files and does not access the network, credentials, or sensitive data.

## Success Metrics

1. **Alignment**: `git log --oneline` output shows all commit descriptions starting at a consistent column across all commit types.
2. **Transparency**: No change to the commit authoring workflow — authors do not need to manually pad messages.
3. **Installability**: `make install-hooks` completes successfully on a fresh clone in a single command.

## Open Questions

No open questions at this time.
