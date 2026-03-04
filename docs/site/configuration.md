# Configuration

git-workspace stores its configuration in `~/.gws/config.json`. The configuration includes:

- **version**: Config file format version
- **workspace**: Root directory of the workspace
- **profiles**: Array of git user profiles for managing identity across repositories
- **repositories**: Array of discovered repositories

## Configuration Example

```json
{
  "version": "1.0.0",
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
      "user_source": "includeif"
    }
  ]
}
```

## Field Reference

### Top-level Fields

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Config file format version (e.g., `"1.0.0"`) |
| `workspace` | string | Absolute path to the root directory of the workspace |
| `profiles` | array | List of git user profile objects (see below) |
| `repositories` | array | List of all discovered repository objects |

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

## Config File Location

The config file is always written to and read from `~/.gws/config.json`. This path is not currently configurable.

To view the raw config at any time:

```bash
cat ~/.gws/config.json
```
