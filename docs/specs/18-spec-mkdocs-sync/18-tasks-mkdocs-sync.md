# 18 Tasks - MkDocs Sync

## Relevant Files

- `docs/site/commands-core.md` - Core commands documentation (list, init, add, refresh, print-workspace, parent)
- `docs/site/commands-user.md` - User management documentation (user list/add/show/remove/assign/sync)
- `docs/site/commands-tagging.md` - Tagging documentation (tag add/remove)
- `docs/site/commands-legacy.md` - Deprecated flags migration reference
- `docs/site/getting-started.md` - Installation, shell setup, and quick start guide
- `docs/site/shell-integration.md` - Shell function setup, navigation, and tab completion
- `docs/site/configuration.md` - Config file structure and field reference
- `docs/site/index.md` - Docs site homepage with feature list and navigation
- `README.md` - Repository README with feature list, installation, and development guide

### Notes

- All changes are documentation-only (markdown files). No Go code modifications.
- Use `docs:` commit prefix for all commits (does not trigger a release).
- Use `python -m mkdocs` (not bare `mkdocs`) for building/previewing the docs site.
- Example CLI output should be regenerated from actual v2.17.0 CLI runs where possible.
- Source of truth for flags: `cmd/git-workspace/*.go` flag registrations.
- Source of truth for deprecated flags: `cmd/git-workspace/deprecated.go` `depWarnings` map (23 entries).
- Source of truth for config structure: `internal/config/config.go` Config/Preferences/Profile/Repository structs.

## Tasks

### [x] 1.0 Update Core Commands Page (commands-core.md)

Rewrite the `list` command section with the complete dual-purpose flag convention (lowercase = filter-only, uppercase = show column + optional filter). Fix the `add` command's wrong `-v` short flag to `-r`. Add the `parent` command. Update `refresh` output. Regenerate all example output from actual v2.17.0 CLI.

#### 1.0 Proof Artifact(s)

- CLI: `gws list --help` output matches all documented filter, show, and display flags
- CLI: `gws add --help` confirms `--recursive`/`-r` short flag
- CLI: `gws parent --help` output matches documented description
- Comparison: each flag table entry in commands-core.md verified against `--help` output

#### 1.0 Tasks

- [x] 1.1 Rewrite the `list` Flags section: replace the existing 7-row flag table with three separate tables — **Filter Flags** (lowercase: `--type`/`-y`, `--visibility`/`-i`, `--tag`/`-t`, `--path`/`-p`, `--status`/`-s`, `--show-user`/`-u`, `--remote`/`-r`, `--remote-raw`/`-b`, `--name`/`-n`), **Show Column Flags** (uppercase: `--show-type`/`-Y`, `--show-visibility`/`-I`, `--show-tag`/`-T`, `--show-path`/`-P`, `--show-status`/`-S`, `--show-user-col`/`-U`, `--show-remote`/`-R`, `--show-remote-raw`/`-B`), and **Display Flags** (`--output`/`-o`, `--verbose`/`-v`, `--workers`, `--color`). Add an introductory paragraph explaining the lowercase/uppercase dual-purpose convention.
- [x] 1.2 Add notes to the `list` section: (a) `--tag`/`-t` accepts a single value, not repeatable; (b) `--type`, `--tag`, and `--visibility` use exact matching; (c) show-column flags accept an optional value to simultaneously show and filter; (d) `--verbose`/`-v` shows stored data columns, `-vv` shows all columns including status/user/remote.
- [x] 1.3 Update the "Basic listing" example: replace the full table output with the compact multi-column names-only layout that `gws list` now produces by default. Add a new example showing `gws list -v` for the stored-data table view.
- [x] 1.4 Update the "With git status" example: change status icon order from `main ✗ ↑2` / `main ✓ ↓1` to the correct order: behind → ahead → clean/dirty (e.g., `↓1 ✓`). Add a note about color output for status icons. Add a note about branch name truncation at 30 characters.
- [x] 1.5 Update the "Filtering" examples: remove the `--tag work --tag backend` repeated-tag example (no longer valid). Update `--status` examples to show it as a filter flag (not a show flag). Add examples for new filter flags (`--visibility`, `--remote`, `--remote-raw`).
- [x] 1.6 Update the "JSON output" example: show that bare `gws list -o json` only includes `name`. Show an example with show-flags (e.g., `gws list -o json -YITP`) that includes the additional fields. Remove the current JSON example that always shows all fields.
- [x] 1.7 Fix the `add` command Flags table: change `--recursive` short flag from `-v` to `-r`. Add a note that `--recursive` scans from the current working directory (not the path argument). Add a note about symlink creation for repositories located outside the workspace root.
- [x] 1.8 Add a new `## Parent Navigation` section after `Print Workspace`: document `gws parent <repo>` with description "Print the parent directory of a repository", the `--quiet`/`-q` flag, and usage examples (`gws parent my-repo`, `cd "$(gws parent my-repo)"`).
- [x] 1.9 Update the `refresh` example output to match actual output: include "Detecting git user configuration...", conditional lines for removed/updated/new repos, and "Repositories with user configuration: N". Add a note about the drift indicator (`⚠` in NAME column when user config drifts).

