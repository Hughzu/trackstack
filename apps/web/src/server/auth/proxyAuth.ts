const resolveAuthApiBaseUrl = () => {
  const rawBase = process.env.API_PROXY_URL ?? import.meta.env?.PUBLIC_API_BASE_URL ?? "http://127.0.0.1:8080/api";
  return rawBase.replace(/\/+$/, "").replace(/\/api$/, "");
};

const buildAuthApiUrl = (path: string) => {
  const normalizedPath = path.startsWith("/api") ? path : `/api${path.startsWith("/") ? path : `/${path}`}`;
  return `${resolveAuthApiBaseUrl()}${normalizedPath}`;
};

export const proxyAuthRequest = async (request: Request, path: string) => {
  const headers = new Headers();
  const contentType = request.headers.get("content-type");
  if (contentType) {
    headers.set("Content-Type", contentType);
  }

  const cookie = request.headers.get("cookie");
  if (cookie) {
    headers.set("Cookie", cookie);
  }

  const bodyText = request.method === "GET" || request.method === "HEAD" ? undefined : await request.text();

  const response = await fetch(buildAuthApiUrl(path), {
    method: request.method,
    headers,
    body: bodyText,
  });

  const outgoingHeaders = new Headers();
  const setCookie = response.headers.get("set-cookie");
  if (setCookie) {
    outgoingHeaders.set("set-cookie", setCookie);
  }

  const location = response.headers.get("location");
  if (location) {
    outgoingHeaders.set("location", location);
  }

  const body = response.status === 204 || response.status === 303 ? null : await response.text();
  const contentTypeOut = response.headers.get("content-type");
  if (contentTypeOut) {
    outgoingHeaders.set("content-type", contentTypeOut);
  }

  return new Response(body, {
    status: response.status,
    headers: outgoingHeaders,
  });
};
