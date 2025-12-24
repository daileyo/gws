# 01-spec-gws-core.md

## Introduction/Overview

Git Workspace (gws) is a lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system. Rather than performing git operations directly, gws serves as an intelligent repository index and navigation layer, allowing developers to quickly find, filter, and access repositories based on metadata such as type (GitHub, GitLab, ADO, Bitbucket), custom tags, status, and location.

## Goals

- Enable rapid discovery and navigation of git repositories across multiple directories
- Provide flexible classification through both automatic detection (by remote URL) and manual tagging
- Support powerful search and filtering based on repository metadata (name, type, tags, status, dates)
- Deliver a cross-platform, dependency-free tool that works identically on Linux, macOS, and Windows
- Establish a robust CI/CD pipeline with automated semantic versioning, security scanning, and artifact publishing
- Create an extensible foundation that can later serve as a backend for editor plugins (neovim, vscode)

## User Stories

- **As a developer with many repositories**, I want to automatically discover all git repos in my workspace directories so that I don't have to manually track each one.

- **As a developer working across personal and professional projects**, I want to classify repositories by type and custom tags so that I can quickly filter to the subset I need.

- **As a developer switching contexts frequently**, I want to quickly navigate to the location of all my repositories or filter to specific categories (like "personal" or "work") so that I can find what I need without searching manually.

- **As a developer managing multiple repository sources**, I want gws to automatically identify whether a repo is from GitHub, GitLab, Azure DevOps, or Bitbucket so that I can filter and organize accordingly.

- **As a developer who values simplicity**, I want a single binary with no dependencies that works the same way on all my machines so that I can install and use it anywhere without hassle.

- **As a contributor or early adopter**, I want each development milestone to produce a usable, tested artifact so that I can verify functionality and provide feedback at every stage.

- **As a security-conscious developer**, I want all builds to be scanned for vulnerabilities and blocked if high-severity issues are found so that I can trust the security of the tool.

## Demoable Units of Work

### Unit 1: Repository Discovery and Workspace Initialization

**Purpose:** Establish the foundational workspace management and repository discovery capabilities that enable gws to find and track git repositories.

**Functional Requirements:**
- The system shall provide a command to initialize a gws workspace location where repositories are organized
- The system shall scan specified directories recursively to discover git repositories (identified by `.git` directory)
- The system shall store discovered repository metadata including: name, absolute path, and remote URL
- The system shall persist repository metadata in a JSON or YAML configuration file in the user's home directory
- The system shall support the `gws` command to navigate to the initialized workspace location
- The system shall display appropriate error messages when no workspace is initialized or no repositories are found

**Proof Artifacts:**
- Terminal recording: Running `gws init <directory>` demonstrates workspace initialization
- Screenshot: `gws` command output showing navigation to workspace location demonstrates workspace access
- Screenshot: Contents of config file showing discovered repositories demonstrates metadata persistence
- GitHub Actions artifact: Binary for Linux, macOS, Windows available for download demonstrates cross-platform build
- Snyk scan results: Clean scan or warnings only (no high-severity blocks) demonstrates security validation
- README section: Manual verification steps for discovery process

### Unit 2: Repository Classification and Metadata Management

**Purpose:** Enable automatic and manual classification of repositories to support filtering and organization.

**Functional Requirements:**
- The system shall automatically detect repository type by parsing remote URL patterns:
  - `github.com` → GitHub
  - `gitlab.com` or self-hosted GitLab → GitLab
  - `dev.azure.com` or `visualstudio.com` → Azure DevOps (ADO)
  - `bitbucket.org` → Bitbucket
- The system shall determine repository visibility (public/private) based on remote URL protocol (https vs ssh)
- The system shall allow users to add custom tags to repositories via CLI command
- The system shall allow users to remove custom tags from repositories
- The system shall store classification metadata (type, visibility, tags) alongside repository information
- The system shall handle repositories with multiple remotes by using the primary/origin remote

