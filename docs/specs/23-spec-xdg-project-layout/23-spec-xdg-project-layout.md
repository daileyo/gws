# 23-spec-xdg-project-layout

## Introduction/Overview

gws currently scatters its own files across the user's home directory in two non-standard
places:

- **Config** lives at `~/.gws/config.json` — a dotdir in `$HOME`, predating any XDG awareness.
- **Worktrees** live at `<repo-path>.wt/`, a sibling directory next to every repo that has one.
  A user with ten repos that each have worktrees ends up with ten extra `*.wt` directories
  interleaved with their actual projects.

Both are noise, and neither follows the XDG Base Directory Specification that most modern CLI
tools now honor. This spec relocates both to standardized locations: config to
`$XDG_CONFIG_HOME/gws/`, and worktrees to `$XDG_DATA_HOME/gws/projects/`.

The scope is deliberately narrow. A **project** is a git repo plus any worktrees associated
with it, and that model is not changing. Discovery stays passive: `gws refresh` continues to
find worktrees via `git worktree list` wherever they are, continues to track them, and
continues to report each one as aligned or unaligned. The only thing that changes is **where
the aligned location is**.

## Goals

- Relocate the config file from `~/.gws/config.json` to `$XDG_CONFIG_HOME/gws/config.json` (default `~/.config/gws/config.json`), migrating existing installs automatically
- Relocate the worktree root from `<repo-path>.wt/` to `$XDG_DATA_HOME/gws/projects/<repo-name>/` (default `~/.local/share/gws/projects/`)
- Redefine "aligned" against the new projects root, so existing `.wt/` worktrees report as unaligned and `gws worktree align` migrates them
- Preserve the existing passive discovery and aligned/unaligned reporting behavior unchanged
- Preserve the existing `-dup-NN` collision handling unchanged
- Produce an identical directory layout on Linux, macOS, and Windows, mimicking XDG on Windows rather than diverging to `%AppData%`
- Never relocate a working tree without the user explicitly asking

## User Stories

- **As a gws user with many worktrees**, I want them collected under one standard directory so that my projects directory contains projects, not a `*.wt` sibling for every repo I have ever branched.
- **As a user who cares about a tidy home directory**, I want gws to follow the XDG Base Directory Specification like the rest of my tooling, so its files live where I expect and my backup and dotfile rules apply to them automatically.
- **As an existing user**, I want to upgrade without losing my configuration or having my working trees moved out from under me — the config should follow me automatically, and my worktrees should move only when I run `align`.
- **As a Windows user**, I want gws to live in the same relative place it does on my Linux and macOS machines, so that one set of documentation applies to all of them and my config is where git already keeps its own.
- **As a developer who works across platforms**, I want cross-platform parity to be a guarantee rather than a best effort, since moving between a work laptop and a WSL environment should not mean learning a second layout.

## Demoable Units of Work

### Unit 1: XDG Path Resolution

**Purpose:** Establish one authoritative place that answers "where does gws keep its config?" and "where does gws keep its worktrees?", producing **the same layout on every platform**, so no other code hardcodes a path and no user has to learn a second convention.

**Functional Requirements:**

- A new `internal/xdg` package shall expose `ConfigDir()`, `ConfigFile()`, and `ProjectsDir()`
- `ConfigDir()` shall resolve to `$XDG_CONFIG_HOME/gws` when `XDG_CONFIG_HOME` is set and absolute, otherwise `<home>/.config/gws`
- `ProjectsDir()` shall resolve to `$XDG_DATA_HOME/gws/projects` when `XDG_DATA_HOME` is set and absolute, otherwise `<home>/.local/share/gws/projects`
- **This resolution shall be identical on Linux, macOS, and Windows.** On Windows, `<home>` is `%USERPROFILE%`, producing `C:\Users\<user>\.config\gws` and `C:\Users\<user>\.local\share\gws\projects`
- The `XDG_CONFIG_HOME` and `XDG_DATA_HOME` environment variables shall be honored on Windows exactly as on Unix, so a developer who has already adopted the convention gets consistent behavior everywhere
- Home resolution shall use `os.UserHomeDir()`, which returns `%USERPROFILE%` on Windows and `$HOME` elsewhere
- `os.UserConfigDir()` shall **not** be used, because it returns `%AppData%` on Windows and would break the cross-platform uniformity this unit exists to provide
- Per the XDG specification, a relative value in either environment variable shall be treated as unset and the default used
- Resolution shall be pure and injectable so tests can drive it with a temporary home and explicit environment variables, without touching the real user directories
- `config.GetConfigPath()` and `config.GetConfigDir()` shall delegate to the new package rather than building `~/.gws` themselves

