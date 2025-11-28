Generate pwa icons : https://www.pwabuilder.com/imageGenerator

# TrackStack Application - Astro Monolith

> A modular monolith designed for future extraction - business domains first, technology second.

## Overview

TrackStack is a personal tracking platform built as an Astro monolith with a clear path to micro-frontend extraction. The architecture prioritizes business domain boundaries over technical layers, making it trivial to split modules into standalone applications when they prove valuable.

**Core Philosophy:** Validate with a monolith, extract when it makes business sense.

## Business Domains

TrackStack manages four distinct tracking capabilities:

1. **Expense Tracker** - Personal finance and spending analysis
2. **Heat Monitor** - Wood pellet consumption and heating efficiency
3. **Recipe Book** - Shared recipe management with nutritional tracking
4. **Calorie Tracker** - Daily calorie and macro logging

Each domain is architecturally isolated and can become an independent application with minimal refactoring.

---

## Repository Structure

```
trackstack/app/
├── data/                           # SQLite databases (git-ignored, backup to S3)
│   ├── auth.db                     # Users, sessions, tokens
│   ├── auth.db-wal                 # Write-Ahead Log
│   ├── auth.db-shm                 # Shared memory
│   ├── expenses.db                 # Expense tracker data
│   ├── heat.db                     # Heat monitoring data
│   ├── recipes.db                  # Recipe book data
│   └── calories.db                 # Calorie tracking data
│
├── src/
│   ├── modules/                    # Business capability modules (extraction candidates)
│   │   ├── expense-tracker/
│   │   │   ├── components/        # Domain-specific UI components (Astro + Angular)
│   │   │   │   ├── ExpenseList.astro          # Static list component
│   │   │   │   ├── ExpenseCard.astro          # Static card component
│   │   │   │   ├── ExpenseForm.component.ts   # Dynamic Angular form
│   │   │   │   └── ExpenseChart.component.ts  # Dynamic Angular chart
│   │   │   ├── services/          # API clients and business logic
│   │   │   │   └── api.ts         # Expense API endpoints
│   │   │   ├── repositories/      # Database access layer
│   │   │   │   └── expenseRepository.ts
│   │   │   ├── migrations/        # Database schema migrations
│   │   │   │   └── index.ts       # SQLite migration definitions
│   │   │   ├── types/             # TypeScript types and interfaces
│   │   │   ├── utils/             # Domain utilities (formatters, validators)
│   │   │   └── index.ts           # Public API - only export what's needed
│   │   │
│   │   ├── heat-monitor/
│   │   │   ├── components/
│   │   │   ├── services/
│   │   │   ├── repositories/
│   │   │   ├── migrations/
│   │   │   ├── types/
│   │   │   ├── utils/
│   │   │   └── index.ts
│   │   │
│   │   ├── recipe-book/
│   │   │   ├── components/
│   │   │   ├── services/
│   │   │   ├── repositories/
│   │   │   ├── migrations/
│   │   │   ├── types/
│   │   │   ├── utils/
│   │   │   └── index.ts
│   │   │
│   │   └── calorie-tracker/
│   │       ├── components/
│   │       ├── services/
│   │       ├── repositories/
│   │       ├── migrations/
│   │       ├── types/
│   │       ├── utils/
│   │       └── index.ts
│   │
│   ├── core/                       # Shared infrastructure (extraction = npm package)
│   │   ├── database/               # Database infrastructure
│   │   │   ├── connection.ts      # SQLite connection management
│   │   │   ├── migrations.ts      # Migration runner
│   │   │   └── init.ts            # Database initialization
│   │   │
│   │   ├── auth/                   # Authentication & authorization
│   │   │   ├── components/        # Login, Logout, AuthGuard components
│   │   │   │   ├── LoginForm.astro
│   │   │   │   ├── LogoutButton.astro
│   │   │   │   └── AuthGuard.astro
│   │   │   ├── middleware/        # Auth route protection
│   │   │   │   └── authMiddleware.ts
│   │   │   ├── services/          # Auth API client
│   │   │   │   ├── authClient.ts  # Core auth logic
│   │   │   │   └── tokenManager.ts # JWT storage & refresh
│   │   │   ├── repositories/      # Auth database access
│   │   │   │   └── userRepository.ts
│   │   │   ├── migrations/        # Auth schema migrations
│   │   │   │   └── index.ts
│   │   │   └── types/
│   │   │       └── auth.ts        # User, Session, Token types
│   │   │
│   │   ├── api/                    # HTTP client configuration
│   │   │   ├── client.ts          # Base fetch wrapper with interceptors
│   │   │   ├── config.ts          # API base URLs and endpoints
│   │   │   └── types.ts           # Common API response types
│   │   │
│   │   └── router/                 # Routing utilities
│   │       └── guards.ts           # Route protection helpers
│   │
│   ├── shared/                     # Design system (extraction = npm package)
│   │   ├── components/             # Reusable UI components
│   │   │   ├── Button/
│   │   │   │   ├── Button.astro
│   │   │   │   └── Button.types.ts
│   │   │   ├── Card/
│   │   │   │   ├── Card.astro
│   │   │   │   └── Card.types.ts
│   │   │   ├── Form/
│   │   │   │   ├── Input.astro
│   │   │   │   ├── Select.astro
│   │   │   │   └── Form.types.ts
│   │   │   ├── Modal/
│   │   │   │   ├── Modal.astro
│   │   │   │   └── Modal.types.ts
│   │   │   └── Table/
│   │   │       ├── Table.astro
│   │   │       └── Table.types.ts
│   │   │
│   │   ├── layouts/                # Common page layouts
│   │   │   ├── AppLayout.astro     # Main authenticated app layout
│   │   │   ├── DashboardLayout.astro # Dashboard with sidebar
│   │   │   └── AuthLayout.astro    # Login/register pages
│   │   │
│   │   ├── styles/                 # Global styles and theming
│   │   │   ├── theme.css           # CSS variables for theming
│   │   │   ├── tailwind.config.ts  # Tailwind configuration
│   │   │   └── global.css          # Global CSS resets
│   │   │
│   │   └── utils/                  # Common utilities
│   │       ├── date.ts             # Date formatting and manipulation
│   │       ├── format.ts           # Currency, number formatters
│   │       └── validation.ts       # Form validation helpers
│   │
│   ├── pages/                      # Astro routing (file-based)
│   │   ├── index.astro             # Landing page (/)
│   │   ├── dashboard.astro         # Main dashboard (/dashboard)
│   │   │
│   │   ├── api/                    # API routes
│   │   │   ├── expenses/          # Expense API endpoints
│   │   │   │   ├── index.ts       # GET/POST /api/expenses
│   │   │   │   └── [id].ts        # GET/PUT/DELETE /api/expenses/:id
│   │   │   ├── heat/              # Heat monitoring API
│   │   │   ├── recipes/           # Recipe API
│   │   │   └── calories/          # Calorie tracking API
│   │   │
│   │   ├── expenses/               # Expense tracker routes
│   │   │   ├── index.astro         # /expenses - list view
│   │   │   ├── new.astro           # /expenses/new - create
│   │   │   └── [id].astro          # /expenses/:id - detail view
│   │   │
│   │   ├── heat/                   # Heat monitor routes
│   │   │   ├── index.astro         # /heat - monitoring dashboard
│   │   │   └── history.astro       # /heat/history - historical data
│   │   │
│   │   ├── recipes/                # Recipe book routes
│   │   │   ├── index.astro         # /recipes - recipe list
│   │   │   ├── new.astro           # /recipes/new - create recipe
│   │   │   └── [id].astro          # /recipes/:id - recipe detail
│   │   │
│   │   └── calories/               # Calorie tracker routes
│   │       ├── index.astro         # /calories - daily log
│   │       └── stats.astro         # /calories/stats - analytics
│   │
│   └── env.d.ts                    # TypeScript environment types
│
├── public/                         # Static assets (images, fonts, etc.)
│   ├── favicon.svg
│   └── images/
│
├── astro.config.mjs                # Astro configuration
├── package.json
├── tsconfig.json                   # TypeScript configuration
├── tailwind.config.ts              # Tailwind CSS configuration
└── .gitignore
```

