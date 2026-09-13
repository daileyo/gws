# 22 Tasks - cd Command

## Relevant Files

- `cmd/git-workspace/cd.go` - **New.** The `cd` subcommand: workspace resolution, stdout/stderr split, quiet mode, TTY hint.
- `cmd/git-workspace/cd_test.go` - **New.** Tests for the `cd` command.
- `cmd/git-workspace/main.go` - Root command. Contains the `cdgws`/`gcd` snippet in `rootCmd.Long` that this spec removes; also the usage template where `cd` should be listed under Navigation.
- `cmd/git-workspace/navigate.go` - `isTerminal` (line 232) and the `isTerminalFunc` indirection (line 19), currently applied to **stdin** at lines 189 and 248. The character-device predicate is extracted here for reuse against stdout. Also the reference for the stdout/stderr contract.
- `cmd/git-workspace/parent.go` - Closest existing analogue for a path-printing navigation subcommand with `-q`.
- `cmd/git-workspace/shellinit.go` - The three shell templates. Each needs a `cd` case in the `gws` function dispatch.
- `cmd/git-workspace/shellinit_test.go` - Structural template tests. New `TestShellTemplatesContainCdNavigation` goes here.
- `README.md` - Command list and the shell-integration snippet containing the helper functions.
- `docs/site/commands-core.md` - Core command reference; needs a `gws cd` entry.
- `docs/site/shell-integration.md` - bash/zsh and PowerShell integration docs; contains the helper snippet.
- `docs/site/getting-started.md` - First-run walkthrough; references reaching the workspace root.

### Notes

- Run tests with `go test ./cmd/git-workspace/`. Build with `make build` (outputs to `build/git-workspace`).
- The binary can never change the parent shell's directory — every task here splits into a binary half (print the path) and a shell half (perform the `cd`).
- Templates use `{BIN}` as the binary placeholder; template tests assert structure, not exact strings.
- Follow conventional commits: `feat(cmd): add cd command`, `feat(shell): route cd through gws function`, `docs: replace cdgws helper with gws cd`.

## Tasks

### [ ] 1.0 `cd` Subcommand

Add `cmd/git-workspace/cd.go` defining a `cd` subcommand that resolves the workspace root from config and prints it to stdout, with the informational line on stderr, `-q` to suppress it, and a shell-integration hint when run directly from a terminal.

#### 1.0 Proof Artifact(s)

- CLI: `make build && ./build/git-workspace cd` prints the workspace root on stdout and `workspace → <path>` on stderr demonstrates the stdout/stderr split
- CLI: `./build/git-workspace cd -q` emits the path only demonstrates quiet mode
- CLI: `HOME=$(mktemp -d) ./build/git-workspace cd` prints the `gws init` guidance demonstrates the uninitialized path
- CLI: `./build/git-workspace cd | cat` emits no TTY hint while a direct terminal run does demonstrates the hint is TTY-gated

#### 1.0 Tasks

- [ ] 1.1 Create `cmd/git-workspace/cd.go` with a `cdCmd` cobra command: `Use: "cd"`, `Args: cobra.NoArgs`, short/long help describing workspace-root navigation and noting that shell integration is required
- [ ] 1.2 Register `cdCmd` on `rootCmd` in an `init()` and add a `-q`/`--quiet` bool flag matching `parent.go`'s pattern
- [ ] 1.3 Implement `runCd(quiet bool, stdout, stderr io.Writer) error`: load config, and on a missing config emit the same two-line `workspace not initialized` + `gws init` guidance the root command uses rather than the raw config error
- [ ] 1.4 Stat `cfg.Workspace`; if it does not exist or is not a directory, return a clear error naming the recorded path and suggesting `gws init` — do not print a dead path
- [ ] 1.5 Print the workspace path to stdout with `fmt.Fprintln`; when `!quiet`, print `workspace → <path>` to stderr first
- [ ] 1.6 Extract the character-device predicate from `isTerminal` (`navigate.go:232`) into a `*os.File`-typed helper, leaving the existing stdin call sites (`navigate.go:189,248`, `worktree_navigate.go:117`) behaviorally unchanged; then, when `!quiet` and stdout is a character device, print a hint to stderr that `gws cd` requires shell integration, naming `git-workspace shell-init <shell>`
- [ ] 1.7 Add `cd` to the Navigation section of the usage template in `main.go` so it is discoverable in `gws --help`
- [ ] 1.8 Create `cd_test.go` covering: path on stdout, informational line on stderr, `-q` suppressing stderr, uninitialized-workspace guidance, missing-directory error, and both branches of the TTY hint via the overridable indirection

