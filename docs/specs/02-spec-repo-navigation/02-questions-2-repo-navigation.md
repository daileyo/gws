# 02 Questions Round 2 - Repository Navigation

Follow-up questions based on your Round 1 answers. Please mark your selections with `[X]`.

---

## 1. Tab Selection for Multiple Matches (follow-up to Q2)

You chose "list all matches and allow tabbing to select." Can you clarify the interactive model?

- [ ] (A) Arrow-key menu — display a navigable list in the terminal (like `fzf` or `gum choose`), user presses up/down arrows and Enter to select
- [ ] (B) Numbered list — display numbered options (1, 2, 3...), user types the number and presses Enter
- [X] (C) Tab-completion style — user presses Tab to cycle through matches inline
- [ ] (D) Other (describe)

## 2. Interactive Selection Dependency

For the interactive selection UI, would you like to:

- [ ] (A) Use a third-party TUI library (e.g., `bubbletea`, `survey`, `promptui`) for a polished experience
- [X] (B) Keep it simple with standard I/O — no external TUI dependencies, just numbered list with input prompt
- [ ] (C) Other (describe)

## 3. Verbose Output Details (follow-up to Q5)

You chose verbose by default with `--quiet` for scripting. What should verbose output include?

- [X] (A) Path + repo name and type (e.g., `my-repo (github) → /home/user/projects/my-repo`)
- [ ] (B) Path + repo name, type, and tags (e.g., `my-repo (github) [personal, go] → /home/user/projects/my-repo`)
- [ ] (C) Path + full repo details including git status (name, type, tags, branch, clean/dirty)
- [ ] (D) Other (describe)

## 4. Flag Name for Navigation

You chose both positional args and a flag. What should the flag be named?

- [X] (A) `--go` / `-g` — intuitive action verb
- [ ] (B) `--navigate` / `-g` — more descriptive
- [ ] (C) `--cd` / `-g` — mirrors the shell action
- [ ] (D) Other (describe)

## 5. Quiet Flag Scope

Should `--quiet` be a navigation-only flag, or a global flag usable with other commands too?

- [X] (A) Navigation-only — `--quiet` / `-q` only applies when navigating
- [ ] (B) Global flag — `--quiet` / `-q` suppresses non-essential output across all commands (future-proof)
- [ ] (C) Other (describe)
