import { test, expect } from '@playwright/test';

const e2eEmail = process.env.E2E_TEST_EMAIL ?? '';
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? '';

async function login(page: import('@playwright/test').Page) {
    await page.goto('/login');
    const responsePromise = page.waitForResponse(response =>
        response.url().includes('/api/auth/login') && response.request().method() === 'POST'
    );
    await page.fill('input[name="email"]', e2eEmail);
    await page.fill('input[name="password"]', e2ePassword);
    await page.click('button[type="submit"]');
    const response = await responsePromise;
    expect(response.ok()).toBeTruthy();
    await page.waitForURL('/');
}

test.describe('Heat Refill Flow', () => {
    test('User can delete a refill from history', async ({ page }) => {
        await login(page);

        const deleteResponsePromise = page.waitForResponse(response =>
            response.url().includes('/api/heat/refills/') && response.request().method() === 'DELETE'
        );

        await page.goto('/heat/new');
        await page.fill('input[name="bags"]', '1');
        await page.fill('input[name="weightKg"]', '15');
        await page.fill('input[name="date"]', '2025-03-10');

        await page.fill('input[name="temperature"]', '12');

        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/heat');

        const initialDeleteCount = await page.locator('[data-refill-delete]').count();
        expect(initialDeleteCount).toBeGreaterThan(0);

        const deleteButtons = page.locator('[data-refill-delete]');
        await deleteButtons.first().click();
        const deleteModal = page.locator('#refill-delete-modal');
        await expect(deleteModal).toBeVisible();

        await deleteModal.locator('[data-confirm-modal]').click();
        const deleteResponse = await deleteResponsePromise;
        expect(deleteResponse.ok()).toBeTruthy();
        await expect(deleteModal).toBeHidden();

        const finalDeleteCount = await page.locator('[data-refill-delete]').count();
        expect(finalDeleteCount).toBe(initialDeleteCount - 1);
    });
});
