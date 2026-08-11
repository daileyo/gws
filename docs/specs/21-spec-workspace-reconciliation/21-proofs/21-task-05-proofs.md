# Spec 21 — Task 5.0 Proofs: `gws refresh` Adoption and Symlink Repair Rules

**Task:** 5.0 `gws refresh` Adoption and Symlink Repair Rules
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

| File | Change |
|---|---|
| `cmd/git-workspace/refresh.go` | `runRefresh` delegates to `ReconcileWorkspace`; deleted `validateExistingRepos`, `scanWorkspaceForNewRepos`, and the local `discoverWorktrees`; added repair-only symlink rule, `rescanRemovedParents`, `writeRefreshSummary` |
| `cmd/git-workspace/refresh_test.go` | Updated 7 existing tests for the new signature; added 10 new tests |
| `cmd/git-workspace/deprecated.go` | Updated the deprecated `--refresh` call site |

Net effect: `refresh.go` lost ~150 lines of duplicated discovery logic to the shared engine and gained the responsibilities that are genuinely its own.

## CLI Output — Real Binary

Starting from the task 4.0 demo workspace, a repository was deleted, and a new one added three container directories deep:

```text
$ gws refresh
Refreshing workspace at: /tmp/…/demo/ws
Detecting git user configuration...
Cleared git status cache

Refresh complete!
Total repositories: 3
Removed 1 repository (path no longer valid)
Found 1 new repository
Repositories with user configuration: 3
Repositories with worktrees: 1
Worktrees: 2 (1 aligned, 1 unaligned)
```

The deeply nested `late-arrival` repository was found, which the previous one-level-deep scan could not have done.

## The Symlink Accumulation Fix — Real Binary

A curated workspace with one symlink at the root and one nested inside a container directory. Under the old behavior, the nested one has an external real path and so looked link-worthy, producing a duplicate at the root on every refresh.

```text
=== workspace BEFORE ===
apps -> /tmp/…/demo2/ext/apps
grouping/kubernetes -> /tmp/…/demo2/ext/kubernetes

$ gws init /tmp/…/demo2/ws
$ gws refresh    (×3)

=== workspace AFTER init + 3 refreshes ===
apps -> /tmp/…/demo2/ext/apps
grouping/kubernetes -> /tmp/…/demo2/ext/kubernetes
```

Byte-identical. No `kubernetes` link appeared at the workspace root.

## Test Results

```text
$ go test ./cmd/git-workspace/ -run 'Refresh|InitAndRefresh' -v
--- PASS: TestRunRefresh_ValidPathRetained (0.00s)
--- PASS: TestRunRefresh_MissingPathRemoved (0.00s)
--- PASS: TestRunRefresh_NewSymlinkDiscovered (0.00s)
--- PASS: TestRunRefresh_RealDirDiscovered (0.00s)
--- PASS: TestRunRefresh_BrokenSymlinkSkipped (0.00s)
--- PASS: TestRunRefresh_DiscoverWorktrees (0.01s)
--- PASS: TestRunRefresh_CorrectSymlinkNoDuplicate (0.00s)
--- PASS: TestRunRefresh_DiscoversDeeplyNestedRepository (0.10s)
--- PASS: TestRunRefresh_RemovesAndCountsDeletedRepository (0.07s)
--- PASS: TestRunRefresh_NewWorktreeVisibleImmediately (0.08s)
--- PASS: TestRunRefresh_IsIdempotent (0.05s)
--- PASS: TestRunRefresh_NoDuplicateLinkForNestedSymlink (0.02s)
--- PASS: TestRunRefresh_RepairsMissingLinkForTrackedExternalRepo (0.01s)
--- PASS: TestRunRefresh_ClearsStatusCache (0.07s)
--- PASS: TestRunRefresh_ParentRescanSafetyNet (0.02s)
--- PASS: TestInitAndRefresh_ReportIdenticalTotals (0.07s)
--- PASS: TestInitAndRefresh_IdenticalScanErrorOutput (0.01s)
ok  	github.com/daileyo/gws/cmd/git-workspace	0.549s
```

The first seven are the pre-existing refresh tests, passing **unmodified** apart from the `io.Writer` signature change — the engine swap preserved their behavior.

