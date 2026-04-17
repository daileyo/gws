# 19 Questions Round 1 - Worktree Management

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Worktree Discovery Scope

During `gws refresh`, how should worktrees be discovered?

- [ ] (A) Only detect worktrees that already follow the `<repo>.wt/` convention
- [ ] (B) Use `git worktree list` on each tracked repo to discover all worktrees regardless of current location
- [ ] (C) Both — discover all via `git worktree list`, but only track/manage ones that follow the convention (flag non-conforming ones as "unaligned")
- [X] (D) Other (describe) C, but track them for being able to naviagate to them

## 2. Worktree Data Storage

Where should discovered worktree data be stored?

- [ ] (A) Add a `Worktrees` field to the existing `config.Repository` struct in `config.json` (e.g., `"worktrees": [{"name": "feature-x", "path": "/path/to/repo.wt/feature-x", "branch": "feature-x"}]`)
- [ ] (B) Store worktree data in a separate file (e.g., `~/.gws/worktrees.json`)
- [ ] (C) Don't persist worktree data at all — always discover dynamically via `git worktree list` at runtime
- [X] (D) Other (describe) discover and track them in config.Repository struct in config.json for now... unless it seems like its getting too bloated.

## 3. The `(wt)` Indicator in List Output

When running `gws list` (or bare `gws`), where should the `(wt)` indicator appear?

- [X] (A) Inline next to the repo name, e.g., `my-repo (wt)`
- [ ] (B) As a dedicated column (similar to type/visibility columns) shown with a flag like `-W`
- [ ] (C) Both — inline in compact mode, column in verbose/table mode
- [ ] (D) Other (describe)

## 4. Worktree Naming Convention

When navigating with `gws <repo> -wt <worktree-name>`, what should `<worktree-name>` refer to?

- [ ] (A) The directory name inside `<repo>.wt/` (e.g., if the worktree is at `repo.wt/feature-x`, the name is `feature-x`)
- [X] (B) The branch name checked out in the worktree
- [ ] (C) Either — match against directory name first, fall back to branch name
- [ ] (D) Other (describe)

## 5. `gws worktree align` — Handling Edge Cases

When `gws worktree align` moves worktrees into the `<repo>.wt/` directory, how should these situations be handled?

- [X] (A) **Naming conflicts**: If two worktrees would end up with the same directory name inside `<repo>.wt/`, should the tool error out, auto-rename (append suffix), or prompt the user?
- [ ] (B) **Bare worktrees outside the workspace**: Should `align` only handle worktrees for repos tracked by gws, or also discover/move worktrees for any git repo it finds?
- [X] (C) **Dry-run mode**: Should `gws worktree align --dry-run` be supported to preview what would be moved before actually moving?

Please indicate your preference for each sub-question above.
For A: I tshould append -dup-xx suffix where xx is 00-99 and incremented based on the number of duplicates.
For C: yes it should be, especially to note potential duplicate worktrees
## 6. `gws worktree align` — Move Strategy

When moving a worktree, the `.git` file inside the worktree and the `worktrees/<name>/gitdir` file in the main repo both need updating. How should the move be performed?

- [X] (A) Use `git worktree move` command (Git 2.17+), which handles the internal bookkeeping automatically
- [ ] (B) Manually move the directory and update the `.git` pointer files (more control, works with older Git versions)
- [ ] (C) Prefer `git worktree move` but fall back to manual if Git version is too old
- [ ] (D) Other (describe)

## 7. Additional Worktree Subcommands

Beyond `align`, should `gws worktree` support other subcommands now or in the future?

- [ ] (A) Just `align` for now — keep it minimal
- [ ] (B) Also include `gws worktree list` to show all worktrees across all repos
- [ ] (C) Also include `gws worktree add <repo> <branch>` to create a new worktree following the convention
- [X] (D) Include both `list` and `add` alongside `align`
- [ ] (E) Other (describe)

## 8. Shell Integration for `-wt` Flag

The `gws` shell function currently intercepts navigation to perform `cd`. For `gws <repo> -wt <name>`, the shell function will need updating. Should the `-wt` flag:

- [ ] (A) Be handled entirely in the shell function (parse `-wt` in shell, call a new binary subcommand to resolve the path)
- [ ] (B) Be handled as a flag on the root command in Go (`git-workspace -g <repo> --worktree <name>`), with the shell function passing it through
- [X] (C) Other (describe) B, if that will result in it navigating to where the worktree is. otherwise A. Also, if it is easy to tell which will be most performant and the simplest code change, then that would also be a factor for choosing.

## 9. Proof Artifacts

What would you consider convincing proof that this feature works? Select all that apply:

- [ ] (A) CLI output screenshot: `gws list` showing `(wt)` indicators on repos with worktrees
- [ ] (B) CLI output: `gws <repo> -wt` displaying the enumerated worktree selection list
- [ ] (C) CLI output: `gws worktree align` moving worktrees and showing before/after
- [ ] (D) CLI output: `gws refresh` discovering and registering worktree data
- [X] (E) Test suite passing for all new worktree-related code
- [ ] (F) Other (describe)
