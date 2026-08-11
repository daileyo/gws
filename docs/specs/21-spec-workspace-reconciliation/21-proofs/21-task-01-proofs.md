# Spec 21 — Task 1.0 Proofs: Boundary-Aware Recursive Scanner

**Task:** 1.0 Boundary-Aware Recursive Scanner
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

| File | Change |
|---|---|
| `internal/config/config.go` | Added `DefaultScanMaxDepth = 6`, `Preferences.ScanMaxDepth`, and `(*Config).EffectiveScanMaxDepth()` |
| `internal/config/config_test.go` | Added `TestEffectiveScanMaxDepth` and `TestScanMaxDepthBackwardCompatibility` |
| `internal/discovery/scanner.go` | Replaced `filepath.Walk` with an explicit recursive walker; added `Options{MaxDepth}`; exported `BuildRepository`; hidden-directory skipping |
| `internal/discovery/scanner_test.go` | Added 5 new tests; updated 2 tests that encoded deliberately changed behavior |
| `cmd/git-workspace/add.go` | Deleted duplicate `buildRepository`; now uses `discovery.BuildRepository`; updated `Scan` call site |
| `cmd/git-workspace/init.go` | Updated `Scan` call site for the new signature |
| `cmd/git-workspace/refresh.go` | Updated `Scan` and `BuildRepository` call sites |

## Configuration

New optional preference, additive and backward compatible:

```go
// internal/config/config.go
const DefaultScanMaxDepth = 6

type Preferences struct {
	StatusWorkers int `json:"status_workers,omitempty"`
	ScanMaxDepth  int `json:"scan_max_depth,omitempty"`
}

func (c *Config) EffectiveScanMaxDepth() int {
	if c != nil && c.Preferences != nil && c.Preferences.ScanMaxDepth > 0 {
		return c.Preferences.ScanMaxDepth
	}
	return DefaultScanMaxDepth
}
```

Example `~/.gws/config.json` fragment (paths are illustrative):

```json
{
  "version": "1.1.0",
  "workspace": "/home/user/gws",
  "preferences": {
    "status_workers": 8,
    "scan_max_depth": 6
  },
  "repositories": []
}
```

## Test Results

### Preference default resolution and backward compatibility

```text
### go test ./internal/config/ -run 'ScanMaxDepth' -v
=== RUN   TestEffectiveScanMaxDepth
=== RUN   TestEffectiveScanMaxDepth/nil_config_falls_back_to_default
=== RUN   TestEffectiveScanMaxDepth/nil_preferences_falls_back_to_default
=== RUN   TestEffectiveScanMaxDepth/preferences_present_but_field_unset_falls_back_to_default
=== RUN   TestEffectiveScanMaxDepth/explicit_value_is_used
=== RUN   TestEffectiveScanMaxDepth/zero_falls_back_to_default
=== RUN   TestEffectiveScanMaxDepth/negative_falls_back_to_default
--- PASS: TestEffectiveScanMaxDepth (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/nil_config_falls_back_to_default (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/nil_preferences_falls_back_to_default (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/preferences_present_but_field_unset_falls_back_to_default (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/explicit_value_is_used (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/zero_falls_back_to_default (0.00s)
    --- PASS: TestEffectiveScanMaxDepth/negative_falls_back_to_default (0.00s)
=== RUN   TestScanMaxDepthBackwardCompatibility
=== RUN   TestScanMaxDepthBackwardCompatibility/config_without_preferences_key_defaults_to_6
=== RUN   TestScanMaxDepthBackwardCompatibility/config_with_existing_preferences_but_no_scan_max_depth_defaults_to_6
=== RUN   TestScanMaxDepthBackwardCompatibility/scan_max_depth_round-trips_and_is_omitted_when_zero
--- PASS: TestScanMaxDepthBackwardCompatibility (0.00s)
    --- PASS: TestScanMaxDepthBackwardCompatibility/config_without_preferences_key_defaults_to_6 (0.00s)
    --- PASS: TestScanMaxDepthBackwardCompatibility/config_with_existing_preferences_but_no_scan_max_depth_defaults_to_6 (0.00s)
    --- PASS: TestScanMaxDepthBackwardCompatibility/scan_max_depth_round-trips_and_is_omitted_when_zero (0.00s)
PASS
ok  	github.com/daileyo/gws/internal/config	0.003s
```

### Scanner: pruning, filtering, and depth

