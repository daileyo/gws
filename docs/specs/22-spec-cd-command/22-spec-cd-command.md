# 22-spec-cd-command

## Introduction/Overview

Navigating to the workspace root is currently a manual, user-authored step. The README and
`main.go` root help both instruct users to hand-write shell helpers:

```sh
function cdgws() { cd "$(gws print-workspace)"; }
alias gcd=cdgws
```

This is friction in an otherwise batteries-included shell integration: `gws <repo>`,
`gws <repo> -p`, and `gws <repo> -wt <branch>` all navigate for you, but "go to the
workspace root" — the most common jump of all — requires copy-pasting a function into your
rc file, and it only exists for bash/zsh users who read that section of the docs.

This spec adds a first-class `gws cd` command that navigates to the workspace root through
the existing shell-integration mechanism, and retires the hand-written helper from the docs.

## Goals

- Add a `cd` subcommand that resolves and prints the workspace root path
- Route `cd` through the `gws` shell function in all three shell templates (zsh, bash, powershell) so the parent shell actually changes directory
- Give users without shell integration a clear, actionable message instead of a silent no-op
- Remove the hand-written `cdgws`/`gcd` helper from the README, root help text, and docs site, replacing it with `gws cd`
- Preserve `print-workspace` unchanged as the scripting-oriented primitive

## User Stories

- **As a gws user**, I want to type `gws cd` to land at my workspace root so that I do not have to maintain a shell alias for the single most common navigation I perform.
- **As a new user**, I want workspace-root navigation to work immediately after running `shell-init`, so that I do not have to read the docs and paste a helper function to get a complete experience.
- **As a PowerShell user**, I want the same `gws cd` behavior as my bash/zsh colleagues, since the hand-written helper in the docs was never provided for PowerShell.

## Demoable Units of Work

### Unit 1: `cd` Subcommand

**Purpose:** Provide the binary-side half of the command — resolve the workspace root and print it on stdout, following the same stdout/stderr contract that `navigate` and `parent` already use so shell command substitution works.

**Functional Requirements:**

- The command shall be registered as `gws cd` with `cobra.NoArgs`
- The command shall print the workspace root path from `config.Workspace` to **stdout**, one line, no trailing decoration
- The command shall accept `-q` / `--quiet` to suppress the informational line, matching `parent` and root navigation
- In non-quiet mode the command shall print a human-readable line to **stderr** (e.g. `workspace → /home/user/gws`) so that `$(gws cd)` still yields only the path
- When the workspace is not initialized, the command shall fail with the same guidance the root command already emits (`Error: workspace not initialized` + `gws init` hint) rather than a bare config error
- When **stdout** is a character device — meaning the user ran the raw binary rather than the shell function, which always captures stdout through a pipe — the command shall additionally print a hint to stderr explaining that shell integration is required for the directory change to take effect, and naming `shell-init`
- The command shall be excluded from repo-name completion (it takes no arguments)

**Proof Artifacts:**

- CLI: `git-workspace cd` prints the workspace root path demonstrates resolution works
- CLI: `git-workspace cd -q` prints only the path with no stderr output demonstrates the quiet contract
- CLI: `git-workspace cd` in an uninitialized environment prints the init guidance demonstrates error handling
- Test: `cd_test.go` covering stdout/stderr split, quiet mode, and the uninitialized path demonstrates correctness

### Unit 2: Shell Integration Routing

**Purpose:** Make `gws cd` actually change directory. A subprocess cannot change its parent's working directory, so the `gws` shell function must intercept `cd`, capture the path, and perform the `cd`/`Set-Location` itself — exactly as it already does for repo navigation.

**Functional Requirements:**

- The zsh, bash, and PowerShell templates shall each add a `cd` case to the `gws` function's dispatch
- The `cd` case shall invoke `{BIN} cd -q`, capture stdout, and `cd`/`Set-Location` to the result only when the result is non-empty
- The `cd` case shall follow the existing stderr/stdin redirection convention (`2>/dev/tty </dev/tty` for bash/zsh) so error output and the init hint remain visible
- A non-empty second argument to `gws cd` shall be passed through to the binary rather than swallowed, so that the `NoArgs` error surfaces to the user instead of silently navigating to the root
- `print-workspace` shall remain in the passthrough subcommand list, unchanged
- `shellinit_test.go` shall be extended with a `TestShellTemplatesContainCdNavigation` covering all three templates, following the structural-validation pattern used by the existing worktree and parent navigation tests

**Proof Artifacts:**

- CLI: `git-workspace shell-init zsh | grep -A3 'cd)'` shows the cd routing block demonstrates the template contains the case
- CLI: In a live shell with the function loaded, `gws cd` changes the shell's working directory to the workspace root demonstrates end-to-end behavior
- Test: `TestShellTemplatesContainCdNavigation` passes for zsh, bash, and powershell demonstrates all three templates are wired
- Test: `go test ./cmd/git-workspace/` passes with no regressions demonstrates existing routing is intact

### Unit 3: Documentation Replacement

**Purpose:** Retire the hand-written helper. Every place that currently tells users to define `cdgws`/`gcd` should tell them to use `gws cd` instead.

**Functional Requirements:**

- The `rootCmd.Long` shell-integration example in `main.go` shall drop the `cdgws`/`gcd` function definition
- `README.md` shall document `gws cd` in the command list and remove the helper-function snippet
- `docs/site/commands-core.md` shall gain a `gws cd` entry alongside the other navigation commands
- `docs/site/shell-integration.md` shall describe `gws cd` in the bash/zsh and PowerShell sections and remove the helper snippet
- `docs/site/getting-started.md` shall use `gws cd` wherever it currently demonstrates reaching the workspace root
- Documentation shall state plainly that `gws cd` requires shell integration, and that `print-workspace` remains the right primitive for scripts
- Documentation shall state that the `gws` function is generated once at shell startup, so upgrading the binary does not update an already-loaded function — `shell-init` must be re-evaluated (or a new shell started) before newly added commands are routed

