## 🛠️ Technical Stack

* **Backend:** Go (Golang) 1.26+ (Hexagonal Architecture).
* **Frontend:** Astro (SSR & Static) + TailwindCSS in `apps/web/`.
* **Database:** Turso (SQLite at the Edge via HTTP/WebSockets).
* **Infra:** Terraform (AWS) structured by environment.
* **Deployment Matrix:**
  * **Prod (`serverless`):** AWS Lambda (Go) + CloudFront + S3 (Astro Static).
  * **Lab (`ecs`/`eks`):** Docker Containers (AWS Fargate / Kubernetes).

---

## 🚀 Development Workflows

### Documentation Contract (Update When You Change)

When you change a contract or boundary, update the corresponding doc:

- Infra changes (Terraform, AWS resources, IAM, CI/CD): update `docs/INFRA.md`.
- App architecture or routes: update `docs/ARCHITECTURE.md`.
- Schema or migrations: update `docs/SCHEMA.md`.
- Frontend/backend behavior and feature map: update `docs/APPLICATION.md`.
- Architecture decisions: update `docs/DECISIONS.md`.
- Testing strategy, test commands, or CI test changes: update `docs/TESTING.md`.

## 📝 Coding Standards for AI

* **Go**: Prefer standard library. Use `slog` for structured logging. Return errors wrapped with context; never `panic()`.
* **Astro**: Use Tailwind utility classes. Favor SSR for dynamic routes and Static for marketing/shell pages.
* **Terraform**: Use Terraform `0.15+` syntax. Provide explicit `description` fields for all variables. Pin provider versions.
* **Cost Awareness**: Before suggesting an AWS service, evaluate its minimum monthly cost. If it breaks the $0 Serverless rule without being in an ephemeral Lab, suggest a cheaper alternative.
