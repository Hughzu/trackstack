import { apiClient, unwrap } from '../../../core/api/client'
import type {
  CalorieLog,
  CalorieTarget,
  CaloriesDashboard,
  CreateCalorieLogRequest,
  UpdateCalorieTargetRequest,
} from '../../../core/api/types'

export const readCaloriesDashboard = async (): Promise<CaloriesDashboard> => {
  const { data, error } = await apiClient.GET('/api/calories/dashboard', {
    params: {
      query: {
        recentLimit: 8,
        logsLimit: 50,
      },
    },
  })

  return unwrap(data, error, 'Unable to read calories dashboard')
}

export const readCalorieTarget = async (): Promise<CalorieTarget> => {
  const { data, error } = await apiClient.GET('/api/calories/target')

  return unwrap(data, error, 'Unable to read calorie target')
}

export const createCalorieLog = async (body: CreateCalorieLogRequest): Promise<CalorieLog> => {
  const { data, error } = await apiClient.POST('/api/calories/log', { body })

  return unwrap(data, error, 'Unable to create calorie log')
}

export const updateCalorieTarget = async (body: UpdateCalorieTargetRequest): Promise<CalorieTarget> => {
  const { data, error } = await apiClient.POST('/api/calories/target', { body })

  return unwrap(data, error, 'Unable to update calorie target')
}

export const deleteCalorieLog = async (id: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/api/calories/logs/{id}', {
    params: {
      path: { id },
    },
  })

  if (error) {
    throw new Error('Unable to delete calorie log')
  }
}
