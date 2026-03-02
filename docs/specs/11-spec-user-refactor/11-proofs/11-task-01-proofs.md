# Task 1.0 Proof Artifacts - Add Short-Flag Aliases to User Subcommand

## Short-Flag Registration

| Flag | Short | Operation | Equivalent |
|------|-------|-----------|------------|
| `--add` | `-a` | Add a profile | `gws user add` |
| `--delete` | `-d` | Remove a profile | `gws user remove` |
| `--list` | `-l` | List all profiles | `gws user list` |
| `--show` | `-s` | Show profile details | `gws user show` |

## Changes

- Added `userFlagAdd`, `userFlagDelete`, `userFlagList`, `userFlagShow` bool vars in `user.go`
- Added `RunE` to `userCmd` with mutual exclusivity check and dispatch routing
- Registered `-a`, `-d`, `-l`, `-s` flags on `userCmd` via `BoolVarP`
- Registered `--email`, `--name`, `--signing-key`, `--sign-commits` on `userCmd` (same backing vars as `userAddCmd`) so they work with `-a`
- Updated `userCmd.Long` help text to document short flags alongside word-form commands
- Reset `userCmd` usage template to Cobra default (same pattern as `tagCmd`)

## RunE Dispatch Logic

- Counts active short flags; returns error if more than one is set
- `-l` → delegates to `userListCmd.RunE`
- `-a` → validates `len(args) >= 1`, delegates to `userAddCmd.RunE`
- `-s` → validates `len(args) == 1`, delegates to `userShowCmd.RunE`
- `-d` → validates `len(args) == 1`, delegates to `userRemoveCmd.RunE`
- No flag + no sub-subcommand → `cmd.Help()`

## Test Results

```
$ go vet ./... && go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.324s
ok  github.com/daileyo/gws/internal/classifier  0.006s
ok  github.com/daileyo/gws/internal/config      0.007s
ok  github.com/daileyo/gws/internal/discovery    0.028s
ok  github.com/daileyo/gws/internal/filter       0.007s
ok  github.com/daileyo/gws/internal/git          0.035s
ok  github.com/daileyo/gws/internal/user         0.115s
```
