import type { APIRoute } from "astro";
import { caloriesService } from "@/modules/calories/services/caloriesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: {
      calories?: number;
      protein?: number;
      carbs?: number;
      fat?: number;
      title?: string;
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

    const parseOptionalNumber = (value?: number | string) => {
      if (value === undefined || value === null) return undefined;
      if (typeof value === "string" && value.trim().length === 0) return undefined;
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    };

    const calories = parseOptionalNumber(data.calories);
    const protein = parseOptionalNumber(data.protein);
    const carbs = parseOptionalNumber(data.carbs);
    const fat = parseOptionalNumber(data.fat);

    if (calories === undefined || protein === undefined) {
      return new Response(JSON.stringify({ error: "Missing required fields" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const date = data.date ?? new Date().toISOString().split("T")[0];

    const log = await caloriesService.addLog({
      userId: getCurrentUserId(),
      date,
      calories,
      proteinG: protein,
      carbsG: carbs,
      fatG: fat,
      title: data.title?.trim() ? data.title.trim() : undefined
    });

    return new Response(JSON.stringify(log), {
      status: 201,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in POST /api/calories/log:", error);
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
      return new Response(JSON.stringify({ error: "Missing log id" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const deleted = await caloriesService.deleteLog(data.id, getCurrentUserId());
    if (!deleted) {
      return new Response(JSON.stringify({ error: "Log not found" }), {
        status: 404,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    return new Response(null, { status: 204 });
  } catch (error) {
    console.error("Error in DELETE /api/calories/log:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
