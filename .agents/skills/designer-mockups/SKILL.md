# Designer Mockups Skill

## Overview
This skill provides guidance for generating low-fidelity HTML mockups using daisyUI v5 with Tailwind CSS v4 via CDN, aligned with the project's layout templates.

## When to Use
- After user flows and sitemaps are approved
- Before the feedback interview process
- When creating navigable HTML prototypes

## Process

### 1. Base Template Setup
- Choose the appropriate base layout template:
  - **`templates/base.html`**: For authenticated layout mockups (includes organization switchers, user profiles, and active app navigation).
  - **`templates/base-loggedout.html`**: For guest pages (landing pages, login forms, sign-up flows).
- Create HTML5 boilerplate with daisyUI CDN and themes
- Include Tailwind CSS v4 browser CDN script
- Add standard responsive viewport meta tags
- Set default theme to `dark` to match the project's styling guidelines
- Add Google Fonts (Inter) and custom body font mappings

### 2. Navigation Integration
- Create shared navbar with all sitemap pages
- Mark current page as active
- Include user avatar dropdown with settings/logout options
- Integrate dynamic theme selector (using `theme.js`)
- Ensure mobile-responsive hamburger collapse layout

### 3. Page Generation
- Create one HTML file per page from the sitemap
- Implement components per page description
- Use semantic HTML elements
- Apply appropriate daisyUI classes
- **Generate launcher index page (`index.html`)**: Replicate the structure of `templates/index.html` at the output root. This page lists all other generated mockup pages with a brief description and launcher link, helping stakeholders review mockups systematically.

### 4. Cross-Linking
- Ensure all pages link together via the shared navigation
- Verify all links function correctly
- Add breadcrumbs for deeper nested pages if needed

## Output
- **Launcher index page** (`.agents/plans/designer-output/mockups/index.html`) acting as the mockups directory hub.
- Multi-page HTML mockups saved in `.agents/plans/designer-output/mockups/`
- All pages navigable via the shared navigation component

## daisyUI Quick Reference

### Layout
- `navbar` - Sticky top navbar
- `hero` - Hero sections
- `footer` - Footers matching starter app
- `container` - Centered wrapper

### Forms & Controls
- `input` - Text inputs
- `select` - Dropdown selects
- `checkbox` / `radio` - Option toggles
- `textarea` - Multi-line inputs

### Feedback & Loading
- `alert` - Feedback alerts
- `badge` - Status indicators
- `skeleton` - Placeholder/loading states
- `progress` - Progress bars

### Data Display
- `table` - Data lists
- `stat` - Statistics dashboard widgets
- `avatar` - User avatars
- `card` - Component blocks

### Actions
- `btn` - Buttons (primary, secondary, outline)
- `dropdown` - Interactive dropdown menus
- `modal` - Toggleable overlay dialogs
- `menu` - Vertical or horizontal menu layouts

## Best Practices
- Keep it low-fidelity but aesthetically aligned with the starter app
- Use realistic placeholder data (no `lorem ipsum`)
- Retain the theme picker in all mockups so stakeholders can test theme color palettes
- Include mock CSRF headers in case page forms simulate API calls
