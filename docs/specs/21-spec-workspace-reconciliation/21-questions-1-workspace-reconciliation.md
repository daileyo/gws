# 21 Questions Round 1 - Workspace Reconciliation

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

Context notes from the code review are included above some questions where the current behavior is not obvious.

## 1. Repository Boundary Pruning

> **Current behavior:** `discovery.Scan` uses `filepath.Walk` and only skips the `.git` directory itself — it keeps descending *inside* a repository, so a vendored clone or submodule nested under a repo is registered as its own top-level repository.

Your request says "once a Git repository root is found, register it and prune traversal below that repository." How strict should that pruning be?

- [X] (A) Strict prune — once a repo root is found, never descend below it. Submodules and nested clones are invisible to GWS.
- [ ] (B) Strict prune, but register submodules declared in `.gitmodules` as their own tracked repositories.
- [ ] (C) Prune, but keep descending to find nested *independent* clones (anything not a submodule).
- [ ] (D) Keep the current behavior (no pruning) and register every nested `.git` directory found.
- [ ] (E) Other (describe)

## 2. Worktree Directories During Repository Discovery

> **Current behavior:** a git worktree has `.git` as a *file*, not a directory. The scanner only matches `.git` directories, so worktrees are never registered as repositories today. Recursive scanning will now walk into `<repo>.wt/`, so this needs an explicit rule.

How should discovery treat directories that are git worktrees rather than repository clones?

- [ ] (A) Detect and skip any directory whose `.git` is a file (a worktree or submodule pointer) — worktrees are only ever discovered in the phase-3 `git worktree list` pass.
- [ ] (B) Skip directories by naming convention only (anything under a `<repo>.wt/` directory).
- [ ] (C) Both A and B — belt and braces.
- [ ] (D) Register worktrees as repositories as well as worktrees.
- [X] (E) Other (describe) Is it not possible to decipher between a worktree and a git repository? I'd prefer to not have to work on assumed directory structure naming for filtering. the idea is that gws can work with the existing structure fine. the .wt directory format is a passive opinion; but if a user for whatever reason his a git repo in a folder named important-repo.wt... i don't want it ignored

## 3. Symlink Resolution Depth

> **Current behavior:** `init` (via `filepath.Walk`) does not follow symlinks at all, so a curated symlink-only workspace yields zero repositories. `refresh` resolves symlinks, but only for direct children of the workspace root.

How deep should the unified scanner resolve symlinks?

- [ ] (A) First level only — resolve symlinks that are direct children of the workspace root; never follow symlinks encountered deeper in the tree. (Matches the curated-workspace model in your request.)
- [ ] (B) Resolve symlinks at any depth, guarded by a visited-real-path set to prevent cycles.
- [ ] (C) Resolve symlinks that are direct children of the workspace root *or* direct children of a container directory (i.e. symlinked repos may be grouped one level down).
- [ ] (D) Never follow symlinks; external repositories stay tracked purely through their stored real paths.
- [X] (E) Other (describe) I'm pretty certain we debugged and explicitly wanted to exclude symlinks from traversal in the past. it is possible for a this to happen
~/gws
  |_repo-1
  |_symlink-to repo like ~/.config/nvim

as long as we don't end up with a recursive loop; that re-registers a symlinked repo; symlink traversal can be good.
I believe we ran into an issue of symlinks to symlinks starting to occur if refresh or init were re-ran. refresh is expected to be reran by design. init should be resiliant enough to be accidentally re-ran

## 4. Repository `branch` Metadata

> **Current behavior:** `config.Repository` has no `branch` field. Branch information exists only per-worktree (`config.Worktree.Branch`) and in the transient git status cache.

Your phase-2 list mentions updating "branch" as mutable repository metadata. What did you mean?

- [ ] (A) Add a `branch` field to `config.Repository`, populated during reconciliation (requires a config schema version bump).
- [X] (B) No new field — "branch" refers to the existing per-worktree branch data and the status cache only.
- [ ] (C) Add the field, but leave it empty during reconciliation and populate it lazily elsewhere.
- [ ] (D) Other (describe)

## 5. Duplicate Repository Names

> **Current behavior:** `Name` is `filepath.Base(path)` and is used for navigation and for workspace symlink creation. Recursive discovery through container directories makes collisions likely (e.g. `org-a/api` and `org-b/api`).

How should reconciliation handle two discovered repositories with the same directory name?

- [ ] (A) Track both under the same name; rely on the existing multi-match interactive selection for navigation. Create a workspace symlink for the first only, and skip the conflicting one.
- [ ] (B) Track both, and qualify the display name as `<parent-dir>/<name>` when — and only when — a collision exists.
- [ ] (C) Track both; skip symlink creation entirely for all colliding names and print a warning.
- [ ] (D) Not a concern for this spec — handle collisions in a future spec, keep whatever falls out of the current code.
- [X] (E) Other (describe) Collisions should already be handled in the existing functionality. the expectation is that the directory+repo combo is the uniqueness. if this happens

