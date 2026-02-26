# 01 Questions Round 2 - Init & Add Refactor

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

---

## 1. `--add-tag` Short Flag — Conflict With `-g`

In Round 1 you noted `-g` as the preferred short form for `--add-tag`. However, `-g` is already used by `--go` (navigate to repository). This would create a conflict.

What would you like to do?

- [ ] (A) **Reassign `-g` from `--go`** — Give `--go` a different short form (e.g. `-n` for navigate) and use `-g` for `--add-tag`. Please note your preferred short form for `--go` below.
- [X] (B) **Use a different letter for `--add-tag`** — Keep `-g` for `--go` and pick a different short form for `--add-tag`. Please note your preferred letter below.
- [ ] (C) **`--add-tag` loses its short form** — Keep `-g` for `--go`. Users type `--add-tag` in full. (Simplest option.)
- [ ] (D) Other (describe)

> Notes (preferred letters or other details):
use -d
---

## 2. Tab Completion — Scope & Interaction Model

You selected "Both" in Round 1 (shell completion scripts + interactive fuzzy search). A clarifying question on how these two fit together:

In Round 1 Q4, you said `gws --add` with no argument defaults to the current directory (no interactive prompt). So when would the interactive fuzzy search be triggered?

- [X] (A) **Shell completion only for this spec** — The `<Tab>` key in the shell handles path completion for `gws --add <Tab>`. The in-app fuzzy search is a separate, future feature. Keeps scope tight.
- [ ] (B) **Shell completion + fuzzy search triggered by a flag** — Shell completion for normal use; e.g. `gws --add --browse` opens an interactive directory picker. The fuzzy search is only shown when explicitly requested.
- [ ] (C) **Shell completion + fuzzy search as fallback** — If the user types `gws --add` with no argument AND the current dir is NOT a git repo (instead of erroring), show the interactive directory picker.
- [ ] (D) Other (describe)

> Notes:
