import { createEffect, createResource, createSignal, Show, Suspense } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { authState } from '../../core/auth/state'
import { readCalorieTarget, updateCalorieTarget } from './api/client'
import { CaloriesSettingsCard } from './components/calories-settings-card'
import { CaloriesSettingsSkeleton } from './components/dashboard-skeleton'
import { createTargetBadge, normalizeCalorieTargetInput } from './display'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

const toNumberOrNull = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return null

  const parsed = Number.parseFloat(trimmed)
  return Number.isFinite(parsed) ? parsed : Number.NaN
}

export default function CaloriesSettingsPage() {
  const navigate = useNavigate()
  const [target, { refetch }] = createResource(readyKey, readCalorieTarget)
  const [targetCalories, setTargetCalories] = createSignal('')
  const [targetProtein, setTargetProtein] = createSignal('')
  const [targetCarbs, setTargetCarbs] = createSignal('')
  const [targetFat, setTargetFat] = createSignal('')
  const [isSaving, setIsSaving] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  createEffect(() => {
    const data = target()
    if (!data) return

    setTargetCalories(String(data.targetCalories))
    setTargetProtein(String(data.targetProteinGrams))
    setTargetCarbs(data.targetCarbGrams == null ? '' : String(data.targetCarbGrams))
    setTargetFat(data.targetFatGrams == null ? '' : String(data.targetFatGrams))
  })

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSaving()) return

    const nextCalories = Number.parseFloat(normalizeCalorieTargetInput(targetCalories()))
    const nextProtein = Number.parseFloat(normalizeCalorieTargetInput(targetProtein()))
    const nextCarbs = toNumberOrNull(targetCarbs())
    const nextFat = toNumberOrNull(targetFat())

    if (!Number.isFinite(nextCalories) || nextCalories <= 0) {
      setErrorMessage('Daily target needs to be above zero.')
      return
    }

    if (!Number.isFinite(nextProtein) || nextProtein <= 0) {
      setErrorMessage('Protein target needs to be above zero.')
      return
    }

    if ((nextCarbs != null && (!Number.isFinite(nextCarbs) || nextCarbs < 0)) || (nextFat != null && (!Number.isFinite(nextFat) || nextFat < 0))) {
      setErrorMessage('Optional carb and fat targets need valid positive numbers.')
      return
    }

    setIsSaving(true)
    setErrorMessage(undefined)

    try {
      await updateCalorieTarget({
        targetCalories: nextCalories,
        targetProteinGrams: nextProtein,
        targetCarbGrams: nextCarbs,
        targetFatGrams: nextFat,
      })

      await refetch()
      void navigate('/calories', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save calorie settings')
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <ContentDeck layout="stacked" animate>
      <DataRow variant="header">
        <FormBackLink href="/calories">Back</FormBackLink>
      </DataRow>

      <Show when={errorMessage()}>
        {(message) => <Notice tone="error" message={message()} />}
      </Show>

      <Suspense fallback={<CaloriesSettingsSkeleton />}>
        <Show when={target()} fallback={target.error ? <Notice tone="error" message="Calories settings failed to load." /> : <CaloriesSettingsSkeleton />}>
          {(data) => (
            <FormStack onSubmit={handleSubmit}>
              <FormSection>
                <AmountHeroField
                  inputId="calories-target"
                  name="targetCalories"
                  label="Daily target"
                  unit="KCAL"
                  value={targetCalories()}
                  placeholder="0"
                  badge={createTargetBadge(data())}
                  required
                  min="0"
                  step="1"
                  onInput={(event) => setTargetCalories(event.currentTarget.value)}
                />

                <CaloriesSettingsCard
                  fields={[
                    {
                      id: 'calories-target-protein',
                      name: 'targetProteinGrams',
                      label: 'Protein',
                      unit: 'g',
                      value: targetProtein(),
                      tone: 'success',
                      onInput: (event) => setTargetProtein(event.currentTarget.value),
                    },
                    {
                      id: 'calories-target-carbs',
                      name: 'targetCarbGrams',
                      label: 'Carbs',
                      unit: 'g',
                      value: targetCarbs(),
                      tone: 'warning',
                      onInput: (event) => setTargetCarbs(event.currentTarget.value),
                    },
                    {
                      id: 'calories-target-fat',
                      name: 'targetFatGrams',
                      label: 'Fat',
                      unit: 'g',
                      value: targetFat(),
                      tone: 'danger',
                      onInput: (event) => setTargetFat(event.currentTarget.value),
                    },
                  ]}
                />
              </FormSection>

              <FormActions>
                <ActionButton type="submit" busy={isSaving()} disabled={isSaving()}>
                  {isSaving() ? 'Saving...' : 'Save settings'}
                </ActionButton>
              </FormActions>
            </FormStack>
          )}
        </Show>
      </Suspense>
    </ContentDeck>
  )
}
