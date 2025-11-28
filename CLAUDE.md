# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

TrackStack is a monorepo containing:
1. **TrackStack App** (`trackstack/app/`) - A modular monolith PWA for personal tracking (expenses, heat monitoring, recipes, calories)
2. **Blog** (`blog/`) - An Astro-based static blog
3. **Infrastructure as Code** (`iac/`) - Terraform modules for AWS deployment

## Development Commands

### Blog (`blog/`)
```bash
cd blog
pnpm install         # Install dependencies
pnpm dev             # Start dev server at localhost:4321
pnpm build           # Build to ./dist/
pnpm preview         # Preview production build
pnpm astro ...       # Run Astro CLI commands
```

### TrackStack App (`trackstack/app/`)
```bash
cd trackstack/app
pnpm install         # Install dependencies
pnpm dev             # Start dev server at localhost:4321
pnpm build           # Build to ./dist/
pnpm preview         # Preview production build
pnpm astro ...       # Run Astro CLI commands
```

### Infrastructure (`iac/`)
```bash
cd iac/blog          # or iac/trackstack
terraform init       # Initialize Terraform
terraform plan       # Plan infrastructure changes
terraform apply      # Apply changes
```

## Architecture

### TrackStack App - Modular Monolith Design

The TrackStack app follows a **modular monolith** architecture designed for future extraction into micro-frontends. Key principles:

**1. Module Independence**
- Each business module (`src/modules/`) is self-contained and isolated
- Modules: `expense`, `heat`, `cook`, `ball`, `kcal`
- Modules can only import from `core/` and `shared/`, NEVER from other modules
- Each module has its own database file in `data/` directory

**2. Directory Structure**
```
src/
├── modules/          # Business capability modules
│   ├── expense/     # Expense tracking
│   ├── heat/        # Heat monitoring
│   ├── cook/        # Recipe management
│   ├── kcal/        # Calorie tracking
│   ├── ball/        # Ball tracking
│
├── core/            # Shared infrastructure (auth, database, api clients)
├── shared/          # Design system (layouts, components, styles)
└── pages/           # Astro file-based routing
```

**3. Module Structure**
Each module follows consistent internal organization:
- `components/` - UI components (Astro + potentially Angular for interactivity)
- `services/` - API clients and business logic
- `repositories/` - Database access layer
- `migrations/` - Database schema migrations
- `types/` - TypeScript definitions
- `utils/` - Module-specific utilities
- `index.ts` - Public API (only exports)

**4. Database Strategy**
- Uses SQLite in production (one database per module)
- Database files stored in `data/` directory (git-ignored)
- Schema uses snake_case convention
- Future migration to .NET + EF Core will transform to PascalCase
- Each module owns its database migrations

**5. Dependency Rules**
```
modules/  → Can import: core/, shared/, self
core/     → Can import: shared/, self
shared/   → Can import: self only (zero dependencies)
```

**6. TypeScript Path Aliases**
Use these import paths:
```typescript
import { ... } from '@/modules/expense-tracker';  // Module public API
import { ... } from '@/core/auth';                // Core infrastructure
import { ... } from '@/shared/components';        // Shared components
```

**7. Technology Stack**
- Framework: Astro 5.x with PWA support
- Styling: Tailwind CSS 4.x
- Database: SQLite with better-sqlite3
- TypeScript: Strict mode enabled
- Package Manager: npm (for trackstack/app), pnpm (for blog)

### Blog Architecture

Standard Astro blog template:
- `src/pages/` - File-based routing
- `src/content/blog/` - MDX blog posts
- `src/components/` - Reusable Astro components
- `src/layouts/` - Page layouts
- MDX and sitemap support enabled

### Infrastructure Architecture

Terraform modules in `iac/common/modules/`:
- `static-website/` - S3 + CloudFront setup with Origin Access Control
- `budget/` - AWS Budget alerts

CloudFront includes URL rewrite function to append `index.html` for directory requests (Astro static site compatibility).

## Deployment

### Blog Deployment
GitHub Actions workflows in `.github/workflows/`:
- `deploy-blog-iac.yml` - Deploy infrastructure (auto-applies on push to main)
- `deploy-blog-to-s3.yml` - Build and deploy blog content to S3, invalidate CloudFront cache

**Infrastructure workflow:**
- Runs on pushes to `iac/blog/**`
- Terraform auto-applies changes on main branch
- Manual workflow dispatch supports plan/apply/destroy

**Application workflow:**
- Runs on pushes to `blog/src/**`, `blog/public/**`, etc.
- Uses npm to build (not pnpm, for CI consistency)
- Syncs to S3 with `--delete` flag
- Invalidates CloudFront and waits for completion

### Secrets Required
- `AWS_ROLE_ARN` - IAM role for infrastructure deployment
- `BLOG_ROLE_ARN` - IAM role for blog deployment
- `AWS_REGION` - AWS region
- `BLOG_S3_BUCKET_NAME` - S3 bucket name
- `BLOG_CLOUDFRONT_DISTRIBUTION_ID` - CloudFront distribution ID

## Key Architectural Decisions

### Why SQLite?
For a 5-user application, SQLite provides:
- Zero configuration and cost
- Fast performance for small datasets
- Simple backups (copy .db files to S3)
- Easy per-module extraction (each module has its own database)

### Why Modular Monolith?
- Validate business domains quickly
- Extract modules to standalone apps when they prove valuable
- Each module is architecturally ready for extraction
- No premature micro-frontend complexity

### Why snake_case Database Schema?
- Standard SQLite convention
- Clear visual distinction from application code
- Intentional: Will migrate to PascalCase when moving to .NET + EF Core
- Accepts one-time migration cost rather than premature optimization

## Important Patterns

### When Adding a New Module
1. Create `src/modules/new-module/` with standard structure
2. Add database in `data/new-module.db` (create in `core/database/connection.ts`)
3. Create migrations in `modules/new-module/migrations/`
4. Export public API via `modules/new-module/index.ts`
5. Add routes in `src/pages/new-module/`

### When Working with Database
- Always use repository pattern for database access
- Never import database connections directly in components
- Use migrations for schema changes
- Enable WAL mode and foreign keys in connection setup

### When Creating Components
- Use Astro components for static/server-rendered content
- Use Angular components for dynamic/interactive features
- Place shared components in `src/shared/components/`
- Place domain-specific components in module's `components/` folder

## Common Pitfalls to Avoid

1. **DO NOT** import from other modules' internals - use public API only
2. **DO NOT** create circular dependencies between modules
3. **DO NOT** put business logic in `shared/` - it belongs in modules or `core/`
4. **DO NOT** use relative imports across major boundaries - use path aliases
5. **DO NOT** commit the `data/` directory - databases are git-ignored

## Future Roadmap Context

**Phase 1 (Current):** Validate monolith, implement all tracking modules
**Phase 2:** Migrate backend to .NET + EF Core, run schema migration
**Phase 3:** Extract highest-value module to standalone app
**Phase 4:** Scale database if needed (probably not for <50 users)

The architecture intentionally optimizes for speed and learning over premature scaling.