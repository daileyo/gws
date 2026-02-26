# 04-spec-brew-install

## Introduction/Overview

This feature makes `git-workspace` installable via Homebrew by creating a custom Homebrew tap (`daileyo/homebrew-gws`) and configuring GoReleaser to automatically publish the formula on each release. The goal is to give macOS users a simple, familiar one-step install experience without requiring them to build from source.

## Goals

- Create the `daileyo/homebrew-gws` GitHub repository as a valid Homebrew tap
- Configure GoReleaser to automatically generate and push the `gws` formula to the tap repo on every release
- Configure the CI release workflow with the GitHub token needed to push to the tap
- Document the Homebrew install process in README.md so new users can discover it
- Give users a working `brew tap daileyo/gws && brew install gws` install path

## User Stories

- **As a macOS developer**, I want to install `git-workspace` with a single `brew install` command so that I can get started without cloning the source repo or running build scripts.
- **As a project maintainer**, I want the Homebrew formula to be automatically updated on every release so that I don't have to manually maintain it.
- **As a new project visitor**, I want to see clear Homebrew install instructions in the README so that I know how to install the tool on my machine.

## Demoable Units of Work

### Unit 1: Tap Repository Setup

**Purpose:** Creates the GitHub repository that hosts the Homebrew formula, making it a valid tap that Homebrew can register.

**Functional Requirements:**
- The system shall have a public GitHub repository named `homebrew-gws` under the `daileyo` account.
- The repository shall contain a `Formula/` directory at its root so Homebrew can locate formula files.
- The repository shall have a description and README that identify it as the Homebrew tap for `git-workspace`.

**Proof Artifacts:**
- Screenshot: `https://github.com/daileyo/homebrew-gws` repository page showing the repo exists, is public, and contains a `Formula/` directory.

---

### Unit 2: GoReleaser Formula Auto-Publishing

**Purpose:** Configures the release pipeline so that every tagged release automatically generates the `gws` Homebrew formula and pushes it to the tap repo — no manual step required.

**Functional Requirements:**
- The system shall have a `brews` stanza added to `.goreleaser.yml` that specifies the formula name as `gws`, the target tap repo as `daileyo/homebrew-gws`, the formula output folder as `Formula/`, and reads the push token from the environment variable `HOMEBREW_TAP_GITHUB_TOKEN`.
- The system shall have a GitHub Personal Access Token (PAT) or fine-grained token with `contents: write` permission on `daileyo/homebrew-gws` created and stored as a repository secret named `HOMEBREW_TAP_GITHUB_TOKEN` in `daileyo/gws`.
- The release workflow (`.github/workflows/release.yml`) shall pass the `HOMEBREW_TAP_GITHUB_TOKEN` secret as an environment variable to the GoReleaser step.
- The user shall be able to verify that GoReleaser processes the brews stanza without errors by running `make snapshot` locally (snapshot does not push; it validates config and generates the formula file locally).

**Proof Artifacts:**
- CLI: `make snapshot` output shows GoReleaser completes successfully and a `Formula/gws.rb` file is generated in the local `dist/` directory, demonstrating the formula generation logic is correct.
- Screenshot: `Formula/gws.rb` file visible in the `daileyo/homebrew-gws` repository on GitHub after a real or test release, demonstrating the auto-publish pipeline works end-to-end.

---

### Unit 3: README Documentation Update

**Purpose:** Makes the Homebrew install path discoverable so new users know they can install with `brew` instead of building from source.

**Functional Requirements:**
- The system shall have a new "Install via Homebrew" subsection added under the existing `## Installation` heading in `README.md`.
- The subsection shall include the exact commands a user needs to run: `brew tap daileyo/gws` followed by `brew install gws`.
- The Homebrew install section shall appear before the existing "Build from Source" subsection so it is the first install option shown.

**Proof Artifacts:**
- Screenshot: `README.md` rendered on the `daileyo/gws` GitHub repository page showing the "Install via Homebrew" section with the correct `brew tap` and `brew install` commands above the "Build from Source" section.

## Non-Goals (Out of Scope)

