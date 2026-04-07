import { createResource, createSignal, createMemo, Show, Suspense } from 'solid-js'

import { ActionButton, CheckToggleButton, FloatingActionGroup, IconButton } from '../../components/ui/ActionButton'
import { BudgetBreakdown, type BudgetBreakdownItem, SkeletonBudgetBreakdown } from '../../components/ui/BudgetBreakdown'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { FilterGroup, type FilterOption } from '../../components/ui/FilterGroup'
import { List, ListItem, ListMeta, ListMetaDivider, SkeletonListItem } from '../../components/ui/List'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { CounterPill, Pill } from '../../components/ui/Pill'
import { SkeletonPanel } from '../../components/ui/Skeleton'
import { Stat } from '../../components/ui/Stat'
import type { ExpenseChecklistItem, ExpenseEntry, ExpensesDashboard } from '../../core/api/types'
import { authState } from '../../core/auth/state'
import { readExpensesDashboard } from './api/client'

const euroFormatter = new Intl.NumberFormat('en-IE', {
  style: 'currency',
  currency: 'EUR',
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

const formatEuro = (value: number) => euroFormatter.format(value)

const formatDateLabel = (dateValue: string) => {
  const date = new Date(`${dateValue}T00:00:00`)
  const today = new Date()
  const startToday = new Date(today.getFullYear(), today.getMonth(), today.getDate())
  const startDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const diffDays = Math.round((startToday.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Yesterday'
  return date.toLocaleDateString('en-GB', { day: '2-digit', month: 'short' })
}

const formatCategory = (category: string) => {
  if (category === 'fun') return 'Fun'
  if (category === 'future') return 'Future'
  return 'Fund.'
}

const getCategoryTone = (category: string) => {
  if (category === 'fun') return 'warning'
  if (category === 'future') return 'success'
  return 'danger'
}

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

function SummaryCard(props: { dashboard: ExpensesDashboard }) {
  const getRatioData = (id: string) => {
    const ratio = props.dashboard.ratios.find(r => r.categoryId === id)
    if (!ratio) return { percent: 0, budget: 0, over: false }
    return ratio
  }

  const fund = getRatioData('fund')
  const fun = getRatioData('fun')
  const future = getRatioData('future')
  const items: BudgetBreakdownItem[] = [
    {
      label: 'Fund.',
      value: formatEuro(props.dashboard.spent.fund),
      subtext: `/ ${formatEuro(fund.budget)} (${fund.percent}%)`,
      percent: fund.percent,
      color: fund.over ? 'danger' : 'main',
      barColor: 'danger',
    },
    {
      label: 'Fun',
      value: formatEuro(props.dashboard.spent.fun),
      subtext: `/ ${formatEuro(fun.budget)} (${fun.percent}%)`,
      percent: fun.percent,
      color: fun.over ? 'danger' : 'main',
      barColor: 'warning',
    },
    {
      label: 'Future',
      value: formatEuro(props.dashboard.spent.future),
      subtext: `/ ${formatEuro(future.budget)} (${future.percent}%)`,
      percent: future.percent,
      color: future.over ? 'danger' : 'main',
      barColor: 'success',
    },
  ]

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

function ObligationsCard(props: { obligations: ExpenseChecklistItem[] }) {
  return (
    <Panel
      title="Obligations"
      description={props.obligations.length > 0 ? <CounterPill value={props.obligations.length} label="Left" /> : undefined}
      collapsibleId="expenses_obligations"
    >
      <List emptyMessage="All obligations paid for this month!" variant="flush">
        {props.obligations.map((item) => (
          <ListItem
            title={item.title}
            value={`-${formatEuro(item.amount)}`}
            valueTone="danger"
            valueStyle="mono"
            prefix={<CheckToggleButton />}
          />
        ))}
      </List>
    </Panel>
  )
}

const DeleteIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-4">
    <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
  </svg>
)

function HistoryCard(props: { history: ExpenseEntry[] }) {
  const [filter, setFilter] = createSignal('all')
  
  const options: FilterOption[] = [
    { value: 'all', label: 'All' },
    { value: 'fund', label: 'Fund.', tone: 'danger' },
    { value: 'fun', label: 'Fun', tone: 'warning' },
    { value: 'future', label: 'Future', tone: 'success' },
  ]

  const visibleHistory = createMemo(() => {
    if (filter() === 'all') return props.history
    return props.history.filter(tx => tx.category === filter())
  })

  return (
    <Panel 
      title="History" 
      headerAction={<FilterGroup options={options} value={filter()} onChange={setFilter} />}
    >
      <List emptyMessage="No expenses match this filter." variant="flush">
        {visibleHistory().map((tx) => (
          <ListItem
            title={tx.title}
            subtitle={
              <ListMeta>
                <span>{formatDateLabel(tx.date)}</span>
                <ListMetaDivider />
                <Pill size="sm" tone={getCategoryTone(tx.category) as any}>{formatCategory(tx.category)}</Pill>
                {tx.type === 'recurring' ? <Pill size="sm" tone="neutral">AUTO</Pill> : null}
                {tx.type === 'checklist' ? <Pill size="sm" tone="neutral">CHECK</Pill> : null}
              </ListMeta>
            }
            value={`-${formatEuro(tx.amount)}`}
            valueTone="danger"
            valueStyle="mono"
            action={
              <IconButton icon={<DeleteIcon />} textDanger />
            }
          />
        ))}
      </List>
    </Panel>
  )
}

function DashboardSkeleton() {
  return (
    <>
      <SkeletonBudgetBreakdown />
      <SkeletonPanel titleVariant="md">
        <List variant="flush">
          <SkeletonListItem />
          <SkeletonListItem />
        </List>
      </SkeletonPanel>
    </>
  )
}

export default function ExpensesPage() {
  const [dashboard] = createResource(readyKey, readExpensesDashboard)

  const mockObligations = [
    { id: 'mock-1', title: 'Car Insurance', amount: 125.50, category: 'fund' },
    { id: 'mock-2', title: 'Gym Membership', amount: 35.00, category: 'fund' }
  ] as ExpenseChecklistItem[]

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <DataRow variant="header">
        <Stat label="Period" value={dashboard()?.periodKey || '...'} variant="md" />
        <ActionButton tone="ghost" busy={dashboard.loading && !dashboard()}>Close month</ActionButton>
      </DataRow>

      <Suspense fallback={<DashboardSkeleton />}>
        <Show when={dashboard()} fallback={<Notice tone="error" message="Expenses dashboard failed to load." />}>
          {(data) => (
            <>
              <SummaryCard dashboard={data()} />
              <ObligationsCard obligations={data().pendingObligations?.length ? data().pendingObligations : mockObligations} />
              <HistoryCard history={data().history ?? []} />
            </>
          )}
        </Show>
      </Suspense>

      <FloatingActionGroup>
        <DataCell flex><ActionButton block tone="ghost">Settings</ActionButton></DataCell>
        <DataCell flex><ActionButton block tone="primary">Add Expense</ActionButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
