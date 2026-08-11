# 21 Questions Round 2 - Workspace Reconciliation

Round 1 settled most of the design. This round closes the gaps your answers opened. As before, tick boxes and add notes freely.

**Settled in round 1 (no need to revisit):**

- Strict pruning below a repository root; submodules and nested clones are invisible (Q1-A)
- No `branch` field on `Repository` (Q4-B)
- Existing summary wording kept; only the numbers become consistent (Q6-D)
- Status cache behavior unchanged — `refresh` clears it, `init` does not (Q7-D)
- Symlink maintenance stays in `refresh`; `init` never creates symlinks (Q8-B)
- `init` stays protective, exit 0, reworded to point at `refresh` (Q9-C)
- Existing skip list + hidden-dir skipping + configurable max depth + progress indicator (Q10-A/B/C/D)
- Go unit tests as proof artifacts (Q11-A)
- `docs/site/commands-core.md` and `README.md` updated in this spec (Q12-A)

---

## 1. Worktree Detection Rule

> **Verified:** a clone's `.git` is a **directory**. A linked worktree's `.git` is a **regular file** containing `gitdir: <main-repo>/.git/worktrees/<name>`. Inside a worktree, `git rev-parse --git-dir` differs from `git rev-parse --git-common-dir`; in a normal clone they are identical. So yes — worktrees are fully distinguishable without any naming convention, and a real clone living in a directory named `important-repo.wt` would be registered normally.

Which detection rule should discovery use?

- [ ] (A) `.git` is a directory → register as a repository. `.git` is a file → resolve its `gitdir:` target; if the target path sits under a `worktrees/` directory, treat it as a worktree and skip it. Pure filesystem reads, no subprocess.
- [ ] (B) Same as A, but confirm with `git rev-parse --git-dir` / `--git-common-dir` for every `.git`-file candidate (authoritative, costs one subprocess per candidate).
- [X] (C) A as the fast path, falling back to B only when the `gitdir:` target is ambiguous.
- [ ] (D) Other (describe)

## 2. Worktrees Whose Main Repository Is Outside the Workspace

If discovery walks into a worktree directory (e.g. `~/gws/proj.wt/feature`) whose main repository is **not** under the workspace root and is **not** already tracked, what should happen?