**Proof Artifacts:**
- Screenshot: Repository list showing auto-detected types demonstrates automatic classification
- Screenshot: Adding and removing tags via CLI demonstrates tagging functionality
- Terminal recording: Commands for classification management demonstrates complete workflow
- GitHub Release: Semantically versioned artifact (e.g., v0.2.0) demonstrates versioning
- Snyk scan results: Passing security scan demonstrates continued security compliance
- README section: Manual verification of classification accuracy

### Unit 3: Search and Filter Capabilities

**Purpose:** Provide powerful search and filtering to help users quickly find specific repositories from a large collection.

**Functional Requirements:**
- The system shall provide a `gws list` command that displays all tracked repositories with their metadata
- The system shall support filtering repositories by:
  - Repository name (partial match, case-insensitive)
  - Repository path (absolute or relative)
  - Repository type (github, gitlab, ado, bitbucket)
  - Custom tags (one or more tags)
  - Repository status (clean, dirty, ahead, behind)
  - Date-based criteria (last commit date, last modified)
  - Current branch name
- The system shall allow combining multiple filter criteria (AND logic)
- The system shall display results in a human-readable table format showing: name, path, type, tags, and status
- The system shall handle zero results gracefully with an informative message

**Proof Artifacts:**
- Terminal recording: Demonstrating various filter combinations shows filtering capabilities
- Screenshot: `gws list --type github --tag personal` output demonstrates multi-criteria filtering
- Screenshot: Search results table format demonstrates readable output
- GitHub Release: Semantically versioned artifact (e.g., v0.3.0) demonstrates continued development
- Snyk scan results: Passing security scan demonstrates security compliance
- README section: Examples of all supported filter flags and combinations

### Unit 4: Navigation and Git Status Integration

**Purpose:** Enable quick navigation to filtered repository sets and provide basic git status information for context.

**Functional Requirements:**
- The system shall support `gws` with no arguments to navigate to the workspace root directory (default behavior)
- The system shall support `gws <filter>` to list repositories matching the filter (e.g., `gws personal` lists repos tagged "personal")
- The system shall navigate to workspace root by default unless a specific repository is selected/specified
- The system shall integrate basic git status for each repository showing:
  - Current branch name
  - Uncommitted changes indicator (clean/dirty)
  - Ahead/behind remote status
- The system shall cache git status information with configurable TTL to avoid slow operations
- The system shall provide a `gws refresh` command to update cached git status and repository metadata
- The system shall display git status in list output using visual indicators (symbols or color if terminal supports it)

**Proof Artifacts:**
- Terminal recording: Using `gws` and `gws personal` demonstrates navigation and filtering
- Screenshot: List output showing git status indicators demonstrates status integration
- Screenshot: Before/after `gws refresh` showing updated status demonstrates cache refresh
- GitHub Release: Semantically versioned artifact (e.g., v0.4.0 or v1.0.0) demonstrates complete feature set
- Snyk scan results: Passing security scan demonstrates continued security compliance
- README section: Documentation of all navigation patterns and status indicators

### Unit 5: CI/CD Pipeline and Release Automation

**Purpose:** Establish automated build, test, security scanning, and release pipeline to ensure quality and enable continuous delivery.

**Functional Requirements:**
- The system shall use GitHub Actions for continuous integration and deployment
- The system shall implement semantic versioning (SemVer) for all releases
- The system shall automatically determine version bumps based on conventional commits
- The system shall run Snyk security scanning on every build (PR and release)
- The system shall block builds if high-severity vulnerabilities are found
- The system shall warn (but not block) on medium and low-severity vulnerabilities
- The system shall build binaries for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64)
- The system shall publish release artifacts to GitHub Releases for each version
- The system shall run automated tests before publishing any artifacts
- The system shall generate and display security scan results in CI output

**Proof Artifacts:**
- GitHub Actions workflow file: `.github/workflows/ci.yml` demonstrates pipeline configuration
- GitHub Actions runs: Successful pipeline executions with security scans demonstrate automation
- GitHub Releases page: Multiple semantically versioned releases with downloadable binaries demonstrate versioning and distribution
- Snyk scan results: Badge or report in README showing security status demonstrates transparency
- README section: Documentation of release process, versioning strategy, and security policy

