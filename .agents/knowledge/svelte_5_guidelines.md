# Svelte 5 & UI Guidelines

This project exclusively uses **Svelte 5 Runes** and **DaisyUI/Tailwind CSS v4**. AI agents must adhere to these standards.

## 1. Runes Syntax vs Options API
The Svelte 4 Options API is strictly forbidden.

**WRONG (Options API):**
```svelte
<script>
  export let title;
  let count = 0;
  $: doubled = count * 2;
</script>
```

**CORRECT (Svelte 5 Runes):**
```svelte
<script>
  let { title } = $props();
  let count = $state(0);
  let doubled = $derived(count * 2);
</script>
```

## 2. Forms & CSRF
If submitting a form directly via `fetch()` to the Buffalo backend, you must extract the CSRF token. Buffalo automatically injects this into the `<meta name="csrf-token">` tag on page load if using the standard `application.plush.html` layout.

```javascript
async function handleSubmit(event) {
  event.preventDefault();

  const csrfToken = document.querySelector('meta[name="csrf-token"]')?.content;

  const response = await fetch('/api/endpoint', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken // Or place it in the payload as authenticity_token
    },
    body: JSON.stringify({ data })
  });
}
```

## 3. Styling with DaisyUI
Always defer to DaisyUI classes before writing custom CSS or raw Tailwind utility classes.
- Use `.btn`, `.btn-primary` instead of `bg-blue-500 text-white rounded px-4 py-2`.
- Use `.card`, `.card-body` instead of building cards manually.
- Use `.form-control`, `.label`, `.input` for all forms.

If you must use custom colors, use the CSS variables provided by DaisyUI, e.g., `hsl(var(--p))` for primary.
