import { createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { ConfirmSheet } from '../../../components/ui/ConfirmSheet'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import type { Refill } from '../../../core/api/types'
import { createHeatDeleteDescription, formatHeatBags, formatHeatDate, formatHeatTemperature, formatHeatWeight } from '../display'

function DeleteIcon() {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" width="16" height="16" aria-hidden="true">
      <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
    </svg>
  )
}

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
                  icon={<DeleteIcon />}
                  textDanger
                  onClick={() => setPendingDeleteRefill(item)}
                />
              }
            />
          ))}
        </List>
      </Panel>

      <ConfirmSheet
        open={Boolean(pendingDeleteRefill())}
        title="Delete refill"
        description={pendingDeleteRefill() ? createHeatDeleteDescription(pendingDeleteRefill()!) : ''}
        confirmLabel="Delete refill"
        confirmTone="danger"
        eyebrow="Refill history"
        busy={props.deletingId === pendingDeleteRefill()?.id}
        onCancel={() => setPendingDeleteRefill(undefined)}
        onConfirm={confirmDelete}
      />
    </>
  )
}
