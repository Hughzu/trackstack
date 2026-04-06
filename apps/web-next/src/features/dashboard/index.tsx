import { createMemo, createResource, Show } from 'solid-js'

import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { ProgressBar } from '../../components/ui/ProgressBar'
import { SkeletonPanel, SkeletonProgressBar } from '../../components/ui/Skeleton'
import { Stat, SkeletonStat } from '../../components/ui/Stat'
import type { CaloriesDashboard, ExpensesDashboard, HeatDashboard } from '../../core/api/types'
import { authState } from '../../core/auth/state'
import { readCaloriesDashboard } from '../calories/api/client'
import { readExpensesDashboard } from '../expenses/api/client'
import { readHeatDashboard } from '../heat/api/client'

const euroFormatter = new Intl.NumberFormat('en-IE', {
  style: 'currency',
  currency: 'EUR',
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

const wholeNumberFormatter = new Intl.NumberFormat('en-IE')

const formatEuro = (value: number) => euroFormatter.format(value)

const formatCount = (value: number) => wholeNumberFormatter.format(value)

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

function ExpensesCard(props: { dashboard: ExpensesDashboard }) {
  return (
    <Panel title="Expenses" href="/expenses">
      <DataRow variant="header">
        <Stat label="Remaining" value={formatEuro(props.dashboard.balance.remaining)} variant="lg" />
        <Stat label="Income" value={formatEuro(props.dashboard.balance.income)} variant="mono" align="right" />
      </DataRow>

      <ProgressBar
        segments={props.dashboard.ratios.map((ratio) => ({
          percent: ratio.percent,
          color: ratio.categoryId === 'fun' ? 'warning' : ratio.categoryId === 'future' ? 'success' : 'danger',
        }))}
      />

      <DataRow variant="divided">
        <DataCell flex>
          <Stat label="Fund" value={formatEuro(props.dashboard.spent.fund)} subtext={`/ ${formatEuro(props.dashboard.budget.fund)}`} variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <Stat label="Fun" value={formatEuro(props.dashboard.spent.fun)} subtext={`/ ${formatEuro(props.dashboard.budget.fun)}`} variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <Stat label="Future" value={formatEuro(props.dashboard.spent.future)} subtext={`/ ${formatEuro(props.dashboard.budget.future)}`} variant="sm" align="center" />
        </DataCell>
      </DataRow>
    </Panel>
  )
}

function CaloriesCard(props: { dashboard: CaloriesDashboard }) {
  const consumedPercent = createMemo(() => {
    const targetCalories = props.dashboard.summary.target.targetCalories

    if (!targetCalories) {
      return 0
    }

    return Math.min((props.dashboard.summary.consumedCalories / targetCalories) * 100, 100)
  })

  return (
    <Panel title="Daily Intake" href="/calories">
      <DataRow variant="header-end">
        <Stat
          label="Consumed"
          value={formatCount(props.dashboard.summary.consumedCalories)}
          unit={`/ ${formatCount(props.dashboard.summary.target.targetCalories)} kcal`}
          variant="lg"
        />
        <Stat
          label="Protein"
          value={`${formatCount(props.dashboard.summary.proteinGrams)}g`}
          unit={`/ ${formatCount(props.dashboard.summary.target.targetProteinGrams)}g`}
          variant="md"
          color="accent"
          align="right"
        />
      </DataRow>
      <ProgressBar segments={[{ percent: consumedPercent(), color: 'accent' }]} />
    </Panel>
  )
}

function HeatCard(props: { dashboard: HeatDashboard }) {
  return (
    <Panel title="Heating" href="/heat">
      <DataRow variant="header-spaced">
        <Stat
          label="Since last refill"
          labelPosition="bottom"
          value={formatCount(props.dashboard.daysSinceRefill)}
          unit="Days"
          variant="hero"
        />
        <Stat label="Season" value={props.dashboard.seasonSnapshot.seasonLabel} variant="mono" align="right" />
      </DataRow>

      <DataRow variant="footer">
        <Stat variant="inline" label="This Season:" value={`${formatCount(props.dashboard.seasonSnapshot.seasonToDate)} bags`} />
        <Stat
          variant="inline"
          label="Last Season:"
          value={`${formatCount(props.dashboard.seasonSnapshot.lastSeasonToDate)} bags`}
          color="muted"
        />
      </DataRow>
    </Panel>
  )
}

function ErrorCard(props: { title: string; message: string }) {
  return (
    <Panel title={props.title}>
      <Notice tone="error" message={props.message} />
    </Panel>
  )
}

function ExpensesSkeleton() {
  return (
    <SkeletonPanel titleVariant="md">
      <DataRow variant="header">
        <SkeletonStat variant="lg" />
        <SkeletonStat variant="mono" align="right" />
      </DataRow>
      <SkeletonProgressBar />

      <DataRow variant="divided">
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
        </DataCell>
      </DataRow>
    </SkeletonPanel>
  )
}

function CaloriesSkeleton() {
  return (
    <SkeletonPanel titleVariant="lg">
      <DataRow variant="header-end">
        <SkeletonStat variant="lg" />
        <SkeletonStat variant="md" align="right" />
      </DataRow>
      <SkeletonProgressBar />
    </SkeletonPanel>
  )
}

function HeatSkeleton() {
  return (
    <SkeletonPanel titleVariant="md">
      <DataRow variant="header-spaced">
        <SkeletonStat variant="hero" labelPosition="bottom" />
        <SkeletonStat variant="mono" align="right" />
      </DataRow>
      <DataRow variant="footer">
        <SkeletonStat variant="inline" />
        <SkeletonStat variant="inline" />
      </DataRow>
    </SkeletonPanel>
  )
}

export default function Dashboard() {
  const [expenses] = createResource(readyKey, readExpensesDashboard)
  const [calories] = createResource(readyKey, readCaloriesDashboard)
  const [heat] = createResource(readyKey, readHeatDashboard)

  const allFailed = createMemo(() => Boolean(expenses.error && calories.error && heat.error))

  return (
    <>
      <Show when={allFailed()}>
        <Notice tone="error" message="Unable to load the dashboard from the monolith right now." />
      </Show>

      <ContentDeck layout="stacked" animate>
        <Show when={expenses()} fallback={expenses.error ? <ErrorCard title="Expenses" message="Expenses dashboard failed to load." /> : <ExpensesSkeleton />}>
          {(dashboard) => <ExpensesCard dashboard={dashboard()} />}
        </Show>

        <Show when={calories()} fallback={calories.error ? <ErrorCard title="Daily Intake" message="Calories dashboard failed to load." /> : <CaloriesSkeleton />}>
          {(dashboard) => <CaloriesCard dashboard={dashboard()} />}
        </Show>

        <Show when={heat()} fallback={heat.error ? <ErrorCard title="Heating" message="Heat dashboard failed to load." /> : <HeatSkeleton />}>
          {(dashboard) => <HeatCard dashboard={dashboard()} />}
        </Show>
      </ContentDeck>
    </>
  )
}
