# Spec 21 Validation Report — Unified Workspace Reconciliation

**Spec:** `docs/specs/21-spec-workspace-reconciliation/21-spec-workspace-reconciliation.md`
**Task List:** `docs/specs/21-spec-workspace-reconciliation/21-tasks-workspace-reconciliation.md`
**Branch:** `feat/workspace-reconciliation`
**Baseline:** `c4eccde` (chore(main): release 2.19.2)
**Head:** `f107744`

---

## 1) Executive Summary

**Overall: PASS** — all six gates satisfied after remediation.

> **Revalidation, 2026-08-10.** The first pass returned **FAIL**: GATE A on issue V-1 (HIGH) and GATE D on issue V-2 (MEDIUM). Both were fixed in commit `f107744` and re-verified. This report reflects the post-fix state; the original findings are retained in full in §3 with their resolutions, because the defect V-1 describes is a real behavior that shipped in no release but is worth the record.

**Implementation Ready: Yes** — all 59 functional requirements and all 9 original acceptance criteria verify against independently re-executed evidence, with the full suite, race detector, linter, and strict docs build clean.

### Key Metrics

| Metric | Result |
|---|---|
| Functional requirements verified | **59 / 59 (100%)** — 0 Failed, 0 Unknown |
| Proof artifacts working | **107 / 107 test claims (100%)** — all exist in source and pass on a cleared cache |
| Original acceptance criteria met | **9 / 9** |
| Files changed vs declared | 29 code/doc changed; **29 declared-and-changed; 0 undeclared** |
| Code delta | 25 Go files, +3,786 / −322 |
| Full suite | **388 tests pass**, `-race` clean, `make vet` and `make lint` exit 0 |

---

## 2) Coverage Matrix

### Functional Requirements

Requirements are grouped by spec section; line numbers refer to the spec file.

#### Unit 1 — Traversal and Boundaries (spec L36–42)

| Requirement | Status | Evidence |
|---|---|---|
| FR-1 Recursive traversal through container directories | Verified | `TestScan_BoundaryPruningThroughContainers`; `scanner.go:113` `walk` |
| FR-2 Register on `.git` **directory**; strict pruning below | Verified | `TestScan_NestedRepositories`, `TestScan_BoundaryPruningThroughContainers` (asserts nested clone absent); `scanner.go:120-123` |
| FR-3 Siblings and ancestors still scanned after pruning | Verified | `TestScan_SiblingsScannedAfterPruning` |
| FR-4 Skip list honored | Verified | `TestScan_SkipsCommonDirectories`, `TestShouldSkipDir`; `scanner.go` `shouldSkipDir` |
| FR-5 Hidden directories skipped at every level | Verified | `TestScan_SkipsFilteredDirectoriesAtEveryLevel` (depths 1 and 2), `TestShouldSkipDir` |
| FR-6 Stop past max depth; repo at exactly max depth registered | Verified | `TestScan_MaxDepth` (depths 1–8, 5 sub-cases) |
| FR-7 Depth from `scan_max_depth`, default 6 | Verified | `TestEffectiveScanMaxDepth` (6 cases), `TestScanMaxDepthBackwardCompatibility`; `config.go` `DefaultScanMaxDepth = 6` |

#### Unit 1 — Worktree Detection (spec L46–49)

| Requirement | Status | Evidence |
|---|---|---|
| FR-8 Worktrees never registered as repositories | Verified | `TestScan_WorktreesAreNotRegisteredAsRepositories` (aligned + unaligned) |
| FR-9 Fast path: read `.git` file, match `worktrees/` shape | Verified | `TestIsLinkedWorktree_LinkedWorktree`, `TestParseGitdirPointer`, `TestMainRepoFromWorktreeGitDir` |
| FR-10 Fallback to `rev-parse --git-dir` vs `--git-common-dir` | Verified | `TestIsLinkedWorktree_MalformedGitFile` (3 cases), `TestIsLinkedWorktree_MalformedGitFileInRealWorktree`, `TestIsLinkedWorktree_SeparateGitDir` |
| FR-11 Structural only; `.wt`-named clone still a repository | Verified | `TestScan_CloneInDotWtDirectoryIsRegistered`, `TestIsLinkedWorktree_CloneInDotWtDirectory`; no name-based branch exists in `scanner.go` |

