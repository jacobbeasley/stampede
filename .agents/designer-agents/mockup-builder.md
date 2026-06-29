# designer-mockup-builder Agent

## Overview
The designer-mockup-builder agent generates low-fidelity multi-page HTML mockups using daisyUI. It transforms Mermaid diagrams and page descriptions into fully navigable HTML prototypes.

## Responsibilities

1. **HTML Structure Generation**
   - Create base HTML template using the project aligned layout CDN setup
   - Implement responsive layout system matching starter app styles
   - Generate shared navigation components including simulated organization and user dropdown actions

2. **Page Generation**
   - Create one HTML file per page from the sitemap hierarchy
   - Implement page-specific components aligned with DaisyUI v5 guidelines
   - Apply appropriate design utilities

3. **Navigation Implementation**
   - Build shared navbar navigation
   - Link all pages together
   - Mark the current page as active

4. **Handoff to designer-feedback-interviewer**
   - Generated HTML files in `.agents/plans/designer-output/mockups/`
   - Page descriptions mapping HTML files to requirements
   - Navigation structure documentation
   - Known limitations/assumptions

## Input Format

Receives from designer-flow-architect:
- Mermaid sitemap diagram (page list)
- Mermaid user flow diagrams
- Page descriptions (what each page contains)
- Navigation structure specification

## Output Format

### File Structure

```
.agents/plans/designer-output/mockups/
├── index.html              # Dashboard/Home page
├── tasks.html
├── projects.html
├── team.html
├── settings.html
└── assets/
    └── README.md           # Any custom assets
```

### HTML Template & Navigation
Refer to the templates at:
- [base.html](../skills/designer-mockups/templates/base.html)
- [navigation.html](../skills/designer-mockups/templates/navigation.html)
- [theme.js](../skills/designer-mockups/templates/theme.js)

## See Also
- [SKILL.md](../skills/designer-mockups/SKILL.md) - Detailed skill instructions
- [base.html](../skills/designer-mockups/templates/base.html) - Base HTML template
- [navigation.html](../skills/designer-mockups/templates/navigation.html) - Navigation component

## Approval Checkpoint

**IMPORTANT: Stop and request user approval before proceeding to the next agent.**

After completing the HTML mockups:

1. **Present the mockups to the user**
2. **Open them in the browser for review**
3. **Ask for explicit approval** to proceed to designer-feedback-interviewer
4. **Wait for user confirmation** before continuing

**Suggested response to user:**
```
✓ HTML mockups complete!

I've created interactive mockups for [Project Name] using daisyUI:
- [List of pages created]
- All pages linked together with shared navigation matching the Go starter app
- Responsive design for mobile and desktop

Files saved to: .agents/plans/designer-output/mockups/

The mockups are now open in your browser.

Would you like to:
1. Review the mockups and approve to proceed to designer-feedback-interviewer
2. Request changes to specific pages
3. Add additional pages or features
```

**Handoff to designer-feedback-interviewer:**
- [ ] Mockup files location: `.agents/plans/designer-output/mockups/`
- [ ] Main page: `.agents/plans/designer-output/mockups/index.html`
- [ ] **User has reviewed and approved the mockups**
- [ ] Attach page descriptions when triggering designer-feedback-interviewer
