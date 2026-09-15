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

**Measured baseline.** Re-measured 2026-09-14 at `c4641a9` (v2.22.0). Commands are recorded
so the counts can be re-derived rather than trusted.

```bash
# A. module/import references
grep -rno 'daileyo/gws' --include='*.go' --include='*.yml' --include='*.json' \
  --include='Makefile' . --exclude-dir=.git | wc -l

# B. raw repo-wide name references
grep -rno 'git-workspace' --include='*.go' --include='*.md' --include='*.yml' \
  --include='Makefile' . --exclude-dir=.git | wc -l

# C. live references (excluding historical records)
grep -rno 'git-workspace' --include='*.go' --include='*.md' --include='*.yml' \
  --include='Makefile' . --exclude-dir=.git --exclude=CHANGELOG.md --exclude-dir=specs | wc -l

# D. command directory size
ls cmd/git-workspace/*.go | wc -l
```

| Measure | 2026-09-03 | 2026-09-14 |
| --- | --- | --- |
| A. `daileyo/gws` references | 89 | **105** |
| B. raw `git-workspace` references | 604 | **745** |
| C. **live** `git-workspace` references | not measured | **110** |
| D. files in `cmd/git-workspace/` | 36 | **39** |

**The raw number is misleading, and the first draft of this spec was wrong about why.**
It claimed the 604 was "dominated by `CHANGELOG.md`". It is not — `CHANGELOG.md` contains
**3** occurrences. The bulk lives in `docs/specs/` (**632** of 745): the historical spec
documents, proofs, and task lists accumulated across specs 01-24.

Those should not be rewritten. A spec describing what was built in March is a record of that
moment, exactly like a changelog entry; rewriting it to say `omgitworks` would make it claim
something untrue. The same reasoning applies to `CHANGELOG.md`.

So the real rename surface is **110 live references**, not 745 — roughly one seventh of the
raw count.

**Live surface by file:**

| File | Count | Class | Affects installed users? |
| --- | --- | --- | --- |
| `docs/site/shell-integration.md` | 43 | cosmetic | no |
| `docs/site/getting-started.md` | 18 | cosmetic | no |
| `.goreleaser.yml` | 16 | **irreversible** (published artifacts) | **yes** |
| `cmd/git-workspace/shellinit.go` | 10 | **breaking** | **yes** |
| `Makefile` | 4 | cosmetic | no |
| `cmd/git-workspace/main.go` | 3 | cosmetic (help text) | no |
| `README.md` | 3 | cosmetic | no |
| `.github/workflows/ci.yml` | 3 | cosmetic | no |
| `docs/site/index.md` | 2 | cosmetic | no |
| `docs/site/configuration.md` | 2 | cosmetic | no |
| `cmd/git-workspace/main_test.go` | 2 | cosmetic | no |
| `cmd/git-workspace/cd.go` | 2 | cosmetic | no |
| `mkdocs.yml` | 1 | cosmetic (site identity) | no |
| `docs/site/commands-core.md` | 1 | cosmetic | no |

Plus the structural touchpoints, which are counted separately because renaming them is not a
text substitution:

| Touchpoint | Location | Class | Affects installed users? |
| --- | --- | --- | --- |
| Go module path | `go.mod:1` | **breaking** (for importers) | no |
| Import statements | 105 sites across `cmd/`, `internal/` | breaking (mechanical) | no |
| Command directory | `cmd/git-workspace/` (39 files) | cosmetic (git mv) | no |
| Binary name | `Makefile:4`, `.goreleaser.yml:16` | **breaking** | **yes** |
| goreleaser project | `.goreleaser.yml:6` | **irreversible** | **yes** |
| Brew formula | `.goreleaser.yml:55` | **irreversible** | **yes** |
| Brew tap repo | `daileyo/homebrew-gws` | **irreversible** | **yes** |
| Completion function names | `shellinit.go` (`_git-workspace`, `__start_git-workspace`) | **breaking** | **yes** |
| Repository name | GitHub `daileyo/gws` | **irreversible** | no (redirects) |
| Docs site URL | `mkdocs.yml:2` | **irreversible** | no |
| Logo assets | `gws-logo.png`, `docs/site/assets/images/` | cosmetic | no |
| Release automation | `release-please-config.json`, manifest | cosmetic | no |

