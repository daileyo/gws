# 08 Task 1.0 Proof Artifacts - CLI Flag Infrastructure and Dispatch

## CLI: `gws --help` output demonstrates new flags in "User Operations" section

```
User Operations (require --user):
      --all                Also remove signing config when deleting (requires --delete)
  -D, --delete             Delete local git user config from repositories (requires --user)
      --git-email string   Inline git user.email for --user --update
      --git-name string    Inline git user.name for --user --update
  -u, --update             Update local git user config for repositories (requires --user)
      --user               Enable user operations (requires --update or --delete)
      --verbose            Show detailed output for user operations
```

## CLI: `gws --user --update` without required args demonstrates validation

```
$ gws --user --update
Error: not yet implemented
```

(Stub correctly dispatches to `runUserUpdate` which returns "not yet implemented" — validation passes, dispatch works.)

## CLI: `gws --user --update --delete` demonstrates mutual exclusivity error

```
$ gws --user --update --delete
Error: --update and --delete are mutually exclusive
```

## CLI: `gws --update` without `--user` demonstrates flag dependency error

```
$ gws --update
Error: --update/-u requires --user to be set
```

## Additional Validation Tests

```
$ gws --delete
Error: --delete/-D requires --user to be set

$ gws --user
Error: --user requires either --update/-u or --delete/-D

$ gws --user --update --all
Error: --all requires --delete/-D to be set
```

## Test Results

```
=== RUN   TestUserFlagValidation
=== RUN   TestUserFlagValidation/update_without_user_returns_error
=== RUN   TestUserFlagValidation/delete_without_user_returns_error
=== RUN   TestUserFlagValidation/user_without_update_or_delete_returns_error
=== RUN   TestUserFlagValidation/user_update_and_delete_are_mutually_exclusive
=== RUN   TestUserFlagValidation/all_without_delete_returns_error
=== RUN   TestUserFlagValidation/user_is_mutually_exclusive_with_list
--- PASS: TestUserFlagValidation (0.00s)

=== RUN   TestUserFlagsRegistered
=== RUN   TestUserFlagsRegistered/user
=== RUN   TestUserFlagsRegistered/update
=== RUN   TestUserFlagsRegistered/delete
=== RUN   TestUserFlagsRegistered/all
=== RUN   TestUserFlagsRegistered/verbose
=== RUN   TestUserFlagsRegistered/git-name
=== RUN   TestUserFlagsRegistered/git-email
--- PASS: TestUserFlagsRegistered (0.00s)

All tests PASS.
```
