import { apiClient, unwrap } from '../../../core/api/client'
import type { ExpensesSettings, ExpensesSettingsView, UpdateExpensesSettingsRequest } from '../../../core/api/types'

export const readExpensesSettings = async (): Promise<ExpensesSettingsView> => {
  const { data, error } = await apiClient.GET('/api/expenses/settings')

  return unwrap(data, error, 'Unable to read expenses settings')
}

export const updateExpensesSettings = async (
  body: UpdateExpensesSettingsRequest,
): Promise<ExpensesSettings> => {
  const { data, error } = await apiClient.POST('/api/expenses/settings', { body })

  return unwrap(data, error, 'Unable to update expenses settings')
}
