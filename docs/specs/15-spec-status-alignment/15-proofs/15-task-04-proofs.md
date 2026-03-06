# Task 4.0 Proof: Apply Colors to Status Icons

## CLI: Colored status icons in terminal

When run in a TTY, `gws list -s` displays:
- `✓` in green (ANSI 32)
- `✗` in red (ANSI 31)
- `↑N` in cyan (ANSI 36)
- `↓N` in magenta (ANSI 35)

All color combinations verified: clean, dirty, ahead, behind, ahead+behind.

## CLI: Colors stripped when piped (auto mode)

```
$ gws list -s | head -5 | cat -v
Found 218 repositories:

NAME                                                     STATUS
-------------------------------------------------------  -----------------------------------------
bgs-feedback-service                                     master                              M-bM-^\M-^S M-bM-^FM-^S37
```

No `^[[` ANSI escape sequences present — only UTF-8 multi-byte sequences shown by `cat -v`.

## CLI: --color=always forces colors when piped

```
$ gws list -s --color=always | head -5 | cat -v
Found 218 repositories:

NAME                                                     STATUS
-------------------------------------------------------  -----------------------------------------
bgs-feedback-service                                     master                              ^[[32mM-bM-^\M-^S^[[0m ^[[35mM-bM-^FM-^S37^[[0m
```

ANSI escape codes (`^[[32m`, `^[[35m`, `^[[0m`) are present when `--color=always` is used, even through a pipe.

## Test Results: Color width parity

```
$ go test ./cmd/git-workspace/... -run TestFormatStatusIconsColorWidthParity -v
=== RUN   TestFormatStatusIconsColorWidthParity
--- PASS: TestFormatStatusIconsColorWidthParity (0.00s)
PASS
```

Colored and uncolored icons have identical display widths, confirming alignment is not affected by ANSI codes.

## Full test suite passes

```
$ go test ./...
ok  github.com/daileyo/gws/cmd/git-workspace    1.594s
ok  github.com/daileyo/gws/internal/classifier  0.649s
ok  github.com/daileyo/gws/internal/config      0.917s
ok  github.com/daileyo/gws/internal/discovery    0.399s
ok  github.com/daileyo/gws/internal/filter       1.431s
ok  github.com/daileyo/gws/internal/git          2.361s
ok  github.com/daileyo/gws/internal/user         2.119s
```
