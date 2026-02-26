# 04 Questions Round 1 - Brew Install

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Homebrew Tap Repository

Homebrew requires a dedicated "tap" repository to host custom formulas (e.g., `github.com/daileyo/homebrew-tap`). Does this repository already exist, or does it need to be created as part of this spec?

- [ ] (A) The tap repo (`daileyo/homebrew-tap` or similar) already exists on GitHub
- [X] (B) The tap repo does not exist and needs to be created as part of this work
- [ ] (C) Other (describe)

## 2. Tap Name / Install Command

What should the user-facing `brew tap` command look like? The tap name is derived from the GitHub repo name (e.g., repo `daileyo/homebrew-tap` → `brew tap daileyo/tap`).

- [ ] (A) `brew tap daileyo/tap` (repo name: `homebrew-tap`) — most concise
- [X] (B) `brew tap daileyo/gws` (repo name: `homebrew-gws`) — matches project name
- [ ] (C) Other (describe)

## 3. Formula Auto-Publishing via GoReleaser

GoReleaser has built-in support for automatically generating and publishing the Homebrew formula to the tap repo on each release. Should we configure GoReleaser to publish the formula automatically?

- [X] (A) Yes — configure GoReleaser to auto-publish the formula on release (recommended; requires a GitHub token with write access to the tap repo added as a CI secret)
- [ ] (B) No — manually maintain the formula in the tap repo; GoReleaser only builds the binaries
- [ ] (C) Other (describe)

## 4. Formula Name / Install Command

The GoReleaser config shows the binary is named `git-workspace` and the GitHub repo is `daileyo/gws`. What should the Homebrew formula be named? This determines how users install it (e.g., `brew install daileyo/tap/git-workspace`).

- [ ] (A) `git-workspace` — matches the binary name, install with `brew install daileyo/tap/git-workspace`
- [X] (B) `gws` — matches the legacy binary / repo name, install with `brew install daileyo/tap/gws`
- [ ] (C) Other (describe)

## 5. GitHub Actions Secret for Tap Publishing

If GoReleaser auto-publishes to the tap repo (option A in question 3), it needs a GitHub Personal Access Token (or fine-grained token) with write access to the tap repo, stored as a secret in the `daileyo/gws` repository. Do you have a token available, or does the spec need to include creating/configuring one?

- [ ] (A) I have a token (or can create one) — the spec just needs to document where to add it as a CI secret
- [X] (B) The spec should include instructions for creating and configuring the token
- [ ] (C) Other (describe)

## 6. Documentation Scope

Which documentation locations should be updated to include the `brew` install instructions?

- [X] (A) README.md only — add a new "Install via Homebrew" section under the existing "Installation" heading
- [ ] (B) README.md + the release footer in `.goreleaser.yml` — so release notes also show the brew install command
- [ ] (C) README.md + `.goreleaser.yml` release footer + any other docs (describe which)
- [ ] (D) Other (describe)
