# Task 4.0: Command Reference Pages

## commands-core.md — Core Commands Documentation

### Repository Classification
Comprehensive documentation of repository type detection:

**Type Detection Table**:
| Type | URL Pattern | Detection |
|------|-------------|-----------|
| GitHub | `github.com/...` | URL contains `github.com` |
| GitLab | `gitlab.com/...` | URL contains `gitlab.com` |
| Gitea | Self-hosted | URL pattern matching |
| Other | Custom | User-specified |

### gws list Command
Complete documentation with all 7 supported flags:

```bash
gws list [OPTIONS]
```

**Flags**:
- `--type, -y <type>`: Filter by repository type (github, gitlab, gitea, other)
- `--tag, -t <tag>`: Filter by tag (supports partial matching)
- `--name, -n <name>`: Filter by repository name
- `--path, -p <path>`: Filter by repository path
- `--output, -o <format>`: Output format (json, table, csv)
- `--status, -s <status>`: Filter by repository status (clean, dirty, unknown)
- `--show-user, -u`: Include user profile information in output

**Examples**:
```bash
gws list --type github
gws list --tag personal --output json
gws list --name kubernetes --show-user
gws list --status dirty
```

### Core Commands
- `gws init`: Initialize a workspace
- `gws add`: Add a single repository
- `gws add --recursive`: Add repositories recursively
- `gws refresh`: Refresh workspace state
- `gws print-workspace`: Print current workspace path

## commands-user.md — User Management Commands

Complete documentation of all 6 user subcommands:

### user list
```bash
gws user list [OPTIONS]
```
Lists all configured user profiles with optional filtering.

### user add
```bash
gws user add <name> [OPTIONS]
```

**Options**:
- `--email <email>`: User's email address
- `--name <name>`: User's display name (git commits)
- `--signing-key <key>`: GPG signing key ID
- `--sign-commits`: Enable commit signing

**Example**:
```bash
gws user add john --email john@example.com --name "John Doe" --signing-key 1A2B3C4D --sign-commits
```

### user show
```bash
gws user show <name>
```
Display detailed information for a specific user profile.

### user remove
```bash
gws user remove <name> [OPTIONS]
```
Remove a user profile from the configuration.

### user assign
```bash
gws user assign <repository> <user> [OPTIONS]
```

**Options**:
- `--use-subdirs`: Automatically assign based on repository subdirectory structure
- `--dry-run`: Preview assignment without applying changes

**Example**:
```bash
gws user assign my-repo john --dry-run
```

### user sync
```bash
gws user sync
```
Synchronize user configurations across repositories.

## commands-tagging.md — Tag Operations

Complete documentation of tag subcommands:

### gws tag add
```bash
gws tag add <repository> <tag>
```
Add a tag to a repository. Supports partial repository name matching.

**Features**:
- Partial name matching for easy tagging
- Tab completion support
- Validates repository existence before tagging

**Examples**:
```bash
gws tag add my-repo personal
gws tag add kubernetes prod
gws tag add project- archive
```

### gws tag remove
```bash
gws tag remove <repository> <tag>
```
Remove a tag from a repository. Supports partial repository name matching.

**Features**:
- Partial name matching for easy removal
- Tab completion support
- Confirms before removing tag

### Tag Filtering
Tags can be used with the `gws list` command for filtering:

```bash
gws list --tag personal        # List all repositories tagged 'personal'
gws list --tag work --output json
```

## commands-legacy.md — Legacy Flags Migration Guide

Comprehensive migration table with all 16 deprecated flags and their modern replacements:

### Migration Reference

| Deprecated Flag | Replacement | Notes |
|---|---|---|
| `gws --init` | `gws init` | Initialize workspace |
| `gws --list` | `gws list` | List repositories |
| `gws --add <repo>` | `gws add <repo>` | Add repository |
| `gws --refresh` | `gws refresh` | Refresh workspace |
| `gws --print-workspace` | `gws print-workspace` | Print workspace path |
| `gws --tag-add <repo> <tag>` | `gws tag add <repo> <tag>` | Add tag |
| `gws --tag-remove <repo> <tag>` | `gws tag remove <repo> <tag>` | Remove tag |
| `gws --list-tags <repo>` | `gws tag list <repo>` | List tags |
| `gws --user-add <name>` | `gws user add <name>` | Add user |
| `gws --user-list` | `gws user list` | List users |
| `gws --user-remove <name>` | `gws user remove <name>` | Remove user |
| `gws --user-show <name>` | `gws user show <name>` | Show user |
| `gws --user-assign <repo> <user>` | `gws user assign <repo> <user>` | Assign user |
| `gws --user-sync` | `gws user sync` | Sync users |
| `gws --completion <shell>` | `gws completion <shell>` | Shell completion |
| `gws --shell-init <shell>` | `gws shell-init <shell>` | Shell initialization |

### Before/After Examples

**Listing repositories (old flag syntax)**:
```bash
gws --list --type github --tag personal
```

**Listing repositories (new subcommand syntax)**:
```bash
gws list --type github --tag personal
```

**Adding a tag (old flag syntax)**:
```bash
gws --tag-add my-repo personal
```

**Adding a tag (new subcommand syntax)**:
```bash
gws tag add my-repo personal
```

### Deprecation Warning Format
When using deprecated flags, users receive a clear deprecation message:

```
WARN: The '--list' flag is deprecated and will be removed in version 3.0
Use 'gws list' instead
```

The warning clearly indicates:
- Which flag is deprecated
- When it will be removed
- The recommended replacement command

## File Reorganization

### Changes Made
- **Old commands.md**: Deleted (replaced by four organized command reference pages)
- **commands-core.md**: Created
- **commands-user.md**: Created
- **commands-tagging.md**: Created
- **commands-legacy.md**: Created

### Navigation Structure
The new structure provides better organization and discoverability:
- Core commands for everyday usage
- User management commands grouped together
- Tagging operations clearly documented
- Legacy flags with migration guide for existing users

## Build Verification

The documentation site was successfully built with the new command reference pages:

- Command executed: `python3 -m mkdocs build`
- Status: **Completed successfully**
- Errors: None
- Warnings: None
- Site generation: All pages rendered correctly
- Navigation: All command reference pages properly indexed

The build confirms that all markdown files are valid, all internal links are correct, and the site is ready for deployment.
