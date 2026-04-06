import type { JSX } from 'solid-js'
import { createMemo, createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { authState, logout } from '../../core/auth/store'
import { appConfig } from '../../core/config/app'
import { resolveTheme } from '../../core/config/theme'
import { ActionButton } from './ActionButton'
import { DomainTabs } from './DomainTabs'

type AppShellProps = {
  currentDomain: string
  children: JSX.Element
}

export function AppShell(props: AppShellProps) {
  const navigate = useNavigate()
  const theme = resolveTheme(import.meta.env.VITE_DEPLOY_TARGET)
  const [isLoggingOut, setIsLoggingOut] = createSignal(false)
  const visibleDomains = createMemo(() => {
    if (authState().status === 'authenticated') {
      return appConfig.domains.filter((domain) => domain.id !== 'auth')
    }

    return appConfig.domains.filter((domain) => domain.id === 'auth')
  })

  const handleLogout = async () => {
    if (isLoggingOut()) {
      return
    }

    setIsLoggingOut(true)

    try {
      await logout()
      void navigate('/login', { replace: true })
    } finally {
      setIsLoggingOut(false)
    }
  }

  return (
    <div class="min-h-screen bg-background text-text-main antialiased selection:bg-accent/30">
      <header class="sticky top-0 z-40 border-b border-transparent bg-background/95 backdrop-blur-md">
        <div class="mx-auto flex max-w-3xl items-center justify-between px-4 py-3">
          <div>
            <div class="text-xl font-bold tracking-tight text-text-main">{appConfig.appName}</div>
            <div class="text-[0.68rem] uppercase tracking-[0.28em] text-text-muted">{theme.label}</div>
          </div>

          <div class="flex items-center gap-3">
            {authState().status === 'authenticated' ? (
              <ActionButton tone="ghost" disabled={isLoggingOut()} busy={isLoggingOut()} onClick={handleLogout}>
                {isLoggingOut() ? 'Logging out...' : 'Logout'}
              </ActionButton>
            ) : null}
          </div>
        </div>

        <DomainTabs currentDomain={props.currentDomain} domains={visibleDomains()} />
      </header>

      <main class="mx-auto flex min-h-[60vh] max-w-3xl flex-col gap-6 px-4 pb-6 pt-6 animate-rise-in">
        {props.children}
      </main>
    </div>
  )
}
