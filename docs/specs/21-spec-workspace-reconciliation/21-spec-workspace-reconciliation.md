# 21-spec-workspace-reconciliation

## Introduction/Overview

`gws init` and `gws refresh` currently discover repositories using two different, drifted implementations. `init` walks the directory tree recursively but never follows symlinks and never discovers worktrees; `refresh` follows symlinks but only reads one directory level deep. The result is that a user's workspace model depends on which command they happened to run last, and a curated symlink-based workspace produces zero repositories on `init`.

This spec extracts a single reconciliation engine used by both commands. The two commands keep their distinct lifecycle roles — `init` creates the metadata library then reconciles it, `refresh` reconciles the existing library — but they share one set of discovery and merge rules, so a user gets the same repository and worktree model whichever command they run.

## Goals

- Replace the two divergent discovery paths with one shared `ReconcileWorkspace` function called by both `gws init` and `gws refresh`.
- Make repository discovery repository-boundary aware: recurse through arbitrary container directories, register a repository at its root, and stop descending below it.
- Discover native git worktrees for every tracked repository during both commands, classified as aligned or unaligned, so `gws worktree list` is correct immediately after either command.
- Follow workspace symlinks safely at any depth, deduplicating by resolved real path and refusing to traverse loops.
- Report repository and worktree totals consistently from both commands, preserving user-set tags on repositories whose tracked path remains valid.

## User Stories

- **As a developer setting up gws for the first time**, I want `gws init` to find every repository beneath my workspace root — including ones reached through symlinks and ones nested inside grouping directories — so that my workspace is complete without running a second command.
- **As a developer who reorganizes checkouts**, I want `gws refresh` to re-scan my workspace with the same rules `init` used so that flattening or nesting my directory hierarchy does not silently drop repositories from gws.
- **As a developer using git worktrees**, I want newly created worktrees to appear in `gws worktree list` after a single `gws refresh` so that I can navigate to them without wondering whether gws has noticed them yet.
- **As a developer who deleted an old checkout**, I want `gws refresh` to remove that repository from my metadata so that `gws list` reflects what actually exists on disk.
- **As a developer with a curated symlink workspace**, I want gws to resolve my workspace symlinks without creating new duplicate symlinks every time I refresh so that my carefully arranged workspace directory stays as I built it.
- **As a developer who has tagged repositories**, I want my tags to survive every refresh so that my organizational work is not lost to a routine maintenance command.

## Demoable Units of Work

### Unit 1: Boundary-Aware Discovery Scanner

**Purpose:** Rebuild `internal/discovery` so a single scan of the workspace root produces the correct set of repositories, regardless of nesting, symlinks, or worktrees. This is the foundation both commands depend on and is fully testable in isolation.

**Functional Requirements:**

*Traversal and boundaries*

- The system shall traverse the workspace root recursively, descending through arbitrary organizational and container directories.
- The system shall register a directory as a repository when it contains a `.git` **directory**, and shall not traverse any path below that repository (strict pruning). Submodules and nested clones inside a registered repository are intentionally invisible to gws.
- The system shall continue scanning sibling and ancestor container directories after pruning below a repository.
- The system shall skip directories named in the existing skip list (`node_modules`, `vendor`, `.venv`, `venv`, `__pycache__`, `.tox`, `target`, `build`, `dist`, `.cache`, `.terraform`).
- The system shall skip any directory whose name begins with `.` at every level below the workspace root. A repository placed at `<workspace>/.dotfiles` will not be discovered.
- The system shall stop descending past a maximum depth, counted in directory levels below the workspace root. A repository found exactly at the maximum depth is registered; nothing below it is traversed.
- The maximum depth shall be read from `config.Preferences` under the key `scan_max_depth`, defaulting to `6` when unset.

*Worktree detection*

