# 03 Questions Round 2 - Command Flag Rework

Follow-up based on your Round 1 answers. There's one area where answers interact that I want to confirm.

## 1. Tag Filter on List — Confirming the Pattern

In Q1 you noted that the list filter `--tag` should become a positional subcommand `tags` (e.g., `gws -l tags personal`). In Q3 you selected (C) to keep current flag patterns for filters. I want to confirm how these work together.

My understanding: all list filters remain as flags **except** `--tag` which becomes the positional subcommand `tags` to avoid collision with the `--tag`/`-t` command. So the full list command surface would be:

```
gws --list                              # List all repos
gws -l                                  # Same (shorthand)
gws -l --type github                    # Filter by type (flag stays)
gws -l -y                               # New shorthand
gws -l --name myproject                 # Filter by name (flag stays)
gws -l -n                               # New shorthand
gws -l --path /home/user                # Filter by path (flag stays)
gws -l -p                               # New shorthand
gws -l --output json                    # JSON output (flag stays)
gws -l -o json                          # Same (shorthand stays)
gws -l --status                         # Show git status (flag stays)
gws -l -s                               # Same (shorthand stays)
gws -l --tag personal                   # Filter by tag (was --tag, now positional)
gws -l -t
gws -l --tag work backend               # Multiple tags (AND logic, was --tag work --tag backend)
gws -l -t work backend
gws -l --tag personal --type github -s   # Combine tag subcommand with other flags
```

- [ ] (A) Yes, this is exactly what I want
- [ ] (B) No, `--tag` should also remain a flag on list — resolve the collision differently
- [X] (C) Other all filter commands should be combine-able. they should also support wildcard. i.e. wi* match any with wi

Notes:

