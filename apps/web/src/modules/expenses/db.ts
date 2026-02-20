import { getDb } from "@/core/db/sqlite";

const DOMAIN = "expenses";

export const getExpensesDb = () => {
  return getDb(DOMAIN);
};
