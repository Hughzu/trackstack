export const prerender = false;

import { getDb } from "@/shared/db/sqlite";
import type { APIRoute } from "astro";

export const GET: APIRoute = async () => {
    try {
        const db = getDb("health");
        // Simple query to verify DB connection
        await db.run("SELECT 1");

        return new Response(
            JSON.stringify({
                status: "ok",
                uptime: process.uptime(),
                timestamp: new Date().toISOString()
            }),
            {
                status: 200,
                headers: {
                    "Content-Type": "application/json"
                }
            }
        );
    } catch (error) {
        console.error("Health check failed:", error);
        return new Response(
            JSON.stringify({
                status: "error",
                message: "Database connection failed"
            }),
            {
                status: 503,
                headers: {
                    "Content-Type": "application/json"
                }
            }
        );
    }
};
