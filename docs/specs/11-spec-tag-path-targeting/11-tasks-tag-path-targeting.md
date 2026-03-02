# 11-tasks-tag-path-targeting

## Relevant Files

- `cmd/git-workspace/tag.go` — Core changes: new `findRepositoriesWithFilters` function, new package-level flag vars, flag registration in `init()`, updated `tagCmd.RunE`, updated `tagAddCmd` definition (Use, Args, RunE), new `runAddTagWithFilters` function, new `completeRepoPaths` function, updated Long help strings
- `cmd/git-workspace/untag.go` — Add `runRemoveTagWithFilters` function; update `tagRemoveCmd` RunE references (note: `tagRemoveCmd` is defined in `tag.go` but `runRemoveTag` lives here)
- `cmd/git-workspace/tag_test.go` — Add `TestFindRepositoriesByPath` and `TestFindRepositoriesWithFilters` table-driven tests; update any existing tests that are affected by `Args` changes on `tagAddCmd`/`tagRemoveCmd`

### Notes

- Unit tests are co-located with source files (e.g., `tag.go` and `tag_test.go` in the same directory).
- Run tests with `go test ./cmd/git-workspace/` or target a specific test with `go test ./cmd/git-workspace/ -run TestName`.
- Run the full quality gate with `make ci` before marking tasks complete.
- Follow existing Go patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests using `[]struct{ name string; ... }`.
- `tagAddCmd` and `tagRemoveCmd` currently use `cobra.ExactArgs(2)` — this must be changed to `cobra.ArbitraryArgs` so Cobra doesn't reject 1-arg invocations before `RunE` runs. Move all arg-count validation into `RunE`.

## Tasks

### [x] 1.0 Implement `findRepositoriesWithFilters` and Unit Tests

#### 1.0 Proof Artifact(s)

- Test: `go test ./cmd/git-workspace/ -run TestFindRepositoriesByPath` passes with all cases green demonstrates path prefix/substring/case-sensitivity logic is correct
- Test: `go test ./cmd/git-workspace/ -run TestFindRepositoriesWithFilters` passes with all cases green demonstrates repo-only, path-only, and AND-logic combinations work
- Test: `go test ./cmd/git-workspace/ -run TestFindRepositories` passes unchanged demonstrates no regression to existing name/path identifier matching

#### 1.0 Tasks

- [ ] 1.1 In `tag.go`, directly below the existing `findRepositories` function, add a new function `findRepositoriesWithFilters(cfg *config.Config, repoFilter, pathFilter string) []*config.Repository`. The function should work as follows:
  - Start with all repos in `cfg.Repositories` as candidates
  - If `repoFilter` is non-empty, keep only repos where `strings.Contains(strings.ToLower(repo.Name), strings.ToLower(repoFilter))` is true (partial, case-insensitive name match)
  - If `pathFilter` is non-empty, first attempt a prefix match: keep repos where `strings.HasPrefix(repo.Path, pathFilter)` is true (case-sensitive). If that produces zero results, fall back to substring match: keep repos where `strings.Contains(repo.Path, pathFilter)` is true (case-sensitive)
  - When both `repoFilter` and `pathFilter` are non-empty, the repo must satisfy both conditions (AND logic — apply name filter first, then path filter to the remaining set)
  - Return the matched repos as `[]*config.Repository` (pointers into the slice, same pattern as `findRepositories`)

- [ ] 1.2 In `tag_test.go`, add a new table-driven test `TestFindRepositoriesByPath` that tests path matching in isolation (call `findRepositoriesWithFilters` with `repoFilter = ""`). Use this shared test config:
  ```
  repos: [
    { Name: "api-service",  Path: "/home/user/work/api-service" }
    { Name: "api-gateway",  Path: "/home/user/work/api-gateway" }
    { Name: "personal-app", Path: "/home/user/personal/personal-app" }
    { Name: "Other-Repo",   Path: "/Home/User/Other-Repo" }
  ]
  ```
  Include these cases:
  - Prefix match: `pathFilter = "/home/user/work"` → returns `api-service` and `api-gateway`
  - Prefix match single: `pathFilter = "/home/user/work/api-service"` → returns `api-service` only
  - Substring fallback: `pathFilter = "user/personal"` (not a prefix of any path, but a substring) → returns `personal-app`
  - Case-sensitive — no match: `pathFilter = "/HOME/USER/WORK"` → returns empty (paths are case-sensitive)
  - Exact case match for casing test: `pathFilter = "/Home/User/Other-Repo"` → returns `Other-Repo`
  - No match: `pathFilter = "/nonexistent"` → returns empty

