# 10 Tasks - Tag Refactor

## Relevant Files

- `cmd/git-workspace/tag.go` - Currently contains `runTag()` (add-tag logic) and `findRepositories()`. Will be heavily refactored: extract parameterized `runAddTag(repo, tag)`, create `tagCmd` with `tagAddCmd`/`tagRemoveCmd` sub-subcommands, register `-a`/`-d` flags, add `ValidArgsFunction` for tab completion.
- `cmd/git-workspace/tag_test.go` - Existing tests for `findRepositories()` and tag management. Will add tests for the new `tagCmd`, sub-subcommands, short flags, mutual exclusivity, and root alias.
- `cmd/git-workspace/untag.go` - Currently contains `runUntag()` (remove-tag logic). Will be refactored to `runRemoveTag(repo, tag)` with parameterized signature. May be consolidated into `tag.go`.
- `cmd/git-workspace/main.go` - Root command. Remove `flagAddTag`/`flagRemoveTag` variables and their flag registrations from `init()`. Remove the tag dispatch block from `RunE`. Remove "Tag Operations" section from usage template. Update `cobra.AddTemplateFunc` calls if tag flags are removed.
- `cmd/git-workspace/main_test.go` - Update or remove tests that reference `flagAddTag`/`flagRemoveTag`. Update mutual exclusivity tests if tag flag counting changes.
- `cmd/git-workspace/deprecated.go` - Add `depAddTag`/`depRemoveTag` hidden flag variables and their registrations. Add deprecation warning mappings. Update `handleDeprecatedFlags()` to count and dispatch deprecated tag flags.
- `cmd/git-workspace/deprecated_test.go` - Add tests for deprecated `--add-tag`/`--remove-tag` flags: hidden status, warning emission, correct routing.

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `tag.go` and `tag_test.go` in the same directory).
- Use `make test` for running tests, `make vet` for vet checks, `make ci` for full CI validation.
- Follow existing Go patterns: `fmt.Errorf("...: %w", err)` for error wrapping, table-driven tests, `t.TempDir()`/`t.Setenv()` for isolation.
- All test files use `package main` (not `_test` package) for access to unexported functions.
- The `user.go` subcommand is the reference pattern for how Cobra subcommands with sub-subcommands are structured in this project.
- The `-t` shorthand is currently used by the deprecated `--tag` filter flag on root (registered in `deprecated.go`). This shorthand must be removed from the deprecated `--tag` registration before `-t` can be assigned to the new tag command root alias. The long-form `--tag` on root continues to work for backward compat.

## Tasks

### [x] 1.0 Extract Tag Business Logic into Parameterized Functions

#### 1.0 Proof Artifact(s)

- Test: `make test` passes with zero changes to external behavior demonstrates refactor is safe
- Test: `make vet` passes demonstrates code quality maintained
- CLI: `gws --add-tag` and `gws --remove-tag` still function identically demonstrates no regressions

#### 1.0 Tasks

- [x] 1.1 Refactor `runTag(_ *cobra.Command, args []string)` in `tag.go` to a parameterized `runAddTag(repoIdentifier, tag string) error` that accepts repo and tag as direct parameters instead of parsing from `args`. Update the error message from `"--add-tag requires exactly 2 arguments"` to a generic `"add-tag requires exactly 2 arguments: <repository> <tag>"`.
- [x] 1.2 Refactor `runUntag(_ *cobra.Command, args []string)` in `untag.go` to a parameterized `runRemoveTag(repoIdentifier, tag string) error` that accepts repo and tag as direct parameters. Update the error message similarly.
- [x] 1.3 Update the call sites in `main.go` RunE: change `runTag(cmd, args)` to `runAddTag(args[0], args[1])` (with the existing arg count validation before the call) and `runUntag(cmd, args)` to `runRemoveTag(args[0], args[1])`. Add arg count validation: if `flagAddTag` or `flagRemoveTag` is set and `len(args) != 2`, return an error.
- [x] 1.4 Verify `make test` and `make vet` pass with the refactored signatures.

