# Shell Integration

## Workspace Navigation

Navigate to your workspace root with shell integration.

**Setup (add to `~/.bashrc` or `~/.zshrc`):**

```bash
# Navigate to workspace root
function cdgws() {
  cd "$(git-workspace --print-workspace)"
}

# Short alias
alias gcd=cdgws
```

**Usage:**

```bash
# Navigate to workspace from anywhere
cdgws
# or
gcd

# You'll be in your workspace root directory
pwd
# Output: /home/user/projects
```

**Print workspace path:**

```bash
git-workspace --print-workspace
# Output: /home/user/projects
```

---

## Repository Navigation

Jump to any tracked repository by name using the `gws` shell function.

**Setup (add to `~/.bashrc` or `~/.zshrc`):**

```bash
# 'gws': navigate to a repo by name, or pass flags through to git-workspace
function gws() {
  if [[ "$1" == -* ]]; then
    git-workspace "$@"
  else
    cd "$(git-workspace -g "$1" -q)"
  fi
}
```

**Usage:**

```bash
# Navigate to a repository by name — changes your directory
gws my-repo

# Pass any flag through to git-workspace as normal
gws --list
gws -l --tag personal --status
gws --refresh
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

# Verbose mode (default) — details to stderr, path to stdout
git-workspace -g my-repo
# stderr: my-repo (github) → /home/user/projects/my-repo
# stdout: /home/user/projects/my-repo
```
