# Configuration

git-workspace stores its configuration in `~/.gws/config.json`. The configuration includes:

- **version**: Config file format version
- **workspace**: Root directory of the workspace
- **repositories**: Array of discovered repositories

## Configuration Example

```json
{
  "version": "1.0.0",
  "workspace": "/home/user/projects",
  "repositories": [
    {
      "name": "gws",
      "path": "/home/user/projects/gws",
      "remote_url": "https://github.com/daileyo/gws.git",
      "type": "github",
      "visibility": "unknown",
      "tags": ["personal", "go"]
    },
    {
      "name": "my-api",
      "path": "/home/user/projects/my-api",
      "remote_url": "git@gitlab.com:user/my-api.git",
      "type": "gitlab",
      "visibility": "private",
      "tags": ["work", "backend"]
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
| `repositories` | array | List of all discovered repository objects |

### Repository Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Repository name (derived from the directory name) |
| `path` | string | Absolute path to the repository on disk |
| `remote_url` | string | Git remote URL (empty string if no remote is configured) |
| `type` | string | Detected hosting provider: `github`, `gitlab`, `ado`, `bitbucket`, or `unknown` |
| `visibility` | string | Inferred visibility: `private` (SSH URL) or `unknown` (HTTPS URL) |
| `tags` | array | List of custom tag strings added via `--add-tag` |

## Config File Location

The config file is always written to and read from `~/.gws/config.json`. This path is not currently configurable.

To view the raw config at any time:

```bash
cat ~/.gws/config.json
```
