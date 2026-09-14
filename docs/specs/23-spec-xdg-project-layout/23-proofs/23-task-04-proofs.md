# 23 Task 4.0 Proofs - Legacy `.wt` Migration via `align`

## Files

- `internal/git/worktree.go` — cross-device branch in `MoveWorktree`, `isCrossDeviceErr`
- `internal/git/worktree_test.go` — `TestIsCrossDeviceErr`
- `cmd/git-workspace/worktree_align.go` — relocation notice, `isLegacyWtPath`, `hasLegacyPlan`, `removeEmptyLegacyDir`
- `cmd/git-workspace/worktree_align_test.go` — legacy migration, conservative cleanup, path classification

## Proof 1: a legacy layout migrates end to end

Starting from a worktree in the pre-XDG `<repo>.wt/` location:

```
$ gws worktree list
REPO         BRANCH    PATH                                   STATUS
legacy-repo  feat-old  .../ws/legacy-repo.wt/feat-old         (unaligned)
```

It reports **unaligned** with no migration-specific code — `IsAligned` was rekeyed in task
3.0 and this simply falls out.

```
$ gws worktree align --dry-run
Dry run — no changes will be made:

Worktrees now live in the projects root: /tmp/.../.local/share/gws/projects

Would move [legacy-repo] feat-old
  from: /tmp/.../ws/legacy-repo.wt/feat-old
  to:   /tmp/.../.local/share/gws/projects/legacy-repo/feat-old

Total: 1 worktree to align

$ ls ~/ws
legacy-repo
legacy-repo.wt                      <- husk still present; dry-run changed nothing

$ gws worktree align
Worktrees now live in the projects root: /tmp/.../.local/share/gws/projects

Moving [legacy-repo] feat-old
  from: /tmp/.../ws/legacy-repo.wt/feat-old
  to:   /tmp/.../.local/share/gws/projects/legacy-repo/feat-old

Removed empty /tmp/.../ws/legacy-repo.wt
Aligned 1 worktree

$ ls ~/ws
legacy-repo                         <- husk gone; the home directory is quieter

$ gws worktree list
REPO         BRANCH    PATH                                              STATUS
legacy-repo  feat-old  .../.local/share/gws/projects/legacy-repo/feat-old  aligned
```

## Proof 2: no migration-specific code (task 4.1)

`TestRunWorktreeAlign_MigratesLegacyWtLayout` asserts the whole path: fixture starts in
`<repo>.wt/`, `align` moves it to the projects root, the husk is removed, and config records
the new path with `Aligned: true`. No branch in `align` special-cases `.wt` as a *source* —
`isLegacyWtPath` is used only to decide whether to print the one-time notice and whether the
emptied directory is a husk worth removing.

## Proof 3: cleanup is conservative (task 4.3)

`TestRunWorktreeAlign_LegacyDirWithOtherContentKept` puts a `notes.txt` in the legacy
directory. After `align`, the file still exists and the directory survives. `removeEmptyLegacyDir`
removes only a directory that `ReadDir` reports empty, and only one whose name ends in `.wt`.

## Proof 4: cross-device moves report usefully (task 4.2)

`git worktree move` is ultimately a rename and cannot cross filesystems. This went from
implausible to likely: the projects root is under `$XDG_DATA_HOME` while the repo may be on a
different mount.

```
$ go test ./internal/git/ -run TestIsCrossDeviceErr -v
--- PASS: TestIsCrossDeviceErr/wrapped_EXDEV
--- PASS: TestIsCrossDeviceErr/git_message
--- PASS: TestIsCrossDeviceErr/lowercase_variant
--- PASS: TestIsCrossDeviceErr/unrelated_failure
```

Both forms are detected because git reports this as message text, not a typed error; the
wrapped `syscall.EXDEV` is checked as well. The resulting message names both paths and
suggests setting `XDG_DATA_HOME` to the same filesystem, rather than surfacing git's raw output.

## Proof 5: existing behavior preserved (tasks 4.4-4.6)

| Behavior | Evidence |
| --- | --- |
| `--dry-run` changes nothing | Proof 1: husk still present after dry-run |
| Locked worktrees skipped | existing `IsWorktreeLocked` branch untouched |
| `Worktree.Path`/`Aligned` updated | `TestRunWorktreeAlign_MigratesLegacyWtLayout` asserts both |
| `-dup-NN` conflict handling | untouched; `TestRunWorktreeAlign_NameConflict` still passes |

## Proof 6: full suite

```
$ go test -count=1 ./...      # 8/8 packages ok
$ go vet ./...                # clean
```

## Output ordering fix

The relocation notice initially printed *after* the "Moving" lines, so the explanation arrived
after the thing it explained. It now prints before the plan listing, and in dry-run mode too —
so a user previewing the change sees where their worktrees are headed before deciding.
