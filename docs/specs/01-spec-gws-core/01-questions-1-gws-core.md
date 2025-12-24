# 01 Questions Round 1 - GWS Core

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Repository Management - What does "manage" mean?

What specific management capabilities do you need for git repositories?

- [ ] (A) Add repositories to tracking (manually specify path or URL)
- [ ] (B) Automatically discover repositories in specified directories
- [ ] (C) Remove repositories from tracking
- [ ] (D) Organize repositories into groups/categories/tags
- [ ] (E) Update repository metadata (last pulled, branch info, status)
- [ ] (F) Clone repositories from remote URLs
- [ ] (G) Other (describe)

## 2. Repository Classification - How should repos be categorized?

How do you want to classify and distinguish between different types of repositories?

- [ ] (A) Automatic detection based on remote URL (github.com, gitlab.com, etc.)
- [ ] (B) Manual tagging by user (personal, work, client-name, etc.)
- [ ] (C) Visibility classification (public/private based on remote access)
- [ ] (D) Both automatic detection AND manual tagging
- [ ] (E) Other (describe)

## 3. Search and Filter - What search capabilities are needed?

What do you want to be able to search and filter by?

- [ ] (A) Repository name or path
- [ ] (B) Repository type (github, gitlab, ado, bitbucket)
- [ ] (C) Tags or categories
- [ ] (D) Repository status (has uncommitted changes, needs pull, ahead of remote)
- [ ] (E) Date-based filters (last commit, last modified)
- [ ] (F) Content search (search within repository files)
- [ ] (G) Branch information
- [ ] (H) Other (describe)

## 4. Analysis - What kind of analysis do you need?

What insights or analysis should gws provide about your repositories?

- [ ] (A) Repository status summary (clean, dirty, ahead, behind)
- [ ] (B) Activity metrics (commit frequency, last commit date)
- [ ] (C) Size and file statistics
- [ ] (D) Dependency analysis (what languages/frameworks)
- [ ] (E) Branch health (stale branches, unmerged work)
- [ ] (F) Just basic git status information
- [ ] (G) Other (describe)

## 5. Data Storage - How should repository information be stored?

Where and how should gws store the repository metadata?

- [ ] (A) Single JSON/YAML file in user's home directory
- [ ] (B) SQLite database for fast querying
- [ ] (C) File per repository in a dedicated config directory
- [ ] (D) Git repository itself (track gws data in git)
- [ ] (E) No preference - you decide based on best practices
- [ ] (F) Other (describe)

## 6. Primary Use Cases - How will you use gws?

What are your main use cases for this tool?

- [ ] (A) Quickly find a specific repository among many
- [ ] (B) Get an overview of all repositories and their status
- [ ] (C) Batch operations (pull all work repos, check status of all personal repos)
- [ ] (D) Track repositories across multiple machines
- [ ] (E) Identify repositories that need attention (uncommitted changes, outdated)
- [ ] (F) Other (describe)

## 7. Command Line Interface - What CLI style do you prefer?

How should the CLI commands be structured?

- [ ] (A) Git-style subcommands (`gws add`, `gws list`, `gws search`)
- [ ] (B) Flag-based (`gws --add`, `gws --list`, `gws --search`)
- [ ] (C) Interactive mode with prompts
- [ ] (D) Both CLI commands AND interactive mode
- [ ] (E) No preference - you decide
- [ ] (F) Other (describe)

## 8. Cross-Platform Requirements - What platforms are critical?

Which operating systems must gws support in the first version?

- [ ] (A) Linux (primary focus)
- [ ] (B) macOS (primary focus)
- [ ] (C) Windows (primary focus)
- [ ] (D) All three equally important
- [ ] (E) Start with one platform, expand later (specify which)
- [ ] (F) Other (describe)

## 9. Installation and Distribution - How should users install gws?

What installation methods are important?

- [ ] (A) Single binary download (no dependencies)
- [ ] (B) Package managers (brew, apt, cargo, npm, etc.)
- [ ] (C) Build from source
- [ ] (D) All of the above
- [ ] (E) Other (describe)

## 10. Technology Stack - Any preferences or constraints?

Do you have preferences for the implementation language or technology?

- [ ] (A) Rust (performance, single binary, cross-platform)
- [ ] (B) Go (simple deployment, good CLI tools)
- [ ] (C) Python (easy to extend, good ecosystem)
- [ ] (D) Node.js/TypeScript (JavaScript ecosystem)
- [ ] (E) Shell script (maximum compatibility)
- [ ] (F) No preference - choose based on requirements
- [ ] (G) Other (describe)

## 11. Proof Artifacts - How will you demo this working?

What would you want to show to prove this feature works?

- [ ] (A) Terminal recording (asciinema) showing CLI commands
- [ ] (B) Screenshots of command output
- [ ] (C) Example repository list with search/filter results
- [ ] (D) Before/after showing repository organization
- [ ] (E) Performance benchmarks (searching across many repos)
- [ ] (F) Other (describe)