### [ ] 2.0 Create Tag Subcommand with Add/Remove Sub-Operations

#### 2.0 Proof Artifact(s)

- CLI: `gws tag add my-repo work` adds the tag demonstrates word-form add works
- CLI: `gws tag -a my-repo work` adds the tag demonstrates short-flag add works
- CLI: `gws tag remove my-repo work` removes the tag demonstrates word-form remove works
- CLI: `gws tag -d my-repo work` removes the tag demonstrates short-flag remove works
- CLI: `gws tag` displays help text demonstrates default behavior
- CLI: `gws -t -a my-repo work` works via root alias demonstrates root alias delegation
- Test: `make ci` passes demonstrates no regressions

#### 2.0 Tasks

- [ ] 2.1 Create the `tagCmd` Cobra command in `tag.go` with `Use: "tag"`, `Short: "Manage repository tags"`. The `RunE` should check if `-a` or `-d` flags are set (see 2.4), and if neither is set and no sub-subcommand was invoked, display help (`cmd.Help()`). Register with `rootCmd.AddCommand(tagCmd)` in an `init()` function.
- [ ] 2.2 Create the `tagAddCmd` sub-subcommand in `tag.go` with `Use: "add <repo> <tag>"`, `Args: cobra.ExactArgs(2)`. Its `RunE` calls `runAddTag(args[0], args[1])`. Register with `tagCmd.AddCommand(tagAddCmd)`.
- [ ] 2.3 Create the `tagRemoveCmd` sub-subcommand in `tag.go` with `Use: "remove <repo> <tag>"`, `Args: cobra.ExactArgs(2)`. Its `RunE` calls `runRemoveTag(args[0], args[1])`. Register with `tagCmd.AddCommand(tagRemoveCmd)`.
- [ ] 2.4 Register `-a` (bool, `tagFlagAdd`) and `-d` (bool, `tagFlagDelete`) flags on `tagCmd`. In `tagCmd.RunE`, if `-a` is set, validate `len(args) == 2` and call `runAddTag(args[0], args[1])`. If `-d` is set, validate `len(args) == 2` and call `runRemoveTag(args[0], args[1])`. Enforce mutual exclusivity: if both `-a` and `-d` are set, return an error.
- [ ] 2.5 Set `tagCmd.Args = cobra.ArbitraryArgs` so that positional args can be passed when `-a` or `-d` flags are used (Cobra needs to accept args on the parent command for the flag-based invocation pattern).
- [ ] 2.6 Remove the `-t` shorthand from the deprecated `--tag` filter flag registration in `deprecated.go`. Change `StringSliceVarP(&filterTags, "tag", "t", ...)` to `StringSliceVar(&filterTags, "tag", ...)` (no shorthand). This frees `-t` for the tag command alias.
- [ ] 2.7 Register `-t` as a hidden boolean flag on the root command that serves as an alias for the `tag` subcommand. In root `RunE`, before the tag dispatch block, check if `-t` alias is set and delegate to `tagCmd` execution by calling `tagCmd.RunE(tagCmd, args)` (passing through remaining args). The `-t` flag variable can be named `flagTagAlias`.
- [ ] 2.8 Write tests in `tag_test.go`: (a) `TestTagSubcommandsRegistered` — verify `add` and `remove` are registered on `tagCmd`, (b) `TestTagFlagAdd` — verify `-a` flag routes to add logic, (c) `TestTagFlagDelete` — verify `-d` flag routes to remove logic, (d) `TestTagMutualExclusivity` — verify `-a` and `-d` together error, (e) `TestTagNoOperationShowsHelp` — verify bare `gws tag` shows help, (f) `TestTagAliasFlag` — verify `-t` is registered and hidden on root.
- [ ] 2.9 Verify `make ci` passes.

