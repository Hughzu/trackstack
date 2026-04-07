import { For } from 'solid-js'

export type FilterOption = {
  value: string
  label: string
  tone?: 'neutral' | 'danger' | 'warning' | 'success'
}

type FilterGroupProps = {
  options: FilterOption[]
  value: string
  onChange: (value: string) => void
}

export function FilterGroup(props: FilterGroupProps) {
  return (
    <div class="flex flex-wrap items-center justify-end gap-1.5">
      <For each={props.options}>
        {(option) => {
          const isActive = () => props.value === option.value
          const getClasses = () => {
            let base = 'rounded-full px-2.5 py-1 text-[10px] uppercase tracking-widest transition-colors outline-none focus-visible:ring-2 focus-visible:ring-accent '
            if (isActive()) {
              if (option.tone === 'danger') return base + 'bg-danger/20 text-danger border border-danger/50'
              if (option.tone === 'warning') return base + 'bg-warning/20 text-warning border border-warning/50'
              if (option.tone === 'success') return base + 'bg-success/20 text-success border border-success/50'
              return base + 'bg-surface text-text-main border border-border/50 shadow-sm'
            }
            // Inactive
            if (option.tone === 'danger') return base + 'text-danger border border-danger/30 opacity-50 hover:bg-danger/10 hover:opacity-100'
            if (option.tone === 'warning') return base + 'text-warning border border-warning/30 opacity-50 hover:bg-warning/10 hover:opacity-100'
            if (option.tone === 'success') return base + 'text-success border border-success/30 opacity-50 hover:bg-success/10 hover:opacity-100'
            return base + 'text-text-muted border border-transparent hover:text-text-main hover:bg-surface/50'
          }
          return (
            <button type="button" class={getClasses()} onClick={() => props.onChange(option.value)}>
              {option.label}
            </button>
          )
        }}
      </For>
    </div>
  )
}
