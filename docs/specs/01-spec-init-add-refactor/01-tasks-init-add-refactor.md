# 01-tasks-init-add-refactor.md

## Relevant Files

- `cmd/git-workspace/main.go` - Root command, flag variable declarations, flag registration, mutual exclusivity check, command dispatch, and usage template. Modified in every task.
- `cmd/git-workspace/main_test.go` - Root command tests for flag validation and mutual exclusivity. Update throughout to cover new flags.
- `cmd/git-workspace/init.go` - `--init` command handler. Refactored in Task 1 to use `os.Getwd()` and add already-initialized guard.
- `cmd/git-workspace/init_test.go` - Tests for the init handler. Updated in Task 1.
- `cmd/git-workspace/add.go` - New `--add` command handler. Created in Task 2, extended in Task 3.
- `cmd/git-workspace/add_test.go` - Tests for the add handler. Created in Task 2, extended in Task 3.
- `cmd/git-workspace/tag.go` - `--add-tag` handler. No code changes required; the short-form reassignment (`-a` → `-d`) is made in `main.go`.
- `internal/config/config.go` - Config struct, `Load()`, `Save()`, `Exists()`, `New()` functions. Read-only reference; no changes needed.
- `internal/discovery/scanner.go` - `Scan()` function used by the recursive add path. Read-only reference; no changes needed.
- `internal/classifier/detector.go` - `Classify()` function used to detect repo type and visibility when adding a single repo. Read-only reference; no changes needed.

### Notes

- Run tests with `go test ./...` or use `make test` for verbose output.
- Use `make ci` (vet + lint + test-race) as the quality gate before marking any task complete.
- All test files live alongside their source files (e.g., `add.go` and `add_test.go` in the same directory).
- Use `t.TempDir()` in tests to create isolated temporary directories; the directory is automatically cleaned up after the test.
- Write table-driven tests (using a `tests := []struct{...}` slice) wherever multiple input/output combinations are tested.

## Tasks

### [x] 1.0 Refactor `--init` Command and Reassign Flag Short Forms

Refactor the `--init` flag from a string flag (requiring a path argument) to a boolean flag (always using the current directory). Reassign `--add-tag`'s short form from `-a` to `-d`. Update init behavior to include an already-initialized guard. All changes to `main.go`, `init.go`, and `init_test.go`.

#### 1.0 Proof Artifact(s)

- CLI: `gws --help` shows `--init` as a boolean flag with no path argument and `--add-tag` using `-d` demonstrates flag signature changes
- CLI: `gws --init` run in a directory containing git repos outputs `Initialized workspace at [path]. Found N repositories.` demonstrates successful boolean-flag init with repo discovery
- CLI: `gws --init` run a second time outputs a notification with the existing workspace path, a suggestion to use `--add`, and a suggestion to use `--refresh` demonstrates the already-initialized guard
- Test: `go test ./cmd/git-workspace/... -run TestInit` passes demonstrates behavioral correctness

#### 1.0 Tasks

