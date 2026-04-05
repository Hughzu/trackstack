import { apiClient, unwrap } from '../../../core/api/client'
import type { CalorieLog, CalorieTarget, CreateCalorieLogRequest } from '../../../core/api/types'

export const readCalorieTarget = async (): Promise<CalorieTarget> => {
  const { data, error } = await apiClient.GET('/api/calories/target')

  return unwrap(data, error, 'Unable to read calorie target')
}

export const createCalorieLog = async (body: CreateCalorieLogRequest): Promise<CalorieLog> => {
  const { data, error } = await apiClient.POST('/api/calories/log', { body })

  return unwrap(data, error, 'Unable to create calorie log')
}
