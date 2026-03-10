import { test, expect } from '@playwright/test';

const e2eEmail = process.env.E2E_TEST_EMAIL ?? '';
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? '';

test.describe('Calories Logging Flow', () => {
    test('User can submit a new calorie entry', async ({ page }) => {
        // 0. Log in as the seeded user
        await page.goto('/login');
        await page.fill('input[name="email"]', e2eEmail);
        await page.fill('input[name="password"]', e2ePassword);
        await page.click('button[type="submit"]');
        await page.waitForURL('/');

        // 1. Visit the new calories page
        await page.goto('/calories/new');

        // 2. Fill out the form
        await page.fill('input[name="calories"]', '500');
        await page.fill('input[name="protein"]', '30');
        await page.fill('input[name="title"]', 'Test Meal');

        // 3. Submit the form
        const responsePromise = page.waitForResponse(response =>
            response.url().includes('/api/calories/log') && response.request().method() === 'POST'
        );

        await page.click('button[type="submit"]');

        // 4. Verify the response is successful
        const response = await responsePromise;
        expect(response.ok()).toBeTruthy();

        // 5. Verify the user is redirected to the calories dashboard
        await expect(page).toHaveURL('/calories');
    });
});