#### Unit 1 — Symlinks and Identity (spec L53–59)

| Requirement | Status | Evidence |
|---|---|---|
| FR-12 Follow symlinked directories at any depth | Verified | `TestScan_CuratedSymlinkWorkspace`, `TestScan_SymlinkedRepoInsideContainer` |
| FR-13 `Repository.Path` is the resolved real path | Verified | `TestScan_CuratedSymlinkWorkspace` (asserts real-path prefix); `scanner.go` `register(dir)` |
| FR-14 Dedupe by real path, order-independent | Verified | `TestScan_DeduplicatesByRealPath` (both orders), `TestScan_DeduplicatesExternalRepoReachedTwice` |
| FR-15 Visited real-path set; no directory twice | Verified | `scanner.go:114-117`; `TestScan_SymlinkLoopSafety` |
| FR-16 Refuse symlink targeting workspace root or above | Verified | `TestScan_SymlinkLoopSafety` sub-cases 1 and 4; `TestIsAncestorOrEqual` |
| FR-17 Refuse symlink targeting an ancestor of current path | Verified | `TestScan_SymlinkLoopSafety` sub-cases 2 and 3 |
| FR-18 Broken/inaccessible symlinks recorded, scan continues | Verified | `TestScan_BrokenSymlinkRecordedAsError`, `TestScan_UnreadableDirectoryRecordedAsError` |

#### Unit 2 — Engine, Merge, Worktrees, Summary (spec L80–114)

| Requirement | Status | Evidence |
|---|---|---|
| FR-19 `internal/reconcile` with `ReconcileWorkspace(root, existing, …)` | Verified | `internal/reconcile/reconcile.go`; nil-existing path covered by `TestReconcileWorkspace_InitPath`. Signature adds an `Options` parameter, permitted by the spec's "conceptually equivalent to" |
| FR-20 `internal/discovery` reduced to scanning; merge lives in reconcile | Verified | Zero `merge`/`Reconcile` symbols in `scanner.go` (grep count 0) |
| FR-21 Three phases in order; returns set + summary | Verified | `reconcile.go:47-89`; `TestReconcileWorkspace_ResultContract` |
| FR-22 Engine touches neither cache nor symlinks | Verified | `TestReconcileWorkspace_NoSideEffects`; grep shows 0 occurrences of `os.Symlink`/`GetCachePath`/`NewCache` across all non-test reconcile files |
| FR-23 Retain metadata for existing paths, incl. tags and user config | Verified | `TestMergeRepositories_PreservesTags`, `TestMergeRepositories_RetainsUnreachableButExistingRepository` |
| FR-24 Preserve user-set tags on every retained repository | Verified | `TestMergeRepositories_PreservesTags`; `TestRunRefresh_ValidPathRetained`; live CLI: tag `work` survived a refresh |
| FR-25 Add discovered repositories not already tracked | Verified | `TestMergeRepositories` case 3, `TestRunRefresh_DiscoversDeeplyNestedRepository` |
| FR-26 Remove entries whose stored path no longer exists | Verified | `TestMergeRepositories_RemovesVanishedPaths`, `TestRunRefresh_RemovesAndCountsDeletedRepository` |
| FR-27 Re-read name/remote/type/visibility; count updates | Verified | `TestMergeRepositories_CountsMetadataUpdates` (5 cases incl. a no-change control) |
| FR-28 Worktrees via `git worktree list --porcelain` | Verified | `reconcile/worktree.go` → `git.ListWorktrees`; `TestDiscoverWorktrees_AlignmentClassification` |
| FR-29 Add new worktrees; remove vanished ones | Verified | `TestDiscoverWorktrees_ClearsStaleEntries`, `TestRunRefresh_NewWorktreeVisibleImmediately` |
| FR-30 Mark aligned when inside `<repo>.wt/` | Verified | `TestDiscoverWorktrees_AlignmentClassification` |
| FR-31 Register main repo of an untracked linked worktree | Verified | `TestRegisterOrphanMainRepos_OutsideWorkspace`, `TestReconcileWorkspace_RegistersOrphanMainRepoEndToEnd` |
| FR-32 Result returns all eight counts | Verified | `TestReconcileWorkspace_ResultContract` |
| FR-33 Scan errors returned, not printed | Verified | `TestReconcileWorkspace_ScanErrorsAreReturnedNotPrinted` |
| FR-34 Single shared scan-error reporter | Verified | `reconcile/errors.go`; `TestInitAndRefresh_IdenticalScanErrorOutput` compares both commands' stderr |
| FR-35 stderr, header, cap 5, remainder line | Verified | `TestReportScanErrors` (5 cases), `TestReportScanErrors_ExactFormat` (byte-exact) |
| FR-36 `init` reports errors exclusively via the reporter | Verified | `TestRunInit_ReportsScanErrorsOnStderr`; inline block removed from `init.go` |
| FR-37 Engine accepts progress using `internal/git/progress.go` | Verified | `ProgressReporter` in `reconcile.go`; `TestReconcileWorkspace_ProgressLifecycle` |
| FR-38 Both commands display progress unconditionally | Verified (note) | `init.go:74` and `refresh.go:61` both construct and pass a reporter with no count threshold. See Observation O-1 on TTY gating |

