# 08 Questions Round 2 - User Update/Delete

Follow-up questions based on your Round 1 answers.

## 1. Interactive Multi-Select for Multiple Matches (from Q2D)

You mentioned wanting an interactive list with spacebar selection when multiple repos match. Should this interactive picker also apply to `--user --delete` (not just `--update`)?

- [X] (A) Yes, use the interactive picker for both --update and --delete when multiple repos match
- [ ] (B) Only for --update; --delete should require exact match for safety
- [ ] (C) Other (describe)

## 2. Interactive Picker Library

Go CLI tools typically use a library for interactive selection. Do you have a preference?

- [ ] (A) Use `charmbracelet/bubbletea` (popular, modern TUI framework)
- [ ] (B) Use `AlecAivazis/survey` (simpler, focused on prompts/surveys)
- [ ] (C) Use `manifoldco/promptui` (lightweight, minimal dependencies)
- [ ] (D) Build a simple custom picker (no new dependency — just arrow keys + spacebar with raw terminal input)
- [X] (E) Other (describe) Let's todo this. I want to think about it more. for now, we will assume that when there are multiple matches, that we want to update the user for all the mateches.

## 3. Relationship to `gws user assign` (Clarified)

Now that we've confirmed `gws user assign` exists (just hidden in help), how should `--user --update` relate to it?

- [ ] (A) `--user --update` replaces `user assign` — deprecate/remove the subcommand approach
- [ ] (B) `--user --update` is the primary interface; keep `user assign` but don't advertise it
- [ ] (C) Both coexist — `--user --update` adds tag-based batch + interactive picker, `user assign` keeps --dry-run and --use-subdirs
- [X] (D) Other (describe) Let's start with exposing user assign. I don't know how it got hidden. it should be --user -u.

## 4. Delete Signing Config Flag (from Q4D)

You said "default to leaving the signing key, but have a flag to remove it." What should that flag be called?

- [ ] (A) `--include-signing` (e.g., `gws --user --delete my-repo --include-signing`)
- [X] (B) `--all` (e.g., `gws --user --delete my-repo --all` — removes everything gws-managed)
- [ ] (C) `--clean` (e.g., `gws --user --delete my-repo --clean`)
- [ ] (D) Other (describe)

## 5. Batch Update Output (from Q3 — you selected B)

For the summary shown after batch update, how verbose should it be?

- [ ] (A) Compact: `Updated 5 repositories with profile 'work'` (count only)
- [ ] (B) Moderate: List each repo name + what changed (e.g., `my-project: user.name "Old" → "New"`)
- [ ] (C) Verbose: Full before/after for each repo (like --dry-run output for each)
- [X] (D) Other (describe) Moderate, but with a --verbose option. I'd actually also like to have a --quite option, which just updates and doesn't print anything.
