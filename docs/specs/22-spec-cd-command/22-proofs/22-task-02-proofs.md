# 22 Task 2.0 Proofs - Shell Integration Routing

## Files

- `cmd/git-workspace/shellinit.go` — `cd` dispatch added to all three templates
- `cmd/git-workspace/shellinit_test.go` — `TestShellTemplatesContainCdNavigation`, `TestShellTemplatesKeepPrintWorkspacePassthrough`

## Proof 1: zsh template contains the cd case

```
$ ./build/git-workspace shell-init zsh | grep -A7 '^    cd)'
    cd)
      if [[ -n "$2" ]]; then
        git-workspace "$@"
        return
      fi
      _dest="$(git-workspace cd -q 2>/dev/tty </dev/tty)"
      [[ -n "$_dest" ]] && cd "$_dest"
      ;;
```

## Proof 2: PowerShell template contains the cd branch

```
$ ./build/git-workspace shell-init powershell | grep -A8 "'\^cd\$'"
        '^cd$' {
            if ($rest.Count -gt 0) {
                & git-workspace @args
                return
            }
            $dest = & git-workspace cd -q 2>&1 | Where-Object { $_ -is [string] }
            if ($dest) { Set-Location $dest }
            return
        }
```

## Proof 3: end-to-end directory change, zsh

```
$ script -qec 'zsh -c "export PATH=.../build:$PATH; eval \"$(git-workspace shell-init zsh)\"; \
                       cd /tmp; gws cd >/dev/null 2>&1; echo PWD=$(pwd)"' /dev/null
PWD=/home/daileyo/gws
```

## Proof 4: end-to-end directory change, bash

```
$ script -qec 'bash -c "... eval \"$(git-workspace shell-init bash)\"; \
                        cd /tmp; gws cd >/dev/null 2>&1; echo PWD=$(pwd)"' /dev/null
PWD=/home/daileyo/gws
```

## Proof 5: `gws cd foo` errors instead of silently navigating

```
$ ... gws cd foo 2>&1 | head -1; echo PWD=$(pwd)
Error: unknown command "foo" for "git-workspace cd"
PWD=/tmp
```

The non-empty `$2` guard passes through to the binary, so the `NoArgs` error surfaces and the
shell stays put rather than treating the stray argument as a request for the workspace root.

## Proof 6: repo navigation unregressed

```
$ ... gws git-workspace >/dev/null 2>&1; echo PWD=$(pwd)
PWD=/home/daileyo/gws/git-workspace
```

## Proof 7: template tests

```
$ go test ./cmd/git-workspace/ -run TestShellTemplates -v
--- PASS: TestShellTemplatesRouteSubcommands
--- PASS: TestShellTemplatesContainBinPlaceholder
--- PASS: TestShellTemplatesContainNavigationFallthrough
--- PASS: TestShellTemplatesContainWorktreeNavigation
--- PASS: TestShellTemplatesContainParentNavigation
--- PASS: TestShellTemplatesContainCdNavigation        (zsh, bash, powershell)
--- PASS: TestShellTemplatesKeepPrintWorkspacePassthrough (zsh, bash, powershell)
PASS
```

## Testing note

Two artifacts of the harness cost time and are worth recording:

1. **No controlling TTY.** The templates redirect `2>/dev/tty </dev/tty`, which fails in a
   non-interactive shell, so the command substitution yields nothing and no `cd` happens. All
   existing navigation cases share this redirection. Verifying end-to-end behavior requires a
   pty — hence `script -qec`.
2. **`{BIN}` resolves through PATH.** The generated function calls bare `git-workspace`, which
   picked up an older `~/.local/bin/git-workspace` with no `cd` subcommand, making `gws cd`
   look like a failed repo lookup. The build directory must precede it on PATH when testing.
