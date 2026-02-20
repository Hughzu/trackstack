import { getDb } from "@/core/db/sqlite";

const DOMAIN = "users";

export type UserRecord = {
  id: string;
  email: string;
  passwordHash: string;
  sessionVersion: number;
  createdAt: string;
  lastLoginAt?: string | null;
};

export type SessionRecord = {
  id: string;
  userId: string;
  createdAt: string;
  expiresAt: string;
  rotatedAt: string;
  lastSeenAt: string;
  absoluteExpiresAt: string;
  parentId?: string | null;
  revokedAt?: string | null;
  userAgentHash?: string | null;
  ipPrefix?: string | null;
};

export const getAuthDb = () => getDb(DOMAIN);

export const authDb = {
  findUserByEmail: async (email: string) => {
    const db = getAuthDb();
    return db.get<UserRecord>(
      `SELECT
        id,
        email,
        password_hash as passwordHash,
        session_version as sessionVersion,
        created_at as createdAt,
        last_login_at as lastLoginAt
      FROM users
      WHERE email = ?
      LIMIT 1`,
      [email]
    );
  },
  updateUserLastLogin: async (userId: string, lastLoginAt: string) => {
    const db = getAuthDb();
    await db.run("UPDATE users SET last_login_at = ? WHERE id = ?", [lastLoginAt, userId]);
  },
  findSessionById: async (id: string) => {
    const db = getAuthDb();
    return db.get<SessionRecord>(
      `SELECT
        id,
        user_id as userId,
        created_at as createdAt,
        expires_at as expiresAt,
        rotated_at as rotatedAt,
        last_seen_at as lastSeenAt,
        absolute_expires_at as absoluteExpiresAt,
        parent_id as parentId,
        revoked_at as revokedAt,
        user_agent_hash as userAgentHash,
        ip_prefix as ipPrefix
      FROM sessions
      WHERE id = ?
      LIMIT 1`,
      [id]
    );
  },
  insertSession: async (session: SessionRecord) => {
    const db = getAuthDb();
    await db.run(
      `INSERT INTO sessions (
        id,
        user_id,
        created_at,
        expires_at,
        rotated_at,
        last_seen_at,
        absolute_expires_at,
        parent_id,
        revoked_at,
        user_agent_hash,
        ip_prefix
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)` ,
      [
        session.id,
        session.userId,
        session.createdAt,
        session.expiresAt,
        session.rotatedAt,
        session.lastSeenAt,
        session.absoluteExpiresAt,
        session.parentId ?? null,
        session.revokedAt ?? null,
        session.userAgentHash ?? null,
        session.ipPrefix ?? null
      ]
    );
  },
  touchSession: async (id: string, lastSeenAt: string, expiresAt: string) => {
    const db = getAuthDb();
    await db.run(
      "UPDATE sessions SET last_seen_at = ?, expires_at = ? WHERE id = ?",
      [lastSeenAt, expiresAt, id]
    );
  },
  rotateOutSession: async (id: string, revokedAt: string, expiresAt: string, rotatedAt: string) => {
    const db = getAuthDb();
    await db.run(
      "UPDATE sessions SET revoked_at = ?, expires_at = ?, rotated_at = ? WHERE id = ?",
      [revokedAt, expiresAt, rotatedAt, id]
    );
  },
  revokeSession: async (id: string, revokedAt: string) => {
    const db = getAuthDb();
    await db.run(
      "UPDATE sessions SET revoked_at = ?, expires_at = ? WHERE id = ?",
      [revokedAt, revokedAt, id]
    );
  }
};
