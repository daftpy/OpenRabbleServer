import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
  server: {
    host: "0.0.0.0", // Exposes the dev server to Docker
    port: 3000, // Forces Vite to use port 3000
    strictPort: true,
    watch: {
      usePolling: true, // Required for Docker to detect file changes
    },
    allowedHosts: ["example.com"],
  },
});
