# 08 Task 2.0 Proof Artifacts - Single Repository User Update

## CLI: `gws --user -u <repo> <profile>` demonstrates single repo update

```
$ gws --user -u dratdns work
dratdns: user.name "daileyo" → "Work User", user.email "dale13@gmail.com" → "work@company.com"
```

## CLI: `git config --local --list` demonstrates actual .git/config changes

```
$ git -C /path/to/dratdns config --local --list
core.repositoryformatversion=0
core.filemode=false
core.bare=false
core.logallrefupdates=true
core.symlinks=false
core.ignorecase=true
user.name=Work User
user.email=work@company.com
```

## CLI: `gws -l --show-user` demonstrates the repo reflects new user info

```
$ gws -l --show-user --name dratdns
Found 1 repository:

NAME     USER               EMAIL             SIGN  TYPE     VISIBILITY  TAGS  PATH
-------  -----------------  ----------------  ----  -------  ----------  ----  ----
dratdns  Work User (local)  work@company.com        unknown  unknown     -     /path/to/dratdns
```

Note the `(local)` indicator showing the user config is set locally.

## CLI: Inline value support

```
$ gws --user -u dratdns --git-name "Custom Name" --git-email "custom@test.com"
dratdns: user.name "Work User" → "Custom Name", user.email "work@company.com" → "custom@test.com"

$ git -C /path/to/dratdns config --local --list
...
user.name=Custom Name
user.email=custom@test.com
```

## Test Results

```
=== RUN   TestResolveProfile
=== RUN   TestResolveProfile/named_profile_found_in_stored_profiles
=== RUN   TestResolveProfile/named_profile_not_found
=== RUN   TestResolveProfile/inline_values_only
=== RUN   TestResolveProfile/inline_email_only_uses_email_prefix_as_name
=== RUN   TestResolveProfile/profile_name_with_inline_overrides
=== RUN   TestResolveProfile/no_profile_name_and_no_inline_values_returns_error
=== RUN   TestResolveProfile/tag_mode_uses_first_arg_as_profile_name
--- PASS: TestResolveProfile (0.00s)

=== RUN   TestRunUserUpdate_SingleRepo
--- PASS: TestRunUserUpdate_SingleRepo (0.02s)
=== RUN   TestRunUserUpdate_InlineValues
--- PASS: TestRunUserUpdate_InlineValues (0.02s)
=== RUN   TestRunUserUpdate_MultipleRepoMatch
--- PASS: TestRunUserUpdate_MultipleRepoMatch (0.04s)
=== RUN   TestRunUserUpdate_NoMatch
--- PASS: TestRunUserUpdate_NoMatch (0.00s)
=== RUN   TestRunUserUpdate_NoArgs
--- PASS: TestRunUserUpdate_NoArgs (0.00s)

All tests PASS.
```
