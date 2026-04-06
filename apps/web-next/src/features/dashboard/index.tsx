import { createSignal, onMount, Show } from 'solid-js'

import { AppShell } from '../../components/ui/AppShell'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { Panel } from '../../components/ui/Panel'
import { DataRow, DataCell } from '../../components/ui/DataRow'
import { Stat, SkeletonStat } from '../../components/ui/Stat'
import { ProgressBar } from '../../components/ui/ProgressBar'
import { SkeletonPanel, SkeletonProgressBar } from '../../components/ui/Skeleton'
import { SkeletonFlexRow, SkeletonTextInline } from '../../components/ui/SkeletonFlexRow'

function ExpensesCard() {
  return (
    <Panel title="Expenses" href="/expenses">
      <DataRow variant="header">
        <Stat label="Remaining" value="€1,240.50" variant="lg" />
        <Stat label="Income" value="€3,500.00" variant="mono" align="right" />
      </DataRow>
      
      <ProgressBar 
        segments={[
          { percent: 30, color: 'danger' },
          { percent: 40, color: 'warning' },
          { percent: 10, color: 'success' },
        ]} 
      />

      <DataRow variant="divided">
        <DataCell flex>
          <Stat label="Fund" value="€450" subtext="/ €1,000" variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <Stat label="Fun" value="€600" subtext="/ €1,500" variant="sm" align="center" />
        </DataCell>
        <DataCell flex>
          <Stat label="Future" value="€150" subtext="/ €1,000" variant="sm" align="center" />
        </DataCell>
      </DataRow>
    </Panel>
  )
}

function CaloriesCard() {
  return (
    <Panel title="Daily Intake" href="/calories">
      <DataRow variant="header-end">
        <Stat label="Consumed" value="1,840" unit="/ 2,400 kcal" variant="lg" />
        <Stat label="Protein" value="145g" unit="/ 160g" variant="md" color="accent" align="right" />
      </DataRow>
      <ProgressBar segments={[{ percent: 76, color: 'accent' }]} />
    </Panel>
  )
}

function HeatCard() {
  return (
    <Panel title="Heating" href="/heat">
      <DataRow variant="header-spaced">
        <Stat label="Since last refill" labelPosition="bottom" value="14" unit="Days" variant="hero" />
        <Stat label="Season" value="2025/2026" variant="mono" align="right" />
      </DataRow>

      <DataRow variant="footer">
        <Stat variant="inline" label="This Season:" value="30 bags" />
        <Stat variant="inline" label="Last Season:" value="28 bags" color="muted" />
      </DataRow>
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
      
      <SkeletonFlexRow>
        <SkeletonTextInline variant="md" />
        <SkeletonTextInline variant="md" />
        <SkeletonTextInline variant="md" />
      </SkeletonFlexRow>

      <DataRow variant="divided">
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
          <SkeletonTextInline variant="sm" />
        </DataCell>
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
          <SkeletonTextInline variant="sm" />
        </DataCell>
        <DataCell flex>
          <SkeletonStat variant="sm" align="center" />
          <SkeletonTextInline variant="sm" />
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
        <SkeletonTextInline variant="lg" />
        <SkeletonTextInline variant="lg" />
      </DataRow>
    </SkeletonPanel>
  )
}

export default function Dashboard() {
  const [isLoading, setIsLoading] = createSignal(true)

  onMount(() => {
    // Mock network request to show off the skeleton loading 
    const timer = setTimeout(() => {
      setIsLoading(false)
    }, 1500)

    return () => clearTimeout(timer)
  })

  return (
    <AppShell currentDomain="home">
      <Show 
        when={!isLoading()} 
        fallback={
          <ContentDeck layout="stacked">
            <ExpensesSkeleton />
            <CaloriesSkeleton />
            <HeatSkeleton />
          </ContentDeck>
        }
      >
        <ContentDeck layout="stacked" animate>
          <ExpensesCard />
          <CaloriesCard />
          <HeatCard />
        </ContentDeck>
      </Show>
    </AppShell>
  )
}
