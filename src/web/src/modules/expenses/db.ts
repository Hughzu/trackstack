import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "expenses";

const schema = `
CREATE TABLE IF NOT EXISTS expense_settings (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE,
  income REAL NOT NULL,
  ratio_fund INTEGER NOT NULL,
  ratio_fun INTEGER NOT NULL,
  ratio_future INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS expense_checklist_templates (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  title TEXT NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL DEFAULT 'fund' CHECK (category IN ('fund', 'fun', 'future')),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_expense_checklist_templates_user ON expense_checklist_templates(user_id);

CREATE TABLE IF NOT EXISTS expense_recurring_templates (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  title TEXT NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL DEFAULT 'fund' CHECK (category IN ('fund', 'fun', 'future')),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_expense_recurring_templates_user ON expense_recurring_templates(user_id);

CREATE TABLE IF NOT EXISTS expense_sheets (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  period_key TEXT NOT NULL,
  created_at TEXT NOT NULL,
  closed_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_expense_sheets_user_open ON expense_sheets(user_id, closed_at);

CREATE TABLE IF NOT EXISTS expense_checklist_items (
  id TEXT PRIMARY KEY,
  sheet_id TEXT NOT NULL,
  template_id TEXT,
  title TEXT NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('fund', 'fun', 'future')),
  created_at TEXT NOT NULL,
  completed_at TEXT,
  expense_id TEXT
);
CREATE INDEX IF NOT EXISTS idx_expense_checklist_items_sheet ON expense_checklist_items(sheet_id, completed_at);

CREATE TABLE IF NOT EXISTS expense_entries (
  id TEXT PRIMARY KEY,
  sheet_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  title TEXT NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('fund', 'fun', 'future')),
  date TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('manual', 'recurring', 'checklist')),
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_expense_entries_sheet_date ON expense_entries(sheet_id, date);
`;

let initialized = false;

export const getExpensesDb = () => {
  const db = getDb(DOMAIN);
  if (!initialized) {
    db.exec(schema);
    initialized = true;
  }

  return db;
};
