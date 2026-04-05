import { createSignal } from 'solid-js'

import type { SessionResponse } from '../api/types'

export type AuthStatus = 'checking' | 'authenticated' | 'guest'

export type AuthState = {
  status: AuthStatus
  session: SessionResponse | null
  initialized: boolean
  error: string | null
}

const [authState, setAuthState] = createSignal<AuthState>({
  status: 'checking',
  session: null,
  initialized: false,
  error: null,
})

export { authState }

export const setAuthChecking = () => {
  setAuthState((current) => ({
    ...current,
    status: 'checking',
    error: null,
  }))
}

export const setAuthenticated = (session: SessionResponse) => {
  setAuthState({
    status: 'authenticated',
    session,
    initialized: true,
    error: null,
  })
}

export const setGuest = (error: string | null = null) => {
  setAuthState({
    status: 'guest',
    session: null,
    initialized: true,
    error,
  })
}
