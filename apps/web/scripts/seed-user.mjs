import crypto from "node:crypto";
import { promisify } from "node:util";
import { createClient } from "@libsql/client";

const scryptAsync = promisify(crypto.scrypt);

const DEFAULT_PARAMS = {
  N: 16384,
  r: 8,
  p: 1,
  keyLength: 64,
  saltLength: 16
};

const readArg = (name) => {
  const index = process.argv.indexOf(`--${name}`);
  if (index === -1) return undefined;
  return process.argv[index + 1];
};

const hashPassword = async (password) => {
  const salt = crypto.randomBytes(DEFAULT_PARAMS.saltLength);
  const derivedKey = await scryptAsync(password, salt, DEFAULT_PARAMS.keyLength, {
    cost: DEFAULT_PARAMS.N,
    blockSize: DEFAULT_PARAMS.r,
    parallelization: DEFAULT_PARAMS.p
  });

  const params = `N=${DEFAULT_PARAMS.N},r=${DEFAULT_PARAMS.r},p=${DEFAULT_PARAMS.p}`;
  const saltEncoded = salt.toString("base64");
  const hashEncoded = Buffer.from(derivedKey).toString("base64");
  return `$scrypt$${params}$${saltEncoded}$${hashEncoded}`;
};

const resolveDbUrl = () => {
  const envUrl = process.env.TURSO_USERS_URL;
  if (envUrl && envUrl.trim().length > 0) return envUrl.trim();
  throw new Error("Missing TURSO_USERS_URL. Seed the same users database configured for the app.");
};

const resolveDbToken = () => {
  const token = process.env.TURSO_USERS_TOKEN;
  return token && token.trim().length > 0 ? token.trim() : undefined;
};

const main = async () => {
  const email = readArg("email");
  const password = readArg("password");

  if (!email || !password) {
    console.error("Usage: node apps/web/scripts/seed-user.mjs --email you@example.com --password yourpass");
    process.exit(1);
  }

  const normalizedEmail = email.trim().toLowerCase();
  const url = resolveDbUrl();
  const authToken = url.startsWith("libsql://") ? resolveDbToken() : undefined;
  const client = createClient({ url, authToken });

  const existing = await client.execute({
    sql: "SELECT id FROM users WHERE email = ? LIMIT 1",
    args: [normalizedEmail]
  });

  const passwordHash = await hashPassword(password);
  const now = new Date().toISOString();

  if (existing.rows.length > 0) {
    const existingId = String(existing.rows[0].id);
    await client.execute({
      sql: "UPDATE users SET password_hash = ? WHERE id = ?",
      args: [passwordHash, existingId]
    });
    console.log(`Updated password for existing user ${normalizedEmail}`);
    process.exit(0);
  }

  const id = crypto.randomUUID();

  await client.execute({
    sql: "INSERT INTO users (id, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
    args: [id, normalizedEmail, passwordHash, now]
  });

  console.log(`Created user ${normalizedEmail} with id ${id}`);
};

main().catch((error) => {
  console.error("Failed to seed user:", error);
  process.exit(1);
});
