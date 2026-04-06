import { For, Show } from 'solid-js'

type Segment = {
  percent: number
  color?: 'danger' | 'warning' | 'success' | 'accent' | 'main' | 'muted'
  label?: string
}

export type ProgressBarProps = {
  segments: Segment[]
}

const colorMap = {
  danger: 'bg-red-500/80',
  warning: 'bg-orange-500/80',
  success: 'bg-emerald-500/80',
  accent: 'bg-accent',
  main: 'bg-text-main',
  muted: 'bg-text-muted',
}

export function ProgressBar(props: ProgressBarProps) {
  const hasLabels = () => props.segments.some((s) => s.label)

  return (
    <>
      <div class="mb-1 flex h-1.5 w-full overflow-hidden rounded-full bg-border/50">
        <For each={props.segments}>
          {(segment) => (
            <div
              class={`h-full transition-all duration-500 ${colorMap[segment.color || 'main']}`}
              style={{ width: `${segment.percent}%` }}
            />
          )}
        </For>
      </div>
      
      <Show when={hasLabels()}>
        <div class="mt-2 flex justify-between text-[0.65rem] font-mono text-text-muted">
          <For each={props.segments}>
            {(segment) => (
              segment.label ? <span>{segment.label} {segment.percent}%</span> : <span />
            )}
          </For>
        </div>
      </Show>
    </>
  )
}
