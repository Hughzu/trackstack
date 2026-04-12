import { For } from 'solid-js'

import { MetricField } from '../../../components/ui/MetricField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
type RatioItem = {
  id: string
  title: string
  caption: string
  tone: 'danger' | 'warning' | 'success'
  value: string
  amount: string
  onInput: (value: string) => void
}

export function SettingsRatioCard(props: { items: RatioItem[] }) {
  return (
    <Panel title="Ratio architecture" description={<Pill tone="neutral">Adjust the split</Pill>}>
      <div class="divide-y divide-border/50">
        <For each={props.items}>
          {(item) => <RatioEditorRow {...item} />}
        </For>
      </div>
    </Panel>
  )
}

function RatioEditorRow(props: RatioItem) {
  const accentClass = props.tone === 'danger' ? 'accent-red-500' : props.tone === 'warning' ? 'accent-orange-500' : 'accent-emerald-500'
  const badgeTone = props.tone === 'danger' ? 'danger' : props.tone === 'warning' ? 'warning' : 'success'

  return (
    <div class="flex flex-col gap-4 py-4 first:pt-0 last:pb-0">
      <div class="flex items-start justify-between gap-3">
        <div>
          <div class="text-base font-bold tracking-tight text-text-main">{props.title}</div>
          <p class="mt-1 text-xs leading-5 text-text-muted">{props.caption}</p>
        </div>
        <Pill tone={badgeTone}>{Math.round(Number(props.value) || 0)}%</Pill>
      </div>

      <input
        type="range"
        min="0"
        max="100"
        step="1"
        value={props.value}
        onInput={(event) => props.onInput(event.currentTarget.value)}
        class={`h-2 w-full cursor-pointer appearance-none rounded-full bg-border/60 ${accentClass}`}
      />

      <div class="flex items-end justify-between gap-3">
        <div>
          <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Projected amount</div>
          <div class="mt-1 text-xl font-bold tracking-tight text-text-main">{props.amount}</div>
        </div>

        <div class="w-24">
          <MetricField
            id={`${props.id}-manual`}
            name={`${props.id}-manual`}
            label="Manual"
            unit="%"
            value={props.value}
            tone={props.tone}
            variant="inline"
            min="0"
            max="100"
            step="1"
            inputMode="decimal"
            onInput={(event) => props.onInput(event.currentTarget.value)}
          />
        </div>
      </div>
    </div>
  )
}
