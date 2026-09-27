import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { mkdirSync, writeFileSync, rmSync } from "node:fs";
import { fileURLToPath } from "node:url";

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    {
      name: "gonertia-hot-file",
      configureServer(server) {
        mkdirSync("tmp", { recursive: true });
        writeFileSync("tmp/vite.hot", "http://localhost:5173");
        server.httpServer?.once("close", () =>
          rmSync("tmp/vite.hot", { force: true }),
        );
      },
    },
  ],
  resolve: {
    alias: { "@": fileURLToPath(new URL("./web/src", import.meta.url)) },
  },
  publicDir: "web/public",
  server: {
    host: "127.0.0.1",
    port: 5173,
    strictPort: true,
    cors: {
      origin: new URL(process.env.PUBLIC_URL ?? "http://localhost:3000").origin,
    },
  },
  build: {
    outDir: "web/build",
    emptyOutDir: true,
    manifest: "manifest.json",
    rollupOptions: { input: "web/src/app.tsx" },
  },
});
