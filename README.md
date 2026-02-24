# git-workspace

[![CI](https://github.com/daileyo/gws/actions/workflows/ci.yml/badge.svg)](https://github.com/daileyo/gws/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/daileyo/gws)](https://goreportcard.com/report/github.com/daileyo/gws)
[![Snyk Security](https://snyk.io/test/github/daileyo/gws/badge.svg)](https://snyk.io/test/github/daileyo/gws)
[![Release](https://img.shields.io/github/v/release/daileyo/gws)](https://github.com/daileyo/gws/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Repository Discovery**: Automatically find all git repositories in a directory tree
- **Automatic Classification**: Detect repository type (GitHub, GitLab, Azure DevOps, Bitbucket) from remote URLs
- **Git Status Integration**: View branch, clean/dirty state, and ahead/behind indicators at a glance
- **Smart Caching**: Fast status display with configurable cache (5-minute TTL)
- **Custom Tagging**: Organize repositories with custom tags (personal, work, archived, etc.)
- **Advanced Filtering**: Search and filter repositories by type, tags, name, or path
- **Repository Navigation**: Jump to any repository instantly with `gws my-repo` shell function
- **Workspace Management**: Track and organize repositories in a centralized configuration
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Lightweight**: Single binary with no external dependencies

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/daileyo/gws.git
cd gws

# Build the binary
make build

# The binary will be in ./build/git-workspace
# Optionally, install to your PATH
make install
```

## Quick Start

### 1. Initialize a Workspace

Initialize git-workspace by scanning a directory for git repositories:

```bash
# Initialize in current directory
git-workspace --init .

# Initialize in a specific directory
git-workspace --init ~/projects

# Initialize with absolute path
git-workspace --init /path/to/your/workspace
```

This command will:
- Recursively scan the directory for git repositories
- Extract repository metadata (name, path, remote URL)
- Save the configuration to `~/.gws/config.json`

### 2. View Workspace Information

Once initialized, run `git-workspace` to see your workspace information:

```bash
git-workspace
```

Output:
```
Workspace: /home/user/projects
Repositories: 15

Use 'git-workspace --help' to see available commands
```

### 3. List Repositories

View all tracked repositories with their metadata:

```bash
git-workspace --list
```

Output:
```
Found 15 repositories:

NAME           TYPE      VISIBILITY  TAGS              PATH
----           ----      ----------  ----              ----
my-project     github    private     personal, web     /home/user/projects/my-project
work-api       gitlab    private     work              /home/user/projects/work-api
client-site    bitbucket unknown     client, archived  /home/user/projects/client-site
```

**Show Git Status:**

```bash
git-workspace --list --status
# or
git-workspace -l -s
```

Output:
```
Found 15 repositories:

NAME           STATUS         TYPE      VISIBILITY  TAGS              PATH
----           ------         ----      ----------  ----              ----
my-project     main ✓         github    private     personal, web     /home/user/projects/my-project
work-api       develop ✗ ↑2   gitlab    private     work              /home/user/projects/work-api
client-site    main ✓ ↓1      bitbucket unknown     client, archived  /home/user/projects/client-site
```

**Status Indicators:**
- `✓` = Clean working tree (no uncommitted changes)
- `✗` = Dirty working tree (uncommitted changes)
- `↑N` = N commits ahead of remote
- `↓N` = N commits behind remote

### 4. Check Version

```bash
git-workspace --version
```

Output:
```
git-workspace version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
```

## Repository Classification and Tagging

### Automatic Classification

When repositories are discovered with `git-workspace --init`, they are automatically classified based on their remote URL:

| Repository Type | Detected From |
|----------------|---------------|
| **GitHub** | `github.com` |
| **GitLab** | `gitlab.com`, `gitlab.*` |
| **Azure DevOps (ADO)** | `dev.azure.com`, `visualstudio.com` |
| **Bitbucket** | `bitbucket.org` |
| **Unknown** | Other or no remote URL |

#### Supported Remote URL Patterns

**GitHub:**
- HTTPS: `https://github.com/user/repo.git`
- SSH: `git@github.com:user/repo.git`

**GitLab:**
- HTTPS: `https://gitlab.com/user/repo.git`
- SSH: `git@gitlab.com:user/repo.git`
- Self-hosted: `https://gitlab.company.com/user/repo.git`

**Azure DevOps:**
- HTTPS: `https://dev.azure.com/org/project/_git/repo`
- HTTPS (legacy): `https://org.visualstudio.com/project/_git/repo`
- SSH: `git@ssh.dev.azure.com:v3/org/project/repo`

**Bitbucket:**
- HTTPS: `https://bitbucket.org/user/repo.git`
- SSH: `git@bitbucket.org:user/repo.git`

### Visibility Detection

Repository visibility is inferred from the remote URL protocol:

- **Private**: SSH URLs (`git@...` or `ssh://...`) typically require authentication
- **Unknown**: HTTPS URLs could be public or private (requires authentication check)

### Custom Tags

Add custom tags to organize repositories. Tags are applied to **all repositories** matching the given name:

```bash
# Add a tag to all repos matching "my-project" (partial match)
git-workspace --add-tag my-project personal

# Add a tag to all repos with "api" in the name
git-workspace --add-tag api backend

# Remove a tag from all matching repos
git-workspace --remove-tag my-project personal

# Tags can be anything: personal, work, client, archived, production, etc.
```

**How Tag Matching Works:**
- Matches by **partial name** (case-insensitive): `--add-tag api work` tags "my-api", "api-gateway", etc.
- Tags are applied to **all matching repositories**
- You'll see a summary of how many repos were tagged

**Examples:**
```bash
# Tag all API services as backend
git-workspace --add-tag api backend
# Output: Added tag 'backend' to 3 repositories
```

### Filtering Repositories

Filter repositories by type, tags, name, or path. All filters can be combined:

```bash
# Filter by repository type
git-workspace --list --type github

# Filter by single tag
git-workspace --list --tag personal

# Filter by multiple tags (AND logic - repo must have ALL tags)
git-workspace --list --tag work --tag backend

# Filter by repository name (partial match, case-insensitive)
git-workspace --list --name project

# Filter by repository path (partial match)
git-workspace --list --path /home/user/projects

# Combine multiple filters
git-workspace --list --type gitlab --tag work --name api
```

### Output Formats

Control how repositories are displayed:

```bash
# Default table format
git-workspace --list

# JSON format for scripting/automation
git-workspace --list --output json
git-workspace --list -o json

# JSON with filters
git-workspace --list --type github -o json
```

**JSON Output Example:**
```json
[
  {
    "name": "my-project",
    "path": "/home/user/projects/my-project",
    "remote_url": "https://github.com/user/my-project.git",
    "type": "github",
    "visibility": "unknown",
    "tags": ["personal", "web"]
  }
]
```

### Refresh Workspace

Update repository metadata and clear status cache:

```bash
git-workspace --refresh
```

**What it does:**
- Re-scans workspace for new repositories
- Updates remote URLs and classification
- Clears and rebuilds git status cache
- Preserves all custom tags

**When to use:**
- After adding new repositories to your workspace
- When remote URLs have changed
- To force update of cached git status
- After bulk repository operations

**Example output:**
```
Refreshing workspace at: /home/user/projects
Re-scanning for repositories...
Cleared git status cache

Refresh complete!
Total repositories: 15
New repositories found: 2
```

### Shell Navigation

Navigate to your workspace root with shell integration:

**Setup (add to ~/.bashrc or ~/.zshrc):**

```bash
# Navigate to workspace root
function cdgws() {
  cd "$(git-workspace --print-workspace)"
}

# Short alias
alias gcd=cdgws
```

**Usage:**

```bash
# Navigate to workspace from anywhere
cdgws
# or
gcd

# You'll be in your workspace root directory
pwd
# Output: /home/user/projects
```

**Print workspace path:**

```bash
git-workspace --print-workspace
# Output: /home/user/projects
```

### Repository Navigation

Jump to any tracked repository by name using the `gws` shell function.

**Setup (add to ~/.bashrc or ~/.zshrc):**

```bash
# 'gws': navigate to a repo by name, or pass flags through to git-workspace
function gws() {
  if [[ "$1" == -* ]]; then
    git-workspace "$@"
  else
    cd "$(git-workspace -g "$1" -q)"
  fi
}
```

**Usage:**

```bash
# Navigate to a repository by name — changes your directory
gws my-repo

# Pass any flag through to git-workspace as normal
gws --list
gws -l --tag personal --status
gws --refresh
```

**Wildcard matching:**

```bash
# Wildcards work too (* = zero or more, ? = single character)
gws "api-*"
gws "?rontend"
```

**Multiple matches:**

When multiple repositories match, git-workspace displays a numbered list for selection:

```
Multiple repositories match 'api':

  1) my-api (github) /home/user/projects/my-api
  2) my-api-v2 (github) /home/user/projects/my-api-v2
  3) work-api (gitlab) /home/user/projects/work-api

Select repository [1-3]:
```

When piped (non-TTY), all matching paths are printed without prompting.

**No match suggestions:**

When no repositories match, git-workspace suggests similar names:

```
No repositories found matching 'aip'

Did you mean:
  my-api
  work-api
```

**Using the binary directly:**

```bash
# Print path without changing directory (useful for scripting)
git-workspace -g my-repo -q
# Output: /home/user/projects/my-repo

# Verbose mode (default) — details to stderr, path to stdout
git-workspace -g my-repo
# stderr: my-repo (github) → /home/user/projects/my-repo
# stdout: /home/user/projects/my-repo
```

## Manual Verification Steps

### Verify Repository Discovery

1. **Create a test workspace** with multiple git repositories:
   ```bash
   mkdir -p /tmp/test-workspace
   cd /tmp/test-workspace

   # Create test repositories
   git clone https://github.com/some/repo1.git
   git clone https://github.com/some/repo2.git
   mkdir subdir
   git clone https://github.com/some/repo3.git subdir/repo3
   ```

2. **Initialize git-workspace** in the test workspace:
   ```bash
   git-workspace --init /tmp/test-workspace
   ```

3. **Verify the output** shows all 3 repositories discovered:
   ```
   Workspace initialized successfully!
   Found 3 git repositories:

     - repo1
       https://github.com/some/repo1.git
     - repo2
       https://github.com/some/repo2.git
     - repo3
       https://github.com/some/repo3.git

   Configuration saved to: /home/user/.gws/config.json
   ```

4. **Inspect the configuration file**:
   ```bash
   cat ~/.gws/config.json
   ```

   Expected structure:
   ```json
   {
     "version": "1.0.0",
     "workspace": "/tmp/test-workspace",
     "repositories": [
       {
         "name": "repo1",
         "path": "/tmp/test-workspace/repo1",
         "remote_url": "https://github.com/some/repo1.git",
         "tags": []
       },
       ...
     ]
   }
   ```

### Verify Error Handling

1. **Test without initialization**:
   ```bash
   # Remove config to simulate uninitialized state
   rm ~/.gws/config.json

   # Run git-workspace command
   git-workspace
   ```

   Expected output:
   ```
   Error: workspace not initialized

   To get started, initialize a workspace:
     git-workspace --init <directory>

   Example:
     git-workspace --init ~/projects
   ```

2. **Test with invalid directory**:
   ```bash
   git-workspace --init /nonexistent/directory
   ```

   Expected output:
   ```
   Error: failed to scan workspace: failed to access directory: ...
   ```

### Verify Nested Repository Discovery

1. **Create nested repository structure**:
   ```bash
   mkdir -p /tmp/nested-test/parent
   cd /tmp/nested-test/parent
   git init

   mkdir child
   cd child
   git init
   ```

2. **Initialize and verify both are found**:
   ```bash
   git-workspace --init /tmp/nested-test
   ```

   Expected: Both `parent` and `child` repositories should be discovered.

### Verify Skipped Directories

1. **Create repositories in directories that should be skipped**:
   ```bash
   mkdir -p /tmp/skip-test/node_modules/repo1
   cd /tmp/skip-test/node_modules/repo1
   git init

   mkdir -p /tmp/skip-test/valid-repo
   cd /tmp/skip-test/valid-repo
   git init
   ```

2. **Initialize and verify only valid-repo is found**:
   ```bash
   git-workspace --init /tmp/skip-test
   ```

   Expected: Only `valid-repo` should be found (node_modules should be skipped).

## Configuration

git-workspace stores its configuration in `~/.gws/config.json`. The configuration includes:

- **version**: Config file format version
- **workspace**: Root directory of the workspace
- **repositories**: Array of discovered repositories

### Configuration Example

```json
{
  "version": "1.0.0",
  "workspace": "/home/user/projects",
  "repositories": [
    {
      "name": "gws",
      "path": "/home/user/projects/gws",
      "remote_url": "https://github.com/daileyo/gws.git",
      "type": "github",
      "visibility": "unknown",
      "tags": ["personal", "go"]
    },
    {
      "name": "my-api",
      "path": "/home/user/projects/my-api",
      "remote_url": "git@gitlab.com:user/my-api.git",
      "type": "gitlab",
      "visibility": "private",
      "tags": ["work", "backend"]
    }
  ]
}
```

## Development

### Build Commands

```bash
# Build the binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/config
go test -v ./internal/discovery
```

## Project Structure

```
.
├── cmd/
│   └── git-workspace/    # Main application entry point
│       ├── main.go       # Root command and CLI setup
│       ├── init.go       # Init command implementation
│       ├── list.go       # List command with status display
│       ├── navigate.go   # Navigation implementation
│       ├── refresh.go    # Refresh command implementation
│       ├── tag.go        # Tag command implementation
│       └── untag.go      # Untag command implementation
├── internal/
│   ├── classifier/       # Repository classification
│   │   ├── detector.go
│   │   └── detector_test.go
│   ├── config/           # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── discovery/        # Repository discovery
│   │   ├── scanner.go
│   │   └── scanner_test.go
│   ├── filter/           # Repository filtering logic
│   │   ├── filter.go
│   │   └── filter_test.go
│   └── git/              # Git status integration
│       ├── cache.go
│       ├── status.go
│       └── status_test.go
├── docs/
│   └── specs/            # Specification documents
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── README.md             # This file
```

## Roadmap

- [x] Repository discovery and workspace initialization
- [x] Repository classification (GitHub, GitLab, Azure DevOps, Bitbucket)
- [x] Manual tagging for organization
- [x] Advanced search and filter capabilities (type, tags, name, path)
- [x] Multiple output formats (table, JSON)
- [x] Git status integration with visual indicators
- [x] Workspace refresh and cache management
- [x] Repository navigation with `gws` shell function
- [x] CI/CD pipeline with automated releases
- [ ] Plugin support for editors (Neovim, VSCode)

## Release Process

This project uses automated releases powered by:
- **[Release Please](https://github.com/googleapis/release-please)**: Automates versioning and changelog generation based on [Conventional Commits](https://www.conventionalcommits.org/)
- **[GoReleaser](https://goreleaser.com/)**: Builds cross-platform binaries and publishes to GitHub Releases
- **[Snyk](https://snyk.io/)**: Security scanning for dependencies

### Versioning

We follow [Semantic Versioning](https://semver.org/) (SemVer):
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for bug fixes (backwards compatible)

### How Releases Work

1. **Commits to `main`** trigger Release Please to create/update a Release PR
2. **Merging the Release PR** creates a new GitHub Release with a version tag
3. **The version tag** triggers GoReleaser to build and publish binaries

### Downloading Releases

Pre-built binaries are available on the [Releases page](https://github.com/daileyo/gws/releases) for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

```bash
# Example: Download and install on Linux/macOS
curl -LO https://github.com/daileyo/gws/releases/latest/download/git-workspace_<version>_<os>_<arch>.tar.gz
tar -xzf git-workspace_*.tar.gz
sudo mv git-workspace /usr/local/bin/
```

## Security

Security scanning is performed on every build using Snyk:
- **High severity** vulnerabilities block the build
- **Medium/Low severity** vulnerabilities are logged as warnings

To report a security vulnerability, please open an issue or contact the maintainers directly.

## License

MIT

## Contributing

Contributions are welcome! Please follow these guidelines:

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat:` New feature (triggers MINOR version bump)
- `fix:` Bug fix (triggers PATCH version bump)
- `docs:` Documentation only
- `style:` Code style (formatting, semicolons, etc.)
- `refactor:` Code refactoring
- `perf:` Performance improvement
- `test:` Adding or updating tests
- `chore:` Maintenance tasks
- `ci:` CI/CD changes

**Breaking Changes:**
Add `!` after the type or include `BREAKING CHANGE:` in the footer to trigger a MAJOR version bump:
```
feat!: remove deprecated API
```

### Development Workflow

1. **Fork** the repository
2. **Create a branch** for your feature: `git checkout -b feat/my-feature`
3. **Make changes** and commit using conventional commits
4. **Run tests**: `make test`
5. **Run linter**: `make lint`
6. **Push** to your fork and create a Pull Request

### Running CI Locally

```bash
# Run all CI checks
make ci

# Individual checks
make vet      # Run go vet
make lint     # Run golangci-lint
make test     # Run tests
make coverage # Run tests with coverage report

# Build snapshot release (for local testing)
make snapshot
```

### Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Ensure `make lint` passes with no errors
