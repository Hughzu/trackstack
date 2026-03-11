import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, test } from 'vitest';

const srcDir = resolve(__dirname, '../src');

const readSource = (relativePath: string) => readFileSync(resolve(srcDir, relativePath), 'utf-8');

describe('client-loaded read paths', () => {
  test('home page no longer depends on SSR auth context', () => {
    const content = readSource('pages/index.astro');
    expect(content).toContain('DashboardOverviewClient');
    expect(content).not.toContain('getCurrentUserId');
    expect(content).not.toContain('dashboardService');
  });

  test('dashboard overview fetches dashboard in the browser', () => {
    const content = readSource('modules/dashboard/components/DashboardOverviewClient.astro');
    expect(content).toContain('resolveBrowserApiUrl("/api/dashboard")');
    expect(content).toContain('waitForAuthReady');
  });

  test('calories dashboard no longer depends on SSR auth context', () => {
    const content = readSource('pages/calories/index.astro');
    expect(content).toContain('CaloriesDashboardClient');
    expect(content).not.toContain('getCurrentUserId');
    expect(content).not.toContain('caloriesService');
  });

  test('calories dashboard fetches browser data from Go', () => {
    const content = readSource('modules/calories/components/CaloriesDashboardClient.astro');
    expect(content).toContain('resolveBrowserApiUrl("/api/calories/dashboard?recentLimit=8")');
    expect(content).toContain('waitForAuthReady');
    expect(content).toContain('trackstack:bind-api-forms');
  });

  test('expenses and heat dashboards no longer depend on SSR auth context', () => {
    const expensesContent = readSource('pages/expenses/index.astro');
    expect(expensesContent).toContain('ExpensesDashboardClient');
    expect(expensesContent).not.toContain('getCurrentUserId');
    expect(expensesContent).not.toContain('expensesService');

    const heatContent = readSource('pages/heat/index.astro');
    expect(heatContent).toContain('HeatDashboardClient');
    expect(heatContent).not.toContain('getCurrentUserId');
    expect(heatContent).not.toContain('heatService');
  });

  test('expenses and heat dashboards fetch browser data from Go', () => {
    const expensesContent = readSource('modules/expenses/components/ExpensesDashboardClient.astro');
    expect(expensesContent).toContain('resolveBrowserApiUrl("/api/expenses/sheet/current")');
    expect(expensesContent).toContain('waitForAuthReady');
    expect(expensesContent).toContain('trackstack:bind-api-forms');

    const heatContent = readSource('modules/heat/components/HeatDashboardClient.astro');
    expect(heatContent).toContain('resolveBrowserApiUrl("/api/heat/dashboard?page=1&limit=20")');
    expect(heatContent).toContain('waitForAuthReady');
  });

  test('calories settings page no longer depends on SSR auth context', () => {
    const caloriesSettings = readSource('pages/calories/settings.astro');
    expect(caloriesSettings).toContain('CaloriesSettingsClient');
    expect(caloriesSettings).not.toContain('getCurrentUserId');
    expect(caloriesSettings).not.toContain('caloriesService');
  });

  test('calories settings client fetches browser data from Go', () => {
    const caloriesSettingsClient = readSource('modules/calories/components/CaloriesSettingsClient.astro');
    expect(caloriesSettingsClient).toContain('resolveBrowserApiUrl("/api/calories/target")');
    expect(caloriesSettingsClient).toContain('waitForAuthReady');
  });
});
