# 20-spec-worktree-management

## Introduction/Overview

This feature adds git worktree awareness and management to gws. Users who work with multiple branches simultaneously via `git worktree` will be able to discover, navigate, organize, and create worktrees through the gws CLI. Worktrees follow a standard directory convention: the main repo lives at `<repo-directory-name>` and all its worktrees live under a sibling directory `<repo-directory-name>.wt/`.

## Goals

- Discover all worktrees for tracked repos during `gws refresh` using `git worktree list`
- Display a `(wt)` indicator inline next to repo names in `gws list` output for repos that have worktrees
- Enable navigation to worktrees via `gws <repo> -wt <branch-name>` with partial matching and interactive selection
- Provide `gws worktree align` to reorganize existing worktrees into the `<repo>.wt/` convention using `git worktree move`
- Provide `gws worktree add` and `gws worktree list` subcommands for creating and viewing worktrees

## User Stories

- **As a developer working on multiple branches**, I want to navigate to a worktree by repo and branch name so that I can switch contexts quickly without remembering directory paths.
- **As a developer with many repos**, I want to see at a glance which repos have worktrees so that I know where I have active parallel work.
- **As a developer adopting gws**, I want to align my existing worktrees to a consistent directory structure so that my workspace is organized and predictable.
- **As a developer**, I want to create new worktrees that automatically follow the directory convention so that I don't have to manage paths manually.
- **As a developer**, I want to list all worktrees across my workspace (or for a specific repo) so that I can see what branches I have checked out where.

## Demoable Units of Work

### Unit 1: Worktree Discovery and Display

**Purpose:** Enable gws to discover worktrees during refresh and show them in list output, giving users visibility into their worktree landscape.

**Functional Requirements:**
- The system shall run `git worktree list --porcelain` for each tracked repo during `gws refresh` and store discovered worktrees in the `config.Repository` struct
- Each worktree entry shall store: path, branch name, and whether it is aligned (i.e., located inside `<repo>.wt/`)
- The system shall persist worktree data in `config.json` as a `Worktrees` field on each `Repository`
- The `gws list` output shall display a `(wt)` indicator inline next to the repo name for any repo that has one or more worktrees
- The `(wt)` indicator shall appear in both compact (multi-column) and table display modes
- The refresh summary shall include the number of repos with worktrees discovered

**Proof Artifacts:**
- Test: worktree discovery correctly parses `git worktree list --porcelain` output and populates the struct
- Test: `(wt)` indicator appears in list output for repos with worktrees and is absent for repos without
- Test: refresh stores and retrieves worktree data from config

### Unit 2: Worktree Navigation

**Purpose:** Let users navigate to worktrees using the familiar `gws <repo>` pattern extended with a `-wt` flag, including interactive selection when no specific branch is given.

**Functional Requirements:**
- The system shall accept a `--worktree` / `-wt` flag when navigating to a repo (e.g., `git-workspace -g <repo> --worktree <branch-name>`)
- The shell function shall detect `-wt` in the argument list and pass it through as `--worktree` to the binary, then `cd` to the resulting path
- When `<branch-name>` is provided, the system shall match it against worktree branch names using the same partial/substring matching logic used for repo name matching (`filter.MatchesPattern`)
- When a single worktree matches, the system shall print the worktree path to stdout (same pattern as repo navigation)
- When multiple worktrees match, the system shall display a numbered selection list showing branch name and path, and prompt the user to choose (same interactive pattern as multi-repo matching)
- When no worktrees match, the system shall display an error with suggestions (same pattern as no-match repo navigation)
- When `-wt` is used without a branch name, the system shall list all worktrees for that repo as a numbered selection list
- The system shall navigate to unaligned worktrees (not in `.wt/` directory) — alignment is not required for navigation

**Proof Artifacts:**
- Test: single worktree match prints correct path to stdout
- Test: multiple worktree matches display numbered selection list and accept input
- Test: no match displays error with suggestions
- Test: bare `-wt` (no branch name) lists all worktrees for the repo
- Test: shell function correctly passes `-wt` flag and performs `cd`

### Unit 3: Worktree Subcommands — `list`, `add`, `align`

**Purpose:** Provide direct worktree management commands for listing, creating, and organizing worktrees.

**Functional Requirements:**

**`gws worktree list [repo]`:**
- The system shall display all worktrees across all tracked repos when no argument is given
- The system shall filter to a single repo's worktrees when a repo name argument is provided
- Each worktree entry shall show: repo name, branch name, path, and an `(unaligned)` indicator if the worktree is not inside the `<repo>.wt/` directory
- The output shall be formatted as a readable table

**`gws worktree add <repo> <branch>`:**
- The system shall create a new worktree at `<repo-path>.wt/<branch>` using `git worktree add`
- The system shall create the `<repo>.wt/` directory if it does not already exist
- The branch name shall be used as-is for the directory name inside `.wt/` (e.g., `feature/auth-flow` becomes `<repo>.wt/feature/auth-flow`)
- The system shall update the stored worktree data in config after creation
- The system shall error if the repo name does not match a tracked repo or if the branch already has a worktree

