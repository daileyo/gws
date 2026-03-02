# 11-spec-tag-path-targeting

## Introduction/Overview

The `gws tag` command currently accepts a single positional identifier that matches repositories by partial name (case-insensitive) or exact path. This spec extends tag targeting with two explicit flags: `--path`/`-p` for matching repositories by path prefix or substring, and `--repo`/`-r` for explicit name-based matching. Both flags work on `tag add`, `tag remove`, and their short-flag equivalents (`tag -a`, `tag -d`). When both flags are combined, AND logic applies — a repo must match both conditions. Tags can be deleted in exactly the same way they are added.

## Goals

- Add `--path`/`-p` flag to match repositories by path prefix or substring (case-sensitive, partial match)
- Add `--repo`/`-r` flag to explicitly target repositories by partial name (case-insensitive), mirroring the current default
- Support combining `--path` and `--repo` with AND logic (repo must satisfy both conditions)
- Preserve existing default behavior: positional `<identifier> <tag>` continues to match by partial name OR exact path
- Display repo paths alongside names in output when `--path` targeting is used so the user can verify what was matched

## User Stories

- **As a developer with many repos in nested directories**, I want to use `gws tag add --path /work/clients backend` so that I can tag every repo under a directory in one command without naming each one.
- **As a power user**, I want to combine `--repo api --path /work` so that only repos named something like "api" AND located under `/work` get tagged — avoiding accidental matches in other directories.
- **As a script author**, I want to use `gws tag add --repo myrepo work` instead of the positional form so that my script's intent is explicit and readable.
- **As a user managing tags**, I want `gws tag remove --path /personal archive` so that I can bulk-remove a tag from all repos in a directory with a single command.

## Demoable Units of Work

### Unit 1: Path-Based Tag Add and Remove

**Purpose:** Allows users to add or remove a tag from all repos whose path starts with or contains a given string, using the `--path`/`-p` flag on `tag add`, `tag remove`, `tag -a`, and `tag -d`.

**Functional Requirements:**
- The system shall add a `--path`/`-p` string flag to the `tag add` subcommand
- The system shall add a `--path`/`-p` string flag to the `tag remove` subcommand
- The system shall add a `--path`/`-p` string flag to the `tag` command (for use with `-a` and `-d`)
- When `--path` is provided, the system shall match repositories whose `path` field starts with the given string (prefix match); if no repos match by prefix, the system shall fall back to matching repos whose `path` field contains the given string (substring match)
- Path matching shall be case-sensitive
- When `--path` is provided, the command shall require exactly 1 positional argument: the tag value (not 2 as in the default mode)
- The system shall return an error if `--path` is provided but the positional tag argument is missing
- The system shall print `"no repositories found matching path: <value>"` and return an error when `--path` is provided but no repos match
- When `--path` is used and up to 5 repos are displayed, the output shall show both the repo name and path on each line (e.g. `  - my-repo  /home/user/my-repo`)
- When `--path` is used and more than 5 repos are matched, only the count summary line is shown (no per-repo lines), consistent with current behavior

**Proof Artifacts:**
- CLI: `gws tag add --path /home/user/work backend` adds the "backend" tag to all repos under `/home/user/work` and prints repo names with paths demonstrates path prefix matching works
- CLI: `gws tag remove --path /home/user/work backend` removes the "backend" tag from matching repos demonstrates path-based remove works
- CLI: `gws tag -a --path /home/user/work backend` produces identical results to `tag add --path` demonstrates short-flag integration
- CLI: `gws tag -d --path /home/user/work backend` produces identical results to `tag remove --path` demonstrates short-flag integration
- CLI: `gws tag add --path /nonexistent backend` returns an error with "no repositories found" demonstrates graceful no-match handling
- Test: `TestFindRepositoriesByPath` table-driven tests cover prefix match, substring fallback, case-sensitivity, and no-match demonstrates path matching logic

### Unit 2: Explicit Repo Flag and AND-Logic Combination

**Purpose:** Adds `--repo`/`-r` for explicit name-based targeting, and supports combining `--repo` and `--path` with AND logic so only repos satisfying both conditions are targeted.

**Functional Requirements:**
- The system shall add a `--repo`/`-r` string flag to the `tag add` subcommand
- The system shall add a `--repo`/`-r` string flag to the `tag remove` subcommand
- The system shall add a `--repo`/`-r` string flag to the `tag` command (for use with `-a` and `-d`)
- When `--repo` is provided, the system shall match repositories whose name contains the given string (partial, case-insensitive), identical to the current default name-matching behavior
- When `--repo` is provided (without `--path`), the command shall require exactly 1 positional argument: the tag value
- When both `--repo` and `--path` are provided, the system shall return only repos that satisfy both conditions (AND logic — repo name matches `--repo` AND path matches `--path`)
- The default positional form `gws tag add <identifier> <tag>` (no flags) shall remain unchanged: matches by partial name OR exact path, requires 2 positional args
- The system shall return an error if `--repo` is provided but the positional tag argument is missing
- The system shall print `"no repositories found matching repo: <value>"` when `--repo` alone finds no matches
- The system shall print `"no repositories found matching repo: <value> and path: <value>"` when the AND combination finds no matches

