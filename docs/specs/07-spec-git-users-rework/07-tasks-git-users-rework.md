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

### [ ] 3.0 Wire User Detection into Init and Add Commands

Currently only `refresh` calls `git.GetUserConfig()` to populate user info. Extend `init` and `add` to also detect user configuration (including includeIf resolution) when discovering/adding repos.

#### 3.0 Proof Artifact(s)

- CLI: `gws --init` output followed by `gws -l --show-user` shows user info detected at init time
- CLI: `gws --add .` followed by `gws -l --show-user` shows user info detected at add time
- Test: Unit tests for init and add verify user detection is called and results are stored

#### 3.0 Tasks

- [ ] 3.1 Extract the user detection loop from `refresh.go` (lines 51-70) into a shared helper function `detectUserForRepos(repos []config.Repository) int` that iterates over repos, calls `git.GetUserConfig()` for each, populates User/Email/SigningEnabled/UserSource fields, and returns the count of repos with detected user info. Place this in a new file `cmd/git-workspace/userdetect.go` or in an existing shared location
- [ ] 3.2 Update `refresh.go` to call the shared helper instead of inline user detection
- [ ] 3.3 In `init.go`, after `cfg.Repositories = result.Repositories` (line 59), add a call to the shared helper to detect user config for all discovered repos. Add an output line reporting the count (e.g., `"Repositories with user configuration: %d"`)
- [ ] 3.4 In `add.go`, after `buildRepository()` returns (line 74), call `git.GetUserConfig(absPath)` and populate the repo's User/Email/SigningEnabled/UserSource fields before appending to config
- [ ] 3.5 In `add.go` `runAddRecursive()`, after discovering new repos (line 153), iterate over `newRepos` and call `git.GetUserConfig()` for each to populate user fields before saving
- [ ] 3.6 Add unit tests verifying that repos added via `init` and `add` have user info populated (User, Email, UserSource fields are non-empty when a git user config exists)
- [ ] 3.7 Run `make test` to verify all tests pass

### [ ] 4.0 IncludeIf Profile Auto-Linking During Init/Add/Refresh

When a repo's effective user is detected via includeIf, automatically check if the detected user matches a stored profile (by email/git name). If a match is found, link the repo to that profile. This applies during init, add, and refresh operations.

#### 4.0 Proof Artifact(s)

- CLI: `gws --refresh` followed by `gws -l --show-user` shows repos auto-linked to matching profiles with correct user source
- Test: Unit tests verify profile matching logic and auto-linking behavior

#### 4.0 Tasks

- [ ] 4.1 Add a new function `MatchProfileByUser(profiles []config.Profile, userName string, email string) *config.Profile` in `internal/user/profile.go` that finds a stored profile matching the given email (primary) or git name (secondary). Return `nil` if no match
- [ ] 4.2 Update the shared user detection helper (from task 3.1) to accept a `profiles []config.Profile` parameter. After detecting a user with `UserSourceIncludeIf`, call `MatchProfileByUser()` to check for a matching stored profile. If matched, the repo's user info is confirmed as linked to that profile
- [ ] 4.3 Update the callers in `init.go`, `add.go`, and `refresh.go` to load profiles from config and pass them to the user detection helper
- [ ] 4.4 Add unit tests for `MatchProfileByUser()` in `internal/user/profile_test.go` with cases: exact email match, email match with different name, no match, empty profiles list
- [ ] 4.5 Add unit tests verifying that during refresh, a repo under an includeIf gitdir path is auto-linked to a matching stored profile
- [ ] 4.6 Run `make test` to verify all tests pass

### [ ] 5.0 Local Config Detection for Display in List Output

Detect local `.git/config` author info at display time for `--show-user` output. Local config values are shown with the `(local)` marker but are not persisted into the gws config file. Ensure display priority follows git's precedence: local > includeIf > global.

#### 5.0 Proof Artifact(s)

- CLI: `gws -l --show-user` output showing repos with local config displaying `(local)` source marker
- Test: Unit tests for local config detection and display priority (local > includeIf > global) pass

#### 5.0 Tasks

- [ ] 5.1 In `cmd/git-workspace/list.go` `displayTable()`, within the `showUser` block (lines 122-161), update the `repoUserInfo` computation to read the repo's local `.git/config` at display time using `git.GetUserConfig()`. If local user.name/email is found, override the stored values and mark `userDisplay` with `(local)` suffix
- [ ] 5.2 Ensure the display priority is correct: if local config exists, use it and show `(local)`; else if stored source is `includeif`, show the stored values (no `(local)` marker); else show stored global values
- [ ] 5.3 Update the `repoUserInfo` struct in `list.go` to track whether the displayed user comes from local config vs. stored config, so the `(local)` marker is only shown when the effective user differs from the stored user or comes from a local `.git/config`
- [ ] 5.4 In `refresh.go` (and the shared user detection helper), when `GetUserConfig()` returns `UserSourceLocal`, do NOT store the local user values in the config. Instead, store the underlying global/includeIf values so that the config reflects the "assigned" identity, not the local override
- [ ] 5.5 Add unit tests for the display priority logic: test that a repo with local config shows `(local)`, a repo with includeIf shows the includeIf user, and a repo with only global config shows the global user
- [ ] 5.6 Run `make test` to verify all tests pass

### [ ] 6.0 Default User Indicator in List Output

Add an asterisk (`*`) marker after the user name in `--show-user` list output for repos whose effective user matches the global `~/.gitconfig` default identity. Repos using local or includeIf-sourced identities should not show the marker.

#### 6.0 Proof Artifact(s)

- CLI: `gws -l --show-user` output showing `*` marker on repos using global default user, no marker on repos with local/includeIf overrides
- Test: Unit tests verifying default user detection, marker display logic, and edge case where no global user is configured

#### 6.0 Tasks

- [ ] 6.1 Add a new exported function `GetGlobalDefaultUser() (*GlobalUserConfig, error)` in `internal/git/user.go` that reads only the `[user]` section from `~/.gitconfig` (plus `[include]` directives but NOT `[includeIf]`) to determine the default global user identity. This leverages the existing `loadGlobalConfig()` function
- [ ] 6.2 In `cmd/git-workspace/list.go` `displayTable()`, before the repo loop, call `git.GetGlobalDefaultUser()` to get the default user name and email
- [ ] 6.3 In the `repoUserInfo` computation block, after determining the displayed user: if the effective user's source is `global` (i.e., not local, not includeIf) and the user name/email matches the global default, append ` *` to `userDisplay` (e.g., `"John Doe *"`)
- [ ] 6.4 Handle the edge case where `GetGlobalDefaultUser()` returns an error or empty values: skip the marker entirely (no `*` shown for any repo)
- [ ] 6.5 Account for the `*` marker in column width calculation: include the ` *` suffix length when computing `maxUserLen`
- [ ] 6.6 Add unit tests for `GetGlobalDefaultUser()` in `internal/git/user_test.go` with cases: gitconfig with [user] section, gitconfig with [include] that sets user, gitconfig with no [user] section, missing gitconfig file
- [ ] 6.7 Add unit tests for the marker display logic: repo with global user gets `*`, repo with local user does not, repo with includeIf user does not, no global user configured means no markers
- [ ] 6.8 Run `make test` to verify all tests pass
