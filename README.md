# gws - Git Workspace

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Repository Discovery**: Automatically find all git repositories in a directory tree
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

### 3. Check Version

```bash
gws version
```

Output:
```
gws version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
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
      "tags": []
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
│       └── init.go       # Init command implementation
├── internal/
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
- [ ] Repository classification (GitHub, GitLab, Azure DevOps, Bitbucket)
- [ ] Manual tagging for organization
- [ ] Search and filter capabilities
- [ ] Git status integration
- [ ] Navigation to repository directories
- [ ] CI/CD pipeline with automated releases
- [ ] Plugin support for editors (Neovim, VSCode)

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
