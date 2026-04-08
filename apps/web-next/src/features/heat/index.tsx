import { ActionButton, ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell } from '../../components/ui/DataRow'
import { FormBackLink } from '../../components/ui/Form'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../components/ui/List'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { Stat } from '../../components/ui/Stat'
import { readHeatMockState } from './mock-state'

const wholeNumberFormatter = new Intl.NumberFormat('en-IE')

const formatCount = (value: number) => wholeNumberFormatter.format(value)
const formatDate = (value: string) => new Intl.DateTimeFormat('en-IE', { month: 'short', day: 'numeric' }).format(new Date(value))
const formatWeight = (value: number) => `${value} kg`
const formatBags = (value: number) => `${formatCount(value)} bag${value === 1 ? '' : 's'}`

export default function HeatPage() {
  const heat = readHeatMockState()
  const isWarning = heat.daysSinceRefill > 14

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <div class="flex items-center justify-between gap-3">
        <FormBackLink href="/">Back</FormBackLink>
        <Pill tone="neutral">{heat.seasonLabel}</Pill>
      </div>

      <Panel title="Heating" description={heat.seasonLabel}>
        <div class="grid gap-4">
          <div class="relative border-b border-border/50 pb-4">
            {isWarning ? (
              <div class="absolute right-0 top-0 flex h-3 w-3">
                <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-warning opacity-75" />
                <span class="relative inline-flex h-3 w-3 rounded-full bg-warning" />
              </div>
            ) : null}

            <Stat
              label="Days since refill"
              labelPosition="bottom"
              value={formatCount(heat.daysSinceRefill)}
              unit="Days"
              align="center"
              variant="hero"
            />
          </div>

          <div class="grid grid-cols-2 gap-3 border-b border-border/50 pb-4">
            <div class="rounded-2xl border border-border/50 bg-panel/50 px-4 py-4">
              <Stat label="This season" value={formatCount(heat.totals.thisSeason)} unit="bags" variant="lg" />
            </div>
            <div class="rounded-2xl border border-border/50 bg-panel/50 px-4 py-4">
              <Stat label="Last season by now" value={formatCount(heat.totals.lastSeasonToDate)} unit="bags" variant="lg" color="muted" />
            </div>
            <div class="col-span-2 rounded-2xl border border-border/50 bg-panel/50 px-4 py-4">
              <Stat label="Last season total" value={formatCount(heat.totals.lastSeasonTotal)} unit="bags" variant="lg" color="muted" />
            </div>
          </div>

        </div>
      </Panel>

      <section class="grid gap-2">
        <div class="px-1 text-sm font-bold text-text-muted">History</div>
        <List emptyMessage="No refills yet.">
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
      </section>

      <FloatingActionGroup>
        <DataCell flex><ActionButton tone="ghost" disabled block>Settings</ActionButton></DataCell>
        <DataCell flex><ActionLinkButton href="/heat/new" block>Add refill</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
