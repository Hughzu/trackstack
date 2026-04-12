import type { ExpensesDashboard } from '../../core/api/types'

export type ExpenseCategoryId = 'fund' | 'fun' | 'future'

type ExpenseCategoryMeta = {
  id: ExpenseCategoryId
  label: string
  compactLabel: string
  tone: 'danger' | 'warning' | 'success'
  barColor: 'danger' | 'warning' | 'success'
}

const expenseCategoryMeta: Record<ExpenseCategoryId, ExpenseCategoryMeta> = {
  fund: { id: 'fund', label: 'Fund', compactLabel: 'Fund.', tone: 'danger', barColor: 'danger' },
  fun: { id: 'fun', label: 'Fun', compactLabel: 'Fun', tone: 'warning', barColor: 'warning' },
  future: { id: 'future', label: 'Future', compactLabel: 'Future', tone: 'success', barColor: 'success' },
}

const expenseCategoryOrder = ['fund', 'fun', 'future'] as const satisfies readonly ExpenseCategoryId[]

export type ExpenseHistoryFilterOption = {
  value: string
  label: string
  tone?: 'danger' | 'warning' | 'success'
}

export type ExpenseBudgetPreviewItem = {
  id: ExpenseCategoryId
  label: string
  compactLabel: string
  description: string
  value: string
  subtext: string
  percent: number
  tone: 'danger' | 'warning' | 'success'
  color: 'main' | 'danger'
  barColor: 'danger' | 'warning' | 'success'
}

export type ExpenseCategoryChoice = {
  value: ExpenseCategoryId
  title: string
  description: string
  tone: 'danger' | 'warning' | 'success'
}

export type ExpenseRatioSummary = {
  fund: number
  fun: number
  future: number
  total: number
  totalTone: 'success' | 'danger' | 'warning'
  totalLabel: string
}

type BudgetBreakdownOptions = {
  compactFundLabel?: boolean
  includePercent?: boolean
  highlightOverBudget?: boolean
  budgetSource?: 'dashboard' | 'ratio'
}

export function getExpenseCategoryMeta(category: string) {
  return expenseCategoryMeta[normalizeExpenseCategory(category)]
}

export function createExpenseHistoryFilterOptions(): ExpenseHistoryFilterOption[] {
  return [
    { value: 'all', label: 'All' },
    ...expenseCategoryOrder.map((category) => {
      const meta = getExpenseCategoryMeta(category)
      return {
        value: meta.id,
        label: meta.compactLabel,
        tone: meta.tone,
      }
    }),
  ]
}

export function createExpenseBudgetBreakdownItems(
  dashboard: ExpensesDashboard,
  formatEuro: (value: number) => string,
  options: BudgetBreakdownOptions = {},
): ExpenseBudgetPreviewItem[] {
  const ratioByCategory = new Map(dashboard.ratios.map((ratio) => [normalizeExpenseCategory(ratio.categoryId), ratio]))
  const spentByCategory = {
    fund: dashboard.spent.fund,
    fun: dashboard.spent.fun,
    future: dashboard.spent.future,
  }
  const budgetByCategory = {
    fund: dashboard.budget.fund,
    fun: dashboard.budget.fun,
    future: dashboard.budget.future,
  }

  return expenseCategoryOrder.map((category) => {
    const meta = getExpenseCategoryMeta(category)
    const ratio = ratioByCategory.get(category)
    const budget = options.budgetSource === 'ratio' ? ratio?.budget || 0 : budgetByCategory[category]
    const percent = ratio?.percent || 0
    const isOver = Boolean(ratio?.over)

    return {
      id: category,
      label: options.compactFundLabel && category === 'fund' ? meta.compactLabel : meta.label,
      compactLabel: meta.compactLabel,
      description: meta.label,
      value: formatEuro(spentByCategory[category]),
      subtext: options.includePercent ? `/ ${formatEuro(budget)} (${percent}%)` : `/ ${formatEuro(budget)}`,
      percent,
      color: options.highlightOverBudget && isOver ? 'danger' : 'main',
      tone: meta.tone,
      barColor: meta.barColor,
    }
  })
}

export function createExpenseRatioSummary(values: { fund: number, fun: number, future: number }): ExpenseRatioSummary {
  const total = values.fund + values.fun + values.future

  return {
    fund: values.fund,
    fun: values.fun,
    future: values.future,
    total,
    totalTone: total === 100 ? 'success' : total > 100 ? 'danger' : 'warning',
    totalLabel: total === 100 ? 'Allocation locked' : total > 100 ? 'Over allocated' : 'Unassigned budget',
  }
}

export function createExpenseSettingsBudgetPreview(
  income: number,
  ratioSummary: ExpenseRatioSummary,
  formatEuro: (value: number) => string,
): ExpenseBudgetPreviewItem[] {
  return expenseCategoryOrder.map((category) => {
    const meta = getExpenseCategoryMeta(category)
    const percent = category === 'fund' ? ratioSummary.fund : category === 'fun' ? ratioSummary.fun : ratioSummary.future
    const amount = (income * percent) / 100

    return {
      id: category,
      label: meta.label,
      compactLabel: meta.compactLabel,
      description: category === 'fund'
        ? 'Bills, rent, insurance, the annoying grown-up stuff.'
        : category === 'fun'
          ? 'Leisure, meals out, random little dopamine hits.'
          : 'Savings, investing, buffers, not screwing over future-you.',
      value: formatEuro(amount),
      subtext: `${formatRatio(percent)} of income`,
      percent,
      tone: meta.tone,
      color: 'main',
      barColor: meta.barColor,
    }
  })
}

export function createExpenseCategoryChoices(): ExpenseCategoryChoice[] {
  return expenseCategoryOrder.map((category) => {
    const meta = getExpenseCategoryMeta(category)

    return {
      value: meta.id,
      title: meta.label === 'Fund' ? 'Fundamentals' : meta.label,
      description: meta.compactLabel,
      tone: meta.tone,
    }
  })
}

export function formatRatio(value: number) {
  return `${Math.round(value)}%`
}

export function toNumberOrZero(value: string) {
  const parsed = Number.parseFloat(value)
  return Number.isFinite(parsed) ? parsed : 0
}

export function clampRatio(value: string) {
  const parsed = Number.parseFloat(value)
  if (!Number.isFinite(parsed)) return '0'
  if (parsed < 0) return '0'
  if (parsed > 100) return '100'
  return String(parsed)
}

function normalizeExpenseCategory(category: string): ExpenseCategoryId {
  if (category === 'fun' || category === 'future') return category
  return 'fund'
}
