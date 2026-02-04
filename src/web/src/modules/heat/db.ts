import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "heat";

const schema = `
CREATE TABLE IF NOT EXISTS refills (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  date TEXT NOT NULL,
  weight_kg REAL NOT NULL,
  bags INTEGER NOT NULL,
  temperature REAL,
  season TEXT
);
CREATE INDEX IF NOT EXISTS idx_refills_user_date ON refills(user_id, date);
`;

const getSeasonLabel = (dateIso: string) => {
  const date = new Date(dateIso);
  const startYear = date.getMonth() >= 8 ? date.getFullYear() : date.getFullYear() - 1;
  return `${startYear}-${startYear + 1}`;
};

let initialized = false;

export const getHeatDb = () => {
  const db = getDb(DOMAIN);
  if (!initialized) {
    db.exec(schema);
    const columns = db.prepare("PRAGMA table_info(refills)").all() as Array<{ name: string }>;
    const hasSeason = columns.some(column => column.name === "season");
    if (!hasSeason) {
      db.exec("ALTER TABLE refills ADD COLUMN season TEXT");
      const rows = db
        .prepare("SELECT id, date FROM refills WHERE season IS NULL")
        .all() as Array<{ id: string; date: string }>;
      const update = db.prepare("UPDATE refills SET season = ? WHERE id = ?");
      rows.forEach(row => {
        update.run(getSeasonLabel(row.date), row.id);
      });
    }
    initialized = true;
  }

  return db;
};
