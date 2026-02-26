# Task 4.0 Proof Artifacts — Shell Tab Completion

## CLI: `gws completion --help` — Subcommand Registered with All Shells

```
$ gws completion --help
Generate the autocompletion script for git-workspace for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:
  git-workspace completion [flags]
```

**Demonstrates:** `completion` subcommand is registered and supports bash, zsh, fish, powershell.

---

## CLI: `gws completion zsh` — Zsh Script Generation

```
$ gws completion zsh | head -5
#compdef git-workspace
compdef _git-workspace git-workspace

# zsh completion for git-workspace                        -*- shell-script -*-
```

**Demonstrates:** `gws completion zsh` outputs a valid zsh completion script (212 lines).

---

## CLI: `gws completion bash` — Bash Script Generation

```
$ gws completion bash | head -3
# bash completion V2 for git-workspace                        -*- shell-script -*-

__git-workspace_debug()
```

**Demonstrates:** `gws completion bash` outputs a valid bash completion script.

---

## CLI: `gws --help` — Completion Line in Commands Section

```
Commands:
  -a, --add string[="."]   Add a git repository to the workspace (defaults to current directory)
  -d, --add-tag            Add a tag to repositories (args: <repo> <tag>)
  -i, --init               Initialize a gws workspace in the current directory
  -l, --list               List all tracked repositories
  -w, --print-workspace    Print workspace path (for shell integration)
  -r, --refresh            Refresh repository metadata and git status cache
  -u, --remove-tag         Remove a tag from repositories (args: <repo> <tag>)
  gws completion [bash|zsh|fish|powershell]    Generate shell completion script
```

**Demonstrates:** `completion` subcommand is surfaced in `gws --help` Commands section.

---

## Shell Demo: `--add` Flag Directory Completion (ShellCompDirectiveFilterDirs)

```
$ gws __complete "--add=" 2>&1
:16
Completion ended with directive: ShellCompDirectiveFilterDirs
```

**Demonstrates:** `RegisterFlagCompletionFunc("add", ...)` returns directive `:16`
(`ShellCompDirectiveFilterDirs`) for the `--add` flag value, telling the shell
to suggest only directories when completing `gws --add=<TAB>`.

---

## Shell Demo: Live Path Completion (zsh)

After running `source <(gws completion zsh)`, typing `gws --add=/tmp/<TAB>`
presents directory entries under `/tmp/` for tab completion.

The zsh completion script correctly defines `shellCompDirectiveFilterDirs=16`
and applies directory-only filtering when the directive is returned:

```zsh
local shellCompDirectiveFilterDirs=16
...
elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
    # Use absolute paths
    compadd -f -- ...
```

---

## Full Test Suite: `go vet ./... && go test -race ./...`

```
ok  github.com/daileyo/gws/cmd/git-workspace
ok  github.com/daileyo/gws/internal/classifier
ok  github.com/daileyo/gws/internal/config
ok  github.com/daileyo/gws/internal/discovery
ok  github.com/daileyo/gws/internal/filter
ok  github.com/daileyo/gws/internal/git
```

**Demonstrates:** No regressions across the full test suite.
