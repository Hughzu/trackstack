import { createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { DeleteConfirmSheet } from '../../../components/ui/DeleteConfirmSheet'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import { TrashIcon } from '../../../components/ui/TrashIcon'
import type { Refill } from '../../../core/api/types'
import { createHeatDeleteDescription, formatHeatBags, formatHeatDate, formatHeatTemperature, formatHeatWeight } from '../display'

export function HeatHistoryCard(props: { history: Refill[], deletingId?: string, onDelete: (refill: Refill) => void | Promise<void> }) {
  const [pendingDeleteRefill, setPendingDeleteRefill] = createSignal<Refill | undefined>()

  const confirmDelete = async () => {
    const refill = pendingDeleteRefill()
    if (!refill) return

    await props.onDelete(refill)
    setPendingDeleteRefill(undefined)
  }

  return (
    <>
      <Panel title="History" description={`${props.history.length} refills logged`}>
        <List variant="flush" emptyMessage="No refills yet.">
          {props.history.map((item) => (
            <ListItem
              id={item.id}
              title={formatHeatDate(item.date)}
              subtitle={
                <ListMeta>
                  <span>{formatHeatWeight(item.weightKg)}</span>
                  <ListMetaDivider />
                  <span>{formatHeatTemperature(item.temperature)}</span>
                  {item.season ? (
                    <>
                      <ListMetaDivider />
                      <span>{item.season}</span>
                    </>
                  ) : null}
                </ListMeta>
              }
              value={formatHeatBags(item.bags)}
              valueStyle="mono"
              action={
                <IconButton
                  ariaLabel={`Delete refill ${formatHeatDate(item.date)}`}
                  disabled={props.deletingId === item.id}
                  icon={<TrashIcon />}
                  textDanger
                  onClick={() => setPendingDeleteRefill(item)}
                />
              }
            />
          ))}
        </List>
      </Panel>

      <DeleteConfirmSheet
        open={Boolean(pendingDeleteRefill())}
        title="Delete refill"
        description={pendingDeleteRefill() ? createHeatDeleteDescription(pendingDeleteRefill()!) : ''}
        eyebrow="Refill history"
        confirmLabel="Delete refill"
        busy={props.deletingId === pendingDeleteRefill()?.id}
        onCancel={() => setPendingDeleteRefill(undefined)}
        onConfirm={confirmDelete}
      />
    </>
  )
}
