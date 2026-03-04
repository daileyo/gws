# 11-tasks-mkdocs-revamp

## Relevant Files

- `mkdocs.yml` - MkDocs site configuration: theme, nav structure, repo_url
- `docs/site/assets/css/custom.css` - Custom CSS overrides for Material theme (logo size)
- `README.md` - Project README with header layout (logo, title, tagline, badges)
- `docs/site/index.md` - MkDocs homepage: feature list and documentation links
- `docs/site/getting-started.md` - Installation and quick start guide
- `docs/site/shell-integration.md` - Shell setup, navigation, and tab completion docs
- `docs/site/configuration.md` - Config file structure and field reference
- `docs/site/commands.md` - Current monolithic commands page (to be replaced by split pages)
- `docs/site/commands-core.md` - New page: core commands (list, init, add, refresh, print-workspace)
- `docs/site/commands-user.md` - New page: user management commands (user list/add/show/remove/assign/sync)
- `docs/site/commands-tagging.md` - New page: tagging commands (tag add/remove)
- `docs/site/commands-legacy.md` - New page: deprecated flag-to-subcommand migration reference

### Notes

- All documentation source files live in `docs/site/` per the mkdocs.yml `docs_dir` setting.
- The `gws` shell function is the primary invocation form in all examples; use `git-workspace` only when referring to the binary directly (e.g., installation, `eval "$(git-workspace shell-init zsh)"`).
- Commit messages should use `docs:` type per Conventional Commits.
- README HTML must use `align` attributes (not CSS) for GitHub Markdown compatibility.
- Refer to source files in `cmd/git-workspace/` for exact flag names, short flags, and default values when writing command documentation.

## Tasks

### [x] 1.0 MkDocs Configuration and Cosmetic Changes

Update `mkdocs.yml` to add the GitHub repo link (`repo_url`/`repo_name`), update the `nav` structure to accommodate the command page split, and add CSS in `custom.css` to increase the nav bar logo size to ~40-48px.

#### 1.0 Proof Artifact(s)

- Screenshot: MkDocs site header showing larger logo (~40-48px) and GitHub icon/link in the top-right corner demonstrates cosmetic and config changes are applied
- Diff: `mkdocs.yml` showing added `repo_url`, `repo_name`, and updated `nav` with split command pages demonstrates configuration changes
- Diff: `custom.css` showing logo size override demonstrates CSS customization

#### 1.0 Tasks

- [x] 1.1 Add `repo_url: https://github.com/daileyo/gws` and `repo_name: gws` to `mkdocs.yml` (top-level, after `docs_dir`)
- [x] 1.2 Update the `nav` section in `mkdocs.yml` to replace the single "Commands Reference" entry with a nested "Commands" group containing four sub-pages: "Core" (`commands-core.md`), "User Management" (`commands-user.md`), "Tagging" (`commands-tagging.md`), and "Legacy Flags" (`commands-legacy.md`)
- [x] 1.3 Add CSS rule in `custom.css` targeting `.md-header__button.md-logo img` to set `height: 44px` (up from the Material theme default of ~24px)
- [x] 1.4 Verify changes by running `python -m mkdocs serve` and confirming the larger logo and GitHub icon appear in the header bar

### [x] 2.0 README Header Redesign

Replace the current `<table>` header layout in `README.md` with a centered logo (~300px) above a centered `<h1>` title and centered tagline. All existing badges, features, and other content remain unchanged.

#### 2.0 Proof Artifact(s)

- Screenshot: GitHub-rendered README showing centered logo (~300px wide), centered `<h1>` title, and centered tagline demonstrates the new layout
- Diff: `README.md` header section showing `<p align="center">` and `<h1 align="center">` replacing the old `<table>` layout demonstrates the implementation

#### 2.0 Tasks

- [x] 2.1 Replace lines 1-4 of `README.md` (the `<table>` block with side-by-side logo and title) with: a `<p align="center">` containing the logo `<img>` at `width="300"`, followed by an `<h1 align="center">git-workspace</h1>`
- [x] 2.2 Replace line 6 (the uncentered tagline `*Your Git workspace, simplified*`) with `<p align="center"><em>Your Git workspace, simplified</em></p>`
- [x] 2.3 Verify the rendered layout on GitHub (or via GitHub Markdown preview) shows the logo centered above the title and tagline

### [x] 3.0 Documentation Content Rewrite — Core and Navigation Pages

Rewrite `index.md`, `getting-started.md`, `shell-integration.md`, and `configuration.md` to use current subcommand syntax (`gws init`, `gws list`, etc.), add missing features (user management, profiles config fields), and update navigation links to match the new page structure.

#### 3.0 Proof Artifact(s)

- File: `index.md` showing updated feature list and documentation links reflecting the split command pages demonstrates homepage update
- File: `getting-started.md` showing `gws init`, `gws list` syntax (no deprecated `--init`/`--list` flags) demonstrates subcommand rewrite
- File: `shell-integration.md` showing `gws print-workspace`, `gws list`, `gws refresh` examples and updated manual zsh snippet demonstrates shell docs update
- File: `configuration.md` showing `profiles` array, repository-level `user`/`email`/`signing_enabled`/`user_source` fields, and `gws tag add`/`gws tag remove` references demonstrates config docs update

#### 3.0 Tasks