#### Unit 2 — `gws init` (spec L118–132)

| Requirement | Status | Evidence |
|---|---|---|
| FR-39 Create library, record root, reconcile with nil, persist, report | Verified | `TestRunInit_FullReconciliationSummary`; live CLI transcript |
| FR-40 Defaults to current directory | Verified | `TestRunInit_HappyPath` (calls `runInit("")`) |
| FR-41 `init` creates no workspace symlinks | Verified | `TestRunInit_CreatesNoSymlinks`; grep: 0 symlink calls in `init.go` |
| FR-42 Guard: stderr, exit 0, recommends `refresh` | Verified | `TestRunInit_GuardMessageRecommendsRefresh` (asserts ordering and byte-identical config); live CLI exit code 0 |
| FR-43 Summary of the specified shape | Verified | Live CLI compared line-for-line against the documented block |
| FR-44 Conditional lines omitted at zero | Verified | `TestRunInit_OmitsZeroCountLines` |

#### Unit 3 — `gws refresh` (spec L154–183)

| Requirement | Status | Evidence |
|---|---|---|
| FR-45 Load config, reconcile with it, persist, report | Verified | `TestRunRefresh_*` (16 tests); live CLI transcript |
| FR-46 Still clears status cache with its message | Verified | `TestRunRefresh_ClearsStatusCache` (cache seeded, then asserted empty) |
| FR-47 Retains parent-directory rescan safety net | Verified | `TestRunRefresh_ParentRescanSafetyNet` |
| FR-48 Summary of the specified shape | Verified | Live CLI compared line-for-line |
| FR-49 Conditional lines omitted at zero | Verified | `writeRefreshSummary` guards each line; live CLI omitted "Updated" when zero |
| FR-50 Totals derived from the same result in both commands | Verified | `TestInitAndRefresh_ReportIdenticalTotals` |
| FR-51 `refresh` reports errors exclusively via shared reporter | Verified | `TestInitAndRefresh_IdenticalScanErrorOutput` (7 errors, exercises truncation; stdout asserted clean) |
| **FR-52 Create a link only when unreachable AND previously tracked with a symlink** | **Verified** (after fix) | `TestRunRefresh_NoSymlinkForOrphanWorktreeMainRepo` (new), `TestRunRefresh_RepairsMissingLinkForTrackedExternalRepo`; combined live probe below shows repair, no-duplicate, and no-orphan-link all holding at once. See resolved Issue **V-1** |
| FR-53 No link for a repo discovered by traversing a symlink | Verified | `TestRunRefresh_NoDuplicateLinkForNestedSymlink` (3 refreshes); live CLI demo byte-identical |
| FR-54 Remove only a symlink whose target matches | Verified | `removeWorkspaceSymlink` unchanged from baseline (`git diff` shows no edit to its body) |
| FR-55 Repeated refresh produces no new symlinks or changes | Verified | `TestRunRefresh_IsIdempotent` (listing + `config.json` byte-compared); live CLI: identical after 3 refreshes |

