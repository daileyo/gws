# 06 Task 1.0 Proofs - commit-msg Hook Script

## Overview

The `commit-msg` hook at `.githooks/commit-msg` pads conventional commit type tokens
so that all descriptions align to column 13 (token width 12 + space). The longest
possible token is `refactor!:` (10 chars); TARGET_WIDTH = 10 + 2 = 12.

## CLI Output

All tests run by invoking the hook script directly against temp files (`sh .githooks/commit-msg <tmpfile>`).

### Test 1: Standard type — `feat:`

```
$ printf 'feat: add new command\n' > /tmp/t && sh .githooks/commit-msg /tmp/t && cat /tmp/t
feat:       add new command
```

`feat:` = 5 chars + 7 spaces = 12 char prefix → description at col 13 ✓

### Test 2: Shortest type — `ci:`

```
$ printf 'ci: fix workflow\n' > /tmp/t && sh .githooks/commit-msg /tmp/t && cat /tmp/t
ci:         fix workflow
```

`ci:` = 3 chars + 9 spaces = 12 char prefix → description at col 13 ✓

### Test 3: Breaking change, longest token — `refactor!:`

```
$ printf 'refactor!: breaking API change\n' > /tmp/t && sh .githooks/commit-msg /tmp/t && cat /tmp/t
refactor!:  breaking API change
```

`refactor!:` = 10 chars + 2 spaces = 12 char prefix → description at col 13 ✓

### Test 4: Longest plain type — `refactor:`

```
$ printf 'refactor: simplify discovery logic\n' > /tmp/t && sh .githooks/commit-msg /tmp/t && cat /tmp/t
refactor:   simplify discovery logic
```

`refactor:` = 9 chars + 3 spaces = 12 char prefix → description at col 13 ✓

### Test 5: Non-conforming message — pass-through with warning

```
$ printf 'WIP: not formatted\n' > /tmp/t && sh .githooks/commit-msg /tmp/t 2>&1 && cat /tmp/t
commit-msg: message does not match conventional commit format, skipping formatting
WIP: not formatted
```

Warning printed to stderr; original message unchanged; exit 0 ✓

### Test 6: Multi-line commit — body preserved

```
$ printf 'fix: correct requirements path\n\nThis fixes the incorrect path shown\nin the project structure tree.\n' > /tmp/t
$ sh .githooks/commit-msg /tmp/t && cat /tmp/t
fix:        correct requirements path

This fixes the incorrect path shown
in the project structure tree.
```

Subject line formatted; blank line and body preserved unchanged ✓

## Alignment Verification

All type tokens aligned side-by-side for visual confirmation:

```
ci:         fix workflow
fix:        correct path
feat:       add new command
docs:       update readme
perf:       improve cache
test:       add unit tests
chore:      update deps
style:      fix formatting
refactor:   simplify logic
refactor!:  breaking change
```

Every description starts at column 13.

## File Permissions

```
$ ls -la .githooks/
-rwxr-xr-x  commit-msg
-rwxr-xr-x  pre-push
```

Script is executable ✓

## Hook Activation

`git config core.hooksPath` returns `.githooks` — hook is active for all commits on this machine ✓
