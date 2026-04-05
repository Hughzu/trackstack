import createClient from 'openapi-fetch'

import { readAccessToken } from '../auth/token'
import { resolveBrowserApiUrl } from '../config/api'
import type { paths } from './generated/schema'

const runtimeFetch: typeof fetch = async (input, init) => {
  const token = readAccessToken()
  const headers = new Headers(init?.headers)

  headers.set('Accept', 'application/json')

  if (token) {
    headers.set('X-Trackstack-Authorization', `Bearer ${token}`)
  }

  if (typeof input === 'string') {
    return fetch(resolveBrowserApiUrl(input), {
      ...init,
      headers,
    })
  }

  if (input instanceof Request) {
    return fetch(
      new Request(resolveBrowserApiUrl(input.url), {
        ...input,
        headers,
      }),
    )
  }

  return fetch(input, {
    ...init,
    headers,
  })
}

export const apiClient = createClient<paths>({
  baseUrl: '',
  fetch: runtimeFetch,
})

const toErrorMessage = (error: unknown, fallback: string) => {
  if (error && typeof error === 'object' && 'error' in error && typeof error.error === 'string') {
    return error.error
  }

  return fallback
}

export const unwrap = <T>(data: T | undefined, error: unknown, fallback: string) => {
  if (error || data === undefined) {
    throw new Error(toErrorMessage(error, fallback))
  }

  return data
}
