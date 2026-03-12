# 19 Questions Round 1 - Status Compact Display

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Compact Status Icons Format

Currently the full status (`-S`) shows icons like `↓2 ↑1 ✗`. For the compact `-s` display right-justified in the NAME column, should the icons use the same format?

- [X] (A) Same icons and spacing as today: `↓2 ↑1 ✗` (e.g., `my-repo          ↓2 ↑1 ✗`)
- [ ] (B) More compact format without spaces: `↓2↑1✗` (e.g., `my-repo             ↓2↑1✗`)
- [ ] (C) Different/simplified symbols (describe below)
- [ ] (D) Other (describe)

## 2. Clean Repos - Show Indicator?

When a repo is clean with no ahead/behind, should the compact `-s` display show anything?

- [X] (A) Show `✓` right-justified (e.g., `my-repo                    ✓`)
- [ ] (B) Show nothing — only show icons when there's something noteworthy (dirty, ahead, behind)
- [ ] (C) Other (describe)

## 3. Color Support

Should the compact status icons in the NAME column support the same ANSI colors as the current full status display?

- [X] (A) Yes, same colors (magenta for behind, cyan for ahead, red for dirty, green for clean)
- [ ] (B) No color — plain text only
- [ ] (C) Other (describe)

## 4. NAME Column Width

When `-s` is used, the NAME column needs extra width for the right-justified icons. How should the width be determined?

- [X] (A) Expand the NAME column to accommodate the longest name + padding + longest icon string
- [ ] (B) Use a fixed minimum total width (e.g., 40 chars) and truncate repo names if needed
- [ ] (C) Other (describe)

## 5. Filtering Behavior with `-s`

Currently `-s` is defined as a filter-by-status flag (e.g., `gws list -s dirty`). Should this compact display be:

- [X] (A) Always on when `-s` is used (bare `-s` shows compact status; `-s dirty` shows compact status AND filters)
- [ ] (B) Only compact display with no filtering — `-s` just means "show compact status" (filtering only via `-S`)
- [ ] (C) Other (describe)

## 6. Combining `-s` and `-S`

If a user specifies both `-s` and `-S`, what should happen?

- [X] (A) `-S` wins — show the full status column (ignore `-s` compact display)
- [ ] (B) Show both — compact icons in NAME column AND full status column
- [ ] (C) Error — tell the user they're mutually exclusive
- [ ] (D) Other (describe)

## 7. Branch Name in Compact Display

The current full status (`-S`) shows the branch name alongside icons. Should the compact `-s` display include the branch name?

- [X] (A) No branch name — only show ahead/behind/clean/dirty icons (as described in your request)
- [ ] (B) Include abbreviated branch name before icons
- [ ] (C) Other (describe)

## 8. JSON Output

When using `--json` output with `-s`, should the compact status data be included in the JSON?

- [ ] (A) Yes, add a compact status field to JSON output
- [X] (B) No, JSON output only gets status data with `-S` (full status)
- [ ] (C) Other (describe)
