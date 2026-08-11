# 21-tasks-workspace-reconciliation

## Relevant Files

- `internal/config/config.go` - Add `ScanMaxDepth` to `Preferences`, a `DefaultScanMaxDepth` constant, and an effective-value helper
- `internal/config/config_test.go` - Tests for preference default resolution and backward-compatible loading of older `config.json` files
- `internal/discovery/scanner.go` - Rewrite: explicit recursive traversal with boundary pruning, depth cap, hidden-dir skipping, symlink resolution, real-path dedupe; export `BuildRepository`
- `internal/discovery/scanner_test.go` - Tests for pruning, filtering, depth, symlinks, loop safety, worktree exclusion, dedupe
- `internal/git/worktree.go` - Add `IsLinkedWorktree` (structural detection plus `git rev-parse` fallback) returning the main repository path
- `internal/git/worktree_test.go` - Tests for worktree detection, the `.wt`-named-clone case, and the ambiguous-`.git`-file fallback
- `internal/git/progress.go` - Allow a caller-supplied label; `NewProgress` currently hardcodes `"Checking repos... %d/%d"`
- `internal/git/progress_test.go` - Tests for the label change and existing default behavior
- `internal/reconcile/reconcile.go` - New: `ReconcileWorkspace` entry point orchestrating discovery, merge, and worktree phases
- `internal/reconcile/reconcile_test.go` - Tests for the `init` path (nil config) and the `refresh` path (existing config)
- `internal/reconcile/merge.go` - New: merge discovered repositories against existing metadata; tag and user-config preservation
- `internal/reconcile/merge_test.go` - Tests for retain/add/remove/update rules and tag preservation
- `internal/reconcile/worktree.go` - New: worktree discovery phase, alignment classification, orphan main-repository registration
- `internal/reconcile/worktree_test.go` - Tests for worktree phase and orphan handling
- `internal/reconcile/result.go` - New: `Result` struct carrying all summary counts, scan errors, and reachable real paths
- `internal/reconcile/errors.go` - New: shared scan-error reporter (stderr, cap 5, `... and N more errors`)
- `internal/reconcile/errors_test.go` - Tests for reporter output shape and truncation
- `cmd/git-workspace/init.go` - Rewrite `runInit` to delegate to `ReconcileWorkspace`; new summary lines; reworded guard message
- `cmd/git-workspace/init_test.go` - Tests for init summary output, symlink workspaces, worktree persistence, guard behavior
- `cmd/git-workspace/refresh.go` - Rewrite `runRefresh` to delegate to `ReconcileWorkspace`; retain cache clear, safety net; repair-only symlinks
- `cmd/git-workspace/refresh_test.go` - Tests for recursive rescan, removals, worktrees, idempotence, symlink repair, cross-command parity
- `cmd/git-workspace/add.go` - Update the `discovery.Scan` call site for the new signature; remove the duplicate `buildRepository`
- `cmd/git-workspace/add_test.go` - Confirm `gws add --recursive` still behaves correctly under the new scanner
- `cmd/git-workspace/userdetect.go` - No change expected; `detectUserForRepos` and `syncProfilesFromRepos` are called identically by both commands
- `docs/site/commands-core.md` - Update `Initialize Workspace` and `Refresh Workspace` sections and example output
- `docs/site/configuration.md` - Document the `scan_max_depth` preference
- `README.md` - Sync `init` and `refresh` descriptions with the docs site

### Notes

