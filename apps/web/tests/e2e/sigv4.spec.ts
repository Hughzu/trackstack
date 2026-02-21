import { test, expect } from '@playwright/test';

// Read credentials from environment variables
const TEST_USER = {
  email: process.env.TEST_USER_EMAIL || '',
  password: process.env.TEST_USER_PASSWORD || ''
};

test.describe('SigV4 Form Submission', () => {
  test.beforeEach(async ({ page }) => {
    if (!TEST_USER.email || !TEST_USER.password) {
      test.skip(true, 'TEST_USER_EMAIL and TEST_USER_PASSWORD must be set in .env');
    }
    
    // Login flow - creates session cookie
    await page.goto('/login');
    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);
    
    // Wait for navigation after login
    await Promise.all([
      page.waitForURL('**/'),
      page.click('button[type="submit"]')
    ]);
  });

  test('calories form submits with SigV4 headers', async ({ page }) => {
    await page.goto('/calories/new');
    
    // Fill form
    await page.fill('input[name="calories"]', '500');
    await page.fill('input[name="protein"]', '30');
    await page.fill('input[name="title"]', 'Test Meal');
    
    // Capture the request
    const requestPromise = page.waitForRequest('**/api/calories/log');
    await page.click('button[type="submit"]');
    
    const request = await requestPromise;
    
    // Verify SigV4 requirements
    expect(request.headers()['content-type']).toBe('application/json');
    expect(request.headers()['x-amz-content-sha256']).toBeDefined();
    expect(request.headers()['x-amz-content-sha256']).toMatch(/^[a-f0-9]{64}$/);
    
    // Verify response handling
    await expect(page).toHaveURL('/calories');
  });

  test('form shows error on API failure', async ({ page }) => {
    await page.goto('/calories/new');
    
    // Submit without required fields
    await page.click('button[type="submit"]');
    
    // Should stay on page and show error
    await expect(page).toHaveURL('/calories/new');
    
    // Error element should be visible
    const errorVisible = await page.locator('[data-error-target]').isVisible().catch(() => false);
    // Or check for error message in DOM
    const hasError = await page.locator('text=/error|failed|missing/i').count() > 0;
    expect(hasError || errorVisible).toBeTruthy();
  });

  test('DELETE request includes SigV4 headers', async ({ page }) => {
    // First create a log
    await page.goto('/calories/new');
    await page.fill('input[name="calories"]', '100');
    await page.fill('input[name="protein"]', '10');
    await page.click('button[type="submit"]');
    await page.waitForURL('/calories');
    
    // Now try to delete it (using the modal)
    await page.click('[data-delete-trigger]');
    
    const requestPromise = page.waitForRequest('**/api/calories/log**');
    await page.click('[data-confirm-delete]');
    
    const request = await requestPromise;
    expect(request.method()).toBe('DELETE');
    expect(request.headers()['x-amz-content-sha256']).toBeDefined();
  });
});

test.describe('Session Handling', () => {
  test('unauthenticated user redirected to login', async ({ page }) => {
    await page.goto('/calories');
    await expect(page).toHaveURL('/login');
  });

  test('API returns 401 without session', async ({ page }) => {
    const response = await page.request.post('/api/calories/log', {
      data: { calories: 500, protein: 30 }
    });
    expect(response.status()).toBe(401);
  });
});
