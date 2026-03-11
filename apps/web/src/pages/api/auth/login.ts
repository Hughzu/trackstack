import type { APIRoute } from "astro";
import { proxyAuthRequest } from "@/server/auth/proxyAuth";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  const contentType = request.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");

  let payload: { email?: string; password?: string };

  if (isJson) {
    try {
      payload = await request.json();
    } catch {
      return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      });
    }
  } else {
    const form = await request.formData();
    const email = form.get("email");
    const password = form.get("password");
    payload = {
      email: typeof email === "string" ? email : undefined,
      password: typeof password === "string" ? password : undefined,
    };
  }

  const response = await proxyAuthRequest(request, "/auth/login", payload);
  const headers = new Headers();

  for (const value of response.setCookieHeaders) {
    headers.append("set-cookie", value);
  }

  if (!isJson) {
    headers.set("Location", response.status >= 200 && response.status < 300 ? "/" : withErrorParam(request.url, "/login"));
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
