# Task 5.0 Proof Artifacts - Windows Install Documentation and Release Notes

## Documentation: `docs/site/getting-started.md`

Windows install section added with:

- `Invoke-WebRequest` download commands with version placeholder
- `curl.exe` alternative download commands
- Checksum verification via `checksums.txt`
- Install target: `$HOME\.local\bin\`
- Complete `$PROFILE` setup snippet (PATH + shell-init)
- `Add-Content` automated approach for profile setup
- Verification command: `& "$HOME\.local\bin\git-workspace.exe" --version`

## Documentation: `docs/site/shell-integration.md`

PowerShell section added with:

- Setup instructions for both PS 5.1 and 7+
- `$PROFILE` integration line: `Invoke-Expression (& git-workspace shell-init powershell | Out-String)`
- Full manual PowerShell function listing (matching `shell-init powershell` output)
- Mirrors existing bash/zsh section structure

## Config: `.goreleaser.yml`

Release footer updated with:

- Separate Linux/macOS and Windows (PowerShell) install sections
- Windows section includes `Invoke-WebRequest` download with versioned URL
- `$PROFILE` integration instructions included in release notes
- Uses goreleaser template variables (`{{ .Tag }}`, `{{ .Version }}`)

## Verification

- `docs/site/getting-started.md` contains both `Invoke-WebRequest` and `curl.exe` download methods
- `docs/site/getting-started.md` targets `$HOME\.local\bin\` as install location
- `docs/site/getting-started.md` includes complete `$PROFILE` snippet with PATH and shell-init
- `docs/site/shell-integration.md` has PowerShell setup section with PS 5.1/7+ note
- `.goreleaser.yml` footer includes Windows PowerShell install instructions
