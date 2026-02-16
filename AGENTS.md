As **Trackstack's Lead Architect**, I've designed this `AGENTS.md` specifically to act as the "brain" for any AI agent (Cursor, Windsurf, GitHub Copilot) entering this codebase. It enforces our **Polylith** boundaries and the **Managed Core** philosophy we just established.

Place this file at the root of your repository.

---

# 🤖 AGENTS.md: Trackstack System Instructions

**Context:** You are an expert Software Engineer and Platform Architect assisting in the development of **Trackstack**.
**Core Philosophy:** "The Managed Polylith" — A single-binary Go application embedding an Astro frontend, deployable anywhere from Lambda to EKS.

---

## 🛠️ Technical Stack

* **Backend:** Go (Golang) 1.22+ (Hexagonal Architecture).
* **Frontend:** Astro (Static Output) + TailwindCSS.
* **Database:** Turso (LibSQL over HTTP).
* **Infra:** Terraform (AWS) managed via the `tstack` CLI.
* **Deployment:** Docker (Universal Image with AWS Lambda Web Adapter).

---

## 🏗️ Architectural Rules (The Law)

### 1. The Core vs. Module Boundary

* **`internal/core/`**: This is the "Managed Layer". AI agents should **never** modify files here unless explicitly asked to patch the framework itself.
* **`internal/modules/`**: This is the "Business Layer". New features (Calories, Heat) go here.
* **Dependency Rule**: Modules can import `core`. Core **MUST NEVER** import Modules.

### 2. The Polylith Pattern

* All business logic must be encapsulated in its own component folder within `internal/modules/`.
* Communication between modules must happen via **Interfaces** or **Events** defined in `internal/core/common`.

### 3. Frontend Embedding

* The frontend is a static Astro build located in `src/web`.
* The Go backend embeds the `dist` folder using `go:embed`.
* **Rule**: Always ensure `pnpm build` is run before compiling the Go binary to update the embedded UI.

---

## 📁 Directory Map

```text
.
├── cmd/                # Artifact Entrypoints (Assembly)
│   └── monolith/       # The Universal Binary (Imports all modules)
├── internal/           # Backend Logic
│   ├── core/           # ⛔ MANAGED: Auth, DB Setup, Server Plumbing
│   └── modules/        # ✅ PRODUCT: Business Logic (Hexagonal)
├── src/web/            # 🎨 Frontend (Astro)
│   ├── src/core/       # ⛔ MANAGED UI: Layouts, Design System
│   └── src/modules/    # ✅ PRODUCT UI: Feature Pages
├── iac/                # Infrastructure as Code
│   ├── _vendor/        # ⛔ MANAGED: Base TF Modules
│   └── environments/   # ✅ USER: Project-specific config
└── tstack.yaml         # Project metadata for the CLI

```

---

## 🚀 Development Workflows

### Adding a New Feature

1. Create a folder in `internal/modules/<feature_name>`.
2. Implement the **Domain**, **Port**, and **Adapter** (Hexagonal).
3. Register the module's routes in `cmd/monolith/main.go`.
4. Create corresponding UI in `src/web/src/modules/<feature_name>`.

### FinOps Guardrails

* Assume all AWS infrastructure is **Stateless**.
* Step 1 (Lambda) is the production target.
* Steps 2 (Fargate) & 3 (EKS) are **Ephemeral Labs**. Always suggest `terraform destroy` logic for these.

---

## 📝 Coding Standards for AI

* **Go**: Prefer standard library. Use `slog` for structured logging. Use `chi` for routing.
* **Astro**: Use Tailwind utility classes. Favor server-side components where possible, but remember the output is **Static** (Client-side fetching for dynamic data).
* **Error Handling**: Wrap errors with context. No `panic()`.
* **Documentation**: Every new Module must include a `README.md` explaining its domain.

---

## 🎯 Current Goal

> **Phase 0: The Prototype.** End to finish the SSR base feature wise.
