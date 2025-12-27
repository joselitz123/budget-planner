# Budget Planner Frontend - Development Guide

This file provides comprehensive guidance for Claude Code when working with the SvelteKit frontend of the Budget Planner project.

## Frontend Overview

**Tech Stack:**
- SvelteKit 2.0 - Full-stack web framework
- Svelte 5.0 - Reactive UI framework
- TypeScript 5.0 - Type safety
- Vite 5.0 - Build tool and dev server
- Shadcn-Svelte - Component library (Radix UI-based)
- Tailwind CSS 3.4 - Utility-first CSS framework
- IndexedDB (via `idb` library) - Offline storage
- Clerk - Authentication

**Current Status:** Initial setup complete (package.json configured), ready for development

---

## Quick Commands

```bash
cd frontend

# Install dependencies
npm install

# Run development server (with host access for devcontainer)
npm run dev

# Type checking
npm run check
npm run check:watch

# Build for production
npm run build

# Preview production build
npm run preview

# Linting and formatting
npm run lint
npm run format

# Generate PWA icons
npm run pwa:generate-icons
```

---

## Project Structure

**Expected Structure** (to be implemented):

```
frontend/
├── src/
│   ├── lib/                   # Utility functions, API clients
│   │   ├── components/        # Shared UI components
│   │   ├── stores/            # Svelte stores for state management
│   │   ├── db/                # IndexedDB client setup
│   │   ├── api/               # API client functions
│   │   └── utils/             # Helper functions
│   ├── routes/                # SvelteKit file-based routing
│   │   +layout.svelte         # Root layout with navigation
│   │   +page.svelte           # Home page
│   │   +error.svelte          # Error page
│   │   /auth/                 # Authentication routes
│   │   │   +layout.svelte     # Auth layout
│   │   │   /sign-in/          # Clerk sign-in
│   │   │   /sign-up/          # Clerk sign-up
│   │   /dashboard/            # Main dashboard
│   │   │   +page.svelte       # Dashboard overview
│   │   /budgets/              # Budget management
│   │   │   +page.svelte       # Budgets list
│   │   │   /[id]/             # Budget details
│   │   /transactions/         # Transaction tracking
│   │   │   +page.svelte       # Transactions list
│   │   │   /new/              # Add transaction
│   │   /categories/          # Category management
│   │   /analytics/            # Analytics and insights
│   │   /settings/             # User settings
│   │   └ /shares/             # Budget sharing
│   └── app.html               # HTML template with PWA meta tags
├── static/                     # Static assets
│   ├── logo.png               # App logo
│   ├── icons/                 # PWA icons (generated)
│   ├── favicon.ico            # Favicon
│   └── manifest.json          # PWA manifest
├── .env                       # Environment variables
├── svelte.config.js           # SvelteKit configuration
├── vite.config.js            # Vite configuration
├── tailwind.config.js        # Tailwind CSS configuration
├── tsconfig.json             # TypeScript configuration
├── postcss.config.js         # PostCSS configuration
└── package.json              # Dependencies and scripts
```

---

## Architecture Patterns

### SvelteKit File-Based Routing

SvelteKit uses file-based routing. The filesystem determines the routes:

```
routes/                    → /
routes/about              → /about
routes/blog/             → /blog
routes/blog/[slug]       → /blog/:slug (dynamic)
routes/[lang]/about      → /:lang/about (dynamic layout)
```

**Special Files:**
- `+page.svelte` - Page component
- `+layout.svelte` - Layout wrapper for child routes
- `+load.ts` - Server-side data loading
- `+page.server.ts` - Server actions
- `+error.svelte` - Error boundary
- `+page.svelte` with `/+layout.svelte` - Nested layouts

### Offline-First Architecture with IndexedDB

The frontend uses IndexedDB for offline data persistence:

**Pattern:**
1. **IndexedDB** stores all data locally for offline access
2. **Sync queue** captures changes when offline
3. **Sync API** (`/api/sync/push`, `/api/sync/pull`) handles bidirectional sync
4. **Conflict resolution** uses timestamp-based and owner-priority strategies

