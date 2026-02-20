import { randomBytes, scrypt, timingSafeEqual } from "node:crypto";
import { promisify } from "node:util";

const scryptAsync = promisify(scrypt);

const DEFAULT_PARAMS = {
  N: 16384,
  r: 8,
  p: 1,
  keyLength: 64,
  saltLength: 16
};

const parseParams = (segment: string) => {
  const parsed: Record<string, number> = {};
  for (const part of segment.split(",")) {
    const [key, value] = part.split("=");
    const numeric = Number(value);
    if (key && Number.isFinite(numeric)) parsed[key.trim()] = numeric;
  }
  return parsed;
};

export const hashPassword = async (password: string) => {
  const salt = randomBytes(DEFAULT_PARAMS.saltLength);
  const derivedKey = (await scryptAsync(password, salt, DEFAULT_PARAMS.keyLength, {
    cost: DEFAULT_PARAMS.N,
    blockSize: DEFAULT_PARAMS.r,
    parallelization: DEFAULT_PARAMS.p
  })) as Buffer;

  const params = `N=${DEFAULT_PARAMS.N},r=${DEFAULT_PARAMS.r},p=${DEFAULT_PARAMS.p}`;
  const saltEncoded = salt.toString("base64");
  const hashEncoded = derivedKey.toString("base64");
  return `$scrypt$${params}$${saltEncoded}$${hashEncoded}`;
};

export const verifyPassword = async (password: string, storedHash: string) => {
  if (!storedHash.startsWith("$scrypt$")) return false;
  const parts = storedHash.split("$");
  if (parts.length < 5) return false;

  const params = parseParams(parts[2] ?? "");
  const saltEncoded = parts[3];
  const hashEncoded = parts[4];
  if (!saltEncoded || !hashEncoded) return false;

  const salt = Buffer.from(saltEncoded, "base64");
  const expected = Buffer.from(hashEncoded, "base64");

  const derivedKey = (await scryptAsync(password, salt, expected.length, {
    cost: params.N ?? DEFAULT_PARAMS.N,
    blockSize: params.r ?? DEFAULT_PARAMS.r,
    parallelization: params.p ?? DEFAULT_PARAMS.p
  })) as Buffer;

  if (derivedKey.length !== expected.length) return false;
  return timingSafeEqual(derivedKey, expected);
};