**Proof Artifacts:**

- Test: `internal/xdg/xdg_test.go` covering set/unset/relative `XDG_CONFIG_HOME` and `XDG_DATA_HOME` demonstrates resolution is correct
- Test: A case asserting the resolved *relative* layout is byte-identical across platforms demonstrates cross-platform uniformity
- CLI: `XDG_DATA_HOME=/tmp/xdgtest git-workspace worktree add <repo> demo` creates the worktree under `/tmp/xdgtest/gws/projects/` demonstrates the override is honored
- CLI: On Windows, `git-workspace init` writes `C:\Users\<user>\.config\gws\config.json` demonstrates the Windows layout matches the Unix one
- Test: `go test ./internal/config/` passes with the delegated path functions demonstrates no regression

### Unit 2: Config Relocation and Migration

**Purpose:** Move the config file to its XDG home and carry existing users across without them having to do anything.

**Functional Requirements:**

- `config.Load()` and `config.Save()` shall read and write `$XDG_CONFIG_HOME/gws/config.json`
- On load, if the XDG config does not exist but `~/.gws/config.json` does, gws shall migrate it: create the XDG config directory, move the file, and print a one-line notice to stderr naming both paths
- Migration shall preserve the existing `0600` file permission and `0755` directory permission
- If the legacy directory contains only the migrated `config.json`, it shall be removed after a successful migration; if it contains anything else, it shall be left in place and the extra contents mentioned in the notice
- If migration fails for any reason, gws shall fall back to reading the legacy path and warn, rather than presenting the user as uninitialized
- Migration shall be attempted at most once per process and shall be a no-op when the XDG config already exists
- The `not initialized` error shall continue to point at `gws init`, and shall not fire merely because a legacy config is present

**Proof Artifacts:**

- CLI: With a seeded `~/.gws/config.json`, running `git-workspace list` migrates the file and prints the notice demonstrates automatic migration
- CLI: `ls ~/.gws` after migration shows the directory removed demonstrates cleanup
- Test: `config_test.go` cases for migrate-on-load, no-op when XDG config exists, permission preservation, non-empty legacy dir retained, and fallback on failure demonstrates the migration contract
- CLI: A second `git-workspace list` prints no notice demonstrates idempotency

### Unit 3: Worktree Root Relocation

**Purpose:** Point worktree creation and alignment at the XDG projects root instead of `<repo>.wt/`, and redefine `IsAligned` accordingly.

**Functional Requirements:**

- `git.IsAligned` shall be redefined to test whether a worktree path is inside `ProjectsDir()/<repo-name>/`, and shall take the repo name (not the repo path) as its second input
- `gws worktree add <repo> <branch>` shall create the worktree at `ProjectsDir()/<repo-name>/<branch>`, creating parent directories as needed
- Branch names containing `/` shall continue to produce nested directories, as they do today
- `gws worktree align` shall target `ProjectsDir()/<repo-name>/` as its destination
- The existing `-dup-NN` conflict-resolution logic in `worktree_align.go` shall be reused as-is; only the base directory it joins against changes
- `gws worktree list` and the `(wt)` indicator in `gws list` shall continue to behave as they do today, reporting unaligned worktrees with the existing indicator
- Locked worktrees shall continue to be skipped by `align`, as they are today
- Help text and examples across `worktree.go`, `worktree_add.go`, `worktree_align.go`, and `worktree_list.go` shall be updated to describe the projects root rather than `<repo>.wt/`

