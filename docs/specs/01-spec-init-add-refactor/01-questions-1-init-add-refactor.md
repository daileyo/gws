# 01 Questions Round 1 - Init & Add Refactor

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

---

## 1. The "gws Directory" — What Does It Actually Look Like?

You mentioned the initialized gws directory "will serve as an actual directory in the filesystem that can be navigated to without gws itself." What does this mean in practice?

- [X] (A) **Symlink directory** — gws creates a single directory (e.g. `~/gws/`) that contains symlinks pointing to each tracked git repository. Users can `cd ~/gws/my-repo` to navigate directly.
- [ ] (B) **Marker directory** — gws creates a plain directory and places a metadata marker file (e.g. `.gws`) inside it. The repos live wherever they already are; the directory is just the "home" of the workspace. Users navigate repos using `gws -g`.
- [ ] (C) **Mirror directory** — gws creates a directory tree that mirrors the structure of tracked repos relative to the workspace root. No symlinks, just a navigable representation.
- [X] (D) Other (describe) The gws directory can have both actual repositories as well as symlinks to ones elswhere on the filesystem.

---

## 2. Init — Interactive Prompt for Location

When `gws --init` is run and no workspace exists yet, you want to prompt for location (default: current working directory). What form should this prompt take?

- [ ] (A) **Simple text prompt** — Print `"Enter workspace path [/current/dir]: "` to stderr and read a line from stdin (similar to the existing navigation prompt pattern).
- [ ] (B) **Flag-based** — Require the user to pass the path as an argument, e.g. `gws --init /path/to/workspace`. No interactive prompt; default is current directory if no argument given.
- [ ] (C) **Both** — Accept an optional path argument; if none given, show an interactive prompt.
- [X] (D) Other (describe) no prompt. if the user types just gws --init it will create it in the current directory. However, it should only create it if there isn't already a gws directory that is known to gws.

---

## 3. The `-a` / `--add-tag` Short Flag Conflict

Currently `-a` maps to `--add-tag`. You want `-a` to map to the new `--add` command instead. What should happen to `--add-tag`?

- [ ] (A) **Remove short form** — `--add-tag` keeps its long form but loses the `-a` shortcut. Users must type `--add-tag` in full.
- [X] (B) **Reassign short form** — Give `--add-tag` a new short form (e.g. `-T` or `-A`). Please specify which letter you prefer.
- [ ] (C) Other (describe)

> **Note**: If you prefer (B), please add your preferred letter below.
-g
---

## 4. `--add` Without a Path — What Does It Add?

When `gws --add` is run with no path argument, the default is to add the current directory. But what should happen in each scenario?

- [X] (A) **Add current dir if it's a git repo; error if not.** Keep it simple.
- [ ] (B) **Add current dir if it's a git repo; silently do nothing if not (exit 0).**
- [ ] (C) **Add current dir if it's a git repo; if not, automatically behave as if `--recursive` was passed and scan subdirectories.**
- [ ] (D) Other (describe)

---

## 5. `--add` With a Path Argument

Should `gws --add /some/path` be supported (adding a specific path rather than defaulting to the current directory)?

- [X] (A) **Yes** — `gws --add /some/path` adds that specific path (must be a git repo). If it's a directory of repos, users should use `--recursive` instead.
- [ ] (B) **No** — The path is always inferred from the current directory. To add a specific path, users `cd` there first.
- [ ] (C) Other (describe)

---

## 6. `--add` — Already-Tracked Repositories

When `--add` encounters a git repo that's already in the gws metadata, what should happen?

- [ ] (A) **Skip silently** — No output, no error. The repo is already tracked; nothing to do.
- [X] (B) **Warn and skip** — Print a message like `"my-repo is already tracked, skipping."` and continue.
- [ ] (C) **Error** — Treat as an error condition and exit non-zero.
- [ ] (D) Other (describe)

---

## 7. Tab Completion — What Kind?

You mentioned tab completion for `--add` when typing a path. What kind of completion did you have in mind?

- [ ] (A) **Shell completion scripts** — Generate bash/zsh/fish completion scripts (via Cobra's built-in `completion` subcommand). This gives native `<Tab>` behavior in the shell for all flags, including path arguments.
- [ ] (B) **In-app fuzzy search prompt** — When `--add` is run interactively, show a filterable list of nearby directories so the user can type to filter and select.
- [X] (C) **Both** — Shell completion scripts for normal shell use, plus an interactive prompt when no argument is provided.
- [ ] (D) **Out of scope for now** — Focus on the init/add behavior; tab completion can be addressed separately.
- [ ] (E) Other (describe)

---

## 8. `--add --recursive` — What Gets Reported?

For `gws --add --recursive`, you mentioned reporting how many repos were added. What's the right level of detail?

- [ ] (A) **Count only** — `"Added 3 repositories."` / `"No repositories found."`
- [X] (B) **Count + names** — `"Added 3 repositories: repo-a, repo-b, repo-c"`
- [ ] (C) **Full table** — List each repo with name and path, similar to how `--init` currently displays results.
- [ ] (D) Other (describe)

---

## 9. Init — Already-Initialized Workspace

When a workspace already exists and the user runs `gws --init`, you want to notify them and mention `--add`. Should the command also do anything else?

- [ ] (A) **Notify only** — Print `"Workspace already initialized at /path. Use --add to add repositories."` and exit 0.
- [X] (B) **Notify + re-scan option** — Show the notification, but also suggest `--refresh` to re-scan the existing workspace.
- [ ] (C) **Notify + show current state** — Display the notification and list current workspace info (path, repo count).
- [ ] (D) Other (describe)
