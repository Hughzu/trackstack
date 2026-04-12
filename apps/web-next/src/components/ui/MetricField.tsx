import type { JSX } from 'solid-js'

import { SkeletonBlock, SkeletonText } from './Skeleton'

type MetricFieldTone = 'accent' | 'success' | 'warning' | 'danger' | 'neutral'

type MetricFieldProps = {
  id: string
  name: string
  label: string
  unit: string
  value: string
  tone?: MetricFieldTone
  variant?: 'card' | 'inline'
  density?: 'default' | 'compact'
  placeholder?: string
  required?: boolean
  min?: number | string
  max?: number | string
  step?: number | string
  inputMode?: 'decimal' | 'numeric' | 'text'
  onInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
}

const toneDotClass = {
  accent: 'bg-accent',
  success: 'bg-emerald-400',
  warning: 'bg-orange-400',
  danger: 'bg-red-400',
  neutral: 'bg-text-muted',
}

export function MetricField(props: MetricFieldProps) {
  const tone = () => toneDotClass[props.tone ?? 'accent']
  const compact = () => props.density === 'compact'

  if (props.variant === 'inline') {
    return (
      <label class={compact()
        ? 'flex items-center justify-between gap-3 rounded-2xl border border-border/40 bg-panel/35 px-4 py-2.5'
        : 'flex items-center justify-between gap-4 rounded-2xl border border-border/40 bg-panel/35 px-4 py-3'} for={props.id}>
        <div class="flex min-w-0 items-center gap-3">
          <span class={`h-2 w-2 rounded-full ${tone()}`} />
          <span class="text-sm font-semibold text-text-main">{props.label}</span>
        </div>

        <div class={compact() ? 'flex w-24 items-center justify-end gap-2 border-b border-border/50 pb-1' : 'flex w-28 items-center justify-end gap-2 border-b border-border/50 pb-1.5'}>
          <input
            id={props.id}
            name={props.name}
            type="number"
            value={props.value}
            placeholder={props.placeholder}
            required={props.required}
            min={props.min}
            max={props.max}
            step={props.step ?? '1'}
            inputmode={props.inputMode ?? 'decimal'}
            onInput={props.onInput}
            class={compact()
              ? 'w-full bg-transparent text-right text-base font-semibold text-text-main outline-none placeholder:text-text-muted/40'
              : 'w-full bg-transparent text-right text-lg font-semibold text-text-main outline-none placeholder:text-text-muted/40'}
          />
          <span class="text-xs font-semibold uppercase tracking-[0.2em] text-text-muted">{props.unit}</span>
        </div>
      </label>
    )
  }

  return (
    <label class={compact()
      ? 'flex flex-col gap-1.5 rounded-2xl border border-border/60 bg-panel/60 px-4 py-2.5 text-sm text-text-main'
      : 'flex flex-col gap-2 rounded-2xl border border-border/60 bg-panel/60 px-4 py-3 text-sm text-text-main'} for={props.id}>
      <div class="flex items-center justify-between gap-3">
        <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">{props.label}</span>
        <span class={`h-2 w-2 rounded-full ${tone()}`} />
      </div>
      <div class={compact() ? 'flex items-end justify-between gap-3 border-b border-border/50 pb-1.5' : 'flex items-end justify-between gap-3 border-b border-border/50 pb-2'}>
        <input
          id={props.id}
          name={props.name}
          type="number"
          value={props.value}
          placeholder={props.placeholder}
          required={props.required}
          min={props.min}
          max={props.max}
            step={props.step ?? '1'}
            inputmode={props.inputMode ?? 'decimal'}
            onInput={props.onInput}
            class={compact()
              ? 'w-full bg-transparent text-lg font-semibold text-text-main outline-none placeholder:text-text-muted/40'
              : 'w-full bg-transparent text-xl font-semibold text-text-main outline-none placeholder:text-text-muted/40'}
          />
          <span class="text-xs font-semibold uppercase tracking-[0.2em] text-text-muted">{props.unit}</span>
        </div>
    </label>
  )
}

export function SkeletonMetricField(props: { variant?: 'card' | 'inline' }) {
  if (props.variant === 'inline') {
    return (
      <div class="flex items-center justify-between gap-4 rounded-2xl border border-border/40 bg-panel/35 px-4 py-3">
        <div class="flex min-w-0 items-center gap-3">
          <SkeletonBlock class="h-2 w-2 rounded-full" />
          <SkeletonText width="w-16" />
        </div>

        <div class="flex w-28 items-center justify-end gap-2 border-b border-border/50 pb-1.5">
          <SkeletonBlock class="h-6 w-14 border-border/20 bg-surface/70" />
          <SkeletonText width="w-5" />
        </div>
      </div>
    )
  }

  return (
    <div class="flex flex-col gap-2 rounded-2xl border border-border/60 bg-panel/60 px-4 py-3">
      <div class="flex items-center justify-between gap-3">
        <SkeletonText width="w-16" />
        <SkeletonBlock class="h-2 w-2 rounded-full" />
      </div>
      <div class="flex items-end justify-between gap-3 border-b border-border/50 pb-2">
        <SkeletonBlock class="h-8 w-20 border-border/20 bg-surface/70" />
        <SkeletonText width="w-5" />
      </div>
    </div>
  )
}
