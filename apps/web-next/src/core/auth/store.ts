import type { LoginRequest } from '../api/types'
import { refreshAccessToken } from './refresh'
import { authState, setAuthenticated, setAuthChecking, setGuest } from './state'
import { clearAccessToken, markLoggedOut, readAccessToken, shouldAttemptRefresh, writeAccessToken } from './token'
import { loginRequest, logoutRequest, readSessionRequest } from './transport'

let bootstrapPromise: Promise<void> | null = null

const toErrorMessage = (error: unknown, fallback: string) => {
  if (error instanceof Error && error.message) {
    return error.message
  }

  return fallback
}

const readLiveSession = async () => {
  const session = await readSessionRequest()
  setAuthenticated(session)
}

export const bootstrapAuth = async () => {
  if (bootstrapPromise) {
    return bootstrapPromise
  }

  setAuthChecking()

  bootstrapPromise = (async () => {
    if (!readAccessToken()) {
      if (!shouldAttemptRefresh()) {
        setGuest()
        return
      }

      const refreshed = await refreshAccessToken()
      if (!refreshed) {
        setGuest()
        return
      }
    }

    try {
      await readLiveSession()
    } catch (error) {
      if (!shouldAttemptRefresh()) {
        clearAccessToken()
        setGuest(toErrorMessage(error, 'Session expired'))
        return
      }

      const refreshed = await refreshAccessToken()
      if (!refreshed) {
        clearAccessToken()
        setGuest(toErrorMessage(error, 'Session expired'))
        return
      }

      try {
        await readLiveSession()
      } catch (retryError) {
        clearAccessToken()
        setGuest(toErrorMessage(retryError, 'Session expired'))
      }
    }
  })().finally(() => {
    bootstrapPromise = null
  })

  return bootstrapPromise
}

export const login = async (body: LoginRequest) => {
  setAuthChecking()

  try {
    const response = await loginRequest(body)
    writeAccessToken(response.accessToken)
    await readLiveSession()
  } catch (error) {
    clearAccessToken()
    setGuest(toErrorMessage(error, 'Unable to log in'))
    throw error
  }
}

export const logout = async () => {
  try {
    await logoutRequest()
  } finally {
    markLoggedOut()
    clearAccessToken()
    setGuest()
  }
}

export { authState }
