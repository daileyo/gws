# Spec 21 — Task 2.0 Proofs: Symlink Resolution and Structural Worktree Detection

**Task:** 2.0 Symlink Resolution and Structural Worktree Detection
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

| File | Change |
|---|---|
| `internal/git/worktree.go` | Added `IsLinkedWorktree`, `parseGitdirPointer`, `mainRepoFromWorktreeGitDir`, `linkedWorktreeViaGit` |
| `internal/git/worktree_test.go` | Added 9 tests covering detection across all four repository kinds plus pointer parsing |
| `internal/discovery/scanner.go` | Symlink resolution with loop guards; visited real-path set; `ReachablePaths` and `WorktreeMainRepos` on `ScanResult` |
| `internal/discovery/scanner_test.go` | Added 11 tests for symlinks, dedupe, loop safety, worktree exclusion, and error resilience |

## Structural Worktree Detection — Empirical Basis

Detection uses no directory naming whatsoever. The rule was verified against live git repositories of every relevant kind before implementation:

```text
=== normal clone ===
.git is a:      directory
git-dir:        [.git]
git-common-dir: [.git]                                   → equal    → NOT a worktree

=== linked worktree ===
.git is a:      regular file
contents:       gitdir: /…/main-repo/.git/worktrees/feature
git-dir:        [/…/main-repo/.git/worktrees/feature]
git-common-dir: [/…/main-repo/.git]                      → differ   → IS a worktree

=== separate-git-dir clone ===
.git is a:      regular file
contents:       gitdir: /…/sep.git
git-dir:        [/…/sep.git]
git-common-dir: [/…/sep.git]                             → equal    → NOT a worktree
```

The fast path reads the `.git` file and accepts it as a worktree only when the `gitdir:` target has the shape `<main>/.git/worktrees/<name>`. Anything else — a submodule pointer at `.git/modules/…`, a `--separate-git-dir` target, malformed content, or an unreadable file — falls through to the `rev-parse` comparison above.

## Loop Safety Design

Three independent mechanisms, each covered by a test:

1. **Visited real-path set** — every directory is resolved with `filepath.EvalSymlinks` before processing, and a directory already visited is skipped. This is also what makes deduplication order-independent.
2. **Workspace-root guard** — a symlink whose target is the scan root or an ancestor of it is not followed.
3. **Current-branch guard** — a symlink whose target is an ancestor of the current traversal path is not followed.

Loop-safety tests run inside a 30-second deadline via `scanWithDeadline`, so a regression fails the test rather than hanging CI.

## Test Results

### `internal/git` — worktree detection

```text
--- PASS: TestIsLinkedWorktree_NormalClone (0.03s)
--- PASS: TestIsLinkedWorktree_LinkedWorktree (0.02s)
--- PASS: TestIsLinkedWorktree_CloneInDotWtDirectory (0.01s)
--- PASS: TestIsLinkedWorktree_SeparateGitDir (0.01s)
--- PASS: TestIsLinkedWorktree_MalformedGitFile (0.00s)
    --- PASS: TestIsLinkedWorktree_MalformedGitFile/garbage_content (0.00s)
    --- PASS: TestIsLinkedWorktree_MalformedGitFile/empty_file (0.00s)
    --- PASS: TestIsLinkedWorktree_MalformedGitFile/gitdir_with_no_worktrees_segment (0.00s)
--- PASS: TestIsLinkedWorktree_MalformedGitFileInRealWorktree (0.02s)
--- PASS: TestIsLinkedWorktree_NoGitEntry (0.00s)
--- PASS: TestParseGitdirPointer (0.00s)
    --- PASS: TestParseGitdirPointer/standard_pointer (0.00s)
    --- PASS: TestParseGitdirPointer/no_space_after_colon (0.00s)
    --- PASS: TestParseGitdirPointer/extra_whitespace (0.00s)
    --- PASS: TestParseGitdirPointer/relative_pointer (0.00s)
    --- PASS: TestParseGitdirPointer/not_a_pointer (0.00s)
    --- PASS: TestParseGitdirPointer/empty (0.00s)
--- PASS: TestMainRepoFromWorktreeGitDir (0.00s)
    --- PASS: TestMainRepoFromWorktreeGitDir/standard_worktree_git_dir (0.00s)
    --- PASS: TestMainRepoFromWorktreeGitDir/trailing_separator (0.00s)
    --- PASS: TestMainRepoFromWorktreeGitDir/submodule_git_dir (0.00s)
    --- PASS: TestMainRepoFromWorktreeGitDir/separate_git_dir (0.00s)
    --- PASS: TestMainRepoFromWorktreeGitDir/plain_git_dir (0.00s)
PASS
ok  	github.com/daileyo/gws/internal/git
```

### `internal/discovery` — full package

