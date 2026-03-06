# Task 4.0 Proof Artifacts - Default No-Args Behavior

## Test Results

All tests pass:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.412s
```

## Changes Made

### `main.go` - Root command RunE

The no-args path now calls `runList(ListOptions{OutputFormat: "table"})` instead of printing the workspace summary. This triggers the multi-column default view.

Before:
```go
fmt.Printf("Workspace: %s\n", cfg.Workspace)
fmt.Printf("Repositories: %d\n", len(cfg.Repositories))
fmt.Println("Use 'gws --help' to see available commands")
```

After:
```go
return runList(ListOptions{OutputFormat: "table"})
```

### Navigation preserved

- `gws my-repo` still navigates (positional arg path unchanged)
- `gws -q my-repo` still works (quiet flag validation unchanged)
- `gws -q` (no args) still errors: "can only be used with navigation"

### Usage template updated

- Long description now mentions `gws` defaults to listing repos
- Commands section shows `gws` as first example with "(default, same as gws list)"
- Added `gws list -v` and `gws list -vv` examples
