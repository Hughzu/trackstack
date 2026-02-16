import { expensesService } from "@/modules/expenses/services/expensesService";
import { caloriesService } from "@/modules/calories/services/caloriesService";
import { heatService } from "@/modules/heat/services/heatService";

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
    consumed: number;
    target: number;
    percent: number;
    protein: number;
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
    // 1. Fetch data in parallel
    const [expenseDashboard, healthSummary, daysSinceRefill, seasonSnapshot] = await Promise.all([
      expensesService.getDashboard(userId),
      caloriesService.getTodaySummary(userId),
      heatService.getDaysSinceLastRefill(userId),
      heatService.getSeasonSnapshot(userId)
    ]);

    // 2. Calculate Expenses Logic
    const balance = {
      remaining: Math.round((expenseDashboard.settings.income - expenseDashboard.totalSpent) * 100) / 100,
      income: expenseDashboard.settings.income
    };

    const spent = {
      fund: Math.round(expenseDashboard.spentByCategory.fund * 100) / 100,
      fun: Math.round(expenseDashboard.spentByCategory.fun * 100) / 100,
      future: Math.round(expenseDashboard.spentByCategory.future * 100) / 100
    };

    const budget = {
      fund: Math.round((balance.income * expenseDashboard.settings.ratioFund) / 100),
      fun: Math.round((balance.income * expenseDashboard.settings.ratioFun) / 100),
      future: Math.round((balance.income * expenseDashboard.settings.ratioFuture) / 100)
    };

    const expenseRatios = [
      {
        percent: balance.income ? Math.round((spent.fund / balance.income) * 100) : 0,
        color: 'bg-red-500',
        label: 'Fund.',
        over: spent.fund > budget.fund
      },
      {
        percent: balance.income ? Math.round((spent.fun / balance.income) * 100) : 0,
        color: 'bg-orange-500',
        label: 'Fun',
        over: spent.fun > budget.fun
      },
      {
        percent: balance.income ? Math.round((spent.future / balance.income) * 100) : 0,
        color: 'bg-emerald-500',
        label: 'Future',
        over: spent.future > budget.future
      }
    ];

    // 3. Calculate Calories Logic
    const calorieTarget = healthSummary.target.targetKcal;
    const caloriePercent = calorieTarget
      ? Math.min(Math.round((healthSummary.consumed / calorieTarget) * 100), 100)
      : 0;

    // 4. Return View Model
    return {
      expenses: {
        balance,
        spent,
        budget,
        ratios: expenseRatios
      },
      calories: {
        consumed: healthSummary.consumed,
        target: calorieTarget,
        percent: caloriePercent,
        protein: healthSummary.protein
      },
      heat: {
        daysSinceRefill,
        season: {
          label: seasonSnapshot.seasonLabel,
          current: seasonSnapshot.seasonToDate,
          last: seasonSnapshot.lastSeasonToDate
        }
      }
    };
  }
};
