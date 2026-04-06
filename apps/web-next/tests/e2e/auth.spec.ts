import { expect, test, type Page } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

const refreshCookieName = 'trackstack_refresh'

const refreshUrl = /\/api\/auth\/refresh$/
const loginUrl = /\/api\/auth\/login$/
const logoutUrl = /\/api\/auth\/logout$/

async function expectDashboardVisible(page: Page) {
  await expect(page.getByRole('button', { name: 'Logout' })).toBeVisible()
  await expect(page.locator('nav').getByRole('link', { name: 'Expenses', exact: true })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Expenses' }).first()).toBeVisible()
  await expect(page.getByRole('link', { name: 'Daily Intake' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Heating' })).toBeVisible()
}

test.describe('Auth flow', () => {
  test('logs in, refreshes the session, reaches dashboard, and logs out cleanly', async ({ page, context }) => {
    await page.goto('/')
    await page.waitForURL('/login')

    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible()

    await page.getByLabel('Email').fill(e2eEmail)
    await page.getByLabel('Password').fill(e2ePassword)

    const [loginResponse] = await Promise.all([
      page.waitForResponse((response) => loginUrl.test(response.url()) && response.request().method() === 'POST'),
      page.getByRole('button', { name: 'Sign in' }).click(),
    ])

    expect(loginResponse.ok()).toBeTruthy()

    await page.waitForURL('/')
    await expectDashboardVisible(page)

    const tokenAfterLogin = await page.evaluate(() => window.sessionStorage.getItem('trackstack_token'))
    expect(tokenAfterLogin).toBeTruthy()

    const cookiesAfterLogin = await context.cookies()
    expect(cookiesAfterLogin.some((cookie) => cookie.name === refreshCookieName && cookie.value.length > 0)).toBeTruthy()

    await page.evaluate(() => {
      window.sessionStorage.removeItem('trackstack_token')
    })

    const [refreshResponse] = await Promise.all([
      page.waitForResponse((response) => refreshUrl.test(response.url()) && response.request().method() === 'POST'),
      page.reload(),
    ])

    expect(refreshResponse.ok()).toBeTruthy()

    await page.waitForURL('/')
    await expectDashboardVisible(page)

    const tokenAfterRefresh = await page.evaluate(() => window.sessionStorage.getItem('trackstack_token'))
    expect(tokenAfterRefresh).toBeTruthy()
    expect(tokenAfterRefresh).not.toEqual(tokenAfterLogin)

    const [logoutResponse] = await Promise.all([
      page.waitForResponse((response) => logoutUrl.test(response.url()) && response.request().method() === 'POST'),
      page.getByRole('button', { name: 'Logout' }).click(),
    ])

    expect(logoutResponse.ok()).toBeTruthy()

    await page.waitForURL('/login')
    await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible()

    const tokenAfterLogout = await page.evaluate(() => window.sessionStorage.getItem('trackstack_token'))
    const logoutMarker = await page.evaluate(() => window.sessionStorage.getItem('trackstack_logged_out'))
    expect(tokenAfterLogout).toBeNull()
    expect(logoutMarker).toBe('1')

    const cookiesAfterLogout = await context.cookies()
    expect(cookiesAfterLogout.some((cookie) => cookie.name === refreshCookieName && cookie.value.length > 0)).toBeFalsy()

    await page.goto('/')
    await page.waitForURL('/login')
  })
})
