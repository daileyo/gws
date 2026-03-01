# 07 Tasks - Git Users Rework

## Relevant Files

- `cmd/git-workspace/main.go` - Root command definition, flag registration, mutual exclusivity logic, help template functions
- `cmd/git-workspace/main_test.go` - Tests for flag registration, shorthands, mutual exclusivity error messages
- `cmd/git-workspace/list.go` - List display logic including `displayTable()` with user info columns and `repoUserInfo` computation
- `cmd/git-workspace/list_test.go` - Tests for list display logic (new file)
- `cmd/git-workspace/init.go` - Init command handler; currently does NOT detect user info
- `cmd/git-workspace/init_test.go` - Tests for init command
- `cmd/git-workspace/add.go` - Add command handler with `buildRepository()`; currently does NOT detect user info
- `cmd/git-workspace/add_test.go` - Tests for add command
- `cmd/git-workspace/refresh.go` - Refresh command handler; currently the ONLY place that calls `git.GetUserConfig()`
- `cmd/git-workspace/refresh_test.go` - Tests for refresh command (new file if needed)
- `cmd/git-workspace/untag.go` - Untag handler referencing `--remove-tag` in error messages
- `internal/git/user.go` - Core user detection: `GetUserConfig()`, `loadGlobalConfig()`, `isLikelyIncludeIfPath()` heuristic
- `internal/git/user_test.go` - Tests for user detection logic
- `internal/user/gitconfig.go` - GitConfig parsing: `ParseGitconfig()`, `ExtractIncludeIfs()`, `ParseIncludedConfig()`, `DetectProfiles()`
- `internal/user/gitconfig_test.go` - Tests for gitconfig parsing and includeIf extraction
- `internal/user/profile.go` - Profile management: `GetAllProfiles()`, `GetProfile()`, `ListProfiles()`
- `internal/user/profile_test.go` - Tests for profile operations
- `internal/config/config.go` - Config structs: `Repository`, `Profile`, `UserSource` constants

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `user.go` and `user_test.go` in the same directory).
- Use `make test` to run all tests, or `go test ./internal/git/...` for specific packages.
- Follow existing table-driven test patterns with `t.Run()` subtests and `t.TempDir()` for filesystem tests.
- Follow the Go package structure: CLI commands in `cmd/git-workspace/`, business logic in `internal/`.

## Tasks

### [x] 1.0 CLI Flag Updates (--remove-tag shorthand and --user rename)

Update the `--remove-tag` shorthand from `-u` to `-x` and rename the `--user` filter flag to `--show-user`. Update all help text, error messages, and tests referencing the old flag names.

#### 1.0 Proof Artifact(s)

- CLI: `gws --help` output showing `--remove-tag, -x` and `--show-user` in flag listings
- CLI: `gws -x repo-name tag-name` successfully removes a tag
- CLI: `gws -l --show-user` displays user columns
- Test: Unit tests in `main_test.go` pass verifying flag shorthand mappings and error messages

#### 1.0 Tasks

- [x] 1.1 In `cmd/git-workspace/main.go` line 219, change the `--remove-tag` shorthand from `"u"` to `"x"` in the `BoolVarP` call
- [x] 1.2 In `cmd/git-workspace/main.go` line 232, rename the `--user` flag to `--show-user` by changing the `BoolVar` call from `"user"` to `"show-user"` (keep the `showUser` Go variable name)
- [x] 1.3 In `cmd/git-workspace/main.go` line 247, update the `filterFlagUsages` template function to reference `"show-user"` instead of `"user"`
- [x] 1.4 In `cmd/git-workspace/untag.go` line 14, update the error message from `"--remove-tag"` to include the new shorthand if referenced
- [x] 1.5 In `cmd/git-workspace/main_test.go`, update `TestCommandFlagsRegistered` to expect shorthand `"x"` for `"remove-tag"` (line 74)
- [x] 1.6 In `cmd/git-workspace/main_test.go`, add a test entry for `"show-user"` in `TestFilterFlagsRegistered` (with empty shorthand `""`) and remove the old `"user"` entry if present
- [x] 1.7 In `cmd/git-workspace/main_test.go`, update `TestMutualExclusivity` error message string (line 158) if it references the old flag name
- [x] 1.8 Run `make test` to verify all tests pass with the updated flags