**Proof Artifacts:**

- CLI: `gws worktree add my-repo feat-x` creates `~/.local/share/gws/projects/my-repo/feat-x` demonstrates the new destination
- CLI: `gws worktree add my-repo feature/nested` creates a nested directory demonstrates slash handling is preserved
- CLI: `gws worktree list` shows a pre-existing `.wt/` worktree marked `(unaligned)` demonstrates the redefined alignment
- Test: `TestIsAligned` updated for the projects root demonstrates the predicate is correct
- Test: `worktree_add_test.go` and `worktree_align_test.go` pass against the new root demonstrates no regressions

### Unit 4: Legacy `.wt` Migration via `align`

**Purpose:** Give users a single, explicit gesture that cleans up their home directory — moving every `.wt/` worktree into the projects root — without ever moving a working tree behind their back.

**Functional Requirements:**

- Because `IsAligned` now tests the projects root, existing `.wt/` worktrees are automatically picked up by `align` with no special-casing
- `gws worktree align --dry-run` shall continue to preview moves, now showing `.wt/` sources and projects-root destinations
- After a successful move, if the `<repo>.wt/` directory is left empty, it shall be removed so the home directory is actually cleaned up
- `align` shall detect when source and destination are on different filesystems — likely, since `~/.local/share` may be a separate mount — and shall report a clear, actionable error naming both paths rather than surfacing a raw `git worktree move` failure
- The config's `Worktree.Path` and `Worktree.Aligned` fields shall be updated after each move, as they are today
- `align` output shall make it visible that worktrees are being relocated to a new standard location the first time it runs after upgrade

**Proof Artifacts:**

- CLI: `gws worktree align --dry-run` on a repo with `.wt/` worktrees lists the planned moves demonstrates preview
- CLI: `gws worktree align` moves them and the emptied `<repo>.wt/` directory is gone demonstrates cleanup
- CLI: `gws worktree list` afterward shows every worktree without the `(unaligned)` marker demonstrates alignment
- Test: `worktree_align_test.go` case migrating a simulated `.wt/` layout into the projects root demonstrates the path
- Test: Cross-device move produces the actionable error demonstrates the failure mode is handled

### Unit 5: Documentation

**Purpose:** Document the new layout, the migration, and what users should expect on upgrade.

**Functional Requirements:**

- `docs/site/configuration.md` shall document the XDG config path, the `XDG_CONFIG_HOME` and `XDG_DATA_HOME` overrides, the Windows equivalents, and the automatic config migration
- `docs/site/commands-core.md` shall be updated wherever it references `<repo>.wt/`
- A short "Upgrading" note shall explain that the config moves automatically, that existing worktrees will report as unaligned until `gws worktree align` is run, and that this is expected
- `README.md` shall reflect the new paths in its layout/description sections
- The term **project** shall be defined once, plainly: a git repo plus any worktrees associated with it

**Proof Artifacts:**

- Documentation: `grep -rn '~/.gws' docs/ README.md` returns only historical/upgrade references demonstrates docs are current
- Documentation: configuration.md documents both XDG variables and Windows fallbacks demonstrates coverage

## Non-Goals (Out of Scope)

1. **Reworking repo tracking**: Repos stay tracked as they are today, in a single `config.json` with a `repositories` array bound to a workspace root. No per-repo record files, no named multi-repo workspaces.
2. **Changing discovery**: `gws refresh`, `gws init`, and `gws add` keep their current passive behavior. Aligned/unaligned reporting after discovery is explicitly preserved.
3. **Renaming `gws worktree` to `gws project`**: The directory root is named `projects/`, but the command surface stays `gws worktree`. A command rename is a separate, user-facing decision.
4. **Moving repos themselves**: Only worktrees relocate. Main clones stay exactly where the user put them.
5. **Automatic worktree migration**: Worktrees move only via explicit `gws worktree align`.
6. **A general `gws migrate` command**: Config migration is automatic and worktree migration is `align`; no third mechanism.
7. **Cache/state directories**: `XDG_CACHE_HOME` and `XDG_STATE_HOME` are not used by this spec. `internal/git/cache.go` is untouched.

