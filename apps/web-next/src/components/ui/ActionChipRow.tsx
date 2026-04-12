import { For } from 'solid-js'

import { SkeletonBlock } from './Skeleton'

export type ActionChipItem = {
  id: string
  label: string
  ariaLabel?: string
  disabled?: boolean
  onClick: () => void
}

type ActionChipRowProps = {
  items: ActionChipItem[]
  emptyMessage?: string
}

export function ActionChipRow(props: ActionChipRowProps) {
  if (props.items.length === 0) {
    return <div class="py-2 text-center text-xs italic text-text-muted">{props.emptyMessage ?? 'Nothing here yet.'}</div>
  }

  return (
    <div class="overflow-x-auto pb-2 [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
      <div class="flex min-w-max gap-2 pt-1">
        <For each={props.items}>
          {(item) => (
            <button
              type="button"
              aria-label={item.ariaLabel}
              disabled={item.disabled}
              onClick={item.onClick}
              class="inline-flex h-[38px] min-w-[8rem] items-center justify-center rounded-full border border-border/70 bg-panel/60 px-4 py-2 text-sm font-medium text-text-main transition hover:border-accent/40 hover:text-accent disabled:cursor-not-allowed disabled:opacity-50"
            >
              {item.label}
            </button>
          )}
        </For>
      </div>
    </div>
  )
}

export function SkeletonActionChipRow(props: { count?: number }) {
  const items = Array.from({ length: props.count ?? 3 }, (_, index) => index)

  return (
    <div class="overflow-x-auto pb-2 [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
      <div class="flex min-w-max gap-2 pt-1">
        <For each={items}>
          {() => <SkeletonBlock class="h-[38px] min-w-[8rem] rounded-full border-border/20 bg-surface/70" />}
        </For>
      </div>
    </div>
  )
}
