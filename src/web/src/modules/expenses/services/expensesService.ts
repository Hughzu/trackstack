import { randomUUID } from "node:crypto";
import { getExpensesDb } from "../db";
import type {
  ChecklistItem,
  ExpenseCategory,
  ExpenseEntry,
  ExpenseSettings,
  ExpenseSheet,
  ExpenseTemplate,
  ExpenseType
} from "../types";

type SettingsRow = {
  id: string;
  user_id: string;
  income: number;
  ratio_fund: number;
  ratio_fun: number;
  ratio_future: number;
  created_at: string;
  updated_at: string;
};

type TemplateRow = {
  id: string;
  user_id: string;
  title: string;
  amount: number;
  category: ExpenseCategory;
  created_at: string;
  updated_at: string;
};

type SheetRow = {
  id: string;
  user_id: string;
  period_key: string;
  created_at: string;
  closed_at: string | null;
};

type ChecklistRow = {
  id: string;
  sheet_id: string;
  template_id: string | null;
  title: string;
  amount: number;
  category: ExpenseCategory;
  created_at: string;
  completed_at: string | null;
  expense_id: string | null;
};

type EntryRow = {
  id: string;
  sheet_id: string;
  user_id: string;
  title: string;
  amount: number;
  category: ExpenseCategory;
  date: string;
  type: ExpenseType;
  created_at: string;
};

const DEFAULT_SETTINGS = {
  income: 2215,
  ratioFund: 60,
  ratioFun: 20,
  ratioFuture: 20
};

const formatPeriodKey = (date: Date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  return `${year}-${month}`;
};

const addPeriodKeyMonth = (periodKey: string) => {
  const [yearPart, monthPart] = periodKey.split("-");
  const year = Number(yearPart);
  const month = Number(monthPart);
  const base = new Date(year, month - 1, 1);
  base.setMonth(base.getMonth() + 1);
  return formatPeriodKey(base);
};

const getTodayDate = () => new Date().toISOString().split("T")[0];

const normalizeCategory = (value?: string): ExpenseCategory => {
  if (value === "fun" || value === "future" || value === "fund") return value;
  return "fund";
};

const mapSettingsRow = (row: SettingsRow): ExpenseSettings => ({
  id: row.id,
  userId: row.user_id,
  income: row.income,
  ratioFund: row.ratio_fund,
  ratioFun: row.ratio_fun,
  ratioFuture: row.ratio_future,
  createdAt: new Date(row.created_at),
  updatedAt: new Date(row.updated_at)
});

const mapTemplateRow = (row: TemplateRow): ExpenseTemplate => ({
  id: row.id,
  userId: row.user_id,
  title: row.title,
  amount: row.amount,
  category: row.category,
  createdAt: new Date(row.created_at),
  updatedAt: new Date(row.updated_at)
});

const mapSheetRow = (row: SheetRow): ExpenseSheet => ({
  id: row.id,
  userId: row.user_id,
  periodKey: row.period_key,
  createdAt: new Date(row.created_at),
  closedAt: row.closed_at ? new Date(row.closed_at) : undefined
});

const mapChecklistRow = (row: ChecklistRow): ChecklistItem => ({
  id: row.id,
  sheetId: row.sheet_id,
  templateId: row.template_id ?? undefined,
  title: row.title,
  amount: row.amount,
  category: row.category,
  createdAt: new Date(row.created_at),
  completedAt: row.completed_at ? new Date(row.completed_at) : undefined,
  expenseId: row.expense_id ?? undefined
});

const mapEntryRow = (row: EntryRow): ExpenseEntry => ({
  id: row.id,
  sheetId: row.sheet_id,
  userId: row.user_id,
  title: row.title,
  amount: row.amount,
  category: row.category,
  date: row.date,
  type: row.type,
  createdAt: new Date(row.created_at)
});

const createSheet = async (userId: string, periodKey: string): Promise<ExpenseSheet> => {
  const db = getExpensesDb();
  const id = randomUUID();
  const now = new Date().toISOString();

  await db.run(
    "INSERT INTO expense_sheets (id, user_id, period_key, created_at, closed_at) VALUES (?, ?, ?, ?, NULL)",
    [id, userId, periodKey, now]
  );

  const templates = await db.all<TemplateRow>(
    "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_checklist_templates WHERE user_id = ? ORDER BY created_at ASC",
    [userId]
  );

  const recurring = await db.all<TemplateRow>(
    "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_recurring_templates WHERE user_id = ? ORDER BY created_at ASC",
    [userId]
  );

  await Promise.all(
    templates.map(template =>
      db.run(
        "INSERT INTO expense_checklist_items (id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id) VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)",
        [
          randomUUID(),
          id,
          template.id,
          template.title,
          template.amount,
          template.category,
          now
        ]
      )
    )
  );

  const recurringDate = `${periodKey}-01`;
  await Promise.all(
    recurring.map(template =>
      db.run(
        "INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
        [
          randomUUID(),
          id,
          userId,
          template.title,
          template.amount,
          template.category,
          recurringDate,
          "recurring",
          now
        ]
      )
    )
  );

  return {
    id,
    userId,
    periodKey,
    createdAt: new Date(now)
  };
};

