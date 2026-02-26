import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/photos': 'http://localhost:8080',
      '/sessions': 'http://localhost:8080',
      '/devices': 'http://localhost:8080',
      '/print-sizes': 'http://localhost:8080',
      '/qr': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
      '/admin': 'http://localhost:8080',
      '/metrics': 'http://localhost:8080',
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
})
