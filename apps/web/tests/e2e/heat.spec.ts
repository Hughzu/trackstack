import { test, expect } from '@playwright/test';

const e2eEmail = process.env.E2E_TEST_EMAIL ?? '';
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? '';

function uniqueSuffix() {
    return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

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

        const uniqueTemperature = 50 + Number(uniqueSuffix().slice(-2).replace(/\D/g, '').padEnd(2, '0').slice(0, 2));

        await page.goto('/heat/new');
        await page.fill('input[name="bags"]', '1');
        await page.fill('input[name="weightKg"]', '15');
        await page.fill('input[name="date"]', '2025-03-10');

        await page.fill('input[name="temperature"]', String(uniqueTemperature));

        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/heat');

        const createdRow = page.locator('[data-refill-history] > div').filter({ hasText: `Avg Temp: ${uniqueTemperature}degC` }).first();
        await expect(createdRow).toBeVisible();

        await createdRow.locator('[data-refill-delete]').click();
        const deleteModal = page.locator('#refill-delete-modal');
        await expect(deleteModal).toBeVisible();

        const [deleteResponse] = await Promise.all([
            page.waitForResponse(response =>
                response.url().includes('/api/heat/refills/') && response.request().method() === 'DELETE'
            ),
            deleteModal.locator('[data-confirm-modal]').click(),
        ]);

        expect(deleteResponse.ok()).toBeTruthy();
        await expect(deleteModal).toBeHidden();
        await expect(page.locator('[data-refill-history] > div').filter({ hasText: `Avg Temp: ${uniqueTemperature}degC` })).toHaveCount(0);
    });
});