- [x] 1.1 In `main.go`, change the `flagInit` variable declaration from `string` to `bool` in the `var` block.
- [x] 1.\1 In `main.go init()`, replace `rootCmd.Flags().StringVarP(&flagInit, "init", "i", "", ...)` with `rootCmd.Flags().BoolVarP(&flagInit, "init", "i", false, "Initialize a gws workspace in the current directory")`.
- [x] 1.\1 In `main.go init()`, change `--add-tag` registration from short form `"a"` to `"d"`: `rootCmd.Flags().BoolVarP(&flagAddTag, "add-tag", "d", false, ...)`.
- [x] 1.\1 In `main.go` `RunE`, update the mutual exclusivity check from `if flagInit != ""` to `if flagInit` (two occurrences: the counter increment and the dispatch call).
- [x] 1.\1 In `main.go` `RunE`, update the "workspace not initialized" error message to reference `gws --init` (without a path argument) to match the new flag signature.
- [x] 1.\1 In `main.go init()`, update the `commandFlagUsages` template group list to reflect the `-d` short form for `--add-tag` (no code change needed — Cobra reads this from the flag definition — but update the Long description string to remove the `--init ~/projects` example and replace it with `gws --init` showing no path).
- [x] 1.\1 In `init.go`, remove the `workspacePath := flagInit` line and the `filepath.Abs(workspacePath)` call. Replace them with `absPath, err := os.Getwd()` (add `"os"` to imports; remove `"path/filepath"` if no longer used).
- [x] 1.\1 In `init.go` `runInit`, add an already-initialized guard at the top of the function: call `config.Exists()` and, if the workspace already exists, load it with `config.Load()`, print a notification to stderr that includes the existing workspace path (`cfg.Workspace`), suggest `gws --add` to add repositories, and suggest `gws --refresh` to re-scan. Return `nil` after printing (exit 0).
- [x] 1.\1 In `init.go`, update the success output messages to match the spec format: `Initialized workspace at /path. Found N repositories.` Replace the current multi-line output block with this concise format.
- [x] 1.10 Update `init_test.go` to cover: (a) happy path — `runInit` creates a workspace in a temp dir containing a git repo and prints the expected confirmation, (b) already-initialized guard — running `runInit` a second time prints the notification and returns nil without overwriting config, (c) empty directory — `runInit` in a directory with no git repos initializes successfully with count of 0.
- [x] 1.11 Update `main_test.go` to reflect the `--init` flag type change: any test that passes a path string to `--init` should be updated to pass the boolean flag only.
- [x] 1.12 Run `make ci` and confirm all tests pass before marking this task complete.

---

### [x] 2.0 Implement `--add` Command — Single Repository with Symlink Support

Create `add.go` with the `runAdd` handler. Register `--add` / `-a` in `main.go` using Cobra's `NoOptDefVal` so the flag works both with and without a path argument. Implement git repo validation, config update, and symlink creation for repositories located outside the workspace root. Wire `--add` into the mutual exclusivity check and dispatch in `main.go`.

#### 2.0 Proof Artifact(s)

- CLI: `gws --add` run inside a git repository outputs `Added [repo-name] to workspace.` demonstrates adding a repo from the current directory
- CLI: `gws --add /path/to/external/repo` outputs `Added [repo-name] to workspace.` followed by `Created symlink: [workspace]/[repo-name] → /path/to/external/repo` demonstrates external repo add with symlink creation
- CLI: `gws --add` run in a non-git directory outputs an error message and exits non-zero demonstrates the git validation check
- CLI: Running `gws --add` twice in the same git repo outputs the already-tracked warning on the second run demonstrates duplicate detection
- Filesystem: `ls -la [workspace-path]` shows a symlink entry for the external repo demonstrates the symlink exists on disk
- Test: `go test ./cmd/git-workspace/... -run TestAdd` passes demonstrates behavioral correctness

#### 2.0 Tasks

