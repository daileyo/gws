# User Management

Manage git user profiles across your workspace. Profiles store identity information (name, email, signing key) that can be assigned to repositories to configure their local git `user.*` settings.

Profiles can be created manually or auto-detected from `~/.gitconfig` `includeIf` directives.

---

## Short-Flag Aliases

The `user` parent command supports short-flag aliases for common operations:

| Short Flag | Equivalent |
|------------|------------|
| `-l` | `gws user list` |
| `-a` | `gws user add` |
| `-s` | `gws user show` |
| `-d` | `gws user remove` |

When using `-a`, you can also pass `--email`, `--name`, `--signing-key`, and `--sign-commits` inline:

```bash
gws user -a work --email work@company.com --name "Jane Doe"
```

---

## List Profiles

```
gws user list
```

List all stored and auto-detected user profiles.

**Example output:**

```
Stored Profiles:

  NAME        GIT NAME     EMAIL                    SIGN
  ----        --------     -----                    ----
  personal    Jane Doe     jane@personal.dev        ABC123 (signing on)
  work        Jane Doe     jane.doe@company.com     —

Auto-Detected Profiles:

  NAME        GIT NAME     EMAIL                    SIGN
  ----        --------     -----                    ----
  oss         Jane Doe     jane@opensource.org       —

(Auto-detected from ~/.gitconfig includeIf directives)
```

---

## Add Profile

```
gws user add <name> [flags]
```

Create a new user profile.

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--email` | Yes | | Email address for git commits |
| `--name` | No | Profile name | Git `user.name` value |
| `--signing-key` | No | | GPG signing key ID |
| `--sign-commits` | No | `false` | Enable commit signing (`commit.gpgsign`) |

**Examples:**

```bash
# Add a basic profile
gws user add personal --email jane@personal.dev

# Add with a different git name
gws user add work --email jane.doe@company.com --name "Jane Doe"

# Add with signing enabled
gws user add secure --email jane@secure.dev --signing-key ABC123DEF456 --sign-commits
```

---

## Show Profile

```
gws user show <name>
```

Display detailed information about a profile, including how many repositories use it.

**Example output:**

```
Profile: personal

  Git Name:       Jane Doe
  Email:          jane@personal.dev
  Signing Key:    ABC123DEF456
  Sign Commits:   true

  Used by 3 repositories
```

---

## Remove Profile

```
gws user remove <name>
```

Remove a stored user profile. If any repositories are currently using the profile, you will be prompted before removal.

---

## Assign Profile

```
gws user assign <repository> <profile> [flags]
```

Assign a user profile to a repository. This sets `user.name` and `user.email` (and optionally signing configuration) in the repository's local `.git/config`.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--use-subdirs` | `false` | Move repository to profile subdirectory |
| `--dry-run` | `false` | Preview changes without applying them |

**Examples:**

```bash
# Assign the "work" profile to a repository
gws user assign my-api work

# Preview what would change
gws user assign my-api work --dry-run
```

---

## Sync Profiles

```
gws user sync
```

Synchronize stored user information with the effective git configuration for all tracked repositories. This re-reads each repository's `.git/config` and updates the cached `user`, `email`, `signing_enabled`, and `user_source` fields in the workspace configuration.

**When to use:**

- After manually editing a repository's `.git/config`
- After changing `~/.gitconfig` `includeIf` rules
- To ensure the workspace config reflects the current state
