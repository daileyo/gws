# 22 Task 3.0 Proofs - Documentation Replacement

## Files

- `cmd/git-workspace/main.go` — helper snippet removed from `rootCmd.Long`
- `README.md` — `cd` added to the subcommand list and a Workspace Navigation feature bullet
- `docs/site/commands-core.md` — new "Navigate to Workspace Root" section
- `docs/site/shell-integration.md` — `gws cd` documented, with the mechanism explained
- `docs/site/getting-started.md` — new "Navigate Your Workspace" quick-start step

## Proof 1: the helper is fully retired

```
$ grep -rn 'cdgws\|alias gcd' README.md docs/ cmd/ Makefile | grep -v '^docs/specs/'
(no matches)
```

The only surviving references are inside `docs/specs/`, where the helper is named as the thing
being replaced.

## Proof 2: root help no longer teaches the helper

Before:

```
  # Navigate to workspace root
  function cdgws() { cd "$(gws print-workspace)"; }
  alias gcd=cdgws
```

After:

```
$ ./build/git-workspace --help
Navigation:
  gws cd                                 # Navigate to the workspace root
  gws my-repo                            # Navigate to repository by name
  gws "api-*"                            # Wildcard match with interactive selection

Shell integration (add to ~/.bashrc or ~/.zshrc):

  # Recommended — sets up 'gws' and tab completion, always up to date:
  export PATH="$HOME/.local/bin:$PATH"
  eval "$(git-workspace shell-init zsh)"   # or: shell-init bash

This provides 'gws cd' for the workspace root — no helper function needed.
```

## Proof 3: `cd` is discoverable in the command list

```
$ ./build/git-workspace --help
Available Commands:
  add             Add a git repository to the workspace
  cd              Navigate to the workspace root
  completion      Generate the autocompletion script for the specified shell
...
```

## Proof 4: documented across the docs site

```
$ grep -c 'gws cd' README.md docs/site/commands-core.md \
      docs/site/shell-integration.md docs/site/getting-started.md
README.md:1
docs/site/commands-core.md:6
docs/site/shell-integration.md:8
docs/site/getting-started.md:1
```

## Proof 5: `print-workspace` positioned as the scripting primitive

`commands-core.md` now carries both sections, with the division stated explicitly:

> `print-workspace` is the scripting primitive; `gws cd` is the interactive command. Because
> `print-workspace` only ever writes the path to stdout, it stays the right choice inside
> scripts and command substitution.

`shell-integration.md` closes its Workspace Navigation section with the same guidance, and
additionally documents *why* `gws cd` needs shell integration — a process cannot change its
parent shell's directory, so the binary prints and the function navigates.

## Proof 6: no regressions

```
$ GIT_CONFIG_GLOBAL=/dev/null go test ./...
ok  	github.com/daileyo/gws/cmd/git-workspace	0.239s
ok  	github.com/daileyo/gws/internal/classifier
ok  	github.com/daileyo/gws/internal/config
ok  	github.com/daileyo/gws/internal/discovery
ok  	github.com/daileyo/gws/internal/filter
ok  	github.com/daileyo/gws/internal/git
ok  	github.com/daileyo/gws/internal/user

$ go vet ./...
(clean)
```

## Note on task 3.7

The Makefile's `install` target prints shell-integration instructions but never referenced the
`cdgws`/`gcd` helper, so there was nothing to update. The task was a check, and it passed.

## Note on running the test suite

`go test ./...` fails in this environment without `GIT_CONFIG_GLOBAL=/dev/null`. The cause is
unrelated to this spec: `core.hooksPath` points at a global `commit-msg` hook that rejects the
plain `-m init` commits the `internal/git` test fixtures create, so fixture setup fails before
the assertions run.

---

## Proof 7 (task 3.8): upgrade note

Added after a real failure: `gws cd` reported `command not found` on a machine where the
binary had been rebuilt but the shell function had not been re-evaluated.

Three compounding causes, all now documented in `docs/site/shell-integration.md`:

1. **The `gws` function is a snapshot taken at shell startup.** Re-running `make build` or even
   `make install` does not update a function already loaded in a running shell.
2. **`make build` does not install.** It writes `./build/git-workspace` only (`Makefile:38`);
   `~/.local/bin/git-workspace` keeps whatever was there before.
3. **Every local build reports `version dev`**, so `--version` cannot distinguish builds. Only
   the `commit:` line does:

   ```
   installed: version dev  commit: 907d3ae  built: 2026-08-11   <- a month old
   build/:    version dev  commit: 9e039f4  built: 2026-09-13
   ```

A pre-existing sentence in the Setup section actively contributed to the confusion:

> ...so your shell integration is always in sync with the installed version. No manual updates
> needed when you upgrade.

That is true only for *new* shells. It has been corrected, and an "Upgrading" section now
states the re-evaluation requirement, names the confusing symptom (`No repositories found
matching 'cd'`), and documents a throwaway-shell recipe for testing a build without touching
an installed setup.

## Note on scope

Only the documentation lands here. Whether the shell function should detect that it is older
than the binary and warn or reload itself is deliberately **not** part of this spec — that
concern predates `cd` and applies equally to spec 20's worktree navigation and spec 21's
PowerShell support. It is recorded as open question 4 in the spec.