#### Unit 3 — Documentation (spec L187–190)

| Requirement | Status | Evidence |
|---|---|---|
| FR-56 `commands-core.md` describes discovery, symlinks, worktrees, `scan_max_depth` | Verified | New `Discovery Rules` section (L331) linked from both commands |
| FR-57 Example output blocks match actual output | Verified | Line-for-line comparison with digits normalized; only singular/plural differ, correctly. See Observation O-2 |
| FR-58 `configuration.md` documents `scan_max_depth`, default 6, effect | Verified | Preferences table row plus dedicated subsection with worked example |
| FR-59 `README.md` synced | Verified | Discovery bullet rewritten; unified-reconciliation bullet added |

### Original Acceptance Criteria (from the initiating request)

| # | Criterion | Status | Evidence |
|---|---|---|---|
| 1 | `init` and `refresh` use the same boundary-aware discovery rules | Verified | Both call `ReconcileWorkspace`; `TestInitAndRefresh_ReportIdenticalTotals` |
| 2 | `refresh` re-scans recursively, not one level | Verified | `TestRunRefresh_DiscoversDeeplyNestedRepository` (depth 4) |
| 3 | `refresh` discovers repositories added after init | Verified | Same test; live CLI found `late-arrival` |
| 4 | `refresh` removes repositories whose paths are gone | Verified | `TestRunRefresh_RemovesAndCountsDeletedRepository` |
| 5 | `refresh` immediately discovers new native worktrees | Verified | `TestRunRefresh_NewWorktreeVisibleImmediately` |
| 6 | `worktree list` correct after either command, no extra step | Verified | `TestRunInit_WorktreesVisibleImmediately` + live CLI |
| 7 | Curated symlink workspaces keep working | Verified | `TestRunInit_CuratedSymlinkWorkspace`, `TestRunRefresh_IsIdempotent` |
| 8 | Totals reported consistently by both commands | Verified | `TestInitAndRefresh_ReportIdenticalTotals` |
| 9 | Tags survive for repositories whose path remains valid | Verified | `TestMergeRepositories_PreservesTags`; live CLI tag survived refresh |

### Repository Standards

| Standard Area | Status | Evidence |
|---|---|---|
| Coding standards | Verified | `make lint` exit 0 (`golangci-lint v1.64.8`, gosec/errcheck/unparam/govet enable-all); `make vet` exit 0 |
| Package layout | Verified | Commands in `cmd/git-workspace/`, logic in `internal/<pkg>/`; new `internal/reconcile` follows the pattern |
| Cobra conventions | Verified | `initCmd`/`refreshCmd` keep `Use`/`Short`/`Long`+Examples and thin `RunE` delegating to `run<Name>` |
| Testing patterns | Verified | Table-driven tests, `t.TempDir()` fixtures, `io.Writer` injection matching `runWorktreeList`; 387 tests pass under `-race` |
| Test isolation | Verified | Every config-touching test redirects `HOME` via `t.Setenv` (`setupWorkspace`, `withTempHome`, `newWorkspaceFixture`). Developer's real `~/.gws/config.json` untouched — mtime still 2026-08-09, 43 repos intact |
| Error handling | Verified | `fmt.Errorf(... %w ...)` throughout; deliberate ignores use `_` |
| Commit conventions | Verified | All 7 commits conventional and type-padded per `.githooks/commit-msg` |
| Documentation sync (spec 18 standard) | Verified | `mkdocs build --strict` exit 0; all 6 cross-reference anchors resolve |
| Security | Verified | GATE F clean — no credential patterns, no emails, no real home paths in proof artifacts |

### Proof Artifacts