- [ ] 1.3 In `tag_test.go`, add a new table-driven test `TestFindRepositoriesWithFilters` that tests all filter combinations. Use the same shared config from 1.2. Include these cases:
  - Repo filter only: `repoFilter = "api"` → returns both `api-service` and `api-gateway`
  - Path filter only: `pathFilter = "/home/user/work"` → returns `api-service` and `api-gateway`
  - AND logic — intersection: `repoFilter = "api"`, `pathFilter = "/home/user/work/api-service"` → returns `api-service` only
  - AND logic — no intersection: `repoFilter = "personal"`, `pathFilter = "/home/user/work"` → returns empty (personal-app is not under /work)
  - Both empty: `repoFilter = ""`, `pathFilter = ""` → returns all 4 repos

- [ ] 1.4 Run `go test ./cmd/git-workspace/ -run TestFindRepositories` (covers all three `TestFindRepositories*` functions) and confirm all pass with zero failures before moving to Task 2.0

---

### [x] 2.0 Register Flags and Update Command `RunE` Logic

#### 2.0 Proof Artifact(s)

- CLI: `gws tag add --help` shows `--path`/`-p` and `--repo`/`-r` flags with their descriptions demonstrates flags are registered on `tagAddCmd`
- CLI: `gws tag remove --help` shows `--path`/`-p` and `--repo`/`-r` flags demonstrates flags are registered on `tagRemoveCmd`
- CLI: `gws tag --help` shows `--path`/`-p` and `--repo`/`-r` flags demonstrates flags are registered on `tagCmd` (for `-a`/`-d` usage)
- CLI: `gws tag add --path /home/user/work backend` reports repos tagged and shows no error demonstrates path-flag add integration
- CLI: `gws tag remove --path /home/user/work backend` reports repos untagged demonstrates path-flag remove integration
- CLI: `gws tag -a --path /home/user/work backend` produces the same result as `tag add --path` demonstrates short-flag path integration
- CLI: `gws tag -d --path /home/user/work backend` produces the same result as `tag remove --path` demonstrates short-flag remove integration
- CLI: `gws tag add --repo api backend` tags repos by name using the explicit flag demonstrates repo-flag add integration
- CLI: `gws tag add --repo api --path /work backend` only tags repos matching both demonstrates AND logic
- CLI: `gws tag add api backend` (legacy positional form) still works and produces the same output as before demonstrates backward compatibility
- CLI: `gws tag add --path /nonexistent backend` prints an error containing `"no repositories found matching path"` demonstrates no-match error for path flag
- CLI: `gws tag add --repo ghost backend` prints an error containing `"no repositories found matching repo"` demonstrates no-match error for repo flag
- CLI: `gws tag add --repo ghost --path /nonexistent backend` prints an error containing both `"repo"` and `"path"` demonstrates AND no-match error message

#### 2.0 Tasks

- [ ] 2.1 In `tag.go`, add two new package-level string variables for `tagCmd`'s flags, grouped with the existing bool vars:
  ```go
  var (
      tagFlagAdd    bool
      tagFlagDelete bool
      tagFlagPath   string // --path/-p on tagCmd (for use with -a/-d)
      tagFlagRepo   string // --repo/-r on tagCmd (for use with -a/-d)
  )
  ```
  Also add four package-level string variables for `tagAddCmd` and `tagRemoveCmd`:
  ```go
  var (
      tagAddPath    string // --path/-p on tagAddCmd
      tagAddRepo    string // --repo/-r on tagAddCmd
      tagRemovePath string // --path/-p on tagRemoveCmd
      tagRemoveRepo string // --repo/-r on tagRemoveCmd
  )
  ```

