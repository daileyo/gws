# 21 Questions Round 1 - Windows PowerShell Support

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. PowerShell Version Target

Which version(s) of PowerShell should we target?

- [ ] (A) PowerShell 5.1 only (Windows PowerShell, ships with Windows 10/11)
- [ ] (B) PowerShell 7+ only (cross-platform "pwsh", installed separately)
- [X] (C) Both PowerShell 5.1 and PowerShell 7+ (maximizes compatibility)
- [ ] (D) Other (describe)

## 2. Shell Function Name

The `gws` wrapper function on bash/zsh intercepts navigation commands to `cd` into repos. Should the PowerShell version use the same name?

- [X] (A) `gws` — same name as bash/zsh, consistent experience
- [ ] (B) A different name (describe)

## 3. Navigation / `cd` Behavior

The bash/zsh `gws` function captures stdout from `git-workspace` and `cd`s into the result. PowerShell can do the same via `Set-Location`. Should we replicate the full navigation behavior including:

- Bare `gws <repo>` navigates to repo
- `gws <repo> -p` / `gws parent <repo>` navigates to parent
- `gws <repo> -wt <branch>` navigates to worktree
- `gws <repo> -wt` lists worktrees for selection

- [X] (A) Yes, replicate all navigation behaviors listed above
- [ ] (B) Start with basic `gws <repo>` navigation only, add the rest later
- [ ] (C) Other (describe)

## 4. Tab Completion Approach

Cobra has built-in PowerShell completion generation (`git-workspace completion powershell`). Should we:

- [X] (A) Use Cobra's built-in PowerShell completion and wire it into the `shell-init powershell` output (consistent with bash/zsh approach)
- [ ] (B) Write custom PowerShell completion using `Register-ArgumentCompleter` for more control
- [ ] (C) Other (describe)

## 5. Install Method for Windows

For the initial release, you mentioned curl-based download by version. On Windows, `curl` is available as `curl.exe` (ships with Windows 10+) and also via `Invoke-WebRequest`. Should we provide:

- [ ] (A) PowerShell `Invoke-WebRequest` commands (most native to Windows users)
- [ ] (B) `curl.exe` commands (similar to Linux/macOS instructions)
- [X] (C) Both PowerShell and curl.exe instructions
- [ ] (D) Other (describe)

## 6. Install Location on Windows

Where should the binary be installed on Windows?

- [ ] (A) `$env:LOCALAPPDATA\Programs\git-workspace\` (follows Windows conventions)
- [X] (B) `$HOME\.local\bin\` (mirrors Linux/macOS XDG layout)
- [ ] (C) Let the user choose, but recommend a specific location
- [ ] (D) Other (describe)

## 7. PATH Management

Should the install instructions include guidance on adding the binary to PATH?

- [ ] (A) Yes, provide the PowerShell command to add to user PATH permanently
- [ ] (B) Just tell users to add it to PATH manually
- [X] (C) Other (describe) I'd like instructions for how to setup the user's $profile with the appropriate adds and what not.

## 8. Profile Integration

For `eval "$(git-workspace shell-init ...)"` equivalent in PowerShell, users need to add a line to their `$PROFILE`. Should we:

- [ ] (A) Output the integration code and tell users to add it to `$PROFILE` manually
- [X] (B) Provide instructions for both manual addition and an automated `Add-Content` command
- [ ] (C) Other (describe)

## 9. Windows-Specific Testing

Are there any Windows-specific concerns we should address in the binary itself (path separators, terminal encoding, etc.)?

- [X] (A) The binary should work on Windows as-is since Go handles cross-platform paths — just test the PowerShell integration layer
- [ ] (B) We should audit the codebase for Unix-specific assumptions (hardcoded `/`, shell-specific code, etc.)
- [ ] (C) Other (describe)
