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

By default, only repository names are shown in a compact multi-column layout. Use flags to filter results and control which columns are displayed.

### Flag Convention

`gws list` uses a **dual-purpose flag convention**:

- **Lowercase flags** (`-t`, `-y`, `-s`, etc.) = **filter only** — narrow results without adding a column
- **Uppercase flags** (`-T`, `-Y`, `-S`, etc.) = **show column** — display the column, optionally filtering when a value is provided

For example:

- `gws list -t work` — filter by tag "work", no tag column shown
- `gws list -T work` — filter by tag "work" AND show the tags column
- `gws list -T` — show the tags column with no filter applied

### Filter Flags (Lowercase)

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-y` | Filter by type (exact match: `github`, `gitlab`, `ado`, `bitbucket`, `unknown`) |
| `--visibility` | `-i` | Filter by visibility (exact match: `private`, `unknown`) |
| `--tag` | `-t` | Filter by tag (exact match, single value) |
| `--path` | `-p` | Filter by path pattern (partial match) |
| `--status` | `-s` | Show compact status in name column, or filter by status pattern (partial match) |
| `--show-user` | `-u` | Filter by user name (partial match) |
| `--remote` | `-r` | Filter by remote URL pattern (partial match) |
| `--remote-raw` | `-b` | Filter by raw remote URL pattern (partial match) |
| `--name` | `-n` | Filter by repository name (partial match, supports wildcards) |

### Show Column Flags (Uppercase)

These flags show a column in the output. When provided with a value, they simultaneously show the column and filter by that value.

| Flag | Short | Description |
|------|-------|-------------|
| `--show-type` | `-Y` | Show type column, or show and filter by type |
| `--show-visibility` | `-I` | Show visibility column, or show and filter by visibility |
| `--show-tag` | `-T` | Show tags column, or show and filter by tag |
| `--show-path` | `-P` | Show path column, or show and filter by path |
| `--show-status` | `-S` | Show status column, or show and filter by status |
| `--show-user-col` | `-U` | Show user column, or show and filter by user |
| `--show-remote` | `-R` | Show remote column, or show and filter by remote |
| `--show-remote-raw` | `-B` | Show raw remote column, or show and filter by raw remote |

### Display Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `table` | Output format: `table` or `json` |
| `--verbose` | `-v` | | Verbose output: `-v` shows type, visibility, tags, path; `-vv` shows all columns |
| `--workers` | | `8` | Number of concurrent workers for status fetching |
| `--color` | | `auto` | Color output: `auto`, `always`, `never` |

!!! note "Filter behavior"
    - `--type`, `--tag`, and `--visibility` use **exact matching**.
    - `--tag`/`-t` accepts a **single value** (not repeatable).
    - All other filter flags use partial/pattern matching.

### Examples

**Default listing (compact multi-column names):**

```bash
gws list
```

```
Found 15 repositories:

my-project       work-api         client-site      my-api
frontend-app     docs-site        infra-tools      backend-svc
mobile-app       shared-libs      auth-service     data-pipeline
ml-models        cli-tools        test-harness
```

**Verbose listing (stored data columns):**

```bash
gws list -v
```

```
Found 15 repositories:

NAME              TYPE       VISIBILITY  TAGS              PATH
----------------  ---------  ----------  ----------------  ------------------------------------
my-project        github     private     personal, web     /home/user/projects/my-project
work-api          gitlab     private     work              /home/user/projects/work-api
client-site       bitbucket  unknown     client, archived  /home/user/projects/client-site
```

**Show specific columns:**

```bash
gws list -YTSP
```

**Compact status (icons in name column):**

```bash
gws list -s
```

```
Found 15 repositories:

NAME
------------------------------
my-project                   ✓
work-api                ↑2   ✗
client-site          ↓1      ✓
```

Use `-s` with a value to filter by status:

```bash
gws list -s dirty         # Show only dirty repos
gws list -s clean         # Show only clean repos
gws list -s ahead         # Show only repos ahead of remote
```

When both `-s` and `-S` are specified, `-S` wins and the full STATUS column is shown.

**Full git status column:**

```bash
gws list -S
```

```
Found 15 repositories:

NAME              STATUS
----------------  -----------------------------------------
my-project        main                                    ✓
work-api          develop                            ↑2   ✗
client-site       main                          ↓1        ✓
```

**Status Indicators:**

- `✓` = Clean working tree (no uncommitted changes)
- `✗` = Dirty working tree (uncommitted changes)
- `↑N` = N commits ahead of remote
- `↓N` = N commits behind remote

Status icons are displayed in order: behind → ahead → clean/dirty. Status icons are colorized when the terminal supports it (controlled by `--color`).

Branch names longer than 30 characters are truncated with `...`.

A `⚠` indicator in the NAME column means the repository's git user configuration has drifted from the assigned profile.

**Filtering:**

```bash
# Filter by repository type (exact match)
gws list -y github

# Filter by visibility
gws list -i private

# Filter by single tag (exact match)
gws list -t personal

# Filter by repository name (partial match, case-insensitive)
gws list -n project

# Filter by remote URL pattern
gws list -r github.com

# Filter by status pattern
gws list -s dirty

# Combine multiple filters
gws list -y gitlab -t work -n api

# Filter and show column simultaneously
gws list -T work -S
```

**JSON output:**

By default, JSON output only includes the `name` field:

```bash
gws list -o json
```

```json
[
  {
    "name": "my-project"
  }
]
```

Add show-column flags to include additional fields:

```bash
gws list -o json -YTP
```

```json
[
  {
    "name": "my-project",
    "type": "github",
    "tags": ["personal", "web"],
    "path": "/home/user/projects/my-project"
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
| `--recursive` | `-r` | `false` | Recursively add all git repositories found under the path |

!!! note
    When `--recursive` is used, the scan always operates from the **current working directory**, regardless of the path argument.

When adding a repository that lives **outside the workspace root**, a symlink is automatically created inside the workspace directory pointing to the external repository.

**Examples:**

```bash
# Add current directory as a repository
gws add

# Add a specific repository
gws add ~/projects/my-repo

# Recursively add all repos in current directory
gws add -r

# Recursively add all repos under a path
gws add ~/projects -r
```

---

## Refresh Workspace

```
gws refresh
```

Re-scan the workspace and update repository metadata.

**What it does:**

- Re-scans workspace for new repositories
- Removes repositories whose paths are no longer valid
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
Detecting git user configuration...
Cleared git status cache

Refresh complete!
Total repositories: 15
Removed 1 repository (path no longer valid)
Found 2 new repositories
Updated 3 repositories
Repositories with user configuration: 12
```

The conditional lines (Removed, Found, Updated, Repositories with user configuration) only appear when their counts are greater than zero.

---

## Print Workspace

```
gws print-workspace
```

Print the workspace root path to stdout. Useful for scripting:

```bash
cd "$(gws print-workspace)"
```

---

## Parent Navigation

```
gws parent <repo> [flags]
```

Print the parent directory path of a repository. Useful for navigating to the directory that contains a repository.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--quiet` | `-q` | `false` | Suppress verbose output, print only the path |

**Examples:**

```bash
# Print parent directory path
gws parent my-repo

# Navigate to parent directory
cd "$(gws parent my-repo)"
```

See [Shell Integration](shell-integration.md) for shorthand navigation forms (`gws -p my-repo`, `gws my-repo -p`).
