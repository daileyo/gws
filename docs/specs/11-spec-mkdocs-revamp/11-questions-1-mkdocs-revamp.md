# 11 Questions Round 1 - MkDocs Revamp

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Documentation Content Scope

The current docs are significantly behind the CLI. The following sections need updating. Which of these should be in scope for this spec?

- [X] (A) Full rewrite of all 4 content pages (getting-started, commands, shell-integration, configuration) to reflect subcommand syntax and all new features (user, tag, add, init, refresh, print-workspace)
- [ ] (B) Only update commands.md to add missing commands/subcommands, leave getting-started and configuration mostly as-is
- [ ] (C) Update commands.md and getting-started.md only (the most user-facing pages), defer configuration.md updates
- [ ] (D) Other (describe)

## 2. Deprecated Flags in Documentation

The old flag syntax (e.g., `--list`, `--init`, `--add-tag`) is deprecated but still works. How should the docs handle these?

- [ ] (A) Remove all references to deprecated flags — only document the new subcommand syntax
- [X] (B) Document the new subcommand syntax as primary, but include a "Legacy Flags" or "Migration" section noting the old flags still work
- [ ] (C) Keep both syntaxes side by side throughout (not recommended — creates confusion)
- [ ] (D) Other (describe)

## 3. MkDocs Nav Bar Icon Size

You mentioned wanting the nav bar icon to be bigger. The current logo is set via `theme.logo` in mkdocs.yml (rendered at Material theme's default size, typically ~24px). How much bigger are you thinking?

- [ ] (A) Slightly larger — roughly 32-36px (subtle increase, still compact in nav bar)
- [X] (B) Noticeably larger — roughly 40-48px (stands out more but still fits nav)
- [ ] (C) I'll provide a specific pixel size or ratio
- [ ] (D) Other (describe)

## 4. README Icon Layout

You want the README icon to be larger and centered above the title instead of the current side-by-side table layout. What size and style are you thinking?

- [X] (A) Centered icon (~300px wide) with the project title as a centered `<h1>` below it, followed by centered tagline
- [ ] (B) Centered icon (~200px wide, same size as now) with centered title and tagline below
- [ ] (C) Centered icon with title embedded in or overlapping the icon image (would need a new image asset)
- [ ] (D) Other (describe)

## 5. GitHub Repo Link in MkDocs

Where should the GitHub repo link appear in the MkDocs site?

- [X] (A) Use Material theme's built-in `repo_url` feature — adds a GitHub icon/link in the top-right header bar automatically
- [ ] (B) Add a manual link in the nav sidebar (e.g., "GitHub" nav entry)
- [ ] (C) Both — repo_url header icon plus a nav entry
- [ ] (D) Other (describe)

## 6. Navigation Structure

The current nav has 5 pages (Home, Getting Started, Commands, Shell Integration, Configuration). With the addition of `user` commands (6 subcommands), the commands page could get long. Should the nav structure change?

- [ ] (A) Keep the current 5-page structure — just expand the Commands page with new sections
- [X] (B) Split Commands into multiple pages (e.g., "Core Commands", "User Management", "Tagging")
- [ ] (C) Add a dedicated "User Management" page, keep everything else on one Commands page
- [ ] (D) Other (describe)

## 7. Home Page (index.md) Updates

The current homepage has quick links to the 4 doc pages and a brief feature overview. Should it be updated as part of this spec?

- [X] (A) Yes — update feature list and links to reflect current capabilities
- [ ] (B) Keep it mostly as-is, just fix any broken or outdated references
- [ ] (C) No changes needed to the homepage
- [ ] (D) Other (describe)