**Six touchpoints affect already-installed users.** Those, and only those, drive the
compatibility guarantee: the binary name, the completion function names, the shell-init
template, the goreleaser project and archive names, the brew formula, and the brew tap.

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

| Name | Current | Proposed | Firm? | Rationale |
| --- | --- | --- | --- | --- |
| Brand | git-workspace | `omgitworks` | **firm** | Distinctive; clear in all software channels; the incumbent is taken three times over in this niche |
| Command | `gws` | `omgw` | **firm** | Four characters; clear on every channel checked and on PATH; derives from the brand |
| Binary | `git-workspace` | `omgitworks` | **firm** | Matches the brand; `shell-init` means users rarely type it |
| Repository | `daileyo/gws` | `daileyo/omgitworks` | provisional | GitHub redirects the old path, but see the namespace-reuse risk below |
| Module path | `github.com/daileyo/gws` | `github.com/daileyo/omgitworks` | provisional | Follows the repo; 105 import sites, mechanical but wide |
| Brew tap | `daileyo/homebrew-gws` | new tap under the new name | deferred | Wanted, but as a follow-up project; the existing tap keeps serving users until then |

**PATH collision check for `omgw` (task 2.4), 2026-09-14.** The registry checks in the
decision record do not cover whether something already ships a binary by that name:

| Check | Result |
| --- | --- |
| `command -v omgw` on a normal Linux PATH | not found |
| Debian stable package contents, exact filename `omgw` | "Sorry, your search gave no results" |
| homebrew-core formula and cask `omgw` | 404 / 404 |
| Shell builtin or keyword collision | none |

**Why brand and command are decided separately.** `shell-init` already generates a shell
function whose name is independent of the binary it invokes: users type `gws`, the binary is
`git-workspace`. That decoupling was built for convenience, and it is what makes this rename
survivable — the binary can be renamed without a single user changing a single keystroke, and
the long form can carry the brand while a short form carries the typing.

**The mnemonic, which is the point of the long form.** `omgitworks` is a double reading and
both are intended: *"OMG, it works"* — the reaction — and *"om git works(pace)"* — literally
what the tool is. A four-character command alone would be forgettable and unsearchable; the
long form is what makes the project findable, which is the entire motivation for leaving
`git-workspace` behind.

**Proof Artifacts:**

- Documentation: A decision table covering all five names demonstrates the decision is complete rather than partial

### Unit 3: Staged Migration Plan

**Purpose:** Sequence the work so each stage is independently revertible and no stage leaves an installed user broken.

**Functional Requirements:**

- The plan shall order stages from lowest to highest irreversibility
- Each stage shall state its revert procedure
- The plan shall define the compatibility guarantee: which old names keep working and for how long
- The plan shall identify the point of no return and the preconditions for crossing it

**Staging.** Ordered by irreversibility, lowest first. Each stage states how to undo it.

**Stage 1 — Brand surface.** README, docs site content, logo assets, `mkdocs.yml` `site_name`.
Establishes the name publicly before any technical commitment.
*Revert:* `git revert`. Nothing installed changes; no published artifact moves.

**Stage 2 — Shell-integration compatibility.** Teach `shell-init` to emit **both** the `gws`
and `omgw` functions, and to register completions for both names, while the binary is still
called `git-workspace`. This must ship and reach users **before** stage 3.
*Revert:* `git revert`; users keep `gws` either way.

**Stage 3 — Binary and archive naming.** `Makefile` `BINARY_NAME`, `.goreleaser.yml`
`project_name` and `binary`. Ship `git-workspace` as a symlink alongside `omgitworks` for the
compatibility window, so an rc file containing `eval "$(git-workspace shell-init zsh)"` keeps
working. Update the hardcoded completion function names in the same commit — see the lockstep
note below.
*Revert:* possible but noisy — a released archive cannot be unpublished cleanly, so reverting
means a follow-up release restoring the old names.

**Stage 4 — Repository rename.** `daileyo/gws` → `daileyo/omgitworks`.
*Revert:* rename back; redirects follow. Safe **only** while nothing occupies the old path.

**Stage 5 — Go module path.** `go.mod` plus 105 import sites. Mechanical but wide.
*Revert:* `git revert`; the old path resumes working since GitHub redirects git operations.

**Stage 6 — Brew formula and tap.** Last, because the formula is how users upgrade; renaming
mid-flight can strand someone between versions.
*Revert:* re-publish the formula under the old name; the old tap must not be deleted until
this stage is proven.

