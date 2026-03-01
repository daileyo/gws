# 07 Questions Round 1 - Git Users Rework

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Default User Indicator in List

What does "default user" mean in this context, and how should it be indicated?

- [X] (A) **Global gitconfig user** - The user defined in `~/.gitconfig` (not via includeIf or local config) is considered the "default" and repos using it get an indicator like `(default)` or `*` next to the user name
- [ ] (B) **A designated default profile** - Add a `default` field to a stored profile in config, and mark repos using that profile with an indicator
- [ ] (C) **The global/fallback source** - Repos where `UserSource` is `global` (no local or includeIf override) get marked to show they're using the fallback/default config
- [ ] (D) Other (describe)

## 2. Default User Indicator Style

How should the default user be visually indicated in `gws -l --show-user` output?

- [X] (A) **Asterisk marker** - `John Doe *` or `* John Doe`
- [ ] (B) **Label suffix** - `John Doe (default)`
- [ ] (C) **Separate column** - Add a `DEFAULT` column with `✓` or blank
- [ ] (D) Other (describe)

## 3. IncludeIf Evaluation During Init/Add/Refresh

Currently, user info is detected via `git.GetUserConfig()`. When you say "includeIf configs should be evaluated," do you mean:

- [ ] (A) **Auto-detect the effective user** - During init/add/refresh, resolve which includeIf applies to each repo based on its path and store the resulting user/email with `UserSource: includeif`. This means the stored config reflects the *effective* user for each repo's location
- [ ] (B) **Match repos to stored profiles** - During init/add/refresh, check if a repo's path falls under an includeIf gitdir pattern and if so, automatically link it to the matching stored profile
- [X] (C) **Both** - Detect the effective user AND auto-link to a matching stored profile if one exists
- [ ] (D) Other (describe)

## 4. Local Config Evaluation During Init/Add/Refresh

For repos with local `.git/config` author info, what should happen during init/add/refresh?

- [ ] (A) **Store local config values** - Read user.name/user.email from `.git/config` and store them in the repo entry with `UserSource: local`, overriding any global/includeIf values
- [ ] (B) **Store and match to profile** - Read local config AND try to match it to a stored profile by email/name
- [X] (C) **Just detect and display** - Read local config for display purposes in `--show-user` but don't persist it in the config
- [ ] (D) Other (describe)

## 5. Shorthand Conflict: --remove-tag -d vs --add-tag -d

There's a conflict: `--add-tag` currently uses `-d` as its shorthand. You want `--remove-tag` to use `-d`. How should this be resolved?

- [ ] (A) **Swap shorthands** - Change `--add-tag` to a different shorthand (e.g., `-A`) and give `-d` to `--remove-tag`
- [ ] (B) **Remove --add-tag shorthand** - Remove the `-d` shorthand from `--add-tag` entirely, give `-d` to `--remove-tag`
- [X] (C) **Use a different shorthand for --remove-tag** - Keep `--add-tag` as `-d` and pick something else for `--remove-tag` (e.g., `-x`, `-D`)
- [ ] (D) Other (describe)

## 6. --user to --show-user Scope

The `--user` flag is currently a filter flag for `--list` (shows USER, EMAIL, SIGN columns). Renaming it to `--show-user`:

- [X] (A) **Simple rename only** - Just rename `--user` to `--show-user`, keep same behavior, no shorthand
- [ ] (B) **Rename with shorthand** - Rename to `--show-user` and add a shorthand like `-u` (freed up from --remove-tag)
- [ ] (C) **Rename and enhance** - Rename and also add the default user indicator behavior from question 1 to this flag's output
- [ ] (D) Other (describe)

## 7. Test Coverage Scope

What level of test coverage do you want for these changes?

- [X] (A) **Unit tests only** - Test individual functions (includeIf evaluation, local config detection, default user marking, flag changes)
- [ ] (B) **Unit + integration tests** - Unit tests plus CLI-level integration tests that verify end-to-end behavior of commands
- [ ] (C) **Match existing patterns** - Follow whatever test patterns already exist in the codebase (appears to be primarily unit tests with table-driven patterns)
- [ ] (D) Other (describe)

## 8. Proof Artifacts

What proof artifacts would best demonstrate these features work?

- [ ] (A) **CLI output screenshots/captures** - Show `gws -l --show-user` output with default indicators, includeIf-resolved users, and local config users
- [ ] (B) **Test results** - Passing test suite output demonstrating all scenarios
- [X] (C) **Both CLI output and test results** - Comprehensive proof with both visual output and test coverage
- [ ] (D) Other (describe)
