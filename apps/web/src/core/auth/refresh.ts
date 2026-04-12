import { resolveBrowserApiUrl } from '../config/api'
import { attachPayloadHashHeader } from '../api/payload-hash'
import { clearAccessToken, writeAccessToken } from './token'

type RefreshResult = {
  accessToken: string
  expiresAt?: string
  tokenType?: string
  userId?: string
}

const isRefreshResult = (value: unknown): value is RefreshResult => {
  if (!value || typeof value !== 'object') {
    return false
  }

  return 'accessToken' in value && typeof value.accessToken === 'string'
}

export const refreshAccessToken = async (): Promise<RefreshResult | null> => {
  try {
    const headers = new Headers({
      Accept: 'application/json',
    })

    await attachPayloadHashHeader(headers, 'POST')

    const response = await fetch(resolveBrowserApiUrl('/api/auth/refresh'), {
      method: 'POST',
      headers,
      credentials: 'include',
    })

    if (!response.ok) {
      clearAccessToken()
      return null
    }

    const payload: unknown = await response.json()

    if (!isRefreshResult(payload)) {
      clearAccessToken()
      return null
    }

    writeAccessToken(payload.accessToken)

    return payload
  } catch {
    clearAccessToken()
    return null
  }
}