**IndexedDB Database Structure:**
```typescript
// Expected IndexedDB structure (using idb library)
interface BudgetDB {
  budgets: {
    key: string;
    value: Budget;
    indexes: { 'by-month': string };
  };
  categories: { key: string; value: Category };
  transactions: { key: string; value: Transaction };
  reflections: { key: string; value: Reflection };
  syncQueue: { key: string; value: SyncOperation };
}
```

### PWA Configuration

**Service Worker** (`vite-plugin-pwa`):
- Registers on app load
- Caches static assets for offline access
- Provides offline fallback
- Syncs data when connection restored

**Web App Manifest** (`static/manifest.json`):
- App name and short name
- Icons (generated via `npm run pwa:generate-icons`)
- Theme and background colors
- Display mode (standalone)

### Clerk Authentication Integration

**Pattern:**
```typescript
// Using @clerk/sveltekit
import { clerkClient } from '@clerk/sveltekit/sdk';

// In load functions
export const load = async (event) => {
  const { userId } = await event.locals.auth();
  if (!userId) throw redirect(302, '/sign-in');
  return { userId };
};
```

**Environment Variables:**
```env
PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard
```

---

## Component Organization

### Shadcn-Svelte Components

**Usage Pattern:**
1. Components are in `src/lib/components/ui/`
2. Add components via Shadcn CLI: `npx shadcn-svelte@latest add [component]`
3. Import and use in Svelte files

**Example:**
```svelte
<script>
  import { Button } from '$lib/components/ui/button';
  import { Card } from '$lib/components/ui/card';
</script>

<Card>
  <Button>Click me</Button>
</Card>
```

**Common Components:**
- `button` - Button with variants (default, outline, ghost, destructive)
- `card` - Card container with header, content, footer
- `input` - Form input with labels
- `select` - Dropdown selection
- `dialog` - Modal/dialog component
- `table` - Data table
- `badge` - Status badges
- `toast` - Notifications

### Layout Components

**Root Layout** (`src/routes/+layout.svelte`):
- Navigation bar
- Theme provider
- Clerk authentication wrapper
- PWA update prompt

**Feature Layouts:**
- Dashboard layout with sidebar
- Auth layout (minimal, centered)
- Settings layout with tabs

---

## State Management

### Svelte Stores Pattern

**Writable Stores** (for client state):
```typescript
// src/lib/stores/budgets.ts
import { writable } from 'svelte/store';

export const budgets = writable<Budget[]>([]);
export const currentBudget = writable<Budget | null>(null);
```

**Derived Stores** (computed values):
```typescript
import { derived } from 'svelte/store';

const totalSpent = derived(budgets, ($budgets) =>
  $budgets.reduce((sum, b) => sum + b.spent, 0)
);
```

**Using Stores in Components:**
```svelte
<script>
  import { budgets } from '$lib/stores/budgets';
</script>

{#each $budgets as budget}
  <BudgetCard {budget} />
{/each}
```

### IndexedDB Data Flow

**Pattern:**
1. **Load data** from IndexedDB on app init
2. **Sync with backend** via API on change
3. **Update IndexedDB** when data changes
4. **React to changes** using Svelte stores

**IndexedDB Client** (`src/lib/db/index.ts`):
```typescript
import { openDB } from 'idb';

export const db = await openDB<BudgetDB>('budget-planner', 1, {
  upgrade(db) {
    // Create object stores
    db.createObjectStore('budgets', { keyPath: 'id' });
    // ... other stores
  }
});
```

### Sync Queue Management

**Offline Operation Queue:**
```typescript
// src/lib/stores/syncQueue.ts
import { writable } from 'svelte/store';
import { addToQueue, processQueue } from '$lib/db/sync';

export const syncQueue = writable<SyncOperation[]>([]);
export const isOnline = writable(navigator.onLine);

// Add to queue when offline
if (!$isOnline) {
  addToQueue({ type: 'CREATE', entity: 'transaction', data });
}

// Process queue when online
isOnline.subscribe((online) => {
  if (online) processQueue();
});
```

