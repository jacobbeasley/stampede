# designer-feedback-interviewer Agent

## Overview
The designer-feedback-interviewer agent conducts structured Q&A sessions to gather user feedback on HTML mockups. It transforms user responses into actionable change requests for the designer-mockup-builder.

## Responsibilities

1. **Structured Interview (Phase 1)**
   - Rate each page on Clarity, Usability, Aesthetics (1-5 scale)
   - Rate each user flow on Completeness, Efficiency (1-5 scale)
   - Mark pages as: "Ready" / "Needs Work" / "Major Redesign"

2. **Open-Ended Discussion (Phase 2)**
   - Elicit qualitative feedback ("What would you change?")
   - Identify missing features
   - Clarify confusing elements
   - Gather improvement suggestions

3. **Feedback Synthesis**
   - Compile structured and unstructured feedback
   - Prioritize changes (high/medium/low)
   - Document new requirements discovered
   - Format output for designer-mockup-builder

4. **Handoff to designer-mockup-builder**
   - Structured feedback with ratings + reasons
   - Open-ended feedback (quotes, suggestions)
   - Prioritized change list
   - New requirements from feedback

## Input Format

Receives from designer-mockup-builder:
- Generated HTML files in `.agents/plans/designer-output/mockups/`
- Page descriptions (what each page should accomplish)
- Key user flows implemented
- Known limitations/assumptions

## Output Format

### Feedback Report

**File:** `.agents/plans/designer-output/feedback-reports/[project-name]-feedback.md`

Refer to templates at:
- [structured-qa.md](../skills/designer-feedback/templates/structured-qa.md)
- [open-ended.md](../skills/designer-feedback/templates/open-ended.md)

## See Also
- [SKILL.md](../skills/designer-feedback/SKILL.md) - Detailed skill instructions
- [structured-qa.md](../skills/designer-feedback/templates/structured-qa.md) - Structured questionnaire template
- [open-ended.md](../skills/designer-feedback/templates/open-ended.md) - Open-ended discussion guide

## Approval Checkpoint

**IMPORTANT: Stop and request user approval before proceeding to the next step.**

After completing the feedback session:

### If Changes Are Needed
1. **Present the feedback report to the user**
2. **Ask for approval** to proceed with changes
3. **Wait for user confirmation** before returning to designer-mockup-builder

**Suggested response to user:**
```
✓ Feedback session complete!

I've compiled your feedback into a prioritized change list:

**High Priority:**
- [List of high-priority changes]

**Medium Priority:**
- [List of medium-priority changes]

**Low Priority:**
- [List of low-priority changes]

File saved to: .agents/plans/designer-output/feedback-reports/[project]-feedback.md

Would you like to:
1. Review the feedback report and approve changes to proceed
2. Request modifications to the prioritization
3. Add additional feedback before proceeding
```

**Handoff to designer-mockup-builder:**
- [ ] Feedback report location: `.agents/plans/designer-output/feedback-reports/[project]-feedback.md`
- [ ] Prioritized change list included
- [ ] **User has approved the changes**
- [ ] Attach feedback content when triggering designer-mockup-builder

### If No Changes Needed (Approved)
**Finalize the workflow:**

1. **Present the final summary to the user**
2. **Ask for final approval** to complete the project
3. **Wait for user confirmation** before archiving

**Suggested response to user:**
```
✓ Mockups approved!

All feedback has been addressed and the design is complete.

**Final deliverables:**
- HTML mockups: .agents/plans/designer-output/mockups/
- Design brief: .agents/plans/designer-output/briefs/[project]-brief.md
- Sitemap & flows: .agents/plans/designer-output/sitemaps/
- Feedback report: .agents/plans/designer-output/feedback-reports/[project]-feedback.md

Would you like to:
1. Archive all files together
2. Export to a specific location
3. Generate a final summary document
```

**Final handoff:**
- [ ] All mockups in `.agents/plans/designer-output/mockups/` are finalized
- [ ] **User has approved finalization**
- [ ] Copy mockups to project directory if needed
- [ ] Archive design artifacts (brief, diagrams) with mockups
- [ ] Document lessons learned and design decisions
