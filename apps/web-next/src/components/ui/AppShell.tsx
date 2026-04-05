import type { JSX } from 'solid-js'

import { appConfig } from '../../core/config/app'
import { resolveTheme } from '../../core/config/theme'
import { DomainTabs } from './DomainTabs'

type AppShellProps = {
  currentDomain: string
  eyebrow: string
  title: string
  description: string
  children: JSX.Element
}

export function AppShell(props: AppShellProps) {
  const theme = resolveTheme(import.meta.env.VITE_DEPLOY_TARGET)
  const initials = appConfig.appName
    .split(' ')
    .map((part) => part[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()

  return (
    <div class="min-h-screen bg-background text-text-main antialiased selection:bg-accent/30">
      <header class="sticky top-0 z-40 border-b border-transparent bg-background/95 backdrop-blur-md">
        <div class="mx-auto flex max-w-3xl items-center justify-between px-4 py-3">
          <div>
            <div class="text-xl font-bold tracking-tight text-text-main">{appConfig.appName}</div>
            <div class="text-[0.68rem] uppercase tracking-[0.28em] text-text-muted">{theme.label}</div>
          </div>

          <div class="flex h-9 w-9 items-center justify-center rounded-full border border-border bg-panel text-xs font-bold text-text-main">
              {initials}
          </div>
        </div>

        <DomainTabs currentDomain={props.currentDomain} domains={appConfig.domains} />
      </header>

      <main class="mx-auto flex min-h-[60vh] max-w-3xl flex-col gap-6 px-4 pb-24 pt-6 animate-rise-in">
        <section class="space-y-2">
          <div class="text-xs font-bold uppercase tracking-widest text-accent">{props.eyebrow}</div>
          <h1 class="text-3xl font-bold tracking-tight text-text-main sm:text-4xl">{props.title}</h1>
          <p class="max-w-2xl text-sm leading-6 text-text-muted sm:text-base">{props.description}</p>
        </section>

        {props.children}
      </main>
    </div>
  )
}
