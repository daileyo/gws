# Getting Started

## Installation

### Install via Homebrew

```bash
brew install daileyo/gws/git-workspace
```

### Install on Windows (PowerShell)

Download the latest release and extract to `$HOME\.local\bin`:

**Using `Invoke-WebRequest`:**

```powershell
# Set the version you want to install
$VERSION = "2.19.1"

# Download and extract
$url = "https://github.com/daileyo/gws/releases/download/v$VERSION/git-workspace_${VERSION}_windows_amd64.zip"
$zip = "$env:TEMP\git-workspace.zip"
Invoke-WebRequest -Uri $url -OutFile $zip
New-Item -ItemType Directory -Force -Path "$HOME\.local\bin" | Out-Null
Expand-Archive -Path $zip -DestinationPath "$HOME\.local\bin" -Force
Remove-Item $zip
```

**Using `curl.exe`:**

```powershell
# Set the version you want to install
$VERSION = "2.19.1"

# Download and extract
curl.exe -Lo "$env:TEMP\git-workspace.zip" "https://github.com/daileyo/gws/releases/download/v$VERSION/git-workspace_${VERSION}_windows_amd64.zip"
New-Item -ItemType Directory -Force -Path "$HOME\.local\bin" | Out-Null
Expand-Archive -Path "$env:TEMP\git-workspace.zip" -DestinationPath "$HOME\.local\bin" -Force
Remove-Item "$env:TEMP\git-workspace.zip"
```

**Verify checksum (optional):**

```powershell
curl.exe -Lo "$env:TEMP\checksums.txt" "https://github.com/daileyo/gws/releases/download/v$VERSION/checksums.txt"
(Get-FileHash "$HOME\.local\bin\git-workspace.exe" -Algorithm SHA256).Hash
Get-Content "$env:TEMP\checksums.txt" | Select-String "windows_amd64"
```

**Verify installation:**

```powershell
& "$HOME\.local\bin\git-workspace.exe" --version
```

### Build from Source

**Linux / macOS:**

```bash
git clone https://github.com/daileyo/gws.git
cd gws
make build
# The binary will be in ./build/git-workspace
# Optionally, install to your PATH
make install
```

**Windows (PowerShell):**

```powershell
git clone https://github.com/daileyo/gws.git
cd gws

# Build and install to ~/.local/bin in one step
New-Item -ItemType Directory -Force -Path "$HOME\.local\bin" | Out-Null
go build -o "$HOME\.local\bin\git-workspace.exe" ./cmd/git-workspace
```

> **Note:** The Makefile requires a Unix shell (`bash`/`zsh`). On Windows, use `go build` directly as shown above.

## Shell Integration

### Linux / macOS (bash / zsh)

Add these lines to your `~/.zshrc` (or `~/.bashrc`):

```bash
export PATH="$HOME/.local/bin:$PATH"
eval "$(git-workspace shell-init zsh)"   # or: shell-init bash
```

### Windows (PowerShell)

Add the following to your PowerShell `$PROFILE`. To find your profile path, run `echo $PROFILE` in PowerShell.

```powershell
# Add git-workspace to PATH
$env:Path = "$HOME\.local\bin;$env:Path"

# Set up the gws function and tab completion
Invoke-Expression (& git-workspace shell-init powershell | Out-String)
```

**To add automatically via `Add-Content`:**

```powershell
# Create profile if it doesn't exist
if (!(Test-Path -Path $PROFILE)) { New-Item -ItemType File -Path $PROFILE -Force | Out-Null }

# Append shell integration
Add-Content -Path $PROFILE -Value "`n# git-workspace shell integration"
Add-Content -Path $PROFILE -Value '$env:Path = "$HOME\.local\bin;$env:Path"'
Add-Content -Path $PROFILE -Value 'Invoke-Expression (& git-workspace shell-init powershell | Out-String)'
```

Restart your PowerShell session (or run `. $PROFILE`) to activate.

This gives you the `gws` function with tab completion. See [Shell Integration](shell-integration.md) for full details.

## Quick Start

### 1. Initialize a Workspace

Scan a directory for git repositories:

```bash
# Initialize in current directory
gws init

# Initialize in a specific directory
gws init ~/projects

# Initialize with absolute path
gws init /path/to/your/workspace
```

This command will:

- Recursively scan the directory for git repositories
- Extract repository metadata (name, path, remote URL)
- Detect repository type and user configuration
- Save the configuration to `~/.gws/config.json`

### 2. List Your Repositories

Once initialized, run `gws` with no arguments to see your repositories:

```bash
gws
```

Output (compact multi-column names):

```
Found 15 repositories:

my-project       work-api         client-site      my-api
frontend-app     docs-site        infra-tools      backend-svc
mobile-app       shared-libs      auth-service     data-pipeline
ml-models        cli-tools        test-harness
```

### 3. View Repository Details

Use verbose mode to see more information:

```bash
gws list -v
```

Output:

```
Found 15 repositories:

NAME              TYPE       VISIBILITY  TAGS              PATH
----------------  ---------  ----------  ----------------  ------------------------------------
my-project        github     private     personal, web     /home/user/projects/my-project
work-api          gitlab     private     work              /home/user/projects/work-api
client-site       bitbucket  unknown     client, archived  /home/user/projects/client-site
```

### 4. Check Version

```bash
gws --version
```

Output:

```
git-workspace version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
```

## Next Steps

- [Core Commands](commands-core.md) — Listing, filtering, refreshing, and more
- [User Management](commands-user.md) — Managing git user profiles
- [Tagging](commands-tagging.md) — Organizing repositories with custom tags
- [Shell Integration](shell-integration.md) — Navigation and tab completion
- [Configuration](configuration.md) — Config file structure and fields
