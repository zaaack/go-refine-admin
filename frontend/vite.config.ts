import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";

const srcDir = path.resolve(import.meta.dirname, "src");

// Vite 8 config. Dev server proxies /api and /login to the Go backend (port 8080).
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": srcDir,
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/login": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/uploads": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/swagger": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
