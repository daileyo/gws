# 23 Task 2.0 Proofs - Config Relocation and Automatic Migration

## Files

- `internal/config/migrate.go` (new) — `ensureMigrated`, `migrateLegacyConfig`, `moveFile`, `reportLegacyDirCleanup`
- `internal/config/migrate_test.go` (new) — 10 tests
- `internal/config/config.go` — `Load` migrates then falls back; `Exists` checks both locations; `ConfigVersion` bumped to `1.2.0`

## Proof 1: real migration, end to end

```
$ H=$(mktemp -d); mkdir -p "$H/.gws" "$H/ws"
$ printf '{"version":"1.1.0","workspace":"%s/ws","repositories":[]}' "$H" > "$H/.gws/config.json"
$ chmod 600 "$H/.gws/config.json"

$ HOME=$H ./build/git-workspace list
note: moved config /tmp/tmp.i1qohnmLua/.gws/config.json -> /tmp/tmp.i1qohnmLua/.config/gws/config.json
note: removed empty /tmp/tmp.i1qohnmLua/.gws
No repositories found. Run 'gws init' to discover repositories.

$ ls -l "$H/.config/gws/config.json"
-rw------- ... config.json          <- 0600 preserved

$ ls "$H/.gws"
ls: cannot access ...: No such file or directory   <- empty legacy dir removed

$ HOME=$H ./build/git-workspace list
No repositories found. Run 'gws init' to discover repositories.   <- no notice, idempotent
```

## Proof 2: test suite

```
$ go test -count=1 ./internal/config/ -run 'TestMigrate|TestExists' -v
--- PASS: TestMigrate_MovesLegacyConfig
--- PASS: TestMigrate_PreservesPermissions
--- PASS: TestMigrate_RemovesEmptyLegacyDir
--- PASS: TestMigrate_KeepsNonEmptyLegacyDir
--- PASS: TestMigrate_NoOpWhenNewConfigExists
--- PASS: TestMigrate_RefusesSymlinkedLegacyConfig
--- PASS: TestMigrate_FallsBackToLegacyOnFailure
--- PASS: TestMigrate_UninitializedWhenNeitherExists
--- PASS: TestExists_FindsLegacyConfig
--- PASS: TestExists_FalseWhenNeitherExists
ok  	github.com/daileyo/gws/internal/config
```

## Proof 3: safety properties

| Property | Mechanism | Test |
| --- | --- | --- |
| `0600` never widens | explicit `Chmod` after rename, `perm` passed to `WriteFile` on the copy path | `TestMigrate_PreservesPermissions` |
| Symlink cannot redirect the move | `os.Lstat` + `ModeSymlink` check, refuse | `TestMigrate_RefusesSymlinkedLegacyConfig` |
| User files are never deleted | legacy dir removed only when `ReadDir` is empty | `TestMigrate_KeepsNonEmptyLegacyDir` |
| Cross-filesystem move works | `Rename`, falling back to copy-then-remove | `moveFile` |
| Runs at most once per process | `sync.Once` | `migrateOnce` |
| A working workspace never reports uninitialized | failure warns and falls back to the legacy path | `TestMigrate_FallsBackToLegacyOnFailure` |

## Proof 4: `Exists` had to change too

`Exists()` runs *before* `Load()` in the root command's initialization check. Had it only
looked at the new location, a user with a pending migration would have been told to run
`gws init` — and `init` would then refuse because a workspace already existed. It now checks
both locations. Covered by `TestExists_FindsLegacyConfig`.

## Proof 5: no regressions

```
$ go test -count=1 ./...      # 8/8 packages ok
$ go vet ./...                # clean
```

## Design note

Migration failure is deliberately non-fatal. The alternative — surfacing the error — would
turn a cosmetic relocation into an outage for anyone whose config directory is read-only or on
a full disk. Instead the user gets a warning and their existing config keeps working from its
old home, so the worst case is that the file has not moved yet.
