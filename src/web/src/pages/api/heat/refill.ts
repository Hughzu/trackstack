import type { APIRoute } from 'astro';
import { heatService } from '@/modules/heat/services/heatService';
import { getCurrentUserId } from '@/shared/auth/currentUser';

export const prerender = false;

export const POST: APIRoute = async ({ request }) => {
  try {
    let data: { date?: string; weightKg?: number; bags?: number; temperature?: number } = {};
    try {
      data = await request.json();
    } catch {
      return new Response(JSON.stringify({ error: 'Invalid JSON body' }), {
        status: 400,
        headers: {
          'Content-Type': 'application/json'
        }
      });
    }
    
    // Basic validation
    if (!data.date || !data.weightKg || !data.bags) {
      return new Response(JSON.stringify({ error: 'Missing required fields' }), {
        status: 400
      });
    }

    const newRefill = heatService.addRefill({
      userId: getCurrentUserId(),
      date: new Date(data.date),
      weightKg: Number(data.weightKg),
      bags: Number(data.bags),
      temperature: data.temperature ? Number(data.temperature) : undefined
    });

    return new Response(JSON.stringify(newRefill), {
      status: 201,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  } catch (e) {
    console.error(e);
    return new Response(JSON.stringify({ error: 'Server Error' }), { status: 500 });
  }
};

export const DELETE: APIRoute = async ({ request }) => {
  let data: { id?: string } = {};
  try {
    data = await request.json();
  } catch {
    return new Response(JSON.stringify({ error: 'Invalid JSON body' }), {
      status: 400,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  if (!data.id) {
    return new Response(JSON.stringify({ error: 'Missing refill id' }), {
      status: 400,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  const deleted = heatService.deleteRefill(data.id, getCurrentUserId());
  if (!deleted) {
    return new Response(JSON.stringify({ error: 'Refill not found' }), {
      status: 404,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  return new Response(null, { status: 204 });
};
