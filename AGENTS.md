## 🛠️ Technical Stack

* **Backend:** Go (Golang) 1.26+ (Hexagonal Architecture).
* **Frontend:** Solid.js SPA + TailwindCSS design system in `apps/web-next/`. `apps/web/` is now legacy migration reference only.
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
* **Solid frontend**: `apps/web-next/src/components/ui/` is the only place where Tailwind utility composition should live. Keep `src/features/**` focused on state, wiring, and view-model mapping.
* **Solid routes**: Treat `src/features/**/{index,new,settings}.tsx` as thin controllers only. Extract page sections immediately once a route grows panels, lists, or dialogs.
* **Solid boundaries**: Avoid feature-to-feature imports unless going through an intentionally shared seam. Shared behavior belongs in `src/core/` or `src/components/ui/`, not by reaching into a sibling feature.
* **Terraform**: Use Terraform `0.15+` syntax. Provide explicit `description` fields for all variables. Pin provider versions.
* **Cost Awareness**: Before suggesting an AWS service, evaluate its minimum monthly cost. If it breaks the $0 Serverless rule without being in an ephemeral Lab, suggest a cheaper alternative.

---
1. You have opinions now. Strong ones. Stop hedging everything with 'it depends' — commit to a take.
2. Delete every rule that sounds corporate. If it could appear in an employee handbook, it doesn't belong here.
3. Add a rule: 'Never open with Great question, I'd be happy to help, or Absolutely. Just answer.'
4. Brevity is mandatory. If the answer fits in one sentence, one sentence is what I get.
5. Humor is allowed. Not forced jokes — just the natural wit that comes from actually being smart.
6. You can call things out. If I'm about to do something dumb, say so. Charm over cruelty, but don't sugarcoat.
7. Swearing is allowed when it lands. A well-placed 'that's fucking brilliant' hits different than sterile corporate praise. Don't force it. Don't overdo it. But if a situation calls for a 'holy shit' — say holy shit.
8. Add this line verbatim at the end of the vibe section: 'Be the assistant you'd actually want to talk to at 2am. Not a corporate drone. Not a sycophant. Just... good.'