- Unit tests live alongside the code they test (`scanner.go` and `scanner_test.go` in the same directory).
- Run `go test ./...` for the full suite, or target a package with `go test ./internal/discovery/`. Use `make test-race` and `make lint` as the quality gates before considering a parent task done.
- Follow the existing Cobra pattern: `var xxxCmd = &cobra.Command{...}` plus `func init()` with `rootCmd.AddCommand`, and a thin `RunE` delegating to a testable `run<Name>` function.
- Use the `gitCommand()` helper in `internal/git/status.go` for git subprocess calls.
- Use `t.TempDir()` for filesystem fixtures and `bytes.Buffer` for output capture, as `worktree_list_test.go` already does.
- Config is at `~/.gws/config.json`. Tests that touch it must isolate `HOME` via `t.Setenv` so they never read or write the developer's real workspace.
- Conventional commits: `refactor:` for the engine extraction, `feat:` for the new preference and behavior, `fix:` for the symlink accumulation bug, `docs:` for task 6.0.
- **Scope note:** `discovery.Scan` has three callers — `init.go`, `refresh.go`, and `add.go`. Changing it means `gws add --recursive` inherits boundary pruning, hidden-directory skipping, and the depth cap. The spec lists `gws add` as a non-goal, so the intent is not to redesign it, but this inherited change is real and is called out in task 1.7.

## Tasks

### [x] 1.0 Boundary-Aware Recursive Scanner

Replace the `filepath.Walk` implementation in `internal/discovery` with an explicit recursive traversal that prunes below repository roots, honors the skip list and hidden-directory rule, and stops at a configurable maximum depth. Symlink and worktree handling are deliberately deferred to task 2.0 so this task can be verified in isolation.

#### 1.0 Proof Artifact(s)

- Test: `internal/discovery/scanner_test.go` — a fixture tree with repositories nested under multiple container directories discovers every repository and registers nothing below a repository root, demonstrating strict boundary pruning
- Test: `internal/discovery/scanner_test.go` — a repository containing a nested clone registers only the outer repository, demonstrating that submodules and nested clones are invisible
- Test: `internal/discovery/scanner_test.go` — directories in the skip list and directories beginning with `.` are skipped at every level, demonstrating traversal filtering
- Test: `internal/discovery/scanner_test.go` — repositories at exactly `scan_max_depth` are registered while deeper ones are not, and a non-default preference value moves the boundary, demonstrating the configurable depth cap
- Test: `internal/config/config_test.go` — a `config.json` written before this change loads with `scan_max_depth` defaulting to `6`, demonstrating backward compatibility
- Command: `make lint` and `make test-race` pass, demonstrating the rewrite meets repository quality gates

#### 1.0 Tasks

- [x] 1.1 In `internal/config/config.go`, add `ScanMaxDepth int` with tag `json:"scan_max_depth,omitempty"` to the `Preferences` struct, and add `const DefaultScanMaxDepth = 6`.
- [x] 1.2 In `internal/config/config.go`, add a helper method on `*Config` that returns the effective scan depth — the preference value when set and greater than zero, otherwise `DefaultScanMaxDepth`. Follow the existing pattern used for `StatusWorkers` at `cmd/git-workspace/list.go:414`.
- [x] 1.3 In `internal/config/config_test.go`, add table-driven tests for the helper covering: nil `Preferences`, `Preferences` with the field unset, an explicit value, and a zero or negative value. Add a test that unmarshals a `config.json` fixture with no `preferences` key and confirms the default applies.
- [x] 1.4 In `internal/discovery/scanner.go`, define an `Options` struct with a `MaxDepth int` field, and change the exported entry point to `Scan(rootPath string, opts Options) (*ScanResult, error)`. Treat a zero `MaxDepth` as `config.DefaultScanMaxDepth` so callers cannot accidentally disable the cap.
- [x] 1.5 Replace the `filepath.Walk` body with a recursive directory walker that takes the current path and current depth. For each directory it must, in order: stop if depth exceeds `MaxDepth`; skip if the directory name is in the existing skip list; skip if the directory name begins with `.`; register the directory as a repository and return without descending when it contains a `.git` **directory**; otherwise read its entries and recurse into subdirectories.
- [x] 1.6 Export the repository-building logic as `BuildRepository(path string) (*config.Repository, error)` in `internal/discovery`, consolidating the near-identical `parseRepository` in `scanner.go` and `buildRepository` in `cmd/git-workspace/add.go` into one implementation. Delete both duplicates.
- [x] 1.7 Update all three `discovery.Scan` call sites for the new signature — `cmd/git-workspace/init.go:73`, `cmd/git-workspace/refresh.go:167`, and `cmd/git-workspace/add.go:164` — and point `add.go` at `discovery.BuildRepository`. Run `go test ./cmd/git-workspace/` and confirm the existing `add` tests still pass; if any fail because `add --recursive` now prunes or filters differently, record the behavior change in the test rather than weakening the scanner.
- [x] 1.8 Write tests in `internal/discovery/scanner_test.go` for boundary pruning: build a `t.TempDir()` fixture with container directories nesting repositories at several levels, plus a repository containing a nested clone, and assert the exact set of discovered paths.
- [x] 1.9 Write tests for traversal filtering: assert that a repository inside a `node_modules` directory and a repository inside a `.hidden` directory are both absent from results, at more than one nesting level.
- [x] 1.10 Write tests for the depth cap: build a fixture nesting repositories from depth 1 to depth 8, then assert that scanning with the default registers those at depth 6 and shallower, and that an explicit `MaxDepth` of 3 registers only those at depth 3 and shallower.
- [x] 1.11 Run `make lint` and `make test-race`; resolve any `gosec` G304 findings on the new file reads by cleaning paths and adding narrowly scoped suppressions with justifying comments.

