# Task 3.0 Proof Artifacts - Deprecate Root-Level User Flags

## Deprecated User Flags

| Flag | Short | Replacement |
|------|-------|-------------|
| `--user` | | `gws user list` or `gws user -l` |
| `--update` | `-u` | `gws user assign <repo> <profile>` |
| `--delete` | `-D` | `gws user assign (remove local config)` |
| `--all` | | `gws user assign (with --all)` |
| `--verbose` | | `gws user --verbose` |
| `--git-name` | | `gws user add --name` |
| `--git-email` | | `gws user add --email` |
| `--list-users` | | `gws user list` or `gws user -l` |

## Changes

- Moved `flagUser`→`depUser`, `flagUpdate`→`depUpdate`, `flagDelete`→`depDelete`, `flagAll`→`depAll`, `flagVerbose`→`depVerbose`, `flagInlineName`→`depInlineName`, `flagInlineEmail`→`depInlineEmail`, `flagListUsers`→`depListUsers` from `main.go` to `deprecated.go`
- Moved flag registrations from `main.go` `init()` to `registerDeprecatedFlags()` in `deprecated.go`
- Added all 8 user flags to `hiddenFlags` slice
- Added deprecation warning mappings to `depWarnings` map
- Added user flag dependency validation in `handleDeprecatedFlags()` (update requires user, delete requires user, mutual exclusivity, etc.)
- Added dispatch logic: `--list-users` → `runListUsers`, `--user` alone → `runListUsers`, `--user --update` → `runUserUpdate`, `--user --delete` → `runUserDelete`
- Removed user flag routing from root `RunE` in `main.go`
- Removed "User Operations" section and `userFlagUsages` template function from root usage template
- Removed `pflag` import from `main.go` (no longer needed)
- Updated `userupdate.go`, `userdelete.go` to reference renamed `dep*` variables
- Updated `userupdate_test.go`, `userdelete_test.go` to reference renamed `dep*` variables
- Updated `main_test.go` to reference renamed `dep*` variables and include user flags in hidden flags test
- Added `TestDeprecatedUserEmitsWarning`, `TestDeprecatedUserUpdateRequiresUser`, `TestDeprecatedUserDeleteRequiresUser`, `TestDeprecatedUserUpdateDeleteMutualExclusivity` to `deprecated_test.go`

## Help Output Verification

`gws --help` shows:
- Available Commands: `add`, `completion`, `help`, `init`, `list`, `print-workspace`, `refresh`, `shell-init`, `tag`, `user`
- **No `--user`, `--update`, `--delete`, `--all`, `--verbose`, `--git-name`, `--git-email`, `--list-users` visible** (all hidden)
- "User Operations" section removed from template

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestDeprecated|TestUserFlag" -v
=== RUN   TestDeprecatedFlagsAreHidden       (25 sub-tests) --- PASS
=== RUN   TestDeprecatedListEmitsWarning      --- PASS
=== RUN   TestDeprecatedMutualExclusivity     --- PASS
=== RUN   TestDeprecatedRecursiveRequiresAdd  --- PASS
=== RUN   TestDeprecatedGoConflictsWithArgs   --- PASS
=== RUN   TestDeprecatedNoFlagsSet            --- PASS
=== RUN   TestDeprecatedAddTagEmitsWarning    --- PASS
=== RUN   TestDeprecatedAddTagDispatch        --- PASS
=== RUN   TestDeprecatedUserEmitsWarning      --- PASS
=== RUN   TestDeprecatedUserUpdateRequiresUser --- PASS
=== RUN   TestDeprecatedUserDeleteRequiresUser --- PASS
=== RUN   TestDeprecatedUserUpdateDeleteMutualExclusivity --- PASS
=== RUN   TestUserFlagValidation              (6 sub-tests) --- PASS
=== RUN   TestUserFlagsRegistered             (7 sub-tests) --- PASS
```

## Full Suite

```
$ go vet ./... && go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.312s
ok  github.com/daileyo/gws/internal/classifier  0.006s
ok  github.com/daileyo/gws/internal/config      0.007s
ok  github.com/daileyo/gws/internal/discovery    0.025s
ok  github.com/daileyo/gws/internal/filter       0.007s
ok  github.com/daileyo/gws/internal/git          0.034s
ok  github.com/daileyo/gws/internal/user         0.111s
```
