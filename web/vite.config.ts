import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import legacy from "@vitejs/plugin-legacy";

import { resolve } from "path";

const __dirname = resolve();

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [
          [
            "babel-plugin-direct-import",
            {
              modules: [
                "@mui/material",
                "@mui/lab",
                "@mui/icons-material",
                "react-icons",
              ],
            },
          ],
        ],
      },
    }),
    legacy(),
  ],
  resolve: {
    alias: {
      "@": resolve(__dirname, "src"),
      "@components": resolve(__dirname, "src/components"),
      "@hooks": resolve(__dirname, "src/hooks"),
      "@api": resolve(__dirname, "src/network/api"),
      "@util": resolve(__dirname, "src/util"),
      "@store": resolve(__dirname, "src/store"),
    },
  },
});
