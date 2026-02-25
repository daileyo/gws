# 01-spec-init-add-refactor.md

## Introduction/Overview

The `--init` command currently serves dual purposes — workspace creation and repository discovery — and takes a required path argument that makes its intent ambiguous. This spec refactors `--init` into a focused, zero-argument workspace creation command and introduces a new `--add` command for adding repositories to an existing workspace. The `--add` command also manages symlinks inside the gws directory for repositories located outside the workspace root, so the workspace directory is fully navigable on the filesystem without using gws.

## Goals

- Clarify `--init` to be a single-purpose workspace creation command that always targets the current directory
- Introduce `--add` as the canonical command for adding repositories to an existing workspace
- Support a gws workspace directory that contains both real repositories and symlinks to repositories located elsewhere
- Enable shell tab completion for path arguments passed to `--add`
- Reassign short flags to reflect new command roles: `-a` → `--add`, `--add-tag` → `-d`, `-v` → `--recursive`
- Establish a flag shorthand priority principle: subcommand flags retain short forms before base/utility flags do

## User Stories

**As a developer**, I want to run `gws --init` in my projects directory so that gws creates a workspace there and automatically discovers any git repositories present, without having to supply a path argument.

**As a developer**, I want `gws --init` to clearly tell me when a workspace is already initialized so that I know to use `--add` or `--refresh` instead of re-running init.

**As a developer**, I want to run `gws --add` while inside a git repository so that it gets tracked in my workspace, with a symlink created in my gws directory if the repo lives outside the workspace root.

**As a developer**, I want to run `gws --add --recursive` in a directory containing multiple git repositories so that all of them are added in one step, with a clear summary of what was added.

**As a developer**, I want `gws --add <Tab>` to autocomplete filesystem paths in my shell so that I can specify repository paths quickly without typing them in full.

## Demoable Units of Work

### Unit 1: Refactored `--init` Command

**Purpose:** Simplifies workspace initialization to always target the current directory, with clear feedback when a workspace is already known to gws.

**Functional Requirements:**
- The system shall check for an existing gws workspace by reading `~/.gws/config.json` before creating one.
- The system shall create the workspace at the current working directory when no workspace exists.
- The system shall scan the current directory for git repositories and add them to `~/.gws/config.json` during initialization.
- The system shall display a confirmation after successful initialization that includes the workspace path and the count of discovered repositories (e.g., `Initialized workspace at /path. Found 5 repositories.`).
- When a workspace already exists, the system shall print a notification that includes: the existing workspace path, a suggestion to use `--add` to add more repositories, and a suggestion to use `--refresh` to re-scan the existing workspace.
- When a workspace already exists, the system shall exit with code 0.
- The `--init` flag shall change from accepting a path string argument to a boolean flag (no value required).

**Proof Artifacts:**
- CLI: `gws --init` run in a directory containing git repos outputs `Initialized workspace at [path]. Found N repositories.` demonstrates successful initialization with repo discovery.
- CLI: Running `gws --init` a second time outputs a notification with the existing workspace path and suggestions for `--add` and `--refresh` demonstrates the already-initialized guard.
- CLI: `gws --help` shows `--init` listed as a boolean flag with no path argument demonstrates the flag signature change.

---

### Unit 2: New `--add` Command — Single Repository

**Purpose:** Provides a focused command for adding an individual repository to an existing workspace, creating a symlink in the gws directory when the repo is located outside the workspace root.

**Functional Requirements:**
- The system shall register `--add` / `-a` as a new string flag with an optional value, defaulting to the current working directory when no value is supplied.
- The system shall reassign `--add-tag` to use `-d` as its short form (replacing the previous `-a`).
- When `gws --add` is run with no argument, the system shall use the current working directory as the target path.
- When `gws --add <path>` is run with a path argument, the system shall resolve and use that path as the target.
- The system shall validate that the target path is a git repository (contains a `.git` directory); if it is not, the system shall print an error and exit non-zero.
- The system shall add the validated repository to `~/.gws/config.json`.
- When the target repository is located outside the gws workspace directory, the system shall create a symlink inside the workspace directory pointing to the target repository. The symlink name shall match the repository's directory name.
- When a file or symlink already exists at the intended symlink path in the workspace, the system shall warn the user and skip symlink creation without exiting non-zero.
- When the target repository is already tracked in `~/.gws/config.json`, the system shall print a warning (e.g., `my-repo is already tracked, skipping.`) and exit 0.
- The system shall print a confirmation after successfully adding a repository (e.g., `Added my-repo to workspace.`).
- The system shall print an additional line if a symlink was created (e.g., `Created symlink: /workspace/my-repo → /external/path/my-repo`).