---

## Architecture Principles

### 1. Module Independence

Each business module is **self-contained** and enforces encapsulation:

```typescript
// ❌ BAD: Direct import from another module's internals
import { ExpenseList } from '@/modules/expense-tracker/components/ExpenseList';

// ✅ GOOD: Import from module's public API
import { ExpenseList } from '@/modules/expense-tracker';
```

**Why this matters:** When you extract a module, you only need to maintain the public API contract. Internal refactoring doesn't break other modules.

### 2. Dependency Rules

Clear hierarchy prevents circular dependencies and makes extraction clean:

```
modules/expense-tracker/  → Can import: core/, shared/, self
modules/heat-monitor/     → Can import: core/, shared/, self
modules/recipe-book/      → Can import: core/, shared/, self
modules/calorie-tracker/  → Can import: core/, shared/, self

core/                     → Can import: shared/, self
shared/                   → Can import: self only (zero dependencies)
```

**Modules NEVER import from other modules.** If shared logic is needed, it goes in `core/` or `shared/`.

### 3. API Layer Abstraction

Each module owns its API client, making backend extraction trivial:

```typescript
// src/modules/expense-tracker/services/api.ts
import { apiClient } from '@/core/api/client';

export const expenseApi = {
  list: () => apiClient.get('/api/expenses'),
  create: (data) => apiClient.post('/api/expenses', data),
  update: (id, data) => apiClient.put(`/api/expenses/${id}`, data),
  delete: (id) => apiClient.delete(`/api/expenses/${id}`)
};
```

