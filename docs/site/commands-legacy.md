# Legacy Flags

Prior to v2.6.0, `gws` used root-level flags for all operations (e.g., `gws --list`, `gws --init`). These flags have been replaced by subcommands but **remain functional** for backward compatibility. Using a deprecated flag will emit a warning pointing you to the new form.

## Migration Reference

### Command Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--list` | `-l` | `gws list` |
| `--init` | `-i` | `gws init` |
| `--add` | `-a` | `gws add [path]` |
| `--recursive` | `-v` | `gws add --recursive` |
| `--refresh` | `-r` | `gws refresh` |
| `--print-workspace` | `-w` | `gws print-workspace` |
| `--go` | `-g` | `gws <repo-name>` |
| `--add-tag` | `-d` | `gws tag add <repo> <tag>` |
| `--remove-tag` | `-x` | `gws tag remove <repo> <tag>` |

### List Filter Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--type` | `-y` | `gws list --type` |
| `--tag` | | `gws list --tag` |
| `--name` | `-n` | `gws list --name` |
| `--path` | `-p` | `gws list --path` |
| `--output` | `-o` | `gws list --output` |
| `--status` | `-s` | `gws list --status` |
| `--show-user` | | `gws list --show-user` |
| `-V` | | `gws list -i` (filter) or `gws list -I` (show visibility column) |

### User Management Flags

| Deprecated Flag | Short | Replacement |
|----------------|-------|-------------|
| `--user` | | `gws user list` |
| `--list-users` | | `gws user list` |
| `--update` | `-u` | `gws user assign <repo> <profile>` |
| `--delete` | `-D` | `gws user assign` (remove local config) |
| `--all` | | `gws user assign` (with `--all`) |
| `--verbose` | | `gws user --verbose` |
| `--git-name` | | `gws user add --name` |
| `--git-email` | | `gws user add --email` |

## Example

**Before (deprecated):**

```bash
gws --list --type github --status
gws --add-tag my-repo personal
gws --init ~/projects
gws --user --update my-repo work
```

**After (current):**

```bash
gws list --type github -S
gws tag add my-repo personal
gws init ~/projects
gws user assign my-repo work
```

When you use a deprecated flag, you will see a warning like:

```
Warning: --list is deprecated, use 'gws list' instead
```

The deprecated flags are hidden from `--help` output but will continue to work for the foreseeable future.
