# 24-spec-rebrand-omgitworks

## Introduction/Overview

`git-workspace` is a crowded name. At least three other actively-published CLI tools occupy it
in this exact problem space, and the short form `gws` is already a homebrew-core formula for a
tool that manages workspaces composed of git repositories — the same one-line description this
project uses. That is not a hypothetical future collision; it is the current state of the
ecosystem, and it costs this project discoverability every time someone searches for it.

`omgitworks` is the candidate replacement. This spec records the research behind that choice,
inventories every place the current name is embedded, and lays out a staged migration plan.

**This spec is deliberately not executable.** Per round 1, it produces the decision record and
the plan; no code, configuration, or documentation is renamed until the plan is explicitly
approved. The deliverable is a document that makes the decision reversible and, if approved,
makes execution mechanical.

## Goals

- Record a durable decision record for the name change, with the supporting research and its date
- Establish availability across the package managers the tool ships through, or may ship through — homebrew, winget, scoop, chocolatey, AUR — as a first-class selection criterion
- Inventory every touchpoint where `git-workspace`, `gws`, or `daileyo/gws` is embedded, with counts and file locations
- Define what the binary, the command, the module path, and the repository are each called after the change
- Stage the migration so each phase is independently revertible and no phase leaves users stranded
- Define the compatibility guarantee for existing users — what keeps working, for how long
- Identify the irreversible steps and what must be true before taking them
- Produce **no** renames as part of this spec

## User Stories

- **As the maintainer**, I want the name decision written down with its evidence so that six months from now I can tell whether it was a good call, and so I do not re-research the same registries.
- **As the maintainer**, I want a complete touchpoint inventory before committing, so I can see the true cost of the rename rather than discovering it halfway through.
- **As an existing user**, I want my installed `gws` command and my shell config to keep working across the rename, rather than silently breaking on an upgrade.
- **As a prospective user**, I want searching for the tool to find this tool, rather than three unrelated projects with the same name.

## Decision Record

### Research conducted 2026-09-03

**Candidate: `omgitworks`**

| Surface | Status | Detail |
| --- | --- | --- |
| npm | Available | 404 on `omgitworks` and `omgitworks-cli` |
| crates.io | Available | 404 |
| PyPI | Available | 404 |
| homebrew-core | Available | No formula |
| GitHub repo name | Effectively clear | One hit: `ThatOctopus/omgitworks`, 0 stars, no pushes since 2016-03-10 |
| GitHub username/org | **Taken** | `github.com/OMGItworks` — created 2021-06-25, 0 public repos, no activity |
| `omgitworks.com` | **Registered** | Parked on GoDaddy nameservers |
| `omgitworks.dev` | Available | NXDOMAIN |
| `omgitworks.io` | Available | NXDOMAIN |
| `omgitworks.sh` | Available | NXDOMAIN |
| `omgitworks.co.uk` | **Taken** | OMG IT Works, a Cornwall-based IT recycling and data-destruction business |

**Package-manager availability, checked 2026-09-04**

Availability across distribution channels is a primary selection criterion, not an
afterthought: a name that cannot be claimed on winget or homebrew is a name that costs the
project an install path.

| Channel | `omgitworks` | `omgw` | `gws` (incumbent) |
| --- | --- | --- | --- |
| homebrew-core | Available | Available | **Taken** — `streakycobra/gws` |
| winget (`microsoft/winget-pkgs`) | Available — no `omg*` publisher among 281 in the `o/` namespace | Available | Not checked; moot |
| scoop (Main + Extras buckets) | Available | Available | — |
| chocolatey | Available | Available | — |
| AUR | Available | Available | **Taken** — `gws`, "Colorful KISS helper for git workspaces" |
| npm | Available | Available | **Taken** |
| crates.io | Available | Available | — |
| PyPI | Available | Available | — |

Both candidate names are clear on every channel checked. The incumbent short name is taken on
three of them by the same unrelated project.

**Mnemonic and etymology**

The name is a double reading, and both readings are intended:

- **"OMG, it works"** — the reaction the tool is named for
- **"om git works(pace)"** — literally *git workspace*, which is what it is

The long form is what makes the name memorable and searchable; the short form is what gets
typed. This is the reason the brand and the command are decided separately below.

**Incumbent: `git-workspace` / `gws`**

| Collision | Detail |
| --- | --- |
| `orf/git-workspace` | "Sync personal and work git repositories from multiple providers" |
| `mariocasciaro/git-workspace` | "CLI util to keep multiple projects in sync with different remote git repos" |
| `c0fec0de/git-ws` | "Git Workspace - Multi Repository Management Tool", also published to PyPI as `git-ws` |
| homebrew-core `gws` | "Manage workspaces composed of git repositories" — `streakycobra/gws`, v0.2.0 |
| AUR `gws` | "Colorful KISS helper for git workspaces" — StreakyCobra again, packaged for Arch |
| GitKraken Workspaces | Commercial product using the same term for the same concept |

