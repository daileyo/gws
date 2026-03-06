# 13 Questions Round 1 - List Refactor

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Default No-Args Behavior

When `gws` is run with no arguments, it currently shows a workspace summary and usage hint. It also handles navigation when a positional arg is passed (`gws my-repo`). If we default to `gws list`, how should navigation work?

- [X] (A) `gws` (no args) shows `gws list` output; `gws <repo-name>` still navigates (detect positional args vs no args)
- [ ] (B) `gws` (no args) shows `gws list` output; navigation requires `gws go <repo-name>` or similar subcommand
- [ ] (C) `gws` (no args) shows `gws list` output; keep navigation as-is with positional args (same as A but clarifying no change to navigation)
- [ ] (D) Other (describe)

## 2. Multi-Column Default View

When `gws list` shows just repo names in multi-column layout, what style do you prefer?

- [X] (A) Simple columns like `ls` — fill columns left-to-right, top-to-bottom, auto-detect terminal width
- [ ] (B) Sorted alphabetically down each column (like `ls -C` default behavior)
- [ ] (C) Single alphabetically sorted list displayed in columns flowing top-to-bottom then left-to-right
- [ ] (D) Other (describe)

## 3. Flag Letter Assignments

You mentioned `-y` for type, `-V` for visibility, `-t` for tag, `-p` for path. Currently `-t` is already used for `--tag` (filter) and `-p` for `--path` (filter). Should we confirm these mappings? Also, what about status (`-s`), user (`-u`), and remote (`-r`)?

- [X] (A) Keep the mappings as you described: `-y`=type, `-V`=visibility, `-t`=tag, `-p`=path. Status/user/remote keep current `-s`/`-u`/`-r` flags with same dual-purpose behavior (show column, or filter if value given)
- [ ] (B) Same as A, but `-s`/`-u`/`-r` remain display-only toggles (no filter capability for status/user/remote)
- [ ] (C) Remap some flags — please specify preferred mappings
- [ ] (D) Other (describe)

## 4. Dual-Purpose Flag Behavior

For flags that act as both display and filter (e.g., `gws list -y github`), how should the interaction work when combining display-only and filter flags?

Example: `gws list -y github -Vtp`
- Filters: type=github
- Displays: name + type + visibility + tag + path

- [_] (A) Exactly as above — a flag with a value filters AND shows that column; a flag without a value just shows the column. Name column is always shown.
- [ ] (B) Same as A, but filtered columns are always shown even if not explicitly requested (so `-y github` would show type column automatically even without `-y` as display flag)
- [X] (C) Other (describe) Same as A, but if -v|--verbose is passed with any other flag/filter; all columns are displayed.

## 5. Name Filter

Currently `--name`/`-n` filters by name. In the new model, should the name filter work differently since names are always shown?

- [X] (A) Keep `-n` as a filter-only flag (filters by name pattern, no display effect since name is always shown)
- [ ] (B) Positional argument after `list` acts as name filter: `gws list api-*` filters by name
- [ ] (C) Both: `-n` flag and positional arg both work as name filters
- [ ] (D) Other (describe)

## 6. Verbose Mode

`gws list -v` / `gws list --verbose` should show all columns. Which columns should "all" include?

- [ ] (A) All data columns: NAME, TYPE, VISIBILITY, TAGS, PATH (the current default view without status/user/remote)
- [ ] (B) Everything: NAME, STATUS, USER, EMAIL, SIGN, TYPE, VISIBILITY, TAGS, PATH, REMOTE
- [ ] (C) All stored data columns: NAME, TYPE, VISIBILITY, TAGS, PATH, REMOTE — but not live-fetched columns like STATUS and USER (which require git operations)
- [X] (D) Other (describe) A for -v|--verbose. everything for -vv|--verbose-verbose.

## 7. Root Help (`gws -h`)

You said "all flags need to be in `gws -h`". Currently `gws -h` shows subcommands and navigation help. Do you want:

- [ ] (A) Show `list` subcommand flags directly on `gws -h` since list is the default behavior (merge list flags into root help)
- [ ] (B) Show all subcommand flags on `gws -h` (comprehensive but potentially long)
- [ ] (C) Just ensure `gws -h` prominently shows that `gws list -h` has detailed flags, and make the root help cleaner
- [X] (D) Other (describe) I mean gws list -h should show all list flags. there are some flags that are not currently shown when running gws list -h

## 8. JSON Output

Currently `gws list -o json` outputs JSON. How should JSON interact with the new flag system?

- [ ] (A) JSON always outputs all fields regardless of display flags; display flags only affect table output. Filter flags still apply to JSON.
- [X] (B) JSON output respects display flags — only includes fields for selected columns
- [ ] (C) Keep JSON behavior unchanged; new flags only affect table output
- [ ] (D) Other (describe)

## 9. Backward Compatibility with Deprecated Flags

The old `gws --list`, `gws --status` etc. deprecated flags still work. Should they:

- [X] (A) Continue working with deprecation warnings, mapping to new flag behavior
- [ ] (B) Be removed as part of this refactor (breaking change)
- [ ] (C) Continue working but update to match new behavior (no deprecation warning)
- [ ] (D) Other (describe)