~/gws
 |_ repo-1
 |_ subdir/repo-1
 |_ symlink-to/repo-1

 it's expecta that gws list would show repo-1 three times. gws repo-1 would result in a selector to pick the correct one. I don't want this changed; but if there is a gap that has been identified during this spec; it's worth discussing.

## 6. Summary Output Format

Your request shows a reporting block for `init`. What exactly should both commands print?

- [ ] (A) Exactly the block from your request for both commands (only the header line differs: "Initialized workspace at:" vs "Refreshed workspace at:").
- [ ] (B) That block plus repository `added` / `removed` / `updated` counts, since the acceptance criteria depend on those being visible.
- [ ] (C) Option B plus the names of the repositories that were added and removed.
- [X] (D) Keep each command's existing summary wording; only make the underlying numbers consistent.
- [ ] (E) Other (describe)

## 7. Git Status Cache Handling

> **Current behavior:** `refresh` clears the whole status cache at the end. `init` does not touch it.

Where does the cache fit in the shared reconciliation function?

- [ ] (A) Reconciliation always clears the status cache at the end — same code path for both commands (a no-op on a fresh init).
- [ ] (B) Cache clearing stays outside reconciliation, in the `refresh` command only.
- [ ] (C) Reconciliation refreshes cache entries in place rather than clearing them.
- [X] (D) Other (describe) Unless it is somehow necessary; existing cache functionality should stay as is.

## 8. Workspace Symlink Maintenance

> **Current behavior:** `refresh` calls `ensureWorkspaceSymlink` / `removeWorkspaceSymlink` to keep `~/gws/<name> → /real/path` links in sync for repositories that live outside the workspace root. `init` does not.

Should symlink maintenance live inside the shared reconciliation function?

- [ ] (A) Yes — reconciliation owns creating and removing workspace symlinks for external repositories, for both commands.
- [ ] (B) No — keep symlink maintenance in `refresh` only; `init` never creates symlinks.
- [ ] (C) Reconciliation creates missing symlinks but never removes existing ones (safer, may leave stale links).
- [X] (D) Other (describe) This should be B. I see no scenario where init should create symlinks

## 9. `gws init` When Metadata Already Exists

> **Current behavior:** prints "Workspace already initialized at: ..." to stderr, suggests `gws add` / `gws refresh`, and exits 0.

You said `init` should "remain protective." What is the exact behavior you want?

- [ ] (A) Keep the current behavior exactly (message to stderr, exit code 0).
- [ ] (B) Same message, but exit with a non-zero status so scripts can detect it.
- [X] (C) Keep current behavior, but reword so `gws refresh` is clearly the recommended next step.
- [ ] (D) Add a `--force` flag that discards and recreates the metadata library.
- [ ] (E) Other (describe)

## 10. Traversal Safeguards and Performance

> **Current behavior:** `shouldSkipDir` skips `node_modules`, `vendor`, `.venv`, `venv`, `__pycache__`, `.tox`, `target`, `build`, `dist`, `.cache`, `.terraform`. There is no depth limit. Reconciliation also runs `git worktree list` (plus repair/prune) and git user detection per repository, so cost scales with repository count.

Which safeguards should this spec include? (Select all that apply.)

- [X] (A) Keep the existing skip list as-is.
- [X] (B) Also skip hidden directories (any directory starting with `.`), except the workspace root itself.
- [X] (C) Add a maximum traversal depth, configurable via `config.Preferences`.
- [X] (D) Show a progress indicator during reconciliation when the workspace is large (`internal/git/progress.go` already exists).
- [ ] (E) Other (describe)

## 11. Proof Artifacts

What should the proof artifacts for this spec be?

- [X] (A) Go unit tests only, following the existing table-driven test conventions in the repo.
- [ ] (B) Unit tests plus captured CLI transcripts (`init` and `refresh` against a fixture workspace) saved in `21-proofs/`.
- [ ] (C) Option B plus a reusable shell script that builds the fixture workspace (nested containers, external symlinks, worktrees, deleted repos) so results are reproducible.
- [ ] (D) Option C plus a before/after `config.json` diff showing tags preserved and worktrees discovered.
- [ ] (E) Other (describe)

## 12. Documentation Scope

> **Repo standard:** spec 18 (`mkdocs-sync`) established that `docs/site/` and `README.md` must stay in sync with CLI behavior.

Should documentation updates be part of this spec?

- [X] (A) Yes — update `docs/site/commands-core.md` (the `Initialize Workspace` and `Refresh Workspace` sections) and `README.md` as part of the work.
- [ ] (B) Only update the example output blocks that change; leave prose alone.
- [ ] (C) No — track documentation in a separate follow-up spec.
- [ ] (D) Other (describe)
