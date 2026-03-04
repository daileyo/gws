# Task 2.0 Proofs — Configure GoReleaser Homebrew Formula Auto-Publishing

## GoReleaser Config Diff

Added `brews` stanza to `.goreleaser.yml` after `archives` section, before `checksum`:

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
    directory: Formula
    homepage: "https://github.com/daileyo/gws"
    description: "A lightweight CLI tool for discovering, organizing, and navigating git repositories"
    license: "MIT"
```

**Note:** The spec originally specified `folder: Formula` but GoReleaser v2 renamed this field to `directory`. Updated accordingly.

## Release Workflow Diff

Added `HOMEBREW_TAP_GITHUB_TOKEN` to `.github/workflows/release.yml` GoReleaser step:

```yaml
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

## CLI Output — `make snapshot`

```
$ goreleaser release --snapshot --clean

  • skipping announce, publish, and validate...
  • cleaning distribution directory
  • loading environment variables
  • getting and validating git state
    • git state   commit=e0f8602 branch=feat-brew-install current_tag=v2.10.1 previous_tag=v2.10.0 dirty=true
  • parsing tag
  • setting defaults
    • brews is being phased out in favor of homebrew_casks
  • snapshotting
    • building snapshot...   version=2.10.2-next
  • running before hooks
  • building binaries (5 targets)
  • archives (5 archives)
  • calculating checksums
  • homebrew formula
    • writing   formula=dist/homebrew/Formula/gws.rb
  • release succeeded after 20s
```

Snapshot completed successfully. Formula written to `dist/homebrew/Formula/gws.rb`.

## Generated Formula Verification

The generated `dist/homebrew/Formula/gws.rb` contains:

- `desc` — "A lightweight CLI tool for discovering, organizing, and navigating git repositories"
- `homepage` — "https://github.com/daileyo/gws"
- `license` — "MIT"
- `url` — Points to release tar.gz archives for macOS (amd64/arm64) and Linux (amd64/arm64)
- `sha256` — Computed checksums for each archive
- `install` — `bin.install "git-workspace"`

## Deprecation Note

GoReleaser v2.10+ shows: `brews is being phased out in favor of homebrew_casks`. The `brews` stanza works correctly in v2.14.1 but should be migrated to `homebrew_casks` in a future release cycle.
