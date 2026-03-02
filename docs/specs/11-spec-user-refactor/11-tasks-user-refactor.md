# 11 Tasks - User Refactor

## Relevant Files

- `cmd/git-workspace/user.go` - User subcommand with sub-subcommands (`list`, `add`, `show`, `remove`, `assign`, `sync`). Short-flag aliases and `RunE` routing are added here.
- `cmd/git-workspace/user_test.go` - New test file for user subcommand short-flag registration, mutual exclusivity, dispatch routing, and tab completion.
- `cmd/git-workspace/main.go` - Root command. User flag variables and registrations are removed from here; root `RunE` user routing is removed; usage template is updated.
- `cmd/git-workspace/main_test.go` - Root command tests. User flag validation tests are migrated to deprecated tests; user flag registration tests updated.
- `cmd/git-workspace/deprecated.go` - Deprecation layer. User flag variables, registrations, warnings, and dispatch logic are added here.
- `cmd/git-workspace/deprecated_test.go` - Deprecation tests. Tests for deprecated user flag dispatch, warnings, and hidden flags added here.
- `cmd/git-workspace/userupdate.go` - Contains `runListUsers()`, `resolveProfile()`, `runUserUpdate()`. Business logic unchanged; `filterTags` references remain (variable stays in `main.go`).
- `cmd/git-workspace/userdelete.go` - Contains `runUserDelete()`. Business logic unchanged; `filterTags` references remain.
- `cmd/git-workspace/userupdate_test.go` - Existing user update tests. May need minor adjustments if deprecated flag variables are renamed.
- `cmd/git-workspace/userdelete_test.go` - Existing user delete tests. May need minor adjustments if deprecated flag variables are renamed.

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `user.go` and `user_test.go` in the same directory).
- Use `go vet ./... && go test ./... -count=1` for validation.
- Follow the repository's existing code organization, naming conventions, and style guidelines (Cobra/pflag patterns, `fmt.Errorf("...: %w", err)` wrapping, table-driven tests).
- Follow the tag subcommand pattern from spec 10 (`tag.go`) for short-flag registration and `RunE` dispatch.
- `filterTags` is shared between `list.go` (list subcommand `--tag`/`-t`) and `deprecated.go` (root `--tag` for backward compat). It stays declared in `main.go` as a package-level variable.

## Tasks

### [x] 1.0 Add short-flag aliases (-a, -d, -l, -s) to the user subcommand

Register `-a` (add), `-d` (remove), `-l` (list), and `-s` (show) as boolean flags on `userCmd`. Add routing logic in `userCmd.RunE` to dispatch to the corresponding sub-subcommand logic based on which flag is set. Enforce mutual exclusivity between the four flags, and between the flags and word-form sub-subcommands. Ensure existing word-form sub-subcommands (`add`, `remove`, `list`, `show`, `assign`, `sync`) continue to work unchanged. When `gws user` is invoked with no flag and no sub-subcommand, display help text.

#### 1.0 Proof Artifact(s)

- CLI: `gws user -a test-profile --email test@example.com` adds a profile demonstrates short-flag add works
- CLI: `gws user -l` lists all profiles demonstrates short-flag list works
- CLI: `gws user -s test-profile` shows profile details demonstrates short-flag show works
- CLI: `gws user -d test-profile` removes a profile demonstrates short-flag remove works
- CLI: `gws user add` still works demonstrates word-form preserved
- CLI: `gws user` shows help demonstrates default behavior
- Test: Short-flag routing tests pass demonstrates mutual exclusivity and dispatch

#### 1.0 Tasks

