import type { APIRoute } from "astro";
import { authDb } from "@/core/auth/db";
import { authConfig } from "@/core/auth/config";
import { hashToken } from "@/core/auth/session";

export const prerender = false;

export const POST: APIRoute = async ({ request, cookies }) => {
  const rawToken = cookies.get(authConfig.cookie.name)?.value;
  if (rawToken) {
    const tokenId = hashToken(rawToken);
    await authDb.revokeSession(tokenId, new Date().toISOString());
  }

  cookies.delete(authConfig.cookie.name, { path: authConfig.cookie.path });
  const contentType = request.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    return new Response(null, { status: 303, headers: { Location: "/login" } });
  }
  return new Response(null, { status: 204 });
};
