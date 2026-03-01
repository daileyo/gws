# 08 Tasks - User Update/Delete

## Notes

- **Delete shorthand**: `-d` is already used by `--add-tag`, so `--delete` uses `-D` as its shorthand.
- `-u` is available and will be used for `--update`.

## Relevant Files

- `cmd/git-workspace/main.go` - Root command flag registration, validation logic, dispatch, and help template grouping
- `cmd/git-workspace/userupdate.go` - **New file.** Handler for `--user --update/-u` operations (single and batch)
- `cmd/git-workspace/userdelete.go` - **New file.** Handler for `--user --delete/-D` operations (single and batch)
- `cmd/git-workspace/userupdate_test.go` - **New file.** Tests for user update handler
- `cmd/git-workspace/userdelete_test.go` - **New file.** Tests for user delete handler
- `cmd/git-workspace/tag.go` - Contains `findRepositories()` helper used by update/delete to match repos
- `cmd/git-workspace/user.go` - Existing `gws user` subcommand tree; reference for profile lookup patterns
- `cmd/git-workspace/userdetect.go` - `detectUserForRepos()` used as reference for config.json update logic
- `internal/user/assign.go` - Contains `AssignLocal()` (reused by update); will add `DeleteLocal()` for delete
- `internal/user/assign_test.go` - Tests for `AssignLocal()` and new `DeleteLocal()`
- `internal/user/profile.go` - `GetProfile()`, `GetAllProfiles()`, `ValidateEmail()` used for profile resolution
- `internal/filter/filter.go` - `filter.Apply()` with `Criteria{Tags: ...}` used for batch tag filtering
- `internal/git/user.go` - `GetUserConfig()`, `GetNonLocalUserConfig()` used for re-detection after delete
- `internal/config/config.go` - `Repository` struct, `UserSource` constants, `Load()`/`Save()`

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `userupdate.go` and `userupdate_test.go` in the same directory).
- Use `make test` to run all tests; `go test ./cmd/git-workspace/ -run TestName` for individual tests.
- Follow the repository's table-driven test pattern with `t.Run()` subtests and `t.TempDir()` for temp git repos.
- Follow Conventional Commits format for commit messages (`feat:`, `fix:`, etc.).
- Use `strings.EqualFold()` for case-insensitive name/profile comparisons.

## Tasks

### [x] 1.0 CLI Flag Infrastructure and Dispatch

Register the new `--user`, `--update/-u`, `--delete/-D`, `--all`, `--verbose` flags on the root command. Add validation logic (mutual exclusivity, flag dependencies) and dispatch skeleton. Update the help template to display the new flags in a grouped section. Extend `--quiet` to work with `--user` operations (it already exists but is restricted to navigation).

#### 1.0 Proof Artifact(s)

- CLI: `gws --help` output demonstrates new flags appear in a grouped section (e.g., "User Operations")
- CLI: `gws --user --update` without required args demonstrates validation error message
- CLI: `gws --user --update --delete` demonstrates mutual exclusivity error
- CLI: `gws --update` without `--user` demonstrates flag dependency error

#### 1.0 Tasks

- [x] 1.1 Add new flag variables to the command flags section in `main.go`: `flagUser` (bool), `flagUpdate` (bool, shorthand `-u`), `flagDelete` (bool, shorthand `-D`), `flagAll` (bool), `flagVerbose` (bool). Also add `flagInlineName` (string) and `flagInlineEmail` (string) for inline value support.
- [x] 1.2 Register the new flags in the `init()` function in `main.go` using `rootCmd.Flags().BoolVarP()` and `rootCmd.Flags().StringVar()`. Group them logically near the existing command flags.
- [x] 1.3 Update the mutual exclusivity validation in `RunE` to include `flagUser` in the `activeCount` check (so `--user` cannot be combined with `--list`, `--init`, `--add`, etc.).
- [x] 1.4 Add flag dependency validation: `--update` and `--delete` require `--user`; `--update` and `--delete` are mutually exclusive; `--all` requires `--delete`; `--name`/`--email` inline flags require `--update`; `--tag` is allowed with `--user` (remove the `--list`-only restriction for `--tag`).
- [x] 1.5 Extend the `--quiet` validation to also allow `--quiet` with `--user` operations (currently only allowed with navigation).
- [x] 1.6 Add dispatch in `RunE`: when `flagUser` is true, call `runUserUpdate(cmd, args)` if `--update` or `runUserDelete(cmd, args)` if `--delete`. Create stub functions in `userupdate.go` and `userdelete.go` that return `fmt.Errorf("not yet implemented")`.
- [x] 1.7 Update the help template in `init()`: add a new "User Operations" section using a new `cobra.AddTemplateFunc("userFlagUsages", flagGroup([]string{"user", "update", "delete", "all", "verbose", "name", "email"}))`. Add the section to `SetUsageTemplate` between "Commands" and "List Filters".
- [x] 1.8 Write tests in `cmd/git-workspace/main_test.go` (or a new test file) that verify: `--update` without `--user` returns error; `--delete` without `--user` returns error; `--user --update --delete` returns error; `--all` without `--delete` returns error; `--user` is mutually exclusive with other command flags.

