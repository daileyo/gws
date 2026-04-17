# 19 Questions Round 2 - Worktree Management

A few follow-ups based on your Round 1 answers.

## 1. Worktree Navigation — Branch Name Matching

You chose branch name for `-wt <worktree-name>` matching. Branches can have slashes (e.g., `feature/auth-flow`). When matching, should:

- [ ] (A) Require exact branch name match (user types full `feature/auth-flow`)
- [X] (B) Support partial/substring matching (e.g., `auth-flow` matches `feature/auth-flow`), same as repo name matching works today
- [ ] (C) Other (describe)

## 2. `gws worktree add` — Behavior

When running `gws worktree add <repo> <branch>`:

- [X] (A) Create the worktree directly in `<repo>.wt/<branch>` (creating the `.wt` directory if needed), using `git worktree add`
- [ ] (B) Same as A, but use the last segment of the branch name as the directory name (e.g., `feature/auth-flow` → `<repo>.wt/auth-flow`)
- [ ] (C) Other (describe)

## 3. `gws worktree list` — Scope

When running `gws worktree list`:

- [ ] (A) Show worktrees for all repos that have them (global view across workspace)
- [ ] (B) Accept an optional repo name argument to filter to one repo: `gws worktree list [repo]`
- [X] (C) Both — no argument shows all, with argument filters to that repo
- [ ] (D) Other (describe)

## 4. Unaligned Worktree Visibility

You want to discover all worktrees (via `git worktree list`) and track them even if unaligned. Should unaligned worktrees:

- [ ] (A) Show an `(unaligned)` indicator in `gws worktree list` output
- [ ] (B) Be navigable with `-wt` even when not in the `.wt/` directory
- [X] (C) Both A and B
- [ ] (D) Other (describe)

## 5. `gws worktree align` — Scope per Run

Should `gws worktree align`:

- [ ] (A) Align all repos at once (workspace-wide)
- [ ] (B) Accept an optional repo name to align just one repo: `gws worktree align [repo]`
- [X] (C) Both — no argument aligns all, with argument targets one repo
- [ ] (D) Other (describe)
