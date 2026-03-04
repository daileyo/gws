# Task 3.0: Core and Navigation Pages

## index.md — Updated Feature List and Links

The index page was updated to reflect the new command structure and subcommand-based CLI:

### Updated Features List
- Core workspace commands (`init`, `list`, `add`, `refresh`, `print-workspace`)
- User profile management with `gws user` subcommands
- Tag operations with `gws tag` subcommands
- Shell integration for automatic workspace switching
- Tab completion for modern shells

### Documentation Links
The documentation now references organized command sections:
- **Core Commands**: `gws init`, `gws list`, `gws add`, `gws refresh`, `gws print-workspace`
- **User Management**: `gws user list`, `gws user add`, `gws user show`, `gws user remove`, `gws user assign`, `gws user sync`
- **Tagging**: `gws tag add`, `gws tag remove`
- **Legacy Flags**: Migration guide for deprecated `--init`, `--list`, `--tag-add`, `--tag-remove` flags

## getting-started.md — Examples Using Subcommand Syntax

All examples in the Getting Started guide use the new subcommand-based syntax:

### Installation & Basic Setup
```bash
gws init              # Initialize a workspace
gws list              # List repositories
gws add <repo>        # Add a repository
gws refresh           # Refresh workspace
gws print-workspace   # Print workspace path
```

### Shell Integration
The guide references:
- Shell integration setup instructions for automatic workspace switching
- Configuration of shell functions to use `gws` commands
- No mention of deprecated `--init` or `--list` flags

### Key Changes
- All command examples use subcommand format: `gws <command>`
- Removed any references to deprecated flag-based syntax
- New user experience focused on intuitive command hierarchy

## shell-integration.md — Updated Shell Integration Examples

Examples demonstrate the modern subcommand interface:

### Bash/Zsh Examples
```bash
gws list              # List all repositories
gws refresh           # Refresh workspace profile
gws print-workspace   # Get current workspace path
```

### Manual Integration Snippet
The provided manual zsh integration routes the following commands:
```
list
init
add
refresh
print-workspace
tag
user
completion
shell-init
help
```

Each command is properly mapped to the gws executable with appropriate argument passing.

## configuration.md — Profile and Repository Configuration

### User Profiles Section
Added comprehensive user profile management documentation with the following fields:

```json
{
  "name": "John Doe",
  "git_name": "John Doe",
  "email": "john@example.com",
  "signing_key": "1A2B3C4D5E6F7890",
  "sign_commits": true
}
```

**Profile Fields**:
- `name`: User's display name
- `git_name`: Name for git commits
- `email`: Email address for commits
- `signing_key`: GPG signing key ID
- `sign_commits`: Boolean flag to enable commit signing

### Repository Configuration
Repository fields table updated to include:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Repository name |
| `path` | string | Local repository path |
| `url` | string | Repository URL |
| `type` | string | Repository type (github, gitlab, etc.) |
| `user` | string | Associated user profile |
| `email` | string | Email for this repository |
| `signing_enabled` | boolean | Enable GPG signing |
| `user_source` | string | Source of user assignment (manual, inherited) |

### Tags Configuration
Tag documentation updated to use new subcommand syntax:
```bash
gws tag add <repo> <tag>      # Add tag to repository
gws tag remove <repo> <tag>   # Remove tag from repository
```

### Example Configuration
Complete example JSON configuration included with sample profiles and updated repository fields demonstrating all new configuration options.
