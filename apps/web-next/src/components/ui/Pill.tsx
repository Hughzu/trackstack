import type { JSX } from 'solid-js'

type PillProps = {
  tone?: 'accent' | 'neutral' | 'success' | 'warning' | 'danger'
  size?: 'sm' | 'md'
  children: JSX.Element
}

export function Pill(props: PillProps) {
  let className = props.size === 'sm' 
    ? 'inline-flex w-fit rounded-full px-2 py-0.5 text-[0.55rem] font-bold uppercase tracking-wider'
    : 'inline-flex w-fit rounded-full px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.24em]'

  switch (props.tone) {
    case 'neutral':
      className += ' border border-border bg-panel text-text-muted'
      break
    case 'success':
      className += ' border border-success/20 bg-success/10 text-success'
      break
    case 'warning':
      className += ' border border-warning/20 bg-warning/10 text-warning'
      break
    case 'danger':
      className += ' border border-danger/20 bg-danger/10 text-danger'
      break
    default:
      className += ' border border-accent/20 bg-accent/10 text-accent'
      break
  }

  return <span class={className}>{props.children}</span>
}

type CounterPillProps = {
  tone?: 'accent' | 'neutral' | 'success' | 'warning' | 'danger'
  value: JSX.Element | string | number
  label: JSX.Element | string
}

export function CounterPill(props: CounterPillProps) {
  return <Pill tone={props.tone} size="sm">{props.value} {props.label}</Pill>
}
