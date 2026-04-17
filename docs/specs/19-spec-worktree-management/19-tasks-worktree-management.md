# 19-tasks-worktree-management

## Relevant Files

- `internal/config/config.go` - Add `Worktree` struct and `Worktrees` field to `Repository`; bump config version
- `internal/config/config_test.go` - Tests for worktree serialization/deserialization in config
- `internal/git/worktree.go` - New file: parse `git worktree list --porcelain`, worktree move/add operations
- `internal/git/worktree_test.go` - Tests for porcelain parsing, alignment detection, move, add
- `cmd/git-workspace/refresh.go` - Add Phase 3: worktree discovery after repo validation/scan
- `cmd/git-workspace/refresh_test.go` - Tests for worktree discovery during refresh
- `cmd/git-workspace/list.go` - Add `(wt)` indicator to repo name in compact and table display
- `cmd/git-workspace/list_test.go` - Tests for `(wt)` indicator presence/absence
- `cmd/git-workspace/navigate.go` - Add `--worktree` flag handling and worktree selection logic
- `cmd/git-workspace/navigate_test.go` - Tests for worktree navigation (single, multi, no match, bare flag)
- `cmd/git-workspace/main.go` - Register `--worktree` flag on root command, pass to `runNavigate`
- `cmd/git-workspace/main_test.go` - Tests for `--worktree` flag registration
- `cmd/git-workspace/worktree.go` - New file: `gws worktree` parent command registration
- `cmd/git-workspace/worktree_list.go` - New file: `gws worktree list [repo]` subcommand
- `cmd/git-workspace/worktree_list_test.go` - Tests for worktree list output and repo filtering
- `cmd/git-workspace/worktree_add.go` - New file: `gws worktree add <repo> <branch>` subcommand
- `cmd/git-workspace/worktree_add_test.go` - Tests for worktree creation in `.wt/` convention
- `cmd/git-workspace/worktree_align.go` - New file: `gws worktree align [repo]` with `--dry-run`
- `cmd/git-workspace/worktree_align_test.go` - Tests for align, dry-run, and `-dup-NN` conflict handling
- `cmd/git-workspace/shellinit.go` - Update zsh/bash templates for `-wt` flag and `worktree` passthrough
- `cmd/git-workspace/shellinit_test.go` - Tests for updated shell templates

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `worktree.go` and `worktree_test.go` in the same directory).
- Use `go test ./...` to run the full test suite; use `go test ./internal/git/` or `go test ./cmd/git-workspace/` to target specific packages.
- Follow the existing Cobra command registration pattern: define `var xxxCmd` + `func init()` with `rootCmd.AddCommand` or `parentCmd.AddCommand`.
- Follow the existing `gitCommand()` helper in `internal/git/status.go` for running git CLI commands.
- Use `t.TempDir()` and `t.Cleanup()` for filesystem-dependent tests, and `bytes.Buffer`/`strings.NewReader` for I/O capture.
- Use conventional commits: `feat:` for new functionality.

## Tasks

### [x] 1.0 Worktree Data Model and Discovery

Add the `Worktree` struct to the config package, implement `git worktree list --porcelain` parsing, and wire discovery into `gws refresh`.

#### 1.0 Proof Artifact(s)

- Test: `internal/git/worktree_test.go` — parsing of `git worktree list --porcelain` output correctly produces worktree entries with path, branch, and aligned status
- Test: `internal/config/config_test.go` — worktree data round-trips through JSON serialization
- Test: `cmd/git-workspace/refresh_test.go` — refresh discovers worktrees and persists them in config
- CLI: `gws refresh` output includes worktree count in the summary line

#### 1.0 Tasks