## Design Considerations

The central simplification is that **`IsAligned` is the only real behavior change**. Every
downstream feature — the `(wt)` indicator, `worktree list`'s `(unaligned)` marker, `align`'s
work list, `refresh`'s tracking — already routes through that one predicate. Redefining it
against the projects root makes existing `.wt/` worktrees fall out as unaligned automatically,
with no special-case code and no third state to thread through the display logic. That is
precisely why option (A) was chosen in round 1: it turns a migration into a no-op at the call
sites.

The asymmetry between config and worktree migration is deliberate. The config file is small,
fully owned by gws, and reconstructible from a rescan, so moving it silently is safe and
saves every user a chore. A worktree may contain hours of uncommitted work; relocating one
without being asked is the kind of thing a tool gets uninstalled for. Hence: config moves
itself, worktrees wait for `align`.

`IsAligned`'s signature changes from `(worktreePath, repoPath)` to a repo-name-keyed form,
because the projects root is no longer derived from where the repo happens to sit on disk.
This is a small but load-bearing change — every call site in `refresh.go`, `worktree_add.go`,
and `worktree_align.go` must be updated together.

### Windows: mimic XDG rather than diverge

Cross-platform parity is a core requirement, so the layout is identical everywhere:
`<home>/.config/gws` and `<home>/.local/share/gws/projects`, with `<home>` being
`%USERPROFILE%` on Windows. Windows has no XDG specification, but it does not need one here —
what matters is that a user moving between machines finds gws in the same relative place.

The decisive precedent is **git itself**. Per `git-config(1)`, git reads
`$XDG_CONFIG_HOME/git/config`, and "when the XDG_CONFIG_HOME environment variable is not set
or empty, `$HOME/.config/` is used as `$XDG_CONFIG_HOME`." Git for Windows sets `$HOME` to
`%USERPROFILE%`, so `C:\Users\<user>\.config\git\config` is already a real, supported,
widely-populated path on Windows machines. A git-adjacent tool that puts its config beside
git's own is following the convention of the ecosystem it lives in, not inventing one.

The cost is that this diverges from Windows-native convention, where application config
belongs in `%AppData%`. That is a deliberate trade: `%AppData%` would give Windows users a
different layout from their colleagues, different documentation, and a path that no XDG-aware
tooling or dotfile manager knows about. The `XDG_CONFIG_HOME` and `XDG_DATA_HOME` variables
are honored on Windows too, so anyone who wants the native location can still point there
explicitly.

A useful side effect: `C:\Users\<user>\.local\share\gws\projects\<repo>\<branch>` is
shorter than the `%LocalAppData%` equivalent, which eases the `MAX_PATH` pressure noted below.

## Repository Standards

- Go source files follow standard `gofmt` formatting
- Tests use the standard `testing` package with table-driven tests where applicable
- Internal packages live under `internal/<name>/` with a matching `_test.go`
- Path handling uses `path/filepath` throughout for cross-platform correctness
- Config is JSON with `omitempty` on optional fields; `ConfigVersion` is bumped on format changes
- Cobra is used for CLI structure
- Conventional commits (`feat(config): ...`, `feat(worktree): ...`, `docs: ...`)

## Technical Considerations

