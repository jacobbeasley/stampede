---
name: landing-page
description: Helps create mockups of single-page landing pages with interactive theme pickers, web-researched options, variations, feedback cycles, and export functions.
---

# Landing Page Mockup Generator

This skill guides the process of researching, planning, and creating interactive, high-fidelity mockups of single-page landing pages. It enforces a 5-phase workflow using DaisyUI, Tailwind CSS v4, and custom features like theme pickers, variant generation, feedback iterations, and code conversions.

## Phase 1: Discovery & Requirements
Always begin by interacting with the user to define the core requirements. Ask questions to define:
1. **Goal/Purpose:** What is the primary purpose of this landing page? (e.g., lead generation, product launch, webinar registration).
2. **Call to Action (CTA):** What is the desired action? (e.g., form submission, button click, purchase, sign-up).
3. **Key Messaging:** Who is the target audience and what are the main value propositions?

Wait for the user's response before proceeding to Phase 2.

## Phase 2: Research & Brainstorming
1. **Web Research:** Use your search tools to find layout trends, copy ideas, and designs of high-converting landing pages in the user's target domain.
2. **Propose Variations:** Present 3 distinct design variations based on the requirements (e.g., Minimalist/Modern, Bold/Creative, Corporate/Professional). Describe the layout structure, color schemes, and tone for each.
3. **Wait for Approval:** Ask the user to approve the concepts before generating any code.

## Phase 3: Generation
Once the concepts are approved, generate the 3 variations as functional HTML files.
1. **Output Directory:** Save the generated mockups into `.agents/plans/landing-page-output/<project-name>/var-[1/2/3]/index.html`.
2. **Base Setup:** Use the template located at `resources/base.html` which includes DaisyUI, Tailwind CSS v4 (via CDN), Google Fonts (Inter/Outfit), and SEO meta tags placeholders.
3. **Theme Picker:** Ensure the interactive theme picker UI from `base.html` is retained and the `resources/theme.js` script is copied into each variation's directory alongside the HTML file.
4. **Design Rules:**
   - Enforce SEO best practices (fill in Open Graph tags and meta descriptions).
   - Ensure the design is mobile-first and responsive.
   - Use semantic HTML5 elements.
   - Ensure you design full sections (Hero, Features, Testimonials, Footer).

## Phase 4: Review & Iteration
1. Present the generated file paths to the user.
2. Ask the user to open them in their browser, review the variations, and pick their favorite.
3. Solicit feedback and apply requested iterations to the selected version.

## Phase 5: Export & Tracking
Once the user is satisfied with the chosen variation, offer the following finalization steps:
1. **PHP Conversion:** Offer to translate the static HTML form into a working PHP script. If requested, use the `resources/form-handler.php` template to build the backend and rename the mockup to `index.php`.
2. **Analytics Integration:** Ask the user if they'd like to include an analytics tracking package to measure campaign performance (e.g., PostHog or Microsoft Clarity). If they say yes, inject the appropriate snippet from `resources/analytics.md` into the `<head>` of the final mockup.