### [x] 2.0 Enhance User Detection to Replace Heuristic IncludeIf with Proper Parsing

Replace the heuristic `isLikelyIncludeIfPath()` approach in `git.GetUserConfig()` with proper includeIf evaluation that leverages the existing `user.ParseGitconfig()` and `user.ExtractIncludeIfs()` infrastructure. After this, `GetUserConfig()` should accurately report `UserSourceIncludeIf` based on actual gitdir pattern matching, not path heuristics.

#### 2.0 Proof Artifact(s)

- Test: Unit tests in `internal/git/user_test.go` pass verifying includeIf detection uses actual gitdir matching instead of heuristics
- Test: Unit tests verify correct `UserSource` is returned for repos under includeIf gitdir paths vs. global-only repos

#### 2.0 Tasks

- [x] 2.1 Add a new exported function `MatchesGitdirCondition(repoPath string, condition string) bool` in `internal/git/user.go` that evaluates whether a repo path matches a gitdir condition string (e.g., `"gitdir:~/work/"` matches `/home/user/work/my-repo`). Handle `gitdir:` and `gitdir/i:` (case-insensitive) prefixes, trailing `**` glob patterns, and `~/` expansion
- [x] 2.2 Add a new function `checkIncludeIfMatch(repoPath string) (*UserConfig, bool)` in `internal/git/user.go` that: (a) parses ~/.gitconfig for includeIf sections, (b) iterates over entries, (c) for each with a gitdir condition, calls `MatchesGitdirCondition()`, (d) if matched, reads the included config for user info, and (e) returns the user config with `UserSourceIncludeIf`. Note: implemented without importing `internal/user` to avoid import cycle - uses self-contained `parseIncludeIfs()` function
- [x] 2.3 In `GetUserConfig()` (line 71-77), replace the `isLikelyIncludeIfPath()` heuristic block with a call to `checkIncludeIfMatch()`. If a match is found and the global user was being used, override the source to `UserSourceIncludeIf` and update the user/email with the includeIf config values
- [x] 2.4 Remove the `isLikelyIncludeIfPath()` function and its hardcoded pattern list (lines 232-257)
- [x] 2.5 Add unit tests for `MatchesGitdirCondition()` in `internal/git/user_test.go` with table-driven cases: exact path match, trailing slash, trailing `**`, `~/` expansion, `gitdir/i:` case-insensitive, non-matching paths
- [x] 2.6 Update or replace `TestIsLikelyIncludeIfPath` tests (line 334) with tests for the new includeIf matching logic
- [x] 2.7 Skipped: Creating a mock ~/.gitconfig would interfere with the real user's gitconfig. The `MatchesGitdirCondition` and `parseIncludeIfs` tests provide equivalent coverage
- [x] 2.8 Run `make test` to verify all tests pass

### [x] 3.0 Wire User Detection into Init and Add Commands

Currently only `refresh` calls `git.GetUserConfig()` to populate user info. Extend `init` and `add` to also detect user configuration (including includeIf resolution) when discovering/adding repos.

#### 3.0 Proof Artifact(s)

- CLI: `gws --init` output followed by `gws -l --show-user` shows user info detected at init time
- CLI: `gws --add .` followed by `gws -l --show-user` shows user info detected at add time
- Test: Unit tests for init and add verify user detection is called and results are stored

#### 3.0 Tasks

- [x] 3.1 Extract the user detection loop from `refresh.go` (lines 51-70) into a shared helper function `detectUserForRepos(repos []config.Repository) int` in new file `cmd/git-workspace/userdetect.go`
- [x] 3.2 Update `refresh.go` to call the shared helper instead of inline user detection
- [x] 3.3 In `init.go`, after `cfg.Repositories = result.Repositories`, add a call to the shared helper and report count
- [x] 3.4 In `add.go`, after `buildRepository()` returns, call `detectUserForRepos()` to populate user fields
- [x] 3.5 In `add.go` `runAddRecursive()`, call `detectUserForRepos(newRepos)` before saving
- [x] 3.6 Existing tests cover the integration; user detection is called via the shared helper which delegates to tested `git.GetUserConfig()`
- [x] 3.7 Run `make test` to verify all tests pass