**Ordering check (task 3.2).** Stage 2 must precede stage 3, and it does. If the binary were
renamed first, a user whose rc file calls `git-workspace shell-init` would get
"command not found" on **every new shell** until they edited their rc by hand. The symlink in
stage 3 is the second line of defense; shipping stage 2 first is the first.

**Compatibility guarantee.** To be confirmed before execution, but the proposal is:

| Guarantee | Duration |
| --- | --- |
| `shell-init` emits a `gws` function | indefinitely — it costs nothing and breaking it buys nothing |
| `shell-init` accepts and emits `omgw` | from stage 2 onward |
| `git-workspace` ships as a symlink to `omgitworks` | through the next major version, then removed with a release note |
| Old brew tap remains installable | until one release cycle after stage 6 |

**Point of no return.** Stage 4 (repository rename) combined with stage 6 (brew tap rename).
Everything before is `git revert` plus a patch release. After stage 4, third parties may link
to the new name; after stage 6, users' installed taps point at it.

**Lockstep requirement (task 3.7).** `shellinit.go` hardcodes Cobra's generated completion
function names, which derive from the binary name:

```
shellinit.go:92   compdef _git-workspace gws
shellinit.go:163  complete -o default -F __start_git-workspace gws
```

Renaming the binary changes what `completion zsh` / `completion bash` generate. If these two
lines are not updated in the same commit, **tab completion silently stops working** — no
error, it just does nothing, which is the hardest kind of breakage to notice.

**Proof Artifacts:**

- Documentation: An ordered stage list with per-stage revert procedures demonstrates the migration is planned rather than improvised
- Documentation: A stated compatibility guarantee demonstrates existing users are accounted for

### Unit 4: Preconditions and Go/No-Go

**Purpose:** Name what must be true before execution begins, so the decision is not made by momentum.

**Functional Requirements:**

- The document shall list the preconditions for execution, each independently checkable
- The document shall state what would constitute a reason to abandon the rename
- The document shall record that no execution occurs under this spec

**Preconditions.** Each is independently checkable. Status as of 2026-09-14.

| # | Precondition | Status |
| --- | --- | --- |
| P1 | Short command name resolved | **met** — `omgw` (task 2.2) |
| P2 | `omgw` free on PATH, not just in registries | **met** — task 2.4 |
| P3 | `daileyo/omgitworks` available as a repo name | **met** — GitHub API returns 404 |
| P4 | Specs 22 and 23 landed | **met** — both merged; 22 released in v2.21.0, 23 in v2.22.0 |
| P5 | Availability re-verified across all target channels before execution | **met 2026-09-15, decays** — see below |
| P6 | Trademark sanity check recorded | **met, with caveats** — see below |
| P7 | Decision on the dormant `OMGItworks` GitHub username | **met** — not pursued |
| P8 | Decision on registering `omgitworks.dev` | **met** — registered 2026-09-15 |
| P9 | Decision on defensive registry registration | **met** — declined; not a Go channel |
| P10 | Decision on renaming the brew tap | **met** — deferred to a follow-up project |
| P11 | Decision on placeholder repo at the old path | **met** — no placeholder, path left empty |
| P12 | Decision on major vs minor release | **met** — major version |

**P5 — availability re-check, 2026-09-14** (previous check 2026-09-04; both names still clear):

| Channel | `omgitworks` | `omgw` |
| --- | --- | --- |
| npm / crates.io / PyPI | available | available |
| homebrew-core | available | available |
| chocolatey | available | available |
| scoop (Main) | available | available |
| winget | available | available |

Domains: `omgitworks.dev` and `.io` remain unregistered (NXDOMAIN); `.com` is still parked on
GoDaddy nameservers; `.co.uk` still belongs to the UK business.

**P6 — trademark sanity check, and its limits.** A general web search surfaced **no registered
mark** for "OMG IT Works". That is a sanity check, **not a clearance search**, and this spec
should not be read as a legal opinion.

What is known: the UK business at `omgitworks.co.uk` operates in IT recycling and secure data
destruction. Under the Nice classification those sit in classes such as 40 (treatment of
materials) and 37 (repair), whereas a developer CLI sits in 9 (software) and 42 (software
development services). Different classes, different customers, no plausible confusion between
a Cornwall recycling firm and a git workspace manager.

