# 21-spec-windows-powershell-support

## Introduction/Overview

Git-workspace (`gws`) currently supports shell integration for bash and zsh on Linux/macOS, but has no PowerShell support for Windows users. While the Go binary already cross-compiles for Windows via goreleaser, there is no `shell-init powershell` command, no PowerShell wrapper function, and no documented install path for Windows users.

This spec adds full PowerShell shell integration — including the `gws` wrapper function, tab completion, and navigation — along with documented install instructions using curl and PowerShell download commands.

## Goals

- Add `powershell` as a supported shell in the `shell-init` command, producing a PowerShell script with equivalent functionality to the existing bash/zsh templates
- Provide tab completion for all commands, subcommands, flags, and repository names in PowerShell using Cobra's built-in `completion powershell` output
- Replicate all navigation behaviors (`gws <repo>`, parent navigation, worktree navigation) in the PowerShell wrapper function
- Support both PowerShell 5.1 (Windows PowerShell) and PowerShell 7+ (pwsh)
- Provide documented install instructions with both `Invoke-WebRequest` and `curl.exe` download commands

## User Stories

- **As a Windows developer**, I want to install git-workspace and use `gws` in PowerShell so that I have the same repo navigation experience as my colleagues on macOS/Linux.
- **As a PowerShell user**, I want tab completion for `gws` commands and repository names so that I can discover and navigate repos efficiently.
- **As a Windows developer**, I want clear install instructions so that I can get git-workspace running without needing to figure out PATH and profile setup on my own.

## Demoable Units of Work

### Unit 1: PowerShell Shell-Init Template

**Purpose:** Add the PowerShell `gws` wrapper function and completion wiring to `shell-init`, giving Windows users feature parity with bash/zsh.

**Functional Requirements:**
- The `shell-init` command shall accept `powershell` as a valid argument (in addition to `zsh` and `bash`)
- The PowerShell template shall define a `gws` function that delegates subcommands (`list`, `init`, `add`, `refresh`, `print-workspace`, `tag`, `user`, `completion`, `shell-init`, `worktree`, `help`) directly to the `git-workspace` binary
- The `gws` function shall handle bare `gws <repo>` navigation by capturing stdout from `git-workspace` and calling `Set-Location` on the result
- The `gws` function shall handle `gws <repo> -p` and `gws parent <repo>` navigation (parent directory)
- The `gws` function shall handle `gws <repo> -wt <branch>` navigation (specific worktree) and `gws <repo> -wt` (worktree selection)
- The template shall source Cobra's `completion powershell` output to provide tab completion for `git-workspace` and register it for the `gws` function alias
- The template shall be compatible with both PowerShell 5.1 and PowerShell 7+
- The `shellinit_test.go` tests shall be updated to cover the new `powershell` template with the same pattern validation as bash/zsh (subcommand routing, navigation fallthrough, worktree navigation, parent navigation)

**Proof Artifacts:**
- CLI: `git-workspace shell-init powershell` outputs valid PowerShell script demonstrates template generation works
- CLI: PowerShell session with `gws` function loaded demonstrates wrapper function works
- CLI: Tab-completing `gws l` shows `list` demonstrates completion works
- CLI: `gws <repo-name>` changes directory demonstrates navigation works
- Test: `shellinit_test.go` passes with powershell template validation demonstrates template correctness

### Unit 2: Windows Install Documentation

**Purpose:** Provide Windows users with clear, versioned install instructions using both PowerShell-native and curl download methods, plus `$PROFILE` integration guidance.

