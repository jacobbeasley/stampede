---
name: create-svelte-component
description: Use this skill when generating or modifying Svelte UI components in the frontend asset pipeline. This enforces modern Svelte 5 Runes syntax, explicit API state handling, and Tailwind/DaisyUI design constraints.
---

# Create Svelte Component

This skill outlines the strict requirements for building frontend components in this repository.

## 1. Syntax Restrictions (Svelte 5 Runes)
This project uses **Svelte 5** exclusively. The legacy Options API is forbidden.
- **Props**: Use `$props()` instead of `export let`.
- **State**: Use `$state()` for local reactive variables.
- **Derived State**: Use `$derived()` to compute values based on `$state`.
- **Side Effects**: Use `$effect()` for DOM side effects.
- **Functions**: Define standard functions inside the `<script>` block for event handlers. Do not write inline complex logic in the template.

## 2. API Communication
Components that fetch or send data to the backend must handle all standard interaction states:
- **Loading State**: Always define a `let isLoading = $state(false);` and use DaisyUI's `.loading` spinner to provide user feedback during the request.
- **Error State**: Always define a `let errorMsg = $state("");` and display errors using DaisyUI's `.alert .alert-error` component.
- **CSRF Tokens**: If sending a `POST`, `PUT`, or `DELETE` request against the Buffalo backend, you must extract the CSRF token from the DOM and include it in the request payload or headers:
  ```javascript
  const token = document.querySelector('meta[name="csrf-token"]')?.content;
  ```

## 3. UI/UX and Aesthetics
You must rely on standard design systems to keep the UI consistent.
- **DaisyUI First**: Use built-in components (e.g., `.btn`, `.card`, `.form-control`, `.input`, `.alert`). Do not rebuild standard components from scratch.
- **Semantic Colors**: Use semantic utility classes (e.g., `text-primary`, `bg-base-100`, `text-error`). Do not use raw Tailwind colors like `text-blue-500` or hex codes.
- **Responsiveness**: Use Tailwind's responsive prefixes (`sm:`, `md:`, `lg:`). Ensure multi-column layouts degrade nicely into single columns on mobile devices.
- **Minimal Custom CSS**: Do not use inline `style="..."` attributes. Only write custom CSS in the `<style>` block if a Tailwind or DaisyUI class genuinely cannot achieve the goal.

## 4. Component Structure
- Each `.svelte` file exports exactly one component.
- Extract large, complex blocks of markup into smaller sub-components stored in `assets/js/components/`.
- Maintain logical grouping of Svelte code: Imports first, Runes/State second, Helper functions third, Template markup last.
