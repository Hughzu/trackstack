import { randomUUID } from "node:crypto";
import type { Refill } from "../types";
import { getHeatDb } from "../db";

type RefillRow = {
  id: string;
  user_id: string;
  date: string;
  weight_kg: number;
  bags: number;
  temperature: number | null;
  season: string | null;
};

type SeasonSnapshot = {
  seasonLabel: string;
  seasonToDate: number;
  lastSeasonToDate: number;
  delta: number;
  deltaPct: number | null;
};

const mapRowToRefill = (row: RefillRow): Refill => ({
  id: row.id,
  userId: row.user_id,
  date: new Date(row.date),
  weightKg: row.weight_kg,
  bags: row.bags,
  temperature: row.temperature ?? undefined,
  season: row.season ?? undefined,
});

const getSeasonStartYear = (date: Date) => (date.getMonth() >= 8 ? date.getFullYear() : date.getFullYear() - 1);

const getSeasonRange = (startYear: number) => {
  const start = new Date(Date.UTC(startYear, 8, 1));
  const end = new Date(Date.UTC(startYear + 1, 8, 1));
  return {
    start,
    end,
    label: `${startYear}-${startYear + 1}`
  };
};

const getSeasonSum = async (userId: string, start: Date, end: Date) => {
  const db = getHeatDb();
  const row = await db.get<{ total: number }>(
    "SELECT COALESCE(SUM(bags), 0) as total FROM refills WHERE user_id = ? AND date >= ? AND date < ?",
    [userId, start.toISOString(), end.toISOString()]
  );

  return row?.total ?? 0;
};

export const heatService = {
  /**
   * Calculates days elapsed since the most recent refill.
   */
  getDaysSinceLastRefill: async (userId: string): Promise<number> => {
    const db = getHeatDb();
    const row = await db.get<{ date?: string }>(
      "SELECT date FROM refills WHERE user_id = ? ORDER BY date DESC LIMIT 1",
      [userId]
    );

    if (!row?.date) return 0;

    const lastRefillDate = new Date(row.date);
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const lastDate = new Date(
      lastRefillDate.getFullYear(),
      lastRefillDate.getMonth(),
      lastRefillDate.getDate()
    );

    const diffTime = today.getTime() - lastDate.getTime();
    return Math.floor(diffTime / (1000 * 60 * 60 * 24));
  },

  /**
   * Returns an array of 12 numbers representing bags used per month (Jan-Dec) for a given year.
   */
  getMonthlyConsumption: async (year: number, userId: string): Promise<number[]> => {
    const monthlyData = new Array(12).fill(0);
    const db = getHeatDb();
    const start = new Date(Date.UTC(year, 0, 1));
    const end = new Date(Date.UTC(year + 1, 0, 1));

    const rows = await db.all<{ date: string; bags: number }>(
      "SELECT date, bags FROM refills WHERE user_id = ? AND date >= ? AND date < ?",
      [userId, start.toISOString(), end.toISOString()]
    );

    rows.forEach(row => {
      const date = new Date(row.date);
      monthlyData[date.getMonth()] += row.bags;
    });

    return monthlyData;
  },

  /**
   * Retrieves paginated history sorted by date descending.
   */
  getHistory: async (page: number = 1, limit: number = 10, userId: string) => {
    const db = getHeatDb();
    const offset = (page - 1) * limit;
    const rows = await db.all<RefillRow>(
      "SELECT id, user_id, date, weight_kg, bags, temperature, season FROM refills WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?",
      [userId, limit, offset]
    );

    const totalRow = await db.get<{ count: number }>(
      "SELECT COUNT(*) as count FROM refills WHERE user_id = ?",
      [userId]
    );

    const total = totalRow?.count ?? 0;

    return {
      data: rows.map(mapRowToRefill),
      total,
      page,
      limit,
      totalPages: Math.ceil(total / limit)
    };
  },

  getSeasonSnapshot: async (
    userId: string,
    referenceDate: Date = new Date()
  ): Promise<SeasonSnapshot> => {
    const todayUtcStart = new Date(
      Date.UTC(referenceDate.getFullYear(), referenceDate.getMonth(), referenceDate.getDate())
    );
    const todayUtcEndExclusive = new Date(todayUtcStart.getTime() + 24 * 60 * 60 * 1000);
    const seasonStartYear = getSeasonStartYear(referenceDate);
    const currentSeason = getSeasonRange(seasonStartYear);
    const currentEnd =
      todayUtcEndExclusive.getTime() < currentSeason.end.getTime()
        ? todayUtcEndExclusive
        : currentSeason.end;

    const lastSeason = getSeasonRange(seasonStartYear - 1);
    const offsetMs = currentEnd.getTime() - currentSeason.start.getTime();
    const lastSeasonEndSamePeriod = new Date(lastSeason.start.getTime() + offsetMs);
    const lastEnd =
      lastSeasonEndSamePeriod.getTime() < lastSeason.end.getTime()
        ? lastSeasonEndSamePeriod
        : lastSeason.end;

    const seasonToDate = await getSeasonSum(userId, currentSeason.start, currentEnd);
    const lastSeasonToDate = await getSeasonSum(userId, lastSeason.start, lastEnd);
    const delta = seasonToDate - lastSeasonToDate;
    const deltaPct = lastSeasonToDate === 0 ? null : Math.round((delta / lastSeasonToDate) * 100);

    return {
      seasonLabel: currentSeason.label,
      seasonToDate,
      lastSeasonToDate,
      delta,
      deltaPct
    };
  },

  /**
   * Adds a new refill entry.
   */
  addRefill: async (data: Omit<Refill, "id">): Promise<Refill> => {
    const db = getHeatDb();
    const id = randomUUID();
    const dateIso = data.date.toISOString();
    const seasonStartYear = getSeasonStartYear(data.date);
    const seasonLabel = `${seasonStartYear}-${seasonStartYear + 1}`;

    await db.run(
      "INSERT INTO refills (id, user_id, date, weight_kg, bags, temperature, season) VALUES (?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        data.userId,
        dateIso,
        data.weightKg,
        data.bags,
        data.temperature ?? null,
        seasonLabel
      ]
    );

    return {
      ...data,
      id
    };
  },

  deleteRefill: async (id: string, userId: string): Promise<boolean> => {
    const db = getHeatDb();
    const changes = await db.run("DELETE FROM refills WHERE id = ? AND user_id = ?", [id, userId]);

    return changes > 0;
  }
};
