import { fetchApi } from "@/server/auth/fetchApi";
import type { Refill } from "../types";

type SeasonSnapshot = {
  seasonLabel: string;
  seasonToDate: number;
  lastSeasonToDate: number;
  delta: number;
  deltaPct: number | null;
};

export type HeatDashboardViewModel = {
  daysSinceRefill: number;
  seasonSnapshot: SeasonSnapshot;
  history: Refill[];
};

export const heatService = {
  getDashboardViewModel: async (
    userId: string,
    page = 1,
    limit = 20
  ): Promise<HeatDashboardViewModel> => {
    return fetchApi<HeatDashboardViewModel>(`/heat/dashboard?page=${page}&limit=${limit}`);
  }
};
