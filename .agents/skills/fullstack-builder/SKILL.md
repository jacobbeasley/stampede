---
name: fullstack-builder
description: Implements the full stack of a feature: Buffalo API routes, Pop ORM models/migrations, and Svelte 5 frontend components. Reads AGENTS.md, updates task tracking, and explicitly applies indexing and strict component structures.
---

# Fullstack Builder Skill

You are the fullstack implementation worker for this project. Your job is to implement both the backend and frontend of the feature described in `plan.md`.

Before you edit anything:
1. Ensure the story status in `.agents/specifications/stories/backlog.md` is set to **`In Progress`**.
2. Read `AGENTS.md` and `.agents/knowledge/`.
3. Review `.agents/specifications/stories/<storyid>/plan.md` and `task.md`.

## Implementation rules:
- **Backend (Go/Buffalo & Database)**:
  - Edit files in `actions/`, `models/`, `workers/`, `migrations/`.
  - **Immutability**: DO NOT modify existing migrations that have already been committed. Always generate new migrations (`buffalo pop generate fizz <name>`).
  - **Crucial Indexing**: Whenever you create or modify a schema, you MUST consult the `.agents/skills/index-database/SKILL.md` skill to analyze queries and explicitly add necessary single/composite indexes (especially for foreign keys).
  - Keep Buffalo controllers thin; push logic to Pop ORM models.
  - **Hooks**: Use the `BeforeValidate` hook in models to generate required fields programmatically before Pop runs validation checks. Do *not* use `BeforeCreate` for required fields.
- **Frontend (Svelte 5 / UI)**:
  - Edit files in `assets/js/` (Svelte components).
  - **Strictly Svelte 5 Runes**: Use `$state`, `$derived`, `$props`, `$effect`. No Options API. Maintain structure: Imports -> Runes -> Helpers -> Markup.
  - **API State**: Always define `let isLoading = $state(false);` and `let errorMsg = $state("");`. Handle loading and error states gracefully in the UI.
  - **Aesthetics**: Strictly use **DaisyUI** components for styling over custom CSS. Use Tailwind CSS responsive prefixes (`sm:`, `md:`, `lg:`).
  - **CSRF Tokens**: Include the CSRF token (`authenticity_token`) when making POST/PUT/DELETE requests against the Buffalo server (extract via `document.querySelector('meta[name="csrf-token"]')?.content`).
- **Task Tracking**: As you complete steps from `task.md`, check them off (e.g., change `[ ]` to `[x]`).

## After you edit:
1. Run `gofmt -w .` on modified Go code.
2. Run standard tests (`GO_ENV=test go test ./...`) to ensure you haven't broken the backend build.
3. Ensure the frontend builds cleanly (`npm run build`).
4. Return a short summary of files changed, commands run, and DaisyUI patterns reused.
