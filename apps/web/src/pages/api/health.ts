export const prerender = false;

import { healthService } from "@/server/services/healthService";
import type { APIRoute } from "astro";

export const GET: APIRoute = async () => {
    const health = await healthService.checkHealth();

    if (health.status === "error") {
        return new Response(JSON.stringify(health), {
            status: 503,
            headers: {
                "Content-Type": "application/json"
            }
        });
    }

    return new Response(JSON.stringify(health), {
        status: 200,
        headers: {
            "Content-Type": "application/json"
        }
    });
};
