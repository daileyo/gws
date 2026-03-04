# Core Commands

## Repository Classification

### Automatic Classification

When repositories are discovered with `gws init`, they are automatically classified based on their remote URL:

| Repository Type | Detected From |
|----------------|---------------|
| **GitHub** | `github.com` |
| **GitLab** | `gitlab.com`, `gitlab.*` |
| **Azure DevOps (ADO)** | `dev.azure.com`, `visualstudio.com` |
| **Bitbucket** | `bitbucket.org` |
| **Unknown** | Other or no remote URL |

#### Supported Remote URL Patterns

**GitHub:**

- HTTPS: `https://github.com/user/repo.git`
- SSH: `git@github.com:user/repo.git`

**GitLab:**

- HTTPS: `https://gitlab.com/user/repo.git`
- SSH: `git@gitlab.com:user/repo.git`
- Self-hosted: `https://gitlab.company.com/user/repo.git`

**Azure DevOps:**

- HTTPS: `https://dev.azure.com/org/project/_git/repo`
- HTTPS (legacy): `https://org.visualstudio.com/project/_git/repo`
- SSH: `git@ssh.dev.azure.com:v3/org/project/repo`

**Bitbucket:**

- HTTPS: `https://bitbucket.org/user/repo.git`
- SSH: `git@bitbucket.org:user/repo.git`

### Visibility Detection

Repository visibility is inferred from the remote URL protocol:

- **Private**: SSH URLs (`git@...` or `ssh://...`) typically require authentication
- **Unknown**: HTTPS URLs could be public or private (requires authentication check)

---

## List Repositories

```
gws list [flags]
```

List all tracked repositories with optional filtering and display options.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-y` | | Filter by repository type: `github`, `gitlab`, `ado`, `bitbucket` |
| `--tag` | `-t` | | Filter by tag(s); repeat for AND logic |
| `--name` | `-n` | | Filter by repository name (partial match, supports wildcards) |
| `--path` | `-p` | | Filter by repository path (partial match, supports wildcards) |
| `--output` | `-o` | `table` | Output format: `table` or `json` |
| `--status` | `-s` | `false` | Show git status (branch, clean/dirty, ahead/behind) |
| `--show-user` | `-u` | `false` | Show git user info (USER, EMAIL, SIGN columns) |

### Examples

**Basic listing:**

```bash
gws list
```

```
Found 15 repositories:

NAME           TYPE      VISIBILITY  TAGS              PATH
----           ----      ----------  ----              ----
my-project     github    private     personal, web     /home/user/projects/my-project
work-api       gitlab    private     work              /home/user/projects/work-api
client-site    bitbucket unknown     client, archived  /home/user/projects/client-site
```

**With git status:**

```bash
gws list --status
```

```
Found 15 repositories:

NAME           STATUS         TYPE      VISIBILITY  TAGS              PATH
----           ------         ----      ----------  ----              ----
my-project     main ✓         github    private     personal, web     /home/user/projects/my-project
work-api       develop ✗ ↑2   gitlab    private     work              /home/user/projects/work-api
client-site    main ✓ ↓1      bitbucket unknown     client, archived  /home/user/projects/client-site
```

**Status Indicators:**

- `✓` = Clean working tree (no uncommitted changes)
- `✗` = Dirty working tree (uncommitted changes)
- `↑N` = N commits ahead of remote
- `↓N` = N commits behind remote

**With user info:**

```bash
gws list --show-user
```

**Filtering:**

```bash
# Filter by repository type
gws list --type github

# Filter by single tag
gws list --tag personal

# Filter by multiple tags (AND logic — repo must have ALL tags)
gws list --tag work --tag backend

# Filter by repository name (partial match, case-insensitive)
gws list --name project

# Filter by repository path (partial match)
gws list --path /home/user/projects

# Combine multiple filters
gws list --type gitlab --tag work --name api
```

**JSON output:**

```bash
gws list --output json
```

```json
[
  {
    "name": "my-project",
    "path": "/home/user/projects/my-project",
    "remote_url": "https://github.com/user/my-project.git",
    "type": "github",
    "visibility": "unknown",
    "tags": ["personal", "web"]
  }
]
```

---

## Initialize Workspace

```
gws init [directory]
```

Initialize a workspace by scanning a directory for git repositories. Defaults to the current directory if no path is given.

**What it does:**

- Recursively scans the directory for git repositories
- Extracts repository metadata (name, path, remote URL)
- Detects repository type and git user configuration
- Saves the configuration to `~/.gws/config.json`

**Examples:**

```bash
# Initialize in current directory
gws init

# Initialize in a specific directory
gws init ~/projects

# Initialize with absolute path
gws init /path/to/your/workspace
```

---

## Add Repository

```
gws add [path] [flags]
```

Add a single git repository to the workspace. Defaults to the current directory if no path is given.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--recursive` | `-v` | `false` | Recursively add all git repositories found under the path |

**Examples:**

```bash
# Add current directory as a repository
gws add

# Add a specific repository
gws add ~/projects/my-repo

# Recursively add all repos under a directory
gws add ~/projects --recursive
```

---

## Refresh Workspace

```
gws refresh
```

Re-scan the workspace and update repository metadata.

**What it does:**

- Re-scans workspace for new repositories
- Updates remote URLs and classification
- Re-detects git user configuration
- Clears and rebuilds git status cache
- Preserves all custom tags

**When to use:**

- After adding new repositories to your workspace
- When remote URLs have changed
- To force update of cached git status
- After bulk repository operations

**Example output:**

```
Refreshing workspace at: /home/user/projects
Re-scanning for repositories...
Cleared git status cache

Refresh complete!
Total repositories: 15
New repositories found: 2
```

---

## Print Workspace

```
gws print-workspace
```

Print the workspace root path to stdout. Useful for scripting:

```bash
cd "$(gws print-workspace)"
```