const getOpenSheet = async (userId: string) => {
  const db = getExpensesDb();
  const row = await db.get<SheetRow>(
    "SELECT id, user_id, period_key, created_at, closed_at FROM expense_sheets WHERE user_id = ? AND closed_at IS NULL ORDER BY created_at DESC LIMIT 1",
    [userId]
  );

  return row ? mapSheetRow(row) : undefined;
};

const getLatestSheet = async (userId: string) => {
  const db = getExpensesDb();
  const row = await db.get<SheetRow>(
    "SELECT id, user_id, period_key, created_at, closed_at FROM expense_sheets WHERE user_id = ? ORDER BY created_at DESC LIMIT 1",
    [userId]
  );

  return row ? mapSheetRow(row) : undefined;
};

export const expensesService = {
  getSettings: async (userId: string): Promise<ExpenseSettings> => {
    const db = getExpensesDb();
    const row = await db.get<SettingsRow>(
      "SELECT id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at FROM expense_settings WHERE user_id = ?",
      [userId]
    );

    if (row) return mapSettingsRow(row);

    const now = new Date().toISOString();
    const id = randomUUID();
    await db.run(
      "INSERT INTO expense_settings (id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        userId,
        DEFAULT_SETTINGS.income,
        DEFAULT_SETTINGS.ratioFund,
        DEFAULT_SETTINGS.ratioFun,
        DEFAULT_SETTINGS.ratioFuture,
        now,
        now
      ]
    );

    return {
      id,
      userId,
      income: DEFAULT_SETTINGS.income,
      ratioFund: DEFAULT_SETTINGS.ratioFund,
      ratioFun: DEFAULT_SETTINGS.ratioFun,
      ratioFuture: DEFAULT_SETTINGS.ratioFuture,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  updateSettings: async (
    data: Omit<ExpenseSettings, "id" | "createdAt" | "updatedAt">
  ): Promise<ExpenseSettings> => {
    const db = getExpensesDb();
    const existing = await db.get<{ id: string; created_at: string }>(
      "SELECT id, created_at FROM expense_settings WHERE user_id = ?",
      [data.userId]
    );

    const now = new Date().toISOString();
    if (existing) {
      await db.run(
        "UPDATE expense_settings SET income = ?, ratio_fund = ?, ratio_fun = ?, ratio_future = ?, updated_at = ? WHERE user_id = ?",
        [
          data.income,
          data.ratioFund,
          data.ratioFun,
          data.ratioFuture,
          now,
          data.userId
        ]
      );

      return {
        id: existing.id,
        userId: data.userId,
        income: data.income,
        ratioFund: data.ratioFund,
        ratioFun: data.ratioFun,
        ratioFuture: data.ratioFuture,
        createdAt: new Date(existing.created_at),
        updatedAt: new Date(now)
      };
    }

    const id = randomUUID();
    await db.run(
      "INSERT INTO expense_settings (id, user_id, income, ratio_fund, ratio_fun, ratio_future, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
      [
        id,
        data.userId,
        data.income,
        data.ratioFund,
        data.ratioFun,
        data.ratioFuture,
        now,
        now
      ]
    );

    return {
      id,
      userId: data.userId,
      income: data.income,
      ratioFund: data.ratioFund,
      ratioFun: data.ratioFun,
      ratioFuture: data.ratioFuture,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  getChecklistTemplates: async (userId: string): Promise<ExpenseTemplate[]> => {
    const db = getExpensesDb();
    const rows = await db.all<TemplateRow>(
      "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_checklist_templates WHERE user_id = ? ORDER BY created_at ASC",
      [userId]
    );

    return rows.map(mapTemplateRow);
  },

  getRecurringTemplates: async (userId: string): Promise<ExpenseTemplate[]> => {
    const db = getExpensesDb();
    const rows = await db.all<TemplateRow>(
      "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_recurring_templates WHERE user_id = ? ORDER BY created_at ASC",
      [userId]
    );

    return rows.map(mapTemplateRow);
  },

  upsertChecklistTemplate: async (data: {
    id?: string;
    userId: string;
    title: string;
    amount: number;
    category?: string;
  }): Promise<ExpenseTemplate> => {
    const db = getExpensesDb();
    const now = new Date().toISOString();
    const category = normalizeCategory(data.category);

    if (data.id) {
      const existing = await db.get<{ id: string }>(
        "SELECT id FROM expense_checklist_templates WHERE id = ? AND user_id = ?",
        [data.id, data.userId]
      );

      if (existing) {
        const updatedChanges = await db.run(
          "UPDATE expense_checklist_templates SET title = ?, amount = ?, category = ?, updated_at = ? WHERE id = ? AND user_id = ?",
          [data.title, data.amount, category, now, data.id, data.userId]
        );

        if (updatedChanges > 0) {
          await db.run(
            "UPDATE expense_checklist_items SET title = ?, amount = ?, category = ? WHERE template_id = ? AND completed_at IS NULL",
            [data.title, data.amount, category, data.id]
          );
        }

        const row = await db.get<TemplateRow>(
          "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_checklist_templates WHERE id = ?",
          [data.id]
        );

        if (!row) {
          throw new Error("Checklist template not found after update");
        }

        return mapTemplateRow(row);
      }
    }

    const id = randomUUID();
    await db.run(
      "INSERT INTO expense_checklist_templates (id, user_id, title, amount, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
      [id, data.userId, data.title, data.amount, category, now, now]
    );

    const sheet = await expensesService.getOrCreateOpenSheet(data.userId);
    await db.run(
      "INSERT INTO expense_checklist_items (id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id) VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)",
      [randomUUID(), sheet.id, id, data.title, data.amount, category, now]
    );

    return {
      id,
      userId: data.userId,
      title: data.title,
      amount: data.amount,
      category,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  deleteChecklistTemplate: async (id: string, userId: string): Promise<boolean> => {
    const db = getExpensesDb();
    await db.run(
      "DELETE FROM expense_checklist_items WHERE template_id = ? AND completed_at IS NULL AND sheet_id IN (SELECT id FROM expense_sheets WHERE user_id = ? AND closed_at IS NULL)",
      [id, userId]
    );

    const changes = await db.run(
      "DELETE FROM expense_checklist_templates WHERE id = ? AND user_id = ?",
      [id, userId]
    );

    return changes > 0;
  },

  upsertRecurringTemplate: async (data: {
    id?: string;
    userId: string;
    title: string;
    amount: number;
    category?: string;
  }): Promise<ExpenseTemplate> => {
    const db = getExpensesDb();
    const now = new Date().toISOString();
    const category = normalizeCategory(data.category);

    if (data.id) {
      const existing = await db.get<{ id: string }>(
        "SELECT id FROM expense_recurring_templates WHERE id = ? AND user_id = ?",
        [data.id, data.userId]
      );

      if (existing) {
        await db.run(
          "UPDATE expense_recurring_templates SET title = ?, amount = ?, category = ?, updated_at = ? WHERE id = ? AND user_id = ?",
          [data.title, data.amount, category, now, data.id, data.userId]
        );

        const row = await db.get<TemplateRow>(
          "SELECT id, user_id, title, amount, category, created_at, updated_at FROM expense_recurring_templates WHERE id = ?",
          [data.id]
        );

        if (!row) {
          throw new Error("Recurring template not found after update");
        }

        return mapTemplateRow(row);
      }
    }

    const id = randomUUID();
    await db.run(
      "INSERT INTO expense_recurring_templates (id, user_id, title, amount, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
      [id, data.userId, data.title, data.amount, category, now, now]
    );

    const sheet = await expensesService.getOrCreateOpenSheet(data.userId);
    const date = getTodayDate();
    await db.run(
      "INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
      [randomUUID(), sheet.id, data.userId, data.title, data.amount, category, date, "recurring", now]
    );

    return {
      id,
      userId: data.userId,
      title: data.title,
      amount: data.amount,
      category,
      createdAt: new Date(now),
      updatedAt: new Date(now)
    };
  },

  deleteRecurringTemplate: async (id: string, userId: string): Promise<boolean> => {
    const db = getExpensesDb();
    const changes = await db.run(
      "DELETE FROM expense_recurring_templates WHERE id = ? AND user_id = ?",
      [id, userId]
    );

    return changes > 0;
  },

  getOrCreateOpenSheet: async (userId: string): Promise<ExpenseSheet> => {
    const open = await getOpenSheet(userId);
    if (open) return open;

    const last = await getLatestSheet(userId);
    const periodKey = last ? addPeriodKeyMonth(last.periodKey) : formatPeriodKey(new Date());
    return createSheet(userId, periodKey);
  },

  closeSheet: async (userId: string): Promise<ExpenseSheet> => {
    const db = getExpensesDb();
    const now = new Date().toISOString();
    const open = await getOpenSheet(userId);
    let basePeriodKey = open?.periodKey;

    if (open) {
      await db.run("UPDATE expense_sheets SET closed_at = ? WHERE id = ?", [now, open.id]);
    }

    if (!basePeriodKey) {
      const last = await getLatestSheet(userId);
      basePeriodKey = last?.periodKey ?? formatPeriodKey(new Date());
    }

    const nextPeriodKey = addPeriodKeyMonth(basePeriodKey);
    return createSheet(userId, nextPeriodKey);
  },

  addExpense: async (data: {
    userId: string;
    title: string;
    amount: number;
    category?: string;
    date?: string;
  }): Promise<ExpenseEntry> => {
    const db = getExpensesDb();
    const sheet = await expensesService.getOrCreateOpenSheet(data.userId);
    const id = randomUUID();
    const now = new Date().toISOString();
    const date = data.date && data.date.trim().length > 0 ? data.date : getTodayDate();
    const category = normalizeCategory(data.category);

    await db.run(
      "INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
      [id, sheet.id, data.userId, data.title, data.amount, category, date, "manual", now]
    );

    return {
      id,
      sheetId: sheet.id,
      userId: data.userId,
      title: data.title,
      amount: data.amount,
      category,
      date,
      type: "manual",
      createdAt: new Date(now)
    };
  },

  deleteExpense: async (id: string, userId: string): Promise<boolean> => {
    const db = getExpensesDb();
    const changes = await db.run(
      "DELETE FROM expense_entries WHERE id = ? AND user_id = ?",
      [id, userId]
    );

    return changes > 0;
  },

  completeChecklistItem: async (data: { id: string; userId: string; date?: string }) => {
    const db = getExpensesDb();
    const row = await db.get<{
      id: string;
      sheet_id: string;
      title: string;
      amount: number;
      category: ExpenseCategory;
    }>(
      "SELECT i.id, i.sheet_id, i.title, i.amount, i.category FROM expense_checklist_items i JOIN expense_sheets s ON i.sheet_id = s.id WHERE i.id = ? AND s.user_id = ? AND i.completed_at IS NULL",
      [data.id, data.userId]
    );

    if (!row) return null;

    const id = randomUUID();
    const now = new Date().toISOString();
    const date = data.date && data.date.trim().length > 0 ? data.date : getTodayDate();

    await db.run(
      "INSERT INTO expense_entries (id, sheet_id, user_id, title, amount, category, date, type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
      [id, row.sheet_id, data.userId, row.title, row.amount, row.category, date, "checklist", now]
    );

    await db.run(
      "UPDATE expense_checklist_items SET completed_at = ?, expense_id = ? WHERE id = ?",
      [now, id, data.id]
    );

    return {
      id,
      sheetId: row.sheet_id,
      userId: data.userId,
      title: row.title,
      amount: row.amount,
      category: row.category,
      date,
      type: "checklist",
      createdAt: new Date(now)
    } as ExpenseEntry;
  },

  getDashboard: async (userId: string) => {
    const db = getExpensesDb();
    const settings = await expensesService.getSettings(userId);
    const sheet = await expensesService.getOrCreateOpenSheet(userId);

    const totalRow = await db.get<{ total: number }>(
      "SELECT COALESCE(SUM(amount), 0) as total FROM expense_entries WHERE sheet_id = ?",
      [sheet.id]
    );

    const spentRows = await db.all<{ category: ExpenseCategory; total: number }>(
      "SELECT category, COALESCE(SUM(amount), 0) as total FROM expense_entries WHERE sheet_id = ? GROUP BY category",
      [sheet.id]
    );

    const spentByCategory = spentRows.reduce(
      (acc, row) => {
        acc[row.category] = row.total;
        return acc;
      },
      { fund: 0, fun: 0, future: 0 } as Record<ExpenseCategory, number>
    );

    const pendingRows = await db.all<ChecklistRow>(
      "SELECT id, sheet_id, template_id, title, amount, category, created_at, completed_at, expense_id FROM expense_checklist_items WHERE sheet_id = ? AND completed_at IS NULL ORDER BY created_at ASC",
      [sheet.id]
    );

    const historyRows = await db.all<EntryRow>(
      "SELECT id, sheet_id, user_id, title, amount, category, date, type, created_at FROM expense_entries WHERE sheet_id = ? ORDER BY date DESC, created_at DESC",
      [sheet.id]
    );

    return {
      settings,
      sheet,
      totalSpent: totalRow?.total ?? 0,
      spentByCategory,
      pendingChecklist: pendingRows.map(mapChecklistRow),
      history: historyRows.map(mapEntryRow)
    };
  }
};
