import { describe, expect, test } from 'vitest';

import { normalizeApiPath, shouldUseDirectBrowserApi } from '../src/config/api';

describe('api routing config', () => {
  test('normalizes relative api paths', () => {
    expect(normalizeApiPath('/expenses/entries')).toBe('/api/expenses/entries');
    expect(normalizeApiPath('/api/expenses/entries')).toBe('/api/expenses/entries');
  });

  test('keeps auth routes on Astro during migration', () => {
    expect(shouldUseDirectBrowserApi('/api/auth/login')).toBe(false);
    expect(shouldUseDirectBrowserApi('/api/auth/logout')).toBe(false);
  });

  test('allows non-auth browser mutations to go direct to Go', () => {
    expect(shouldUseDirectBrowserApi('/api/calories/log')).toBe(true);
    expect(shouldUseDirectBrowserApi('/api/expenses/entries')).toBe(true);
    expect(shouldUseDirectBrowserApi('/api/heat/refills')).toBe(true);
  });
});