- **Do not use `os.UserConfigDir()`.** It returns `%AppData%` on Windows, which would break the cross-platform uniformity this spec requires. Both `ConfigDir()` and `ProjectsDir()` must be written by hand: honor the `XDG_*` variable if set and absolute, else join `.config` / `.local/share` onto `os.UserHomeDir()`. `os.UserHomeDir()` correctly returns `%USERPROFILE%` on Windows.
- **Spec 21 just added Windows support**, so none of this may regress it. The PowerShell template and install docs assume a working config path; both should be exercised on Windows before this ships.
- **Cross-filesystem moves.** `git worktree move` is ultimately a rename and will fail with `EXDEV` if `~/.local/share` is on a different mount than the repo — plausible on setups with a separate `/home` or a small root partition. `MoveWorktree` in `internal/git/worktree.go` already has partial-move recovery; it needs an explicit cross-device branch that fails clearly rather than half-moving.
- **Windows path length.** `C:\Users\<user>\.local\share\gws\projects\<repo>\<branch>` plus a deep branch name can approach the 260-character `MAX_PATH` limit on systems without long-path support. The `%USERPROFILE%`-rooted layout is shorter than the `%LocalAppData%` alternative, but still longer than the old sibling-directory layout. Worth a documented note, and a candidate for a clear error rather than an obscure git failure.
- **`ConfigVersion` bump.** The on-disk schema does not change shape, but `Worktree.Path` values become stale relative to the new alignment rules. Bump `ConfigVersion` from `1.1.0` so the change is legible in the file itself.
- **Repo-name keying.** The projects root is keyed by `repo.Name`, which is a directory basename and is not guaranteed unique across a workspace. Two tracked repos both named `api` would share `projects/api/`. Per round 1 this is not being designed for now; see Open Questions.
- **Symlinked home directories.** `ListWorktrees` already calls `filepath.EvalSymlinks`; the alignment check against the projects root needs the same treatment, since `~/.local/share` is a symlink on some setups.
- **Test isolation.** Tests must never touch the real `~/.config` or `~/.local/share`. Every test that exercises path resolution should set `HOME`/`XDG_*` to a `t.TempDir()`.

## Security Considerations

- The config file must keep its `0600` permission across migration; it records repository paths and git identity metadata. A migration that widens permissions would leak that to other local users.
- The config directory should be created `0755` and the projects directory `0755`, consistent with existing behavior and with what git itself creates.
- `XDG_CONFIG_HOME` and `XDG_DATA_HOME` are attacker-influenceable in a shared environment. Per the XDG specification, relative values must be rejected rather than resolved against the current working directory, which would otherwise let a hostile `cd` redirect where gws writes.
- Migration must not follow a symlink at the legacy config path in a way that writes outside the intended directory; stat and verify before moving.
- No network access, credentials, or elevated permissions are involved.

## Success Metrics

1. **Home directory is quieter**: After `gws worktree align`, no `*.wt` directories remain in the user's project tree.
2. **XDG compliance**: Config resolves under `$XDG_CONFIG_HOME`, worktrees under `$XDG_DATA_HOME`, with both overrides honored.
3. **Cross-platform parity**: The relative layout is identical on Linux, macOS, and Windows; a user moving between machines finds gws in the same place, and the documentation has one set of paths rather than two.
4. **Silent upgrade**: An existing user upgrades, runs any command, and their config is migrated with a single notice and no lost data.
5. **No behavior regressions**: Discovery, tracking, and aligned/unaligned reporting behave exactly as before; the full test suite passes.
6. **No surprise moves**: No working tree is relocated except by an explicit `gws worktree align`.

## Open Questions

1. **Repo-name collisions at the projects root.** Two tracked repos sharing a basename (`~/work/api` and `~/oss/api`) would share `projects/api/`. Round 1 deferred this. Should `align` detect and refuse the collision, or is the existing `-dup-NN` scheme sufficient once branch names differ?
2. Should `gws cd` gain `--config` and `--projects` flags to reach these new directories? (Carried over from spec 22.)
3. Should the command surface eventually rename from `gws worktree` to `gws project` to match the directory vocabulary, or does `worktree` stay because it matches git's own term?
4. Is automatic config migration the right call, or should it be opt-in via a prompt on first run?
5. Should `align` offer to remove an emptied legacy `.wt` directory, or just do it silently as specified here?
