import { apiClient, unwrap } from '../../../core/api/client'
import type {
  CompleteChecklistItemRequest,
  CreateExpenseEntryRequest,
  ExpenseEntry,
  ExpenseSheet,
  ExpenseTemplate,
  ExpensesDashboard,
  ExpensesSettings,
  ExpensesSettingsView,
  UpsertTemplateRequest,
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

export const createExpenseEntry = async (body: CreateExpenseEntryRequest): Promise<ExpenseEntry> => {
  const { data, error } = await apiClient.POST('/api/expenses/entries', { body })

  return unwrap(data, error, 'Unable to create expense entry')
}

export const deleteExpenseEntry = async (id: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/api/expenses/entries/{id}', {
    params: {
      path: { id },
    },
  })

  if (error) {
    throw new Error('Unable to delete expense entry')
  }
}

export const upsertChecklistTemplate = async (body: UpsertTemplateRequest): Promise<ExpenseTemplate> => {
  const { data, error } = await apiClient.POST('/api/expenses/checklists', { body })

  return unwrap(data, error, 'Unable to upsert checklist template')
}

export const deleteChecklistTemplate = async (id: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/api/expenses/checklists/{id}', {
    params: {
      path: { id },
    },
  })

  if (error) {
    throw new Error('Unable to delete checklist template')
  }
}

export const completeChecklistItem = async (body: CompleteChecklistItemRequest): Promise<ExpenseEntry> => {
  const { data, error } = await apiClient.POST('/api/expenses/checklists/complete', { body })

  return unwrap(data, error, 'Unable to complete checklist item')
}

export const upsertRecurringTemplate = async (body: UpsertTemplateRequest): Promise<ExpenseTemplate> => {
  const { data, error } = await apiClient.POST('/api/expenses/recurring', { body })

  return unwrap(data, error, 'Unable to upsert recurring template')
}

export const deleteRecurringTemplate = async (id: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/api/expenses/recurring/{id}', {
    params: {
      path: { id },
    },
  })

  if (error) {
    throw new Error('Unable to delete recurring template')
  }
}

export const closeExpenseSheet = async (): Promise<ExpenseSheet> => {
  const { data, error } = await apiClient.POST('/api/expenses/sheet/close')

  return unwrap(data, error, 'Unable to close expenses sheet')
}
