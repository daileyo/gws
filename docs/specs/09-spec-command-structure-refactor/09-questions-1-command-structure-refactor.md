# 09 Questions Round 1 - Command Structure Refactor

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Base Command Syntax

Currently all operations use `--` flags (`gws --list`, `gws --init`). You want to move back to subcommand-style. What syntax should base commands use?

- [ ] (A) Bare subcommands: `gws list`, `gws init`, `gws tag`, `gws add`, `gws refresh`, `gws user`
- [ ] (B) Single-dash shorthand only: `gws -l`, `gws -i`, `gws -t` (keep current short flags as the primary interface)
- [ ] (C) Both: `gws list` and `gws -l` are equivalent (subcommand with a single-letter alias)
- [X] (D) Other (describe) C. however, I want to ensure this doesn't break tab completions for repo names (i.e. if a user types gws user + tab, it should be able to try to match both the command and any repos that might also start with user... like user-experience.. as an example repo name)

## 2. Tab Completion Differentiation

The original reason for `--` prefixed commands was so tab completion could distinguish commands from repo names. How should this be handled now?

- [ ] (A) Use Cobra's subcommand completion natively — repo names complete after navigation-specific commands (e.g., `gws go <tab>`) but not after base commands
- [X] (B) Keep positional `gws <repo>` for navigation and rely on Cobra to differentiate subcommands from repo names automatically (Cobra handles this since subcommands take priority)
- [ ] (C) Add a dedicated `go` or `cd` subcommand for navigation (`gws go <repo>`) and remove bare positional navigation
- [ ] (D) Other (describe)

## 3. Reusable Sub-Flags

You mentioned `-a` (add) and `-d` (delete) should be reusable across commands like `user` and `tag`. What is the full set of reusable sub-flags you envision?

- [ ] (A) Minimal set: `-a` (add/create), `-d` (delete/remove), `-l` (list)
- [ ] (B) Extended set: `-a` (add/create), `-d` (delete/remove), `-l` (list), `-u` (update), `-s` (show/details)
- [ ] (C) Full CRUD: `-a` (add/create), `-d` (delete/remove), `-l` (list), `-u` (update), `-s` (show), plus domain-specific ones
- [X] (D) Other (describe) It can and will vary. what I dont' want to have to do is keep coming up with different names and non-intuitive short-hand substitutions for things like crud that may be part of multiple different base commands. (verbose, quite, etc are other ones that may end up getting reused)

## 4. Tag Command Consolidation

Currently tagging is split: `--add-tag` / `-d` (add) and `--remove-tag` / `-x` (remove). With the new structure, how should tag management work?

- [ ] (A) `gws tag <repo> <tag-name>` (add by default), `gws tag -d <repo> <tag-name>` (delete)
- [ ] (B) `gws tag -a <repo> <tag-name>` (add), `gws tag -d <repo> <tag-name>` (delete), `gws tag` (list all tags)
- [ ] (C) `gws tag <repo> <tag-name>` (add), `gws tag -d <repo> <tag-name>` (delete), `gws tag -l` (list tags)
- [X] (D) Other (describe) I always wanted tag to be the base command, with add, and remove as subcommands. they just somehow got broken out during the rework. This is also complicated by the fact that there is a tag option as a filter for the list base command.

## 5. User Command Structure

The `user` command already has a proper subcommand tree (`user add`, `user list`, etc.) plus overlapping root-level flags (`--user`, `--update`, `--delete`). How should this be handled?

- [ ] (A) Keep the current `user` subcommand tree as-is (it's already well-structured) and remove the root-level `--user`/`--update`/`--delete` flags
- [X] (B) Refactor `user` to match the new reusable sub-flag pattern: `gws user -a` (add), `gws user -d` (delete), `gws user -l` (list), etc.
- [ ] (C) Support both: keep the subcommand tree (`user add`, `user list`) AND add shorthand flags (`user -a`, `user -l`) as aliases
- [ ] (D) Other (describe)

## 6. Backward Compatibility

The `--` flag style was introduced in v2.0. How should backward compatibility be handled?

- [ ] (A) Clean break: remove all `--` flag forms, this is a new major version (v3.0)
- [X] (B) Deprecation period: keep `--` flags working with deprecation warnings for one release cycle, then remove
- [ ] (C) Maintain both: keep `--` flags as hidden aliases alongside new subcommands permanently
- [ ] (D) Other (describe)

## 7. Filter Flags with List

Currently filter flags (`--type`, `--name`, `--tag`, `--status`, etc.) only work with `--list`. How should filters work with subcommand-style `list`?

- [X] (A) Filters become flags on the `list` subcommand: `gws list -t work -y github -s` (same short flags, scoped to `list`)
- [ ] (B) Filters stay as-is but work with `gws list` instead of `gws --list`: `gws list --type github --tag work`
- [ ] (C) Mix: long-form filters on `list` subcommand, reusable short flags: `gws list -t work --type github`
- [ ] (D) Other (describe)

## 8. Stacking Example Clarification

You mentioned sub-flags should be "stackable." Can you clarify what stacking means in your vision?

- [ ] (A) Stacking means combining sub-flags: `gws user -a --name "John" --email "j@x.com"` (the `-a` triggers "add mode" and subsequent flags provide values)
- [ ] (B) Stacking means combining base commands: `gws -l -s` (list + show status in one call)
- [ ] (C) Both: sub-flags combine with their base command, AND some base command short flags can combine
- [X] (D) Other (describe) I would like to be able to do something like gws -luv where u could be show-users v could be a verbose mode and l was the base command alternatively, this could also be something like gws list -uv.
