import fs from "node:fs";
import path from "node:path";
import Database from "better-sqlite3";

const dbCache = new Map<string, Database.Database>();

const resolveDataDir = () => {
  const envDir = process.env.DATA_DIR;
  if (envDir && envDir.trim().length > 0) {
    return path.isAbsolute(envDir) ? envDir : path.resolve(process.cwd(), envDir);
  }

  return path.resolve(process.cwd(), "..", "data");
};

export const getDb = (domain: string): Database.Database => {
  const existing = dbCache.get(domain);
  if (existing) return existing;

  const dataDir = resolveDataDir();
  fs.mkdirSync(dataDir, { recursive: true });

  const dbPath = path.join(dataDir, `${domain}.sqlite`);
  const db = new Database(dbPath);
  db.pragma("journal_mode = WAL");

  dbCache.set(domain, db);
  return db;
};
