# 22 Questions Round 1 - cd Command

Answers captured from the planning session on 2026-09-03.

## 1. Scope of `gws cd`

How far should `gws cd` go?

- [X] (A) **Workspace root only** — `gws cd` cds to the workspace root, replacing the `cdgws`/`gcd` helper functions the README currently tells users to hand-write. Small and focused.
- [ ] (B) Root + gws-owned dirs — also `gws cd --config` / `gws cd --projects` to reach the XDG directories introduced in spec 23.
- [ ] (C) Root + repo targets — `gws cd <repo>` as an explicit alias of bare `gws <repo>`, plus `gws cd -` for the previous directory.

**Notes:** Keep it minimal. Bare `gws <repo>` already covers repo navigation; `gws cd` is specifically the "take me to the workspace root" verb. `--config` / `--projects` are recorded as open questions for a follow-up once spec 23 lands.
