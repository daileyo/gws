# 23 Task 5.0 Proofs - Documentation

## Files

- `docs/site/configuration.md` — "What is a project?", File Locations, resolution order, cross-platform section, Windows path-length note, Upgrading
- `docs/site/commands-core.md` — worktree layout, aligned/unaligned wording, sample paths
- `docs/site/getting-started.md` — config path
- `docs/site/shell-integration.md` — sample worktree path
- `README.md` — feature bullet and project structure

## Proof 1: "project" is defined once, plainly (task 5.4)

> A **project** is a git repository plus any worktrees associated with it. The repository is
> the main checkout, wherever you keep it; its worktrees live together under the projects root.

## Proof 2: locations and resolution order (task 5.1)

`configuration.md` documents both paths and their defaults, and states that a relative value in
either variable is treated as unset — "otherwise the location would depend on your current
working directory", which is the actual reason rather than a bare rule.

## Proof 3: cross-platform parity and its precedent (tasks 5.1, 5.7)

The docs show the Windows paths explicitly and justify them:

> Windows has no XDG specification, but this is not an invention — **git does the same thing**.
> Per `git-config(1)`, when `XDG_CONFIG_HOME` is unset git uses `$HOME/.config`, and Git for
> Windows sets `$HOME` to `%USERPROFILE%`.

It also documents the escape hatch for users who prefer `%AppData%`: the XDG variables are
honored on Windows too.

## Proof 4: Windows path length (task 5.6)

An admonition covers the 260-character `MAX_PATH` limit, giving both the `LongPathsEnabled`
registry command and the simpler alternative of pointing `XDG_DATA_HOME` at a short path.

## Proof 5: upgrade path (task 5.2)

The Upgrading section states the asymmetry plainly — config moves itself, worktrees do not —
and explains *why*: a worktree can hold uncommitted work. It shows the actual migration output,
shows `(unaligned)` as the expected post-upgrade state rather than an error, and gives the
`--dry-run` then `align` sequence. It closes with the cross-filesystem warning.

## Proof 6: stale paths retired (tasks 5.3, 5.5)

```
$ grep -rn '\.wt/\|~/\.gws' docs/site/*.md README.md | grep -vi 'legacy\|Earlier versions\|old location'
(no matches)
```

Every remaining reference is inside the Upgrading section or the sentence pointing users there,
where naming the old layout is the point.

## Proof 7: the docs site builds

```
$ mkdocs build --strict
INFO - Cleaning site directory
INFO - Building documentation to directory: .../site
INFO - Documentation built in 0.21 seconds
```

`--strict` turns warnings into errors, so the new cross-references
(`configuration.md#file-locations`, `configuration.md#upgrading-from-earlier-versions`) all
resolve.

## Correction made during this task

The Upgrading section originally read "Versions before 2.21 kept the config at `~/.gws`".
That was wrong: 2.21.0 is the release that shipped `gws cd` (spec 22), and the XDG relocation
is not in it. Changed to "Earlier versions", since the release number for this work is not yet
known.
