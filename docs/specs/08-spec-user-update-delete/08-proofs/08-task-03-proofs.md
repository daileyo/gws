# 08 Task 3.0 Proof Artifacts - Single Repository User Delete

## CLI: `gws --user --delete <repo>` demonstrates local config removal

```
$ gws --user -u dratdns work
dratdns: user.name "daileyo" → "Work User", user.email "dale13@gmail.com" → "work@company.com"

$ gws --user --delete dratdns
dratdns: removed local user config (now using global: "daileyo" <dale13@gmail.com>)
```

## CLI: `git config --local --list` confirms user section removed

```
$ git -C /path/to/dratdns config --local --list | grep -E "user\.|commit\."
(empty — no local user config)
```

## CLI: `gws -l --show-user` demonstrates repo shows default user

```
$ gws -l --show-user --name dratdns
Found 1 repository:

NAME     USER       EMAIL             SIGN  TYPE     VISIBILITY  TAGS  PATH
-------  ---------  ----------------  ----  -------  ----------  ----  ----
dratdns  daileyo *  dale13@gmail.com        unknown  unknown     -     /path/to/dratdns
```

No `(local)` indicator — repo falls back to global default with `*` marker.

## CLI: `gws --user --delete <repo> --all` demonstrates signing config removal

```
$ git -C /path/to/dratdns config --local --list | grep -E "user\.|commit\."
user.name=Work User
user.email=work@company.com
user.signingkey=ABCD1234
commit.gpgsign=true

$ gws --user --delete dratdns --all
dratdns: removed local user config + signing config (now using global: "daileyo" <dale13@gmail.com>)

$ git -C /path/to/dratdns config --local --list | grep -E "user\.|commit\."
(empty — all user and signing config removed)
```

## Test Results

```
=== RUN   TestRemoveGitConfigKey
=== RUN   TestRemoveGitConfigKey/remove_key_from_section
=== RUN   TestRemoveGitConfigKey/remove_last_key_removes_section_header
=== RUN   TestRemoveGitConfigKey/remove_key_preserves_other_sections
=== RUN   TestRemoveGitConfigKey/remove_nonexistent_key_is_no-op
=== RUN   TestRemoveGitConfigKey/remove_from_nonexistent_section_is_no-op
=== RUN   TestRemoveGitConfigKey/case_insensitive_key_match
--- PASS: TestRemoveGitConfigKey (0.00s)

=== RUN   TestDeleteLocal
=== RUN   TestDeleteLocal/remove_name_and_email_only
=== RUN   TestDeleteLocal/remove_all_including_signing
=== RUN   TestDeleteLocal/preserves_non-user_sections
=== RUN   TestDeleteLocal/error_for_non-git_repo
--- PASS: TestDeleteLocal (0.05s)

=== RUN   TestRunUserDelete_SingleRepo
--- PASS: TestRunUserDelete_SingleRepo (0.02s)
=== RUN   TestRunUserDelete_WithAll
--- PASS: TestRunUserDelete_WithAll (0.02s)
=== RUN   TestRunUserDelete_NoLocalConfig
--- PASS: TestRunUserDelete_NoLocalConfig (0.02s)
=== RUN   TestRunUserDelete_NoMatch
--- PASS: TestRunUserDelete_NoMatch (0.00s)
=== RUN   TestRunUserDelete_NoArgs
--- PASS: TestRunUserDelete_NoArgs (0.00s)
=== RUN   TestRunUserDelete_ConfigJsonUpdated
--- PASS: TestRunUserDelete_ConfigJsonUpdated (0.02s)

All tests PASS.
```
