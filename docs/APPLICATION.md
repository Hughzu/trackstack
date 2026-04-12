# Application

This is the practical guide for working in the current app.

If you are new here, pair this with `docs/ARCHITECTURE.md` and `docs/MASTERPLAN.md`.

## Current Truth

- Active frontend: `apps/web`
- Active backend consumed by the frontend: `apps/server`
- Browser contract source: `apps/server/internal/app/monolithapi/openapi.yaml`
- Generated frontend types: `apps/web/src/core/api/generated/schema.ts`


## Frontend Layout

```text
apps/web/src/
├── components/ui/     # shared UI primitives and composed surfaces
├── core/              # auth, API client, config, formatting
├── features/
│   ├── auth/
│   ├── dashboard/
│   ├── calories/
│   ├── expenses/
│   └── heat/
├── routes.tsx         # lazy route table + guards
└── index.tsx          # SPA mount, theme apply, auth bootstrap
```

## Frontend Operating Model

### `src/core/`

Shared runtime behavior lives here:

- API client and fetch policy
- auth state, bootstrap, refresh, token storage, route guards
- deploy-target config and theme selection
- shared formatters and generated contract types

Important files:

- `apps/web/src/core/api/client.ts`
- `apps/web/src/core/auth/store.ts`
- `apps/web/src/core/auth/guards.tsx`
- `apps/web/src/core/auth/token.ts`
- `apps/web/src/core/config/api.ts`
- `apps/web/src/core/config/theme.ts`

### `src/components/ui/`

This is the shared UI library. It contains primitives and repeatable interaction patterns such as:

- shell and tabs
- panels and list surfaces
- form layout primitives
- buttons, pills, notices, stats, skeletons
- confirm/delete sheets

Important files:

- `apps/web/src/components/ui/AppRoot.tsx`
- `apps/web/src/components/ui/AppShell.tsx`
- `apps/web/src/components/ui/Panel.tsx`
- `apps/web/src/components/ui/Form.tsx`
- `apps/web/src/components/ui/ActionButton.tsx`
- `apps/web/src/components/ui/ConfirmSheet.tsx`

### `src/features/<domain>/`

Each domain owns its route entry files, local components, local view-model mapping, and domain-specific API wrappers.

Current feature layout:

- `apps/web/src/features/auth`
- `apps/web/src/features/dashboard`
- `apps/web/src/features/calories`
- `apps/web/src/features/expenses`
- `apps/web/src/features/heat`

Typical feature shape:

```text
features/<domain>/
├── api/         # typed wrappers over openapi-fetch
├── components/  # domain-specific UI sections
├── display.ts   # view-model and formatting helpers
├── index.tsx    # main route
├── new.tsx      # create route when needed
└── settings.tsx # settings route when needed
```

## Routing And App Shell

- `apps/web/src/routes.tsx` lazy-loads all routes.
- Route access is enforced through `ProtectedRoute` and `PublicOnlyRoute`.
- `apps/web/src/components/ui/AppRoot.tsx` keeps the `AppShell` mounted while route content swaps.
- `apps/web/src/index.tsx` applies the deploy theme and runs `bootstrapAuth()` on mount.

That shell-persistence choice matters. It keeps navigation from flashing the whole app like a cheap carnival ride.

## Auth And API Transport

The frontend auth model is:

- access token stored in `sessionStorage`
- refresh token stored in an `HttpOnly` cookie
- browser auth bootstraps once on SPA start
- protected requests retry once after `401` by calling refresh
- explicit logout drops the access token and writes a logout marker so the app does not silently rehydrate itself from a still-valid cookie in the same session

Important files:

- `apps/web/src/core/auth/store.ts`
- `apps/web/src/core/auth/token.ts`
- `apps/web/src/core/auth/refresh.ts`
- `apps/web/src/core/auth/transport.ts`
- `apps/web/src/core/api/client.ts`

Deployment-specific transport rule:

- The browser sends bearer auth in `X-Trackstack-Authorization`, not `Authorization`, because CloudFront needs `Authorization` for SigV4 signing to the Lambda Function URL origin.
- Mutating requests may need `x-amz-content-sha256`; `apps/web/src/core/api/payload-hash.ts` and `apps/web/src/core/api/client.ts` handle that.

## Design System Reality

The shared UI library is real and heavily used.

But here is the honest version: boundary enforcement is still convention-based.

What is true today:

