import type { APIRoute } from "astro";
import { authDb } from "@/server/auth/db";
import { authConfig } from "@/server/auth/config";
import { getClientContext } from "@/server/auth/client";
import { createSession, getCookieOptions } from "@/server/auth/session";
import { verifyPassword } from "@/server/auth/password";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request, cookies }) => {
  let payload: { email?: string; password?: string } = {};
  const contentType = request.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");

  if (isJson) {
    try {
      payload = await request.json();
    } catch {
      return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
        status: 400,
        headers: { "Content-Type": "application/json" }
      });
    }
  } else {
    const form = await request.formData();
    const email = form.get("email");
    const password = form.get("password");
    payload = {
      email: typeof email === "string" ? email : undefined,
      password: typeof password === "string" ? password : undefined
    };
  }

  const email = payload.email?.trim().toLowerCase();
  const password = payload.password ?? "";
  if (!email || !password) {
    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: withErrorParam(request.url, "/login") } });
    }
    return new Response(JSON.stringify({ error: "Missing credentials" }), {
      status: 400,
      headers: { "Content-Type": "application/json" }
    });
  }

  const user = await authDb.findUserByEmail(email);
  if (!user) {
    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: withErrorParam(request.url, "/login") } });
    }
    return new Response(JSON.stringify({ error: "Unauthorized" }), {
      status: 401,
      headers: { "Content-Type": "application/json" }
    });
  }

  const valid = await verifyPassword(password, user.passwordHash);
  if (!valid) {
    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: withErrorParam(request.url, "/login") } });
    }
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

  if (!isJson) {
    return new Response(null, { status: 303, headers: { Location: "/" } });
  }

  return new Response(null, { status: 204 });
};
