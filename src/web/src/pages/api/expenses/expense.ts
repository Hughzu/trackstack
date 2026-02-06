import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  let data: {
    title?: string;
    amount?: number | string;
    category?: string;
    date?: string;
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
  if (!Number.isFinite(amount)) {
    return new Response(JSON.stringify({ error: "Missing required fields" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  const title = data.title?.trim() ? data.title.trim() : "Untitled";

  const entry = expensesService.addExpense({
    userId: getCurrentUserId(),
    title,
    amount: Number(amount),
    category: data.category,
    date: data.date
  });

  return new Response(JSON.stringify(entry), {
    status: 201,
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
    return new Response(JSON.stringify({ error: "Missing expense id" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  const deleted = expensesService.deleteExpense(data.id, getCurrentUserId());
  if (!deleted) {
    return new Response(JSON.stringify({ error: "Expense not found" }), {
      status: 404,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }

  return new Response(null, { status: 204 });
};
