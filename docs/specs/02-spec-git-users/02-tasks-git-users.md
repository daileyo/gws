# 02-tasks-git-users.md

## Relevant Files

### Config Package
- `internal/config/config.go` - Extend with Profile struct and Repository user fields (User, Email, SigningEnabled, UserSource)
- `internal/config/config_test.go` - Tests for config serialization with new fields

### Git Package
- `internal/git/user.go` - NEW: Functions to detect effective git user config from repositories
- `internal/git/user_test.go` - NEW: Tests for git user detection logic

### User Package (NEW)
- `internal/user/profile.go` - NEW: Profile struct and profile management functions
- `internal/user/profile_test.go` - NEW: Tests for profile management logic
- `internal/user/gitconfig.go` - NEW: Parse ~/.gitconfig and includeIf directives to detect existing profiles
- `internal/user/gitconfig_test.go` - NEW: Tests for gitconfig parsing
- `internal/user/assign.go` - NEW: Repository user assignment logic (local config and subdirectory modes)
- `internal/user/assign_test.go` - NEW: Tests for assignment logic

### Commands
- `cmd/gws/user.go` - NEW: User parent command with list, add, show, remove, assign, sync subcommands
- `cmd/gws/list.go` - Add --user flag and USER, EMAIL, SIGN columns
- `cmd/gws/refresh.go` - Update to populate user fields during refresh

### Documentation
- `README.md` - Add user management section with examples

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `user.go` and `user_test.go` in the same directory)
- Use `go test ./... -v` to run all tests, or `go test ./internal/user/... -v` for specific package
- Follow existing table-driven test patterns (see `internal/git/status_test.go`, `internal/filter/filter_test.go`)
- Use `t.TempDir()` for creating temporary git repositories in tests
- Follow existing command patterns from `cmd/gws/tag.go` for new subcommands
- Use wrapped errors with context: `fmt.Errorf("failed to X: %w", err)`

## Tasks

### [ ] 1.0 Extend Config Model for User Profiles and Repository User Data

#### 1.0 Proof Artifact(s)

- File: `internal/config/config.go` contains Profile struct and extended Repository struct demonstrates model is defined
- Test: `go test ./internal/config/... -v` passes demonstrates model serialization works
- File: `~/.gws/config.json` can store and load profiles and per-repo user data demonstrates persistence works

#### 1.0 Tasks

- [ ] 1.1 Add Profile struct to `internal/config/config.go` with fields: Name (string), GitName (string), Email (string), SigningKey (string), SignCommits (bool)
- [ ] 1.2 Add Profiles slice field to Config struct: `Profiles []Profile`
- [ ] 1.3 Add user fields to Repository struct: User (string), Email (string), SigningEnabled (bool), UserSource (string)
- [ ] 1.4 Add UserSource constants: UserSourceGlobal, UserSourceLocal, UserSourceIncludeIf, UserSourceUnknown
- [ ] 1.5 Update ConfigVersion constant to "1.1.0" to indicate schema change
- [ ] 1.6 Write unit tests in `internal/config/config_test.go` for Profile serialization (JSON marshal/unmarshal)
- [ ] 1.7 Write unit tests for Repository with user fields serialization
- [ ] 1.8 Write unit tests for Config with Profiles slice serialization
- [ ] 1.9 Verify backward compatibility: loading config without profiles/user fields should work (empty defaults)

### [ ] 2.0 Implement Git User Detection from Repositories

#### 2.0 Proof Artifact(s)

- CLI: `gws refresh` updates user.name, user.email, and signing config for all repositories demonstrates detection works
- Test: `go test ./internal/git/... -v` passes demonstrates user detection logic works
- File: `~/.gws/config.json` contains populated User, Email, SigningEnabled, UserSource fields after refresh demonstrates storage works

#### 2.0 Tasks

- [ ] 2.1 Create `internal/git/user.go` with UserConfig struct containing Name, Email, SigningKey, SignCommits, Source fields
- [ ] 2.2 Implement `GetUserConfig(repoPath string) (*UserConfig, error)` function using go-git to read effective config
- [ ] 2.3 Implement config cascade detection: determine if user config comes from global, local .git/config, or includeIf
- [ ] 2.4 Handle edge cases: repo with no user configured (returns global default), repo with no commits yet
- [ ] 2.5 Create `internal/git/user_test.go` with test helper to create git repos with specific user configs
- [ ] 2.6 Write unit tests for GetUserConfig with global-only user config
- [ ] 2.7 Write unit tests for GetUserConfig with local .git/config override
- [ ] 2.8 Write unit tests for GetUserConfig with signing configuration
- [ ] 2.9 Update `cmd/gws/refresh.go` to call GetUserConfig for each repository during refresh
- [ ] 2.10 Store detected user info in Repository struct fields (User, Email, SigningEnabled, UserSource)
- [ ] 2.11 Save updated config after refresh completes