| Task | Proof Artifact | Status | Verification Result |
|---|---|---|---|
| 1.0 | `21-task-01-proofs.md` + 8 named tests | Verified | File exists; all tests exist in source and pass on cleared cache |
| 2.0 | `21-task-02-proofs.md` + 20 named tests | Verified | All pass; empirical `rev-parse` table reproduced independently |
| 3.0 | `21-task-03-proofs.md` + 32 named tests | Verified | All pass; `internal/reconcile` purity re-verified by grep |
| 4.0 | `21-task-04-proofs.md` + 10 named tests, CLI transcript | Verified | Tests pass; CLI transcript reproduced with the built binary |
| 5.0 | `21-task-05-proofs.md` + 17 named tests, CLI transcript | Verified | Tests pass; accumulation-fix demo reproduced byte-identical |
| 6.0 | `21-task-06-proofs.md` + strict build | Verified | `mkdocs build --strict` exit 0; anchors verified against headings |
| All | 107 distinct test names claimed across proofs | Verified | 0 missing from source; 107/107 PASS when executed by name on a cleared cache; 0 FAIL, 0 SKIP |

---

## 3) Validation Issues

### V-1 · HIGH · ✅ RESOLVED · `gws refresh` created a symlink the spec says it must not

**Location:** `cmd/git-workspace/refresh.go` — `workspaceLinkOwnedRepos` and `repairWorkspaceSymlinks`

**Requirement violated:** FR-52 (spec L180) — "`gws refresh` shall create a workspace symlink for a repository only when that repository has no reachable path under the workspace root **and was previously tracked with a workspace symlink**."

**Also contradicts** the spec's Technical Considerations (L238): "Orphan worktree main repositories may live outside the workspace root. They are tracked by real path and, **per the repair-only rule, receive no workspace symlink**."

**Evidence — reproduced with the built binary.** A workspace containing only a worktree whose main repository lives outside it:

```text
$ gws init <ws>                    # main repo at <ext>/proj registered via FR-31
$ ls <ws>
proj.wt                            # correct: no symlink

$ gws refresh
$ ls <ws>
proj -> /…/probe/ext/proj          # symlink created — never previously linked
proj.wt
```

**Root cause.** `workspaceLinkOwnedRepos` treats *tracked + outside the workspace root* as proof that gws previously created a link. That inference holds for repositories added by `gws add`, but not for main repositories auto-registered by orphan-worktree handling (FR-31), which are tracked and external yet never linked. Those satisfy all three implemented conditions on the next refresh.

**Impact:** Functionality. gws mutates the user's curated workspace directory with a link they never requested — the precise class of behavior round-2 Q4-C set out to prevent. Self-limiting: exactly one link, then stable (verified stable across 4 refreshes), and no accumulation loop.

*Note: this was flagged during task 5.0 as an interpretation call on the task-5.5 wording and surfaced for review. This validation confirmed the deviation was observable, not merely theoretical.*

#### Resolution (commit `f107744`)

`reconcile.Result` now carries `WorktreeMainRepos`, the set of main repository paths inferred from linked worktrees during the pass, resolved to real paths. `repairWorkspaceSymlinks` gained a third condition: a repository listed there is tracked because one of its worktrees lives in the workspace, not because the user asked for it, so gws does not link it.

```go
// cmd/git-workspace/refresh.go
if result.ReachablePaths[repo.Path] { continue }   // already findable
if !linkOwned[repo.Path]           { continue }    // not ours to repair
if result.WorktreeMainRepos[repo.Path] { continue } // registered on gws's own behalf
ensureWorkspaceSymlink(cfg, repo.Path, repo.Name)
```

**Regression test:** `TestRunRefresh_NoSymlinkForOrphanWorktreeMainRepo` asserts that an orphan main repository outside the workspace gains no symlink after `init` and across three refreshes, while remaining tracked with its worktree attached.

**Verified with the real binary** — all three symlink behaviors holding simultaneously in one workspace:

```text
workspace: added -> ext/added        (user-added external, linked)
           grouping/nested -> ext/nested   (reached via nested symlink)
           orphan.wt/feature               (worktree; main repo outside)

$ rm ws/added && gws refresh
  added link repaired?  YES   (want YES)
  nested duplicated?    NO    (want NO)
  orphan linked?        NO    (want NO)

$ gws refresh ×2  →  workspace unchanged; 3 repositories tracked; worktree still listed
```

