const resolveAuthApiBaseUrl = () => {
  const rawBase = process.env.API_PROXY_URL ?? import.meta.env?.PUBLIC_API_BASE_URL ?? "http://127.0.0.1:8080/api";
  return rawBase.replace(/\/+$/, "").replace(/\/api$/, "");
};

const buildAuthApiUrl = (path: string) => {
  const normalizedPath = path.startsWith("/api") ? path : `/api${path.startsWith("/") ? path : `/${path}`}`;
  return `${resolveAuthApiBaseUrl()}${normalizedPath}`;
};

const collectSetCookieHeaders = (response: Response) => {
  const maybeGetSetCookie = response.headers as Headers & {
    getSetCookie?: () => string[];
  };

  if (typeof maybeGetSetCookie.getSetCookie === "function") {
    return maybeGetSetCookie.getSetCookie();
  }

  const single = response.headers.get("set-cookie");
  return single ? [single] : [];
};

type ProxiedAuthResponse = {
  status: number;
  contentType: string | null;
  body: string | null;
  setCookieHeaders: string[];
};

export const proxyAuthRequest = async (request: Request, path: string, payload?: unknown): Promise<ProxiedAuthResponse> => {
  const headers = new Headers();
  headers.set("Content-Type", "application/json");

  const cookie = request.headers.get("cookie");
  if (cookie) {
    headers.set("Cookie", cookie);
  }

  const bodyText = request.method === "GET" || request.method === "HEAD"
    ? undefined
    : payload === undefined
      ? undefined
      : JSON.stringify(payload);

  const response = await fetch(buildAuthApiUrl(path), {
    method: request.method,
    headers,
    body: bodyText,
  });

  return {
    status: response.status,
    contentType: response.headers.get("content-type"),
    body: response.status == 204 ? null : await response.text(),
    setCookieHeaders: collectSetCookieHeaders(response),
  };
};
