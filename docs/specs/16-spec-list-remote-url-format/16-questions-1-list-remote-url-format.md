# 16 Questions Round 1 - List Remote URL Format

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. URL Format Output

What format should the "clean" URL display use? For example, given a raw remote of `git@github.com:daileyo/gws.git` or `https://user@github.com/daileyo/gws.git`, what should the formatted output look like?

- [ ] (A) Browser-friendly HTTPS URL: `https://github.com/daileyo/gws` (strip `.git` suffix, strip user info, normalize SSH to HTTPS)
- [X] (B) HTTPS URL preserving `.git`: `https://github.com/daileyo/gws.git` (strip user info, normalize SSH to HTTPS, keep `.git`)
- [ ] (C) Just strip user info but preserve protocol: `git@github.com:daileyo/gws.git` stays as-is (SSH has no user info to strip), `https://github.com/daileyo/gws.git` strips only `user@`
- [ ] (D) Other (describe)

## 2. Raw Flag Design

How should the `--raw` flag work with the existing `-r|--remote` flag?

- [ ] (A) `--raw` is a standalone modifier: `-r` shows formatted URL, `-r --raw` shows raw remote. `--raw` without `-r` is an error or ignored.
- [ ] (B) `--raw` is a global modifier that affects any column that has a "raw" variant (future-proofed)
- [X] (C) Instead of `--raw`, use a long-form variant like `--remote-raw` to keep it scoped to remote
- [ ] (D) Other (describe)

## 3. Short Flag for Raw

Should `--raw` have a short flag?

- [X] (A) Yes, `-R` (uppercase R, complements lowercase `-r` for remote)
- [ ] (B) Yes, a different short flag (specify which)
- [ ] (C) No, long-form `--raw` only is fine
- [ ] (D) Other (describe)

## 4. JSON Output Behavior

How should the JSON output (`-o json`) handle this change when `-r` is used?

- [ ] (A) Always include both `remote_url` (formatted) and `remote_url_raw` (original) in JSON regardless of `--raw` flag
- [X] (B) JSON follows the same flag behavior: `remote_url` shows formatted by default, raw with `--raw`
- [ ] (C) JSON always shows raw (current behavior unchanged), only table display changes
- [ ] (D) Other (describe)

## 5. Unfomattable URLs

What should happen if the remote URL can't be cleanly formatted (e.g., unusual protocols like `file://`, custom SSH aliases, or non-standard formats)?

- [X] (A) Fall back to displaying the raw URL as-is (best effort formatting)
- [ ] (B) Display the raw URL with a visual indicator that it couldn't be formatted
- [ ] (C) Other (describe)
