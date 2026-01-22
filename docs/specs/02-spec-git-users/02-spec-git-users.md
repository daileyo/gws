# 02-spec-git-users.md

## Introduction/Overview

This feature adds git user/identity management capabilities to gws, allowing developers to track and manage multiple git user profiles (name, email, signing keys) across their repositories. The primary goal is to help developers who work across personal and professional projects maintain correct git identity configuration, detect when the effective git user differs from expectations, and optionally reorganize repositories to support multiple identities.

## Goals

- Enable gws to detect and display the effective git user (name, email) for each repository
- Track whether commit signing is configured for each repository
- Support defining and managing multiple user profiles within gws
- Provide commands to assign user profiles to repositories
- Detect and indicate when a repository's effective user differs from the stored/expected user
- Minimize directory reorganization; prefer flat workspace structure when possible
- Optionally manage gitconfig files (includeIf directives) when subdirectory organization is required

## User Stories

- **As a developer with multiple git identities**, I want to see which user/email is configured for each repository so that I can verify I'm committing with the correct identity before making changes.

- **As a developer working across personal and work projects**, I want to define named profiles (e.g., "personal", "work") so that I can easily assign and switch identities across repositories.

- **As a security-conscious developer**, I want to see whether commit signing is configured for each repository so that I can ensure my commits are signed where required.

- **As a developer managing my workspace**, I want gws to detect when my effective git user differs from what's expected so that I can identify and fix configuration drift.

- **As a developer who values simplicity**, I want gws to keep my repositories in a flat directory structure whenever possible, only creating subdirectories when absolutely necessary for multi-user support.

- **As a developer with existing gitconfig setup**, I want gws to auto-detect my existing profiles from includeIf configurations so that I don't have to manually re-enter them.

## Demoable Units of Work

### Unit 1: User Profile Detection and Storage

**Purpose:** Enable gws to detect git user information from repositories and store it in configuration, providing the foundation for all user management features.

**Functional Requirements:**
- The system shall detect the effective git user.name and user.email for each repository during discovery/refresh
- The system shall detect whether commit signing is configured (user.signingkey and commit.gpgsign settings)
- The system shall store detected user information in the gws config alongside repository metadata
- The system shall parse ~/.gitconfig to auto-detect existing user profiles from includeIf directives
- The system shall store discovered profiles in gws config with name, email, and optional signing key
- The system shall mark repositories where the effective user differs from the stored user with a drift indicator
- The system shall provide a sync command to update stored user info to match effective config

**Proof Artifacts:**
- CLI: `gws refresh` updates user information for all repositories demonstrates detection works
- File: `~/.gws/config.json` contains user profiles and per-repo user data demonstrates storage works
- CLI: `gws user list` shows auto-detected profiles from existing gitconfig demonstrates profile detection
- Test: `internal/git/user_test.go` passes demonstrates user detection logic works

### Unit 2: User Display in List Command

**Purpose:** Allow developers to see user identity information when listing repositories, making it easy to verify correct configuration at a glance.

**Functional Requirements:**
- The system shall add optional USER and EMAIL columns to `gws list` output
- The system shall display user columns when `--user` flag is provided
- The system shall show the inherited global user for repositories without local user config
- The system shall mark repositories with local .git/config user overrides as "local"
- The system shall indicate repositories where effective user differs from stored user (drift indicator)
- The system shall show signing status indicator (e.g., "signed" or icon) when signing is configured
- The system shall include user information in JSON output when `--user` flag is combined with `--output json`

**Proof Artifacts:**
- CLI: `gws list --user` shows USER, EMAIL columns with correct values demonstrates display works
- CLI: `gws list --user` shows "(local)" marker for repos with local config override demonstrates override detection
- CLI: `gws list --user` shows drift indicator when effective differs from stored demonstrates drift detection
- Screenshot: Table output with user columns, signing indicators, and drift markers demonstrates UX

### Unit 3: User Profile Management Commands

**Purpose:** Provide CLI commands for developers to manage user profiles and view profile information.

**Functional Requirements:**
- The system shall provide `gws user list` to display all known profiles with name, email, and signing key status
- The system shall provide `gws user add <name> --email <email> --name <gitname>` to create new profiles
- The system shall provide `gws user add <name> --email <email> --signing-key <key>` to include signing key
- The system shall validate that profile names are unique
- The system shall validate email format when adding profiles
- The system shall provide `gws user remove <name>` to delete a profile (with confirmation if repos use it)
- The system shall provide `gws user show <name>` to display detailed profile information

**Proof Artifacts:**
- CLI: `gws user list` displays all profiles with their configuration demonstrates listing works
- CLI: `gws user add work --email work@company.com --name "John Doe"` creates profile demonstrates creation works
- CLI: `gws user add work --email duplicate@test.com` fails with error demonstrates validation works
- Test: `internal/user/profile_test.go` passes demonstrates profile management logic works

### Unit 4: Repository User Assignment

**Purpose:** Allow developers to assign user profiles to repositories, with intelligent handling of directory structure based on necessity.

**Functional Requirements:**
- The system shall provide `gws user assign <repo> <profile>` to assign a profile to a repository
- The system shall determine if assignment can be done via repo-local .git/config (preferred, maintains flat structure)
- The system shall set user.name, user.email, and optionally user.signingkey and commit.gpgsign in repo-local config
- The system shall detect when subdirectory organization is required (when user wants automatic inheritance for new repos)
- The system shall provide `gws user assign <repo> <profile> --use-subdirs` to opt into subdirectory organization
- When using subdirectories, the system shall create profile subdirectory if it doesn't exist
- When using subdirectories, the system shall move the repository to the appropriate subdirectory
- When using subdirectories, the system shall create/update the profile's .gitconfig-<profile> file
- When using subdirectories, the system shall add/update includeIf directive in ~/.gitconfig
- The system shall update repository path in gws config after moving
- The system shall provide `gws user sync` to update stored user info to match current effective config for all repos

