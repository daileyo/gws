# 01-tasks-gws-core.md

## Tasks

### [x] 1.0 Project Setup and Basic CLI Framework

#### 1.0 Proof Artifact(s)

- Screenshot: `gws --version` output shows version information demonstrates CLI is functional
- Screenshot: `gws --help` output shows available commands demonstrates command structure
- CLI: `make build` successfully compiles binary demonstrates build system works
- CLI: `make test` runs unit tests demonstrates testing framework is configured
- File: `go.mod` with project dependencies demonstrates Go module initialization
- File: `.gitignore` configured for Go artifacts demonstrates proper version control setup

#### 1.0 Tasks

- [x] 1.1 Initialize Go module and create standard Go project structure (cmd/gws/, internal/, pkg/)
- [x] 1.2 Install and configure Cobra CLI framework with root command
- [x] 1.3 Implement version command with version information
- [x] 1.4 Create Makefile with build, test, and clean targets
- [x] 1.5 Create .gitignore configured for Go artifacts (binaries, vendor/, coverage files)
- [x] 1.6 Write initial unit tests for version command
- [x] 1.7 Verify `make build` compiles successfully and `make test` runs tests

### [x] 2.0 Repository Discovery and Workspace Initialization

#### 2.0 Proof Artifact(s)

- Terminal recording: Running `gws init <directory>` discovers repositories demonstrates initialization works
- Screenshot: `~/.gws/config.json` file contents showing discovered repositories demonstrates metadata persistence
- CLI: `gws init /path/to/workspace` returns success message with count of discovered repos demonstrates discovery accuracy
- Screenshot: Error message when running `gws` without initialization demonstrates error handling
- Test: `internal/discovery/scanner_test.go` passes demonstrates repository scanning logic works
- README section: Manual verification steps for repository discovery

#### 2.0 Tasks

- [x] 2.1 Create config package to handle workspace configuration (load, save, create)
- [x] 2.2 Create discovery package with repository scanner (recursive directory traversal, .git detection)
- [x] 2.3 Implement `gws init` command with directory scanning and config persistence
- [x] 2.4 Update root command to validate workspace initialization and show helpful errors
- [x] 2.5 Write unit tests for config package operations
- [x] 2.6 Write unit tests for repository discovery scanner
- [x] 2.7 Create README with manual verification steps and usage examples
- [x] 2.8 Verify end-to-end: init workspace, check config file, test error handling

### [x] 3.0 Repository Classification and Metadata Management

#### 3.0 Proof Artifact(s)

- Screenshot: `gws list` output showing auto-detected repository types (GitHub, GitLab, ADO, Bitbucket) demonstrates automatic classification
- Terminal recording: `gws tag <repo> personal` and `gws untag <repo> personal` demonstrates tagging functionality
- Screenshot: Configuration file showing repository metadata including type, visibility, and tags demonstrates metadata storage
- CLI: `gws list --type github` filters only GitHub repos demonstrates type filtering works
- Test: `internal/classifier/detector_test.go` passes demonstrates URL pattern matching works
- README section: Classification examples and supported remote patterns

#### 3.0 Tasks

- [x] 3.1 Update repository model to include type, visibility, and tags fields
- [x] 3.2 Create classifier package with URL pattern matching for repository type detection (GitHub, GitLab, ADO, Bitbucket)
- [x] 3.3 Implement visibility detection logic based on remote URL protocol (https vs ssh)
- [x] 3.4 Update config package to persist and load classification metadata (type, visibility, tags)
- [x] 3.5 Implement `gws tag <repo> <tag>` command to add custom tags to repositories
- [x] 3.6 Implement `gws untag <repo> <tag>` command to remove tags from repositories
- [x] 3.7 Update `gws list` command to display repository type, visibility, and tags
- [x] 3.8 Write unit tests for classifier package (URL pattern matching and visibility detection)
- [x] 3.9 Write unit tests for tag management operations
- [x] 3.10 Update README with classification examples and supported remote URL patterns
- [x] 3.11 Verify end-to-end: auto-detection on init/refresh, tag/untag operations, list output

### [x] 4.0 Search and Filter Capabilities

#### 4.0 Proof Artifact(s)

- Terminal recording: Multiple filter combinations (by name, type, tags, status) demonstrates filtering works
- Screenshot: `gws list --type github --tag personal` output demonstrates multi-criteria filtering
- Screenshot: Table format output with aligned columns demonstrates readable output formatting
- Screenshot: `gws list --json` output demonstrates JSON output for scripting
- CLI: `gws list --name "my-project"` returns matching repositories demonstrates name search
- Test: `internal/filter/filter_test.go` passes demonstrates filter logic works
- README section: All supported filter flags and usage examples

#### 4.0 Tasks

- [x] 4.1 Add JSON output format support (--json flag) to list command
- [x] 4.2 Add path filtering support (--path flag for absolute or partial path matching)
- [x] 4.3 Add support for multiple tags filtering (allow multiple --tag flags or comma-separated)
- [x] 4.4 Create filter package to centralize and test filtering logic
- [x] 4.5 Write comprehensive unit tests for filter package
- [x] 4.6 Add output formatting options (table, json, simple)
- [x] 4.7 Update README with complete filter examples and output format options
- [x] 4.8 Verify end-to-end: multiple filter combinations, JSON output, path filtering

