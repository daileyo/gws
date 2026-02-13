# 02 Questions Round 1 - Repository Navigation

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Conflict Resolution: Tag vs Repository Name

Currently `gws personal` lists repos with tag "personal". Your new feature would have `gws test-repo` navigate to that repo. How should gws decide whether an argument is a tag or a repository name?

- [ ] (A) Repository name takes priority - if a repo matches, navigate to it; otherwise treat as tag filter
- [ ] (B) Tag takes priority - if a tag matches, list repos; otherwise try to navigate to repo
- [ ] (C) Use a flag to differentiate (e.g., `gws go test-repo` or `gws -g test-repo` for navigation)
- [ ] (D) Remove tag shorthand entirely - always require `gws list --tag personal`
- [ ] (E) Other (describe)

## 2. Multiple Repository Matches

What should happen if the name matches multiple repositories? For example, if you have `my-api` and `my-api-v2`:

- [ ] (A) Exact match only - `gws my-api` only matches "my-api", not "my-api-v2"
- [ ] (B) Partial match with selection - show a list and let user pick (interactive)
- [ ] (C) Partial match navigates to first result (alphabetical or by config order)
- [ ] (D) Partial match shows list but doesn't navigate - user must be more specific
- [ ] (E) Other (describe)

## 3. No Match Behavior

What should happen when the repository name doesn't match any tracked repositories?

- [ ] (A) Show error message with suggestions (similar names)
- [ ] (B) Show error message only, no suggestions
- [ ] (C) Fall back to treating input as tag filter (current behavior)
- [ ] (D) Other (describe)

## 4. Shell Integration Approach

Since a CLI tool cannot change the parent shell's directory directly, how would you like this to work?

- [ ] (A) Print the path only (like `--print-workspace`) - user wraps in shell function: `cd "$(gws go my-repo)"`
- [ ] (B) Print a `cd` command that user can eval: `eval "$(gws go my-repo)"`
- [ ] (C) Provide a shell function in docs that wraps the navigation command
- [ ] (D) All of the above - print path by default, document shell wrapper for convenience
- [ ] (E) Other (describe)

## 5. Output Format

When navigation is successful, what should be displayed?

- [ ] (A) Just the path (silent navigation, path only for shell integration)
- [ ] (B) Path with brief confirmation message (e.g., "Navigating to: /path/to/repo")
- [ ] (C) Path with repo details (name, type, tags, status)
- [ ] (D) Configurable via flag (e.g., `--quiet` for just path, verbose by default)
- [ ] (E) Other (describe)

## 6. Case Sensitivity

Should repository name matching be case-sensitive?

- [ ] (A) Case-insensitive (recommended) - `gws My-Repo` matches "my-repo"
- [ ] (B) Case-sensitive - exact case match required
- [ ] (C) Case-sensitive by default, with flag for insensitive (e.g., `-i`)
- [ ] (D) Other (describe)
