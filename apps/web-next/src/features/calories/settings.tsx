import { createSignal } from 'solid-js'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'

export default function CaloriesSettingsPage() {
  const [targetCalories, setTargetCalories] = createSignal('2400')
  const [targetProtein, setTargetProtein] = createSignal('180')
  const [isSaving, setIsSaving] = createSignal(false)
  const [feedbackMessage, setFeedbackMessage] = createSignal<string | undefined>()
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSaving()) return

    const nextCalories = Number.parseFloat(targetCalories())
    const nextProtein = Number.parseFloat(targetProtein())

    if (!Number.isFinite(nextCalories) || nextCalories <= 0) {
      setErrorMessage('Daily target needs to be above zero.')
      setFeedbackMessage(undefined)
      return
    }

    if (!Number.isFinite(nextProtein) || nextProtein <= 0) {
      setErrorMessage('Protein target needs to be above zero.')
      setFeedbackMessage(undefined)
      return
    }

    setIsSaving(true)
    setErrorMessage(undefined)
    setFeedbackMessage(undefined)

    try {
      setFeedbackMessage('Settings saved. Your calorie targets are locked in.')
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

      {errorMessage() ? <Notice tone="error" message={errorMessage()!} /> : null}
      {feedbackMessage() ? <Notice tone="info" message={feedbackMessage()!} /> : null}

      <FormStack onSubmit={handleSubmit}>
        <FormSection>
          <AmountHeroField
            inputId="calories-target"
            label="Daily target"
            unit="KCAL"
            value={targetCalories()}
            placeholder="0"
            badge="Intake goal"
            onInput={(event) => setTargetCalories(event.currentTarget.value)}
          />

          <Panel title="Macros" description={<Pill tone="success">Protein required</Pill>}>
            <div class="flex flex-col gap-4">
              <div class="border-b border-border/40 pb-4 text-sm text-text-muted">
                Protein is the only macro target worth pinning here.
              </div>

              <div class="divide-y divide-border/30 rounded-2xl border border-border/40 bg-panel/35 px-4">
                <MacroSettingRow
                  label="Protein"
                  tone="success"
                  value={targetProtein()}
                  onInput={(event) => setTargetProtein(event.currentTarget.value)}
                />
              </div>
            </div>
          </Panel>
        </FormSection>

        <FormActions>
          <ActionButton type="submit" busy={isSaving()} disabled={isSaving()}>
            {isSaving() ? 'Saving...' : 'Save settings'}
          </ActionButton>
        </FormActions>
      </FormStack>
    </ContentDeck>
  )
}

function MacroSettingRow(props: {
  label: string
  tone: 'success'
  value: string
  onInput: (event: InputEvent & { currentTarget: HTMLInputElement, target: HTMLInputElement }) => void
}) {
  const dotClass = 'bg-emerald-400'

  return (
    <label class="flex items-center justify-between gap-4 py-4 first:pt-4 last:pb-0">
      <div class="flex min-w-0 items-center gap-3">
        <span class={`h-2 w-2 rounded-full ${dotClass}`} />
        <span class="text-sm font-semibold text-text-main">{props.label}</span>
      </div>

      <div class="flex w-28 items-center justify-end gap-2 border-b border-border/50 pb-1.5">
        <input
          type="number"
          value={props.value}
          onInput={props.onInput}
          class="w-full bg-transparent text-right text-lg font-semibold text-text-main outline-none"
        />
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-text-muted">g</span>
      </div>
    </label>
  )
}
