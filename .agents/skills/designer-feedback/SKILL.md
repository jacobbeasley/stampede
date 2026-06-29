# Designer Feedback Interview Skill

## Overview
This skill provides guidance for conducting structured Q&A sessions to gather user feedback on HTML mockups.

## When to Use
- After mockups are generated and ready for review
- Before finalizing design decisions
- When validating assumptions with stakeholders

## Process

### Phase 1: Structured Interview
**Purpose:** Quantitative assessment

1. Rate each page on Clarity, Usability, Aesthetics (1-5 scale)
2. Rate each user flow on Completeness, Efficiency (1-5 scale)
3. Mark pages as: "Ready" / "Needs Work" / "Major Redesign"

### Phase 2: Open-Ended Discussion
**Purpose:** Qualitative insights

1. Ask "What works well?"
2. Ask "What would you change?"
3. Ask "What's missing?"
4. Ask "What's confusing?"
5. Ask prioritization questions

### Phase 3: Synthesis
**Purpose:** Prepare for mockup-builder

1. Compile ratings into summary table
2. Extract verbatim quotes for feedback
3. Prioritize changes (high/medium/low)
4. Document new requirements discovered
5. Format output for mockup-builder

## Output
Feedback report in `.agents/plans/designer-output/feedback-reports/[project]-feedback.md`

## Best Practices
- Always start with structured assessment before open-ended
- Capture ratings even if user skips (mark as N/A)
- Quote user verbatim when possible
- Flag contradictions (e.g., "confusing" but rated "5")
- Prioritize by impact vs. effort
- Document assumptions from ambiguous feedback