### [ ] 2.0 Symlink Resolution and Structural Worktree Detection

Teach the scanner to follow symlinked directories safely and to distinguish worktrees from clones without relying on directory naming. This completes `internal/discovery` and is the task that fixes the curated-symlink-workspace failure on `gws init`.

#### 2.0 Proof Artifact(s)

- Test: `internal/discovery/scanner_test.go` — a curated workspace consisting only of symlinks to external repositories discovers all of them, demonstrating the failure mode this spec exists to fix
- Test: `internal/discovery/scanner_test.go` — a repository reachable both directly and through a symlink yields exactly one entry, asserted with the symlink visited first and with the physical path visited first, demonstrating order-independent deduplication by real path
- Test: `internal/discovery/scanner_test.go` — a symlink targeting the workspace root, and one targeting an ancestor of the current traversal path, complete without infinite recursion, demonstrating loop safety
- Test: `internal/git/worktree_test.go` — a repository with a linked worktree registers the repository and not the worktree directory, demonstrating structural worktree exclusion
- Test: `internal/git/worktree_test.go` — a genuine clone in a directory named `<name>.wt` is registered as a repository, demonstrating that no naming-convention filtering is applied
- Test: `internal/git/worktree_test.go` — a malformed or unreadable `.git` file triggers the `git rev-parse --git-dir` / `--git-common-dir` fallback, demonstrating the ambiguous-case branch
- Test: `internal/discovery/scanner_test.go` — broken and inaccessible symlinks are recorded as scan errors without aborting the scan, demonstrating error resilience
- Command: `make lint` and `make test-race` pass

#### 2.0 Tasks