## Non-Goals (Out of Scope)

1. **Git operations**: gws will NOT perform git operations like clone, pull, push, commit, merge, etc. It is a navigation and discovery tool only.
2. **Repository modification**: gws will NOT modify repository contents, configuration, or git history.
3. **Editor integration**: Editor plugins for neovim, vscode, etc. are future enhancements, not part of the initial release.
4. **Package manager distribution**: Initial release focuses on single binary and build-from-source; package managers (brew, apt, cargo) are future enhancements.
5. **Content search**: Searching within repository file contents is not included; only metadata-based search.
6. **Multi-machine synchronization**: Tracking repositories across multiple machines is not supported in the initial version.
7. **Advanced analytics**: Deep analysis like commit frequency, dependency detection, branch health are out of scope; only basic git status.
8. **Docker images**: Docker containerization is not required for a CLI tool; focus is on native binary distribution.

## Design Considerations

**Command Line Interface:**
- Use git-style subcommands: `gws <subcommand> [flags]`
- Support both subcommand mode and interactive mode
- Common commands:
  - `gws init <directory>` - Initialize workspace and discover repos
  - `gws list [--filter-options]` - List repositories with filters
  - `gws tag <repo> <tag>` - Add tag to repository
  - `gws untag <repo> <tag>` - Remove tag from repository
  - `gws refresh` - Refresh repository metadata and git status
  - `gws` - Navigate to workspace root

**Output Format:**
- Use table format for list views with aligned columns
- Include visual indicators for status (e.g., `✓` for clean, `*` for dirty, `↑` for ahead)
- Support optional JSON output for scripting via `--json` flag
- Color output when terminal supports it (with `--no-color` override)

## Repository Standards

**Project Structure:**
- Use Go modules for dependency management (`go.mod`, `go.sum`)
- Follow standard Go project layout:
  - `cmd/gws/` - Main application entry point
  - `internal/` - Private application code
  - `pkg/` - Public libraries (if any)
  - `docs/` - Documentation including this spec
  - `tests/` - Integration tests
  - `.github/workflows/` - GitHub Actions pipeline definitions

**Version Control:**
- Use conventional commits for all commit messages (e.g., `feat:`, `fix:`, `docs:`, `chore:`)
- Follow semantic versioning (SemVer) for releases: MAJOR.MINOR.PATCH
- Tag releases with `v` prefix (e.g., `v1.0.0`)

**CI/CD Pipeline:**
- Use GitHub Actions for all automation
- Implement the following workflows:
  - **CI Workflow**: Run on every PR and push (build, test, lint, security scan)
  - **Release Workflow**: Run on tag push (build artifacts, create release, publish binaries)
- Integrate Snyk for security vulnerability scanning:
  - Run on every build
  - Block pipeline on high-severity vulnerabilities
  - Warn (log but allow) on medium and low-severity vulnerabilities
- Publish binaries to GitHub Releases
- Include security scan badge in README

**Documentation:**
- Include comprehensive README with installation and usage instructions
- Provide manual verification steps for each feature
- Document all CLI commands and flags
- Include contributing guidelines
- Display security scan status badge

**Quality Standards:**
- Include `.gitignore` for Go artifacts
- Provide `Makefile` for common build/test tasks
- Write unit tests for core functionality
- Write integration tests for CLI commands
- Maintain >80% code coverage where practical

## Technical Considerations

**Language and Technology:**
- Implement in **Go** for:
  - Excellent cross-platform support (Linux, macOS, Windows)
  - Single binary distribution with no runtime dependencies
  - Strong standard library for file I/O, CLI parsing, and system interaction
  - Fast compilation and execution
  - Good ecosystem for CLI tools (cobra, viper, etc.)

**Configuration Storage:**
- Store metadata in `~/.gws/config.json` or `~/.gws/config.yaml`
- Use JSON for simplicity and wide tool support
- Structure: `{ "workspace": "path", "repositories": [...] }`
- Support one root directory for workspace (multi-directory support is a future enhancement)
- Include config version for future migration support

