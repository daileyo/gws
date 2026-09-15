# 24 Task 1.0 Proofs - Touchpoint Inventory

Re-measured 2026-09-14 at `c4641a9` (v2.22.0). The 2026-09-03 baseline predates specs 22 and
23, both now merged.

## Proof 1: measurements are reproducible

Commands are recorded in the spec. Results:

| Measure | 2026-09-03 | 2026-09-14 |
| --- | --- | --- |
| A. `daileyo/gws` references | 89 | 105 |
| B. raw `git-workspace` references | 604 | 745 |
| C. **live** `git-workspace` references | not measured | **110** |
| D. files in `cmd/git-workspace/` | 36 | 39 |

## Proof 2: the first draft was wrong about where the bulk lives

The original spec asserted the raw count was "dominated by `CHANGELOG.md`". Measured:

```
$ grep -o 'git-workspace' CHANGELOG.md | wc -l
3

$ grep -rno 'git-workspace' docs/specs --include='*.md' | wc -l
632
```

`CHANGELOG.md` has **3**. The bulk is `docs/specs/` — **632 of 745** — the historical spec
documents, task lists, and proofs from specs 01 through 24.

This matters beyond pedantry. Those documents should not be rewritten: a spec describing what
was built in March is a record of that moment, and editing it to say `omgitworks` would make it
assert something untrue. Same for the changelog.

So the actual rename surface is **110 live references**, about one seventh of the raw number.
The rename is substantially more tractable than the headline count suggested.

## Proof 3: risk classification and user impact

Every row is classified irreversible / breaking / cosmetic, and flagged for whether it affects
an already-installed user. **Six touchpoints do**, and they are exactly what the compatibility
guarantee has to cover:

1. Binary name (`Makefile:4`, `.goreleaser.yml:16`)
2. Completion function names (`shellinit.go:92,163`)
3. The `shell-init` template itself
4. goreleaser `project_name` and archive names
5. Brew formula name
6. Brew tap repository

## Proof 4: line references verified at HEAD

`go.mod:1`, `Makefile:4`, `.goreleaser.yml:6,16,55`, `shellinit.go:92,163`, `mkdocs.yml:2` all
confirmed against the current tree.

## Proof 5: availability matrix carries its check date

Recorded in the spec with both the 2026-09-04 and 2026-09-14 checks, since availability decays
and a stale matrix is worse than none.
