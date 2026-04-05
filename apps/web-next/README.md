## TrackStack Frontend Migration

Solid.js SPA scaffold for the Astro -> Solid migration described in `docs/FRONT_MIGRATION.md`.

## Commands

```bash
pnpm dev
pnpm openapi:generate
pnpm build
```

- `pnpm openapi:generate`: regenerates the typed client schema from `apps/server/internal/app/monolithapi/openapi.yaml` into `src/core/api/generated/schema.ts`; run it whenever the backend contract changes.

## Structure

- `src/core/` shared config, auth helpers, and typed API client
- `src/components/ui/` design-system building blocks
- `src/features/` domain entry pages and domain API wrappers
- `src/styles/global.css` global theme tokens and base styles

## Env

- `VITE_DEPLOY_TARGET`: `serverless`, `vps`, `k8s`, `ecs`, or `eks`
- `VITE_API_BASE_URL`: optional direct API origin for static deployments
