# Task 2.0 Proof Artifacts — GitHub Actions Deployment Workflow

## Workflow File

File: `.github/workflows/docs.yml`

```yaml
name: Deploy Docs

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    name: Deploy Documentation
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'

      - name: Install MkDocs Material
        run: pip install -r requirements.txt

      - name: Build documentation
        run: mkdocs build

      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: site/

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## GitHub Pages Repository Setting

> **Manual action required**: In the `daileyo/gws` repository, navigate to
> **Settings → Pages → Build and deployment → Source** and select **"GitHub Actions"**.
> This one-time change enables the workflow to deploy to GitHub Pages.

## Screenshots

> **Action required**: Capture the following screenshots after the first release tag triggers the workflow:

- **Screenshot 1**: `.github/workflows/docs.yml` file shown on the `daileyo/gws` GitHub repository page (Actions tab or Code view)
- **Screenshot 2**: GitHub Actions run for `Deploy Docs` showing all steps with green checkmarks
- **Screenshot 3**: `https://daileyo.github.io/gws` rendered in a browser showing the MkDocs site home page