```text
--- PASS: TestScan_SingleRepository (0.00s)
--- PASS: TestScan_MultipleRepositories (0.00s)
--- PASS: TestScan_NestedRepositories (0.00s)
--- PASS: TestScan_NonExistentDirectory (0.00s)
--- PASS: TestScan_EmptyDirectory (0.00s)
--- PASS: TestScan_SkipsCommonDirectories (0.01s)
--- PASS: TestScan_RepositoryWithoutRemote (0.00s)
--- PASS: TestScan_BoundaryPruningThroughContainers (0.01s)
--- PASS: TestScan_SiblingsScannedAfterPruning (0.00s)
--- PASS: TestScan_SkipsFilteredDirectoriesAtEveryLevel (0.01s)
--- PASS: TestScan_MaxDepth (0.01s)
    --- PASS: TestScan_MaxDepth/default_depth_registers_through_depth_6 (0.00s)
    --- PASS: TestScan_MaxDepth/explicit_depth_3_registers_through_depth_3 (0.00s)
    --- PASS: TestScan_MaxDepth/explicit_depth_1_registers_only_the_top_level (0.00s)
    --- PASS: TestScan_MaxDepth/zero_falls_back_to_the_default (0.00s)
    --- PASS: TestScan_MaxDepth/negative_falls_back_to_the_default (0.00s)
--- PASS: TestScan_RepositoryAtScanRoot (0.00s)
--- PASS: TestScan_CuratedSymlinkWorkspace (0.00s)
--- PASS: TestScan_SymlinkedRepoInsideContainer (0.00s)
--- PASS: TestScan_DeduplicatesByRealPath (0.00s)
    --- PASS: TestScan_DeduplicatesByRealPath/physical_visited_first (0.00s)
    --- PASS: TestScan_DeduplicatesByRealPath/symlink_visited_first (0.00s)
--- PASS: TestScan_DeduplicatesExternalRepoReachedTwice (0.00s)
--- PASS: TestScan_SymlinkLoopSafety (0.01s)
    --- PASS: TestScan_SymlinkLoopSafety/symlink_pointing_at_the_workspace_root (0.00s)
    --- PASS: TestScan_SymlinkLoopSafety/symlink_pointing_at_its_own_parent (0.00s)
    --- PASS: TestScan_SymlinkLoopSafety/two_directories_symlinked_to_each_other (0.00s)
    --- PASS: TestScan_SymlinkLoopSafety/symlink_to_an_ancestor_above_the_workspace_root (0.00s)
--- PASS: TestScan_WorktreesAreNotRegisteredAsRepositories (0.02s)
--- PASS: TestScan_CloneInDotWtDirectoryIsRegistered (0.00s)
--- PASS: TestScan_BrokenSymlinkRecordedAsError (0.00s)
--- PASS: TestScan_UnreadableDirectoryRecordedAsError (0.00s)
--- PASS: TestIsAncestorOrEqual (0.00s)
--- PASS: TestShouldSkipDir (0.00s)
PASS
ok  	github.com/daileyo/gws/internal/discovery
```

### Full suite under the race detector

```text
$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace	1.473s
ok  	github.com/daileyo/gws/internal/classifier	(cached)
ok  	github.com/daileyo/gws/internal/config	(cached)
ok  	github.com/daileyo/gws/internal/discovery	1.232s
ok  	github.com/daileyo/gws/internal/filter	(cached)
ok  	github.com/daileyo/gws/internal/git	1.493s
ok  	github.com/daileyo/gws/internal/user	(cached)
```

## Quality Gates

```text
$ make lint
Running linter...
(exit 0, no findings)
```

The one file read introduced in this task (`os.ReadFile` of a `.git` file) carries a narrowly scoped `#nosec G304` with a justification: the path is cleaned and only the fixed `.git` name is appended to a directory the scanner is already traversing.

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| Curated symlink-only workspace discovers all repositories | `TestScan_CuratedSymlinkWorkspace` (3 external repos via symlinks), `TestScan_SymlinkedRepoInsideContainer` | ✅ PASS |
| Repository reachable directly and via symlink yields one entry, in both traversal orders | `TestScan_DeduplicatesByRealPath` (2 sub-cases with names forcing each order), `TestScan_DeduplicatesExternalRepoReachedTwice` | ✅ PASS |
| Symlink to workspace root and to an ancestor complete without infinite recursion | `TestScan_SymlinkLoopSafety` (4 sub-cases, 30s deadline) | ✅ PASS |
| Repository with a linked worktree registers the repository, not the worktree | `TestScan_WorktreesAreNotRegisteredAsRepositories` (aligned + unaligned worktrees) | ✅ PASS |
| Clone in a `<name>.wt` directory is registered as a repository | `TestScan_CloneInDotWtDirectoryIsRegistered`, `TestIsLinkedWorktree_CloneInDotWtDirectory` | ✅ PASS |
| Malformed `.git` file triggers the `rev-parse` fallback | `TestIsLinkedWorktree_MalformedGitFile` (3 sub-cases), `TestIsLinkedWorktree_MalformedGitFileInRealWorktree` | ✅ PASS |
| Broken and inaccessible symlinks recorded as errors without aborting | `TestScan_BrokenSymlinkRecordedAsError`, `TestScan_UnreadableDirectoryRecordedAsError` | ✅ PASS |
| `make lint` and `make test-race` pass | Quality Gates section above | ✅ PASS |

## Notes for Downstream Tasks

- **`ScanResult.ReachablePaths`** is populated with the resolved path of every registered repository. Task 5.5 uses it for the repair-only symlink rule.
- **`ScanResult.WorktreeMainRepos`** collects the main repository path inferred from every linked worktree encountered, including main repositories outside the scan root. Task 3.7 uses it for orphan main-repository registration, avoiding a second filesystem pass.
- The interim state noted in task 1.0 is resolved: directories whose `.git` is a file are now classified, and only genuine linked worktrees are excluded. `--separate-git-dir` clones register normally, covered by `TestIsLinkedWorktree_SeparateGitDir`.
- Hidden directories are skipped at every level per round-2 Q7-A, which means a **hidden symlink** at the workspace root (for example `~/gws/.config-nvim`) is also skipped. This follows directly from the chosen rule; making it configurable was recorded as future work.
