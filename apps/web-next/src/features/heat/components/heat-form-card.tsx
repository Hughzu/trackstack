import { FormFieldRow, FormSection } from '../../../components/ui/Form'
import { MetricField } from '../../../components/ui/MetricField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { TextField } from '../../../components/ui/TextField'

export function HeatFormCard(props: {
  weightKg: string
  date: string
  temperature: string
  compact?: boolean
  onWeightInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
  onDateInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
  onTemperatureInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
}) {
  const density = props.compact ? 'compact' : 'default'

  return (
    <Panel title="Refill details" description={<Pill tone="neutral">15 kg per bag baseline</Pill>} density={density}>
      <FormSection density={density}>
        <FormFieldRow density={density}>
          <MetricField
            id="heat-weight"
            name="weightKg"
            label="Total weight"
            unit="kg"
            value={props.weightKg}
            density={density}
            min="0"
            step="0.1"
            inputMode="decimal"
            onInput={props.onWeightInput}
          />

          <TextField
            id="heat-date"
            name="date"
            label="Date"
            type="date"
            value={props.date}
            density={density}
            onInput={props.onDateInput}
          />
        </FormFieldRow>

        <MetricField
          id="heat-temperature"
          name="temperature"
          label="Average temperature"
          unit="C"
          value={props.temperature}
          density={density}
          tone="warning"
          min="-30"
          step="0.1"
          inputMode="decimal"
          onInput={props.onTemperatureInput}
        />
      </FormSection>
    </Panel>
  )
}
