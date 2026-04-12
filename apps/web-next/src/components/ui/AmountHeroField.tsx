import type { JSX } from 'solid-js'

import { SkeletonBlock, SkeletonText } from './Skeleton'

type AmountHeroFieldProps = {
  label: string
  unit: string
  value: string
  name?: string
  badge?: string
  badgeTone?: 'default' | 'danger' | 'warning' | 'success'
  inputId?: string
  placeholder?: string
  type?: 'number' | 'text'
  required?: boolean
  min?: number | string
  max?: number | string
  step?: number | string
  inputMode?: 'decimal' | 'numeric' | 'text'
  density?: 'default' | 'compact'
  onInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
}

const badgeToneClass = {
  default: 'text-accent border-accent/40 bg-accent/10',
  danger: 'text-danger border-danger/40 bg-danger/10',
  warning: 'text-warning border-warning/40 bg-warning/10',
  success: 'text-success border-success/40 bg-success/10',
}

export function AmountHeroField(props: AmountHeroFieldProps) {
  const compact = () => props.density === 'compact'

  return (
    <div class={compact()
      ? 'flex flex-col items-center gap-1.5 rounded-[1.5rem] border border-border/50 bg-panel/70 px-4 py-4 shadow-lg shadow-black/10'
      : 'flex flex-col items-center gap-2 rounded-[1.5rem] border border-border/50 bg-panel/70 px-4 py-5 shadow-lg shadow-black/10'}>
      <div class="text-[0.65rem] font-bold uppercase tracking-[0.28em] text-text-muted">{props.label}</div>
      <div class={compact() ? 'flex w-full max-w-[15rem] flex-col items-center border-b-2 border-border pb-1' : 'flex w-full max-w-[16rem] flex-col items-center border-b-2 border-border pb-1.5'}>
        <div class="mb-0.5 text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">{props.unit}</div>
        <label class="flex w-full justify-center" for={props.inputId}>
          <input
            id={props.inputId}
            name={props.name}
            type={props.type ?? 'number'}
            inputmode={props.inputMode ?? 'decimal'}
            required={props.required}
            min={props.min}
            max={props.max}
            step={props.step ?? '0.01'}
            value={props.value}
            placeholder={props.placeholder}
            onInput={props.onInput}
            class={compact()
              ? 'w-full bg-transparent text-center text-[2rem] font-bold tracking-tight text-text-main outline-none placeholder:text-text-muted/30 sm:text-[2.65rem]'
              : 'w-full bg-transparent text-center text-[2.35rem] font-bold tracking-tight text-text-main outline-none placeholder:text-text-muted/30 sm:text-5xl'}
          />
        </label>
      </div>
      {props.badge ? <div class={`inline-flex rounded-full border px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.24em] ${badgeToneClass[props.badgeTone || 'default']}`}>{props.badge}</div> : null}
    </div>
  )
}

export function SkeletonAmountHeroField() {
  return (
    <div class="flex flex-col items-center gap-2 rounded-[1.5rem] border border-border/50 bg-panel/70 px-4 py-5 shadow-lg shadow-black/10 animate-pulse">
      <SkeletonText width="w-16" />
      <div class="flex w-full max-w-[16rem] flex-col items-center border-b-2 border-border pb-1.5">
        <SkeletonText width="w-10" />
        <SkeletonBlock class="mt-1 h-12 w-32 border-border/20 bg-surface/70 sm:h-14 sm:w-40" />
      </div>
      <SkeletonBlock class="h-7 w-28 rounded-full border-border/20 bg-surface/70" />
    </div>
  )
}
