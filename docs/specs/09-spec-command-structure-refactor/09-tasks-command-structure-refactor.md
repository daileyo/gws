# 09 Tasks - Command Structure Refactor

## Relevant Files

- `cmd/git-workspace/main.go` - Root command definition, flag registration, RunE dispatch, ValidArgsFunction, usage template. Major changes: remove command flags, simplify RunE to handle only navigation/default/alias dispatch, update usage template.
- `cmd/git-workspace/main_test.go` - Tests for root command flags, mutual exclusivity, validation. Must be rewritten for new subcommand structure.
- `cmd/git-workspace/list.go` - List logic and display functions. Must be refactored to accept filter params instead of reading globals, then wrapped in a `listCmd` Cobra subcommand.
- `cmd/git-workspace/list_test.go` - New test file for the `list` subcommand and its scoped filter flags.
- `cmd/git-workspace/init.go` - Init logic. Must be wrapped in an `initCmd` Cobra subcommand.
- `cmd/git-workspace/init_test.go` - Tests for init. Shared test helpers (`createInitTestRepo`, `withTempHome`, `withTempWorkdir`) live here — these are reused across test files.
- `cmd/git-workspace/add.go` - Add logic. Must be refactored to accept path param instead of reading `flagAdd` global, then wrapped in an `addCmd` Cobra subcommand.
- `cmd/git-workspace/add_test.go` - Tests for add. Must be updated to use new subcommand invocation.
- `cmd/git-workspace/refresh.go` - Refresh logic. Must be wrapped in a `refreshCmd` Cobra subcommand.
- `cmd/git-workspace/navigate.go` - Navigation logic. Already clean (takes params). No changes needed to the function itself.
- `cmd/git-workspace/navigate_test.go` - Navigation tests. No changes needed (tests call `runNavigate` directly with params).
- `cmd/git-workspace/shellinit.go` - Shell integration templates. Must update case statement to route subcommand names to binary.
- `cmd/git-workspace/shellinit_test.go` - New or updated test file for shell template assertions.
- `cmd/git-workspace/deprecated.go` - New file for the deprecation layer: hidden flags, warning messages, delegation.
- `cmd/git-workspace/deprecated_test.go` - New test file for deprecation layer.
- `cmd/git-workspace/tag.go` - Tag logic. Not modified in this spec, but tag flags (`--add-tag`, `--remove-tag`) are removed from root and NOT re-registered (deferred to spec 10).
- `cmd/git-workspace/tag_test.go` - Tag tests. May need minor updates if tag flag removal affects test setup.
- `cmd/git-workspace/user.go` - User subcommand tree. Not modified in this spec, but root-level `--user`/`--update`/`--delete` flags are removed and NOT re-registered (deferred to spec 11).
- `cmd/git-workspace/userupdate.go` - User update/delete logic. Not modified but root flag references removed.
- `cmd/git-workspace/userdelete.go` - User delete logic. Not modified but root flag references removed.
- `cmd/git-workspace/userupdate_test.go` - User update tests. May need minor updates.
- `cmd/git-workspace/userdelete_test.go` - User delete tests. May need minor updates.

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `list.go` and `list_test.go` in the same directory).
- Use `make test` for running tests, `make lint` for linting, `make ci` for full CI validation.
- Follow existing Go patterns: `fmt.Errorf("...: %w", err)` for error wrapping, table-driven tests, `t.TempDir()`/`t.Setenv()` for isolation.
- All test files use `package main` (not `_test` package) for access to unexported functions.
- The `user` subcommand in `user.go` is the reference pattern for how Cobra subcommands are structured in this project.

## Tasks

### [x] 1.0 Extract Business Logic from Root RunE into Standalone Functions

#### 1.0 Proof Artifact(s)

- Test: `make test` passes with zero changes to external behavior demonstrates refactor is safe
- Test: `make lint` passes demonstrates code quality maintained
- CLI: `gws --list`, `gws --init`, `gws --add .`, `gws --refresh`, `gws --print-workspace` all produce identical output to before demonstrates no regressions

