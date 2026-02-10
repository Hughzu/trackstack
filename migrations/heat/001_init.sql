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
