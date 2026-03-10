# Task 2.0 Proofs - User, Tagging, and Legacy Pages Update

## CLI: gws user --help

```
Short flags:
  gws user -l                                      # List all profiles
  gws user -a work --email work@company.com        # Add a new profile
  gws user -s work                                 # Show profile details
  gws user -d work                                 # Remove a profile

Flags:
  -a, --add                  Add a profile (equivalent to 'user add')
  -d, --delete               Remove a profile (equivalent to 'user remove')
      --email string         Email address for git commits (required with -a)
  -l, --list                 List all profiles (equivalent to 'user list')
      --name string          Git user name (defaults to profile name)
  -s, --show                 Show profile details (equivalent to 'user show')
      --sign-commits         Enable commit signing
      --signing-key string   GPG signing key ID
```

All 4 short-flag aliases (`-l`, `-a`, `-s`, `-d`) documented in commands-user.md. Inline add flags (`--email`, `--name`, `--signing-key`, `--sign-commits`) also documented.

## CLI: gws tag add --help

```
Flags:
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

Both `--path`/`-p` and `--repo`/`-r` targeting flags documented in commands-tagging.md for both `tag add` and `tag remove`.

## CLI: gws tag --help (parent command short flags)

```
Flags:
  -a, --add           Add a tag (equivalent to 'tag add')
  -d, --delete        Remove a tag (equivalent to 'tag remove')
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

All 4 short-flag aliases (`-a`, `-d`, `-p`, `-r`) documented in commands-tagging.md.

## Comparison: Legacy flags vs deprecated.go

`deprecated.go` contains 23 `depWarnings` entries. The legacy flags table now includes:
- 9 command flags (list, init, add, recursive, refresh, print-workspace, go, add-tag, remove-tag)
- 7 list filter flags (type, tag, name, path, output, status, show-user) + `-V`
- 8 user management flags (user, list-users, update, delete, all, verbose, git-name, git-email)

Total: 25 table rows covering all 23 depWarnings entries (some entries like `-V` map to separate rows).
