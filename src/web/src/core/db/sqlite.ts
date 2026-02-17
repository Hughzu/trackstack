import { createClient, type Client, type ResultSet } from "@libsql/client/web";

type StatementArgs = Array<string | number | boolean | null>;

export type DbClient = {
  get: <T>(sql: string, args?: StatementArgs) => Promise<T | undefined>;
  all: <T>(sql: string, args?: StatementArgs) => Promise<T[]>;
  run: (sql: string, args?: StatementArgs) => Promise<number>;
  execute: (sql: string, args?: StatementArgs) => Promise<ResultSet>;
};

const dbCache = new Map<string, DbClient>();

const readEnv = (key: string) => {
  const metaEnv = (import.meta as { env?: Record<string, string | undefined> }).env;
  return metaEnv?.[key] ?? process.env[key];
};

const resolveDomainUrl = (domain: string) => {
  const upper = domain.toUpperCase();
  const tursoUrl = readEnv(`TURSO_${upper}_URL`);
  if (tursoUrl && tursoUrl.trim().length > 0) return tursoUrl.trim();

  throw new Error(`Missing TURSO_${upper}_URL for ${domain} database`);
};

const resolveDomainToken = (domain: string) => {
  const upper = domain.toUpperCase();
  const token = readEnv(`TURSO_${upper}_TOKEN`);
  return token && token.trim().length > 0 ? token.trim() : undefined;
};

const createDbClient = (domain: string): DbClient => {
  const url = resolveDomainUrl(domain);
  const authToken = url.startsWith("libsql://") ? resolveDomainToken(domain) : undefined;

  if (url.startsWith("libsql://") && !authToken) {
    throw new Error(`Missing TURSO_${domain.toUpperCase()}_TOKEN for ${domain} database`);
  }

  logInfo("db_connect", { domain, url: url.replace(authToken || "", "***") });

  const client: Client = createClient({ url, authToken });

  return {
    execute: async (sql: string, args?: StatementArgs) => {
      const start = performance.now();
      try {
        const result = await client.execute({ sql, args: args ?? [] });
        logDebug("db_query", { domain, sql, durationMs: performance.now() - start });
        return result;
      } catch (err) {
        logError("db_query_error", { domain, sql, error: err instanceof Error ? err.message : String(err) });
        throw err;
      }
    },
    get: async <T>(sql: string, args?: StatementArgs) => {
      const start = performance.now();
      try {
        const result = await client.execute({ sql, args: args ?? [] });
        logDebug("db_query_get", { domain, sql, durationMs: performance.now() - start });
        return (result.rows?.[0] as T | undefined) ?? undefined;
      } catch (err) {
        logError("db_query_error", { domain, sql, error: err instanceof Error ? err.message : String(err) });
        throw err;
      }
    },
    all: async <T>(sql: string, args?: StatementArgs) => {
      const start = performance.now();
      try {
        const result = await client.execute({ sql, args: args ?? [] });
        logDebug("db_query_all", { domain, sql, rows: result.rows.length, durationMs: performance.now() - start });
        return (result.rows as T[]) ?? [];
      } catch (err) {
        logError("db_query_error", { domain, sql, error: err instanceof Error ? err.message : String(err) });
        throw err;
      }
    },
    run: async (sql: string, args?: StatementArgs) => {
      const start = performance.now();
      try {
        const result = await client.execute({ sql, args: args ?? [] });
        logDebug("db_query_run", { domain, sql, affected: result.rowsAffected, durationMs: performance.now() - start });
        return result.rowsAffected ?? 0;
      } catch (err) {
        logError("db_query_error", { domain, sql, error: err instanceof Error ? err.message : String(err) });
        throw err;
      }
    }
  };
};

// Simple Structured Logger
const logInfo = (event: string, data: Record<string, any>) => {
  console.info(JSON.stringify({ level: "info", event, ...data, timestamp: new Date().toISOString() }));
};

const logError = (event: string, data: Record<string, any>) => {
  console.error(JSON.stringify({ level: "error", event, ...data, timestamp: new Date().toISOString() }));
};

const logDebug = (event: string, data: Record<string, any>) => {
  if (process.env.NODE_ENV === "development") {
    // console.debug(JSON.stringify({ level: "debug", event, ...data, timestamp: new Date().toISOString() }));
  }
};

export const getDb = (domain: string): DbClient => {
  const existing = dbCache.get(domain);
  if (existing) return existing;

  const db = createDbClient(domain);
  dbCache.set(domain, db);
  return db;
};
