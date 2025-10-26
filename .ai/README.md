# TrackStack - AI Context

## Project Overview
Personal tracking platform with 4 domains: Expenses, Heat Consumption, Recipes, Calories.

## Current Phase
Phase 1: MVP with Astro + Angular micro-frontends

## Key Architectural Decisions
- Astro as shell (SSG + routing)
- Angular apps as micro-frontends (isolated domains)
- Design for extraction (minimal effort to make standalone)
- SQLite for initial data storage

## Project Structure
- `/blog` - Blog where I write my journey
- `/trackstack` - My app
- `/trackstack/frontend` - Astro shell + Angular apps
- `/trackstack/backend` - Will be added in Phase 2 (C# monolith)
- `/trackstack/iac` - Will be added in Phase 4 (Terraform)

## Tech Stack
- Frontend: Astro, Angular
- Backend: TBD (Phase 2: C# + SQLite)
- Infrastructure: TBD - AWS (S3 + CloudFront for now)