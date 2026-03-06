# 15 Questions Round 1 - Status Alignment

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Status Column Layout

Currently the status column shows: `main ✓` or `feature ✗ ↑3↓2` (all left-justified as one string). How should the new layout work?

- [X] (A) Split into two sub-columns: branch name left-justified on the left, then icons right-justified on the right, within the same STATUS column
- [ ] (B) Keep as one STATUS column but pad so icons always align to the right edge of the column (branch name left, icons pushed right with spaces)
- [ ] (C) Split STATUS into separate columns: BRANCH (left-justified) and STATUS (right-justified, icons only)
- [ ] (D) Other (describe)

## 2. Color for Clean/Dirty Icons

What colors should the clean (`✓`) and dirty (`✗`) icons use?

- [X] (A) Green for `✓` (clean), Red for `✗` (dirty)
- [ ] (B) Green for `✓` (clean), Yellow for `✗` (dirty)
- [ ] (C) Green for `✓` (clean), Red for `✗` (dirty), with bold
- [ ] (D) Other (describe specific colors)

## 3. Color for Ahead/Behind Arrows

What colors should the ahead (`↑N`) and behind (`↓N`) indicators use?

- [ ] (A) Yellow for both ahead and behind
- [X] (B) Cyan/Blue for ahead (`↑`), Magenta/Red for behind (`↓`)
- [ ] (C) Yellow for ahead (`↑`), Red for behind (`↓`)
- [ ] (D) No color on arrows, only color on `✓`/`✗`
- [ ] (E) Other (describe)

## 4. Color Library Preference

The project currently has no color library. Which approach do you prefer?

- [ ] (A) Use `fatih/color` — popular, simple API, auto-detects terminal support
- [ ] (B) Use `charmbracelet/lipgloss` — modern, part of the Charm ecosystem
- [X] (C) Use raw ANSI escape codes — no dependency, minimal approach
- [ ] (D) Other (describe)

## 5. Non-TTY / Piped Output Behavior

When output is piped (not a TTY), should colors be automatically disabled?

- [ ] (A) Yes, auto-detect TTY and disable colors when piped (standard practice)
- [X] (B) Yes, auto-detect, but also add a `--color` flag to force on/off
- [ ] (C) Other (describe)

## 6. Branch Name Truncation

If a branch name is very long (e.g., `feature/JIRA-1234-implement-user-authentication-flow`), should it be truncated to keep the column manageable?

- [ ] (A) No truncation — show full branch name, let the column grow
- [X] (B) Truncate with ellipsis after N characters (e.g., `feature/JIRA-1234-impl...`)
- [ ] (C) Only truncate if terminal width would be exceeded
- [ ] (D) Other (describe)
