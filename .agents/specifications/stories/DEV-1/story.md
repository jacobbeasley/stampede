# DEV-1: Scaffold Repository

## Story

> *As a developer, I want to scaffold the initial Buffalo repository with Tailwind CSS, DaisyUI, and Svelte integration so I have a solid, modern foundation to build the To-Do Manager application on.*

## Background

The `buffalo-app` workspace is currently empty (only an `AGENTS.md` and `.gitignore` exist). This story covers everything needed to stand up a fully working Buffalo application skeleton that is wired up to the frontend toolchain specified in the application spec.

## Acceptance Criteria

- [x] A working Buffalo application is initialized in the repo root.
- [x] A PostgreSQL `database.yml` is configured with sensible defaults for `development`, `test`, and `production` environments using environment variables.
- [x] The Buffalo app boots with `buffalo dev` without errors.
- [x] Tailwind CSS and DaisyUI are installed and configured.
- [x] Svelte is installed and a Vite-based build pipeline is integrated with Buffalo's asset serving.
- [x] A basic landing page (`/`) renders using a Plush template that loads the Tailwind/DaisyUI CSS.
- [x] A hello-world Svelte component (`HelloWorld.svelte`) is mounted on the landing page and renders correctly.
- [x] `npm run build` compiles Svelte + Tailwind assets to `public/assets/`.
- [x] `buffalo dev` + `npm run dev` serves the app on `localhost:3000` with the Vite dev server on `localhost:3001`.
- [x] The repo is committed in a clean, logical state.
