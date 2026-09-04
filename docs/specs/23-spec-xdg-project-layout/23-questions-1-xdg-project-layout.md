# 23 Questions Round 1 - XDG Project Layout

Answers captured from the planning session on 2026-09-03.

## 1. Scope — what is actually changing?

- [X] (A) **Only the location of the worktree directory.** A "project" is a git repo plus any
      associated worktrees. The existing passive discovery model stays exactly as it is —
      including reporting aligned vs. unaligned after discovery. Nothing about how repos are
      tracked gets reworked.
- [ ] (B) Rework repo tracking into per-repo record files
- [ ] (C) Introduce named multi-repo workspaces

**Notes:** "It's already handled and I don't want to change how we already handle it in this
work." The goal is to reduce noise in the user's home space and to shift toward standardized
locations for this kind of thing.

## 2. Which XDG base directory should hold the worktrees?

- [X] (A) **`XDG_DATA_HOME`** — `~/.local/share/gws/projects/`. Worktrees are real user data;
      losing one loses work, so DATA is the correct XDG class.
- [ ] (B) `XDG_STATE_HOME` — `~/.local/state/gws/projects/`
- [ ] (C) DATA by default but user-overridable via config + env var

## 3. Directory naming and collisions under the projects root

- [X] (A) **Reuse what already exists.** The `-dup-NN` suffix scheme in `worktree_align.go`
      already handles naming conflicts. Do not redesign it.
- [ ] (B) Hash-suffix on repo-name collision
- [ ] (C) Mirror the remote URL (`<host>/<owner>/<repo>/`)
- [ ] (D) Always hash-suffixed

**Notes:** Don't overthink this part — it's already handled. The repo-name-level collision
(two tracked repos sharing a basename) is recorded as an open question rather than designed
for in this spec.

## 4. What happens to worktrees still in an old `<repo>.wt/` directory?

- [X] (A) **They become unaligned, and `gws worktree align` moves them.** Only the XDG projects
      root counts as aligned. One meaning of "aligned", no new states.
- [ ] (B) A distinct `(legacy)` state separate from `(unaligned)`
- [ ] (C) Both locations continue to count as aligned

## 5. How should existing installs move over?

- [X] (A) **Config auto-migrates; worktrees only move on explicit `align`.** The config file is
      small and gws-owned, so relocating it on first run is safe. A worktree may hold
      uncommitted work, so it is never relocated without the user asking.
- [ ] (B) Auto-migrate both on first run
- [ ] (C) Explicit `gws migrate` for both

**Notes:** Assumed default, carried forward from the round-1 discussion. Flagged in Open
Questions so it can be revisited before implementation.
