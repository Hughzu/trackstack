import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";

export const prerender = false;

export const POST: APIRoute = async () => {
  try {
    const sheet = await fetchApi("/expenses/sheet/close", {
      method: "POST",
    });
    return new Response(JSON.stringify(sheet), {
      status: 200,
      headers: {
        "Content-Type": "application/json",
      },
    });
  } catch (error) {
    console.error("Error in POST /api/expenses/sheet/close:", error);
    return new Response(JSON.stringify({ error: error instanceof Error ? error.message : "Server Error" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
};
