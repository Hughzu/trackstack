import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "heat";

const schema = `
CREATE TABLE IF NOT EXISTS refills (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  date TEXT NOT NULL,
  weight_kg REAL NOT NULL,
  bags INTEGER NOT NULL,
  temperature REAL
);
CREATE INDEX IF NOT EXISTS idx_refills_user_date ON refills(user_id, date);
`;

let initialized = false;

export const getHeatDb = () => {
  const db = getDb(DOMAIN);
  if (!initialized) {
    db.exec(schema);
    initialized = true;
  }

  return db;
};
