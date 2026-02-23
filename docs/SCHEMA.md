# Schema

This document is the single source of truth for Trackstack database structure and mutation rules. It reflects the current Turso (LibSQL) schema and how the app accesses it.

## Overview

Trackstack uses **four separate Turso databases**, one per domain:

- `users`
- `calories`
- `expenses`
- `heat`

Each database has its own migration directory under `packages/db/migrations/` and is managed with Atlas (`packages/db/migrations/atlas.hcl`).

## Connection & Access Rules

- **All DB access goes through** `apps/web/src/server/db/sqlite.ts` via `getDb(domain)`.
- **Never instantiate a LibSQL client directly.**
- Runtime secrets are resolved from SSM when env values start with `/trackstack/`.

Domains are mapped in the modules:

- `calories` → `apps/web/src/modules/calories/db.ts`
- `expenses` → `apps/web/src/modules/expenses/db.ts`
- `heat` → `apps/web/src/modules/heat/db.ts`
- `users` → `apps/web/src/server/auth` (sessions + auth logic)

## Migration Locations

- `packages/db/migrations/atlas.hcl` defines the Atlas environments.
- Domain migrations live in:
  - `packages/db/migrations/users/`
  - `packages/db/migrations/calories/`
  - `packages/db/migrations/expenses/`
  - `packages/db/migrations/heat/`

## Migrations (Atlas + Turso)

### Prerequisites

- Atlas CLI: https://atlasgo.io
- Turso CLI: https://docs.turso.tech/quickstart
- Local SQLite CLI (for data import): `sqlite3`

### Turso setup (one DB per domain)

Create databases:

```bash
turso db create users
turso db create calories
turso db create expenses
turso db create heat
```

Fetch URLs:

```bash
turso db show users --url
turso db show calories --url
turso db show expenses --url
turso db show heat --url
```

Create tokens:

```bash
turso db tokens create users
turso db tokens create calories
turso db tokens create expenses
turso db tokens create heat
```

### Environment variables

Set these in your local `.env` (not committed) or CI secrets. For the Astro app, use `apps/web/.env`:

```bash
TURSO_USERS_URL=libsql://<users-db>.turso.io
TURSO_USERS_TOKEN=<token>

TURSO_CALORIES_URL=libsql://<calories-db>.turso.io
TURSO_CALORIES_TOKEN=<token>

TURSO_EXPENSES_URL=libsql://<expenses-db>.turso.io
TURSO_EXPENSES_TOKEN=<token>

TURSO_HEAT_URL=libsql://<heat-db>.turso.io
TURSO_HEAT_TOKEN=<token>
```

Atlas uses `packages/db/migrations/atlas.hcl` to map these env vars to each domain.

### Apply migrations to Turso

```bash
atlas migrate apply --env calories
atlas migrate apply --env expenses
atlas migrate apply --env heat
atlas migrate apply --env users
```

### Apply migrations locally (SQLite)

```bash
atlas migrate apply \
  --dir file://migrations/calories \
  --url "sqlite://<absolute-path-to>/src/data/calories.sqlite"
```

Repeat for `expenses` and `heat` with the matching migration directory.

If the local database already has tables, Atlas will consider it dirty:

```bash
atlas migrate apply \
  --dir file://migrations/calories \
  --url "sqlite://<absolute-path-to>/src/data/calories.sqlite" \
  --allow-dirty
```

### Seed user

The seed script creates a single email/password user in the `users` database.

```bash
node src/web/scripts/seed-user.mjs --email you@example.com --password yourpass
```

Database target precedence:

1) `TURSO_USERS_URL` + `TURSO_USERS_TOKEN`
2) Local SQLite file at `src/data/users.sqlite`

Notes:

- If the user already exists, the script exits without changes.
- The password is hashed using scrypt with the same parameters as the app.

## CI/CD Migrations

The workflow `.github/workflows/deploy-serverless.yml` runs Atlas migrations on pushes that change `packages/db/migrations/**`.
It applies migrations sequentially and only deploys the app if migrations succeed (or were skipped).

Current order:
- users → expenses → heat → calories

Rollback policy:
- Roll back **only** the failed domain using `atlas migrate down --step 1`.
- **Every migration must include a down script** (`*.down.sql`).

## Database Schemas

### users

Migration: `packages/db/migrations/users/001_init.sql`

**users**
- `id` TEXT PK
- `email` TEXT UNIQUE NOT NULL
- `password_hash` TEXT NOT NULL
- `session_version` INTEGER NOT NULL DEFAULT 0
- `created_at` TEXT NOT NULL
- `last_login_at` TEXT

**sessions**
- `id` TEXT PK
- `user_id` TEXT NOT NULL (FK → users.id)
- `created_at` TEXT NOT NULL
- `expires_at` TEXT NOT NULL
- `rotated_at` TEXT NOT NULL
- `last_seen_at` TEXT NOT NULL
- `absolute_expires_at` TEXT NOT NULL
- `parent_id` TEXT
- `revoked_at` TEXT
- `user_agent_hash` TEXT
- `ip_prefix` TEXT

