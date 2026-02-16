import type { MiddlewareHandler } from "astro";
import { authDb } from "@/core/auth/db";
import { authConfig } from "@/core/auth/config";
import { getClientContext } from "@/core/auth/client";
import { evaluateSession, getCookieOptions, hashToken, rotateSession, touchSession } from "@/core/auth/session";
import { runWithAuthContext } from "@/core/auth/currentUser";

const isAllowedApiPath = (pathname: string) => {
  if (pathname === "/api/health") return true;
  if (pathname.startsWith("/api/auth")) return true;
  return false;
};

const isPublicPage = (pathname: string) => {
  if (pathname === "/login") return true;
  return false;
};

export const onRequest: MiddlewareHandler = async (context, next) => {
  const url = new URL(context.request.url);
  const pathname = url.pathname;
  const isApiRequest = pathname.startsWith("/api/");

  if (!isApiRequest && isPublicPage(pathname)) {
    return next();
  }

  if (isApiRequest && isAllowedApiPath(pathname)) {
    return next();
  }

  const rawToken = context.cookies.get(authConfig.cookie.name)?.value;
  const { userAgent, ipPrefix } = getClientContext(context.request);

  let authContext: { userId?: string; sessionId?: string } = {};

  if (rawToken) {
    const tokenId = hashToken(rawToken);
    const session = await authDb.findSessionById(tokenId);
    if (session) {
      const evaluation = evaluateSession(session);
      if (evaluation.valid) {
        authContext = { userId: session.userId, sessionId: session.id };

        if (evaluation.needsRotation) {
          const rotated = await rotateSession(session, { userAgent, ipPrefix });
          const now = new Date();
          const cookieOptions = getCookieOptions(now, rotated.expiresAt);
          context.cookies.set(authConfig.cookie.name, rotated.rawToken, cookieOptions);
          authContext.sessionId = rotated.sessionId;
        } else if (evaluation.needsTouch) {
          await touchSession(session);
        }
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