**Proof Artifacts:**
- CLI: `gws user assign myrepo work` sets local git config demonstrates local assignment works
- CLI: Running `git config user.email` in assigned repo shows correct email demonstrates git config is set
- CLI: `gws user assign myrepo work --use-subdirs` moves repo and updates gitconfig demonstrates subdir mode works
- File: `~/.gitconfig` contains includeIf for profile subdirectory demonstrates gitconfig integration
- CLI: `gws user sync` updates stored user info for repos with drift demonstrates sync works
- Test: `internal/user/assign_test.go` passes demonstrates assignment logic works

## Non-Goals (Out of Scope)

1. **SSH key management**: gws will not manage SSH keys or their association with git profiles; only git config settings (user.name, user.email, signingkey)
2. **GPG key generation or management**: gws will reference existing signing keys but will not create, import, or manage GPG/SSH signing keys
3. **Remote authentication**: gws will not manage credentials, tokens, or authentication for git remotes
4. **Commit author rewriting**: gws will not modify existing commit history or authors
5. **Multi-machine synchronization**: User profiles are local to each machine's gws installation
6. **Automatic user detection from remotes**: gws will not infer which user "should" be used based on remote URL patterns (e.g., GitHub org → work email)

## Design Considerations

**CLI Output Format:**

Extended `gws list --user` output:
```
NAME         STATUS       TYPE    USER              EMAIL                   SIGN  PATH
----         ------       ----    ----              -----                   ----  ----
my-project   main ✓       github  John Doe          john@personal.com             /path/to/my-project
work-api     main ✗ ↑2    gitlab  John Doe (local)  john@work.com           ✓     /path/to/work-api
old-repo     main ✓ ⚠     github  Jane Doe          jane@old.com                  /path/to/old-repo
```

- `(local)` indicates user is set in repo's local .git/config
- `⚠` (or similar) indicates drift between stored and effective user
- `✓` in SIGN column indicates commit signing is configured

**Profile List Output:**
```
gws user list

NAME       GIT NAME    EMAIL                 SIGNING
----       --------    -----                 -------
personal   John Doe    john@personal.com     No
work       John Doe    john@work.com         Yes (GPG)
```

## Repository Standards

Follow established repository patterns and conventions:

- **Package Organization**: Create new `internal/user/` package for profile management logic
- **Command Structure**: Add `cmd/gws/user.go` with subcommands following existing patterns (tag.go, untag.go)
- **Config Extension**: Extend `internal/config/config.go` with Profile and per-repo user fields
- **Git Integration**: Extend `internal/git/` package with user detection functions
- **Testing**: Follow table-driven test patterns with test helpers for creating repos with specific user configs
- **Error Handling**: Use wrapped errors with context (`fmt.Errorf("...: %w", err)`)
- **Documentation**: Update README.md with user management section

## Technical Considerations

**Git User Detection:**
- Use `go-git` library to read effective git config (already a project dependency)
- Must handle config cascade: system → global → local → worktree
- Use `git config --get user.name` equivalent via go-git's config APIs
- Cache user info following existing status cache patterns (TTL-based)

**Profile Storage in Config:**
```go
type Config struct {
    Version      string
    Workspace    string
    Profiles     []Profile     // NEW: user profiles
    Repositories []Repository
}

type Profile struct {
    Name       string  // e.g., "work", "personal"
    GitName    string  // user.name value
    Email      string  // user.email value
    SigningKey string  // user.signingkey value (optional)
    SignCommits bool   // commit.gpgsign setting
}

type Repository struct {
    // ... existing fields ...
    User        string  // NEW: stored user.name
    Email       string  // NEW: stored user.email
    SigningEnabled bool // NEW: whether signing is configured
    UserSource  string  // NEW: "global", "local", "includeif"
}
```

**Gitconfig Modification:**
- Parse ~/.gitconfig using go-git's config format parser
- Preserve existing content and comments when adding includeIf
- Create profile gitconfig files in profile subdirectories
- Handle path expansion (~ to home directory)

**Repository Movement:**
- Use `os.Rename()` for same-filesystem moves
- Update gws config with new path after move
- Verify no path conflicts before moving
- Handle case where repo is currently the working directory

## Security Considerations

**Sensitive Data Handling:**
- Signing keys may contain key IDs that should not be logged verbosely
- Profile configurations should have user-only file permissions (0600)
- Do not log or display full signing key values in normal output

**Gitconfig Modification Safety:**
- Create backup of ~/.gitconfig before modification
- Validate gitconfig syntax after modification
- Provide `--dry-run` flag for assignment commands to preview changes

**Proof Artifact Security:**
- Do not commit actual ~/.gitconfig files to the repository
- Use sanitized examples in documentation
- Test with generated/temporary gitconfig files

## Success Metrics

1. **Detection Accuracy**: 100% of repositories show correct effective user after refresh
2. **Profile Auto-Detection**: All existing includeIf-based profiles are discovered automatically
3. **Assignment Success**: User assignment via local config works without directory changes in 90%+ of cases
4. **Performance**: User detection adds <100ms overhead to refresh operation for 100 repositories
5. **User Experience**: Developer can verify correct identity across all repos with single `gws list --user` command

## Open Questions

1. Should `gws user assign` without `--use-subdirs` warn if the repo is in a directory that has an includeIf rule that would override the local config?
2. When using `--use-subdirs`, should gws offer to migrate all repos for a profile at once, or only the specified repo?
3. Should there be a `gws user default <profile>` command to set which profile is used for repos without explicit assignment?