```text
### go test ./internal/discovery/ -v
=== RUN   TestScan_SingleRepository
--- PASS: TestScan_SingleRepository (0.00s)
=== RUN   TestScan_MultipleRepositories
--- PASS: TestScan_MultipleRepositories (0.00s)
=== RUN   TestScan_NestedRepositories
--- PASS: TestScan_NestedRepositories (0.00s)
=== RUN   TestScan_NonExistentDirectory
--- PASS: TestScan_NonExistentDirectory (0.00s)
=== RUN   TestScan_EmptyDirectory
--- PASS: TestScan_EmptyDirectory (0.00s)
=== RUN   TestScan_SkipsCommonDirectories
--- PASS: TestScan_SkipsCommonDirectories (0.01s)
=== RUN   TestScan_RepositoryWithoutRemote
--- PASS: TestScan_RepositoryWithoutRemote (0.00s)
=== RUN   TestScan_BoundaryPruningThroughContainers
--- PASS: TestScan_BoundaryPruningThroughContainers (0.01s)
=== RUN   TestScan_SiblingsScannedAfterPruning
--- PASS: TestScan_SiblingsScannedAfterPruning (0.01s)
=== RUN   TestScan_SkipsFilteredDirectoriesAtEveryLevel
--- PASS: TestScan_SkipsFilteredDirectoriesAtEveryLevel (0.01s)
=== RUN   TestScan_MaxDepth
=== RUN   TestScan_MaxDepth/default_depth_registers_through_depth_6
=== RUN   TestScan_MaxDepth/explicit_depth_3_registers_through_depth_3
=== RUN   TestScan_MaxDepth/explicit_depth_1_registers_only_the_top_level
=== RUN   TestScan_MaxDepth/zero_falls_back_to_the_default
=== RUN   TestScan_MaxDepth/negative_falls_back_to_the_default
--- PASS: TestScan_MaxDepth (0.01s)
    --- PASS: TestScan_MaxDepth/default_depth_registers_through_depth_6 (0.00s)
    --- PASS: TestScan_MaxDepth/explicit_depth_3_registers_through_depth_3 (0.00s)
    --- PASS: TestScan_MaxDepth/explicit_depth_1_registers_only_the_top_level (0.00s)
    --- PASS: TestScan_MaxDepth/zero_falls_back_to_the_default (0.00s)
    --- PASS: TestScan_MaxDepth/negative_falls_back_to_the_default (0.00s)
=== RUN   TestScan_RepositoryAtScanRoot
--- PASS: TestScan_RepositoryAtScanRoot (0.00s)
=== RUN   TestShouldSkipDir
--- PASS: TestShouldSkipDir (0.00s)
PASS
ok  	github.com/daileyo/gws/internal/discovery	0.055s
```

### Full suite under the race detector

```text
$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace
ok  	github.com/daileyo/gws/internal/classifier
ok  	github.com/daileyo/gws/internal/config
ok  	github.com/daileyo/gws/internal/discovery
ok  	github.com/daileyo/gws/internal/filter
ok  	github.com/daileyo/gws/internal/git
ok  	github.com/daileyo/gws/internal/user
Tests complete
```

## Quality Gates

```text
$ golangci-lint --version
golangci-lint has version v1.64.8 built with go1.26.5

$ make lint
Running linter...
(exit 0, no findings)
```

No `gosec` G304 findings were produced: task 1.0 reads directory entries and `os.Lstat` metadata only. File reads of `.git` files arrive in task 2.0 and will be re-checked there.

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| Repositories nested under container directories all discovered; nothing below a repository root registered | `TestScan_BoundaryPruningThroughContainers`, `TestScan_SiblingsScannedAfterPruning` | ✅ PASS |
| Repository containing a nested clone registers only the outer repository | `TestScan_NestedRepositories`, `TestScan_BoundaryPruningThroughContainers` (asserts `embedded` absent) | ✅ PASS |
| Skip-list and hidden directories skipped at every level | `TestScan_SkipsFilteredDirectoriesAtEveryLevel`, `TestScan_SkipsCommonDirectories`, `TestShouldSkipDir` | ✅ PASS |
| Repositories at exactly `scan_max_depth` registered, deeper ones not; non-default value moves the boundary | `TestScan_MaxDepth` (5 sub-cases across depths 1–8) | ✅ PASS |
| Pre-existing `config.json` loads with `scan_max_depth` defaulting to 6 | `TestScanMaxDepthBackwardCompatibility` | ✅ PASS |
| `make lint` and `make test-race` pass | Quality Gates section above | ✅ PASS |

## Behavior Changes Recorded

Two pre-existing tests asserted behavior this spec deliberately changes. Both were updated rather than worked around:

1. **`TestScan_NestedRepositories`** previously expected 2 repositories when one clone was nested inside another. Strict boundary pruning (spec Unit 1, round-1 Q1-A) means only the outer repository is registered. The test now asserts 1 and names it.
2. **`TestShouldSkipDir`** previously expected `shouldSkipDir(".git") == false` with the comment "`.git` is handled separately". Hidden-directory skipping (round-2 Q7-A) makes every dot-prefixed directory skippable, so the expectation is now `true`. Cases for `.dotfiles`, `.`, and `not.hidden` were added to pin the boundary of the rule.

### Inherited change to `gws add --recursive`

`discovery.Scan` is shared with `gws add --recursive` (`cmd/git-workspace/add.go`). That command therefore inherits boundary pruning, hidden-directory skipping, and the depth cap. The spec lists `gws add` as a non-goal, and no `add` behavior was modified directly — all 12 existing `add` tests pass unchanged:

```text
--- PASS: TestRunAdd_CurrentDirectory (0.01s)
--- PASS: TestRunAdd_ExplicitExternalPath (0.01s)
--- PASS: TestRunAdd_NonGitDirectory (0.00s)
--- PASS: TestRunAdd_AlreadyTracked (0.00s)
--- PASS: TestRunAdd_SymlinkPathConflict (0.00s)
--- PASS: TestRunAdd_NoWorkspace (0.00s)
--- PASS: TestCreateSymlinkIfExternal_InsideWorkspace (0.00s)
--- PASS: TestRunAdd_ConfigMetadata (0.01s)
--- PASS: TestRunAddRecursive_MultipleRepos (0.02s)
--- PASS: TestRunAddRecursive_AllAlreadyTracked (0.01s)
--- PASS: TestRunAddRecursive_NoRepos (0.00s)
--- PASS: TestRunAddRecursive_RequiresAdd (0.00s)
```

`add --recursive` now also reads `EffectiveScanMaxDepth()` from the loaded config rather than scanning without a cap.

## Known Interim State

`classifyDir` currently returns `dirGitLinked` for any directory whose `.git` is a **file**, and such directories are neither registered nor traversed. This is correct for linked worktrees but temporarily also excludes clones created with `git clone --separate-git-dir`. Task 2.1–2.4 adds `IsLinkedWorktree` to distinguish the two, after which only genuine worktrees are excluded. No test in this task asserts the interim behavior.
