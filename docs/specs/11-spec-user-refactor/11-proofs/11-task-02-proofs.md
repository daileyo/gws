# Task 2.0 Proof Artifacts - Short-Flag Tests and Tab Completion

## Tests Added

| Test | What It Verifies |
|------|------------------|
| `TestUserShortFlagsRegistered` (4 sub-tests) | `-a`, `-d`, `-l`, `-s` flags registered on `userCmd` |
| `TestUserShortFlagMutualExclusivity` (7 sub-tests) | All pairwise and all-four flag combinations return error |
| `TestUserSubcommandsRegistered` (6 sub-tests) | `list`, `add`, `show`, `remove`, `assign`, `sync` on `userCmd` |
| `TestUserNoOperationShowsHelp` | `gws user` with no flags shows help (no error) |
| `TestUserShortFlagAddRequiresArgs` | `-a` without args returns usage error |
| `TestUserShortFlagShowRequiresOneArg` | `-s` with 0 or 2+ args returns error |
| `TestUserShortFlagDeleteRequiresOneArg` | `-d` with 0 or 2+ args returns error |
| `TestUserAddFlagsOnUserCmd` (4 sub-tests) | `--email`, `--name`, `--signing-key`, `--sign-commits` on `userCmd` |
| `TestUserProfileCompletion` | Profile name completion returns correct matches |

## Tab Completion

- Added `completeProfileNames()` helper in `user.go` — returns stored + auto-detected profile names matching prefix
- Added `completeProfileThenNone()` as `ValidArgsFunction` on `userCmd` — completes first arg, nothing after
- Same pattern as `completeRepoNames()` / `completeRepoThenNone()` in `tag.go`

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestUser" -v
=== RUN   TestUserShortFlagsRegistered       (4 sub-tests) --- PASS
=== RUN   TestUserShortFlagMutualExclusivity  (7 sub-tests) --- PASS
=== RUN   TestUserSubcommandsRegistered       (6 sub-tests) --- PASS
=== RUN   TestUserNoOperationShowsHelp        --- PASS
=== RUN   TestUserShortFlagAddRequiresArgs    --- PASS
=== RUN   TestUserShortFlagShowRequiresOneArg --- PASS
=== RUN   TestUserShortFlagDeleteRequiresOneArg --- PASS
=== RUN   TestUserAddFlagsOnUserCmd           (4 sub-tests) --- PASS
=== RUN   TestUserProfileCompletion           --- PASS
```

## Full Suite

```
$ go vet ./... && go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.355s
ok  github.com/daileyo/gws/internal/classifier  0.006s
ok  github.com/daileyo/gws/internal/config      0.009s
ok  github.com/daileyo/gws/internal/discovery    0.033s
ok  github.com/daileyo/gws/internal/filter       0.011s
ok  github.com/daileyo/gws/internal/git          0.044s
ok  github.com/daileyo/gws/internal/user         0.128s
```
