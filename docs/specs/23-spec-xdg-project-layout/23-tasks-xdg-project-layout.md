# 23 Tasks - XDG Project Layout

## Relevant Files

- `internal/xdg/xdg.go` - **New.** Authoritative resolution of `ConfigDir()`, `ConfigFile()`, `ProjectsDir()` — identical layout on every platform, XDG variables honored everywhere including Windows.
- `internal/xdg/xdg_test.go` - **New.** Resolution tests across set / unset / relative env vars, plus a cross-platform parity assertion.
- `internal/config/config.go` - `GetConfigPath()` and `GetConfigDir()` currently build `~/.gws` directly (lines ~88-105); `Load()`/`Save()` use them. Migration logic lands here. `ConfigVersion` const is here.
- `internal/config/config_test.go` - Config load/save tests; migration cases go here.
- `internal/git/worktree.go` - `IsAligned()` at line ~215 derives `repoPath + ".wt"`. `MoveWorktree()` at line ~113 needs the cross-device branch.
- `internal/git/worktree_test.go` - `TestIsAligned` and move tests.
- `cmd/git-workspace/worktree_add.go` - `wtDir := repo.Path + ".wt"` at line 65; destination construction and `IsAligned` call.
- `cmd/git-workspace/worktree_align.go` - `wtDir` at lines 77 and 166; the `-dup-NN` conflict logic at lines 110-119 stays as-is; `IsAligned` call at line 208.
- `cmd/git-workspace/refresh.go` - `IsAligned` call at line 130 during discovery.
- `cmd/git-workspace/worktree.go` - Help text referencing `<repo>.wt/` at lines 29-30.
- `cmd/git-workspace/worktree_list.go` - Help text at line 22; `(unaligned)` display at line 107.
- `docs/site/configuration.md` - Config path documentation.
- `docs/site/commands-core.md` - Command reference with `.wt/` references.
- `README.md` - Layout and description sections.

### Notes

- Run tests with `go test ./...`. Build with `make build`.
- **Tests must never touch the real `~/.config` or `~/.local/share`** — set `HOME` and the `XDG_*` variables to a `t.TempDir()` in every test that exercises path resolution.
- `IsAligned`'s signature changes from `(worktreePath, repoPath)` to a repo-name-keyed form. All four call sites must be updated in the same change or the build breaks.
- The `-dup-NN` collision logic is explicitly out of scope for modification — only the base directory it joins against changes.
- Follow conventional commits: `feat(config): ...`, `feat(worktree): ...`, `docs: ...`.

## Tasks

### [ ] 1.0 XDG Path Resolution Package

Create `internal/xdg` as the single source of truth for where gws keeps its files, correct on both Unix and Windows, and repoint `internal/config` at it.

#### 1.0 Proof Artifact(s)

- Test: `go test ./internal/xdg/` passes covering set, unset, and relative `XDG_CONFIG_HOME`/`XDG_DATA_HOME` plus platform defaults demonstrates resolution is correct
- CLI: `XDG_CONFIG_HOME=/tmp/xdgc git-workspace init` writes `/tmp/xdgc/gws/config.json` demonstrates the config override is honored
- CLI: On Windows, `git-workspace init` writes `C:\Users\<user>\.config\gws\config.json` demonstrates the Windows layout matches Unix
- CLI: `XDG_DATA_HOME=/tmp/xdgd git-workspace worktree add <repo> demo` creates `/tmp/xdgd/gws/projects/<repo>/demo` demonstrates the data override is honored
- Test: `go test ./internal/config/` passes after delegation demonstrates no regression

#### 1.0 Tasks

