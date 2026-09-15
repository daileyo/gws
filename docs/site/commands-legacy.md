# Legacy Flags

Prior to v2.6.0, `omgw` used root-level flags for all operations (e.g., `omgw --list`, `omgw --init`). These flags have been replaced by subcommands but **remain functional** for backward compatibility. Using a deprecated flag will emit a warning pointing you to the new form.

## Migration Reference

### Command Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--list` | `-l` | `omgw list` |
| `--init` | `-i` | `omgw init` |
| `--add` | `-a` | `omgw add [path]` |
| `--recursive` | `-v` | `omgw add --recursive` |
| `--refresh` | `-r` | `omgw refresh` |
| `--print-workspace` | `-w` | `omgw print-workspace` |
| `--go` | `-g` | `omgw <repo-name>` |
| `--add-tag` | `-d` | `omgw tag add <repo> <tag>` |
| `--remove-tag` | `-x` | `omgw tag remove <repo> <tag>` |

### List Filter Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--type` | `-y` | `omgw list --type` |
| `--tag` | | `omgw list --tag` |
| `--name` | `-n` | `omgw list --name` |
| `--path` | `-p` | `omgw list --path` |
| `--output` | `-o` | `omgw list --output` |
| `--status` | `-s` | `omgw list --status` |
| `--show-user` | | `omgw list --show-user` |
| `-V` | | `omgw list -i` (filter) or `omgw list -I` (show visibility column) |

### User Management Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--user` | | `omgw user list` |
| `--list-users` | | `omgw user list` |
| `--update` | `-u` | `omgw user assign <repo> <profile>` |
| `--delete` | `-D` | `omgw user assign` (remove local config) |
| `--all` | | `omgw user assign` (with `--all`) |
| `--verbose` | | `omgw user --verbose` |
| `--git-name` | | `omgw user add --name` |
| `--git-email` | | `omgw user add --email` |

## Example

**Before (deprecated):**

```bash
omgw --list --type github --status
omgw --add-tag my-repo personal
omgw --init ~/projects
omgw --user --update my-repo work
```

**After (current):**

```bash
omgw list --type github -S
omgw tag add my-repo personal
omgw init ~/projects
omgw user assign my-repo work
```

When you use a deprecated flag, you will see a warning like:

```
Warning: --list is deprecated, use 'omgw list' instead
```

The deprecated flags are hidden from `--help` output but will continue to work for the foreseeable future.
