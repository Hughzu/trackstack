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
  getTarget: (userId: string): CalorieTarget => {
    const db = getCaloriesDb();
    const row = db
      .prepare(
        "SELECT id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at FROM calorie_targets WHERE user_id = ?"
      )
      .get(userId) as CalorieTargetRow | undefined;

    if (row) return mapTargetRow(row);

    const now = new Date().toISOString();
    const id = randomUUID();
    db.prepare(
      "INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
    ).run(
      id,
      userId,
      DEFAULT_TARGET.targetKcal,
      DEFAULT_TARGET.targetProteinG,
      DEFAULT_TARGET.targetCarbsG ?? null,
      DEFAULT_TARGET.targetFatG ?? null,
      now,
      now
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

  updateTarget: (data: Omit<CalorieTarget, "id" | "createdAt" | "updatedAt">): CalorieTarget => {
    const db = getCaloriesDb();
    const existing = db
      .prepare(
        "SELECT id, created_at FROM calorie_targets WHERE user_id = ?"
      )
      .get(data.userId) as { id: string; created_at: string } | undefined;

    const now = new Date().toISOString();
    if (existing) {
      db.prepare(
        "UPDATE calorie_targets SET target_kcal = ?, target_protein_g = ?, target_carbs_g = ?, target_fat_g = ?, updated_at = ? WHERE user_id = ?"
      ).run(
        data.targetKcal,
        data.targetProteinG,
        data.targetCarbsG ?? null,
        data.targetFatG ?? null,
        now,
        data.userId
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
    db.prepare(
      "INSERT INTO calorie_targets (id, user_id, target_kcal, target_protein_g, target_carbs_g, target_fat_g, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
    ).run(
      id,
      data.userId,
      data.targetKcal,
      data.targetProteinG,
      data.targetCarbsG ?? null,
      data.targetFatG ?? null,
      now,
      now
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

  getTodaySummary: (userId: string) => {
    const db = getCaloriesDb();
    const { start, end } = getLocalDayRange();
    const row = db
      .prepare(
        "SELECT COALESCE(SUM(calories), 0) as calories, COALESCE(SUM(protein_g), 0) as protein, COALESCE(SUM(carbs_g), 0) as carbs, COALESCE(SUM(fat_g), 0) as fat FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ?"
      )
      .get(userId, start.toISOString(), end.toISOString()) as {
      calories: number;
      protein: number;
      carbs: number;
      fat: number;
    };

    const target = caloriesService.getTarget(userId);
    return {
      consumed: row.calories ?? 0,
      protein: row.protein ?? 0,
      carbs: row.carbs ?? 0,
      fat: row.fat ?? 0,
      target
    };
  },

  getTodayLogs: (userId: string) => {
    const db = getCaloriesDb();
    const { start, end } = getLocalDayRange();
    const rows = db
      .prepare(
        "SELECT id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title FROM calorie_logs WHERE user_id = ? AND date_time >= ? AND date_time < ? ORDER BY date_time DESC"
      )
      .all(userId, start.toISOString(), end.toISOString()) as CalorieLogRow[];

    return rows.map(mapLogRow);
  },

  addLog: (data: Omit<CalorieLog, "id" | "dateTime"> & { date: string; time?: string }) => {
    const db = getCaloriesDb();
    const id = randomUUID();
    const dateTimeIso = buildDateTimeIso(data.date, data.time);

    db.prepare(
      "INSERT INTO calorie_logs (id, user_id, date_time, calories, protein_g, carbs_g, fat_g, title) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
    ).run(
      id,
      data.userId,
      dateTimeIso,
      data.calories,
      data.proteinG,
      data.carbsG ?? null,
      data.fatG ?? null,
      data.title ?? null
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

  deleteLog: (id: string, userId: string): boolean => {
    const db = getCaloriesDb();
    const result = db
      .prepare("DELETE FROM calorie_logs WHERE id = ? AND user_id = ?")
      .run(id, userId);

    return result.changes > 0;
  }
};
