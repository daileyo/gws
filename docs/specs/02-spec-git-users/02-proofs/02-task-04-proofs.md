# Task 4.0 Proof Artifacts: Display User Information in List Command

## Summary

Implemented the `--user` flag for `gws list` command to display USER, EMAIL, and SIGN columns with git user information for each repository.

## Changes Made

### `cmd/gws/list.go`

Added user display functionality:

```go
// New flag variable
var showUser bool

// New flag in init()
listCmd.Flags().BoolVarP(&showUser, "user", "u", false, "Show git user info (USER, EMAIL, SIGN columns)")

// Updated displayTable function with:
// - Dynamic column building based on showUser flag
// - USER column with "(local)" marker for local overrides
// - EMAIL column
// - SIGN column with "✓" for signing enabled
// - Drift detection with "⚠" indicator in NAME column
```

### `internal/git/user.go`

Enhanced global config loading to follow `[include]` directives:

```go
// New function to parse config with includes
func parseGitConfigWithIncludes(configPath, home string) (*GlobalUserConfig, error)

// Updated parseGitConfigAndIncludes to extract include paths
func parseGitConfigAndIncludes(content string) (*GlobalUserConfig, []string)
```

## CLI Output Examples

### Basic User Display (`gws list --user`)

```
Found 100 repositories:

NAME                          USER           EMAIL             SIGN  TYPE     VISIBILITY  TAGS  PATH
----------------------------  -------------  ----------------  ----  -------  ----------  ----  ----
boeing-thb-deployments        Robert Dailey  dale13@gmail.com        unknown  unknown     -     /Users/daileyo/gws/liatrio/boeing-thb-deployments
copilot-agents                Robert Dailey  dale13@gmail.com        github   private     -     /Users/daileyo/gws/liatrio/copilot-agents
```

### Combined with Status (`gws list --user --status`)

```
NAME                          STATUS      USER           EMAIL             SIGN  TYPE     VISIBILITY  TAGS  PATH
----------------------------  ----------  -------------  ----------------  ----  -------  ----------  ----  ----
copilot-agents                main ✗      Robert Dailey  dale13@gmail.com        github   private     -     /Users/daileyo/gws/liatrio/copilot-agents
```

### JSON Output (`gws list -o json`)

```json
[
  {
    "name": "boeing-thb-deployments",
    "path": "/Users/daileyo/gws/liatrio/boeing-thb-deployments",
    "type": "unknown",
    "visibility": "unknown",
    "user": "Robert Dailey",
    "email": "dale13@gmail.com",
    "user_source": "includeif"
  }
]
```

## Features Implemented

1. **USER Column**: Displays `user.name`, appends "(local)" if UserSource is "local"
2. **EMAIL Column**: Displays `user.email`
3. **SIGN Column**: Shows "✓" if signing is enabled, empty otherwise
4. **Drift Detection**: Compares stored user/email with current effective config, shows "⚠" in NAME column when drift detected
5. **Dynamic Columns**: Column headers and widths adjust based on flags (`--status`, `--user`)
6. **JSON Integration**: User fields automatically included in JSON output via Repository struct

## Test Results

```
$ go test ./...
ok      github.com/daileyo/gws/cmd/gws          0.309s
ok      github.com/daileyo/gws/internal/classifier      (cached)
ok      github.com/daileyo/gws/internal/config  (cached)
ok      github.com/daileyo/gws/internal/discovery       (cached)
ok      github.com/daileyo/gws/internal/filter  (cached)
ok      github.com/daileyo/gws/internal/git     0.534s
ok      github.com/daileyo/gws/internal/user    (cached)
```

## Bug Fix: Global Config Include Support

During implementation, discovered that user detection wasn't working because global gitconfig used `[include]` directive to load default user from a separate file (`~/.gitconfig-default-user`).

Fixed by updating `loadGlobalConfig()` to:
1. Parse main gitconfig for `[include]` sections
2. Read and merge included config files
3. Support `~/` path expansion in include paths

## Sub-tasks Completed

- [x] 4.1 Add `showUser` flag variable and `--user` flag to list command
- [x] 4.2 Update `displayTable` function to accept showUser parameter
- [x] 4.3 Add USER, EMAIL, SIGN column headers when showUser is true
- [x] 4.4 Display user.name in USER column, append "(local)" if UserSource is "local"
- [x] 4.5 Display user.email in EMAIL column
- [x] 4.6 Display "✓" in SIGN column if SigningEnabled is true
- [x] 4.7 Implement drift detection: compare stored User/Email with current effective config
- [x] 4.8 Display "⚠" indicator when drift is detected
- [x] 4.9 Update column width calculations to account for new columns
- [x] 4.10 Update `displayJSON` function to include User, Email, SigningEnabled, UserSource fields
- [x] 4.11 Update list command help text with --user flag documentation
- [x] 4.12 Test display with various repository configurations
