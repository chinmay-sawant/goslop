import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

const rootDir = import.meta.dirname

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(rootDir, './src'),
    },
  },
  // Relative base so the site works from GitHub Pages /docs or file:// previews.
  base: './',
  build: {
    // Emit static site into the repo-root docs/ folder (GitHub Pages).
    outDir: path.resolve(rootDir, '../docs'),
    emptyOutDir: true,
  },
})
