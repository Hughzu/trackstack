import { createMemo, createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { DeleteConfirmSheet } from '../../../components/ui/DeleteConfirmSheet'
import { FilterGroup } from '../../../components/ui/FilterGroup'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { TrashIcon } from '../../../components/ui/TrashIcon'
import type { ExpenseEntry } from '../../../core/api/types'
import { formatEuro } from '../../../core/format/money'
import { createExpenseHistoryFilterOptions, getExpenseCategoryMeta } from '../display'

const formatDateLabel = (dateValue: string) => {
  const date = new Date(`${dateValue}T00:00:00`)
  const today = new Date()
  const startToday = new Date(today.getFullYear(), today.getMonth(), today.getDate())
  const startDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const diffDays = Math.round((startToday.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Yesterday'
  return date.toLocaleDateString('en-GB', { day: '2-digit', month: 'short' })
}

export function HistoryCard(props: { history: ExpenseEntry[], deletingId?: string, onDelete: (entry: ExpenseEntry) => void | Promise<void> }) {
  const [filter, setFilter] = createSignal('all')
  const [pendingDeleteEntry, setPendingDeleteEntry] = createSignal<ExpenseEntry | undefined>()
  const options = createExpenseHistoryFilterOptions()

  const visibleHistory = createMemo(() => {
    if (filter() === 'all') return props.history
    return props.history.filter((entry) => entry.category === filter())
  })

  const confirmDeleteEntry = async () => {
    const entry = pendingDeleteEntry()
    if (!entry) return

    await props.onDelete(entry)
    setPendingDeleteEntry(undefined)
  }

  return (
    <>
      <Panel
        title="History"
        headerAction={<FilterGroup options={options} value={filter()} onChange={setFilter} />}
      >
        <List emptyMessage="No expenses match this filter." variant="flush">
          {visibleHistory().map((entry) => {
            const category = getExpenseCategoryMeta(entry.category)

            return (
              <ListItem
                id={entry.id}
                title={entry.title}
                subtitle={
                  <ListMeta>
                    <span>{formatDateLabel(entry.date)}</span>
                    <ListMetaDivider />
                    <Pill size="sm" tone={category.tone}>{category.compactLabel}</Pill>
                    {entry.type === 'recurring' ? <Pill size="sm" tone="neutral">AUTO</Pill> : null}
                    {entry.type === 'checklist' ? <Pill size="sm" tone="neutral">CHECK</Pill> : null}
                  </ListMeta>
                }
                value={`-${formatEuro(entry.amount)}`}
                valueTone="danger"
                valueStyle="mono"
                action={
                  <IconButton
                    ariaLabel={`Delete ${entry.title}`}
                    disabled={props.deletingId === entry.id}
                    icon={<TrashIcon />}
                    textDanger
                    onClick={() => setPendingDeleteEntry(entry)}
                  />
                }
              />
            )
          })}
        </List>
      </Panel>

      <DeleteConfirmSheet
        open={Boolean(pendingDeleteEntry())}
        title="Delete expense"
        description={pendingDeleteEntry() ? `Delete ${pendingDeleteEntry()!.title} for ${formatEuro(pendingDeleteEntry()!.amount)} from ${formatDateLabel(pendingDeleteEntry()!.date)}? This removes it from the current history.` : ''}
        eyebrow="History entry"
        confirmLabel="Delete expense"
        busy={props.deletingId === pendingDeleteEntry()?.id}
        onCancel={() => setPendingDeleteEntry(undefined)}
        onConfirm={confirmDeleteEntry}
      />
    </>
  )
}
