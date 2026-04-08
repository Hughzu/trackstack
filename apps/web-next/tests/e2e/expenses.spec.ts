import { expect, test, type Locator, type Page } from '@playwright/test'

const e2eEmail = process.env.E2E_TEST_EMAIL ?? ''
const e2ePassword = process.env.E2E_TEST_PASSWORD ?? ''

const loginUrl = /\/api\/auth\/login$/
const closeMonthUrl = /\/api\/expenses\/sheet\/close$/
const completeChecklistUrl = /\/api\/expenses\/checklists\/complete$/
const createEntryUrl = /\/api\/expenses\/entries$/
const deleteEntryUrl = /\/api\/expenses\/entries\//
const dashboardUrl = /\/api\/expenses\/sheet\/current(\?|$)/

function confirmSheet(page: Page) {
  return page.locator('.fixed.inset-0.z-50').last()
}

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

async function createExpenseFromAddPage(page: Page, options?: { title?: string, amount?: string, category?: 'Fund.' | 'Fun' | 'Future' }) {
  const title = options?.title ?? `E2E expense ${Date.now()}`
  const amount = options?.amount ?? '27.40'
  const category = options?.category ?? 'Future'

  await page.getByRole('link', { name: 'Add Expense' }).click()
  await expect(page.getByRole('button', { name: 'Save Expense' })).toBeVisible()

  await page.locator('#expense-amount').fill(amount)
  await page.locator('[data-testid="expense-category-choice"]').getByRole('button', { name: new RegExp(`^${category}`) }).click()
  await page.getByLabel('Title').fill(title)

  await Promise.all([
    page.waitForResponse((response) => createEntryUrl.test(response.url()) && response.request().method() === 'POST'),
    page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
    page.getByRole('button', { name: 'Save Expense' }).click(),
  ])

  await expect(page).toHaveURL(/\/expenses$/)
  return title
}

test.describe('Expenses page', () => {
  test.describe.configure({ mode: 'serial' })

  test('filters by category', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const historyPanel = page.locator('section').filter({ has: page.getByText('History') }).first()
    await expect(historyPanel.locator('[data-list-item-id]').first()).toBeVisible()

    await historyPanel.getByRole('button', { name: 'Fun', exact: true }).click()
    await expect(historyPanel.getByRole('button', { name: 'Fun', exact: true })).toHaveAttribute('aria-pressed', 'true')
    await expect(historyPanel.locator('[data-list-item-id]').filter({ hasText: 'Fund.' })).toHaveCount(0)
    await expect(historyPanel.locator('[data-list-item-id]').filter({ hasText: 'Future' })).toHaveCount(0)
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

    const title = await createExpenseFromAddPage(page, { category: 'Future' })

    const historyPanel = page.locator('section').filter({ has: page.getByText('History') }).first()
    const row = historyPanel.locator('[data-list-item-id]').filter({ hasText: title }).first()
    await expect(row).toBeVisible()

    const rowId = await row.getAttribute('data-list-item-id')
    expect(rowId).toBeTruthy()

    await row.getByRole('button', { name: `Delete ${title}` }).click()

    const sheet = confirmSheet(page)
    await expect(sheet.getByRole('heading', { name: 'Delete expense' })).toBeVisible()
    await expect(sheet.getByText(title, { exact: false })).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => deleteEntryUrl.test(response.url()) && response.request().method() === 'DELETE'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      sheet.getByRole('button', { name: 'Delete expense' }).click(),
    ])

    await expect(page.locator(`[data-list-item-id="${rowId}"]`)).toHaveCount(0)
  })

  test('closes the current month', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const periodStat = page.getByTestId('expenses-period')
    const initialPeriod = await periodStat.textContent()

    await page.getByRole('button', { name: 'Close month' }).click()

    const sheet = confirmSheet(page)
    await expect(sheet.getByRole('heading', { name: 'Close month' })).toBeVisible()
    await expect(sheet.getByText('Period rollover')).toBeVisible()

    await Promise.all([
      page.waitForResponse((response) => closeMonthUrl.test(response.url()) && response.request().method() === 'POST'),
      page.waitForResponse((response) => dashboardUrl.test(response.url()) && response.request().method() === 'GET'),
      sheet.getByRole('button', { name: 'Yes, close month' }).click(),
    ])

    await expect(page.getByRole('button', { name: 'Close month' })).toBeVisible()
    if (initialPeriod) {
      await expect(periodStat).not.toHaveText(initialPeriod)
    }
  })

  test('creates a new expense from the add page', async ({ page }) => {
    await loginAndOpenExpenses(page)

    const title = await createExpenseFromAddPage(page, { category: 'Future' })
    await expect(page.locator('[data-list-item-id]').filter({ hasText: title }).first()).toBeVisible()
  })
})
