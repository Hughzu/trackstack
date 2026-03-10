import { authConfig } from "@/server/auth/config";

type SessionResponse = {
  userId: string;
  sessionId: string;
};

const resolveAuthApiBaseUrl = () => {
  const rawBase = process.env.API_PROXY_URL ?? import.meta.env?.PUBLIC_API_BASE_URL ?? "http://127.0.0.1:8080/api";
  return rawBase.replace(/\/+$/, "").replace(/\/api$/, "");
};

const sessionUrl = () => `${resolveAuthApiBaseUrl()}/api/auth/session`;

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

type VerifySessionResult = {
  authenticated: boolean;
  userId?: string;
  sessionId?: string;
  rawToken?: string;
  setCookieHeaders: string[];
};

export const verifySession = async (rawToken?: string): Promise<VerifySessionResult> => {
  if (!rawToken) {
    return { authenticated: false, setCookieHeaders: [] };
  }

  const response = await fetch(sessionUrl(), {
    headers: {
      Cookie: `${authConfig.cookie.name}=${rawToken}`,
    },
  });

  const setCookieHeaders = collectSetCookieHeaders(response);

  if (!response.ok) {
    return { authenticated: false, setCookieHeaders };
  }

  const payload = (await response.json()) as SessionResponse;
  const replacementCookie = setCookieHeaders.find((value) => value.startsWith(`${authConfig.cookie.name}=`));
  const replacementRawToken = replacementCookie?.match(new RegExp(`^${authConfig.cookie.name}=([^;]+)`))?.[1];

  return {
    authenticated: true,
    userId: payload.userId,
    sessionId: payload.sessionId,
    rawToken: replacementRawToken ?? rawToken,
    setCookieHeaders,
  };
};
