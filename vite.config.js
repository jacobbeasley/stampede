import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    tailwindcss(),
    svelte(),
  ],
  // root is the project root; Vite resolves all paths from here
  // Buffalo already serves public/ directly; don't let Vite copy it into the outDir
  publicDir: false,
  build: {
    outDir: 'public/assets',
    assetsDir: '',
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      input: {
        main: 'assets/js/main.js',
      },
    },
  },
  server: {
    port: 3001,
    cors: true,
    origin: 'http://localhost:3001',
  },
})
