import type { CaloriesDashboard, CalorieLog, CalorieTarget } from '../../core/api/types'
import { formatCount } from '../../core/format/number'

export type CaloriesQuickAddItem = {
  id: string
  label: string
  ariaLabel: string
}

export type MacroTone = 'success' | 'warning' | 'danger'

type MacroDefinition = {
  key: 'proteinGrams' | 'carbGrams' | 'fatGrams'
  title: string
  shortTitle: string
  tone: MacroTone
}

const macroDefinitions: MacroDefinition[] = [
  { key: 'proteinGrams', title: 'Protein', shortTitle: 'Prot.', tone: 'success' },
  { key: 'carbGrams', title: 'Carbs', shortTitle: 'Carbs', tone: 'warning' },
  { key: 'fatGrams', title: 'Fat', shortTitle: 'Fat', tone: 'danger' },
]

const timeFormatter = new Intl.DateTimeFormat('en-GB', {
  hour: '2-digit',
  minute: '2-digit',
})

export function createCaloriesProgressSegments(dashboard: CaloriesDashboard) {
  const targetCalories = dashboard.summary.target.targetCalories
  const consumedCalories = dashboard.summary.consumedCalories
  const consumedPercent = targetCalories > 0 ? Math.min((consumedCalories / targetCalories) * 100, 100) : 0

  return [{ percent: consumedPercent, color: 'accent' as const }]
}

export function createCaloriesHeroModel(dashboard: CaloriesDashboard) {
  const consumedCalories = dashboard.summary.consumedCalories
  const targetCalories = dashboard.summary.target.targetCalories
  const targetProteinGrams = dashboard.summary.target.targetProteinGrams
  const remainingCalories = Math.max(targetCalories - consumedCalories, 0)
  const progressPercent = Math.round(createCaloriesProgressSegments(dashboard)[0]?.percent ?? 0)

  return {
    consumedCalories: formatCount(consumedCalories),
    targetCalories: formatCount(targetCalories),
    proteinGrams: formatCount(dashboard.summary.proteinGrams),
    targetProteinGrams: formatCount(targetProteinGrams),
    remainingLabel: `${formatCount(remainingCalories)} kcal left`,
    progressPercent,
  }
}

export function createCaloriesMacroItems(dashboard: CaloriesDashboard) {
  return macroDefinitions.map((macro) => ({
    key: macro.key,
    title: macro.title,
    shortTitle: macro.shortTitle,
    tone: macro.tone,
    value: `${formatCount(Number(dashboard.summary[macro.key] ?? 0))}g`,
  }))
}

export function createRecentMealChips(dashboard: CaloriesDashboard): CaloriesQuickAddItem[] {
  return dashboard.recentMeals.map((meal) => ({
    id: meal.id,
    label: meal.title?.trim() || 'Untitled',
    ariaLabel: `Quick add ${meal.title?.trim() || 'meal'}`,
  }))
}

export function formatCalorieLogTime(dateTime: string) {
  return timeFormatter.format(new Date(dateTime))
}

export function createCalorieLogDeleteDescription(log: CalorieLog) {
  const title = normalizeMealTitle(log.title) || 'Untitled'
  return `Delete ${title} for ${formatCount(log.calories)} kcal from ${formatCalorieLogTime(log.dateTime)}? This removes it from today.`
}

export function createTargetBadge(target: CalorieTarget) {
  return `${formatCount(target.targetProteinGrams)}g protein`
}

export function normalizeCalorieTargetInput(value: string) {
  return value.trim()
}

function normalizeMealTitle(title: string | null | undefined) {
  return title?.trim() || ''
}
