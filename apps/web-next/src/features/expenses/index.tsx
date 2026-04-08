import { createResource, createSignal, Show, Suspense } from 'solid-js'

import { ActionButton, ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ConfirmSheet } from '../../components/ui/ConfirmSheet'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { Notice } from '../../components/ui/Notice'
import { Stat } from '../../components/ui/Stat'
import type { ExpenseChecklistItem, ExpenseEntry } from '../../core/api/types'
import { authState } from '../../core/auth/state'
import { closeExpenseSheet, completeChecklistItem, deleteExpenseEntry, readExpensesDashboard } from './api/client'
import { DashboardSkeleton } from './components/dashboard-skeleton'
import { HistoryCard } from './components/history-card'
import { ObligationsCard } from './components/obligations-card'
import { SummaryCard } from './components/summary-card'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

export default function ExpensesPage() {
  const [dashboard, { refetch }] = createResource(readyKey, readExpensesDashboard)
  const [isClosingMonth, setIsClosingMonth] = createSignal(false)
  const [isConfirmingCloseMonth, setIsConfirmingCloseMonth] = createSignal(false)
  const [busyObligationId, setBusyObligationId] = createSignal<string | undefined>()
  const [deletingEntryId, setDeletingEntryId] = createSignal<string | undefined>()
  const [actionError, setActionError] = createSignal<string | undefined>()

  const refreshDashboard = async () => {
    await refetch()
  }

  const handleCloseMonth = async () => {
    if (isClosingMonth()) return

    setActionError(undefined)
    setIsClosingMonth(true)

    try {
      await closeExpenseSheet()
      setIsConfirmingCloseMonth(false)
      await refreshDashboard()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to close month')
    } finally {
      setIsClosingMonth(false)
    }
  }

  const handleCompleteObligation = async (item: ExpenseChecklistItem) => {
    if (busyObligationId()) return

    setActionError(undefined)
    setBusyObligationId(item.id)

    try {
      await completeChecklistItem({ id: item.id })
      await refreshDashboard()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to complete obligation')
    } finally {
      setBusyObligationId(undefined)
    }
  }

  const handleDeleteEntry = async (entry: ExpenseEntry) => {
    if (deletingEntryId()) return

    setActionError(undefined)
    setDeletingEntryId(entry.id)

    try {
      await deleteExpenseEntry(entry.id)
      await refreshDashboard()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to delete expense entry')
    } finally {
      setDeletingEntryId(undefined)
    }
  }

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <DataRow variant="header">
        <div data-testid="expenses-period">
          <Stat label="Period" value={dashboard()?.periodKey || '...'} variant="md" />
        </div>
        <ActionButton tone="ghost" busy={isClosingMonth()} onClick={() => setIsConfirmingCloseMonth(true)}>Close month</ActionButton>
      </DataRow>

      <Show when={actionError()}>
        {(message) => <Notice tone="error" message={message()} />}
      </Show>

      <ConfirmSheet
        open={isConfirmingCloseMonth()}
        title="Close month"
        description={`This rolls ${dashboard()?.periodKey ?? 'the current period'} into a new sheet. Make sure every obligation and last-minute expense is handled before you close it.`}
        confirmLabel="Yes, close month"
        confirmTone="danger"
        eyebrow="Period rollover"
        busy={isClosingMonth()}
        onCancel={() => setIsConfirmingCloseMonth(false)}
        onConfirm={handleCloseMonth}
      />

      <Suspense fallback={<DashboardSkeleton />}>
        <Show when={dashboard()} fallback={dashboard.error ? <Notice tone="error" message="Expenses dashboard failed to load." /> : <DashboardSkeleton />}> 
          {(data) => (
            <>
              <SummaryCard dashboard={data()} />
              <ObligationsCard
                obligations={data().pendingObligations ?? []}
                busyId={busyObligationId()}
                onComplete={handleCompleteObligation}
              />
              <HistoryCard history={data().history ?? []} deletingId={deletingEntryId()} onDelete={handleDeleteEntry} />
            </>
          )}
        </Show>
      </Suspense>

      <FloatingActionGroup>
        <DataCell flex><ActionLinkButton href="/expenses/settings" block tone="ghost">Settings</ActionLinkButton></DataCell>
        <DataCell flex><ActionLinkButton href="/expenses/new" block>Add Expense</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
