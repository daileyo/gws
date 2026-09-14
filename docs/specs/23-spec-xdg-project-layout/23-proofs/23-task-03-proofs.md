# 23 Task 3.0 Proofs - Worktree Root Relocation

## Files

- `internal/git/worktree.go` — `IsAligned` rekeyed to the projects root; `resolvePath` helper
- `internal/git/worktree_test.go` — `TestIsAligned` rewritten
- `cmd/git-workspace/worktree_add.go` — destination is the projects root
- `cmd/git-workspace/worktree_align.go` — both `wtDir` derivations relocated
- `cmd/git-workspace/refresh.go` — `IsAligned` call updated
- `cmd/git-workspace/projects_path_test.go` (new) — `projectsPath` test helper
- `worktree_add_test.go`, `worktree_align_test.go`, `refresh_test.go` — fixtures no longer build `.wt` paths

## Proof 1: worktrees land in the projects root

```
$ git-workspace worktree add demo-repo feat-x
Created worktree for branch 'feat-x' at /tmp/.../.local/share/gws/projects/demo-repo/feat-x

$ find $HOME/.local/share/gws/projects -mindepth 2 -maxdepth 3 -type d
.../projects/demo-repo/feature
.../projects/demo-repo/feature/auth      <- slashes still nest
.../projects/demo-repo/feat-x

$ ls $HOME/ws
demo-repo                                 <- no .wt sibling; the workspace stays clean
```

## Proof 2: `.wt` is no longer aligned

`TestIsAligned` case `unaligned - legacy .wt directory` asserts
`/workspace/my-repo.wt/feature-x` is **not** aligned. This is the hinge of task 4.0: `align`
picks up legacy worktrees with no special-casing, purely because the predicate changed.

Other cases: the repo's own projects dir (aligned), nested branch paths (aligned), another
repo's projects dir (unaligned), a sibling with a similar prefix (unaligned).

## Proof 3: no `.wt` left in non-test source

```
$ grep -rn '\.wt' cmd/git-workspace/*.go internal/**/*.go | grep -v _test
(no matches)
```

Help text in `worktree.go`, `worktree_add.go`, `worktree_align.go`, and `worktree_list.go` now
describes the projects root, and `worktree add --help` prints the resolved default.

## Proof 4: full suite

```
$ go test -count=1 ./...      # 8/8 packages ok
$ go vet ./...                # clean
```

## The risk this task carried

`IsAligned`'s signature changed from `(worktreePath, repoPath)` to
`(worktreePath, repoName)` — **both parameters are strings**, so the compiler accepted every
existing call site unchanged:

```
$ go build ./...              # succeeded with all three call sites still passing repo.Path
```

The build passing proved nothing. All three (`worktree_add.go:86`, `refresh.go:130`,
`worktree_align.go:208`) were located by grep and updated by hand. Had they been missed,
`IsAligned` would have compared worktree paths against a projects directory named after a
filesystem path — every worktree silently unaligned, and `align` would have tried to move all
of them.

A `type RepoName string` would have made this class of change compiler-checked. Not adopted
here: it would ripple through `config.Repository` and well beyond this spec's scope. Recorded
as a note rather than a silent judgment call.

## Fixture change worth noting

`refresh_test.go` previously created its worktree in `<repo>.wt/` and asserted
`Aligned == true`. Under the new semantics that combination is contradictory. The fixture now
creates the worktree in the projects root and still asserts aligned, so the test keeps
verifying discovery rather than being weakened to accommodate the change.