### [ ] 3.0 Auto-Detect Profiles from Existing Gitconfig

#### 3.0 Proof Artifact(s)

- CLI: `gws user list` shows profiles auto-detected from `~/.gitconfig` includeIf directives demonstrates profile detection works
- Test: `go test ./internal/user/... -v` passes demonstrates gitconfig parsing logic works
- CLI: Running `gws user list` on system with existing includeIf setup shows correct profiles demonstrates real-world detection

#### 3.0 Tasks

- [ ] 3.1 Create `internal/user/` package directory
- [ ] 3.2 Create `internal/user/gitconfig.go` with types for representing gitconfig sections
- [ ] 3.3 Implement `ParseGitconfig(path string) (*GitConfig, error)` to parse gitconfig file format
- [ ] 3.4 Implement `ExtractIncludeIfs(cfg *GitConfig) []IncludeIf` to find all includeIf directives
- [ ] 3.5 Implement `ParseIncludedConfig(includeIf IncludeIf) (*Profile, error)` to read referenced .gitconfig file and extract user info
- [ ] 3.6 Implement `DetectProfiles() ([]Profile, error)` to orchestrate: parse ~/.gitconfig, find includeIfs, parse each, return profiles
- [ ] 3.7 Derive profile name from includeIf path or gitconfig filename (e.g., `.gitconfig-work` → "work")
- [ ] 3.8 Create `internal/user/gitconfig_test.go` with test fixtures for various gitconfig formats
- [ ] 3.9 Write unit tests for ParseGitconfig with basic key-value pairs
- [ ] 3.10 Write unit tests for ExtractIncludeIfs with gitdir patterns
- [ ] 3.11 Write unit tests for ParseIncludedConfig extracting user.name, user.email, user.signingkey
- [ ] 3.12 Write integration test for DetectProfiles with temporary gitconfig files

### [ ] 4.0 Display User Information in List Command

#### 4.0 Proof Artifact(s)

- CLI: `gws list --user` shows USER, EMAIL, SIGN columns with correct values demonstrates display works
- CLI: `gws list --user` shows "(local)" marker for repos with local .git/config override demonstrates source detection
- CLI: `gws list --user` shows drift indicator (⚠) when effective user differs from stored demonstrates drift detection
- CLI: `gws list --user -o json` includes user fields in JSON output demonstrates JSON integration

#### 4.0 Tasks

- [ ] 4.1 Add `showUser` flag variable and `--user` flag to list command in `cmd/gws/list.go`
- [ ] 4.2 Update `displayTable` function to accept showUser parameter
- [ ] 4.3 Add USER, EMAIL, SIGN column headers when showUser is true
- [ ] 4.4 Display user.name in USER column, append "(local)" if UserSource is "local"
- [ ] 4.5 Display user.email in EMAIL column
- [ ] 4.6 Display "✓" in SIGN column if SigningEnabled is true, empty otherwise
- [ ] 4.7 Implement drift detection: compare stored User/Email with current effective config
- [ ] 4.8 Display "⚠" indicator in STATUS column when drift is detected
- [ ] 4.9 Update column width calculations to account for new columns
- [ ] 4.10 Update `displayJSON` function to include User, Email, SigningEnabled, UserSource fields in output
- [ ] 4.11 Update list command help text with --user flag documentation
- [ ] 4.12 Test display with various repository configurations (global, local, signed, unsigned, drift)

### [ ] 5.0 Implement User Profile Management Commands

#### 5.0 Proof Artifact(s)

- CLI: `gws user list` displays all profiles with name, email, signing status demonstrates listing works
- CLI: `gws user add work --email work@company.com --name "John Doe"` creates profile demonstrates creation works
- CLI: `gws user add work --email duplicate@test.com` fails with duplicate error demonstrates validation works
- CLI: `gws user show work` displays detailed profile information demonstrates show works
- CLI: `gws user remove work` deletes profile (with confirmation if repos use it) demonstrates removal works
- Test: `go test ./internal/user/... -v` passes demonstrates profile management logic works

