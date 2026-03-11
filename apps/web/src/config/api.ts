const publicApiBaseUrl = (import.meta.env.PUBLIC_API_BASE_URL ?? "").trim();

const AUTH_API_PREFIX = "/api/auth/";

export const apiConfig = {
  publicBaseUrl: publicApiBaseUrl,
  directBrowserCalls: publicApiBaseUrl.length > 0,
};

export const normalizeApiPath = (input?: string | null) => {
  if (!input || typeof input !== "string") return input ?? "";
  if (input.startsWith("http://") || input.startsWith("https://")) return input;

  return input.startsWith("/api")
    ? input
    : `/api${input.startsWith("/") ? input : `/${input}`}`;
};

export const shouldUseDirectBrowserApi = (input?: string | null) => {
  const normalizedPath = normalizeApiPath(input);
  if (!normalizedPath || normalizedPath.startsWith("http://") || normalizedPath.startsWith("https://")) {
    return false;
  }

  return !normalizedPath.startsWith(AUTH_API_PREFIX);
};

export const resolveBrowserApiUrl = (input?: string | null) => {
  const normalizedPath = normalizeApiPath(input);
  if (!normalizedPath || normalizedPath.startsWith("http://") || normalizedPath.startsWith("https://")) {
    return normalizedPath;
  }

  if (!apiConfig.directBrowserCalls || !shouldUseDirectBrowserApi(normalizedPath)) {
    return normalizedPath;
  }

  const base = apiConfig.publicBaseUrl.replace(/\/+$/, "").replace(/\/api$/, "");
  return `${base}${normalizedPath}`;
};

export const resolveApiUrl = (input?: string | null) => {
  return resolveBrowserApiUrl(input);
};