### Assessment

The case for changing is stronger than the case for the specific replacement. `git-workspace`
is genuinely ambiguous in its own niche, and `brew install gws` today installs a different
tool that does a similar job. This project's own formula is `daileyo/gws/git-workspace`, so
there is no hard install conflict — but the discovery confusion is real and permanent.

`omgitworks` is clear in every channel that matters for a Go CLI: all three major package
registries, homebrew-core, and the practical GitHub namespace. The two encumbrances are
soft ones. The GitHub username is held by a dormant account with no public repositories, which
constrains an org name but not a repo name under the existing `daileyo` account. The UK
business operates in IT recycling — a different class of goods and services entirely — so
brand confusion with a developer CLI is unlikely, though it does mean the obvious `.com` and
`.co.uk` are unavailable and `omgitworks.dev` is the natural home.

The genuine cost is ergonomic: `omgitworks` is ten characters, versus three for `gws`. Nobody
will type it — and nobody has to. `shell-init` already decouples the typed command from the
binary name, so the long form can carry the brand while a short form carries the typing.

**`omgw` is the chosen command name.** It is four characters, available on every package
manager and registry checked, and derives transparently from the brand. The pairing gives the
project a distinctive, searchable name and a short thing to type, which is the combination the
incumbent `git-workspace`/`gws` pair fails at from the other direction: `gws` is short but
taken, and `git-workspace` is descriptive but generic.

## Demoable Units of Work

### Unit 1: Touchpoint Inventory

**Purpose:** Establish the true, measured scope of the rename before anyone commits to it.

**Functional Requirements:**

- The inventory shall enumerate every category of touchpoint with file paths and occurrence counts, measured rather than estimated
- The inventory shall separate touchpoints by risk class: **irreversible** (published artifacts, repository rename), **breaking** (module path, binary name), and **cosmetic** (docs, help text)
- The inventory shall explicitly identify which touchpoints affect already-installed users
- The inventory shall record the measurement commands so the counts can be reproduced as the codebase changes

**Measured baseline as of 2026-09-03 (commit `bd8b210`, v2.20.0):**

| Category | Location | Count / Detail |
| --- | --- | --- |
| Go module path | `go.mod:1` | `github.com/daileyo/gws` |
| Import statements | across `cmd/`, `internal/` | 89 occurrences of `daileyo/gws` in `.go`/`.yml`/`.json`/`Makefile` |
| Command directory | `cmd/git-workspace/` | 36 Go files |
| Binary name | `Makefile:4`, `.goreleaser.yml:16` | `BINARY_NAME=git-workspace` |
| goreleaser project | `.goreleaser.yml:6` | `project_name: git-workspace` |
| Brew formula | `.goreleaser.yml:55` | `name: git-workspace` in tap `daileyo/homebrew-gws` |
| Brew tap repo | `.goreleaser.yml:56-58` | `daileyo/homebrew-gws` |
| Shell templates | `shellinit.go` | `{BIN}` placeholder plus `_git-workspace` / `__start_git-workspace` completion function names |
| Docs site | `mkdocs.yml`, `docs/site/*.md` | 8 markdown files; `site_name`, `site_url`, `repo_url`, `repo_name` |
| Docs site URL | `mkdocs.yml:2` | `https://daileyo.github.io/gws` |
| Logo assets | `gws-logo.png`, `docs/site/assets/images/` | logo, nav logo, favicon |
| README | `README.md` | title, install instructions, examples |
| CI | `.github/workflows/` | build and release workflows |
| Release automation | `release-please-config.json`, `.release-please-manifest.json` | Go release type |
| Total textual | repo-wide | 604 occurrences of `git-workspace` across `.go`/`.md`/`.yml`/`Makefile` |

**Proof Artifacts:**

- Documentation: The inventory table above, with reproducible measurement commands, demonstrates scope is known
- Documentation: Risk classification separating irreversible from cosmetic changes demonstrates the decision is informed

### Unit 2: Naming Decisions

**Purpose:** Decide what each distinct name becomes, since "the project name" is really five separate names that need not all change together.

**Functional Requirements:**

- The document shall decide, for each of the following, what it becomes and why:
  - **Project / brand name** — what appears in the README title, docs site, and logo
  - **Repository name** — the GitHub repo path
  - **Go module path** — the import prefix
  - **Binary name** — what ships in the release archive and lands on `PATH`
  - **Command name** — what the user actually types, i.e. the shell function from `shell-init`
- The document shall state explicitly that the command name and the binary name need not match, since `shell-init` already generates a `gws` function that wraps a `git-workspace` binary — this decoupling already exists and is the mechanism that makes a rename survivable
- The document shall record the recommendation and its rationale for each
- The command name shall be `omgw`, and the document shall state how long `gws` continues to be emitted alongside it