- [ ] 2.1 In `main.go`, add a new `flagAdd string` variable to the command flags `var` block.
- [ ] 2.2 In `main.go init()`, register the `--add` flag: `rootCmd.Flags().StringVarP(&flagAdd, "add", "a", "", "Add a git repository to the workspace (defaults to current directory)")`. Then immediately after, set `rootCmd.Flags().Lookup("add").NoOptDefVal = "."` so that `gws --add` without a value defaults to the current directory.
- [ ] 2.3 In `main.go init()`, add `"add"` to the `commandFlagUsages` template group list so it appears in the Commands section of `--help`.
- [ ] 2.4 In `main.go` `RunE`, add `flagAdd` to the mutual exclusivity counter: `if flagAdd != "" { activeCount++ }`. Update the error message string to include `--add`.
- [ ] 2.5 In `main.go` `RunE`, add an `--add` dispatch call after the `--init` dispatch: `if flagAdd != "" { return runAdd(cmd, args) }`. Place this before the workspace existence check so `--add` can give its own clear error when no workspace exists.
- [ ] 2.6 Create `cmd/git-workspace/add.go` with `package main`. Define `func runAdd(_ *cobra.Command, _ []string) error`.
- [ ] 2.7 In `runAdd`, check whether a workspace is initialized using `config.Exists()`. If not, print an error message directing the user to run `gws --init` first and return an error.
- [ ] 2.8 In `runAdd`, resolve the target path: use `flagAdd` as the input. If `flagAdd` is `"."`, replace it with `os.Getwd()`. Convert to absolute path using `filepath.Abs()`.
- [ ] 2.9 In `runAdd`, validate the target is a git repository by checking that a `.git` entry exists inside the resolved path using `os.Stat(filepath.Join(absPath, ".git"))`. If the stat fails or the path is not a directory, print `"[path] is not a git repository"` to stderr and return an error.
- [ ] 2.10 In `runAdd`, load the config with `config.Load()`. Check whether the target path is already tracked by comparing `absPath` against each `repo.Path` in `cfg.Repositories`. If found, print `"[repo-name] is already tracked, skipping."` and return nil.
- [ ] 2.11 In `runAdd`, extract the repository metadata for the new repo: open the git repo using `git.PlainOpen(absPath)` from `go-git`, retrieve the `origin` remote URL, and use `classifier.Classify()` to get the repo type and visibility. Construct a `config.Repository` struct from this data (set `Name` to `filepath.Base(absPath)`, `Path` to `absPath`).
- [ ] 2.12 In `runAdd`, determine whether the target repo is inside the workspace directory using `strings.HasPrefix(absPath, cfg.Workspace)`. If it is outside the workspace, create a symlink: compute `symlinkPath := filepath.Join(cfg.Workspace, repoName)`, then call `os.Symlink(absPath, symlinkPath)`.
- [ ] 2.13 Before creating the symlink, check if a file or symlink already exists at `symlinkPath` using `os.Lstat()`. If it exists, print a warning (`"Symlink path [symlinkPath] already exists, skipping symlink creation."`) and continue without returning an error.
- [ ] 2.14 Append the new `config.Repository` to `cfg.Repositories` and call `config.Save(cfg)`. Print `"Added [repo-name] to workspace."`. If a symlink was created, also print `"Created symlink: [symlinkPath] → [absPath]"`.
- [ ] 2.15 Create `cmd/git-workspace/add_test.go`. Write table-driven tests covering: (a) adding the current directory (a valid git repo), (b) adding an explicit external path (valid git repo outside workspace — verify symlink is created in workspace), (c) adding a non-git directory (expect error), (d) adding an already-tracked repo (expect warning, no error), (e) symlink path conflict (a file already exists at the symlink destination — expect warning, no error).
- [ ] 2.16 Run `make ci` and confirm all tests pass before marking this task complete.

---

### [x] 3.0 Implement `--add --recursive` / `-v` — Batch Repository Add

Register `--recursive` / `-v` as a boolean modifier flag in `main.go`. Add validation that `--recursive` requires `--add`. Extend `runAdd` in `add.go` to handle recursive scanning using the existing `discovery.Scan()` function, batch config updates, symlink creation for external repos, and count + names reporting.

#### 3.0 Proof Artifact(s)

- CLI: `gws --add --recursive` in a directory containing multiple git repos outputs `Added N repositories: repo-a, repo-b, repo-c.` demonstrates batch add with count and names
- CLI: `gws -a -v` produces the same result as above demonstrates short-form flag equivalence
- CLI: `gws --add --recursive` when all repos are already tracked outputs `No new repositories found.` demonstrates idempotent behavior
- CLI: `gws --recursive` without `--add` returns an error message and exits non-zero demonstrates modifier flag validation
- Test: `go test ./cmd/git-workspace/... -run TestAddRecursive` passes demonstrates behavioral correctness

#### 3.0 Tasks

