# 05-tasks-github-pages-docs

## Relevant Files

- `mkdocs.yml` — New file. MkDocs configuration: site name, site URL, `docs_dir: docs/site`, Material theme, and explicit `nav` order.
- `requirements.txt` — New file. Pins `mkdocs-material` to a specific version so the build is reproducible.
- `docs/site/index.md` — New file. Documentation site home page with project description and links to all sections.
- `docs/site/getting-started.md` — New file. Installation and first-time setup guide, migrated from `README.md`.
- `docs/site/commands.md` — New file. Full commands reference with all flags and output examples, migrated from `README.md`.
- `docs/site/shell-integration.md` — New file. Shell function setup and usage guide, migrated from `README.md`.
- `docs/site/configuration.md` — New file. Config file location, structure, and all fields, migrated from `README.md`.
- `.github/workflows/docs.yml` — New file. GitHub Actions workflow that builds and deploys the MkDocs site to GitHub Pages on `v*` tag pushes.
- `.gitignore` — Modified. Add `site/` so the MkDocs build output is not committed.
- `README.md` — Modified. Simplified to under 150 lines; removes all how-to usage sections; adds docs site link; updates Contributing section.

### Notes

- There are no automated tests to run for this feature — it is documentation-only. Verification is done by building and previewing the site locally with `mkdocs serve`.
- Follow the existing GitHub Actions workflow style (`actions/checkout@v4`, named steps, `permissions` block at the top of the file).
- Use Conventional Commits: `docs:` for documentation files and README changes, `ci:` for the workflow file.
- The GitHub Pages source setting ("GitHub Actions") is a one-time manual change in the repository **Settings → Pages → Source**. It cannot be done via a file commit — it must be done in the browser.

## Tasks

### [x] 1.0 Set Up MkDocs Site Foundation and Populate Documentation Pages

#### 1.0 Proof Artifact(s)

- Screenshot: Terminal showing `mkdocs serve` running without errors (e.g., `INFO - Serving on http://127.0.0.1:8000`) demonstrates the site builds successfully.
- Screenshot: Browser at `http://127.0.0.1:8000` showing the site home page with the site title and navigation links to Getting Started, Commands Reference, Shell Integration, and Configuration demonstrates the MkDocs configuration and all four pages are correct.

#### 1.0 Tasks

- [x] 1.1 Add `site/` on its own line to `.gitignore` (under the `# Build directory` section, alongside `build/` and `dist/`) so MkDocs build output is never committed.
- [x] 1.2 Create `requirements.txt` at the repository root. Look up the latest stable version of `mkdocs-material` on PyPI and pin it (e.g., `mkdocs-material==9.5.50`). One dependency, one line.
- [x] 1.3 Create `mkdocs.yml` at the repository root with the following configuration:
  - `site_name: git-workspace`
  - `site_url: https://daileyo.github.io/gws`
  - `docs_dir: docs/site`
  - `theme: name: material`
  - `nav:` listing the five pages in order: `Home: index.md`, `Getting Started: getting-started.md`, `Commands Reference: commands.md`, `Shell Integration: shell-integration.md`, `Configuration: configuration.md`
- [x] 1.4 Create `docs/site/index.md` as the site home page. Include: a one-sentence description of `git-workspace`, the status badges (CI, Go Report Card, Snyk, Release, License) copied from the README, and a short bullet list linking to each of the four documentation sections.
- [x] 1.5 Create `docs/site/getting-started.md`. Migrate the following content from `README.md`: the Installation section (Homebrew install commands from `04-spec-brew-install` and the "Build from Source" block), and the Quick Start section (steps 1–4 covering `--init`, viewing workspace info, `--list`, and `--version`). Use `bash` fenced code blocks throughout.
- [x] 1.6 Create `docs/site/commands.md`. Migrate the following content from `README.md`: the Repository Classification and Tagging section (automatic classification table, visibility detection, custom tags and tag matching, filtering repositories, output formats, the refresh workspace section). Include all flag examples and sample output blocks.
- [x] 1.7 Create `docs/site/shell-integration.md`. Migrate the following content from `README.md`: the Shell Navigation section (`cdgws`/`gcd` setup and usage) and the Repository Navigation section (the `gws` shell function setup, usage, wildcard matching, multiple-match selection, no-match suggestions, and binary direct usage). Use `bash` fenced code blocks throughout.
- [x] 1.8 Create `docs/site/configuration.md`. Migrate the following content from `README.md`: the Configuration section (config file path, field descriptions, and the full JSON example).
- [x] 1.9 Install dependencies and verify the site builds: run `pip install -r requirements.txt` then `mkdocs build`. Confirm no errors appear and a `site/` directory is created at the repo root.
- [x] 1.10 Run `mkdocs serve` and open `http://127.0.0.1:8000` in a browser. Click through all five pages (Home, Getting Started, Commands Reference, Shell Integration, Configuration) and confirm each page loads with the correct content and navigation. Capture the proof artifact screenshots.