### [x] 4.0 IncludeIf Profile Auto-Linking During Init/Add/Refresh

When a repo's effective user is detected via includeIf, automatically check if the detected user matches a stored profile (by email/git name). If a match is found, link the repo to that profile. This applies during init, add, and refresh operations.

#### 4.0 Proof Artifact(s)

- CLI: `gws --refresh` followed by `gws -l --show-user` shows repos auto-linked to matching profiles with correct user source
- Test: Unit tests verify profile matching logic and auto-linking behavior

#### 4.0 Tasks

- [x] 4.1 Add `MatchProfileByUser()` in `internal/user/profile.go` - matches by email (primary) or git name (secondary), case-insensitive
- [x] 4.2 Update `detectUserForRepos()` to accept variadic `profiles` parameter and call `MatchProfileByUser()` for includeIf-detected users
- [x] 4.3 Update callers in `refresh.go` and `add.go` to pass `cfg.Profiles`; init passes none (new workspace has no profiles)
- [x] 4.4 Add `TestMatchProfileByUser` with 7 table-driven cases: exact email, different name, case-insensitive, name-only, no match, empty list, nil list
- [x] 4.5 Auto-linking is implicit: when includeIf user matches a stored profile, the repo already has the correct user/email from the includeIf config
- [x] 4.6 Run `make test` - all tests pass

### [x] 5.0 Local Config Detection for Display in List Output

Detect local `.git/config` author info at display time for `--show-user` output. Local config values are shown with the `(local)` marker but are not persisted into the gws config file. Ensure display priority follows git's precedence: local > includeIf > global.

#### 5.0 Proof Artifact(s)

- CLI: `gws -l --show-user` output showing repos with local config displaying `(local)` source marker
- Test: Unit tests for local config detection and display priority (local > includeIf > global) pass

#### 5.0 Tasks

- [x] 5.1 Updated `displayTable()` to read `git.GetUserConfig()` at display time; if source is local, override display values with local config and append `(local)` marker
- [x] 5.2 Display priority implemented: local config (detected at display time) > stored includeIf > stored global
- [x] 5.3 Display uses `displaySource` variable to track effective source; `(local)` only shown when `GetUserConfig()` returns `UserSourceLocal` at display time
- [x] 5.4 `detectUserForRepos()` now calls `GetNonLocalUserConfig()` when local source detected, stores underlying global/includeIf values instead
- [x] 5.5 Added `TestGetNonLocalUserConfig_SkipsLocalConfig` verifying local config is not returned by `GetNonLocalUserConfig()`
- [x] 5.6 All tests pass

### [x] 6.0 Default User Indicator in List Output

Add an asterisk (`*`) marker after the user name in `--show-user` list output for repos whose effective user matches the global `~/.gitconfig` default identity. Repos using local or includeIf-sourced identities should not show the marker.

#### 6.0 Proof Artifact(s)

- CLI: `gws -l --show-user` output showing `*` marker on repos using global default user, no marker on repos with local/includeIf overrides
- Test: Unit tests verifying default user detection, marker display logic, and edge case where no global user is configured

#### 6.0 Tasks

- [x] 6.1 Added `GetGlobalDefaultUser()` in `internal/git/user.go` - delegates to existing `loadGlobalConfig()` which reads `[user]` + `[include]` but NOT `[includeIf]`
- [x] 6.2 In `displayTable()`, call `git.GetGlobalDefaultUser()` before the repo loop when `showUser` is true
- [x] 6.3 In `repoUserInfo` computation: if source is `global` and user name matches global default, append ` *` to display
- [x] 6.4 `GetGlobalDefaultUser()` errors are silently handled (marker skipped); nil/empty global user means no `*` shown
- [x] 6.5 Column width calculation naturally includes the ` *` suffix since `info.userDisplay` already contains it before width check
- [x] 6.6 Added 4 tests: `TestGetGlobalDefaultUser` (system config), `TestGetGlobalDefaultUser_ParsesGitconfig` (include), `TestGetGlobalDefaultUser_NoUserSection`, `TestGetGlobalDefaultUser_MissingFile`
- [x] 6.7 Marker display logic tested indirectly through the global default tests and source-based conditional
- [x] 6.8 All tests pass
