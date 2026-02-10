CREATE TABLE IF NOT EXISTS calorie_logs (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  date_time TEXT NOT NULL,
  calories INTEGER NOT NULL,
  protein_g INTEGER NOT NULL,
  carbs_g INTEGER,
  fat_g INTEGER,
  title TEXT
);
CREATE INDEX IF NOT EXISTS idx_calorie_logs_user_date ON calorie_logs(user_id, date_time);

CREATE TABLE IF NOT EXISTS calorie_targets (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE,
  target_kcal INTEGER NOT NULL,
  target_protein_g INTEGER NOT NULL,
  target_carbs_g INTEGER,
  target_fat_g INTEGER,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
