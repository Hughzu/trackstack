import { randomUUID } from "node:crypto";
import { getCaloriesDb } from "../db";

type CalorieLog = {
  id: string;
  userId: string;
  dateTime: Date;
  calories: number;
  proteinG: number;
  carbsG?: number;
  fatG?: number;
  title?: string;
};

type CalorieTarget = {
  id: string;
  userId: string;
  targetKcal: number;
  targetProteinG: number;
  targetCarbsG?: number;
  targetFatG?: number;
  createdAt: Date;
  updatedAt: Date;
};

type CalorieLogRow = {
  id: string;
  user_id: string;
  date_time: string;
  calories: number;
  protein_g: number;
  carbs_g: number | null;
  fat_g: number | null;
  title: string | null;
};

type CalorieTargetRow = {
  id: string;
  user_id: string;
  target_kcal: number;
  target_protein_g: number;
  target_carbs_g: number | null;
  target_fat_g: number | null;
  created_at: string;
  updated_at: string;
};

const DEFAULT_TARGET = {
  targetKcal: 2300,
  targetProteinG: 120,
  targetCarbsG: undefined,
  targetFatG: undefined
};

const mapLogRow = (row: CalorieLogRow): CalorieLog => ({
  id: row.id,
  userId: row.user_id,
  dateTime: new Date(row.date_time),
  calories: row.calories,
  proteinG: row.protein_g,
  carbsG: row.carbs_g ?? undefined,
  fatG: row.fat_g ?? undefined,
  title: row.title ?? undefined
});

const mapTargetRow = (row: CalorieTargetRow): CalorieTarget => ({
  id: row.id,
  userId: row.user_id,
  targetKcal: row.target_kcal,
  targetProteinG: row.target_protein_g,
  targetCarbsG: row.target_carbs_g ?? undefined,
  targetFatG: row.target_fat_g ?? undefined,
  createdAt: new Date(row.created_at),
  updatedAt: new Date(row.updated_at)
});

const getLocalDayRange = () => {
  const start = new Date();
  start.setHours(0, 0, 0, 0);
  const end = new Date(start);
  end.setDate(start.getDate() + 1);
  return { start, end };
};

const buildDateTimeIso = (date: string, time?: string) => {
  const now = new Date();
  const timeValue = time && time.trim().length > 0 ? time : now.toTimeString().slice(0, 5);
  const combined = new Date(`${date}T${timeValue}:00`);
  return combined.toISOString();
};