**Accepted trade-off**, documented in the function comment: a user-added external repository that *also* has a worktree inside the workspace will not have a deleted link auto-repaired, because it is indistinguishable from the orphan case by path alone. gws declines to act rather than risk creating a link it does not own. Recording link ownership explicitly in `config.json` would remove the ambiguity, but that is a schema change and was out of scope here.

### V-2 · MEDIUM · ✅ RESOLVED · Three changed files were neither declared nor justified in commit messages

**GATE D.** Files changed outside the task list's "Relevant Files":

| File | Justified in a commit message? |
|---|---|
| `cmd/git-workspace/deprecated.go` | No |
| `cmd/git-workspace/reconcilefixture_test.go` | No |
| `internal/reconcile/testhelpers_test.go` | No |
| `docs/site/getting-started.md` | Yes — commit `bc64772` |

**Evidence:** `comm -13` of declared vs changed file lists; `git log --format='%h: %s%n%b'` searched for each filename.

All three changes are benign and reviewable: `deprecated.go` is a two-line mechanical signature propagation (`runInit("")` → `runInit("", os.Stdout)`, same for `runRefresh`), and the two `*_test.go` files are shared fixture helpers described in proof files 04 and 03 respectively. The gate failure is a bookkeeping gap, not scope creep.

**Impact:** Traceability. A reviewer reading only the task list would not expect these files in the diff.

#### Resolution (commit `f107744`)

All three were added to the task list's "Relevant Files" with descriptions, alongside `docs/site/getting-started.md` for completeness. Re-checked:

```text
$ comm -13 <declared> <changed>
(empty — every changed code/doc file is now declared)
```

The only remaining asymmetry is two files declared but deliberately unchanged (`add_test.go`, `userdetect.go`), covered by Observation O-3.

### V-3 · LOW · Spec Open Questions remain unresolved

Spec Open Questions 1 and 2 were never closed:

1. *Should the `Worktrees: N (X aligned, Y unaligned)` line appear in `refresh` as well as `init`?* — implemented in **both**, on the reasoning that acceptance criterion 8 requires consistent totals. Not confirmed by the spec owner.
2. *Does adding the optional `scan_max_depth` preference warrant a `ConfigVersion` bump from `1.1.0`?* — assumed **no**; `ConfigVersion` is unchanged. The field is additive and `omitempty`, and `TestScanMaxDepthBackwardCompatibility` proves older configs load unchanged, so no defect follows — but the decision is unrecorded.

**Recommendation.** Record both decisions in the spec, or close them explicitly before merge.

---

## Observations (non-blocking)

**O-1 · FR-38 "unconditionally".** Both commands pass a progress reporter with no repository-count threshold, satisfying round-2 Q8-D. However `git.Progress.Start()` renders only when stderr is a TTY (`progress.go:36-39`) — pre-existing behavior of the facility FR-37 mandates reusing. So progress is unconditional with respect to workspace size, but still absent in non-interactive contexts. Reading FR-38 as "no size threshold" makes this compliant; reading it as "always renders" does not.

**O-2 · FR-57 example numbers.** Every documented output *line format* was verified verbatim against the built binary. The counts in the examples (15 repositories, etc.) are illustrative of a realistic workspace, consistent with surrounding documentation, rather than transcribed from a single run.

**O-3 · Two declared files unchanged.** `cmd/git-workspace/add_test.go` was declared for verification only ("confirm `add --recursive` still behaves") — its 12 tests pass unmodified. `cmd/git-workspace/userdetect.go` was declared "No change expected" and correctly went unchanged. Neither is a defect.

**O-4 · `gws add --recursive` inherited a behavior change**, as anticipated in the task-list scope note: it now prunes at repository boundaries, skips hidden directories, and honors `scan_max_depth`. The spec lists `gws add` as a non-goal; a live smoke test confirmed it still discovers and adds correctly, and all 12 `add` tests pass unchanged.

---

## 4) Evidence Appendix

### Commits analyzed