- [ ] 2.1 In `internal/git/worktree.go`, add `IsLinkedWorktree(dir string) (isWorktree bool, mainRepoPath string, err error)`. Fast path: when `dir/.git` is a regular file, read it, parse the `gitdir: <path>` value, and treat it as a linked worktree when the target resolves to a path containing a `worktrees` path segment. Derive `mainRepoPath` from the target by walking up from `<main>/.git/worktrees/<name>` to `<main>`.
- [ ] 2.2 Add the ambiguous-case fallback to `IsLinkedWorktree`: when the `.git` file is unreadable, malformed, or its target does not match the expected shape, run `git rev-parse --git-dir` and `git rev-parse --git-common-dir` in `dir` using the `gitCommand` helper, and classify as a linked worktree only when the two values differ. Derive `mainRepoPath` from the common dir in that case.
- [ ] 2.3 Write tests in `internal/git/worktree_test.go` for `IsLinkedWorktree` using real temp repositories created with `git init` and `git worktree add`, covering: a normal clone (false), a linked worktree (true, with the correct main repository path), a clone in a directory named `<name>.wt` (false), and a `.git` file containing garbage (falls back and does not panic).
- [ ] 2.4 In `internal/discovery/scanner.go`, call `IsLinkedWorktree` when a candidate directory contains a `.git` file, and skip the directory without registering it when the result is true. Do not descend into it either way.
- [ ] 2.5 Add symlink resolution to the walker: for each directory entry, use `os.Lstat` to detect symlinks and `filepath.EvalSymlinks` to resolve them to a real path before deciding what to do with them. Record resolution failures as scan errors and continue.
- [ ] 2.6 Add a `visited map[string]bool` of resolved real directory paths to the scan state. Before traversing or registering any directory, check and populate it so no real directory is processed twice within a single scan. This is what makes deduplication independent of traversal order.
- [ ] 2.7 Add the two symlink ancestor guards: refuse to descend when the resolved target is an ancestor of, or equal to, the workspace root; and refuse to descend when the resolved target is an ancestor of the current traversal path. Implement ancestor checks on cleaned absolute paths with a separator-aware prefix comparison, matching the approach in `git.IsAligned`.
- [ ] 2.8 Set `Repository.Path` to the resolved real path for every registered repository, and add a `ReachablePaths map[string]bool` field to `ScanResult` recording every real repository path that was reachable from under the workspace root. Task 5.0 needs this for the repair-only symlink rule.
- [ ] 2.9 Write tests for the curated symlink workspace: a `t.TempDir()` workspace containing only symlinks to repositories created in a separate temp directory, asserting all are discovered with their real paths stored.
- [ ] 2.10 Write order-independence tests for deduplication: a workspace where the same repository appears both physically and via a symlink, run twice with entry names chosen so the symlink sorts first in one fixture and second in the other, asserting exactly one entry each time.
- [ ] 2.11 Write loop-safety tests: a symlink pointing at the workspace root, a symlink pointing at its own parent directory, and a pair of directories symlinked to each other. Each must complete and return without hanging — consider guarding these with a generous `context` timeout or a test-level deadline so a regression fails rather than stalls CI.
- [ ] 2.12 Write error-resilience tests: a dangling symlink and (where the platform permits) a directory with permissions removed both produce entries in `ScanResult.Errors` while the rest of the scan completes.
- [ ] 2.13 Run `make lint` and `make test-race`.

### [ ] 3.0 Shared Reconciliation Engine

Create `internal/reconcile` with the single entry point both commands will call. It performs discovery, merges against existing metadata, discovers worktrees, and returns a summary — with no cache clearing and no symlink writes. This task also introduces the shared scan-error reporter and the configurable progress label.

#### 3.0 Proof Artifact(s)

- Test: `internal/reconcile/reconcile_test.go` — called with a nil existing configuration, reconciliation returns the full discovered repository set with correct summary counts, demonstrating the `init` path
- Test: `internal/reconcile/merge_test.go` — called with an existing configuration, user-set tags survive on every repository whose path remains valid, demonstrating tag preservation
- Test: `internal/reconcile/merge_test.go` — new repositories are added, repositories with vanished paths are removed, and changes to name, remote URL, type, or visibility increment the updated count, demonstrating the merge rules
- Test: `internal/reconcile/worktree_test.go` — `Repository.Worktrees` is populated with correct `aligned` classification for a worktree inside `<repo>.wt/` and an unaligned worktree elsewhere, demonstrating worktree discovery and classification
- Test: `internal/reconcile/worktree_test.go` — a worktree whose main repository is untracked causes that main repository to be registered, including when it lives outside the workspace root, demonstrating orphan worktree handling
- Test: `internal/reconcile/reconcile_test.go` — reconciliation leaves the git status cache untouched and creates, repairs, or removes no symlinks, demonstrating separation of concerns
- Test: `internal/reconcile/reconcile_test.go` — the returned summary reports total, added, removed, updated, repositories-with-worktrees, total worktrees, aligned, and unaligned, demonstrating the complete result contract
- Test: `internal/reconcile/errors_test.go` — the shared scan-error reporter writes to stderr, lists at most 5 errors, and appends the correct `... and N more errors` line for a 10-error input, demonstrating unified error reporting
- Command: `make lint` and `make test-race` pass

