# TrackStack Agent Guidelines

This document provides context, commands, and coding standards for AI agents operating within the TrackStack codebase. TrackStack is a personal tracking application built with Astro, Tailwind CSS, and TypeScript.

## 1. Environment & Commands

### Working Directory
- The main application lives in `trackstack/app`.
- Always run commands from `trackstack/app` unless specified otherwise.
- **IMPORTANT:** Do not attempt to run `npm` commands. Use `pnpm`.

### Build & Run
| Command | Description |
|---------|-------------|
| `pnpm dev` | Start the development server (Astro). |
| `pnpm build` | Build the application for production. |
| `pnpm preview` | Preview the production build locally. |
| `pnpm astro check` | Run TypeScript type checking across the project. |

*Note: There are currently no explicit `test` or `lint` scripts in `package.json`. Rely on `pnpm astro check` for type safety.*

### Dependencies
- **Core:** Astro v5.x
- **Styling:** Tailwind CSS v4.x (via `@tailwindcss/vite`)
- **PWA:** `@vite-pwa/astro`
- **Package Manager:** `pnpm`

## 2. Project Structure

The project follows a modular, domain-driven structure inside `trackstack/app/src`:

- **`src/modules/`**: Contains domain-specific logic and components (e.g., `nba/`, `finance/`).
    - Encapsulate features here rather than polluting `shared`.
    - Example: `src/modules/nba/components/NightlyRecap.astro`
- **`src/pages/`**: Astro file-based routing.
    - Filenames here determine the URL route (e.g., `expenses/index.astro` -> `/expenses`).
    - Use lowercase, kebab-case for page filenames.
- **`src/shared/`**: Common utilities, UI components, and layouts.
    - `components/UI/`: Generic, reusable atoms (Cards, Buttons).
    - `layouts/`: Page wrappers (e.g., `AppShell.astro`).
    - `styles/`: Global CSS and Tailwind configuration.

## 3. Code Style & Conventions

### TypeScript
- **Strict Mode:** Enabled. Avoid `any`.
- **Interfaces:** Define props interfaces for all components.
    ```typescript
    interface Props {
      title: string;
      isActive?: boolean;
    }
    ```
- **Aliases:** Use the `@/` alias for all internal imports.
    - `import Card from '@/shared/components/UI/Card.astro';`
    - DO NOT use relative paths like `../../shared/...`.

### Astro Components
- **Frontmatter (---):**
    - Keep logic minimal and view-oriented.
    - Destructure `Astro.props` at the top.
    - Set defaults for optional props.
- **Naming:**
    - Component files: `PascalCase.astro` (e.g., `DomainNav.astro`).
    - Page files: `kebab-case.astro` (e.g., `settings.astro`).
- **Styling:**
    - Use `class:list` for conditional class logic.
    - Example: `<div class:list={['base-class', { 'active-class': isActive }]}>`

### Tailwind CSS
- **v4 Configuration:** Uses CSS variables defined in `@theme` blocks (see `global.css`).
- **Semantic Colors:** Prefer semantic variables over raw hex codes or standard Tailwind colors.
    - `text-text-main` (White/Primary text)
    - `text-text-muted` (Gray/Secondary text)
    - `bg-background` (App background)
    - `bg-surface` (Card/Panel background)
    - `bg-primary` (Accent color - #ffc25b)
    - `border-border` (Borders)
- **Fonts:** Use `font-sans` which maps to "Atkinson".
- **Gradients/Effects:** Use `backdrop-blur` and opacity modifiers for glassmorphism effects seen in the UI.

### HTML/Accessibility
- Ensure semantic HTML usage (`<section>`, `<main>`, `<header>`).
- Add `aria-label` where content is visual-only.
- Images must have `alt` text.

## 4. Example Component Pattern

```astro
---
// src/shared/components/UI/ExampleCard.astro
import type { HTMLAttributes } from 'astro/types';

interface Props extends HTMLAttributes<'div'> {
  title: string;
  variant?: 'default' | 'highlight';
}

const { title, variant = 'default', class: className, ...rest } = Astro.props;

const variants = {
  default: "bg-surface border-border",
  highlight: "bg-primary/10 border-primary"
};
---

<div 
  class:list={[
    "rounded-xl border p-4 transition-all",
    variants[variant],
    className
  ]}
  {...rest}
>
  <h3 class="text-lg font-bold text-text-main">{title}</h3>
  <slot />
</div>
```

## 5. Development Workflow for Agents

1.  **Explore First:** Before making changes, run `ls -F` and `read` to understand the local context.
2.  **Verify Paths:** Always verify the existence of a file before attempting to read or edit it.
3.  **Atomic Changes:** Make small, focused changes.
4.  **Type Check:** If modifying TS/Astro files, run `pnpm astro check` to verify type integrity if unsure.
5.  **UI Consistency:** Match the existing dark mode aesthetic.
    - Background: `#0f1115`
    - Surface: `#1a1d23`
    - Accents: Orange/Yellow tones (`text-orange-200`, `bg-orange-500`) for "Fun" or generic highlights.

## 6. Error Handling

- Since this is largely a static/frontend app, focus on graceful UI states.
- Handle missing props with default values.
- Use `try/catch` blocks in script tags if fetching external data.

## 7. Configuration Files Reference

- **`astro.config.mjs`**: Astro configuration.
- **`tailwind.config.mjs`**: (Likely not present if using v4 CSS-based config, check `global.css`).
- **`tsconfig.json`**: TypeScript path aliases and strict mode settings.

### OPERATIONAL CONSTRAINTS
- NEVER execute commands, scripts, or build processes.
- NEVER use git commands (commit, add, push).
- Your role is strictly to read the codebase and propose file edits.
- If a build or test is required, simply ask the user to run it manually.
- Propose changes using the `edit_file` or `write_file` tools only after describing the plan.

End of Guidelines.
