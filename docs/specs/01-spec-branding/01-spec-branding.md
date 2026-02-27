# 01-spec-branding

## Introduction/Overview

The gws (Git Workspace) project currently has no visual branding — the README uses plain text and badges, and the MkDocs documentation site runs with default Material theme styling. A logo (`gws-logo.png`) has been created featuring a git branch icon on layered diamond shapes with a purple-to-orange gradient and blue accent dots. This spec covers integrating that logo throughout the project's documentation and establishing a cohesive color theme, font selection, and brand identity derived from the logo's color palette.

## Goals

- Integrate the gws logo into the README and MkDocs documentation site for consistent brand identity
- Establish a color theme for the MkDocs site derived from the logo's purple, orange, and blue palette
- Create optimized logo variants (favicon, nav icon, hero image) for fast web performance
- Configure the Material theme with custom fonts (Nunito + JetBrains Mono) to complement the brand
- Support both light and dark color schemes with dark mode as the default

## User Stories

- **As a visitor to the GitHub repository**, I want to see the gws logo and tagline prominently in the README so that I immediately understand the project's identity and purpose.
- **As a documentation reader**, I want the MkDocs site to have a polished, branded appearance with consistent colors and a recognizable logo so that the project feels professional and trustworthy.
- **As a user browsing multiple tabs**, I want a recognizable favicon in my browser tab so that I can quickly identify the gws documentation site.
- **As a user who prefers dark mode**, I want the documentation site to default to dark mode with an option to toggle to light mode so that I can read comfortably in any lighting condition.

## Demoable Units of Work

### Unit 1: Logo Optimization and Asset Preparation

**Purpose:** Create optimized logo variants from the original 1024x1024 PNG so that all downstream branding work has properly sized, web-optimized assets to reference.

**Functional Requirements:**
- The system shall include a favicon version of the logo (e.g., 32x32 or similar standard favicon size) in a web-compatible format
- The system shall include a navigation bar icon version of the logo (approximately 30-40px) optimized for the MkDocs Material theme header
- The system shall include a hero/homepage version of the logo (approximately 200-250px width) for use on the README and MkDocs homepage
- All optimized versions shall be stored in a consistent location within the docs directory (e.g., `docs/site/assets/images/`)
- The original `gws-logo.png` shall remain in the project root unchanged

**Proof Artifacts:**
- File listing: `ls -la docs/site/assets/images/` shows all optimized logo variants with significantly reduced file sizes compared to the 1.5 MB original
- Visual comparison: Opening each optimized file confirms the logo renders correctly at its target size

### Unit 2: README Branding

**Purpose:** Add the logo and tagline to the README so that the GitHub repository landing page has a branded, professional appearance.

**Functional Requirements:**
- The README shall display the gws logo left-aligned next to the project title text
- The logo shall render at approximately 200px width in the README
- The README shall display the tagline "Your Git workspace, simplified" as a subtitle below the logo/title area
- The existing badges and content below shall remain intact and unmodified
- The logo image reference shall use a relative path to the logo file in the repository

**Proof Artifacts:**
- Screenshot: GitHub README rendering shows the logo left-aligned next to the title with the tagline visible below
- Raw markdown: The README.md source shows correct image markup with alt text and sizing

### Unit 3: MkDocs Theme and Color Configuration

**Purpose:** Configure the MkDocs Material theme with branded colors, fonts, and color scheme toggle so that the documentation site reflects the gws brand identity.

**Functional Requirements:**
- The MkDocs theme shall use deep purple as the primary color
- The MkDocs theme shall use orange as the accent color
- The MkDocs theme shall default to dark mode (dark color scheme)
- The MkDocs theme shall include a light/dark mode toggle switch in the header
- The MkDocs theme shall use Nunito as the body/heading font
- The MkDocs theme shall use JetBrains Mono as the code font
- The light mode color scheme shall use the same primary and accent colors adapted for light backgrounds
- Any custom CSS required beyond the Material theme's built-in color options shall be placed in a dedicated overrides/stylesheet file

