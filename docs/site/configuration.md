# Configuration

git-workspace stores its configuration in `~/.gws/config.json`. The configuration includes:

- **version**: Config file format version
- **workspace**: Root directory of the workspace
- **profiles**: Array of git user profiles for managing identity across repositories
- **repositories**: Array of discovered repositories
- **preferences**: Optional preferences for controlling CLI behavior

## Configuration Example

```json
{
  "version": "1.2.0",
  "workspace": "/home/user/projects",
  "profiles": [
    {
      "name": "personal",
      "git_name": "Jane Doe",
      "email": "jane@personal.dev",
      "signing_key": "ABC123DEF456",
      "sign_commits": true
    },
    {
      "name": "work",
      "git_name": "Jane Doe",
      "email": "jane.doe@company.com"
    }
  ],
  "repositories": [
    {
      "name": "gws",
      "path": "/home/user/projects/gws",
      "remote_url": "https://github.com/daileyo/gws.git",
      "type": "github",
      "visibility": "unknown",
      "tags": ["personal", "go"],
      "user": "Jane Doe",
      "email": "jane@personal.dev",
      "signing_enabled": true,
      "user_source": "local"
    },
    {
      "name": "my-api",
      "path": "/home/user/projects/my-api",
      "remote_url": "git@gitlab.com:user/my-api.git",
      "type": "gitlab",
      "visibility": "private",
      "tags": ["work", "backend"],
      "user": "Jane Doe",
      "email": "jane.doe@company.com",
      "user_source": "includeif",
      "worktrees": [
        {
          "path": "/home/user/projects/my-api.wt/feat-auth",
          "branch": "feat-auth",
          "aligned": true
        }
      ]
    }
  ],
  "preferences": {
    "status_workers": 8,
    "scan_max_depth": 6
  }
}
```

## Field Reference

### Top-level Fields

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Config file format version (currently `"1.2.0"`). Stamped on every write; no migration is required from earlier versions. |
| `workspace` | string | Absolute path to the root directory of the workspace |
| `profiles` | array | List of git user profile objects (see below) |
| `repositories` | array | List of all discovered repository objects |
| `preferences` | object | Optional preferences for controlling CLI behavior (see below) |

### Profile Fields

Profiles are managed via `gws user add`, `gws user remove`, and related commands. See [User Management](commands-user.md) for details.

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Profile identifier (e.g., `"work"`, `"personal"`) |
| `git_name` | string | Git `user.name` value |
| `email` | string | Git `user.email` value |
| `signing_key` | string | GPG signing key ID (optional) |
| `sign_commits` | boolean | Whether to enable `commit.gpgsign` (optional, defaults to `false`) |

### Repository Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Repository name (derived from the directory name) |
| `path` | string | Absolute path to the repository on disk |
| `remote_url` | string | Git remote URL (empty string if no remote is configured) |
| `type` | string | Detected hosting provider: `github`, `gitlab`, `ado`, `bitbucket`, or `unknown` |
| `visibility` | string | Inferred visibility: `private` (SSH URL) or `unknown` (HTTPS URL) |
| `tags` | array | List of custom tag strings managed via `gws tag add` / `gws tag remove` |
| `user` | string | Git `user.name` configured for this repository |
| `email` | string | Git `user.email` configured for this repository |
| `signing_enabled` | boolean | Whether commit signing is configured for this repository |
| `user_source` | string | Where the user config comes from: `global`, `local`, `includeif`, or `unknown` |
| `worktrees` | array | List of git worktrees for this repository (omitted when empty). See Worktree Fields below. |

### Worktree Fields

Each entry in the `worktrees` array represents a git worktree associated with a repository. Worktree data is populated during `gws refresh` and updated by `gws worktree add` and `gws worktree align`.

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Absolute path to the worktree directory on disk |
| `branch` | string | Branch checked out in this worktree |
| `aligned` | boolean | Whether the worktree is inside the `<repo>.wt/` directory convention |

### Preferences Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `status_workers` | integer | `8` | Number of concurrent workers for fetching git status. Also configurable per-invocation with `gws list --workers`. |
| `scan_max_depth` | integer | `6` | How many directory levels below the workspace root repository discovery traverses. |

#### `scan_max_depth`

Depth is counted in directory levels below the workspace root: a repository at `~/gws/service-api` is depth 1, and one at `~/gws/org-a/team-one/service-api` is depth 3. A repository found at exactly `scan_max_depth` is registered; nothing below it is traversed.

The default of `6` comfortably covers host / organization / team / repository layouts. Raise it if you group repositories more deeply than that and `gws refresh` is not finding them:

```json
{
  "preferences": {
    "scan_max_depth": 8
  }
}
```

Lowering it speeds up scanning of very large trees at the risk of missing deeply nested repositories. Values of zero or below are ignored and the default applies.

Note that traversal also stops at every repository boundary regardless of depth, so the limit only applies to the container directories above your repositories. See [Discovery Rules](commands-core.md#discovery-rules) for the full set of traversal rules.

## Config File Location

The config file is always written to and read from `~/.gws/config.json`. This path is not currently configurable.

To view the raw config at any time:

```bash
cat ~/.gws/config.json
```
