# Task 5.0 Proof Artifacts - Add Deprecation Layer for Old Flag Forms

## Deprecated Flags

Created `deprecated.go` with:
- `registerDeprecatedFlags()` — registers hidden flags on root for backward compatibility
- `handleDeprecatedFlags()` — dispatches deprecated flag usage with warnings
- `emitDeprecationWarnings()` — prints stderr warnings for each used deprecated flag
- `depWarnings` map — maps each deprecated flag to its replacement form

### Deprecated Command Flags
| Flag | Short | Replacement |
|------|-------|-------------|
| `--list` | `-l` | `gws list` |
| `--init` | `-i` | `gws init` |
| `--add` | `-a` | `gws add [path]` |
| `--recursive` | `-v` | `gws add --recursive` |
| `--refresh` | `-r` | `gws refresh` |
| `--print-workspace` | `-w` | `gws print-workspace` |
| `--go` | `-g` | `gws <repo-name>` |

### Deprecated Filter Flags (compound usage with --list)
| Flag | Short | Replacement |
|------|-------|-------------|
| `--type` | `-y` | `gws list --type` |
| `--tag` | `-t` | `gws list --tag` |
| `--name` | `-n` | `gws list --name` |
| `--path` | `-p` | `gws list --path` |
| `--output` | `-o` | `gws list --output` |
| `--status` | `-s` | `gws list --status` |
| `--show-user` | | `gws list --show-user` |

## Flags NOT Deprecated (remain on root for specs 10/11)
- `--add-tag`, `--remove-tag` (spec 10)
- `--user`, `--update`, `--delete`, `--all`, `--verbose`, `--git-name`, `--git-email`, `--list-users` (spec 11)

## Test Results

```
$ go test ./cmd/git-workspace/ -run "TestDeprecated|TestDepWarnings" -v

=== RUN   TestDeprecatedFlagsAreHidden        (14 sub-tests) --- PASS
=== RUN   TestDeprecatedListEmitsWarning       --- PASS
=== RUN   TestDeprecatedMutualExclusivity      --- PASS
=== RUN   TestDeprecatedRecursiveRequiresAdd   --- PASS
=== RUN   TestDeprecatedGoConflictsWithArgs    --- PASS
=== RUN   TestDeprecatedNoFlagsSet             --- PASS
=== RUN   TestDepWarningsMap                   (7 sub-tests) --- PASS
PASS
```

## Architecture

```
main.go RunE
  └── handleDeprecatedFlags(cmd, args)    # First: check deprecated flags
       ├── mutual exclusivity check
       ├── emitDeprecationWarnings()      # Print warnings to stderr
       └── dispatch to runInit/runAdd/runList/etc.
  └── user flag validation                # Then: validate remaining root flags
  └── tag/user dispatch                   # Then: dispatch tag/user operations
  └── positional arg navigation           # Finally: navigation fallthrough
```

## Cleanup

- Removed `flagList`, `flagInit`, `flagAdd`, `flagRecursive`, `flagRefresh`, `flagPrintWorkspace`, `flagGo` from `main.go`
- Removed their registrations from `main.go` `init()`
- All old flag references in tests updated to `dep*` variables
- `filterTags` remains shared between deprecated.go and user operations

## Full Suite

```
$ make vet && go test ./... -count=1
go vet ./...  ✓
ok  github.com/daileyo/gws/cmd/git-workspace   0.362s
ok  github.com/daileyo/gws/internal/classifier  0.005s
ok  github.com/daileyo/gws/internal/config      0.006s
ok  github.com/daileyo/gws/internal/discovery    0.041s
ok  github.com/daileyo/gws/internal/filter       0.008s
ok  github.com/daileyo/gws/internal/git          0.045s
ok  github.com/daileyo/gws/internal/user         0.138s
```