**Proof Artifacts:**

- Documentation: `grep -rn 'cdgws\|alias gcd' README.md docs/ cmd/` returns no results demonstrates the helper is fully retired
- Documentation: `gws cd` appears in commands-core.md and shell-integration.md demonstrates the replacement is documented
- Documentation: shell-integration.md explains that an upgraded binary needs `shell-init` re-evaluated demonstrates the stale-function failure mode is documented

## Non-Goals (Out of Scope)

1. **`gws cd <repo>`**: Repo navigation is already `gws <repo>`; adding a second spelling creates two ways to do one thing. Deliberately excluded.
2. **`gws cd -`**: Previous-directory tracking requires shell-side state that the current template design does not carry. Out of scope.
3. **`gws cd --config` / `--projects`**: These depend on the XDG locations introduced in spec 23. Recorded as an open question there.
4. **Removing `print-workspace`**: It remains the scripting primitive and is depended on by existing user shell configs. It is not deprecated by this spec.
5. **Multi-workspace switching**: There is one workspace root; selecting between several is not part of this work.

## Design Considerations

The stdout/stderr split is the load-bearing detail. `parent.go` and `navigate.go` already
establish it: the path goes to stdout so `$(...)` captures it cleanly, everything human-facing
goes to stderr so it reaches the terminal without polluting the capture. `cd` follows the
same contract, which is also why the shell function passes `-q`.

The TTY hint exists because `gws cd` is the one command whose failure mode is invisible.
Running the raw binary prints a path and appears to do nothing, with no indication that the
shell function is missing. Detecting a character device on stdout is a reliable proxy for
"a human ran this directly," since the shell function always captures stdout through a pipe.

Note that the existing helper checks the wrong stream for this purpose: `isTerminal(r io.Reader)`
in `navigate.go:232` is used to decide whether *stdin* can drive an interactive selection
prompt. Its underlying check — type-assert to `*os.File`, `Stat()`, test `os.ModeCharDevice` —
is exactly what `cd` needs, but applied to stdout. Rather than passing `os.Stdout` into a
reader-typed function, extract the shared predicate into a small `*os.File`-typed helper that
both call sites delegate to, and keep an overridable indirection so tests can force both
branches.

## Repository Standards

- Go source files follow standard `gofmt` formatting
- Tests use the standard `testing` package with table-driven tests where applicable
- Shell templates are Go `const` strings in `shellinit.go` using the `{BIN}` placeholder
- Template tests validate structural properties, not exact string matches
- Cobra is used for CLI structure; commands live in `cmd/git-workspace/<name>.go` with a matching `_test.go`
- Navigation commands print paths to stdout and human-facing output to stderr
- Conventional commits (`feat(cmd): ...`, `docs: ...`)

## Technical Considerations

- A child process cannot change its parent shell's working directory. The binary can only ever print a path; the `cd` must happen in the shell function. This is the same mechanism repo navigation already uses and requires no new machinery.
- `cd` is a shell builtin in bash/zsh. Because the dispatch happens inside the `gws` function and matches on `$1`, there is no shadowing risk — `gws cd` is unambiguous and the builtin is untouched.
- In PowerShell, `cd` is an alias for `Set-Location`. The same reasoning applies: the match is on the first argument to `gws`, not on a command name.
- `isTerminal` in `navigate.go:232` takes an `io.Reader` and is called with stdin throughout (`navigate.go:189,248`, `worktree_navigate.go:117`). `cd` needs the same character-device test against **stdout**, so the predicate should be extracted to a `*os.File`-typed helper and shared, rather than passing `os.Stdout` through the reader-shaped signature.
- `config.Workspace` is currently an absolute path recorded at `init` time. If the directory has since been deleted or renamed, `gws cd` will print a path that no longer exists. The command should stat the path and report a clear error rather than emitting a dead path for the shell to fail on.

## Security Considerations

- The workspace path is read from the user's own config and printed verbatim; no new input is accepted, so there is no injection surface.
- A workspace path containing spaces or shell metacharacters must remain safe. The templates already quote captured paths (`cd "$dest"`); the `cd` case must follow the same quoting.
- No credentials, network access, or elevated permissions are involved.

## Success Metrics

1. **Zero-config navigation**: A user who has run `shell-init` can type `gws cd` and reach their workspace root with no rc-file editing.
2. **Helper retired**: No occurrence of `cdgws` or `alias gcd` remains in the repo.
3. **Parity**: `gws cd` behaves identically in zsh, bash, and PowerShell.
4. **No regressions**: All existing tests pass; `print-workspace` behavior is byte-for-byte unchanged.

## Open Questions

1. Should `gws cd` gain `--config` and `--projects` flags once spec 23 relocates those directories to XDG paths? (Deferred to spec 23 by decision in round 1.)
2. Should the workspace root also be offered as a completion candidate somewhere, or is `gws cd` discoverable enough on its own?
3. If `config.Workspace` points at a missing directory, should `gws cd` fail outright or fall back to the nearest existing ancestor?
4. Should the `gws` function detect that it is older than the binary it calls and warn, or reload itself? Deferred — this is a general shell-function concern that predates this spec (it applies equally to spec 20's worktree navigation and spec 21's PowerShell support), and it carries its own design trade-offs around startup cost. This spec documents the failure mode; a separate spec should decide the mechanism.
