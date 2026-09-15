# 24 Task 3.0 Proofs - Staged Migration Plan

## Proof 1: six stages, each with a revert procedure

| Stage | Revert |
| --- | --- |
| 1. Brand surface | `git revert` — nothing installed changes |
| 2. Shell-integration compatibility | `git revert` — users keep `gws` either way |
| 3. Binary and archive naming | noisy — a released archive cannot be unpublished; needs a follow-up release |
| 4. Repository rename | rename back; safe **only** while the old path stays empty |
| 5. Go module path | `git revert`; the old path still works via redirects |
| 6. Brew formula and tap | re-publish under the old name; keep the old tap until proven |

## Proof 2: the ordering constraint holds (task 3.2)

Stage 2 precedes stage 3. This is the constraint that protects existing users, and it is not
theoretical: a user whose rc file contains `eval "$(git-workspace shell-init zsh)"` would get
`command not found` on **every new shell** if the binary were renamed first. Shipping
dual-function `shell-init` before the rename is the first defense; the `git-workspace` symlink
in stage 3 is the second.

## Proof 3: compatibility guarantee is stated, not implied

`shell-init` emits `gws` indefinitely; `omgw` from stage 2; the old binary ships as a symlink
through the next major; the old tap stays installable one release cycle past stage 6.

## Proof 4: point of no return identified

Stage 4 plus stage 6. Everything earlier is `git revert` and a patch release. After stage 4
third parties link to the new name; after stage 6 users' installed taps point at it.

## Proof 5: GitHub redirect behavior verified against the docs

Confirmed from GitHub's own documentation on renaming a repository, which surfaced a
consequence the first draft of this spec missed:

- Web traffic and git operations (clone, fetch, push) **do** redirect.
- **GitHub Pages URLs do not.** The docs site is at `https://daileyo.github.io/gws`
  (`mkdocs.yml:2`); after a rename that URL breaks outright, taking existing links, bookmarks,
  and search results with it.
- Redirects lapse permanently if anything is later created at the old path.

The Pages finding changes the domain decision from a branding nicety to the mechanism that
keeps documentation reachable across the rename. It is now the top recommendation in
Preconditions.

## Proof 6: lockstep requirement recorded (task 3.7)

```
$ grep -n '_git-workspace\|__start_git-workspace' cmd/git-workspace/shellinit.go
92:  compdef _git-workspace gws
163: complete -o default -F __start_git-workspace gws
```

These are Cobra's generated completion function names, derived from the binary name. Renaming
the binary changes what `completion zsh` / `completion bash` emit. If these two lines are not
updated in the same commit, tab completion **silently stops working** — no error, it simply
does nothing.
