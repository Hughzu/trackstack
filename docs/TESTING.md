# Testing Guide

The only thing that matters is whether the app actually works for a user and runs on the infrastructure.

## 1. End-to-End Testing (Playwright)

Playwright is the sole arbiter of truth for domain functionality. It catches regressions on specific features directly through the UI.

**Setup:**
You need a test user. Define it in `apps/web/.env.local`:
```env
E2E_TEST_EMAIL=your-e2e-email@test.com
E2E_TEST_PASSWORD=your-secret
```

**Run:**
```bash
cd apps/web
pnpm test:e2e
```
*Note: Make sure your Solid dev server and Go backend are running first.*

## 2. Infrastructure Validation (The Split Test)

The architecture guarantees the exact same backend logic runs locally as a monolith OR as a microservice cluster. We validate this using Compose.

Whenever a new domain is implemented, you must verify both setups boot and serve traffic seamlessly:

**Verify the Monolith:**
```bash
docker compose up --build -d
```

**Verify the Microservice Split:**
```bash
docker compose -f docker-compose.microservices.yml up --build -d
```
If the microservice split fails but the monolith works, your domain boundaries are bleeding. Fix it.

## 3. Regression Coverage

Playwright handles functional regressions. If you touch a domain, ensure these critical paths still work in the e2e suite:

- **Auth:** Guest redirects, successful login, token injection, and logout wipes.
- **Calories:** Form submission, target updates, and log deletion.
- **Expenses:** Entry creation, monthly checklist application, and template generation.
- **Heat:** Refill creation and history deletion.

If the E2E browser flows pass and both Compose stacks boot, the feature is good to ship.
