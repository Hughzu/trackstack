import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/server/auth/currentUser";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: {
      title?: string;
      amount?: number | string;
      category?: string;
      date?: string;
    } = {};

    const contentType = request.headers.get("content-type") ?? "";
    const isJson = contentType.includes("application/json");
    if (isJson) {
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
    } else {
      const form = await request.formData();
      const title = form.get("title");
      const amount = form.get("amount");
      const category = form.get("category");
      const date = form.get("date");
      data = {
        title: typeof title === "string" ? title : undefined,
        amount: typeof amount === "string" ? amount : undefined,
        category: typeof category === "string" ? category : undefined,
        date: typeof date === "string" ? date : undefined
      };
    }

    const amount = typeof data.amount === "string" ? Number(data.amount) : data.amount;
    if (!Number.isFinite(amount)) {
      if (!isJson) {
        const fallback = request.headers.get("referer") ?? "/expenses/new";
        return new Response(null, { status: 303, headers: { Location: withErrorParam(fallback) } });
      }
      return new Response(JSON.stringify({ error: "Missing required fields" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const title = data.title?.trim() ? data.title.trim() : "Untitled";

    const entry = await expensesService.addExpense({
      userId: getCurrentUserId(),
      title,
      amount: Number(amount),
      category: data.category,
      date: data.date
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/expenses" } });
    }

    return new Response(JSON.stringify(entry), {
      status: 201,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in POST /api/expenses/expense:", error);
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
      return new Response(JSON.stringify({ error: "Missing expense id" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const deleted = await expensesService.deleteExpense(id, getCurrentUserId());
    if (!deleted) {
      return new Response(JSON.stringify({ error: "Expense not found" }), {
        status: 404,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    return new Response(null, { status: 204 });
  } catch (error) {
    console.error("Error in DELETE /api/expenses/expense:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
