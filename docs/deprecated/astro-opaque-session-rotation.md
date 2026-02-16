# Astro Opaque Session Auth (httpOnly Cookie + Rotation)

This is a minimal, drop-in blueprint for implementing opaque session tokens with rotation in Astro endpoints.
It avoids JWTs entirely and stores only a random session token in a secure cookie.

This approach aligns well with Trackstack's "Managed Polylith" goals: minimal dependencies, easy to port to Go, and a clean core/module split.

## Flow Overview

- Login: create session record -> set cookie
- Auth: read cookie -> look up session in DB -> allow/deny
- Rotate: replace token on a schedule (e.g. every request or every N minutes)
- Logout: delete session -> clear cookie

## Table Schema (example)

```sql
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,              -- random token or token hash
  user_id TEXT NOT NULL,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  rotated_at TEXT,
  parent_id TEXT,                    -- for rotation chain (optional)
  revoked_at TEXT                    -- set if invalidated
);
```

Notes:
- You can store `id` as a hashed token for better security.
- `parent_id` lets you track rotation chains if you want to debug/inspect.

Optional fields (low-cost hardening):
- `absolute_expires_at`: hard cap to prevent perpetual sessions.
- `last_seen_at`: for idle timeout and lower rotation churn.
- `user_agent_hash` / `ip_prefix`: optional device fingerprinting (audit only).

## Cookie Settings (recommended defaults)

- Name: `session`
- `httpOnly: true`
- `secure: true` (true in production)
- `sameSite: "lax"`
- `path: "/"`
- `maxAge`: matches session expiry
- If you add an absolute expiry, set `maxAge` to the shorter of idle expiry and absolute expiry.

## Token Creation Helpers

```ts
import crypto from "node:crypto";

export function newToken() {
  return crypto.randomBytes(32).toString("hex");
}

export function hashToken(token: string) {
  return crypto.createHash("sha256").update(token).digest("hex");
}
```

## Password Hashing (Go-Compatible)

Never store passwords in plain text. Store a strong password hash that Go can verify later.

Recommended:
- Argon2id with explicit parameters
- Store the full PHC-formatted hash string (includes algorithm + params + salt + hash)

This lets a future Go service verify the same hashes without rehashing or migration.

Minimal approach:
- On signup: hash the password -> store PHC string
- On login: verify password against the stored PHC string

Example PHC string format:

```
$argon2id$v=19$m=65536,t=3,p=1$<salt>$<hash>
```

If you prefer bcrypt, store the bcrypt hash string and verify in Go with the same cost.

## Login Endpoint (Astro)

```ts
// src/pages/api/auth/login.ts
import type { APIRoute } from "astro";
import { newToken, hashToken } from "../../lib/session";

export const POST: APIRoute = async ({ cookies, request, locals }) => {
  const body = await request.json();
  const { email, password } = body;

  // 1) Verify credentials
  const user = await locals.db.users.verify(email, password);
  if (!user) return new Response("Unauthorized", { status: 401 });

  // 2) Create session
  const rawToken = newToken();
  const tokenId = hashToken(rawToken);
  const now = new Date();
  const expiresAt = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

  await locals.db.sessions.insert({
    id: tokenId,
    userId: user.id,
    createdAt: now.toISOString(),
    expiresAt: expiresAt.toISOString(),
    rotatedAt: now.toISOString(),
    parentId: null,
    revokedAt: null,
  });

  // 3) Set cookie
  cookies.set("session", rawToken, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    maxAge: 7 * 24 * 60 * 60,
  });

  return new Response(null, { status: 204 });
};
```

## Auth Middleware (Astro Server Hooks)

```ts
// src/middleware.ts
import type { MiddlewareHandler } from "astro";
import { hashToken } from "./lib/session";

export const onRequest: MiddlewareHandler = async (context, next) => {
  const raw = context.cookies.get("session")?.value;
  if (!raw) return next();

  const tokenId = hashToken(raw);
  const session = await context.locals.db.sessions.findById(tokenId);

  if (!session) return next();
  if (session.revokedAt) return next();
  if (new Date(session.expiresAt) < new Date()) return next();

  // Optional: global logout / password change invalidation
  // if (session.sessionVersion !== user.sessionVersion) return next();

  context.locals.userId = session.userId;
  context.locals.sessionId = session.id;
  return next();
};
```

## Rotation Endpoint (Replace Token)

```ts
// src/pages/api/auth/rotate.ts
import type { APIRoute } from "astro";
import { newToken, hashToken } from "../../lib/session";

export const POST: APIRoute = async ({ cookies, locals }) => {
  const raw = cookies.get("session")?.value;
  if (!raw) return new Response("Unauthorized", { status: 401 });

  const currentId = hashToken(raw);
  const current = await locals.db.sessions.findById(currentId);
  if (!current || current.revokedAt) return new Response("Unauthorized", { status: 401 });
  if (new Date(current.expiresAt) < new Date()) return new Response("Unauthorized", { status: 401 });

  // Create replacement
  const newRaw = newToken();
  const newId = hashToken(newRaw);
  const now = new Date();
  const expiresAt = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

  await locals.db.sessions.insert({
    id: newId,
    userId: current.userId,
    createdAt: now.toISOString(),
    expiresAt: expiresAt.toISOString(),
    rotatedAt: now.toISOString(),
    parentId: current.id,
    revokedAt: null,
  });

  await locals.db.sessions.revoke(current.id, now.toISOString());

  cookies.set("session", newRaw, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    maxAge: 7 * 24 * 60 * 60,
  });

  return new Response(null, { status: 204 });
};
```

## Logout Endpoint

```ts
// src/pages/api/auth/logout.ts
import type { APIRoute } from "astro";
import { hashToken } from "../../lib/session";

export const POST: APIRoute = async ({ cookies, locals }) => {
  const raw = cookies.get("session")?.value;
  if (raw) {
    const tokenId = hashToken(raw);
    await locals.db.sessions.revoke(tokenId, new Date().toISOString());
  }

  cookies.delete("session", { path: "/" });
  return new Response(null, { status: 204 });
};
```

## Rotation Policy (Minimal)

Option A: rotate on every authenticated request.
- Simple, but creates more DB writes.

Option B: rotate every N minutes.
- Store `rotatedAt` and only rotate if older than N minutes.

Pseudo-check:

```ts
const shouldRotate = (rotatedAt: string) => {
  const last = new Date(rotatedAt).getTime();
  return Date.now() - last > 15 * 60 * 1000;
};
```

Recommended additions (still minimal):

- **Absolute expiry:** enforce a hard cap (e.g. 30 days) even with rotation.
- **Idle timeout:** expire if `last_seen_at` is older than N hours.
- **Global logout:** store `users.session_version` and compare on auth; bump to revoke all sessions.

## Notes

- Prefer storing token hashes in DB.
- Always set `secure: true` in production.
- `SameSite=Lax` is a good default and still allows normal navigation.
- If you need CSRF protection for state-changing requests, add a CSRF token.
- Beware of rotation race conditions: multiple concurrent requests can revoke a token that was just rotated. Allow a short overlap window or make revocation conditional.
- If rotating on every request, DB writes can grow quickly. Time-based rotation + idle timeout is usually enough.

## Minimal DB API (what you need)

```ts
sessions.insert({ id, userId, createdAt, expiresAt, rotatedAt, parentId, revokedAt })
sessions.findById(id)
sessions.revoke(id, revokedAt)
```

Optional low-cost helpers:

```ts
sessions.touch(id, { lastSeenAt, rotatedAt })
sessions.revokeByUser(userId)
```
