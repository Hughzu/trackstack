import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    const contentType = request.headers.get("content-type") ?? "";
    const isJson = contentType.includes("application/json");
    let data: {
      id?: string;
      title?: string;
      amount?: number | string;
      category?: string;
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
      const title = form.get("title");
      const amount = form.get("amount");
      const category = form.get("category");
      data = {
        title: typeof title === "string" ? title : undefined,
        amount: typeof amount === "string" ? amount : undefined,
        category: typeof category === "string" ? category : undefined
      };
    }

    const amount = typeof data.amount === "string" ? Number(data.amount) : data.amount;
    if (!data.title?.trim() || !Number.isFinite(amount)) {
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

		const payload = {
			id: data.id,
			title: data.title.trim(),
			amount: Number(amount),
			category: data.category
		};

		const template = await fetchApi("/expenses/checklists", {
			method: "POST",
			body: JSON.stringify(payload)
		});

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/expenses/settings" } });
    }

    return new Response(JSON.stringify(template), {
      status: 200,
      headers: {
        "Content-Type": "application/json"
      }
    });
  } catch (error) {
    console.error("Failed to proxy expenses checklist POST:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 400,
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

		await fetchApi(`/expenses/checklists?id=${encodeURIComponent(id)}`, { method: "DELETE" });
		return new Response(null, { status: 204 });
	} catch (error) {
		console.error("Failed to proxy expenses checklist DELETE:", error);
    return new Response(JSON.stringify({ error: "Server Error" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json"
      }
    });
  }
};
