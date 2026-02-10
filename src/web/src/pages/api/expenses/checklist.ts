import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  let data: {
    id?: string;
    title?: string;
    amount?: number | string;
    category?: string;
  } = {};

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

  const amount = typeof data.amount === "string" ? Number(data.amount) : data.amount;
  if (!data.title?.trim() || !Number.isFinite(amount)) {
    return new Response(JSON.stringify({ error: "Missing required fields" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  const template = await expensesService.upsertChecklistTemplate({
    id: data.id,
    userId: getCurrentUserId(),
    title: data.title.trim(),
    amount: Number(amount),
    category: data.category
  });

  return new Response(JSON.stringify(template), {
    status: 200,
    headers: {
      "Content-Type": "application/json"
    }
  });
};

export const DELETE: APIRoute = async ({ request }) => {
  let data: { id?: string } = {};
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
    return new Response(JSON.stringify({ error: "Missing template id" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  const deleted = await expensesService.deleteChecklistTemplate(data.id, getCurrentUserId());
  if (!deleted) {
    return new Response(JSON.stringify({ error: "Template not found" }), {
      status: 404,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  return new Response(null, { status: 204 });
};
