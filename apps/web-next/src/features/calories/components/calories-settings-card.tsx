import { For, type JSX } from 'solid-js'

import { FormSection } from '../../../components/ui/Form'
import { MetricField } from '../../../components/ui/MetricField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import type { MacroTone } from '../display'

type MacroTargetField = {
  id: string
  name: string
  label: string
  unit: string
  value: string
  tone: MacroTone
  onInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
}

export function CaloriesSettingsCard(props: { fields: MacroTargetField[] }) {
  return (
    <Panel title="Macros" description={<Pill tone="success">Protein required</Pill>}>
      <p>Protein matters most. Carbs and fat stay optional instead of turning this page into nutrition bureaucracy.</p>

      <FormSection>
        <For each={props.fields}>
          {(field) => (
            <MetricField
              id={field.id}
              name={field.name}
              label={field.label}
              unit={field.unit}
              value={field.value}
              tone={field.tone}
              variant="inline"
              inputMode="decimal"
              min="0"
              step="1"
              onInput={field.onInput}
            />
          )}
        </For>
      </FormSection>
    </Panel>
  )
}
