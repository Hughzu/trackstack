import { BudgetBreakdown, type BudgetBreakdownItem } from '../../../components/ui/BudgetBreakdown'
import { Panel } from '../../../components/ui/Panel'
import type { ExpensesDashboard } from '../../../core/api/types'
import { formatEuro } from '../../../core/format/money'
import { createExpenseBudgetBreakdownItems } from '../display'

export function SummaryCard(props: { dashboard: ExpensesDashboard }) {
  const items: BudgetBreakdownItem[] = createExpenseBudgetBreakdownItems(props.dashboard, formatEuro, {
    compactFundLabel: true,
    includePercent: true,
    highlightOverBudget: true,
    budgetSource: 'ratio',
  })

  return (
    <Panel title="Summary">
      <BudgetBreakdown
        remaining={formatEuro(props.dashboard.balance.remaining)}
        income={formatEuro(props.dashboard.balance.income)}
        items={items}
      />
    </Panel>
  )
}
