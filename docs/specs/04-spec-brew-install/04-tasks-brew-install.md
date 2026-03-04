# 04 Tasks - Brew Install

## Relevant Files

- `.goreleaser.yml` - Add the `brews` stanza after the `archives` section to configure Homebrew formula generation and publishing.
- `.github/workflows/release.yml` - Add `HOMEBREW_TAP_GITHUB_TOKEN` environment variable to the GoReleaser step so CI can push to the tap repo on release.
- `README.md` - Add the "Install via Homebrew" subsection under `## Installation`, before the existing "Build from Source" subsection.

> **Note — External tap repository files (not in this repo):**
> The following files are created in the separate `daileyo/homebrew-gws` GitHub repository during Task 1.0. They are not part of the `daileyo/gws` codebase.
> - `daileyo/homebrew-gws/README.md` — Identifies the repo as a Homebrew tap and shows install commands.
> - `daileyo/homebrew-gws/Formula/.gitkeep` — Placeholder file that creates the required `Formula/` directory. GoReleaser will replace this with `gws.rb` on first release.

### Notes

- Task 1.0 is done entirely in the GitHub web UI (creating a new repository). No local code changes are required.
- Task 2.1 and 2.2 are manual steps in the GitHub web UI (token creation and secret storage). They must be completed before the first real release so CI can push the formula.
- All code changes to `daileyo/gws` should be made on the current `feat-brew-install` branch and merged to `main` via PR.
- `make snapshot` validates GoReleaser config locally without pushing anything. It is safe to run at any time.
- The GoReleaser `brews` stanza should be placed after the `archives` section in `.goreleaser.yml`, following the existing file structure.
- Secrets are referenced in workflow YAML using `${{ secrets.SECRET_NAME }}` syntax, consistent with how `SNYK_TOKEN` is already used.
- Commit messages must follow Conventional Commits: `feat:` for the GoReleaser config change, `docs:` for the README update.

## Tasks

### [X] 1.0 Create the `daileyo/homebrew-gws` Tap Repository

#### 1.0 Proof Artifact(s)

- Screenshot: `https://github.com/daileyo/homebrew-gws` repository page showing the repo is public, has the correct description, and contains a visible `Formula/` directory demonstrates the tap repo is valid and Homebrew-ready.

#### 1.0 Tasks

- [X] 1.1 In the GitHub web UI, create a new **public** repository named `homebrew-gws` under the `daileyo` account. Set the description to: `Homebrew tap for git-workspace — a CLI tool for discovering and navigating git repositories.`
- [X] 1.2 Initialize the repository with a `README.md` (use the GitHub UI option during creation, or commit one manually). The README content should identify this as the tap for `git-workspace` and show the install commands:
  ```
  brew tap daileyo/gws
  brew install gws
  ```
- [X] 1.3 Create a `Formula/` directory in the tap repo by adding a placeholder file `Formula/.gitkeep`. This file can be added via the GitHub web UI ("Create new file" → type `Formula/.gitkeep` as the filename). GoReleaser will replace its contents with `gws.rb` on the first release.

---

### [x] 2.0 Configure GoReleaser Homebrew Formula Auto-Publishing

#### 2.0 Proof Artifact(s)

- CLI: `make snapshot` completes without errors and a `Formula/gws.rb` file is present in the `dist/` directory demonstrates the GoReleaser brews stanza is correctly configured and formula generation works.
- Diff: `.goreleaser.yml` includes a `brews` stanza with `name: gws`, target repo `daileyo/homebrew-gws`, and token from `HOMEBREW_TAP_GITHUB_TOKEN` demonstrates the formula publishing config is in place.
- Diff: `.github/workflows/release.yml` passes `HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}` to the GoReleaser step demonstrates the CI workflow will supply the token on release.

#### 2.0 Tasks

- [X] 2.1 In the GitHub web UI, create a fine-grained Personal Access Token (PAT) scoped only to the `daileyo/homebrew-gws` repository:
  1. Go to **GitHub Settings → Developer settings → Personal access tokens → Fine-grained tokens**.
  2. Click **Generate new token**. Give it a name like `gws-homebrew-tap-publisher`. Set an expiration date (e.g., 1 year).
  3. Under **Repository access**, choose "Only select repositories" and select `daileyo/homebrew-gws`.
  4. Under **Permissions → Repository permissions**, set **Contents** to `Read and write`. Leave all other permissions at `No access`.
  5. Click **Generate token** and copy the token value immediately (it will not be shown again).
- [X] 2.2 Add the token as a repository secret in `daileyo/gws`:
  1. Go to the `daileyo/gws` repository → **Settings → Secrets and variables → Actions**.
  2. Click **New repository secret**.
  3. Set the name to `HOMEBREW_TAP_GITHUB_TOKEN` and paste the token value from step 2.1.
  4. Click **Add secret**.
- [x] 2.3 Add the `brews` stanza to `.goreleaser.yml` after the closing of the `archives` section and before the `checksum` section. The stanza should be:
  ```yaml
  brews:
    - name: gws
      repository:
        owner: daileyo
        name: homebrew-gws
        token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
      commit_author:
        name: goreleaserbot
        email: bot@goreleaser.com
      folder: Formula
      homepage: "https://github.com/daileyo/gws"
      description: "A lightweight CLI tool for discovering, organizing, and navigating git repositories"
      license: "MIT"
  ```
- [x] 2.4 Update the `Run GoReleaser` step in `.github/workflows/release.yml` to pass the `HOMEBREW_TAP_GITHUB_TOKEN` secret as an environment variable alongside the existing `GITHUB_TOKEN`:
  ```yaml
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
  ```
- [x] 2.5 Run `make snapshot` from the repository root to validate the GoReleaser configuration locally. Confirm that the command completes without errors and that the file `dist/Formula/gws.rb` exists and contains a valid Homebrew formula with `url`, `sha256`, `homepage`, and `desc` fields.

---

### [x] 3.0 Update README.md with Homebrew Install Instructions

#### 3.0 Proof Artifact(s)

- Screenshot: `README.md` rendered on the `daileyo/gws` GitHub repository page showing the "Install via Homebrew" section with `brew tap daileyo/gws` and `brew install gws` commands appearing above the "Build from Source" section demonstrates the install path is discoverable by new users.

#### 3.0 Tasks

- [x] 3.1 In `README.md`, locate the `## Installation` heading (currently at line 24). Immediately after the `## Installation` line and before the `### Build from Source` line, insert a new subsection with the Homebrew install commands:
  ```markdown
  ### Install via Homebrew

  ```bash
  brew tap daileyo/gws
  brew install gws
  ```
  ```
- [x] 3.2 Verify the README renders correctly by previewing it in a Markdown viewer or pushing to the branch and checking the GitHub preview. Confirm that:
  - The "Install via Homebrew" section appears before "Build from Source"
  - The code block uses `bash` syntax highlighting
  - The two commands (`brew tap` and `brew install`) are on separate lines
