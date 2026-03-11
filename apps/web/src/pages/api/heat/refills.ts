import type { APIRoute } from 'astro';
import { fetchApi } from '@/server/auth/fetchApi';
import { withErrorParam } from '@/server/http/redirects';

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  const contentType = request.headers.get('content-type') ?? '';
  const isJson = contentType.includes('application/json');
  const parseOptionalNumber = (value: unknown) => {
    if (typeof value === 'number' && Number.isFinite(value)) return value;
    if (typeof value === 'string' && value.trim() !== '') {
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : undefined;
    }
    return undefined;
  };

  let bodyData;
  if (isJson) {
    let payload;
    try {
      payload = await request.json();
    } catch {
      return new Response(JSON.stringify({ error: 'Invalid JSON body' }), {
        status: 400,
        headers: {
          'Content-Type': 'application/json'
        }
      });
    }
    bodyData = {
      date: typeof payload?.date === 'string' ? payload.date : undefined,
      weightKg: parseOptionalNumber(payload?.weightKg),
      bags: parseOptionalNumber(payload?.bags),
      temperature: parseOptionalNumber(payload?.temperature),
    };
  } else {
    const form = await request.formData();
    bodyData = {
      date: typeof form.get('date') === 'string' ? form.get('date') : undefined,
      weightKg: parseOptionalNumber(form.get('weightKg')),
      bags: parseOptionalNumber(form.get('bags')),
      temperature: parseOptionalNumber(form.get('temperature')),
    };
  }

  try {
    const refill = await fetchApi('/heat/refills', {
      method: 'POST',
      body: JSON.stringify(bodyData),
    });

    if (!isJson) {
      return new Response(null, { status: 303, headers: { Location: '/heat' } });
    }

    return new Response(JSON.stringify(refill), {
      status: 201,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  } catch (error: any) {
    console.error('Failed to proxy heat POST:', error);
    if (!isJson) {
      const fallback = request.headers.get('referer') ?? '/heat/new';
      return new Response(null, { status: 303, headers: { Location: withErrorParam(request.url, fallback) } });
    }
    return new Response(JSON.stringify({ error: error.message }), {
      status: 400,
      headers: {
        'Content-Type': 'application/json'
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
      id = url.searchParams.get('id');
    }

    if (!id) {
      return new Response(JSON.stringify({ error: 'Missing refill id' }), {
        status: 400,
        headers: {
          'Content-Type': 'application/json'
        }
      });
    }

    await fetchApi(`/heat/refills?id=${encodeURIComponent(id)}`, { method: 'DELETE' });
    return new Response(null, { status: 204 });
  } catch (error: any) {
    console.error('Failed to proxy heat DELETE:', error);
    return new Response(JSON.stringify({ error: error?.message ?? 'Server Error' }), {
      status: 400,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }
};
