# Testing Guide for Trackstack

This document describes the testing strategy and setup for the Trackstack application.

## Overview

Our testing strategy focuses on **fast local feedback** while preventing regressions, especially for the critical SigV4 form submission flow. Tests are designed to be lightweight since the Astro SSR backend will be rewritten in Go.

## Test Structure

```
apps/web/
├── tests/
│   ├── forms.test.ts         # Unit tests: form attribute validation
│   └── e2e/
│       └── sigv4.spec.ts     # E2E tests: SigV4 submission flow
├── vitest.config.ts          # Unit test configuration
└── playwright.config.ts      # E2E test configuration
```

## Quick Start

### Installation

Playwright browsers auto-install on `pnpm install`:

```bash
cd apps/web
pnpm install
```

### Running Tests

```bash
# Terminal 1: Start dev server (required for E2E)
pnpm dev

# Terminal 2: Fast unit tests (< 300ms)
pnpm test           # Run once
pnpm test:watch     # Watch mode for development

# Terminal 3: E2E tests (requires dev server)
pnpm test:e2e       # Headless mode
pnpm test:e2e:ui    # Interactive UI mode
```

## Unit Tests (Vitest)

**Purpose:** Verify forms have correct SigV4 attributes without rendering Astro components.

**What they check:**
- ✅ `data-api-form` attribute present
- ✅ Correct `action` URLs
- ✅ Proper `method` attributes
- ✅ `data-redirect` for success navigation
- ✅ No inline scripts bypassing ApiFormHandler

**Example:**
```typescript
// tests/forms.test.ts
test('calories/new form has required data attributes', () => {
  const content = getPageContent('calories/new.astro');
  expect(content).toContain('data-api-form');
  expect(content).toContain('action="/api/calories/log"');
});
```

**When to run:** Before every commit. These are fast and catch missing attributes.

## E2E Tests (Playwright)

**Purpose:** Verify SigV4 signing actually works in the browser.

**What they check:**
- ✅ Login flow creates session
- ✅ Forms submit with `x-amz-content-sha256` header
- ✅ SHA-256 hash is valid (64 hex characters)
- ✅ DELETE requests are also signed
- ✅ Error handling works correctly
- ✅ Auth redirects function properly

**Configuration:**

1. Add test credentials to `apps/web/.env`:
```bash
TEST_USER_EMAIL=your-test@email.com
TEST_USER_PASSWORD=your-password
```

2. Create a test user in your Turso database if needed.

**When to run:** Before pushing to main, or when modifying:
- `ApiFormHandler.astro`
- Form components
- Authentication flow
- API routes

## Test Strategy by Phase

### Current Phase (Astro SSR)

Focus on **boundary testing**:
1. Unit tests verify form attributes (fast feedback)
2. E2E tests verify SigV4 integration (regression prevention)
3. Skip testing business logic deeply (being rewritten in Go)

### Future Phase (Go Backend + SSG)

When migrating:
1. Port E2E tests to call Go API endpoints
2. Add contract tests between Astro SSG and Go API
3. Heavy testing shifts to Go backend unit tests
4. Astro tests become minimal (build verification only)

## Adding New Tests

### For a New Form

1. **Add unit test** in `tests/forms.test.ts`:
```typescript
test('new-domain/new form has required data attributes', () => {
  const content = getPageContent('new-domain/new.astro');
  expect(content).toContain('data-api-form');
  expect(content).toContain('action="/api/new-domain/endpoint"');
});
```

2. **Add E2E test** in `tests/e2e/sigv4.spec.ts`:
```typescript
test('new domain form submits with SigV4 headers', async ({ page }) => {
  await page.goto('/new-domain/new');
  await page.fill('input[name="field"]', 'value');
  
  const requestPromise = page.waitForRequest('**/api/new-domain/endpoint');
  await page.click('button[type="submit"]');
  
  const request = await requestPromise;
  expect(request.headers()['x-amz-content-sha256']).toBeDefined();
});
```

### For New Components

Components receiving data via props don't need tests (Astro handles rendering). Focus tests on:
- Interactive behaviors (via `ClientRuntime.astro` data attributes)
- Form submissions (SigV4 contract)
- API responses (error handling)

## CI/CD Integration

```yaml
# .github/workflows/test.yml
jobs:
  test:
    steps:
      - name: Run unit tests
        run: pnpm test
        working-directory: apps/web
      
      - name: Run E2E tests
        run: pnpm test:e2e
        working-directory: apps/web
        env:
          TEST_USER_EMAIL: ${{ secrets.TEST_USER_EMAIL }}
          TEST_USER_PASSWORD: ${{ secrets.TEST_USER_PASSWORD }}
```

## Troubleshooting

### Playwright browsers not found
```bash
pnpm exec playwright install chromium
```

### Tests fail with "TEST_USER_EMAIL must be set"
Add credentials to `apps/web/.env` file (see Configuration above).

### E2E tests timeout
Ensure dev server is running on `http://localhost:4321` before running E2E tests.

### Unit tests fail on Windows
Paths use forward slashes. Tests use `path.resolve()` for cross-platform compatibility.

## Best Practices

1. **Run unit tests frequently** - they're sub-second
2. **Run E2E before major changes** - they catch integration issues
3. **Don't test Astro internals** - test the contracts (SigV4, API responses)
4. **Use `test.skip()` for flaky tests** - fix them before merging
5. **Keep credentials in `.env`** - never commit real credentials