**Git Integration:**
- Use `go-git` library for reading git metadata without shelling out to git CLI
- This provides better cross-platform support and avoids git dependency
- Cache git status results with configurable TTL (default 5 minutes)

**Cross-Platform Considerations:**
- Use `filepath.Join()` for all path operations to handle OS-specific separators
- Test path handling on Windows (backslashes, drive letters, UNC paths)
- Ensure home directory detection works on all platforms (`os.UserHomeDir()`)
- Handle symlinks appropriately
- Build for multiple architectures: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64

**Performance:**
- Lazy-load git status (don't fetch on discovery, only when displayed/requested)
- Use goroutines for parallel git status checks when listing many repos
- Index repositories by name/path for fast lookup
- Consider mmap for large config files if needed in future

**CI/CD Technical Details:**
- Use GitHub Actions matrix builds for multi-platform compilation
- Use semantic-release for automated versioning based on conventional commits
- Use goreleaser for building and publishing cross-platform binaries
- Use Snyk GitHub Action for security scanning with custom policy:
  ```yaml
  # Snyk configuration
  severity-threshold: high  # Block on high, warn on medium/low
  fail-on: upgradable       # Only fail if there's an upgrade available
  ```
- Cache Go modules in CI for faster builds
- Store build artifacts for each workflow run
- Ensure every tagged release produces downloadable binaries

## Security Considerations

**Sensitive Data:**
- Remote URLs may contain credentials (HTTPS URLs with embedded tokens); ensure these are not logged or displayed
- Consider redacting sensitive portions of URLs in output
- Repository metadata config file should use user-only permissions (0600)

**Git Operations:**
- Since gws doesn't perform git operations, attack surface is limited
- Reading git metadata is low-risk but validate all file paths to prevent directory traversal
- Sanitize all user inputs to prevent command injection

**Dependency Security:**
- Use Snyk to scan dependencies for known vulnerabilities on every build
- **Block builds on high-severity vulnerabilities** (pipeline fails)
- **Warn on medium and low-severity vulnerabilities** (logged but pipeline continues)
- Keep dependencies up to date with automated PRs (Dependabot)
- Display security scan status in README via badge

**Proof Artifacts:**
- Do NOT commit any `.gws` configuration files to this repository as they may contain local paths and potentially sensitive remote URLs
- Use example/sanitized configs for documentation
- Ensure Snyk scan results don't expose internal implementation details
- Security scan reports should be viewable in GitHub Actions logs

## Success Metrics

1. **Discovery accuracy**: Successfully discovers 100% of git repositories in scanned directories
2. **Cross-platform compatibility**: Identical behavior on Linux, macOS, and Windows with automated builds for all platforms
3. **Performance**: Can index and filter 1000+ repositories in under 1 second (excluding git status fetch)
4. **Usability**: Users can find a specific repository with 1-2 commands without consulting documentation
5. **Reliability**: Zero data loss or corruption of repository metadata
6. **Automation**: Every commit triggers CI with security scans, every tag produces a release with binaries, zero manual steps for publishing
7. **Security**: Zero high-severity vulnerabilities in dependencies and code; medium/low vulnerabilities documented and tracked

## Open Questions

1. Should gws support watching directories for new repositories automatically, or require manual `gws refresh`?
2. For repositories with multiple remotes, should users be able to specify which remote to use for type detection?
3. What should the default cache TTL be for git status information (5 minutes, 10 minutes, 1 hour)?

## Resolved Design Decisions

The following questions have been resolved and are reflected in the specification:
- **Workspace root directories**: Support one root directory initially; multi-directory support is a future enhancement
- **Navigation behavior**: By default, `gws` navigates to workspace root unless a specific repository is selected/specified
- **Versioning tool**: Use semantic-release for automated versioning based on conventional commits

---

**Next Steps**: Once this spec is approved, run `/generate-task-list-from-spec` to create the implementation task list.
