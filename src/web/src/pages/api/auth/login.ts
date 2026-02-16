import type { APIRoute } from "astro";
import { authDb } from "@/core/auth/db";
import { authConfig } from "@/core/auth/config";
import { getClientContext } from "@/core/auth/client";
import { createSession, getCookieOptions } from "@/core/auth/session";
import { verifyPassword } from "@/core/auth/password";

export const prerender = false;

export const POST: APIRoute = async ({ request, cookies }) => {
  let payload: { email?: string; password?: string } = {};
  try {
    payload = await request.json();
  } catch {
    return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
      status: 400,
      headers: { "Content-Type": "application/json" }
    });
  }

  const email = payload.email?.trim().toLowerCase();
  const password = payload.password ?? "";
  if (!email || !password) {
    return new Response(JSON.stringify({ error: "Missing credentials" }), {
      status: 400,
      headers: { "Content-Type": "application/json" }
    });
  }

  const user = await authDb.findUserByEmail(email);
  if (!user) {
    return new Response(JSON.stringify({ error: "Unauthorized" }), {
      status: 401,
      headers: { "Content-Type": "application/json" }
    });
  }

  const valid = await verifyPassword(password, user.passwordHash);
  if (!valid) {
    return new Response(JSON.stringify({ error: "Unauthorized" }), {
      status: 401,
      headers: { "Content-Type": "application/json" }
    });
  }

  const context = getClientContext(request);
  const { rawToken, session } = await createSession(user.id, context);
  const now = new Date();
  const cookieOptions = getCookieOptions(now, new Date(session.expiresAt));
  cookies.set(authConfig.cookie.name, rawToken, cookieOptions);

  await authDb.updateUserLastLogin(user.id, now.toISOString());

  return new Response(null, { status: 204 });
};
