# Task 1.0 Proof Artifacts - PowerShell `gws` Wrapper Function Template

## CLI Output: `git-workspace shell-init powershell`

```powershell
$ ./git-workspace shell-init powershell

# git-workspace shell integration — do not edit, managed via shell-init
function gws {
    if ($args.Count -eq 0) {
        & git-workspace
        return
    }

    $first = $args[0]
    $rest = @()
    if ($args.Count -gt 1) {
        $rest = $args[1..($args.Count - 1)]
    }

    switch -Regex ($first) {
        '^(list|init|add|refresh|print-workspace|tag|user|completion|shell-init|help|__.*)$' {
            & git-workspace @args
            return
        }
        '^worktree$' {
            if ($rest.Count -eq 0) {
                & git-workspace @args
                return
            }
            $second = $rest[0]
            switch ($second) {
                { $_ -in 'list', 'align', 'add' } {
                    & git-workspace @args
                    return
                }
                default {
                    $dest = & git-workspace @args -q 2>&1 | Where-Object { $_ -is [string] }
                    if ($dest) { Set-Location $dest }
                    return
                }
            }
        }
        '^(-p|--parent|parent)$' {
            $second = if ($rest.Count -gt 0) { $rest[0] } else { $null }
            $dest = & git-workspace parent $second -q 2>&1 | Where-Object { $_ -is [string] }
            if ($dest) { Set-Location $dest }
            return
        }
        '^-' {
            & git-workspace @args
            return
        }
        default {
            $second = if ($rest.Count -gt 0) { $rest[0] } else { $null }

            if ($second -eq '-p' -or $second -eq '--parent') {
                $dest = & git-workspace parent $first -q 2>&1 | Where-Object { $_ -is [string] }
                if ($dest) { Set-Location $dest }
                return
            }
            elseif ($second -eq '-wt') {
                if ($rest.Count -gt 1) {
                    $branch = $rest[1]
                    $dest = & git-workspace $first --worktree $branch -q 2>&1 | Where-Object { $_ -is [string] }
                } else {
                    $dest = & git-workspace $first --worktree -q 2>&1 | Where-Object { $_ -is [string] }
                }
                if ($dest) { Set-Location $dest }
                return
            }
            else {
                $dest = & git-workspace $first -q 2>&1 | Where-Object { $_ -is [string] }
                if ($dest) { Set-Location $dest }
                return
            }
        }
    }
}
& git-workspace completion powershell | Invoke-Expression
(& git-workspace completion powershell) -replace 'git-workspace', 'gws' | Invoke-Expression
```

## Verification

- Template outputs valid PowerShell with `function gws` defined
- All subcommands routed via `switch -Regex` passthrough
- Worktree routing: `list`, `align`, `add` pass through; others navigate via `Set-Location`
- Parent navigation: `-p`, `--parent`, `parent` all call `{BIN} parent`
- `-wt` flag: handles both `gws <repo> -wt <branch>` and bare `gws <repo> -wt`
- Default navigation: captures stdout with `-q` and calls `Set-Location`
- Stderr handling: `2>&1 | Where-Object { $_ -is [string] }` filters ErrorRecords (PS 5.1/7+ compatible)
- `{BIN}` placeholder correctly replaced with `git-workspace` at runtime
- Completion registered for both `git-workspace` and `gws` via `-replace` approach
- `pwsh` accepted as alias for `powershell`
- Build compiles successfully with `go build ./cmd/git-workspace/`
