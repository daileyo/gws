# 22 Task 1.0 Proofs - `cd` Subcommand

## Files

- `cmd/git-workspace/cd.go` (new) — `cdCmd` and `runCd`
- `cmd/git-workspace/cd_test.go` (new) — 10 tests
- `cmd/git-workspace/navigate.go` — `isCharDevice` extracted from `isTerminal`
- `cmd/git-workspace/main.go` — `cd` added to the usage template's Navigation section

## Proof 1: stdout/stderr split

`gws cd` prints the path on stdout and the informational line on stderr.

```
$ ./build/git-workspace cd
workspace → /home/daileyo/gws
/home/daileyo/gws

$ ./build/git-workspace cd 2>/dev/null
/home/daileyo/gws
```

Discarding stderr leaves only the path, so `$(gws cd)` captures cleanly.

## Proof 2: quiet mode

```
$ ./build/git-workspace cd -q
/home/daileyo/gws
```

No stderr output at all — this is the form the shell function invokes.

## Proof 3: uninitialized workspace

```
$ HOME=$(mktemp -d) ./build/git-workspace cd
Error: workspace not initialized

To get started, navigate to your projects directory and run:
  gws init
```

Matches the guidance the root command emits rather than surfacing a raw config error.

## Proof 4: TTY-gated shell-integration hint

Piped stdout (the shell-function case) — no hint, as shown in Proof 1.

Real terminal (user ran the binary directly), via a pty:

```
$ script -qec "./build/git-workspace cd" /dev/null
workspace → /home/daileyo/gws
/home/daileyo/gws

Note: changing directory requires shell integration; this printed the path only.
Add to your shell config:  eval "$(git-workspace shell-init zsh)"   # or: bash
PowerShell ($PROFILE):     Invoke-Expression (& git-workspace shell-init powershell | Out-String)
```

The one command whose failure mode is otherwise invisible now explains itself.

## Proof 5: positional arguments rejected

```
$ ./build/git-workspace cd foo
Error: unknown command "foo" for "git-workspace cd"
```

`cobra.NoArgs` means `gws cd foo` surfaces an error instead of silently navigating to the root.

## Proof 6: discoverability

```
$ ./build/git-workspace --help
Available Commands:
  cd              Navigate to the workspace root
...
Navigation:
  gws cd                                            # Navigate to the workspace root
  gws <repo-name>                                   # Navigate to repository by name
```

## Proof 7: test suite

```
$ go test ./cmd/git-workspace/ -run 'TestRunCd|TestCdCmd|TestIsCharDevice' -v
--- PASS: TestRunCd_PrintsPathToStdout
--- PASS: TestRunCd_StdoutHoldsOnlyThePath
--- PASS: TestRunCd_QuietSuppressesStderr
--- PASS: TestRunCd_HintShownWhenStdoutIsTerminal
--- PASS: TestRunCd_NoHintWhenStdoutIsPiped
--- PASS: TestRunCd_UninitializedWorkspace
--- PASS: TestRunCd_WorkspaceDirectoryMissing
--- PASS: TestRunCd_WorkspaceIsNotADirectory
--- PASS: TestCdCmd_RejectsArguments
--- PASS: TestIsCharDevice_NonTerminalFile
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	0.011s
```

## Note on the predicate extraction

`isTerminal(r io.Reader)` was checking **stdin** for interactive selection prompts
(`navigate.go:189,248`, `worktree_navigate.go:117`). `cd` needs the same character-device
test against **stdout**. Rather than pass `os.Stdout` through a reader-shaped signature, the
predicate was extracted to `isCharDevice(f *os.File)` and both call sites delegate to it.
`isTerminal` keeps its signature and behavior, so the existing stdin callers are unchanged.