- The system shall not register a git worktree as a repository. Worktrees are discovered only in the worktree phase of reconciliation (Unit 2).
- When a candidate directory contains a `.git` **file**, the system shall read the file's `gitdir:` target and treat the directory as a linked worktree when that target resolves to a path under a `worktrees/` directory. This is the fast path and performs filesystem reads only.
- When the `gitdir:` target is ambiguous — unreadable, malformed, or not matching the expected `worktrees/` shape — the system shall fall back to running `git rev-parse --git-dir` and `git rev-parse --git-common-dir` in that directory, and shall classify it as a linked worktree only when the two values differ.
- The system shall determine worktree status structurally and shall never use directory naming (such as the `.wt` suffix) to decide whether something is a repository. A genuine clone in a directory named `important-repo.wt` shall be registered as a normal repository.

*Symlinks and identity*

- The system shall follow symlinked directories encountered anywhere in the traversal, resolving each to its real path.
- The system shall record `Repository.Path` as the resolved real path of the repository.
- The system shall deduplicate discovered repositories by resolved real path, so a repository reachable through multiple workspace paths is registered exactly once. Deduplication shall be independent of traversal order: whether the symlinked path or the physical path is visited first, the result is a single entry.
- The system shall maintain a visited set of resolved real directory paths and shall never traverse the same real directory twice within a single scan.
- The system shall refuse to descend into a symlink whose target is an ancestor of, or equal to, the workspace root.
- The system shall refuse to descend into a symlink whose target is an ancestor of the current traversal path.
- The system shall record broken or inaccessible symlinks as scan errors and continue scanning.

**Proof Artifacts:**

- Test: `internal/discovery` scan of a fixture tree with nested container directories discovers all repositories and registers none of the paths below a repository root — demonstrates strict boundary pruning.
- Test: scan of a fixture repository with a linked worktree registers the repository once and does not register the worktree directory — demonstrates structural worktree exclusion.
- Test: scan of a fixture where a real clone lives in a directory named `<name>.wt` registers it as a repository — demonstrates no naming-convention filtering.
- Test: scan of a fixture with a malformed `.git` file exercises the `git rev-parse` fallback path — demonstrates the ambiguous-case branch.
- Test: scan of a fixture where the same repository is reachable both directly and through a symlink yields one entry, asserted in both traversal orders — demonstrates order-independent deduplication by real path.
- Test: scan of a fixture containing a symlink pointing to the workspace root, and one pointing to an ancestor of the current path, completes without infinite recursion — demonstrates loop safety.
- Test: scan of a fixture nested deeper than `scan_max_depth` registers repositories at the limit and none below it; a non-default preference value changes the boundary — demonstrates the configurable depth cap.
- Test: scan skips hidden directories and skip-list directories at every level — demonstrates traversal filtering.

### Unit 2: Shared Reconciliation Engine and `gws init` Adoption

**Purpose:** Create the single reconciliation function both commands call, and make `gws init` the first consumer. After this unit, a fresh `gws init` produces a complete workspace model including worktrees.

**Functional Requirements:**

*The reconciliation engine*

- The system shall provide a new `internal/reconcile` package exposing a reconciliation entry point conceptually equivalent to `ReconcileWorkspace(workspaceRoot string, existing *config.Config)`, where a nil or empty existing configuration represents the `init` case.
- `internal/discovery` shall be reduced to the filesystem scan it owns; the merge logic shall live in `internal/reconcile`.
- The reconciliation function shall perform, in order: (1) discovery via Unit 1's scanner, (2) merge of discovered repositories with existing metadata, (3) worktree discovery for every resulting repository, and shall return the reconciled repository set together with a result summary.
- The reconciliation function shall not clear the git status cache and shall not create, repair, or remove workspace symlinks. Those responsibilities remain with the calling command.

*Merge rules*

- The system shall retain existing metadata for repositories whose stored real path still exists, including user-set tags and user/profile configuration.
- The system shall preserve user-set tags on every retained repository.
- The system shall add repositories that were discovered but are not present in the existing metadata.
- The system shall remove entries whose stored path no longer exists on disk.
- The system shall re-read mutable metadata for retained repositories — name, remote URL, type, and visibility — and count an entry as updated when any of those values changed.

*Worktree phase*

