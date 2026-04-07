import type { JSX } from 'solid-js'

type ListProps = {
  children: JSX.Element
  emptyMessage?: string
  variant?: 'card' | 'flush'
}

export function List(props: ListProps) {
  const baseClass = props.variant === 'flush'
    ? 'flex flex-col divide-y divide-border/20 -mx-5 px-5'
    : 'flex flex-col overflow-hidden rounded-2xl border border-border/50 bg-panel shadow-lg shadow-black/20 divide-y divide-border/30'

  return (
    <div class={baseClass}>
      {props.children || (props.emptyMessage ? <div class="p-4 text-center text-xs italic text-text-muted">{props.emptyMessage}</div> : null)}
    </div>
  )
}

type ListItemProps = {
  title: string
  subtitle?: JSX.Element | string
  value?: JSX.Element | string
  action?: JSX.Element
  prefix?: JSX.Element
}

export function ListItem(props: ListItemProps) {
  return (
    <div class="group flex items-center justify-between p-4 transition-colors hover:bg-surface/50">
      <div class="flex min-w-0 items-center gap-3">
        {props.prefix}
        <div class="flex flex-col min-w-0">
          <span class="text-[0.95rem] font-bold text-text-main truncate tracking-tight">{props.title}</span>
          {props.subtitle ? <span class="text-[0.68rem] text-text-muted truncate">{props.subtitle}</span> : null}
        </div>

      </div>
      <div class="flex items-center gap-3 pl-3">
        {props.value ? <span class="text-sm font-mono whitespace-nowrap">{props.value}</span> : null}
        {props.action}
      </div>
    </div>
  )
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
