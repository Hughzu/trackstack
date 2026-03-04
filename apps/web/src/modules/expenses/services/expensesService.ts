import { fetchApi } from "@/server/auth/fetchApi";
import type {
  ChecklistItem,
  ExpenseCategory,
  ExpenseEntry,
  ExpenseSettings,
  ExpenseSheet,
  ExpenseTemplate,
  ExpenseType
} from "../types";

export type ExpenseDashboardViewModel = {
  periodKey: string;
  balance: {
    remaining: number;
    income: number;
  };
  spent: {
    fund: number;
    fun: number;
    future: number;
  };
  budget: {
    fund: number;
    fun: number;
    future: number;
  };
  ratios: Array<{
    percent: number;
    color: string;
    label: string;
    value: number;
    budget: number;
    target: number;
    over: boolean;
  }>;
  pendingObligations: ChecklistItem[];
  history: ExpenseEntry[];
};

type ViewSettingsResponse = {
  settings: ExpenseSettings;
  checklist: ExpenseTemplate[];
  recurring: ExpenseTemplate[];
};

export const expensesService = {
  getDashboardViewModel: async (userId: string): Promise<ExpenseDashboardViewModel> => {
    return fetchApi<ExpenseDashboardViewModel>("/expenses/sheet/current");
  },

  getSettings: async (userId: string): Promise<ExpenseSettings> => {
    const res = await fetchApi<ViewSettingsResponse>("/expenses/settings");
    return res.settings;
  },

  getChecklistTemplates: async (userId: string): Promise<ExpenseTemplate[]> => {
    const res = await fetchApi<ViewSettingsResponse>("/expenses/settings");
    return res.checklist ?? [];
  },

  getRecurringTemplates: async (userId: string): Promise<ExpenseTemplate[]> => {
    const res = await fetchApi<ViewSettingsResponse>("/expenses/settings");
    return res.recurring ?? [];
  }
};
