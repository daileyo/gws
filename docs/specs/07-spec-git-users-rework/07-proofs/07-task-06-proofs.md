# 07 Task 6.0 Proof Artifacts - Default User Indicator in List Output

## Test Results

```
=== RUN   TestGetGlobalDefaultUser
    user_test.go:506: Global default user: name="daileyo" email="[REDACTED]"
--- PASS: TestGetGlobalDefaultUser (0.00s)
=== RUN   TestGetGlobalDefaultUser_ParsesGitconfig
--- PASS: TestGetGlobalDefaultUser_ParsesGitconfig (0.00s)
=== RUN   TestGetGlobalDefaultUser_NoUserSection
--- PASS: TestGetGlobalDefaultUser_NoUserSection (0.00s)
=== RUN   TestGetGlobalDefaultUser_MissingFile
--- PASS: TestGetGlobalDefaultUser_MissingFile (0.00s)

All tests pass: make test exits cleanly.
```

## Changes Summary

| File | Change |
|------|--------|
| `internal/git/user.go` | Added `GetGlobalDefaultUser()` - public wrapper for `loadGlobalConfig()` |
| `internal/git/user_test.go` | Added 4 tests for global default user detection |
| `cmd/git-workspace/list.go` | Call `GetGlobalDefaultUser()` before repo loop; append ` *` to display when source is global and name matches default |

## Display Logic

```
if displaySource == global AND globalDefaultUser != nil AND displayUser == globalDefaultUser.Name:
    userDisplay += " *"
```

- Repos using global default identity: `daileyo *`
- Repos with local config override: `Local User (local)`
- Repos with includeIf override: `Work User`
- No global user configured: no markers shown
