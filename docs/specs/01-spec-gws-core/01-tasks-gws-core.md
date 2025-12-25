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

### [ ] 2.0 Repository Discovery and Workspace Initialization

#### 2.0 Proof Artifact(s)

- Terminal recording: Running `gws init <directory>` discovers repositories demonstrates initialization works
- Screenshot: `~/.gws/config.json` file contents showing discovered repositories demonstrates metadata persistence
- CLI: `gws init /path/to/workspace` returns success message with count of discovered repos demonstrates discovery accuracy
- Screenshot: Error message when running `gws` without initialization demonstrates error handling
- Test: `internal/discovery/scanner_test.go` passes demonstrates repository scanning logic works
- README section: Manual verification steps for repository discovery

#### 2.0 Tasks

TBD

### [ ] 3.0 Repository Classification and Metadata Management

#### 3.0 Proof Artifact(s)

- Screenshot: `gws list` output showing auto-detected repository types (GitHub, GitLab, ADO, Bitbucket) demonstrates automatic classification
- Terminal recording: `gws tag <repo> personal` and `gws untag <repo> personal` demonstrates tagging functionality
- Screenshot: Configuration file showing repository metadata including type, visibility, and tags demonstrates metadata storage
- CLI: `gws list --type github` filters only GitHub repos demonstrates type filtering works
- Test: `internal/classifier/detector_test.go` passes demonstrates URL pattern matching works
- README section: Classification examples and supported remote patterns

#### 3.0 Tasks

TBD

### [ ] 4.0 Search and Filter Capabilities

#### 4.0 Proof Artifact(s)

- Terminal recording: Multiple filter combinations (by name, type, tags, status) demonstrates filtering works
- Screenshot: `gws list --type github --tag personal` output demonstrates multi-criteria filtering
- Screenshot: Table format output with aligned columns demonstrates readable output formatting
- Screenshot: `gws list --json` output demonstrates JSON output for scripting
- CLI: `gws list --name "my-project"` returns matching repositories demonstrates name search
- Test: `internal/filter/filter_test.go` passes demonstrates filter logic works
- README section: All supported filter flags and usage examples

#### 4.0 Tasks

TBD

### [ ] 5.0 Git Status Integration and Navigation

#### 5.0 Proof Artifact(s)

- Terminal recording: `gws` command navigates to workspace root demonstrates navigation behavior
- Screenshot: `gws list` output showing git status indicators (branch, clean/dirty, ahead/behind) demonstrates status integration
- Screenshot: Before/after `gws refresh` showing updated git status demonstrates cache refresh
- CLI: `gws personal` lists repos tagged "personal" demonstrates filter shorthand
- Screenshot: Status output with visual indicators (symbols/colors) demonstrates enhanced UX
- Test: `internal/git/status_test.go` passes demonstrates git status reading works
- README section: Navigation patterns and status indicator documentation

#### 5.0 Tasks

TBD

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

TBD
