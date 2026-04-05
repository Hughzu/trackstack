const configuredApiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').trim()

const normalizeConfiguredBaseUrl = (value: string) => {
  return value.replace(/\/+$/, '').replace(/\/api$/, '')
}

export const apiConfig = {
  baseUrl: normalizeConfiguredBaseUrl(configuredApiBaseUrl),
}

export const resolveBrowserApiUrl = (path: string) => {
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }

  if (!apiConfig.baseUrl) {
    return path
  }

  return `${apiConfig.baseUrl}${path}`
}
