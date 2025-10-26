# Architecture Overview

## Phase 1: MVP
Work in progress

## Design Principles
- Micro-frontends: Each Angular app is isolated
- Shell pattern: Astro manages routing, auth, layout
- Extraction-ready: Apps can become standalone with minimal effort
- Event-driven foundation: Even in MVP, thinking about events if it make sense

## Domain Boundaries
1. **Expenses**: Financial tracking
2. **Energy**: Pellet consumption monitoring
3. **Recipes**: Shared recipe book with nutrition
4. **Calories**: Daily calorie logging

Each domain is independent at the frontend level.