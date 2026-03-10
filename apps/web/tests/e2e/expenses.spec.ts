import { test, expect } from '@playwright/test';

const e2eEmail = process.env.E2E_TEST_EMAIL ?? '';
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? '';

async function login(page: import('@playwright/test').Page) {
    await page.goto('/login');
    await page.fill('input[name="email"]', e2eEmail);
    await page.fill('input[name="password"]', e2ePassword);
    await page.click('button[type="submit"]');
    await page.waitForURL('/');
}

test.describe('Expenses Logging Flow', () => {
    test('User can submit a new expense entry', async ({ page }) => {
        await login(page);

        // 1. Visit the new expense page
        await page.goto('/expenses/new');

        // 2. Fill out the form
        await page.fill('input[name="amount"]', '42.50');
        // Click visible category tile to select the underlying radio input.
        await page.getByText('Fundamentals').click();
        await page.fill('input[name="title"]', 'Groceries');

        // 3. Submit the form
        const responsePromise = page.waitForResponse(response =>
            response.url().includes('/api/expenses/expense') && response.request().method() === 'POST'
        );

        await page.click('button[type="submit"]');

        // 4. Verify the response is successful
        const response = await responsePromise;
        expect(response.ok()).toBeTruthy();

        // 5. Verify the user is redirected to the expenses dashboard
        await expect(page).toHaveURL('/expenses');
    });

    test('User can save expense settings', async ({ page }) => {
        await login(page);

        await page.goto('/expenses/settings');
        await page.fill('input[name="income"]', '2800');
        await page.fill('input[name="ratioFund"]', '55');
        await page.fill('input[name="ratioFun"]', '20');
        await page.fill('input[name="ratioFuture"]', '25');

        const responsePromise = page.waitForResponse(response =>
            response.url().includes('/api/expenses/settings') && response.request().method() === 'POST'
        );

        await page.click('#expenses-settings-form button[type="submit"]');

        const response = await responsePromise;
        expect(response.ok()).toBeTruthy();
        await expect(page).toHaveURL('/expenses');
    });

    test('User can add monthly checklist and recurring templates', async ({ page }) => {
        await login(page);

        await page.goto('/expenses/settings');

        const checklistCount = await page.locator('[data-checklist-delete]').count();
        const recurringCount = await page.locator('[data-recurring-delete]').count();

        const checklistPostPromise = page.waitForResponse(response =>
            response.url().includes('/api/expenses/checklist') && response.request().method() === 'POST'
        );

        await page.locator('input[form="expense-checklist-form"][name="title"]').fill('Checklist Test');
        await page.locator('input[form="expense-checklist-form"][name="amount"]').fill('18.75');
        await page.locator('select[form="expense-checklist-form"][name="category"]').selectOption('fun');
        await page.locator('button[form="expense-checklist-form"][type="submit"]').click();

        const checklistResponse = await checklistPostPromise;
        expect(checklistResponse.ok()).toBeTruthy();
        await expect(page).toHaveURL('/expenses/settings');
        await expect(page.locator('[data-checklist-delete]')).toHaveCount(checklistCount + 1);

        const recurringPostPromise = page.waitForResponse(response =>
            response.url().includes('/api/expenses/recurring') && response.request().method() === 'POST'
        );

        await page.locator('input[form="expense-recurring-form"][name="title"]').fill('Recurring Test');
        await page.locator('input[form="expense-recurring-form"][name="amount"]').fill('44.10');
        await page.locator('select[form="expense-recurring-form"][name="category"]').selectOption('future');
        await page.locator('button[form="expense-recurring-form"][type="submit"]').click();

        const recurringResponse = await recurringPostPromise;
        expect(recurringResponse.ok()).toBeTruthy();
        await expect(page).toHaveURL('/expenses/settings');
        await expect(page.locator('[data-recurring-delete]')).toHaveCount(recurringCount + 1);
    });
});
