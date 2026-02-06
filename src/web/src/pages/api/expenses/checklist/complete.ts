import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  let data: { id?: string; date?: string } = {};
  try {
    data = await request.json();
  } catch {
    return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  if (!data.id) {
    return new Response(JSON.stringify({ error: "Missing checklist item id" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  const entry = expensesService.completeChecklistItem({
    id: data.id,
    userId: getCurrentUserId(),
    date: data.date
  });

  if (!entry) {
    return new Response(JSON.stringify({ error: "Checklist item not found" }), {
      status: 404,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  return new Response(JSON.stringify(entry), {
    status: 201,
    headers: {
      "Content-Type": "application/json"
    }
  });
};