**Recommendation:**

| Name | Current | Proposed | Rationale |
| --- | --- | --- | --- |
| Brand | git-workspace | omgitworks | Distinctive; clear in all software channels |
| Repository | `daileyo/gws` | `daileyo/omgitworks` | GitHub redirects the old path indefinitely |
| Module path | `github.com/daileyo/gws` | `github.com/daileyo/omgitworks` | Follows the repo; a `retract`-free major-version-free rename is possible since v2 tags are already in use |
| Binary | `git-workspace` | `omgitworks` | Matches the brand; users rarely type it directly |
| Command | `gws` | `omgw` | Four characters; available on every channel checked; derives from the brand. `gws` continues to be emitted by `shell-init` for compatibility |

**Proof Artifacts:**

- Documentation: A decision table covering all five names demonstrates the decision is complete rather than partial

### Unit 3: Staged Migration Plan

**Purpose:** Sequence the work so each stage is independently revertible and no stage leaves an installed user broken.

**Functional Requirements:**

- The plan shall order stages from lowest to highest irreversibility
- Each stage shall state its revert procedure
- The plan shall define the compatibility guarantee: which old names keep working and for how long
- The plan shall identify the point of no return and the preconditions for crossing it

**Proposed staging:**

1. **Documentation and brand surface** — README, docs site, logo, `mkdocs.yml`. Fully revertible; nothing installed changes. Establishes the name publicly before any technical commitment.
2. **Binary and archive naming** — `Makefile`, `.goreleaser.yml` `project_name` and `binary`. Ships a differently-named binary; the brew formula must install both names, or the old one as a symlink, for at least one release cycle.
3. **Shell integration compatibility** — `shell-init` continues to define `gws` regardless of the binary name, and additionally defines the new command name if one is chosen. This is the stage that protects existing users, and it must ship *before* stage 2 reaches anyone.
4. **Repository rename** — GitHub redirects the old path indefinitely for both web and git operations, so this is safer than it appears, but it does break anything pinned to the old URL that does not follow redirects.
5. **Go module path** — the genuinely breaking change for anyone importing the packages. Requires updating 89 import sites plus `go.mod`. Given this is a CLI rather than a library, the practical blast radius is small.
6. **Brew formula rename** — last, because the formula is how users upgrade, and a rename mid-flight can strand someone between versions.

**Compatibility guarantee to define:** how long `shell-init` keeps emitting a `gws` function, and whether the old binary name ships as a symlink or is dropped at a stated version.

**Proof Artifacts:**

- Documentation: An ordered stage list with per-stage revert procedures demonstrates the migration is planned rather than improvised
- Documentation: A stated compatibility guarantee demonstrates existing users are accounted for

### Unit 4: Preconditions and Go/No-Go

**Purpose:** Name what must be true before execution begins, so the decision is not made by momentum.

**Functional Requirements:**

- The document shall list the preconditions for execution, each independently checkable
- The document shall state what would constitute a reason to abandon the rename
- The document shall record that no execution occurs under this spec

**Preconditions:**

- ~~The short-command-name question is resolved~~ — **done**, `omgw`; remaining check is that it collides with no command already on a default PATH
- Package-manager availability for `omgitworks` and `omgw` is re-verified immediately before execution, since the 2026-09-04 results decay
- `omgitworks.dev` is registered, or a decision is made to forgo a domain
- The GitHub repo name `daileyo/omgitworks` is confirmed available
- A decision is recorded on whether to pursue the dormant `OMGItworks` GitHub username
- Specs 22 and 23 have landed, so the rename does not collide with in-flight work
- A trademark sanity check is done on the UK business, confirming no overlap in class of goods

**Proof Artifacts:**

- Documentation: A precondition checklist demonstrates execution is gated
- Documentation: An explicit "no changes made under this spec" statement demonstrates scope is respected

## Non-Goals (Out of Scope)

1. **Any actual rename.** No file, module path, binary, repository, or published artifact is renamed under this spec. This is the defining constraint.
2. **Registering domains or GitHub namespaces.** These are the maintainer's decisions to make outside the repo; the spec only records that they are preconditions.
3. **Trademark filing.** Out of scope; a sanity check is a precondition, a filing is not.
4. **Choosing between alternative names.** `omgitworks` is the candidate under evaluation. Broadening to a general naming exercise is a different piece of work.
5. **Deprecating the `gws` command.** Whatever else changes, `shell-init` continues to define `gws`; removing it would be a separate, later decision.
6. **Rewriting git history** to change the module path retroactively.

## Design Considerations

