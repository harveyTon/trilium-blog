import vue from "@vitejs/plugin-vue";
import { defineConfig } from "vite";
import viteCompression from "vite-plugin-compression";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), viteCompression()],
  build: {
    cssCodeSplit: true,
    reportCompressedSize: false,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ["vue", "vue-router"],
        },
      },
    },
  },
});
