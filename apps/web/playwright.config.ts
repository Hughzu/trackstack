import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
    testDir: './tests/e2e',
    timeout: 10 * 1000,
    expect: { timeout: 5000 },
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 2 : 0,
    workers: process.env.CI ? 1 : undefined,
    reporter: process.env.CI ? [['html', { open: 'never' }]] : 'list',
    use: {
        baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://127.0.0.1:4321',
        trace: 'on-first-retry',
    },

    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],
});
