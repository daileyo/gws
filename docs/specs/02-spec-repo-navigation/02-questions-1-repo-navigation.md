# 02 Questions Round 1 - Repository Navigation

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## Context

Since the original questions were drafted, the CLI has been reworked from subcommands to a flag-based pattern (Spec 03). Key changes relevant to navigation:

- All commands are now flags: `gws --list`, `gws --init <dir>`, `gws --add-tag`, etc.
- The tag filter shorthand (`gws personal`) has been removed
- Positional arguments are now available for new functionality
- Shell integration currently only supports workspace root via `gws --print-workspace`

---

## 1. Navigation Command Pattern

How should repository navigation be invoked?

- [ ] (A) Positional argument: `gws my-repo` — bare positional args trigger navigation (since tag shorthand was removed, positional args are available)
- [ ] (B) Flag-based: `gws --go my-repo` or `gws -g my-repo` — consistent with the new flag pattern
- [X] (C) Both: positional arg navigates by default, flag available as explicit alternative
- [ ] (D) Other (describe)

## 2. Multiple Repository Matches

What should happen if the name matches multiple repositories? For example, if you have `my-api` and `my-api-v2`:

- [ ] (A) Exact match only — `gws my-api` only matches "my-api", not "my-api-v2"
- [ ] (B) Exact match preferred, partial match fallback — if exact match exists use it, otherwise show partial matches as a list
- [ ] (C) Partial match with interactive selection — show a list and let user pick
- [ ] (D) Partial match shows list but doesn't navigate — user must be more specific
- [X] (E) list all matches and allow tabbing to select

## 3. No Match Behavior

What should happen when the repository name doesn't match any tracked repositories?

- [X] (A) Show error message with suggestions (similar names using fuzzy matching or partial match)
- [ ] (B) Show error message only, no suggestions
- [ ] (C) Other (describe)

## 4. Shell Integration Approach

Since a CLI tool cannot change the parent shell's directory directly, how would you like this to work?

- [ ] (A) Print the path only — user wraps in shell function: `cd "$(gws my-repo)"` or `cd "$(gws -g my-repo)"`
- [ ] (B) Print a `cd` command that user can eval: `eval "$(gws my-repo)"`
- [X] (C) All of the above — print path by default, document shell wrapper for convenience
- [ ] (D) Other (describe)

## 5. Output When Navigation Succeeds

When a single repository matches, what should be printed?

- [ ] (A) Just the path (clean for shell integration, e.g., `/home/user/projects/my-repo`)
- [ ] (B) Path with brief confirmation to stderr (path to stdout for piping, message to stderr for user feedback)
- [X] (C) Configurable via `--quiet` flag (verbose by default, quiet for scripting)
- [ ] (D) Other (describe)

## 6. Wildcard Support

Should navigation support the same wildcard patterns (`*`, `?`) added to filters in Spec 03?

- [X] (A) Yes — `gws "my-*"` shows all matching repos (consistent with filter behavior)
- [ ] (B) No — navigation should only match by exact or partial name (keep it simple)
- [ ] (C) Other (describe)
