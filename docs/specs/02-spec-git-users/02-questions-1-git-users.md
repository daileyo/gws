# 02 Questions Round 1 - Git Users

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. User Profile Management

How should gws manage the git user profiles (user.name, user.email, signing key)?

- [ ] (A) **gws manages profiles internally only** - Store user profiles in gws config (~/.gws/config.json) and read-only detect which profile is active for each repo based on existing gitconfig/includeIf setup
- [ ] (B) **gws manages gitconfig files** - gws creates/updates the supplemental .gitconfig-* files in subdirectories AND updates ~/.gitconfig with includeIf directives
- [X] (C) **Hybrid approach** - gws stores profile definitions in its config, and can optionally write gitconfig files on user request
- [ ] (D) Other (describe)

## 2. Profile Detection vs Assignment

When gws discovers/scans a repository, what user information should it track?

- [ ] (A) **Effective config only** - Detect what user.name/user.email git would actually use for that repo (the resolved config)
- [ ] (B) **Last commit author** - Track the author of the most recent commit in the repo
- [ ] (C) **Both** - Store both the effective git config AND the last commit author (they may differ)
- [ ] (D) **Explicit assignment only** - Only track user profiles that have been explicitly assigned via gws commands
- [X] (E) Other - Track effective, but mark it with an indicator if the effective differs from what was previously stored in the config. also include an option to sync the stored with the effective.

## 3. Commit Signing Detection

What commit signing information should gws track?

- [X] (A) **Signing configuration** - Whether commit signing is configured (user.signingkey, commit.gpgsign) in the effective config
- [ ] (B) **Actual signing status** - Whether recent commits in the repo are actually signed
- [ ] (C) **Both** - Track both the configuration AND actual signing status of commits
- [ ] (D) **None** - Signing information is out of scope for this feature
- [ ] (E) Other (describe)

## 4. Changing User Assignment

When a user wants to change the git user/email for a repository, what should happen?

- [ ] (A) **Move repo to appropriate subdirectory** - Physically move the repo folder to match the includeIf directory structure (e.g., ~/gws/work/ for work profile)
- [ ] (B) **Set repo-local config** - Run `git config user.name/email` in the repo's local .git/config (no directory movement)
- [ ] (C) **User choice** - Ask the user which approach they prefer each time
- [ ] (D) **Recommend only** - gws tells the user what to do but doesn't make changes
- [X] (E) Other - I would prefer sub directories to only be used if absolutely required. I.E. best case scenario is that the gws directory should only have repos in it; no sub directories of repos. However, if sub directories are absolutely required in order to support having multiple users, then lets do that and move the repos to the sub directories when a user wants to change the associated signer. If sub directories are required, then gws should handle creating them and updating the gitconfig accordingly.

## 5. Subdirectory Organization

You mentioned "only create sub directories to support multiple users if absolutely required." Can you clarify?

- [ ] (A) **Flat by default, prompt for subdirs** - Keep all repos in main workspace; only suggest subdirectory organization if user explicitly has multiple profiles
- [ ] (B) **Profile-based organization** - Organize into subdirectories by profile/identity (e.g., ~/gws/personal/, ~/gws/work/)
- [ ] (C) **Source-based (current)** - Keep current organization by source (ado/, gh/, gitlab/) which happens to align with some profiles
- [ ] (D) **No reorganization** - Never move repos; use repo-local config exclusively for profile assignment
- [X] (E) Other Do a, if at all possible. otherwise, create subdirectories for each git user profile, and move repos to the appropriate subdirectories.

## 6. Profile Definition

How should user profiles be defined in gws?

- [ ] (A) **Auto-detect from gitconfig** - Parse ~/.gitconfig and any includeIf files to discover existing profiles
- [ ] (B) **Explicit creation** - User must explicitly create profiles via `gws user add "Work" --email work@company.com --name "John Doe"`
- [X] (C) **Hybrid** - Auto-detect existing profiles AND allow creating new ones
- [ ] (D) Other (describe)

## 7. Display and Filtering

How should user information be displayed and filtered?

- [X] (A) **New columns in list** - Add USER and EMAIL columns to `gws list` output (optional via flag)
- [ ] (B) **New filter flags** - Add `--user` and `--email` filter flags to find repos by profile
- [ ] (C) **Both columns and filters** - Include both display columns and filter capabilities
- [ ] (D) **Separate command** - New `gws users` command to manage and view user/profile information
- [ ] (E) Other (describe)

## 8. Edge Cases

How should gws handle these edge cases?

### 8a. Repository with no user configured (inherits global default)
- [X] (A) Show the inherited global user explicitly
- [ ] (B) Show "default" or "global" as the profile name
- [ ] (C) Show as empty/unassigned

### 8b. Repository with local .git/config user override
- [X] (A) Show the local config user, mark as "local override"
- [ ] (B) Treat as a separate profile
- [ ] (C) Warn user about inconsistency with directory-based profile

### 8c. Multiple remotes with potentially different user expectations
- [ ] (A) Use origin remote's typical user (if detectable)
- [ ] (B) Show all users that have committed
- [X] (C) Ignore remote-based expectations, only show configured user

## 9. Commands

What new commands should be added?

- [ ] (A) **Minimal** - Just extend existing commands (list --user, tag-like assignment)
- [X] (B) **User subcommand** - `gws user list`, `gws user add`, `gws user assign <repo> <profile>`
- [ ] (C) **Profile subcommand** - `gws profile list`, `gws profile create`, `gws profile assign`
- [ ] (D) Other (describe)

## 10. Proof Artifacts

What would demonstrate this feature is working correctly?

- [ ] (A) **CLI output** - `gws list --user` shows user/email columns with correct values
- [ ] (B) **Profile management** - `gws user list` shows all detected/configured profiles
- [ ] (C) **Assignment** - `gws user assign <repo> <profile>` successfully changes repo's user
- [ ] (D) **Detection accuracy** - Making a commit with a specific user shows correctly in gws
- [X] (E) All of the above
- [ ] (F) Other (describe)

---

## Additional Context

Please add any additional context, constraints, or preferences here:

```
[Your additional notes]
```
