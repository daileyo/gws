# 24 Questions Round 1 - Rebrand to omgitworks

Answers captured from the planning session on 2026-09-03.

## 1. How much should the rebrand spec commit to?

- [X] (A) **Spec it, don't execute.** Produce the full decision record with research evidence,
      a complete touchpoint inventory, and a staged migration and alias plan — but make zero
      code changes until explicitly approved. Keeps the decision reversible.
- [ ] (B) Full rebrand, planned for execution — module path, `cmd/` directory, binary, repo,
      docs site, goreleaser, brew tap, with `gws` kept as the shell function name
- [ ] (C) Brand surface only — README, docs site, logo, release notes, formula; module path and
      repo name unchanged

## 2. Name-saturation research (findings, not a question)

Conducted 2026-09-03 against package registries, GitHub, DNS, and web search. Full results and
sources are recorded in the Decision Record section of the spec. Summary:

- `omgitworks` is **clear** on npm, crates.io, PyPI, and homebrew-core
- The GitHub username `OMGItworks` is **taken** but dormant (0 public repos since 2021)
- `omgitworks.com` is **registered and parked**; `.dev`, `.io`, `.sh` are **available**
- `omgitworks.co.uk` belongs to an unrelated UK IT-recycling business
- The **current** name is the more saturated one: at least three other multi-repo CLI tools
  are named `git-workspace` or `git-ws`, and `gws` is already a homebrew-core formula for a
  tool in the same problem space
