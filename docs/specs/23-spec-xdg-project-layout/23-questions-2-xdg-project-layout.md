# 23 Questions Round 2 - XDG Project Layout

Follow-up captured 2026-09-04.

## 1. Windows path convention

Windows has no XDG Base Directory Specification. What should gws do there?

- [X] (A) **Mimic XDG under `%USERPROFILE%`.** Identical layout on every platform:
      `<home>/.config/gws` and `<home>/.local/share/gws/projects`, where `<home>` is
      `%USERPROFILE%` on Windows. `XDG_CONFIG_HOME` / `XDG_DATA_HOME` honored on Windows too.
- [ ] (B) Use Windows-native locations — `%AppData%\gws` and `%LocalAppData%\gws\projects`
- [ ] (C) Native by default, XDG when the variables are set

**Rationale from the requester:** "I want to retain cross platform functionality. That is a
core desire. I know it doesn't have XDG; but if we can mimic it, or be opinionated on a
structure like `C:\Users\<user>\.config` or something in the Windows ecosystem that is
equivalent, that would be ideal."

**Supporting precedent found during research:** git itself already does this. Per
`git-config(1)`: "When the XDG_CONFIG_HOME environment variable is not set or empty,
`$HOME/.config/` is used as `$XDG_CONFIG_HOME`." Git for Windows sets `$HOME` to
`%USERPROFILE%`, so `C:\Users\<user>\.config\git\config` is already a real and supported path
on Windows machines. A git-adjacent tool storing its config beside git's own is following the
convention of its ecosystem rather than inventing one.

**Consequences recorded in the spec:**

- `os.UserConfigDir()` must **not** be used — it returns `%AppData%` on Windows and would
  break parity. Resolution is hand-written on top of `os.UserHomeDir()`.
- No `runtime.GOOS` branch is needed in path resolution at all.
- Cross-platform parity is promoted to a success metric with its own test.
- The `%USERPROFILE%`-rooted path is shorter than the `%LocalAppData%` alternative, which
  slightly eases `MAX_PATH` pressure.
- Users who prefer the native Windows location can still set `XDG_CONFIG_HOME` /
  `XDG_DATA_HOME` explicitly; this is documented rather than special-cased.

## 2. Acknowledged trade-off

Choosing (A) diverges from Windows-native convention, where application configuration belongs
in `%AppData%`. This is deliberate: `%AppData%` would give Windows users a different layout
from their colleagues, a second set of documentation, and a path that XDG-aware tooling and
dotfile managers do not know about. Cross-platform parity was stated as a core requirement and
outranks native-convention conformance here.
