import createClient from 'openapi-fetch'

import { refreshAccessToken } from '../auth/refresh'
import { setGuest } from '../auth/state'
import { readAccessToken } from '../auth/token'
import { resolveBrowserApiUrl } from '../config/api'
import { attachPayloadHashHeader } from './payload-hash'
import type { paths } from './generated/schema'

let refreshPromise: Promise<string | null> | null = null

const authPaths = new Set(['/api/auth/login', '/api/auth/logout', '/api/auth/session', '/api/auth/refresh'])

const resolvePathname = (value: string) => {
  try {
    return new URL(value, window.location.origin).pathname
  } catch {
    return value
  }
}

const createHeaders = (init: HeadersInit | undefined, token: string | null) => {
  const headers = new Headers(init)

  headers.set('Accept', 'application/json')

  if (token) {
    headers.set('X-Trackstack-Authorization', `Bearer ${token}`)
  }

  return headers
}

const executeRequest = async (input: RequestInfo | URL, init: RequestInit | undefined, token: string | null) => {
  const headers = createHeaders(init?.headers, token)

  if (typeof input === 'string') {
    const method = init?.method ?? 'GET'
    await attachPayloadHashHeader(headers, method, init?.body)

    return fetch(resolveBrowserApiUrl(input), {
      ...init,
      headers,
      credentials: init?.credentials ?? 'include',
    })
  }

  if (input instanceof Request) {
    const request = input.clone()
    const method = init?.method ?? request.method
    let body = init?.body

    if (body === undefined && method !== 'GET' && method !== 'HEAD') {
      const requestBody = await request.arrayBuffer()
      body = requestBody.byteLength > 0 ? requestBody : null
    }

    await attachPayloadHashHeader(headers, method, body)

    return fetch(resolveBrowserApiUrl(request.url), {
      method,
      headers,
      body,
      credentials: init?.credentials ?? request.credentials ?? 'include',
      cache: init?.cache ?? request.cache,
      mode: init?.mode ?? request.mode,
      redirect: init?.redirect ?? request.redirect,
      referrer: init?.referrer ?? request.referrer,
      referrerPolicy: init?.referrerPolicy ?? request.referrerPolicy,
      integrity: init?.integrity ?? request.integrity,
      keepalive: init?.keepalive ?? request.keepalive,
      signal: init?.signal ?? request.signal,
    })
  }

  const method = init?.method ?? 'GET'
  await attachPayloadHashHeader(headers, method, init?.body)

  return fetch(input, {
    ...init,
    headers,
    credentials: init?.credentials ?? 'include',
  })
}

const refreshOnce = async () => {
  if (!refreshPromise) {
    refreshPromise = refreshAccessToken()
      .then((payload) => payload?.accessToken ?? null)
      .finally(() => {
        refreshPromise = null
      })
  }

  return refreshPromise
}

const runtimeFetch: typeof fetch = async (input, init) => {
  const path = resolvePathname(typeof input === 'string' ? input : input instanceof Request ? input.url : String(input))
  const response = await executeRequest(input, init, readAccessToken())

  if (response.status !== 401 || authPaths.has(path)) {
    return response
  }

  const refreshedToken = await refreshOnce()

  if (!refreshedToken) {
    setGuest('Session expired')
    return response
  }

  return executeRequest(input, init, refreshedToken)
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
