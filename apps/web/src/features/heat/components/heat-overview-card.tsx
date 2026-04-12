import { DataCell, DataRow } from '../../../components/ui/DataRow'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { Stat } from '../../../components/ui/Stat'
import type { HeatDashboard } from '../../../core/api/types'
import { createHeatHeroModel, createHeatSeasonComparisonModel } from '../display'

export function HeatOverviewCard(props: { dashboard: HeatDashboard }) {
  const hero = () => createHeatHeroModel(props.dashboard)
  const comparison = () => createHeatSeasonComparisonModel(props.dashboard)

  return (
    <Panel title="Heating" description={<Pill tone={hero().statusTone}>{hero().statusLabel}</Pill>}>
      <DataRow variant="header-spaced">
        <Stat
          label="Days since refill"
          labelPosition="bottom"
          value={hero().daysSinceRefill}
          unit="days"
          variant="hero"
        />
        <Stat label="Season" value={hero().seasonLabel} variant="mono" align="right" />
      </DataRow>

      <DataRow variant="footer">
        <Stat variant="inline" label="Last refill:" value={hero().latestRefillDate} />
        <Stat variant="inline" label="Temp:" value={hero().latestRefillTemperature} color="muted" />
      </DataRow>

      <DataRow variant="divided">
        <DataCell flex>
          <Stat label={comparison().currentLabel} value={comparison().currentValue} variant="lg" align="center" />
        </DataCell>
        <DataCell flex>
          <Stat label={comparison().previousLabel} value={comparison().previousValue} color="muted" variant="lg" align="center" />
        </DataCell>
      </DataRow>

      <DataRow variant="footer">
        <Stat variant="inline" label="Gap:" value={comparison().deltaValue} color={comparison().deltaColor} />
        <Stat variant="inline" label="Trend:" value={comparison().deltaTrend} color="muted" />
      </DataRow>
    </Panel>
  )
}
