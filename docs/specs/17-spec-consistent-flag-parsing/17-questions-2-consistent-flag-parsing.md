# 17 Questions Round 2 - Consistent Flag Parsing

Follow-up questions based on your answers in Round 1, particularly the lowercase=filter / uppercase=show+filter pattern.

## 1. Verbose Flag Conflict

Currently `-v` is a count flag for verbose levels (`-v`, `-vv`, `-vvv` progressively reveal more columns). You suggested `-V` (visibility) become `-vv`. This conflicts because `-vv` currently means "verbose level 2." How should we handle this?

- [ ] (A) Remove the `-v` count flag entirely and rely on the new uppercase pattern for showing columns (e.g., `-T` shows tags, `-Y` shows type, etc.). Verbose levels are no longer needed.
- [ ] (B) Keep `-v` as verbose BUT change its short flag to something else (e.g., `--verbose` only, no short flag), freeing `-vv` for visibility.
- [ ] (C) Keep `-v` verbose as-is and find a different short flag for visibility (not `-vv`). What letter?
- [X] (D) Other (describe) For now let's change the flag associated to visibility column to be -i. -I would then become how it is displayed.

## 2. Remote Raw Flag

You suggested `-R` (remote-raw) become `-rr`. Since `-r` is the remote flag, would this pattern mean:

- `-r github` = filter by remote (no column shown)
- `-R github` = filter by remote AND show remote column
- `-rr` or `--remote-raw` = show raw remote URL column (no filter)
- `-rr github` = filter by raw remote AND show raw remote column

Is that the intent, or did you mean something different by `-rr`?

- [X] (A) Yes, that's roughly right. `-rr` is shorthand for "remote raw" as a double-letter convention.
- [ ] (B) Actually, `-rr` should just be the short flag equivalent of `--remote-raw`, following the same lower/upper pattern: `-rr` filters, and some uppercase variant shows the column.
- [ ] (C) Drop the remote-raw short flag entirely; just use `--remote-raw` / `--remote-raw=value`. The remote flag `-r`/`-R` is sufficient for most use cases.
- [ ] (D) Other (describe)

## 3. Complete Flag Mapping

Based on the new pattern (lowercase=filter, UPPERCASE=filter+show), here's what the full mapping would look like. Please confirm or correct:

| Long Flag | Filter (lowercase) | Show+Filter (uppercase) | Notes |
|-----------|-------------------|------------------------|-------|
| `--name` | `-n value` | `-N value` | Name has no column today; should `-N` add one? |
| `--tag` | `-t value` | `-T value` | |
| `--type` | `-y value` | `-Y value` | |
| `--status` | `-s value` | `-S value` | |
| `--path` | `-p value` | `-P value` | |
| `--show-user` | `-u value` | `-U value` | |
| `--remote` | `-r value` | `-R value` | |
| `--remote-raw` | `-rr` | `-RR` | Needs resolution per Q2 |
| `--visibility` | `-I` | `-I` | Needs resolution per Q1 |

- [ ] (A) This mapping looks correct. Answer Q1 and Q2 to fill in the gaps.
- [ ] (B) Some flags need different letters (describe which ones).
- [X] (C) Other (describe) A, but with the added not that -N is not necessary as current expected functionality is that name is always displayed (it is the minimum displayed value)

## 4. Showing Columns Without Filtering

Currently `gws list -t` (bare) shows the tags column without filtering. In the new pattern, how do you show a column WITHOUT filtering?

- [X] (A) Use the uppercase flag with no value: `-T` alone shows the tags column. `-T ai` shows it AND filters.
- [ ] (B) The uppercase flag always requires a value. Use `-v` verbose levels or a separate `--show-columns` flag to add columns without filtering.
- [ ] (C) Other (describe)
