export type ExpenseCategory = "fund" | "fun" | "future";
export type ExpenseType = "manual" | "recurring" | "checklist";

export type ExpenseSettings = {
  id: string;
  userId: string;
  income: number;
  ratioFund: number;
  ratioFun: number;
  ratioFuture: number;
  createdAt: Date;
  updatedAt: Date;
};

export type ExpenseTemplate = {
  id: string;
  userId: string;
  title: string;
  amount: number;
  category: ExpenseCategory;
  createdAt: Date;
  updatedAt: Date;
};

export type ExpenseSheet = {
  id: string;
  userId: string;
  periodKey: string;
  createdAt: Date;
  closedAt?: Date;
};

export type ChecklistItem = {
  id: string;
  sheetId: string;
  templateId?: string;
  title: string;
  amount: number;
  category: ExpenseCategory;
  createdAt: Date;
  completedAt?: Date;
  expenseId?: string;
};

export type ExpenseEntry = {
  id: string;
  sheetId: string;
  userId: string;
  title: string;
  amount: number;
  category: ExpenseCategory;
  date: string;
  type: ExpenseType;
  createdAt: Date;
};
