import { defineConfig } from "@playwright/test";
export default defineConfig({
  testDir: "tests/browser",
  workers: 1,
  retries: 0,
  use: { baseURL: "http://127.0.0.1:3100", trace: "retain-on-failure" },
  webServer: {
    command: "go run ./cmd/server",
    url: "http://127.0.0.1:3100/health/ready",
    reuseExistingServer: false,
    env: {
      APP_ENV: "test",
      HTTP_ADDR: "127.0.0.1:3100",
      PUBLIC_URL: "http://127.0.0.1:3100",
    },
  },
});