- [ ] 2.2 In `tag.go`'s `init()` function, register `--path`/`-p` and `--repo`/`-r` as string flags on `tagCmd`, `tagAddCmd`, and `tagRemoveCmd`. Add these lines after the existing `-a` and `-d` registrations:
  ```go
  tagCmd.Flags().StringVarP(&tagFlagPath, "path", "p", "", "Match repositories by path prefix or substring (case-sensitive)")
  tagCmd.Flags().StringVarP(&tagFlagRepo, "repo", "r", "", "Match repositories by name (partial, case-insensitive)")
  tagAddCmd.Flags().StringVarP(&tagAddPath, "path", "p", "", "Match repositories by path prefix or substring (case-sensitive)")
  tagAddCmd.Flags().StringVarP(&tagAddRepo, "repo", "r", "", "Match repositories by name (partial, case-insensitive)")
  tagRemoveCmd.Flags().StringVarP(&tagRemovePath, "path", "p", "", "Match repositories by path prefix or substring (case-sensitive)")
  tagRemoveCmd.Flags().StringVarP(&tagRemoveRepo, "repo", "r", "", "Match repositories by name (partial, case-insensitive)")
  ```

- [ ] 2.3 In `tag.go`, change `tagAddCmd.Args` from `cobra.ExactArgs(2)` to `cobra.ArbitraryArgs`, and replace its `RunE` body with branching logic:
  - If `tagAddPath != ""` or `tagAddRepo != ""` (flag mode): require exactly 1 positional arg (the tag value); if `len(args) != 1`, return `fmt.Errorf("tag add with --path or --repo requires exactly 1 argument: <tag>")`; otherwise call `runAddTagWithFilters(tagAddRepo, tagAddPath, args[0])`
  - If neither flag is set (default mode): require exactly 2 positional args; if `len(args) != 2`, return `fmt.Errorf("tag add requires exactly 2 arguments: <repo> <tag>")`; otherwise call `runAddTag(args[0], args[1])`

- [ ] 2.4 In `tag.go`, change `tagRemoveCmd.Args` from `cobra.ExactArgs(2)` to `cobra.ArbitraryArgs`, and replace its `RunE` body with the same branching logic as 2.3 but for remove:
  - Flag mode: call `runRemoveTagWithFilters(tagRemoveRepo, tagRemovePath, args[0])`
  - Default mode: call `runRemoveTag(args[0], args[1])`
  - Error messages mirror those in 2.3 but say `"tag remove"` instead of `"tag add"`

- [ ] 2.5 In `tag.go`, update `tagCmd.RunE` to handle `tagFlagPath` and `tagFlagRepo` when `-a` or `-d` is set. Inside the `if tagFlagAdd` block, add a branch: if `tagFlagPath != ""` or `tagFlagRepo != ""`, require `len(args) == 1` and call `runAddTagWithFilters(tagFlagRepo, tagFlagPath, args[0])`; otherwise keep the existing `len(args) == 2` check and call `runAddTag(args[0], args[1])`. Do the same inside the `if tagFlagDelete` block, calling `runRemoveTagWithFilters` or `runRemoveTag` respectively. Error message for wrong arg count in flag mode: `"tag -a with --path or --repo requires exactly 1 argument: <tag>"` (and same for `-d`).

- [ ] 2.6 In `tag.go`, add a new function `runAddTagWithFilters(repoFilter, pathFilter, tag string) error`. This function should follow the same structure as `runAddTag`, but:
  - Call `findRepositoriesWithFilters(cfg, repoFilter, pathFilter)` instead of `findRepositories`
  - Use a context-aware no-match error message:
    - If both filters are non-empty: `fmt.Errorf("no repositories found matching repo: %s and path: %s", repoFilter, pathFilter)`
    - If only `pathFilter` is non-empty: `fmt.Errorf("no repositories found matching path: %s", pathFilter)`
    - If only `repoFilter` is non-empty: `fmt.Errorf("no repositories found matching repo: %s", repoFilter)`
  - Pass `pathFilter` through to the reporting section (needed in Task 3.0 for output formatting — for now, report the same as `runAddTag`)

- [ ] 2.7 In `untag.go`, add a new function `runRemoveTagWithFilters(repoFilter, pathFilter, tag string) error`. Follow the same structure as `runRemoveTag`, with the same no-match error message logic as described in 2.6. Pass `pathFilter` through to the reporting section.