### Full suite, and under the race detector

```text
$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace	2.258s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	(cached)
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	(cached)
ok  	github.com/daileyo/gws/internal/reconcile	(cached)
ok  	github.com/daileyo/gws/internal/user	(cached)

$ go test ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	1.147s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	0.138s
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	0.436s
ok  	github.com/daileyo/gws/internal/reconcile	0.644s
ok  	github.com/daileyo/gws/internal/user	0.063s
```

No regressions in `list`, `navigate`, `tag`, `user`, or the worktree subcommands.

## Quality Gates

```text
$ make lint
Running linter...
(exit 0, no findings)
```

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| Discovers repositories added after init, including several containers deep | `TestRunRefresh_DiscoversDeeplyNestedRepository` (depth 4); CLI output above | ✅ PASS |
| Removes repositories whose stored paths are gone, and counts them | `TestRunRefresh_RemovesAndCountsDeletedRepository` | ✅ PASS |
| Worktree created after the last refresh appears after one refresh | `TestRunRefresh_NewWorktreeVisibleImmediately` | ✅ PASS |
| Two consecutive refreshes leave the workspace unchanged, no new symlinks | `TestRunRefresh_IsIdempotent` (listing + `config.json` byte-compared); CLI demo above | ✅ PASS |
| Missing link for a tracked external repo is repaired; nested-symlink repo gets no new link | `TestRunRefresh_RepairsMissingLinkForTrackedExternalRepo`, `TestRunRefresh_NoDuplicateLinkForNestedSymlink` (3 refreshes) | ✅ PASS |
| Status cache still cleared; parent-rescan safety net still recovers a replaced repository | `TestRunRefresh_ClearsStatusCache` (seeded then asserted empty), `TestRunRefresh_ParentRescanSafetyNet` | ✅ PASS |
| `init` and `refresh` report identical repository and worktree totals | `TestInitAndRefresh_ReportIdenticalTotals` | ✅ PASS |
| `init` and `refresh` emit byte-identical stderr for the same scan errors; stdout clean | `TestInitAndRefresh_IdenticalScanErrorOutput` (7 errors, exercises truncation) | ✅ PASS |
| `make lint` and `make test-race` pass | Quality Gates section above | ✅ PASS |

## Interpretation Note on the Repair-Only Rule (task 5.5)

Task 5.5 suggested determining "previously tracked with a symlink" by checking the pre-reconcile config **together with the presence of a matching symlink entry**. Implemented literally, that rule can never fire: the only case worth repairing is the one where the link is *missing*, so requiring the link to be present makes repair a no-op.

The implemented rule keeps the intent of round-2 Q4-C — repair only, never discovery-driven creation — using the signal that is actually available:

```text
create a link only when ALL of:
  1. the repository is not reachable from the workspace root  (not in ReachablePaths)
  2. the repository was already tracked before this reconcile  (in the pre-reconcile config)
  3. the repository lives outside the workspace root
```

Condition 2 is what makes it repair rather than creation: a tracked external repository can only have entered the library through `gws add` or a previous refresh, both of which create the workspace symlink. Condition 1 is what stops accumulation: anything the scan could reach — including a repository behind a nested symlink — already has a route from the workspace and is left alone.

`workspaceLinkOwnedRepos` and `repairWorkspaceSymlinks` in `refresh.go` carry this reasoning in comments. **This is a small interpretation call on the task wording and is flagged for review.**

## Design Notes

- **Retention vs. reachability.** The engine retains a repository that still exists on disk even when the scan cannot reach it (task 3.0). That is what makes repair possible: "exists but unreachable" is exactly the repair trigger.
- **`removeWorkspaceSymlink` is unchanged** per task 5.6 — it still removes only a symlink whose target matches the removed repository's path, leaving non-symlink and differently-targeted entries alone.
- **Parent rescan kept** as a safety net (round-2 Q11-C), now deduplicated by parent directory so several removals under one parent trigger a single scan.
- **Summary counts come from `Result`**, with recovered repositories from the safety net folded into the "Found N new" line so the report matches what was actually saved.
- **Two behavior changes in scan-error reporting**, as specified: the preview limit rose from 3 to 5, and warnings moved from stdout to stderr.
