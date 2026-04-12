import { createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { ConfirmSheet } from '../../../components/ui/ConfirmSheet'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import type { CalorieLog } from '../../../core/api/types'
import { formatCount } from '../../../core/format/number'
import {
  createCalorieLogDeleteDescription,
  formatCalorieLogTime,
} from '../display'

function DeleteIcon() {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" width="16" height="16" aria-hidden="true">
      <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
    </svg>
  )
}

export function CaloriesLogListCard(props: { dashboardLogs: CalorieLog[], deletingId?: string, onDelete: (log: CalorieLog) => void | Promise<void> }) {
  const [pendingDeleteLog, setPendingDeleteLog] = createSignal<CalorieLog | undefined>()

  const confirmDelete = async () => {
    const log = pendingDeleteLog()
    if (!log) return

    await props.onDelete(log)
    setPendingDeleteLog(undefined)
  }

  return (
    <>
      <Panel title="Meal log" description={`${props.dashboardLogs.length} entries today`}>
        <List variant="flush" emptyMessage="No logs yet for today.">
          {props.dashboardLogs.map((log) => (
            <ListItem
              id={log.id}
              title={log.title?.trim() || 'Untitled'}
              subtitle={
                <ListMeta>
                  <span>{formatCalorieLogTime(log.dateTime)}</span>
                  <ListMetaDivider />
                  <span>{formatCount(log.proteinGrams)}g protein</span>
                </ListMeta>
              }
              value={`${formatCount(log.calories)} kcal`}
              valueStyle="mono"
              action={
                <IconButton
                  ariaLabel={`Delete ${log.title?.trim() || 'meal log'}`}
                  disabled={props.deletingId === log.id}
                  icon={<DeleteIcon />}
                  textDanger
                  onClick={() => setPendingDeleteLog(log)}
                />
              }
            />
          ))}
        </List>
      </Panel>

      <ConfirmSheet
        open={Boolean(pendingDeleteLog())}
        title="Delete meal"
        description={pendingDeleteLog() ? createCalorieLogDeleteDescription(pendingDeleteLog()!) : ''}
        confirmLabel="Delete meal"
        confirmTone="danger"
        eyebrow="Meal log"
        busy={props.deletingId === pendingDeleteLog()?.id}
        onCancel={() => setPendingDeleteLog(undefined)}
        onConfirm={confirmDelete}
      />
    </>
  )
}
