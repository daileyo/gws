# 12 Questions Round 1 - List Remote

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Remote URL Source

The config already stores a `RemoteURL` (captured at discovery time from origin). Should the `--remote` flag display the stored URL from config, or should it do a live git inspection of each repo's remotes?

- [ ] (A) Use stored `RemoteURL` from config (faster, but may be stale if remotes changed since last `refresh`)
- [ ] (B) Live inspect each repo's git remotes (slower, but always current — similar to how `--status` works)
- [ ] (C) Use stored URL but fall back to live inspection if missing
- [X] (D) Other (describe) If you mean the remote that is part of the git info for the repo; then use the one from config. If you mean gws refresh stores the remote url somehow; please explain more.

## 2. Asterisk Indicator for Multiple Remotes

You mentioned an asterisk when there are remotes beyond origin. Should the asterisk appear:

- [ ] (A) Appended to the URL (e.g., `https://github.com/user/repo.git *`)
- [X] (B) As a prefix (e.g., `* https://github.com/user/repo.git`)
- [ ] (C) In a separate column (e.g., `REMOTE` column + `MULTI` column)
- [ ] (D) Other (describe)

## 3. No-Origin Display

When a repo has no origin but has other remotes, you want just an asterisk displayed. Should there be any additional info in this case?

- [X] (A) Just `*` — the user can investigate manually
- [ ] (B) `* (no origin)` — with a label explaining why
- [ ] (C) Show the first non-origin remote URL with the asterisk (e.g., `* https://gitlab.com/user/repo.git`)
- [ ] (D) Other (describe)

## 4. JSON Output Support

The list command supports `--output json`. Should the remote info be included in JSON output when `--remote` is used?

- [X] (A) Yes, add `remote_url` and `has_multiple_remotes` fields to JSON output
- [ ] (B) Yes, add a full `remotes` object with all remote names and URLs
- [ ] (C) No, remote flag only affects table output
- [ ] (D) Other (describe)

## 5. Interaction with Existing Flags

Should `--remote` be stackable with `--status` and `--show-user` (e.g., `gws list -rsu`)?

- [X] (A) Yes, fully stackable — remote column appears alongside status and user columns
- [ ] (B) Yes, but remote replaces the PATH column to keep output width manageable
- [ ] (C) Other (describe)

## 6. Column Placement

Where should the REMOTE column appear in the table output?

- [X] (A) After PATH (last column)
- [ ] (B) Before PATH (second to last)
- [ ] (C) After TYPE/VISIBILITY columns, before TAGS
- [ ] (D) Other (describe)
