import type { APIRoute } from "astro";
import { fetchApi } from "@/server/auth/fetchApi";
import { withErrorParam } from "@/server/http/redirects";

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  let bodyData;
  const contentType = request.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");

  if (!isJson) {
    const form = await request.formData();
    const calories = form.get("calories");
    const protein = form.get("protein");
    const carbs = form.get("carbs");
    const fat = form.get("fat");
    const title = form.get("title");
    const date = form.get("date");

    const pNum = (s: any) => typeof s === "string" && s.trim() !== "" ? Number(s) : undefined;

    bodyData = {
      calories: pNum(calories),
      protein: pNum(protein),
      carbs: pNum(carbs),
      fat: pNum(fat),
      title: typeof title === "string" ? title : undefined,
      date: typeof date === "string" ? date : undefined
    };
  } else {
    bodyData = await request.json();
  }

  try {
    const entry = await fetchApi("/calories/log", {
      method: "POST",
      body: JSON.stringify(bodyData)
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: "/calories" } });
    }

    return new Response(JSON.stringify(entry), {
      status: 201,
      headers: { "Content-Type": "application/json" }
    });
  } catch (err: any) {
    console.error("Failed to proxy calorie POST:", err);
    if (!isJson) {
      const fallback = request.headers.get("referer") ?? "/calories";
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
    return new Response(JSON.stringify({ error: "Missing log id" }), {
      status: 400,
      headers: { "Content-Type": "application/json" }
    });
  }

  try {
    await fetchApi(`/calories/log?id=${encodeURIComponent(id)}`, { method: "DELETE" });
    return new Response(null, { status: 204 });
  } catch (err: any) {
    return new Response(JSON.stringify({ error: err.message }), { status: 400 });
  }
};
