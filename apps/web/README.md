## TrackStack Frontend Migration

Solid.js SPA scaffold for the Astro -> Solid migration described in `docs/FRONT_MIGRATION.md`.

## Commands

```bash
pnpm dev
pnpm openapi:generate
pnpm build
pnpm test:e2e
```

- `pnpm openapi:generate`: regenerates the typed client schema from `apps/server/internal/app/monolithapi/openapi.yaml` into `src/core/api/generated/schema.ts`; run it whenever the backend contract changes.
- `pnpm test:e2e`: runs the Solid auth smoke flow against a running frontend/backend pair. It expects `E2E_TEST_EMAIL` and `E2E_TEST_PASSWORD` in `apps/web/.env.local` or your shell.

## Structure

- `src/core/` shared config, auth helpers, and typed API client
- `src/components/ui/` design-system building blocks
- `src/features/` domain entry pages and domain API wrappers
- `src/styles/global.css` global theme tokens and base styles

## Guardrails

- Route files under `src/features/**/{index,new,settings}.tsx` are controllers only: load data, own route-level mutation state, and compose extracted sections.
- Tailwind utility composition belongs in `src/components/ui/`; feature files should avoid ad-hoc styling and focus on domain mapping/wiring.
- Do not import sibling feature internals just because it is convenient. Shared behavior moves to `src/core/` or `src/components/ui/`.
- Backend contract changes should follow this order: update Go -> regenerate OpenAPI types -> update API wrappers -> update Playwright coverage.

## Env

- `VITE_DEPLOY_TARGET`: `serverless`, `vps`, `k8s`, `ecs`, or `eks`
- `VITE_API_BASE_URL`: optional direct API origin for static deployments
