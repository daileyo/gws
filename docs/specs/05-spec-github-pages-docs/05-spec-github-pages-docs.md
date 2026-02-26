# 05-spec-github-pages-docs

## Introduction/Overview

This feature creates a published documentation site for `git-workspace` using MkDocs with the Material theme, hosted on GitHub Pages at `https://daileyo.github.io/gws`. How-to usage content (commands, shell integration, configuration) will be moved from `README.md` into organized documentation pages, and the README will be simplified to a developer-focused brief overview. The goal is to give users a clean, searchable reference manual while keeping the README focused on development and contribution.

## Goals

- Set up a MkDocs site with the Material theme, sourced from `docs/site/`, buildable locally with a single command
- Publish four user-facing documentation pages: Getting Started, Commands Reference, Shell Integration, and Configuration
- Automate deployment to GitHub Pages on every release tag via a new GitHub Actions workflow
- Simplify `README.md` to a brief overview (badges, description, features list, install snippet, docs link, Development, Contributing)
- Update the Contributing section to clearly state the project is not yet open for external contributions, with an invitation to open an issue or reach out

## User Stories

- **As a new user**, I want a clear documentation site so that I can quickly understand how to install and use `git-workspace` without reading a long README.
- **As a returning user**, I want a searchable commands reference so that I can quickly look up flags and examples without running `--help`.
- **As a developer evaluating the project**, I want a simplified README that links to full docs so that I can find what I need without wading through CLI usage examples.
- **As a potential contributor**, I want to know upfront that the project is not yet open for contributions so that I know how to engage with the project.

## Demoable Units of Work

### Unit 1: MkDocs Site Foundation and Content

**Purpose:** Establishes the documentation site structure locally, with all four documentation pages populated from existing README content. This proves the site builds and renders correctly before any deployment is configured.

**Functional Requirements:**
- The system shall have a `mkdocs.yml` file at the repository root that configures the site name, site URL (`https://daileyo.github.io/gws`), the Material theme, and `docs_dir: docs/site`.
- The system shall have a `docs/site/` directory containing the following markdown source files:
  - `index.md` — site home page (brief project description and links to main sections)
  - `getting-started.md` — installation (Homebrew, build from source) and first-time setup (`--init`)
  - `commands.md` — all flags and commands with examples and output samples
  - `shell-integration.md` — setup and usage of the `gws` shell function and `cdgws`/`gcd` aliases
  - `configuration.md` — config file location (`~/.gws/config.json`), structure, and all fields
- The system shall have a `requirements.txt` (or equivalent) at the repository root pinning `mkdocs-material` so the site can be built with `pip install -r requirements.txt && mkdocs build`.
- The MkDocs build output directory (`site/`) shall be listed in `.gitignore` so generated files are not committed.
- The user shall be able to run `mkdocs serve` locally and view all four pages in a browser with correct content and navigation.

**Proof Artifacts:**
- Screenshot: `mkdocs serve` output in the terminal showing the local development server running without errors, demonstrating the site builds successfully.
- Screenshot: Browser showing the rendered `https://127.0.0.1:8000` home page with the site title and navigation links to all four sections, demonstrating the MkDocs configuration is correct.

---

### Unit 2: GitHub Actions Deployment Workflow

**Purpose:** Configures automated deployment so the documentation site is published to GitHub Pages whenever a release tag is pushed, requiring no manual steps.

**Functional Requirements:**
- The system shall have a new GitHub Actions workflow file at `.github/workflows/docs.yml` that triggers on `push` events matching tags of the pattern `v*` (e.g., `v2.1.0`).
- The workflow shall install Python, install `mkdocs-material` from `requirements.txt`, build the site with `mkdocs build`, and deploy the contents of the `site/` output directory to GitHub Pages using the `actions/upload-pages-artifact` and `actions/deploy-pages` actions.
- The workflow shall use `pages: write` and `id-token: write` permissions, consistent with the official GitHub Pages Actions deployment pattern.
- The GitHub repository's Pages settings shall be configured to use "GitHub Actions" as the deployment source (not a branch).
- The user shall be able to visit `https://daileyo.github.io/gws` after a tagged release and see the published documentation site.

**Proof Artifacts:**
- Screenshot: `.github/workflows/docs.yml` file shown in the repository on GitHub, demonstrating the workflow file exists and is committed.
- Screenshot: GitHub Actions run for the `docs.yml` workflow showing all steps completed successfully (green checkmarks), demonstrating the deploy pipeline runs end-to-end.
- Screenshot: `https://daileyo.github.io/gws` rendered in a browser, demonstrating the site is live and accessible.

---

### Unit 3: README Simplification

**Purpose:** Updates `README.md` so it serves developers and contributors, not end users. Removes how-to usage content (now covered by the docs site) and adds a clear link to the published documentation.

