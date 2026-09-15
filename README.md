<p align="center">
  <img src="docs/site/assets/images/omgitworks-logo-hero.png" alt="omgitworks logo" width="400">
</p>
<h1 align="center">omgitworks</h1>
<p align="center"><em>Your Git workspace, simplified</em></p>

[![CI](https://github.com/daileyo/omgitworks/actions/workflows/ci.yml/badge.svg)](https://github.com/daileyo/omgitworks/actions/workflows/ci.yml)
[![Snyk Security](https://snyk.io/test/github/daileyo/omgitworks/badge.svg)](https://snyk.io/test/github/daileyo/omgitworks)
[![Release](https://img.shields.io/github/v/release/daileyo/omgitworks)](https://github.com/daileyo/omgitworks/releases/latest)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Subcommand-Based CLI**: Clean command structure with `list`, `init`, `add`, `refresh`, `cd`, `tag`, `user`, `parent`, and `worktree` subcommands
- **Compact Default View**: Multi-column repository names at a glance, with verbose modes for detailed tables
- **Repository Discovery**: Automatically find all git repositories in a directory tree
- **Automatic Classification**: Detect repository type (GitHub, GitLab, Azure DevOps, Bitbucket) from remote URLs
- **Dual-Purpose Filtering**: Lowercase flags filter results; uppercase flags show columns with optional filtering
- **Git Status Integration**: View branch, clean/dirty state, and ahead/behind indicators with color output
- **Smart Caching**: Fast status display with configurable cache and concurrent workers (`--workers`)
- **Custom Tagging**: Organize repositories with `omgw tag add` / `omgw tag remove`
- **Advanced Filtering**: Search and filter repositories by type, tags, name, path, status, user, or remote URL
- **User Profile Management**: Manage git user profiles across repositories with `omgw user`
- **Repository Navigation**: Jump to any repository instantly with `omgw <repo-name>`
- **Workspace Navigation**: Jump to the workspace root with `omgw cd`
- **Parent Navigation**: Navigate to a repository's parent directory with `omgw parent` or `omgw -p`
- **Worktree Management**: Discover, navigate, create, and organize git worktrees with `omgw worktree`
- **XDG Layout**: Config in `~/.config/gws/`, worktrees in `~/.local/share/gws/projects/` — the same paths on Linux, macOS, and Windows
- **Remote URL Display**: View formatted or raw remote URLs with `--show-remote` / `--show-remote-raw`
- **External Repo Symlinks**: Automatically creates workspace symlinks for repositories outside the workspace root
- **Workspace Management**: Track and organize repositories in a centralized configuration
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Lightweight**: Single binary with no external dependencies

## Documentation

Full documentation is available at **https://omgitworks.dev**

## Installation

### Install via Homebrew

```bash
brew install daileyo/tap/omgitworks
```

> Homebrew 7.0 refuses to load formulae from untrusted third-party taps. The command above is
> unaffected — naming the tap explicitly trusts that formula. If you tap first and install by
> bare name, run `brew trust daileyo/tap` in between.

See the [Getting Started](https://omgitworks.dev/getting-started/) guide for additional installation options.

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
│   └── omgitworks/    # Main application entry point
├── internal/
│   ├── classifier/       # Repository classification
│   ├── config/           # Configuration management
│   ├── discovery/        # Repository discovery
│   ├── filter/           # Repository filtering logic
│   ├── git/              # Git status integration
│   ├── user/             # Git user profile management
│   └── xdg/              # XDG base directory resolution
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

If you're interested in contributing or have ideas to share, please [open an issue](https://github.com/daileyo/omgitworks/issues) or reach out directly — feedback is always welcome.

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
