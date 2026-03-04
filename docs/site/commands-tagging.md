# Tagging

Organize repositories with custom tags. Tags can be anything meaningful to your workflow — `personal`, `work`, `archived`, `backend`, `frontend`, etc.

Tags are used for filtering with `gws list --tag`.

---

## Add a Tag

```
gws tag add <repo> <tag>
```

Add a tag to all repositories matching the given identifier.

**How matching works:**

- Matches by **partial name** (case-insensitive): `gws tag add api work` tags "my-api", "api-gateway", etc.
- Also matches by **exact path**
- Tags are applied to **all matching repositories**

**Examples:**

```bash
# Tag a specific repo
gws tag add my-project personal

# Tag all API services as backend
gws tag add api backend
# Output: Added tag 'backend' to 3 repositories

# Tags can be anything
gws tag add old-service archived
```

---

## Remove a Tag

```
gws tag remove <repo> <tag>
```

Remove a tag from all repositories matching the given identifier.

**Examples:**

```bash
# Remove a tag from a specific repo
gws tag remove my-project personal

# Remove a tag from all matching repos
gws tag remove api backend
# Output: Removed tag 'backend' from 3 repositories
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

# Show repos with BOTH tags (AND logic)
gws list --tag work --tag backend

# Combine with other filters
gws list --tag work --type github --status
```
