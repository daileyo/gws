# Task 4.0 Proof Artifacts - Deprecate Old Tag Flags and Clean Up Root Command

## Deprecated Tag Flags

| Flag | Short | Replacement |
|------|-------|-------------|
| `--add-tag` | `-d` | `gws tag add <repo> <tag>` |
| `--remove-tag` | `-x` | `gws tag remove <repo> <tag>` |

## Changes

- Removed `flagAddTag`/`flagRemoveTag` variables from `main.go`
- Removed flag registrations from `main.go` `init()`
- Removed tag dispatch block from root `RunE`
- Added `depAddTag`/`depRemoveTag` to `deprecated.go`
- Registered as hidden flags with deprecation warnings
- Added dispatch in `handleDeprecatedFlags()` with arg validation
- Removed "Tag Operations" section from root usage template
- `tag` subcommand now appears in "Available Commands" naturally

## Help Output Verification

`gws --help` shows:
- Available Commands: `add`, `completion`, `help`, `init`, `list`, `print-workspace`, `refresh`, `shell-init`, `tag`, `user`
- **No `--add-tag` or `--remove-tag` visible** (all hidden)
- "Tag Operations" section removed from template

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestDeprecated" -v
=== RUN   TestDeprecatedFlagsAreHidden       (16 sub-tests) --- PASS
=== RUN   TestDeprecatedListEmitsWarning      --- PASS
=== RUN   TestDeprecatedMutualExclusivity     --- PASS
=== RUN   TestDeprecatedRecursiveRequiresAdd  --- PASS
=== RUN   TestDeprecatedGoConflictsWithArgs   --- PASS
=== RUN   TestDeprecatedNoFlagsSet            --- PASS
=== RUN   TestDeprecatedAddTagEmitsWarning    --- PASS
=== RUN   TestDeprecatedAddTagDispatch        --- PASS
=== RUN   TestDepWarningsMap                  (9 sub-tests) --- PASS
```

## Full Suite

```
$ go vet ./... && go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.349s
ok  github.com/daileyo/gws/internal/classifier  0.006s
ok  github.com/daileyo/gws/internal/config      0.008s
ok  github.com/daileyo/gws/internal/discovery    0.033s
ok  github.com/daileyo/gws/internal/filter       0.008s
ok  github.com/daileyo/gws/internal/git          0.048s
ok  github.com/daileyo/gws/internal/user         0.137s
```
