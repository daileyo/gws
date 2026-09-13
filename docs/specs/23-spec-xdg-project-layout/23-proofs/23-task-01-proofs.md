# 23 Task 1.0 Proofs - XDG Path Resolution Package

## Files

- `internal/xdg/xdg.go` (new) — `ConfigDir`, `ConfigFile`, `DataDir`, `ProjectsDir`, `RepoProjectsDir`, `LegacyConfigDir`, `LegacyConfigFile`
- `internal/xdg/xdg_test.go` (new) — 8 tests
- `internal/config/config.go` — `GetConfigPath`/`GetConfigDir` now delegate to `xdg`
- `internal/config/config_test.go` — `TestGetConfigDir` updated for the new location

## Proof 1: resolution table

```
$ go test ./internal/xdg/ -v
--- PASS: TestConfigDir/unset_falls_back_to_home
--- PASS: TestConfigDir/absolute_value_is_honored
--- PASS: TestConfigDir/relative_value_is_treated_as_unset
--- PASS: TestProjectsDir/unset_falls_back_to_home
--- PASS: TestProjectsDir/absolute_value_is_honored
--- PASS: TestProjectsDir/relative_value_is_treated_as_unset
--- PASS: TestLayoutIsPlatformIndependent/config
--- PASS: TestLayoutIsPlatformIndependent/data
--- PASS: TestLayoutIsPlatformIndependent/projects
--- PASS: TestConfigFile
--- PASS: TestRepoProjectsDir
--- PASS: TestLegacyConfig
--- PASS: TestHomeDirFailurePropagates
ok  	github.com/daileyo/gws/internal/xdg
```

## Proof 2: cross-platform parity is enforced by a test

`TestLayoutIsPlatformIndependent` asserts the path *relative to home* is exactly
`.config/gws`, `.local/share/gws`, and `.local/share/gws/projects`. A regression would most
likely arrive as a `runtime.GOOS` branch or a switch to `os.UserConfigDir` — either would fail
this test on Windows, where `os.UserConfigDir` returns `%AppData%`.

There is no `runtime.GOOS` branch in the package:

```
$ grep -c 'runtime.GOOS' internal/xdg/xdg.go
0
```

`os.UserHomeDir()` supplies `%USERPROFILE%` on Windows, so `C:\Users\<user>\.config\gws`
falls out without a platform branch.

## Proof 3: relative env values are rejected

Per the XDG specification a relative value is treated as unset. Honoring one would resolve it
against the working directory, letting a stray `cd` change where gws writes. Covered by the
`relative value is treated as unset` cases above.

## Proof 4: legacy path ignores XDG_CONFIG_HOME

`TestLegacyConfig` sets `XDG_CONFIG_HOME=/somewhere/else` and asserts `LegacyConfigDir()` still
returns `<home>/.gws`. The legacy location is a fixed historical path, not a configurable one —
migration must look where old versions actually wrote.

## Proof 5: delegation, no regressions

```
$ go build ./... && go test -count=1 ./...
ok  	github.com/daileyo/gws/cmd/git-workspace
ok  	github.com/daileyo/gws/internal/classifier
ok  	github.com/daileyo/gws/internal/config
ok  	github.com/daileyo/gws/internal/discovery
ok  	github.com/daileyo/gws/internal/filter
ok  	github.com/daileyo/gws/internal/git
ok  	github.com/daileyo/gws/internal/user
ok  	github.com/daileyo/gws/internal/xdg
```

All 8 packages pass. `internal/config` no longer builds `~/.gws` itself; both accessors
delegate, so nothing else in the tree hardcodes a path.

`TestGetConfigDir` asserted the directory was named `.gws`; it now asserts `gws` under
`.config`, which is the behavior change this task exists to make.
