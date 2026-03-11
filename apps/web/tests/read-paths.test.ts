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
});
