# Task 2.0 Proof Artifacts — README Branding

## Diff Output

```diff
$ git diff README.md

-# git-workspace
+<table><tr>
+<td><img src="docs/site/assets/images/gws-logo-hero.png" alt="gws logo" width="200"></td>
+<td><h1>git-workspace</h1></td>
+</tr></table>
+
+*Your Git workspace, simplified*

 [![CI](https://github.com/daileyo/gws/actions/workflows/ci.yml/badge.svg)]...
```

Only the heading was replaced. All badges and content below remain unchanged.

## Raw Markdown Verification

```html
<table><tr>
<td><img src="docs/site/assets/images/gws-logo-hero.png" alt="gws logo" width="200"></td>
<td><h1>git-workspace</h1></td>
</tr></table>

*Your Git workspace, simplified*
```

- Image path: `docs/site/assets/images/gws-logo-hero.png` (relative path)
- Alt text: `alt="gws logo"`
- Width: `width="200"`
- Tagline: Italicized, placed between logo/title and badges

## Verification Checklist

- [x] Logo left-aligned next to title via HTML table
- [x] Logo rendered at ~200px width
- [x] Tagline "Your Git workspace, simplified" present as italicized subtitle
- [x] All 5 badges (CI, Go Report Card, Snyk, Release, License) intact
- [x] All content below badges unchanged
- [x] Image uses relative path with alt text and width attribute
