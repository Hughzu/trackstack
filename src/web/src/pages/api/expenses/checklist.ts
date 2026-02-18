import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/core/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
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
  } catch (error) {
    console.error("Error in POST /api/expenses/checklist:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};

export const DELETE: APIRoute = async ({ request }) => {
  try {
    let id: string | null = null;
    try {
      const data = await request.json();
      if (data?.id) id = String(data.id);
    } catch {
      // ignore missing/invalid body
    }

    if (!id) {
      const url = new URL(request.url);
      id = url.searchParams.get("id");
    }

    if (!id) {
      return new Response(JSON.stringify({ error: "Missing template id" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const deleted = await expensesService.deleteChecklistTemplate(id, getCurrentUserId());
    if (!deleted) {
      return new Response(JSON.stringify({ error: "Template not found" }), {
        status: 404,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    return new Response(null, { status: 204 });
  } catch (error) {
    console.error("Error in DELETE /api/expenses/checklist:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
