import { createSignal } from 'solid-js'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormFieldRow, FormSection, FormStack } from '../../components/ui/Form'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { TextField } from '../../components/ui/TextField'

export default function NewCaloriesPage() {
  const [calories, setCalories] = createSignal('520')
  const [title, setTitle] = createSignal('Oats + whey')
  const [protein, setProtein] = createSignal('38')
  const [carbs, setCarbs] = createSignal('56')
  const [fat, setFat] = createSignal('12')

  return (
    <ContentDeck layout="stacked" animate>
      <DataRow variant="header">
        <FormBackLink href="/calories">Back</FormBackLink>
      </DataRow>

      <FormStack>
        <FormSection>
          <AmountHeroField
            inputId="calorie-amount"
            label="Calories"
            unit="KCAL"
            value={calories()}
            placeholder="0"
            badge="Meal entry"
            onInput={(event) => setCalories(event.currentTarget.value)}
          />

          <Panel title="Meal log">
            <div class="flex flex-col gap-5">
              <TextField
                id="calorie-title"
                name="title"
                label="Title"
                value={title()}
                onInput={(event) => setTitle(event.currentTarget.value)}
              />

              <div class="flex flex-col gap-3 border-t border-border/40 pt-4">
                <div class="flex items-center justify-between gap-3">
                  <div class="text-sm font-semibold text-text-main">Macros</div>
                  <Pill tone="success">Protein first</Pill>
                </div>

                <FormFieldRow>
                  <MacroField
                    id="calorie-protein"
                    name="proteinGrams"
                    label="Protein"
                    value={protein()}
                    unit="g"
                    accentClass="bg-emerald-400/70"
                    onInput={(event) => setProtein(event.currentTarget.value)}
                  />

                  <MacroField
                    id="calorie-carbs"
                    name="carbGrams"
                    label="Carbs"
                    value={carbs()}
                    unit="g"
                    accentClass="bg-orange-400/70"
                    onInput={(event) => setCarbs(event.currentTarget.value)}
                  />
                </FormFieldRow>

                <MacroField
                  id="calorie-fat"
                  name="fatGrams"
                  label="Fat"
                  value={fat()}
                  unit="g"
                  accentClass="bg-yellow-400/70"
                  onInput={(event) => setFat(event.currentTarget.value)}
                />
              </div>

              <div class="border-t border-border/50 pt-4">
                <FormActions>
                  <ActionButton type="submit">Save meal</ActionButton>
                </FormActions>
              </div>
            </div>
          </Panel>
        </FormSection>
      </FormStack>
    </ContentDeck>
  )
}

function MacroField(props: {
  id: string
  name: string
  label: string
  value: string
  unit: string
  accentClass: string
  onInput: (event: InputEvent & { currentTarget: HTMLInputElement, target: HTMLInputElement }) => void
}) {
  return (
    <label class="flex flex-col gap-2 rounded-2xl border border-border/60 bg-panel/60 px-4 py-3 text-sm text-text-main" for={props.id}>
      <div class="flex items-center justify-between gap-3">
        <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">{props.label}</span>
        <span class={`h-2 w-2 rounded-full ${props.accentClass}`} />
      </div>
      <div class="flex items-end justify-between gap-3 border-b border-border/50 pb-2">
        <input
          id={props.id}
          name={props.name}
          type="number"
          value={props.value}
          onInput={props.onInput}
          class="w-full bg-transparent text-xl font-semibold text-text-main outline-none"
        />
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-text-muted">{props.unit}</span>
      </div>
    </label>
  )
}