### [~] 2.0 Single Repository User Update (`--user -u`)

Implement the `runUserUpdate()` handler that sets local git user config on one or more matching repositories. Support both named profiles and inline `--name`/`--email` values (with inline taking precedence). Display moderate output by default, with `--verbose` for full detail and `--quiet` for silent operation. Update `config.json` after successful operations.

#### 2.0 Proof Artifact(s)

- CLI: `gws --user -u <repo> <profile>` output demonstrates single repo update with moderate summary (repo name + changes)
- CLI: `git config --local --list` in the updated repo demonstrates actual `.git/config` changes (user.name, user.email set)
- CLI: `gws -l --show-user` demonstrates the repo reflects the new user info with "(local)" indicator
- CLI: `gws --user -u <repo> --name "Inline Name" --email "inline@example.com"` demonstrates inline value support

#### 2.0 Tasks

- [x] 2.1 Implement `resolveProfile()` in `userupdate.go`: given args and inline flags, resolve the target profile. If a profile name is provided as a positional arg, look it up via `user.GetProfile()` then fall back to auto-detected profiles via `user.DetectProfiles()`. If inline `--name`/`--email` are provided, create an ad-hoc `config.Profile` from them. If both profile name and inline values are provided, start from the named profile and override with inline values. Return error if neither profile name nor inline email is provided.
- [x] 2.2 Implement `runUserUpdate()` in `userupdate.go`: validate args (need at least a repo identifier unless `--tag` is used — tag-based batch is Task 4.0), call `findRepositories()` to match repos, call `resolveProfile()` to get the target profile, then for each matched repo call `user.AssignLocal(repo.Path, profile)` and update the repo's config.json fields (User, Email, SigningEnabled, UserSource = "local").
- [x] 2.3 Implement output formatting in `runUserUpdate()`: before calling `AssignLocal`, capture current config via `git.GetUserConfig()`. After update, compare old vs new. In default (moderate) mode, print `<repo-name>: user.name "old" → "new", user.email "old" → "new"`. In `--verbose` mode, print full before/after including signing config. In `--quiet` mode, print nothing. When multiple repos are updated, print a summary count at the end (e.g., `Updated 3 repositories`).
- [x] 2.4 Add `config.Save(cfg)` call after all repos are updated to persist config.json changes.
- [x] 2.5 Write tests in `userupdate_test.go`: test profile resolution (named profile, inline values, profile + inline override); test single repo update writes to `.git/config`; test multiple repo match updates all; test `--quiet` suppresses output; test error when no profile or inline values provided; test error when no repos match identifier.

### [ ] 3.0 Single Repository User Delete (`--user --delete`)

Implement `DeleteLocal()` in `internal/user/assign.go` to remove user config keys from `.git/config`. Implement the `runUserDelete()` handler. By default, remove only `[user]` name and email; with `--all`, also remove signing config. After deletion, re-detect the effective user config and update `config.json` so the repo shows its fallback (global/includeIf) user.

#### 3.0 Proof Artifact(s)

- CLI: `gws --user --delete <repo>` output demonstrates local config removal with moderate summary
- CLI: `git config --local --list` confirms `[user]` name/email removed from `.git/config`
- CLI: `gws -l --show-user` demonstrates the repo now shows the default user (no "(local)" indicator)
- CLI: `gws --user --delete <repo> --all` demonstrates signing config is also removed

