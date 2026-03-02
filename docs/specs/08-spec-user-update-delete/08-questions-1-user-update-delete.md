# 08 Questions Round 1 - User Update/Delete

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. CLI Interface for --user --update

How should the `--user --update` command identify **which user profile** to apply?

- [ ] (A) Require a named profile from `gws user list` (e.g., `gws --user --update <repo> <profile-name>`)
- [ ] (B) Accept inline values (e.g., `gws --user --update <repo> --name "John" --email "john@co.com"`)
- [X] (C) Both: accept a profile name OR inline values, with inline values taking precedence
- [ ] (D) Other (describe)

## 2. Repository Targeting for Single Update

How should the target repository be specified for `--user --update` (single repo, no tag)?

- [ ] (A) By repo name (partial match, like existing tag commands): `gws --user --update my-project work`
- [ ] (B) By exact path: `gws --user --update /path/to/repo work`
- [X] (C) Both name and path matching (consistent with existing `findRepositories` logic)
- [X] (D) Other (describe) When there are multipe repo matches, I'd like ot have a list that can be marked via the space bar, and one option in the list that is for selecting all.

## 3. Batch Update Behavior (--tag)

When using `--user --update --tag <tag> <profile>`, what should happen if some repos already have a local config?

- [ ] (A) Overwrite all matching repos with the new profile (no confirmation)
- [ ] (B) Overwrite all but show a summary of what changed
- [ ] (C) Prompt for confirmation before overwriting repos that already have local config
- [ ] (D) Skip repos that already have local config (only update repos without local config)
- [ ] (E) Other (describe)

## 4. Delete Behavior Details

When using `--user --delete`, what exactly should be removed from the repo's local `.git/config`?

- [ ] (A) Remove only `[user]` section (name, email) — leave signing config alone
- [ ] (B) Remove `[user]` section AND `[commit] gpgsign` signing config
- [ ] (C) Remove all gws-managed settings (user, email, signing key, gpgsign)
- [X] (D) Other (describe) default to leaving the siging key, but have a flag to remove it.

## 5. Delete with --tag (Batch Delete)

Should batch delete via tags also be supported? (e.g., `gws --user --delete --tag <tag>`)

- [X] (A) Yes, support batch delete with the same tag filtering as update
- [ ] (B) No, delete should only work on individual repos for safety
- [ ] (C) Yes, but require confirmation before batch delete
- [ ] (D) Other (describe)

## 6. Config.json Update After Operations

After `--user --update` or `--user --delete`, should the stored config.json also be updated?

- [X] (A) Yes, update the stored User/Email/UserSource in config.json to reflect the change (consistent with `user assign`)
- [ ] (B) No, only modify the local `.git/config` — let `--show-user` detect changes at display time
- [ ] (C) Update config.json for `--update`, but for `--delete` just clear the stored values so detection falls back to display-time
- [X] (D) Other (describe) What user assign are you referring to? 

## 7. Relationship to Existing `gws user assign`

The existing `gws user assign <repo> <profile>` already sets local user config. How should `--user --update` relate to it?

- [ ] (A) `--user --update` is a convenience alias — same underlying logic as `user assign`, just flag-based interface
- [ ] (B) `--user --update` replaces `user assign` entirely (deprecate assign)
- [ ] (C) They coexist independently — `--user --update` is simpler (no subdirs/dry-run), `user assign` keeps advanced features
- [X] (D) Other (describe) Are you sure this exists? if so can you demo the functionality to me? I do not see any commands that we have implemented that follow the format of gws user assign. When I manually test by doing a make build and make install, then running, it doesn't exist. Doing a grep of the repo also suggests it doesn't exist.

## 8. Proof Artifacts

What proof artifacts best demonstrate this feature works?

- [ ] (A) CLI output showing update/delete operations with before/after user config
- [ ] (B) CLI output + `git config --local --list` showing actual git config changes
- [ ] (C) CLI output + `gws --list --show-user` showing the repos reflect the changes
- [X] (D) All of the above
- [ ] (E) Other (describe)
