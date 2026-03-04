# 11-spec-mkdocs-revamp

## Introduction/Overview

The MkDocs documentation site and README for `gws` have fallen significantly behind the current CLI. Major features added in v2.2.0-v2.7.0 — including subcommand-based syntax (`list`, `init`, `add`, `refresh`, `print-workspace`), the `tag` subcommand with `add`/`remove` operations, the full `user` command family (6 subcommands), and the deprecation of all legacy flag forms — are either missing or documented with outdated syntax. This spec covers a full documentation rewrite to align with the current CLI, cosmetic improvements to the nav bar icon and README header layout, and adding a GitHub repo link to the MkDocs site.

## Goals

- Rewrite all MkDocs content pages to reflect the current subcommand-based CLI syntax and all features through v2.7.0
- Split the commands reference into multiple focused pages to improve navigation and readability
- Include a "Legacy Flags" migration reference so users of older syntax can find their way to the new commands
- Increase the MkDocs nav bar logo size to ~40-48px for better visibility
- Redesign the README header to center the logo (~300px) above a centered title and tagline
- Add a GitHub repository link to the MkDocs site header using Material theme's `repo_url` feature
- Update the homepage (index.md) feature list and navigation links to reflect current capabilities

## User Stories

- **As a new user**, I want the Getting Started guide to show the current subcommand syntax so that I can set up my workspace correctly on the first try.
- **As an existing user**, I want clear documentation for all commands — including `user`, `tag`, `add`, and `refresh` subcommands — so that I can use the full feature set.
- **As a user upgrading from an older version**, I want a legacy flags reference so that I can understand how my existing commands map to the new subcommand syntax.
- **As a visitor to the docs site**, I want to easily navigate to the GitHub repository so that I can view the source code, file issues, or check releases.
- **As a visitor to the GitHub repo**, I want a visually appealing README with a prominent centered logo so that the project looks professional and polished.

## Demoable Units of Work

### Unit 1: MkDocs Configuration and Cosmetic Changes

**Purpose:** Update the mkdocs.yml configuration and CSS to add the GitHub repo link and increase the nav bar logo size. These are quick, foundational changes that affect the entire site.

**Functional Requirements:**
- The mkdocs.yml shall include `repo_url` and `repo_name` settings pointing to `https://github.com/daileyo/gws` so that the Material theme renders a GitHub icon/link in the top-right header bar
- The custom.css shall override the Material theme's nav logo size to approximately 40-48px (from the default ~24px)
- The nav structure in mkdocs.yml shall be updated to reflect the new page split (see Unit 3 for the pages themselves)

**Proof Artifacts:**
- Screenshot: MkDocs site header showing larger logo and GitHub icon link demonstrates cosmetic and config changes are applied
- File: `mkdocs.yml` diff showing `repo_url`, `repo_name`, and updated `nav` structure demonstrates configuration changes

### Unit 2: README Header Redesign

**Purpose:** Redesign the README.md header to center the logo above the title instead of the current side-by-side table layout, making it more visually prominent.

**Functional Requirements:**
- The README.md shall replace the `<table>` header layout with a centered `<p align="center">` block containing the logo image at ~300px width
- The README.md shall display the project title as a centered `<h1>` element below the logo
- The README.md shall display the tagline as centered text below the title
- All existing badges, feature list, and other README content shall remain unchanged

**Proof Artifacts:**
- Screenshot: GitHub-rendered README showing centered logo (~300px), centered title, and centered tagline demonstrates the new layout
- File: `README.md` diff showing the header HTML change demonstrates the implementation

### Unit 3: Documentation Content Rewrite — Core and Navigation Pages

**Purpose:** Rewrite the homepage, Getting Started, Shell Integration, and Configuration pages to use current subcommand syntax and document all features. Split the commands reference into focused pages.

**Functional Requirements:**
- The `index.md` (homepage) shall update the feature list to include user management, tag subcommands, and subcommand-based syntax
- The `index.md` shall update the documentation links section to reflect the new page structure (split commands pages)
- The `getting-started.md` shall replace all deprecated flag examples (`--init`, `--list`) with current subcommand syntax (`gws init`, `gws list`)
- The `getting-started.md` shall reference the `gws` shell function as the primary invocation method after shell integration setup
- The `shell-integration.md` shall replace `--print-workspace` with `gws print-workspace`
- The `shell-integration.md` shall update the "Pass any flag through" examples to use subcommand syntax (`gws list`, `gws refresh` instead of `gws --list`, `gws --refresh`)
- The `shell-integration.md` manual zsh setup snippet shall reflect the current shell-init output (routing subcommand names like `list`, `init`, `add`, `refresh`, `print-workspace`, `tag`, `user`)
- The `configuration.md` shall document the `profiles` top-level config array (added in v2.2.0)
- The `configuration.md` shall document the repository-level fields `user`, `email`, `signing_enabled`, and `user_source`
- The `configuration.md` shall update the `tags` field reference to note `gws tag add`/`gws tag remove` as the management commands (replacing `--add-tag`)