- [x] 1.1 Add short-flag variables (`userFlagAdd`, `userFlagDelete`, `userFlagList`, `userFlagShow`) as bool vars in `user.go`, following the pattern from `tag.go` (`tagFlagAdd`, `tagFlagDelete`).
- [x] 1.2 Register the four flags on `userCmd` in `user.go`'s `init()`: `-a`/`--add`, `-d`/`--delete`, `-l`/`--list`, `-s`/`--show`. Register them with `BoolVarP` on `userCmd.Flags()`.
- [x] 1.3 Add a `RunE` function to `userCmd` that:
  - Checks mutual exclusivity among `-a`, `-d`, `-l`, `-s` (only one may be set)
  - If `-l` is set: delegates to `userListCmd.RunE(cmd, args)`
  - If `-a` is set: validates `len(args) >= 1` (profile name), then calls `userAddCmd.RunE(cmd, args)` (the add command's RunE handles `--email` validation)
  - If `-s` is set: validates `len(args) == 1` (profile name), then calls `userShowCmd.RunE(cmd, args)`
  - If `-d` is set: validates `len(args) == 1` (profile name), then calls `userRemoveCmd.RunE(cmd, args)`
  - If no flag and no sub-subcommand: calls `cmd.Help()`
- [x] 1.4 Register the `--email`, `--name`, `--signing-key`, `--sign-commits` flags on `userCmd` (in addition to `userAddCmd`) so they are available when using `gws user -a <name> --email <email>`. Use the same backing variables (`addEmail`, `addGitName`, `addSigningKey`, `addSignCommit`).
- [x] 1.5 Update the `userCmd.Long` help text to document the short-flag aliases alongside the word-form sub-subcommands.
- [x] 1.6 Reset `userCmd`'s usage template to Cobra's default (same pattern as `tagCmd` in `tag.go:101-124`) so it doesn't inherit root's custom template.
- [x] 1.7 Run `go vet ./... && go test ./... -count=1` to verify no regressions. Fix any issues.

### [x] 2.0 Add short-flag tests and tab completion for user short flags

Add comprehensive tests for the new short-flag interface: flag registration, mutual exclusivity, dispatch routing, and coexistence with word-form sub-subcommands. Add tab completion functions for the positional arguments used with `-a`, `-d`, and `-s` (profile name completion).

#### 2.0 Proof Artifact(s)

- Test: `TestUserShortFlags*` tests pass demonstrates all short-flag routing and validation
- Test: `TestUserShortFlagMutualExclusivity` passes demonstrates flags are mutually exclusive
- CLI: Tab completion for `gws user -s <TAB>` returns profile names demonstrates completion works
- Test: Full test suite (`go test ./... -count=1`) passes demonstrates no regressions

#### 2.0 Tasks

- [x] 2.1 Create `cmd/git-workspace/user_test.go` with tests for:
  - `TestUserShortFlagsRegistered`: Verify `-a`, `-d`, `-l`, `-s` flags are registered on `userCmd` with correct shorthand letters.
  - `TestUserShortFlagMutualExclusivity`: Verify setting two or more short flags returns an error.
  - `TestUserSubcommandsRegistered`: Verify `list`, `add`, `show`, `remove`, `assign`, `sync` are registered as sub-subcommands of `userCmd`.
  - `TestUserNoOperationShowsHelp`: Verify `gws user` with no flags and no sub-subcommand returns help (no error).
- [x] 2.2 Add a `completeProfileNames` helper function in `user.go` that returns stored and auto-detected profile names matching a prefix (similar to `completeRepoNames` in `tag.go`).
- [x] 2.3 Set `userCmd.ValidArgsFunction` to complete profile names for the first positional argument (used with `-a`, `-s`, `-d`).
- [x] 2.4 Add a test `TestUserProfileCompletion` in `user_test.go` that verifies the completion function returns profile names.
- [x] 2.5 Run `go vet ./... && go test ./... -count=1` to verify all tests pass with no regressions.

### [x] 3.0 Deprecate root-level user flags to deprecated.go

Move the root-level user flag variables (`flagUser`, `flagUpdate`, `flagDelete`, `flagAll`, `flagListUsers`, `flagVerbose`, `flagInlineName`, `flagInlineEmail`) and their registrations from `main.go` to `deprecated.go`. Register them as hidden flags with deprecation warnings. Add dispatch logic in `handleDeprecatedFlags()` for the compound flag combinations (`--user` alone → `user list`, `--user --update` → `runUserUpdate`, `--user --delete` → `runUserDelete`, `--list-users` → `runListUsers`). Remove the user flag routing from root `RunE` and the user flag section from the root usage template.

#### 3.0 Proof Artifact(s)

- CLI: `gws --user` lists profiles plus prints deprecation warning demonstrates backward compat
- CLI: `gws --list-users` lists profiles plus prints deprecation warning demonstrates backward compat
- CLI: `gws --help` does NOT show deprecated user flags demonstrates clean help output
- Test: Deprecated user flag dispatch tests pass demonstrates deprecation layer works
- Test: Full test suite passes demonstrates no regressions

#### 3.0 Tasks

- [x] 3.1 Move user flag variables from `main.go` to `deprecated.go`: rename `flagUser` → `depUser`, `flagUpdate` → `depUpdate`, `flagDelete` → `depDelete`, `flagAll` → `depAll`, `flagListUsers` → `depListUsers`, `flagVerbose` → `depVerbose`, `flagInlineName` → `depInlineName`, `flagInlineEmail` → `depInlineEmail`. Add them to the deprecated var block with a comment section for spec 11.
- [x] 3.2 Move the user flag registrations from `main.go`'s `init()` to `registerDeprecatedFlags()` in `deprecated.go`. Register them using the same flag names and shorthands. Add them to the `hiddenFlags` slice so they are hidden from `--help`.
- [x] 3.3 Add deprecation warning mappings to the `depWarnings` map for all user flags: `"user"` → `"gws user list"`, `"update"` → `"gws user assign <repo> <profile>"`, `"delete"` → `"gws user assign (then remove local config)"`, `"all"` → `"gws user assign (with --all)"`, `"list-users"` → `"gws user list"`, `"git-name"` → `"gws user add --name"`, `"git-email"` → `"gws user add --email"`, `"verbose"` → `"gws user --verbose"`.
- [x] 3.4 Add user flag dispatch logic to `handleDeprecatedFlags()`:
  - Include `depUser || depListUsers` in the mutual exclusivity count (replace the existing `flagUser || flagListUsers` reference).
  - Add user flag dependency validation (same checks currently in root `RunE`): `--update` requires `--user`, `--delete` requires `--user`, `--update`/`--delete` mutually exclusive, `--all` requires `--delete`, `--git-name`/`--git-email` require `--user --update`, `--verbose` requires `--user`, `--quiet` with `--user` validation.
  - Dispatch: `depListUsers` → `emitDeprecationWarnings` + `runListUsers(cmd, args)`. `depUser` alone → `emitDeprecationWarnings` + `runListUsers(cmd, args)`. `depUser + depUpdate` → `emitDeprecationWarnings` + `runUserUpdate(cmd, args)`. `depUser + depDelete` → `emitDeprecationWarnings` + `runUserDelete(cmd, args)`.
- [x] 3.5 Remove user flag routing from root `RunE` in `main.go`: remove the `flagUpdate && !flagUser` check, `flagDelete && !flagUser` check, `flagUpdate && flagDelete` check, `flagAll && !flagDelete` check, inline name/email checks, `flagVerbose` check, `flagQuiet` user check, the `flagTagAlias` block, `flagListUsers` dispatch, and `flagUser` dispatch block. Keep navigation and default workspace info logic.
- [x] 3.6 Remove the user flag group template functions from `main.go`'s `init()` (`userFlagUsages`). Update the root usage template to remove the "User Operations" section. Keep the "Navigation" and "Other" sections.
- [x] 3.7 Update `main_test.go`: move user flag validation tests (`TestUserFlagValidation`, `TestUserFlagsRegistered`) to `deprecated_test.go` with updated variable names. Update `TestAliasFlagsAreHidden` to include the newly hidden user flags. Update `TestTagRequiresListOrUser` to use `depUser` instead of `flagUser`.
- [x] 3.8 Update `userupdate.go` and `userdelete.go` to reference the renamed variables: `flagQuiet` → `depQuiet` (if moved) or keep `flagQuiet` since it's still a root navigation flag. Update references to `flagVerbose` → `depVerbose`, `flagAll` → `depAll`, `flagInlineName` → `depInlineName`, `flagInlineEmail` → `depInlineEmail`.
- [x] 3.9 Update `userupdate_test.go` and `userdelete_test.go` to use the renamed deprecated variable names.
- [x] 3.10 Add deprecated user flag tests to `deprecated_test.go`: `TestDeprecatedUserFlagsAreHidden`, `TestDeprecatedUserEmitsWarning`, `TestDeprecatedUserListUsersDispatch`.
- [x] 3.11 Run `go vet ./... && go test ./... -count=1` to verify all tests pass. Fix any issues.

### [x] 4.0 Clean up shared state and update help text

Update `filterTags` comment in `main.go` to reflect that it's shared between `list.go` and `deprecated.go` (no longer used by user ops on root). Remove the `flagTagAlias` routing from root `RunE` if not already removed. Verify `assign` and `sync` work unchanged. Run final quality gates.

#### 4.0 Proof Artifact(s)

- CLI: `gws user --help` shows short-flag documentation demonstrates updated help
- CLI: `gws user assign my-repo dev-profile` still works demonstrates assign unchanged
- CLI: `gws user sync` still works demonstrates sync unchanged
- CLI: `gws --help` shows clean output with no user flags demonstrates cleanup
- Test: Full test suite (`go vet ./... && go test ./... -count=1`) passes demonstrates no regressions

#### 4.0 Tasks

- [x] 4.1 Update the `filterTags` comment in `main.go` to reflect its current usage: shared between `list.go` (list subcommand `--tag`/`-t`) and `deprecated.go` (root `--tag` for backward compat with `--list` and deprecated user operations).
- [x] 4.2 Verify the `--tag` flag on root (registered in `deprecated.go`) still works correctly for `gws --list --tag work` (deprecated list with tag filter) and `gws --user -u --tag work dev-profile` (deprecated user update with tag batch).
- [x] 4.3 Verify `gws user assign <repo> <profile>` and `gws user sync` work unchanged (no regressions from the refactor).
- [x] 4.4 Verify `gws user --help` displays the short-flag aliases (`-a`, `-d`, `-l`, `-s`) alongside the word-form sub-subcommands.
- [x] 4.5 Verify `gws --help` shows a clean output with no user flags visible (all hidden in deprecated layer).
- [x] 4.6 Run `go vet ./... && go test ./... -count=1` as final quality gate. Verify all tests pass across all packages.