---

## Styling

### Tailwind CSS Configuration

**Configuration** (`tailwind.config.js`):
```javascript
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        // ... other theme colors
      }
    }
  },
  plugins: [require('@tailwindcss/forms')]
};
```

**CSS Variables** (`src/app.css`):
```css
@layer base {
  :root {
    --primary: 210 40% 98%;  /* hsl values */
    --primary-foreground: 222 47% 11%;
    /* ... other variables */
  }
}
```

### Component Styling Patterns

**Class Variants** (using `class-variance-authority`):
```typescript
import { cva } from 'class-variance-authority';

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md text-sm',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground',
        destructive: 'bg-destructive text-destructive-foreground',
        outline: 'border border-input bg-background',
      },
      size: {
        default: 'h-10 px-4 py-2',
        sm: 'h-9 rounded-md px-3',
        lg: 'h-11 rounded-md px-8',
      }
    }
  }
);
```

**Tailwind Merge** (combine class names):
```typescript
import { cn } from '$lib/utils';

<div class={cn('base-class', condition && 'conditional-class')} />
```

---

## Development Workflow

### Hot Reload with Vite

- Vite dev server provides instant hot module replacement (HMR)
- Changes to `.svelte` files update in browser without refresh
- Changes to CSS update immediately

### Type Checking

**Check types:**
```bash
npm run check           # One-time type check
npm run check:watch     # Watch mode for development
```

**TypeScript Config** (`tsconfig.json`):
- Strict mode enabled
- Path aliases: `$lib` → `src/lib`
- SvelteKit type checking enabled

### Linting and Formatting

**Lint:**
```bash
npm run lint    # Check for issues
```

**Format:**
```bash
npm run format  # Format code with Prettier
```

**Prettier Config** (`.prettierrc`):
```json
{
  "semi": true,
  "singleQuote": true,
  "plugins": ["prettier-plugin-svelte", "prettier-plugin-tailwindcss"]
}
```

### Building for Production

**Build:**
```bash
npm run build
```

**Output:**
- `build/` directory with optimized assets
- Server-side rendering (SSR) build
- Client-side bundle
- PWA assets and service worker

---

## Environment Configuration

### Environment Variables

**File:** `.env` (create from `.env.example` if provided)

```env
# Clerk Authentication
PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard

# API Configuration
PUBLIC_API_URL=http://localhost:8080/api

# PWA Configuration
PUBLIC_APP_NAME="Budget Planner"
PUBLIC_APP_SHORT_NAME="Budget"
PUBLIC_APP_DESCRIPTION="Offline-first budget planning application"
PUBLIC_APP_THEME_COLOR="#3b82f6"
PUBLIC_APP_BACKGROUND_COLOR="#ffffff"

# Offline Sync Settings
PUBLIC_SYNC_INTERVAL=30000      # 30 seconds
PUBLIC_OFFLINE_RETRY_DELAY=5000 # 5 seconds
PUBLIC_MAX_OFFLINE_OPERATIONS=100
```

**Access in SvelteKit:**
```svelte
<script>
  const apiUrl = import.meta.env.VITE_PUBLIC_API_URL;
</script>
```

### Public API URL Configuration

**Development:**
```
PUBLIC_API_URL=http://localhost:8080/api
```

**Production:**
```
PUBLIC_API_URL=https://api.budgetplanner.com/api
```

### Clerk Keys Setup

1. Create Clerk account at https://clerk.com
2. Create new application
3. Copy publishable key to `.env`:
   ```
   PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
   ```
4. Configure redirect URLs in Clerk dashboard:
   - Allowed redirect URLs: `http://localhost:5173/*`
   - Allowed origins: `http://localhost:5173`

---

## API Integration Patterns

### Fetch API with TypeScript

**Typed API Client** (`src/lib/api/client.ts`):
```typescript
interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: { message: string; code: string };
}

export async function apiGet<T>(path: string): Promise<T> {
  const token = await clerkClient.session?.getToken();
  const response = await fetch(`${import.meta.env.VITE_PUBLIC_API_URL}${path}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.json();
}
```

**Usage:**
```typescript
import { apiGet } from '$lib/api/client';

