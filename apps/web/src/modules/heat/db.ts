import { getDb } from "@/server/db/sqlite";

const DOMAIN = "heat";

export const getHeatDb = () => {
  return getDb(DOMAIN);
};