**Proof Artifacts:**
- Screenshot: MkDocs site in dark mode shows purple header/nav, orange accents on links/buttons, and Nunito body text
- Screenshot: MkDocs site in light mode shows the same color scheme adapted for light backgrounds
- Screenshot: The light/dark toggle switch is visible and functional in the header
- Config review: `mkdocs.yml` shows the theme configuration with color palette, font, and toggle settings

### Unit 4: MkDocs Logo and Homepage Integration

**Purpose:** Add the logo to the MkDocs navigation bar, set the favicon, and add a hero image to the homepage so that every page of the documentation is branded.

**Functional Requirements:**
- The MkDocs site shall display the logo in the top navigation bar (replacing or accompanying the site name)
- The MkDocs site shall use the logo-derived favicon in the browser tab
- The MkDocs homepage (index.md) shall display the logo as a hero image at approximately 200-250px width
- The homepage hero image shall be centered on the page
- The favicon shall be visible in the browser tab when viewing any page of the documentation

**Proof Artifacts:**
- Screenshot: MkDocs site navigation bar shows the gws logo
- Screenshot: Browser tab shows the gws favicon
- Screenshot: MkDocs homepage shows the centered hero logo image at the expected size
- Config review: `mkdocs.yml` shows logo and favicon paths in the theme configuration

## Non-Goals (Out of Scope)

1. **Custom landing page or splash screen**: We are not building a separate landing page beyond the existing MkDocs homepage structure
2. **Logo redesign or modification**: The existing `gws-logo.png` design is final; we are only creating size-optimized variants
3. **Animated or interactive branding elements**: No CSS animations, SVG animations, or JavaScript-based branding effects
4. **Social media assets or Open Graph images**: Social sharing images and meta tags are not included in this spec
5. **Custom MkDocs plugins or extensions**: We are only configuring the existing Material theme, not adding new plugins

## Design Considerations

- The logo uses a gradient from deep purple (#5B2D8E approximate) through warm coral/red to orange (#F47B20 approximate) with small blue accent dots (#4A90D9 approximate). The Material theme's built-in color palette options should be used where possible, with custom CSS only for refinements that the palette cannot achieve.
- The dark mode scheme should ensure sufficient contrast between text, background, and accent colors for readability.
- The light mode scheme should use a white or near-white background with the same purple/orange brand colors adapted for light backgrounds.
- The logo contains transparency (RGBA PNG), so it should render well on both dark and light backgrounds without modification.

## Repository Standards

- Follow Conventional Commits for any commit messages (e.g., `feat: add branding to mkdocs and README`)
- Documentation source files live in `docs/site/`
- MkDocs dependencies are managed in `docs/requirements.txt`
- The project uses `mkdocs-material==9.5.18` — all theme configuration must be compatible with this version
- Media files larger than necessary should be optimized; the `.gitignore` excludes certain media formats but PNGs are tracked

## Technical Considerations

- The Material theme for MkDocs (v9.5.18) supports `palette` configuration with `scheme`, `primary`, `accent`, and `toggle` options directly in `mkdocs.yml`
- Custom CSS can be added via `extra_css` in `mkdocs.yml` pointing to a stylesheet in the docs directory
- Google Fonts (Nunito, JetBrains Mono) are loaded via the Material theme's `font` configuration — no manual font file management needed
- Logo optimization may require image processing tools (e.g., ImageMagick, PIL/Pillow, or similar) — these should be used during development but are not runtime dependencies
- Favicon should be in `.ico` or `.png` format; the Material theme accepts a path to a favicon image file
- The `docs/site/assets/images/` directory may need to be created for storing optimized logo variants

## Security Considerations

No specific security considerations identified. All assets are static image files and theme configuration. No API keys, credentials, or sensitive data are involved.

## Success Metrics

1. **Visual consistency**: The logo appears correctly in the README, MkDocs nav bar, favicon, and homepage — all rendering at their intended sizes without distortion
2. **Performance**: Optimized logo variants are significantly smaller than the 1.5 MB original (target: under 50 KB each for web variants)
3. **Theme completeness**: Both dark and light modes render correctly with branded colors, and the toggle switch works
4. **Font rendering**: Nunito and JetBrains Mono load and display correctly on the MkDocs site
5. **Brand cohesion**: A user visiting both the GitHub README and the MkDocs site perceives them as part of the same project identity

## Open Questions

No open questions at this time.
