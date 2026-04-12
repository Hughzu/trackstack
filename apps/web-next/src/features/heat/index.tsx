import { createResource, createSignal, Show, Suspense } from 'solid-js'

import { ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { FormBackLink } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { authState } from '../../core/auth/state'
import type { Refill } from '../../core/api/types'
import { deleteRefill, readHeatDashboard } from './api/client'
import { HeatDashboardSkeleton } from './components/dashboard-skeleton'
import { HeatHistoryCard } from './components/heat-history-card'
import { HeatOverviewCard } from './components/heat-overview-card'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

export default function HeatPage() {
  const [dashboard, { refetch }] = createResource(readyKey, readHeatDashboard)
  const [deletingId, setDeletingId] = createSignal<string | undefined>()
  const [actionError, setActionError] = createSignal<string | undefined>()

  const handleDeleteRefill = async (refill: Refill) => {
    if (deletingId()) return

    setActionError(undefined)
    setDeletingId(refill.id)

    try {
      await deleteRefill(refill.id)
      await refetch()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to delete refill')
    } finally {
      setDeletingId(undefined)
    }
  }

  return (
    <ContentDeck layout="stacked" animate hasFloatingActions>
      <DataRow variant="header">
        <FormBackLink href="/">Back</FormBackLink>
      </DataRow>

      <Show when={actionError()}>
        {(message) => <Notice tone="error" message={message()} />}
      </Show>

      <Suspense fallback={<HeatDashboardSkeleton />}>
        <Show when={dashboard()} fallback={dashboard.error ? <Notice tone="error" message="Heat dashboard failed to load." /> : <HeatDashboardSkeleton />}>
          {(data) => (
            <>
              <HeatOverviewCard dashboard={data()} />
              <HeatHistoryCard history={data().history ?? []} deletingId={deletingId()} onDelete={handleDeleteRefill} />
            </>
          )}
        </Show>
      </Suspense>

      <FloatingActionGroup>
        <DataCell flex><ActionLinkButton href="/heat/new" block>Add refill</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
