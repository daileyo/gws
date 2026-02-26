# 06 Tasks - Commit Formatting

## Architecture Note

The spec described storing the hook at `scripts/hooks/commit-msg` with a symlink via `make install-hooks`. On review of the existing codebase, the repo already uses a different (cleaner) pattern: hooks live in `.githooks/` (tracked by git) and are activated by `git config core.hooksPath .githooks` via the existing `make setup-hooks` target. The `commit-msg` hook will follow this established pattern — placed directly in `.githooks/commit-msg`, piggybacking on `make setup-hooks` with no new Makefile target required.

## Relevant Files

- `.githooks/commit-msg` — CREATE: the POSIX sh commit-msg hook script that formats conventional commit messages.
- `.githooks/pre-push` — REFERENCE: existing hook; follow its file header comment style.
- `Makefile` — MODIFY: update the `## setup-hooks:` help comment and the target's echo messages to mention both hooks.
- `README.md` — MODIFY: add a one-time setup subsection to the Development section documenting `make setup-hooks`.

### Notes

- The hook script must use `#!/bin/sh` (not `#!/bin/bash`) for POSIX portability.
- Once `git config core.hooksPath .githooks` is set (via `make setup-hooks`), git automatically picks up any executable file in `.githooks/` by name — no additional registration step needed.
- The Makefile's `make help` output is generated from `## target:` comment lines via `sed` — update the comment line, not a separate help block.
- Follow the existing Makefile pattern: `@echo "..."` with feedback messages inside targets.

## Tasks

### [x] 1.0 Implement the `commit-msg` hook script

#### 1.0 Proof Artifact(s)

- `CLI`: `git commit -m "feat: add new command"` produces a commit whose subject reads `feat:       add new command` — demonstrates auto-formatting is transparent and correct.
- `CLI`: `git commit -m "refactor!: breaking API change"` produces `refactor!:  breaking API change` — demonstrates breaking-change token aligns to the same column.
- `CLI`: `git commit -m "ci: fix workflow"` produces `ci:         fix workflow` — demonstrates the shortest type gets the most padding.
- `CLI`: `git commit -m "WIP: something"` proceeds without error and prints a warning line to stderr — demonstrates non-conforming messages pass through unblocked.
- `CLI`: `git log --oneline` shows multiple commit types with descriptions starting at the same column — demonstrates visual alignment.

#### 1.0 Tasks

- [x] 1.1 Create `.githooks/commit-msg` with `#!/bin/sh` as the first line, followed by a header comment block that states the hook's purpose, the types it handles, and `# Activated by: make setup-hooks`. Look at `.githooks/pre-push` lines 1–5 for the expected comment style.

- [x] 1.2 Define the project type list as a space-separated shell variable (e.g., `TYPES="ci fix feat docs perf test chore style refactor"`) and implement a loop that computes the maximum token width by checking the length of each `type!:` string (e.g., `refactor!:` = 10 chars). Set `TARGET_WIDTH` to that maximum plus 2 (minimum two spaces of separation between token and description). This ensures all descriptions align to the same column even as types are added to the list.

- [x] 1.3 Implement first-line extraction and type matching: read the first line of the file at `$1` using `head -n 1`. Loop over `$TYPES`; for each type, first try to match the `type!:` prefix using `grep -q`, then try `type:`. When a match is found, record the full token string (e.g., `feat!:` or `feat:`) and extract the description by stripping the token and any following whitespace from the start of the line. Break out of the loop on the first match. Set a `matched` flag.

- [x] 1.4 Implement padding and file rewrite: compute `padding_len = TARGET_WIDTH - len(token)` using POSIX arithmetic (`$(( ))`). Build a padding string of that many spaces using `printf '%*s' "$padding_len" ''`. Construct the new first line as `token + padding + description`. Read the remainder of the commit message (all lines after line 1) using `tail -n +2 "$1"`. Write the new first line followed by the remainder back to `$1` using `printf`. Preserve any blank lines or footer content exactly as-is.

- [x] 1.5 Implement non-conforming message handling: if `matched` is still 0 after the loop, print a warning to stderr in the format `commit-msg: message does not match conventional commit format, skipping formatting` and `exit 0`. Do not use emoji (POSIX sh portability). Do not exit non-zero — the commit must be allowed to proceed.

- [x] 1.6 Make the script executable with `chmod +x .githooks/commit-msg`. Verify it is active by running `make setup-hooks` (if not already done). Then manually test by creating and immediately amending or dropping test commits using several type formats (`ci`, `feat`, `refactor`, `refactor!`, `WIP`) to confirm all alignment and pass-through behaviors match the proof artifacts above. Do not push these test commits.

### [x] 2.0 Update Makefile messaging and README documentation

#### 2.0 Proof Artifact(s)

- `CLI`: `make setup-hooks` prints output that names both the `pre-push` linting hook and the `commit-msg` formatting hook — demonstrates the install step self-documents what it activates.
- `CLI`: `make help` shows an updated description for `setup-hooks` that mentions commit formatting — demonstrates discoverability via help.
- `Diff`: `README.md` Development section includes a "One-Time Setup" subsection with `make setup-hooks` documented and explained — demonstrates contributor discoverability.

#### 2.0 Tasks

- [x] 2.1 In `Makefile`, update the `## setup-hooks:` comment line (the line used by `make help`) to reflect both hooks. Current value: `Install git hooks for pre-push linting`. New value should mention both pre-push linting and commit-msg formatting, e.g.: `Install git hooks (pre-push linting, commit-msg formatting)`.

- [x] 2.2 In `Makefile`, update the echo messages inside the `setup-hooks` target body to name both hooks so that `make setup-hooks` output is informative. Keep the same `@echo` style as the rest of the target. Example: replace the current single "Git hooks installed!" message with two lines — one naming the pre-push hook and one naming the commit-msg hook — followed by a brief description of what each does.

- [x] 2.3 In `README.md`, add a `### One-Time Setup` subsection to the `## Development` section, placed before the existing `### Build Commands` subsection. It should contain a brief explanation that git hooks are provided in `.githooks/` and must be activated once after cloning, followed by the command `make setup-hooks` in a `bash` fenced code block. Keep the tone and style consistent with the surrounding Development content.
