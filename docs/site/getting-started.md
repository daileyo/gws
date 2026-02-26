# Getting Started

## Installation

### Install via Homebrew

The easiest way to install `git-workspace` on macOS:

```bash
brew tap daileyo/gws
brew install gws
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

## Quick Start

### 1. Initialize a Workspace

Initialize git-workspace by scanning a directory for git repositories:

```bash
# Initialize in current directory
git-workspace --init .

# Initialize in a specific directory
git-workspace --init ~/projects

# Initialize with absolute path
git-workspace --init /path/to/your/workspace
```

This command will:

- Recursively scan the directory for git repositories
- Extract repository metadata (name, path, remote URL)
- Save the configuration to `~/.gws/config.json`

### 2. View Workspace Information

Once initialized, run `git-workspace` to see your workspace information:

```bash
git-workspace
```

Output:

```
Workspace: /home/user/projects
Repositories: 15

Use 'git-workspace --help' to see available commands
```

### 3. List Repositories

View all tracked repositories with their metadata:

```bash
git-workspace --list
```

Output:

```
Found 15 repositories:

NAME           TYPE      VISIBILITY  TAGS              PATH
----           ----      ----------  ----              ----
my-project     github    private     personal, web     /home/user/projects/my-project
work-api       gitlab    private     work              /home/user/projects/work-api
client-site    bitbucket unknown     client, archived  /home/user/projects/client-site
```

### 4. Check Version

```bash
git-workspace --version
```

Output:

```
git-workspace version dev
  commit: abc1234
  built:  2025-12-25T21:00:00Z
```
