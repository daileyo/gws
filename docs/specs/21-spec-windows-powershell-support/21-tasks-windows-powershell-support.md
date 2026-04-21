# 21 Tasks - Windows PowerShell Support

## Relevant Files

- `cmd/git-workspace/shellinit.go` - Contains shell-init command, zsh/bash templates, and `runShellInit` switch. PowerShell template and command updates go here.
- `cmd/git-workspace/shellinit_test.go` - Structural validation tests for shell templates. PowerShell template tests go here.
- `cmd/git-workspace/main.go` - Root command definition. No changes expected but useful reference for command structure.
- `.goreleaser.yml` - Release configuration including build targets, archive formats, and release footer. Windows install instructions go in the footer.
- `docs/site/getting-started.md` - Getting started guide. Needs Windows/PowerShell install section.
- `docs/site/shell-integration.md` - Shell integration docs. Needs PowerShell setup section.

### Notes

- Tests use the standard `testing` package with table-driven tests. Run with `go test -v ./cmd/git-workspace/`.
- Shell templates are Go `const` strings using `{BIN}` as a placeholder, replaced at runtime with `git-workspace`.
- The bash/zsh templates serve as the reference implementation — the PowerShell template must replicate the same routing logic using PowerShell syntax.
- Follow conventional commit messages (e.g., `feat(shell): ...`).
- goreleaser already builds `windows/amd64` and produces `.zip` archives.

## Tasks

### [x] 1.0 PowerShell `gws` Wrapper Function Template

Add a `powershellInitTemplate` const to `shellinit.go` that defines a PowerShell `gws` function with the same routing logic as the bash/zsh templates: subcommands delegate to the binary, navigation args capture stdout and call `Set-Location`, and flags like `-p`/`--parent` route through the parent subcommand. The template must use the `{BIN}` placeholder and be compatible with both PowerShell 5.1 and 7+.

#### 1.0 Proof Artifact(s)

- CLI: `go build ./cmd/git-workspace && ./build/git-workspace shell-init powershell` outputs a valid PowerShell script with `function gws` defined demonstrates template generation works
- CLI: In a PowerShell session, pasting the output defines the `gws` function and `Get-Command gws` returns it demonstrates the function is loadable

#### 1.0 Tasks

- [x] 1.1 Add `powershellInitTemplate` const to `shellinit.go` with a PowerShell `function gws` that handles the no-args case (calls `{BIN}` with no arguments, equivalent to `gws` → `gws list`)
- [x] 1.2 Add subcommand passthrough routing in the `gws` function using a PowerShell `switch` statement: delegate `list`, `init`, `add`, `refresh`, `print-workspace`, `tag`, `user`, `completion`, `shell-init`, `help`, and `__*` patterns directly to `& {BIN} @args`
- [x] 1.3 Add `worktree` subcommand routing: pass through `list`, `align`, `add`, and empty second arg to the binary; for all other worktree args, capture output with `-q` and call `Set-Location`
- [x] 1.4 Add parent navigation routing: handle `-p`, `--parent`, and `parent` as the first arg by calling `& {BIN} parent $secondArg -q` and `Set-Location` on the result
- [x] 1.5 Add `-*` flag passthrough: any first arg starting with `-` (other than `-p`/`--parent`) delegates to the binary
- [x] 1.6 Add default navigation case: capture `& {BIN} $repoName -q` output and `Set-Location`, with sub-cases for `-p`/`--parent` as second arg (parent nav) and `-wt` as second arg (worktree nav with optional third arg)
- [x] 1.7 Ensure all `{BIN}` invocations redirect stderr to the console (PowerShell equivalent of `2>/dev/tty`) so interactive prompts display correctly. Use `2>&1` or `$ErrorActionPreference` as needed for PS 5.1/7+ compat.

### [x] 2.0 PowerShell Tab Completion Integration

Wire Cobra's built-in `completion powershell` output into the PowerShell template so that tab completion works for both `git-workspace` and the `gws` alias. Register the completer for the `gws` function name in addition to `git-workspace`.

#### 2.0 Proof Artifact(s)

- CLI: `git-workspace shell-init powershell` output contains `Register-ArgumentCompleter` for both `git-workspace` and `gws` demonstrates completion is registered for both names
- CLI: In a PowerShell session with the function loaded, pressing Tab after `gws l` completes to `list` demonstrates completion works end-to-end

#### 2.0 Tasks

- [x] 2.1 Add a line to the PowerShell template that invokes Cobra's completion output: `& {BIN} completion powershell | Invoke-Expression` (this registers `Register-ArgumentCompleter` for `git-workspace`)
- [x] 2.2 Add a second `Register-ArgumentCompleter` call that copies the `git-workspace` completer to the `gws` function name, so tab completion works when the user types `gws` (not just `git-workspace`)
- [x] 2.3 Verify the completion script is compatible with PowerShell 5.1 (no features that require PS 7+)

### [x] 3.0 Update `shell-init` Command to Accept `powershell`

Update the `shellInitCmd` cobra command definition to accept `powershell` as a valid argument, update the `Use` string, help text, and `runShellInit` switch statement. Also consider accepting `pwsh` as an alias.

#### 3.0 Proof Artifact(s)

