import type { APIRoute } from "astro";
import { expensesService } from "@/modules/expenses/services/expensesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const GET: APIRoute = async () => {
  try {
    const userId = getCurrentUserId();
    const settings = await expensesService.getSettings(userId);
    const checklist = await expensesService.getChecklistTemplates(userId);
    const recurring = await expensesService.getRecurringTemplates(userId);

    return new Response(JSON.stringify({ settings, checklist, recurring }), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in GET /api/expenses/settings:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: {
      income?: number | string;
      ratioFund?: number | string;
      ratioFun?: number | string;
      ratioFuture?: number | string;
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

    const parseNumber = (value?: number | string) => {
      if (value === undefined || value === null) return undefined;
      if (typeof value === "string" && value.trim().length === 0) return undefined;
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    };

    const income = parseNumber(data.income);
    const ratioFund = parseNumber(data.ratioFund);
    const ratioFun = parseNumber(data.ratioFun);
    const ratioFuture = parseNumber(data.ratioFuture);

    if (income === undefined || ratioFund === undefined || ratioFun === undefined || ratioFuture === undefined) {
      return new Response(JSON.stringify({ error: "Missing required fields" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const settings = await expensesService.updateSettings({
      userId: getCurrentUserId(),
      income,
      ratioFund,
      ratioFun,
      ratioFuture
    });

    return new Response(JSON.stringify(settings), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in POST /api/expenses/settings:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