1. **Submission to `homebrew-core`**: This spec only covers a custom tap. Getting into the official homebrew-core registry requires separate criteria (download thresholds, usage metrics, review process) and is not included.
2. **Linux package managers**: No `.deb`, `.rpm`, Snap, Flatpak, or other Linux packaging is in scope.
3. **Windows package managers**: No Chocolatey, Scoop, or Winget packaging is in scope.
4. **`gws` binary rename**: The formula is named `gws` for user convenience, but the installed binary remains `git-workspace`. Renaming or aliasing the binary inside the formula is out of scope.
5. **Cask creation**: This spec covers a Homebrew formula (CLI tool), not a Homebrew Cask (GUI app).

## Design Considerations

No specific UI/UX design requirements. The only user-facing design element is the README section layout: the Homebrew install block should use a fenced code block with `bash` syntax highlighting, consistent with the existing README style.

## Repository Standards

- **Commit convention**: Conventional Commits. The tap repo creation commit should use `feat:` type. GoReleaser config changes use `feat:` or `chore:`. README updates use `docs:`.
- **Branch and PR workflow**: All changes to `daileyo/gws` should be made on a feature branch and merged via PR to `main` following the existing workflow.
- **GoReleaser config**: Follow the existing `.goreleaser.yml` structure. Add the `brews` stanza after the `archives` section.
- **CI secrets**: Secrets are referenced in workflow files using `${{ secrets.SECRET_NAME }}` syntax, consistent with how `SNYK_TOKEN` is already used in the existing workflows.
- **Makefile**: No changes to the Makefile are required unless a convenience target is needed for local formula validation.

## Technical Considerations

- **Tap repo naming**: Homebrew derives the tap name from the repo name by stripping the `homebrew-` prefix. The repo `daileyo/homebrew-gws` becomes the tap `daileyo/gws`, so the user runs `brew tap daileyo/gws`.
- **GoReleaser `brews` stanza**: GoReleaser generates the formula from the release archives it already builds (macOS amd64 and arm64 `.tar.gz` files). It computes the SHA256 checksums automatically. The formula name (`gws`) determines the file written to `Formula/gws.rb` in the tap repo.
- **Token scope**: The GitHub token must have `contents: write` access to `daileyo/homebrew-gws` only. A fine-grained PAT scoped to that single repository is preferred over a broad classic PAT.
- **Token creation steps** (to be followed during implementation):
  1. In GitHub, go to **Settings → Developer settings → Personal access tokens → Fine-grained tokens**.
  2. Click **Generate new token**. Set an expiration. Under **Repository access**, select "Only select repositories" and choose `daileyo/homebrew-gws`. Under **Permissions**, grant **Contents: Read and write**.
  3. Copy the generated token value.
  4. In the `daileyo/gws` repository, go to **Settings → Secrets and variables → Actions → New repository secret**. Name it `HOMEBREW_TAP_GITHUB_TOKEN` and paste the token value.
- **`make snapshot` validation**: GoReleaser's `--snapshot` flag builds archives and generates formula files locally without pushing. Running `make snapshot` is a safe, non-destructive way to verify the brews config is correct before a real release.
- **Formula homepage and description**: The formula should set `homepage` to `https://github.com/daileyo/gws` and `description` to `"A lightweight CLI tool for discovering, organizing, and navigating git repositories"`.

## Security Considerations

- **Token storage**: The `HOMEBREW_TAP_GITHUB_TOKEN` value must never be committed to any file in either repository. It exists only as a GitHub Actions secret.
- **Token scope minimization**: Use a fine-grained PAT limited to `contents: write` on `daileyo/homebrew-gws` only, not a broad classic PAT with full repo access.
- **Tap repo visibility**: The tap repo must be public for Homebrew to access it. No sensitive information should be committed to the tap repo (it will only contain `Formula/gws.rb` and a README).
- **Proof artifact security**: Do not include the token value in any proof artifact screenshots or documents.

## Success Metrics

1. **Installable via brew**: `brew tap daileyo/gws && brew install gws` completes without errors on a macOS machine and produces a working `git-workspace` binary.
2. **Auto-publishing verified**: The `Formula/gws.rb` file in `daileyo/homebrew-gws` is updated automatically within the CI release pipeline after a tagged release, with no manual steps.
3. **Documentation visible**: The "Install via Homebrew" section is present in `README.md` on the main branch and rendered correctly on the GitHub repository page.

## Open Questions

No open questions at this time.