**Proof Artifacts:**
- CLI: `gws tag add --repo api backend` adds the "backend" tag to all repos with "api" in the name demonstrates explicit repo-flag add
- CLI: `gws tag remove --repo api backend` removes the tag demonstrates explicit repo-flag remove
- CLI: `gws tag add --repo api --path /work backend` only tags repos matching both conditions demonstrates AND logic
- CLI: `gws tag -a --repo api backend` and `gws tag -d --repo api backend` work correctly demonstrates short-flag integration
- CLI: Legacy positional form `gws tag add api backend` still works unchanged demonstrates backward compatibility
- Test: `TestFindRepositoriesWithFilters` covers: repo-only, path-only, and combined AND logic, plus no-match cases demonstrates all filter combinations

## Non-Goals (Out of Scope)

1. **Tag listing or search**: This spec does not add a way to list or query tags by path; that remains `gws list --tag`
2. **OR logic for combined flags**: When both `--repo` and `--path` are given, only AND logic is supported — union matching is not included
3. **Glob or regex patterns**: Path matching is limited to prefix/substring string matching; no wildcards or regex
4. **Changes to `gws list` filtering**: The `list --tag`/`-t` filter behavior is unchanged
5. **Removing all tags from a matched set**: No bulk "remove all tags" operation; the tag value is always required

## Design Considerations

No UI/visual design changes. Output format change is limited to: when `--path` targeting is used and up to 5 repos are listed, the path is shown next to the repo name on each line:

```
Added tag 'backend' to 3 repositories
  - my-api  /home/user/work/my-api
  - other-api  /home/user/work/other-api
  - gateway  /home/user/work/gateway
```

When name-based targeting is used (default or `--repo`), the output remains as today:

```
Added tag 'backend' to 3 repositories
  - my-api
  - other-api
  - gateway
```

## Repository Standards

- Written in Go; follow existing patterns: error wrapping with `fmt.Errorf("...: %w", err)`, table-driven tests in `_test.go` files co-located with source
- Use `cobra.Command` and `pflag` conventions for flag registration, consistent with how `-a` and `-d` are registered on `tagCmd`
- New flags (`--path`/`-p` and `--repo`/`-r`) are string flags (not bool), registered on `tagAddCmd`, `tagRemoveCmd`, and `tagCmd`
- All business logic lives in `cmd/git-workspace/`; no changes to `internal/config/` or other packages needed
- Commit messages follow Conventional Commits (`feat:` for this feature)
- Validate with `make test` and `make lint`; CI runs `make ci`

## Technical Considerations

- **Argument count changes**: When `--path` or `--repo` is provided, commands expect exactly 1 positional arg (the tag). Without flags, 2 positional args are required (identifier + tag). The `RunE` on `tagAddCmd`, `tagRemoveCmd`, and `tagCmd` (for `-a`/`-d`) must branch on whether targeting flags are set.
- **`findRepositories()` extension**: The existing `findRepositories(cfg, identifier)` function handles the default (no-flag) case unchanged. A new internal function — e.g. `findRepositoriesWithFilters(cfg, repoFilter, pathFilter string)` — handles the flag-based targeting. When `pathFilter` is non-empty, apply prefix match first; if that returns no results, fall back to substring match. When `repoFilter` is non-empty, apply partial case-insensitive name match. When both are non-empty, a repo must satisfy both.
- **Flag variable scope**: `--path` and `--repo` flags need package-level variables for the `tag` command (to support `tag -a --path ...`), and local variables for `tagAddCmd` / `tagRemoveCmd`.
- **Shell completion for `--path`**: Implement a `completeRepoPaths` function that reads `cfg.Repositories` and returns the list of known `repo.Path` values as completion candidates. Register it via `tagAddCmd.RegisterFlagCompletionFunc("path", ...)` and `tagRemoveCmd.RegisterFlagCompletionFunc("path", ...)`.
- **No changes to `internal/config`**: The `Repository` struct and `Config` struct are unchanged; all new logic is in command-layer files.

## Security Considerations

No specific security considerations identified. Path values are used only for string comparison against in-memory config data — they are not used in file system operations, shell execution, or external calls.

## Success Metrics

1. **All existing tag tests pass** — no regressions to current behavior
2. **New tests pass** — `TestFindRepositoriesByPath` and `TestFindRepositoriesWithFilters` cover all matching modes
3. **`make ci` passes** — linting and full test suite green
4. **`gws tag add --help` and `gws tag remove --help`** show `--path` and `--repo` flags with descriptions
5. **Path tab completion** suggests known repo paths from config when `--path` is typed

## Open Questions

No open questions at this time.
