---
name: feature-orchestrator
description: Orchestrates the full lifecycle of a feature build. Manages subagents, human approvals, and finalizing specifications and documentation.
---

# Feature Orchestrator Skill

Process:

1. Invoke the `codebase-researcher` subagent. Wait for findings.

2. Invoke the `story-writer` subagent. Pass findings. Wait for the `story.md` and backlog creation.

3. **ASK HUMAN**: "Does `story.md` match what you want? Reply 'approved' to continue, describe what to change, or reply 'reject' to stop the chain."
   - If changes requested, invoke `story-writer` again with feedback. Repeat until approved or rejected.

4. Invoke the `spec-writer` subagent. Wait for `plan.md` and `task.md`.

5. **ASK HUMAN**: "Any design red flags in `plan.md`? Reply 'approved' to continue, describe what to change, or reply 'reject' to stop the chain."
   - If changes requested, invoke `spec-writer` again with feedback. Repeat until approved or rejected.

6. Invoke the `fullstack-builder` subagent. Wait for the implementation and its summary.

7. Invoke the `test-verifier` subagent. Wait for the acceptance tests and the verifier's report.

8. Invoke the `implementation-validator` subagent. Wait for findings.

9. If the validator reports critical findings, route them back to the `fullstack-builder`. Then re-run `test-verifier` and the validator.

10. **Documentation Phase (Finalization)**:
    - Update `.agents/specifications/application_spec.md` with the new data models, routing, and features implemented.
    - Generate `.agents/specifications/stories/<storyid>/walkthrough.md` summarizing the changes and providing testing instructions.
    - Update the story's status in `.agents/specifications/stories/backlog.md` to **`Complete`**.

11. Show the validator findings and walkthrough to the user. **ASK HUMAN**: "Ready to open the PR?"

## Rules:
- Never skip the human approval points.
- Never invoke `test-verifier` before the builder has finished.
- Never invoke the validator before the chain has produced some implementation and the verifier has run.
- Do not mark the story as `Complete` until the implementation is fully validated.
