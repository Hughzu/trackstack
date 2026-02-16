import { getDb } from "@/core/db/sqlite";

const DOMAIN = "heat";

export const getHeatDb = () => {
  return getDb(DOMAIN);
};