### [ ] 3.0 Add Tab Completion for Tag Operations

#### 3.0 Proof Artifact(s)

- CLI: `git-workspace __complete tag add ""` returns repo names demonstrates repo completion
- CLI: `git-workspace __complete tag remove my-repo ""` returns existing tags for that repo demonstrates tag completion
- Test: Completion function tests pass demonstrates correct completion behavior

#### 3.0 Tasks

- [ ] 3.1 Add a `ValidArgsFunction` to `tagAddCmd` that suggests repo names for the first positional arg (load config, return repo names matching the prefix) and no completions for the second arg. Follow the same pattern as `rootCmd.ValidArgsFunction`.
- [ ] 3.2 Add a `ValidArgsFunction` to `tagRemoveCmd` that suggests repo names for the first positional arg, and for the second positional arg, suggests existing tags for the specified repo (load config, find repo by first arg, return its tags).
- [ ] 3.3 Add a `ValidArgsFunction` to `tagCmd` itself for the flag-based invocation pattern (`tag -a <repo> <tag>`) — same logic as `tagAddCmd` (repo names for first arg, no completions for second).
- [ ] 3.4 Write tests that verify the completion functions return expected repo names and tags. Use table-driven tests with mock config.
- [ ] 3.5 Verify `make ci` passes.

### [ ] 4.0 Deprecate Old Tag Flags and Clean Up Root Command

#### 4.0 Proof Artifact(s)

- CLI: `gws --add-tag my-repo work` adds tag plus prints deprecation warning demonstrates backward compat
- CLI: `gws --remove-tag my-repo work` removes tag plus prints deprecation warning demonstrates backward compat
- CLI: `gws --help` does NOT show `--add-tag` or `--remove-tag` demonstrates clean help output
- Test: Deprecated tag flag tests pass demonstrates deprecation layer works
- Test: `make ci` passes demonstrates full suite green

#### 4.0 Tasks

- [ ] 4.1 Remove `flagAddTag` and `flagRemoveTag` variables from `main.go`. Remove their flag registrations (`--add-tag`/`-d` and `--remove-tag`/`-x`) from `main.go`'s `init()`.
- [ ] 4.2 Remove the tag dispatch block from root `RunE` (`if flagAddTag { ... }` and `if flagRemoveTag { ... }`).
- [ ] 4.3 In `deprecated.go`, add `depAddTag` (bool) and `depRemoveTag` (bool) variables. Register them as hidden flags in `registerDeprecatedFlags()`: `--add-tag` (bool, hidden) and `--remove-tag` (bool, hidden). Add them to the hidden flags list.
- [ ] 4.4 Add deprecation warning mappings to `depWarnings`: `"add-tag" → "gws tag add <repo> <tag>"` and `"remove-tag" → "gws tag remove <repo> <tag>"`.
- [ ] 4.5 Update `handleDeprecatedFlags()` to count `depAddTag`/`depRemoveTag` in the mutual exclusivity check (replacing the old `flagAddTag`/`flagRemoveTag` counting). Add dispatch logic: if `depAddTag` is set, validate 2 args, emit warnings, call `runAddTag(args[0], args[1])`; same for `depRemoveTag` calling `runRemoveTag(args[0], args[1])`.
- [ ] 4.6 Remove the "Tag Operations" section from the root usage template in `main.go`. The tag subcommand will appear naturally in the "Available Commands" list.
- [ ] 4.7 Update `deprecated_test.go`: (a) add `"add-tag"` and `"remove-tag"` to `TestDeprecatedFlagsAreHidden`, (b) add a test that `--add-tag` emits a deprecation warning, (c) add a test that deprecated `--add-tag` dispatches to `runAddTag`.
- [ ] 4.8 Update `main_test.go` to remove any references to `flagAddTag`/`flagRemoveTag`. Update the mutual exclusivity test if it references these old variables.
- [ ] 4.9 Verify `make ci` passes with the complete refactor in place.
