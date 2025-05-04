import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());

  return {
    plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
    server: {
      host: "0.0.0.0", // Exposes the dev server to Docker
      port: 3000,
      strictPort: true,
      watch: {
        usePolling: true,
      },
      allowedHosts: [env.VITE_HOSTNAME ?? "localhost"],
    },
  };
});