export type HeatHistoryItem = {
  id: string
  date: string
  bags: number
  weightKg: number
  temperature: number | null
  season: string
}

export type HeatOverviewMock = {
  seasonLabel: string
  daysSinceRefill: number
  totals: {
    thisSeason: number
    lastSeasonToDate: number
    lastSeasonTotal: number
  }
  history: HeatHistoryItem[]
}

type CreateHeatMockRefillInput = {
  date: string
  bags: number
  weightKg: number
  temperature: number | null
}

const state: HeatOverviewMock = {
  seasonLabel: '2025/26',
  daysSinceRefill: 18,
  totals: {
    thisSeason: 132,
    lastSeasonToDate: 118,
    lastSeasonTotal: 164,
  },
  history: [
    {
      id: 'refill-1',
      date: '2026-03-21',
      bags: 2,
      weightKg: 30,
      temperature: 8,
      season: '2025/26',
    },
    {
      id: 'refill-2',
      date: '2026-03-03',
      bags: 3,
      weightKg: 45,
      temperature: 6,
      season: '2025/26',
    },
    {
      id: 'refill-3',
      date: '2026-02-14',
      bags: 2,
      weightKg: 30,
      temperature: 4,
      season: '2025/26',
    },
    {
      id: 'refill-4',
      date: '2026-01-29',
      bags: 4,
      weightKg: 60,
      temperature: 2,
      season: '2025/26',
    },
  ],
}

const cloneHistoryItem = (item: HeatHistoryItem): HeatHistoryItem => ({ ...item })

export const readHeatMockState = (): HeatOverviewMock => ({
  seasonLabel: state.seasonLabel,
  daysSinceRefill: state.daysSinceRefill,
  totals: { ...state.totals },
  history: state.history.map(cloneHistoryItem),
})

export const createHeatMockRefill = (input: CreateHeatMockRefillInput) => {
  state.history.unshift({
    id: `mock-${Date.now()}`,
    date: input.date,
    bags: input.bags,
    weightKg: input.weightKg,
    temperature: input.temperature,
    season: state.seasonLabel,
  })

  state.daysSinceRefill = 0
  state.totals.thisSeason += input.bags
}
