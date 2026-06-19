import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// В dev все запросы /api проксируются на Go-сервис (по умолчанию :8080),
// поэтому CORS на бэке в разработке не нужен.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: process.env.VITE_API_TARGET || "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
