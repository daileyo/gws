# 24 Task 2.0 Proofs - Naming Decisions

## Proof 1: all five names decided, with firmness marked

| Name | Decision | Firm? |
| --- | --- | --- |
| Brand | `omgitworks` | firm |
| Command | `omgw` | firm |
| Binary | `omgitworks` | firm |
| Repository | `daileyo/omgitworks` | provisional |
| Module path | `github.com/daileyo/omgitworks` | provisional |
| Brew tap | undecided | open |

Marking firmness matters: a future execution spec needs to know which entries are settled and
which are still arguable, rather than treating the whole table as equally decided.

## Proof 2: PATH collision check (task 2.4)

Registry availability says nothing about whether a binary by that name already exists on a
typical system. Checked 2026-09-14:

```
$ command -v omgw
(not found)

$ curl -s "https://packages.debian.org/search?searchon=contents&keywords=omgw\
&mode=exactfilename&suite=stable&arch=any" | grep -o 'no results'
Sorry, your search gave no results

$ curl -s -o /dev/null -w '%{http_code}' https://formulae.brew.sh/api/formula/omgw.json
404
$ curl -s -o /dev/null -w '%{http_code}' https://formulae.brew.sh/api/cask/omgw.json
404
```

No Debian package ships a file named `omgw`, no brew formula or cask claims it, and it is not
a shell builtin or keyword.

## Proof 3: the decoupling is recorded (task 2.6)

The spec now states explicitly that `shell-init` generates a shell function whose name is
independent of the binary it invokes — users type `gws`, the binary is `git-workspace`. This
is the single fact that makes the rename survivable at the user-facing layer, and it is why
brand and command are decided separately.

## Proof 4: the mnemonic is recorded (task 2.5)

Both readings are captured: *"OMG, it works"* and *"om git works(pace)"*. This is the
justification for keeping a ten-character brand alongside a four-character command — the long
form is what makes the project searchable, which is the entire motivation for the change.
