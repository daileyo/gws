# Task 4.0 Proof Artifacts - Clean Up Shared State and Update Help Text

## filterTags Comment Update

Updated `main.go:26-28`:
```go
// filterTags is shared between the list subcommand (list.go --tag/-t) and
// the deprecated root --tag flag (deprecated.go) for backward compat with
// --list and deprecated user operations (--user -u --tag).
var filterTags []string
```

## --tag Flag Verification

The `--tag` flag on root (registered in `deprecated.go:64`) writes to the shared `filterTags` variable. Code paths verified:
- `deprecated.go:178` validates `--tag` requires `--list` or `--user`
- `deprecated.go:250` passes `filterTags` to list dispatch via `FilterTags` field
- `userupdate.go` reads `filterTags` for batch operations with tag filter

## user assign Help Output

```
$ gws user assign --help
Assign a user profile to a repository, setting user.name and user.email
in the repository's local .git/config file.

Usage:
  git-workspace user assign <repository> <profile> [flags]

Flags:
      --dry-run       Preview changes without applying
  -h, --help          help for assign
      --use-subdirs   Move repository to profile subdirectory
```

## user sync Help Output

```
$ gws user sync --help
Update stored user information for all repositories to match their
current effective git configuration.

Usage:
  git-workspace user sync [flags]

Flags:
  -h, --help   help for sync
```

## user --help Output

```
$ gws user --help
Manage git user profiles for different contexts (work, personal, etc).

Commands:
  gws user list                                    # List all profiles
  gws user add work --email work@company.com       # Add a new profile
  gws user show work                               # Show profile details
  gws user remove work                             # Remove a profile
  gws user assign my-repo work                     # Assign profile to repo
  gws user sync                                    # Sync user info

Short flags:
  gws user -l                                      # List all profiles
  gws user -a work --email work@company.com        # Add a new profile
  gws user -s work                                 # Show profile details
  gws user -d work                                 # Remove a profile

Available Commands:
  add, assign, list, remove, show, sync

Flags:
  -a, --add       -d, --delete       -l, --list       -s, --show
      --email         --name             --signing-key    --sign-commits
```

Short flags `-a`, `-d`, `-l`, `-s` are displayed alongside word-form sub-subcommands.

## gws --help Output

```
$ gws --help
Available Commands:
  add, completion, help, init, list, print-workspace, refresh, shell-init, tag, user

Flags:
  -q, --quiet     -h, --help      --version
```

**No deprecated user flags visible** (`--user`, `--update`, `--delete`, `--all`, `--verbose`, `--git-name`, `--git-email`, `--list-users` all hidden).

## Full Test Suite

```
$ go vet ./... && go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.328s
ok  github.com/daileyo/gws/internal/classifier  0.007s
ok  github.com/daileyo/gws/internal/config      0.008s
ok  github.com/daileyo/gws/internal/discovery    0.027s
ok  github.com/daileyo/gws/internal/filter       0.007s
ok  github.com/daileyo/gws/internal/git          0.032s
ok  github.com/daileyo/gws/internal/user         0.113s
```

All 7 packages pass with no issues.
