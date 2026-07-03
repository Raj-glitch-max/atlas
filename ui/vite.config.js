import { defineConfig } from "vite";

export default defineConfig({
  base: "./",
  server: { host: "127.0.0.1", port: 5173, open: false },
  build: {
    target: "es2020",
    // keep the initial bundle lean; the heavy 3D (three) is split out and
    // lazy-loaded, so it never blocks first paint / the reduced-motion path.
    rollupOptions: {
      input: {
        main: "index.html",
        console: "console.html",
      },
      output: {
        manualChunks: {
          three: ["three"],
          gsap: ["gsap"],
        },
      },
    },
  },
});
