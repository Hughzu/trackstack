import { getDb } from "@/server/db/sqlite";

const DOMAIN = "expenses";

export const getExpensesDb = () => {
  return getDb(DOMAIN);
};
