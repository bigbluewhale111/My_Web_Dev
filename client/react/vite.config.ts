import { defineConfig } from 'vite'
import { resolve } from 'path'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname,'index.html'),
        login: resolve(__dirname,'login.html'),
        redirect: resolve(__dirname, 'redirect.html')
      }
    }
  }
})
