const accessTokenKey = 'trackstack_token'
const logoutMarkerKey = 'trackstack_logged_out'

let accessTokenCache: string | null = null

const readSessionStorage = (key: string) => {
  if (typeof window === 'undefined') {
    return null
  }

  return window.sessionStorage.getItem(key)
}

const writeSessionStorage = (key: string, value: string) => {
  if (typeof window === 'undefined') {
    return
  }

  window.sessionStorage.setItem(key, value)
}

const removeSessionStorage = (key: string) => {
  if (typeof window === 'undefined') {
    return
  }

  window.sessionStorage.removeItem(key)
}

export const readAccessToken = () => {
  if (accessTokenCache) {
    return accessTokenCache
  }

  accessTokenCache = readSessionStorage(accessTokenKey)

  return accessTokenCache
}

export const writeAccessToken = (token: string) => {
  accessTokenCache = token
  writeSessionStorage(accessTokenKey, token)
  clearLogoutMarker()
}

export const clearAccessToken = () => {
  accessTokenCache = null
  removeSessionStorage(accessTokenKey)
}

export const markLoggedOut = () => {
  writeSessionStorage(logoutMarkerKey, '1')
}

export const clearLogoutMarker = () => {
  removeSessionStorage(logoutMarkerKey)
}

export const shouldAttemptRefresh = () => readSessionStorage(logoutMarkerKey) !== '1'
