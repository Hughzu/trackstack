import { createMemo, createSignal } from 'solid-js'

import { IconButton } from '../../../components/ui/ActionButton'
import { ConfirmSheet } from '../../../components/ui/ConfirmSheet'
import { FilterGroup } from '../../../components/ui/FilterGroup'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
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

function DeleteIcon() {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" width="16" height="16" aria-hidden="true">
      <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
    </svg>
  )
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
                    icon={<DeleteIcon />}
                    textDanger
                    onClick={() => setPendingDeleteEntry(entry)}
                  />
                }
              />
            )
          })}
        </List>
      </Panel>

      <ConfirmSheet
        open={Boolean(pendingDeleteEntry())}
        title="Delete expense"
        description={pendingDeleteEntry() ? `Delete ${pendingDeleteEntry()!.title} for ${formatEuro(pendingDeleteEntry()!.amount)} from ${formatDateLabel(pendingDeleteEntry()!.date)}? This removes it from the current history.` : ''}
        confirmLabel="Delete expense"
        confirmTone="danger"
        eyebrow="History entry"
        busy={props.deletingId === pendingDeleteEntry()?.id}
        onCancel={() => setPendingDeleteEntry(undefined)}
        onConfirm={confirmDeleteEntry}
      />
    </>
  )
}
