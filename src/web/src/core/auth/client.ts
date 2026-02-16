type ClientContext = {
  userAgent?: string | null;
  ipPrefix?: string | null;
};

const getClientIp = (request: Request) => {
  const forwarded = request.headers.get("x-forwarded-for");
  if (forwarded && forwarded.trim().length > 0) return forwarded.split(",")[0]?.trim();
  const realIp = request.headers.get("x-real-ip");
  return realIp?.trim();
};

const hashIpPrefix = (ip?: string | null) => {
  if (!ip) return null;
  if (ip.includes(".")) {
    const parts = ip.split(".").slice(0, 3);
    return parts.length === 3 ? `${parts.join(".")}.0` : ip;
  }
  if (ip.includes(":")) {
    const parts = ip.split(":").slice(0, 4);
    return parts.length > 0 ? `${parts.join(":")}::` : ip;
  }
  return ip;
};

export const getClientContext = (request: Request): ClientContext => {
  const userAgent = request.headers.get("user-agent");
  const ipPrefix = hashIpPrefix(getClientIp(request));
  return { userAgent, ipPrefix };
};
