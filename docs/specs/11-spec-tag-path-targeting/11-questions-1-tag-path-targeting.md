# 11 Questions Round 1 - Tag Path Targeting

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Default Behavior (No Flag)

Currently `gws tag add <identifier> <tag>` matches repos by partial name OR exact path. With this change, what should the default behavior be when no flag is given?

- [ ] (A) Match by partial name only — path matching moves exclusively to `--path` (clean separation)
- [X] (B) Keep current behavior — match by partial name OR exact path (backward compatible, no change to default)
- [ ] (C) Other (describe)

## 2. Explicit `--repo` / `-r` Flag

You mentioned adding a `--repo`/`-r` flag. What should it do relative to the current default name-matching behavior?

- [A] (A) Same as the default (partial, case-insensitive name match) — just makes targeting explicit; useful for scripts/clarity
- [ ] (B) Exact name match only — stricter than the default partial match
- [ ] (C) It's not needed — the default (no flag) already handles name matching; skip the `--repo` flag and only add `--path`
- [ ] (D) Other (describe)

## 3. Path Matching Strategy

For `--path`/`-p`, what kind of matching should it perform?

- [ ] (A) Substring match — any repo whose path contains the given string (e.g. `/daileyo` matches `/home/daileyo/gws` and `/home/daileyo/api`)
- [ ] (B) Prefix match — repo path must start with the given string (e.g. `/home/daileyo` matches repos under that directory)
- [X] (C) Both A and B — try prefix first, fall back to substring
- [ ] (D) Other (describe)

## 4. Path Match Case Sensitivity

Should `--path` matching be case-sensitive or case-insensitive?

- [ ] (A) Case-insensitive — consistent with how name matching works today
- [X] (B) Case-sensitive — paths are usually case-sensitive on Linux; match literally
- [ ] (C) Other (describe)

## 5. Combining `--repo` and `--path`

What should happen if someone passes both `--repo` and `--path` in the same command?

- [ ] (A) Error — they are mutually exclusive; print a clear message
- [ ] (B) Union — match repos that satisfy either condition (OR logic)
- [X] (C) Other (describe) If both are supplied, it should apply AND logic (both need to be true)

## 6. Placement of Flags

Where should `--path`/`-p` (and `--repo`/`-r` if included) be available?

- [ ] (A) On `tag add` and `tag remove` sub-subcommands only (e.g. `gws tag add --path /foo bar`)
- [ ] (B) Also on the `tag` command's short flags (`gws tag -a --path /foo bar`)
- [A] (C) Both A and B — available everywhere the targeting identifier can be specified
- [ ] (D) Other (describe)

## 7. Output / Feedback

When `--path` matching finds (or fails to find) repos, what feedback style is preferred?

- [X] (A) Same as current — show count and list up to 5 repo names (e.g. `Added tag 'work' to 3 repositories`)
- [X] (B) Also show the matched path next to the repo name (e.g. `  - my-repo  /home/daileyo/my-repo`)
- [ ] (C) Other (describe)

## 8. Shell Completion for `--path`

For shell tab completion when using `--path`, what should be offered?

- [X] (A) Complete from known repo paths in the gws config (suggests actual tracked repo paths)
- [ ] (B) No completion for `--path` — filesystem paths are too varied to complete usefully
- [ ] (C) Complete directory paths from the filesystem (standard file completion)
- [ ] (D) Other (describe)