- The system shall discover worktrees for every tracked repository using `git worktree list --porcelain`.
- The system shall add newly created worktrees and remove worktree entries whose paths no longer exist.
- The system shall mark a worktree as `aligned` when its path is inside `<repo>.wt/`, and as unaligned otherwise.
- When discovery encounters a linked worktree whose main repository is not already tracked, the system shall resolve the main repository from the worktree's `gitdir:` target and register that main repository as a tracked repository, after which its worktrees are discovered normally. This applies whether or not the main repository lives under the workspace root.

*Result summary*

- The reconciliation function shall return counts for: total repositories, repositories added, repositories removed, repositories updated, repositories with worktrees, total worktrees, aligned worktrees, and unaligned worktrees.
- The reconciliation function shall return the collected scan errors rather than printing them.

*Scan error reporting*

- The system shall provide a single shared scan-error reporter used by both commands, so identical scan errors are presented identically regardless of which command produced them.
- The reporter shall write to **stderr**, print a `Warning: %d errors occurred during scanning:` header, list at most **5** individual errors, and append `... and %d more errors` when the total exceeds 5.
- `gws init` shall report scan errors exclusively through this reporter, replacing its inline error-printing block.

*Progress reporting*

- The reconciliation function shall accept a progress reporting mechanism and emit progress during the per-repository phases, using the existing `internal/git/progress.go` facility.
- Both `gws init` and `gws refresh` shall display progress during reconciliation, unconditionally.

*`gws init` behavior*

- `gws init [directory]` shall create the metadata library, record the resolved absolute workspace root, call the shared reconciliation function with no existing configuration, persist the result, and report the summary.
- `gws init` shall continue to default to the current directory when no argument is given.
- `gws init` shall not create workspace symlinks.
- When metadata already exists, `gws init` shall print a message to stderr, exit with status 0, and word that message so that `gws refresh` is clearly the recommended next step.
- `gws init` shall print a summary of this shape:

  ```text
  Initialized workspace at: /home/daileyo/gws
  Found 8 repositories.
  Repositories with user configuration: 6
  Repositories with worktrees: 2
  Worktrees: 3 (1 aligned, 2 unaligned)
  ```

- The `Repositories with user configuration` line shall continue to appear only when the count is greater than zero; the worktree lines shall likewise be omitted when no worktrees were found.

**Proof Artifacts:**

- Test: `ReconcileWorkspace` called with a nil existing configuration against a fixture workspace returns the full discovered repository set and correct summary counts — demonstrates the `init` path.
- Test: `ReconcileWorkspace` called with an existing configuration retains tags on repositories whose paths remain valid — demonstrates tag preservation.
- Test: reconciliation adds new repositories, removes repositories whose paths are gone, and counts metadata updates for changed remote URL, type, or visibility — demonstrates the merge rules.
- Test: reconciliation populates `Repository.Worktrees` with correct `aligned` classification for a worktree inside `<repo>.wt/` and an unaligned worktree elsewhere — demonstrates worktree discovery and classification.
- Test: reconciliation registers the main repository of an otherwise-untracked worktree discovered under the workspace root — demonstrates orphan worktree handling.
- Test: reconciliation neither clears the status cache nor touches workspace symlinks — demonstrates separation of concerns.
- Test: `runInit` against a fixture workspace produces the specified summary lines with correct counts — demonstrates `init` reporting.
- Test: `runInit` with existing metadata prints the reworded guard message to stderr and returns nil — demonstrates protective behavior.
- Test: the shared scan-error reporter writes to stderr, lists at most 5 errors, and appends the correct `... and N more errors` line for a 10-error input — demonstrates unified error reporting.

### Unit 3: `gws refresh` Adoption, Symlink Rules, and Documentation

**Purpose:** Move `gws refresh` onto the shared engine, keeping the responsibilities that are genuinely refresh-specific, and fix the symlink accumulation that deep traversal would otherwise cause. Documentation is updated so the site matches behavior.

**Functional Requirements:**

*`gws refresh` behavior*

