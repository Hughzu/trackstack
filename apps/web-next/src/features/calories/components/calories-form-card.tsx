import { For, type JSX } from 'solid-js'

import { FormFieldRow, FormSection } from '../../../components/ui/Form'
import { MetricField } from '../../../components/ui/MetricField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { TextField } from '../../../components/ui/TextField'
import type { MacroTone } from '../display'

type MacroInputField = {
  id: string
  name: string
  label: string
  value: string
  tone: MacroTone
  onInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
}

export function CaloriesFormCard(props: {
  title: string
  titleValue: string
  compact?: boolean
  onTitleInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
  macroFields: MacroInputField[]
}) {
  const leadMacroFields = () => props.macroFields.slice(0, 2)
  const finalMacroField = () => props.macroFields[2]

  return (
    <Panel title="Meal log" description={<Pill tone="success">Protein first</Pill>} density={props.compact ? 'compact' : 'default'}>
      <FormSection density={props.compact ? 'compact' : 'default'}>
        <TextField
          id="calorie-title"
          name="title"
          label="Title"
          value={props.titleValue}
          placeholder={props.title}
          density={props.compact ? 'compact' : 'default'}
          onInput={props.onTitleInput}
        />

        <p class="text-sm font-semibold text-text-main">Macros</p>
        <FormFieldRow density={props.compact ? 'compact' : 'default'}>
          <For each={leadMacroFields()}>
            {(field) => (
              <MetricField
                id={field.id}
                name={field.name}
                label={field.label}
                unit="g"
                value={field.value}
                tone={field.tone}
                density={props.compact ? 'compact' : 'default'}
                min="0"
                step="1"
                inputMode="decimal"
                onInput={field.onInput}
              />
            )}
          </For>
        </FormFieldRow>

        <MetricField
          id={finalMacroField().id}
          name={finalMacroField().name}
          label={finalMacroField().label}
          unit="g"
          value={finalMacroField().value}
          tone={finalMacroField().tone}
          density={props.compact ? 'compact' : 'default'}
          min="0"
          step="1"
          inputMode="decimal"
          onInput={finalMacroField().onInput}
        />
      </FormSection>
    </Panel>
  )
}
