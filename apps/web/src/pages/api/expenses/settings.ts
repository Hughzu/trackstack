import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const GET: APIRoute = async () => {
  try {
    const settings = await fetchApi("/expenses/settings");
    return new Response(JSON.stringify(settings), {
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
    const contentType = request.headers.get("content-type") ?? "";
    const isJson = contentType.includes("application/json");
    let data: {
      income?: number | string;
      ratioFund?: number | string;
      ratioFun?: number | string;
      ratioFuture?: number | string;
    } = {};

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
      const income = form.get("income");
      const ratioFund = form.get("ratioFund");
      const ratioFun = form.get("ratioFun");
      const ratioFuture = form.get("ratioFuture");
      data = {
        income: typeof income === "string" ? income : undefined,
        ratioFund: typeof ratioFund === "string" ? ratioFund : undefined,
        ratioFun: typeof ratioFun === "string" ? ratioFun : undefined,
        ratioFuture: typeof ratioFuture === "string" ? ratioFuture : undefined
      };
    }

    const parseNumber = (value?: number | string) => {
      if (value === undefined || value === null) return undefined;
      if (typeof value === "string" && value.trim().length === 0) return undefined;
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    };

    const payload = {
      income: parseNumber(data.income),
      ratioFund: parseNumber(data.ratioFund),
      ratioFun: parseNumber(data.ratioFun),
      ratioFuture: parseNumber(data.ratioFuture)
    };

    if (payload.income === undefined || payload.ratioFund === undefined || payload.ratioFun === undefined || payload.ratioFuture === undefined) {
      if (!isJson) {
        const fallback = request.headers.get("referer") ?? "/expenses/settings";
        return new Response(null, { status: 303, headers: { Location: withErrorParam(fallback) } });
      }
      return new Response(JSON.stringify({ error: "Missing required fields" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json"
        }
      });
    }

    const settings = await fetchApi("/expenses/settings", {
      method: "POST",
      body: JSON.stringify(payload)
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/expenses" } });
    }

    return new Response(JSON.stringify(settings), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Failed to proxy expenses settings POST:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
