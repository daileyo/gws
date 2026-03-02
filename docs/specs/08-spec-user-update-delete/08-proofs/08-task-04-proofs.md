# 08 Task 4.0 Proof Artifacts - Batch Update and Delete via Tags

## CLI: `gws --user -u --tag <tag> <profile>` demonstrates batch update

```
$ gws -l --show-user --tag demo-tag
Found 2 repositories:

NAME        USER       EMAIL             SIGN  TYPE     VISIBILITY  TAGS      PATH
----------  ---------  ----------------  ----  -------  ----------  --------  ----
dratdns     daileyo *  dale13@gmail.com        unknown  unknown     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratdns
dratreprox  daileyo *  dale13@gmail.com        github   private     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratreprox

$ gws --user -u --tag demo-tag --git-name "Demo User" --git-email "demo@example.com"
dratdns: user.name "daileyo" → "Demo User", user.email "dale13@gmail.com" → "demo@example.com"
dratreprox: user.name "daileyo" → "Demo User", user.email "dale13@gmail.com" → "demo@example.com"

Updated 2 repositories
```

## CLI: `gws --user --delete --tag <tag>` demonstrates batch delete

```
$ gws --user --delete --tag demo-tag
dratdns: removed local user config (now using global: "daileyo" <dale13@gmail.com>)
dratreprox: removed local user config (now using global: "daileyo" <dale13@gmail.com>)

Deleted local user config from 2 repositories
```

## CLI: `gws -l --show-user --tag <tag>` demonstrates all tagged repos reflect changes

After batch update:

```
$ gws -l --show-user --tag demo-tag
Found 2 repositories:

NAME        USER               EMAIL             SIGN  TYPE     VISIBILITY  TAGS      PATH
----------  -----------------  ----------------  ----  -------  ----------  --------  ----
dratdns     Demo User (local)  demo@example.com        unknown  unknown     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratdns
dratreprox  Demo User (local)  demo@example.com        github   private     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratreprox
```

After batch delete:

```
$ gws -l --show-user --tag demo-tag
Found 2 repositories:

NAME        USER       EMAIL             SIGN  TYPE     VISIBILITY  TAGS      PATH
----------  ---------  ----------------  ----  -------  ----------  --------  ----
dratdns     daileyo *  dale13@gmail.com        unknown  unknown     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratdns
dratreprox  daileyo *  dale13@gmail.com        github   private     demo-tag  /mnt/c/Users/daileyo/gws/personal/dratreprox
```

No `(local)` indicator — repos fall back to global default with `*` marker.

## CLI: No-match tag error message

```
$ gws --user -u --tag nonexistent-tag --git-email "test@test.com"
Error: no repositories found with tag(s): nonexistent-tag
```

## Test Results

```
=== RUN   TestRunUserUpdate_BatchByTag
--- PASS: TestRunUserUpdate_BatchByTag (0.05s)
=== RUN   TestRunUserUpdate_BatchByMultipleTags
--- PASS: TestRunUserUpdate_BatchByMultipleTags (0.03s)
=== RUN   TestRunUserUpdate_BatchNoTagMatch
--- PASS: TestRunUserUpdate_BatchNoTagMatch (0.00s)
=== RUN   TestRunUserDelete_BatchByTag
--- PASS: TestRunUserDelete_BatchByTag (0.05s)
=== RUN   TestRunUserDelete_BatchNoTagMatch
--- PASS: TestRunUserDelete_BatchNoTagMatch (0.00s)
=== RUN   TestRunUserDelete_BatchQuietSuppressesOutput
--- PASS: TestRunUserDelete_BatchQuietSuppressesOutput (0.02s)

All batch tests PASS.
```

## Full Test Suite

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.406s
ok  	github.com/daileyo/gws/internal/classifier	0.006s
ok  	github.com/daileyo/gws/internal/config	0.007s
ok  	github.com/daileyo/gws/internal/discovery	0.026s
ok  	github.com/daileyo/gws/internal/filter	0.007s
ok  	github.com/daileyo/gws/internal/git	0.032s
ok  	github.com/daileyo/gws/internal/user	0.138s
```
