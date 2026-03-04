import { fetchApi } from "@/server/auth/fetchApi";

type CalorieLog = {
  id: string;
  userId: string;
  dateTime: Date;
  calories: number;
  proteinG: number;
  carbsG?: number;
  fatG?: number;
  title?: string;
};

type CalorieTarget = {
  id: string;
  userId: string;
  targetKcal: number;
  targetProteinG: number;
  targetCarbsG?: number;
  targetFatG?: number;
  createdAt: Date;
  updatedAt: Date;
};

export type CaloriesDashboardViewModel = {
  summary: {
    consumed: number;
    protein: number;
    carbs: number;
    fat: number;
    target: CalorieTarget;
  };
  logs: CalorieLog[];
  recentMeals: CalorieLog[];
};

export const caloriesService = {
  getDashboardViewModel: async (
    userId: string,
    recentLimit = 8
  ): Promise<CaloriesDashboardViewModel> => {
    return fetchApi<CaloriesDashboardViewModel>(`/calories/dashboard?recentLimit=${recentLimit}`);
  },

  getTarget: async (userId: string): Promise<CalorieTarget> => {
    return fetchApi<CalorieTarget>("/calories/target");
  }
};
