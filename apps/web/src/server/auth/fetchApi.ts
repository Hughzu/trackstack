import { getCurrentToken } from "@/server/auth/currentUser";
import { authConfig } from "@/server/auth/config";

const API_BASE_URL = (typeof process !== 'undefined' && process.env.API_PROXY_URL ? `${process.env.API_PROXY_URL}/api` : null) || import.meta.env?.PUBLIC_API_BASE_URL || "http://127.0.0.1:8080/api";

export const fetchApi = async <T>(path: string, options: RequestInit = {}): Promise<T> => {
    const token = getCurrentToken();
    const headers = new Headers(options.headers || {});

    if (token) {
        headers.set("Cookie", `${authConfig.cookie.name}=${token}`);
    }

    headers.set("Content-Type", "application/json");

    // Format path
    const normalizedPath = path.startsWith("/api") ? path : `/api${path.startsWith("/") ? path : `/${path}`}`;

    // Ensure base URL doesn't end with slash, and remove /api if present to avoid duplication
    const baseUrl = API_BASE_URL.replace(/\/+$/, "").replace(/\/api$/, "");
    const finalUrl = `${baseUrl}${normalizedPath}`;

    const response = await fetch(finalUrl, {
        ...options,
        headers,
    });

    if (!response.ok) {
        let errorMsg = `API Request failed: ${response.status} ${response.statusText}`;
        try {
            const data = await response.json();
            if (data.error) errorMsg = data.error;
        } catch { }
        throw new Error(errorMsg);
    }

    if (response.status === 204) {
        return null as any;
    }

    return response.json();
};
