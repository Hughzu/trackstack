import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "calories";

export const getCaloriesDb = () => {
  return getDb(DOMAIN);
};
