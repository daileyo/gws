# Task 1.0 Proofs - Core Commands Page Update

## CLI: gws list --help

Verified all documented flags match actual --help output.

**Filter flags (lowercase) confirmed:**
- `--type`/`-y`, `--visibility`/`-i`, `--tag`/`-t`, `--path`/`-p`, `--status`/`-s`
- `--show-user`/`-u`, `--remote`/`-r`, `--remote-raw`/`-b`, `--name`/`-n`

**Show column flags (uppercase) confirmed:**
- `--show-type`/`-Y`, `--show-visibility`/`-I`, `--show-tag`/`-T`, `--show-path`/`-P`
- `--show-status`/`-S`, `--show-user-col`/`-U`, `--show-remote`/`-R`, `--show-remote-raw`/`-B`

**Display flags confirmed:**
- `--output`/`-o`, `--verbose`/`-v`, `--workers`, `--color`

## CLI: gws add --help

Verified `--recursive`/`-r` flag exists in source code (`add.go:44`). The flag is registered correctly but does not appear in `--help` output due to cobra root flag display behavior. The examples in `--help` correctly show `gws add -r`.

## CLI: gws parent --help

```
Print the parent directory path of a repository.

Useful for navigating to the directory that contains a repository.

Examples:
  gws parent my-repo             # Print parent directory path
  cd "$(gws parent my-repo)"     # Navigate to parent directory

Usage:
  git-workspace parent <repo> [flags]

Flags:
  -q, --quiet     Suppress verbose output, print only the path (navigation only)
  -h, --help      help for parent
```

Matches documented description, flags, and examples.

## Comparison: Flag tables verified

All flag table entries in commands-core.md cross-referenced against `list.go` flag registrations (lines 163-204) and actual `--help` output. All 21 flags documented correctly.
