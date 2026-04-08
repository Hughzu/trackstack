import { For, type JSX } from 'solid-js'

type ChoiceCardTone = 'default' | 'danger' | 'warning' | 'success'

export type ChoiceCardOption = {
  value: string
  title: string
  description?: string
  icon?: JSX.Element
  tone?: ChoiceCardTone
}

type ChoiceCardsProps = {
  options: ChoiceCardOption[]
  value: string
  testId?: string
  onChange: (value: string) => void
}

export function ChoiceCards(props: ChoiceCardsProps) {
  return (
    <div class="grid gap-2 sm:grid-cols-3" data-testid={props.testId}>
      <For each={props.options}>
        {(option) => {
          const active = () => props.value === option.value
          const classes = () => {
            if (!active()) return 'border-border/70 bg-panel/80 text-text-muted'
            if (option.tone === 'warning') return 'border-warning/50 bg-warning/10 text-warning shadow-lg shadow-warning/5'
            if (option.tone === 'success') return 'border-success/50 bg-success/10 text-success shadow-lg shadow-success/5'
            if (option.tone === 'danger') return 'border-danger/50 bg-danger/10 text-danger shadow-lg shadow-danger/5'
            return 'border-accent/50 bg-accent/10 text-accent shadow-lg shadow-accent/5'
          }

          return (
            <button
              type="button"
              onClick={() => props.onChange(option.value)}
              class={`rounded-2xl border px-4 py-2.5 text-left transition-all duration-200 hover:border-accent/40 ${classes()}`}
              aria-pressed={active()}
            >
              <div class="mb-1 flex items-center justify-between gap-2">
                <span class="text-[0.95rem] font-bold tracking-tight">{option.title}</span>
                {option.icon ? <span class="flex h-7 w-7 items-center justify-center rounded-full border border-current/15 bg-black/5">{option.icon}</span> : null}
              </div>
              {option.description ? <div class="text-[0.66rem] uppercase tracking-[0.18em] opacity-70">{option.description}</div> : null}
            </button>
          )
        }}
      </For>
    </div>
  )
}