- [ ] 2.8 Run `go test ./cmd/git-workspace/` and confirm all tests pass. Pay attention to `TestTagFlagAdd`, `TestTagFlagDelete`, and `TestTagSubcommandsRegistered` — the `Args` change may require updating any test that directly validates argument count error messages.

---

### [x] 3.0 Path Output Formatting, Shell Completion, Help Text, and CI

#### 3.0 Proof Artifact(s)

- CLI: `gws tag add --path /home/user/work backend` output shows per-repo lines as `  - my-repo  /home/user/work/my-repo` (name and path on the same line) demonstrates path-aware output format
- CLI: `gws tag add my-repo backend` output shows per-repo lines as `  - my-repo` (name only) demonstrates default output format is unchanged
- CLI: `gws tag remove --path /home/user/work backend` output shows `  - name  /path` format demonstrates remove also uses path-aware format
- CLI: Tab-completing `gws tag add --path <TAB>` offers known repo paths from the gws config demonstrates shell completion works for `--path`
- CLI: `gws tag add --help` includes an example using `--path` and an example using `--repo` in the `Long` description demonstrates updated help text
- CLI: `gws tag remove --help` includes `--path` and `--repo` examples demonstrates remove help is also updated
- Test: `go test ./cmd/git-workspace/` passes with all tests green demonstrates no regressions from formatting changes
- CI: `make ci` completes with exit code 0 and no lint warnings demonstrates full quality gate passes

#### 3.0 Tasks

- [ ] 3.1 In `tag.go`, update `runAddTagWithFilters` to pass `pathFilter` into the reporting section. In the per-repo listing block (the `if taggedCount <= 5` section), change the line `fmt.Printf("  - %s\n", repo.Name)` to: if `pathFilter != ""`, print `fmt.Printf("  - %s  %s\n", repo.Name, repo.Path)`; otherwise print `fmt.Printf("  - %s\n", repo.Name)`. Leave `runAddTag`'s output block completely unchanged.

- [ ] 3.2 In `untag.go`, update `runRemoveTagWithFilters` with the same output change: when `pathFilter != ""`, print `  - name  path`; otherwise print `  - name` only.

- [ ] 3.3 In `tag.go`, add a new function `completeRepoPaths(toComplete string) ([]string, cobra.ShellCompDirective)`. It should load the config, iterate `cfg.Repositories`, and return paths where `strings.HasPrefix(repo.Path, toComplete)` is true. If config fails to load, return `nil, cobra.ShellCompDirectiveNoFileComp`. Follow the same pattern as `completeRepoNames`.

- [ ] 3.4 In `tag.go`'s `init()` function, register the `completeRepoPaths` function as the completion handler for the `path` flag on `tagAddCmd` and `tagRemoveCmd`. Add after the flag registrations from Task 2.2:
  ```go
  _ = tagAddCmd.RegisterFlagCompletionFunc("path", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
      return completeRepoPaths(toComplete)
  })
  _ = tagRemoveCmd.RegisterFlagCompletionFunc("path", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
      return completeRepoPaths(toComplete)
  })
  ```

- [ ] 3.5 Update the `Long` help strings and `Use` field on `tagCmd`, `tagAddCmd`, and `tagRemoveCmd` in `tag.go` to document the new flags. Replace the existing `Long` block on `tagCmd` with one that includes examples for `--path` and `--repo`:
  ```
  gws tag add --path /home/user/work backend   # Add tag to all repos under a path
  gws tag add --repo api backend               # Add tag to repos matching "api" by name
  gws tag add --repo api --path /work backend  # Add tag to repos matching both conditions
  gws tag remove --path /home/user/work backend
  ```
  Update `tagAddCmd.Long` and `tagRemoveCmd.Long` similarly. Also update `tagAddCmd.Use` to `"add [--path <path>] [--repo <repo>] <tag> | add <repo> <tag>"` to reflect that the positional form is flexible.

- [ ] 3.6 Run `go test ./cmd/git-workspace/` to confirm all tests pass.

- [ ] 3.7 Run `make ci` and confirm it completes with exit code 0 and no lint warnings. Fix any issues before marking this task done.
