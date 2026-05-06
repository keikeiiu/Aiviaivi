import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  timeout: 30000,
  retries: 1,
  use: {
    baseURL: "http://localhost:19000",
    headless: true,
  },
  webServer: {
    command: "echo 'Using running server'",
    port: 19000,
    reuseExistingServer: true,
  },
});
