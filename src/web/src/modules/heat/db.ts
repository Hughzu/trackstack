import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "heat";

export const getHeatDb = () => {
  return getDb(DOMAIN);
};