### [x] 5.0 Git Status Integration and Navigation

#### 5.0 Proof Artifact(s)

- Terminal recording: `gws` command navigates to workspace root demonstrates navigation behavior
- Screenshot: `gws list` output showing git status indicators (branch, clean/dirty, ahead/behind) demonstrates status integration
- Screenshot: Before/after `gws refresh` showing updated git status demonstrates cache refresh
- CLI: `gws personal` lists repos tagged "personal" demonstrates filter shorthand
- Screenshot: Status output with visual indicators (symbols/colors) demonstrates enhanced UX
- Test: `internal/git/status_test.go` passes demonstrates git status reading works
- README section: Navigation patterns and status indicator documentation

#### 5.0 Tasks

- [x] 5.1 Create git package for reading git status information (branch, clean/dirty, ahead/behind)
- [x] 5.2 Implement status caching with configurable TTL (default 5 minutes)
- [x] 5.3 Update list command to display git status indicators
- [x] 5.4 Implement `gws refresh` command to update cached status and rediscover repos
- [x] 5.5 Add navigation to workspace root when running `gws` with no arguments
- [x] 5.6 Add filter shorthand support (e.g., `gws personal` lists repos tagged "personal")
- [x] 5.7 Write unit tests for git status package
- [x] 5.8 Add visual indicators (symbols/colors) for status display
- [x] 5.9 Update README with navigation patterns and status indicators documentation
- [x] 5.10 Verify end-to-end: navigation, status display, refresh, filter shorthand

### [ ] 6.0 CI/CD Pipeline and Release Automation

#### 6.0 Proof Artifact(s)

- File: `.github/workflows/ci.yml` workflow configuration demonstrates pipeline setup
- Screenshot: GitHub Actions successful build for all platforms (Linux, macOS, Windows) demonstrates cross-platform builds
- Screenshot: Snyk security scan results in PR demonstrates security scanning integration
- Screenshot: GitHub Release with semantically versioned binaries demonstrates automated releases
- Badge: Snyk security badge in README demonstrates security transparency
- File: `.releaserc.json` or similar semantic-release configuration demonstrates versioning automation
- README section: Release process, versioning strategy, and contribution guidelines

#### 6.0 Tasks

- [ ] 6.1 Create `.github/workflows/ci.yml` for continuous integration (build, test, lint on every PR and push)
  - Configure Go setup with version matrix (e.g., 1.21.x, 1.22.x)
  - Run `go build`, `go test -race -coverprofile`, `go vet`
  - Add golangci-lint for static analysis
  - Cache Go modules for faster builds
  - Run on `push` to main and all `pull_request` events
- [ ] 6.2 Integrate Snyk security scanning into CI workflow
  - Add Snyk GitHub Action to scan Go dependencies
  - Configure severity threshold: block on high, warn on medium/low
  - Display scan results in PR checks
  - Store Snyk token as repository secret
- [ ] 6.3 Create `.goreleaser.yml` configuration for cross-platform builds
  - Configure build targets: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
  - Embed version, commit, and date via ldflags (matching Makefile pattern)
  - Configure archive formats (tar.gz for Unix, zip for Windows)
  - Add checksums generation (sha256)
  - Include README.md and LICENSE in archives
- [ ] 6.4 Create `.github/workflows/release.yml` for automated releases
  - Trigger on version tags (v*)
  - Use GoReleaser GitHub Action to build and publish binaries
  - Publish artifacts to GitHub Releases
  - Generate release notes from conventional commits
  - Run tests before release to ensure quality
- [ ] 6.5 Set up semantic versioning automation with release-please or semantic-release
  - Create `.release-please-manifest.json` or `.releaserc.json` configuration
  - Configure conventional commit parsing (feat → minor, fix → patch, breaking → major)
  - Automate changelog generation
  - Create release PR workflow for version bumps
- [ ] 6.6 Add Snyk security badge to README.md
  - Include badge showing current security scan status
  - Link badge to Snyk dashboard or scan results
- [ ] 6.7 Update Makefile with additional CI-friendly targets
  - Add `lint` target for golangci-lint
  - Add `coverage` target for test coverage report
  - Add `snapshot` target for local GoReleaser testing
- [ ] 6.8 Write unit tests for version command to ensure ldflags injection works
  - Verify version, commit, and date are correctly embedded
- [ ] 6.9 Update README with release process and contribution guidelines
  - Document conventional commit format requirements
  - Document versioning strategy (SemVer)
  - Document how releases are automated
  - Add contributing guidelines section
- [ ] 6.10 Verify end-to-end: trigger CI on PR, verify security scan, create test release
  - Push a feature branch and verify CI runs successfully
  - Verify Snyk scan results appear in checks
  - Create and push a version tag to verify release workflow
  - Download and verify binaries from GitHub Releases for at least one platform