#### 1.0 Tasks

- [x] 1.1 Refactor `runList` to accept a `ListOptions` struct parameter (containing `filterType`, `filterTags`, `filterName`, `filterPath`, `outputFormat`, `showStatus`, `showUser`) instead of reading package-level globals. Update the call site in `main.go` RunE to build the struct from globals and pass it in.
- [x] 1.2 Refactor `runAdd` to accept a path string and recursive bool as parameters instead of reading `flagAdd` and `flagRecursive` globals. Update the call site in `main.go` RunE to pass the current global values.
- [x] 1.3 Refactor `runInit` to accept an optional directory path parameter (empty string means current directory) instead of always using `os.Getwd()`. Update the call site in `main.go` RunE.
- [x] 1.4 Refactor `runRefresh` to require no global state (it currently has none beyond config loading — verify and confirm it's already clean).
- [x] 1.5 Extract the inline `print-workspace` logic from RunE into a `runPrintWorkspace()` function in a new or existing file, so it can be called from both the subcommand and deprecation layer.
- [x] 1.6 Update all existing tests (`init_test.go`, `add_test.go`, `main_test.go`) to use the new function signatures. Verify `make test` and `make lint` pass.

### [x] 2.0 Create Core Subcommands (list, init, add, refresh, print-workspace)

#### 2.0 Proof Artifact(s)

- CLI: `gws list` produces the same output as current `gws --list` demonstrates subcommand works
- CLI: `gws -l` produces the same output as `gws list` demonstrates root alias works
- CLI: `gws init` in a new directory creates workspace demonstrates init subcommand works
- CLI: `gws add` in a git repo adds it to config demonstrates add subcommand works
- CLI: `gws refresh` refreshes metadata demonstrates refresh subcommand works
- CLI: `gws print-workspace` prints workspace path demonstrates print-workspace subcommand works
- Test: `make ci` passes demonstrates no regressions

#### 2.0 Tasks

- [x] 2.1 Create the `listCmd` Cobra subcommand in `list.go`. Register it with `rootCmd.AddCommand(listCmd)` in an `init()` function. The `RunE` should build a `ListOptions` struct from its own flags and call the refactored `runList`. Register filter flags on `listCmd` (this is completed in task 3.0 — for now, register them on both root and list to avoid breaking anything).
- [x] 2.2 Create the `initCmd` Cobra subcommand in `init.go`. It should accept an optional positional arg for directory path (default `.`). Register with `rootCmd.AddCommand(initCmd)` in `init()`. The `RunE` calls the refactored `runInit` with the path argument.
- [x] 2.3 Create the `addCmd` Cobra subcommand in `add.go`. It should accept an optional positional arg for repo path (default `.`). Register `--recursive`/`-v` flag on `addCmd`. Register with `rootCmd.AddCommand(addCmd)`. The `RunE` calls the refactored `runAdd` with the path and recursive flag.
- [x] 2.4 Create the `refreshCmd` Cobra subcommand in `refresh.go`. No arguments. Register with `rootCmd.AddCommand(refreshCmd)`. The `RunE` calls `runRefresh`.
- [x] 2.5 Create the `printWorkspaceCmd` Cobra subcommand (command name `print-workspace`) in a new file or in `main.go`. No arguments. Register with `rootCmd.AddCommand`. The `RunE` calls `runPrintWorkspace`.
- [x] 2.6 Register root-level hidden boolean alias flags for each subcommand: `-l` (list), `-i` (init), `-a` (add), `-r` (refresh), `-w` (print-workspace). In the root `RunE`, check if any alias flag is set and delegate to the corresponding subcommand's logic. Enforce mutual exclusivity (only one alias flag at a time, and no alias flag if a subcommand was invoked).
- [x] 2.7 Update the root command's `Long` description and usage template to reflect the new subcommand structure. Replace the flag-based examples with subcommand examples (e.g., `gws list` instead of `gws --list`).
- [x] 2.8 Update `TestNoSubcommandsRegistered` in `main_test.go` — invert it to `TestSubcommandsRegistered` that verifies `list`, `init`, `add`, `refresh`, `print-workspace` ARE registered. Update `TestCommandFlagsRegistered` to verify the alias flags (`-l`, `-i`, `-a`, `-r`, `-w`) are registered and hidden.
- [x] 2.9 Write tests for each new subcommand: verify `listCmd`, `initCmd`, `addCmd`, `refreshCmd`, `printWorkspaceCmd` can be invoked and produce correct results. Follow existing test patterns (save/restore globals, `withTempHome`, `withTempWorkdir`).
- [x] 2.10 Verify `make ci` passes (vet, lint, test-race).

### [x] 3.0 Scope Filter Flags to the List Subcommand

#### 3.0 Proof Artifact(s)

- CLI: `gws list -t work -y github -su` produces filtered output with status and user columns demonstrates scoped filters and stacking
- CLI: `gws list --tag "wo*" --status` produces wildcard-filtered output demonstrates long-form filters
- CLI: `gws list -n "my*" -o json` produces JSON output demonstrates all filter types work
- Test: `make ci` passes demonstrates no regressions

#### 3.0 Tasks

- [x] 3.1 Remove filter flag registrations (`--type`, `--tag`, `--name`, `--path`, `--output`, `--status`, `--show-user`) from the root command's `init()` in `main.go`. Also remove the corresponding package-level variables (`filterType`, `filterTags`, etc.) from `main.go`.
- [x] 3.2 Register filter flags on the `listCmd` in `list.go`'s `init()`: `--tag`/`-t` (string slice), `--type`/`-y` (string), `--name`/`-n` (string), `--path`/`-p` (string), `--output`/`-o` (string, default "table"), `--status`/`-s` (bool), `--show-user`/`-u` (bool). Define the corresponding variables as package-level vars in `list.go` (or as struct fields).
- [x] 3.3 Remove the `hasNonTagFilterFlags` helper and the filter validation logic from root `RunE` (no longer needed since filters are scoped to `list`).
- [x] 3.4 Remove `TestFilterFlagsRegistered` and `TestFilterFlagsRequireList` from `main_test.go`. Add equivalent tests in `list_test.go` that verify filter flags are registered on `listCmd` with correct names and shorthands.
- [x] 3.5 Verify POSIX flag stacking works: write a test that invokes `listCmd` with `-su` (stacked booleans) and confirms both `showStatus` and `showUser` are set to true.
- [x] 3.6 Remove the `filterFlagUsages` template function from `main.go`'s usage template (filters now appear in `gws list --help`, not `gws --help`). Update the root usage template to remove the "List Filters" section.
- [x] 3.7 Verify `make ci` passes.

### [x] 4.0 Update Shell Integration and Tab Completion

#### 4.0 Proof Artifact(s)

- CLI: `gws my-repo` navigates to the repo (via shell function `cd`) demonstrates positional navigation unchanged
- CLI: Tab completion after `gws ` shows both commands (`list`, `init`, etc.) and repo names demonstrates dual completion
- CLI: `gws shell-init zsh` outputs updated function with subcommand routing demonstrates template updated
- Test: Shell init tests pass with updated template assertions demonstrates correctness

#### 4.0 Tasks

- [x] 4.1 Update the `zshInitTemplate` in `shellinit.go`: change the `case` pattern to route known subcommand names to the binary. Replace `-*|__*|completion|shell-init|help)` with a pattern that includes all subcommands: `list|init|add|refresh|print-workspace|tag|user|completion|shell-init|help|-*)`. The `*` catch-all continues to handle repo name navigation via `cd`.
- [x] 4.2 Update the `bashInitTemplate` in `shellinit.go` with the same subcommand routing pattern as zsh.
- [x] 4.3 Verify the root command's `ValidArgsFunction` still returns repo names. With subcommands registered, Cobra automatically includes subcommand names in completions. Test that both subcommand names and repo names appear in completion output by running `git-workspace __complete ""` (Cobra's internal completion protocol).
- [x] 4.4 Write or update tests for `shellinit.go` that assert the shell templates contain the expected subcommand names in the routing pattern. Verify both zsh and bash templates are correct.
- [x] 4.5 Verify `make ci` passes.

