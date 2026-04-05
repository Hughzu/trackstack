import { apiClient, unwrap } from '../../../core/api/client'
import type { CreateRefillRequest, Refill } from '../../../core/api/types'

export const listRefills = async (): Promise<Refill[]> => {
  const { data, error } = await apiClient.GET('/api/heat/refills')

  return unwrap(data, error, 'Unable to read heat refills')
}

export const createRefill = async (body: CreateRefillRequest): Promise<Refill> => {
  const { data, error } = await apiClient.POST('/api/heat/refills', { body })

  return unwrap(data, error, 'Unable to create heat refill')
}