- [ ] (A) Skip it entirely. GWS only knows worktrees of repositories it already tracks; the user can `gws add` the main repo.
- [ ] (B) Register the main repository (resolved from the worktree's `gitdir:`) as a tracked repository, then its worktrees are discovered normally in phase 3.
- [ ] (C) Skip it, but print a note in the summary suggesting `gws add <main-repo-path>`.
- [X] (D) Other (describe) B feels like the intended flow. I don't want worktress of repos that aren't tracked. but I'd like the worktrees of all tracked repos.

## 3. Repository Identity — What Makes Two Entries Distinct

> Your round-1 note said directory+repo is the uniqueness key and that this layout should list `repo-1` three times:
>
> ```text
> ~/gws
>  |_ repo-1
>  |_ subdir/repo-1
>  |_ symlink-to/repo-1
> ```
>
> I need to know which of two readings you meant, because the third entry is the ambiguous one. If `symlink-to/repo-1` is a symlink resolving to the **same directory on disk** as `~/gws/repo-1`, it is the same clone reached by two paths. If it resolves to a **different** directory (an external clone that also happens to be named `repo-1`), it is genuinely a third repository.

- [ ] (A) Deduplicate by resolved real path. Three entries only when there are three distinct directories on disk; a symlink alias to an already-tracked repo does not create a second entry. (This is what the current code does — `Repository.Path` stores the resolved real path.)
- [ ] (B) Treat each reachable workspace path as its own entry, so a symlink alias to an already-tracked clone appears twice in `gws list`.
- [ ] (C) Deduplicate by real path, but remember the alias path so `gws list` can show how it was reached.
- [X] (D) Other (describe) A. gws should only care about the real path. It needs to be able to handle symlinks; but it it resolves to an already tracked repo, then skip. it should also be skipped if, say the symlink one was traversed first; then the physical one is seen later.

## 4. Preventing Symlink Accumulation Across Re-Runs

> This is the failure you remembered from before. `refresh` currently calls `ensureWorkspaceSymlink` for every repository whose real path is outside the workspace root. Once discovery follows symlinks at depth, a repo found via `~/gws/subdir/link → /elsewhere/repo` has an external real path, so refresh would create a **second** link at `~/gws/repo`. Re-running keeps producing aliases the user never asked for. `init` no longer creates symlinks (Q8-B), so this only affects `refresh`.

What rule should stop it?

- [ ] (A) Before creating a link, check whether the repository's real path is already reachable from anywhere under the workspace root (directly or through an existing symlink). If it is, create nothing.
- [ ] (B) `refresh` never creates symlinks at all — it only removes stale ones. Symlink creation belongs solely to `gws add`.
- [X] (C) Only create a link for repositories that have no reachable path under the workspace root **and** were previously tracked with a link (i.e. repair only).
- [ ] (D) Other (describe)

## 5. Traversal Loop and Depth Safeguards

Which safeguards should the scanner apply? (Select all that apply.)

- [X] (A) Maintain a visited set of resolved real paths; never traverse the same real directory twice within one scan.
- [X] (B) Refuse to descend into a symlink whose target is an ancestor of (or equal to) the workspace root.
- [X] (C) Refuse to descend into a symlink whose target is an ancestor of the current traversal path.
- [ ] (D) Cap total scan time and abort with a warning if exceeded.
- [ ] (E) Other (describe)

## 6. Max Traversal Depth — Default and Naming

You chose a configurable max depth in `config.Preferences`. What default and key?

- [ ] (A) Default `6`, key `scan_max_depth`. Depth counted in directory levels below the workspace root; a repo found at depth 6 is registered, nothing below is traversed.
- [ ] (B) Default `4`, key `scan_max_depth`.
- [ ] (C) Default `10`, key `scan_max_depth` — effectively unlimited for real workspaces but still loop-safe.
- [ ] (D) Unlimited by default; the preference exists for users who want to constrain it.
- [ ] (E) Other (describe)

## 7. Hidden Directory Skipping — Exact Scope

You chose to skip hidden directories. Where does that apply?

- [ ] (A) Skip any directory whose name starts with `.` at every level below the workspace root. A repo at `~/gws/.dotfiles` would **not** be discovered.
- [ ] (B) Skip hidden directories except direct children of the workspace root, so a deliberately placed `~/gws/.dotfiles` is still found.
- [ ] (C) Skip hidden directories, but always follow a hidden **symlink** in the workspace root (covers `~/gws/.config-nvim → ~/.config/nvim`).
- [ ] (D) B and C combined.
- [X] (E) Other (describe) A. but may like to have a configurable option to allow b and c. that is future functionality.

## 8. Progress Indicator Threshold

When should reconciliation show progress?

- [ ] (A) When the workspace has more than 25 tracked repositories.
- [ ] (B) When it has more than 50.
- [ ] (C) Time-based — show progress only once the operation has run for ~1 second, regardless of count.
- [X] (D) Always show it for `refresh` and `init`.
- [ ] (E) Other (describe)

## 9. Confirm Exact Command Output

You asked to keep existing wording and only make the numbers consistent. `init` now also discovers worktrees, so it needs lines it does not currently have. Please confirm or edit these blocks.

**`gws init ~/gws`:**

```text
Initialized workspace at: /home/daileyo/gws
Found 8 repositories.
Repositories with user configuration: 6
Repositories with worktrees: 2
```

**`gws refresh`** (unchanged from today except that the repository numbers now come from the same recursive scan):

```text
Refreshing workspace at: /home/daileyo/gws
Detecting git user configuration...
Cleared git status cache

Refresh complete!
Total repositories: 8
Removed 1 repository (path no longer valid)
Found 2 new repositories
Updated 3 repositories
Repositories with user configuration: 6
Repositories with worktrees: 2
```

- [ ] (A) Both blocks are correct as shown.
- [X] (B) Correct, but `init` should also report aligned/unaligned worktree counts.
- [ ] (C) Correct, but `init` should report new/removed counts too, for symmetry with `refresh`.
- [ ] (D) Other (edit the blocks above or describe)

## 10. Where the Shared Function Lives

Your request sketched `reconcileWorkspace(workspaceRoot, existingConfig?)`. Given that cache clearing and symlink maintenance stay in `refresh`, this function is pure discovery + metadata merge. Where should it live?

- [X] (A) New `internal/reconcile` package, with `internal/discovery` reduced to the filesystem scan it already owns.
- [ ] (B) Extend `internal/discovery` with the reconcile entry point — one package owns "figure out what is really there."
- [ ] (C) Keep it as a shared helper in `cmd/git-workspace`, alongside the existing `buildRepository` helper.
- [ ] (D) Other (describe)

## 11. Removing the Parent-Directory Rescan Special Case

> **Current behavior:** when `refresh` finds a tracked repository whose path is gone, it scans that repository's *parent directory* to pick up replacements. With a full recursive scan of the workspace root, this is redundant — and it can reach outside the workspace root when the missing repo was external.

- [ ] (A) Delete it. The recursive scan covers everything under the workspace root.
- [ ] (B) Delete it, but confirm through a test that repositories replaced in place are still discovered.
- [X] (C) Keep it as a safety net.
- [ ] (D) Other (describe)