The single most useful fact about this codebase's rebrand-ability is that **the command name
is already decoupled from the binary name**. `shell-init` generates a shell function called
`gws` that invokes a binary called `git-workspace`. Users type `gws`; the binary's name is an
implementation detail they see only at install time. That decoupling was built for convenience
and it turns out to be the mechanism that makes a rename nearly free at the user-facing layer:
the binary can be renamed without a single user changing a single keystroke.

This is why staging matters more than speed. The stages that touch what users type are cheap
and revertible; the stages that touch published artifacts are neither. Doing them in the wrong
order — renaming the brew formula before the shell compatibility shim ships, say — creates a
window where an upgrading user has neither the old command nor the new one.

## Repository Standards

- Specs live under `docs/specs/NN-spec-<slug>/` with spec, tasks, questions, and a proofs directory
- Documentation is Markdown, rendered by mkdocs with the material theme
- Conventional commits (`docs: ...` for everything in this spec)
- Release automation is release-please with goreleaser; version tags are `vN.N.N`

## Technical Considerations

- **Go module rename**: Go does not require a major-version bump for a module path change, but every importer breaks at once. Since this is a CLI and not a published library, the practical impact is limited to this repo's own 89 import sites.
- **GitHub repository rename**: GitHub maintains redirects from the old repo path indefinitely, for web, API, and git remote operations. This makes stage 4 far less risky than it appears. Redirects do *not* apply if someone later creates a new repo at the old path.
- **pkg.go.dev**: The old module path remains indexed. A `retract` directive in `go.mod` is not appropriate for a rename; the old path simply stops receiving updates.
- **Homebrew tap**: The tap repository `daileyo/homebrew-gws` would ideally be renamed too, which changes the tap name users type (`brew tap daileyo/gws`). GitHub redirects cover the git operation, but `brew tap` output would show the old name. Consider whether the tap rename is worth the churn.
- **Completion function names**: `shellinit.go` templates hardcode `_git-workspace` and `__start_git-workspace`, which are derived from the binary name by Cobra. These change automatically with the binary rename but the templates reference them literally, so they must be updated in lockstep or completion silently breaks.
- **Existing installs**: A user who has `eval "$(git-workspace shell-init zsh)"` in their rc file will get a "command not found" on shell startup the moment the binary is renamed, unless the old name ships as a symlink for a transition period. This is the most likely way to break people and deserves an explicit decision.
- **The 604 textual occurrences** are dominated by `CHANGELOG.md`, which is historical and should not be rewritten. The live count is much smaller and should be re-measured excluding the changelog before execution.

## Security Considerations

- **Namespace squatting after release**: once the old GitHub repo path is vacated, another party could create a repo there. GitHub's redirect stops working the moment they do. Consider retaining a placeholder repository at the old path rather than leaving it open.
- **Package registry squatting**: `omgitworks` is currently free on npm, PyPI, and crates.io. Publishing the name publicly before claiming those namespaces invites a squatter to take them, even though this project does not ship to any of them. Defensive registration is cheap; the decision should be conscious.
- **Install-instruction integrity**: any rename changes the release-artifact URLs users curl. Checksum verification via `checksums.txt` must remain documented and correct throughout the transition, so a user following stale instructions fails loudly rather than fetching an unexpected artifact.
- **Brew tap trust**: renaming the tap changes the path users add to their trusted tap list. The old tap should be left intact until the new one is proven.

## Success Metrics

1. **Decision is documented**: the name choice, its evidence, and its date are recorded and reproducible.
2. **Scope is measured**: the touchpoint inventory is derived from commands run against the repo, not estimated.
3. **Distribution is unblocked**: both the brand and command names are confirmed available on every package manager the project ships through or may ship through.
4. **The plan is revertible**: every stage has a stated revert procedure, and the point of no return is identified.
5. **Existing users are accounted for**: the compatibility guarantee is explicit about what keeps working and for how long.
6. **Nothing was renamed**: `git diff` for this spec touches only `docs/specs/24-spec-rebrand-omgitworks/`.

## Open Questions

1. ~~What does the user actually type?~~ **Resolved 2026-09-04: `omgw`.** Four characters, available on every package manager and registry checked, transparently derived from the brand. `gws` continues to be emitted by `shell-init` during the compatibility window; question 5 covers for how long.
2. Should the dormant `OMGItworks` GitHub username be pursued via GitHub's name-release process, or is a repo under `daileyo` sufficient?
3. Is `omgitworks.dev` worth registering, given `.com` is parked and `.co.uk` belongs to an unrelated business?
4. Should the brew tap `daileyo/homebrew-gws` be renamed, or left as-is to avoid churn in what users type?
5. Should `git-workspace` ship as a symlink alongside the new binary, and if so, until which version?
6. Does the rename wait for a major version bump (v3.0.0) to signal the change, or ship in a minor release since no user-facing command changes?
7. Is the UK business's `omgitworks` usage in a sufficiently different class of goods to pose no trademark risk, and does that need more than a sanity check?
