import { getDb } from "@/core/db/sqlite";

const DOMAIN = "calories";

export const getCaloriesDb = () => {
  return getDb(DOMAIN);
};