**Functional Requirements:**
- The documentation shall provide `Invoke-WebRequest` commands to download a specific version of the `git-workspace` Windows zip from GitHub releases
- The documentation shall provide equivalent `curl.exe` commands for users who prefer that approach
- The install instructions shall target `$HOME\.local\bin\` as the recommended binary location (mirroring the Linux/macOS XDG layout)
- The documentation shall include instructions for adding `$HOME\.local\bin` to the user's PATH via `$PROFILE`
- The documentation shall include the `shell-init powershell` invocation line to add to `$PROFILE` (manual addition and `Add-Content` command)
- The documentation shall include a complete example `$PROFILE` snippet showing PATH setup and shell-init integration together
- The goreleaser release footer shall be updated to include Windows-specific install instructions alongside the existing Linux/macOS instructions

**Proof Artifacts:**
- Documentation: Install guide with PowerShell and curl.exe download commands demonstrates clear install path exists
- CLI: Following the install instructions on a Windows machine results in a working `gws` command demonstrates end-to-end install works
- Config: Updated `.goreleaser.yml` release footer demonstrates release notes include Windows instructions

## Non-Goals (Out of Scope)

1. **Windows package manager support**: No Chocolatey, Scoop, or winget integration in this spec — that is planned for a future spec
2. **cmd.exe support**: Only PowerShell is targeted; no support for the legacy Windows Command Prompt
3. **Auditing the Go codebase for Windows path issues**: The binary is assumed to work on Windows as-is since Go handles cross-platform paths; any bugs found will be addressed separately
4. **Automated installer or MSI**: This spec covers manual download-and-extract installation only
5. **Fish shell support**: While Cobra supports fish completion, adding fish is out of scope for this spec

## Design Considerations

No specific design requirements identified. The PowerShell function output follows the same pattern as the existing bash/zsh templates — it is plain text output from `shell-init` that the user evals in their profile.

## Repository Standards

- Go source files follow standard `gofmt` formatting
- Tests use the standard `testing` package with table-driven tests where applicable
- Shell templates are defined as Go `const` strings in `shellinit.go`
- Template tests validate structural properties (subcommand routing, placeholder presence, navigation patterns) rather than exact string matching
- Cobra is used for CLI structure; built-in completion generation is preferred over custom solutions
- goreleaser handles cross-platform builds and release artifact packaging
- Conventional commits are used for commit messages

## Technical Considerations

- The PowerShell template must work in both PowerShell 5.1 (uses `powershell.exe`, no `$PSNativeCommandArgumentPassing`) and PowerShell 7+ (uses `pwsh.exe`, different argument handling). Test both if possible.
- Cobra's `completion powershell` generates a script block that registers completions via `Register-ArgumentCompleter`. This registers for `git-workspace`; the template must also register the same completer for the `gws` alias.
- PowerShell does not have `eval` — the equivalent is `Invoke-Expression` or dot-sourcing. The `$PROFILE` integration should use `Invoke-Expression (& git-workspace shell-init powershell)` or a similar pattern.
- The `{BIN}` placeholder pattern used in bash/zsh templates should be reused in the PowerShell template for consistency.
- Interactive selection (for multiple repo matches or worktree selection) relies on stdin/stderr, which should work the same way in PowerShell since the Go binary handles the I/O directly.

## Security Considerations

- Download instructions should include checksum verification using the `checksums.txt` file from the GitHub release
- The `$PROFILE` modification instructions should warn users to review the shell-init output before adding it to their profile
- No credentials or API keys are involved in this feature

## Success Metrics

1. **Feature parity**: `git-workspace shell-init powershell` produces a working `gws` function with all navigation modes (bare, parent, worktree)
2. **Tab completion**: All subcommands, flags, and repository names complete correctly in PowerShell
3. **Install success**: A Windows user can follow the install documentation to go from zero to working `gws` in under 5 minutes
4. **Test coverage**: All existing `shellinit_test.go` test patterns extended to cover the PowerShell template

## Open Questions

1. Should `shell-init` also accept `pwsh` as an alias for `powershell` to match the executable name on PowerShell 7+?
2. Does Cobra's PowerShell completion handle the `gws` alias automatically, or do we need to explicitly duplicate the `Register-ArgumentCompleter` call for both `git-workspace` and `gws`?