Indexes:
- `sessions_user_id_idx`
- `sessions_expires_at_idx`
- `sessions_absolute_expires_at_idx`

### calories

Migration: `packages/db/migrations/calories/001_init.sql`

**calorie_logs**
- `id` TEXT PK
- `user_id` TEXT NOT NULL
- `date_time` TEXT NOT NULL
- `calories` INTEGER NOT NULL
- `protein_g` INTEGER NOT NULL
- `carbs_g` INTEGER
- `fat_g` INTEGER
- `title` TEXT

Index:
- `idx_calorie_logs_user_date` on (`user_id`, `date_time`)

**calorie_targets**
- `id` TEXT PK
- `user_id` TEXT NOT NULL UNIQUE
- `target_kcal` INTEGER NOT NULL
- `target_protein_g` INTEGER NOT NULL
- `target_carbs_g` INTEGER
- `target_fat_g` INTEGER
- `created_at` TEXT NOT NULL
- `updated_at` TEXT NOT NULL

Type mappings live in `apps/web/src/modules/calories/services/caloriesService.ts`.

### expenses

Migration: `packages/db/migrations/expenses/001_init.sql`

**expense_settings**
- `id` TEXT PK
- `user_id` TEXT NOT NULL UNIQUE
- `income` REAL NOT NULL
- `ratio_fund` INTEGER NOT NULL
- `ratio_fun` INTEGER NOT NULL
- `ratio_future` INTEGER NOT NULL
- `created_at` TEXT NOT NULL
- `updated_at` TEXT NOT NULL

**expense_checklist_templates**
- `id` TEXT PK
- `user_id` TEXT NOT NULL
- `title` TEXT NOT NULL
- `amount` REAL NOT NULL
- `category` TEXT NOT NULL (fund|fun|future)
- `created_at` TEXT NOT NULL
- `updated_at` TEXT NOT NULL

Index:
- `idx_expense_checklist_templates_user`

**expense_recurring_templates**
- `id` TEXT PK
- `user_id` TEXT NOT NULL
- `title` TEXT NOT NULL
- `amount` REAL NOT NULL
- `category` TEXT NOT NULL (fund|fun|future)
- `created_at` TEXT NOT NULL
- `updated_at` TEXT NOT NULL

Index:
- `idx_expense_recurring_templates_user`

**expense_sheets**
- `id` TEXT PK
- `user_id` TEXT NOT NULL
- `period_key` TEXT NOT NULL
- `created_at` TEXT NOT NULL
- `closed_at` TEXT

Index:
- `idx_expense_sheets_user_open` on (`user_id`, `closed_at`)

**expense_checklist_items**
- `id` TEXT PK
- `sheet_id` TEXT NOT NULL
- `template_id` TEXT
- `title` TEXT NOT NULL
- `amount` REAL NOT NULL
- `category` TEXT NOT NULL (fund|fun|future)
- `created_at` TEXT NOT NULL
- `completed_at` TEXT
- `expense_id` TEXT

Index:
- `idx_expense_checklist_items_sheet` on (`sheet_id`, `completed_at`)

**expense_entries**
- `id` TEXT PK
- `sheet_id` TEXT NOT NULL
- `user_id` TEXT NOT NULL
- `title` TEXT NOT NULL
- `amount` REAL NOT NULL
- `category` TEXT NOT NULL (fund|fun|future)
- `date` TEXT NOT NULL
- `type` TEXT NOT NULL (manual|recurring|checklist)
- `created_at` TEXT NOT NULL

Index:
- `idx_expense_entries_sheet_date` on (`sheet_id`, `date`)

Type mappings live in `apps/web/src/modules/expenses/types.ts` and service logic in `apps/web/src/modules/expenses/services/expensesService.ts`.

### heat

Migrations:
- `packages/db/migrations/heat/001_init.sql`
- `packages/db/migrations/heat/002_backfill_season.sql`

**refills**
- `id` TEXT PK
- `user_id` TEXT NOT NULL
- `date` TEXT NOT NULL
- `weight_kg` REAL NOT NULL
- `bags` INTEGER NOT NULL
- `temperature` REAL
- `season` TEXT

Index:
- `idx_refills_user_date` on (`user_id`, `date`)

Type mappings live in `apps/web/src/modules/heat/types.ts` and service logic in `apps/web/src/modules/heat/services/heatService.ts`.

## Mutation Rules

- Database mutations **must** occur in API routes under `apps/web/src/pages/api/`.
- UI components and pages must call domain services, not raw SQL.
- All DB access must go through `getDb(domain)` (`apps/web/src/server/db/sqlite.ts`).
- Migrations are the only source of truth for schema changes; update SQL + this doc together.