**Proof Artifacts:**
- CLI: `gws --add` run inside a git repository outputs `Added [repo-name] to workspace.` demonstrates adding a repo from the current directory.
- CLI: `gws --add /path/to/external/repo` outputs `Added [repo-name] to workspace.` followed by `Created symlink: [workspace]/[repo-name] → /path/to/external/repo` demonstrates external repo add with symlink creation.
- CLI: `gws --add` run in a non-git directory outputs an error message and exits non-zero demonstrates the git validation check.
- CLI: Running `gws --add` twice in the same git repo outputs the already-tracked warning on the second run demonstrates duplicate detection.
- Filesystem: `ls -la [workspace-path]` shows a symlink entry for the external repo demonstrates the symlink exists on disk.

---

### Unit 3: `--add --recursive` — Batch Repository Add

**Purpose:** Enables adding all git repositories found in the current directory in a single command, with a clear summary of results.

**Functional Requirements:**
- The system shall register `--recursive` / `-v` as a boolean modifier flag. `-v` is used as the short form because `-r` is reserved for `--refresh`; per the project's shorthand priority principle, the less-frequently-used base command (`--refresh`) should yield its short form to subcommand modifiers before subcommands do — this is deferred to a future refactor, and `-v` is used in the interim.
- The system shall return an error if `--recursive` / `-v` is used without `--add`.
- When `gws --add --recursive` is run, the system shall scan the current directory for all git repositories using the same non-nested scan logic as the existing discovery scanner.
- The system shall add all discovered repositories to `~/.gws/config.json`.
- For any discovered repository located outside the gws workspace directory, the system shall create a symlink in the workspace directory (same logic as Unit 2).
- Already-tracked repositories encountered during the recursive scan shall be warned and skipped (same behavior as Unit 2 single add).
- The system shall report the count and names of newly added repositories (e.g., `Added 3 repositories: repo-a, repo-b, repo-c.`).
- When no new repositories are found (all already tracked or none present), the system shall output `No new repositories found.`

**Proof Artifacts:**
- CLI: `gws --add --recursive` in a directory containing multiple git repos outputs `Added N repositories: [names]` demonstrates batch add with count and names.
- CLI: `gws --add --recursive` when all repos are already tracked outputs `No new repositories found.` demonstrates idempotent behavior.
- CLI: `gws --recursive` / `gws -v` without `--add` returns an error message and exits non-zero demonstrates modifier flag validation.

---

### Unit 4: Shell Tab Completion

**Purpose:** Enables native shell tab completion for gws commands, particularly filesystem path arguments for `--add`.

**Functional Requirements:**
- The system shall register a `completion` subcommand using Cobra's built-in completion support.
- The `completion` subcommand shall support generating scripts for bash, zsh, and fish shells.
- The `--add` flag shall be registered with a `ValidArgsFunction` that provides filesystem directory path completion.
- Shell completion for `--add` shall suggest directories when the user presses `<Tab>` after the flag.

**Proof Artifacts:**
- CLI: `gws completion zsh` outputs a valid zsh completion script demonstrates completion script generation.
- Shell demo: After sourcing the zsh completion script, typing `gws --add ~/pro<Tab>` auto-completes to matching directories demonstrates live path completion.

## Non-Goals (Out of Scope)

1. **In-app interactive fuzzy directory picker**: An interactive TUI for browsing and selecting directories when running `--add` is a future feature and not included in this spec.
2. **Symlink lifecycle management in `--refresh`**: Updating or removing stale symlinks during `gws --refresh` is not addressed by this spec.
3. **Workspace relocation**: Moving or renaming an existing workspace directory is out of scope.
4. **Nested repository support**: Discovering git repositories nested inside other git repositories is not supported, consistent with existing scanner behavior.
5. **Multiple workspace support**: gws continues to support a single workspace per user; this spec does not introduce multi-workspace support.

## Design Considerations

