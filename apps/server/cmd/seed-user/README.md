# Seed User

Backend-owned helper for creating or updating a test user in the `users` Turso database.

## Usage

From `apps/server/`:

```bash
go run ./cmd/seed-user --email you@example.com --password yourpass
```

Or rely on environment variables:

```bash
E2E_TEST_EMAIL=you@example.com E2E_TEST_PASSWORD=yourpass go run ./cmd/seed-user
```

## Env Loading

- Loads `apps/web/.env` first so the frontend e2e workflow keeps working.
- Loads `apps/server/.env` second for backend-local configuration.
- Explicit shell environment variables still win.

## Behavior

- Requires `TURSO_USERS_URL`
- Uses `TURSO_USERS_TOKEN` when needed for remote Turso
- Creates the user if missing
- Updates the password hash in place if the user already exists
- Uses the Go auth password hashing implementation so seeded credentials match production auth