### [x] 2.0 Update User, Tagging, and Legacy Pages

Add missing short-flag aliases to `commands-user.md` (`-l`, `-a`, `-s`, `-d`). Fix `user list` output format (separate Stored/Auto-Detected sections, `SIGN` column header). Add `--path`/`-p` and `--repo`/`-r` targeting flags to `commands-tagging.md`. Add all missing deprecated flags to `commands-legacy.md` (~9 flags including user-related deprecations and `-V`).

#### 2.0 Proof Artifact(s)

- CLI: `gws user --help` output matches documented short-flag aliases
- CLI: `gws tag add --help` output matches documented `--path`/`--repo` flags
- Comparison: `commands-legacy.md` flag table verified against `deprecated.go` source entries

#### 2.0 Tasks

- [x] 2.1 Add a "Short-Flag Aliases" section to `commands-user.md` (after the intro, before List Profiles): document that the `user` parent command supports `-l` (list), `-a` (add), `-s` (show), `-d` (remove) as shorthand. Note that `-a` also accepts `--email`, `--name`, `--signing-key`, `--sign-commits` inline.
- [x] 2.2 Update the `user list` example output in `commands-user.md`: change column header from `SIGNING` to `SIGN`. Show separate "Stored Profiles" and "Auto-Detected Profiles" table sections with a footer note "(Auto-detected from ~/.gitconfig includeIf directives)".
- [x] 2.3 Add `--path`/`-p` and `--repo`/`-r` flags to both `tag add` and `tag remove` sections in `commands-tagging.md`. Add a flags table for each subcommand showing these targeting flags. Add examples: `gws tag add --path /work backend`, `gws tag add --repo api backend`.
- [x] 2.4 Add a "Short-Flag Aliases" section to `commands-tagging.md`: document that the `tag` parent command supports `-a` (add), `-d` (delete), `-p` (path), `-r` (repo) as shorthand. Add examples: `gws tag -a my-repo work`, `gws tag -d my-repo work`.
- [x] 2.5 Update the "Using Tags for Filtering" section in `commands-tagging.md`: remove the `--tag work --tag backend` repeated-tag example (no longer valid for AND logic via repetition). Update to reflect that `--tag` accepts a single value.
- [x] 2.6 Add all missing deprecated flags to `commands-legacy.md` Migration Reference table. Add these 9 entries: `--user` → `gws user list`, `--update`/`-u` → `gws user assign <repo> <profile>`, `--delete`/`-D` → `gws user assign (remove local config)`, `--all` → `gws user assign (with --all)`, `--verbose` → `gws user --verbose`, `--git-name` → `gws user add --name`, `--git-email` → `gws user add --email`, `--list-users` → `gws user list`, `-V` → `gws list -i (filter) or -I (show)`.