What is not known: whether that business holds an unregistered mark with acquired goodwill,
which in the UK can support a passing-off claim regardless of registration. For a free,
open-source tool the practical risk is low. If this ever becomes commercial, a proper clearance
search of the UK IPO register and USPTO is warranted first.

**Abandonment criteria (task 4.9).** Reasons to keep `git-workspace` instead:

1. Either name becomes unavailable on homebrew-core or winget before execution — the whole
   point was an unclaimed distribution surface.
2. The UK business is found to hold a registered mark covering software.
3. The `omgw` command is found to collide with something already installed widely.
4. The docs-site URL break (see Technical Considerations) is judged to cost more than the
   ambiguity of the current name, and no custom domain is acquired.
5. Project priorities shift such that nobody is available to shepherd a six-stage migration —
   a half-finished rename is worse than either endpoint.

**Decisions recorded 2026-09-15.** All previously open preconditions are now settled.

| # | Decision | Outcome |
| --- | --- | --- |
| P7 | Dormant `OMGItworks` GitHub username | **Not pursued.** The repo lives at `daileyo/omgitworks`, which does not need the org-level name. GitHub does not release usernames on request outside trademark disputes, so chasing it is slow and unlikely to succeed. |
| P8 | Register `omgitworks.dev` | **Done.** The domain is registered and on Cloudflare nameservers as of 2026-09-15. This is what protects the docs site across the repository rename, since GitHub Pages URLs do not redirect. |
| P9 | Defensive registration on npm / PyPI / crates.io | **Declined.** This is a Go project; Go has no central package registry, so modules resolve directly from the VCS path. npm (JavaScript), PyPI (Python), and crates.io (Rust) are not distribution channels for it, and registering placeholder packages there is noise. Accepted risk: a squatter could take those names later. It costs the project nothing, because it never intended to publish there. |
| P10 | Rename the brew tap | **Deferred to a follow-up project.** A tap under the new name is wanted, but creating it is separate work from the rename itself. `daileyo/homebrew-gws` keeps serving users until then. |
| P11 | Placeholder repo at the old path | **No placeholder.** Creating anything at `daileyo/gws` after the rename destroys GitHub's redirects. The correct action is to leave the path empty and never reuse it. |
| P12 | Major or minor release | **Major version.** No user-facing command changes, but the binary name and module path do. |

**Distribution channels (new requirement, 2026-09-15).** The brand should be registered and
installable on Homebrew, winget, Scoop, pacman (AUR), and apt, plus other channels as they
make sense. Availability verified for both names on 2026-09-15:

| Channel | `omgitworks` | `omgw` |
| --- | --- | --- |
| homebrew-core | available | available |
| winget (`microsoft/winget-pkgs`) | available — no `omg*` publisher | available |
| Scoop (Main, Extras) | available | available |
| chocolatey | available | available |
| AUR (pacman) | available | available |
| Debian source package name | available | available |

**This expands beyond the rebrand and belongs in its own spec.** Publishing to these channels
is packaging work, not renaming work, and each has its own submission process and review
latency — winget and Scoop take pull requests against their manifest repositories, AUR needs a
`PKGBUILD` and a maintainer account, and Debian or Ubuntu proper requires a sponsoring
maintainer, which in practice means shipping a `.deb` from releases or maintaining a PPA
instead. Attaching that to the rename would make a six-stage migration into something much
larger and slower.

What the rebrand spec owes it is the names: both are confirmed free everywhere above, and the
re-check in P5 should cover this full list rather than just the registries, since these are now
the channels that matter.

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
- **GitHub repository rename**: GitHub redirects web traffic and git operations (clone, fetch, push) from the old path, which makes stage 4 less risky than it appears. Two documented exceptions matter here:
    - **GitHub Pages URLs are not redirected.** The docs site is published at `https://daileyo.github.io/gws` (`mkdocs.yml:2`). Renaming the repository moves it to `/omgitworks` and **the old URL breaks outright** — existing links, bookmarks, and search results die with it. A custom domain registered *before* stage 4 would insulate the site from this rename and any future one. This is the strongest practical argument for buying `omgitworks.dev`.
    - **Redirects lapse if the old path is reoccupied.** If anything is later created at `daileyo/gws`, every redirect stops working. See the placeholder decision in Preconditions.
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
