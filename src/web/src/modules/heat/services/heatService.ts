import type { Refill } from '../types';

// Initial Mock Data
let refills: Refill[] = [
  { id: '1', date: new Date('2026-10-28'), weightKg: 15, bags: 1, temperature: 12 },
  { id: '2', date: new Date('2026-10-25'), weightKg: 15, bags: 1, temperature: 10 },
  { id: '3', date: new Date('2026-10-21'), weightKg: 30, bags: 2, temperature: 8 },
  { id: '4', date: new Date('2026-10-15'), weightKg: 15, bags: 1, temperature: 14 },
  { id: '5', date: new Date('2025-01-10'), weightKg: 45, bags: 3, temperature: 2 },
  { id: '6', date: new Date('2025-12-15'), weightKg: 60, bags: 4, temperature: 0 },
  { id: '7', date: new Date('2025-11-20'), weightKg: 30, bags: 2, temperature: 5 },
];

export const heatService = {
  /**
   * Calculates days elapsed since the most recent refill.
   */
  getDaysSinceLastRefill: (): number => {
    if (refills.length === 0) return 0;
    
    // Sort by date descending
    const sorted = [...refills].sort((a, b) => b.date.getTime() - a.date.getTime());
    const lastRefill = sorted[0];
    
    const now = new Date();
    // Reset time components for accurate day calculation
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const lastDate = new Date(lastRefill.date.getFullYear(), lastRefill.date.getMonth(), lastRefill.date.getDate());
    
    const diffTime = today.getTime() - lastDate.getTime();
    return Math.floor(diffTime / (1000 * 60 * 60 * 24));
  },

  /**
   * Returns an array of 12 numbers representing bags used per month (Jan-Dec) for a given year.
   */
  getMonthlyConsumption: (year: number): number[] => {
    const monthlyData = new Array(12).fill(0);
    
    refills.forEach(r => {
      if (r.date.getFullYear() === year) {
        monthlyData[r.date.getMonth()] += r.bags;
      }
    });
    
    return monthlyData;
  },

  /**
   * Retrieves paginated history sorted by date descending.
   */
  getHistory: (page: number = 1, limit: number = 10) => {
    const sorted = [...refills].sort((a, b) => b.date.getTime() - a.date.getTime());
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;
    
    return {
      data: sorted.slice(startIndex, endIndex),
      total: refills.length,
      page,
      limit,
      totalPages: Math.ceil(refills.length / limit)
    };
  },

  /**
   * Adds a new refill entry.
   */
  addRefill: (data: Omit<Refill, 'id'>): Refill => {
    const newRefill: Refill = {
      ...data,
      id: Math.random().toString(36).substring(2, 9)
    };
    refills.push(newRefill);
    return newRefill;
  }
};
