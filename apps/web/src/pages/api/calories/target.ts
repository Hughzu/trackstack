import type { APIRoute } from "astro";
import { caloriesService } from "@/modules/calories/services/caloriesService";
import { getCurrentUserId } from "@/server/auth/currentUser";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const GET: APIRoute = async () => {
  try {
    const target = await caloriesService.getTarget(getCurrentUserId());
    return new Response(JSON.stringify(target), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in GET /api/calories/target:", error);
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
      targetKcal?: number | string;
      targetProtein?: number | string;
      targetCarbs?: number | string;
      targetFat?: number | string;
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
      const targetKcal = form.get("targetKcal");
      const targetProtein = form.get("targetProtein");
      const targetCarbs = form.get("targetCarbs");
      const targetFat = form.get("targetFat");
      data = {
        targetKcal: typeof targetKcal === "string" ? targetKcal : undefined,
        targetProtein: typeof targetProtein === "string" ? targetProtein : undefined,
        targetCarbs: typeof targetCarbs === "string" ? targetCarbs : undefined,
        targetFat: typeof targetFat === "string" ? targetFat : undefined
      };
    }

    const parseOptionalNumber = (value?: number | string) => {
      if (value === undefined || value === null) return undefined;
      if (typeof value === "string" && value.trim().length === 0) return undefined;
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    };

    const targetKcal = parseOptionalNumber(data.targetKcal);
    const targetProtein = parseOptionalNumber(data.targetProtein);
    const targetCarbs = parseOptionalNumber(data.targetCarbs);
    const targetFat = parseOptionalNumber(data.targetFat);

    if (targetKcal === undefined || targetProtein === undefined) {
      if (!isJson) {
        const fallback = request.headers.get("referer") ?? "/calories/settings";
        return new Response(null, { status: 303, headers: { Location: withErrorParam(fallback) } });
      }
      return new Response(JSON.stringify({ error: "Missing required fields" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const target = await caloriesService.updateTarget({
      userId: getCurrentUserId(),
      targetKcal,
      targetProteinG: targetProtein,
      targetCarbsG: targetCarbs,
      targetFatG: targetFat
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/calories" } });
    }

    return new Response(JSON.stringify(target), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Error in POST /api/calories/target:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 500,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
