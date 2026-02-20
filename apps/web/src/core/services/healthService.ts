import { getDb } from "@/core/db/sqlite";

export type HealthStatus = {
    status: "ok" | "error";
    uptime?: number;
    timestamp: string;
    message?: string;
};

export const healthService = {
    checkHealth: async (): Promise<HealthStatus> => {
        return {
            status: "ok",
            uptime: process.uptime(),
            timestamp: new Date().toISOString()
        };
    }
};
