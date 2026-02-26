# 06 Questions Round 1 - Commit Formatting

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Hook Trigger Point

At what point in the commit workflow should the formatting be applied?

- [X] (A) **`commit-msg` hook** — auto-reformats the message after you write it, before the commit is finalized. Transparent: you type `feat: my message`, it becomes `feat:      my message` silently.
- [ ] (B) **`prepare-commit-msg` hook** — fires before the editor opens. Can insert a formatted template, but since you haven't written the message yet, it can't pad the type until you've typed it.
- [ ] (C) **Standalone script only** — a `gws-fmt-commit` script (or similar) you run manually against a draft message file, with no automatic hook behavior.
- [ ] (D) Other (describe)

## 2. Hook Installation

Git hooks live in `.git/hooks/` which is not tracked by git. How should the hook be distributed and installed?

- [X] (A) **Script in repo + Makefile target** — store the hook script at `scripts/hooks/commit-msg` (tracked in git), add a `make install-hooks` target that symlinks or copies it into `.git/hooks/`. Contributors run `make install-hooks` once after cloning.
- [ ] (B) **Script in repo + automatic install via existing hook** — the `pre-push` hook (which already runs linting and tests) also checks if the `commit-msg` hook is installed and installs it if not.
- [ ] (C) **Script only, no Makefile target** — just commit the script file; users install it manually by copying/symlinking to `.git/hooks/commit-msg`.
- [ ] (D) Other (describe)

## 3. Breaking Change Syntax (`feat!:`)

The project uses `feat!:` for breaking changes (pattern 1, no scope). Should the `!` be included in the padding calculation?

- [X] (A) **Yes, treat `feat!:` as its own token** — pad it separately from `feat:`. Both align their descriptions to the same column, with `feat!:` padded to match `refactor:` width (e.g., `feat!:    description`).
- [ ] (B) **Yes, but `!` shifts the column** — since `feat!:` is one char longer than `feat:`, accept that breaking-change commits have descriptions one position further right. No special handling needed.
- [ ] (C) **Out of scope for now** — skip formatting for any message containing `!`. Leave it as-is.
- [ ] (D) Other (describe)

## 4. Non-Conforming Messages

What should happen when a commit message does NOT match the conventional commit pattern (e.g., a merge commit, a WIP note, or a typo in the type)?

- [ ] (A) **Skip silently** — if the message doesn't parse as a conventional commit, leave it unchanged and let the commit proceed normally.
- [X] (B) **Warn but proceed** — print a warning to stderr that the message wasn't formatted, but allow the commit to complete.
- [ ] (C) **Error and block** — reject the commit if it doesn't conform to the pattern. (Note: this turns the hook into a linter, not just a formatter.)
- [ ] (D) Other (describe)

## 5. Padding Target

Which set of commit types should define the alignment column width?

- [ ] (A) **Full conventional commits spec** — use `refactor:` (9 chars) as the fixed maximum, since it's the longest valid type. Hardcode this width.
- [X] (B) **Project-defined types only** — derive the width dynamically from the types listed in README.md (`ci`, `fix`, `feat`, `docs`, `perf`, `test`, `chore`, `style`, `refactor`). Same result today, but automatically adjusts if types are added/removed.
- [ ] (C) Other (describe)
