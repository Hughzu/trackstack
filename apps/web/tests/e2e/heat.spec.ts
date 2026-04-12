import { expect, test, type Page } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

const loginUrl = /\/api\/auth\/login$/
const dashboardUrl = /\/api\/heat\/dashboard(\?|$)/
const createRefillUrl = /\/api\/heat\/refills$/
const deleteRefillUrl = /\/api\/heat\/refills\//

function uniqueSuffix() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function confirmSheet(page: Page) {
  return page.locator('.fixed.inset-0.z-50').last()
}

async function loginAndOpenHeat(page: Page) {
  await page.goto('/login')
  await page.getByLabel('Email').fill(e2eEmail)
  await page.getByLabel('Password').fill(e2ePassword)

  const [loginResponse] = await Promise.all([
    page.waitForResponse((response) => loginUrl.test(response.url()) && response.request().method() === 'POST'),
    page.getByRole('button', { name: 'Sign in' }).click(),
  ])

  expect(loginResponse.ok()).toBeTruthy()

  await page.goto('/heat')
  await page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET')
  await expect(page.getByRole('link', { name: 'Add refill' })).toBeVisible()
}

async function createRefill(page: Page, options?: { bags?: string, date?: string, temperature?: string }) {
  await page.getByRole('link', { name: 'Add refill' }).click()
  await expect(page.getByRole('button', { name: 'Log refill' })).toBeVisible()

  const suffix = uniqueSuffix()
  const temperature = options?.temperature ?? String((Number(suffix.slice(-2).replace(/\D/g, '').padEnd(2, '0').slice(0, 2)) % 20) + 5)
  const date = options?.date ?? '2025-03-10'

  await page.locator('#heat-bags').fill(options?.bags ?? '2')
  await page.locator('#heat-date').fill(date)
  await page.locator('#heat-temperature').fill(temperature)

  await Promise.all([
    page.waitForResponse((response) => createRefillUrl.test(response.url()) && response.request().method() === 'POST'),
    page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
    page.getByRole('button', { name: 'Log refill' }).click(),
  ])

  await expect(page).toHaveURL(/\/heat$/)
  return { temperature, date }
}

test.describe('Heat page', () => {
  test.describe.configure({ mode: 'serial' })

  test('creates a refill from the add page', async ({ page }) => {
    await loginAndOpenHeat(page)

    const { temperature, date } = await createRefill(page)
    const expectedDate = new Date(`${date}T00:00:00`).toLocaleDateString('en-IE', { month: 'short', day: 'numeric' })

    const row = page.locator('[data-list-item-id]').filter({ hasText: expectedDate }).filter({ hasText: `${temperature} C` }).first()
    await expect(row).toBeVisible()
  })

  test('deletes a refill from history', async ({ page }) => {
    await loginAndOpenHeat(page)

    const { temperature, date } = await createRefill(page, { bags: '1' })
    const expectedDate = new Date(`${date}T00:00:00`).toLocaleDateString('en-IE', { month: 'short', day: 'numeric' })
    const row = page.locator('[data-list-item-id]').filter({ hasText: expectedDate }).filter({ hasText: `${temperature} C` }).first()
    await expect(row).toBeVisible()

    const rowId = await row.getAttribute('data-list-item-id')
    expect(rowId).toBeTruthy()

    await row.getByRole('button', { name: `Delete refill ${expectedDate}` }).click()
    const sheet = confirmSheet(page)
    await expect(sheet.getByRole('heading', { name: 'Delete refill' })).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => deleteRefillUrl.test(response.url()) && response.request().method() === 'DELETE'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      sheet.getByRole('button', { name: 'Delete refill' }).click(),
    ])

    await expect(page.locator(`[data-list-item-id="${rowId}"]`)).toHaveCount(0)
  })
})
