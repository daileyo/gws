# Getting Started

## Installation

### Install via Homebrew

```bash
brew install daileyo/gws/git-workspace
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/daileyo/gws.git
cd gws

# Build the binary
make build

# The binary will be in ./build/git-workspace
# Optionally, install to your PATH
make install
```

## Shell Integration

Before using the `gws` shorthand, set up shell integration. Add these lines to your `~/.zshrc` (or `~/.bashrc`):

```bash
export PATH="$HOME/.local/bin:$PATH"
eval "$(git-workspace shell-init zsh)"   # or: shell-init bash
```

This gives you the `gws` function with tab completion. See [Shell Integration](shell-integration.md) for full details.

## Quick Start

### 1. Initialize a Workspace

Scan a directory for git repositories:

```bash
# Initialize in current directory
gws init

# Initialize in a specific directory
gws init ~/projects

# Initialize with absolute path
gws init /path/to/your/workspace
```

This command will:

- Recursively scan the directory for git repositories, descending through container directories and following workspace symlinks
- Extract repository metadata (name, path, remote URL)
- Detect repository type and user configuration
- Discover git worktrees for every repository it finds
- Save the configuration to `~/.gws/config.json`

Later, when your checkouts change, [`gws refresh`](commands-core.md#refresh-workspace) re-runs exactly the same discovery against the workspace you initialized.

### 2. List Your Repositories

Once initialized, run `gws` with no arguments to see your repositories:

```bash
gws
```

Output (compact multi-column names):

```
Found 15 repositories:

my-project       work-api         client-site      my-api
frontend-app     docs-site        infra-tools      backend-svc
mobile-app       shared-libs      auth-service     data-pipeline
ml-models        cli-tools        test-harness
```

### 3. View Repository Details

Use verbose mode to see more information:

```bash
gws list -v
```

Output:

```
Found 15 repositories:

NAME              TYPE       VISIBILITY  TAGS              PATH
----------------  ---------  ----------  ----------------  ------------------------------------
my-project        github     private     personal, web     /home/user/projects/my-project
work-api          gitlab     private     work              /home/user/projects/work-api
client-site       bitbucket  unknown     client, archived  /home/user/projects/client-site
```

### 4. Check Version

```bash
gws --version
```

Output:

```
git-workspace version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
```

## Next Steps

- [Core Commands](commands-core.md) — Listing, filtering, refreshing, and more
- [User Management](commands-user.md) — Managing git user profiles
- [Tagging](commands-tagging.md) — Organizing repositories with custom tags
- [Shell Integration](shell-integration.md) — Navigation and tab completion
- [Configuration](configuration.md) — Config file structure and fields