**When you extract:** Change the base URL to point to a dedicated microservice:

```typescript
// After extraction
const apiClient = createClient('https://expenses-api.trackstack.com');
```

### 4. TypeScript Path Aliases

Clean imports via `tsconfig.json`:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/modules/*": ["src/modules/*"],
      "@/core/*": ["src/core/*"],
      "@/shared/*": ["src/shared/*"]
    }
  }
}
```

---

## Module Structure Deep Dive

### Anatomy of a Module

Every module follows the same internal structure for consistency:

```
src/modules/expense-tracker/
├── components/                      # UI components (Astro + Angular)
│   ├── ExpenseList.astro           # Static list component (Astro)
│   ├── ExpenseCard.astro           # Static card component (Astro)
│   ├── ExpenseForm.component.ts    # Dynamic form (Angular)
│   └── ExpenseChart.component.ts   # Dynamic chart (Angular)
│
├── services/                        # Business logic and API
│   ├── api.ts                      # HTTP calls to backend
│   └── expenseCalculations.ts     # Domain logic
│
├── repositories/                    # Database access layer
│   └── expenseRepository.ts        # SQLite CRUD operations
│
├── migrations/                      # Database schema migrations
│   └── index.ts                    # Migration definitions
│
├── types/                           # TypeScript definitions
│   ├── expense.ts                  # Expense, ExpenseCategory types
│   └── api.ts                      # API request/response types
│
├── utils/                           # Module-specific utilities
│   ├── formatters.ts               # Currency, date formatting
│   └── validators.ts               # Form validation rules
│
└── index.ts                         # Public API
    // Only export what other parts of the app need
    export { ExpenseList, ExpenseCard } from './components';
    export { ExpenseFormComponent, ExpenseChartComponent } from './components';
    export { expenseApi } from './services/api';
    export type { Expense, ExpenseCategory } from './types/expense';
```

### Example: Expense API Service

```typescript
// src/modules/expense-tracker/services/api.ts
import { apiClient } from '@/core/api/client';
import type { Expense, CreateExpenseDto, UpdateExpenseDto } from '../types/expense';
import type { PaginatedResponse } from '@/core/api/types';

export const expenseApi = {
  /**
   * Fetch all expenses with optional filtering
   */
  list: async (params?: {
    page?: number;
    limit?: number;
    category?: string;
  }): Promise<PaginatedResponse<Expense>> => {
    return apiClient.get('/api/expenses', { params });
  },

  /**
   * Get a single expense by ID
   */
  getById: async (id: string): Promise<Expense> => {
    return apiClient.get(`/api/expenses/${id}`);
  },

  /**
   * Create a new expense
   */
  create: async (data: CreateExpenseDto): Promise<Expense> => {
    return apiClient.post('/api/expenses', data);
  },

  /**
   * Update an existing expense
   */
  update: async (id: string, data: UpdateExpenseDto): Promise<Expense> => {
    return apiClient.put(`/api/expenses/${id}`, data);
  },

  /**
   * Delete an expense
   */
  delete: async (id: string): Promise<void> => {
    return apiClient.delete(`/api/expenses/${id}`);
  }
};
```

---

## Core Infrastructure

### Authentication (`core/auth/`)

The auth module will eventually become `@trackstack/auth` npm package.

**Key components:**

```typescript
// src/core/auth/services/authClient.ts
export class AuthClient {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    });
    
    if (!response.ok) throw new Error('Login failed');
    
    const data = await response.json();
    await this.tokenManager.setTokens(data.accessToken, data.refreshToken);
    
    return data;
  }

  async logout(): Promise<void> {
    await fetch('/api/auth/logout', { method: 'POST' });
    await this.tokenManager.clearTokens();
  }

  async getCurrentUser(): Promise<User | null> {
    const token = await this.tokenManager.getAccessToken();
    if (!token) return null;
    
    // Decode JWT or fetch from API
    return this.fetchUserProfile();
  }

  async refreshAccessToken(): Promise<string> {
    const refreshToken = await this.tokenManager.getRefreshToken();
    if (!refreshToken) throw new Error('No refresh token');
    
    const response = await fetch('/api/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refreshToken })
    });
    
    const { accessToken } = await response.json();
    await this.tokenManager.setAccessToken(accessToken);
    
    return accessToken;
  }
}
```

**Usage in pages:**

```astro
---
// src/pages/expenses/index.astro
import { authClient } from '@/core/auth/services/authClient';
import AppLayout from '@/shared/layouts/AppLayout.astro';
import { ExpenseList } from '@/modules/expense-tracker';