- shared patterns like buttons, panels, shell, list surfaces, notices, confirm sheets, and form layout already live in `src/components/ui/`
- features compose those primitives instead of rebuilding the whole world every time
- feature-local components still contain some raw Tailwind composition when the shared seam is not abstracted enough yet

So the rule is not "all styling only in `ui`" because the codebase does not actually honor that yet. The real rule is:

- repeated interaction patterns belong in `src/components/ui/`
- one-off domain rendering can stay inside feature components until repetition becomes obvious

## Current Feature Boundary Rules

These are the rules you should follow now, because they reflect the repo as it exists and the direction it should keep:

- Features may own domain-specific API wrappers under `features/<domain>/api/`.
- Features may own domain-specific mapping helpers under `features/<domain>/display.ts`.
- Shared logic used across multiple domains should move to `src/core/`.
- Shared UI patterns used across multiple domains should move to `src/components/ui/`.
- Do not reach into another feature's random internals from a sibling feature.

Current exception you should know about:

- `apps/web/src/features/dashboard/index.tsx` imports sibling feature API and display helpers directly.

Treat that as tolerated technical debt, not the pattern to copy.

If a cross-feature seam is needed, make it deliberate. A small shared module is better than secret tunnel imports.

## Route Thickness Rules

Routes are controllers, but not every current route is as thin as the ideal.

What route files should do:

- gate on auth readiness
- load route-level resources
- manage route-level mutation state
- call feature API helpers
- compose extracted sections and shared UI primitives

What route files should not turn into:

- giant mixed files containing fetch logic, validation, multiple forms, modal flows, and bespoke rendering all at once

Reality check:

- `apps/web/src/features/expenses/settings.tsx` is still heavy and is the best example of where extraction should continue.
- `apps/web/src/features/dashboard/index.tsx` still contains card implementations inline.

If a route starts to feel like a junk drawer, that feeling is usually correct.

## Data Loading Rules

The current pattern is straightforward:

- wait for auth state to become authenticated
- create route resources with `createResource`
- refetch after successful mutations
- prefer independent resource loading when one slow domain should not block another

Examples:

- dashboard loads expenses, calories, and heat independently
- calories, expenses, and heat route pages refetch their own dashboard resource after mutations

This is simple and maintainable. It is not the most optimized thing on earth, but for a side project it is a good default until perf data says otherwise.

## Theme And Deploy Targets

The same SPA can brand itself by deploy target.

- Theme selection lives in `apps/web/src/core/config/theme.ts`.
- Current targets: `serverless`, `vps`, `k8s`.
- `ecs`, `eks`, `kubernetes`, and `lambda` are normalized aliases.

This is a branding/runtime skin seam, not a different app.

## API Client Rules

Use the generated contract, not handwritten fantasy types.

- Source OpenAPI: `apps/server/internal/app/monolithapi/openapi.yaml`
- Generated schema: `apps/web/src/core/api/generated/schema.ts`
- Friendly aliases: `apps/web/src/core/api/types.ts`
- Runtime client: `apps/web/src/core/api/client.ts`

Feature API wrappers should stay thin. They should call `apiClient`, unwrap the response, and return typed data.

## What A New Coworker Know First

Before editing anything, understand these seams:

1. `apps/web` is the active frontend.
2. `apps/server` owns the backend contract.
3. Auth is bearer access token + refresh cookie, with SPA bootstrap and one refresh retry.
4. The CloudFront/Lambda deployment uses `X-Trackstack-Authorization` and payload hashing rules.
5. Repeated UI patterns belong in `src/components/ui/`.
6. Cross-feature imports are mostly a smell unless the seam is intentional.
7. The dashboard currently cheats a bit by importing sibling feature clients; do not multiply that pattern.

## Current Gaps And Debt

- Feature boundary rules are not enforced by tooling yet.
- Some route files are thicker than they should be.
- Some feature components still contain raw Tailwind composition that should eventually be promoted into better shared primitives.
- `Router preload={true}` may be worth revisiting if route chunk eagerness starts fighting your performance goals.

## Safe Workflow For Changes

If you are adding or changing a frontend feature:

1. Start from the Go contract or update it first.
2. Regenerate frontend types if the contract changed.
3. Keep route files focused on orchestration.
4. Extract reusable sections before the route becomes gross.
5. Promote repeated patterns into `src/components/ui/`.
6. Update this doc if you changed a boundary, workflow, or rule another engineer would need.
