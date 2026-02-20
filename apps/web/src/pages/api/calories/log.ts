import type { APIRoute } from "astro";
import { caloriesService } from "@/modules/calories/services/caloriesService";
import { getCurrentUserId } from "@/server/auth/currentUser";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: {
      calories?: number | string;
      protein?: number | string;
      carbs?: number | string;
      fat?: number | string;
      title?: string;
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
      const calories = form.get("calories");
      const protein = form.get("protein");
      const carbs = form.get("carbs");
      const fat = form.get("fat");
      const title = form.get("title");
      const date = form.get("date");
      data = {
        calories: typeof calories === "string" ? calories : undefined,
        protein: typeof protein === "string" ? protein : undefined,
        carbs: typeof carbs === "string" ? carbs : undefined,
        fat: typeof fat === "string" ? fat : undefined,
        title: typeof title === "string" ? title : undefined,
        date: typeof date === "string" ? date : undefined
      };
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
      if (!isJson) {
        const fallback = request.headers.get("referer") ?? "/calories";
        return new Response(null, { status: 303, headers: { Location: withErrorParam(fallback) } });
      }
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

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/calories" } });
    }

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
      return new Response(JSON.stringify({ error: "Missing log id" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const deleted = await caloriesService.deleteLog(id, getCurrentUserId());
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
