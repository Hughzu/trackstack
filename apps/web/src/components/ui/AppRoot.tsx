import { Show, Suspense, createMemo } from 'solid-js'
import { useIsRouting, useLocation, type RouteSectionProps } from '@solidjs/router'

import { AppShell } from './AppShell'
import { RouteStatus } from './RouteStatus'

const resolveCurrentDomain = (pathname: string) => {
  if (pathname === '/login') {
    return 'auth'
  }

  if (pathname.startsWith('/calories')) {
    return 'calories'
  }

  if (pathname.startsWith('/expenses')) {
    return 'expenses'
  }

  if (pathname.startsWith('/heat')) {
    return 'heat'
  }

  return 'home'
}

export function AppRoot(props: RouteSectionProps) {
  const location = useLocation()
  const isRouting = useIsRouting()
  const currentDomain = createMemo(() => resolveCurrentDomain(location.pathname))

  return (
    <AppShell currentDomain={currentDomain()}>
      <Show when={!isRouting()} fallback={<RouteStatus />}>
        <Suspense fallback={<RouteStatus />}>
          {props.children}
        </Suspense>
      </Show>
    </AppShell>
  )
}