```text
bc64772 docs:       sync init and refresh documentation with unified discovery
f1b572d fix:        unify gws refresh discovery and stop symlink accumulation
a948ed7 feat:       move gws init onto the shared reconciliation engine
a835a47 feat:       add shared workspace reconciliation engine
9bd4076 feat:       resolve symlinks and detect worktrees structurally
0143771 feat:       add boundary-aware recursive scanner with depth cap
3a56b8a docs:       add spec 21 for unified workspace reconciliation
```

Each commit maps to exactly one parent task and carries a `Related to T<n>.0 in Spec 21` trailer. Progression is coherent: scanner → symlinks/worktrees → engine → init → refresh → docs.

### Commands executed

```text
$ go clean -testcache && go test ./...
ok  cmd/git-workspace 1.103s   ok internal/classifier 0.004s
ok  internal/config   0.004s   ok internal/discovery  0.126s
ok  internal/filter   0.006s   ok internal/git        0.436s
ok  internal/reconcile 0.624s  ok internal/user       0.059s

$ go test ./... -run "^(<107 claimed test names>)$" -v | grep -c "^--- PASS"
107
$ ... | grep -E "^--- (FAIL|SKIP)"
none

$ make vet    → exit 0
$ make lint   → exit 0 (golangci-lint v1.64.8, no findings)
$ mkdocs build --strict → exit 0

$ git diff --shortstat c4eccde..HEAD -- '*.go'
25 files changed, 3786 insertions(+), 322 deletions(-)
```

### Proof-artifact source verification

All 107 test names appearing as `--- PASS` in the six proof files were extracted and checked against the source tree: **0 missing**. All were then executed explicitly by name on a cleared cache: **107 PASS, 0 FAIL, 0 SKIP**. No proof artifact cites a test that does not exist or does not pass.

### Live CLI verification performed independently of proof files

| Check | Result |
|---|---|
| `gws init` on nested + symlinked + worktree workspace | Correct summary; 3 repositories, 2 worktrees (1 aligned, 1 unaligned) |
| `gws worktree list` immediately after `init` | Both worktrees listed with correct alignment |
| `gws init` on already-initialized workspace | Reworded guard, `refresh` listed first, exit 0 |
| `gws refresh` after deletion + deep addition | Removed 1, found 1 new, correct totals |
| Symlink accumulation demo (init + 3 refreshes) | Workspace byte-identical; no duplicate link |
| Orphan main-repo probe | Initially failed (Issue V-1); after fix, no symlink across 3 refreshes |
| Smoke: `list`, `tag add`, `list -t`, `refresh`, `add -r` | All correct; tag survived refresh |
| Developer's real `~/.gws/config.json` | Untouched (mtime 2026-08-09, 43 repos) |

### Gate results

| Gate | Initial | Final | Notes |
|---|---|---|---|
| **A** — No CRITICAL/HIGH | FAIL | **PASS** | V-1 resolved in `f107744` and re-verified |
| **B** — No `Unknown` in matrix | PASS | **PASS** | 59 Verified, 0 Failed, 0 Unknown |
| **C** — Proof artifacts accessible and functional | PASS | **PASS** | 6/6 files present; 107/107 test claims pass |
| **D** — Changed files declared or justified | FAIL | **PASS** | V-2 resolved; zero undeclared files |
| **E** — Repository standards followed | PASS | **PASS** | vet, lint, test-race, conventions, docs build |
| **F** — No credentials in proof artifacts | PASS | **PASS** | Credential, email, and home-path scans clean |

### Post-fix re-verification

```text
$ go clean -testcache && go test ./...          → all 8 packages ok (388 tests)
$ make test-race                                 → all 8 packages ok
$ make vet                                       → exit 0
$ make lint                                      → exit 0
$ comm -13 <declared> <changed>                  → empty (GATE D)
$ live 3-case symlink probe                      → repair YES / duplicate NO / orphan NO
```

---

## Remaining Non-Blocking Item

**V-3** — record decisions on the two spec Open Questions (the `Worktrees:` line on `refresh`, and whether `ConfigVersion` should bump for the additive `scan_max_depth` preference). Neither affects correctness; both are currently implemented under stated assumptions.

The implementation is ready for final human code review and merge.

---

**Validation Completed:** 2026-08-10
**Validation Performed By:** Claude Opus 5 (claude-opus-5)
