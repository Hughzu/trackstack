# Story 1.2: Persistent User Session & SQLite Foundation

Status: done

## Story

As a User,
I want the system to remember who I am across sessions using a persistent cookie,
So that I don't have to log in every time I open the app.

## Acceptance Criteria

1. **Given** a running server and a defined SQLite schema for users
2. **When** I access the application
3. **Then** the system should issue/recognize a session cookie
4. **And** my user context should be available to the backend for data scoping.

## Tasks / Subtasks

- [x] Initialize SQLite database with CGO enabled (AC: 1)
  - [x] Set up `internal/common/db` with `sql.DB` connection pool
  - [x] Implement basic migrations or schema initialization for `users` table
- [x] Implement Session Middleware (AC: 2, 3)
  - [x] Create a middleware that checks for a session cookie
  - [x] Issue a persistent cookie if none exists (Auto-provision for MVP)
- [x] User Context Injection (AC: 4)
  - [x] Inject user object/ID into request context
  - [x] Create helper to retrieve user from context

## Dev Notes

- **Relevant architecture patterns:** SQLite CGO, Middleware, Context-based scoping.
- **Source tree components to touch:** `internal/common/db`, `internal/common/server`.
- **Database:** Use `modernc.org/sqlite` (pure Go) or `github.com/mattn/go-sqlite3` (CGO). Architecture says "SQLite CGO enabled" for performance.
- **Session:** For MVP, keep it simple. A secure cookie with a user ID or a simple session store in SQLite.

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Database]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.2]
- [Source: _bmad-output/planning-artifacts/project-context.md#Technology Stack]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5

### Debug Log References

N/A - Implementation completed without issues

### Completion Notes List

1. **SQLite Database Initialization**: Implemented `internal/common/db` with CGO-enabled SQLite using `github.com/mattn/go-sqlite3`
2. **Database Configuration**:
   - Enabled WAL mode for better concurrency
   - Foreign keys enabled per connection
   - Connection pool limited to 1 for SQLite single-writer constraint
   - Database file: `./data/trackstack.db`
3. **Schema Implementation**:
   - `users` table with id (UUID), created_at, last_seen_at
   - `sessions` table with id (UUID), user_id, created_at, expires_at
   - Indexes on foreign keys (sessions.user_id) and expiry (sessions.expires_at)
4. **User & Session Models**: 
   - CRUD operations with context propagation
   - Helper functions `NewUser()` and `NewSession()` with UUID generation
   - 30-day session expiry
5. **Session Middleware**:
   - Auto-provisioning: creates user + session on first visit
   - Cookie validation and session lookup
   - User context injection via `context.WithValue`
   - Updates `last_seen_at` on every request
6. **Cookie Security**:
   - HttpOnly flag (XSS protection)
   - Secure flag (production only, false for localhost)
   - SameSite=Strict (CSRF protection)
   - 30-day MaxAge
7. **Session Endpoint**: `/api/session` returns current user JSON for verification
8. **All Acceptance Criteria Met**:
   - ✅ SQLite schema created with users and sessions tables
   - ✅ Application issues persistent cookies on first access
   - ✅ System recognizes returning users via cookie
   - ✅ User context available via `GetUserFromContext(ctx)`

### File List

Created:
- `app/internal/common/db/db.go` - Database connection, WAL mode, schema initialization (~75 lines)
- `app/internal/common/db/user.go` - User/Session models and CRUD operations (~170 lines)
- `app/internal/common/server/middleware.go` - Session middleware and context helpers (~95 lines)
- `app/data/.gitkeep` - Data directory placeholder

Modified:
- `app/internal/common/server/server.go` - Added database dependency, session middleware, `/api/session` endpoint
- `app/cmd/monolith/main.go` - Database initialization and connection
- `.gitignore` - Added WAL file patterns (*.db-wal, *.db-shm)
- `app/go.mod` - Added dependencies: go-sqlite3, google/uuid
