# TrackStack & Platform Engineering Agent Context

## Project Overview
**TrackStack** is a "Life Observability Platform" designed to track personal drift (calories, finance) while serving as a high-value **Platform Engineering** showcase.

**Core Philosophy:** "Architecture Evolution on Unified Infrastructure."
Instead of moving from VPS to Cloud, we demonstrate how to evolve software architecture (Starter -> Monolith -> Microservices) on a fixed, production-grade Kubernetes foundation (K3s on VPS) using a "Zero Config" CLI.

## The Strategy
1.  **The User Need:** A fast, "Lazy UX" app built with Astro.
2.  **The Engineering Need:** A "Golden Path" CLI that abstracts infrastructure complexity, allowing the user to deploy sophisticated architectures on a budget (Single VPS).

## Technology Stack

### Application Layer
-   **Frontend:** Astro (SSR) + TailwindCSS.
-   **Backend:** Go (Golang).
-   **Database:** **Turso (LibSQL)** over HTTP/WebSockets.
    -   *Crucial:* Replaces local SQLite to allow multiple pods/services to access data simultaneously without file-locking issues.
    -   *Migration:* Managed via "Maintenance Mode" checks during deployment.

### Platform Layer (The "Golden Paths")
-   **Infrastructure:** Single VPS (e.g., Hetzner) running **K3s**.
-   **Orchestration:** K3s (Lightweight Kubernetes).
-   **Registry:** GitHub Container Registry (GHCR).
-   **CLI:** Custom **Go CLI** (`trackstack-cli`) acting as the "Glue" and Orchestrator.

## The Golden Paths (Concurrent Deployment)
All three paths run simultaneously on the same K3s cluster using distinct Namespaces:

1.  **Starter Path:**
    -   *Architecture:* Astro SSR (Node.js) containerized directly.
    -   *Goal:* Demonstrate containerization and basic K8s deployment.
2.  **Monolith Path:**
    -   *Architecture:* Astro + Modular Monolith (Go).
    -   *Goal:* Demonstrate frontend/backend networking and service discovery.
3.  **Microservices Path:**
    -   *Architecture:* Astro + Distributed Microservices (Go).
    -   *Goal:* Demonstrate complex distributed patterns (Ingress, heavy scaling).

## Repository Structure (Mono-Repo)
```text
trackstack/
├── src/
│   ├── web/                  # 📦 Astro Project (Frontend/BFF)
│   └── go/                   # 📦 Go Project (Monolith & Microservices)
├── cli/                      # 🛠️ Go CLI (The Platform Tool)
├── infra/
│   ├── k8s/                  # Kubernetes Manifests / Helm Charts
│   └── ansible/              # Provisioning Playbooks (called by CLI)
├── .github/workflows/        # 🚀 CI/CD (Build -> Push GHCR -> SSH Trigger)
└── AGENTS.md                 # 🧠 Context for AI Agents