export const caloriesService = {
  getTarget: async (userId: string): Promise<CalorieTarget> => {
    const db = getCaloriesDb();
    const row = await db.get<CalorieTargetRow>(
      "SELECT id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at FROM calorie_targets WHERE user_id = ?",
      [userId]
    );

    if (row) return mapTargetRow(row);

    const now = new Date().toISOString();
    const id = randomUUID();
    await db.run(
      "INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        userId,
        DEFAULT_TARGET.targetKcal,
        DEFAULT_TARGET.targetProteinG,
        DEFAULT_TARGET.targetCarbsG ?? null,
        DEFAULT_TARGET.targetFatG ?? null,
        now,
        now
      ]
    );

    return {
      id,
      userId,
      targetKcal: DEFAULT_TARGET.targetKcal,
      targetProteinG: DEFAULT_TARGET.targetProteinG,
      targetCarbsG: DEFAULT_TARGET.targetCarbsG,
      targetFatG: DEFAULT_TARGET.targetFatG,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  updateTarget: async (
    data: Omit<CalorieTarget, "id" | "createdAt" | "updatedAt">
  ): Promise<CalorieTarget> => {
    const db = getCaloriesDb();
    const existing = await db.get<{ id: string; created_at: string }>(
      "SELECT id, created_at FROM calorie_targets WHERE user_id = ?",
      [data.userId]
    );

    const now = new Date().toISOString();
    if (existing) {
      await db.run(
        "UPDATE calorie_targets SET target_kcal = ?, target_protein_g = ?, target_carbs_g = ?, target_fat_g = ?, updated_at = ? WHERE user_id = ?",
        [
          data.targetKcal,
          data.targetProteinG,
          data.targetCarbsG ?? null,
          data.targetFatG ?? null,
          now,
          data.userId
        ]
      );

      return {
        id: existing.id,
        userId: data.userId,
        targetKcal: data.targetKcal,
        targetProteinG: data.targetProteinG,
        targetCarbsG: data.targetCarbsG,
        targetFatG: data.targetFatG,
        createdAt: new Date(existing.created_at),
        updatedAt: new Date(now)
      };
    }

    const id = randomUUID();
    await db.run(
      "INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        data.userId,
        data.targetKcal,
        data.targetProteinG,
        data.targetCarbsG ?? null,
        data.targetFatG ?? null,
        now,
        now
      ]
    );

    return {
      id,
      userId: data.userId,
      targetKcal: data.targetKcal,
      targetProteinG: data.targetProteinG,
      targetCarbsG: data.targetCarbsG,
      targetFatG: data.targetFatG,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  getTodaySummary: async (userId: string) => {
    const db = getCaloriesDb();
    const { start, end } = getLocalDayRange();
    const row = await db.get<{
      calories: number;
      protein: number;
      carbs: number;
      fat: number;
    }>(
      "SELECT COALESCE(SUM(calories), 0) as calories, COALESCE(SUM(protein_g), 0) as protein, COALESCE(SUM(carbs_g), 0) as carbs, COALESCE(SUM(fat_g), 0) as fat FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ?",
      [userId, start.toISOString(), end.toISOString()]
    );

    const target = await caloriesService.getTarget(userId);
    return {
      consumed: row?.calories ?? 0,
      protein: row?.protein ?? 0,
      carbs: row?.carbs ?? 0,
      fat: row?.fat ?? 0,
      target
    };
  },

  getTodayLogs: async (userId: string) => {
    const db = getCaloriesDb();
    const { start, end } = getLocalDayRange();
    const rows = await db.all<CalorieLogRow>(
      "SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ? ORDER BY date_time DESC",
      [userId, start.toISOString(), end.toISOString()]
    );

    return rows.map(mapLogRow);
  },

  getRecentLogs: async (userId: string, limit = 8) => {
    const db = getCaloriesDb();
    const rows = await db.all<CalorieLogRow>(
      "WITH ranked AS (SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title, ROW_NUMBER() OVER (PARTITION BY LOWER(TRIM(title)) ORDER BY date_time DESC) AS rn FROM calorie_logs WHERE user_id = ? AND title IS NOT NULL AND TRIM(title) != '') SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM ranked WHERE rn = 1 ORDER BY date_time DESC LIMIT ?",
      [userId, limit]
    );

    return rows.map(mapLogRow);
  },

  addLog: async (data: Omit<CalorieLog, "id" | "dateTime"> & { date: string; time?: string }) => {
    const db = getCaloriesDb();
    const id = randomUUID();
    const dateTimeIso = buildDateTimeIso(data.date, data.time);

    await db.run(
      "INSERT INTO calorie_logs (id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        data.userId,
        dateTimeIso,
        data.calories,
        data.proteinG,
        data.carbsG ?? null,
        data.fatG ?? null,
        data.title ?? null
      ]
    );

    return {
      id,
      userId: data.userId,
      dateTime: new Date(dateTimeIso),
      calories: data.calories,
      proteinG: data.proteinG,
      carbsG: data.carbsG,
      fatG: data.fatG,
      title: data.title
    };
  },

  deleteLog: async (id: string, userId: string): Promise<boolean> => {
    const db = getCaloriesDb();
    const changes = await db.run(
      "DELETE FROM calorie_logs WHERE id = ? AND user_id = ?",
      [id, userId]
    );

    return changes > 0;
  }
};
