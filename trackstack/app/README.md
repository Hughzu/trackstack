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
├── src/
│   ├── modules/                    # Business capability modules (extraction candidates)
│   │   ├── expense-tracker/
│   │   │   ├── components/        # Domain-specific UI components
│   │   │   ├── pages/             # Module routing pages
│   │   │   ├── services/          # API clients and business logic
│   │   │   │   └── api.ts         # Expense API endpoints
│   │   │   ├── types/             # TypeScript types and interfaces
│   │   │   ├── utils/             # Domain utilities (formatters, validators)
│   │   │   └── index.ts           # Public API - only export what's needed
│   │   │
│   │   ├── heat-monitor/
│   │   │   ├── components/
│   │   │   ├── pages/
│   │   │   ├── services/
│   │   │   │   └── api.ts         # Heat monitoring API
│   │   │   ├── types/
│   │   │   ├── utils/
│   │   │   └── index.ts
│   │   │
│   │   ├── recipe-book/
│   │   │   ├── components/
│   │   │   ├── pages/
│   │   │   ├── services/
│   │   │   │   └── api.ts         # Recipe API
│   │   │   ├── types/
│   │   │   ├── utils/
│   │   │   └── index.ts
│   │   │
│   │   └── calorie-tracker/
│   │       ├── components/
│   │       ├── pages/
│   │       ├── services/
│   │       │   └── api.ts         # Calorie tracking API
│   │       ├── types/
│   │       ├── utils/
│   │       └── index.ts
│   │
│   ├── core/                       # Shared infrastructure (extraction = npm package)
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
├── components/              # UI components
│   ├── ExpenseList.astro   # List all expenses
│   ├── ExpenseForm.astro   # Create/edit form
│   └── ExpenseCard.astro   # Single expense display
│
├── pages/                   # Module-specific pages (optional)
│   └── ExpenseDetail.astro # Could be used in routing
│
├── services/                # Business logic and API
│   ├── api.ts              # HTTP calls to backend
│   └── expenseCalculations.ts # Domain logic
│
├── types/                   # TypeScript definitions
│   ├── expense.ts          # Expense, ExpenseCategory types
│   └── api.ts              # API request/response types
│
├── utils/                   # Module-specific utilities
│   ├── formatters.ts       # Currency, date formatting
│   └── validators.ts       # Form validation rules
│
└── index.ts                 # Public API
    // Only export what other parts of the app need
    export { ExpenseList, ExpenseForm } from './components';
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

## Shared Design System (`shared/`)

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

## Development Workflow

### Getting Started

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Adding a New Module

1. Create folder: `src/modules/my-new-module/`
2. Add required subfolders: `components/`, `services/`, `types/`, `utils/`
3. Create `index.ts` with public API exports
4. Add routes in `src/pages/my-module/`
5. Update this README

### Adding a New Component to Design System

1. Create folder: `src/shared/components/MyComponent/`
2. Add `MyComponent.astro` and `MyComponent.types.ts`
3. Export from `src/shared/components/index.ts`
4. Document usage in component file

---

## Tech Stack

- **Framework:** Astro 5.x (static-first with SSR capabilities)
- **Styling:** Tailwind CSS 4.x
- **Type Safety:** TypeScript with strict mode
- **Package Manager:** pnpm (fast, disk-efficient)
- **Hosting:** AWS S3 + CloudFront (via Terraform)

---

## Future Roadmap

### Phase 1: Monolith Validation (Current)
- ✅ Set up Astro monolith structure
- ✅ Implement module boundaries
- ⏳ Build core authentication
- ⏳ Implement all four tracking modules
- ⏳ Use daily for 3+ months to validate

### Phase 2: First Extraction
- Extract authentication to `@trackstack/auth` npm package
- Extract design system to `@trackstack/ui` npm package
- Identify first module candidate for standalone app

### Phase 3: Microservices Backend
- Replace embedded backend with C# monolith
- Implement hexagonal architecture
- Add observability (logging, metrics, tracing)

### Phase 4: Module Extraction
- Extract highest-value module to standalone app
- Migrate to dedicated backend microservice
- Implement event-driven communication

---

## Cost Considerations

**Current monthly cost:** <€1
- S3 storage: ~€0.02
- CloudFront: Free tier covers typical usage
- No compute costs (static hosting)

**Post-extraction estimated:** €5-15/month
- Additional S3 buckets: +€0.10
- CloudFront distributions: +€0.50
- Backend hosting (Fargate/EC2): €5-10
- RDS database: €3-5

**Philosophy:** Optimize for learning and agility, not premature cost optimization.

---

## Contributing

This is a personal project, but feedback is welcome:

1. Open an issue for bugs or suggestions
2. Fork and submit PRs for improvements
3. Follow existing code structure and patterns

---

## License

MIT - Build whatever you want with this structure.

---

## Questions?

Refer to the [TrackStack blog](https://blog.trackstack.com) for detailed architectural decisions and evolution journey.