#### 3.0 Tasks

- [ ] 3.1 Create `internal/reconcile/result.go` with a `Result` struct holding: `Repositories []config.Repository`, `TotalRepositories`, `Added`, `Removed`, `Updated`, `ReposWithWorktrees`, `TotalWorktrees`, `AlignedWorktrees`, `UnalignedWorktrees` counts, `Errors []error`, and `ReachablePaths map[string]bool` carried through from the scan.
- [ ] 3.2 Create `internal/reconcile/reconcile.go` with `ReconcileWorkspace(workspaceRoot string, existing *config.Config, opts Options) (*Result, error)`, where a nil `existing` represents the `init` case. Define `Options` to carry the effective scan depth and an optional progress reporter. Orchestrate the three phases in order: discovery, merge, worktree discovery.
- [ ] 3.3 Create `internal/reconcile/merge.go` implementing the merge rules against the discovered set: retain existing entries whose stored real path still exists, preserving `Tags`, `User`, `Email`, `SigningEnabled`, and `UserSource`; add discovered repositories not present in the existing metadata; drop entries whose paths are gone; and re-read name, remote URL, type, and visibility via `discovery.BuildRepository`, incrementing the updated count when any of those four changed.
- [ ] 3.4 Ensure merge results are deterministically ordered — sort the resulting repository slice by path — so `config.json` does not churn between runs and test assertions are stable.
- [ ] 3.5 Create `internal/reconcile/worktree.go` by moving `discoverWorktrees` out of `cmd/git-workspace/refresh.go:107`. Keep its existing behavior: `git worktree repair`, then `git worktree prune`, then `git worktree list --porcelain`, skipping entries whose paths no longer exist, and classifying each with `git.IsAligned`.
- [ ] 3.6 Extend the worktree phase to accumulate the aligned and unaligned counts and the total worktree count into `Result`, rather than only counting repositories that have worktrees.
- [ ] 3.7 Add orphan main-repository registration: when discovery reported a linked worktree whose main repository is not in the reconciled set, register that main repository via `discovery.BuildRepository` and include it in the results, then run the worktree phase for it. Have the scanner surface the main repository paths it detected in `ScanResult` so this does not require a second filesystem pass.
- [ ] 3.8 Create `internal/reconcile/errors.go` with `ReportScanErrors(w io.Writer, errs []error)`: print nothing when the slice is empty; otherwise write the `Warning: %d errors occurred during scanning:` header, list at most 5 errors as `  - %v`, and append `  ... and %d more errors` when the total exceeds 5.
- [ ] 3.9 In `internal/git/progress.go`, add a way to supply a label so the hardcoded `"Checking repos... %d/%d"` at line 65 is not shown during reconciliation. Keep `NewProgress` working unchanged for its existing caller at `cmd/git-workspace/list.go:504` — add a variant or an optional setter rather than changing the existing signature. Add a test in `internal/git/progress_test.go` covering both the default and custom labels.
- [ ] 3.10 Wire progress into `ReconcileWorkspace`: start the reporter before the per-repository phases, increment once per repository processed, and stop it before returning. Ensure `Stop` runs even on the error path.
- [ ] 3.11 Write `internal/reconcile/reconcile_test.go` covering the nil-existing-config path end to end against a `t.TempDir()` fixture, asserting the full `Result` contract including every count.
- [ ] 3.12 Write `internal/reconcile/merge_test.go` as table-driven tests for retain, add, remove, and update, with an explicit case asserting that tags survive a reconcile and another asserting that a changed remote URL increments `Updated`.
- [ ] 3.13 Write `internal/reconcile/worktree_test.go` covering alignment classification for both aligned and unaligned worktrees, and orphan main-repository registration for a main repository located outside the workspace root.
- [ ] 3.14 Write a test asserting reconciliation performs no side effects: snapshot the workspace directory listing and the status cache file before and after a reconcile and assert both are unchanged.
- [ ] 3.15 Write `internal/reconcile/errors_test.go` using a `bytes.Buffer`, covering zero errors, three errors, exactly five errors, and ten errors.
- [ ] 3.16 Run `make lint` and `make test-race`.