- [ ] 1.1 Create `internal/xdg/xdg.go` with `ConfigDir() (string, error)`: return `$XDG_CONFIG_HOME/gws` when the variable is set and absolute, else `os.UserHomeDir()/.config/gws`. Do **not** use `os.UserConfigDir()` — it returns `%AppData%` on Windows and would break cross-platform parity
- [ ] 1.2 Add `ConfigFile() (string, error)` returning `ConfigDir()/config.json`
- [ ] 1.3 Add `ProjectsDir() (string, error)`: return `$XDG_DATA_HOME/gws/projects` when set and absolute, else `os.UserHomeDir()/.local/share/gws/projects`. No `runtime.GOOS` branch — the layout is identical on every platform, with `os.UserHomeDir()` supplying `%USERPROFILE%` on Windows
- [ ] 1.4 Treat a relative value in `XDG_CONFIG_HOME` or `XDG_DATA_HOME` as unset, per the XDG specification, using `filepath.IsAbs`
- [ ] 1.5 Add `LegacyConfigDir()` returning `~/.gws`, used only by the migration path in task 2.0
- [ ] 1.6 Create `internal/xdg/xdg_test.go` with table-driven cases: env set and absolute, env set but relative, env unset. Add a parity case asserting the resolved path relative to home is identical regardless of platform. Use `t.Setenv` and `t.TempDir()` so the real home directory is never read
- [ ] 1.7 Rewrite `config.GetConfigPath()` and `config.GetConfigDir()` to delegate to `xdg.ConfigFile()` and `xdg.ConfigDir()`, removing the hardcoded `.gws` join
- [ ] 1.8 Run `go build ./... && go test ./internal/...` and confirm the delegation compiles and existing config tests pass

### [ ] 2.0 Config Relocation and Automatic Migration

Move the config to its XDG home, migrating an existing `~/.gws/config.json` on first load, safely and exactly once.

#### 2.0 Proof Artifact(s)

- CLI: With `~/.gws/config.json` seeded and no XDG config, `git-workspace list` prints a migration notice naming both paths and the file now exists at the XDG location demonstrates automatic migration
- CLI: A second `git-workspace list` prints no notice demonstrates idempotency
- CLI: `stat -c %a ~/.config/gws/config.json` returns `600` demonstrates permissions are preserved
- CLI: `ls ~/.gws` after migration shows the directory removed demonstrates cleanup
- Test: `config_test.go` migration cases pass demonstrates the full contract

#### 2.0 Tasks

- [ ] 2.1 Add `migrateLegacyConfig() (migrated bool, err error)` to `internal/config`: return early if the XDG config already exists or the legacy file does not
- [ ] 2.2 Before moving, `os.Lstat` the legacy path and refuse to migrate if it is a symlink, so migration cannot be redirected outside the intended directory
- [ ] 2.3 Create the XDG config directory with `0755`, then move the file, preserving the `0600` permission; fall back to copy-then-remove if the rename crosses a filesystem
- [ ] 2.4 After a successful move, remove `~/.gws` only if it is empty; if other files remain, leave it and name them in the notice
- [ ] 2.5 Print a single-line notice to stderr naming the source and destination paths
- [ ] 2.6 Call the migration from `config.Load()` before the read, guarded by a `sync.Once` so it runs at most once per process
- [ ] 2.7 On migration failure, log a warning to stderr and fall back to reading the legacy path, so a user is never told their initialized workspace is uninitialized
- [ ] 2.8 Verify the `workspace not initialized` error still fires only when neither location has a config
- [ ] 2.9 Bump `ConfigVersion` from `1.1.0` to `1.2.0` to mark the layout change
- [ ] 2.10 Add `config_test.go` cases: migrate-on-load, no-op when XDG config exists, `0600` preserved, non-empty legacy dir retained, symlinked legacy config refused, and fallback-on-failure

### [ ] 3.0 Worktree Root Relocation

Repoint worktree creation and the alignment predicate at `ProjectsDir()/<repo-name>/`, updating every call site together.

#### 3.0 Proof Artifact(s)

- CLI: `gws worktree add my-repo feat-x` creates `~/.local/share/gws/projects/my-repo/feat-x` demonstrates the new destination
- CLI: `gws worktree add my-repo feature/nested` creates the nested path demonstrates slash handling is preserved
- CLI: `gws worktree list` marks a pre-existing `.wt/` worktree `(unaligned)` demonstrates the redefined predicate
- Test: `go test ./internal/git/ ./cmd/git-workspace/` passes demonstrates all call sites are consistent

#### 3.0 Tasks

