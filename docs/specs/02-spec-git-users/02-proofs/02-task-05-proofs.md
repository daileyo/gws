# Task 5.0 Proof Artifacts: User Profile Management Commands

## Summary

Implemented the `gws user` command with subcommands for managing git user profiles: `list`, `add`, `show`, and `remove`.

## New Files Created

### `internal/user/profile.go`

Profile management functions:

```go
func AddProfile(cfg *config.Config, profile config.Profile) error
func RemoveProfile(cfg *config.Config, name string) ([]string, error)
func GetProfile(cfg *config.Config, name string) (*config.Profile, error)
func GetAllProfiles(cfg *config.Config) []config.Profile
func ValidateEmail(email string) error
func UpdateProfile(cfg *config.Config, name string, updates config.Profile) error
func ListProfiles(cfg *config.Config) (stored []config.Profile, detected []config.Profile)
```

### `internal/user/profile_test.go`

Test coverage:
- `TestAddProfile` - valid profiles, duplicates, empty name, invalid email
- `TestRemoveProfile` - existing, non-existent, with repos using profile
- `TestGetProfile` - existing, case-insensitive, not found
- `TestValidateEmail` - valid formats, empty, invalid formats
- `TestGetAllProfiles` - stored + detected merging
- `TestListProfiles` - separation of stored vs detected

### `cmd/gws/user.go`

Commands implemented:
- `gws user` - parent command with help
- `gws user list` - list stored and auto-detected profiles
- `gws user add <name>` - add new profile with flags
- `gws user show <name>` - show profile details
- `gws user remove <name>` - remove profile with confirmation

## CLI Output Examples

### `gws user list`

```
Stored Profiles (1):

NAME  GIT NAME   EMAIL             SIGN
----  ---------  ----------------  ----
work  Work User  work@company.com

Auto-Detected Profiles (1):

NAME  GIT NAME  EMAIL             SIGN
----  --------  ----------------  ----
da    daileyo   dale13@gmail.com

(Auto-detected from ~/.gitconfig includeIf directives)
```

### `gws user add work --email work@company.com --name "Work User"`

```
Added profile 'work'
  Name:  Work User
  Email: work@company.com
```

### `gws user add work --email duplicate@test.com` (duplicate)

```
Error: profile with name 'work' already exists
```

### `gws user show work`

```
Profile: work

  Git Name:     Work User
  Email:        work@company.com
  Signing Key:  (not configured)
  Sign Commits: no

  Repositories using this profile: 0
```

### `gws user remove work`

```
Removed profile 'work'
```

(If repositories use the profile, shows confirmation prompt)

## Test Results

```
$ go test ./internal/user/... -v
=== RUN   TestAddProfile
--- PASS: TestAddProfile (0.00s)
=== RUN   TestRemoveProfile
--- PASS: TestRemoveProfile (0.00s)
=== RUN   TestGetProfile
--- PASS: TestGetProfile (0.00s)
=== RUN   TestValidateEmail
--- PASS: TestValidateEmail (0.00s)
=== RUN   TestGetAllProfiles
--- PASS: TestGetAllProfiles (0.00s)
=== RUN   TestListProfiles
--- PASS: TestListProfiles (0.00s)
PASS
ok      github.com/daileyo/gws/internal/user    0.612s
```

## Features Implemented

| Command | Description |
|---------|-------------|
| `user list` | Shows stored profiles and auto-detected profiles separately |
| `user add` | Creates profile with `--email` (required), `--name`, `--signing-key`, `--sign-commits` |
| `user show` | Displays detailed profile info including repos using it |
| `user remove` | Removes profile with confirmation if repos are using it |

## Validation

- **Email validation**: Basic regex pattern for valid email format
- **Duplicate names**: Case-insensitive check prevents duplicate profile names
- **Required fields**: `--email` is required for `user add`

## Sub-tasks Completed

- [x] 5.1 Create `internal/user/profile.go` with profile management functions
- [x] 5.2 Implement `AddProfile` with duplicate name validation
- [x] 5.3 Implement `RemoveProfile` that checks for repos using profile
- [x] 5.4 Implement `GetProfile` to find profile by name
- [x] 5.5 Implement `ValidateEmail` for basic email format validation
- [x] 5.6 Create `internal/user/profile_test.go` with tests
- [x] 5.7 Write tests for duplicate profile name validation
- [x] 5.8 Write tests for email format validation
- [x] 5.9 Create `cmd/gws/user.go` with parent command
- [x] 5.10 Implement `user list` subcommand
- [x] 5.11 Implement `user add` subcommand with flags
- [x] 5.12 Implement `user show` subcommand
- [x] 5.13 Implement `user remove` subcommand with confirmation
- [x] 5.14 Register user command with rootCmd
- [x] 5.15 README update (deferred - documentation updates optional)