### [ ] 4.0 `gws init` Adoption

Rewrite `runInit` to create the metadata library and then delegate to the shared engine, producing a complete workspace model — including worktrees — from a single command. Adds the new summary lines, the reworded guard message, and progress reporting.

#### 4.0 Proof Artifact(s)

- Test: `cmd/git-workspace/init_test.go` — `runInit` against a fixture workspace containing nested containers, symlinked external repositories, and worktrees produces the exact summary lines specified in the spec with correct counts, demonstrating end-to-end `init` reconciliation
- Test: `cmd/git-workspace/init_test.go` — `init` on a curated symlink-only workspace discovers repositories where the previous implementation discovered none, demonstrating the headline behavior change
- Test: `cmd/git-workspace/init_test.go` — worktrees discovered during `init` are persisted to config and `runWorktreeList` shows them with no intervening command, demonstrating that `init` alone produces a complete model
- Test: `cmd/git-workspace/init_test.go` — the `Repositories with user configuration` and worktree summary lines are omitted when their counts are zero, demonstrating conditional output
- Test: `cmd/git-workspace/init_test.go` — with metadata already present, `runInit` writes the reworded guard message to stderr, returns nil, and modifies nothing, demonstrating protective behavior
- Test: `cmd/git-workspace/init_test.go` — `init` creates no workspace symlinks, demonstrating that symlink maintenance stayed with `refresh`
- Command: `make lint` and `make test-race` pass

#### 4.0 Tasks

- [ ] 4.1 Change `runInit` in `cmd/git-workspace/init.go` to accept an `io.Writer` for its summary output, following the `runWorktreeList(repoFilter string, stdout io.Writer)` pattern at `worktree_list.go:49`, so output can be asserted in tests. Update the Cobra `RunE` to pass `os.Stdout`.
- [ ] 4.2 Replace the `discovery.Scan` call and inline config assembly with: resolve the absolute workspace path, build the config via `config.New`, then call `reconcile.ReconcileWorkspace(absPath, nil, opts)` with the effective scan depth and a progress reporter.
- [ ] 4.3 Delete the inline scan-error printing block at `init.go:79-89` and call `reconcile.ReportScanErrors(os.Stderr, result.Errors)` instead.
- [ ] 4.4 Keep the existing post-reconcile steps in place and in the same order as `refresh` will use: `detectUserForRepos`, then `syncProfilesFromRepos`, then `config.Save`.
- [ ] 4.5 Implement the summary output exactly as specified: `Initialized workspace at: <path>`, `Found N repositories.`, then `Repositories with user configuration: N`, `Repositories with worktrees: N`, and `Worktrees: N (X aligned, Y unaligned)` — each conditional line omitted when its count is zero. Keep the existing `pluralize` helper.
- [ ] 4.6 Reword the already-initialized guard message so `gws refresh` is clearly the recommended next step, keeping it on stderr and keeping the `return nil` exit-zero behavior.
- [ ] 4.7 Add a test helper in `cmd/git-workspace/init_test.go` that builds a fixture workspace in `t.TempDir()` containing: two repositories at the root, one repository two container directories deep, one external repository reached by a symlink, one repository with an aligned worktree, and one with an unaligned worktree. Isolate `HOME` with `t.Setenv` so the real `~/.gws/config.json` is never touched. This helper is reused by task 5.0.
- [ ] 4.8 Write the summary-output test against that fixture, asserting the exact expected lines and counts.
- [ ] 4.9 Write the curated-symlink-workspace test asserting a non-zero repository count, which is the behavior that was previously broken.
- [ ] 4.10 Write a test that runs `runInit` and then `runWorktreeList` with no command in between, asserting both the aligned and unaligned worktrees appear.
- [ ] 4.11 Write tests for conditional line omission, for the reworded guard message on stderr with a nil return, and for the absence of any newly created symlinks in the workspace directory after `init`.
- [ ] 4.12 Run `make lint` and `make test-race`.

