import { readFileSync } from 'node:fs';
import { existsSync } from 'node:fs';
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

  test('expenses settings page no longer depends on SSR auth context', () => {
    const expensesSettings = readSource('pages/expenses/settings.astro');
    expect(expensesSettings).toContain('ExpensesSettingsClient');
    expect(expensesSettings).not.toContain('getCurrentUserId');
    expect(expensesSettings).not.toContain('expensesService');
  });

  test('expenses settings client fetches browser data from Go', () => {
    const expensesSettingsClient = readSource('modules/expenses/components/ExpensesSettingsClient.astro');
    expect(expensesSettingsClient).toContain('resolveBrowserApiUrl("/api/expenses/settings")');
    expect(expensesSettingsClient).toContain('waitForAuthReady');
  });

  test('astro middleware is removed and auth bootstrap owns page guarding', () => {
    expect(existsSync(resolve(srcDir, 'middleware.ts'))).toBe(false);

    const authBootstrap = readSource('layouts/AuthBootstrap.astro');
    expect(authBootstrap).toContain('bootstrapAuthMode === "required"');
    expect(authBootstrap).toContain('window.location.replace("/login")');
  });

  test('astro auth adapter routes are removed', () => {
    expect(existsSync(resolve(srcDir, 'pages/api/auth/login.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'pages/api/auth/logout.ts'))).toBe(false);
  });

  test('legacy SSR auth helper files are removed', () => {
    expect(existsSync(resolve(srcDir, 'server/auth/currentUser.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'server/auth/fetchApi.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'server/auth/verifySession.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'server/auth/config.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'server/auth/client.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'server/http/redirects.ts'))).toBe(false);
  });

  test('legacy SSR service wrappers are removed', () => {
    expect(existsSync(resolve(srcDir, 'modules/dashboard/services/dashboardService.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'modules/calories/services/caloriesService.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'modules/expenses/services/expensesService.ts'))).toBe(false);
    expect(existsSync(resolve(srcDir, 'modules/heat/services/heatService.ts'))).toBe(false);
  });

  test('astro frontend config uses static output', () => {
    const astroConfig = readFileSync(resolve(srcDir, '../astro.config.mjs'), 'utf-8');
    expect(astroConfig).toContain('output: "static"');
    expect(astroConfig).not.toContain('@astro-aws/adapter');
  });
});
