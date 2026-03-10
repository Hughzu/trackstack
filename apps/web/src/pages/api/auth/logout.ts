import type { APIRoute } from "astro";
import { proxyAuthRequest } from "@/server/auth/proxyAuth";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  return proxyAuthRequest(request, "/auth/logout");
};
