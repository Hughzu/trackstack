import { createEffect, createSignal, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { createRefill } from './api/client'
import { HeatFormCard } from './components/heat-form-card'

const today = new Date().toISOString().slice(0, 10)

const toWholeNumber = (value: string) => {
  const parsed = Number.parseInt(value, 10)
  return Number.isFinite(parsed) ? parsed : Number.NaN
}

const toDecimalNumber = (value: string) => {
  const parsed = Number.parseFloat(value)
  return Number.isFinite(parsed) ? parsed : Number.NaN
}

export default function NewHeatRefillPage() {
  const navigate = useNavigate()
  const [bags, setBags] = createSignal('')
  const [weightKg, setWeightKg] = createSignal('')
  const [temperature, setTemperature] = createSignal('')
  const [date, setDate] = createSignal(today)
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  createEffect(() => {
    const nextBags = toWholeNumber(bags())
    if (!Number.isFinite(nextBags) || nextBags <= 0) {
      setWeightKg('')
      return
    }

    setWeightKg(String(nextBags * 15))
  })

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSubmitting()) return

    const nextBags = toWholeNumber(bags())
    const nextWeight = toDecimalNumber(weightKg())
    const nextTemperature = temperature().trim() ? toDecimalNumber(temperature()) : null

    if (!Number.isFinite(nextBags) || nextBags <= 0) {
      setErrorMessage('Bags must be above zero.')
      return
    }

    if (!Number.isFinite(nextWeight) || nextWeight <= 0) {
      setErrorMessage('Weight must be above zero.')
      return
    }

    if (nextTemperature != null && (!Number.isFinite(nextTemperature) || nextTemperature < -50 || nextTemperature > 50)) {
      setErrorMessage('Temperature needs a real-world value.')
      return
    }

    setIsSubmitting(true)
    setErrorMessage(undefined)

    try {
      await createRefill({
        date: date(),
        bags: nextBags,
        weightKg: nextWeight,
        temperature: nextTemperature,
      })

      void navigate('/heat', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save refill')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <ContentDeck layout="stacked" animate density="compact">
      <DataRow variant="header">
        <FormBackLink href="/heat">Back</FormBackLink>
      </DataRow>

      <FormStack onSubmit={handleSubmit} density="compact">
        <FormSection density="compact">
          <AmountHeroField
            inputId="heat-bags"
            name="bags"
            label="Bags added"
            unit="BAGS"
            value={bags()}
            placeholder="0"
            density="compact"
            required
            min="0"
            step="1"
            onInput={(event) => setBags(event.currentTarget.value)}
          />

          <HeatFormCard
            compact
            weightKg={weightKg()}
            date={date()}
            temperature={temperature()}
            onWeightInput={(event) => setWeightKg(event.currentTarget.value)}
            onDateInput={(event) => setDate(event.currentTarget.value)}
            onTemperatureInput={(event) => setTemperature(event.currentTarget.value)}
          />

          <Show when={errorMessage()}>
            {(message) => <Notice tone="error" message={message()} />}
          </Show>
        </FormSection>

        <FormActions>
          <ActionButton type="submit" disabled={isSubmitting()} busy={isSubmitting()}>
            {isSubmitting() ? 'Saving...' : 'Log refill'}
          </ActionButton>
        </FormActions>
      </FormStack>
    </ContentDeck>
  )
}
