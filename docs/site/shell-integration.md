# Shell Integration

## Setup

Add two lines to your `~/.zshrc` (or `~/.bashrc`):

```bash
export PATH="$HOME/.local/bin:$PATH"
eval "$(git-workspace shell-init zsh)"   # or: shell-init bash
```

The `shell-init` command outputs the `gws` function and tab completion setup directly from the binary, so your shell integration is always in sync with the installed version. No manual updates needed when you upgrade.

---

## Repository Navigation

Once set up, use `gws` to jump to any tracked repository by name:

```bash
# Navigate to a repository by name — changes your directory
gws my-repo

# Use subcommands directly through gws
gws list
gws list --tag personal -S
gws refresh
```

**Wildcard matching:**

```bash
# Wildcards work too (* = zero or more, ? = single character)
gws "api-*"
gws "?rontend"
```

**Multiple matches:**

When multiple repositories match, git-workspace displays a numbered list for selection:

```
Multiple repositories match 'api':

  1) my-api (github) /home/user/projects/my-api
  2) my-api-v2 (github) /home/user/projects/my-api-v2
  3) work-api (gitlab) /home/user/projects/work-api

Select repository [1-3]:
```

When piped (non-TTY), all matching paths are printed without prompting.

**No match suggestions:**

When no repositories match, git-workspace suggests similar names:

```
No repositories found matching 'aip'

Did you mean:
  my-api
  work-api
```

**Using the binary directly:**

```bash
# Print path without changing directory (useful for scripting)
git-workspace -g my-repo -q
# Output: /home/user/projects/my-repo
```

---

## Parent Navigation

Navigate to the **parent directory** of a repository (the directory that contains it):

```bash
# These all navigate to the parent directory of "my-repo"
gws parent my-repo
gws -p my-repo
gws my-repo -p
```

This is useful when you want to work in the directory that contains a repository, rather than inside the repository itself.

**Using the binary directly:**

```bash
# Print parent path without changing directory
git-workspace parent my-repo -q
# Output: /home/user/projects
```

---

## Worktree Navigation

Navigate to git worktrees directly from the shell:

```bash
# Navigate to a worktree by branch name (searches all repos)
gws worktree feat-auth

# Wildcard matching with interactive selection
gws worktree "feat-*"

# Navigate to a worktree within a specific repo
gws my-repo -wt feat-auth

# List all worktrees for a repo and choose interactively
gws my-repo -wt
```

The shell function handles the `cd` automatically — `gws worktree <branch>` and `gws <repo> -wt <branch>` both change your working directory to the matched worktree path.

The `worktree` subcommands that don't navigate (`list`, `align`, `add`) are passed through to the binary without `cd`:

```bash
gws worktree list              # Lists worktrees (no cd)
gws worktree align --dry-run   # Previews alignment (no cd)
gws worktree add my-repo feat  # Creates worktree (no cd)
```

**Using the binary directly:**

```bash
# Print worktree path without changing directory
git-workspace worktree feat-auth -q
# Output: /home/user/projects/my-repo.wt/feat-auth
```

See [Core Commands](commands-core.md#worktree-management) for full worktree command reference.

---

## Tab Completion

Tab completion is set up automatically by `shell-init` — no separate step needed if you followed the setup above.

For other shells:

**fish:**

```bash
git-workspace completion fish > ~/.config/fish/completions/git-workspace.fish
```

**Manual zsh setup (without shell-init):**

```bash
# git-workspace shell integration — do not edit, managed via shell-init
if ! type compdef &>/dev/null; then
  autoload -U compinit && compinit
fi
function gws() {
  local _dest
  if [[ $# -eq 0 ]]; then
    git-workspace
    return
  fi
  case "$1" in
    list|init|add|refresh|print-workspace|tag|user|completion|shell-init|help|__*) git-workspace "$@" ;;
    worktree)
      case "$2" in
        list|align|add|"") git-workspace "$@" ;;
        *)
          _dest="$(git-workspace "$@" -q 2>/dev/tty </dev/tty)"
          [[ -n "$_dest" ]] && cd "$_dest"
          ;;
      esac
      ;;
    -p|--parent|parent)
      _dest="$(git-workspace parent "$2" -q 2>/dev/tty </dev/tty)"
      [[ -n "$_dest" ]] && cd "$_dest"
      ;;
    -*)
      git-workspace "$@"
      ;;
    *)
      if [[ "$2" == "-p" || "$2" == "--parent" ]]; then
        _dest="$(git-workspace parent "$1" -q 2>/dev/tty </dev/tty)"
      elif [[ "$2" == "-wt" ]]; then
        if [[ -n "$3" ]]; then
          _dest="$(git-workspace "$1" --worktree "$3" -q 2>/dev/tty </dev/tty)"
        else
          _dest="$(git-workspace "$1" --worktree -q 2>/dev/tty </dev/tty)"
        fi
        [[ -n "$_dest" ]] && cd "$_dest"
        return
      else
        _dest="$(git-workspace "$1" -q 2>/dev/tty </dev/tty)"
      fi
      [[ -n "$_dest" ]] && cd "$_dest"
      ;;
  esac
}
source <(git-workspace completion zsh)
compdef _git-workspace gws
```

---

## Workspace Navigation

Print the workspace root path (useful for scripting):

```bash
gws print-workspace
# Output: /home/user/projects
```
