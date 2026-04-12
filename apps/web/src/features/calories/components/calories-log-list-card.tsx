import { createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { DeleteConfirmSheet } from '../../../components/ui/DeleteConfirmSheet'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import { TrashIcon } from '../../../components/ui/TrashIcon'
import type { CalorieLog } from '../../../core/api/types'
import { formatCount } from '../../../core/format/number'
import {
  createCalorieLogDeleteDescription,
  formatCalorieLogTime,
} from '../display'

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
                  icon={<TrashIcon />}
                  textDanger
                  onClick={() => setPendingDeleteLog(log)}
                />
              }
            />
          ))}
        </List>
      </Panel>

      <DeleteConfirmSheet
        open={Boolean(pendingDeleteLog())}
        title="Delete meal"
        description={pendingDeleteLog() ? createCalorieLogDeleteDescription(pendingDeleteLog()!) : ''}
        eyebrow="Meal log"
        confirmLabel="Delete meal"
        busy={props.deletingId === pendingDeleteLog()?.id}
        onCancel={() => setPendingDeleteLog(undefined)}
        onConfirm={confirmDelete}
      />
    </>
  )
}
