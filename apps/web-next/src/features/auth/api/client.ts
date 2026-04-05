import { apiClient, unwrap } from '../../../core/api/client'
import { clearAccessToken, writeAccessToken } from '../../../core/auth/token'
import type { LoginRequest, LoginResponse, SessionResponse } from '../../../core/api/types'

export const login = async (body: LoginRequest): Promise<LoginResponse> => {
  const { data, error } = await apiClient.POST('/api/auth/login', { body })
  const response = unwrap(data, error, 'Unable to log in')

  writeAccessToken(response.accessToken)

  return response
}

export const logout = async () => {
  const { error } = await apiClient.POST('/api/auth/logout')

  if (error) {
    throw new Error('Unable to log out')
  }

  clearAccessToken()
}

export const readSession = async (): Promise<SessionResponse> => {
  const { data, error } = await apiClient.GET('/api/auth/session')

  return unwrap(data, error, 'Unable to read session')
}
