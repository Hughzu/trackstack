import { expect, test, type Page } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

const loginUrl = /\/api\/auth\/login$/
const dashboardUrl = /\/api\/calories\/dashboard(\?|$)/
const createLogUrl = /\/api\/calories\/log$/
const updateTargetUrl = /\/api\/calories\/target$/
const deleteLogUrl = /\/api\/calories\/logs\//

function uniqueSuffix() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function confirmSheet(page: Page) {
  return page.locator('.fixed.inset-0.z-50').last()
}

async function loginAndOpenCalories(page: Page) {
  await page.goto('/login')
  await page.getByLabel('Email').fill(e2eEmail)
  await page.getByLabel('Password').fill(e2ePassword)

  const [loginResponse] = await Promise.all([
    page.waitForResponse((response) => loginUrl.test(response.url()) && response.request().method() === 'POST'),
    page.getByRole('button', { name: 'Sign in' }).click(),
  ])

  expect(loginResponse.ok()).toBeTruthy()

  await page.goto('/calories')
  await page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET')
  await expect(page.getByRole('heading', { name: 'Today' })).toBeVisible()
}

async function createMeal(page: Page, options?: { title?: string, calories?: string, protein?: string, carbs?: string, fat?: string }) {
  const title = options?.title ?? `E2E meal ${uniqueSuffix()}`

  await page.getByRole('link', { name: 'Add Meal' }).click()
  await expect(page.getByRole('button', { name: 'Save meal' })).toBeVisible()

  await page.locator('#calorie-amount').fill(options?.calories ?? '610')
  await page.getByLabel('Title').fill(title)
  await page.locator('#calorie-protein').fill(options?.protein ?? '42')
  await page.locator('#calorie-carbs').fill(options?.carbs ?? '58')
  await page.locator('#calorie-fat').fill(options?.fat ?? '14')

  await Promise.all([
    page.waitForResponse((response) => createLogUrl.test(response.url()) && response.request().method() === 'POST'),
    page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
    page.getByRole('button', { name: 'Save meal' }).click(),
  ])

  await expect(page).toHaveURL(/\/calories$/)
  return title
}

test.describe('Calories page', () => {
  test.describe.configure({ mode: 'serial' })

  test('creates a new meal from the add page', async ({ page }) => {
    await loginAndOpenCalories(page)

    const title = await createMeal(page)
    await expect(page.locator('[data-list-item-id]').filter({ hasText: title }).first()).toBeVisible()
  })

  test('keeps macro input focus while typing on the add meal page', async ({ page }) => {
    await loginAndOpenCalories(page)

    await page.getByRole('link', { name: 'Add Meal' }).click()
    await expect(page.getByRole('button', { name: 'Save meal' })).toBeVisible()

    const proteinInput = page.locator('#calorie-protein')
    await proteinInput.click()
    await expect(proteinInput).toBeFocused()

    await page.keyboard.press('4')
    await expect(proteinInput).toBeFocused()
    await expect(proteinInput).toHaveValue('4')

    await page.keyboard.press('2')
    await expect(proteinInput).toBeFocused()
    await expect(proteinInput).toHaveValue('42')
  })

  test('updates calorie targets from settings', async ({ page }) => {
    await loginAndOpenCalories(page)

    await page.getByRole('link', { name: 'Settings' }).click()
    await expect(page.getByRole('button', { name: 'Save settings' })).toBeVisible()

    await page.locator('#calories-target').fill('2750')
    await page.locator('#calories-target-protein').fill('205')
    await page.locator('#calories-target-carbs').fill('260')
    await page.locator('#calories-target-fat').fill('85')

    const [updateResponse] = await Promise.all([
      page.waitForResponse((response) => updateTargetUrl.test(response.url()) && response.request().method() === 'POST'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      page.getByRole('button', { name: 'Save settings' }).click(),
    ])

    expect(updateResponse.ok()).toBeTruthy()
    await expect(page).toHaveURL(/\/calories$/)
  })

  test('quick-adds a recent meal from the dashboard', async ({ page }) => {
    await loginAndOpenCalories(page)

    const title = await createMeal(page, { title: `Quick Add Seed ${uniqueSuffix()}` })
    const quickAddButton = page.getByRole('button', { name: `Quick add ${title}` }).first()
    await expect(quickAddButton).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => createLogUrl.test(response.url()) && response.request().method() === 'POST'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      quickAddButton.click(),
    ])

    await expect(page.locator('[data-list-item-id]').filter({ hasText: title })).toHaveCount(2)
  })

  test('deletes a calorie log from the dashboard', async ({ page }) => {
    await loginAndOpenCalories(page)

    const title = await createMeal(page, { title: `Delete Meal ${uniqueSuffix()}` })
    const row = page.locator('[data-list-item-id]').filter({ hasText: title }).first()
    await expect(row).toBeVisible()

    const rowId = await row.getAttribute('data-list-item-id')
    expect(rowId).toBeTruthy()

    await row.getByRole('button', { name: `Delete ${title}` }).click()

    const sheet = confirmSheet(page)
    await expect(sheet.getByRole('heading', { name: 'Delete meal' })).toBeVisible()
    await expect(sheet.getByText(title, { exact: false })).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => deleteLogUrl.test(response.url()) && response.request().method() === 'DELETE'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      sheet.getByRole('button', { name: 'Delete meal' }).click(),
    ])

    await expect(page.locator(`[data-list-item-id="${rowId}"]`)).toHaveCount(0)
  })
})
