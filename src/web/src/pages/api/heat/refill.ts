import type { APIRoute } from 'astro';
import { heatService } from '@/modules/heat/services/heatService';

export const POST: APIRoute = async ({ request }) => {
  try {
    const data = await request.json();
    
    // Basic validation
    if (!data.date || !data.weightKg || !data.bags) {
      return new Response(JSON.stringify({ error: 'Missing required fields' }), {
        status: 400
      });
    }

    const newRefill = heatService.addRefill({
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
