# gws - Git Workspace

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Repository Discovery**: Automatically find all git repositories in a directory tree
- **Automatic Classification**: Detect repository type (GitHub, GitLab, Azure DevOps, Bitbucket) from remote URLs
- **Custom Tagging**: Organize repositories with custom tags (personal, work, archived, etc.)
- **Advanced Filtering**: Search and filter repositories by type, tags, or name
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

# The binary will be in ./build/gws
# Optionally, install to your PATH
make install
```

## Quick Start

### 1. Initialize a Workspace

Initialize gws by scanning a directory for git repositories:

```bash
# Initialize in current directory
gws init .

# Initialize in a specific directory
gws init ~/projects

# Initialize with absolute path
gws init /path/to/your/workspace
```

This command will:
- Recursively scan the directory for git repositories
- Extract repository metadata (name, path, remote URL)
- Save the configuration to `~/.gws/config.json`

### 2. View Workspace Information

Once initialized, run `gws` to see your workspace information:

```bash
gws
```

Output:
```
Workspace: /home/user/projects
Repositories: 15

Use 'gws --help' to see available commands
```

### 3. List Repositories

View all tracked repositories with their metadata:

```bash
gws list
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

### 4. Check Version

```bash
gws version
```

Output:
```
gws version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
```

## Repository Classification and Tagging

### Automatic Classification

When repositories are discovered with `gws init`, they are automatically classified based on their remote URL:

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

Add custom tags to organize repositories:

```bash
# Add a tag
gws tag my-project personal
gws tag work-api work

# Remove a tag
gws untag my-project personal

# Tags can be anything: personal, work, client, archived, production, etc.
```

### Filtering Repositories

Filter repositories by type, tag, or name:

```bash
# List only GitHub repositories
gws list --type github

# List repositories tagged as "personal"
gws list --tag personal

# List repositories matching a name pattern
gws list --name project

# Combine multiple filters
gws list --type gitlab --tag work
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

2. **Initialize gws** in the test workspace:
   ```bash
   gws init /tmp/test-workspace
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

   # Run gws command
   gws
   ```

   Expected output:
   ```
   Error: workspace not initialized

   To get started, initialize a workspace:
     gws init <directory>

   Example:
     gws init ~/projects
   ```

2. **Test with invalid directory**:
   ```bash
   gws init /nonexistent/directory
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
   gws init /tmp/nested-test
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
   gws init /tmp/skip-test
   ```

   Expected: Only `valid-repo` should be found (node_modules should be skipped).

## Configuration

gws stores its configuration in `~/.gws/config.json`. The configuration includes:

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
│   └── gws/              # Main application entry point
│       ├── main.go       # Root command and CLI setup
│       ├── init.go       # Init command implementation
│       ├── list.go       # List command implementation
│       ├── tag.go        # Tag command implementation
│       └── untag.go      # Untag command implementation
├── internal/
│   ├── classifier/       # Repository classification
│   │   ├── detector.go
│   │   └── detector_test.go
│   ├── config/           # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   └── discovery/        # Repository discovery
│       ├── scanner.go
│       └── scanner_test.go
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
- [x] Search and filter capabilities
- [ ] Git status integration
- [ ] Navigation to repository directories
- [ ] CI/CD pipeline with automated releases
- [ ] Plugin support for editors (Neovim, VSCode)

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