**`gws worktree align [repo]`:**
- The system shall move all unaligned worktrees into `<repo-path>.wt/` using `git worktree move` (requires Git 2.17+)
- When a repo name argument is provided, the system shall only align worktrees for that repo
- When no argument is provided, the system shall align worktrees for all tracked repos
- If two worktrees would produce the same directory name inside `.wt/`, the system shall append a `-dup-NN` suffix (where NN is 00-99, incremented per duplicate)
- The system shall support `--dry-run` to preview moves without executing them, including flagging potential duplicate names
- The system shall display a summary of moves performed (or previewed in dry-run mode): source path, destination path, and any renames due to duplicates
- The system shall update stored worktree data in config after moving

**Proof Artifacts:**
- Test: `worktree list` shows all worktrees with correct repo grouping and `(unaligned)` indicator
- Test: `worktree list <repo>` filters to specified repo
- Test: `worktree add` creates worktree in correct `.wt/` location and updates config
- Test: `worktree align` moves worktrees and updates config
- Test: `worktree align --dry-run` previews moves without executing
- Test: duplicate directory names receive `-dup-NN` suffix

### Unit 4: Shell Integration Update

**Purpose:** Update the shell function to support `-wt` navigation so that `gws <repo> -wt <branch>` performs a `cd` to the worktree directory.

**Functional Requirements:**
- The zsh shell function shall detect `-wt` as $2 or $3 and route through `git-workspace -g <repo> --worktree <branch> -q` with `cd` on the result
- The bash shell function shall have the same behavior
- The `worktree` subcommand shall be added to the passthrough list in the shell function's `case` statement (alongside `list`, `init`, `add`, etc.)
- Tab completion shall be updated to include `worktree` as a subcommand and `--worktree` as a flag

**Proof Artifacts:**
- Test: shell-init output contains `-wt` handling logic
- Test: `worktree` appears in the passthrough command list

## Non-Goals (Out of Scope)

1. **Worktree deletion/removal**: `gws worktree remove` is not included in this spec. Users can use `git worktree remove` directly.
2. **Worktree status in `gws list`**: Showing git status (clean/dirty/branch ahead-behind) for individual worktrees in list output is not included. The `(wt)` indicator only shows presence.
3. **Auto-alignment on refresh**: `gws refresh` discovers and tracks worktrees but does not move them. Moving requires explicit `gws worktree align`.
4. **Support for Git versions below 2.17**: The `git worktree move` command requires Git 2.17+. No fallback to manual file manipulation is provided.
5. **Worktree-aware filtering in `gws list`**: No `-w` filter flag to show only repos with worktrees. This can be added later.

## Design Considerations

No specific UI/UX design requirements beyond matching existing gws patterns:
- The `(wt)` indicator follows the same inline style as the existing `(type)` display in navigation output
- Worktree selection lists follow the exact same numbered-list-with-prompt pattern used in `navigate.go` for multi-repo matches
- `worktree list` table output should follow the same formatting conventions as `gws list` table mode

## Repository Standards

- **Go conventions**: Follow existing patterns in `cmd/git-workspace/` for command files and `internal/` for packages
- **Cobra command structure**: New commands follow the same registration pattern (`init()` + `rootCmd.AddCommand` or parent `.AddCommand`)
- **Testing**: Use Go standard `testing` package, `bytes.Buffer` for I/O capture, `t.TempDir()` for filesystem tests, `t.Cleanup()` for teardown
- **Shell templates**: Updates to `shellinit.go` templates follow the existing `{BIN}` placeholder pattern
- **Config versioning**: If the config struct changes shape, bump the config version appropriately
- **Commit messages**: Use conventional commits (`feat:`, `fix:`, etc.) — `feat:` triggers a minor version bump via Release Please

## Technical Considerations

- **`git worktree list --porcelain`**: Use porcelain format for reliable parsing. Output format is well-defined: `worktree <path>\nHEAD <sha>\nbranch <ref>\n\n` per entry.
- **Config struct growth**: Adding a `Worktrees []Worktree` field to `Repository` is reasonable. Each worktree entry is small (path, branch, aligned bool). If a user has dozens of worktrees per repo this could grow, but that's an unusual case.
- **`git worktree move` availability**: Requires Git 2.17+ (released April 2018). The spec assumes this is available and does not provide a fallback. The `align` command should check Git version and error clearly if too old.
- **Shell function complexity**: The `-wt` flag adds a new branch to the shell `case` logic. Option B (handle as a Go flag with shell passthrough) is preferred as it keeps logic in Go where it's testable and avoids duplicating matching logic in shell.
- **Workspace symlinks**: For repos outside the workspace that are tracked via symlink, worktree discovery should use the real (resolved) path from the repo, not the symlink path.

## Security Considerations

No specific security considerations identified. Worktree operations are local filesystem operations using standard git commands. No credentials, tokens, or sensitive data are involved.

## Success Metrics

1. **All tracked repos with worktrees show `(wt)` indicator** in `gws list` output after `gws refresh`
2. **Navigation via `gws <repo> -wt <branch>` works** with partial matching, interactive selection, and `cd`
3. **`gws worktree align` correctly moves worktrees** into `<repo>.wt/` and updates git internal pointers via `git worktree move`
4. **`gws worktree add` creates worktrees** in the correct convention location
5. **Test coverage** for all new worktree-related code (discovery, navigation, subcommands, shell integration)

## Open Questions

No open questions at this time.
