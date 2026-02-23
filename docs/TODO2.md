# 🧭 Master Plan: Regaining Control (Astro / AWS / Turso)

## 🧠 The Philosophy: From "Vibe Coder" to "System Architect"
*Vibe Coding* (AI-driven development) generates immense speed, but also silent technical debt. To maintain a **craftsman's** standard, the goal is no longer to read and polish every single line of generated code, but to **lock down the contracts and boundaries** between system components.

**Golden Rules of the Audit:**
1. **Flows over Files:** Focus your effort on the integration points (Front -> Back -> DB).
2. **The Black Box Principle:** If the inputs (Props/Types) and outputs (UI/Data) of a component are solid, tolerate temporary internal messiness.
3. **Documentation as a Constraint:** Write documentation to set strict boundaries for the AI in future iterations, not just to explain past code.

---

## 📂 The "Control Room": Markdown Files to Create

To take back ownership of the codebase and guide the AI, create the following 6 documents (use the AI to generate the first drafts based on your existing code).

### 1. `ARCHITECTURE.md` (The Macro View) - DONE
* **Goal:** Understand how a user request travels to the data layer.
* **Content:**
  * Mermaid.js diagram of the complete flow (Astro -> CloudFront -> S3/Lambda -> Turso).
  * List of required environment variables and their scope (Frontend vs. Server-side).
  * High-level deployment process.

### 2. `INFRASTRUCTURE.md` (The IaC Vault) - DONE
* **Goal:** Audit and freeze Terraform/AWS resources.
* **Content:**
  * Inventory of active resources (e.g., 1 CDN, 1 S3 Bucket, 1 Lambda).
  * **Security Boundaries (CRITICAL):** Strict definition of IAM roles (e.g., "The Lambda can only read from S3 and execute code; no admin rights").
  * Manual infrastructure steps not covered by IaC.

### 3. `SCHEMA.md` (The Single Source of Truth) - DONE
* **Goal:** Safeguard the Turso database structure.
* **Content:**
  * Schema of core tables and their relationships.
  * Location of TypeScript / Zod types (e.g., `src/types/...`).
  * Mutation rules (e.g., "Only API Routes inside `src/pages/api/` are allowed to mutate the database").

### 4. `APPLICATION.md` (Business Logic & Astro) - DONE
* **Goal:** Prevent the AI from reinventing your frontend architecture with every new feature.
* **Content:**
  * **Feature Map:** The main business domains (Auth, Dashboard, User Profile, etc.).
  * **Folder Organization:** Strict rules for `src/pages/` (SSR routing) vs. `src/components/ui/` (Dumb/Presentational) vs. `src/lib/` (Turso/DB logic).
  * **Philosophy:** Strict separation between data fetching (Astro Frontmatter) and UI rendering (Template).

### 5. `DECISIONS.md` (The Architecture Decision Log / ADR)
* **Goal:** Record architectural choices to prevent repeating the same debates.
* **Content:**
  * Why Turso? (Distributed SQLite, Serverless-friendly, cost-effective).
  * Why AWS Lambda + S3 instead of Vercel? (Total infrastructure control via Terraform, predictable costs).