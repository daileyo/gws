<p align="center">
  <img src="docs/site/assets/images/gws-logo-hero.png" alt="gws logo" width="300">
</p>
<h1 align="center">git-workspace</h1>
<p align="center"><em>Your Git workspace, simplified</em></p>

[![CI](https://github.com/daileyo/gws/actions/workflows/ci.yml/badge.svg)](https://github.com/daileyo/gws/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/daileyo/gws)](https://goreportcard.com/report/github.com/daileyo/gws)
[![Snyk Security](https://snyk.io/test/github/daileyo/gws/badge.svg)](https://snyk.io/test/github/daileyo/gws)
[![Release](https://img.shields.io/github/v/release/daileyo/gws)](https://github.com/daileyo/gws/releases/latest)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Subcommand-Based CLI**: Clean command structure with `list`, `init`, `add`, `refresh`, `tag`, `user`, and `parent` subcommands
- **Compact Default View**: Multi-column repository names at a glance, with verbose modes for detailed tables
- **Repository Discovery**: Automatically find all git repositories in a directory tree
- **Automatic Classification**: Detect repository type (GitHub, GitLab, Azure DevOps, Bitbucket) from remote URLs
- **Dual-Purpose Filtering**: Lowercase flags filter results; uppercase flags show columns with optional filtering
- **Git Status Integration**: View branch, clean/dirty state, and ahead/behind indicators with color output
- **Smart Caching**: Fast status display with configurable cache and concurrent workers (`--workers`)
- **Custom Tagging**: Organize repositories with `gws tag add` / `gws tag remove`
- **Advanced Filtering**: Search and filter repositories by type, tags, name, path, status, user, or remote URL
- **User Profile Management**: Manage git user profiles across repositories with `gws user`
- **Repository Navigation**: Jump to any repository instantly with `gws <repo-name>`
- **Parent Navigation**: Navigate to a repository's parent directory with `gws parent` or `gws -p`
- **Remote URL Display**: View formatted or raw remote URLs with `--show-remote` / `--show-remote-raw`
- **External Repo Symlinks**: Automatically creates workspace symlinks for repositories outside the workspace root
- **Workspace Management**: Track and organize repositories in a centralized configuration
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Lightweight**: Single binary with no external dependencies

## Documentation

Full documentation is available at **https://daileyo.github.io/gws**

## Installation

### Install via Homebrew

```bash
brew install daileyo/gws/git-workspace
```

See the [Getting Started](https://daileyo.github.io/gws/getting-started/) guide for additional installation options.

## Development

### One-Time Setup

Git hooks are provided in `.githooks/` for pre-push checks and commit message formatting. Activate them once after cloning:

```bash
make setup-hooks
```

### Build Commands

```bash
# Build the binary
make build

# Run tests
make test

# Run all CI checks (vet, lint, test with race detector)
make ci

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
```

### Building the Documentation Site

```bash
# Install Python dependencies
pip install -r docs/requirements.txt

# Preview locally
python -m mkdocs serve

# Build static site
python -m mkdocs build
```

### Project Structure

```
.
├── cmd/
│   └── git-workspace/    # Main application entry point
├── internal/
│   ├── classifier/       # Repository classification
│   ├── config/           # Configuration management
│   ├── discovery/        # Repository discovery
│   ├── filter/           # Repository filtering logic
│   └── git/              # Git status integration
├── docs/
│   ├── site/             # MkDocs documentation source
│   ├── specs/            # Specification documents
│   └── requirements.txt  # Python dependencies (MkDocs)
├── Makefile              # Build automation
├── mkdocs.yml            # Documentation site configuration
├── go.mod                # Go module definition
└── README.md             # This file
```

## Contributing

This project is not currently open for external contributions. This may change in the future as the project matures.

If you're interested in contributing or have ideas to share, please [open an issue](https://github.com/daileyo/gws/issues) or reach out directly — feedback is always welcome.

### Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning:

```
<type>[optional scope]: <description>
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

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
