# 24 Tasks - Rebrand to omgitworks

## Relevant Files

**No source files are modified by this spec.** The files below are *read* to produce and verify
the inventory, and are listed so a future execution spec has the map already drawn.

- `go.mod` - Module path `github.com/daileyo/gws` (line 1).
- `Makefile` - `BINARY_NAME=git-workspace` (line 4); build target and shell-integration help text.
- `.goreleaser.yml` - `project_name` (line 6), `binary` (line 16), `brews.name` (line 55), tap `daileyo/homebrew-gws` (lines 56-58), release footer.
- `mkdocs.yml` - `site_name`, `site_url`, `repo_url`, `repo_name`, logo and favicon asset paths.
- `cmd/git-workspace/` - 36 Go files; the directory name itself is a touchpoint.
- `cmd/git-workspace/shellinit.go` - `{BIN}` placeholder plus literal `_git-workspace` and `__start_git-workspace` completion function names.
- `README.md` - Title, install instructions, examples.
- `docs/site/*.md` - 8 files.
- `.github/workflows/` - CI and release workflows.
- `release-please-config.json`, `.release-please-manifest.json` - Release automation.
- `docs/specs/24-spec-rebrand-omgitworks/` - **The only directory this spec writes to.**

### Notes

- The defining constraint: **this spec produces documents, not renames.** Task 5.0 verifies that.
- All counts must be measured with the recorded commands, not estimated, and re-measured before any future execution.
- `CHANGELOG.md` dominates the raw `git-workspace` count and is historical — exclude it when measuring live touchpoints.
- Follow conventional commits: `docs: ...` for everything here.

## Tasks

### [ ] 1.0 Verify and Record the Touchpoint Inventory

Reproduce the measured inventory in the spec, confirm the counts against the current tree, and record the exact commands so they can be re-run before execution.

#### 1.0 Proof Artifact(s)

- Documentation: Inventory table in the spec with counts matching a fresh run of the recorded commands demonstrates the measurement is reproducible
- Documentation: A live-count measurement excluding `CHANGELOG.md` demonstrates the real scope is separated from historical noise

#### 1.0 Tasks

- [ ] 1.1 Record the measurement commands verbatim in the spec, so counts can be re-derived: the `daileyo/gws` import count, the repo-wide `git-workspace` count, and the `cmd/git-workspace/` file count
- [ ] 1.2 Re-measure excluding `CHANGELOG.md` and record the live touchpoint count separately from the 604 raw total
- [ ] 1.3 Confirm each inventory row's file path and line number against the current tree at `HEAD`
- [ ] 1.4 Classify every row as irreversible, breaking, or cosmetic, and confirm the classification is stated in the spec
- [ ] 1.5 Flag every row that affects an already-installed user, since those drive the compatibility guarantee

### [ ] 2.0 Complete the Naming Decision Table

Resolve or explicitly defer each of the five names, so no future execution has to guess.

#### 2.0 Proof Artifact(s)

- Documentation: A decision table covering brand, repository, module path, binary, and command demonstrates completeness
- Documentation: The command-name question is either resolved or explicitly recorded as blocking demonstrates the decision is honest about what is unresolved

#### 2.0 Tasks

- [ ] 2.1 Confirm the recommendation for brand, repository, module path, and binary name, or record a different decision
- [ ] 2.2 Resolve open question 1 — what the user actually types. Evaluate keeping `gws` indefinitely, adopting `omg` or `omgw`, or shipping both
- [ ] 2.3 Verify short-name availability for whichever candidate is chosen: npm, crates.io, homebrew-core, and any conflict with an existing command on a default PATH
- [ ] 2.4 Record explicitly that the command name and binary name are decoupled by `shell-init`, and that this is what makes the binary rename invisible to users
- [ ] 2.5 Note in the decision table which entries are firm and which remain provisional

### [ ] 3.0 Finalize the Staged Migration Plan

Turn the proposed six stages into a plan with per-stage revert procedures and a stated compatibility guarantee.

#### 3.0 Proof Artifact(s)

- Documentation: Each stage carries an explicit revert procedure demonstrates reversibility is designed, not assumed
- Documentation: A stated compatibility guarantee naming versions demonstrates existing users are accounted for
- Documentation: The point of no return is identified demonstrates the irreversible step is known in advance

#### 3.0 Tasks

- [ ] 3.1 Write the revert procedure for each of the six stages
- [ ] 3.2 Verify the stage ordering holds — specifically that shell-integration compatibility (stage 3) ships before the binary rename (stage 2) reaches any user, and reorder if it does not
- [ ] 3.3 Define the compatibility guarantee: how long `shell-init` keeps emitting a `gws` function, and whether `git-workspace` ships as a symlink and until which version
- [ ] 3.4 Identify the point of no return and record the preconditions for crossing it
- [ ] 3.5 Confirm GitHub's repository-redirect behavior still covers web, API, and git-remote operations, and note that redirects lapse if a new repo is created at the old path
- [ ] 3.6 Decide whether to retain a placeholder repository at the old path to prevent namespace reuse
- [ ] 3.7 Note the lockstep requirement between the binary rename and the hardcoded `_git-workspace` / `__start_git-workspace` completion function names in `shellinit.go`

### [ ] 4.0 Preconditions and Go/No-Go Checklist

Produce the gate that execution has to pass, so the decision is made deliberately.

#### 4.0 Proof Artifact(s)

- Documentation: A precondition checklist with independently verifiable items demonstrates execution is gated
- Documentation: Stated abandonment criteria demonstrate the decision can still go the other way

#### 4.0 Tasks

- [ ] 4.1 Confirm `daileyo/omgitworks` is available as a repository name
- [ ] 4.2 Record a decision on pursuing the dormant `OMGItworks` GitHub username, and on whether it is needed at all given the repo lives under `daileyo`
- [ ] 4.3 Record a decision on registering `omgitworks.dev`, given `.com` is parked and `.co.uk` is held by an unrelated business
- [ ] 4.4 Record a decision on defensive registration of `omgitworks` on npm, PyPI, and crates.io, weighed against the squatting risk noted in Security Considerations
- [ ] 4.5 Do a trademark sanity check against the UK business and record the finding, including whether the class of goods differs sufficiently
- [ ] 4.6 Record the dependency that specs 22 and 23 land first, so the rename does not collide with in-flight work
- [ ] 4.7 Write the abandonment criteria — what would make keeping `git-workspace` the better call
- [ ] 4.8 Decide whether the rename ships in a major version bump or a minor release

### [ ] 5.0 Verify No Execution Occurred

Confirm the defining constraint held.

#### 5.0 Proof Artifact(s)

- CLI: `git diff --stat` for this spec's commits touches only `docs/specs/24-spec-rebrand-omgitworks/` demonstrates nothing was renamed
- CLI: `grep -rn 'omgitworks' --include='*.go' --include='*.yml' --include='Makefile' --include='go.mod' .` returns nothing demonstrates no source touchpoint was altered
- CLI: `make build && ./build/git-workspace --version` still reports the current binary demonstrates the tool is unchanged

#### 5.0 Tasks

- [ ] 5.1 Run `git diff --stat` across this spec's commits and confirm only the spec directory is touched
- [ ] 5.2 Grep the source tree for `omgitworks` outside `docs/specs/` and confirm no hits
- [ ] 5.3 Build and run the binary, confirming name and version are unchanged
- [ ] 5.4 Confirm the spec states plainly that no changes were made under it
