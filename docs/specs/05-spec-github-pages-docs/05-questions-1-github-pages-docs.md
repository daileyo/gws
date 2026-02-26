# 05 Questions Round 1 - GitHub Pages Docs

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Static Site Generator

What tool should be used to build and publish the documentation site on GitHub Pages?

- [X] (A) **MkDocs with Material theme** — Python-based, widely used for CLI/tool docs, beautiful out-of-the-box design, easy navigation and search. Requires a `mkdocs.yml` config file and a `docs/` folder of markdown pages. Needs a GitHub Actions step to build and deploy.
- [ ] (B) **Jekyll (GitHub Pages native)** — No build step required; GitHub Pages renders Jekyll sites automatically. Limited theme options compared to MkDocs, but zero external tooling.
- [ ] (C) **Docusaurus** — React/Node.js-based, feature-rich (versioning, blog, search). More setup overhead than MkDocs or Jekyll.
- [ ] (D) **Plain Markdown (no generator)** — Just a structured set of `.md` files in a `docs/` folder rendered by GitHub's default markdown view. No custom site, no deploy workflow needed.
- [ ] (E) Other (describe)

## 2. Documentation Site Structure

Which pages should the documentation site include? (Select all that apply)

- [X] (A) **Getting Started** — Installation methods (Homebrew, build from source), first-time setup (`--init`)
- [X] (B) **Commands Reference** — All flags and commands with examples (`--list`, `--status`, `--add-tag`, `--refresh`, `--go`, etc.)
- [X] (C) **Shell Integration** — Setup and usage of the `gws` shell function and `cdgws`/`gcd` aliases
- [X] (D) **Configuration** — Config file location, structure, and fields (`~/.gws/config.json`)
- [ ] (E) Other (describe)

## 3. README Simplified Content

After moving usage docs to the site, what should remain in the README?

- [ ] (A) **Minimal**: Badges, one-line description, quick install snippet, link to docs site, Development section (build commands, running tests, project structure), Contributing section
- [X] (B) **Brief overview**: Badges, description, features bullet list, quick install snippet, link to docs site, Development section, Contributing section
- [ ] (C) **Other (describe)**

## 4. Contributing Section Content

The current README says "Contributions are welcome!" and has a full workflow guide. What should the updated Contributing section say?

- [ ] (A) **Closed with roadmap note**: State that the project is not yet open for external contributions, and that the roadmap will be updated when it is. Keep commit message conventions for context.
- [ ] (B) **Closed, brief**: Just a short note that the project is not currently accepting contributions, with no further detail.
- [X] (C) **Closed with contact**: Note it's closed, but invite interested contributors to open an issue or reach out.
- [ ] (D) Other (describe)

## 5. GitHub Pages Deployment Trigger

When should the docs site be rebuilt and deployed?

- [ ] (A) **On every push to `main`** — Docs stay in sync with the latest code at all times
- [X] (B) **On release tags only** — Docs update only when a new version is released; matches user-facing releases
- [ ] (C) **Manual trigger only** — Docs deploy only when explicitly triggered via GitHub Actions `workflow_dispatch`
- [ ] (D) Other (describe)

## 6. Docs Source Location

Where should the documentation source files (markdown pages for the site) live in the repository?

- [ ] (A) **`docs/`** — Top-level `docs/` folder (the `docs/specs/` path would become `docs/specs/` inside it, staying as-is)
- [ ] (B) **`site/`** — A dedicated `site/` folder separate from the existing `docs/specs/` directory
- [X] (C) **`docs/site/`** — A subfolder within `docs/` to keep specs and site source clearly separated
- [ ] (D) Other (describe)

## 7. Custom Domain

Should the docs site use a custom domain, or the default GitHub Pages URL?

- [X] (A) **Default GitHub Pages URL** — `https://daileyo.github.io/gws` (no DNS changes needed)
- [ ] (B) **Custom domain** — Specify the domain in notes below
- [ ] (C) Other (describe)