---

### [ ] 2.0 Add GitHub Actions Deployment Workflow and Enable GitHub Pages

#### 2.0 Proof Artifact(s)

- Screenshot: `.github/workflows/docs.yml` file on the `daileyo/gws` GitHub repository showing the file contents, demonstrating the workflow is committed.
- Screenshot: GitHub Actions run for `docs.yml` on the `daileyo/gws` repository showing all steps with green checkmarks, demonstrating end-to-end deployment succeeds.
- Screenshot: `https://daileyo.github.io/gws` rendered in a browser showing the MkDocs site home page, demonstrating the site is publicly accessible.

#### 2.0 Tasks

- [ ] 2.1 Create `.github/workflows/docs.yml`. The workflow should:
  - Be named `Deploy Docs`
  - Trigger on `push` with tag pattern `v*`
  - Declare `permissions` at the top level: `contents: read`, `pages: write`, `id-token: write`
  - Have a single job named `deploy` with `runs-on: ubuntu-latest` and an `environment` block setting `name: github-pages` and `url: ${{ steps.deployment.outputs.page_url }}`
  - Include these steps in order:
    1. `actions/checkout@v4` (Checkout code)
    2. `actions/setup-python@v5` with `python-version: '3.12'` (Set up Python)
    3. `run: pip install -r requirements.txt` (Install MkDocs Material)
    4. `run: mkdocs build` (Build documentation)
    5. `actions/upload-pages-artifact@v3` with `path: site/` (Upload Pages artifact)
    6. `actions/deploy-pages@v4` with `id: deployment` (Deploy to GitHub Pages)
- [ ] 2.2 In the `daileyo/gws` GitHub repository, go to **Settings → Pages → Build and deployment → Source** and select **"GitHub Actions"**. This is a one-time manual change that allows the workflow to deploy to GitHub Pages. (Note: this cannot be done via a file commit.)
- [ ] 2.3 Commit and push `docs.yml` (and all Task 1.0 files if not already pushed) on the feature branch, then open or update the pull request and merge to `main`.
- [ ] 2.4 After merging to `main`, push a release tag (e.g., the next version tag following the existing `v2.0.0`) to trigger the `docs.yml` workflow. Monitor the Actions tab on GitHub to confirm all steps complete with green checkmarks. Capture the proof artifact screenshots.

---

### [ ] 3.0 Simplify README and Update Contributing Section

#### 3.0 Proof Artifact(s)

- Screenshot: `README.md` rendered on `https://github.com/daileyo/gws` showing badges, description, features list, quick install snippet, a prominent link to the docs site, and the Development section — with no CLI how-to usage sections — demonstrating the README is concise and developer-focused.
- Screenshot: The Contributing section of the rendered README showing the "not currently open for contributions" message and invitation to open an issue or reach out, demonstrating the correct contribution posture.

#### 3.0 Tasks

- [ ] 3.1 Remove the following sections entirely from `README.md`: Quick Start (all four numbered steps), Repository Classification and Tagging, Manual Verification Steps, Configuration, Output Formats (this is embedded in the list section — remove it there), Shell Navigation, Repository Navigation, Roadmap, Release Process, and Security.
- [ ] 3.2 After the features bullet list and before the Development section, add a new `## Documentation` section (or similar heading) with a single line linking to the docs site: `Full documentation is available at https://daileyo.github.io/gws`.
- [ ] 3.3 Update the `## Installation` section to show only the Homebrew install command (the primary method) with a brief note that build-from-source instructions are in the documentation. Remove the detailed `Build from Source` block from the README.
- [ ] 3.4 Rewrite the `## Contributing` section. Replace the current "Contributions are welcome!" opening and the Development Workflow, Running CI Locally, and Code Style sub-sections with the following:
  - A clear statement that the project is not currently open for external contributions
  - A brief note that this may change in the future
  - An invitation for interested contributors to open a GitHub issue or reach out directly
  - Retain the **Commit Message Format** sub-section (types table and breaking change example) so the conventions are documented
- [ ] 3.5 Review the final `README.md` and confirm: (a) it is under 150 lines, (b) it contains no CLI how-to usage content, (c) the docs link is visible near the top, and (d) the Contributing section has the correct closed-with-contact message. Capture the proof artifact screenshots after merging to `main`.
