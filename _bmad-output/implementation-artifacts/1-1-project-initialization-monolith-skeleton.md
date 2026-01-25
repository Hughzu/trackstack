# Story 1.1: Project Initialization & Monolith Skeleton

Status: ready-for-dev

## Story

As a Developer,
I want to initialize the repository using the Go-HTMX starter template and set up the modular monolith directory structure,
so that I have a clean foundation that supports future module isolation.

## Acceptance Criteria

1. **Given** a clean repository
2. **When** I run the initialization script and set up `internal/common`, `internal/calories`, and `cmd/monolith`
3. **Then** the project should compile and run a basic health check at `/health`
4. **And** the `golangci-lint` check should pass with zero warnings.

## Tasks / Subtasks

- [ ] Clone starter template (AC: 1)
  - [ ] `git clone https://github.com/unkemptantlin/go-htmx-template.git .`
- [ ] Refactor into Modular Monolith structure (AC: 2)
  - [ ] Create `internal/common`, `internal/calories`, `internal/ingest`
  - [ ] Create `cmd/monolith`
- [ ] Implement Basic Health Check (AC: 3)
  - [ ] Add `/health` endpoint returning `{"status": "ok"}`
- [ ] Configure Linting (AC: 4)
  - [ ] Ensure `.golangci.yml` is present and passing

## Dev Notes

- **Relevant architecture patterns:** Strict Modular Monolith, Interface-based DI.
- **Source tree components to touch:** `internal/`, `cmd/`, `go.mod`.
- **Testing standards summary:** Co-located `_test.go` files.

### Project Structure Notes

- Alignment with unified project structure (paths, modules, naming)
- Follow the "Schrödinger's Binary" setup: logic in `internal/`, entry point in `cmd/`.

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Implementation Handoff]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.1]
- [Source: _bmad-output/planning-artifacts/project-context.md#Technology Stack]

## Dev Agent Record

### Agent Model Used

opencode/gemini-3-flash

### Debug Log References

### Completion Notes List

### File List
