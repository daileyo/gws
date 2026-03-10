# 18-spec-mkdocs-sync

## Introduction/Overview

The MkDocs documentation site and README.md for `gws` (git-workspace) have fallen significantly out of sync with the current CLI functionality (v2.17.0). Multiple commands have new flags, changed behavior, and entirely new features that are undocumented or incorrectly documented. This spec covers a comprehensive update of all documentation to accurately reflect the current state of the project.

## Goals

- Ensure every CLI command, flag, and behavior documented in the MkDocs site matches the v2.17.0 codebase
- Update all example output blocks to reflect actual CLI output
- Document the `parent` command and its navigation usage
- Document the complete dual-purpose (lowercase filter / uppercase show) flag system for `gws list`
- Update the README.md to stay in sync with the docs site
- Add all missing deprecated flags to the legacy flags page

## User Stories

- **As a new user**, I want the getting-started page to show accurate default output so that I know what to expect when I first run `gws`.
- **As a user learning the list command**, I want the docs to accurately describe the lowercase/uppercase flag convention so that I can filter and display columns correctly.
- **As a user navigating repos**, I want the `parent` command documented so that I can discover and use parent-directory navigation.
- **As a user migrating from old flags**, I want a complete legacy flags table so that I know the modern equivalent of every deprecated flag.
- **As a contributor reading the README**, I want it to accurately reflect the project's current feature set so that I understand what the tool does.

## Demoable Units of Work

### Unit 1: Core Commands Page Update (commands-core.md)

**Purpose:** Fix the most impactful documentation gap — the `list` command flags are entirely wrong, the `add` command has a wrong short flag, and `parent` is missing.

**Functional Requirements:**
- The `list` command section shall document the complete dual-purpose flag convention: lowercase flags (filter-only) and uppercase flags (show column + optional filter)
- The `list` command section shall document all filter flags: `--type`/`-y`, `--visibility`/`-i`, `--tag`/`-t`, `--path`/`-p`, `--status`/`-s`, `--show-user`/`-u`, `--remote`/`-r`, `--remote-raw`/`-b`, `--name`/`-n`
- The `list` command section shall document all show flags: `--show-type`/`-Y`, `--show-visibility`/`-I`, `--show-tag`/`-T`, `--show-path`/`-P`, `--show-status`/`-S`, `--show-user-col`/`-U`, `--show-remote`/`-R`, `--show-remote-raw`/`-B`
- The `list` command section shall document display flags: `--output`/`-o`, `--verbose`/`-v` (with `-v` and `-vv` levels), `--workers`, `--color`
- The `list` command section shall document that default output is a compact multi-column names-only layout (not a full table)
- The `list` command section shall document that JSON output only includes fields for shown columns
- The `list` command section shall show correct status icon order: behind → ahead → clean/dirty
- The `list` command section shall document branch name truncation at 30 characters
- The `list` command section shall document color output for status icons
- The `list` command section shall note that `--tag`/`-t` accepts a single value (not repeatable)
- The `list` command section shall note that `--type`, `--tag`, and `--visibility` use exact matching
- The `add` command section shall show the correct short flag `-r` for `--recursive` (not `-v`)
- The `add` command section shall document that `--recursive` scans from the current working directory
- The `add` command section shall document symlink creation for repos outside the workspace root
- The `parent` command shall be documented as a core command for navigating to a repo's parent directory
- All example output blocks shall be regenerated from actual v2.17.0 CLI output

**Proof Artifacts:**
- CLI: `gws list --help` output matches documented flags
- CLI: `gws add --help` output matches documented flags
- CLI: `gws parent --help` output matches documented description
- Comparison: each flag table entry verified against `--help` output

### Unit 2: User, Tagging, and Legacy Pages Update

**Purpose:** Fix missing flags and incorrect details across the secondary command documentation pages.

**Functional Requirements:**
- The `commands-user.md` shall document short-flag aliases on the `user` parent command: `-l` (list), `-a` (add), `-s` (show), `-d` (remove)
- The `commands-user.md` shall show correct `user list` output with separate "Stored Profiles" and "Auto-Detected Profiles" sections
- The `commands-user.md` shall use correct column header `SIGN` (not `SIGNING`)
- The `commands-tagging.md` shall document `--path`/`-p` and `--repo`/`-r` targeting flags on `tag add` and `tag remove`
- The `commands-tagging.md` shall document short-flag aliases on the `tag` parent command: `-a` (add), `-d` (delete), `-p` (path), `-r` (repo)
- The `commands-legacy.md` shall include all missing deprecated flags: `--user`, `--update`/`-u`, `--delete`/`-D`, `--all`, `--verbose` (user), `--git-name`, `--git-email`, `--list-users`, `-V` (visibility)
- All example output blocks shall be regenerated from actual v2.17.0 CLI output

