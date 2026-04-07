import { expect, test, type Locator, type Page } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

const loginUrl = /\/api\/auth\/login$/
const closeMonthUrl = /\/api\/expenses\/sheet\/close$/
const completeChecklistUrl = /\/api\/expenses\/checklists\/complete$/
const deleteEntryUrl = /\/api\/expenses\/entries\//
const dashboardUrl = /\/api\/expenses\/sheet\/current(\?|$)/

async function loginAndOpenExpenses(page: Page) {
  await page.goto('/login')
  await page.getByLabel('Email').fill(e2eEmail)
  await page.getByLabel('Password').fill(e2ePassword)

  await Promise.all([
    page.waitForResponse((response) => loginUrl.test(response.url()) && response.request().method() === 'POST'),
    page.getByRole('button', { name: 'Sign in' }).click(),
  ])

  await page.goto('/expenses')
  await page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET')
  await expect(page.getByText('Summary')).toBeVisible()
}

async function openObligationsPanel(page: Page) {
  const button = page.getByRole('button', { name: /Obligations/i })
  await expect(button).toBeVisible()
  const panel = page.locator('section').filter({ has: page.getByText('Obligations') }).first()
  if (await panel.locator('[data-list-item-id]').count()) {
    return
  }
  await button.click()
}

async function firstHistoryRow(page: Page): Promise<Locator> {
  const historyPanel = page.locator('section').filter({ has: page.getByText('History') }).first()
  return historyPanel.locator('[data-list-item-id]').first()
}

test.describe('Expenses page', () => {
  test('filters by category', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const historyPanel = page.locator('section').filter({ has: page.getByText('History') }).first()
    await expect(historyPanel.locator('[data-list-item-id]').first()).toBeVisible()

    await historyPanel.getByRole('button', { name: 'Fun' }).click()
    await expect(historyPanel.getByRole('button', { name: 'Fun' })).toHaveAttribute('aria-pressed', 'true')
    await expect(historyPanel.locator('text=Fund.')).toHaveCount(0)
    await expect(historyPanel.locator('text=Future')).toHaveCount(0)
  })

  test('completes an obligation', async ({ page }) => {
    await loginAndOpenExpenses(page)
    await openObligationsPanel(page)

    const obligationsPanel = page.locator('section').filter({ has: page.getByText('Obligations') }).first()
    if (await obligationsPanel.locator('[data-list-item-id]').count() === 0) {
      return
    }

    const firstItem = obligationsPanel.locator('[data-list-item-id]').first()
    await expect(firstItem).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => completeChecklistUrl.test(response.url()) && response.request().method() === 'POST'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      firstItem.getByRole('button', { name: 'Mark as complete' }).click(),
    ])
  })

  test('deletes a history row', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const row = await firstHistoryRow(page)
    const rowId = await row.getAttribute('data-list-item-id')
    expect(rowId).toBeTruthy()

    await Promise.all([
      page.waitForResponse((response) => deleteEntryUrl.test(response.url()) && response.request().method() === 'DELETE'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      row.getByRole('button', { name: /Delete / }).click(),
    ])

    await expect(page.locator(`[data-list-item-id="${rowId}"]`)).toHaveCount(0)
  })

  test('closes the current month', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const periodStat = page.getByTestId('expenses-period')
    const initialPeriod = await periodStat.textContent()

    await Promise.all([
      page.waitForResponse((response) => closeMonthUrl.test(response.url()) && response.request().method() === 'POST'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      page.getByRole('button', { name: 'Close month' }).click(),
    ])

    await expect(page.getByRole('button', { name: 'Close month' })).toBeVisible()
    if (initialPeriod) {
      await expect(periodStat).not.toHaveText(initialPeriod)
    }
  })
})
