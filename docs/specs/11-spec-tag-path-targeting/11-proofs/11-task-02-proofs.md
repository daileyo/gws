# 11-task-02-proofs: Register Flags and Update Command RunE Logic

## CLI: Help Output

### `gws tag add --help`

```
Add a tag to all repositories matching the identifier.

With no flags, the repository identifier matches by partial name (case-insensitive) or exact path.
Use --path to match by path prefix or substring (case-sensitive).
Use --repo to explicitly match by partial name (case-insensitive).
Combine --path and --repo to require both conditions (AND logic).

Examples:
  gws tag add my-repo work
  gws tag add api backend
  gws tag add --path /home/user/work backend
  gws tag add --repo api backend
  gws tag add --repo api --path /work backend

Usage:
  git-workspace tag add [--path <path>] [--repo <repo>] <tag> | add <repo> <tag> [flags]

Flags:
  -h, --help          help for add
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

### `gws tag remove --help`

```
Remove a tag from all repositories matching the identifier.

With no flags, the repository identifier matches by partial name (case-insensitive) or exact path.
Use --path to match by path prefix or substring (case-sensitive).
Use --repo to explicitly match by partial name (case-insensitive).
Combine --path and --repo to require both conditions (AND logic).

Examples:
  gws tag remove my-repo work
  gws tag remove api backend
  gws tag remove --path /home/user/work backend
  gws tag remove --repo api backend
  gws tag remove --repo api --path /work backend

Usage:
  git-workspace tag remove [--path <path>] [--repo <repo>] <tag> | remove <repo> <tag> [flags]

Flags:
  -h, --help          help for remove
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

### `gws tag --help` (shows --path and --repo on tagCmd for use with -a/-d)

```
Manage tags on tracked repositories.

Examples:
  gws tag add my-repo work                          # Add the "work" tag by name
  gws tag -a my-repo work                           # Same, using short flag
  gws tag add --path /home/user/work backend        # Add tag to all repos under a path
  gws tag add --repo api backend                    # Add tag to repos matching "api" by name
  gws tag add --repo api --path /work backend       # Add tag matching both conditions
  gws tag remove my-repo work                       # Remove the "work" tag
  gws tag -d my-repo work                           # Same, using short flag
  gws tag remove --path /home/user/work backend     # Remove tag from all repos under a path
  gws tag                                           # Show this help

Flags:
  -a, --add           Add a tag (equivalent to 'tag add')
  -d, --delete        Remove a tag (equivalent to 'tag remove')
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

## CLI: Functional Behavior

Test config:
```json
{
  "repositories": [
    {"name": "api-service",  "path": "/home/user/work/api-service",         "tags": []},
    {"name": "api-gateway",  "path": "/home/user/work/api-gateway",         "tags": []},
    {"name": "personal-app", "path": "/home/user/personal/personal-app",    "tags": []}
  ]
}
```

### Path-based add (prefix match, shows path+name in output)

```
$ gws tag add --path /home/user/work backend
Added tag 'backend' to 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

### Path-based remove (shows path+name in output)

```
$ gws tag remove --path /home/user/work backend
Removed tag 'backend' from 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

### Explicit --repo flag (name-only output, name match)

```
$ gws tag add --repo api work
Added tag 'work' to 2 repositories
  - api-service
  - api-gateway
```

### Combined --repo + --path (AND logic)

```
$ gws tag add --repo api --path /home/user/work/api-service combined
Added tag 'combined' to 1 repository
  - api-service  /home/user/work/api-service
```

### Legacy positional form (backward compatibility)

```
$ gws tag add personal-app personal
Added tag 'personal' to 1 repository
  - personal-app
```

### Short flag -a with --path

```
$ gws tag -a --path /home/user/work myshortflag
Added tag 'myshortflag' to 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

### Short flag -d with --path

```
$ gws tag -d --path /home/user/work myshortflag
Removed tag 'myshortflag' from 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

### Error: no match on --path

```
$ gws tag add --path /nonexistent backend
Error: no repositories found matching path: /nonexistent
exit: 1
```

### Error: no match on --repo

```
$ gws tag add --repo ghost backend
Error: no repositories found matching repo: ghost
exit: 1
```

### Error: AND no match (--repo + --path)

```
$ gws tag add --repo ghost --path /nonexistent backend
Error: no repositories found matching repo: ghost and path: /nonexistent
exit: 1
```

## Test Results

```
$ go test ./cmd/git-workspace/
ok  	github.com/daileyo/gws/cmd/git-workspace	0.618s
```

All tests pass with no failures.

## Verification

- `--path`/`-p` and `--repo`/`-r` appear in `tag add --help`, `tag remove --help`, and `tag --help`
- Path-based add tags only repos under the given path, shows `name  path` format
- Path-based remove works symmetrically with same output format
- `--repo` flag performs partial name match, shows name-only output (no path)
- `--repo` + `--path` combined applies AND logic (only 1 of 2 api-* repos matched)
- Legacy positional form `tag add <repo> <tag>` is fully unchanged
- Short flags `-a`/`-d` combined with `--path` work correctly
- All three no-match error messages are contextually correct
