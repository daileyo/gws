# Spec 21 — Task 6.0 Proofs: Documentation Sync

**Task:** 6.0 Documentation Sync
**Branch:** `feat/workspace-reconciliation`
**Date:** 2026-08-10

## Summary of Changes

```text
$ git diff --stat docs/site/ README.md
 README.md                    |   3 +-
 docs/site/commands-core.md   | 110 +++++++++++++++++++++++++++++++++++++++----
 docs/site/configuration.md   |  22 ++++++++-
 docs/site/getting-started.md |   5 +-
 4 files changed, 127 insertions(+), 13 deletions(-)
```

| File | Change |
|---|---|
| `docs/site/commands-core.md` | Rewrote `Initialize Workspace` and `Refresh Workspace`; added a shared `Discovery Rules` section; documented workspace symlink repair and scan warnings; fixed the classification note |
| `docs/site/configuration.md` | Added `scan_max_depth` to the example config, the Preferences table, and a dedicated subsection |
| `docs/site/getting-started.md` | Corrected the "what init does" list — it omitted worktree discovery and symlink following |
| `README.md` | Updated the discovery feature bullet; added a unified-reconciliation bullet |

## New: Shared `Discovery Rules` Section

Rather than describe traversal twice and risk the two descriptions drifting — the documentation equivalent of the bug this spec fixes — both commands now link to one `Discovery Rules` section covering traversal, symlinks, and the repository-versus-worktree distinction.

## Build Verification

`make docs` starts a blocking MkDocs **dev server**, so it cannot serve as a build check. A strict build was run instead, which additionally turns warnings such as broken internal links into errors:

```text
$ mkdocs build --strict --site-dir <tmp>
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: <tmp>
INFO    -  Documentation built in 0.19 seconds
EXIT=0
```

### Cross-reference anchors verified

| Reference | Target heading | Resolves |
|---|---|---|
| `commands-core.md` → `#discovery-rules` | `## Discovery Rules` (commands-core.md:331) | ✅ |
| `commands-core.md` → `#refresh-workspace` | `## Refresh Workspace` (commands-core.md:403) | ✅ |
| `commands-core.md` → `#add-repository` | `## Add Repository` (commands-core.md:366) | ✅ |
| `commands-core.md` → `configuration.md#preferences-fields` | `### Preferences Fields` (configuration.md:120) | ✅ |
| `configuration.md` → `commands-core.md#discovery-rules` | `## Discovery Rules` | ✅ |
| `getting-started.md` → `commands-core.md#refresh-workspace` | `## Refresh Workspace` | ✅ |

## Documented Output Verified Against the Real Binary

Each documented example was compared line-for-line against output from the built binary, with digits normalized to `N` so only the line shapes are compared.

### `gws init`

```text
=== documented ===                          === real binary ===
Initialized workspace at: /home/user/…      Initialized workspace at: WORKSPACE
Found N repositories.                       Found N repositories.
Repositories with user configuration: N     Repositories with user configuration: N
Repositories with worktrees: N              Repositories with worktrees: N
Worktrees: N (N aligned, N unaligned)       Worktrees: N (N aligned, N unaligned)
```

### `gws refresh`

```text
=== documented ===                          === real binary ===
Refreshing workspace at: /home/user/…       Refreshing workspace at: WORKSPACE
Detecting git user configuration...         Detecting git user configuration...
Cleared git status cache                    Cleared git status cache
                                            
Refresh complete!                           Refresh complete!
Total repositories: N                       Total repositories: N
Removed N repository (path no longer…)      Removed N repository (path no longer…)
Found N new repositories                    Found N new repository
Updated N repositories                      Updated N repository
Repositories with user configuration: N     Repositories with user configuration: N
Repositories with worktrees: N              Repositories with worktrees: N
Worktrees: N (N aligned, N unaligned)       Worktrees: N (N aligned, N unaligned)
```

Every line matches. The only differences are singular/plural word forms, which is the `pluralize` helper working correctly: the documented example describes a workspace with 2 new and 3 updated repositories, while the verification run had exactly 1 of each.

**Note on example numbers:** the counts in the documented examples are illustrative of a realistic workspace (15 repositories), consistent with the surrounding documentation, rather than transcribed from the small verification fixture. Every *line format* is verbatim from real output, as shown above.

## Quality Gates

```text
$ make lint
Running linter...
(exit 0, no findings)

$ make test-race
ok  	github.com/daileyo/gws/cmd/git-workspace
ok  	github.com/daileyo/gws/internal/classifier
ok  	github.com/daileyo/gws/internal/config
ok  	github.com/daileyo/gws/internal/discovery
ok  	github.com/daileyo/gws/internal/filter
ok  	github.com/daileyo/gws/internal/git
ok  	github.com/daileyo/gws/internal/reconcile
ok  	github.com/daileyo/gws/internal/user
```

## Verification

| Required Proof Artifact | Evidence | Status |
|---|---|---|
| `Initialize Workspace` / `Refresh Workspace` describe recursive boundary-aware discovery, symlinks, worktrees, `scan_max_depth` | Both sections rewritten; shared `Discovery Rules` section covers all four | ✅ PASS |
| Example output blocks match actual command output | Line-for-line comparison above against the built binary | ✅ PASS |
| `scan_max_depth` documented with default `6` and its effect | `configuration.md` Preferences table plus a dedicated subsection with a worked example | ✅ PASS |
| `README.md` consistent with the docs site | Discovery bullet rewritten; unified-reconciliation bullet added | ✅ PASS |
| Site builds without errors or broken links | `mkdocs build --strict` exit 0; all six cross-references verified against their headings | ✅ PASS |

## Additional Correction Found

`docs/site/getting-started.md` was not named in the task list, but its "This command will:" list for `gws init` was inaccurate after this spec — it omitted worktree discovery and symlink following. Since spec 18 established that the docs site must stay in sync with CLI behavior, it was corrected here rather than left as a known-wrong page.

`docs/site/commands-core.md` also stated that repositories are classified "when repositories are discovered with `gws init`". Both commands classify, so it now says `gws init` or `gws refresh`.
