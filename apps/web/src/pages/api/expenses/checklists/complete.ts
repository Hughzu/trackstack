import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: { id?: string; date?: string } = {};

    const contentType = request.headers.get("content-type") ?? "";
    const isJson = contentType.includes("application/json");
    if (isJson) {
      try {
        data = await request.json();
      } catch {
        return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
          status: 400,
          headers: {
            "Content-Type": "application/json",
          },
        });
      }
    } else {
      const form = await request.formData();
      const id = form.get("id");
      const date = form.get("date");
      data = {
        id: typeof id === "string" ? id : undefined,
        date: typeof date === "string" ? date : undefined,
      };
    }

    if (!data.id) {
      if (!isJson) {
        const fallback = request.headers.get("referer") ?? "/expenses";
        return new Response(null, { status: 303, headers: { Location: withErrorParam(fallback) } });
      }
      return new Response(JSON.stringify({ error: "Missing checklist item id" }), {
        status: 400,
        headers: {
          "Content-Type": "application/json",
        },
      });
    }

    const entry = await fetchApi("/expenses/checklists/complete", {
      method: "POST",
      body: JSON.stringify({
        id: data.id,
        date: data.date,
      }),
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/expenses" } });
    }

    return new Response(JSON.stringify(entry), {
      status: 201,
      headers: {
        "Content-Type": "application/json",
      },
    });
  } catch (error) {
    console.error("Error in POST /api/expenses/checklists/complete:", error);
    return new Response(JSON.stringify({ error: error instanceof Error ? error.message : "Server Error" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
};
