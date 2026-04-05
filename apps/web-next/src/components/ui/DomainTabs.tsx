import { For } from 'solid-js'

import type { AppDomain } from '../../core/config/app'

type DomainTabsProps = {
  currentDomain: string
  domains: AppDomain[]
}

export function DomainTabs(props: DomainTabsProps) {
  return (
    <nav class="mx-auto flex max-w-3xl gap-6 overflow-x-auto px-4 pb-2 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
      <div class="flex min-w-max gap-6">
        <For each={props.domains}>
          {(domain) => (
            <a
              href={domain.href}
              aria-current={props.currentDomain === domain.id ? 'page' : undefined}
              class={
                props.currentDomain === domain.id
                  ? 'select-none whitespace-nowrap border-b-2 border-accent pb-1 text-sm font-medium text-accent transition-all duration-200'
                  : 'select-none whitespace-nowrap pb-1 text-sm font-medium text-text-muted transition-all duration-200 hover:text-text-main'
              }
            >
              {domain.label}
            </a>
          )}
        </For>
      </div>
    </nav>
  )
}
