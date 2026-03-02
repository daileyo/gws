# Task 1.0 Proof Artifacts - Extract Tag Business Logic into Parameterized Functions

## Refactored Function Signatures

### Before
```go
func runTag(_ *cobra.Command, args []string) error     // tag.go
func runUntag(_ *cobra.Command, args []string) error    // untag.go
```

### After
```go
func runAddTag(repoIdentifier, tag string) error        // tag.go
func runRemoveTag(repoIdentifier, tag string) error     // untag.go
```

## Call Site Updates

In `main.go` RunE, the call sites now validate arg count and pass params directly:

```go
if flagAddTag {
    if len(args) != 2 {
        return fmt.Errorf("--add-tag requires exactly 2 arguments: <repository> <tag>")
    }
    return runAddTag(args[0], args[1])
}
if flagRemoveTag {
    if len(args) != 2 {
        return fmt.Errorf("--remove-tag requires exactly 2 arguments: <repository> <tag>")
    }
    return runRemoveTag(args[0], args[1])
}
```

## Test Results

```
$ go vet ./...  ✓
$ go test ./... -count=1
ok  github.com/daileyo/gws/cmd/git-workspace   0.319s
ok  github.com/daileyo/gws/internal/classifier  0.006s
ok  github.com/daileyo/gws/internal/config      0.006s
ok  github.com/daileyo/gws/internal/discovery    0.031s
ok  github.com/daileyo/gws/internal/filter       0.007s
ok  github.com/daileyo/gws/internal/git          0.041s
ok  github.com/daileyo/gws/internal/user         0.129s
```
