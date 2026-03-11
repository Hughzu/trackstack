import { fetchApi } from "@/server/auth/fetchApi";

export type DashboardViewModel = {
  expenses: {
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
      over: boolean;
    }>;
  };
  calories: {
    summary: {
      consumed: number;
      protein: number;
      carbs: number;
      fat: number;
      target: {
        targetKcal: number;
        targetProteinG: number;
      };
    };
  };
  heat: {
    daysSinceRefill: number;
    season: {
      label: string;
      current: number;
      last: number;
    };
  };
};

export const dashboardService = {
  getDashboardViewModel: async (userId: string): Promise<DashboardViewModel> => {
    return fetchApi<DashboardViewModel>("/dashboard");
  }
};
