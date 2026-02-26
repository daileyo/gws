# 06 Task 2.0 Proofs - Makefile Messaging and README Documentation

## CLI Output

### `make setup-hooks` names both hooks

```
$ make setup-hooks
Setting up git hooks...
git config core.hooksPath .githooks
Git hooks installed:
  pre-push   — runs go vet, golangci-lint, and tests before each push
  commit-msg — pads conventional commit types for aligned git log output
```

Both hooks named with descriptions ✓

### `make help` shows updated setup-hooks description

```
$ make help
Usage: make [target]

Targets:
  setup-hooks   Install git hooks (pre-push linting, commit-msg formatting)
  build         Build the binary
  test          Run tests
  ...
```

`setup-hooks` description mentions both hooks ✓

## README Diff

`README.md` Development section now includes a One-Time Setup subsection
before Build Commands:

```markdown
## Development

### One-Time Setup

Git hooks are provided in `.githooks/` for pre-push checks and commit message
formatting. Activate them once after cloning:

    make setup-hooks

### Build Commands
...
```

`make setup-hooks` documented as a one-time setup step ✓
