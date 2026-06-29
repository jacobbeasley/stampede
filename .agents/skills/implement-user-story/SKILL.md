---
name: implement-user-story
description: Use this skill whenever a user provides a developer story, user story, or requests a new feature. This skill orchestrates the full development lifecycle in a Buffalo + Svelte stack, including planning, code generation, database scaffolding, UI/UX aesthetics, testing, validation, and documentation.
---

# Implement User Story

This skill guides you through implementing a new feature or user story end-to-end. You must follow these 7 phases systematically. 

## 1. Context Gathering & Research
Before writing any code or plans, understand the landscape:
- **Review Specifications**: Read relevant files in `.agents/specifications/` and check the `.agents/knowledge/` directory for any relevant Knowledge Items (KIs).
- **Assess Existing Codebase**: Search `models/`, `actions/`, and the Svelte frontend (`assets/` or `src/`) for existing code that could be reused or extended. Ensure you aren't duplicating work.

## 2. Artifact Management & Planning
Do not start coding until a clear plan is made.
- **Story Folder**: Create a directory for this story at `.agents/specifications/stories/<storyid>/` (replace `<storyid>` with the assigned ID from the backlog).
- **Story Spec**: Create `story.md` in that folder containing the full description of the story.
- **Create Implementation Plan**: Generate a `plan.md` file in the story folder outlining the technical design.
- **Data Model**: Determine necessary database tables, columns, and relationships (Pop ORM).
- **Pages & APIs**: Identify what new Svelte pages or Buffalo API endpoints are needed.
- **Scaffolding Strategy**: Define the exact `buffalo generate` commands you will use (e.g., `buffalo pop.g model User email:string`, `buffalo generate resource Users`).
- **UI/UX & Layout**: Define field requirements, aesthetics, and responsiveness strategies. You MUST adhere to "rich aesthetics" using Tailwind CSS and DaisyUI, ensuring it works on mobile and desktop viewports.
- **Acceptance Criteria**: Clearly list conditions that must be met.
- **Task Tracker**: Generate a `task.md` file in the story folder to track your progress step-by-step.
- **Approval**: Stop and wait for the user to approve the implementation plan before proceeding to Execution.
- **Backlog Entry**: If the story does not already exist in `.agents/specifications/stories/backlog.md`, append it to the bottom of the backlog list.
- **Backlog Update**: Once the user approves and before writing any code, update the story's status in `.agents/specifications/stories/backlog.md` to **`In Progress`**.

## 3. Execution Phase (Writing Code)
Once approved, execute the plan:
- **Database**: Run your planned scaffolding commands using Buffalo/Soda CLIs (e.g., `buffalo pop generate fizz` or `soda generate fizz`). Run migrations (`buffalo pop migrate` or `soda migrate`).
- **Backend**: Implement Buffalo actions and routes in `actions/`. Keep controllers minimal; push complex logic to models or service layers. If `buffalo dev` fails to start in an environment, use `go run cmd/app/main.go`.
- **Frontend**: Create small, reusable Svelte components using the design system. Integrate with the backend using the standard `fetch` API. Handle loading and error states gracefully. If submitting forms against Buffalo from Svelte or external tools, extract the CSRF token (`authenticity_token`) from the `<meta name="csrf-token">` tag.
- **Iterative Commits**: Commit your changes to Git logically after each working piece is completed. Write descriptive commit messages.

## 4. Code Review & Refinement
Do a self-review of the code you just wrote.
- **Formatting**: Run `gofmt -w .` or equivalent to ensure Go code is formatted properly. Check Svelte/JS formatting.
- **Aesthetics Check**: Does the UI look generic? If so, you have failed the aesthetics requirement. Add modern typography, harmonious colors (via Tailwind/DaisyUI), hover effects, and micro-animations.

## 5. Testing & Validation Phase
Validate your work.
- **Backlog Update**: At the start of this phase, update the story's status in `.agents/specifications/stories/backlog.md` to **`Internal Testing`**.
- **Test Database Setup**: Ensure you run `buffalo pop create -a` and `buffalo pop migrate` (or equivalent `soda` commands) in the test environment before running tests, as tests require an active database connection.
- **Functional Testing**: Test the UI in a browser using the `browser_subagent` if available, or instruct the user to do so. Verify behavior on both desktop and mobile viewports.
- **Acceptance Validation**: Cross-reference the original Acceptance Criteria to ensure all points are satisfied.
- **Unit Tests**: Write Go unit tests for backend logic.

## 6. Regression Analysis Phase
Ensure you haven't broken anything else.
- Contextually analyze your changes against the existing codebase.
- Did you modify global state or shared components? Ensure those changes don't negatively impact other pages.
- Run any existing unit tests (`buffalo test`).

## 7. Documentation Phase
Wrap up the user story.
- Generate or update documentation where relevant.
- Summarize the changes in a `walkthrough.md` artifact.
- Provide clear, concise testing instructions for the user, detailing how they can manually verify the feature.
- **Spec Update**: Update `.agents/specifications/application_spec.md` (if not already done) so that future work factors in the features, data models, or routing that have been implemented.
- **Backlog Update**: Update the story's status in `.agents/specifications/stories/backlog.md` to **`Complete`**.
