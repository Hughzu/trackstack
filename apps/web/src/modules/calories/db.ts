import { getDb } from "@/server/db/sqlite";

const DOMAIN = "calories";

export const getCaloriesDb = () => {
  return getDb(DOMAIN);
};
