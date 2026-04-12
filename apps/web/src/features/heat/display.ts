import type { HeatDashboard, Refill } from '../../core/api/types'
import { formatCount } from '../../core/format/number'

const dayDateFormatter = new Intl.DateTimeFormat('en-IE', {
  month: 'short',
  day: 'numeric',
})

export function formatHeatDate(value: string) {
  const normalizedValue = value.includes('T') ? value : `${value}T00:00:00`
  return dayDateFormatter.format(new Date(normalizedValue))
}

export function formatHeatTemperature(value: number | null | undefined) {
  return value == null ? 'No temperature' : `${formatCount(value)} C`
}

export function formatHeatBags(value: number) {
  return `${formatCount(value)} bag${value === 1 ? '' : 's'}`
}

export function formatHeatWeight(value: number) {
  return `${formatCount(value)} kg`
}

export function createHeatHeroModel(dashboard: HeatDashboard) {
  const daysSinceRefill = dashboard.daysSinceRefill
  const latestRefill = dashboard.history[0]
  const seasonSnapshot = dashboard.seasonSnapshot
  const needsRefillSoon = daysSinceRefill > 14

  return {
    seasonLabel: seasonSnapshot.seasonLabel,
    daysSinceRefill: formatCount(daysSinceRefill),
    latestRefillDate: latestRefill ? formatHeatDate(latestRefill.date) : 'No history',
    latestRefillTemperature: latestRefill ? formatHeatTemperature(latestRefill.temperature) : 'No temperature',
    statusTone: needsRefillSoon ? 'warning' as const : 'success' as const,
    statusLabel: needsRefillSoon ? 'Refill soon' : 'Fresh refill pace',
  }
}

export function createHeatSeasonComparisonModel(dashboard: HeatDashboard) {
  const snapshot = dashboard.seasonSnapshot
  const delta = snapshot.delta
  const deltaLabel = `${delta > 0 ? '+' : ''}${formatCount(delta)} bag${Math.abs(delta) === 1 ? '' : 's'}`
  const deltaTrend = snapshot.deltaPct == null ? 'vs last year' : `${delta > 0 ? '+' : ''}${snapshot.deltaPct}%`

  return {
    currentLabel: 'This season',
    currentValue: formatHeatBags(snapshot.seasonToDate),
    previousLabel: 'Last year pace',
    previousValue: formatHeatBags(snapshot.lastSeasonToDate),
    deltaValue: deltaLabel,
    deltaTrend,
    deltaColor: delta > 0 ? 'warning' as const : delta < 0 ? 'success' as const : 'muted' as const,
  }
}

export function createHeatDeleteDescription(refill: Refill) {
  return `Delete the ${formatHeatDate(refill.date)} refill for ${formatHeatBags(refill.bags)}? This removes it from the heating history.`
}
