# Task 3.0 Proof Artifacts: Auto-Detect Profiles from Existing Gitconfig

## Summary

Implemented the `internal/user/` package for parsing gitconfig files and auto-detecting user profiles from `includeIf` directives in `~/.gitconfig`.

## New Files Created

### `internal/user/gitconfig.go`

Core types and functions for gitconfig parsing:

```go
// Types
type GitConfig struct {
    Path       string
    Sections   map[string]Section
    IncludeIfs []IncludeIf
}

type Section struct {
    Name   string
    Values map[string]string
}

type IncludeIf struct {
    Condition string  // e.g., "gitdir:~/work/"
    Path      string  // e.g., "~/.gitconfig-work"
}

// Functions
func ParseGitconfig(path string) (*GitConfig, error)
func ExtractIncludeIfs(cfg *GitConfig) []IncludeIf
func ParseIncludedConfig(includeIf IncludeIf) (*config.Profile, error)
func DetectProfiles() ([]config.Profile, error)
```

### `internal/user/gitconfig_test.go`

Comprehensive test suite with 12 test functions:

- `TestParseGitconfig_BasicKeyValuePairs` - parses [user] and [core] sections
- `TestParseGitconfig_QuotedValues` - handles double and single quoted values
- `TestParseGitconfig_Comments` - skips # and ; comments
- `TestExtractIncludeIfs_GitdirPatterns` - extracts includeIf directives
- `TestParseIncludedConfig_UserInfo` - reads user.name, email, signingkey, gpgsign
- `TestParseIncludedConfig_NoUserSection` - returns error for missing [user] section
- `TestDeriveProfileName_FromGitdirCondition` - derives names from gitdir paths
- `TestExtractNameFromConfigPath` - derives names from config filenames
- `TestDetectProfiles_Integration` - end-to-end test with temp gitconfig files
- `TestParseGitconfig_FileNotFound` - handles missing files
- `TestExtractIncludeIfs_NilConfig` - handles nil input
- `TestExpandPath` - expands ~ to home directory

## Test Output

```
$ go test ./internal/user/... -v
=== RUN   TestParseGitconfig_BasicKeyValuePairs
--- PASS: TestParseGitconfig_BasicKeyValuePairs (0.00s)
=== RUN   TestParseGitconfig_QuotedValues
--- PASS: TestParseGitconfig_QuotedValues (0.00s)
=== RUN   TestParseGitconfig_Comments
--- PASS: TestParseGitconfig_Comments (0.00s)
=== RUN   TestExtractIncludeIfs_GitdirPatterns
--- PASS: TestExtractIncludeIfs_GitdirPatterns (0.00s)
=== RUN   TestParseIncludedConfig_UserInfo
--- PASS: TestParseIncludedConfig_UserInfo (0.00s)
=== RUN   TestParseIncludedConfig_NoUserSection
--- PASS: TestParseIncludedConfig_NoUserSection (0.00s)
=== RUN   TestDeriveProfileName_FromGitdirCondition
--- PASS: TestDeriveProfileName_FromGitdirCondition (0.00s)
=== RUN   TestExtractNameFromConfigPath
--- PASS: TestExtractNameFromConfigPath (0.00s)
=== RUN   TestDetectProfiles_Integration
--- PASS: TestDetectProfiles_Integration (0.00s)
=== RUN   TestParseGitconfig_FileNotFound
--- PASS: TestParseGitconfig_FileNotFound (0.00s)
=== RUN   TestExtractIncludeIfs_NilConfig
--- PASS: TestExtractIncludeIfs_NilConfig (0.00s)
=== RUN   TestExpandPath
--- PASS: TestExpandPath (0.00s)
PASS
ok      github.com/daileyo/gws/internal/user    0.277s
```

## Key Implementation Details

### Profile Name Derivation

The `deriveProfileName` function uses a priority-based approach:

1. First checks if config file path has a clear pattern (`.gitconfig-work` → "work")
2. Falls back to extracting from gitdir condition (`gitdir:~/work/` → "work")
3. Uses `isGenericDirName` to avoid generic names like "repos", "home", "users"

### Supported Gitconfig Patterns

- `[includeIf "gitdir:~/work/"]` - standard gitdir condition
- `[includeIf "gitdir:~/personal/**"]` - with glob pattern
- `[includeIf "gitdir/i:~/GitHub/"]` - case-insensitive

### Config File Naming Patterns Supported

- `.gitconfig-work` → "work"
- `gitconfig-ado` → "ado"
- `.gitconfig.github` → "github"
- `work.gitconfig` → "work"

### Path Expansion

The `expandPath` function handles:
- `~/` prefix → expands to home directory
- Cleans paths with `filepath.Clean`

## All Tests Pass

```
$ go test ./...
ok      github.com/daileyo/gws/cmd/gws          0.253s
ok      github.com/daileyo/gws/internal/classifier      0.415s
ok      github.com/daileyo/gws/internal/config  0.580s
ok      github.com/daileyo/gws/internal/discovery       0.995s
ok      github.com/daileyo/gws/internal/filter  0.740s
ok      github.com/daileyo/gws/internal/git     1.334s
ok      github.com/daileyo/gws/internal/user    1.113s
```

## Sub-tasks Completed

- [x] 3.1 Create `internal/user/` package directory
- [x] 3.2 Create `internal/user/gitconfig.go` with types for representing gitconfig sections
- [x] 3.3 Implement `ParseGitconfig(path string) (*GitConfig, error)`
- [x] 3.4 Implement `ExtractIncludeIfs(cfg *GitConfig) []IncludeIf`
- [x] 3.5 Implement `ParseIncludedConfig(includeIf IncludeIf) (*Profile, error)`
- [x] 3.6 Implement `DetectProfiles() ([]Profile, error)`
- [x] 3.7 Derive profile name from includeIf path or gitconfig filename
- [x] 3.8 Create `internal/user/gitconfig_test.go` with test fixtures
- [x] 3.9 Write unit tests for ParseGitconfig with basic key-value pairs
- [x] 3.10 Write unit tests for ExtractIncludeIfs with gitdir patterns
- [x] 3.11 Write unit tests for ParseIncludedConfig extracting user.name, user.email, user.signingkey
- [x] 3.12 Write integration test for DetectProfiles with temporary gitconfig files
