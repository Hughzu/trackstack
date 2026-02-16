# Database Migrations (Atlas + Turso)

Each domain database has its own migration directory and Turso database.
Current domains are derived from `src/data/*.sqlite`:

- `users`
- `calories`
- `expenses`
- `heat`

## Prerequisites

- Install Atlas CLI: https://atlasgo.io
- Install Turso CLI: https://docs.turso.tech/quickstart
- Local SQLite CLI (for data import): `sqlite3`

## Turso setup (one DB per domain)

Create databases (names match domain):

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

Create tokens (one per DB):

```bash
turso db tokens create users
turso db tokens create calories
turso db tokens create expenses
turso db tokens create heat
```

## Environment variables

Set these in your local `.env` (not committed) or CI secrets.
For the Astro app, use `src/web/.env`:

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

Atlas uses `atlas.hcl` in the repo root to wire these env vars to each
domain migration directory.

## Apply migrations to Turso

```bash
atlas migrate apply --env calories
atlas migrate apply --env expenses
atlas migrate apply --env heat
atlas migrate apply --env users
```

## Apply migrations locally (SQLite)

Use Atlas with your local SQLite file URL:

```bash
atlas migrate apply \
  --dir file://migrations/calories \
  --url "sqlite://<absolute-path-to>/src/data/calories.sqlite"
```

Repeat for `expenses` and `heat` with the matching migration directory.
Refer to Atlas URL docs for SQLite connection formats:
https://atlasgo.io/concepts/url

If the local database already has tables, Atlas will consider it "dirty".
You can either baseline or allow dirty once:

```bash
atlas migrate apply \
  --dir file://migrations/calories \
  --url "sqlite://<absolute-path-to>/src/data/calories.sqlite" \
  --allow-dirty
```

## Notes

- The app no longer initializes schema at runtime; run migrations first.
- Importing small datasets can be done via WebStorm using the libSQL JDBC
  driver and the Turso URL + token.

## Data migration (local SQLite -> Turso)

For small datasets, the CLI route is quick and repeatable. This is how the
initial data migration was done:

```bash
# calories
tables=$(sqlite3 "src/data/calories.sqlite" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'atlas_%' AND name NOT LIKE 'sqlite_%';" | tr '\n' ' ')
sqlite3 "src/data/calories.sqlite" ".dump --data-only $tables" | turso db shell calories

# expenses
tables=$(sqlite3 "src/data/expenses.sqlite" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'atlas_%' AND name NOT LIKE 'sqlite_%';" | tr '\n' ' ')
sqlite3 "src/data/expenses.sqlite" ".dump --data-only $tables" | turso db shell expenses

# heat
tables=$(sqlite3 "src/data/heat.sqlite" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'atlas_%' AND name NOT LIKE 'sqlite_%';" | tr '\n' ' ')
sqlite3 "src/data/heat.sqlite" ".dump --data-only $tables" | turso db shell heat
```

Warnings:

- Re-running the import will duplicate rows (no upsert). Clear tables first if needed.
- The import excludes `atlas_%` and `sqlite_%` tables by design.

## Golden path notes

If/when the CLI is implemented, the golden path can automate:

1) `turso db create <domain>`
2) `turso db show <domain> --url`
3) `turso db tokens create <domain>`
4) Write `.env` entries for each domain
5) `atlas migrate apply --env <domain>`
6) Optional: data import via the CLI flow above
