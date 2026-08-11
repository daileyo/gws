# Spec 21 — Task 3.0 Proofs: Shared Reconciliation Engine

**Task:** 3.0 Shared Reconciliation Engine
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

| File | Change |
|---|---|
| `internal/reconcile/result.go` | New — package doc and the `Result` contract |
| `internal/reconcile/reconcile.go` | New — `ReconcileWorkspace`, `Options`, `ProgressReporter` |
| `internal/reconcile/merge.go` | New — merge rules, `repositoryExists`, `refreshMetadata` |
| `internal/reconcile/worktree.go` | New — worktree phase and orphan main-repository registration |
| `internal/reconcile/errors.go` | New — shared `ReportScanErrors` |
| `internal/reconcile/*_test.go` | New — 32 tests across 5 files |
| `internal/git/progress.go` | Added `NewProgressWithLabel` and `DefaultProgressLabel`; `NewProgress` unchanged for its existing caller |
| `internal/git/progress_test.go` | Added label tests |

## The Engine Contract

```go
func ReconcileWorkspace(workspaceRoot string, existing *config.Config, opts Options) (*Result, error)
```

`existing == nil` is the `init` case; a populated config is the `refresh` case. Three phases run in order: **discover → merge → worktrees**.

The engine is deliberately pure. It does not clear the status cache and does not create, repair, or remove symlinks — those stay with the calling command, per round-1 Q7-D and Q8-B.

### Retention rule

Retention is decided by **whether a tracked path still exists on disk**, not by whether the scan reached it. A repository can remain valid while becoming unreachable — for example when a curated symlink is deleted — and dropping it would silently discard the user's tags. This is also what makes the repair-only symlink rule in task 5.0 possible: `ReachablePaths` distinguishes "exists but unreachable" from "exists and reachable".

## Test Results

### `internal/reconcile` — full package (32 tests)

```text
--- PASS: TestReportScanErrors (0.00s)
    --- PASS: TestReportScanErrors/no_errors_writes_nothing (0.00s)
    --- PASS: TestReportScanErrors/three_errors_are_all_listed (0.00s)
    --- PASS: TestReportScanErrors/exactly_five_errors_are_all_listed (0.00s)
    --- PASS: TestReportScanErrors/six_errors_list_five_and_note_one_more (0.00s)
    --- PASS: TestReportScanErrors/ten_errors_list_five_and_note_five_more (0.00s)
--- PASS: TestReportScanErrors_NilSlice (0.00s)
--- PASS: TestReportScanErrors_ExactFormat (0.00s)
--- PASS: TestMergeRepositories (0.08s)
    --- PASS: TestMergeRepositories/empty_existing_adds_everything_discovered (0.02s)
    --- PASS: TestMergeRepositories/already_tracked_repository_is_retained,_not_re-added (0.02s)
    --- PASS: TestMergeRepositories/newly_discovered_repository_is_added_alongside_tracked_ones (0.02s)
    --- PASS: TestMergeRepositories/duplicate_entries_in_stored_config_collapse_to_one (0.02s)
--- PASS: TestMergeRepositories_RemovesVanishedPaths (0.02s)
--- PASS: TestMergeRepositories_PreservesTags (0.01s)
--- PASS: TestMergeRepositories_CountsMetadataUpdates (0.06s)
    --- PASS: TestMergeRepositories_CountsMetadataUpdates/unchanged_metadata_is_not_counted (0.01s)
    --- PASS: TestMergeRepositories_CountsMetadataUpdates/changed_remote_URL_is_counted (0.01s)
    --- PASS: TestMergeRepositories_CountsMetadataUpdates/changed_type_is_counted (0.01s)
    --- PASS: TestMergeRepositories_CountsMetadataUpdates/changed_visibility_is_counted (0.01s)
    --- PASS: TestMergeRepositories_CountsMetadataUpdates/changed_name_is_counted (0.01s)
--- PASS: TestMergeRepositories_RetainsUnreachableButExistingRepository (0.01s)
--- PASS: TestRepositoryExists (0.01s)
--- PASS: TestReconcileWorkspace_InitPath (0.04s)
--- PASS: TestReconcileWorkspace_RequiresWorkspaceRoot (0.00s)
--- PASS: TestReconcileWorkspace_ResultsAreSortedByPath (0.04s)
--- PASS: TestReconcileWorkspace_DeterministicAcrossRuns (0.05s)
--- PASS: TestReconcileWorkspace_NoSideEffects (0.03s)
--- PASS: TestReconcileWorkspace_ResultContract (0.06s)
--- PASS: TestReconcileWorkspace_ProgressLifecycle (0.03s)
--- PASS: TestReconcileWorkspace_ScanErrorsAreReturnedNotPrinted (0.01s)
--- PASS: TestDiscoverWorktrees_AlignmentClassification (0.03s)
--- PASS: TestDiscoverWorktrees_NoWorktrees (0.01s)
--- PASS: TestDiscoverWorktrees_ClearsStaleEntries (0.02s)
--- PASS: TestDiscoverWorktrees_IncrementsProgressPerRepository (0.04s)
--- PASS: TestRegisterOrphanMainRepos_InsideWorkspace (0.01s)
--- PASS: TestRegisterOrphanMainRepos_OutsideWorkspace (0.01s)
--- PASS: TestRegisterOrphanMainRepos_SkipsAlreadyTracked (0.01s)
--- PASS: TestRegisterOrphanMainRepos_SkipsMissingPaths (0.00s)
--- PASS: TestRegisterOrphanMainRepos_EmptyInput (0.00s)
--- PASS: TestReconcileWorkspace_RegistersOrphanMainRepoEndToEnd (0.02s)
PASS
ok  	github.com/daileyo/gws/internal/reconcile	0.628s
```

