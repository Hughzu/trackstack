import type { components } from './generated/schema'

type SchemaMap = components['schemas']

export type LoginRequest = SchemaMap['LoginRequest']
export type LoginResponse = SchemaMap['LoginResponse']
export type SessionResponse = SchemaMap['SessionResponse']
export type Refill = SchemaMap['Refill']
export type CreateRefillRequest = SchemaMap['CreateRefillRequest']
export type CalorieTarget = SchemaMap['CalorieTarget']
export type CreateCalorieLogRequest = SchemaMap['CreateCalorieLogRequest']
export type CalorieLog = SchemaMap['CalorieLog']
export type ExpensesSettings = SchemaMap['ExpensesSettings']
export type ExpensesSettingsView = SchemaMap['ExpensesSettingsView']
export type UpdateExpensesSettingsRequest = SchemaMap['UpdateExpensesSettingsRequest']
