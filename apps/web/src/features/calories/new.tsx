import { createSignal, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { createCalorieLog } from './api/client'
import { CaloriesFormCard } from './components/calories-form-card'

const toNumber = (value: string) => Number.parseFloat(value)

export default function NewCaloriesPage() {
  const navigate = useNavigate()
  const [calories, setCalories] = createSignal('')
  const [title, setTitle] = createSignal('')
  const [protein, setProtein] = createSignal('')
  const [carbs, setCarbs] = createSignal('')
  const [fat, setFat] = createSignal('')
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSubmitting()) return

    const nextCalories = toNumber(calories())
    const nextProtein = toNumber(protein())
    const nextCarbs = carbs().trim() ? toNumber(carbs()) : null
    const nextFat = fat().trim() ? toNumber(fat()) : null

    if (!Number.isFinite(nextCalories) || nextCalories <= 0) {
      setErrorMessage('Calories need to be above zero.')
      return
    }

    if (!Number.isFinite(nextProtein) || nextProtein < 0) {
      setErrorMessage('Protein needs a valid number.')
      return
    }

    if ((nextCarbs != null && (!Number.isFinite(nextCarbs) || nextCarbs < 0)) || (nextFat != null && (!Number.isFinite(nextFat) || nextFat < 0))) {
      setErrorMessage('Carbs and fat need valid positive numbers when filled in.')
      return
    }

    setIsSubmitting(true)
    setErrorMessage(undefined)

    try {
      await createCalorieLog({
        title: title().trim() || null,
        calories: nextCalories,
        proteinGrams: nextProtein,
        carbGrams: nextCarbs,
        fatGrams: nextFat,
      })

      void navigate('/calories', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save meal')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <ContentDeck layout="stacked" animate density="compact">
      <DataRow variant="header">
        <FormBackLink href="/calories">Back</FormBackLink>
      </DataRow>

      <FormStack onSubmit={handleSubmit} density="compact">
        <FormSection density="compact">
          <AmountHeroField
            inputId="calorie-amount"
            name="calories"
            label="Calories"
            unit="KCAL"
            value={calories()}
            placeholder="0"
            badge="Meal entry"
            density="compact"
            required
            min="0"
            step="1"
            onInput={(event) => setCalories(event.currentTarget.value)}
          />

          <CaloriesFormCard
            title="Oats + whey"
            titleValue={title()}
            compact
            onTitleInput={(event) => setTitle(event.currentTarget.value)}
            macroFields={[
              {
                id: 'calorie-protein',
                name: 'proteinGrams',
                label: 'Protein',
                value: protein(),
                tone: 'success',
                onInput: (event) => setProtein(event.currentTarget.value),
              },
              {
                id: 'calorie-carbs',
                name: 'carbGrams',
                label: 'Carbs',
                value: carbs(),
                tone: 'warning',
                onInput: (event) => setCarbs(event.currentTarget.value),
              },
              {
                id: 'calorie-fat',
                name: 'fatGrams',
                label: 'Fat',
                value: fat(),
                tone: 'danger',
                onInput: (event) => setFat(event.currentTarget.value),
              },
            ]}
          />

          <Show when={errorMessage()}>
            {(message) => <Notice tone="error" message={message()} />}
          </Show>
        </FormSection>

        <FormActions>
          <ActionButton type="submit" busy={isSubmitting()} disabled={isSubmitting()}>
            {isSubmitting() ? 'Saving...' : 'Save meal'}
          </ActionButton>
        </FormActions>
      </FormStack>
    </ContentDeck>
  )
}