**Proof Artifacts:**
- Screenshot: Updated homepage showing revised feature list and navigation links demonstrates index.md changes
- File: `getting-started.md` showing all subcommand syntax demonstrates the rewrite
- File: `configuration.md` showing profiles and new repository fields demonstrates config docs update
- File: `shell-integration.md` showing updated examples demonstrates shell docs update

### Unit 4: Documentation Content Rewrite — Command Reference Pages

**Purpose:** Split the monolithic commands.md into focused command reference pages covering all current subcommands, and include a legacy flags migration section.

**Functional Requirements:**
- The commands reference shall be split into separate pages organized by command group (e.g., core commands, user management, tagging) rather than one monolithic page
- A core commands page shall document `gws list` (with all flags: `--type`, `--tag`, `--name`, `--path`, `--output`, `--status`, `--show-user`), `gws init`, `gws add` (with `--recursive`), `gws refresh`, and `gws print-workspace` using the current subcommand syntax
- A user management page shall document the full `gws user` command family: `list`, `add` (with `--email`, `--name`, `--signing-key`, `--sign-commits`), `show`, `remove`, `assign` (with `--use-subdirs`, `--dry-run`), and `sync`
- A tagging page shall document `gws tag add <repo> <tag>` and `gws tag remove <repo> <tag>`, including partial name matching behavior and tab completion
- The repository classification section (type detection, visibility inference) shall be preserved in the appropriate commands page
- A "Legacy Flags" or "Migration" section shall be included documenting the deprecated flag-to-subcommand mappings (e.g., `--list` → `gws list`, `--init` → `gws init`, `--add-tag` → `gws tag add`, etc.) and noting these flags still work but emit deprecation warnings
- All command examples shall use the `gws` shell function form as primary (not `git-workspace` binary form)

**Proof Artifacts:**
- File: New command reference pages showing full subcommand documentation demonstrates the content split and rewrite
- File: Legacy flags section showing deprecated-to-new mapping table demonstrates migration reference
- Screenshot: MkDocs nav showing the split command pages demonstrates the new navigation structure

## Non-Goals (Out of Scope)

1. **CLI code changes**: No changes to Go source code, commands, or flags — this is documentation only
2. **New image assets**: The existing logo PNGs (nav, favicon, hero) will be reused; no new images will be created
3. **MkDocs theme changes**: The Material theme, color palette (deep purple/orange), fonts (Nunito/JetBrains Mono), and dark/light toggle will remain unchanged
4. **Docs deployment pipeline**: The GitHub Pages deployment workflow is not being modified
5. **API or developer documentation**: This spec does not cover generating godoc or internal API docs

## Design Considerations

- **Nav bar logo**: Override Material theme default logo size via `custom.css` targeting the `.md-header__button.md-logo img` selector, setting height to approximately 40-48px
- **README header**: Use `<p align="center">` and `<h1 align="center">` HTML tags for GitHub-compatible centering (GitHub Markdown does not support CSS-based centering)
- **Command page split**: The exact page names and nav labels should be clear and concise (e.g., "Commands — Core", "Commands — User Management", "Commands — Tagging") to maintain a logical hierarchy in the sidebar

## Repository Standards

- Documentation uses MkDocs with Material theme; source files live in `docs/site/`
- Commit messages follow Conventional Commits format (`docs:` type for documentation changes)
- CSS customizations go in `docs/site/assets/css/custom.css`
- The `gws` shell function is the primary user-facing invocation; `git-workspace` is the binary name used when discussing direct binary execution

## Technical Considerations

- The Material theme's `repo_url` feature automatically detects GitHub URLs and renders the GitHub icon — only `repo_url` and `repo_name` need to be added to `mkdocs.yml`
- The nav bar logo size override requires targeting the correct Material theme CSS selector; the selector may vary slightly between Material theme versions
- The nav structure in `mkdocs.yml` supports nested sections which can be used for grouping command pages under a "Commands" parent heading
- README HTML must be compatible with GitHub's Markdown renderer (limited subset of HTML, no `<style>` tags, `align` attributes work)

## Security Considerations

No specific security considerations identified. This spec covers documentation changes only.

## Success Metrics

1. **Documentation accuracy**: All CLI commands documented in MkDocs match the actual v2.7.0 command structure (no deprecated flag syntax used as primary examples)
2. **Page completeness**: Every subcommand (`list`, `init`, `add`, `refresh`, `print-workspace`, `tag add/remove`, `user list/add/show/remove/assign/sync`) has documented syntax, flags, and at least one example
3. **Visual improvement**: Nav bar logo is visibly larger; README logo is centered and ~300px wide above a centered title
4. **GitHub link**: Clicking the GitHub icon in the MkDocs header navigates to the correct repository URL

## Open Questions

No open questions at this time.
