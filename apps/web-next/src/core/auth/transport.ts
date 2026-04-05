import { apiClient, unwrap } from '../api/client'
import type { LoginRequest, LoginResponse, SessionResponse } from '../api/types'

export const loginRequest = async (body: LoginRequest): Promise<LoginResponse> => {
  const { data, error } = await apiClient.POST('/api/auth/login', { body })

  return unwrap(data, error, 'Unable to log in')
}

export const logoutRequest = async () => {
  const { error } = await apiClient.POST('/api/auth/logout')

  if (error) {
    throw new Error('Unable to log out')
  }
}

export const readSessionRequest = async (): Promise<SessionResponse> => {
  const { data, error } = await apiClient.GET('/api/auth/session')

  return unwrap(data, error, 'Unable to read session')
}
