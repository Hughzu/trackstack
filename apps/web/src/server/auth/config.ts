type CookieSameSite = "lax" | "strict" | "none";

const readEnv = (key: string) => {
  const metaEnv = (import.meta as { env?: Record<string, string | undefined> }).env;
  return metaEnv?.[key] ?? process.env[key];
};

const readNumber = (key: string, fallback: number) => {
  const value = readEnv(key);
  if (!value) return fallback;
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
};

const readString = (key: string, fallback: string) => {
  const value = readEnv(key);
  return value && value.trim().length > 0 ? value.trim() : fallback;
};

const readSameSite = (key: string, fallback: CookieSameSite) => {
  const value = readEnv(key);
  if (!value) return fallback;
  const normalized = value.trim().toLowerCase();
  if (normalized === "strict" || normalized === "none" || normalized === "lax") return normalized;
  return fallback;
};

export const authConfig = {
  cookie: {
    name: readString("AUTH_COOKIE_NAME", "trackstack_session"),
    secure: readString("AUTH_COOKIE_SECURE", process.env.NODE_ENV === "production" ? "true" : "false") === "true",
    sameSite: readSameSite("AUTH_COOKIE_SAMESITE", "lax"),
    path: "/",
    httpOnly: true
  },
  session: {
    idleSeconds: readNumber("AUTH_SESSION_IDLE_SECONDS", 1800),
    absoluteSeconds: readNumber("AUTH_SESSION_ABSOLUTE_SECONDS", 86400),
    rotateAfterSeconds: readNumber("AUTH_SESSION_ROTATE_AFTER_SECONDS", 1800),
    graceSeconds: readNumber("AUTH_SESSION_ROTATION_GRACE_SECONDS", 300),
    touchAfterSeconds: readNumber("AUTH_SESSION_TOUCH_SECONDS", 300)
  }
};
