# 01-tasks-branding

## Relevant Files

- `gws-logo.png` - Original logo file (1024x1024, 1.5 MB) at project root. Source for all optimized variants. Do not modify.
- `docs/site/assets/images/gws-logo-favicon.png` - Favicon variant (~32x32). To be created.
- `docs/site/assets/images/gws-logo-nav.png` - Navigation bar icon variant (~40px). To be created.
- `docs/site/assets/images/gws-logo-hero.png` - Hero/homepage image variant (~250px). To be created.
- `README.md` - Project README. Will be modified to add logo and tagline at the top.
- `mkdocs.yml` - MkDocs configuration. Will be modified to add theme palette, fonts, logo, favicon, toggle, and extra_css.
- `docs/site/index.md` - MkDocs homepage. Will be modified to add centered hero logo image.
- `docs/site/assets/css/custom.css` - Custom CSS overrides for branding refinements. To be created.

### Notes

- Image optimization requires a tool like ImageMagick (`convert`), Pillow (`python3 -c "from PIL import Image; ..."`), or similar. These are dev-time tools only.
- Documentation source files live in `docs/site/` (configured as `docs_dir` in mkdocs.yml).
- The project uses `mkdocs-material==9.5.18` — all configuration must be compatible with this version.
- PNG files are tracked by git (not in `.gitignore`).
- Follow Conventional Commits for any commit messages.
- Preview changes locally with `python -m mkdocs serve` before finalizing.

## Tasks

### [x] 1.0 Logo Optimization and Asset Preparation

Create optimized logo variants from the original 1024x1024 PNG (1.5 MB) for use as favicon, navigation bar icon, and hero/homepage image. Store all variants in `docs/site/assets/images/` while keeping the original `gws-logo.png` unchanged at the project root.

#### 1.0 Proof Artifact(s)

- CLI: `ls -la docs/site/assets/images/` shows favicon (~32x32), nav icon (~40px), and hero image (~250px) variants, each under 50 KB
- Visual: Opening each variant confirms the logo renders correctly without distortion at its target size

#### 1.0 Tasks

- [x] 1.1 Create the `docs/site/assets/images/` directory
- [x] 1.2 Generate a favicon variant (`gws-logo-favicon.png`) at 32x32 pixels from the original `gws-logo.png`. Use an image processing tool (e.g., ImageMagick: `convert gws-logo.png -resize 32x32 docs/site/assets/images/gws-logo-favicon.png`). Verify the file is under 50 KB.
- [x] 1.3 Generate a navigation bar icon variant (`gws-logo-nav.png`) at 40x40 pixels from the original. Verify the file is under 50 KB.
- [x] 1.4 Generate a hero image variant (`gws-logo-hero.png`) at 250x250 pixels from the original. Verify the file is under 50 KB.
- [x] 1.5 Visually inspect each generated variant to confirm the logo renders correctly without distortion or artifacts at its target size.
- [x] 1.6 Run `ls -la docs/site/assets/images/` to confirm all three variants exist with appropriate file sizes.

### [x] 2.0 README Branding

Add the gws logo and tagline to the README. The logo should be left-aligned next to the project title at ~200px width, with the tagline "Your Git workspace, simplified" displayed as a subtitle. All existing badges and content must remain intact.

#### 2.0 Proof Artifact(s)

- Screenshot: GitHub README rendering shows logo left-aligned next to the title with tagline visible below
- Diff: `git diff README.md` shows the added image markup with alt text, sizing, and tagline — no other content changed

#### 2.0 Tasks

