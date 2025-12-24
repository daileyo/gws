# 01 Questions Round 1 - GWS Core

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Repository Management - What does "manage" mean?

What specific management capabilities do you need for git repositories?

- [ ] (A) Add repositories to tracking (manually specify path or URL)
- [X] (B) Automatically discover repositories in specified directories
- [ ] (C) Remove repositories from tracking
- [X] (D) Organize repositories into groups/categories/tags
- [ ] (E) Update repository metadata (last pulled, branch info, status)
- [ ] (F) Clone repositories from remote URLs
- [X] (G) Other gws will not do any direct interactions with the repositories themselves at this time. This may be a future enhancement. It instead will initialize a location for repositoires to exist, navigate there when a command like gws is ran from the command line. For example: gws alone wil immidately drop into the location that all of the repose on the system are. gws personal (or something like this) will list all repose that have been identified as personal repositories, etc.

## 2. Repository Classification - How should repos be categorized?

How do you want to classify and distinguish between different types of repositories?

- [X] (A) Automatic detection based on remote URL (github.com, gitlab.com, etc.)
- [X] (B) Manual tagging by user (personal, work, client-name, etc.)
- [X] (C) Visibility classification (public/private based on remote access)
- [X] (D) Both automatic detection AND manual tagging
- [ ] (E) Other (describe)

## 3. Search and Filter - What search capabilities are needed?

What do you want to be able to search and filter by?

- [X] (A) Repository name or path
- [X] (B) Repository type (github, gitlab, ado, bitbucket)
- [X] (C) Tags or categories
- [X] (D) Repository status (has uncommitted changes, needs pull, ahead of remote)
- [X] (E) Date-based filters (last commit, last modified)
- [ ] (F) Content search (search within repository files)
- [X] (G) Branch information
- [ ] (H) Other (describe)

## 4. Analysis - What kind of analysis do you need?

What insights or analysis should gws provide about your repositories?

- [ ] (A) Repository status summary (clean, dirty, ahead, behind)
- [ ] (B) Activity metrics (commit frequency, last commit date)
- [ ] (C) Size and file statistics
- [ ] (D) Dependency analysis (what languages/frameworks)
- [ ] (E) Branch health (stale branches, unmerged work)
- [X] (F) Just basic git status information
- [ ] (G) Other (describe)

## 5. Data Storage - How should repository information be stored?

Where and how should gws store the repository metadata?

- [X] (A) Single JSON/YAML file in user's home directory
- [ ] (B) SQLite database for fast querying
- [ ] (C) File per repository in a dedicated config directory
- [ ] (D) Git repository itself (track gws data in git)
- [X] (E) No preference - you decide based on best practices
- [ ] (F) Other (describe)

## 6. Primary Use Cases - How will you use gws?

What are your main use cases for this tool?

- [X] (A) Quickly find a specific repository among many
- [ ] (B) Get an overview of all repositories and their status
- [ ] (C) Batch operations (pull all work repos, check status of all personal repos)
- [ ] (D) Track repositories across multiple machines
- [ ] (E) Identify repositories that need attention (uncommitted changes, outdated)
- [X] (F) Quickly filter, find, and analyze repositories based on the metadata 

## 7. Command Line Interface - What CLI style do you prefer?

How should the CLI commands be structured?

- [ ] (A) Git-style subcommands (`gws add`, `gws list`, `gws search`)
- [ ] (B) Flag-based (`gws --add`, `gws --list`, `gws --search`)
- [ ] (C) Interactive mode with prompts
- [X] (D) Both CLI commands AND interactive mode
- [ ] (E) No preference - you decide
- [X] (F) Other git-style subcommands with flag-based options (when needed) for the subcommands.

## 8. Cross-Platform Requirements - What platforms are critical?

Which operating systems must gws support in the first version?

- [ ] (A) Linux (primary focus)
- [ ] (B) macOS (primary focus)
- [ ] (C) Windows (primary focus)
- [X] (D) All three equally important
- [ ] (E) Start with one platform, expand later (specify which)
- [ ] (F) Other (describe)

## 9. Installation and Distribution - How should users install gws?

What installation methods are important?

- [X] (A) Single binary download (no dependencies)
- [ ] (B) Package managers (brew, apt, cargo, npm, etc.)
- [X] (C) Build from source
- [ ] (D) All of the above
- [ ] (E) Other (describe)

## 10. Technology Stack - Any preferences or constraints?

Do you have preferences for the implementation language or technology?

- [ ] (A) Rust (performance, single binary, cross-platform)
- [X] (B) Go (simple deployment, good CLI tools)
- [ ] (C) Python (easy to extend, good ecosystem)
- [ ] (D) Node.js/TypeScript (JavaScript ecosystem)
- [ ] (E) Shell script (maximum compatibility)
- [ ] (F) No preference - choose based on requirements
- [X] (G) Prefer Go or possibliy rust, but if other options are more appropriate, prompt me and describe why

## 11. Proof Artifacts - How will you demo this working?

What would you want to show to prove this feature works?

- [X] (A) Terminal recording (asciinema) showing CLI commands
- [X] (B) Screenshots of command output
- [ ] (C) Example repository list with search/filter results
- [ ] (D) Before/after showing repository organization
- [ ] (E) Performance benchmarks (searching across many repos)
- [X] (F) Readme with manual steps to verify each expected functionality
