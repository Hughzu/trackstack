import { AsyncLocalStorage } from "node:async_hooks";

type AuthContext = {
  userId?: string;
  sessionId?: string;
};

const authStorage = new AsyncLocalStorage<AuthContext>();

export const runWithAuthContext = <T>(context: AuthContext, callback: () => T) => {
  return authStorage.run(context, callback);
};

export const getCurrentUserId = (): string => {
  const userId = authStorage.getStore()?.userId;
  if (!userId) throw new Error("Missing authenticated user");
  return userId;
};

export const getCurrentSessionId = (): string | undefined => {
  return authStorage.getStore()?.sessionId;
};
