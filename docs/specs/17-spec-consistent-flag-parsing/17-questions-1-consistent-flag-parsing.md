# 17 Questions Round 1 - Consistent Flag Parsing

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Consistent Direction

The inconsistency is that dual-purpose flags (`-t`, `-y`, `-s`, `-u`, `-p`, `-r`, `-R`, `-V`) require `=` for values (e.g., `-t=ai`) while the filter-only flag (`-n`) accepts a space (e.g., `-n dks`). You've expressed a preference for NOT requiring `=`. To confirm:

- [X] (A) Make all flags work WITHOUT requiring `=` (e.g., `-t ai` works just like `-n dks`). The `=` syntax should still work too (backward compatible).
- [ ] (B) Make all flags REQUIRE `=` for values (e.g., `-n=dks` required, `-n dks` no longer works). This is more strict but fully consistent.
- [ ] (C) Other (describe)

## 2. Bare Flag Behavior (Show Column)

Currently, dual-purpose flags can be used bare to show a column without filtering (e.g., `gws list -t` shows the tags column). If we remove the `=` requirement, we need a way to distinguish "show column" from "filter by value". How should bare flags work?

- [ ] (A) Keep bare flag behavior: `-t` alone still means "show tags column". The parser should look ahead to determine if the next argument is a value or a different flag/positional arg.
- [ ] (B) Bare flags always show the column AND if a value follows with a space, it also filters. Essentially `-t` = show column, `-t ai` = show column + filter by "ai" (same as current `-t=ai` behavior).
- [ ] (C) Drop the "show column only" bare flag feature entirely; require a value or use `-v` verbose levels to control column visibility.
- [X] (D) Other (describe) this question is funny to me becuase it touches on another change I've been thinking about. I Would like to be able to filter with out asserting that the filter must be displayed in the output. It is somewhat related to this uestion in that currently both -t and -t= result in the tag column desplaying. It's just that one filters and the other doesn't. is it possible to have it setup the pattern that lower case flag is filter only and upper case flags are used to filter and show (i.e. -t ai would filter on ai but not display the tag column -T ai would filter and show the column) I know this would require a rework of -R and -V which brings me to another question. can we make -R be -rr, -V be -vv and then. apply the pattern to them as well?

## 3. Flag Stacking with Values

Currently you can stack bare dual-purpose flags: `gws list -yVt` shows type, visibility, and tags columns. With the new space-based value syntax, how should stacking work when a value is provided?

- [X] (A) Only the LAST flag in a stack can take a value: `-yt ai` means show type column + filter tags by "ai". This follows POSIX conventions.
- [ ] (B) Flag stacking should only work for bare (show column) usage. If you want to filter, use the flag separately: `-y -t ai`.
- [ ] (C) Keep current stacking behavior unchanged; stacked flags are always bare/show-column.
- [ ] (D) Other (describe)

## 4. Backward Compatibility

The `=` syntax currently works for dual-purpose flags. Should it continue to work after this change?

- [X] (A) Yes, both syntaxes should work: `-t ai` and `-t=ai` are equivalent (full backward compatibility).
- [ ] (B) Deprecate `=` syntax with a warning, remove in a future release.
- [ ] (C) Other (describe)

## 5. Scope of Commands Affected

Should this consistency fix apply only to the `list` command, or should it extend to other commands that use flags (e.g., `tag --repo`, `tag --path`)?

- [ ] (A) Only the `list` command - that's where the inconsistency is most visible.
- [ ] (B) All commands with filter flags should be made consistent.
- [X] (C) Other (describe) B, but only if the scope is reasonable.