- [x] 3.1 Update `index.md` feature list to include: user profile management, tag subcommands, subcommand-based CLI, and `--show-user` flag on list
- [x] 3.2 Update `index.md` documentation links section to reflect the new split command pages (Core Commands, User Management, Tagging, Legacy Flags) instead of the single "Commands Reference" link
- [x] 3.3 Rewrite `getting-started.md` Quick Start section: replace `git-workspace --init .` with `gws init` (and `gws init ~/projects`), replace `git-workspace --list` with `gws list`, and note that `gws` requires shell integration setup (link to shell-integration.md)
- [x] 3.4 Update `getting-started.md` "View Workspace Information" section to show `gws` (no args) output and update the help reference from `--help` to `gws --help` or `gws help`
- [x] 3.5 Update `shell-integration.md` "Repository Navigation" examples: replace `gws --list` with `gws list`, `gws -l --tag personal --status` with `gws list --tag personal --status`, and `gws --refresh` with `gws refresh`
- [x] 3.6 Update `shell-integration.md` "Workspace Navigation" section: replace `git-workspace --print-workspace` with `gws print-workspace`
- [x] 3.7 Update the manual zsh setup snippet in `shell-integration.md` to match the current `shell-init` output: the `case` statement should route `list|init|add|refresh|print-workspace|tag|user|completion|shell-init|help` to the binary
- [x] 3.8 Update `configuration.md` to add a "Profiles" section documenting the top-level `profiles` array with fields: `name`, `git_name`, `email`, `signing_key` (optional), `sign_commits` (optional)
- [x] 3.9 Update `configuration.md` Repository Fields table to add: `user` (string, git user.name), `email` (string, git user.email), `signing_enabled` (boolean), `user_source` (string: `global`, `local`, `includeif`, `unknown`)
- [x] 3.10 Update `configuration.md` to change the `tags` field description to reference `gws tag add`/`gws tag remove` instead of `--add-tag`
- [x] 3.11 Update `configuration.md` example JSON to include a sample `profiles` entry and the new repository-level fields

### [x] 4.0 Documentation Content Rewrite — Command Reference Pages

Split the monolithic `commands.md` into focused pages: a core commands page (`list`, `init`, `add`, `refresh`, `print-workspace` with all flags), a user management page (full `gws user` family with 6 subcommands), a tagging page (`gws tag add`/`remove`), and a legacy flags migration reference. Remove or replace the old `commands.md`.

#### 4.0 Proof Artifact(s)

- File: Core commands page documenting `gws list` (with `--type`, `--tag`, `--name`, `--path`, `--output`, `--status`, `--show-user`), `gws init`, `gws add --recursive`, `gws refresh`, `gws print-workspace` demonstrates core command documentation
- File: User management page documenting `gws user list/add/show/remove/assign/sync` with all flags demonstrates user command documentation
- File: Tagging page documenting `gws tag add <repo> <tag>` and `gws tag remove <repo> <tag>` with partial matching and tab completion demonstrates tag command documentation
- File: Legacy flags section showing deprecated-to-new mapping table (e.g., `--list` → `gws list`, `--add-tag` → `gws tag add`) demonstrates migration reference
- Screenshot: MkDocs nav sidebar showing split command pages under a "Commands" group demonstrates the new navigation structure

#### 4.0 Tasks

- [x] 4.1 Create `docs/site/commands-core.md` with a "Core Commands" heading and sections for: Repository Classification (preserve the type detection table and URL patterns from the current `commands.md`), `gws list` (document all flags: `--type`/`-y`, `--tag`/`-t`, `--name`/`-n`, `--path`/`-p`, `--output`/`-o`, `--status`/`-s`, `--show-user`/`-u` with examples and output samples), `gws init [directory]`, `gws add [path]` (with `--recursive`/`-v`), `gws refresh`, and `gws print-workspace`
- [x] 4.2 Create `docs/site/commands-user.md` with a "User Management" heading and sections for each subcommand: `gws user list`, `gws user add <name>` (flags: `--email` required, `--name`, `--signing-key`, `--sign-commits`), `gws user show <name>`, `gws user remove <name>`, `gws user assign <repository> <profile>` (flags: `--use-subdirs`, `--dry-run`), and `gws user sync` — each with usage syntax, flag table, and example output
- [x] 4.3 Create `docs/site/commands-tagging.md` with a "Tagging" heading and sections for: `gws tag add <repo> <tag>` and `gws tag remove <repo> <tag>`, including how partial name matching works, that tags apply to all matching repos, and that tab completion is available for repo names and existing tags
- [x] 4.4 Create `docs/site/commands-legacy.md` with a "Legacy Flags" heading, an introductory paragraph explaining these flags are deprecated but still functional, and a mapping table with columns: Deprecated Flag, Short, Replacement — covering all 16 deprecated flags from `depWarnings` in `deprecated.go` (`--list`→`gws list`, `--init`→`gws init`, `--add`→`gws add [path]`, `--recursive`→`gws add --recursive`, `--refresh`→`gws refresh`, `--print-workspace`→`gws print-workspace`, `--go`→`gws <repo-name>`, `--add-tag`→`gws tag add <repo> <tag>`, `--remove-tag`→`gws tag remove <repo> <tag>`, `--type`→`gws list --type`, `--name`→`gws list --name`, `--path`→`gws list --path`, `--output`→`gws list --output`, `--status`→`gws list --status`, `--show-user`→`gws list --show-user`)
- [x] 4.5 Delete the old `docs/site/commands.md` file (its content has been fully replaced by the four new pages)
- [x] 4.6 Verify the full site builds and navigates correctly by running `python -m mkdocs serve` and clicking through all nav links and command pages
