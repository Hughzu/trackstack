import type { JSX } from 'solid-js'
import { children } from 'solid-js'

type ListProps = {
  children: JSX.Element
  emptyMessage?: string
  variant?: 'card' | 'flush'
}

export function List(props: ListProps) {
  const resolvedChildren = children(() => props.children)
  const baseClass = props.variant === 'flush'
    ? 'flex flex-col divide-y divide-border/20 -mx-5 px-5'
    : 'flex flex-col overflow-hidden rounded-2xl border border-border/50 bg-panel shadow-lg shadow-black/20 divide-y divide-border/30'
  const hasContent = () => resolvedChildren.toArray().length > 0

  return (
    <div class={baseClass}>
      {hasContent() ? resolvedChildren() : (props.emptyMessage ? <div class="p-4 text-center text-xs italic text-text-muted">{props.emptyMessage}</div> : null)}
    </div>
  )
}

type ListItemProps = {
  title: string
  subtitle?: JSX.Element | string
  value?: JSX.Element | string
  valueTone?: 'default' | 'danger' | 'muted' | 'success' | 'warning'
  valueStyle?: 'default' | 'mono'
  action?: JSX.Element
  prefix?: JSX.Element
}

const valueToneClass = {
  default: 'text-text-main',
  danger: 'text-danger',
  muted: 'text-text-muted',
  success: 'text-success',
  warning: 'text-warning',
}

export function ListItem(props: ListItemProps) {
  const subtitle = () => {
    if (!props.subtitle) return null
    if (typeof props.subtitle === 'string') {
      return <span class="text-[0.68rem] text-text-muted truncate">{props.subtitle}</span>
    }

    return props.subtitle
  }

  const value = () => {
    if (!props.value) return null
    if (typeof props.value !== 'string') return props.value

    const tone = valueToneClass[props.valueTone || 'default']
    const font = props.valueStyle === 'mono' ? 'font-mono' : 'font-semibold'

    return <span class={`text-sm whitespace-nowrap ${font} ${tone}`}>{props.value}</span>
  }

  return (
    <div class="group flex items-center justify-between p-4 transition-colors hover:bg-surface/50">
      <div class="flex min-w-0 items-center gap-3">
        {props.prefix}
        <div class="flex flex-col min-w-0">
          <span class="text-[0.95rem] font-bold text-text-main truncate tracking-tight">{props.title}</span>
          {subtitle()}
        </div>

      </div>
      <div class="flex items-center gap-3 pl-3">
        {value()}
        {props.action}
      </div>
    </div>
  )
}

export function ListMeta(props: { children: JSX.Element }) {
  return <div class="mt-1.5 flex flex-wrap items-center gap-1.5 text-[0.68rem] text-text-muted">{props.children}</div>
}

export function ListMetaDivider() {
  return <span class="text-border/50">•</span>
}

export function SkeletonListItem() {
  return (
    <div class="flex items-center justify-between p-4">
      <div class="flex flex-col gap-2">
        <div class="h-4 w-32 rounded bg-white/10 animate-pulse" />
        <div class="h-3 w-16 rounded bg-white/10 animate-pulse" />
      </div>
      <div class="h-4 w-12 rounded bg-white/10 animate-pulse" />
    </div>
  )
}