// Protect the route
const user = await authClient.getCurrentUser();
if (!user) {
  return Astro.redirect('/login');
}
---

<AppLayout user={user}>
  <h1>My Expenses</h1>
  <ExpenseList />
</AppLayout>
```

### API Client (`core/api/`)

Centralized HTTP client with interceptors for auth tokens:

```typescript
// src/core/api/client.ts
import { authClient } from '@/core/auth/services/authClient';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    // Inject auth token
    const token = await authClient.getAccessToken();
    const headers = {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers
    };

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers
    });

    // Handle 401 - refresh token and retry
    if (response.status === 401) {
      await authClient.refreshAccessToken();
      return this.request(endpoint, options); // Retry with new token
    }

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
  }

  get<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'GET' });
  }

  post<T>(endpoint: string, data: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  put<T>(endpoint: string, data: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }

  delete<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'DELETE' });
  }
}

export const apiClient = new ApiClient(import.meta.env.PUBLIC_API_URL || '/api');
```

---

## Database & Persistence Strategy

### SQLite: The Pragmatic Choice

TrackStack uses **SQLite in production** with one database per business capability. For a 5-user application, SQLite provides:

- **Zero configuration** - No database server to manage
- **Zero cost** - No RDS fees
- **Fast performance** - Often faster than network-based databases for small datasets
- **Simple backups** - Just copy the `.db` files to S3
- **Perfect for extraction** - Each module has its own database ready to move

**Database per module:**
```
data/
├── auth.db         # Users, sessions, tokens (shared infrastructure)
├── expenses.db     # Expense tracker (domain-specific)
├── heat.db         # Heat monitoring (domain-specific)
├── recipes.db      # Recipe book (domain-specific)
└── calories.db     # Calorie tracking (domain-specific)
```

### Schema Strategy: snake_case Now, Migration Later

**Current approach (Phase 1): snake_case**
```sql
CREATE TABLE expenses (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL,
  description TEXT,
  date TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
```

**Why snake_case now?**
- Standard SQLite convention
- Readable in database tools
- Clear visual distinction from application code

**Future approach (Phase 2): Migration to PascalCase for EF Core**

When migrating to .NET with Entity Framework Core, we'll run a migration script to transform the schema:

```sql
-- Migration script (executed once during backend transition)
CREATE TABLE Expenses (
  Id TEXT PRIMARY KEY,
  UserId TEXT NOT NULL,
  Amount REAL NOT NULL,
  Category TEXT NOT NULL,
  Description TEXT,
  Date TEXT NOT NULL,
  CreatedAt TEXT NOT NULL,
  UpdatedAt TEXT NOT NULL
);

INSERT INTO Expenses (Id, UserId, Amount, Category, Description, Date, CreatedAt, UpdatedAt)
SELECT id, user_id, amount, category, description, date, created_at, updated_at
FROM expenses;

DROP TABLE expenses;
```

**This is intentional** - we accept a one-time migration cost rather than prematurely optimizing for a future that may not happen.

### Database Connection Management

```typescript
// src/core/database/connection.ts
import Database from 'better-sqlite3';
import path from 'path';

const DB_DIR = path.join(process.cwd(), 'data');

export function getDatabase(name: string): Database.Database {
  const dbPath = path.join(DB_DIR, `${name}.db`);
  const db = new Database(dbPath);
  
  // Enable WAL mode for better concurrent read performance
  db.pragma('journal_mode = WAL');
  
  // Enable foreign keys (not enabled by default in SQLite!)
  db.pragma('foreign_keys = ON');
  
  return db;
}

// Pre-initialized connections for each domain
export const authDb = getDatabase('auth');
export const expensesDb = getDatabase('expenses');
export const heatDb = getDatabase('heat');
export const recipesDb = getDatabase('recipes');
export const caloriesDb = getDatabase('calories');
```

### Migration System

Each module owns its schema migrations:

```typescript
// src/modules/expense-tracker/migrations/index.ts
import type { Migration } from '@/core/database/migrations';

export const expenseMigrations: Migration[] = [
  {
    version: 1,
    up: (db) => {
      db.exec(`
        CREATE TABLE expenses (
          id TEXT PRIMARY KEY,
          user_id TEXT NOT NULL,
          amount REAL NOT NULL,
          category TEXT NOT NULL,
          description TEXT,
          date TEXT NOT NULL,
          created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
          updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX idx_expenses_user_id ON expenses(user_id);
        CREATE INDEX idx_expenses_date ON expenses(date);
        CREATE INDEX idx_expenses_category ON expenses(category);
      `);
    },
    down: (db) => {
      db.exec('DROP TABLE expenses');
    }
  }
];
```

Migrations run automatically on startup:

```typescript
// src/core/database/init.ts
import { authDb, expensesDb, heatDb, recipesDb, caloriesDb } from './connection';
import { runMigrations } from './migrations';

import { authMigrations } from '@/core/auth/migrations';
import { expenseMigrations } from '@/modules/expense-tracker/migrations';
// ... other migrations

export function initializeDatabases() {
  console.log('Initializing databases...');
  
  runMigrations(authDb, authMigrations);
  runMigrations(expensesDb, expenseMigrations);
  runMigrations(heatDb, heatMigrations);
  runMigrations(recipesDb, recipeMigrations);
  runMigrations(caloriesDb, calorieMigrations);
  
  console.log('Databases initialized successfully');
}
```

### Repository Pattern

Each module uses repositories to abstract database access:

```typescript
// src/modules/expense-tracker/repositories/expenseRepository.ts
import { expensesDb } from '@/core/database/connection';
import type { Expense, CreateExpenseDto } from '../types/expense';
import { randomUUID } from 'crypto';

export class ExpenseRepository {
  create(userId: string, data: CreateExpenseDto): Expense {
    const id = randomUUID();
    const now = new Date().toISOString();
    
    const stmt = expensesDb.prepare(`
      INSERT INTO expenses (id, user_id, amount, category, description, date, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `);
    
    stmt.run(id, userId, data.amount, data.category, data.description, data.date, now, now);
    
    return this.findById(id)!;
  }

  findById(id: string): Expense | null {
    const stmt = expensesDb.prepare('SELECT * FROM expenses WHERE id = ?');
    return stmt.get(id) as Expense | null;
  }

  findByUserId(userId: string): Expense[] {
    const stmt = expensesDb.prepare('SELECT * FROM expenses WHERE user_id = ? ORDER BY date DESC');
    return stmt.all(userId) as Expense[];
  }

  update(id: string, data: Partial<CreateExpenseDto>): Expense {
    const now = new Date().toISOString();
    
    const updates: string[] = [];
    const values: any[] = [];
    
    if (data.amount !== undefined) {
      updates.push('amount = ?');
      values.push(data.amount);
    }
    if (data.category !== undefined) {
      updates.push('category = ?');
      values.push(data.category);
    }
    
    updates.push('updated_at = ?');
    values.push(now);
    values.push(id);
    
    const stmt = expensesDb.prepare(`
      UPDATE expenses 
      SET ${updates.join(', ')}
      WHERE id = ?
    `);
    
    stmt.run(...values);
    return this.findById(id)!;
  }

  delete(id: string): void {
    const stmt = expensesDb.prepare('DELETE FROM expenses WHERE id = ?');
    stmt.run(id);
  }
}

export const expenseRepository = new ExpenseRepository();
```

### Backup Strategy

Simple daily backups to S3:

```bash
#!/bin/bash
# scripts/backup-databases.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/tmp/trackstack-backup-$DATE"

mkdir -p $BACKUP_DIR

# Copy all databases
cp data/*.db $BACKUP_DIR/

# Upload to S3
aws s3 sync $BACKUP_DIR s3://trackstack-backups/$DATE/

# Clean up local backup
rm -rf $BACKUP_DIR

# Keep only last 30 days in S3
aws s3 ls s3://trackstack-backups/ | \
  awk '{print $2}' | \
  head -n -30 | \
  xargs -I {} aws s3 rm --recursive s3://trackstack-backups/{}
```

Run via cron:
```cron
0 2 * * * /app/scripts/backup-databases.sh
```

### Migration to .NET + EF Core (Future Phase 2)

When transitioning to .NET backend:

1. **Export current SQLite data** to JSON/CSV
2. **Create EF Core models** matching the domain
3. **Run migration script** to transform snake_case → PascalCase
4. **Point EF Core at the transformed databases**
5. **Verify data integrity**
6. **Update Astro API clients** to call .NET endpoints

**Estimated effort:** 2-3 days per module

The snake_case → PascalCase transformation is intentional technical debt that keeps us moving fast now and costs us a predictable, one-time migration later.

---

Will eventually become `@trackstack/ui` npm package.

**Component Example:**

```astro
---
// src/shared/components/Button/Button.astro
import type { HTMLAttributes } from 'astro/types';

interface Props extends HTMLAttributes<'button'> {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
}

const { variant = 'primary', size = 'md', class: className, ...props } = Astro.props;

const variantClasses = {
  primary: 'bg-blue-600 hover:bg-blue-700 text-white',
  secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900',
  danger: 'bg-red-600 hover:bg-red-700 text-white'
};

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg'
};
---

<button 
  class={`rounded font-medium transition-colors ${variantClasses[variant]} ${sizeClasses[size]} ${className || ''}`}
  {...props}
>
  <slot />
</button>
```

---

## Extraction Strategy

### When to Extract a Module

Extract when you see these signals:

1. **Different scaling needs** - Heat monitoring needs real-time processing, recipes don't
2. **Independent deployment cycles** - Expense tracker needs weekly updates, others are stable
3. **Team boundaries** - Different developers/teams working on different modules
4. **Revenue potential** - Expense tracker could be a standalone SaaS product

### Extraction Checklist

**Example: Extracting Expense Tracker**

**Step 1: Create new repository**
```bash
mkdir trackstack-expenses
cd trackstack-expenses
npm init astro
```

**Step 2: Copy module**
```bash
cp -r ../trackstack/app/src/modules/expense-tracker/* ./src/
```

**Step 3: Copy shared dependencies**
```bash
# Install auth package (if published)
npm install @trackstack/auth

# Or copy temporarily
cp -r ../trackstack/app/src/core/auth ./src/core/
cp -r ../trackstack/app/src/shared ./src/shared/
```

**Step 4: Update API configuration**
```typescript
// src/core/api/config.ts
export const API_BASE_URL = 'https://expenses-api.trackstack.com';
```

**Step 5: Deploy independently**
```bash
npm run build
# Deploy to expenses.trackstack.com
```

**Estimated effort:** 1-2 days

---

## Tech Stack

- **Framework:** Astro 5.x (static-first with SSR capabilities)
- **Interactive Components:** Angular (for dynamic forms, charts, filters)
- **Styling:** Tailwind CSS 4.x
- **Type Safety:** TypeScript with strict mode
- **Database:** SQLite with better-sqlite3 driver
- **Package Manager:** pnpm (fast, disk-efficient)
- **Hosting:** AWS S3 + CloudFront (via Terraform)


## Cost Considerations

**Current monthly cost:** <€1
- S3 storage: ~€0.02 (database backups + static assets)
- CloudFront: Free tier covers typical usage
- No database hosting costs (SQLite is serverless)
- No compute costs (static hosting with API routes)

**Post .NET migration estimated:** €3-8/month
- S3 storage: ~€0.05 (slightly more backup data)
- CloudFront: ~€0.50
- Backend hosting (AWS Fargate): €2-5
- Still using SQLite: €0 (no RDS needed!)

**Post module extraction estimated:** €8-15/month
- Additional S3 buckets: +€0.10
- Additional CloudFront distributions: +€0.50
- Multiple backend services: €5-10
- Still using SQLite: €0

**When to migrate from SQLite:**
- If concurrent writes become a bottleneck (unlikely with <50 users)
- If horizontal scaling is needed (multiple servers)
- If database replication is required for HA
- **For 5 users: SQLite is perfect, potentially forever**

**Philosophy:** Optimize for learning and speed. SQLite saves ~€30-50/month vs. managed databases while being faster and simpler. Migrate to PostgreSQL only when actual load demands it.

---