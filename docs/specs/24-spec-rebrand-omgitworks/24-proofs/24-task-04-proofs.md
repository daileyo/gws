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
