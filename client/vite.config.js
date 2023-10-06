import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:5500",
        secure: false,
      },
    },
  },
  base:
	process.env.NODE_ENV === 'production'
      ? '/home/' // prod
	  : '/', // dev
  plugins: [react()],
})
