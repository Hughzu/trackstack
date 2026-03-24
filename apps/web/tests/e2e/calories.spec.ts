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

test.describe('Calories Logging Flow', () => {
    test('User can submit a new calorie entry', async ({ page }) => {
        await login(page);

        // 1. Visit the new calories page
        await page.goto('/calories/new');

        // 2. Fill out the form
        await page.fill('input[name="calories"]', '500');
        await page.fill('input[name="proteinGrams"]', '30');
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

    test('User can update calorie targets', async ({ page }) => {
        await login(page);

        await page.goto('/calories/settings');
        await page.fill('input[name="targetCalories"]', '2600');
        await page.fill('input[name="targetProteinGrams"]', '190');
        await page.fill('input[name="targetCarbGrams"]', '240');
        await page.fill('input[name="targetFatGrams"]', '80');

        const responsePromise = page.waitForResponse(response =>
            response.url().includes('/api/calories/target') && response.request().method() === 'POST'
        );

        await page.click('button[type="submit"]');

        const response = await responsePromise;
        expect(response.ok()).toBeTruthy();
        await expect(page).toHaveURL('/calories');
    });

    test('User can quick-add a recent meal from the dashboard', async ({ page }) => {
        await login(page);

        await page.goto('/calories/new');
        await page.fill('input[name="calories"]', '640');
        await page.fill('input[name="proteinGrams"]', '42');
        await page.fill('input[name="title"]', 'Quick Add Seed');
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/calories');

        const quickAddButton = page.getByRole('button', { name: 'Quick add Quick Add Seed' }).first();
        await expect(quickAddButton).toBeVisible();

        const responsePromise = page.waitForResponse(response =>
            response.url().includes('/api/calories/log') && response.request().method() === 'POST'
        );

        await quickAddButton.click();

        const response = await responsePromise;
        expect(response.ok()).toBeTruthy();
        await expect(page).toHaveURL('/calories');
    });

    test('User can delete a calorie log', async ({ page }) => {
        await login(page);

        const mealTitle = `Delete Me ${uniqueSuffix()}`;

        await page.goto('/calories/new');
        await page.fill('input[name="calories"]', '510');
        await page.fill('input[name="proteinGrams"]', '35');
        await page.fill('input[name="title"]', mealTitle);
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/calories');

        const createdRow = page.locator('[data-calorie-logs-list] > div').filter({ hasText: mealTitle }).first();
        await expect(createdRow).toBeVisible();

        await createdRow.locator('button[data-calorie-delete]').click();
        const deleteModal = page.locator('#calorie-delete-modal');
        await expect(deleteModal).toBeVisible();

        const [deleteResponse] = await Promise.all([
            page.waitForResponse(response =>
                response.url().includes('/api/calories/logs/') && response.request().method() === 'DELETE'
            ),
            deleteModal.locator('[data-confirm-modal]').click(),
        ]);

        expect(deleteResponse.ok()).toBeTruthy();
        await expect(deleteModal).toBeHidden();
        await expect(page.locator('[data-calorie-logs-list] > div').filter({ hasText: mealTitle })).toHaveCount(0);
    });
});
