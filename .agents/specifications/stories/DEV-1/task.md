# DEV-1 Task Tracker

## Phase 1 — Buffalo Scaffolding
- [x] Run `buffalo new` with `--skip-webpack`, `--db-type postgres`, `--module buffalo-app`
- [x] Patch `database.yml` to use env-variable-based URLs

## Phase 2 — Frontend Toolchain
- [x] `npm init` and install Vite, Svelte, Tailwind, DaisyUI
- [x] Create `vite.config.js`
- [x] Create `assets/css/application.css`
- [x] Create `assets/js/main.js`
- [x] Create `assets/js/components/HelloWorld.svelte`

## Phase 3 — Buffalo Template Integration
- [x] Patch `templates/application.plush.html` with DaisyUI layout + Vite script tags
- [x] Patch `templates/home/index.plush.html` with modern landing page + `<div id="app">`

## Phase 4 — Database Setup
- [ ] `buffalo pop create -a`
- [ ] `buffalo pop migrate`

## Phase 5 — Validation
- [ ] `npm run build` passes
- [ ] `buffalo test` passes
- [ ] Browser smoke test — landing page + Svelte component renders

## Phase 6 — Git Commit
- [ ] `gofmt -w .`
- [ ] Commit all scaffold changes