### `internal/git` — progress label

```text
--- PASS: TestProgress_Increment (0.00s)
    --- PASS: TestProgress_Increment/concurrent_increments_are_safe (0.00s)
--- PASS: TestNewProgress_UsesDefaultLabel (0.00s)
--- PASS: TestNewProgressWithLabel (0.00s)
    --- PASS: TestNewProgressWithLabel/custom_label_is_used (0.00s)
    --- PASS: TestNewProgressWithLabel/empty_label_falls_back_to_the_default (0.00s)
--- PASS: TestProgress_StopWithoutStart (0.00s)
--- PASS: TestProgress_NonTTY (0.00s)
--- PASS: TestProgress_StartStop (0.00s)
ok  	github.com/daileyo/gws/internal/git
```

### Full suite under the race detector

```text
$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace	1.536s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	1.296s
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	1.525s
ok  	github.com/daileyo/gws/internal/reconcile	1.715s
ok  	github.com/daileyo/gws/internal/user	1.080s
```

## Quality Gates

```text
$ make lint
Running linter...
(exit 0, no findings)
```

## Shared Scan-Error Reporter — Exact Output

`TestReportScanErrors_ExactFormat` pins the byte-level shape both commands will share, so one cannot silently drift from the other:

```text
Warning: 7 errors occurred during scanning:
  - scan problem 1
  - scan problem 2
  - scan problem 3
  - scan problem 4
  - scan problem 5
  ... and 2 more errors
```

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| Nil existing config returns the full discovered set with correct counts | `TestReconcileWorkspace_InitPath` | ✅ PASS |
| Tags survive on repositories whose paths remain valid | `TestMergeRepositories_PreservesTags` (tags, user, email, signing, source) | ✅ PASS |
| Add / remove / update rules and update counting | `TestMergeRepositories` (4 cases), `TestMergeRepositories_RemovesVanishedPaths`, `TestMergeRepositories_CountsMetadataUpdates` (5 cases) | ✅ PASS |
| Worktrees populated with correct aligned classification | `TestDiscoverWorktrees_AlignmentClassification` | ✅ PASS |
| Orphan worktree causes its main repository to be registered, including outside the workspace root | `TestRegisterOrphanMainRepos_OutsideWorkspace`, `TestReconcileWorkspace_RegistersOrphanMainRepoEndToEnd` | ✅ PASS |
| No status-cache or symlink side effects | `TestReconcileWorkspace_NoSideEffects` (recursive directory snapshot + cache sentinel) | ✅ PASS |
| Complete `Result` contract reported | `TestReconcileWorkspace_ResultContract` (all 8 counts in one workspace) | ✅ PASS |
| Scan-error reporter: stderr, cap 5, correct remainder | `TestReportScanErrors` (5 cases), `TestReportScanErrors_ExactFormat` | ✅ PASS |
| `make lint` and `make test-race` pass | Quality Gates section above | ✅ PASS |

## Design Notes

- **Determinism (task 3.4).** Merged repositories are sorted by path. `TestReconcileWorkspace_DeterministicAcrossRuns` reconciles twice and asserts identical ordering with `added=0, removed=0` on the second pass. Without this, `config.json` would churn and the task 5.0 idempotence test would be flaky for the wrong reason.
- **Progress ownership.** `ProgressReporter` is a three-method interface that `*git.Progress` already satisfies, so the engine never imports terminal concerns. `Stop` is deferred immediately after `Start`, so it runs on the error path too.
- **`NewProgress` is untouched** for its existing caller at `cmd/git-workspace/list.go:504`; the label arrives through the new `NewProgressWithLabel`.
- **Duplicate config entries** are collapsed during merge. Not a spec requirement, but the merge builds a path-keyed set anyway, and tolerating a hand-edited `config.json` costs nothing.

## Two Test-Authoring Corrections

Both were faults in the tests, not the implementation, and are recorded because they show what the assertions actually mean:

1. **Metadata-update baseline.** The first version of `TestMergeRepositories_CountsMetadataUpdates` built its "unchanged" baseline from a hand-made `config.Repository` that lacked the fields `BuildRepository` fills in, so every case read as changed. The baseline now comes from `discovery.BuildRepository` — the real metadata on disk.
2. **Visibility mutation.** The visibility case mutated the stored copy to `VisibilityUnknown`, which is exactly what an `https://` remote already classifies as (`internal/classifier/detector.go:54`), making it a no-op. It now mutates to `VisibilityPrivate`.