### [x] 3.0 Update Getting Started, Shell Integration, Configuration, and Index Pages

Fix `getting-started.md` default `gws` output (compact multi-column names, not workspace summary). Update `shell-integration.md` manual zsh snippet to include `parent` routing and document parent navigation forms. Add `preferences`/`status_workers` to `configuration.md` and update config version to `1.1.0`. Update `index.md` feature list with parent navigation, remote URL display, color output, compact default view, symlink management, and configurable workers.

#### 3.0 Proof Artifact(s)

- CLI: `gws` (no args) output matches getting-started.md example
- CLI: `gws shell-init zsh` output matches shell-integration.md manual snippet routing
- Comparison: `configuration.md` JSON example and field table verified against `config.go` Config struct
- Build: `python -m mkdocs build` completes with no warnings

#### 3.0 Tasks

- [x] 3.1 Update `getting-started.md` section "2. View Workspace Information": replace the old workspace summary output (`Workspace: ...`, `Repositories: 15`, `Use 'gws help'...`) with the compact multi-column names view that `gws` now produces. Update the description text from "workspace summary" to "list your repositories".
- [x] 3.2 Update `getting-started.md` section "3. List Repositories": replace the full-table example output with the compact multi-column names output (matching default `gws list`). Add an example showing `gws list -v` for the stored-data table view.
- [x] 3.3 Replace the manual zsh snippet in `shell-integration.md` (lines 93-98) with the actual shell function from `shellinit.go`: include the `parent` routing arm (`-p|--parent|parent`), the `$2 == "-p"` check in the default arm, and the `-*` passthrough arm. Use `git-workspace` as the binary name (matching `shellinit.go`'s replacement).
- [x] 3.4 Add a "Parent Navigation" section to `shell-integration.md` (under Repository Navigation): document the three forms — `gws parent my-repo`, `gws -p my-repo`, `gws my-repo -p` — and explain that all navigate to the parent directory containing the repository.
- [x] 3.5 Update `configuration.md`: (a) change config version from `"1.0.0"` to `"1.1.0"` in both the JSON example and the Top-level Fields table; (b) add a `preferences` object to the JSON example with `"status_workers": 8`; (c) add a new "Preferences Fields" section to the field reference documenting `status_workers` (int, default 8, number of concurrent workers for status fetching).
- [x] 3.6 Add `preferences` to the top-level fields table in `configuration.md`: add row `preferences | object | Optional preferences for controlling CLI behavior (see below)`.
- [x] 3.7 Update `index.md` feature list: add missing features — parent directory navigation (`gws parent`), remote URL display (`--show-remote`), configurable concurrent status workers, color status output, compact multi-column default view, symlink management for external repos. Add `parent` to the subcommand list in the "Subcommand-Based CLI" bullet.

### [x] 4.0 Update README.md

Sync the README feature list with the docs site. Add missing features: subcommand-based CLI (add `parent` to the list), user profile management, tagging, parent directory navigation, remote URL display, color status output. Ensure no contradictions with docs site index.md.

#### 4.0 Proof Artifact(s)

- Comparison: README feature list cross-referenced against docs site `index.md` feature list — all features present in both
- Visual: README renders correctly on GitHub (verified via `gh` or local preview)

#### 4.0 Tasks

- [x] 4.1 Update the README Features section: add missing bullet points for subcommand-based CLI (add `parent` to the command list), user profile management (`gws user`), parent directory navigation (`gws parent`), remote URL display, and color status output. Ensure feature descriptions match the updated `index.md` feature list.
- [x] 4.2 Update the README Project Structure tree: verify it matches actual directory structure (check for any new directories or missing entries). Add `parent.go` if individual command files are listed (they aren't currently, so this may be a no-op).
- [x] 4.3 Final cross-reference: compare every feature bullet in README against `index.md` feature list to ensure parity. Fix any wording inconsistencies between the two files.
