import { expect, test } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

test.describe('Auth smoke flow', () => {
  test('redirects guests, logs in, and logs out', async ({ page }) => {
    await page.goto('/expenses')
    await page.waitForURL('/login')

    await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible()

    await page.fill('input[name="email"]', e2eEmail)
    await page.fill('input[name="password"]', e2ePassword)

    const [loginResponse] = await Promise.all([
      page.waitForResponse(
        (response) => response.url().includes('/api/auth/login') && response.request().method() === 'POST',
      ),
      page.click('button[type="submit"]'),
    ])

    expect(loginResponse.ok()).toBeTruthy()

    await page.waitForURL('/')
    await expect(page.getByRole('button', { name: 'Logout' })).toBeVisible()
    await expect(page.locator('nav').getByRole('link', { name: 'Expenses', exact: true })).toBeVisible()

    const [logoutResponse] = await Promise.all([
      page.waitForResponse(
        (response) => response.url().includes('/api/auth/logout') && response.request().method() === 'POST',
      ),
      page.getByRole('button', { name: 'Logout' }).click(),
    ])

    expect(logoutResponse.ok()).toBeTruthy()

    await page.waitForURL('/login')
    await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible()

    await page.goto('/heat')
    await page.waitForURL('/login')
  })
})