- [x] 2.1 Replace the current `# git-workspace` heading in `README.md` with a layout that places the logo (~200px width) left-aligned next to the title text. Use an HTML table or inline HTML `<img>` tag with `width="200"` and appropriate `alt` text (e.g., `alt="gws logo"`). Reference the hero image variant at `docs/site/assets/images/gws-logo-hero.png`.
- [x] 2.2 Add the tagline *"Your Git workspace, simplified"* as an italicized subtitle line immediately below the logo/title area, above the badges.
- [x] 2.3 Verify that all existing badges (CI, Go Report Card, Snyk, Release, License) and all content below them remain unchanged.
- [x] 2.4 Review the raw markdown to confirm the image uses a relative path, has alt text, and specifies width.

### [x] 3.0 MkDocs Theme and Color Configuration

Configure the Material theme in `mkdocs.yml` with deep purple primary, orange accent, dark mode default, light/dark toggle, Nunito body font, and JetBrains Mono code font. Add custom CSS if needed for refinements beyond the built-in palette.

#### 3.0 Proof Artifact(s)

- Screenshot: MkDocs site in dark mode shows purple header/nav, orange accent links, Nunito body text, and JetBrains Mono code blocks
- Screenshot: MkDocs site in light mode shows adapted color scheme on white background
- Screenshot: Light/dark toggle switch visible and functional in the header
- Config: `mkdocs.yml` contains palette, font, and toggle configuration

#### 3.0 Tasks

- [x] 3.1 Add `palette` configuration under `theme` in `mkdocs.yml` with two scheme entries: (1) **Dark mode (default)**: `scheme: slate`, `primary: deep purple`, `accent: orange`, with a toggle icon (`material/brightness-7`) and name "Switch to light mode"; (2) **Light mode**: `scheme: default`, `primary: deep purple`, `accent: orange`, with a toggle icon (`material/brightness-4`) and name "Switch to dark mode". Dark mode must be listed first so it is the default.
- [x] 3.2 Add `font` configuration under `theme` in `mkdocs.yml`: `text: Nunito`, `code: JetBrains Mono`.
- [x] 3.3 Create the directory `docs/site/assets/css/` and create a `custom.css` file for any branding refinements (e.g., fine-tuning link hover colors, header gradient, or accent shades that the built-in palette doesn't cover). If no custom CSS is needed beyond the built-in palette, create the file with a comment placeholder for future use.
- [x] 3.4 Add `extra_css` to `mkdocs.yml` pointing to `assets/css/custom.css`.
- [x] 3.5 Run `python -m mkdocs serve` locally and verify: (a) dark mode loads by default, (b) the toggle switch appears and works, (c) purple header/nav and orange accents are visible, (d) Nunito renders for body text, (e) JetBrains Mono renders for code blocks.

### [x] 4.0 MkDocs Logo and Homepage Integration

Add the logo to the MkDocs navigation bar, set the favicon, and add a centered hero image (~200-250px) to the homepage (`docs/site/index.md`). Every page should show the nav logo and favicon.

#### 4.0 Proof Artifact(s)

- Screenshot: MkDocs navigation bar displays the gws logo
- Screenshot: Browser tab shows the gws favicon
- Screenshot: MkDocs homepage shows centered hero logo at expected size
- Config: `mkdocs.yml` contains logo and favicon paths in theme configuration

#### 4.0 Tasks

- [x] 4.1 Add `logo` and `favicon` under `theme` in `mkdocs.yml`. Set `logo` to `assets/images/gws-logo-nav.png` and `favicon` to `assets/images/gws-logo-favicon.png`. (Paths are relative to `docs_dir`.)
- [x] 4.2 Update `docs/site/index.md` to add the hero logo image centered on the page. Use an HTML block like `<p align="center"><img src="assets/images/gws-logo-hero.png" alt="gws logo" width="250"></p>` at the top of the page content, above the existing text.
- [x] 4.3 Run `python -m mkdocs serve` locally and verify: (a) the logo appears in the navigation bar on every page, (b) the favicon is visible in the browser tab, (c) the homepage displays the centered hero logo at ~250px width.
- [x] 4.4 Navigate to at least two different documentation pages to confirm the nav logo and favicon are consistent across the site.
