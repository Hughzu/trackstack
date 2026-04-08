import { For } from 'solid-js'

import { ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { FormBackLink } from '../../components/ui/Form'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../components/ui/List'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { ProgressBar } from '../../components/ui/ProgressBar'

const quickAddMeals: Array<{ id: string, title: string }> = [
  { id: '1', title: 'Oats + whey' },
  { id: '2', title: 'Chicken rice' },
  { id: '3', title: 'Greek yogurt' },
]

const todayLogs: Array<{
  id: string
  title: string
  time: string
  calories: number
  proteinGrams: number
}> = [
  { id: '1', title: 'Oats + whey', time: '08:10', calories: 520, proteinGrams: 38 },
  { id: '2', title: 'Chicken rice', time: '13:05', calories: 760, proteinGrams: 54 },
  { id: '3', title: 'Greek yogurt', time: '16:40', calories: 240, proteinGrams: 20 },
]

export default function CaloriesPage() {
  const consumedCalories = 1520
  const targetCalories = 2400
  const proteinGrams = 112
  const targetProteinGrams = 180
  const progressPercent = Math.min(Math.round((consumedCalories / targetCalories) * 100), 100)

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <DataRow variant="header">
        <FormBackLink href="/">Back</FormBackLink>
      </DataRow>

      <Panel title="Today" description={<Pill tone="neutral">{targetCalories - consumedCalories} kcal left</Pill>}>
        <div class="flex flex-col gap-5">
          <div class="flex flex-col gap-4 border-b border-border/50 pb-4 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Daily intake</div>
              <div class="mt-1 flex items-baseline gap-2">
                <div class="text-4xl font-bold tracking-tight text-text-main sm:text-5xl">{consumedCalories}</div>
                <div class="text-sm font-mono text-text-muted">/ {targetCalories} kcal</div>
              </div>
            </div>

            <div class="min-w-0 sm:text-right">
              <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Protein target</div>
              <div class="mt-1 text-lg font-semibold text-text-main">
                {proteinGrams}g
                <span class="ml-1 text-sm font-normal text-text-muted">/ {targetProteinGrams}g</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col gap-3 border-b border-border/50 pb-4">
            <ProgressBar segments={[{ percent: progressPercent, color: 'accent' }]} />
            <div class="flex items-center justify-between gap-3 text-xs font-semibold uppercase tracking-[0.2em] text-text-muted">
              <span>Energy progress</span>
              <span>{progressPercent}%</span>
            </div>
          </div>
        </div>
      </Panel>

      <section class="flex flex-col gap-2">
        <div class="px-1 text-sm font-semibold text-text-muted">Quick add</div>
        <div class="quick-add-scroll relative overflow-x-auto pb-2">
          <div class="flex min-w-max gap-2 px-1 pb-0 pt-1">
          <For each={quickAddMeals}>
            {(meal) => (
              <button
                type="button"
                class="inline-flex h-[38px] min-w-[8rem] items-center justify-center rounded-full border border-border/70 bg-panel/60 px-4 py-2 text-sm font-medium text-text-main transition hover:border-accent/40 hover:text-accent"
              >
                {meal.title}
              </button>
            )}
          </For>
          </div>
        </div>
      </section>

      <Panel title="Meal log" description={`${todayLogs.length} entries today`}>
        <List variant="flush">
          <For each={todayLogs}>
            {(log) => (
              <ListItem
                id={log.id}
                title={log.title}
                subtitle={
                  <ListMeta>
                    <span>{log.time}</span>
                    <ListMetaDivider />
                    <span>{log.proteinGrams}g protein</span>
                  </ListMeta>
                }
                value={`${log.calories} kcal`}
                valueStyle="mono"
              />
            )}
          </For>
        </List>
      </Panel>

      <FloatingActionGroup>
        <DataCell flex><ActionLinkButton href="/calories/settings" block tone="ghost">Settings</ActionLinkButton></DataCell>
        <DataCell flex><ActionLinkButton href="/calories/new" block>Add Meal</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
