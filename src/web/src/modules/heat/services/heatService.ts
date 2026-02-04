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
};

const mapRowToRefill = (row: RefillRow): Refill => ({
  id: row.id,
  userId: row.user_id,
  date: new Date(row.date),
  weightKg: row.weight_kg,
  bags: row.bags,
  temperature: row.temperature ?? undefined,
});

export const heatService = {
  /**
   * Calculates days elapsed since the most recent refill.
   */
  getDaysSinceLastRefill: (userId: string): number => {
    const db = getHeatDb();
    const row = db
      .prepare("SELECT date FROM refills WHERE user_id = ? ORDER BY date DESC LIMIT 1")
      .get(userId) as { date?: string } | undefined;

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
  getMonthlyConsumption: (year: number, userId: string): number[] => {
    const monthlyData = new Array(12).fill(0);
    const db = getHeatDb();
    const start = new Date(Date.UTC(year, 0, 1));
    const end = new Date(Date.UTC(year + 1, 0, 1));

    const rows = db
      .prepare("SELECT date, bags FROM refills WHERE user_id = ? AND date >= ? AND date < ?")
      .all(userId, start.toISOString(), end.toISOString()) as Array<{ date: string; bags: number }>;

    rows.forEach(row => {
      const date = new Date(row.date);
      monthlyData[date.getMonth()] += row.bags;
    });

    return monthlyData;
  },

  /**
   * Retrieves paginated history sorted by date descending.
   */
  getHistory: (page: number = 1, limit: number = 10, userId: string) => {
    const db = getHeatDb();
    const offset = (page - 1) * limit;
    const rows = db
      .prepare(
        "SELECT id, user_id, date, weight_kg, bags, temperature FROM refills WHERE user_id = ? ORDER BY date DESC LIMIT ? OFFSET ?"
      )
      .all(userId, limit, offset) as RefillRow[];

    const totalRow = db
      .prepare("SELECT COUNT(*) as count FROM refills WHERE user_id = ?")
      .get(userId) as { count: number };

    const total = totalRow?.count ?? 0;

    return {
      data: rows.map(mapRowToRefill),
      total,
      page,
      limit,
      totalPages: Math.ceil(total / limit)
    };
  },

  /**
   * Adds a new refill entry.
   */
  addRefill: (data: Omit<Refill, "id">): Refill => {
    const db = getHeatDb();
    const id = randomUUID();
    const dateIso = data.date.toISOString();

    db.prepare(
      "INSERT INTO refills (id, user_id, date, weight_kg, bags, temperature) VALUES (?, ?, ?, ?, ?, ?)"
    ).run(
      id,
      data.userId,
      dateIso,
      data.weightKg,
      data.bags,
      data.temperature ?? null
    );

    return {
      ...data,
      id
    };
  },

  deleteRefill: (id: string, userId: string): boolean => {
    const db = getHeatDb();
    const result = db
      .prepare("DELETE FROM refills WHERE id = ? AND user_id = ?")
      .run(id, userId);

    return result.changes > 0;
  }
};