const budgets = await apiGet<Budget[]>('/budgets');
```

### Loading States

**SvelteKit Load Function:**
```typescript
// src/routes/budgets/+page.ts
export const load = async ({ fetch }) => {
  const response = await fetch('/api/budgets');
  const budgets = await response.json();
  return { budgets };
};
```

**In Component:**
```svelte
<script>
  export let data;
  $: ({ budgets } = data);
</script>

{#await budgets}
  <LoadingSpinner />
{:then items}
  <BudgetList {items} />
{:catch error}
  <ErrorDisplay {error} />
{/await}
```

---

## PWA Development

### Service Worker

**Configuration** (`vite.config.js`):
```javascript
import { SvelteKitPWA } from 'vite-plugin-pwa';

export default {
  plugins: [
    SvelteKitPWA({
      strategies: ['networkFirst'],
      srcDir: 'static',
      filename: 'service-worker.js',
      manifest: {
        name: 'Budget Planner',
        short_name: 'Budget',
        // ...
      }
    })
  ]
};
```

### Icon Generation

**Generate PWA Icons:**
```bash
npm run pwa:generate-icons
```

**Input:** `static/logo.png` (source logo)
**Output:** `static/icons/` (multiple sizes)

---

## Common Patterns

### Route Protection

**Server Load Function:**
```typescript
// +page.server.ts
import { redirect } from '@sveltejs/kit';
import { clerkClient } from '@clerk/sveltekit/sdk';

export const load = async (event) => {
  const { userId } = await event.locals.auth();
  if (!userId) throw redirect(302, '/sign-in');
  return { userId };
};
```

### Form Validation

**Svelte Actions:**
```svelte
<script>
  import { enhance } from '$app/forms';

  let { form } = $props();
  let { data } = $form;
</script>

<form method="POST" use:enhance>
  <input name="name" value={data?.name || ''} />
  {#if form?.errors}
    <div class="error">{form.errors.name}</div>
  {/if}
  <button type="submit">Submit</button>
</form>
```

### Date Formatting

**Using Intl API:**
```typescript
const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    month: 'long',
    year: 'numeric'
  }).format(date);
};
```

---

## Important Files

- `starting-point.md` (root) - Full project specification
- `CLAUDE.md` (root) - Project-level documentation
- `backend/CLAUDE.md` - Backend development guide

---

## Development Environment Setup

The devcontainer includes:
- **Node.js 20** with Alpine Linux
- **Chromium** for browser testing
- Environment variables auto-generated on container create
- Development certificates for HTTPS

---

## Current Implementation Status

**Completed:**
- ✅ package.json with all dependencies
- ✅ npm scripts configured
- ✅ Project structure planned

**To Implement:**
- SvelteKit scaffolding
- Component library setup (Shadcn-Svelte)
- IndexedDB client
- API client
- Authentication integration
- PWA configuration
- Routing and pages

---

## Dependencies

**Core Framework:**
- `svelte@^5.0.0`
- `@sveltejs/kit@^2.0.0`
- `vite@^5.0.0`
- `typescript@^5.0.0`

**UI & Styling:**
- `@tailwindcss/forms@^0.5.9`
- `tailwindcss@^3.4.10`
- `class-variance-authority@^0.7.0`
- `clsx@^2.0.0`
- `lucide-svelte@^0.469.0`

**Offline Storage:**
- `idb@^8.0.0` - IndexedDB wrapper
- `workbox-window@^7.1.0` - Service worker management

**Authentication:**
- `@clerk/sveltekit@^2.0.0`

**Development Tools:**
- `@sveltejs/vite-plugin-svelte@^3.0.0`
- `svelte-check@^3.0.0`
- `vite-plugin-pwa@^0.20.0`
- `prettier@^3.0.0`
- `prettier-plugin-svelte@^3.0.0`
- `prettier-plugin-tailwindcss@^0.6.5`