### [ ] 5.0 Add Deprecation Layer for Old Flag Forms

#### 5.0 Proof Artifact(s)

- CLI: `gws --list` produces correct output plus deprecation warning to stderr demonstrates backward compatibility
- CLI: `gws --list --type github` produces filtered output plus deprecation warnings demonstrates filter compat
- CLI: `gws --help` does NOT show `--list`, `--init`, or other deprecated flags demonstrates clean help
- Test: Deprecated flags route to correct logic and emit warnings demonstrates deprecation layer works
- Test: `make ci` passes demonstrates full suite green

#### 5.0 Tasks

- [ ] 5.1 Create `deprecated.go` with a `registerDeprecatedFlags(rootCmd)` function. This function registers hidden flags for: `--list` (bool), `--init` (bool), `--add` (string with `NoOptDefVal = "."`), `--refresh` (bool), `--print-workspace` (bool), `--go` (string). Each flag should use `cmd.Flags().MarkHidden()` after registration.
- [ ] 5.2 In `deprecated.go`, register hidden filter flags on root: `--type` (string), `--tag` (string slice), `--name` (string), `--path` (string), `--output` (string), `--status` (bool), `--show-user` (bool). These are needed so `gws --list --type github` continues to work during the deprecation period.
- [ ] 5.3 Create a `handleDeprecatedFlags(cmd, args)` function in `deprecated.go` that checks if any deprecated flag is set (via `cmd.Flags().Changed()`). For each deprecated flag that is set: (a) print a warning to stderr (`fmt.Fprintf(os.Stderr, "Warning: --%s is deprecated, use '%s' instead\n", flagName, newForm)`), and (b) delegate to the corresponding new subcommand logic (building params from the deprecated flag values).
- [ ] 5.4 Call `registerDeprecatedFlags(rootCmd)` from `main.go`'s `init()`. In root `RunE`, call `handleDeprecatedFlags` early in the dispatch chain (before the default/navigation fallthrough) to intercept deprecated flag usage.
- [ ] 5.5 Handle the compound deprecation case: `gws --list --type github --tag work`. The deprecation handler must detect `--list` plus any deprecated filter flags, build a `ListOptions` from them, and call the list logic with those options. Each deprecated flag used should emit its own warning.
- [ ] 5.6 Ensure all deprecated flags are hidden: verify they do NOT appear in `gws --help` output. Write a test that parses help output and confirms none of the deprecated flag names are present.
- [ ] 5.7 Create `deprecated_test.go` with tests: (a) each deprecated flag routes to correct logic, (b) each deprecated flag emits a warning to stderr, (c) compound `--list` + filter flags work correctly, (d) deprecated flags are hidden from help.
- [ ] 5.8 Remove the old command flag variables (`flagList`, `flagInit`, `flagAdd`, `flagRefresh`, `flagPrintWorkspace`) and their registrations from `main.go` — these are now fully replaced by subcommands + deprecation layer. Also remove the old filter flag variables from `main.go` if not already removed in task 3.1.
- [ ] 5.9 Update or remove any remaining tests in `main_test.go` that reference the old flag variables. The `TestMutualExclusivity` test should be rewritten to test the new alias flags and/or subcommand mutual exclusivity.
- [ ] 5.10 Note: Do NOT register deprecated forms for `--add-tag`, `--remove-tag`, `--user`, `--update`, `--delete`, `--all`, `--verbose`, `--git-name`, `--git-email`, `--list-users`. Those are handled in specs 10 and 11. For now, these flags remain on root as-is (they will be migrated in later specs).
- [ ] 5.11 Verify `make ci` passes with the complete refactor in place.
