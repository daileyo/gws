# 24 Questions Round 2 - Rebrand to omgitworks

Follow-up captured 2026-09-04.

## 1. What is driving the rebrand?

The requester identified two distinct motivations, both now recorded as goals in the spec:

1. **Distribution availability.** Land on a name that is not heavily used and that can actually
   be claimed on brew, winget, and other package managers. Availability is a selection
   criterion, not a nice-to-have — a name that cannot be claimed on a channel is a name that
   costs the project an install path.
2. **Memorability.** `omgitworks` is a double reading and both are intended: *"OMG it works"*
   and *"om git works(pace)"* — literally what the tool is. The long form is what makes the
   name stick and stay searchable.

## 2. What does the user type?

Previously open question 1; now resolved.

- [X] (A) **`omgw`** — "default omgw I like." Four characters, derives transparently from the
      brand, available on every channel checked.
- [ ] (B) `omg` — shorter, but taken on npm
- [ ] (C) Keep `gws` indefinitely and treat `omgitworks` purely as a brand
- [ ] (D) Ship both a short and long command

**Resolution:** brand is `omgitworks`, command is `omgw`. These are decided separately because
`shell-init` already decouples the typed command from the binary name — the long form carries
the brand and the short form carries the typing.

## 3. Package-manager availability, checked 2026-09-04

| Channel | `omgitworks` | `omgw` | `gws` (incumbent) |
| --- | --- | --- | --- |
| homebrew-core | Available | Available | **Taken** |
| winget | Available | Available | — |
| scoop (Main + Extras) | Available | Available | — |
| chocolatey | Available | Available | — |
| AUR | Available | Available | **Taken** |
| npm | Available | Available | **Taken** |
| crates.io | Available | Available | — |
| PyPI | Available | Available | — |

Both candidates are clear everywhere checked. The incumbent short name is taken on three
channels, in every case by StreakyCobra's `gws` — a tool with the same one-line description as
this one.

**Caveat:** availability decays. These results must be re-checked immediately before execution,
which is recorded as both a precondition and task 4.5.

## 4. Still open

- Whether `omgw` collides with a command already on a default PATH — the registry checks do not
  cover this (task 2.4)
- Whether to defensively register both names on npm, PyPI, and crates.io (task 4.4)
- Whether to claim winget and scoop manifests as part of the rebrand, given neither is a
  shipping channel today (task 4.6)
