import { For, createMemo } from 'solid-js'

import { DataCell, DataRow } from '../../../components/ui/DataRow'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { ProgressBar } from '../../../components/ui/ProgressBar'
import { Stat } from '../../../components/ui/Stat'
import type { CaloriesDashboard } from '../../../core/api/types'
import { createCaloriesHeroModel, createCaloriesMacroItems, createCaloriesProgressSegments } from '../display'

export function CaloriesOverviewCard(props: { dashboard: CaloriesDashboard }) {
  const hero = createMemo(() => createCaloriesHeroModel(props.dashboard))
  const macroItems = createMemo(() => createCaloriesMacroItems(props.dashboard))
  const progressSegments = createMemo(() => createCaloriesProgressSegments(props.dashboard))

  return (
    <Panel title="Today" description={<Pill tone="neutral">{hero().remainingLabel}</Pill>}>
      <DataRow variant="header-end">
        <Stat
          label="Daily intake"
          value={hero().consumedCalories}
          unit={`/ ${hero().targetCalories} kcal`}
          variant="hero"
        />
        <Stat
          label="Protein target"
          value={`${hero().proteinGrams}g`}
          unit={`/ ${hero().targetProteinGrams}g`}
          variant="md"
          align="right"
        />
      </DataRow>

      <ProgressBar segments={progressSegments()} />

      <DataRow variant="footer">
        <Stat variant="inline" label="Energy progress:" value={`${hero().progressPercent}%`} />
        <Stat variant="inline" label="Target:" value={`${hero().targetCalories} kcal`} color="muted" />
      </DataRow>

      <DataRow variant="divided">
        <For each={macroItems()}>
          {(item) => (
            <DataCell flex>
              <Stat
                label={item.title}
                value={item.value}
                color={item.tone === 'success' ? 'success' : item.tone === 'warning' ? 'warning' : 'danger'}
                variant="sm"
                align="center"
              />
            </DataCell>
          )}
        </For>
      </DataRow>
    </Panel>
  )
}
