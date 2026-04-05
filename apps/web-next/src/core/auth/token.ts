// Bootstrap-only token storage. Replace with an httpOnly session strategy before shipping real auth.
const accessTokenKey = 'trackstack_token'

export const readAccessToken = () => window.localStorage.getItem(accessTokenKey)

export const writeAccessToken = (token: string) => {
  window.localStorage.setItem(accessTokenKey, token)
}

export const clearAccessToken = () => {
  window.localStorage.removeItem(accessTokenKey)
}
