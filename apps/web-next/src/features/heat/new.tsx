import { createEffect, createSignal, Show } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormFieldRow, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { createHeatMockRefill } from './mock-state'
import { TextField } from '../../components/ui/TextField'

const today = new Date().toISOString().slice(0, 10)

const toWholeNumber = (value: string) => {
  const parsed = Number.parseInt(value, 10)
  return Number.isFinite(parsed) ? parsed : 0
}

const toDecimalNumber = (value: string) => {
  const parsed = Number.parseFloat(value)
  return Number.isFinite(parsed) ? parsed : 0
}

export default function NewHeatRefillPage() {
  const navigate = useNavigate()
  const [bags, setBags] = createSignal('2')
  const [weightKg, setWeightKg] = createSignal('30')
  const [temperature, setTemperature] = createSignal('')
  const [date, setDate] = createSignal(today)
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  createEffect(() => {
    const nextBags = toWholeNumber(bags())
    setWeightKg(String(nextBags > 0 ? nextBags * 15 : 0))
  })

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSubmitting()) return

    const nextBags = toWholeNumber(bags())
    const nextWeight = toDecimalNumber(weightKg())
    const nextTemperature = temperature().trim() ? toDecimalNumber(temperature()) : null

    if (nextBags <= 0) {
      setErrorMessage('Bags must be above zero.')
      return
    }

    if (nextWeight <= 0) {
      setErrorMessage('Weight must be above zero.')
      return
    }

    setIsSubmitting(true)
    setErrorMessage(undefined)

    try {
      createHeatMockRefill({
        date: date(),
        bags: nextBags,
        weightKg: nextWeight,
        temperature: nextTemperature,
      })

      void navigate('/heat', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save refill')
      setIsSubmitting(false)
    }
  }

  return (
    <ContentDeck layout="stacked" animate>
      <DataRow variant="header">
        <FormBackLink href="/heat">Back</FormBackLink>
      </DataRow>

      <FormStack onSubmit={handleSubmit}>
        <FormSection>
          <AmountHeroField
            inputId="heat-bags"
            label="Bags added"
            unit="Bags"
            value={bags()}
            placeholder="0"
            onInput={(event) => setBags(event.currentTarget.value)}
          />

          <Panel title="Refill details" description={<Pill tone="neutral">15 kg per bag baseline</Pill>}>
            <div class="flex flex-col gap-5">
              <FormFieldRow>
                <TextField
                  id="heat-weight"
                  name="weightKg"
                  label="Total weight"
                  type="number"
                  value={weightKg()}
                  onInput={(event) => setWeightKg(event.currentTarget.value)}
                />

                <TextField
                  id="heat-date"
                  name="date"
                  label="Date"
                  type="date"
                  value={date()}
                  onInput={(event) => setDate(event.currentTarget.value)}
                />
              </FormFieldRow>

              <div class="border-t border-border/40 pt-4">
                <TextField
                  id="heat-temperature"
                  name="temperature"
                  label="Average temperature"
                  type="number"
                  value={temperature()}
                  onInput={(event) => setTemperature(event.currentTarget.value)}
                />
              </div>
            </div>
          </Panel>

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