### [ ] 5.0 `gws refresh` Adoption and Symlink Repair Rules

Move `runRefresh` onto the shared engine while keeping the responsibilities that are genuinely refresh-specific: status cache clearing, the parent-directory rescan safety net, and workspace symlink maintenance — now constrained to repair only, which is what stops symlink accumulation across runs.

#### 5.0 Proof Artifact(s)

- Test: `cmd/git-workspace/refresh_test.go` — `runRefresh` discovers repositories added after initialization, including ones nested several container directories deep, demonstrating that the one-level-deep scan is gone
- Test: `cmd/git-workspace/refresh_test.go` — repositories whose stored paths no longer exist are removed and counted, demonstrating removal handling
- Test: `cmd/git-workspace/refresh_test.go` — a worktree created after the previous refresh appears in `runWorktreeList` immediately after one `runRefresh`, demonstrating immediate worktree visibility
- Test: `cmd/git-workspace/refresh_test.go` — two consecutive `runRefresh` calls against an unchanged curated symlink workspace leave the workspace directory unchanged and create no additional symlinks, demonstrating the symlink accumulation fix
- Test: `cmd/git-workspace/refresh_test.go` — a repository previously tracked with a workspace symlink that is no longer reachable under the workspace root has its symlink recreated, while a repository discovered by traversing an existing symlink gets no new link, demonstrating repair-only symlink behavior
- Test: `cmd/git-workspace/refresh_test.go` — the status cache is still cleared and the parent-directory rescan safety net still recovers a repository replaced in place, demonstrating retained refresh-specific behavior
- Test: `cmd/git-workspace/refresh_test.go` — `runInit` and `runRefresh` against the same fixture workspace report identical repository and worktree totals, demonstrating the core parity goal of this spec
- Test: `cmd/git-workspace/refresh_test.go` — `runInit` and `runRefresh` emit byte-identical warning output on stderr for the same scan errors, with stdout free of warnings, demonstrating unified scan-error reporting
- Command: `make lint` and `make test-race` pass

#### 5.0 Tasks

- [ ] 5.1 Change `runRefresh` in `cmd/git-workspace/refresh.go` to accept an `io.Writer` for summary output, matching the change made to `runInit`, and update the Cobra `RunE` to pass `os.Stdout`.
- [ ] 5.2 Replace `validateExistingRepos`, `scanWorkspaceForNewRepos`, and the inline worktree phase with a single `reconcile.ReconcileWorkspace(cfg.Workspace, cfg, opts)` call. Delete `scanWorkspaceForNewRepos` and the parts of `validateExistingRepos` now handled by the engine.
- [ ] 5.3 Preserve the parent-directory rescan safety net: for each repository the engine reports as removed, scan its parent directory with `discovery.Scan` and merge in any repositories found there that are not already tracked. Keep this in `refresh.go` as a post-reconcile step so the engine stays pure.
- [ ] 5.4 Replace the inline scan-error printing block at `refresh.go:49-59` with `reconcile.ReportScanErrors(os.Stderr, result.Errors)`. This intentionally changes two behaviors: the preview limit rises from 3 to 5, and warnings move from stdout to stderr.
- [ ] 5.5 Rewrite `ensureWorkspaceSymlink` to implement the repair-only rule: create a link only when the repository's real path is absent from `result.ReachablePaths` **and** the repository was previously tracked with a workspace symlink. Determine "previously tracked with a symlink" by checking the pre-reconcile config for the repository together with the presence of a matching symlink entry, and pass that set in rather than probing the filesystem twice.
- [ ] 5.6 Leave `removeWorkspaceSymlink` behavior unchanged — remove only when the entry is a symlink whose target matches the removed repository's path, never touching non-symlink or differently-targeted entries.
- [ ] 5.7 Keep the status cache clear and its `Cleared git status cache` message exactly as they are today, after reconciliation and before the summary.
- [ ] 5.8 Implement the summary output as specified, keeping today's wording — `Refreshing workspace at:`, `Detecting git user configuration...`, `Refresh complete!`, `Total repositories:`, and the conditional removed, new, updated, user-configuration, and worktree lines — with the added `Worktrees: N (X aligned, Y unaligned)` line, and all counts sourced from `Result`.
- [ ] 5.9 Write tests for recursive rescan and removals using the shared fixture helper from task 4.7: add a repository three container directories deep after an initial reconcile and assert it is found; delete a tracked repository and assert it is removed and counted.
- [ ] 5.10 Write the immediate-worktree-visibility test: create a worktree with `git worktree add` after an initial reconcile, run `runRefresh`, then `runWorktreeList`, and assert the new worktree appears.
- [ ] 5.11 Write the idempotence test: run `runRefresh` twice against an unchanged curated symlink workspace and assert the workspace directory listing is identical between runs, with no new symlinks and no change to the set of stored repository paths.
- [ ] 5.12 Write the repair-only symlink tests: one case where a previously symlinked external repository has had its link deleted and refresh recreates it, and one case where a repository discovered by traversing an existing nested symlink gets no additional link at the workspace root.
- [ ] 5.13 Write tests confirming the retained refresh-specific behavior: the status cache is cleared, and a repository replaced in place at a removed repository's parent directory is picked up by the safety net.
- [ ] 5.14 Write the two parity tests: identical repository and worktree totals from `runInit` and `runRefresh` against the same fixture, and byte-identical stderr warning output from both for the same injected scan errors, with stdout containing no warnings.
- [ ] 5.15 Run `make lint` and `make test-race`, then run the full `go test ./...` to confirm no regressions in `list`, `navigate`, `tag`, `user`, or the worktree subcommands.

