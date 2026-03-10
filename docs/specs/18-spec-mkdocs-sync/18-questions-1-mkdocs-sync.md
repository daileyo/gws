# 18 Questions Round 1 - MkDocs Sync

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Scope of Update

The gap analysis found 24 discrepancies across critical/high/medium/low priority. How comprehensive should this update be?

- [X] (A) Full update — address all 24 items across all priority levels (critical, high, medium, low)
- [ ] (B) Critical + High only — fix the 12 most impactful items (wrong flags, missing commands, misleading examples)
- [ ] (C) Critical only — fix the 5 completely missing items (parent command, list dual-purpose flags, remote/verbose/visibility/workers/color flags)
- [ ] (D) Other (describe)

## 2. Homebrew Install Reference

The getting-started.md currently documents `brew install daileyo/gws/git-workspace`. Per project notes, the Homebrew tap is in-progress on `feat-brew-install` and shouldn't be referenced until it lands. What should we do?

- [ ] (A) Remove the Homebrew install section entirely from getting-started.md (add back when brew work lands)
- [ ] (B) Keep it but add a "coming soon" note indicating it's not yet available
- [ ] (C) Leave it as-is (it's fine for now)
- [X] (D) Other (describe) The homebrew work is compleate and is actively in use. I test regularly using brew install. any notes implying it is in progress should be corrected.

## 3. README.md Updates

The README.md is also out of sync (missing user/tag/parent commands, subcommand structure, etc.). Should this spec include README updates, or keep it docs-site-only?

- [X] (A) Include README.md updates in this spec (keep README and docs site in sync)
- [ ] (B) Docs site only — README is a separate concern
- [ ] (C) Other (describe)

## 4. Legacy Flags Page

The legacy flags table is missing ~8 deprecated user-related flags. How should we handle this?

- [X] (A) Add all missing deprecated flags to the legacy table for completeness
- [ ] (B) Remove the legacy flags page entirely (users should just use `--help`)
- [ ] (C) Keep current entries but add a note saying "see `gws --help` for complete deprecated flag list"
- [ ] (D) Other (describe)

## 5. Example Output in Docs

Many docs pages show example CLI output that is now wrong (e.g., getting-started default output, list flag behavior, status icon order). Should we regenerate all examples from actual CLI output?

- [X] (A) Yes — all example output blocks should be regenerated from actual v2.17.0 CLI output
- [ ] (B) Update only the examples that are clearly wrong; leave cosmetic differences
- [ ] (C) Other (describe)

## 6. Parent Command Documentation

The `gws parent` command is completely undocumented. Where should it live in the docs?

- [ ] (A) Add to commands-core.md alongside list, init, add, refresh, print-workspace
- [ ] (B) Add to shell-integration.md since its primary use is navigation via the shell function
- [ ] (C) Both — brief entry in commands-core.md, detailed usage in shell-integration.md
- [X] (D) Other (describe) It should be documented; but parent (-p) is a navigation feature that allows for navigating to the parent of a specified repo. the scope of the documentation should be limited to this.

## 7. Proof Artifacts

What would demonstrate the docs are correct and complete?

- [ ] (A) Side-by-side comparison of each docs page section vs actual `--help` output for that command
- [ ] (B) Local MkDocs build (`python -m mkdocs build`) with no warnings + visual review of deployed site
- [C] (C) Both A and B
- [ ] (D) Other (describe)
