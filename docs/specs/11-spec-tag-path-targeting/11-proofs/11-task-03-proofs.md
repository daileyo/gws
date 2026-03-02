# 11-task-03-proofs: Path Output Formatting, Shell Completion, Help Text, and CI

## CLI: Output Format — Path Targeting (name + path per line)

```
$ gws tag add --path /home/user/work backend
Added tag 'backend' to 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

## CLI: Output Format — Default / --repo (name only per line, unchanged)

```
$ gws tag add my-repo backend
Added tag 'backend' to 1 repository
  - my-repo
```

## CLI: Output Format — Remove with path targeting

```
$ gws tag remove --path /home/user/work backend
Removed tag 'backend' from 2 repositories
  - api-service  /home/user/work/api-service
  - api-gateway  /home/user/work/api-gateway
```

## CLI: Updated Help Text — `gws tag add --help`

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

Flags:
  -h, --help          help for add
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

## CLI: Updated Help Text — `gws tag remove --help`

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

Flags:
  -h, --help          help for remove
  -p, --path string   Match repositories by path prefix or substring (case-sensitive)
  -r, --repo string   Match repositories by name (partial, case-insensitive)
```

## Shell Completion: `completeRepoPaths` registered for `--path` flag

`completeRepoPaths` is registered on `tagAddCmd` and `tagRemoveCmd` via
`RegisterFlagCompletionFunc("path", ...)`. It loads config and returns paths with
`strings.HasPrefix(repo.Path, toComplete)`, matching the same pattern as `completeRepoNames`.

## Test Results

```
$ go test ./cmd/git-workspace/
ok  	github.com/daileyo/gws/cmd/git-workspace	0.618s
```

## CI Results

```
Running go vet...
go vet ./...
Running linter...
Running tests with race detector...
go test -v -race ./...
[... all tests PASS ...]
Tests complete
All CI checks passed!
make ci exit: 0
```

## Verification

- Output shows `  - name  path` when `--path` is used, `  - name` otherwise — format verified above
- `completeRepoPaths` function implemented and registered on both subcommands
- `Long` help strings and `Use` fields updated on `tagCmd`, `tagAddCmd`, `tagRemoveCmd`
- All tests pass with no failures
- `make ci` exits 0 with no lint warnings — full quality gate passes
