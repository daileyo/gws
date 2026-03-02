# 09 Questions Round 2 - Command Structure Refactor

Based on your Round 1 answers, I need to clarify a few things before writing the spec.

## 1. Short Flag Stacking and Routing

You want `gws -luv` where `-l` is the base command (list) and `-u`/`-v` are modifiers. In Cobra, single-letter flags on the root command can be stacked POSIX-style. But this means ALL stackable flags must be registered on the **same command** (root). This creates a tension:

If filters are scoped to `list` as a subcommand (your Q7 answer), then `gws list -uv` works (flags on the `list` subcommand), but `gws -luv` would NOT work because `-l` is on root and `-u`/`-v` are on the `list` subcommand.

Which trade-off do you prefer?

- [ ] (A) Register common modifier flags (`-v`, `-u`, `-s`, etc.) on BOTH root and relevant subcommands. Root short flags route to the right subcommand. `gws -lv` and `gws list -v` both work.
- [ ] (B) Keep all flags on root (current pattern) but use subcommand words as aliases for the primary flag. `gws list` is sugar for `gws -l`, and all flags live on root. Stacking always works: `gws -luv`.
- [X] (C) Accept that stacking only works within a subcommand: `gws list -uv` works, but `gws -luv` does not. Use subcommand words as the primary interface.
- [ ] (D) Other (describe)

## 2. The `-t` / Tag Conflict

You want `tag` as a base command AND `-t` is currently the tag filter on `list`. These conflict since `-t` can't mean both "invoke the tag command" and "filter by tag in list."

How should this be resolved?

- [ ] (A) `-t` stays as the `list` filter for tags. The `tag` command gets a different shorthand (e.g., no single-letter shorthand, or something like `-T`)
- [ ] (B) `-t` becomes the `tag` base command shorthand. The list tag filter uses a different short flag or only the long form `--tag`
- [ ] (C) `tag` has no single-letter root alias. It's always `gws tag ...`. `-t` remains the list filter.
- [X] (D) Other (describe) I would like to better understand why this wouldn't work to have -t as the shorthand for the tag command; while still being able to use -t within the list for tag filtering (gws list -tu) for example) Now that we've committed to option C for Qestion 1; does that not make this possible?

## 3. Tag Sub-Flag Pattern

You said tag should be a base command with add/remove as sub-operations. Should the sub-operations use the same reusable flags (`-a`/`-d`) or the word forms (`add`/`remove`)?

- [ ] (A) Flag form only: `gws tag -a <repo> <tag>`, `gws tag -d <repo> <tag>`
- [ ] (B) Word form only: `gws tag add <repo> <tag>`, `gws tag remove <repo> <tag>`
- [X] (C) Both: `gws tag add` and `gws tag -a` are equivalent
- [ ] (D) Other (describe): bare `gws tag <repo> <tag>` adds (default action), `-d` removes

## 4. User Sub-Flag Mapping

Refactoring `user` to reusable sub-flags. The current `user` subcommand has: `add`, `list`, `show`, `remove`, `assign`, `sync`. How do these map?

- [X] (A) Direct mapping: `-a` (add), `-d` (remove), `-l` (list), `-s` (show). `assign` and `sync` stay as word subcommands since they don't fit the pattern.
- [ ] (B) All become flags: `-a` (add), `-d` (remove), `-l` (list), `-s` (show), `--assign` (assign), `--sync` (sync)
- [ ] (C) Only common CRUD gets flags, everything else is a word subcommand: `gws user -a`, `gws user -d`, `gws user -l`, but `gws user assign`, `gws user sync`
- [ ] (D) Other (describe)

## 5. Default Behavior When No Sub-Flag

What should happen when a base command is invoked without a sub-flag?

- [X] (A) Show help for that command: `gws tag` shows tag help, `gws user` shows user help
- [ ] (B) Default to list: `gws tag` lists tags, `gws user` lists users
- [ ] (C) Default varies by command: `gws tag <repo> <tag>` adds a tag (most common action), `gws user` lists users
- [ ] (D) Other (describe)

## 6. Scope of This Spec

This refactor touches the entire command layer. Should we tackle it all in one spec, or break it up?

- [ ] (A) One spec covering the full refactor (all commands, tab completion, deprecation, shell integration)
- [ ] (B) Split into 2 specs: (1) Core command structure + list/init/add/refresh, (2) Tag + User refactor
- [X] (C) Split into 3 specs: (1) Core structure + deprecation layer, (2) Tag refactor, (3) User refactor
- [ ] (D) Other (describe)
