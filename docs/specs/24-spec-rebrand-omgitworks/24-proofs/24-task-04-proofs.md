# 24 Task 4.0 Proofs - Preconditions and Go/No-Go

## Proof 1: twelve preconditions, six met, six needing the maintainer

| Met | Open (maintainer decision) |
| --- | --- |
| P1 command name resolved (`omgw`) | P7 dormant GitHub username |
| P2 `omgw` free on PATH | P8 register `omgitworks.dev` |
| P3 `daileyo/omgitworks` available (404) | P9 defensive registry registration |
| P4 specs 22 and 23 landed | P10 brew tap rename |
| P5 availability re-verified | P11 placeholder repo at old path |
| P6 trademark sanity check | P12 major vs minor release |

The open six are genuinely the maintainer's calls — they involve spending money, contacting
GitHub support, and setting release policy. Each carries a recommendation so the decision is a
yes/no rather than an open-ended question.

## Proof 2: availability re-verified, with decay acknowledged

Both names remain clear on npm, crates.io, PyPI, homebrew-core, chocolatey, scoop, and winget
as of 2026-09-14. `omgitworks.dev` and `.io` remain unregistered.

The matrix carries its check date because this is exactly the kind of fact that silently goes
stale; P5 requires re-running it immediately before execution.

## Proof 3: the trademark check is honest about what it is not

Recorded as a **sanity check, not a clearance search**, with the limits stated: no registered
mark surfaced, the UK business operates in different Nice classes (40/37 vs 9/42), but an
unregistered mark with acquired goodwill could still support a passing-off claim in the UK. A
proper register search is flagged as necessary only if the project ever becomes commercial.

This is deliberately not written as a legal opinion.

## Proof 4: abandonment criteria exist

Five concrete reasons to keep `git-workspace`, including one uncovered while writing this
spec — the docs-site URL break — and one about execution capacity, since a half-finished
six-stage rename is worse than either endpoint.

## Proof 5: a recommendation that inverts the obvious intuition

P11 asks whether to create a placeholder repo at `daileyo/gws` to prevent namespace reuse. The
intuitive answer is yes. It is wrong: creating **anything** at the old path is precisely what
kills GitHub's redirects. The correct action is to leave it empty and never reuse it. Recorded
explicitly because the intuition runs backwards.

---

## Update 2026-09-15: all twelve preconditions met

The six items that needed the maintainer are decided.

| # | Decision |
| --- | --- |
| P7 | Dormant GitHub username **not pursued** — `daileyo/omgitworks` does not need the org-level name |
| P8 | `omgitworks.dev` **registered**, on Cloudflare nameservers |
| P9 | Defensive npm/PyPI/crates registration **declined** |
| P10 | Brew tap rename **deferred** to a follow-up project |
| P11 | **No placeholder** at the old path |
| P12 | **Major version** release |

### P9 is the one with a real trade-off

Go has no central package registry — modules resolve directly from the VCS path — so npm
(JavaScript), PyPI (Python), and crates.io (Rust) are not distribution channels for this
project. Registering placeholder packages there would be purely defensive.

The accepted risk is that a squatter takes `omgitworks` or `omgw` on those registries later.
That costs this project nothing, because it never intended to publish there. Recorded as a
conscious trade rather than an oversight.

### New requirement: distribution channels

The brand should be registered and installable on Homebrew, winget, Scoop, pacman (AUR), and
apt. Verified 2026-09-15 — both names are free on every one:

| Channel | `omgitworks` | `omgw` |
| --- | --- | --- |
| homebrew-core | available | available |
| winget | available | available |
| Scoop (Main, Extras) | available | available |
| chocolatey | available | available |
| AUR (pacman) | available | available |
| Debian source package name | available | available |

**This belongs in its own spec.** Publishing to these channels is packaging work with its own
submission processes and review latency: winget and Scoop take PRs against manifest
repositories, AUR needs a `PKGBUILD` and maintainer account, and Debian proper requires a
sponsoring maintainer — which in practice means shipping a `.deb` from releases or maintaining
a PPA. Folding that into a six-stage rename would make it much larger and slower.

What spec 24 owes it is the names, and those are confirmed free. P5's re-check now covers this
full channel list rather than just the registries, since these are the channels that matter.