No specific UI/UX design requirements beyond CLI output. Status messages and informational output shall be written to stderr; repository paths and machine-readable output shall go to stdout. This is consistent with existing behavior in `navigate.go`.

## Repository Standards

- **Language and toolchain**: Go 1.23.0; use `go test ./...` for all tests
- **CLI framework**: Cobra v1.10.2 — new flags registered on the root command via `rootCmd.Flags()`; `--add` and `--recursive` / `-v` follow the same flag-based pattern as existing commands; the `completion` subcommand is the only Cobra subcommand and follows Cobra convention
- **Git operations**: Use `go-git/go-git` for repository detection (consistent with `internal/discovery/scanner.go`)
- **File organization**: New command logic goes in `cmd/git-workspace/add.go`; the refactored init goes in the existing `cmd/git-workspace/init.go`; colocated test files (`add_test.go`, updated `init_test.go`)
- **Test patterns**: Standard Go test files using `t.TempDir()` for filesystem isolation; table-driven tests where appropriate; cover happy path, error cases, and edge cases
- **Config access**: All reads and writes use the existing `internal/config` package (`config.Load()`, `config.Save()`, `config.Exists()`)
- **Commit message format**: Follow existing conventional commit style (e.g., `feat: add --add command with recursive and symlink support`)
- **Build**: All changes must pass `make ci` (vet + lint + test-race)

## Technical Considerations

- **`--init` flag type change**: `--init` must change from `StringVarP` (takes a path value) to `BoolVarP` (boolean). This is a clean break — no deprecation shim. Document in release notes.
- **`--add` as an optional-value string flag**: Use Cobra's `NoOptDefVal` to allow `--add` to function both as `gws --add` (defaults to `.`) and `gws --add /some/path`. Example:
  ```go
  rootCmd.Flags().StringVarP(&flagAdd, "add", "a", "", "Add a repository to the workspace")
  rootCmd.Flags().Lookup("add").NoOptDefVal = "."
  ```
- **`--recursive` short flag**: `-v` is assigned as the short form for `--recursive`. `-r` is unavailable (reserved for `--refresh`). Future refactor may reassign `-r` to `--recursive` and give `--refresh` a less prominent short form, per the shorthand priority principle.
- **Mutual exclusivity**: `--add` must be added to the existing mutual exclusivity check in `main.go` alongside `--list`, `--init`, `--add-tag`, etc. `--recursive` must be validated to only be used in combination with `--add` (similar to how filter flags are gated on `--list`).
- **Symlink creation**: Use `os.Symlink(targetPath, symlinkPath)` where `symlinkPath` is `filepath.Join(workspacePath, repoName)`. Check for existing file/symlink at that path before creating.
- **Cobra completion subcommand**: Register using `rootCmd.AddCommand(rootCmd.NewDefaultCompletionCmd(rootCmd.Name()))` or Cobra's `GenBashCompletion`/`GenZshCompletion` methods. Register `ValidArgsFunction` on the root command to provide directory completion for the `--add` flag value.
- **Workspace check on `--add`**: Before adding a repository, verify a workspace is initialized using `config.Exists()`. If not, print a helpful error directing the user to run `gws --init` first.

## Security Considerations

- **Symlink targets**: Symlinks shall only be created to paths explicitly provided by the user. No automatic symlink chain traversal or resolution beyond the immediate target.
- **File permissions**: The gws workspace directory, if newly created during `--init`, shall use mode `0755`. Symlinks inherit the target's permissions.
- No API keys, tokens, or credentials are involved in this feature.
- **Proof artifacts**: No sensitive data (tokens, private paths) should appear in CLI output used as proof artifacts.

## Success Metrics

1. **Behavioral clarity**: `gws --init` and `gws --add` each have a single, unambiguous purpose with no behavioral overlap.
2. **Idempotency**: Running `gws --init` or `gws --add` multiple times on the same target produces consistent, non-destructive output without corrupting workspace state.
3. **Test coverage**: `init.go` (refactored) and `add.go` (new) each have colocated test files covering the happy path, error cases, and the following edge cases: already-tracked repo, non-git directory, external repo symlink creation, symlink path conflict.
4. **Shell completion functional**: After sourcing the generated completion script, `<Tab>` on `gws --add` provides directory path suggestions in at least zsh and bash.

## Open Questions

No open questions at this time.