### [ ] 2.0 Shell Integration Routing

Add a `cd` case to the `gws` function dispatch in all three shell templates so the parent shell actually changes directory, and extend the template tests to cover it.

#### 2.0 Proof Artifact(s)

- CLI: `./build/git-workspace shell-init zsh` output contains a `cd)` case invoking `{BIN} cd -q` demonstrates zsh routing
- CLI: `./build/git-workspace shell-init bash` and `shell-init powershell` show the equivalent blocks demonstrates all three shells are wired
- CLI: `eval "$(./build/git-workspace shell-init zsh)"; cd /tmp; gws cd; pwd` prints the workspace root demonstrates end-to-end directory change
- Test: `go test -run TestShellTemplates ./cmd/git-workspace/` passes demonstrates structural validation

#### 2.0 Tasks

- [ ] 2.1 Add a `cd)` case to `zshInitTemplate`'s dispatch: capture `_dest="$({BIN} cd -q 2>/dev/tty </dev/tty)"` and `[[ -n "$_dest" ]] && cd "$_dest"`, placed before the `-*` catch-all
- [ ] 2.2 In the same zsh case, pass through to the binary when `$2` is non-empty, so `gws cd foo` surfaces the `NoArgs` error instead of silently navigating to the root
- [ ] 2.3 Add the equivalent `cd)` case to `bashInitTemplate` using the local `dest` variable and the same redirection convention
- [ ] 2.4 Add the equivalent `'^cd$'` branch to `powershellInitTemplate`: invoke `& {BIN} cd -q`, and `Set-Location` on a non-empty result, matching how the existing navigation branches filter output
- [ ] 2.5 Confirm `print-workspace` remains in each template's passthrough subcommand list, unchanged
- [ ] 2.6 Add `TestShellTemplatesContainCdNavigation` to `shellinit_test.go`, table-driven over all three templates, asserting the `cd` dispatch and the `{BIN} cd` invocation are present
- [ ] 2.7 Run `go test ./cmd/git-workspace/` and confirm no regressions in existing template tests

### [ ] 3.0 Documentation Replacement

Retire the hand-written `cdgws`/`gcd` helper everywhere it appears and document `gws cd` in its place.

#### 3.0 Proof Artifact(s)

- Documentation: `grep -rn 'cdgws\|alias gcd' README.md docs/ cmd/ Makefile` returns nothing demonstrates the helper is fully retired
- Documentation: `gws cd` documented in commands-core.md, shell-integration.md, and README demonstrates the replacement is in place
- Documentation: an upgrade note in shell-integration.md explains the stale-function failure mode demonstrates it is documented
- CLI: `./build/git-workspace --help` lists `cd` under Navigation demonstrates in-tool discoverability

#### 3.0 Tasks

- [ ] 3.1 Remove the `function cdgws()` / `alias gcd=cdgws` lines from `rootCmd.Long` in `main.go` and reference `gws cd` in the shell-integration example instead
- [ ] 3.2 Add `gws cd` to the command list in `README.md` and delete the helper-function snippet from its shell-integration section
- [ ] 3.3 Add a `gws cd` section to `docs/site/commands-core.md` alongside the other navigation commands, documenting `-q` and the shell-integration requirement
- [ ] 3.4 Update `docs/site/shell-integration.md`: document `gws cd` in both the bash/zsh and PowerShell sections and remove the helper snippet
- [ ] 3.5 Update `docs/site/getting-started.md` to use `gws cd` wherever it demonstrates reaching the workspace root
- [ ] 3.6 Add a short note to the docs stating that `print-workspace` remains the primitive for scripts, while `gws cd` is the interactive command
- [ ] 3.7 Check `Makefile`'s shell-integration help output (around the `eval` lines) for helper references and update if present
- [ ] 3.8 Add an upgrade note to `docs/site/shell-integration.md`: the `gws` function is generated once at shell startup, so an upgraded binary is not enough — `shell-init` must be re-evaluated before newly added commands route. Note that `make build` writes only to `./build/` and that every local build reports `version dev`, so the commit line is what distinguishes binaries
