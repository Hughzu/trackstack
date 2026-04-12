import type { JSX } from 'solid-js'
import { Match, Switch } from 'solid-js'
import { Navigate } from '@solidjs/router'

import { RouteStatus } from '../../components/ui/RouteStatus'
import { authState } from './state'

type GuardProps = {
  children: JSX.Element
}

export function ProtectedRoute(props: GuardProps) {
  return (
    <Switch>
      <Match when={authState().status === 'authenticated'}>{props.children}</Match>
      <Match when={authState().status === 'guest'}>
        <Navigate href="/login" />
      </Match>
      <Match when>
        <RouteStatus />
      </Match>
    </Switch>
  )
}

export function PublicOnlyRoute(props: GuardProps) {
  return (
    <Switch>
      <Match when={authState().status === 'guest'}>{props.children}</Match>
      <Match when={authState().status === 'authenticated'}>
        <Navigate href="/" />
      </Match>
      <Match when>
        <RouteStatus />
      </Match>
    </Switch>
  )
}
