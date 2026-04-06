import { apiClient, unwrap } from '../../../core/api/client'
import type {
  ExpensesDashboard,
  ExpensesSettings,
  ExpensesSettingsView,
  UpdateExpensesSettingsRequest,
} from '../../../core/api/types'

export const readExpensesDashboard = async (): Promise<ExpensesDashboard> => {
  const { data, error } = await apiClient.GET('/api/expenses/sheet/current', {
    params: {
      query: {
        historyLimit: 50,
      },
    },
  })

  return unwrap(data, error, 'Unable to read expenses dashboard')
}

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
