import type { BudgetBreakdownItem } from '../../components/ui/BudgetBreakdown'
import type { FilterOption } from '../../components/ui/FilterGroup'
import type { ExpensesDashboard } from '../../core/api/types'

type ExpenseCategoryId = 'fund' | 'fun' | 'future'

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

type BudgetBreakdownOptions = {
  compactFundLabel?: boolean
  includePercent?: boolean
  highlightOverBudget?: boolean
  budgetSource?: 'dashboard' | 'ratio'
}

export function getExpenseCategoryMeta(category: string) {
  return expenseCategoryMeta[normalizeExpenseCategory(category)]
}

export function createExpenseHistoryFilterOptions(): FilterOption[] {
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
): BudgetBreakdownItem[] {
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
      label: options.compactFundLabel && category === 'fund' ? meta.compactLabel : meta.label,
      value: formatEuro(spentByCategory[category]),
      subtext: options.includePercent ? `/ ${formatEuro(budget)} (${percent}%)` : `/ ${formatEuro(budget)}`,
      percent,
      color: options.highlightOverBudget && isOver ? 'danger' : 'main',
      barColor: meta.barColor,
    }
  })
}

function normalizeExpenseCategory(category: string): ExpenseCategoryId {
  if (category === 'fun' || category === 'future') return category
  return 'fund'
}
