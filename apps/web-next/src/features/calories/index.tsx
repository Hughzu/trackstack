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
  const carbGrams = 138
  const fatGrams = 44
  const progressPercent = Math.min(Math.round((consumedCalories / targetCalories) * 100), 100)
  const macroTotal = carbGrams + fatGrams + proteinGrams

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <DataRow variant="header">
        <FormBackLink href="/">Back</FormBackLink>
      </DataRow>

      <Panel title="Daily intake">
        <div class="flex flex-col gap-5">
          <div class="flex items-start justify-between gap-4">
            <div>
              <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Today</div>
              <div class="mt-1 flex items-end gap-2">
                <div class="text-4xl font-bold tracking-tight text-text-main sm:text-5xl">{consumedCalories}</div>
                <div class="pb-1 text-sm font-mono text-text-muted">/ {targetCalories} kcal</div>
              </div>
            </div>

            <div class="text-right">
              <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Protein</div>
              <div class="mt-1 text-lg font-semibold text-text-main">
                {proteinGrams}g
                <span class="ml-1 text-sm font-normal text-text-muted">/ {targetProteinGrams}g</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col gap-2">
            <ProgressBar segments={[{ percent: progressPercent, color: 'accent' }]} />
            <div class="flex justify-end">
              <Pill tone="neutral">{targetCalories - consumedCalories} left</Pill>
            </div>
          </div>

          <div class="grid grid-cols-3 divide-x divide-border/30 border-y border-border/40 py-3 text-center">
            <MacroStat
              label="Carbs"
              value={`${carbGrams}g`}
              colorClass="bg-orange-400/70"
              percent={macroTotal ? Math.round((carbGrams / macroTotal) * 100) : 0}
            />
            <MacroStat
              label="Fat"
              value={`${fatGrams}g`}
              colorClass="bg-yellow-400/70"
              percent={macroTotal ? Math.round((fatGrams / macroTotal) * 100) : 0}
            />
            <MacroStat
              label="Protein"
              value={`${proteinGrams}g`}
              colorClass="bg-emerald-400/70"
              percent={macroTotal ? Math.round((proteinGrams / macroTotal) * 100) : 0}
            />
          </div>
        </div>
      </Panel>

      <section class="flex flex-col gap-3">
        <div class="px-1 text-sm font-semibold text-text-muted">Quick add</div>
        <div class="flex flex-wrap gap-2">
          <For each={quickAddMeals}>
            {(meal) => (
              <button
                type="button"
                class="inline-flex min-h-10 items-center rounded-full border border-border/70 bg-panel/60 px-4 py-2 text-sm font-medium text-text-main transition hover:border-accent/40 hover:text-accent"
              >
                {meal.title}
              </button>
            )}
          </For>
        </div>
      </section>

      <Panel title="Today's logs">
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

function MacroStat(props: { label: string, value: string, colorClass: string, percent: number }) {
  return (
    <div class="px-3">
      <div class="text-[0.65rem] font-bold uppercase tracking-[0.2em] text-text-muted">{props.label}</div>
      <div class="mt-1 text-sm font-semibold text-text-main">{props.value}</div>
      <div class="mx-auto mt-2 h-1 w-full max-w-14 rounded-full bg-border/50">
        <div class={`h-1 rounded-full ${props.colorClass}`} style={{ width: `${props.percent}%` }} />
      </div>
    </div>
  )
}
