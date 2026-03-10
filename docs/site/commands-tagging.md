# Tagging

Organize repositories with custom tags. Tags can be anything meaningful to your workflow — `personal`, `work`, `archived`, `backend`, `frontend`, etc.

Tags are used for filtering with `gws list --tag`.

---

## Short-Flag Aliases

The `tag` parent command supports short-flag aliases for common operations:

| Short Flag | Description |
|------------|-------------|
| `-a` | Add a tag (equivalent to `tag add`) |
| `-d` | Remove a tag (equivalent to `tag remove`) |
| `-p` | Match repositories by path |
| `-r` | Match repositories by name |

**Examples:**

```bash
# Add a tag using short flag
gws tag -a my-repo work

# Remove a tag using short flag
gws tag -d my-repo work
```

---

## Add a Tag

```
gws tag add <repo> <tag>
```

Add a tag to all repositories matching the given identifier.

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--path` | `-p` | Match repositories by path prefix or substring (case-sensitive) |
| `--repo` | `-r` | Match repositories by name (partial, case-insensitive) |

**How matching works:**

- With no flags: matches by **partial name** (case-insensitive) or **exact path**
- With `--path`: matches by path prefix or substring (case-sensitive)
- With `--repo`: matches by partial name (case-insensitive)
- Combine `--path` and `--repo` to require both conditions (AND logic)
- Tags are applied to **all matching repositories**

**Examples:**

```bash
# Tag a specific repo
gws tag add my-project personal

# Tag all API services as backend
gws tag add api backend
# Output: Added tag 'backend' to 3 repositories

# Tag by path
gws tag add --path /home/user/work backend

# Tag by repo name
gws tag add --repo api backend

# Tag matching both path and name
gws tag add --repo api --path /work backend

# Tags can be anything
gws tag add old-service archived
```

---

## Remove a Tag

```
gws tag remove <repo> <tag>
```

Remove a tag from all repositories matching the given identifier.

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--path` | `-p` | Match repositories by path prefix or substring (case-sensitive) |
| `--repo` | `-r` | Match repositories by name (partial, case-insensitive) |

**Examples:**

```bash
# Remove a tag from a specific repo
gws tag remove my-project personal

# Remove a tag from all matching repos
gws tag remove api backend
# Output: Removed tag 'backend' from 3 repositories

# Remove by path
gws tag remove --path /home/user/work backend

# Remove by repo name
gws tag remove --repo api backend
```

---

## Tab Completion

Tab completion is available for tag operations when shell integration is set up (see [Shell Integration](shell-integration.md)):

- **First argument** (repo): Completes repository names from your workspace
- **Second argument** (tag, for `tag remove`): Completes existing tags on the matched repository

---

## Using Tags for Filtering

Once tagged, use `gws list --tag` to filter:

```bash
# Show all personal repos
gws list --tag personal

# Filter by tag and show the tags column
gws list -T personal

# Combine with other filters
gws list --tag work --type github -S
```

!!! note
    `--tag`/`-t` accepts a **single value**. To filter by tag and display other columns, combine with uppercase show-column flags.
