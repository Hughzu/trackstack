# Testing Guide for Trackstack

This document describes the testing strategy and setup for the Trackstack application.

## Overview

Our testing strategy focuses on **fast local feedback** while preventing regressions, especially for the critical SigV4 form submission flow. Tests are intentionally lightweight while the Astro SSR frontend transitions to SSG and the backend moves to Go.

## Test Structure

```
apps/web/
├── tests/
│   ├── forms.test.ts         # Unit tests: form attribute validation
├── vitest.config.ts          # Unit test configuration
```

## Quick Start

### Installation

```bash
cd apps/web
pnpm install
```

### Running Tests

```bash
# Fast unit tests (< 300ms)
pnpm test           # Run once
pnpm test:watch     # Watch mode for development
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
  const content = getSourceContent('pages/calories/new.astro');
  expect(content).toContain('data-api-form');
  expect(content).toContain('action="/api/calories/log"');
});
```

**When to run:** Before every commit. These are fast and catch missing attributes.

## Test Strategy by Phase

### Current Phase (Astro SSR)

Focus on **boundary testing**:
1. Unit tests verify form attributes (fast feedback)
2. Skip testing business logic deeply (being rewritten in Go)

### Future Phase (Go Backend + SSG)

When migrating:
1. Add contract tests between Astro SSG and Go API
2. Heavy testing shifts to Go backend unit tests
3. Astro tests become minimal (build verification only)

## Adding New Tests

### For a New Form

1. **Add unit test** in `tests/forms.test.ts`:
```typescript
test('new-domain/new form has required data attributes', () => {
  const content = getSourceContent('pages/new-domain/new.astro');
  expect(content).toContain('data-api-form');
  expect(content).toContain('action="/api/new-domain/endpoint"');
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
```

## Troubleshooting

### Unit tests fail on Windows
Paths use forward slashes. Tests use `path.resolve()` for cross-platform compatibility.

## Best Practices

1. **Run unit tests frequently** - they're sub-second
2. **Don't test Astro internals** - test the contracts (SigV4 form attributes)
3. **Use `test.skip()` for flaky tests** - fix them before merging
4. **Keep credentials out of the repo** - never commit real credentials