- [x] 3.1 In `main.go`, add a new `flagRecursive bool` variable to the command flags `var` block.
- [x] 3.2 In `main.go init()`, register the flag: `rootCmd.Flags().BoolVarP(&flagRecursive, "recursive", "v", false, "Recursively add all git repositories found in the current directory (use with --add)")`.
- [x] 3.3 In `main.go` `RunE`, add a validation block for `--recursive` after the filter flags validation block: `if flagRecursive && flagAdd == "" { return fmt.Errorf("--recursive/-v requires --add/-a to be set") }`. This mirrors the pattern used to gate filter flags on `--list`.
- [x] 3.4 In `add.go` `runAdd`, add a branch at the top of the function: `if flagRecursive { return runAddRecursive() }`. Implement `runAddRecursive` as a separate function in the same file.
- [x] 3.5 In `runAddRecursive`, resolve the scan root: use `os.Getwd()` as the directory to scan (recursive always operates on the current directory).
- [x] 3.6 In `runAddRecursive`, call `discovery.Scan(scanRoot)` to find all git repositories in the current directory (same scanner used by `--init`).
- [x] 3.7 In `runAddRecursive`, load the existing config with `config.Load()`. For each discovered repository, check if it is already tracked (compare `repo.Path` against existing entries). Collect a list of already-tracked repo names to warn about.
- [x] 3.8 In `runAddRecursive`, for each newly discovered repository (not already tracked): append it to `cfg.Repositories`, and if its path is outside the workspace directory, create a symlink using the same `os.Symlink` logic from Task 2 (extract this into a shared helper function `createSymlinkIfExternal(cfg, repo)` in `add.go`).
- [x] 3.9 In `runAddRecursive`, after processing all repos: if already-tracked repos were encountered, print a warning for each one (`"[name] is already tracked, skipping."`). If new repos were added, call `config.Save(cfg)` and print `"Added N repositories: repo-a, repo-b, repo-c."`. If no new repos were added, print `"No new repositories found."`.
- [x] 3.10 Refactor Task 2's symlink creation logic in `runAdd` into the shared `createSymlinkIfExternal(cfg *config.Config, repoPath, repoName string) (created bool, err error)` helper so both single and recursive paths use the same code.
- [x] 3.11 Update `add_test.go` with recursive-specific tests: (a) recursive add in a directory with 3 git repos — verify all 3 are added and output includes count + names, (b) recursive add when all repos already tracked — verify `"No new repositories found."` output, (c) recursive add in a directory with no git repos — verify `"No new repositories found."` output, (d) `gws --recursive` without `--add` — verify error returned from `main.go` validation.
- [x] 3.12 Run `make ci` and confirm all tests pass before marking this task complete.

---

### [x] 4.0 Add Shell Tab Completion

Register Cobra's built-in `completion` subcommand on the root command. Register a `ValidArgsFunction` on the `--add` flag that returns directory path completions. Update the usage template to surface the `completion` subcommand. Manually verify completion works in at least one shell (zsh or bash).

#### 4.0 Proof Artifact(s)

- CLI: `gws completion zsh` outputs a valid zsh completion script demonstrates zsh script generation
- CLI: `gws completion bash` outputs a valid bash completion script demonstrates bash script generation
- CLI: `gws completion --help` shows available shells (bash, zsh, fish, powershell) demonstrates subcommand is fully registered
- Shell demo: After running `source <(gws completion zsh)`, typing `gws --add ~/pro` and pressing `<Tab>` auto-completes to matching directories demonstrates live path completion in the shell
- Build: `make ci` passes demonstrates no regressions introduced

#### 4.0 Tasks

- [x] 4.1 In `main.go init()`, after all flags are registered, add: `rootCmd.InitDefaultCompletionCmd()`. This registers Cobra's built-in `completion` subcommand (supporting bash, zsh, fish, and powershell) on the root command.
- [x] 4.2 In `main.go init()`, register a completion function for the `--add` flag value by calling: `rootCmd.RegisterFlagCompletionFunc("add", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { return nil, cobra.ShellCompDirectiveFilterDirs })`. This tells the shell completion system to suggest only directories when completing the `--add` flag value.
- [x] 4.3 In `main.go`, update the custom `SetUsageTemplate` string to add a `completion` line in the Commands section so users can discover it via `gws --help`. Add a line like `  gws completion [bash|zsh|fish|powershell]  # Generate shell completion script` to the Long description string.
- [x] 4.4 Build the binary (`make build` or `go build ./cmd/git-workspace/`) and run `gws completion zsh` to verify it produces output. Run `gws completion bash` as well. Both should print a shell script to stdout.
- [x] 4.5 Manually test completion in your shell: run `source <(gws completion zsh)` (or the bash equivalent), then type `gws --add ~/` and press `<Tab>`. Verify that directory names are suggested.
- [x] 4.6 Run `make ci` and confirm all tests pass before marking this task complete.
