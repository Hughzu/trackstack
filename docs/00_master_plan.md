# 🗺️ TrackStack & Platform Engineering: Master Plan v6

**Date:** 2026-02-16
**Strategy:** "The Managed Polylith" (Framework-based Approach)
**Goal:** Build "TrackStack" while creating `tstack`, a CLI that generates and **maintains** the core infrastructure and boilerplate for future apps.

---

## 1. The Dual Vision

1. **The Product (TrackStack):** A "Lazy UX" personal tracking app.
2. **The Platform (The Career Flex):** A custom CLI (`tstack`) that acts as a **Framework**. It injects a managed "Core" (Auth, DB, Infra) into apps so the developer only focuses on "Modules" (Business Logic).

---

## 2. Application Architecture: "The Managed Polylith"

We split the application into two distinct zones.

### Zone A: The "Core" (⛔ Managed by CLI)

Technical plumbing that rarely changes between apps.

* **Backend (Go):** `internal/core/` (Auth OIDC, Turso DB Setup, Logger, Server Config).
* **Frontend (Astro):** `src/core/` (Base Layouts, UI Kit like Buttons/Cards, Global Hooks).
* **Infrastructure:** `iac/_vendor/` (Complex Terraform Modules).

### Zone B: The "Modules" (✅ Owned by User)

The unique business value of the specific application.

* **Backend (Go):** `internal/modules/` (e.g., Calories logic, Heat logic).
* **Frontend (Astro):** `src/modules/` (Specific pages and components).
* **Configuration:** `cmd/main.go` (Wires Core and Modules together via Dependency Injection).

---

## 3. The Platform Engine: `tstack` CLI

The CLI is a "Fat Binary" (Go Embed) containing the Golden Paths.

### The Lifecycle

1. **Init:** `tstack init --name budget-app`
* Generates the repo structure.
* Injects the **Core** (Read-Only recommendation).
* Scaffolds the first **Module**.


2. **Dev:** User writes code in `internal/modules` and `src/modules`.
* *Rule:* Modules imports Core. Core **never** imports Modules.


3. **Upgrade:** `tstack upgrade`
* The CLI checks for updates in the "Core" layer (e.g., a security fix in Auth).
* It overwrites `internal/core` and `src/core` (unless the user has explicitly "ejected").
* It ensures all apps benefit from platform improvements without rewriting code.



---

## 4. The 4 Golden Paths (Infrastructure)

The CLI configures the infrastructure to support the Polylith at different scales.

| Step | Path Name | Deployment Target | FinOps Strategy |
| --- | --- | --- | --- |
| **0** | `local` | **Localhost** | Free. SQLite/Turso. |
| **1** | `serverless` | **Lambda + CloudFront** | **Production.** Scale-to-Zero. Cost: **~$0.60/mo**. |
| **2** | `container` | **Fargate Spot** | **Lab.** Ephemeral testing of Networking. |
| **3** | `cloud-native` | **EKS (Kubernetes)** | **Lab.** Distributed Polylith (Feature Flags). |

---

## 5. Repository Structure (The Vendor Pattern)

This structure enforces the separation between "Platform Code" and "Product Code".

```text
my-app/
├── cmd/
│   └── main.go                  # 🔌 WIRING (Injects Core into Modules)
│
├── internal/
│   ├── core/                    # ⛔ MANAGED (Do not edit)
│   │   ├── auth/                # OIDC/Session logic
│   │   ├── db/                  # Turso connection
│   │   └── server/              # HTTP Server (Chi/Echo)
│   │
│   └── modules/                 # ✅ USER LAND
│       ├── calories/            # Hexagonal Module
│       └── heat/                # Hexagonal Module
│
├── src/web/src/
│   ├── core/                    # ⛔ MANAGED (Do not edit)
│   │   ├── layouts/             # AppShell.astro, AuthLayout.astro
│   │   └── ui/                  # Button.astro, Card.astro (Design System)
│   │
│   ├── modules/                 # ✅ USER LAND
│   │   └── calories/            # Calorie Dashboard Components
│   │
│   └── pages/                   # ✅ ROUTING
│
├── iac/
│   ├── _vendor/                 # ⛔ MANAGED Terraform Modules
│   └── environments/            # ✅ USER Configuration
│
├── tstack.yaml                  # CLI Config (Version tracking)
└── Dockerfile                   # 🐳 Universal Builder

```

---

## 6. Execution Roadmap

### 🏁 Phase 1: The Prototype (Manual Mode)

* **Goal:** Build the first app manually to define the "Core" vs "Module" boundary.
* **Tasks:**
1. **Go:** Create `internal/core` (Auth/DB) and `internal/modules/calories`.
2. **Astro:** Create `src/core` (Layouts) and `src/modules`.
3. **Docker:** Write the Universal Dockerfile (Embed).
4. **Deploy:** Ship to AWS Lambda (Step 1).



### 🛠️ Phase 2: The Tooling (The CLI)

* **Goal:** Extract the "Core" into the CLI.
* **Tasks:**
1. Initialize `tstack` (Go Cobra).
2. Implement `tstack init` (Embeds the Core files).
3. Implement `tstack upgrade` (Overwrite logic for Core folders).



### 🚀 Phase 3: The Scale (The Labs)

* **Goal:** Add Infrastructure complexity.
* **Tasks:**
1. Add Step 2 (Fargate) and Step 3 (EKS) templates to the CLI.
2. Demonstrate deploying the same App Artifact to K8s.



---

## 7. Technical Guardrails

1. **Dependency Rule:** `Core` -> `Modules` is **FORBIDDEN**. `Modules` -> `Core` is **REQUIRED**.
2. **Eject Strategy:** If a user *must* modify a file in `core/`, they accept that `tstack upgrade` will no longer update that specific file (or they manage merge conflicts).
3. **Stateless:** All state lives in Turso or S3. The app is ephemeral.