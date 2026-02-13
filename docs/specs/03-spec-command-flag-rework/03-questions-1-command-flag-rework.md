# 03 Questions Round 1 - Command Flag Rework

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Command-to-Flag Mapping

Based on your description (`gws --list all` or `gws -l all`), here's my understanding of how existing subcommands would map to flags. Please confirm or adjust:

| Current (subcommand) | Proposed (flag) | Shorthand | Notes |
|---|---|---|---|
| `gws list` | `gws --list` | `gws -l` | Arguments after flag become subcommands |
| `gws init <dir>` | `gws --init <dir>` | `gws -i <dir>` | |
| `gws tag <repo> <tag>` | `gws --tag <repo> <tag>` | `gws -t <repo> <tag>` | Conflicts with existing `--tag` filter on list. the tag subcommand for list should just become tags rather than --tag. i.e gws -l --tags should become gws -l tags|
| `gws untag <repo> <tag>` | `gws --untag <repo> <tag>` | gws -u <repo> <tag> | |
| `gws refresh` | `gws --refresh` | `gws -r` | |
| `gws version` | `gws --version` | `gws -v` | Cobra provides `--version` natively |

- [x] (A) The mapping above looks correct, adjust any conflicts as needed
- [ ] (B) Some commands should remain as subcommands (specify which)
- [ ] (C) Different shorthands preferred (specify below)
- [ ] (D) Other (describe)

Notes:


## 2. Tag Command Conflict

Currently `gws list --tag personal` uses `--tag` as a filter flag. If we convert the `tag` subcommand to `--tag`, there's a naming collision. How should we resolve this?

- [X] (A) Rename the tag command flag to something else (e.g., `--add-tag` / `--remove-tag` for tag/untag)
- [ ] (B) Rename the list filter to something else (e.g., `--filter-tag` or `--with-tag`)
- [ ] (C) Keep `tag` and `untag` as subcommands (don't convert them to flags)
- [ ] (D) Other (describe)

Notes:


## 3. List Command Filter Flags

Currently, `list` has its own flags: `--type`, `--tag`, `--name`, `--path`, `--output`/`-o`, `--status`/`-s`. In the new pattern, how should these work?

Example: `gws --list --type github --status` or `gws -l --type github -s`

- [ ] (A) All existing list filter flags move to the root command and only apply when `--list` is used
- [ ] (B) Filter flags become positional subcommands (e.g., `gws --list type github`)
- [X] (C) Keep current flag pattern for filters, only the top-level command changes (preferred Cobra pattern)
- [ ] (D) Other (describe)

Notes:


## 4. Tag Filter Shorthand

Currently `gws personal` is shorthand for `gws list --tag personal`. How should this work in the new pattern?

- [ ] (A) Keep it as-is: `gws personal` still works as a tag filter shorthand
- [X] (B) Remove the shorthand; require explicit `gws --list --tag personal` or `gws -l --tag personal`
- [ ] (C) Change shorthand to work with new pattern (e.g., `gws -l personal`)
- [ ] (D) Other (describe)

Notes:


## 5. Help and Version Behavior

Cobra automatically provides `--help`/`-h` and can provide `--version`/`-v`. Should we:

- [X] (A) Use Cobra's built-in `--version`/`-v` flag (replaces `gws version` subcommand entirely)
- [ ] (B) Keep a custom `--version` flag that shows the same detailed output (version, commit, build date)
- [ ] (C) Keep `version` as a subcommand (don't convert it to a flag)
- [ ] (D) Other (describe)

Notes:


## 6. Backward Compatibility

How should we handle the transition from subcommands to flags?

- [X] (A) Clean break - remove old subcommand pattern entirely, no backward compatibility
- [ ] (B) Deprecation period - support both patterns with deprecation warnings for old style
- [ ] (C) Aliases - keep subcommands as hidden aliases that redirect to the new flag pattern
- [ ] (D) Other (describe)

Notes:


## 7. Proof Artifacts

What would you like to see as proof that this rework is complete?

- [ ] (A) CLI output screenshots showing new flag patterns working (e.g., `gws -h`, `gws --list`, `gws -l --type github`)
- [X] (B) All existing tests passing with updated assertions
- [ ] (C) Both A and B
- [ ] (D) Other (describe)

Notes:

