import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  // Proxy POST directly to Go backend
  let bodyData;
  const contentType = request.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");
  const parseAmount = (value: unknown) => {
    if (typeof value === "number" && Number.isFinite(value)) return value;
    if (typeof value === "string" && value.trim() !== "") {
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    }
    return undefined;
  };

  if (!isJson) {
    const form = await request.formData();
    const title = form.get("title");
    const amount = form.get("amount");
    const category = form.get("category");
    const date = form.get("date");
    bodyData = {
      title: typeof title === "string" ? title : undefined,
      amount: parseAmount(amount),
      category: typeof category === "string" ? category : undefined,
      date: typeof date === "string" ? date : undefined
    };
  } else {
    const payload = await request.json();
    bodyData = {
      title: typeof payload?.title === "string" ? payload.title : undefined,
      amount: parseAmount(payload?.amount),
      category: typeof payload?.category === "string" ? payload.category : undefined,
      date: typeof payload?.date === "string" ? payload.date : undefined
    };
  }

  try {
    const entry = await fetchApi("/expenses/expense", {
      method: "POST",
      body: JSON.stringify(bodyData)
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/expenses" } });
    }

    return new Response(JSON.stringify(entry), {
      status: 201,
      headers: { "Content-Type": "application/json" }
    });
  } catch (err: any) {
    console.error("Failed to proxy expense POST:", err);
    if (!isJson) {
      const fallback = request.headers.get("referer") ?? "/expenses/new";
      return new Response(null, { status: 303, headers: { Location: withErrorParam(request.url, fallback) } });
    }
    return new Response(JSON.stringify({ error: err.message }), { status: 400 });
  }
};

export const DELETE: APIRoute = async ({ request }) => {
	let id: string | null = null;
	try {
		const data = await request.json();
    if (data?.id) id = String(data.id);
  } catch {
    // ignore
  }

	if (!id) {
		const url = new URL(request.url);
		id = url.searchParams.get("id");
	}

	if (!id) {
		return new Response(JSON.stringify({ error: "Missing expense id" }), {
			status: 400,
			headers: { "Content-Type": "application/json" }
		});
	}

	try {
		await fetchApi(`/expenses/expense?id=${encodeURIComponent(id)}`, { method: "DELETE" });
		return new Response(null, { status: 204 });
	} catch (err: any) {
		return new Response(JSON.stringify({ error: err.message }), { status: 400 });
  }
};
