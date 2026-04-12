import { apiClient, unwrap } from '../../../core/api/client'
import type { CreateRefillRequest, HeatDashboard, Refill } from '../../../core/api/types'

export const readHeatDashboard = async (): Promise<HeatDashboard> => {
  const { data, error } = await apiClient.GET('/api/heat/dashboard', {
    params: {
      query: {
        page: 1,
        limit: 20,
      },
    },
  })

  return unwrap(data, error, 'Unable to read heat dashboard')
}

export const listRefills = async (): Promise<Refill[]> => {
  const { data, error } = await apiClient.GET('/api/heat/refills')

  return unwrap(data, error, 'Unable to read heat refills')
}

export const createRefill = async (body: CreateRefillRequest): Promise<Refill> => {
  const { data, error } = await apiClient.POST('/api/heat/refills', { body })

  return unwrap(data, error, 'Unable to create heat refill')
}

export const deleteRefill = async (id: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/api/heat/refills/{id}', {
    params: {
      path: { id },
    },
  })

  if (error) {
    throw new Error('Unable to delete heat refill')
  }
}