- [ ] 3.1 Change `git.IsAligned` to take the repo name and test containment within `ProjectsDir()/<repo-name>/`, resolving symlinks on both sides as `ListWorktrees` already does
- [ ] 3.2 Update `TestIsAligned` in `internal/git/worktree_test.go` for the new signature and semantics, including a case asserting a `.wt/` path is now *not* aligned
- [ ] 3.3 In `worktree_add.go`, replace `wtDir := repo.Path + ".wt"` (line 65) with the projects-root destination and `os.MkdirAll` the repo-scoped parent
- [ ] 3.4 Update the `IsAligned` call in `worktree_add.go` (line ~86) for the new signature
- [ ] 3.5 In `worktree_align.go`, replace both `wtDir` derivations (lines 77 and 166) with the projects-root destination, leaving the `-dup-NN` logic at lines 110-119 untouched
- [ ] 3.6 Update the `IsAligned` calls in `worktree_align.go` (line ~208) and `refresh.go` (line 130)
- [ ] 3.7 Update help text and examples in `worktree.go` (lines 29-30), `worktree_add.go` (lines 17-20), `worktree_align.go` (lines 20-21), and `worktree_list.go` (line 22) to describe the projects root instead of `<repo>.wt/`
- [ ] 3.8 Run `go test ./...` and fix any test fixtures that construct `.wt` paths directly

### [ ] 4.0 Legacy `.wt` Migration via `align`

Make `gws worktree align` the single explicit gesture that relocates existing `.wt/` worktrees and cleans up the emptied directories.

#### 4.0 Proof Artifact(s)

- CLI: `gws worktree align --dry-run` on a repo with `.wt/` worktrees lists each planned move with source and destination demonstrates preview
- CLI: `gws worktree align` performs the moves and the emptied `<repo>.wt/` directory is gone demonstrates cleanup
- CLI: `gws worktree list` afterward shows no `(unaligned)` markers demonstrates alignment
- CLI: An align across a filesystem boundary reports a clear error naming both paths demonstrates the failure mode
- Test: `worktree_align_test.go` migration and cross-device cases pass demonstrates coverage

#### 4.0 Tasks

- [ ] 4.1 Confirm no special-casing is needed for `.wt/` sources — the redefined `IsAligned` should make them fall out as unaligned automatically. Add a test asserting this rather than writing migration-specific code
- [ ] 4.2 Add a cross-device branch to `MoveWorktree` in `internal/git/worktree.go`: detect the `EXDEV`/rename failure and return an actionable error naming source, destination, and the likely cause, rather than surfacing raw git output
- [ ] 4.3 After each successful move, if the source's parent `<repo>.wt/` directory is now empty, remove it
- [ ] 4.4 Verify `--dry-run` output shows the new destinations and does not create or remove anything
- [ ] 4.5 Confirm locked worktrees are still skipped with the existing message
- [ ] 4.6 Confirm `Worktree.Path` and `Worktree.Aligned` are updated in config after each move, as today
- [ ] 4.7 Add a one-time note to `align` output when it relocates worktrees out of a legacy `.wt/` directory, so the change of location is visible the first time
- [ ] 4.8 Add `worktree_align_test.go` cases: a simulated `.wt/` layout migrating to the projects root, emptied-directory cleanup, and the cross-device error

### [ ] 5.0 Documentation

Document the new layout, the XDG overrides, the Windows equivalents, and what to expect on upgrade.

#### 5.0 Proof Artifact(s)

- Documentation: `grep -rn '~/.gws\|\.wt/' docs/ README.md` returns only intentional upgrade/historical references demonstrates docs are current
- Documentation: configuration.md documents both XDG variables and their Windows fallbacks demonstrates coverage
- Documentation: An upgrade note explains the automatic config migration and the `align` step demonstrates users are prepared

#### 5.0 Tasks

- [ ] 5.1 Update `docs/site/configuration.md` with the XDG config path, the `XDG_CONFIG_HOME` / `XDG_DATA_HOME` overrides, and a resolution-order table. State that the layout is identical on all platforms and that Windows resolves `<home>` to `%USERPROFILE%`, citing git's own `$HOME/.config/git/config` behavior as the precedent
- [ ] 5.2 Add an "Upgrading from earlier versions" section explaining that the config migrates automatically, that existing worktrees will report as unaligned until `gws worktree align` is run, and that this is expected
- [ ] 5.3 Update `docs/site/commands-core.md` wherever it references `<repo>.wt/`
- [ ] 5.4 Define **project** once, plainly — a git repo plus any worktrees associated with it — in configuration.md or the docs index
- [ ] 5.5 Update `README.md` layout and description sections for the new paths
- [ ] 5.6 Add a note about Windows path length when `C:\Users\<user>\.local\share\gws\projects\` plus a deep branch name approaches the 260-character `MAX_PATH` limit, and how to enable long-path support
- [ ] 5.7 Document that `XDG_CONFIG_HOME` / `XDG_DATA_HOME` are honored on Windows, for users who prefer the native `%AppData%` location
