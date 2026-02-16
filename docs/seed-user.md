# Seed User Script

This script creates a single email/password user in the `users` database.

## Usage

```bash
node src/web/scripts/seed-user.mjs --email you@example.com --password yourpass
```

## Database Target

The script uses the following precedence:

1) `TURSO_USERS_URL` + `TURSO_USERS_TOKEN` (if set)
2) Local SQLite file at `src/data/users.sqlite`

## Notes

- If the user already exists, the script exits without changes.
- The password is hashed using scrypt with the same parameters as the app.
