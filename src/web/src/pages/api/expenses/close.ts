import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async () => {
  const sheet = expensesService.closeSheet(getCurrentUserId());
  return new Response(JSON.stringify(sheet), {
    status: 200,
    headers: {
      "Content-Type": "application/json"
    }
  });
};