#### 3.0 Tasks

- [ ] 3.1 Implement `DeleteLocal(repoPath string, removeAll bool) error` in `internal/user/assign.go`. Read the repo's `.git/config` file, remove `name` and `email` keys from the `[user]` section. If `removeAll` is true, also remove `signingkey` from `[user]` and `gpgsign` from `[commit]`. If the `[user]` or `[commit]` section becomes empty after key removal, remove the entire section header. Write the modified config back to the file.
- [ ] 3.2 Implement a helper `removeGitConfigKey(content, section, key string) string` in `internal/user/assign.go` (companion to existing `setGitConfigValue`). It should find the `[section]` and remove the line matching the key. If the section has no remaining keys, remove the section header line too.
- [ ] 3.3 Write tests for `DeleteLocal()` in `internal/user/assign_test.go`: test removing only name/email (default); test removing name/email/signingkey/gpgsign (`removeAll=true`); test that non-user sections are preserved; test that empty sections are cleaned up; test error when path is not a git repo.
- [ ] 3.4 Implement `runUserDelete()` in `userdelete.go`: validate args (need a repo identifier unless `--tag` is used — tag-based batch is Task 4.0), call `findRepositories()`, for each matched repo capture current config via `git.GetUserConfig()`, call `user.DeleteLocal(repo.Path, flagAll)`, then re-detect effective config via `git.GetUserConfig()` (which will now return global/includeIf) and update the repo's config.json fields.
- [ ] 3.5 Implement output formatting in `runUserDelete()`: in moderate mode, print `<repo-name>: removed local user config (now using <source>: "name" <email>)`. In `--verbose` mode, print full details of removed config and new effective config. In `--quiet` mode, print nothing. Print summary count at end for multiple repos.
- [ ] 3.6 Add `config.Save(cfg)` call after all repos are processed.
- [ ] 3.7 Write tests in `userdelete_test.go`: test single repo delete removes local config; test `--all` removes signing config; test config.json is updated with fallback user; test `--quiet` suppresses output; test error when no repos match.

### [ ] 4.0 Batch Update and Delete via Tags (`--tag`)

Extend `--tag` flag validation to allow its use with `--user` operations (currently restricted to `--list`). Implement batch update (`gws --user -u --tag <tag> <profile>`) and batch delete (`gws --user --delete --tag <tag>`) that apply to all repos matching the tag filter. Display per-repo summary with `--verbose`/`--quiet` support.

#### 4.0 Proof Artifact(s)

- CLI: `gws --user -u --tag <tag> <profile>` output demonstrates batch update with summary listing each affected repo
- CLI: `gws --user --delete --tag <tag>` output demonstrates batch delete with summary
- CLI: `gws -l --show-user --tag <tag>` demonstrates all tagged repos reflect the changes after update
- CLI: `gws -l --show-user --tag <tag>` demonstrates all tagged repos show default user after delete

#### 4.0 Tasks

- [ ] 4.1 Update the `hasFilterFlags()` function in `main.go` to NOT include `--tag` when `--user` is active. The `--tag` flag should be allowed with both `--list` and `--user` operations. Update the validation error message to reflect this (e.g., "filter flags (--type, --name, --path, --output, --status) require --list/-l").
- [ ] 4.2 Update `runUserUpdate()` to handle tag-based batch mode: when `filterTags` is non-empty, load all repos from config, apply `filter.Apply(repos, filter.Criteria{Tags: filterTags})` to get matching repos, then process all matching repos. The repo identifier positional arg is not required when `--tag` is used. Validate that either a repo identifier or `--tag` is provided (not neither).
- [ ] 4.3 Update `runUserDelete()` with the same tag-based batch logic: when `filterTags` is non-empty, filter repos by tag and process all matches. Validate that either a repo identifier or `--tag` is provided.
- [ ] 4.4 Handle the edge case where `--tag` returns zero matching repos: print a clear message like `No repositories found with tag(s): <tags>` and return without error.
- [ ] 4.5 Write tests for batch operations: test `--user -u --tag <tag> <profile>` updates all tagged repos; test `--user --delete --tag <tag>` deletes from all tagged repos; test multiple `--tag` values use AND logic; test no repos match tag; test `--tag` with `--quiet` suppresses output.
