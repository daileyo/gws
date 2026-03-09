# Task 3.0 Proof Artifacts - Help Text Updates

## CLI Output - Help

```
$ gws list --help

...
Examples:
  gws list -r                       # Show formatted remote URL column
  gws list -R                       # Show raw remote URL column
  gws list -rR                      # Show raw remote URL (override)
...

Flags:
  -r, --remote  Show formatted remote URL column, or filter by remote URL pattern
  -R, --remote-raw  Show raw remote URL column, or filter by raw remote URL pattern
```

## Verification

- `-r` flag description updated to "Show formatted remote URL column"
- `-R`/`--remote-raw` flag appears in help with correct description
- Examples section includes `-r`, `-R`, and `-rR` usage
