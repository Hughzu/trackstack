# Story 1.2: Persistent User Session & SQLite Foundation

Status: ready-for-dev

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

- [ ] Initialize SQLite database with CGO enabled (AC: 1)
  - [ ] Set up `internal/common/db` with `sql.DB` connection pool
  - [ ] Implement basic migrations or schema initialization for `users` table
- [ ] Implement Session Middleware (AC: 2, 3)
  - [ ] Create a middleware that checks for a session cookie
  - [ ] Issue a persistent cookie if none exists (Auto-provision for MVP)
- [ ] User Context Injection (AC: 4)
  - [ ] Inject user object/ID into request context
  - [ ] Create helper to retrieve user from context

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

### Debug Log References

### Completion Notes List

### File List
