import type { MiddlewareHandler } from "astro";
import { authConfig } from "@/server/auth/config";
import { runWithAuthContext } from "@/server/auth/currentUser";
import { verifySession } from "@/server/auth/verifySession";

const isAllowedApiPath = (pathname: string) => {
  if (pathname === "/api/health") return true;
  if (pathname.startsWith("/api/auth")) return true;
  return false;
};

const isPublicPage = (pathname: string) => {
  if (pathname === "/login") return true;
  return false;
};

const originHeaderName = process.env.ORIGIN_VERIFY_HEADER;
const originHeaderValue = process.env.ORIGIN_VERIFY_VALUE;
const shouldEnforceOrigin =
  !import.meta.env.DEV && !!originHeaderName && !!originHeaderValue;

export const onRequest: MiddlewareHandler = async (context, next) => {
  const url = new URL(context.request.url);
  const pathname = url.pathname;
  const isApiRequest = pathname.startsWith("/api/");

  if (shouldEnforceOrigin) {
    const provided = context.request.headers.get(originHeaderName);
    if (provided !== originHeaderValue) {
      return new Response("Forbidden", { status: 403 });
    }
  }

  if (!isApiRequest && isPublicPage(pathname)) {
    return next();
  }

  if (isApiRequest && isAllowedApiPath(pathname)) {
    return next();
  }

  const rawToken = context.cookies.get(authConfig.cookie.name)?.value;

  let authContext: { userId?: string; sessionId?: string; rawToken?: string } = {};

  if (rawToken) {
    const verification = await verifySession(rawToken);
    if (verification.authenticated) {
      authContext = {
        userId: verification.userId,
        sessionId: verification.sessionId,
        rawToken: verification.rawToken,
      };

      for (const cookieHeader of verification.setCookieHeaders) {
        context.response.headers.append("set-cookie", cookieHeader);
      }
    }
  }

  if (authContext.userId) {
    context.locals.userId = authContext.userId;
    context.locals.sessionId = authContext.sessionId;
  }

  if (isApiRequest && !authContext.userId) {
    return new Response("Unauthorized", { status: 401 });
  }

  if (!isApiRequest && !authContext.userId) {
    return context.redirect("/login");
  }

  return runWithAuthContext(authContext, () => next());
};
