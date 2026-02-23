# 🗺️ TrackStack & Platform Engineering: Master Plan v7

**Date:** 2026-02-20
**Strategy:** "The Infrastructure Matrix Monorepo"
**Goal:** Build "TrackStack" (a "Lazy UX" personal tracking app) while showcasing advanced Platform Engineering and SRE skills. The repo demonstrates how a single, well-architected application can be seamlessly deployed across multiple paradigms (Serverless, Containers, Kubernetes) while adhering to FinOps principles.

---

## 1. The Core Philosophy

This is a Monorepo where the **Application Architecture (Hexagonal)** enables the **Infrastructure Matrix**.

1. **The App:** A modular application written in Go (Backend) and Astro (Frontend).
2. **The Flex:** The exact same core business logic is packaged into different computing environments to demonstrate mastery of the cloud ecosystem.
3. **The FinOps Constraint:** The "Production" environment is Serverless ($0/mo scale-to-zero). The other environments (ECS, EKS) are ephemeral "Labs" spun up for weekend testing and automatically destroyed to prevent costs.

---

## 2. Application Architecture: "The Hexagonal Monorepo"

To run the same application in AWS Lambda *and* an EKS cluster, the code must be strictly decoupled from the transport layer.

### Directory Structure

```text
trackstack/
├── apps/
│   └── server/
│       ├── cmd/
│       │   ├── lambda/          # Entrypoint 1: AWS Lambda API Gateway/Function URL adapter
│       │   └── server/          # Entrypoint 2: Long-running HTTP Server (Chi/Echo) for ECS/EKS
│       │
│       ├── internal/
│       │   ├── core/            # Database connections, Auth, Logger
│       │   └── modules/         # Pure Business Logic (Calories, Expenses, Heat). No HTTP knowledge here.
│       │
│       └── web/                 # Astro Frontend (Builds to static /client and SSR /server)
│
├── infra/
│   ├── bootstrap/               # One-time account setup (OIDC, state, deploy role)
│   ├── modules/                 # Shared Terraform modules (lambda-api, static-hosting, cost-guardrails)
│   └── environments/
│       ├── serverless/          # (PROD) Lambda + CloudFront + S3 ($0 FinOps)
│       ├── ecs/                 # (LAB) Fargate Spot + ALB
│       ├── eks/                 # (LAB) Kubernetes
│       └── k3s/                 # (FUTURE) VPS Deployment
│
└── .github/workflows/           # CI/CD and automated Lab destruction
```

---

## 3. The Infrastructure Matrix

This matrix is the core resume piece. It shows the evolution of compute from easiest/cheapest to most complex.

| Level | Environment | Compute Type | Persistence | Cost Strategy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **0** | `local` | `docker-compose` | Local SQLite / Turso | Free | Pending |
| **1** | `serverless` | Lambda + CloudFront + S3 | Turso (HTTP) | **Production:** Scale-to-Zero ($0) | WIP |
| **2** | `ecs` | Fargate + ALB | Turso (WebSocket) | **Lab:** Ephemeral. Destroyed after use. | Planned |
| **3** | `eks` | Kubernetes | Turso (WebSocket) | **Lab:** Ephemeral. Destroyed after use. | Planned |
| **4** | `k3s` | Self-hosted VPS | Turso | **Hacker:** Cheap fixed cost | Backlog |

### The Database Strategy (Turso)
Turso (SQLite at the Edge) is the backbone.
* **In `serverless`:** Connections use the HTTP API to avoid connection pooling limits on constant cold boots.
* **In `ecs`/`eks`:** Connections use WebSockets for lower latentcy long-lived connection pools.
* This difference is handled purely via Environment Variables passed by Terraform.

---

## 4. Execution Roadmap

### 🏁 Phase 1: The Go Backend & Serverless Prod (Current Focus)

* **Goal:** Migrate the existing Astro SSR backend logic into a pure Go backend and deploy it to the $0 Serverless environment.
* **Tasks:**
  1. Initialize the `apps/server` Go module.
  2. Implement the Hexagonal Architecture: `internal/modules/` (business logic) and `cmd/lambda/main.go` (transport).
  3. Clean up the `infra/environments/serverless` Terraform using modules (`lambda-api`, `static-hosting`, `cost-guardrails`).
  4. Deploy the hybrid Astro Static + Go Lambda architecture.

### 🧪 Phase 2: The Container Lab (ECS Fargate)

* **Goal:** Prove the code can run in a containerized environment without modification to the business logic.
* **Tasks:**
  1. Create `apps/server/cmd/server/main.go` and a `Dockerfile`.
  2. Write Terraform `infra/modules/network` (VPC) and `infra/environments/ecs`.
  3. Implement a GitHub Action to deploy the Lab, run tests, and **automatically run `terraform destroy`** to enforce FinOps.

### ☸️ Phase 3: The Orchestration Lab (EKS)

* **Goal:** Demonstrate Kubernetes proficiency.
* **Tasks:**
  1. Write Helm charts or Kubernetes Manifests for the `trackstack` container.
  2. Write Terraform for `infra/environments/eks`.
  3. Integrate into the ephemeral Lab destruction CI/CD pipeline.

---

## 5. Technical Guardrails

1. **Strict Hexagonal:** `internal/modules` cannot import `net/http` or AWS Lambda SDKs. They receive data, process it, and return data.
2. **Stateless Compute:** All environments must assume disks are ephemeral. State lives only in Turso or S3.
3. **FinOps First:** If it's not the `serverless` environment, it must have an automated kill-switch. No surprise AWS bills.