- CLI: `git-workspace shell-init powershell` succeeds (exit code 0) demonstrates the command accepts the new argument
- CLI: `git-workspace shell-init --help` shows `powershell` in the usage demonstrates discoverability
- CLI: `git-workspace shell-init invalid` returns an error listing all supported shells demonstrates error messaging is updated

#### 3.0 Tasks

- [x] 3.1 Update `shellInitCmd.Use` from `"shell-init [zsh|bash]"` to `"shell-init [zsh|bash|powershell]"`
- [x] 3.2 Update `shellInitCmd.ValidArgs` to include `"powershell"`
- [x] 3.3 Update `shellInitCmd.Long` help text to add PowerShell profile instructions alongside the existing bash/zsh examples (e.g., `Invoke-Expression (& git-workspace shell-init powershell)` added to `$PROFILE`)
- [x] 3.4 Add `case "powershell":` to the `runShellInit` switch statement, setting `tmpl = powershellInitTemplate`
- [x] 3.5 Update the error message in the `default` case from `"supported: zsh, bash"` to `"supported: zsh, bash, powershell"`
- [x] 3.6 (Optional) Add `case "pwsh":` as an alias that falls through to the `"powershell"` case

### [x] 4.0 PowerShell Template Tests

Extend `shellinit_test.go` with the same structural validation tests for the PowerShell template: subcommand routing, `{BIN}` placeholder, navigation fallthrough, parent navigation, and worktree navigation patterns (using PowerShell syntax equivalents).

#### 4.0 Proof Artifact(s)

- Test: `go test -v -run TestShell ./cmd/git-workspace/` passes with all PowerShell template tests demonstrates template correctness
- CLI: `go test -v ./cmd/git-workspace/` passes all tests (no regressions) demonstrates no existing functionality is broken

#### 4.0 Tasks

- [x] 4.1 Add `{"powershell", powershellInitTemplate}` to the `templates` slice in `TestShellTemplatesRouteSubcommands` so subcommand routing is validated (note: PowerShell uses `switch` not `case`, so the Contains checks may need PowerShell-specific patterns — see 4.5)
- [x] 4.2 Add a `{BIN}` placeholder check for `powershellInitTemplate` in `TestShellTemplatesContainBinPlaceholder`
- [x] 4.3 Add `powershellInitTemplate` to `TestShellTemplatesContainNavigationFallthrough` — check for PowerShell-equivalent navigation pattern (e.g., `Set-Location` and `& {BIN}`)
- [x] 4.4 Add `powershellInitTemplate` to `TestShellTemplatesContainWorktreeNavigation` — check for `-wt` flag detection, `--worktree` with branch argument, bare `--worktree` invocation, `worktree` case block, and worktree subcommand passthrough using PowerShell syntax patterns
- [x] 4.5 Add `powershellInitTemplate` to `TestShellTemplatesContainParentNavigation` — check for `-p`, `--parent`, `parent` keyword routing, and `{BIN} parent` invocation. If the existing string-match patterns don't work for PowerShell syntax (e.g., different quoting), add PowerShell-specific sub-tests alongside the bash/zsh ones rather than forcing identical patterns.
- [x] 4.6 Run `go test -v ./cmd/git-workspace/` and confirm all tests pass with no regressions

### [x] 5.0 Windows Install Documentation and Release Notes

Create a Windows install guide with versioned download instructions using both `Invoke-WebRequest` and `curl.exe`, targeting `$HOME\.local\bin\`. Include complete `$PROFILE` setup instructions (PATH addition, `shell-init powershell` integration, checksum verification). Update the `.goreleaser.yml` release footer with Windows-specific instructions.

#### 5.0 Proof Artifact(s)

- Documentation: Windows install section in `docs/site/getting-started.md` with both download methods, PATH setup, and `$PROFILE` integration demonstrates install guide is complete
- Documentation: PowerShell section in `docs/site/shell-integration.md` with setup and usage examples demonstrates shell integration is documented
- Config: `.goreleaser.yml` release footer includes Windows install section demonstrates release notes will include Windows instructions
- CLI: Following the documented steps on a Windows machine results in `gws --version` working demonstrates instructions are accurate

#### 5.0 Tasks

- [x] 5.1 Add a "Windows (PowerShell)" install section to `docs/site/getting-started.md` with `Invoke-WebRequest` download commands targeting `$HOME\.local\bin\git-workspace.exe`, including version placeholder and checksum verification via `checksums.txt`
- [x] 5.2 Add equivalent `curl.exe` download commands to the same section as an alternative method
- [x] 5.3 Add `$PROFILE` setup instructions to the getting-started guide: a complete example snippet showing PATH addition (`$env:Path` update) and `Invoke-Expression (& git-workspace shell-init powershell)`, with both manual and `Add-Content` approaches
- [x] 5.4 Add a "PowerShell" section to `docs/site/shell-integration.md` covering setup, navigation examples (bare, parent, worktree), and tab completion — mirroring the structure of the existing bash/zsh content
- [x] 5.5 Update the `.goreleaser.yml` release footer to add Windows-specific install instructions (PowerShell download, extract zip, add to PATH) alongside the existing Linux/macOS section