#### 5.0 Tasks

- [ ] 5.1 Create `internal/user/profile.go` with profile management functions
- [ ] 5.2 Implement `AddProfile(cfg *config.Config, profile config.Profile) error` with duplicate name validation
- [ ] 5.3 Implement `RemoveProfile(cfg *config.Config, name string) error` that checks for repos using profile
- [ ] 5.4 Implement `GetProfile(cfg *config.Config, name string) (*config.Profile, error)` to find profile by name
- [ ] 5.5 Implement `ValidateEmail(email string) error` for basic email format validation
- [ ] 5.6 Create `internal/user/profile_test.go` with tests for AddProfile, RemoveProfile, GetProfile
- [ ] 5.7 Write tests for duplicate profile name validation
- [ ] 5.8 Write tests for email format validation
- [ ] 5.9 Create `cmd/gws/user.go` with parent `user` command and help text
- [ ] 5.10 Implement `user list` subcommand: load config, call DetectProfiles to merge auto-detected, display table
- [ ] 5.11 Implement `user add <name>` subcommand with `--email`, `--name`, `--signing-key`, `--sign-commits` flags
- [ ] 5.12 Implement `user show <name>` subcommand to display detailed profile info
- [ ] 5.13 Implement `user remove <name>` subcommand with confirmation prompt if repos use profile
- [ ] 5.14 Register user command and subcommands with rootCmd in init()
- [ ] 5.15 Update README.md with user profile management examples

### [ ] 6.0 Implement Repository User Assignment

#### 6.0 Proof Artifact(s)

- CLI: `gws user assign myrepo work` sets user.name, user.email in repo's local .git/config demonstrates local assignment works
- CLI: Running `git config user.email` in assigned repo shows correct email demonstrates git config is set
- CLI: `gws user assign myrepo work --use-subdirs` moves repo to profile subdirectory demonstrates subdir mode works
- File: `~/.gitconfig` contains includeIf for profile subdirectory after --use-subdirs demonstrates gitconfig integration
- CLI: `gws user sync` updates stored user info to match effective config for all repos demonstrates sync works
- Test: `go test ./internal/user/... -v` passes demonstrates assignment logic works

#### 6.0 Tasks

- [ ] 6.1 Create `internal/user/assign.go` with assignment functions
- [ ] 6.2 Implement `AssignLocal(repoPath string, profile config.Profile) error` to set user.name, user.email in repo's .git/config
- [ ] 6.3 Implement setting user.signingkey and commit.gpgsign in repo's .git/config if profile has signing config
- [ ] 6.4 Implement `AssignWithSubdirs(cfg *config.Config, repoPath string, profile config.Profile, workspace string) (newPath string, error)` for subdirectory mode
- [ ] 6.5 Implement `CreateProfileSubdir(workspace string, profileName string) (string, error)` to create profile subdirectory
- [ ] 6.6 Implement `MoveRepository(oldPath, newPath string) error` using os.Rename with conflict detection
- [ ] 6.7 Implement `CreateProfileGitconfig(subdirPath string, profile config.Profile) error` to write .gitconfig-<profile> file
- [ ] 6.8 Implement `AddIncludeIf(gitconfigPath string, subdirPath string, profileGitconfigPath string) error` to update ~/.gitconfig
- [ ] 6.9 Create backup of ~/.gitconfig before modification
- [ ] 6.10 Implement `SyncUserInfo(cfg *config.Config) (updated int, error)` to refresh stored user info for all repos
- [ ] 6.11 Create `internal/user/assign_test.go` with tests for AssignLocal
- [ ] 6.12 Write tests for AssignWithSubdirs including repo movement
- [ ] 6.13 Write tests for gitconfig modification (CreateProfileGitconfig, AddIncludeIf)
- [ ] 6.14 Write tests for SyncUserInfo
- [ ] 6.15 Implement `user assign <repo> <profile>` subcommand in `cmd/gws/user.go`
- [ ] 6.16 Add `--use-subdirs` flag to assign subcommand
- [ ] 6.17 Add `--dry-run` flag to assign subcommand to preview changes without applying
- [ ] 6.18 Implement `user sync` subcommand to update stored info from effective config
- [ ] 6.19 Update README.md with user assignment examples including --use-subdirs workflow
- [ ] 6.20 Add warning when assigning in subdirectory mode about repo path changes
