# Commands Reference

## Repository Classification

### Automatic Classification

When repositories are discovered with `git-workspace --init`, they are automatically classified based on their remote URL:

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

## Tagging

### Custom Tags

Add custom tags to organize repositories. Tags are applied to **all repositories** matching the given name:

```bash
# Add a tag to all repos matching "my-project" (partial match)
git-workspace --add-tag my-project personal

# Add a tag to all repos with "api" in the name
git-workspace --add-tag api backend

# Remove a tag from all matching repos
git-workspace --remove-tag my-project personal

# Tags can be anything: personal, work, client, archived, production, etc.
```

**How Tag Matching Works:**

- Matches by **partial name** (case-insensitive): `--add-tag api work` tags "my-api", "api-gateway", etc.
- Tags are applied to **all matching repositories**
- You'll see a summary of how many repos were tagged

**Examples:**

```bash
# Tag all API services as backend
git-workspace --add-tag api backend
# Output: Added tag 'backend' to 3 repositories
```

---

## Listing and Filtering

### List Repositories

View all tracked repositories with their metadata:

```bash
git-workspace --list
```

**Show Git Status:**

```bash
git-workspace --list --status
# or
git-workspace -l -s
```

Output:

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

### Filtering Repositories

Filter repositories by type, tags, name, or path. All filters can be combined:

```bash
# Filter by repository type
git-workspace --list --type github

# Filter by single tag
git-workspace --list --tag personal

# Filter by multiple tags (AND logic - repo must have ALL tags)
git-workspace --list --tag work --tag backend

# Filter by repository name (partial match, case-insensitive)
git-workspace --list --name project

# Filter by repository path (partial match)
git-workspace --list --path /home/user/projects

# Combine multiple filters
git-workspace --list --type gitlab --tag work --name api
```

### Output Formats

Control how repositories are displayed:

```bash
# Default table format
git-workspace --list

# JSON format for scripting/automation
git-workspace --list --output json
git-workspace --list -o json

# JSON with filters
git-workspace --list --type github -o json
```

**JSON Output Example:**

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

## Refresh

### Refresh Workspace

Update repository metadata and clear status cache:

```bash
git-workspace --refresh
```

**What it does:**

- Re-scans workspace for new repositories
- Updates remote URLs and classification
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
