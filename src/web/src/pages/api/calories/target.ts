import type { APIRoute } from "astro";
import { caloriesService } from "@/modules/calories/services/caloriesService";
import { getCurrentUserId } from "@/shared/auth/currentUser";

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
      targetKcal?: number;
      targetProtein?: number;
      targetCarbs?: number;
      targetFat?: number;
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

    const targetKcal = parseOptionalNumber(data.targetKcal);
    const targetProtein = parseOptionalNumber(data.targetProtein);
    const targetCarbs = parseOptionalNumber(data.targetCarbs);
    const targetFat = parseOptionalNumber(data.targetFat);

    if (targetKcal === undefined || targetProtein === undefined) {
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
