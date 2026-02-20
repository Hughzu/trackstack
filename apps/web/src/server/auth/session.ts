import crypto from "node:crypto";
import { authConfig } from "@/server/auth/config";
import { authDb, type SessionRecord } from "@/server/auth/db";

type SessionEvaluation = {
  valid: boolean;
  needsRotation: boolean;
  needsTouch: boolean;
};

type ClientContext = {
  userAgent?: string | null;
  ipPrefix?: string | null;
};

const toIso = (date: Date) => date.toISOString();

const addSeconds = (date: Date, seconds: number) => new Date(date.getTime() + seconds * 1000);

const resolveIdleExpiry = (now: Date, absoluteExpiresAt: Date) => {
  const idleExpiry = addSeconds(now, authConfig.session.idleSeconds);
  return idleExpiry < absoluteExpiresAt ? idleExpiry : absoluteExpiresAt;
};

export const newToken = () => crypto.randomBytes(32).toString("hex");

export const hashToken = (token: string) =>
  crypto.createHash("sha256").update(token).digest("hex");

export const evaluateSession = (session: SessionRecord, now = new Date()): SessionEvaluation => {
  const expiresAt = new Date(session.expiresAt);
  const absoluteExpiresAt = new Date(session.absoluteExpiresAt);
  if (absoluteExpiresAt <= now) return { valid: false, needsRotation: false, needsTouch: false };
  if (expiresAt <= now) return { valid: false, needsRotation: false, needsTouch: false };

  const rotatedAt = new Date(session.rotatedAt);
  const lastSeenAt = new Date(session.lastSeenAt);
  const needsRotation =
    Boolean(session.revokedAt) || now.getTime() - rotatedAt.getTime() > authConfig.session.rotateAfterSeconds * 1000;
  const needsTouch = now.getTime() - lastSeenAt.getTime() > authConfig.session.touchAfterSeconds * 1000;

  return { valid: true, needsRotation, needsTouch };
};

export const createSession = async (userId: string, context: ClientContext) => {
  const rawToken = newToken();
  const tokenId = hashToken(rawToken);
  const now = new Date();
  const absoluteExpiresAt = addSeconds(now, authConfig.session.absoluteSeconds);
  const expiresAt = resolveIdleExpiry(now, absoluteExpiresAt);

  const session: SessionRecord = {
    id: tokenId,
    userId,
    createdAt: toIso(now),
    expiresAt: toIso(expiresAt),
    rotatedAt: toIso(now),
    lastSeenAt: toIso(now),
    absoluteExpiresAt: toIso(absoluteExpiresAt),
    parentId: null,
    revokedAt: null,
    userAgentHash: context.userAgent ? hashToken(context.userAgent) : null,
    ipPrefix: context.ipPrefix ?? null
  };

  await authDb.insertSession(session);

  return { rawToken, session };
};

export const rotateSession = async (session: SessionRecord, context: ClientContext) => {
  const now = new Date();
  const rawToken = newToken();
  const tokenId = hashToken(rawToken);
  const absoluteExpiresAt = new Date(session.absoluteExpiresAt);
  const expiresAt = resolveIdleExpiry(now, absoluteExpiresAt);

  const replacement: SessionRecord = {
    id: tokenId,
    userId: session.userId,
    createdAt: toIso(now),
    expiresAt: toIso(expiresAt),
    rotatedAt: toIso(now),
    lastSeenAt: toIso(now),
    absoluteExpiresAt: session.absoluteExpiresAt,
    parentId: session.id,
    revokedAt: null,
    userAgentHash: context.userAgent ? hashToken(context.userAgent) : session.userAgentHash ?? null,
    ipPrefix: context.ipPrefix ?? session.ipPrefix ?? null
  };

  await authDb.insertSession(replacement);

  const graceExpiry = addSeconds(now, authConfig.session.graceSeconds);
  const shortenedExpiry = graceExpiry < absoluteExpiresAt ? graceExpiry : absoluteExpiresAt;
  await authDb.rotateOutSession(session.id, toIso(now), toIso(shortenedExpiry), toIso(now));

  return { rawToken, sessionId: replacement.id, expiresAt };
};

export const touchSession = async (session: SessionRecord) => {
  const now = new Date();
  const absoluteExpiresAt = new Date(session.absoluteExpiresAt);
  const expiresAt = resolveIdleExpiry(now, absoluteExpiresAt);
  await authDb.touchSession(session.id, toIso(now), toIso(expiresAt));
};

export const getCookieOptions = (now: Date, expiresAt: Date) => {
  const maxAge = Math.max(0, Math.floor((expiresAt.getTime() - now.getTime()) / 1000));
  return {
    httpOnly: authConfig.cookie.httpOnly,
    secure: authConfig.cookie.secure,
    sameSite: authConfig.cookie.sameSite,
    path: authConfig.cookie.path,
    maxAge
  } as const;
};
