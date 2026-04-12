import { createMemo, createResource, createSignal, Show, Suspense } from 'solid-js'

import { ActionLinkButton, FloatingActionGroup } from '../../components/ui/ActionButton'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataCell, DataRow } from '../../components/ui/DataRow'
import { FormBackLink } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import type { CalorieLog } from '../../core/api/types'
import { authState } from '../../core/auth/state'
import { createCalorieLog, deleteCalorieLog, readCaloriesDashboard } from './api/client'
import { CaloriesLogListCard } from './components/calories-log-list-card'
import { CaloriesOverviewCard } from './components/calories-overview-card'
import { CaloriesQuickAddCard } from './components/calories-quick-add-card'
import { CaloriesDashboardSkeleton } from './components/dashboard-skeleton'
import { createRecentMealChips } from './display'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

export default function CaloriesPage() {
  const [dashboard, { refetch }] = createResource(readyKey, readCaloriesDashboard)
  const [quickAddBusyId, setQuickAddBusyId] = createSignal<string | undefined>()
  const [deletingLogId, setDeletingLogId] = createSignal<string | undefined>()
  const [actionError, setActionError] = createSignal<string | undefined>()

  const recentMealChips = createMemo(() => {
    const data = dashboard()
    return data ? createRecentMealChips(data) : []
  })

  const refreshDashboard = async () => {
    const refreshed = await refetch()
    return refreshed ?? dashboard.latest
  }

  const handleQuickAdd = async (mealId: string) => {
    if (quickAddBusyId()) return

    const meal = dashboard()?.recentMeals.find((entry) => entry.id === mealId)
    if (!meal) return

    setActionError(undefined)
    setQuickAddBusyId(mealId)

    try {
      await createCalorieLog({
        title: meal.title ?? null,
        calories: meal.calories,
        proteinGrams: meal.proteinGrams,
        carbGrams: meal.carbGrams ?? null,
        fatGrams: meal.fatGrams ?? null,
      })
      await refreshDashboard()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to quick add meal')
    } finally {
      setQuickAddBusyId(undefined)
    }
  }

  const handleDeleteLog = async (log: CalorieLog) => {
    if (deletingLogId()) return

    setActionError(undefined)
    setDeletingLogId(log.id)

    try {
      await deleteCalorieLog(log.id)
      await refreshDashboard()
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Unable to delete meal')
    } finally {
      setDeletingLogId(undefined)
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

      <Suspense fallback={<CaloriesDashboardSkeleton />}>
        <Show when={dashboard()} fallback={dashboard.error ? <Notice tone="error" message="Calories dashboard failed to load." /> : <CaloriesDashboardSkeleton />}>
          {(data) => (
            <>
              <CaloriesOverviewCard dashboard={data()} />
              <CaloriesQuickAddCard items={recentMealChips()} submittingId={quickAddBusyId()} onSelect={(item) => void handleQuickAdd(item.id)} />
              <CaloriesLogListCard dashboardLogs={data().logs ?? []} deletingId={deletingLogId()} onDelete={handleDeleteLog} />
            </>
          )}
        </Show>
      </Suspense>

      <FloatingActionGroup>
        <DataCell flex><ActionLinkButton href="/calories/settings" block tone="ghost">Settings</ActionLinkButton></DataCell>
        <DataCell flex><ActionLinkButton href="/calories/new" block>Add Meal</ActionLinkButton></DataCell>
      </FloatingActionGroup>
    </ContentDeck>
  )
}
