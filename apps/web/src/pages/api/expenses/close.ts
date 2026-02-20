import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/server/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async () => {
  try {
    const sheet = await expensesService.closeSheet(getCurrentUserId());
    return new Response(JSON.stringify(sheet), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in POST /api/expenses/close:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
