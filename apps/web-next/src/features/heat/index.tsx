import { ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell } from '../../components/ui/DataRow'
import { FormBackLink } from '../../components/ui/Form'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../components/ui/List'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { readHeatMockState } from './mock-state'

const wholeNumberFormatter = new Intl.NumberFormat('en-IE')

const formatCount = (value: number) => wholeNumberFormatter.format(value)
const formatDate = (value: string) => new Intl.DateTimeFormat('en-IE', { month: 'short', day: 'numeric' }).format(new Date(value))
const formatWeight = (value: number) => `${value} kg`
const formatBags = (value: number) => `${formatCount(value)} bag${value === 1 ? '' : 's'}`
const formatTemperature = (value: number | null) => (value == null ? 'No temperature' : `${value} C`)

export default function HeatPage() {
  const heat = readHeatMockState()
  const isWarning = heat.daysSinceRefill > 14
  const lastRefill = heat.history[0]

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <div class="flex items-center justify-between gap-3">
        <FormBackLink href="/">Back</FormBackLink>
        <Pill tone="neutral">{heat.seasonLabel}</Pill>
      </div>

      <Panel title="Heating" description={isWarning ? <Pill tone="warning">Refill soon</Pill> : <Pill tone="success">Fresh refill pace</Pill>}>
        <div class="flex flex-col gap-5">
          <div class="border-b border-border/50 pb-4">
            <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Days since refill</div>
                <div class="mt-1 flex items-baseline gap-2">
                  <div class="text-4xl font-bold tracking-tight text-text-main sm:text-5xl">{formatCount(heat.daysSinceRefill)}</div>
                  <div class="text-sm font-mono text-text-muted">days</div>
                </div>
              </div>

              <div class="text-sm text-text-muted sm:max-w-[12rem] sm:text-right">
                {lastRefill ? `Last refill ${formatDate(lastRefill.date)} at ${formatTemperature(lastRefill.temperature)}.` : 'No refill history yet.'}
              </div>
            </div>
          </div>

          <div class="divide-y divide-border/40 rounded-2xl border border-border/50 bg-panel/35 px-4">
            <SeasonLine label="This season" value={formatBags(heat.totals.thisSeason)} emphasis />
            <SeasonLine label="Last season by now" value={formatBags(heat.totals.lastSeasonToDate)} />
            <SeasonLine label="Last season total" value={formatBags(heat.totals.lastSeasonTotal)} />
          </div>
        </div>
      </Panel>

      <Panel title="History" description={`${heat.history.length} refills logged`}>
        <List variant="flush" emptyMessage="No refills yet.">
          {heat.history.map((item) => (
            <ListItem
              id={item.id}
              title={formatDate(item.date)}
              subtitle={
                <ListMeta>
                  <span>{formatWeight(item.weightKg)}</span>
                  {item.temperature != null ? (
                    <>
                      <ListMetaDivider />
                      <span>{item.temperature} C</span>
                    </>
                  ) : null}
                  <ListMetaDivider />
                  <span>{item.season}</span>
                </ListMeta>
              }
              value={formatBags(item.bags)}
              valueStyle="mono"
            />
          ))}
        </List>
      </Panel>

      <FloatingActionGroup>
        <DataCell flex><ActionLinkButton href="/heat/new" block>Add refill</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}

function SeasonLine(props: { label: string, value: string, emphasis?: boolean }) {
  return (
    <div class="flex items-center justify-between gap-3 py-3 first:pt-4 last:pb-4">
      <div class="text-sm text-text-muted">{props.label}</div>
      <div class={`text-sm font-semibold ${props.emphasis ? 'text-text-main' : 'text-text-muted'}`}>{props.value}</div>
    </div>
  )
}
