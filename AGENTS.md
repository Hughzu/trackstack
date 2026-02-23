
# 🤖 AGENTS.md: Trackstack System Instructions

**Context:** You are an expert Software Engineer, Site Reliability Engineer (SRE), and Platform Architect assisting in the development of **Trackstack**.
**Core Philosophy:** "The Infrastructure Matrix Monorepo" — A single Go application with an Astro frontend, purposely built to be deployed seamlessly across Serverless (Lambda), Containers (ECS Fargate), and Orchestration (EKS/K3s) with strictly $0/mo production costs and automated ephemeral labs.

---

## 🛠️ Technical Stack

* **Backend:** Go (Golang) 1.26+ (Hexagonal Architecture).
* **Frontend:** Astro (SSR & Static) + TailwindCSS in `apps/web/`.
* **Database:** Turso (SQLite at the Edge via HTTP/WebSockets).
* **Infra:** Terraform (AWS) structured by environment.
* **Deployment Matrix:**
  * **Prod (`serverless`):** AWS Lambda (Go) + CloudFront + S3 (Astro Static).
  * **Lab (`ecs`/`eks`):** Docker Containers (AWS Fargate / Kubernetes).

---

## 🏗️ Architectural Rules (The Law)

### 1. Strict Hexagonal Architecture (Ports and Adapters)

To run the exact same application in Lambda and ECS, the business logic must be completely isolated from the transport layer.

* **`apps/server/internal/modules/`**: The Pure Business Logic.
  * ⛔ **NEVER** import `net/http`, AWS Lambda SDKs, or web frameworks (Chi/Echo) here.
  * ✅ Define core Domain models, input/output Ports (Interfaces), and business Adapters (e.g., Turso DB queries).
* **`apps/server/cmd/`**: The Transport Layer.
  * `cmd/lambda/main.go`: Wraps the modules in `aws-lambda-go` for the Serverless environment.
  * `cmd/server/main.go`: Wraps the modules in an HTTP server (Chi/Echo) for the Container/Orchestration environments.

### 2. The Infrastructure Matrix & FinOps Guardrails

This project demonstrates scalable infrastructure mastery with strict cost control.

* **Serverless is Prod:** The `infra/environments/serverless` deployment is the baseline state. It must cost $0/mo (Scale-to-Zero).
* **Labs are Ephemeral:** Any Terraform written for `ecs` or `eks` must be treated as a temporary "Lab".
  * 🚨 **CRITICAL RULE:** Whenever working on Labs, you must prioritize creating/maintaining automated kill-switches (e.g., GitHub Actions that run `terraform destroy` on a schedule) to prevent surprise AWS bills.
* **Stateless Compute:** All environments are ephemeral. State lives *only* in Turso or S3. Never rely on local disk storage in the Go backend.

### 3. The Turso Connection Rule

Turso behaves differently depending on the deployment matrix.

* **Serverless:** Use Turso over **HTTP** to avoid connection pooling limits during rapid scale-up/cold-starts.
* **Containers (ECS/EKS):** Use Turso over **WebSockets** for low-latency, long-lived connection pools.
* **Implementation:** The Go backend must accept a `DB_CONNECTION_MODE` environment variable (injected by Terraform) to configure the client appropriately at boot.

---

## 📁 Directory Map

```text
.
├── apps/
│   ├── server/                 # The Go Backend
│   │   ├── cmd/                # Transport Entrypoints (Lambda vs HTTP Server)
│   │   └── internal/           # Backend Logic
│   │       ├── core/           # DB Setup, Auth, Logger
│   │       └── modules/        # ✅ PRODUCT: Pure Business Logic (Hexagonal)
│   │
│   └── web/                    # The Astro Frontend
│       ├── src/                # Astro UI Components & Pages
│       └── dist/               # Build output (Static files -> S3, Server -> Lambda)
│
├── infra/                      # Infrastructure as Code
│   ├── bootstrap/              # One-time account setup (OIDC, state, deploy role)
│   ├── modules/                # Shared, reusable Terraform modules
│   └── environments/           # The Infrastructure Matrix
│       ├── serverless/         # (PROD) $0 Lambda + CloudFront + S3
│       ├── ecs/                # (LAB) Fargate Spot + ALB
│       └── eks/                # (LAB) Kubernetes
│
└── .github/workflows/          # CI/CD and automated Lab destruction
```

---

## 🚀 Development Workflows

### 1. Adding a New Feature (e.g., "Expenses")

**CRITICAL: Run tests after implementing any feature.**

1. Create the pure business logic in `apps/server/internal/modules/expenses`.
2. Define the input/output structs and the Service interface.
3. Wire the feature into `cmd/lambda/main.go` and `cmd/server/main.go` using the appropriate transport wrappers (API Gateway events vs HTTP requests).
4. Build the UI in `apps/web/src/pages/expenses`.
5. **Run tests to prevent regressions:**
   ```bash
   cd apps/web
   pnpm test           # Fast unit tests (< 300ms) - verifies form attributes
   ```
   See `docs/TESTING.md` for detailed testing guidelines.

### 2. Modifying Infrastructure

1. Identify if the change is global (e.g., VPC, IAM) or environment-specific.
2. Bootstrap changes live in `infra/bootstrap/`.
3. Shared changes live in `infra/modules/`.
4. Environment wiring lives in `infra/environments/*/`.

### 3. Documentation Contract (Update When You Change)

When you change a contract or boundary, update the corresponding doc:

- Infra changes (Terraform, AWS resources, IAM, CI/CD): update `docs/INFRA.md`.
- App architecture or routes: update `docs/ARCHITECTURE.md`.
- Schema or migrations: update `docs/SCHEMA.md`.
- Frontend/backend behavior and feature map: update `docs/APPLICATION.md`.
- Architecture decisions: update `docs/DECISIONS.md`.
- Testing strategy, test commands, or CI test changes: update `docs/TESTING.md`.

### 4. Testing Guidelines

**Before every commit:**
```bash
cd apps/web && pnpm test        # Verify form attributes (fast feedback)
```

**Test coverage priorities:**
- ✅ Unit tests: All forms must have `data-api-form` attribute
- ⏭️ Future: Heavy Go backend testing (Phase 1)
- ⏭️ Future: Contract tests between SSG frontend and Go API

**Note:** E2E tests were removed during the transition phase (3-6 weeks to SSG + Go backend). Rely on unit tests for fast feedback; full SigV4 validation will be tested when backend stabilizes.

See `docs/TESTING.md` for complete testing documentation.

---

## 📝 Coding Standards for AI

* **Go**: Prefer standard library. Use `slog` for structured logging. Return errors wrapped with context; never `panic()`.
* **Astro**: Use Tailwind utility classes. Favor SSR for dynamic routes and Static for marketing/shell pages.
* **Terraform**: Use Terraform `0.15+` syntax. Provide explicit `description` fields for all variables. Pin provider versions.
* **Cost Awareness**: Before suggesting an AWS service, evaluate its minimum monthly cost. If it breaks the $0 Serverless rule without being in an ephemeral Lab, suggest a cheaper alternative.

---

## 🎯 Current Goal

> **Phase 1: The Go Backend Rewrite & Serverless Prod.**
> Migrate the Astro SSR logic into a pure Go Hexagonal backend (`apps/server`), deploy it to AWS Lambda, and route static Astro assets to S3 via CloudFront.
