import { createSignal } from 'solid-js'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormSection, FormStack } from '../../components/ui/Form'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'

export default function CaloriesSettingsPage() {
  const [targetCalories, setTargetCalories] = createSignal('2400')
  const [targetProtein, setTargetProtein] = createSignal('180')

  return (
    <ContentDeck layout="stacked" animate>
      <DataRow variant="header">
        <FormBackLink href="/calories">Back</FormBackLink>
      </DataRow>

      <FormStack>
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

          <Panel title="Macros">
            <div class="flex flex-col gap-4">
              <div class="flex items-center justify-between gap-3">
                <div class="text-sm text-text-muted">Protein is the only macro target worth pinning here.</div>
                <Pill tone="success">Protein required</Pill>
              </div>

              <div class="divide-y divide-border/30 border-t border-border/40">
                <MacroSettingRow
                  label="Protein"
                  tone="success"
                  value={targetProtein()}
                  onInput={(event) => setTargetProtein(event.currentTarget.value)}
                />
              </div>

              <div class="border-t border-border/50 pt-4">
                <FormActions>
                  <ActionButton type="submit">Save settings</ActionButton>
                </FormActions>
              </div>
            </div>
          </Panel>
        </FormSection>
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
