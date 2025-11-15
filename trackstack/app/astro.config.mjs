// @ts-check
import { defineConfig } from "astro/config";
import tailwindcss from "@tailwindcss/vite";
import AstroPWA from "@vite-pwa/astro";

// https://astro.build/config
export default defineConfig({
  integrations: [
    AstroPWA({
      registerType: "autoUpdate",
      manifest: {
        name: "TrackStack",
        short_name: "TrackStack",
        description: "Track your habits and goals",
        theme_color: "#ffffff",
        background_color: "#ffffff",
        display: "standalone",
        icons: [
          // Android icons
          {
            src: "/android/android-launchericon-512-512.png",
            sizes: "512x512",
            type: "image/png",
            purpose: "any maskable",
          },
          {
            src: "/android/android-launchericon-192-192.png",
            sizes: "192x192",
            type: "image/png",
            purpose: "any maskable",
          },
          {
            src: "/android/android-launchericon-144-144.png",
            sizes: "144x144",
            type: "image/png",
          },
          {
            src: "/android/android-launchericon-96-96.png",
            sizes: "96x96",
            type: "image/png",
          },
          {
            src: "/android/android-launchericon-72-72.png",
            sizes: "72x72",
            type: "image/png",
          },
          {
            src: "/android/android-launchericon-48-48.png",
            sizes: "48x48",
            type: "image/png",
          },
          // iOS icons
          {
            src: "/ios/180.png",
            sizes: "180x180",
            type: "image/png",
          },
          {
            src: "/ios/152.png",
            sizes: "152x152",
            type: "image/png",
          },
          {
            src: "/ios/120.png",
            sizes: "120x120",
            type: "image/png",
          },
          // Favicons
          {
            src: "/ios/32.png",
            sizes: "32x32",
            type: "image/png",
          },
          {
            src: "/ios/16.png",
            sizes: "16x16",
            type: "image/png",
          },
        ],
      },
      workbox: {
        navigateFallback: "/",
        globPatterns: ["**/*.{css,js,html,svg,png,ico,txt}"],
      },
      devOptions: {
        enabled: true, // Enable PWA in dev mode for testing
      },
    }),
  ],
  vite: {
    plugins: [tailwindcss()],
  },
});