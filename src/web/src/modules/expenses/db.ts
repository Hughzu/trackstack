import { getDb } from "@/shared/db/sqlite";

const DOMAIN = "expenses";

export const getExpensesDb = () => {
  return getDb(DOMAIN);
};