**Functional Requirements:**
- The `README.md` shall retain: status badges, one-line project description, features bullet list, a quick install snippet (Homebrew command), a prominent link to `https://daileyo.github.io/gws` as the full documentation, the Development section (build commands, running tests, project structure), and the Contributing section.
- The `README.md` shall remove: the Quick Start walkthrough, the Repository Classification and Tagging section, the Manual Verification Steps section, the Configuration section, the Output Formats section, the Shell Navigation and Repository Navigation sections, and the Roadmap, Release Process, and Security sections.
- The Contributing section shall state that the project is not currently open for external contributions, include a brief explanation that this may change in the future, and invite interested contributors to open a GitHub issue or reach out directly.
- The Contributing section shall retain the commit message format conventions (Conventional Commits types and examples) so that any future contributors have the context they need.

**Proof Artifacts:**
- Screenshot: `README.md` rendered on `https://github.com/daileyo/gws` showing the simplified layout with the docs site link visible near the top and the updated Contributing section, demonstrating the README no longer contains how-to usage details.

## Non-Goals (Out of Scope)

1. **Versioned documentation**: The docs site will not support multiple versions of the documentation (e.g., v1.x, v2.x). Only the latest released docs will be published.
2. **Custom domain**: No custom domain configuration. The site will be served at the default `https://daileyo.github.io/gws` URL.
3. **Search backend customization**: The built-in MkDocs Material search will be used without customization or third-party search providers.
4. **Blog or changelog page**: No blog section or changelog page will be added to the docs site. The existing `CHANGELOG.md` remains as-is.
5. **Docs-as-code PR previews**: No per-PR preview deployments of the docs site will be set up.
6. **Content edits beyond migration**: Documentation content will be moved from `README.md` to the docs site largely as-is. Rewriting, expanding, or adding new content is out of scope.

## Design Considerations

- **Theme**: MkDocs Material theme with default settings. No custom color palette, logo overrides, or CSS changes are required for the initial launch.
- **Navigation**: The `mkdocs.yml` `nav` key should define the page order explicitly: Home → Getting Started → Commands Reference → Shell Integration → Configuration.
- **Code blocks**: All CLI examples should use fenced code blocks with `bash` syntax highlighting, consistent with the existing README style.
- **Admonitions**: MkDocs Material supports callout boxes (`!!! note`, `!!! tip`). These may be used where helpful (e.g., noting shell setup requirements), but are not required.

## Repository Standards

- **Commit convention**: Conventional Commits. Documentation setup files use `docs:` type (e.g., `docs: add MkDocs site foundation`). Workflow additions use `ci:` type. README updates use `docs:` type.
- **Branch and PR workflow**: All changes should be made on a feature branch and merged via PR to `main`, following the existing workflow.
- **File placement**: `mkdocs.yml` and `requirements.txt` go at the repository root alongside `go.mod` and `Makefile`. Documentation source files go in `docs/site/`.
- **`.gitignore`**: The MkDocs output directory (`site/`) must be added to `.gitignore` to prevent generated files from being committed.
- **GitHub Actions**: New workflow files follow the existing patterns in `.github/workflows/`. Use `actions/checkout@v4`, `actions/setup-python@v5`, and the official GitHub Pages Actions (`actions/upload-pages-artifact`, `actions/deploy-pages`).

## Technical Considerations

- **`docs_dir` configuration**: Because `docs/` is already used for internal specs (`docs/specs/`), `mkdocs.yml` must explicitly set `docs_dir: docs/site` so MkDocs reads from the correct subfolder.
- **`site_dir` configuration**: MkDocs defaults to outputting the built site to `site/` at the repository root. This is fine for the initial setup; it must be added to `.gitignore`. Alternatively, `site_dir` can be set explicitly in `mkdocs.yml` to make it clear.
- **GitHub Pages source**: The GitHub repository Pages settings must be switched to "GitHub Actions" as the source (not a `gh-pages` branch). This is a one-time manual change in the repository settings under **Settings → Pages → Source**.
- **Deployment action**: Use the official GitHub Pages deployment actions (`actions/upload-pages-artifact@v3`, `actions/deploy-pages@v4`) rather than third-party alternatives like `peaceiris/actions-gh-pages`, for consistency with GitHub's recommended patterns.
- **Python version**: Use Python 3.12 (or latest stable) in the workflow. Pin `mkdocs-material` to a specific version in `requirements.txt` to ensure reproducible builds (e.g., `mkdocs-material==9.x.x`).
- **Manual Verification Steps**: The existing "Manual Verification Steps" section in `README.md` is developer/QA-oriented content. It should be removed from the README as part of Unit 3. It does not need to be migrated to the user-facing docs site.

## Security Considerations

- **No secrets required**: The docs deployment workflow uses the default `GITHUB_TOKEN` with `pages: write` and `id-token: write` permissions. No additional secrets or PATs are needed.
- **No sensitive content in docs**: Documentation pages contain only CLI usage examples and configuration field descriptions. No API keys, tokens, or private paths should appear in any doc page content.

## Success Metrics

1. **Site is live**: `https://daileyo.github.io/gws` loads without errors and displays all four documentation pages after a release tag is pushed.
2. **README is concise**: `README.md` is reduced from ~760 lines to under 150 lines, retaining only developer-facing content.
3. **Automated deployment works**: The `docs.yml` workflow runs to completion (green) on a release tag push with no manual intervention.

## Open Questions

No open questions at this time.