**Proof Artifacts:**
- CLI: `gws user --help` output matches documented flags
- CLI: `gws tag add --help` output matches documented flags
- Comparison: legacy flags table verified against `deprecated.go` source

### Unit 3: Getting Started, Shell Integration, Configuration, and Index Pages Update

**Purpose:** Fix foundational pages that new users encounter first — wrong default output, outdated shell snippets, missing config fields.

**Functional Requirements:**
- The `getting-started.md` shall show correct default `gws` output (compact multi-column names view, not workspace summary)
- The `getting-started.md` shall keep the Homebrew install section as-is (it is complete and working)
- The `shell-integration.md` manual zsh snippet shall include `parent` command routing
- The `shell-integration.md` shall document `parent` navigation forms: `gws my-repo -p`, `gws -p my-repo`, `gws parent my-repo`
- The `configuration.md` shall document the `preferences` object and `status_workers` field
- The `configuration.md` JSON example shall show config version `1.1.0` (not `1.0.0`)
- The `index.md` feature list shall include: `parent` navigation, remote URL display, configurable worker pool, color status output, compact multi-column default view, symlink management for external repos
- The `refresh` command documentation (in commands-core.md or getting-started.md) shall reflect actual output: user config detection, stale repo removal, updated repo count, user configuration count
- The drift indicator (`⚠` in NAME column) shall be documented
- All example output blocks shall be regenerated from actual v2.17.0 CLI output

**Proof Artifacts:**
- CLI: `gws` (no args) output matches getting-started.md example
- CLI: `gws shell-init` output matches shell-integration.md manual snippet
- CLI: actual config.json structure matches configuration.md documentation
- Build: `python -m mkdocs build` completes with no warnings

### Unit 4: README.md Update

**Purpose:** Bring the repository README in sync with the docs site so that GitHub visitors see accurate feature information.

**Functional Requirements:**
- The README shall mention subcommand-based CLI structure (`list`, `init`, `add`, `refresh`, `tag`, `user`, `parent`)
- The README shall mention user profile management (`gws user`)
- The README shall mention tagging system (`gws tag`)
- The README shall mention parent directory navigation (`gws parent`)
- The README shall mention remote URL display capability
- The README shall mention color status output
- The README shall accurately reflect all current features without contradicting the docs site

**Proof Artifacts:**
- Comparison: README feature list cross-referenced against docs site index.md feature list
- Visual: README renders correctly on GitHub

## Non-Goals (Out of Scope)

1. **No code changes**: This spec covers documentation only; no CLI behavior, flags, or output will be modified
2. **No new docs pages**: All updates go into existing pages; no new navigation entries added to mkdocs.yml
3. **No docs site redesign**: The MkDocs Material theme, layout, and navigation structure remain unchanged
4. **No tutorial or cookbook content**: This is a reference documentation sync, not new educational content
5. **No automated doc generation**: We are manually updating markdown; no tooling to auto-generate docs from code

## Design Considerations

No specific design requirements identified. The existing MkDocs Material theme and page layout are retained. All updates are content-only within the current structure.

## Repository Standards

- **Commit convention**: `docs:` prefix for all commits (no release triggered)
- **MkDocs build**: Use `python -m mkdocs` (not bare `mkdocs` which is Homebrew's v2.0)
- **Python deps**: Located in `docs/requirements.txt` (not repo root)
- **Docs source**: All site content in `docs/site/`
- **Markdown style**: Follow existing patterns in docs pages (tables, code blocks, admonitions)

## Technical Considerations

- Example output should be generated from actual CLI runs against a representative workspace to ensure accuracy
- The `gws list` flags section is the largest single change — the entire dual-purpose flag paradigm needs clear explanation with examples
- Shell integration manual snippets must be verified against actual `gws shell-init` output
- Config JSON example should be verified against the `Config` struct in `config.go`

## Security Considerations

No specific security considerations identified. All changes are to public documentation.

## Success Metrics

1. **Completeness**: All 24 identified gaps are addressed across all docs pages
2. **Accuracy**: Every flag table, example output, and command description matches `--help` output and actual CLI behavior
3. **Build clean**: `python -m mkdocs build` completes with zero warnings
4. **README parity**: README feature list is consistent with docs site index.md

## Open Questions

No open questions at this time.