- `gws refresh` shall load the existing configuration, call the shared reconciliation function with that configuration and the recorded workspace root, persist the result, and report the summary.
- `gws refresh` shall continue to clear the git status cache after reconciliation, retaining today's behavior and its `Cleared git status cache` message.
- `gws refresh` shall retain the parent-directory rescan safety net for repositories whose stored path has disappeared.
- `gws refresh` shall print a summary of this shape:

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
  Worktrees: 3 (1 aligned, 2 unaligned)
  ```

- Conditional lines (removed, new, updated, user configuration, worktrees) shall continue to be omitted when their count is zero.
- Repository and worktree totals shall be derived from the same reconciliation result in both commands, so the two commands never disagree about the same workspace.
- `gws refresh` shall report scan errors exclusively through the shared reporter introduced in Unit 2, replacing its inline error-printing block. This changes two existing behaviors: the preview limit rises from 3 to 5, and the warnings move from stdout to stderr.

*Workspace symlink rules*

- `gws refresh` shall create a workspace symlink for a repository only when that repository has no reachable path under the workspace root **and** was previously tracked with a workspace symlink. This is a repair operation, not a discovery-driven creation.
- `gws refresh` shall not create a workspace symlink for a repository that was discovered by traversing an existing symlink, regardless of whether its real path lies outside the workspace root.
- `gws refresh` shall continue to remove a workspace symlink only when it is a symlink whose target matches the removed repository's path, leaving non-symlink and differently-targeted entries untouched.
- Running `gws refresh` repeatedly against an unchanged workspace shall produce no new symlinks and no changes to the workspace directory.

*Documentation*

- The `Initialize Workspace` and `Refresh Workspace` sections of `docs/site/commands-core.md` shall be updated to describe recursive, boundary-aware discovery, symlink handling, worktree discovery, and the `scan_max_depth` preference.
- The example output blocks in those sections shall match the actual command output.
- `docs/site/configuration.md` shall document the `scan_max_depth` preference, its default of `6`, and its effect.
- `README.md` shall be updated where it describes `init` or `refresh` behavior, keeping it in sync with the docs site.

**Proof Artifacts:**

- Test: `runRefresh` against a fixture workspace discovers repositories added after initialization, including repositories nested several container directories deep — demonstrates recursive re-scan.
- Test: `runRefresh` removes repositories whose stored paths no longer exist and reports the removal count — demonstrates removal handling.
- Test: `runRefresh` discovers a worktree created after the previous refresh, and `runWorktreeList` shows it without any intervening command — demonstrates immediate worktree visibility.
- Test: `runInit` and `runRefresh` against the same fixture workspace report identical repository and worktree totals — demonstrates consistent reporting across commands.
- Test: a curated workspace of external-repository symlinks yields the correct repositories on both `init` and `refresh` — demonstrates the curated-workspace model still works.
- Test: two consecutive `runRefresh` calls against an unchanged curated workspace leave the workspace directory byte-identical, creating no additional symlinks — demonstrates the symlink accumulation fix.
- Test: `runRefresh` recreates a missing workspace symlink for a repository that was previously tracked with one and is no longer reachable under the workspace root — demonstrates repair-only symlink behavior.
- Test: `runInit` and `runRefresh` given the same set of scan errors emit byte-identical warning output on stderr, with stdout free of warnings — demonstrates unified scan-error reporting across both commands.

## Non-Goals (Out of Scope)

1. **Metadata recovery**: recovering from deleted or corrupted `~/.gws/config.json` is explicitly out of scope. Reinitialization, migration, corruption recovery, and metadata reset remain candidates for separate future commands.
2. **A `--force` flag for `init`**: `init` stays protective and does not gain a way to discard and recreate existing metadata.
3. **Submodule tracking**: strict pruning means submodules and clones nested inside a registered repository are never registered as separate repositories.
4. **A `branch` field on `Repository`**: branch information remains per-worktree and in the transient status cache. No config schema change to `Repository`.
5. **Changing repository name-collision behavior**: two distinct clones sharing a directory name continue to produce two entries resolved by the existing interactive selector. No name qualification is introduced.
6. **Configurable hidden-directory exceptions**: hidden directories are skipped at every level. Making that behavior configurable is future work.
7. **Changing status cache semantics**: the cache continues to be cleared by `refresh` and untouched by `init`. No in-place cache updating.
8. **New output formats**: no `--json` or machine-readable summary output is added.
9. **Changing `gws add` or `gws worktree` subcommands**: these consume the reconciled metadata but their own behavior is unchanged.

## Design Considerations

No graphical interface is involved. The command-line output requirements are specified verbatim in Units 2 and 3, and follow the existing convention in this codebase: a header line, conditional detail lines omitted when their count is zero, warnings and guard messages to stderr, and results to stdout. Progress reporting reuses the existing `internal/git/progress.go` facility rather than introducing a new visual style.

## Repository Standards

- **Language and layout**: Go, with commands in `cmd/git-workspace/` (one file per command, `<name>_test.go` alongside) and reusable logic in `internal/<package>/`. The new `internal/reconcile` package follows that pattern.
- **CLI framework**: Cobra. Commands define a `var <name>Cmd = &cobra.Command{...}` with `Use`, `Short`, a `Long` description ending in an `Examples:` block, and a thin `RunE` delegating to a testable `run<Name>` function.
- **Testing**: standard library `testing` with table-driven tests, run via `make test` and `make test-race`. Fixture workspaces are built in `t.TempDir()`. Tests exercise the `run<Name>` functions directly with an injected `io.Writer` where output is asserted, as `runWorktreeList` already does.
- **Linting**: `golangci-lint` per `.golangci.yml` — `errcheck`, `gosec`, `unparam`, `govet` with `enable-all` (minus `fieldalignment`), `goimports` with local prefix `github.com/daileyo/gws`. Run through `make lint`; `make setup-hooks` installs a pre-push hook that runs vet, lint, and tests.
- **Error handling**: wrap errors with `fmt.Errorf("...: %w", err)`; intentionally ignored errors use the blank identifier, which `errcheck` is configured to allow.
- **Commits**: conventional commits, with the `commit-msg` hook padding the type for aligned git log output. Release automation is release-please, so commit types drive versioning.
- **Documentation**: per spec 18, `docs/site/` and `README.md` must stay in sync with CLI behavior; `make docs` builds the MkDocs site.

## Technical Considerations

- **Shared entry point**: `internal/reconcile` exposes the single reconciliation function. `internal/discovery` keeps ownership of the filesystem scan only. `init` calls reconcile with no existing configuration; `refresh` calls it with the loaded configuration.
- **Replacing `filepath.Walk`**: `filepath.Walk` cannot express "stop below this repository" cleanly and does not follow symlinks. The scanner needs an explicit recursive walk with per-directory decisions, a visited real-path set, and a depth counter. `filepath.WalkDir` is an option for the non-symlink portions but the symlink resolution and pruning logic must be hand-rolled either way.
- **Real-path canonicalization**: every candidate directory is resolved with `filepath.EvalSymlinks` before being used as a deduplication key, so aliasing is resolved regardless of traversal order.
- **Subprocess cost**: the `git rev-parse` fallback runs only for ambiguous `.git` files, keeping the common case free of subprocesses. The worktree phase already runs `git worktree repair`, `git worktree prune`, and `git worktree list --porcelain` per repository; that cost now applies to `init` as well as `refresh`, which is why progress reporting is unconditional.
- **Configuration compatibility**: `scan_max_depth` is a new optional field on the existing `config.Preferences` struct, which is already `omitempty` and nil-able. Adding it is additive and backward compatible, so existing `config.json` files load unchanged. No `Repository` schema change is required. Assumption: `ConfigVersion` does not need a bump for an additive optional preference — flag this during implementation if the project prefers otherwise.
- **Existing helpers to reuse**: `buildRepository` (remote/type/visibility), `classifier.Classify`, `git.ListWorktrees`, `git.IsAligned`, `git.RepairWorktrees`, `git.PruneWorktrees`, `detectUserForRepos`, and `syncProfilesFromRepos`. Helpers currently living in `cmd/git-workspace` that reconcile needs — notably `buildRepository` — must move into an `internal` package to be importable.
- **Symlink reachability check**: the repair-only symlink rule needs a "is this real path reachable from under the workspace root" predicate. The scanner already computes that set during discovery, so it should be returned in the reconcile result rather than recomputed by `refresh`.
- **Orphan worktree main repositories** may live outside the workspace root. They are tracked by real path and, per the repair-only rule, receive no workspace symlink.

## Security Considerations

- **Symlink traversal escape**: following symlinks at arbitrary depth means a symlink under the workspace root can point anywhere on the filesystem, including system directories. The ancestor checks, the visited real-path set, and the `scan_max_depth` cap together bound traversal. Discovery only reads directory entries and `.git` files; it never writes outside the workspace or the metadata library.
- **Symlink creation**: the repair-only rule limits symlink writes to a narrow, previously-established case, reducing the chance of gws creating links in unexpected locations.
- **`gosec` compliance**: the new scanner performs file reads on paths derived from directory traversal, which `gosec` flags (G304). Paths must be cleaned and constrained to the traversal root, and any necessary suppressions must be narrowly scoped with a justifying comment, consistent with existing code.
- **Metadata permissions**: `~/.gws/config.json` continues to be written with `0600`. It contains git user names and email addresses, so nothing about this change may loosen those permissions or copy that content elsewhere.
- **Proof artifacts**: proofs are Go unit tests using `t.TempDir()` fixtures. No real `config.json`, no real repository paths, and no git identities from the developer's machine may be committed to `21-proofs/`.

## Success Metrics

1. **Behavioral parity**: `gws init` and `gws refresh` run against the same workspace produce identical repository sets, identical worktree sets, and identical totals — verified by a test asserting equality of both results.
2. **Discovery completeness**: every acceptance criterion from the original request passes as an automated test — recursive rescan, post-init repository discovery, removal of vanished repositories, immediate worktree discovery, `gws worktree list` correctness after either command, curated symlink workspaces, consistent totals, and tag preservation.
3. **Idempotence**: two consecutive `gws refresh` runs against an unchanged workspace produce an identical `config.json` and an unchanged workspace directory, with no new symlinks.
4. **No regressions**: `make test-race` and `make lint` pass, and the existing `init`, `refresh`, and worktree test suites pass unmodified except where a requirement in this spec deliberately changes behavior.
5. **Documentation accuracy**: every example output block in the `Initialize Workspace` and `Refresh Workspace` documentation sections matches actual command output.

## Open Questions

Both questions were resolved by the spec owner on 2026-08-10. No open questions remain.

1. **Should the `Worktrees: N (X aligned, Y unaligned)` line appear in `gws refresh` output as well as `gws init`?**
   **Resolved: yes, both commands.** Parity is the point of this spec, and the original acceptance criteria require repository and worktree totals to be reported consistently by both commands. Both commands emit the line through the shared `writeWorktreeSummary` helper, so they cannot drift. Covered by `TestInitAndRefresh_ReportIdenticalTotals`.

2. **Does adding the optional `scan_max_depth` preference warrant a `ConfigVersion` bump from `1.1.0`?**
   **Resolved: yes, bump to `1.2.0`.** Although the change is additive and older files load unchanged, the stored format did gain a field, and the version should say so.

   Two consequences follow, both intentional:
   - `Save` now stamps `ConfigVersion` on every write, so a file always declares the format actually written to it. Without this, a workspace initialized under `1.1.0` would keep reporting `1.1.0` forever even after newer code rewrote it. Nothing branches on the value today; it exists to make future migrations possible.
   - No migration step is required or provided. A pre-`1.2.0` file loads unchanged, `scan_max_depth` defaults to `6`, and the version is upgraded the next time the file is written.

   Covered by `TestConfigVersion` (4 sub-cases, including an in-place upgrade that asserts repositories, tags, and workspace survive) and `TestScanMaxDepthBackwardCompatibility`.