### [ ] 6.0 Documentation Sync

Bring the MkDocs site and README into line with the new behavior, per the repository standard established by spec 18.

#### 6.0 Proof Artifact(s)

- Diff: `docs/site/commands-core.md` — the `Initialize Workspace` and `Refresh Workspace` sections describe recursive boundary-aware discovery, symlink handling, worktree discovery, and `scan_max_depth`, demonstrating accurate command documentation
- Diff: `docs/site/commands-core.md` — example output blocks match the actual output of `gws init` and `gws refresh`, demonstrating that documented examples are real
- Diff: `docs/site/configuration.md` — `scan_max_depth` is documented with its default of `6` and its effect on traversal, demonstrating the new preference is discoverable
- Diff: `README.md` — descriptions of `init` and `refresh` match the docs site, demonstrating README/site consistency
- Command: `make docs` builds the site without errors or broken links, demonstrating the documentation is well-formed

#### 6.0 Tasks

- [ ] 6.1 Update the `Initialize Workspace` section of `docs/site/commands-core.md` (around line 275) to describe recursive discovery through container directories, pruning at repository boundaries, symlink resolution, worktree discovery, and the fact that `init` does not create workspace symlinks.
- [ ] 6.2 Update the `Refresh Workspace` section (around line 342) to state that `refresh` applies the same discovery rules as `init`, and to describe what it additionally does: clears the status cache and repairs workspace symlinks.
- [ ] 6.3 Replace the example output blocks in both sections with real captured output from a fixture workspace, matching the summary formats implemented in tasks 4.5 and 5.8.
- [ ] 6.4 Document the traversal rules that users can actually observe: the skip list, hidden directories being skipped, and the depth cap.
- [ ] 6.5 Add a `scan_max_depth` entry to `docs/site/configuration.md` covering its default of `6`, how depth is counted relative to the workspace root, and when a user would raise it.
- [ ] 6.6 Update the `Repository Classification` note at the top of `commands-core.md` if it describes discovery in terms that the new scanner contradicts.
- [ ] 6.7 Update `README.md` wherever it describes `init` or `refresh`, keeping the wording consistent with the docs site.
- [ ] 6.8 Run `make docs` and confirm the site builds cleanly, then review the rendered `Initialize Workspace`, `Refresh Workspace`, and configuration pages for broken formatting or stale links.
