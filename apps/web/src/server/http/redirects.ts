export const withErrorParam = (requestUrl: string, fallback: string) => {
  const url = new URL(fallback, requestUrl);
  url.searchParams.set("error", "1");
  return url.toString();
};
