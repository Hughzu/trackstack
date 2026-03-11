import type { APIRoute } from "astro";
import { proxyAuthRequest } from "@/server/auth/proxyAuth";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  const contentType = request.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");
  const response = await proxyAuthRequest(request, "/auth/logout");
  const headers = new Headers();

  for (const value of response.setCookieHeaders) {
    headers.append("set-cookie", value);
  }

  if (!isJson) {
    headers.set("Location", "/login");
    return new Response(null, { status: 303, headers });
  }

  if (response.contentType) {
    headers.set("Content-Type", response.contentType);
  }

  return new Response(response.body, {
    status: response.status,
    headers,
  });
};
