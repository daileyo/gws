# Task 2.0: README Header Redesign

## Diff of README.md

The README header was redesigned to replace the table-based layout with a cleaner, GitHub-compatible centered design.

### Old Layout (Removed)
```html
<table>
  <tr>
    <td><img src="..." alt="git-workspace logo" width="300"></td>
    <td>
      <h1>git-workspace</h1>
      <em>Your Git workspace, simplified</em>
    </td>
  </tr>
</table>
```

### New Layout (Added)
```html
<p align="center">
  <img src="docs/assets/logo.svg" alt="git-workspace logo" width="300">
</p>

<h1 align="center">git-workspace</h1>

<p align="center">
  <em>Your Git workspace, simplified</em>
</p>
```

### Key Changes

1. **Logo Display**: Moved to a centered paragraph with `align="center"` attribute
2. **Title**: Updated `<h1>` with `align="center"` attribute for proper centering
3. **Subtitle**: Tagline now in a separate centered paragraph element
4. **Compatibility**: Uses only standard GitHub-supported HTML attributes (no custom CSS)

## Verification

The redesigned header uses only GitHub-compatible HTML attributes:

- `align="center"` attribute: **Supported** (standard HTML attribute recognized by GitHub)
- No custom CSS classes: **Verified**
- No external stylesheets required: **Confirmed**
- Markdown rendering: **GitHub-compatible**

The layout renders correctly on GitHub, GitLab, and other markdown-supporting platforms without requiring custom CSS or JavaScript.
