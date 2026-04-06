import type { JSX } from 'solid-js'
import { SkeletonBlock, SkeletonText } from './Skeleton'

export type StatProps = {
  label: string | JSX.Element
  labelPosition?: 'top' | 'bottom' | 'inline'
  value: string | JSX.Element
  unit?: string | JSX.Element   // Rendered inline next to value
  subtext?: string | JSX.Element // Rendered below or inline depending on variant
  align?: 'left' | 'right' | 'center'
  variant?: 'hero' | 'lg' | 'md' | 'sm' | 'mono' | 'inline'
  color?: 'main' | 'accent' | 'muted'
}

const colorMap = {
  main: 'text-text-main',
  accent: 'text-accent',
  muted: 'text-text-muted',
}

export function Stat(props: StatProps) {
  const alignClass = props.align === 'right' ? 'items-end text-right' : props.align === 'center' ? 'items-center text-center' : 'items-start text-left'
  const textColor = colorMap[props.color || 'main']
  const isInlineLabel = props.labelPosition === 'inline' || props.variant === 'inline'
  const labelClass = isInlineLabel 
    ? 'text-xs text-text-muted' 
    : 'mb-1 text-[0.65rem] uppercase tracking-widest text-text-muted'

  const labelUI = <div class={labelClass}>{props.label}</div>

  let valueUI: JSX.Element

  if (props.variant === 'inline') {
    return (
      <div class="flex items-center gap-2 text-xs">
        <span class="text-text-muted">{props.label}</span>
        <span class={`font-semibold ${textColor}`}>{props.value}</span>
      </div>
    )
  }

  if (props.variant === 'hero') {
    valueUI = (
      <div class={`flex items-baseline gap-2 ${props.labelPosition === 'bottom' ? 'mb-1' : ''}`}>
        <span class={`text-4xl font-bold tracking-tight ${textColor}`}>{props.value}</span>
        {props.unit ? <span class="text-xs uppercase tracking-widest text-text-muted">{props.unit}</span> : null}
      </div>
    )
  } else if (props.variant === 'lg') {
    valueUI = (
      <div class={`flex items-baseline gap-2 ${props.labelPosition === 'bottom' ? 'mb-1' : ''}`}>
        <span class={`text-2xl font-bold tracking-tight ${textColor}`}>{props.value}</span>
        {props.unit ? <span class="text-xs font-mono text-text-muted">{props.unit}</span> : null}
      </div>
    )
  } else if (props.variant === 'sm') {
    valueUI = (
      <div class={`mt-0.5 text-xs font-bold ${textColor}`}>
        {props.value}
        {props.unit ? <span class="ml-1">{props.unit}</span> : null}
        {props.subtext ? <div class="mt-0.5 text-[0.65rem] font-normal text-text-muted">{props.subtext}</div> : null}
      </div>
    )
  } else if (props.variant === 'mono') {
    valueUI = (
      <div class={`text-xs font-mono ${textColor}`}>
        {props.value}
      </div>
    )
  } else {
    // md fallback
    valueUI = (
      <div class={`text-sm font-bold ${textColor}`}>
        {props.value}
        {props.unit ? <span class="ml-1 text-xs font-mono font-normal text-text-muted">{props.unit}</span> : null}
      </div>
    )
  }

  return (
    <div class={`flex flex-col ${alignClass}`}>
      {props.labelPosition !== 'bottom' ? labelUI : null}
      {valueUI}
      {props.labelPosition === 'bottom' ? labelUI : null}
    </div>
  )
}

export type SkeletonStatProps = Omit<StatProps, 'label' | 'value'>

export function SkeletonStat(props: SkeletonStatProps) {
  let h = 'h-6'
  let w = 'w-16'
  let lw = 'w-16'
  
  switch(props.variant) {
    case 'hero': h = 'h-10'; w = 'w-24'; lw = 'w-20'; break;
    case 'lg': h = 'h-8'; w = 'w-24'; lw = 'w-16'; break;
    case 'md': h = 'h-5'; w = 'w-20'; lw = 'w-14'; break;
    case 'sm': h = 'h-4'; w = 'w-8'; lw = 'w-8'; break;
    case 'mono': h = 'h-4'; w = 'w-16'; lw = 'w-12'; break;
    case 'inline': h = 'h-4'; w = 'w-16'; lw = 'w-16'; break;
  }

  const lbl = <SkeletonText class={`h-3 ${lw}`} />
  
  return (
    <Stat
      {...props}
      label={lbl}
      value={<SkeletonBlock class={`${h} ${w}`} />}
    />
  )
}
