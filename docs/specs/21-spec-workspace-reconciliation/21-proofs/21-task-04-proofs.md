# Spec 21 — Task 4.0 Proofs: `gws init` Adoption

**Task:** 4.0 `gws init` Adoption
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

| File | Change |
|---|---|
| `cmd/git-workspace/init.go` | `runInit` delegates to `ReconcileWorkspace`; accepts an `io.Writer`; new summary lines; reworded guard; shared error reporter |
| `cmd/git-workspace/init_test.go` | Updated 3 existing tests for the new signature; added 7 new tests |
| `cmd/git-workspace/reconcilefixture_test.go` | New — shared workspace fixture used by tasks 4.0 and 5.0 |
| `cmd/git-workspace/deprecated.go` | Updated the deprecated `--init` call site |
| `internal/git/progress.go` | Added `SetTotal`, since the repository count is not known until after the merge |
| `internal/reconcile/reconcile.go` | `ProgressReporter` gained `SetTotal`; the engine calls it before `Start` |

## CLI Output — Real Binary

A workspace was built on disk with a top-level repository, a repository nested two container directories deep, an external repository reached by symlink, and a repository owning one aligned and one unaligned worktree.

```text
$ gws init /tmp/…/demo/ws
Initialized workspace at: /tmp/…/demo/ws
Found 3 repositories.
Repositories with user configuration: 3
Repositories with worktrees: 1
Worktrees: 2 (1 aligned, 1 unaligned)
```

### Worktrees are visible with no intervening command

This is the property the spec exists to guarantee — `gws worktree list` run directly after `gws init`, with no `gws refresh` between them:

```text
$ gws worktree list
REPO         BRANCH   PATH                                              STATUS
-----------  -------  ------------------------------------------------  ------
service-api  hotfix   /tmp/…/demo/ws/loose/hotfix                       (unaligned)
service-api  feature  /tmp/…/demo/ws/org-a/team-one/service-api.wt/feature  aligned
```

### Protective guard, reworded

```text
$ gws init /tmp/…/demo/ws
Workspace already initialized at: /tmp/…/demo/ws

To re-scan the workspace and pick up changes:  gws refresh
To add a single repository:                    gws add
exit code: 0
```

`gws refresh` is now listed first as the recommended next step, per round-2 Q9-C. Exit code remains 0, and the message stays on stderr.

## Test Results

```text
$ go test ./cmd/git-workspace/ -run TestRunInit -v
--- PASS: TestRunInit_HappyPath (0.00s)
--- PASS: TestRunInit_AlreadyInitialized (0.00s)
--- PASS: TestRunInit_EmptyDirectory (0.00s)
--- PASS: TestRunInit_FullReconciliationSummary (0.08s)
--- PASS: TestRunInit_CuratedSymlinkWorkspace (0.03s)
--- PASS: TestRunInit_WorktreesVisibleImmediately (0.07s)
--- PASS: TestRunInit_OmitsZeroCountLines (0.01s)
--- PASS: TestRunInit_GuardMessageRecommendsRefresh (0.01s)
--- PASS: TestRunInit_CreatesNoSymlinks (0.07s)
--- PASS: TestRunInit_ReportsScanErrorsOnStderr (0.01s)
ok  	github.com/daileyo/gws/cmd/git-workspace	0.301s
```

### Full suite under the race detector

```text
$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace	1.765s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	1.488s
ok  	github.com/daileyo/gws/internal/reconcile	1.662s
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Quality Gates

```text
$ make lint
Running linter...
(exit 0, no findings)
```

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| Fixture workspace with nested containers, symlinked externals, and worktrees produces the exact summary lines | `TestRunInit_FullReconciliationSummary`; CLI output above | ✅ PASS |
| `init` on a curated symlink-only workspace discovers repositories where the previous implementation found none | `TestRunInit_CuratedSymlinkWorkspace` (3 symlinked repos, real paths stored) | ✅ PASS |
| Worktrees persisted by `init` and visible in `runWorktreeList` with no intervening command | `TestRunInit_WorktreesVisibleImmediately`; CLI transcript above | ✅ PASS |
| Conditional lines omitted when counts are zero | `TestRunInit_OmitsZeroCountLines` (also pins the singular "1 repository") | ✅ PASS |
| Guard writes reworded message to stderr, returns nil, modifies nothing | `TestRunInit_GuardMessageRecommendsRefresh` (byte-compares config before/after; asserts refresh precedes add) | ✅ PASS |
| `init` creates no workspace symlinks | `TestRunInit_CreatesNoSymlinks` (symlink listing compared before/after) | ✅ PASS |
| `make lint` and `make test-race` pass | Quality Gates section above | ✅ PASS |

## Design Notes

- **`SetTotal` added to the progress contract.** `git.NewProgress` takes its total at construction, but the repository count is not known until the merge finishes. Rather than have the command guess, `ProgressReporter` gained `SetTotal`, which the engine calls immediately before `Start`. Setting it before the rendering goroutine exists keeps the read ordered behind goroutine creation, so there is no race — confirmed by `make test-race`.
- **`writeWorktreeSummary` is shared.** The worktree lines live in one helper that task 5.0 will also call, so `init` and `refresh` cannot describe the same workspace differently.
- **Post-reconcile order is fixed** as `detectUserForRepos` → `syncProfilesFromRepos` → `config.Save`, matching what `refresh` will do in task 5.0.
- **`runInit` now takes an `io.Writer`**, following the `runWorktreeList(repoFilter, stdout)` pattern already in the codebase, which is what makes the summary assertions possible.

## Shared Fixture for Task 5.0

`newWorkspaceFixture` (in `reconcilefixture_test.go`) builds the workspace both this task and task 5.0 use, so the parity assertions in 5.14 compare two commands against a byte-identical starting state. It redirects `HOME` to a temp directory via `t.Setenv`, so no test reads or writes the developer's real `~/.gws/config.json`.

Fixture contents: 2 top-level repositories, 1 repository nested two containers deep, 1 external repository reached by symlink, and 1 repository owning an aligned worktree under `<repo>.wt/` plus an unaligned worktree elsewhere — 5 repositories and 2 worktrees.
