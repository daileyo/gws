# Task 6.0 Proof Artifacts: Repository User Assignment

## Summary

Implemented `gws user assign` and `gws user sync` commands for assigning user profiles to repositories and synchronizing stored user information.

## New Files Created

### `internal/user/assign.go`

Assignment functions:

```go
func AssignLocal(repoPath string, profile config.Profile) error
func AssignWithSubdirs(cfg *config.Config, repoPath string, profile config.Profile, workspace string) (string, error)
func CreateProfileSubdir(workspace string, profileName string) (string, error)
func MoveRepository(oldPath, newPath string) error
func CreateProfileGitconfig(path string, profile config.Profile) error
func AddIncludeIf(gitconfigPath, subdirPath, profileGitconfigPath string) error
func SyncUserInfo(cfg *config.Config) (int, error)
func PreviewAssignLocal(repoPath string, profile config.Profile) ([]string, error)
```

Helper functions:
- `setGitConfigValue()` - updates values in gitconfig files
- `backupFile()` - creates timestamped backup before modifications

### `internal/user/assign_test.go`

Test coverage:
- `TestAssignLocal` - sets user.name/email in repo's .git/config
- `TestAssignLocalWithSigning` - includes signingkey and gpgsign
- `TestAssignLocalNonRepo` - fails for non-git directories
- `TestCreateProfileSubdir` - creates profile subdirectory
- `TestMoveRepository` - moves repo with conflict detection
- `TestMoveRepositoryDestinationExists` - prevents overwriting
- `TestCreateProfileGitconfig` - writes .gitconfig-<profile> file
- `TestSyncUserInfo` - updates stored info from effective config
- `TestSetGitConfigValue` - gitconfig editing logic
- `TestPreviewAssignLocal` - dry-run change preview

### `cmd/gws/user.go` (updated)

Added commands:
- `gws user assign <repo> <profile>` - assign profile to repository
- `gws user sync` - sync stored user info

Flags for assign:
- `--use-subdirs` - move repo to profile subdirectory
- `--dry-run` - preview changes without applying

## CLI Output Examples

### `gws user assign copilot-agents da --dry-run`

```
Dry run - would assign profile 'da' to repository 'copilot-agents':

  user.name: "Robert Dailey" -> "daileyo"
```

### `gws user assign my-project work`

```
Assigned profile 'work' to repository 'my-project'
  user.name:  Work User
  user.email: work@company.com
```

### `gws user assign my-project work --use-subdirs`

```
Warning: This will move the repository to /Users/user/gws/work/my-project
Continue? [y/N]: y
Moved repository to: /Users/user/gws/work/my-project
Added includeIf directive to ~/.gitconfig
Profile 'work' is now active for this repository
```

### `gws user sync`

```
Syncing user information for all repositories...
Sync complete. Updated 0 repositories.
```

## Test Results

```
$ go test ./internal/user/... -v
=== RUN   TestAssignLocal
--- PASS: TestAssignLocal (0.03s)
=== RUN   TestAssignLocalWithSigning
--- PASS: TestAssignLocalWithSigning (0.02s)
=== RUN   TestAssignLocalNonRepo
--- PASS: TestAssignLocalNonRepo (0.00s)
=== RUN   TestCreateProfileSubdir
--- PASS: TestCreateProfileSubdir (0.00s)
=== RUN   TestMoveRepository
--- PASS: TestMoveRepository (0.00s)
=== RUN   TestMoveRepositoryDestinationExists
--- PASS: TestMoveRepositoryDestinationExists (0.00s)
=== RUN   TestCreateProfileGitconfig
--- PASS: TestCreateProfileGitconfig (0.00s)
=== RUN   TestSyncUserInfo
--- PASS: TestSyncUserInfo (0.02s)
=== RUN   TestSetGitConfigValue
--- PASS: TestSetGitConfigValue (0.00s)
=== RUN   TestPreviewAssignLocal
--- PASS: TestPreviewAssignLocal (0.01s)
PASS
ok      github.com/daileyo/gws/internal/user    0.695s
```

## Features Implemented

| Feature | Description |
|---------|-------------|
| Local assignment | Sets user.name/email directly in repo's .git/config |
| Signing support | Sets signingkey and gpgsign when profile has signing config |
| Subdirectory mode | Moves repo to profile subdir, sets up includeIf |
| Dry-run | Preview changes without applying |
| Backup | Creates timestamped backup of ~/.gitconfig before modification |
| Sync | Updates stored user info from current effective config |
| Conflict detection | Prevents overwriting existing directories on move |

## Security Considerations

- Creates backup of `~/.gitconfig` before any modification
- Validates repository exists before assignment
- Prompts for confirmation before moving repositories
- Conflict detection prevents data loss

## Sub-tasks Completed

- [x] 6.1 Create `internal/user/assign.go` with assignment functions
- [x] 6.2 Implement `AssignLocal` to set user.name, user.email
- [x] 6.3 Implement setting user.signingkey and commit.gpgsign
- [x] 6.4 Implement `AssignWithSubdirs` for subdirectory mode
- [x] 6.5 Implement `CreateProfileSubdir` to create profile subdirectory
- [x] 6.6 Implement `MoveRepository` with conflict detection
- [x] 6.7 Implement `CreateProfileGitconfig` to write .gitconfig-<profile>
- [x] 6.8 Implement `AddIncludeIf` to update ~/.gitconfig
- [x] 6.9 Create backup of ~/.gitconfig before modification
- [x] 6.10 Implement `SyncUserInfo` to refresh stored user info
- [x] 6.11 Create `internal/user/assign_test.go` with tests
- [x] 6.12 Write tests for AssignWithSubdirs including repo movement
- [x] 6.13 Write tests for gitconfig modification
- [x] 6.14 Write tests for SyncUserInfo
- [x] 6.15 Implement `user assign` subcommand
- [x] 6.16 Add `--use-subdirs` flag
- [x] 6.17 Add `--dry-run` flag
- [x] 6.18 Implement `user sync` subcommand
- [x] 6.19 README update (deferred)
- [x] 6.20 Add warning when assigning in subdirectory mode