- [x] 1.1 Add `Worktree` struct to `internal/config/config.go` with fields: `Path string`, `Branch string`, `Aligned bool`. Add `Worktrees []Worktree` field (with `json:"worktrees,omitempty"`) to the `Repository` struct.
- [x] 1.2 Add serialization tests in `internal/config/config_test.go` verifying that a `Repository` with `Worktrees` round-trips through `json.Marshal`/`json.Unmarshal` correctly, and that `Worktrees` is omitted from JSON when empty.
- [x] 1.3 Create `internal/git/worktree.go` with a `ListWorktrees(repoPath string) ([]WorktreeEntry, error)` function that runs `git worktree list --porcelain` and parses the output. Each `WorktreeEntry` should contain `Path` and `Branch` (extracted from `branch refs/heads/...`). Skip the main worktree entry (the one whose path matches `repoPath`).
- [x] 1.4 Add an `IsAligned(worktreePath, repoPath string) bool` function in `internal/git/worktree.go` that checks whether the worktree path is inside `<repoPath>.wt/`.
- [x] 1.5 Write tests in `internal/git/worktree_test.go` for `ListWorktrees` using a real temp git repo with `git worktree add`, verifying correct path/branch extraction and that the main worktree is excluded.
- [x] 1.6 Write tests for `IsAligned` covering: worktree inside `.wt/` (aligned), worktree elsewhere (unaligned), edge cases (similar directory names that shouldn't match).
- [x] 1.7 Add a worktree discovery phase to `runRefresh()` in `cmd/git-workspace/refresh.go`: after `detectUserForRepos`, iterate all repos, call `ListWorktrees`, map results to `config.Worktree` entries (setting `Aligned` via `IsAligned`), and assign to each `repo.Worktrees`.
- [x] 1.8 Add a worktree count to the refresh summary output (e.g., `"Repositories with worktrees: %d"`), printed only when count > 0.
- [x] 1.9 Write tests in `cmd/git-workspace/refresh_test.go` verifying that after refresh, repos with worktrees have populated `Worktrees` fields and repos without remain empty.

### [x] 2.0 Worktree Indicator in List Output

Display a `(wt)` indicator inline next to repo names in `gws list` for repos that have worktrees.

#### 2.0 Proof Artifact(s)

- Test: `cmd/git-workspace/list_test.go` — repos with worktrees show `(wt)` suffix in compact multi-column mode
- Test: `cmd/git-workspace/list_test.go` — repos with worktrees show `(wt)` suffix in table mode
- Test: repos without worktrees do not show the indicator

#### 2.0 Tasks

- [x] 2.1 In the compact multi-column path of `displayTable` (`list.go`), when building the `names` slice, append ` (wt)` to `repo.Name` if `len(repo.Worktrees) > 0`.
- [x] 2.2 In the table display path of `displayTable`, when rendering the NAME column for each row, append ` (wt)` to the name if the repo has worktrees. Account for the extra width in `maxNameLen` calculation.
- [x] 2.3 In `displayJSON`, include the `(wt)` indicator in the name field if the repo has worktrees (or add a separate `"has_worktrees": true` field — match whichever is more consistent with existing JSON patterns).
- [x] 2.4 Write tests in `cmd/git-workspace/list_test.go`: create test repos where some have `Worktrees` populated and verify the output contains `repo-name (wt)` for those and plain `repo-name` for others, in both compact and table modes.

### [x] 3.0 Worktree Navigation (`-wt` flag)

Add a `--worktree` / `-wt` flag to the root command for navigating to worktrees, with partial matching and interactive selection.

#### 3.0 Proof Artifact(s)

- Test: `cmd/git-workspace/navigate_test.go` — single worktree match prints correct path to stdout
- Test: multiple worktree matches display numbered selection list and accept input
- Test: no match displays error with suggestions
- Test: bare `--worktree` (no branch name) lists all worktrees for selection

#### 3.0 Tasks

- [x] 3.1 Add a `--worktree` string flag (no shorthand — `-w` taken by deprecated flag) to the root command in `main.go` `init()`. Use `NoOptDefVal` set to a sentinel (e.g., `"\x00wt"`) so bare `--worktree` (no value) is distinguishable from not provided.
- [x] 3.2 Update the navigation block in `rootCmd.RunE` (`main.go`): when `--worktree` is set and `len(args) > 0`, call a new `runWorktreeNavigate` function instead of `runNavigate`.
- [x] 3.3 Create `runWorktreeNavigate(repoQuery, worktreeQuery string, quiet bool, repos []config.Repository, stderr, stdout io.Writer, stdin io.Reader) error` in `navigate.go`. This function should: find the matching repo (reuse existing `filter.MatchesPattern` logic, error if no match or multiple repo matches), then match `worktreeQuery` against the repo's `Worktrees` branch names using `filter.MatchesPattern`.
- [x] 3.4 Implement the match cases in `runWorktreeNavigate`: (a) single match → print worktree path to stdout, (b) multiple matches → display numbered list with branch name and path, prompt for selection (reuse `handleMultipleMatches` pattern), (c) no match → display error with branch name suggestions, (d) bare flag (sentinel value) → list all worktrees for the repo as a numbered selection.
- [x] 3.5 Write tests in `navigate_test.go` for `runWorktreeNavigate`: single match, multiple matches with TTY selection, no match with suggestions, bare flag listing all worktrees, and non-TTY behavior.

### [ ] 4.0 Worktree Subcommands (`list`, `add`, `align`)

Implement the `gws worktree` parent command with `list`, `add`, and `align` subcommands.

#### 4.0 Proof Artifact(s)

- Test: `cmd/git-workspace/worktree_list_test.go` — lists worktrees with `(unaligned)` indicator, filters by repo name
- Test: `cmd/git-workspace/worktree_add_test.go` — creates worktree in `.wt/` directory, updates config
- Test: `cmd/git-workspace/worktree_align_test.go` — moves worktrees, handles duplicates with `-dup-NN`, `--dry-run` previews without moving
- CLI: `gws worktree list` shows table output with repo, branch, path, alignment status

#### 4.0 Tasks

- [ ] 4.1 Create `cmd/git-workspace/worktree.go` with a `worktreeCmd` parent command (`Use: "worktree"`, `Short: "Manage git worktrees"`). Register it in `init()` with `rootCmd.AddCommand(worktreeCmd)`.
- [ ] 4.2 Create `cmd/git-workspace/worktree_list.go` with `worktreeListCmd` (`Use: "list [repo]"`, `Args: cobra.MaximumNArgs(1)`). In `runWorktreeList(repoFilter string)`: load config, iterate repos, collect worktrees, filter by repo name if provided. Display as a table with columns: REPO, BRANCH, PATH, and `(unaligned)` indicator. Register with `worktreeCmd.AddCommand`.
- [ ] 4.3 Write tests in `worktree_list_test.go`: all worktrees listed when no filter, filtered to one repo with argument, `(unaligned)` shown for worktrees outside `.wt/`, empty output when no worktrees exist.
- [ ] 4.4 Add `AddWorktree(repoPath, branch, destPath string) error` function to `internal/git/worktree.go` that runs `git -C <repoPath> worktree add <destPath> <branch>`.
- [ ] 4.5 Create `cmd/git-workspace/worktree_add.go` with `worktreeAddCmd` (`Use: "add <repo> <branch>"`, `Args: cobra.ExactArgs(2)`). In `runWorktreeAdd(repoName, branch string)`: find repo by name, compute dest as `<repoPath>.wt/<branch>`, create `.wt/` dir if needed via `os.MkdirAll`, call `AddWorktree`, re-discover worktrees for that repo, save config. Register with `worktreeCmd.AddCommand`.
- [ ] 4.6 Write tests in `worktree_add_test.go`: worktree created at correct path, `.wt/` directory created, config updated with new worktree entry, error on unknown repo, error on duplicate branch.
- [ ] 4.7 Add `MoveWorktree(repoPath, currentPath, newPath string) error` function to `internal/git/worktree.go` that runs `git -C <repoPath> worktree move <currentPath> <newPath>`.
- [ ] 4.8 Create `cmd/git-workspace/worktree_align.go` with `worktreeAlignCmd` (`Use: "align [repo]"`, `Args: cobra.MaximumNArgs(1)`). Add `--dry-run` bool flag. In `runWorktreeAlign(repoFilter string, dryRun bool)`: load config, iterate repos (filtered if argument provided), for each unaligned worktree compute target path as `<repoPath>.wt/<branch>`. Handle naming conflicts by appending `-dup-NN` (00-99). If `dryRun`, print planned moves without executing. Otherwise, create `.wt/` dir, call `MoveWorktree`, re-discover worktrees, save config.
- [ ] 4.9 Write tests in `worktree_align_test.go`: unaligned worktree gets moved to `.wt/`, aligned worktree is skipped, `--dry-run` prints plan without moving, duplicate names get `-dup-NN` suffix, filtered by repo name when argument provided.
- [ ] 4.10 Write tests in `internal/git/worktree_test.go` for `AddWorktree` and `MoveWorktree` using real temp git repos.

### [ ] 5.0 Shell Integration Update

Update the shell function templates to support `-wt` navigation and pass through the `worktree` subcommand.

#### 5.0 Proof Artifact(s)

- Test: `cmd/git-workspace/shellinit_test.go` — generated shell function contains `-wt` handling logic
- Test: `worktree` appears in the passthrough case statement
- Test: both zsh and bash templates produce correct output

#### 5.0 Tasks

- [ ] 5.1 Add `worktree` to the passthrough command list in both `zshInitTemplate` and `bashInitTemplate` in `shellinit.go` (the `case "$1"` line alongside `list|init|add|refresh|...`).
- [ ] 5.2 Update the `*` (navigation) case in the zsh template to detect `-wt` as `$2`: if `$2` is `-wt`, route through `{BIN} -g "$1" --worktree "$3" -q 2>/dev/tty </dev/tty` with `cd` on the result. If `$2` is `-wt` and `$3` is empty, route through `{BIN} -g "$1" --worktree -q 2>/dev/tty </dev/tty` (bare flag triggers selection).
- [ ] 5.3 Apply the same `-wt` handling to the bash template.
- [ ] 5.4 Write tests in `shellinit_test.go` verifying: (a) `worktree` is in the passthrough list for both zsh and bash, (b) the `-wt` handling block exists in both templates, (c) the generated output includes `--worktree` flag in the navigation case.
- [ ] 5.5 Update the `rootCmd` usage template in `main.go` to include worktree navigation example (e.g., `gws my-repo -wt feature-branch    # Navigate to worktree`).
